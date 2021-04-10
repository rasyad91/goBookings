package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/rasyad91/goBookings/internal/config"
	"github.com/rasyad91/goBookings/internal/driver"
	"github.com/rasyad91/goBookings/internal/handlers"
	"github.com/rasyad91/goBookings/internal/helpers"
	"github.com/rasyad91/goBookings/internal/models"
	"github.com/rasyad91/goBookings/internal/render"
)

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

const port = ":8080"

func main() {

	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	defer close(app.MailChan)

	fmt.Println("Starting mail listener")
	listenForMail()

	fmt.Printf("Starting at port %s\n", port)
	//log.Fatalln(http.ListenAndServe(port, nil))

	srv := &http.Server{
		Addr:    port,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	// what am i going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.Room{})
	gob.Register(models.User{})
	gob.Register(models.RoomRestriction{})

	mailChain := make(chan models.MailData)
	app.MailChan = mailChain

	// change this to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	//connect to database
	app.InfoLog.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=rasyad password=")
	if err != nil {
		app.ErrorLog.Fatalf("fail to connect to DB: %v", err)
	}
	app.InfoLog.Println("Connected to database...")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		return db, fmt.Errorf("fail to create Template cache: %w", err)
	}

	app.TemplateCache = tc
	app.UseCache = false

	r := handlers.NewRepo(&app, db)
	handlers.NewHandlers(r)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
