package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"ignaciofp.es/web-service-portfolio/config"
	"ignaciofp.es/web-service-portfolio/router"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		// log.Fatalf("Error loading .env file: %s", err)
		log.Println("No .env file found")
	}

	mode := os.Getenv("MODE")
	gin.SetMode(mode)
}

func main() {
	// Initializing dependency injection
	var init *config.Initialization = config.Init()

	// Init gin and it's mappings
	var app *gin.Engine = router.Init(init)

	port := os.Getenv("PORT")
	app.Run(":" + port)
}
