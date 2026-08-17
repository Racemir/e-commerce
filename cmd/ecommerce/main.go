package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Racemir/e-commerce/internal/database"
	"github.com/Racemir/e-commerce/internal/server"
)

func main() {

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Eğer dışarıdan adres gelmezse varsayılan adresi kullan (Geliştirme ortamı için hayat kurtarır)
		dbURL = "postgres://admin:secretpassword@localhost:5432/ecommerce"
	}

	pool, connectionPoolErr := database.ConnectionPool(dbURL)
	if connectionPoolErr != nil {
		log.Fatalf("Error ConnectionPool: %v", connectionPoolErr)
	}
	defer pool.Close()

	mux := server.SetupRoutes()
	fmt.Println("Server 8080 portunda çalışıyor")

	listenAndServeErr := http.ListenAndServe(":8080", mux)
	if listenAndServeErr != nil {
		fmt.Print("Sunucu Başlatılamadı:", listenAndServeErr)
	}
}
