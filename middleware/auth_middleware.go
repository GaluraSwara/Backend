package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/GaluraSwara/Backend/config"
	"github.com/GaluraSwara/Backend/model"
)

// AuthMiddleware memverifikasi PASETO v4.public token dari header Authorization
func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		// Cek fallback header 'token' (standar Gocroot)
		authHeader = c.Get("token")
	}

	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(model.Response{
			Status:  false,
			Message: "Token otentikasi tidak ditemukan, silakan login",
		})
	}

	// Ambil token dari format "Bearer <token>" atau raw token
	tokenString := authHeader
	if strings.HasPrefix(authHeader, "Bearer ") {
		tokenString = strings.TrimPrefix(authHeader, "Bearer ")
	}

	// Verifikasi PASETO v4.public token menggunakan Ed25519 public key
	payload, err := config.VerifyPasetoV4Public(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(model.Response{
			Status:  false,
			Message: "Token PASETO tidak valid atau kadaluarsa: " + err.Error(),
		})
	}

	userID, err := primitive.ObjectIDFromHex(payload.Sub)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(model.Response{
			Status:  false,
			Message: "Format User ID pada token tidak valid",
		})
	}

	// Simpan identitas user ke konteks Fiber
	c.Locals("userID", userID)
	c.Locals("username", payload.Username)

	return c.Next()
}
