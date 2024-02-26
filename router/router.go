package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"ignaciofp.es/web-service-portfolio/config"
)

// Init initializes a gin router with routes and controllers
func Init(init *config.Initialization) *gin.Engine {
	router := gin.Default()

	// Default config:
	// - No origin allowed by default
	// - GET, POST, PUT, HEAD methods
	// - Credentials share disabled
	// - Preflight requests cached for 12 hours
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = append(config.AllowMethods, "OPTIONS")

	router.Use(cors.New(config))

	router.GET("/ping", init.UserCtrl.Ping)
	router.GET("/users/:username", init.UserCtrl.GetUser)
	router.POST("/users", init.UserCtrl.CreateUser)
	router.PUT("/users/:username", init.UserCtrl.UpdateUser)
	router.DELETE("/users/:username", init.UserCtrl.DeleteUser)

	router.Group("/auth")
	{
		router.POST("/login", init.UserCtrl.Authenticate)
		// TODO: Sign-up
	}

	return router
}
