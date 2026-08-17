CREATE TABLE product_media(
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE, -- Bu medyanın hangi ürüne ait olduğu (Ürün silinirse, o ürüne ait tüm resim ve videolar da otomatik silinir)
    media_type VARCHAR(20) NOT NULL, -- Yüklenen dosyanın türü ('image' yani resim mi, yoksa 'video' mu olduğu)
    url TEXT UNIQUE NOT NULL, -- Dosyanın sunucuda veya bulutta tutulduğu adres/link (Aynı dosyanın iki kere eklenmesini engellemek için UNIQUE yapıldı)
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'ready', 'failed')), -- Arka planda resim/video işlenme durumu (Sırasıyla: bekliyor, işleniyor, kullanıma hazır, hata verdi)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);