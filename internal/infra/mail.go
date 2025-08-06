package infra

import (
	"fmt"
	config "github.com/mahopon/notification-service/internal/config"
	dto "github.com/mahopon/notification-service/internal/dto"
	"github.com/wneessen/go-mail"
	"log"
)

type EmailNotifier struct {
	client *mail.Client
}

func NewMailNotifier(host, email, password string) *EmailNotifier {
	client, err := mail.NewClient(host, mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(email), mail.WithPassword(password))
	if err != nil {
		log.Fatalf("Failed to create mail client: %s", err)
	}
	return &EmailNotifier{
		client: client,
	}
}

func (eN *EmailNotifier) Send(notifyUserDTO *dto.NotifyUserRequest) error {
	message := mail.NewMsg()
	from := fmt.Sprintf("Botto <%s>", config.Cfg.Mail.Email)
	if err := message.From(from); err != nil {
		log.Fatalf("Failed to set from address: %s", err)
	}
	if err := message.To(notifyUserDTO.To); err != nil {
		log.Fatalf("Failed to set to address: %s", err)
	}
	client := eN.client
	var t mail.ContentType
	switch notifyUserDTO.BodyType {
	case "html":
		t = mail.TypeTextHTML
	case "plain":
		t = mail.TypeTextPlain
	default:
		{
			t = mail.TypeTextPlain
		}
	}
	message.Subject(notifyUserDTO.Sub)
	message.SetBodyString(t, notifyUserDTO.Body)
	if err := client.DialAndSend(message); err != nil {
		log.Fatalf("failed to send mail: %s", err)
	} else {
		log.Printf("Sent email to: %s, with body:\n%s", notifyUserDTO.To, notifyUserDTO.Body)
	}
	return nil
}
