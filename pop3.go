package main

import (
	"bytes"
	"fmt"
	"net/mail"
	"os"

	"github.com/bytbox/go-pop3"
)

type EmailRetriever interface {
	Retr(int) (string, error)
}

type EmailFetcher interface {
	FetchEmails() error
}

type tlsEmailFetcher struct {
	username string
	password string
	popUrl   string
	popPort  int
}

// Pop3Account pop3 服务器对象
type Pop3Account struct {
	Account  string
	Password string
	Domain   string
	Port     int
}

//Run 获取pop3账号，从邮箱中获取邮件并解析，交给handle处理
func Run(pop3Accounts []Pop3Account, handle func(msg *mail.Message, emailText string)) {
	for _, pop3Account := range pop3Accounts {
		fetcher := NewTlsEmailFetcher(pop3Account.Account, pop3Account.Password, pop3Account.Domain, pop3Account.Port)
		err := fetcher.FetchEmails(handle)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

//NewTlsEmailFetcher 带ssl链接对象
func NewTlsEmailFetcher(username string, password string, url string, port int) *tlsEmailFetcher {
	return &tlsEmailFetcher{
		username: username,
		password: password,
		popUrl:   url,
		popPort:  port,
	}
}

func (f *tlsEmailFetcher) FetchEmails(handle func(*mail.Message, string)) error {
	uri := fmt.Sprintf("%s:%d", f.popUrl, f.popPort)
	client, err := pop3.DialTLS(uri)
	if err != nil {
		return fmt.Errorf("could not dial server: %v", err)
	}
	defer client.Quit()

	err = client.Auth(f.username, f.password)
	if err != nil {
		return fmt.Errorf("could not authenticate: %v", err)
	}

	msgIds, _, err := client.ListAll()
	if err != nil {
		return fmt.Errorf("could not list messages: %v", err)
	}

	return f.harvestMessages(client, msgIds, handle)
}

func (f *tlsEmailFetcher) harvestMessages(retriever EmailRetriever, msgIds []int, handle func(*mail.Message, string)) error {

	for _, id := range msgIds {
		func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println("panic error  %v", err)
				}
			}()

			text, err := retriever.Retr(id)
			if err != nil {

				fmt.Println("could not retrieve message (id=%d): %v", id, err)
				return
			}
			msg, err := mail.ReadMessage(bytes.NewBufferString(text))
			if err != nil {
				fmt.Println("could not read message (id=%d): %v", id, err)
				return
			}
			handle(msg, text)

		}()
	}
	return nil
}
