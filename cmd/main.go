package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mahopon/notification-service/internal/config"
	"github.com/mahopon/notification-service/internal/handler"
	infra "github.com/mahopon/notification-service/internal/infra"
	route "github.com/mahopon/notification-service/internal/routes"
	service "github.com/mahopon/notification-service/internal/services"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	router := mux.NewRouter()

	// Notifiers
	emailNotifier := infra.NewMailNotifier(cfg.Mail)
	telegramNotifier := infra.NewTelegramNotifier(cfg.Telegram)

	notificationService := service.Setup(emailNotifier, telegramNotifier)
	notificationHandler := handler.NewNotificationHandler(notificationService)
	route.Setup(router, notificationHandler)
	PORT := 8080
	log.Printf("Server started on localhost: %d", PORT)
	http.ListenAndServe(fmt.Sprintf(":%d", PORT), router)
}
