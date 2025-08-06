package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mahopon/notification-service/internal/handler"
	infra "github.com/mahopon/notification-service/internal/infra"
	route "github.com/mahopon/notification-service/internal/routes"
	service "github.com/mahopon/notification-service/internal/services"
)

func main() {
	infra.Setup()
	router := mux.NewRouter()
	notificationService := service.Setup(infra.EmailNotif)
	notificationHandler := handler.NewNotificationHandler(notificationService)
	route.Setup(router, notificationHandler)
	PORT := 8080
	log.Printf("Server started on localhost: %d", PORT)
	http.ListenAndServe(fmt.Sprintf(":%d", PORT), router)
}
