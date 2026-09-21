package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	Auth "nonamich/image-processing-service/internal/auth"
	UserRepository "nonamich/image-processing-service/internal/user"
	"os"
	"strings"
)

func HttpServe() {
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		file, header, _ := r.FormFile("test")
		fileBytes, _ := io.ReadAll(file)

		os.WriteFile(header.Filename, fileBytes, 0644)
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
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

		existedUser, err := UserRepository.FindUserByUsername(username)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		if existedUser.ID != 0 {
			http.Error(w, "User already existed", http.StatusConflict)

			return
		}

		newUser, err := UserRepository.InsertUser(username, password)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return

		}

		w.Header().Set("Content-Type", "application/json")

		token, err := Auth.JwtEncode(newUser)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		json.NewEncoder(w).Encode(map[string]any{
			"access_token": token,
		})
	})

	http.HandleFunc("/me", func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")

		if authorization == "" {
			http.Error(w, "No authorization", http.StatusUnauthorized)

			return
		}

		tokenString := strings.TrimPrefix(authorization, "Bearer ")

		if tokenString == authorization {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)

			return
		}

		token, err := Auth.JwtDecode(tokenString)

		if err != nil {
			http.Error(w, "Unauthorized.", http.StatusUnauthorized)

			return
		}

		w.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(w).Encode(token.Claims)

		if err != nil {
			http.Error(w, "Failed to encode claims", http.StatusInternalServerError)

			return
		}
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
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

		user, err := UserRepository.FindUserByUsername(username)

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

		token, _ := Auth.JwtEncode(user)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"token": token,
		})

	})

	log.Print("Server stared http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
