package render

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/justinas/nosurf"
	"github.com/rasyad91/goBookings/internal/config"
	"github.com/rasyad91/goBookings/internal/models"
)

var function = template.FuncMap{}

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

// AddDefaultData adds default data for models.TemplateData
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")

	td.CSRFToken = nosurf.Token(r)
	return td
}

// Templates using templates using text/template
func Templates(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	// get the template cache from the app config

	if !app.UseCache {
		tc, err := CreateTemplateCache()
		if err != nil {
			log.Fatal("cannot create template cache")
		}
		app.TemplateCache = tc
	}

	tc := app.TemplateCache

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal(fmt.Fprintf(w, "Error 404 Page not found"))
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	err := t.Execute(buf, td)
	if err != nil {
		log.Fatal(err)
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser", err)
	}
}

// CreateTemplateCache creates templates for render and store in cache
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}
	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil {
		return myCache, err
	}
	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(function).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.html")
		if err != nil {
			return myCache, err
		}
		if len(matches) == 0 {
			return myCache, fmt.Errorf("Base template not found")
		}

		ts, err = ts.ParseGlob("./templates/*.layout.html")
		myCache[name] = ts

	}
	return myCache, nil
}
