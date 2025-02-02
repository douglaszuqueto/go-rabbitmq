include .env

.EXPORT_ALL_VARIABLES:

CGO_ENABLED=0

dev:
	go run main.go

dev-setup:
	go run main.go -setup

build:
	go build -o ./bin/rabbitmq main.go

run: build
	./bin/rabbitmq

prod: build
	upx ./bin/rabbitmq
	./bin/rabbitmq

.PHONY: dev build run prod