package forms

import (
	"testing"
)

func TestAdd(t *testing.T) {
	var e errors
	e = make(errors)

	e.Add("field", "message")
	if _, ok := e["field"]; !ok {
		t.Errorf("Test failed, expect true")
		t.Fatal()
	}
	// fmt.Println(e)
	if e["field"][0] != "message" {
		t.Errorf("Test failed, expect \"message\", actual: %s", e["field"][0])
	}
}

func TestGet(t *testing.T) {
	var e errors
	e = make(errors)
	e.Add("field", "message")

	if x := e.Get("field"); x != "message" {
		t.Errorf("Test failed, expect \"message\", actual: %s", x)
	}
}
