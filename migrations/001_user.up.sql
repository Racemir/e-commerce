CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL, -- Kullanıcının tam adı
    email VARCHAR(255) UNIQUE NOT NULL, -- Sisteme giriş ve iletişim adresi (UNIQUE ile aynı e-posta iki kez kaydolamaz)
    password_hash VARCHAR(255) NOT NULL, -- Güvenlik kuralı: Şifreler asla düz metin (123456) olarak tutulmaz, şifrelenmiş (hash'lenmiş) hali kaydedilir
    role VARCHAR(50) NOT NULL DEFAULT 'customer', -- Kullanıcının yetkisi (Varsayılan müşteri, yetkililer için 'admin' yapılacak)
    email_verified_at TIMESTAMP WITH TIME ZONE, -- E-posta doğrulama tarihi (Eğer burası boş/NULL ise kullanıcı henüz mailini onaylamamış demektir)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);