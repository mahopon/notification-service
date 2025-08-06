package infra

import (
	cfg "github.com/mahopon/notification-service/internal/config"
	"github.com/wneessen/go-mail"
	"log"
)

var mailClient *mail.Client

func Setup() {
	config, err := cfg.Load()
	if err != nil {
		log.Printf("ERROR LOADING CONFIG: %v", err)
		return
	} else {
		log.Println("Loaded config")
	}
	mailClient = initMailClient(config.Mail.Host, config.Mail.Email, config.Mail.Password)
}
