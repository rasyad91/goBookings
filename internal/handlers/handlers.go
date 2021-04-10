package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/rasyad91/goBookings/internal/config"
	"github.com/rasyad91/goBookings/internal/driver"
	"github.com/rasyad91/goBookings/internal/forms"
	"github.com/rasyad91/goBookings/internal/helpers"
	"github.com/rasyad91/goBookings/internal/models"
	"github.com/rasyad91/goBookings/internal/render"
	"github.com/rasyad91/goBookings/internal/repository"
	"github.com/rasyad91/goBookings/internal/repository/dbrepo"
)

var datelayout = "2006-01-02"

// Repository Pattern start
// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// Repo the repository used by the handlers
var Repo *Repository

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewTestRepo creates a new repository
func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestRepo(a),
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

	res, ok := rp.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		rp.App.Session.Put(r.Context(), "error", "Cannot get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	room, err := rp.DB.GetRoomByID(res.RoomID)
	if err != nil {
		rp.App.Session.Put(r.Context(), "error", "Cannot find room")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.Room = room

	rp.App.Session.Put(r.Context(), "reservation", res)

	sd := res.StartDate.Format(datelayout)
	ed := res.EndDate.Format(datelayout)

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Templates(w, r, "make-reservations.page.html",
		&models.TemplateData{
			Form:      forms.New(nil),
			Data:      data,
			StringMap: stringMap,
		})
}

// PostReservation handles the posting of a reservation form
func (rp *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	reservation, ok := rp.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, fmt.Errorf("cannot get reservation from session"))
		return
	}

	err := r.ParseForm()
	if err != nil {
		rp.App.Session.Put(r.Context(), "error", fmt.Sprint(err))
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse(datelayout, r.Form.Get("start_date"))
	if err != nil {
		rp.App.Session.Put(r.Context(), "error", fmt.Sprint(err))
		http.Redirect(w, r, "/make-reservations", http.StatusBadRequest)
		return
	}

	endDate, err := time.Parse(datelayout, r.Form.Get("end_date"))
	if err != nil {
		rp.App.Session.Put(r.Context(), "error", fmt.Sprint(err))
		http.Redirect(w, r, "/make-reservations", http.StatusBadRequest)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		rp.App.Session.Put(r.Context(), "error", fmt.Sprint(err))
		http.Redirect(w, r, "/make-reservations", http.StatusBadRequest)
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")
	reservation.StartDate = startDate
	reservation.EndDate = endDate
	reservation.RoomID = roomID

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email", "phone")
	form.MinLength("first_name", 3)
	form.IsEmail()

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		http.Error(w, "invalid form submission", http.StatusSeeOther)

		render.Templates(w, r, "make-reservations.page.html",
			&models.TemplateData{
				Form: form,
				Data: data,
			})
		return
	}

	newReservationID, err := rp.DB.InsertReservation(reservation)
	if err != nil {
		rp.App.Session.Put(r.Context(), "error", fmt.Sprint("Can't insert reservation into database: ", err))
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	restriction := models.RoomRestriction{
		RoomID:        reservation.RoomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
	}

	if err := rp.DB.InsertRoomRestriction(restriction); err != nil {
		rp.App.Session.Put(r.Context(), "error", fmt.Sprint("Can't insert room-restriction into database: ", err))
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// send email notification - to guest

	msg := models.MailData{
		To:       reservation.Email,
		From:     "me@here.com",
		Subject:  "Reservation Confirmation",
		Content:  "",
		Template: "basic.html",
	}
	msg.Content = fmt.Sprintf(`
		<strong>Reservation Confirmation</strong><br>
		Dear Mr/Ms %s, <br>
		This is to confirm your reservation from %s to %s. 
	`,
		reservation.LastName,
		reservation.StartDate.Format(datelayout),
		reservation.EndDate.Format(datelayout),
	)

	rp.App.MailChan <- msg

	// send email notification - to owner

	msg = models.MailData{
		To:       "me@here.com",
		From:     "me@here.com",
		Subject:  "Reservation Confirmation",
		Content:  "",
		Template: "basic.html",
	}
	msg.Content = fmt.Sprintf(`
			Guest: %s %s <br>
			RoomID: %d <br>
			Room: %s<br>
			From: %s<br>
			To: %s<br>
		`,
		reservation.FirstName,
		reservation.LastName,
		reservation.RoomID,
		reservation.Room.RoomName,
		reservation.StartDate.Format(datelayout),
		reservation.EndDate.Format(datelayout),
	)

	rp.App.MailChan <- msg

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
	fmt.Println("Reservation", reservation)

	fmt.Println("Reservation room: ", reservation.Room.RoomName)
	data := make(map[string]interface{})
	data["reservation"] = reservation

	sd := reservation.StartDate.Format(datelayout)
	ed := reservation.EndDate.Format(datelayout)

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Templates(w, r, "reservation-summary.page.html",
		&models.TemplateData{
			Data:      data,
			StringMap: stringMap,
		})
	rp.App.Session.Remove(r.Context(), "reservation")

}

// Availability renders the search availability page and displays form
func (rp *Repository) Availability(w http.ResponseWriter, r *http.Request) {

	render.Templates(w, r, "search-availability.page.html", &models.TemplateData{})
}

// PostAvailability renders the search availability page and displays form
func (rp *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	fmt.Println("Start", start)
	fmt.Println("End", end)

	startDate, err := time.Parse(datelayout, start)
	if err != nil {
		fmt.Println(err)
		// helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(datelayout, end)
	if err != nil {
		fmt.Println(err)
		// helpers.ServerError(w, err)
		return
	}

	rooms, err := rp.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
	}

	if len(rooms) == 0 {
		rp.App.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	rp.App.Session.Put(r.Context(), "reservation", res)

	render.Templates(w, r, "choose-room.page.html", &models.TemplateData{
		Data: data,
	})

}

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// AvailabilityJSON renders the search availability page and displays form
func (rp *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		// cant parse form, so return appropriate json
		resp := jsonResponse{
			OK:      false,
			Message: "Internal server error",
		}

		out, _ := json.MarshalIndent(resp, "", "    ")

		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))
	startDate, _ := time.Parse(datelayout, sd)
	endDate, _ := time.Parse(datelayout, ed)

	available, err := rp.DB.IsAvailableByDatesByRoomID(startDate, endDate, roomID)
	if err != nil {
		// cant parse form, so return appropriate json
		resp := jsonResponse{
			OK:      false,
			Message: "Error connecting to database",
		}

		out, _ := json.MarshalIndent(resp, "", "    ")

		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}
	resp := jsonResponse{
		OK:        available,
		Message:   "",
		StartDate: sd,
		EndDate:   ed,
		RoomID:    strconv.Itoa(roomID),
	}

	out, _ := json.MarshalIndent(resp, "", "    ")

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

// ChooseRoom displays list of available rooms
func (rp *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
	}

	res, ok := rp.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, fmt.Errorf("cannot get reservation from session"))
		return
	}

	res.RoomID = roomID
	rp.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservations", http.StatusSeeOther)
}

// BookRoom takes URL parameters, builds a sessional vaiable, and takes user to make reservation screen
func (rp *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	// id, s, e
	roomID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		helpers.ServerError(w, err)
	}
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	startDate, err := time.Parse(datelayout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(datelayout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	room, err := rp.DB.GetRoomByID(roomID)
	if err != nil {
		helpers.ServerError(w, err)
	}

	fmt.Println("Choose room: room:", room)

	reservation := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomID,
		Room:      room,
	}

	rp.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/make-reservations", http.StatusSeeOther)

}

func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.Templates(w, r, "login.page.html", &models.TemplateData{Form: forms.New(nil)})
}

func (m *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())
	err := r.ParseForm()
	if err != nil {
		m.App.ErrorLog.Println(err)
	}

	form := forms.New(r.PostForm)

	email := r.Form.Get("email")
	password := r.Form.Get("password")
	form.IsEmail()

	form.Required("email", "password")
	if !form.Valid() {

		render.Templates(w, r, "login.page.html", &models.TemplateData{Form: form})
		return
	}

	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "invalid login credentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "Logged in successfully")

	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)

}

// Logout logs a user out
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (m *Repository) AdminDashBoard(w http.ResponseWriter, r *http.Request) {
	render.Templates(w, r, "admin-dashboard.page.html", &models.TemplateData{Form: forms.New(nil)})

}
