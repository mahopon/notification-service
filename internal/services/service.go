package services

import (
	infra "github.com/mahopon/notification-service/internal/infra"
)

func Setup(emailNotifier *infra.EmailNotifier, telegramNotifier *infra.TelegramNotifier) NotificationService {
	notifierMux := NewNotifierMux()
	if emailNotifier != nil {
		notifierMux.Register("email", emailNotifier)
	}
	if telegramNotifier != nil {
		notifierMux.Register("telegram", telegramNotifier)
	}
	return initService(notifierMux)
}
