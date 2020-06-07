package main

import (
	"fmt"
	"net/mail"
	"strings"
)

func main() {
	var pop3Accounts []Pop3Account
	var pop3Account Pop3Account
	pop3Account.Domain = "pop.qq.com"
	pop3Account.Port = 995
	pop3Account.Account = "2912150017@qq.com"
	pop3Account.Password = "dystpuamkpendghc"
	pop3Accounts = append(pop3Accounts, pop3Account)

	Run(pop3Accounts, Handle)
}

//Handle 邮件处理类
func Handle(msg *mail.Message, emailText string) {
	emailBody, err := NewEmail(msg)
	if err != nil {
		panic(err)
	}
	// 解析邮件
	emailBody.Header.Parse()
	emailBody.Parse()
	if strings.Contains(emailBody.Text, "发票") {
		fmt.Print(emailBody.Text)
		fmt.Print(emailText)
	}

	//todo do some think
	return
}
