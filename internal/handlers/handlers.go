package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rasyad91/goBookings/internal/config"
	"github.com/rasyad91/goBookings/internal/models"
	"github.com/rasyad91/goBookings/internal/render"
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
	render.Templates(rw, r, "home.page.html", &models.TemplateData{})
}

// About is the about page handler
func (rp *Repository) About(rw http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again"

	remoteIP := rp.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	render.Templates(rw, r, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Contact renders the make-a-reservation page and displays form
func (rp *Repository) Contact(rw http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	rp.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	render.Templates(rw, r, "contact.page.html", &models.TemplateData{})
}

// Reservation renders the make-a-reservation page and displays form
func (rp *Repository) Reservation(rw http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	rp.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	render.Templates(rw, r, "make-reservations.page.html", &models.TemplateData{})
}

// Availability renders the search availability page and displays form
func (rp *Repository) Availability(rw http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	rp.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	render.Templates(rw, r, "search-availability.page.html", &models.TemplateData{})
}

// PostAvailability renders the search availability page and displays form
func (rp *Repository) PostAvailability(rw http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	rw.Write([]byte(fmt.Sprintf("start date is : %s | end date is : %s", start, end)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON renders the search availability page and displays form
func (rp *Repository) AvailabilityJSON(rw http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{OK: true, Message: "Available"}

	out, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(out))
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(out)
}

// Generals is the General's Quarters page handler
func (rp *Repository) Generals(rw http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	rp.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	render.Templates(rw, r, "generals.page.html", &models.TemplateData{})
}

// Majors is the Major's Suites page handler
func (rp *Repository) Majors(rw http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	rp.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	render.Templates(rw, r, "majors.page.html", &models.TemplateData{})
}
