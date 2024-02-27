package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"ignaciofp.es/web-service-portfolio/model"
	"ignaciofp.es/web-service-portfolio/model/request"
	"ignaciofp.es/web-service-portfolio/util"
)

type AuthService interface {
	Authenticate(ctx context.Context, authReq request.Auth) (string, error)
	Register(ctx context.Context, registerReq request.Register) (string, error)
	Logout(ctx context.Context, token string) error
	assignNewToken(ctx context.Context, oldToken string) (string, error)
}

type AuthServiceImpl struct {
	service UserService
}

func AuthServiceInit(service UserService) *AuthServiceImpl {
	return &AuthServiceImpl{service: service}
}

// Authenticate checks if username and password are valid and correct and returns newly
// generated token. If a token with a username is provided, and they are valid also returns the new token
// params: username, password, token
func (s AuthServiceImpl) Authenticate(ctx context.Context, authReq request.Auth) (string, error) {
	username := authReq.Username
	password := authReq.Password
	token := authReq.Token

	if token != "" && username != "" {
		// Login with token and username
		return s.authenticateWithToken(ctx, token)
	}

	if username == "" || password == "" {
		return "", util.ErrNoUsernameOrPasswordProvided
	}

	// Login with username and password
	return s.authenticateWithPassword(ctx, username, password)
}

// Register sets all the required data for the user and creates it. then returns a token for auth
func (s AuthServiceImpl) Register(ctx context.Context, registerReq request.Register) (string, error) {
	var user model.User

	user.Username = registerReq.Username
	user.Password = registerReq.Password
	user.Email = registerReq.Email
	user.Name = registerReq.Name
	user.Role = registerReq.Role

	// If no role is specified then assign to normal user
	if registerReq.Role == "" {
		user.Role = "user"
	}

	// Starting points
	user.Points = 500

	user.Token = generateRandomToken()

	// Hashing password
	hashPassword, err := hashPassword(user.Password)
	if err != nil {
		return "", err
	}
	user.Password = hashPassword

	// Setting account creation date and last login time (Now)
	user.Since = time.Now()
	user.LastSeen = time.Now()

	err = s.service.CreateUser(ctx, user)
	if err != nil {
		return "", err
	}

	return user.Token, nil
}

func (s AuthServiceImpl) Logout(ctx context.Context, token string) error {
	s.service.UpdateUser(ctx, token, request.Update{Token: ""})
	return nil
}

// authenticateWithToken authenticates the user with the provided token and username. Returns a new token
func (s AuthServiceImpl) authenticateWithToken(ctx context.Context, token string) (string, error) {
	// When updating the user token it already
	// ensures that the user exist and the token is valid
	// therefore no need to run more checks
	// Generating new token and returning it
	return s.assignNewToken(ctx, token)
}

// authenticateWithToken authenticates the user with the provided username and password. Returns a new token
func (s AuthServiceImpl) authenticateWithPassword(ctx context.Context, username string, password string) (string, error) {
	// user, err := s.service.GetUserWithPass(ctx, username)
	// Quick test

	filter := bson.D{
		{
			Key: "username", Value: username,
		},
	}

	projection := bson.D{
		{
			Key: "last_seen", Value: false,
		},
		{
			Key: "since", Value: false,
		},
	}

	user, err := s.service.GetUserByFilterAndProjection(ctx, filter, projection)
	if err != nil {
		return "", err
	}

	if !checkPasswordHash(password, user.Password) {
		return "", util.ErrInvalidUsernameOrPassword
	}

	// Generating new token and returning it
	return s.assignNewToken(ctx, user.Token)
}

// Creates a new token for the user and updates it, then return the updated token
func (s AuthServiceImpl) assignNewToken(ctx context.Context, oldToken string) (string, error) {
	// Update the user token
	newToken := generateRandomToken()
	// Updates the user in database
	err := s.service.UpdateUser(ctx, oldToken, request.Update{Token: newToken})
	if err != nil {
		return "", err
	}
	// Return the token
	return newToken, nil
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
