package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

// NoSurf is a middleware adds CRSF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
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
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
