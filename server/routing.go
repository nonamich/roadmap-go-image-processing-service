package server

import (
	"bytes"
	"encoding/json"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"nonamich/image-processing-service/internal/auth"
	"nonamich/image-processing-service/internal/uploader"
	"nonamich/image-processing-service/internal/user"
	"os"
	"strings"
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

	file, header, err := r.FormFile("file")

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	fileBytes, _ := io.ReadAll(file)
	mime := http.DetectContentType(fileBytes)
	fileType := strings.Split(mime, "/")[0]
	ext := strings.Split(mime, "/")[1]

	if fileType != "image" {
		http.Error(w, "File must be image", http.StatusBadRequest)

		return
	}

	config, _, err := image.DecodeConfig(bytes.NewReader(fileBytes))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	metadata := uploader.FileMetadata{
		Width:    config.Width,
		Height:   config.Height,
		Original: header.Filename,
		Size:     uint64(header.Size),
	}
	uploadedFile, err := uploader.SaveFile(user, mime, metadata)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	os.WriteFile("./storage/public/images/"+uploadedFile.Uuid+"."+ext, fileBytes, 0644)

	json.NewEncoder(w).Encode(uploadedFile)
}
