package services

import (
	"errors"
	"log"

	dto "github.com/mahopon/notification-service/internal/dto"
)

type NotificationService interface {
	Notify(notifyUserDTO *dto.NotifyUserRequest)
}

type DefaultNotificationService struct {
	notifierMux *NotifierMux
}

type Notifier interface {
	Send(notifyUserDTO *dto.NotifyUserRequest) error
}

type NotifierMux struct {
	notifiers map[string]Notifier
}

func (m *NotifierMux) Register(channel string, n Notifier) {
	m.notifiers[channel] = n
}

func NewNotifierMux() *NotifierMux {
	return &NotifierMux{
		notifiers: make(map[string]Notifier),
	}
}

func initService(notifierMux *NotifierMux) NotificationService {
	return &DefaultNotificationService{
		notifierMux: notifierMux,
	}
}

func (s *DefaultNotificationService) Notify(notifyUserDTO *dto.NotifyUserRequest) {
	channel := notifyUserDTO.Channel
	notifier, ok := s.notifierMux.notifiers[channel]
	if ok {
		err := notifier.Send(notifyUserDTO)
		if err != nil {
			log.Fatalf("Error sending notification :%v", err)
		}
	} else {
		log.Fatalf("Error sending notification: %v", errors.New("channel doesn't exist"))
	}
}
