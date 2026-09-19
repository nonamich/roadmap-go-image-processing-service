package main

import (
	"nonamich/image-processing-service/server"

	env "github.com/joho/godotenv"
)

func main() {
	env.Load()

	server.HttpServe()
}
