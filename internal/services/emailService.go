package services

import (
	"fmt"
	"log"

	config "github.com/mahopon/notification-service/internal/config"
	infra "github.com/mahopon/notification-service/internal/infra"
	"github.com/wneessen/go-mail"
)

func SendEmail(to, sub, bodyType, body string) {
	// Add regex check for from and to
	message := mail.NewMsg()
	from := fmt.Sprintf("Botto <%s>", config.Cfg.Mail.Email)
	if err := message.From(from); err != nil {
		log.Fatalf("Failed to set from address: %s", err)
	}
	if err := message.To(to); err != nil {
		log.Fatalf("Failed to set from address: %s", err)
	}
	var t mail.ContentType
	switch bodyType {
	case "html":
		t = mail.TypeTextHTML
	case "plain":
		t = mail.TypeTextPlain
	default:
		{
			t = mail.TypeTextPlain
		}
	}

	client := infra.MailClient

	message.Subject(sub)
	message.SetBodyString(t, body)

	if err := client.DialAndSend(message); err != nil {
		log.Fatalf("failed to send mail: %s", err)
	} else {
		log.Printf("Sent email to: %s, with body:\n%s", to, body)
	}
}
