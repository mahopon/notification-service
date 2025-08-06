package infra

import (
	"github.com/wneessen/go-mail"
	"log"
	"sync"
)

var MailClient *mail.Client
var onlyOnce sync.Once

func initMailClient(host string, email, password string) *mail.Client {
	onlyOnce.Do(func() {
		client, err := mail.NewClient(host, mail.WithPort(587), mail.WithSMTPAuth(mail.SMTPAuthPlain), mail.WithUsername(email), mail.WithPassword(password))
		if err != nil {
			log.Fatalf("Failed to create mail client: %s", err)
		} else {
			log.Println("Mail client started")
		}
		MailClient = client
	})
	return MailClient
}
