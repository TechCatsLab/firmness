/*
 * Revision History:
 *     Initial: 2018/05/26        Li Zebang
 */

package main

import (
	"fmt"

	"github.com/TechCatsLab/firmness/mail"
)

func main() {

	// Please make sure to enable smtp service before use
	config := &mail.Config{
		From: mail.Account{
			Email: "xxx@163.com",
		},
		Host: "smtp.163.com",
		Port: "25",
		Credentials: mail.Credentials{
			Username: "xxx@163.com",
			Password: "xxx",
		},
	}

	client, err := mail.NewClient(config)
	if err != nil {
		fmt.Println(err)
	}

	err = client.PostMessage("Test", "Test", mail.Account{Email: "xxx@gmail.com"})
	if err != nil {
		fmt.Println(err)
	}
}
