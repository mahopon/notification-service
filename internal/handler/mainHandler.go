package handler

import (
	"encoding/json"
	emailSvc "github.com/mahopon/notification-service/internal/services"
	"net/http"
)

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	reply := struct {
		Response string `json:"response"`
	}{
		Response: "Server up!",
	}

	err := json.NewEncoder(w).Encode(reply)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func EmailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var emailReq struct {
		To   string `json:"to"`
		Sub  string `json:"sub"`
		Body string `json:"body"`
	}

	err := json.NewDecoder(r.Body).Decode(&emailReq)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	emailSvc.SendEmail(emailReq.To, emailReq.Sub, "plain", emailReq.Body)
	//log.Printf("To: %s, Sub: %s, Body: %s", emailReq.To, emailReq.Sub, emailReq.Body)

}
