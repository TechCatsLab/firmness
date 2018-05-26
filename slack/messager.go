/*
 * Revision History:
 *     Initial: 2018/05/24        Li Zebang
 */

package slack

import (
	"fmt"

	"github.com/nlopes/slack"
)

// Messager contains required information
type Messager struct {
	client  *slack.Client
	channel string
	Config  slack.PostMessageParameters
}

// NewMessager creates a new messager for channel messaging.
func NewMessager(token string, channel string, config *slack.PostMessageParameters) (*Messager, error) {
	if token == "" {
		return nil, fmt.Errorf("token cann't be null")
	}

	if channel == "" {
		return nil, fmt.Errorf("channel cann't be null")
	}

	var messager = &Messager{
		client:  slack.New(token),
		channel: channel,
	}

	if config == nil {
		messager.Config = slack.NewPostMessageParameters()
	} else {
		messager.Config = *config
	}

	return messager, nil
}

// PostMessage send message to the specified channel.
func (m *Messager) PostMessage(message string) error {
	_, _, err := m.client.PostMessage(m.channel, message, m.Config)
	return err
}
