package models

import (
	"database/sql"
	"time"
)

type Models struct {
	DB DBModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		DB: DBModel{DB: db},
	}
}

type Movie struct {
	ID          int            `json:"id" bson:"id"`
	Title       string         `json:"title" bson:"title"`
	Description string         `json:"description" bson:"description"`
	Year        int            `json:"year" bson:"year"`
	ReleaseDate time.Time      `json:"release_date" bson:"release_date"`
	Runtime     int            `json:"runtime" bson:"runtime"`
	Rating      string         `json:"rating" bson:"rating"`
	CreatedAt   time.Time      `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" bson:"updated_at"`
	MovieGenre  map[int]string `json:"genres" bson:"genres"`
	Poster      string         `json:"poster" bson:"poster"`
}

type Genre struct {
	ID        int       `json:"id" bson:"id"`
	Name      string    `json:"name" bson:"name"`
	CreatedAt time.Time `json:"-" bson:"-"`
	UpdatedAt time.Time `json:"-" bson:"-"`
}

type MovieGenre struct {
	ID        int       `json:"id" bson:"id"`
	MovieID   int       `json:"-" bson:"-"`
	GenreID   int       `json:"-" bson:"-"`
	Genre     Genre     `json:"genre" bson:"genre"`
	CreatedAt time.Time `json:"-" bson:"-"`
	UpdatedAt time.Time `json:"-" bson:"-"`
}

type User struct {
	ID       int    `json:"id" bson:"id"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}
