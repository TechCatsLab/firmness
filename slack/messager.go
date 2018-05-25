/*
 * Revision History:
 *     Initial: 2018/05/24        Li Zebang
 */

package slack

import (
	"github.com/nlopes/slack"

	"github.com/TechCatslab/firmness"
)

// Messager implements the Client interface.
type Messager struct {
	client *slack.Client
	config *firmness.Config
}

// NewMessager creates a new messager for channel messaging.
func NewMessager(config *Config) (*Messager, error) {
	return nil, nil
}
