package main

import (
	"context"
	"log"
	db "movie_booking_system/database"
	handlers "movie_booking_system/handler"
	models "movie_booking_system/model"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	config, err := models.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	gin.SetMode(config.GinMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := db.ConnectDB(ctx, config); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := db.Close(shutdownCtx); err != nil {
			log.Printf("database shutdown error: %v", err)
		}
	}()

	r.GET("/", handlers.Root)
	r.GET("/healthz", handlers.HealthCheck)

	r.GET("/movies", handlers.GetMovies)
	r.GET("/movies/:id", handlers.GetMovieById)
	r.POST("/movies", handlers.AddMovie)
	r.PUT("/movies/:id", handlers.UpdateMovie)
	r.DELETE("/movies/:id", handlers.DeleteMovie)
	r.POST("/movies/:id/book", handlers.BookSeat)

	server := &http.Server{
		Addr:              ":" + config.Port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("server shutdown error: %v", err)
		}
	}()

	log.Printf("movie booking service listening on :%s", config.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
