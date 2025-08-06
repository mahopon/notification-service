package handler

import (
	"encoding/json"
	"net/http"

	dto "github.com/mahopon/notification-service/internal/dto"
	notifySvc "github.com/mahopon/notification-service/internal/services"
)

type MainHandler struct {
	service notifySvc.NotificationService
}

func NewNotificationHandler(service notifySvc.NotificationService) *MainHandler {
	return &MainHandler{service: service}
}

type Reply struct {
	Response string `json:"message"`
}

func (h *MainHandler) StatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	reply := &Reply{
		Response: "Server up!",
	}

	err := json.NewEncoder(w).Encode(reply)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *MainHandler) NotifyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var incomingReq *dto.NotifyUserRequest

	err := json.NewDecoder(r.Body).Decode(&incomingReq)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	err = h.service.Notify(incomingReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	reply := &Reply{
		Response: "Email sent",
	}

	err = json.NewEncoder(w).Encode(reply)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
