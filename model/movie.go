package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Movie struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string             `bson:"title" json:"title" binding:"required"`
	Director    string             `bson:"director" json:"director" binding:"required"`
	Year        int                `bson:"year" json:"year" binding:"required,min=1888"`
	Description string             `bson:"description" json:"description" binding:"required"`
	Showtimes   []Showtime         `bson:"showtimes" json:"showtimes" binding:"required,min=1,dive"`
}

type Showtime struct {
	Time  string `bson:"time" json:"time" binding:"required"`
	Seats []Seat `bson:"seats" json:"seats" binding:"required,min=1,dive"`
}

type Seat struct {
	Number   int  `bson:"number" json:"number" binding:"required,min=1"`
	Booked   bool `bson:"booked" json:"booked"`
	Reserved bool `bson:"reserved" json:"reserved"`
}

type BookSeatRequest struct {
	ShowtimeIndex int `json:"showtime_index" binding:"min=0"`
	SeatNumber    int `json:"seat_number" binding:"required,min=1"`
}
