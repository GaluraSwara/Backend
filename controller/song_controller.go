package controller

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/GaluraSwara/Backend/config"
	"github.com/GaluraSwara/Backend/model"
)

// GetAllSongs mengembalikan semua lagu
func GetAllSongs(c *fiber.Ctx) error {
	collection := config.GetCollection("songs")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Optional: filter by genre atau artist
	filter := bson.M{}
	if genre := c.Query("genre"); genre != "" {
		filter["genre"] = genre
	}
	if artist := c.Query("artist"); artist != "" {
		filter["artist"] = bson.M{"$regex": artist, "$options": "i"}
	}
	if search := c.Query("search"); search != "" {
		filter["$or"] = []bson.M{
			{"title": bson.M{"$regex": search, "$options": "i"}},
			{"artist": bson.M{"$regex": search, "$options": "i"}},
			{"album": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Status:  false,
			Message: "Gagal mengambil data lagu",
		})
	}
	defer cursor.Close(ctx)

	var songs []model.Song
	if err := cursor.All(ctx, &songs); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Status:  false,
			Message: "Gagal memproses data lagu",
		})
	}

	if songs == nil {
		songs = []model.Song{}
	}

	return c.JSON(model.Response{
		Status:  true,
		Message: "Daftar lagu",
		Data:    songs,
	})
}

// GetSongByID mengembalikan lagu berdasarkan ID
func GetSongByID(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Status:  false,
			Message: "ID lagu tidak valid",
		})
	}

	collection := config.GetCollection("songs")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var song model.Song
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&song)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Status:  false,
			Message: "Lagu tidak ditemukan",
		})
	}

	return c.JSON(model.Response{
		Status:  true,
		Message: "Detail lagu",
		Data:    song,
	})
}

// CreateSong menambahkan lagu baru
func CreateSong(c *fiber.Ctx) error {
	var req model.SongRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Status:  false,
			Message: "Format request tidak valid",
		})
	}

	if req.Title == "" || req.Artist == "" || req.AudioURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Status:  false,
			Message: "Title, artist, dan audio_url wajib diisi",
		})
	}

	now := primitive.NewDateTimeFromTime(time.Now())
	song := model.Song{
		Title:     req.Title,
		Artist:    req.Artist,
		Album:     req.Album,
		Genre:     req.Genre,
		Duration:  req.Duration,
		CoverURL:  req.CoverURL,
		AudioURL:  req.AudioURL,
		CreatedAt: now,
		UpdatedAt: now,
	}

	collection := config.GetCollection("songs")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, song)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Status:  false,
			Message: "Gagal menambahkan lagu",
		})
	}

	song.ID = result.InsertedID.(primitive.ObjectID)

	return c.Status(fiber.StatusCreated).JSON(model.Response{
		Status:  true,
		Message: "Lagu berhasil ditambahkan",
		Data:    song,
	})
}

// UpdateSong memperbarui data lagu
func UpdateSong(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Status:  false,
			Message: "ID lagu tidak valid",
		})
	}

	var req model.SongRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Status:  false,
			Message: "Format request tidak valid",
		})
	}

	update := bson.M{
		"$set": bson.M{
			"title":      req.Title,
			"artist":     req.Artist,
			"album":      req.Album,
			"genre":      req.Genre,
			"duration":   req.Duration,
			"cover_url":  req.CoverURL,
			"audio_url":  req.AudioURL,
			"updated_at": primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	collection := config.GetCollection("songs")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Status:  false,
			Message: "Gagal memperbarui lagu",
		})
	}

	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Status:  false,
			Message: "Lagu tidak ditemukan",
		})
	}

	return c.JSON(model.Response{
		Status:  true,
		Message: "Lagu berhasil diperbarui",
	})
}

// DeleteSong menghapus lagu
func DeleteSong(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Status:  false,
			Message: "ID lagu tidak valid",
		})
	}

	collection := config.GetCollection("songs")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Status:  false,
			Message: "Gagal menghapus lagu",
		})
	}

	if result.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Status:  false,
			Message: "Lagu tidak ditemukan",
		})
	}

	return c.JSON(model.Response{
		Status:  true,
		Message: "Lagu berhasil dihapus",
	})
}
