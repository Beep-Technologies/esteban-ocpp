FROM golang:1.16-alpine

WORKDIR /opt/app

COPY . .

CMD go run cmd/bb3-ocpp-ws/main.go
