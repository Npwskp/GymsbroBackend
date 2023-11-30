.PHONY: run build

run:
	swag init -g server.go
	go run server.go

build:
	swag init -g server.go
	go build server.go
	server.exe