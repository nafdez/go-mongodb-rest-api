package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"ignaciofp.es/web-service-portfolio/model/request"
	"ignaciofp.es/web-service-portfolio/service"
	"ignaciofp.es/web-service-portfolio/util"
)

type UserController interface {
	Ping(ctx *gin.Context)
	GetUser(ctx *gin.Context)
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

// GetUser gets a token and returns the user associated with the token
func (s UserControllerImpl) GetUser(ctx *gin.Context) {
	token := ctx.GetHeader("Token")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing token"})
		return
	}

	user, err := s.service.GetUserByToken(ctx, token)
	if err != nil {
		if errors.Is(err, util.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, user)
}

// UpdateUser takes a token and updates the user who owns the token
// The only fields that are allowed to update are "points" and "token"
func (s UserControllerImpl) UpdateUser(ctx *gin.Context) {
	token := ctx.GetHeader("Token")

	// Checking if provided body is valid and binding to update request
	var updateReq request.Update
	if err := ctx.ShouldBindJSON(&updateReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err := s.service.UpdateUser(ctx, token, updateReq)
	if err != nil {
		if errors.Is(err, util.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, util.ErrNoValidTokenProvided) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

// DeleteUser deletes the user who the token belongs to
func (s UserControllerImpl) DeleteUser(ctx *gin.Context) {
	token := ctx.GetHeader("Token")
	username := ctx.Param("username")

	if username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid argument username"})
		return
	}
	if err := s.service.DeleteUser(ctx, token); err != nil {
		if errors.Is(err, util.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, util.ErrNoValidTokenProvided) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{})
}
