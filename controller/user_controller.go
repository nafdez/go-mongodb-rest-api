package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"ignaciofp.es/web-service-portfolio/model"
	"ignaciofp.es/web-service-portfolio/repository"
	"ignaciofp.es/web-service-portfolio/service"
)

type UserController interface {
	Ping(ctx *gin.Context)
	Authenticate(ctx *gin.Context)
	GetUser(ctx *gin.Context)
	CreateUser(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
}

type UserControllerImpl struct {
	service service.UserService
}

func UserControllerInit(service service.UserService) *UserControllerImpl {
	return &UserControllerImpl{service: service}
}

// Ping just sends a "Hello!" string back to the client
func (s UserControllerImpl) Ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Hello!")
}

// Authenticate takes a username and password and checks it against the user that has
// the username provided and returns the user if both match
// Also accepts receiving a token for login
func (s UserControllerImpl) Authenticate(ctx *gin.Context) {
}

// GetUser gets a username and returns the user associated with the username
// TODO: GetUser gets a token from the client and returns the user who has the given token
func (s UserControllerImpl) GetUser(ctx *gin.Context) {
	username := ctx.Param("username")
	if username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid argument username"})
	}

	user, err := s.service.GetUser(ctx, username)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
		ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

// CreateUser creates a new user. It needs the following fields to
// be set: username, email and password.
// It returns a user with a token to use on update and delete.
func (s UserControllerImpl) CreateUser(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	user, err := s.service.CreateUser(ctx, &user)
	if err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			ctx.JSON(http.StatusConflict, err.Error())
			return
		}
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

// UpdateUser takes a token and updates the user who owns the token
// The only fields that are allowed to update are "points" and "token"
func (s UserControllerImpl) UpdateUser(ctx *gin.Context) {
	username := ctx.Param("username")
	token := ctx.GetHeader("Token")

	// Checking username is not empty
	if username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid argument username"})
		return
	}

	// Checking if provided body is valid and binding to model.User struct
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Setting username to user object in case isn't provided on body, so the service
	// can actually find the user to update
	user.Username = username
	user, err := s.service.UpdateUser(ctx, token, &user)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, service.ErrNoValidTokenProvided) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

// DeleteUser deletes the user who the token belongs to
func (s UserControllerImpl) DeleteUser(ctx *gin.Context) {
	token := ctx.GetHeader("Token")
	username := ctx.Param("username")

	if username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid argument username"})
		return
	}
	if err := s.service.DeleteUser(ctx, token, username); err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, service.ErrNoValidTokenProvided) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{})
}
