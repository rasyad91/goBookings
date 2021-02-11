package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/rasyad91/goBookings/pkg/config"
	"github.com/rasyad91/goBookings/pkg/handlers"
	"github.com/rasyad91/goBookings/pkg/render"
)

var app config.AppConfig

const port = ":8080"

var session *scs.SessionManager

func main() {
	// change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	render.NewTemplates(&app)

	r := handlers.NewRepo(&app)
	handlers.NewHandlers(r)

	//http.HandleFunc("/", handlers.Repo.Home)
	//http.HandleFunc("/about", handlers.Repo.About)

	fmt.Printf("Starting at port %s", port)
	//log.Fatalln(http.ListenAndServe(port, nil))

	srv := &http.Server{
		Addr:    port,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}