package handlers

import (
	"net/http/httptest"
	"testing"
)

func TestHandle404(t *testing.T) {

	// let's create a new http response recorder which satisfies http.ResponseWriter
	rw := httptest.NewRecorder()

	// http.Request doesn't need to be filled out since 404 handler doesn't really seem to care about it
	Handle404(rw, nil)

	// check if the status code is 404
	if rw.Code != 404 {
		t.Errorf("Handler404 returned non-404 status code: %v\n", rw.Code)
	}
}
