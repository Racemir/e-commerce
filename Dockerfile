#Kodları derlememiz için Go kurulu bir sanal makine getiri. Adı builder olur.
FROM golang:1.26.5-alpine AS builder

#Konteynırın çalışma dizini (yeni bir dosyada app de olsun /bin veya /etc dağılmasın)
WORKDIR /app

# Sadece kütüphaneyi kopyaladım. Kodda bir değişiklik olursa bütün kütüphaneyi baştan
# indirmesin hafızaadki (cache) kullansın diye önden içeri aldım.
COPY go.mod ./

# go.mod dosyamıza bakarak projemizin ihtiyacı olan kütüphaneleri internetten indirip 
# bilgisayarımız kuracak
RUN go mod download

# Projedieki kodlarımı bilgisayardan alıp sanal bilgisayara aktardım
COPY . .

# CGO_ENABLED=0 Dışarıdan hiçbir C dili kütüphanesi bağımlılığı katmaması için yazdım.
# Linux da çalışacak, program adı ecommerce, derlenecek dosya yeri şeklinde ayarladım.
RUN CGO_ENABLED=0 GOOS=linux go build -o ecommerce ./cmd/ecommerce/main.go

# Hafif ve güvenilir, çalıştırmak için ideal
FROM alpine:latest

# --no-cache paket arşivlerini yerel dizinde saklamıyorum İmaj boyutunu küçük tutuyorum.
# ca-certificates güvenli web (HTTPS/SSL) bağlantıları için gerekli kök sertifikaları yüklüyorum.
# ffmpeg ses video dönüştürme, kaydetme , işleme
RUN apk --no-cache add ca-certificates ffmpeg

WORKDIR /app

# Bilgisayarın (builder) içinden sadece derlenmiş olan ecommerce dosyasını çekip alıp, 
# bu yeni boş bilgisayara koymak için yazdım.
COPY --from=builder /app/ecommerce .

# 8080 portundan veri bekliyorum.
EXPOSE 8080

# Sunucu ayağa kalktığında porgramı çalıştırıyorum.
CMD ["./ecommerce"]