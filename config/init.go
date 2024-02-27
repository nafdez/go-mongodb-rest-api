package config

import (
	"ignaciofp.es/web-service-portfolio/controller"
	"ignaciofp.es/web-service-portfolio/repository"
	"ignaciofp.es/web-service-portfolio/service"
)

type Initialization struct {
	userRepo repository.UserRepository
	userSvc  service.UserService
	UserCtrl controller.UserController
	authSvc  service.AuthService
	AuthCtrl controller.AuthController
}

func NewInitialization(
	userRepo repository.UserRepository,
	userSvc service.UserService,
	userCtrl controller.UserController,
	authSvc service.AuthService,
	authCtrl controller.AuthController,
) *Initialization {
	return &Initialization{
		userRepo: userRepo,
		userSvc:  userSvc,
		UserCtrl: userCtrl,
		authSvc:  authSvc,
		AuthCtrl: authCtrl,
	}
}
