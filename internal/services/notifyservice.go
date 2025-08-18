package services

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	dto "github.com/mahopon/notification-service/internal/dto"
	"github.com/mahopon/notification-service/internal/infra"
)

type NotificationService interface {
	Notify(notifyUserDTO *dto.NotifyUserRequest) (string, error)
	HandleUpdate(update tgbotapi.Update) error
}

type NotifierRegistry interface {
	Register(channel string, n Notifier)
	Get(channel string) (Notifier, bool)
}

type Notifier interface {
	Send(notifyUserDTO *dto.NotifyUserRequest) (string, error)
}

type DefaultNotificationService struct {
	NotifierMux NotifierRegistry
	DB          infra.Database
}

type NotifierMux struct {
	Notifiers map[string]Notifier
}

func (m *NotifierMux) Register(channel string, n Notifier) {
	m.Notifiers[channel] = n
}

func (m *NotifierMux) Get(channel string) (Notifier, bool) {
	n, ok := m.Notifiers[channel]
	return n, ok
}

func NewNotifierMux() *NotifierMux {
	return &NotifierMux{
		Notifiers: make(map[string]Notifier),
	}
}

func initService(dbConfig *infra.DatabaseConfig, notifierMux *NotifierMux) NotificationService {
	return &DefaultNotificationService{
		NotifierMux: notifierMux,
		DB:          dbConfig,
	}
}

func (s *DefaultNotificationService) Notify(notifyUserDTO *dto.NotifyUserRequest) (string, error) {
	channel := notifyUserDTO.Channel
	notifier, ok := s.NotifierMux.Get(channel)
	originalTarget := notifyUserDTO.To
	var reply string
	var err error
	if channel == "telegram" {
		bucket := "user_chat"
		to := notifyUserDTO.To
		target, err := s.DB.Get(bucket, to)
		if err != nil {
			log.Printf("ERROR: %v", err)
			return "", err
		}
		notifyUserDTO.To = target
	}

	if ok {
		reply, err = notifier.Send(notifyUserDTO)
		if err != nil {
			log.Printf("Error sending notification :%v", err)
			return "", err
		}
	} else {
		err = errors.New("channel doesn't exist")
		log.Printf("Error sending notification: %v", err)
		return "", err
	}
	log.Printf("NotifyService: Successfully sent message to user: %s", originalTarget)
	return reply, nil
}

func (s *DefaultNotificationService) HandleUpdate(update tgbotapi.Update) error {
	if update.Message == nil || !update.Message.IsCommand() {
		return fmt.Errorf("ERROR: %v", "invalid message")
	}
	bucket := "user_chat"
	switch update.Message.Command() {
	case "start":
		userID := update.Message.Chat.UserName
		var chatID string
		chatID, err := s.DB.Get(bucket, userID)
		if err != nil {
			chatID = strconv.FormatInt(update.Message.Chat.ID, 10)
			err = s.DB.Set(bucket, userID, chatID)
			if err != nil {
				return err
			}
			notifyReq := &dto.NotifyUserRequest{
				To:       userID,
				Sub:      "User registered for service",
				Body:     "You will now be able to receive notifications from this bot.",
				Channel:  "telegram",
				BodyType: "MarkdownV2",
			}
			_, err = s.Notify(notifyReq)
			if err != nil {
				log.Printf("ERROR: Unable to reply to registration: %v", err)
				return err
			}
		} else {
			log.Printf("NotifyService: User already registered, userID: %s, chatID: %s", userID, chatID)
			notifyReq := &dto.NotifyUserRequest{
				To:       userID,
				Sub:      "User already registered",
				Body:     "You are already registered!",
				Channel:  "telegram",
				BodyType: "MarkdownV2",
			}
			_, err = s.Notify(notifyReq)
			if err != nil {
				log.Printf("ERROR: Unable to reply to existing registration: %v", err)
				return err
			}
		}
	}
	return nil
}
