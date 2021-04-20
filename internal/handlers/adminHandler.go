package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/rasyad91/goBookings/internal/forms"
	"github.com/rasyad91/goBookings/internal/helpers"
	"github.com/rasyad91/goBookings/internal/models"
	"github.com/rasyad91/goBookings/internal/render"
)

const (
	dateFormatMonth = "01"
	dateFormatYear  = "2006"
)

func (m *Repository) AdminCalendarReservation(w http.ResponseWriter, r *http.Request) {
	// assume theres no month/year specified
	now := time.Now()

	if r.URL.Query().Get("y") != "" {
		year, err := strconv.Atoi(r.URL.Query().Get("y"))
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		month, err := strconv.Atoi(r.URL.Query().Get("m"))
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		now = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	}

	data := make(map[string]interface{})
	data["now"] = now

	next := now.AddDate(0, 1, 0)
	prev := now.AddDate(0, -1, 0)

	nextMonth := next.Format(dateFormatMonth)
	nextMonthYear := next.Format(dateFormatYear)

	prevMonth := prev.Format(dateFormatMonth)
	prevMonthYear := prev.Format(dateFormatYear)

	stringMap := make(map[string]string)
	stringMap["next_month"] = nextMonth
	stringMap["next_month_year"] = nextMonthYear
	stringMap["prev_month"] = prevMonth
	stringMap["prev_month_year"] = prevMonthYear

	stringMap["this_month"] = now.Format(dateFormatMonth)
	stringMap["this_month_year"] = now.Format(dateFormatYear)

	// how to get the days in the month
	// get the first and last days of the month
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	intMap := make(map[string]int)
	intMap["days_in_month"] = lastOfMonth.Day()

	// get rooms
	rooms, err := m.DB.GetAllRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data["rooms"] = rooms

	for _, x := range rooms {
		reservationMap := make(map[string]int)
		blockMap := make(map[string]int)

		for d := firstOfMonth; !d.After(lastOfMonth); d = d.AddDate(0, 0, 1) {
			reservationMap[d.Format(datelayout)] = 0
			blockMap[d.Format(datelayout)] = 0
		}

		// get all the restrictions of the current room
		restrictions, err := m.DB.GetRestrictionsForRoomByDate(x.ID, firstOfMonth, lastOfMonth)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		for _, r := range restrictions {
			if r.ReservationID < 1 {
				blockMap[r.StartDate.Format(datelayout)] = r.ID
			} else {
				for d := r.StartDate; !d.After(r.EndDate); d = d.AddDate(0, 0, 1) {
					reservationMap[d.Format(datelayout)] = r.ReservationID
				}
			}
		}

		data[fmt.Sprintf("reservation_map_%d", x.ID)] = reservationMap
		data[fmt.Sprintf("block_map_%d", x.ID)] = blockMap

		m.App.Session.Put(r.Context(), fmt.Sprintf("block_map_%d", x.ID), blockMap)
	}

	render.Templates(w, r, "admin-calendar-reservations.page.html", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
		IntMap:    intMap,
	})

}

// AdminAllReservation shows all reservations in admin tool
func (m *Repository) AdminAllReservation(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.GetAllReservations()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "error getting reservations from db")
		http.Redirect(w, r, "/reservations-all", http.StatusInternalServerError)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations

	render.Templates(w, r, "admin-all-reservations.page.html", &models.TemplateData{
		Data: data,
		Form: forms.New(nil)})
}

// AdminNewReservation shows all new reservations in admin tool
func (m *Repository) AdminNewReservation(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.GetNewReservations()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "error getting reservations from db")
		http.Redirect(w, r, "/reservations-all", http.StatusInternalServerError)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations

	render.Templates(w, r, "admin-new-reservation.page.html", &models.TemplateData{
		Data: data,
		Form: forms.New(nil)})
}

// AdminShowReservation shows the reservation in the admin tool
func (m *Repository) AdminShowReservation(w http.ResponseWriter, r *http.Request) {
	// src := chi.URLParam(r, "src")
	split := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(split[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	// id, err := strconv.Atoi(chi.URLParam(r, "id"))
	src := split[3]

	stringMap := make(map[string]string)
	stringMap["src"] = src

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	stringMap["month"] = month
	stringMap["year"] = year

	reservation, err := m.DB.GetReservationByID(id)
	if err != nil {
		helpers.ServerError(w, err)
	}

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Templates(w, r, "admin-show-reservations.page.html", &models.TemplateData{
		Form:      forms.New(nil),
		StringMap: stringMap,
		Data:      data,
	})

}

func (m *Repository) AdminUpdateReservation(w http.ResponseWriter, r *http.Request) {

	split := strings.Split(r.RequestURI, "/")
	src := split[3]

	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}

	id, err := strconv.Atoi(r.PostFormValue("id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation, err := m.DB.GetReservationByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Phone = r.Form.Get("phone")
	reservation.Email = r.Form.Get("email")

	if err := m.DB.UpdateReservation(reservation); err != nil {
		helpers.ServerError(w, err)
		return
	}
	m.App.Session.Put(r.Context(), "flash", "Changes saved!")

	month := r.PostFormValue("month")
	year := r.PostFormValue("year")

	var s string
	if year == "" {
		s = fmt.Sprintf("/admin/reservations-%s", src)
	} else {
		s = fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month)
	}

	http.Redirect(w, r, s, http.StatusSeeOther)

}

func (m *Repository) AdminProcessReservation(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	src := chi.URLParam(r, "src")

	if err := m.DB.ProcessReservation(id); err != nil {
		helpers.ServerError(w, err)
		return
	}

	month := r.URL.Query().Get("m")
	year := r.URL.Query().Get("y")

	var s string
	if year == "" {
		s = fmt.Sprintf("/admin/reservations-%s", src)
	} else {
		s = fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month)
	}

	m.App.Session.Put(r.Context(), "flash", "Reservation processed!")
	http.Redirect(w, r, s, http.StatusSeeOther)
}

// AdminDeleteReservation deletes reservation
func (m *Repository) AdminDeleteReservation(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	src := chi.URLParam(r, "src")

	if err := m.DB.DeleteReservationByID(id); err != nil {
		helpers.ServerError(w, err)
		return
	}

	month := r.URL.Query().Get("m")
	year := r.URL.Query().Get("y")

	var s string
	if year == "" {
		s = fmt.Sprintf("/admin/reservations-%s", src)
	} else {
		s = fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month)
	}

	m.App.Session.Put(r.Context(), "flash", "Reservation deleted!")
	http.Redirect(w, r, s, http.StatusSeeOther)
}

func (m *Repository) AdminCalendarPostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	year, _ := strconv.Atoi(r.PostFormValue("y"))
	month, _ := strconv.Atoi(r.PostFormValue("m"))

	//Process blocks
	rooms, err := m.DB.GetAllRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	form := forms.New(r.PostForm)

	for _, room := range rooms {
		// get block map from the session.
		// then loop through the entire map,
		// if we have an entry in the map that does not exist in the posted data
		// and resrtiction id > 0, then it is a block we need to remove
		currentMap := m.App.Session.Get(r.Context(), fmt.Sprintf("block_map_%d", room.ID)).(map[string]int)

		for date, value := range currentMap {
			// ok will be false if the value is not in the map
			if val, ok := currentMap[date]; ok {
				// only pay attention to values > 0, and that are not in the form post
				// the rest are just placeholders for days without blocks
				if val > 0 {
					if !form.Has(fmt.Sprintf("remove_block_%d_%s", room.ID, date)) {
						//delete the restriction by id
						err := m.DB.DeleteBlockByID(value)
						if err != nil {
							log.Println(err)
						}
					}
				}
			}
		}
	}

	for name := range r.PostForm {
		if strings.HasPrefix(name, "add_block") {
			split := strings.Split(name, "_")
			roomID, _ := strconv.Atoi(split[2])
			date, err := time.Parse("2006-01-02", split[3])
			if err != nil {
				log.Println(err)
			}

			fmt.Println()
			if err := m.DB.InsertBlockForRoom(roomID, date); err != nil {
				log.Println(err)
			}
		}
	}

	m.App.Session.Put(r.Context(), "flash", "Changes saved")

	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%d&m=%d", year, month), http.StatusSeeOther)

}
