package services

import (
	infra "github.com/mahopon/notification-service/internal/infra"
)

func Setup(emailNotifier *infra.EmailNotifier) NotificationService {
	notifierMux := NewNotifierMux()
	notifierMux.Register("email", emailNotifier)
	return initService(notifierMux)
}
