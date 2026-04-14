FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /movie-booking-api .

FROM alpine:3.21

RUN apk add --no-cache ca-certificates \
	&& adduser -D -g '' appuser

WORKDIR /app

COPY --from=builder /movie-booking-api /usr/local/bin/movie-booking-api

EXPOSE 8080

USER appuser

ENTRYPOINT ["movie-booking-api"]
