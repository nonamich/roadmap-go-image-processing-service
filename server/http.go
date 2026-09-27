package server

import (
	"log"
	"net/http"
)

func HttpServe() {
	http.HandleFunc("/register", registerRoute)
	http.HandleFunc("/me", meRoute)
	http.HandleFunc("/login", loginRoute)
	http.HandleFunc("/images", uploadRoute)

	log.Print("Server stared http://localhost:8080")

	http.ListenAndServe(":8080", nil)
}
