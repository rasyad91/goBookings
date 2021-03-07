package forms

import (
	"net/http"
	"net/url"
)

// Form creates a custom form struct, it embeds a url.Values object
type Form struct {
	data   url.Values
	Errors errors
}

// New initializes a form struct
func New(data url.Values) *Form {
	return &Form{data: data, Errors: errors(map[string][]string{})}
}

// Has checks if form fields is in post and not empty
func (f *Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	if x == "" {
		return false
	}
	return true
}
