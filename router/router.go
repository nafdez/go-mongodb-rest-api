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
	config.AllowHeaders = append(config.AllowHeaders, "Token")
	config.AllowMethods = append(config.AllowMethods, "OPTIONS")

	router.Use(cors.New(config))

	router.GET("/ping", init.UserCtrl.Ping)

	// Defining groups and int's mappings
	// Routes are duplicated because a weird error where if the
	// route for example is /users/ and client sends a request to
	// /users throws a CORS error.
	var userGroup *gin.RouterGroup = router.Group("/users")
	{
		userGroup.GET("", init.UserCtrl.GetUser)
		userGroup.PUT("", init.UserCtrl.UpdateUser)
		userGroup.DELETE("", init.UserCtrl.DeleteUser)
		userGroup.GET("/", init.UserCtrl.GetUser)
		userGroup.PUT("/", init.UserCtrl.UpdateUser)
		userGroup.DELETE("/", init.UserCtrl.DeleteUser)
	}

	var authGroup *gin.RouterGroup = router.Group("/auth")
	{
		authGroup.POST("/login", init.AuthCtrl.Authenticate)
		authGroup.POST("/register", init.AuthCtrl.Register)
		authGroup.POST("/logout", init.AuthCtrl.Logout)
		authGroup.POST("/login/", init.AuthCtrl.Authenticate)
		authGroup.POST("/register/", init.AuthCtrl.Register)
		authGroup.POST("/logout/", init.AuthCtrl.Logout)
	}

	return router
}
