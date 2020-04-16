package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func HandleCss(rw http.ResponseWriter, req *http.Request) {
	// get the file
	file, err := os.Open("html" + req.RequestURI)
	if file == nil {
		return
	}
	defer file.Close()
	if err != nil {
		// TODO solve missing fonts
		// fmt.Println("Error while opening css: ", err.Error())
		rw.WriteHeader(404)
		return
	}

	// https://mrwaggel.be/post/golang-transmit-files-over-a-nethttp-server-to-clients/

	// get file size
	fileStat, _ := file.Stat()
	fileSize := strconv.FormatInt(fileStat.Size(), 10)

	// send the headers
	rw.Header().Set("Content-Type", "text/css")
	rw.Header().Set("Content-Length", fileSize)

	// send the file
	_, err = file.Seek(0, 0)
	if err != nil {
		fmt.Println("Error while sending css file to client: ", err.Error())
	}

	_, err = io.Copy(rw, file)
	if err != nil {
		fmt.Println("Error while sending css file to client: ", err.Error())
	}
}
