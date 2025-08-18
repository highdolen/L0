# --- Stage 1: build ---
FROM golang:1.24.2-alpine AS builder

WORKDIR /app

# Установим зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o service ./cmd/server

# --- Stage 2: run ---
FROM alpine:latest

WORKDIR /root/

# Копируем бинарник из builder stage
COPY --from=builder /app/service .

# Копируем веб-файлы
COPY --from=builder /app/web ./web
COPY --from=builder /app/migrations ./migrations

# Устанавливаем переменные окружения (можно перенести в compose)
ENV SERVER_PORT=:8080

EXPOSE 8080

CMD ["./service"]
