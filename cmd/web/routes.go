package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/rasyad91/goBookings/internal/config"
	"github.com/rasyad91/goBookings/internal/handlers"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/generals-quarters", handlers.Repo.Generals)
	mux.Get("/majors-suites", handlers.Repo.Majors)

	mux.Get("/make-reservations", handlers.Repo.Reservation)
	mux.Post("/make-reservations", handlers.Repo.PostReservation)

	mux.Get("/search-availability", handlers.Repo.Availability)
	mux.Post("/search-availability-json", handlers.Repo.AvailabilityJSON)

	mux.Post("/search-availability", handlers.Repo.PostAvailability)

	mux.Get("/contact", handlers.Repo.Contact)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
