FROM golang:1.21.4-alpine

WORKDIR /app

# Install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Generate swagger docs
RUN swag init -g server.go

RUN go build -o main server.go

EXPOSE 8080

CMD ["./main"]