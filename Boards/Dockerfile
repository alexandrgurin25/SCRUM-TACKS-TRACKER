FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /board-service ./cmd/main.go

FROM alpine

WORKDIR /app

COPY --from=builder /board-service .
COPY ./db/migrations ./db/migrations
COPY .env .env

EXPOSE 50052

CMD [ "./board-service" ]
