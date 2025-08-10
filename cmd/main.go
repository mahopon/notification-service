package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/mahopon/notification-service/internal/config"
	"github.com/mahopon/notification-service/internal/handler"
	"github.com/mahopon/notification-service/internal/infra"
	route "github.com/mahopon/notification-service/internal/routes"
	service "github.com/mahopon/notification-service/internal/services"
)

func main() {
	c := make(chan os.Signal, 2)

	cfg, err := config.Load(false)
	if err != nil {
		log.Fatal(err)
	}

	// Database
	db := infra.NewDatabaseConfig(cfg.Database)

	router := mux.NewRouter()

	// Notifiers
	emailNotifier := infra.NewMailNotifier(cfg.Mail)
	telegramNotifier := infra.NewTelegramNotifier(cfg.Telegram)

	notificationService := service.Setup(db, emailNotifier, telegramNotifier)
	if telegramNotifier != nil {
		go func() {
			updates := telegramNotifier.GetUpdatesChan()
			for update := range updates {
				notificationService.HandleUpdate(update)
			}
		}()
	}
	notificationHandler := handler.NewNotificationHandler(notificationService)
	route.Setup(router, notificationHandler)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		db.Database.Close()
		log.Printf("Service terminated...")
		os.Exit(1)
	}()

	PORT := 8080
	log.Printf("Server started on localhost: %d", PORT)
	http.ListenAndServe(fmt.Sprintf(":%d", PORT), router)
}
