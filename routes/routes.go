package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/GaluraSwara/Backend/controller"
	"github.com/GaluraSwara/Backend/middleware"
)

// Setup mendaftarkan semua route API
func Setup(app *fiber.App) {
	api := app.Group("/api")

	// ===== Auth Routes (Public) =====
	auth := api.Group("/auth")
	auth.Post("/register", controller.Register)
	auth.Post("/login", controller.Login)

	// ===== Song Routes =====
	songs := api.Group("/songs")
	songs.Get("/", controller.GetAllSongs)       // Public: daftar semua lagu
	songs.Get("/:id", controller.GetSongByID)     // Public: detail lagu

	// Song routes yang butuh autentikasi
	songsAuth := songs.Group("", middleware.AuthMiddleware)
	songsAuth.Post("/", controller.CreateSong)    // Tambah lagu baru
	songsAuth.Put("/:id", controller.UpdateSong)  // Update lagu
	songsAuth.Delete("/:id", controller.DeleteSong) // Hapus lagu

	// ===== Playlist Routes =====
	playlists := api.Group("/playlists")
	playlists.Get("/", controller.GetAllPlaylists)       // Public: daftar playlist publik
	playlists.Get("/:id", controller.GetPlaylistByID)    // Public: detail playlist

	// Playlist routes yang butuh autentikasi
	playlistsAuth := playlists.Group("", middleware.AuthMiddleware)
	playlistsAuth.Post("/", controller.CreatePlaylist)              // Buat playlist
	playlistsAuth.Put("/:id", controller.UpdatePlaylist)            // Update playlist
	playlistsAuth.Delete("/:id", controller.DeletePlaylist)         // Hapus playlist
	playlistsAuth.Post("/:id/songs", controller.AddSongToPlaylist)  // Tambah lagu ke playlist
	playlistsAuth.Delete("/:id/songs/:songId", controller.RemoveSongFromPlaylist) // Hapus lagu dari playlist

	// ===== User Routes (Protected) =====
	user := api.Group("/user", middleware.AuthMiddleware)
	user.Get("/profile", controller.GetProfile)          // Profil user
	user.Get("/playlists", controller.GetMyPlaylists)     // Playlist milik user
}
