package services

import (
	"errors"

	"github.com/mahopon/notification-service/internal/dto"
)

type MockNotifier struct {
	Calledwith *dto.NotifyUserRequest
	ShouldPass bool
}

func (m *MockNotifier) Send(req *dto.NotifyUserRequest) (string, error) {
	m.Calledwith = req
	if !m.ShouldPass {
		return "", errors.New("send failed")
	}
	return "mock_reply", nil
}
