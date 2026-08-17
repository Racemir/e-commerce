CREATE TABLE email_verification_tokens(
    id SERIAL PRIMARY KEY, -- Tablonun kendi benzersiz kimliği
    token VARCHAR(255) UNIQUE NOT NULL, -- Kullanıcıya e-postayla gidecek olan o rastgele ve eşsiz şifre
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE, -- Bu şifre hangi kullanıcıya ait? (Kullanıcı silinirse bu şifre de otomatik silinir)
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL, -- Şifrenin son kullanma tarihi (ölüm vakti)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Sisteme eklendiği an
    used_at TIMESTAMP WITH TIME ZONE -- Kullanıcı linke tıklarsa buraya tarih düşülür (tek kullanımlık kilidi)
);