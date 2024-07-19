package main

import (
	db "movie_booking_system/database"
	handlers "movie_booking_system/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Connect to MongoDB
	err := db.ConnectDB()
	if err != nil {
		panic(err)
	}

	// Define routes
	r.GET("/movies", handlers.GetMovies)
	r.GET("/movies/:id", handlers.GetMovieById)
	r.POST("/movies", handlers.AddMovie)
	r.PUT("/movies/:id", handlers.UpdateMovie)
	r.DELETE("/movies/:id", handlers.DeleteMovie)
	r.POST("/movies/:id/book", handlers.BookSeat)

	// Start server
	r.Run(":8080")
}
