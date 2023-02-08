api:
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pb/hello.proto
mod:
	@go mod tidy
c:
	@go run client/main.go
s:
	@go run server/main.go
