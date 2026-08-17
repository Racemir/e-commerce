CREATE TABLE bank_transfer_submissions(
    id SERIAL PRIMARY KEY,
    payment_id INTEGER NOT NULL REFERENCES payments(id) ON DELETE CASCADE, -- Bu havale bildirimi hangi ödeme kaydı için yapıldı?
    transfer_date TIMESTAMP WITH TIME ZONE NOT NULL, -- Müşterinin parayı bankadan gönderdiği tarih
    transferred_amount INTEGER NOT NULL, -- Gönderilen tutar
    sender_name VARCHAR(255) NOT NULL, -- Parayı gönderen kişinin adı (Örn: Fehmi Emir Özcan)
    note TEXT, -- Müşterinin eklemek istediği not (Zorunlu değil, o yüzden NOT NULL yok)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Madde 111 Kuralı: Hızlı Arama Fihristi
CREATE INDEX idx_bank_transfer_submissions_payment_id ON bank_transfer_submissions(payment_id);