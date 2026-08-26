package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Racemir/e-commerce/internal/database"
	"github.com/Racemir/e-commerce/internal/server"
	gracefulshutdown "github.com/quii/go-graceful-shutdown"
)

func main() {

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Eğer dışarıdan adres gelmezse varsayılan adresi kullan (Geliştirme ortamı için hayat kurtarır)
		dbURL = "postgres://admin:secretpassword@localhost:5433/ecommerce?sslmode=disable"
	}

	// Veritabanı ile uygulama arasında bir bağlantı havuzu oluşturur.
	// Her istekte yeni bağlantı açmak yerine bu havuzdaki hazır bağlantılar kullanılır.
	dbPool, connectionPoolErr := database.ConnectionPool(dbURL)
	if connectionPoolErr != nil {
		log.Fatalf("Error ConnectionPool: %v", connectionPoolErr)
	}
	// Veritabanı bağlantılarının güvenlice kapatılması
	defer dbPool.Close()

	rdb, connectionRedisError := database.ConnectionRedis()
	if connectionRedisError != nil {
		log.Fatalf("Error connectionRedis: %v", connectionRedisError)
	}
	defer rdb.Close()

	// Argüman kontolü
	if len(os.Args) < 2 {
		fmt.Println("Kullanım: go run main.go [server|migrate|migrate-down]")
		return
	}

	// /api/health, / gibi yolların bağlı olduğu yönlendiriciyi (mux) alıyoruz.
	switch os.Args[1] {
	case "server":
		mux := server.SetupRoutes(dbPool, rdb)
		fmt.Println("Server 8080 portunda çalışıyor")
		ctx := context.Background()
		httpServer := &http.Server{
			Addr:    ":8080",
			Handler: mux,
		}
		gserver := gracefulshutdown.NewServer(httpServer)
		// http.ListenAndServe: Sunucuyu belirtilen portta (8080) başlatır ve gelen istekleri dinlemeye başlar.
		// İkinci parametre olarak hazırladığımız yönlendiriciyi (mux) veriyoruz ki istekler doğru yerlere gitsin.
		// Bu satır bloklayıcıdır (blocking). Yani program burada sürekli bekler ve çalışmaya devam eder.
		listenAndServeErr := gserver.ListenAndServe(ctx)
		// 8080 portu başka bir uygulama tarafından kullanılıyorsa
		if listenAndServeErr != nil {
			log.Fatalf("An error occurred while the server was shutting down (some requests may have been left incomplete)/Sunucu kapanırken hata oluştu (bazı istekler yarım kalmış olabilir): %v", listenAndServeErr)
		}
		fmt.Println("The server was shut down safely/Sunucu güvenli bir şekilde kapatıldı")

	case "migrate":

		runMigrationsUpError := database.RunMigrationsUp(dbURL)
		if runMigrationsUpError != nil {
			log.Fatalf("Migration up hatası: %v", runMigrationsUpError)
		}

		fmt.Println("Migration işlemi başarıyla tamamlandı.")

	case "migrate-down":

		runMigrationsDownError := database.RunMigrationsDown(dbURL)
		if runMigrationsDownError != nil {
			log.Fatalf("Migration down hatası: %v", runMigrationsDownError)
		}
		fmt.Println("Sistem başarıyla geri alındı")
	default:
		log.Fatalf("Unknown command")
	}

}
