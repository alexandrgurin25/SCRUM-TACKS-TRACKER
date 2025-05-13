default: run

build:
	@go build -o ./bin/app ./cmd/main.go

run: build
	@./bin/app

test:
	@go test ./...

compose:
	@docker compose up -d