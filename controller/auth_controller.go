package controller

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"github.com/GaluraSwara/Backend/config"
	"github.com/GaluraSwara/Backend/model"
)

// Register mendaftarkan user baru
func Register(c *fiber.Ctx) error {
	var req model.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.AuthResponse{
			Status:  false,
			Message: "Format request tidak valid",
		})
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.AuthResponse{
			Status:  false,
			Message: "Username, email, dan password wajib diisi",
		})
	}

	if len(req.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(model.AuthResponse{
			Status:  false,
			Message: "Password minimal 6 karakter",
		})
	}

	collection := config.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Cek apakah username atau email sudah terdaftar
	var existingUser model.User
	err := collection.FindOne(ctx, bson.M{
		"$or": []bson.M{
			{"username": req.Username},
			{"email": req.Email},
		},
	}).Decode(&existingUser)

	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(model.AuthResponse{
			Status:  false,
			Message: "Username atau email sudah terdaftar",
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.AuthResponse{
			Status:  false,
			Message: "Gagal memproses password",
		})
	}

	now := primitive.NewDateTimeFromTime(time.Now())
	user := model.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		FullName:  req.FullName,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.AuthResponse{
			Status:  false,
			Message: "Gagal mendaftarkan user",
		})
	}

	// Generate JWT token
	userID := result.InsertedID.(primitive.ObjectID)
	token, err := generateToken(userID, req.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.AuthResponse{
			Status:  false,
			Message: "Berhasil daftar, tapi gagal membuat token",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.AuthResponse{
		Status:  true,
		Message: "Registrasi berhasil",
		Token:   token,
	})
}

// Login mengautentikasi user dan mengembalikan token
func Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.AuthResponse{
			Status:  false,
			Message: "Format request tidak valid",
		})
	}

	if req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.AuthResponse{
			Status:  false,
			Message: "Username dan password wajib diisi",
		})
	}

	collection := config.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user model.User
	err := collection.FindOne(ctx, bson.M{"username": req.Username}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(model.AuthResponse{
			Status:  false,
			Message: "Username atau password salah",
		})
	}

	// Verifikasi password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(model.AuthResponse{
			Status:  false,
			Message: "Username atau password salah",
		})
	}

	// Generate token
	token, err := generateToken(user.ID, user.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.AuthResponse{
			Status:  false,
			Message: "Gagal membuat token",
		})
	}

	return c.JSON(model.AuthResponse{
		Status:  true,
		Message: "Login berhasil",
		Token:   token,
	})
}

// GetProfile mengembalikan profil user yang sedang login
func GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(primitive.ObjectID)

	collection := config.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user model.User
	err := collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Status:  false,
			Message: "User tidak ditemukan",
		})
	}

	return c.JSON(model.Response{
		Status:  true,
		Message: "Profil user",
		Data: model.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FullName:  user.FullName,
			AvatarURL: user.AvatarURL,
			CreatedAt: user.CreatedAt,
		},
	})
}

// generateToken membuat JWT token
func generateToken(userID primitive.ObjectID, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID.Hex(),
		"username": username,
		"exp":      time.Now().Add(config.GetJWTExpiry()).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GetJWTSecret()))
}
