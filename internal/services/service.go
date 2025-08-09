package services

import (
	infra "github.com/mahopon/notification-service/internal/infra"
)

func Setup(emailNotifier *infra.EmailNotifier, telegramNotifier *infra.TelegramNotifier) NotificationService {
	notifierMux := NewNotifierMux()
	notifierMux.Register("email", emailNotifier)
	notifierMux.Register("telegram", telegramNotifier)
	return initService(notifierMux)
}
