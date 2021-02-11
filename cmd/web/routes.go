package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/rasyad91/goBookings/pkg/config"
	"github.com/rasyad91/goBookings/pkg/handlers"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)

	return mux
}
