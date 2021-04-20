package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Errorf("Got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields are missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	r, _ = http.NewRequest("POST", "/", nil)
	r.PostForm = postedData
	form = New(r.PostForm)

	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("shows does not have required fields when it does")
	}
}

func TestForm_IsEmail(t *testing.T) {
	r, _ := http.NewRequest("POST", "/", nil)
	postedData := url.Values{}
	validEmail := "rasyad@gmail.com"
	postedData.Add("email", validEmail)
	r.PostForm = postedData
	form := New(r.PostForm)

	form.IsEmail()
	if _, ok := form.Errors["email"]; ok {
		t.Error("Expected no error for: ", validEmail, form.Errors)
	}

	invalidEmail := "rasyadgmailcom"
	postedData.Del("email")
	postedData.Add("email", invalidEmail)
	form.IsEmail()
	if _, ok := form.Errors["email"]; !ok {
		t.Errorf("Expected error for: %s, e: %v", postedData.Get("email"), form.Errors)
	}

}

func TestForm_MinLength(t *testing.T) {
	r, _ := http.NewRequest("POST", "/", nil)

	postedData := url.Values{}
	postedData.Add("name", "AAA")
	r.PostForm = postedData
	form := New(r.PostForm)

	if ok := form.MinLength("name", 2); ok == false {
		t.Error("Expected no error, but theres error")
	}

	if ok := form.MinLength("name", 4); ok == true {
		t.Error("Expected error", r.PostForm.Get("name"), len(r.PostForm.Get("name")))
	}

}

func TestForm_Has(t *testing.T) {
	r, _ := http.NewRequest("POST", "/", nil)
	form := New(r.PostForm)

	if x := form.Has("Hello"); x == true {
		t.Error("Expected false")
	}

	postedData := url.Values{}
	postedData.Add("Hello", "Value")
	form = New(postedData)

	if x := form.Has("Hello"); x == false {
		t.Error("Expected false")
	}
}
