package model

import "time"

type User struct {
	ID       string    `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string    `json:"username" bson:"username"`
	Password string    `json:"password,omitempty" bson:"password"`
	Name     string    `json:"name,omitempty" bson:"name"`
	Email    string    `json:"email" bson:"email"`
	Role     string    `json:"role" bson:"role"`
	Token    string    `json:"token" bson:"token"`
	LastSeen time.Time `json:"last_seen" bson:"last_seen"`
	Since    time.Time `json:"since" bson:"since"`
	Points   int32     `json:"points" bson:"points"`
}
