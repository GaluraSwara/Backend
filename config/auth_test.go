package config

import (
	"testing"
	"time"
)

func TestPasetoV4Public(t *testing.T) {
	userID := "67300a12b9d03820fa49d910"
	username := "asep_sunda"

	// Sign token
	token, err := SignPasetoV4Public(userID, username, 1*time.Hour)
	if err != nil {
		t.Fatalf("Gagal sign PASETO token: %v", err)
	}

	if token == "" {
		t.Fatal("Token kosong")
	}

	if len(token) < 10 || token[:10] != "v4.public." {
		t.Fatalf("Format token bukan v4.public.: %s", token)
	}

	// Verify token
	payload, err := VerifyPasetoV4Public(token)
	if err != nil {
		t.Fatalf("Gagal verifikasi PASETO token: %v", err)
	}

	if payload.Sub != userID {
		t.Errorf("User ID tidak cocok: dapat %s, ingin %s", payload.Sub, userID)
	}

	if payload.Username != username {
		t.Errorf("Username tidak cocok: dapat %s, ingin %s", payload.Username, username)
	}
}
