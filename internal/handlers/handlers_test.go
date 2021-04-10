package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/rasyad91/goBookings/internal/models"
)

var varTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"generalsQuarters", "/generals-quarters", "GET", http.StatusOK},
	{"majorsSuites", "/majors-suites", "GET", http.StatusOK},
	{"searchAvailability", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	{"makeReservation", "/make-reservations", "GET", http.StatusOK},
}

func TestHandlers(t *testing.T) {

	routes := getRoutes()
	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()

	for _, e := range varTests {
		response, err := testServer.Client().Get(testServer.URL + e.url)
		if err != nil {
			t.Fatal("test failed:", err)
		}
		if response.StatusCode != e.expectedStatusCode {
			t.Errorf("For: %s | Expected error code: %d | Actual error code: %d", e.name, e.expectedStatusCode, response.StatusCode)
		}

	}
}

func TestRepository_Reservation(t *testing.T) {
	// test case where reservation is in session (reset everything)
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	r, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(r)
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Reservation)

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code: actual %d, expected %d", w.Code, http.StatusOK)
	}
	//------------------------------------------------------------------------------------------------------------//
	// Test when room ID is valid,
	r, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(r)
	r = r.WithContext(ctx)

	w = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	reservation.RoomID = -1
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: actual %d, expected %d", w.Code, http.StatusTemporaryRedirect)
	}

	//------------------------------------------------------------------------------------------------------------//

	// test case where reservation is not in session (reset everything)
	r, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(r)
	r = r.WithContext(ctx)

	w = httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: actual %d, expected %d", w.Code, http.StatusTemporaryRedirect)
	}
	//------------------------------------------------------------------------------------------------------------//
}

func TestRepository_PostReservation(t *testing.T) {
	// test case : Sanity postReservation
	startDate := "start_date=2050-01-01"
	endDate := "end_date=2050-01-05"
	firstname := "first_name=John"
	lastname := "last_name=Boo"
	email := "email=john@bo.co"
	phone := "phone=12313213221"
	roomID := "room_id=1"

	postedData := url.Values{}
	postedData.Add("start_date", "2050-01-01")
	postedData.Add("end_date", "2050-01-04")
	postedData.Add("first_name", "John")
	postedData.Add("last_name", "Doe")
	postedData.Add("email", "john@bo.co")
	postedData.Add("phone", "12313213221")
	postedData.Add("room_id", "1")

	r, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx := getCtx(r)
	r = r.WithContext(ctx)

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	handlerfunc := http.HandlerFunc(Repo.PostReservation)
	handlerfunc.ServeHTTP(w, r)

	if w.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code: actual %d, expected %d", w.Code, http.StatusSeeOther)
	}

	//------------------------------------------------------------------------------------------------------------//
	//test case: Invalid start date
	startDate = "start_date=invalid"
	requestBody := fmt.Sprintf("%s&%s&%s&%s&%s&%s&%s",
		startDate,
		endDate,
		firstname,
		lastname,
		email,
		phone,
		roomID)

	r, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(requestBody))
	ctx = getCtx(r)
	r = r.WithContext(ctx)

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w = httptest.NewRecorder()
	handlerfunc.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("PostReservation handler returned wrong response code for missing post body: actual %d, expected %d", w.Code, http.StatusBadRequest)
	}
	//------------------------------------------------------------------------------------------------------------//
	//test case: Invalid end date
	startDate = "start_date=2050-01-01"
	endDate = "end_date=invalid"
	requestBody = fmt.Sprintf("%s&%s&%s&%s&%s&%s&%s",
		startDate,
		endDate,
		firstname,
		lastname,
		email,
		phone,
		roomID)

	r, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(requestBody))
	ctx = getCtx(r)
	r = r.WithContext(ctx)

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w = httptest.NewRecorder()
	handlerfunc.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("PostReservation handler returned wrong response code for missing post body: actual %d, expected %d", w.Code, http.StatusBadRequest)
	}
	//------------------------------------------------------------------------------------------------------------//
	//test case: Invalid room_id date
	endDate = "end_date=2050-01-05"
	roomID = "room_id=invalid"
	requestBody = fmt.Sprintf("%s&%s&%s&%s&%s&%s&%s",
		startDate,
		endDate,
		firstname,
		lastname,
		email,
		phone,
		roomID)

	r, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(requestBody))
	ctx = getCtx(r)
	r = r.WithContext(ctx)

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w = httptest.NewRecorder()
	handlerfunc.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("PostReservation handler returned wrong response code for missing post body: actual %d, expected %d", w.Code, http.StatusBadRequest)
	}
	//------------------------------------------------------------------------------------------------------------//
	//test case: Invalid data
	firstname = "first_name=J"
	roomID = "room_id=1"
	requestBody = fmt.Sprintf("%s&%s&%s&%s&%s&%s&%s",
		startDate,
		endDate,
		firstname,
		lastname,
		email,
		phone,
		roomID)

	r, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(requestBody))
	ctx = getCtx(r)
	r = r.WithContext(ctx)

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w = httptest.NewRecorder()
	handlerfunc.ServeHTTP(w, r)

	if w.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code for missing post body: actual %d, expected %d", w.Code, http.StatusSeeOther)
	}
	//------------------------------------------------------------------------------------------------------------//
	// test case: Missing post body

	r, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(r)
	r = r.WithContext(ctx)

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w = httptest.NewRecorder()

	handlerfunc = http.HandlerFunc(Repo.PostReservation)
	handlerfunc.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("PostReservation handler returned wrong response code for missing post body: actual %d, expected %d", w.Code, http.StatusBadRequest)
	}
	//test case: test insertReservation
	firstname = "first_name=Joe"
	roomID = "room_id=2"
	requestBody = fmt.Sprintf("%s&%s&%s&%s&%s&%s&%s",
		startDate,
		endDate,
		firstname,
		lastname,
		email,
		phone,
		roomID)

	r, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(requestBody))
	ctx = getCtx(r)
	r = r.WithContext(ctx)

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w = httptest.NewRecorder()
	handlerfunc.ServeHTTP(w, r)

	if w.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for missing post body: actual %d, expected %d", w.Code, http.StatusTemporaryRedirect)
	}
	//------------------------------------------------------------------------------------------------------------//
	//test case: Insert restriction
	roomID = "room_id=100"
	requestBody = fmt.Sprintf("%s&%s&%s&%s&%s&%s&%s",
		startDate,
		endDate,
		firstname,
		lastname,
		email,
		phone,
		roomID)

	r, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(requestBody))
	ctx = getCtx(r)
	r = r.WithContext(ctx)

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w = httptest.NewRecorder()
	handlerfunc.ServeHTTP(w, r)

	if w.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for missing post body: actual %d, expected %d", w.Code, http.StatusTemporaryRedirect)
	}
	//------------------------------------------------------------------------------------------------------------//

}

func getCtx(r *http.Request) context.Context {
	ctx, err := session.Load(r.Context(), r.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}

func TestRepository_AvailabilityJSON(t *testing.T) {

	startDate := "start_date=2050-01-01"
	endDate := "end_date=2050-01-05"
	roomID := "room_id=1"
	requestBody := fmt.Sprintf("%s&%s&%s",
		startDate,
		endDate,
		roomID)

	r, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(requestBody))
	ctx := getCtx(r)
	r = r.WithContext(ctx)

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	handlerfunc := http.HandlerFunc(Repo.AvailabilityJSON)
	handlerfunc.ServeHTTP(w, r)

	var j jsonResponse

	err := json.Unmarshal([]byte(w.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json")
	}
}
