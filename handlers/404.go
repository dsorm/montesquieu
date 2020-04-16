package handlers

import (
	"fmt"
	"net/http"
)

// Generic 404 page, for use by other handlers in cases of invalid URL
// TODO make nicer 404 page
func Handle404(rw http.ResponseWriter, _ *http.Request) {
	// better than nothing, i guess
	rw.WriteHeader(404)
	_, err := fmt.Fprintf(rw, "Error 404: Not Found\n")

	if err != nil {
		fmt.Println("Error while writing a 404 response:", err.Error())
	}
}
