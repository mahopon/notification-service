package infra

import (
	"log"
	"strconv"

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

	bot.Debug = cfg.Debug

	log.Printf("Authorized on account %s", bot.Self.UserName)
	return &TelegramNotifier{
		Client: bot,
	}
}

func (tN *TelegramNotifier) Send(notifyUserDTO *dto.NotifyUserRequest) (string, error) {
	target, _ := strconv.ParseInt(notifyUserDTO.To, 10, 64)
	body := notifyUserDTO.Body
	msg := tgbotapi.NewMessage(target, body)
	tN.Client.Send(msg)
	return "Message sent", nil
}

func (tN *TelegramNotifier) GetUpdatesChan() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	return tN.Client.GetUpdatesChan(u)
}
