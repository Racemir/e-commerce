package auth

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CheckDuplicateEmail(ctx context.Context, db *pgxpool.Pool, email string) (bool, error) {

	var isEmailTaken bool

	query := `
	SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)
	`

	queryRowError := db.QueryRow(ctx, query, email).Scan(&isEmailTaken)
	if queryRowError != nil {
		return false, queryRowError
	}

	return isEmailTaken, nil
}

// CreateUser, yeni bir kullanıcıyı veritabanına ekler ve oluşan ID'yi döndürür
func CreateUser(ctx context.Context, db *pgxpool.Pool, name, email, hashedPassword string) (int, error) {
	var newID int

	query := `
		INSERT INTO users (name, email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, 'customer', NOW(), NOW())
		RETURNING id;
	`
	// Veritabanına bir sorgu gönderip geriye sadece bir satır cevap beklediğimizde kullanılır
	// Scan: Dönen 'id' değerini newID değişkeninin hafıza adresine yaz
	queryRowError := db.QueryRow(ctx, query, name, email, hashedPassword).Scan(&newID)
	if queryRowError != nil {
		return 0, queryRowError
	}
	return newID, nil
}

func CreateVerificationToken(ctx context.Context, db *pgxpool.Pool, userID int, token string, expires time.Time) error {

	query := `
	INSERT INTO email_verification_tokens (user_id,token,expires_at) VALUES ($1,$2,$3)
	`

	_, execError := db.Exec(ctx, query, userID, token, expires)
	if execError != nil {
		return execError
	}

	return nil
}
