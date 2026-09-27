package server

import (
	"encoding/json"
	"io"
	"net/http"
	"nonamich/image-processing-service/internal/auth"
	"nonamich/image-processing-service/internal/user"
	"os"
)

func registerRoute(w http.ResponseWriter, r *http.Request) {
	var body map[string]any

	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)

		return
	}

	username, ok := body["username"].(string)

	if !ok {
		http.Error(w, "Username is required", http.StatusBadRequest)

		return
	}

	password, ok := body["password"].(string)

	if !ok {
		http.Error(w, "Password is required", http.StatusBadRequest)

		return
	}

	existedUser, err := user.FindUserByUsername(username)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if existedUser.ID != 0 {
		http.Error(w, "User already existed", http.StatusConflict)

		return
	}

	newUser, err := user.SaveUser(username, password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return

	}

	w.Header().Set("Content-Type", "application/json")

	token, err := auth.JwtEncode(newUser)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"access_token": token,
	})
}

func meRoute(w http.ResponseWriter, r *http.Request) {
	user, err := auth.CurrentUser(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(user)
}

func loginRoute(w http.ResponseWriter, r *http.Request) {
	var body map[string]any

	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)

		return
	}

	username, ok := body["username"].(string)

	if !ok {
		http.Error(w, "Username is required", http.StatusBadRequest)

		return
	}

	password, ok := body["password"].(string)

	if !ok {
		http.Error(w, "Password is required", http.StatusBadRequest)

		return
	}

	user, err := user.FindUserByUsername(username)

	if user.ID == 0 {
		http.Error(w, "user not found", http.StatusNotFound)

		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if user.Password != password {
		http.Error(w, "Wrong password", http.StatusForbidden)

		return
	}

	token, _ := auth.JwtEncode(user)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})

}

func uploadRoute(w http.ResponseWriter, r *http.Request) {
	user, err := auth.CurrentUser(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	file, header, _ := r.FormFile("file")
	fileBytes, _ := io.ReadAll(file)

	os.WriteFile(header.Filename, fileBytes, 0644)

	json.NewEncoder(w).Encode(user)
}
