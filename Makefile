# Proje Değişkenleri
APP_NAME=ecommerce
MAIN_FILE=cmd/ecommerce/main.go

# Bu hedefler gerçek dosya değil.
.PHONY: help dev build clean up down

## clean: Derlenmiş eski dosyaları temizler.
clean:
	@echo "Derlenmiş dosyalar temizleniyor."
	rm -f $(APP_NAME)

## help: Kullanılabilecek komutları listeler.
help:
	@echo "Kullanılabilecek Komutlar:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## dev: Uygulamayı geliştirme modunda başlatır.
dev:
	@echo "Uygulama Başlatılıyor."
	go run $(MAIN_FILE)

## build: Uygulamayı tek bir çalıştırılabilir dosya olarak derler.
build:
	@echo "Uygulama Derleniyor."
	go build -o $(APP_NAME) $(MAIN_FILE)

## clean: Derlenmiş eski dosyaları temizler.
	@echo "Derlenmiş dosyalar temizleniyor."
	rm -f $(APP_NAME)

## up: Docker servislerini (Veritabanı, Redis, Mailpit vb.) ayağa kaldırır.
up:
	@echo "Application stack ayağa kalkıyor."
	docker compose up -d --build

## down: Çalışan tüm Docker servislerini durdurur ve kapatır.
down:
	@echo "Servisler Durduruluyor."
	docker compose down