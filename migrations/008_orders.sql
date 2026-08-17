CREATE TABLE orders(
    id  serial PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE, 
    status VARCHAR(50) NOT NULL DEFAULT 'pending_payment' CHECK (status IN ('pending_payment', 'payment_review', 'paid', 'processing', 'shipped', 'completed', 'cancelled')), -- Siparişin aşamaları (Sırasıyla: ödeme bekliyor, ödeme incelemede, ödendi, sipariş hazırlanıyor, kargolandı, tamamlandı, iptal edildi)
    total_amount INTEGER NOT NULL, -- Sepetteki tüm ürünlerin hesaplanmış genel toplam tutarı
    currency TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_orders_user_id ON orders(user_id);

CREATE TABLE order_items(
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE, -- Bu ürünün hangi ana siparişe/faturaya ait olduğu (Sipariş silinirse içindeki bu ürün kalemleri de silinir)
    product_id INTEGER NOT NULL REFERENCES products(id), -- Hangi ürün olduğu. (Dikkat: Yanında CASCADE yok! Yani satıcı ürünü silse bile, müşterinin geçmiş faturasındaki bu kayıt asla silinmez)
    product_name TEXT NOT NULL, -- Satın alındığı andaki ürünün adı (Snapshot - Satıcı yarın ürünün adını değiştirse bile bu faturadaki isim sabit kalır)
    quantity INTEGER NOT NULL, -- Müşterinin bu üründen kaç adet aldığı (Miktar)
    unit_price INTEGER NOT NULL, -- Ürünün o anki tek bir adet fiyatı (Birim fiyat)
    total_price INTEGER NOT NULL -- Adet x Birim Fiyat (Sadece bu ürün kaleminin toplam tutarı)
);