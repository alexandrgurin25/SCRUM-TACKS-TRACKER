default: run

build:
	@go build -o ./bin/app ./cmd/main.go

run: build
	@./bin/app

gen:
	@go run github.com/99designs/gqlgen generate