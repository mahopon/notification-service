package routes

import (
	"github.com/gorilla/mux"
	handler "github.com/mahopon/notification-service/internal/handler"
)

func Setup(r *mux.Router, mainHandler *handler.MainHandler) {
	r.StrictSlash(true)
	r.HandleFunc("/", mainHandler.StatusHandler)
	r.HandleFunc("/notify", mainHandler.NotifyHandler)
}
