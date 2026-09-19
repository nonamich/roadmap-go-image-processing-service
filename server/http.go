package server

import (
	"io"
	"net/http"
	"os"
)

func HttpServe() {
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		file, header, _ := r.FormFile("test")
		fileBytes, _ := io.ReadAll(file)

		os.WriteFile(header.Filename, fileBytes, 0644)
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {

	})

	http.ListenAndServe(":8080", nil)
}
