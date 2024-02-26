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

	// Defining groups and int's mappings
	var userGroup *gin.RouterGroup = router.Group("/users")
	{
		userGroup.GET("/:username", init.UserCtrl.GetUser)
		userGroup.POST("", init.UserCtrl.CreateUser)
		userGroup.PUT("/:username", init.UserCtrl.UpdateUser)
		userGroup.DELETE("/:username", init.UserCtrl.DeleteUser)
	}

	var authGroup *gin.RouterGroup = router.Group("/auth")
	{
		authGroup.POST("/login", init.UserCtrl.Authenticate)
		// TODO: Sign-up
	}

	return router
}
