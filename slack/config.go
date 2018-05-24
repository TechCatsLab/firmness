/*
 * Revision History:
 *     Initial: 2018/05/24        Li Zebang
 */

package slack

import (
	"github.com/TechCatslab/firmness"
)

// SlackConfig -
type SlackConfig struct {
}

// NewConfig -
func NewConfig() firmness.Config {
	return &SlackConfig{}
}
