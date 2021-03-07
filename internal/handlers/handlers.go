package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rasyad91/goBookings/internal/config"
	"github.com/rasyad91/goBookings/internal/forms"
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

	render.Templates(rw, r, "contact.page.html", &models.TemplateData{})
}

// Reservation renders the make-a-reservation page and displays form
func (rp *Repository) Reservation(rw http.ResponseWriter, r *http.Request) {

	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	fmt.Println("Reservation: ", data["reservation"])
	render.Templates(rw, r, "make-reservations.page.html",
		&models.TemplateData{
			Form: forms.New(nil),
			Data: data,
		})
}

// PostReservation handles the posting of a reservation form
func (rp *Repository) PostReservation(rw http.ResponseWriter, r *http.Request) {
	fmt.Println("PostReservation")
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		return
	}
	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email", "phone")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		fmt.Println("PostReservation: ", data["reservation"])

		render.Templates(rw, r, "make-reservations.page.html",
			&models.TemplateData{
				Form: form,
				Data: data,
			})
		return
	}

}

// Availability renders the search availability page and displays form
func (rp *Repository) Availability(rw http.ResponseWriter, r *http.Request) {

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

	render.Templates(rw, r, "generals.page.html", &models.TemplateData{})
}

// Majors is the Major's Suites page handler
func (rp *Repository) Majors(rw http.ResponseWriter, r *http.Request) {

	render.Templates(rw, r, "majors.page.html", &models.TemplateData{})
}
