package infra

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	config "github.com/mahopon/notification-service/internal/config"
	dto "github.com/mahopon/notification-service/internal/dto"
)

type TelegramNotifier struct {
	Client *tgbotapi.BotAPI
}

func NewTelegramNotifier(cfg *config.TGConfig) *TelegramNotifier {
	if cfg == nil {
		return nil
	}

	bot, err := tgbotapi.NewBotAPI(cfg.Key)
	if err != nil {
		log.Fatalf("Error connecting to Telegram API: %v", err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	return &TelegramNotifier{
		Client: bot,
	}
}

func (tN *TelegramNotifier) Send(notifyUserDTO *dto.NotifyUserRequest) (string, error) {
	log.Printf("Received message, but send not implemented for Telegram, %v", notifyUserDTO.Body)
	return "Message sent", nil
}
