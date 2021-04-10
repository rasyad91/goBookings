package handlers

import (
	"net/http"

	"github.com/rasyad91/goBookings/internal/forms"
	"github.com/rasyad91/goBookings/internal/models"
	"github.com/rasyad91/goBookings/internal/render"
)

func (m *Repository) AdminNewReservation(w http.ResponseWriter, r *http.Request) {
	render.Templates(w, r, "admin-new-reservation.page.html", &models.TemplateData{Form: forms.New(nil)})

}

func (m *Repository) AdminAllReservation(w http.ResponseWriter, r *http.Request) {
	render.Templates(w, r, "admin-all-reservations.page.html", &models.TemplateData{Form: forms.New(nil)})

}

func (m *Repository) AdminCalendarReservation(w http.ResponseWriter, r *http.Request) {
	render.Templates(w, r, "admin-calendar-reservations.page.html", &models.TemplateData{Form: forms.New(nil)})

}
