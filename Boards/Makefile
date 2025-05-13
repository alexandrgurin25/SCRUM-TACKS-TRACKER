generate-proto-boards:
	protoc -I api/boards --go_out=. --go-grpc_out=. boards.proto
generate-proto-types:
	protoc -I api/types --go_out=. --go-grpc_out=. types.proto
run:
	go run cmd/main.go