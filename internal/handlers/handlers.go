package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rasyad91/goBookings/internal/config"
	"github.com/rasyad91/goBookings/internal/forms"
	"github.com/rasyad91/goBookings/internal/helpers"
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
func (rp *Repository) Home(w http.ResponseWriter, r *http.Request) {

	render.Templates(w, r, "home.page.html", &models.TemplateData{})
}

// About is the about page handler
func (rp *Repository) About(w http.ResponseWriter, r *http.Request) {

	render.Templates(w, r, "about.page.html", &models.TemplateData{})
}

// Contact renders the make-a-reservation page and displays form
func (rp *Repository) Contact(w http.ResponseWriter, r *http.Request) {

	render.Templates(w, r, "contact.page.html", &models.TemplateData{})
}

// Reservation renders the make-a-reservation page and displays form
func (rp *Repository) Reservation(w http.ResponseWriter, r *http.Request) {

	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	fmt.Println("Reservation: ", data["reservation"])
	render.Templates(w, r, "make-reservations.page.html",
		&models.TemplateData{
			Form: forms.New(nil),
			Data: data,
		})
}

// PostReservation handles the posting of a reservation form
func (rp *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PostReservation")
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
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
	form.IsEmail()

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		fmt.Println("PostReservation: ", data["reservation"])

		render.Templates(w, r, "make-reservations.page.html",
			&models.TemplateData{
				Form: form,
				Data: data,
			})
		return
	}
	rp.App.Session.Put(r.Context(), "reservation", reservation)

	// http.statusseeother for redirecting
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

}

// ReservationSummary shows summary after user submitted request for reservation at PostReservation
func (rp *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {

	// needs to type assert
	reservation, ok := rp.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		rp.App.ErrorLog.Println("Cant get error from session")
		rp.App.Session.Put(r.Context(), "error", "Cannot get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	rp.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Templates(w, r, "reservation-summary.page.html",
		&models.TemplateData{
			Data: data,
		})

}

// Availability renders the search availability page and displays form
func (rp *Repository) Availability(w http.ResponseWriter, r *http.Request) {

	render.Templates(w, r, "search-availability.page.html", &models.TemplateData{})
}

// PostAvailability renders the search availability page and displays form
func (rp *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	w.Write([]byte(fmt.Sprintf("start date is : %s | end date is : %s", start, end)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON renders the search availability page and displays form
func (rp *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{OK: true, Message: "Available"}

	out, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	fmt.Println(string(out))
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Generals is the General's Quarters page handler
func (rp *Repository) Generals(w http.ResponseWriter, r *http.Request) {

	render.Templates(w, r, "generals.page.html", &models.TemplateData{})
}

// Majors is the Major's Suites page handler
func (rp *Repository) Majors(w http.ResponseWriter, r *http.Request) {

	render.Templates(w, r, "majors.page.html", &models.TemplateData{})
}
