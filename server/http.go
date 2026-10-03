package server

import (
	"log"
	"net/http"
)

func HttpServe() {
	http.HandleFunc("POST /register", registerRoute)
	http.HandleFunc("GET /me", meRoute)
	http.HandleFunc("POST /login", loginRoute)
	http.HandleFunc("POST /images", uploadFileRoute)
	http.HandleFunc("GET /images/{uuid}", getFileRoute)
	http.HandleFunc("GET /images", getFilesRoute)

	log.Print("Server stared http://localhost:8080")

	http.ListenAndServe(":8080", nil)
}
