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
}

func NewInitialization(
	userRepo repository.UserRepository,
	userSvc service.UserService,
	userCtrl controller.UserController,
) *Initialization {
	return &Initialization{
		userRepo: userRepo,
		userSvc:  userSvc,
		UserCtrl: userCtrl,
	}
}
