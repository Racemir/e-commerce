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

// QueueVerificationEmail asenkron olarak doğrulama e-postası göndermek için
// gerekli olan işi (job) Redis kuyruğuna (email_queue) ekler.
// Worker uygulamasının bu kuyruktan işi alıp e-postayı göndermesini sağlar.
func QueueVerificationEmail(ctx context.Context, rdb *redis.Client, email, token string) error {

	emailJob := EmailQueuePayload{
		Type:  "send_verification_email",
		Email: email,
		Token: token,
	}

	jobJSON, marshalErr := json.Marshal(emailJob)
	if marshalErr != nil {
		return fmt.Errorf("failed to marshal email email job/E-posta işi hazırlanamadı: %w", marshalErr)
	}

	// "email_queue" isimli Redis listesinin solundan veriyi ekliyoruz.
	// Worker (işçi) uygulaması ise bu listeyi BRPOP ile sağdan okuyacak.
	pushErr := rdb.LPush(ctx, "email_queue", jobJSON).Err()
	if pushErr != nil {
		return fmt.Errorf("failed to push to Redis queue/Redis kuyruğuna gönderim başarısız oldu: %w", pushErr)
	}

	return nil
}
