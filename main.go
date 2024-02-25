package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"ignaciofp.es/web-service-portfolio/config"
	"ignaciofp.es/web-service-portfolio/router"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	// gin.SetMode(gin.ReleaseMode)
}

func main() {
	// Initializing dependency injection
	init := config.Init()

	// Init gin and it's mappings
	app := router.Init(init)

	port := os.Getenv("PORT")
	app.Run(":" + port)
}
