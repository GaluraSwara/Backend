package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// Song merepresentasikan data lagu
type Song struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title     string             `json:"title" bson:"title" validate:"required"`
	Artist    string             `json:"artist" bson:"artist" validate:"required"`
	Album     string             `json:"album" bson:"album"`
	Genre     string             `json:"genre" bson:"genre"`
	Duration  int                `json:"duration" bson:"duration"` // dalam detik
	CoverURL  string             `json:"cover_url" bson:"cover_url"`
	AudioURL  string             `json:"audio_url" bson:"audio_url" validate:"required"`
	CreatedAt primitive.DateTime `json:"created_at" bson:"created_at"`
	UpdatedAt primitive.DateTime `json:"updated_at" bson:"updated_at"`
}

// SongRequest body untuk membuat/update lagu
type SongRequest struct {
	Title    string `json:"title" validate:"required"`
	Artist   string `json:"artist" validate:"required"`
	Album    string `json:"album"`
	Genre    string `json:"genre"`
	Duration int    `json:"duration"`
	CoverURL string `json:"cover_url"`
	AudioURL string `json:"audio_url" validate:"required"`
}
