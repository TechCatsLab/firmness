/*
 * Revision History:
 *     Initial: 2018/05/24        Li Zebang
 */

package slack

import (
	"github.com/nlopes/slack"

	"github.com/TechCatslab/firmness"
)

// SlackClient implements the Client interface.
type SlackClient struct {
	client *slack.Client
	config *firmness.Config
}

// AddConfig -
func (sc *SlackClient) AddConfig(config *SlackConfig) error {
	return nil
}
