package routes

import (
	"github.com/gorilla/mux"
	"github.com/mahopon/notification-service/internal/handler"
	"net/http"
)

func Setup() http.Handler {
	r := mux.NewRouter()
	r.StrictSlash(true)

	r.HandleFunc("/", handler.StatusHandler)
	r.HandleFunc("/email", handler.EmailHandler)
	return r
}
