package models

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Movie struct {
    ID          primitive.ObjectID `bson:"_id,omitempty"`
    Title       string             `bson:"title"`
    Director    string             `bson:"director"`
    Year        int                `bson:"year"`
    Description string             `bson:"description"`
    Showtimes   []Showtime         `bson:"showtimes"`
}

type Showtime struct {
    Time  string `bson:"time"`
    Seats []Seat `bson:"seats"`
}

type Seat struct {
    Number   int  `bson:"number"`
    Booked   bool `bson:"booked"`
    Reserved bool `bson:"reserved"`
}
