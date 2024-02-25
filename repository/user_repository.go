package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"ignaciofp.es/web-service-portfolio/model"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user aleady exist")
)

type UserRepository interface {
	GetUser(ctx context.Context, username string) (model.User, error)
	CreateUser(ctx context.Context, user *model.User) (model.User, error)
	UpdateUser(ctx context.Context, user *model.User) (model.User, error)
	DeleteUser(ctx context.Context, username string) error
}

type UserRepositoryImpl struct {
	db             *mongo.Database
	userCollection *mongo.Collection
}

func UserRepositoryInit(db *mongo.Database) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db, userCollection: db.Collection("users")}
}

// GetUser finds a user in the database using the specified filter and returns it
func (r UserRepositoryImpl) GetUser(ctx context.Context, username string) (model.User, error) {
	var result model.User
	if err := r.userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&result); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.User{}, ErrUserNotFound
		}
		return model.User{}, err
	}
	return result, nil
}

// CreateUser creates a new user in the database and returns it
func (r UserRepositoryImpl) CreateUser(ctx context.Context, user *model.User) (model.User, error) {
	result, err := r.userCollection.InsertOne(ctx, user)
	if err != nil {
		return model.User{}, err
	}
	user.ID = result.InsertedID.(primitive.ObjectID).String()
	return *user, nil
}

// UpdateUser updates a user in the database and returns it. Only accepts updates to points and token
func (r UserRepositoryImpl) UpdateUser(ctx context.Context, user *model.User) (model.User, error) {
	in := bson.M{}
	if user.Points != 0 {
		in["points"] = user.Points
	}
	if user.Token != "" {
		in["token"] = user.Token
	}

	result, err := r.userCollection.UpdateOne(ctx, bson.M{"username": user.Username}, bson.M{"$set": in})
	if err != nil {
		return model.User{}, err
	}
	if result.MatchedCount == 0 {
		return model.User{}, ErrUserNotFound
	}

	return *user, nil
}

// DeleteUser deletes a user in the database
func (r UserRepositoryImpl) DeleteUser(ctx context.Context, username string) error {
	result, err := r.userCollection.DeleteOne(ctx, bson.M{"username": username})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return ErrUserNotFound
	}
	return nil
}
