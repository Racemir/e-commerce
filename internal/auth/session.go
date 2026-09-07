package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// SessionData, Redis'te tutulan oturum verisidir.
type SessionData struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
}

// CreateSession, Redis üzerinde yeni bir oturum oluşturur.
func CreateSession(ctx context.Context, rdb *redis.Client, sessionID string, data SessionData, duration time.Duration) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("session marshal error: %w", err)
	}
	return rdb.Set(ctx, "session:"+sessionID, bytes, duration).Err()
}

// GetSession, Redis'ten oturum bilgisini okur.
func GetSession(ctx context.Context, rdb *redis.Client, sessionID string) (*SessionData, error) {
	val, err := rdb.Get(ctx, "session:"+sessionID).Result()
	if err != nil {
		return nil, err
	}
	var data SessionData
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// DeleteSession, Redis'ten oturum kaydını siler (logout için).
func DeleteSession(ctx context.Context, rdb *redis.Client, sessionID string) error {
	return rdb.Del(ctx, "session:"+sessionID).Err()
}
