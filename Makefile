run:
	fresh -c runner.conf
start:
	go run cmd/main.go
build:
	go build cmd/main.go
cover:
	go test -coverprofile=coverage.out ./...
fcover:
	go tool cover -func cover.out
swagger:
	swag init -g cmd/main.go