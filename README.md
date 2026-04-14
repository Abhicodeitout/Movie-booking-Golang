# Movie Booking Golang API

Movie Booking Golang API is a small but production-shaped REST service for managing movies, showtimes, and seat bookings with Go, Gin, and MongoDB.

## Highlights

- Clean REST endpoints for movie CRUD and seat booking
- MongoDB-backed persistence with configurable database and collection names
- Graceful shutdown and explicit database health checks
- Request validation for payload shape and duplicate seat numbers
- Docker, Docker Compose, Make targets, and GitHub Actions CI
- Unit tests for configuration and core handler behavior

## Stack

- Go 1.23+
- Gin
- MongoDB
- Docker and Docker Compose for local orchestration
- GitHub Actions for CI

## Project Structure

```text
.
├── .github/workflows/ci.yml
├── database/
├── handler/
├── model/
├── .env.example
├── .gitignore
├── Dockerfile
├── Makefile
├── docker-compose.yml
├── go.mod
├── main.go
└── README.md
```

## Configuration

Copy the sample file for local development:

```bash
cp .env.example .env
```

Available environment variables:

```bash
PORT=8080
GIN_MODE=release
MONGO_URI=mongodb://localhost:27017
MONGO_DATABASE=movie_booking
MONGO_COLLECTION=movies
```

`MONGO_URI` is required. The application loads `.env` automatically when present.

## Run Locally

### Option 1: Go + local MongoDB

```bash
go mod tidy
go run ./...
```

### Option 2: Docker Compose

```bash
docker compose up --build
```

The API is available at `http://localhost:8080`.

## Make Commands

```bash
make build
make run
make test
make fmt
make up
make down
make logs
```

## API Overview

| Method | Route | Description |
| --- | --- | --- |
| GET | `/` | Service metadata |
| GET | `/healthz` | Database health check |
| GET | `/movies` | List all movies |
| GET | `/movies/:id` | Get one movie |
| POST | `/movies` | Create a movie |
| PUT | `/movies/:id` | Update a movie |
| DELETE | `/movies/:id` | Delete a movie |
| POST | `/movies/:id/book` | Book a seat |

## Sample Payloads

### Create or Update Movie

```json
{
   "title": "Inception",
   "director": "Christopher Nolan",
   "year": 2010,
   "description": "A science fiction heist thriller.",
   "showtimes": [
      {
         "time": "2026-04-14T19:00:00Z",
         "seats": [
            { "number": 1, "booked": false, "reserved": false },
            { "number": 2, "booked": false, "reserved": false }
         ]
      }
   ]
}
```

### Book Seat

```json
{
   "showtime_index": 0,
   "seat_number": 2
}
```

## Example Requests

### Health Check

```bash
curl http://localhost:8080/healthz
```

### Create Movie

```bash
curl -X POST http://localhost:8080/movies \
   -H 'Content-Type: application/json' \
   -d '{
      "title":"Inception",
      "director":"Christopher Nolan",
      "year":2010,
      "description":"A science fiction heist thriller.",
      "showtimes":[
         {
            "time":"2026-04-14T19:00:00Z",
            "seats":[
               {"number":1,"booked":false,"reserved":false},
               {"number":2,"booked":false,"reserved":false}
            ]
         }
      ]
   }'
```

### Book Seat

```bash
curl -X POST http://localhost:8080/movies/<movie-id>/book \
   -H 'Content-Type: application/json' \
   -d '{"showtime_index":0,"seat_number":2}'
```

## Validation Notes

- Movie payloads require title, director, year, description, and at least one showtime.
- Each showtime requires at least one seat.
- Duplicate seat numbers within the same showtime are rejected.
- The booking endpoint rejects already booked seats and invalid seat numbers.

## Quality Checks

```bash
go test ./...
go build ./...
```

CI runs these checks automatically on pushes and pull requests.

## License

This project is licensed under the MIT License. See the LICENSE file for details.





 

