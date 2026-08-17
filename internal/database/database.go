package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectionPool(dbURL string) (*pgxpool.Pool, error) {
	pool, pgxpoolErr := pgxpool.New(context.Background(), dbURL)
	if pgxpoolErr != nil {
		return nil, fmt.Errorf("Connection pool could not be created./Bağlantı havuzu oluşturulamadı: %v", pgxpoolErr)
	}

	pingErr := pool.Ping(context.Background())
	if pingErr != nil {
		return nil, fmt.Errorf("Ping, database error./Ping, veritabanı hatası: %v", pingErr)
	}

	return pool, nil
}
