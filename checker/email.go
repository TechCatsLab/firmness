/*
 * Revision History:
 *     Initial: 2018/05/26        Li Zebang
 */

package checker

import (
	"regexp"
)

// IsEmail return ture if email is valid.
func IsEmail(email string) bool {
	rgx := regexp.MustCompile(`\w+[\w.]*@[\w.]+\.\w+`)
	return rgx.MatchString(email)
}
