gen: 
	protoc -I api/proto --go_out=. --go-grpc_out=. auth.proto
run: 
	go run cmd/app/main.go

cover:
	go test -coverprofile=coverage.out ./internal/...
	go tool cover -html=coverage.out
	go tool cover -func=coverage.out | grep -E "^total:" 
	rm -rf coverage.out

