FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /api-gateway ./cmd/main.go

FROM alpine

WORKDIR /app

COPY --from=builder /api-gateway .
COPY .env .env

EXPOSE 50055

CMD [ "./api-gateway", "--congig=./.env" ]
