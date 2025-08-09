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
	NotifierMux NotifierRegistry
}

type NotifierRegistry interface {
	Register(channel string, n Notifier)
	Get(channel string) (Notifier, bool)
}

type Notifier interface {
	Send(notifyUserDTO *dto.NotifyUserRequest) (string, error)
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

func initService(notifierMux *NotifierMux) NotificationService {
	return &DefaultNotificationService{
		NotifierMux: notifierMux,
	}
}

func (s *DefaultNotificationService) Notify(notifyUserDTO *dto.NotifyUserRequest) (string, error) {
	channel := notifyUserDTO.Channel
	notifier, ok := s.NotifierMux.Get(channel)
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
