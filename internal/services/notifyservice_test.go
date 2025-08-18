package services

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mahopon/notification-service/internal/dto"
)

var bucket string = "user_chat"

func setupTestData(shouldPass bool, channelName string) (*DefaultNotificationService, *MockDB, *MockNotifier) {
	db := NewMockDB()
	db.Set("user_chat", "alice", "12345")

	notifier := &MockNotifier{ShouldPass: shouldPass}
	mux := NewNotifierMux()
	mux.Register(channelName, notifier)

	svc := &DefaultNotificationService{
		NotifierMux: mux,
		DB:          db,
	}
	return svc, db, notifier
}

func TestNotify_Success_Telegram(t *testing.T) {
	channelName := "telegram"
	svc, _, notifier := setupTestData(true, channelName)
	req := &dto.NotifyUserRequest{
		To:      "alice",
		Sub:     "Test",
		Body:    "Hello",
		Channel: channelName,
	}

	reply, err := svc.Notify(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reply != "mock_reply" {
		t.Errorf("expected mock_reply, got %s", reply)
	}
	if notifier.Calledwith.To != "12345" {
		t.Errorf("expected mapped chatID 12345, got %s", notifier.Calledwith.To)
	}
}

func TestNotify_Success(t *testing.T) {
	channelName := "test"
	svc, _, notifier := setupTestData(true, channelName)
	req := &dto.NotifyUserRequest{
		To:      "alice",
		Sub:     "Test",
		Body:    "Hello",
		Channel: channelName,
	}

	reply, err := svc.Notify(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reply != "mock_reply" {
		t.Errorf("expected mock_reply, got %s", reply)
	}
	if notifier.Calledwith.To != "alice" {
		t.Errorf("expected mapped chatID alice, got %s", notifier.Calledwith.To)
	}
}

func TestNotify_Failure(t *testing.T) {
	channelName := "test"
	svc, _, notifier := setupTestData(false, channelName)
	req := &dto.NotifyUserRequest{
		To:      "alice",
		Sub:     "Test",
		Body:    "Hello",
		Channel: "test",
	}

	reply, err := svc.Notify(req)

	if err == nil {
		t.Fatalf("unexpected no error: %v", err)
	}
	if reply != "" {
		t.Errorf("unexpected reply, got %s", reply)
	}
	if notifier.Calledwith.To != "alice" {
		t.Errorf("expected mapped chatID alice, got %s", notifier.Calledwith.To)
	}
}

func TestHandleUpdate_NewUser(t *testing.T) {
	svc, db, notifier := setupTestData(true, "telegram")
	update := tgbotapi.Update{
		Message: &tgbotapi.Message{
			Text: "/start",
			Chat: &tgbotapi.Chat{
				ID:       12344,
				UserName: "bob",
			},
			Entities: []tgbotapi.MessageEntity{{
				Type:   "bot_command",
				Offset: 0,
				Length: 6,
			},
			},
		},
	}

	err := svc.HandleUpdate(update)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	chatID, err := db.Get(bucket, update.Message.Chat.UserName)
	if err != nil {
		t.Errorf("not key not found in db")
	}
	if chatID != "12344" {
		t.Errorf("chat id stored incorrectly")
	}

	if notifier.Calledwith == nil {
		t.Fatal("expected notifier to be called")
	}
	if notifier.Calledwith.Body != "You will now be able to receive notifications from this bot." {
		t.Errorf("unepected body text, got %s", notifier.Calledwith.Body)
	}
}

func TestHandleUpdate_ExistingUser(t *testing.T) {
	svc, db, notifier := setupTestData(true, "telegram")
	update := tgbotapi.Update{
		Message: &tgbotapi.Message{
			Text: "/start",
			Chat: &tgbotapi.Chat{
				ID:       12345,
				UserName: "alice",
			},
			Entities: []tgbotapi.MessageEntity{{
				Type:   "bot_command",
				Offset: 0,
				Length: 6,
			},
			},
		},
	}

	err := svc.HandleUpdate(update)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	chatID, err := db.Get(bucket, update.Message.Chat.UserName)
	if err != nil {
		t.Errorf("not key not found in db")
	}
	if chatID != "12345" {
		t.Errorf("chat id stored incorrectly")
	}

	if notifier.Calledwith == nil {
		t.Fatal("expected notifier to be called")
	}
	if notifier.Calledwith.Body != "You are already registered!" {
		t.Errorf("unepected body text, got %s", notifier.Calledwith.Body)
	}
}
