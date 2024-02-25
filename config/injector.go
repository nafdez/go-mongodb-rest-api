// go:build wireinject
//go:build wireinject
// +build wireinject

package config

import (
	"github.com/google/wire"
	"ignaciofp.es/web-service-portfolio/controller"
	"ignaciofp.es/web-service-portfolio/repository"
	"ignaciofp.es/web-service-portfolio/service"
)

var db = wire.NewSet(ConnectToDB)

var userServiceSet = wire.NewSet(service.UserServiceInit,
	wire.Bind(new(service.UserService), new(*service.UserServiceImpl)),
)

var userRepoSet = wire.NewSet(repository.UserRepositoryInit,
	wire.Bind(new(repository.UserRepository), new(*repository.UserRepositoryImpl)),
)

var userCtrlSet = wire.NewSet(controller.UserControllerInit,
	wire.Bind(new(controller.UserController), new(*controller.UserControllerImpl)),
)

func Init() *Initialization {
	wire.Build(NewInitialization, db, userCtrlSet, userServiceSet, userRepoSet)
	return nil
}
