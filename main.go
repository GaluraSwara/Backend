package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	"github.com/GaluraSwara/Backend/config"
	"github.com/GaluraSwara/Backend/routes"
)

func main() {
	// Load .env file jika ada (development)
	_ = godotenv.Load()

	// Koneksi ke MongoDB
	config.ConnectDB()

	// Inisialisasi Fiber
	app := fiber.New(fiber.Config{
		ServerHeader: "GaluraSwara",
		AppName:      "GaluraSwara Backend API v1.0.0",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Health check
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  true,
			"message": "Wilujeng Sumping di GaluraSwara API 🎵",
		})
	})

	// Setup Routes
	routes.Setup(app)

	// Port configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🎵 GaluraSwara Backend running on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
