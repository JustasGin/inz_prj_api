package main

type Response struct {
	message string `json:"message" bson:"message"`
}

type MoviePayload struct {
	ID          string `json:"id" bson:"id"`
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
	Year        string `json:"year" bson:"year"`
	ReleaseDate string `json:"release_date" bson:"release_date"`
	Runtime     string `json:"runtime" bson:"runtime"`
	Rating      string `json:"rating" bson:"rating"`
	CreatedAt   string `json:"-" bson:"-"`
	UpdatedAt   string `json:"-" bson:"-"`
}
