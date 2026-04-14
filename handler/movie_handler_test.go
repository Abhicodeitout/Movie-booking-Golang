package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestRootReturnsServiceMetadata(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/", Root)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var payload map[string]string
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if payload["service"] != "movie-booking-api" {
		t.Fatalf("expected service name movie-booking-api, got %q", payload["service"])
	}
}

func TestHealthCheckWithoutDatabaseReturnsServiceUnavailable(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/healthz", HealthCheck)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status %d, got %d", http.StatusServiceUnavailable, recorder.Code)
	}
}

func TestGetMovieByIDRejectsInvalidObjectID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/movies/:id", GetMovieById)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/movies/not-a-valid-id", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}
}

func TestBookSeatRejectsInvalidPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/movies/:id/book", BookSeat)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/movies/"+primitive.NewObjectID().Hex()+"/book",
		strings.NewReader(`{"showtime_index":0}`),
	)
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}
}

func TestAddMovieRejectsDuplicateSeatNumbers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/movies", AddMovie)

	body := `{
		"title":"Duplicate Seats",
		"director":"Test Director",
		"year":2026,
		"description":"Validation check",
		"showtimes":[
			{
				"time":"2026-04-14T19:00:00Z",
				"seats":[
					{"number":1,"booked":false,"reserved":false},
					{"number":1,"booked":false,"reserved":false}
				]
			}
		]
	}`

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/movies", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}

	if !strings.Contains(recorder.Body.String(), "duplicate seat number") {
		t.Fatalf("expected duplicate seat error, got %s", recorder.Body.String())
	}
}
