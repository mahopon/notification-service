package infra

import (
	cfg "github.com/mahopon/notification-service/internal/config"
	"log"
)

var EmailNotif *EmailNotifier

func Setup() {
	config, err := cfg.Load()
	if err != nil {
		log.Printf("ERROR LOADING CONFIG: %v", err)
		return
	} else {
		log.Println("Loaded config")
	}
	EmailNotif = NewMailNotifier(config.Mail.Host, config.Mail.Email, config.Mail.Password)
}
