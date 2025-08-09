package services

import (
	"errors"
	"log"

	dto "github.com/mahopon/notification-service/internal/dto"
)

type NotificationService interface {
	Notify(notifyUserDTO *dto.NotifyUserRequest) (string, error)
}

type DefaultNotificationService struct {
	notifierMux *NotifierMux
}

type Notifier interface {
	Send(notifyUserDTO *dto.NotifyUserRequest) (string, error)
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

func (s *DefaultNotificationService) Notify(notifyUserDTO *dto.NotifyUserRequest) (string, error) {
	channel := notifyUserDTO.Channel
	notifier, ok := s.notifierMux.notifiers[channel]
	var reply string
	var err error
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
	return reply, nil
}
