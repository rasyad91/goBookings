package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var varTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"generalsQuarters", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"majorsSuites", "/majors-suites", "GET", []postData{}, http.StatusOK},
	{"searchAvailability", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"makeReservation", "/make-reservations", "GET", []postData{}, http.StatusOK},
	{"postSearchAvailabilityOK", "/search-availability", "POST", []postData{
		{key: "start", value: "2021-05-05"},
		{key: "end", value: "2021-05-07"},
	}, http.StatusOK},
	// {"postSearchAvailabilityFAIL", "/search-availability", "POST", []postData{
	// 	{key: "start", value: "2021-05-05"},
	// 	{key: "end", value: "2021-05-02"},
	// }, http.StatusNotFound},
	{"postSearchAvailabilityJson", "/search-availability-json", "POST", []postData{
		{key: "start", value: "2021-05-05"},
		{key: "end", value: "2021-05-07"},
	}, http.StatusOK},
	{"make-reservations", "/make-reservations", "POST", []postData{
		{key: "first_name", value: "John"},
		{key: "last_name", value: "Smith"},
		{key: "email", value: "me@here.com"},
		{key: "phone", value: "555-55-555"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {

	routes := getRoutes()
	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()

	for _, e := range varTests {
		if e.method == "GET" {
			response, err := testServer.Client().Get(testServer.URL + e.url)
			if err != nil {
				t.Fatal("test failed:", err)
			}
			if response.StatusCode != e.expectedStatusCode {
				t.Errorf("For: %s | Expected error code: %d | Actual error code: %d", e.name, e.expectedStatusCode, response.StatusCode)
			}
		} else {
			values := url.Values{}
			for _, x := range e.params {
				values.Add(x.key, x.value)
			}
			response, err := testServer.Client().PostForm(testServer.URL+e.url, values)
			if err != nil {
				t.Fatal("test failed:", err)
			}
			if response.StatusCode != e.expectedStatusCode {
				t.Errorf("For: %s | Expected error code: %d | Actual error code: %d", e.name, e.expectedStatusCode, response.StatusCode)
			}
		}
	}
}
