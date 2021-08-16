package main

import (
	"context"
	"net/http"

	"github.com/justinas/alice"

	"github.com/julienschmidt/httprouter"
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

	router.HandlerFunc(http.MethodPost, "/v1/graphql", app.moviesGraphQL)

	router.HandlerFunc(http.MethodGet, "/v1/movie/:id", app.getOneMovie)
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.getAllMovies)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:genre_id", app.getAllMoviesByGenre)

	router.POST("/v1/admin/editmovie", app.wrap(secure.ThenFunc(app.editMovie)))
	// router.HandlerFunc(http.MethodPost, "/v1/admin/editmovie", app.editMovie)

	router.HandlerFunc(http.MethodGet, "/v1/genres", app.getAllGenres)
	router.HandlerFunc(http.MethodDelete, "/v1/admin/movie/:id", app.deleteMovie)
	return app.enableCORS(router)
}
