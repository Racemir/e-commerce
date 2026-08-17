/* 
  KATEGORİ ANA TABLOSU
  İsim barındırmaz, sadece kimlik tutar. Çünkü isimler kullanıcının seçtiği dile göre değişecektir.
*/
CREATE TABLE categories(
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE category_translations(
    id SERIAL PRIMARY KEY,
    category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    locale VARCHAR(10) NOT NULL, -- Hangi dilde olduğu (Örn: 'tr' veya 'en')
    name VARCHAR(255) NOT NULL, -- Kategorinin o dildeki tam adı (Örn: locale 'tr' ise 'Elektronik', 'en' ise 'Electronics')
    UNIQUE(category_id, locale) -- Güvenlik Kilidi: Aynı kategorinin aynı dilde (Örn: tr) birden fazla çevirisi olmasını kesin olarak engeller
);