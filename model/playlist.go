package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// Playlist merepresentasikan data playlist
type Playlist struct {
	ID          primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Name        string               `json:"name" bson:"name" validate:"required"`
	Description string               `json:"description" bson:"description"`
	CoverURL    string               `json:"cover_url" bson:"cover_url"`
	UserID      primitive.ObjectID   `json:"user_id" bson:"user_id"`
	Songs       []primitive.ObjectID `json:"songs" bson:"songs"`
	IsPublic    bool                 `json:"is_public" bson:"is_public"`
	CreatedAt   primitive.DateTime   `json:"created_at" bson:"created_at"`
	UpdatedAt   primitive.DateTime   `json:"updated_at" bson:"updated_at"`
}

// PlaylistRequest body untuk membuat/update playlist
type PlaylistRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	CoverURL    string `json:"cover_url"`
	IsPublic    bool   `json:"is_public"`
}

// PlaylistDetail playlist dengan detail lagu-lagunya
type PlaylistDetail struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	CoverURL    string             `json:"cover_url" bson:"cover_url"`
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`
	Songs       []Song             `json:"songs" bson:"songs"`
	IsPublic    bool               `json:"is_public" bson:"is_public"`
	CreatedAt   primitive.DateTime `json:"created_at" bson:"created_at"`
	UpdatedAt   primitive.DateTime `json:"updated_at" bson:"updated_at"`
}

// AddSongRequest body untuk menambah lagu ke playlist
type AddSongRequest struct {
	SongID string `json:"song_id" validate:"required"`
}
