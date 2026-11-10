package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/GaluraSwara/Backend/config"
	"github.com/GaluraSwara/Backend/model"
)

// AuthMiddleware memverifikasi JWT token dari header Authorization
func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(model.Response{
			Status:  false,
			Message: "Token tidak ditemukan, silakan login",
		})
	}

	// Ambil token dari "Bearer <token>"
	tokenString := ""
	if strings.HasPrefix(authHeader, "Bearer ") {
		tokenString = strings.TrimPrefix(authHeader, "Bearer ")
	} else {
		// Support header token langsung tanpa prefix Bearer
		tokenString = authHeader
	}

	// Parse dan validasi token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Metode signing tidak valid")
		}
		return []byte(config.GetJWTSecret()), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(model.Response{
			Status:  false,
			Message: "Token tidak valid atau sudah kadaluarsa",
		})
	}

	// Ambil claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(model.Response{
			Status:  false,
			Message: "Token claims tidak valid",
		})
	}

	// Set user info ke context
	userIDHex, ok := claims["user_id"].(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(model.Response{
			Status:  false,
			Message: "User ID tidak ditemukan di token",
		})
	}

	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(model.Response{
			Status:  false,
			Message: "User ID tidak valid",
		})
	}

	c.Locals("userID", userID)
	c.Locals("username", claims["username"])

	return c.Next()
}
