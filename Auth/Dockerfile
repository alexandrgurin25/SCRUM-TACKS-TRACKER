FROM golang:1.24.2-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /auth

COPY . .
RUN go mod download

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -o auth ./cmd/app/main.go


FROM alpine

WORKDIR /auth

COPY --from=builder /auth/auth .
COPY ./db/migrations ./db/migrations
COPY config/.env config/.env


# Открываем порт
EXPOSE 50051

# Запускаем приложение
CMD ["./auth"]
