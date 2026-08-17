CREATE TABLE products(
    id SERIAL PRIMARY KEY,
    category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE CASCADE, -- Ürünün hangi kategoriye ait olduğu
    slug VARCHAR(255) UNIQUE NOT NULL, -- Ürünün URL/Link adı (Örn: 'oyuncu-klavyesi-siyah'). SEO ve temiz linkler için kullanılır
    price INTEGER NOT NULL, -- Ürün fiyatı (Finansal doğruluk için FLOAT yerine INTEGER kullanıldı, kuruş cinsinden tutulacak)
    stock INTEGER NOT NULL DEFAULT 0, -- Depodaki güncel ürün adedi (Sipariş onaylandıkça veritabanı işlemiyle -Transaction- düşürülecek)
    status VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'active', 'archived')), -- Ürünün sitede görünürlüğü (Taslak, Satışa Açık, Arşivlenmiş)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE product_translations(
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    locale VARCHAR(10) NOT NULL, -- Dil kodu ('tr', 'en')
    name VARCHAR(255) NOT NULL, -- Ürünün o dildeki başlığı
    description TEXT NOT NULL, -- Ürünün o dildeki uzun açıklaması ve detayları
    UNIQUE(product_id, locale) -- Güvenlik Kilidi: Bir ürünün aynı dilde iki farklı çeviri kaydı olmasını engeller
);

-- Hızlı Arama Fihristi (Madde 111 Kuralı)
CREATE INDEX idx_products_slug ON products(slug); -- Link üzerinden (slug) bir ürün arandığında saniyeler içinde bulmak için fihrist oluşturur