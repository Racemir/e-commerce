package auth

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// Veritabanı bağlantısı taşıyan struct
type Handler struct {
	DB  *pgxpool.Pool
	RDB *redis.Client
}

// Server içinde çalışan veritabanı bağlantısını "db" parametre olarak alır
// ve Handler struct'ının içindeki DB alanına yerleştirir hazır olan struct'ı geri döndürür
func NewHandler(db *pgxpool.Pool, rdb *redis.Client) *Handler {
	return &Handler{
		DB:  db,
		RDB: rdb,
	}
}
