APP_NAME=movie-booking-api

.PHONY: build run test fmt up down logs clean

build:
	go build ./...

run:
	go run ./...

test:
	go test ./...

fmt:
	gofmt -w main.go database/*.go handler/*.go model/*.go

up:
	docker compose up --build -d

down:
	docker compose down -v

logs:
	docker compose logs -f app

clean:
	rm -rf bin coverage.out dist
