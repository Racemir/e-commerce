CREATE TABLE payments(
    id SERIAL PRIMARY KEY, -- Ödemenin kendi benzersiz kimliği (Ödeme id'si)
    order_id INTEGER NOT NULL REFERENCES orders(id), -- Hangi siparişe ait olduğu (Sipariş id'si)
    user_id INTEGER NOT NULL REFERENCES users(id), -- Ödemeyi yapan müşterinin kimliği (Kullanıcı id'si)
    method VARCHAR(50) NOT NULL DEFAULT 'bank_transfer', -- Ödeme yöntemi (Varsayılan: banka havalesi)
    amount INTEGER NOT NULL, -- Ödenen toplam miktar / tutar
    currency TEXT NOT NULL, -- Para birimi (Örn: TRY, USD)
    reference_code VARCHAR(255) UNIQUE NOT NULL, -- Müşterinin havale yaparken açıklamaya yazacağı benzersiz referans kodu (Örn: ORD-X9K24P)
    status VARCHAR(50) NOT NULL DEFAULT 'awaiting_transfer' CHECK (status IN ('awaiting_transfer','submitted','approved','rejected','cancelled')), -- Ödemenin anlık durumu (havale bekliyor, bildirildi, onaylandı, reddedildi, iptal edildi)
    reviewed_by INTEGER REFERENCES users(id), -- Ödemeyi inceleyip onaylayan veya reddeden yetkilinin (adminin) id'si
    reviewed_at TIMESTAMP WITH TIME ZONE, -- Adminin bu ödemeyi incelediği tarih ve saat
    review_note TEXT, -- Adminin inceleme sonrası bıraktığı not (Örn: "Eksik tutar gönderilmiş")
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Bu ödeme kaydının sisteme ilk oluşturulduğu tarih
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP -- Bu ödeme kaydının en son güncellendiği tarih
);

-- Hızlı Arama Motorlarımız (Fihristler)
CREATE INDEX idx_payments_order_id ON payments(order_id); -- Sipariş id'sine göre hızlı arama fihristi
CREATE INDEX idx_payments_reference_code ON payments(reference_code); -- Referans koduna göre hızlı arama fihristi
CREATE INDEX idx_payments_status ON payments(status); -- Ödeme durumuna göre hızlı arama fihristi