package main

import (
	"fmt"
	infra "github.com/mahopon/notification-service/internal/infra"
	router "github.com/mahopon/notification-service/internal/routes"
	"log"
	"net/http"
)

func main() {
	mux := router.Setup()
	infra.Setup()
	PORT := 8080
	log.Printf("Server started on localhost: %d", PORT)
	http.ListenAndServe(fmt.Sprintf(":%d", PORT), mux)
}
