package handler

import (
	"net/http"
	"os"
	"io"
)

func fileUploadHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20) // // 32 MB is the maximum file size
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    defer file.Close()

	// Create a new file in the uploads directory
    f, err := os.OpenFile("./files/" +handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer f.Close()

	// Copy the contents of the file to the new file
    _, err = io.Copy(f, file)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	w.Write([]byte("File uploaded successfully"))
}