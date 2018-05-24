/*
 * Revision History:
 *     Initial: 2018/05/24        Li Zebang
 */

package mail

import (
	"github.com/TechCatslab/firmness"
)

// MailConfig -
type MailConfig struct {
}

// NewConfig -
func NewConfig() firmness.Config {
	return &MailConfig{}
}
