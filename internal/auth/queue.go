package auth

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// EmailQueuePayload, Redis'e eklenecek e-posta işinin (job) JSON formatındaki yapısıdır.
type EmailQueuePayload struct {
	Type  string `json:"type"`
	Email string `json:"email"`
	Token string `json:"token"`
}

func QueueVerificationEmail(ctx context.Context, rdb *redis.Client, email, token string) error {

	emailJob := EmailQueuePayload{
		Type:  "send_verification_email",
		Email: email,
		Token: token,
	}

	jobJSON, marshalErr := json.Marshal(emailJob)
	if marshalErr != nil {
		return fmt.Errorf("E-posta işi hazırlanamadı/failed to marshal email email job: %w", marshalErr)
	}

	// "email_queue" isimli Redis listesinin solundan veriyi ekliyoruz.
	// Worker (işçi) uygulaması ise bu listeyi BRPOP ile sağdan okuyacak.
	pushErr := rdb.LPush(ctx, "email_queue", jobJSON).Err()
	if pushErr != nil {
		return fmt.Errorf("Redis kuyruğuna gönderim başarısız oldu/failed to push to Redis queue: %w", pushErr)
	}

	return nil
}
