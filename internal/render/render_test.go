package render

import (
	"net/http"
	"testing"

	"github.com/rasyad91/goBookings/internal/models"
)

var varTest = []struct {
	name         string
	tmpl         string
	templateData *models.TemplateData
}{
	{"home", "home.page.html", &models.TemplateData{}},
	{"about", "about.page.html", &models.TemplateData{}},
	{"contact", "contact.page.html", &models.TemplateData{}},
	{"availability", "search-availability.page.html", &models.TemplateData{}},
	{"generals", "generals.page.html", &models.TemplateData{}},
	{"major", "majors.page.html", &models.TemplateData{}},
}

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	session.Put(r.Context(), "flash", "123")
	session.Put(r.Context(), "warning", "123")
	session.Put(r.Context(), "error", "123")

	result := AddDefaultData(&td, r)
	if result.Flash != "123" {
		t.Error("Flash value of 123 not found in session")
	}
	if result.Warning != "123" {
		t.Error("Warning value of 123 not found in session")
	}
	if result.Error != "123" {
		t.Error("Error value of 123 not found in session")
	}
}

func TestRenderTemplate(t *testing.T) {
	pathToTemplates = "../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = tc

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}
	var w testWriter
	for _, e := range varTest {
		err = Templates(&w, r, e.tmpl, e.templateData)
		if err != nil {
			t.Errorf("test: %s | error: %v", e.name, err)
		}

	}
}

func TestNewTemplates(t *testing.T) {
	NewRenderer(app)
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)
	return r, nil

}
