package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var mh testHandler
	h := NoSurfMiddleware(&mh)

	switch h.(type) {
	case http.Handler: // do nothing means h is type httphandler, therefore pass
	default:
		t.Error(fmt.Sprintf("type is not http.Handler, but is type : %T", h))
	}
}

func TestSessionLoad(t *testing.T) {
	var mh testHandler
	h := SessionLoadMiddleware(&mh)
	switch h.(type) {
	case http.Handler: // do nothing means h is type httphandler, therefore pass
	default:
		t.Error(fmt.Sprintf("type is not http.Handler, but is type : %T", h))
	}
}
