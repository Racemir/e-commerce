package mail

import (
	"fmt"
	"net/smtp"
)

// E-posta göndermek için gerekli olan ayarlar
type MailSender struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

func NewMailSender(host, port, username, password, from string) *MailSender {
	return &MailSender{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		From:     from, // Gönderenin e-posta adresi
	}
}

// Aldığı parametreleri standart SMTP metin kurallarına göre (başlıklar alt alta, body'den önce boş satır) formatlar.
func buildMessage(to, subject, body string) []byte {
	message := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body)

	return []byte(message)
}

// veriyi ağ (network) üzerinden E-postayı SMTP sunucusuna ileten ana metot
func (m MailSender) SendEmail(to, subject, body string) error {
	msg := buildMessage(to, subject, body)

	// Kullanıcı adı ve şifre kullanarak basit bir kimlik doğrulama mekanizması sağlar
	auth := smtp.PlainAuth("", m.Username, m.Password, m.Host)

	// Sunucu adresi formatını "Host:Port" ,localhost:1025
	address := fmt.Sprintf("%s:%s", m.Host, m.Port)

	// address , auth , gönderen E-posta , Alıcı E-posta , E-posta gövdesini içeren byte
	err := smtp.SendMail(address, auth, m.From, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("SMTP sunucusuna e-posta iletilemedi/Failed to deliver email to SMTP server: %w", err)
	}

	return nil
}
