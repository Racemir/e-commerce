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

	// /api/health, / gibi yolların bağlı olduğu yönlendiriciyi (mux) alıyoruz.
	switch os.Args[1] {
	case "server":
		mux := server.SetupRoutes(dbPool, rdb)
		fmt.Println("Server 8080 portunda çalışıyor")

		// http.ListenAndServe: Sunucuyu belirtilen portta (8080) başlatır ve gelen istekleri dinlemeye başlar.
		// İkinci parametre olarak hazırladığımız yönlendiriciyi (mux) veriyoruz ki istekler doğru yerlere gitsin.
		// Bu satır bloklayıcıdır (blocking). Yani program burada sürekli bekler ve çalışmaya devam eder.
		listenAndServeErr := http.ListenAndServe(":8080", mux)
		// 8080 portu başka bir uygulama tarafından kullanılıyorsa
		if listenAndServeErr != nil {
			fmt.Print("Sunucu Başlatılamadı: %v", listenAndServeErr)
		}

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
