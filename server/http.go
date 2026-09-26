package server

import (
	"io"
	"log"
	"net/http"
	"os"
)

func HttpServe() {
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		file, header, _ := r.FormFile("test")
		fileBytes, _ := io.ReadAll(file)

		os.WriteFile(header.Filename, fileBytes, 0644)
	})

	http.HandleFunc("/register", registerRoute)
	http.HandleFunc("/me", meRoute)
	http.HandleFunc("/login", loginRoute)

	log.Print("Server stared http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
