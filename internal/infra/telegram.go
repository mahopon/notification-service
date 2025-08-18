package infra

import (
	"fmt"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	config "github.com/mahopon/notification-service/internal/config"
	dto "github.com/mahopon/notification-service/internal/dto"
	util "github.com/mahopon/notification-service/internal/util"
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
	var cleanedBody string = "\n\n"
	if notifyUserDTO.Body != "" {
		if notifyUserDTO.BodyType == "MarkdownV2" {
			cleanedBody = cleanedBody + util.EscapeMarkdownV2(notifyUserDTO.Body)
		} else {
			cleanedBody = cleanedBody + notifyUserDTO.Body
		}
	}
	body := fmt.Sprintf("*%s*%s", notifyUserDTO.Sub, cleanedBody)
	msg := tgbotapi.NewMessage(target, body)
	if notifyUserDTO.BodyType == "html" {
		msg.ParseMode = "html"
	} else {
		msg.ParseMode = "MarkdownV2"
	}
	_, err := tN.Client.Send(msg)
	if err != nil {
		log.Printf("TeleNotifier: Error sending message: %v", err)
		return "", err
	}
	return "Message sent", nil
}

func (tN *TelegramNotifier) GetUpdatesChan() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	return tN.Client.GetUpdatesChan(u)
}
