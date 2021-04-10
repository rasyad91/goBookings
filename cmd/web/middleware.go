package main

import (
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/rasyad91/goBookings/internal/helpers"
)

// NoSurf is a middleware adds CRSF protection to all POST requests
func NoSurfMiddleware(next http.Handler) http.Handler {
	crsfHandler := nosurf.New(next)
	crsfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return crsfHandler
}

// SessionLoad is a middleware that loads and saves the session on every request
func SessionLoadMiddleware(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.IsAuthenticated(r) {
			session.Put(r.Context(), "error", "Log in first!")
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
