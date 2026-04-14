package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	db "movie_booking_system/database"
	models "movie_booking_system/model"
)

const dbTimeout = 5 * time.Second

func dbContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), dbTimeout)
}

func HealthCheck(c *gin.Context) {
	ctx, cancel := dbContext()
	defer cancel()

	if err := db.Ping(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func GetMovies(c *gin.Context) {
	collection := db.GetCollection()
	ctx, cancel := dbContext()
	defer cancel()

	var movies []models.Movie
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get movies"})
		return
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var movie models.Movie
		err := cur.Decode(&movie)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode movie"})
			return
		}
		movies = append(movies, movie)
	}

	if err := cur.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to iterate movies"})
		return
	}

	c.JSON(http.StatusOK, movies)
}

func GetMovieById(c *gin.Context) {
	movieID := c.Param("id")
	oid, err := primitive.ObjectIDFromHex(movieID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	collection := db.GetCollection()
	ctx, cancel := dbContext()
	defer cancel()

	var movie models.Movie
	err = collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&movie)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load movie"})
		return
	}

	c.JSON(http.StatusOK, movie)
}

func AddMovie(c *gin.Context) {
	var movie models.Movie
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie data"})
		return
	}

	if err := validateMovie(movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := db.GetCollection()
	ctx, cancel := dbContext()
	defer cancel()

	movie.ID = primitive.NilObjectID

	result, err := collection.InsertOne(ctx, movie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert movie"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Movie added", "id": result.InsertedID})
}

func UpdateMovie(c *gin.Context) {
	movieID := c.Param("id")
	oid, err := primitive.ObjectIDFromHex(movieID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	var updatedMovie models.Movie
	if err := c.ShouldBindJSON(&updatedMovie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie data"})
		return
	}

	if err := validateMovie(updatedMovie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := db.GetCollection()
	ctx, cancel := dbContext()
	defer cancel()

	updatedMovie.ID = primitive.NilObjectID

	result, err := collection.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{
		"$set": updatedMovie,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update movie"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movie updated"})
}

func DeleteMovie(c *gin.Context) {
	movieID := c.Param("id")
	oid, err := primitive.ObjectIDFromHex(movieID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	collection := db.GetCollection()
	ctx, cancel := dbContext()
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete movie"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movie deleted"})
}

func BookSeat(c *gin.Context) {
	movieID := c.Param("id")

	oid, err := primitive.ObjectIDFromHex(movieID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	var request models.BookSeatRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking request"})
		return
	}

	collection := db.GetCollection()
	ctx, cancel := dbContext()
	defer cancel()

	var movie models.Movie
	if err := collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&movie); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load movie"})
		return
	}

	if request.ShowtimeIndex >= len(movie.Showtimes) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid showtime index"})
		return
	}

	seatIndex := -1
	for index, seat := range movie.Showtimes[request.ShowtimeIndex].Seats {
		if seat.Number == request.SeatNumber {
			seatIndex = index
			if seat.Booked {
				c.JSON(http.StatusConflict, gin.H{"error": "Seat already booked"})
				return
			}
			break
		}
	}

	if seatIndex == -1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Seat number not found for selected showtime"})
		return
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{
			"_id": oid,
			"showtimes": bson.M{
				"$elemMatch": bson.M{
					"seats": bson.M{
						"$elemMatch": bson.M{
							"number": request.SeatNumber,
							"booked": false,
						},
					},
				},
			},
		},
		bson.M{
			"$set": bson.M{
				"showtimes." + strconv.Itoa(request.ShowtimeIndex) + ".seats." + strconv.Itoa(seatIndex) + ".booked": true,
			},
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to book seat"})
		return
	}

	if result.ModifiedCount == 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Seat could not be booked"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Seat booked"})
}

func validateMovie(movie models.Movie) error {
	for showtimeIndex, showtime := range movie.Showtimes {
		seenSeatNumbers := make(map[int]struct{}, len(showtime.Seats))

		for _, seat := range showtime.Seats {
			if _, exists := seenSeatNumbers[seat.Number]; exists {
				return fmt.Errorf("duplicate seat number %d in showtime %d", seat.Number, showtimeIndex)
			}

			seenSeatNumbers[seat.Number] = struct{}{}
		}
	}

	return nil
}
