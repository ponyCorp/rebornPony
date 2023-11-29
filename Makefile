run:
	fresh -c runner.conf
start:
	go run cmd/main.go
build:
	go build cmd/main.go
swagger:
	swag init -g cmd/main.go