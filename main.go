package main

import (
	"nonamich/image-processing-service/internal/database"
	"nonamich/image-processing-service/server"

	env "github.com/joho/godotenv"
)

func main() {
	env.Load()
	defer database.InitDatabase().Close()
	server.HttpServe()
}
