FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /task-service ./cmd/main.go

FROM alpine

WORKDIR /app

COPY --from=builder /task-service .
COPY ./migrations ./migrations
COPY .env .env

EXPOSE 50050

CMD [ "./task-service" ]
