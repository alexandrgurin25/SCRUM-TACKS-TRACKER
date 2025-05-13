FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /project-service ./cmd/main.go

FROM alpine

WORKDIR /app

COPY --from=builder /project-service .
COPY ./db/migrations ./db/migrations
COPY .env .env

EXPOSE 50054

CMD ["./project-service"]