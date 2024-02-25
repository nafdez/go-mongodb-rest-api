package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"ignaciofp.es/web-service-portfolio/model"
	"ignaciofp.es/web-service-portfolio/repository"
)

var (
	ErrNoUsernameOrPasswordProvided = errors.New("no username or password provided")
	ErrInvalidUsernameOrPassword    = errors.New("username or password are wrong")
	ErrNoValidTokenProvided         = errors.New("no valid token provided")
)

type UserService interface {
	Authenticate(ctx context.Context, params map[string]string) (model.User, error)
	GetUser(ctx context.Context, username string) (model.User, error)
	CreateUser(ctx context.Context, user *model.User) (model.User, error)
	UpdateUser(ctx context.Context, token string, user *model.User) (model.User, error)
	DeleteUser(ctx context.Context, token string, username string) error
}

type UserServiceImpl struct {
	repository repository.UserRepository
}

func UserServiceInit(repository repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{repository: repository}
}

// Authenticate checks if username and password are valid and correct and returns the user with a newly
// generated token. If a token with a username is provided, and they are valid also returns the user
// params: username, password, token
func (s UserServiceImpl) Authenticate(ctx context.Context, params map[string]string) (model.User, error) {
	username := params["username"]
	password := params["password"]
	token := params["token"]

	if token != "" && username != "" {
		// Login with token and username
		return s.authenticateWithToken(ctx, username, token)
	}

	if username == "" || password == "" {
		return model.User{}, ErrNoUsernameOrPasswordProvided
	}

	// Login with username and password
	return s.authenticateWithPassword(ctx, username, password)
}

// authenticateWithToken authenticates the user with the provided token and username
func (s UserServiceImpl) authenticateWithToken(ctx context.Context, username string, token string) (model.User, error) {
	user, err := s.GetUser(ctx, username)
	if err != nil {
		return model.User{}, err
	}

	// Update the user token
	user.Token = generateRandomToken()
	s.UpdateUser(ctx, token, &user)

	return user, nil
}

// authenticateWithToken authenticates the user with the provided username and password
func (s UserServiceImpl) authenticateWithPassword(ctx context.Context, username string, password string) (model.User, error) {
	user, err := s.getUserInternal(ctx, username)
	if err != nil {
		return model.User{}, err
	}

	if !checkPasswordHash(password, user.Password) {
		return model.User{}, ErrInvalidUsernameOrPassword
	}

	// Update the user token
	user.Token = generateRandomToken()
	s.UpdateUser(ctx, user.Token, &user)

	// Needs to be sanitized for not sending the password
	return sanitizeUser(user), nil
}

// GetUser finds a user in the database using the specified filter and returns it
func (s UserServiceImpl) GetUser(ctx context.Context, username string) (model.User, error) {
	user, err := s.repository.GetUser(ctx, username)
	user = sanitizeUser(user) // Removing password
	return user, err
}

// GetUser finds a user in the database using the specified filter and
// returns it with sensitive information included
func (s UserServiceImpl) getUserInternal(ctx context.Context, username string) (model.User, error) {
	return s.repository.GetUser(ctx, username)
}

// CreateUser creates a new user in the database and returns it
func (s UserServiceImpl) CreateUser(ctx context.Context, user *model.User) (model.User, error) {
	if user.Password == "" || user.Username == "" {
		return model.User{}, ErrNoUsernameOrPasswordProvided
	}

	// Making sure no other user with the same username exists
	_, err := s.repository.GetUser(ctx, user.Username)
	if !errors.Is(err, repository.ErrUserNotFound) {
		return model.User{}, repository.ErrUserAlreadyExists
	}

	if user.Role == "" {
		user.Role = "user"
	}
	user.Points = 500
	user.Token = generateRandomToken()

	// Hashing password
	user.Password, err = hashPassword(user.Password)
	if err != nil {
		return model.User{}, err
	}

	return s.repository.CreateUser(ctx, user)
}

// UpdateUser updates a user in the database and returns it
func (s UserServiceImpl) UpdateUser(ctx context.Context, token string, user *model.User) (model.User, error) {
	// Making sure user exists and token is valid before updating anything
	userDb, err := s.repository.GetUser(ctx, user.Username)
	if err != nil {
		return model.User{}, err
	}

	if userDb.Token != token {
		return model.User{}, ErrNoValidTokenProvided
	}

	return s.repository.UpdateUser(ctx, user)
}

// DeleteUser deletes a user in the database
func (s UserServiceImpl) DeleteUser(ctx context.Context, token string, username string) error {
	// Making sure user exists and token is valid before updating anything
	userDb, err := s.repository.GetUser(ctx, username)
	if err != nil {
		return err
	}

	if userDb.Token != token {
		return ErrNoValidTokenProvided
	}
	return s.repository.DeleteUser(ctx, username)
}

// hashPassword hashes the password provided and returns it
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// checkPasswordHash takes a password and a hashed password and returns if the hashed
// password comes from the password
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// generateRandomToken generates a random token and returns it
func generateRandomToken() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

// sanitizeUser takes a user and returns the same user
// but without password
func sanitizeUser(user model.User) model.User {
	sanitized := user
	sanitized.Password = ""
	return sanitized
}
