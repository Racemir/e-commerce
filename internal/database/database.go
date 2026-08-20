package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func ConnectionPool(dbURL string) (*pgxpool.Pool, error) {

	// pgxpool.New: Verilen adres (dbURL) ile yeni bir bağlantı havuzu ayarlarını yapar.
	// context.Background(): Go'da işlemlerin yaşam süresini (timeout, iptal vb.) yöneten bağlamdır.
	// Background(), "özel bir kısıtlaması olmayan standart bir işlem" anlamına gelir.
	pool, pgxpoolError := pgxpool.New(context.Background(), dbURL)
	if pgxpoolError != nil {
		return nil, fmt.Errorf("Connection pool could not be created./Bağlantı havuzu oluşturulamadı: %v", pgxpoolError)
	}

	// pool.Ping: Veritabanına ufak bir "Orada mısın?" sinyali gönderir ve cevap bekler.
	pingError := pool.Ping(context.Background())
	if pingError != nil {
		return nil, fmt.Errorf("Ping, database error./Ping, veritabanı hatası: %v", pingError)
	}

	return pool, nil
}

func ConnectionRedis() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Docker-compose dosyasında belirlediğimiz port
		Password: "",               // Şifre belirlemediğimiz için boş
		DB:       0,                // Varsayılan veritabanı numarası
	})

	// context.Background() kullanarak basit bir ping atıyoruz
	if pingErr := rdb.Ping(context.Background()).Err(); pingErr != nil {
		return nil, fmt.Errorf("Ping, redis error./Ping, redis hatası: %v", pingErr)
	}

	return rdb, nil
}
