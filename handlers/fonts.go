package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

//noinspection SpellCheckingInspection
func HandleFonts(rw http.ResponseWriter, req *http.Request) {

	// get the file
	file, err := os.Open("html" + req.RequestURI)
	defer file.Close()
	if err != nil {
		fmt.Println("Error while opening css: ", err.Error())
		rw.WriteHeader(404)
		return
	}

	// send the content-type header
	// css
	if strings.HasSuffix(req.RequestURI, ".css") {
		rw.Header().Set("Content-Type", "text/css")
		// ttf
	} else if strings.HasSuffix(req.RequestURI, ".ttf") {
		rw.Header().Set("Content-Type", "font/ttf")
		// woff
	} else if strings.HasSuffix(req.RequestURI, ".woff") {
		rw.Header().Set("Content-Type", "font/woff")
		// woff2
	} else if strings.HasSuffix(req.RequestURI, ".woff2") {
		rw.Header().Set("Content-Type", "font/woff2")
	}

	// https://mrwaggel.be/post/golang-transmit-files-over-a-nethttp-server-to-clients/

	// get file size
	fileStat, _ := file.Stat()
	fileSize := strconv.FormatInt(fileStat.Size(), 10)

	// send the content-length header
	rw.Header().Set("Content-Length", fileSize)

	// send the file
	file.Seek(0, 0)
	io.Copy(rw, file)
}
