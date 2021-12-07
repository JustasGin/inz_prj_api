package main

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) wrap(next http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.WithValue(r.Context(), "params", ps)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (app *application) routes() http.Handler {
	router := httprouter.New()
	secure := alice.New(app.checkToken)

	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)
	router.HandlerFunc(http.MethodPost, "/v1/signin", app.signIn)
	router.HandlerFunc(http.MethodPost, "/v1/checkup", app.checkUp)

	router.HandlerFunc(http.MethodGet, "/v1/movies", app.getMovies)
	router.HandlerFunc(http.MethodGet, "/v1/movie/:id", app.getMovie)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:genre_id", app.getMoviesByGenre)

	router.HandlerFunc(http.MethodGet, "/v1/genres", app.getGenres)

	router.POST("/v1/admin/movie/function", app.wrap(secure.ThenFunc(app.addOrEditMovie)))
	router.POST("/v1/admin/movie/delete", app.wrap(secure.ThenFunc(app.deleteMovie)))

	return app.enableCORS(router)
}
