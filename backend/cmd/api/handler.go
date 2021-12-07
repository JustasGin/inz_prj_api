package main

import (
	"backend/models"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

func (app *application) statusHandler(w http.ResponseWriter, r *http.Request) {
	currentStatus := AppStatus{
		Status:      "Available",
		Environment: app.config.env,
		Version:     version,
	}

	err := app.writeJSON(w, http.StatusOK, currentStatus, "status")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *application) getMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := app.models.DB.GetMoviesDB()
	if err != nil {
		return
	}
	err = app.writeJSON(w, http.StatusOK, movies, "movies")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *application) getMovie(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.errorJSON(w, err, http.StatusConflict)
		return
	}

	movie, err := app.models.DB.GetMovieDB(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, movie, "movie")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *application) addOrEditMovie(w http.ResponseWriter, r *http.Request) {
	var payload MoviePayload
	var movie models.Movie

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	movie.ID, err = strconv.Atoi(payload.ID)
	if err != nil {
		app.errorJSON(w, err, http.StatusConflict)
		return
	}

	if payload.ID != "0" {
		id, err := strconv.Atoi(payload.ID)
		if err != nil {
			app.errorJSON(w, err, http.StatusConflict)
			return
		}
		m, err := app.models.DB.GetMovieDB(id)
		if err != nil {
			app.errorJSON(w, err)
			return
		}
		movie = *m
		movie.UpdatedAt = time.Now()
	}

	movie.Title = payload.Title
	movie.Description = payload.Description
	movie.ReleaseDate, err = time.Parse("2006-01-02", payload.ReleaseDate)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	movie.Year = movie.ReleaseDate.Year()
	movie.Runtime, err = strconv.Atoi(payload.Runtime)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	movie.Rating = payload.Rating
	movie.CreatedAt = time.Now()
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	movie.UpdatedAt = time.Now()

	movie = getPoster(movie)

	if movie.ID == 0 {
		err = app.models.DB.InsertMovieDB(movie)
		if err != nil {
			app.errorJSON(w, err)
			return
		}

		err = app.writeJSON(w, http.StatusOK, Response{message: fmt.Sprintf("Added the movie %s", movie.Title)}, "response")
		if err != nil {
			app.errorJSON(w, err)
			return
		}
	} else {
		err = app.models.DB.EditMovieDB(movie)
		if err != nil {
			app.errorJSON(w, err)
			return
		}

		err = app.writeJSON(w, http.StatusOK, Response{message: fmt.Sprintf("Updated the movie %s", movie.Title)}, "response")
		if err != nil {
			app.errorJSON(w, err)
			return
		}
	}
}

func (app *application) deleteMovie(w http.ResponseWriter, r *http.Request) {
	var payload MoviePayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	id, err := strconv.Atoi(payload.ID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.models.DB.DeleteMovieDB(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, Response{message: "OK"}, "response")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *application) getMoviesByGenre(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	genreId, err := strconv.Atoi(params.ByName("genre_id"))
	if err != nil {
		app.errorJSON(w, err, http.StatusConflict)
		return
	}

	movies, err := app.models.DB.GetMoviesDB(genreId)
	if err != nil {
		return
	}
	err = app.writeJSON(w, http.StatusOK, movies, "movies")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *application) getGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := app.models.DB.GetGenresDB()
	if err != nil {
		return
	}
	err = app.writeJSON(w, http.StatusOK, genres, "genres")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func getPoster(movie models.Movie) models.Movie {
	type MovieDB struct {
		Page    int `json:"page"`
		Results []struct {
			Adult            bool    `json:"adult"`
			BackdropPath     string  `json:"backdrop_path"`
			GenreIds         []int   `json:"genre_ids"`
			ID               int     `json:"id"`
			OriginalLanguage string  `json:"original_language"`
			OriginalTitle    string  `json:"original_title"`
			Overview         string  `json:"overview"`
			Popularity       float64 `json:"popularity"`
			PosterPath       string  `json:"poster_path"`
			ReleaseDate      string  `json:"release_date"`
			Title            string  `json:"title"`
			Video            bool    `json:"video"`
			VoteAverage      float64 `json:"vote_average"`
			VoteCount        int     `json:"vote_count"`
		} `json:"results"`
		TotalPages   int `json:"total_pages"`
		TotalResults int `json:"total_results"`
	}

	client := &http.Client{}
	key := os.Getenv("MOVIEDB_KEY")
	dbUrl := fmt.Sprintf("https://api.themoviedb.org/3/search/movie?api_key=%s&query=%s", key, url.QueryEscape(movie.Title))

	req, err := http.NewRequest("GET", dbUrl, nil)
	if err != nil {
		log.Println(err)
		return movie
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return movie
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return movie
	}

	var responseObject MovieDB
	json.Unmarshal(bodyBytes, &responseObject)
	if len(responseObject.Results) > 0 {
		movie.Poster = responseObject.Results[0].PosterPath
	}

	return movie
}
