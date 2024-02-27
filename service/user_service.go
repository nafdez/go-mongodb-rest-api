package service

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"ignaciofp.es/web-service-portfolio/model"
	"ignaciofp.es/web-service-portfolio/model/request"
	"ignaciofp.es/web-service-portfolio/repository"
	"ignaciofp.es/web-service-portfolio/util"
)

type UserService interface {
	GetUserByToken(ctx context.Context, token string) (model.User, error)
	GetUserWithPass(ctx context.Context, token string) (model.User, error)
	GetUserByFilter(ctx context.Context, filter bson.D) (model.User, error)
	CreateUser(ctx context.Context, user model.User) error
	UpdateUser(ctx context.Context, token string, updateReq request.Update) error
	DeleteUser(ctx context.Context, token string) error
}

type UserServiceImpl struct {
	repository repository.UserRepository
}

func UserServiceInit(repository repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{repository: repository}
}

// GetUser finds a user in the database using the specified filter and returns it
func (s UserServiceImpl) GetUserByToken(ctx context.Context, token string) (model.User, error) {
	user, err := s.repository.GetUser(ctx, bson.D{{Key: "token", Value: token}})
	user.Password = "" // Easy way to remove password
	return user, err
}

func (s UserServiceImpl) GetUserWithPass(ctx context.Context, token string) (model.User, error) {
	return s.repository.GetUser(ctx, bson.D{{Key: "token", Value: token}})
}

func (s UserServiceImpl) GetUserByFilter(ctx context.Context, filter bson.D) (model.User, error) {
	return s.repository.GetUser(ctx, filter)
}

// UpdateUser updates a user in the database and returns it
func (s UserServiceImpl) UpdateUser(ctx context.Context, token string, updateReq request.Update) error {
	// Making sure user exists and token is valid before updating anything
	user, err := s.GetUserByToken(ctx, token)
	if err != nil {
		return err
	}

	if user.Token != token {
		return util.ErrNoValidTokenProvided
	}

	return s.repository.UpdateUser(ctx, token, updateReq)
}

// CreateUser creates a user in the database
func (s UserServiceImpl) CreateUser(ctx context.Context, user model.User) error {
	// Username and password are required
	if user.Password == "" || user.Username == "" {
		return util.ErrNoUsernameOrPasswordProvided
	}

	// TODO: making sure username is unique

	return s.repository.CreateUser(ctx, user)
}

// DeleteUser deletes a user in the database
func (s UserServiceImpl) DeleteUser(ctx context.Context, token string) error {
	return s.repository.DeleteUser(ctx, token)
}
