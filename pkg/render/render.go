package render

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/rasyad91/goBookings/pkg/config"
	"github.com/rasyad91/goBookings/pkg/models"
)

var function = template.FuncMap{}

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

// AddDefaultData adds default data for models.TemplateData
func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// Templates using templates using text/template
func Templates(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	// get the template cache from the app config
	tc := app.TemplateCache

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal(fmt.Fprintf(w, "Error 404 Page not found"))
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td)

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
