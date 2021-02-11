package handlers

import (
	"net/http"

	"github.com/rasyad91/goBookings/pkg/config"
	"github.com/rasyad91/goBookings/pkg/models"
	"github.com/rasyad91/goBookings/pkg/render"
)

// Repository Pattern start
// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// Repo the repository used by the handlers
var Repo *Repository

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

//End of Repository pattern

// Home is the home page handler
func (rp *Repository) Home(rw http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	rp.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	render.Templates(rw, "home.page.html", &models.TemplateData{})

}

// About is the about page handler
func (rp *Repository) About(rw http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again"

	remoteIP := rp.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	render.Templates(rw, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}
