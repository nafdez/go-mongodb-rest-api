package model

type User struct {
	ID       string `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password,omitempty" bson:"password"`
	Email    string `json:"email" bson:"email"`
	Role     string `json:"role" bson:"role"`
	Token    string `json:"token" bson:"token"`
	Points   int32  `json:"points" bson:"points"`
}
