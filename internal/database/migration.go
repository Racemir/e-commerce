package database

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrationsUp(dbUrl string) error {

	sourceUrl := "file://migrations"

	// golang-migrate kütüphanesinden yeni bir Migrate nesnesi (instance) oluşturur.
	// sourceUrl migrations klasörümüz , dbURL PostgreSQL veritabanına kurulan bağlantı.
	m, newError := migrate.New(sourceUrl, dbUrl)
	if newError != nil {
		return fmt.Errorf("Error migrate not return: %v", newError)
	}
	defer m.Close()

	// migrations klasörümüzde bulunan bütün .up.sql dosyalarını sırasıyla PostgreSQL veritabanında
	// çalıştırır(execute eder). veritabanı şemasını (schema) kodumuzdaki yapıya göre günceller.
	upError := m.Up()
	// hata varsa ve bu hata her şey güncel hatası değilse
	if upError != nil && upError != migrate.ErrNoChange {
		return fmt.Errorf("The migration up operation failed/migration up işlemi başarısız oldu: %v", upError)
	}

	// Sadece bilgi amaçlı log atıyoruz.
	if upError == migrate.ErrNoChange {
		log.Println("Veritabanı zaten güncel, yeni basılacak tablo yok.")
	} else {
		log.Println("Migration işlemi başarıyla tamamlandı! Tablolar oluşturuldu.")
	}

	return nil
}

func RunMigrationsDown(dbURL string) error {

	sourceUrl := "file://migrations"

	m, newError := migrate.New(sourceUrl, dbURL)
	if newError != nil {
		return fmt.Errorf("Error migrate not return/Hata migrate başlatılamadı: %v", newError)
	}
	defer m.Close()

	downError := m.Down()
	if downError != nil && downError != migrate.ErrNoChange {
		return fmt.Errorf("The migration down operation failed/migration down işlemi başarısız oldu: %v", downError)
	}

	if downError == migrate.ErrNoChange {
		log.Println("Geri alınacak bir migration bulunamadı.")
	} else {
		log.Println("Migration DOWN işlemi başarılı! Tablolar silindi.")
	}

	return nil
}
