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

// GetAllPlaylists mengembalikan semua playlist publik
func GetAllPlaylists(c *fiber.Ctx) error {
	collection := config.GetCollection("playlists")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"is_public": true}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Status:  false,
			Message: "Gagal mengambil data playlist",
		})
	}
	defer cursor.Close(ctx)

	var playlists []model.Playlist
	if err := cursor.All(ctx, &playlists); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Status:  false,
			Message: "Gagal memproses data playlist",
		})
	}

	if playlists == nil {
		playlists = []model.Playlist{}
	}

	return c.JSON(model.Response{
		Status:  true,
		Message: "Daftar playlist",
		Data:    playlists,
	})
}

// GetMyPlaylists mengembalikan playlist milik user yang sedang login
func GetMyPlaylists(c *fiber.Ctx) error {
	userID := c.Locals("userID").(primitive.ObjectID)

	collection := config.GetCollection("playlists")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Status:  false,
			Message: "Gagal mengambil data playlist",
		})
	}
	defer cursor.Close(ctx)

	var playlists []model.Playlist
	if err := cursor.All(ctx, &playlists); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Status:  false,
			Message: "Gagal memproses data playlist",
		})
	}

	if playlists == nil {
		playlists = []model.Playlist{}
	}

	return c.JSON(model.Response{
		Status:  true,
		Message: "Playlist saya",
		Data:    playlists,
	})
}

// GetPlaylistByID mengembalikan detail playlist beserta lagu-lagunya
func GetPlaylistByID(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Status:  false,
			Message: "ID playlist tidak valid",
		})
	}

	collection := config.GetCollection("playlists")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var playlist model.Playlist
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&playlist)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Status:  false,
			Message: "Playlist tidak ditemukan",
		})
	}

	// Ambil detail lagu-lagu di playlist
	songCollection := config.GetCollection("songs")
	var songs []model.Song
	if len(playlist.Songs) > 0 {
		cursor, err := songCollection.Find(ctx, bson.M{
			"_id": bson.M{"$in": playlist.Songs},
		})
		if err == nil {
			cursor.All(ctx, &songs)
			cursor.Close(ctx)
		}
	}

	if songs == nil {
		songs = []model.Song{}
	}

	detail := model.PlaylistDetail{
		ID:          playlist.ID,
		Name:        playlist.Name,
		Description: playlist.Description,
		CoverURL:    playlist.CoverURL,
		UserID:      playlist.UserID,
		Songs:       songs,
		IsPublic:    playlist.IsPublic,
		CreatedAt:   playlist.CreatedAt,
		UpdatedAt:   playlist.UpdatedAt,
	}

	return c.JSON(model.Response{
		Status:  true,
		Message: "Detail playlist",
		Data:    detail,
	})
}

// CreatePlaylist membuat playlist baru
func CreatePlaylist(c *fiber.Ctx) error {
	userID := c.Locals("userID").(primitive.ObjectID)

	var req model.PlaylistRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Status:  false,
			Message: "Format request tidak valid",
		})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Status:  false,
			Message: "Nama playlist wajib diisi",
		})
	}

	now := primitive.NewDateTimeFromTime(time.Now())
	playlist := model.Playlist{
		Name:        req.Name,
		Description: req.Description,
		CoverURL:    req.CoverURL,
		UserID:      userID,
		Songs:       []primitive.ObjectID{},
		IsPublic:    req.IsPublic,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	collection := config.GetCollection("playlists")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, playlist)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Status:  false,
			Message: "Gagal membuat playlist",
		})
	}

	playlist.ID = result.InsertedID.(primitive.ObjectID)

	return c.Status(fiber.StatusCreated).JSON(model.Response{
		Status:  true,
		Message: "Playlist berhasil dibuat",
		Data:    playlist,
	})
}

// UpdatePlaylist memperbarui data playlist
func UpdatePlaylist(c *fiber.Ctx) error {
	userID := c.Locals("userID").(primitive.ObjectID)
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Status:  false,
			Message: "ID playlist tidak valid",
		})
	}

	var req model.PlaylistRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Status:  false,
			Message: "Format request tidak valid",
		})
	}

	collection := config.GetCollection("playlists")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Pastikan playlist milik user yang sedang login
	filter := bson.M{"_id": objectID, "user_id": userID}
	update := bson.M{
		"$set": bson.M{
			"name":        req.Name,
			"description": req.Description,
			"cover_url":   req.CoverURL,
			"is_public":   req.IsPublic,
			"updated_at":  primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Status:  false,
			Message: "Gagal memperbarui playlist",
		})
	}

	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Status:  false,
			Message: "Playlist tidak ditemukan atau bukan milik Anda",
		})
	}

	return c.JSON(model.Response{
		Status:  true,
		Message: "Playlist berhasil diperbarui",
	})
}

// DeletePlaylist menghapus playlist
func DeletePlaylist(c *fiber.Ctx) error {
	userID := c.Locals("userID").(primitive.ObjectID)
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Status:  false,
			Message: "ID playlist tidak valid",
		})
	}

	collection := config.GetCollection("playlists")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Pastikan playlist milik user yang sedang login
	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectID, "user_id": userID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Status:  false,
			Message: "Gagal menghapus playlist",
		})
	}

	if result.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Status:  false,
			Message: "Playlist tidak ditemukan atau bukan milik Anda",
		})
	}

	return c.JSON(model.Response{
		Status:  true,
		Message: "Playlist berhasil dihapus",
	})
}

// AddSongToPlaylist menambahkan lagu ke playlist
func AddSongToPlaylist(c *fiber.Ctx) error {
	userID := c.Locals("userID").(primitive.ObjectID)
	playlistID := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(playlistID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Status:  false,
			Message: "ID playlist tidak valid",
		})
	}

	var req model.AddSongRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Status:  false,
			Message: "Format request tidak valid",
		})
	}

	songID, err := primitive.ObjectIDFromHex(req.SongID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Status:  false,
			Message: "ID lagu tidak valid",
		})
	}

	// Verifikasi lagu ada
	songCollection := config.GetCollection("songs")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := songCollection.CountDocuments(ctx, bson.M{"_id": songID})
	if err != nil || count == 0 {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Status:  false,
			Message: "Lagu tidak ditemukan",
		})
	}

	// Tambahkan lagu ke playlist (hindari duplikat)
	collection := config.GetCollection("playlists")
	filter := bson.M{"_id": objectID, "user_id": userID}
	update := bson.M{
		"$addToSet": bson.M{"songs": songID},
		"$set":      bson.M{"updated_at": primitive.NewDateTimeFromTime(time.Now())},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Status:  false,
			Message: "Gagal menambahkan lagu ke playlist",
		})
	}

	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Status:  false,
			Message: "Playlist tidak ditemukan atau bukan milik Anda",
		})
	}

	return c.JSON(model.Response{
		Status:  true,
		Message: "Lagu berhasil ditambahkan ke playlist",
	})
}

// RemoveSongFromPlaylist menghapus lagu dari playlist
func RemoveSongFromPlaylist(c *fiber.Ctx) error {
	userID := c.Locals("userID").(primitive.ObjectID)
	playlistID := c.Params("id")
	songIDParam := c.Params("songId")

	objectID, err := primitive.ObjectIDFromHex(playlistID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Status:  false,
			Message: "ID playlist tidak valid",
		})
	}

	songID, err := primitive.ObjectIDFromHex(songIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Status:  false,
			Message: "ID lagu tidak valid",
		})
	}

	collection := config.GetCollection("playlists")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": objectID, "user_id": userID}
	update := bson.M{
		"$pull": bson.M{"songs": songID},
		"$set":  bson.M{"updated_at": primitive.NewDateTimeFromTime(time.Now())},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{
			Status:  false,
			Message: "Gagal menghapus lagu dari playlist",
		})
	}

	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Status:  false,
			Message: "Playlist tidak ditemukan atau bukan milik Anda",
		})
	}

	return c.JSON(model.Response{
		Status:  true,
		Message: "Lagu berhasil dihapus dari playlist",
	})
}
