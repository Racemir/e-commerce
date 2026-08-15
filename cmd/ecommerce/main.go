package main

import (
	"fmt"
	"net/http"

	"github.com/Racemir/e-commerce/internal/server"
)

func main() {
	mux := server.SetupRoutes()
	fmt.Println("Server 8080 portunda çalışıyor")

	listenAndServeErr := http.ListenAndServe(":8080", mux)
	if listenAndServeErr != nil {
		fmt.Print("Sunucu Başlatılamadı:", listenAndServeErr)
	}
}
