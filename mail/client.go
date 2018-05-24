/*
 * Revision History:
 *     Initial: 2018/05/24        Li Zebang
 */

package mail

import (
	"github.com/TechCatslab/firmness"
)

// MailClient implements the Client interface.
type MailClient struct {
	config *firmness.Config
}

// AddConfig -
func (mc *MailClient) AddConfig(config *MailConfig) error {
	return nil
}
