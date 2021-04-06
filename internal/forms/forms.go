package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	valid "github.com/asaskevich/govalidator"
)

// Form creates a custom form struct, it embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// New initializes a form struct
func New(data url.Values) *Form {
	return &Form{data, errors(map[string][]string{})}
}

// Has checks if form fields is in post and not empty
func (f *Form) Has(field string, _ *http.Request) bool {
	x := f.Get(field)
	if x == "" {
		f.Errors.Add(field, "This field is mandatory")
		fmt.Println(field, " is mandatory")
		return false
	}
	return true
}

// Valid returns true if there are no errors, otherwise returns false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// Required checks if field is populated
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field is mandatory")
			fmt.Println(field, ": is mandatory")
		}
	}
}

// MinLength checks for strings minimum length
func (f *Form) MinLength(field string, length int) bool {
	x := f.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		fmt.Println(field, ": minimum length:", length)
		return false
	}
	return true
}

// IsEmail checks if field is email
func (f *Form) IsEmail() {
	x := f.Get("email")
	if !valid.IsEmail(x) {
		f.Errors.Add("email", "Invalid email")
		fmt.Println("Invalid email")
	}
}
