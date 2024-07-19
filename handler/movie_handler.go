package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	db "movie_booking_system/database"
	models "movie_booking_system/model"
)

func GetMovies(c *gin.Context) {
	collection := db.GetCollection()

	var movies []models.Movie
	cur, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get movies"})
		return
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var movie models.Movie
		err := cur.Decode(&movie)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode movie"})
			return
		}
		movies = append(movies, movie)
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

	var movie models.Movie
	err = collection.FindOne(context.TODO(), bson.M{"_id": oid}).Decode(&movie)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
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

	collection := db.GetCollection()

	result, err := collection.InsertOne(context.TODO(), movie)
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

	collection := db.GetCollection()

	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": oid}, bson.M{
		"$set": updatedMovie,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update movie"})
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

	_, err = collection.DeleteOne(context.TODO(), bson.M{"_id": oid})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete movie"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movie deleted"})
}
func BookSeat(c *gin.Context) {
	movieID := c.Param("id")
	showtimeIndex := c.Query("showtime_index")
	seatNumber := c.Query("seat_number")

	oid, err := primitive.ObjectIDFromHex(movieID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	collection := db.GetCollection()

	var movie models.Movie
	filter := bson.M{"_id": oid}
	update := bson.M{
		"$set": bson.M{
			"showtimes." + showtimeIndex + ".seats." + seatNumber + ".booked": true,
		},
	}

	err = collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&movie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to book seat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Seat booked"})
}
