/*
 * Revision History:
 *     Initial: 2018/06/25        Li Zebang
 */

package github

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Channel -
type Channel struct {
	Ready   chan *Client
	Prepare chan *Client
}

// ChannelPool -
type ChannelPool map[string]Channel

// Timeout -
var Timeout = 30 * time.Second

// NewChannelPool -
func NewChannelPool(tokens ...*Token) (ChannelPool, error) {
	if len(tokens) <= 0 {
		return nil, errors.New("the number of tokens cannot be less than 1")
	}

	for index := range tokens {
		if tokens[index] == nil {
			return nil, errors.New("token cannot be nil")
		}
	}

	cs := make(map[string][]*Client)
	for _, token := range tokens {
		if _, exist := cs[token.Tag]; !exist {
			cs[token.Tag] = make([]*Client, 0)
		}
		client, _ := NewClient(token)
		rl, _, err := client.RateLimits(context.Background())
		if err != nil {
			continue
		}
		client.Remaining = rl.Core.Remaining
		client.Reset = rl.Core.Reset.Unix()
		cs[token.Tag] = append(cs[token.Tag], client)
	}

	cp := make(map[string]Channel)
	for key, clients := range cs {
		if len(clients) == 0 {
			continue
		}
		if _, exist := cp[key]; !exist {
			cp[key] = Channel{
				Ready:   make(chan *Client, len(clients)),
				Prepare: make(chan *Client, len(clients)),
			}
		}
		for _, client := range clients {
			cp[key].Ready <- client
			cp[key].Prepare <- client
		}
	}

	return cp, nil
}

// Prepare -
func (cp ChannelPool) Prepare() {
	for k := range cp {
		key := k
		go func() {
			for client := range cp[key].Prepare {
				<-time.NewTimer(time.Unix(client.Reset, 0).Sub(time.Now())).C
				rl, _, _ := client.RateLimits(context.Background())
				old := client.Remaining
				client.Remaining = rl.Core.Remaining
				client.Reset = rl.Core.Reset.Unix()
				cp[key].Prepare <- client
				if old == 0 {
					cp[key].Ready <- client
				}
			}
		}()
	}
	select {}
}

// Get -
func (cp ChannelPool) Get(tag string) (*Client, error) {
	if _, exist := cp[tag]; !exist {
		return nil, fmt.Errorf("the client of label %s does not exist", tag)
	}

	select {
	case client := <-cp[tag].Ready:
		return client, nil
	case <-time.NewTimer(Timeout).C:
		return nil, fmt.Errorf("there is no client available")
	}
}

// Put -
func (cp ChannelPool) Put(client *Client) {
	if client.Remaining == 0 {
		cp[client.Tag].Prepare <- client
	} else {
		cp[client.Tag].Ready <- client
	}
}
