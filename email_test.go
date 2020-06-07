package main

import (
	"fmt"
	"testing"
)

func TestParseBody(t *testing.T) {
	filename := "../example/email/email.log"
	emailText := ReadEmailFromFile(filename)
	msg := MockEmail(emailText)
	emailBody, err := NewEmail(msg)
	if err != nil {
		panic(err)
	}
	// 解析邮件
	emailBody.Header.Parse()
	emailBody.Parse()
	fmt.Println(emailBody)
}
