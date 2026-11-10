package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// User merepresentasikan data pengguna
type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username  string             `json:"username" bson:"username" validate:"required"`
	Email     string             `json:"email" bson:"email" validate:"required,email"`
	Password  string             `json:"password,omitempty" bson:"password" validate:"required"`
	FullName  string             `json:"full_name" bson:"full_name"`
	AvatarURL string             `json:"avatar_url" bson:"avatar_url"`
	CreatedAt primitive.DateTime `json:"created_at" bson:"created_at"`
	UpdatedAt primitive.DateTime `json:"updated_at" bson:"updated_at"`
}

// UserResponse adalah response tanpa password
type UserResponse struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Username  string             `json:"username" bson:"username"`
	Email     string             `json:"email" bson:"email"`
	FullName  string             `json:"full_name" bson:"full_name"`
	AvatarURL string             `json:"avatar_url" bson:"avatar_url"`
	CreatedAt primitive.DateTime `json:"created_at" bson:"created_at"`
}

// LoginRequest body untuk login
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// RegisterRequest body untuk registrasi
type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	FullName string `json:"full_name"`
}

// AuthResponse response setelah login/register
type AuthResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
}
