package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"ignaciofp.es/web-service-portfolio/model/request"
	"ignaciofp.es/web-service-portfolio/service"
	"ignaciofp.es/web-service-portfolio/util"
)

type AuthController interface {
	Authenticate(ctx *gin.Context)
	Register(ctx *gin.Context)
	Logout(ctx *gin.Context)
}

type AuthControllerImpl struct {
	service service.AuthService
}

func AuthControllerInit(service service.AuthService) *AuthControllerImpl {
	return &AuthControllerImpl{service: service}
}

// Authenticate takes a username and password and checks if the password matches
// with the usernames user password and returns a newly generated token
// Also accepts receiving a token for login and still returns a new token
func (s AuthControllerImpl) Authenticate(ctx *gin.Context) {
	token := ctx.GetHeader("Token")

	// Bind json body to loginReq to retrieve username and/or password
	var loginReq request.Auth
	err := ctx.ShouldBindJSON(&loginReq)
	if err != nil && token == "" { // Check if token is provided to not return invalid request when login with token
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// If token is provided then login with that token
	if token != "" {
		loginReq.Token = token
	}

	newToken, err := s.service.Authenticate(ctx, loginReq)
	if err != nil {
		if errors.Is(err, util.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, util.ErrNoUsernameOrPasswordProvided) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, util.ErrInvalidUsernameOrPassword) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"token": newToken})
}

// Register creates a new user. It needs the following fields to
// be set: username, email and password.
// optional: Name and role
// It returns a token used for auth.
func (s AuthControllerImpl) Register(ctx *gin.Context) {
	var registerReq request.Register
	if err := ctx.ShouldBindJSON(&registerReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := s.service.Register(ctx, registerReq)
	if err != nil {
		if errors.Is(err, util.ErrUserAlreadyExists) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (s AuthControllerImpl) Logout(ctx *gin.Context) {
	token := ctx.GetHeader("Token")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing token"})
		return
	}

	err := s.service.Logout(ctx, token)
	if err != nil {
		if errors.Is(err, util.ErrNoValidTokenProvided) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
