#/bin/bash

swag init -g server.go
go build server.go
./server