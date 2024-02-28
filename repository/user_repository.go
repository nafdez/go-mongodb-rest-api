package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"ignaciofp.es/web-service-portfolio/model"
	"ignaciofp.es/web-service-portfolio/model/request"
	"ignaciofp.es/web-service-portfolio/util"
)

type UserRepository interface {
	GetUser(ctx context.Context, filter bson.D) (model.User, error)
	CreateUser(ctx context.Context, user model.User) error
	UpdateUser(ctx context.Context, token string, updateReq request.Update) error
	DeleteUser(ctx context.Context, token string) error
}

type UserRepositoryImpl struct {
	db             *mongo.Database
	userCollection *mongo.Collection
}

func UserRepositoryInit(db *mongo.Database) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db, userCollection: db.Collection("users")}
}

// GetUser finds a user in the database using the specified filter and returns it
func (r UserRepositoryImpl) GetUser(ctx context.Context, filter bson.D) (model.User, error) {
	var result model.User
	if err := r.userCollection.FindOne(ctx, filter).Decode(&result); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.User{}, util.ErrUserNotFound
		}
		return model.User{}, err
	}
	return result, nil
}

// CreateUser creates a new user in the database
func (r UserRepositoryImpl) CreateUser(ctx context.Context, user model.User) error {
	_, err := r.userCollection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return util.ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

// UpdateUser updates a user in the database and returns it. Only accepts updates to points or token
func (r UserRepositoryImpl) UpdateUser(ctx context.Context, token string, updateReq request.Update) error {
	// Ugly temporary hack
	// When using struct directly it updates
	// both the token and the points
	// Doing this workaround avoids that
	// but not the best.
	// Error: when specifically setting points to 0
	// doesn't update them
	in := bson.M{}
	if updateReq.Points != 0 {
		in["points"] = updateReq.Points
	}
	if updateReq.Token != "" {
		in["token"] = updateReq.Token
	}

	result, err := r.userCollection.UpdateOne(ctx, bson.M{"token": token}, bson.M{"$set": in})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return util.ErrUserNotFound
	}

	return nil
}

// DeleteUser deletes a user in the database
func (r UserRepositoryImpl) DeleteUser(ctx context.Context, token string) error {
	result, err := r.userCollection.DeleteOne(ctx, bson.M{"token": token})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return util.ErrUserNotFound
	}
	return nil
}
