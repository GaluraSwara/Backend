package config

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"
)

// PasetoPayload payload standar token PASETO v4
type PasetoPayload struct {
	Sub       string    `json:"sub"`                  // User ID
	Username  string    `json:"username,omitempty"`  // Username
	Role      string    `json:"role,omitempty"`      // Role user
	IssuedAt  time.Time `json:"iat"`                  // Waktu pembuatan
	ExpiresAt time.Time `json:"exp"`                  // Waktu kadaluarsa
}

var (
	privateKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey
)

func init() {
	initKeys()
}

// Inisialisasi key dari ENV atau generate deterministic pair jika belum ada
func initKeys() {
	privHex := os.Getenv("PASETO_PRIVATE_KEY")
	pubHex := os.Getenv("PASETO_PUBLIC_KEY")

	if privHex != "" && pubHex != "" {
		privBytes, err := hex.DecodeString(privHex)
		if err == nil && len(privBytes) == ed25519.PrivateKeySize {
			privateKey = ed25519.PrivateKey(privBytes)
		}
		pubBytes, err := hex.DecodeString(pubHex)
		if err == nil && len(pubBytes) == ed25519.PublicKeySize {
			publicKey = ed25519.PublicKey(pubBytes)
		}
	}

	// Fallback jika belum di-set di environment: generate key pair
	if privateKey == nil || publicKey == nil {
		pub, priv, err := ed25519.GenerateKey(rand.Reader)
		if err == nil {
			privateKey = priv
			publicKey = pub
		}
	}
}

// GetPublicKeyHex mengembalikan public key dalam format hex string
func GetPublicKeyHex() string {
	return hex.EncodeToString(publicKey)
}

// GetPrivateKeyHex mengembalikan private key dalam format hex string
func GetPrivateKeyHex() string {
	return hex.EncodeToString(privateKey)
}

// LE64 konversi 64-bit uint ke 8-byte little-endian dengan MSB dibersihkan (standar PASETO RFC)
func le64(n uint64) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, n)
	buf[7] &= 0x7F
	return buf
}

// PAE (Pre-Authentication Encoding) resmi standar PASETO v4: PAE(pieces) = LE64(count) || LE64(len1) || piece1 ...
func pae(pieces ...[]byte) []byte {
	count := uint64(len(pieces))
	res := le64(count)
	for _, p := range pieces {
		res = append(res, le64(uint64(len(p)))...)
		res = append(res, p...)
	}
	return res
}

// SignPasetoV4Public menandatangani payload menggunakan PASETO v4.public (Ed25519 Asymmetric)
func SignPasetoV4Public(userID, username string, duration time.Duration) (string, error) {
	if privateKey == nil {
		return "", errors.New("private key belum terinisialisasi")
	}

	now := time.Now()
	payload := PasetoPayload{
		Sub:       userID,
		Username:  username,
		Role:      "user",
		IssuedAt:  now,
		ExpiresAt: now.Add(duration),
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// Format PASETO v4.public: header = "v4.public."
	header := []byte("v4.public.")
	footer := []byte("")
	implicit := []byte("")

	// Pre-Authentication Encoding
	m2 := pae(header, payloadJSON, footer, implicit)

	// Sign menggunakan Ed25519
	sig := ed25519.Sign(privateKey, m2)

	// Body = payloadJSON || sig (64 bytes)
	body := append(payloadJSON, sig...)

	// Base64URL encode tanpa padding
	bodyEncoded := base64.RawURLEncoding.EncodeToString(body)

	return "v4.public." + bodyEncoded, nil
}

// VerifyPasetoV4Public memverifikasi token PASETO v4.public dengan public key
func VerifyPasetoV4Public(tokenStr string) (*PasetoPayload, error) {
	if publicKey == nil {
		return nil, errors.New("public key belum terinisialisasi")
	}

	const prefix = "v4.public."
	if !strings.HasPrefix(tokenStr, prefix) {
		return nil, errors.New("token bukan format paseto v4.public")
	}

	tokenBody := strings.TrimPrefix(tokenStr, prefix)
	// Jika ada footer opsional, split titik
	parts := strings.Split(tokenBody, ".")
	rawBody, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, errors.New("gagal decode token base64url")
	}

	if len(rawBody) < ed25519.SignatureSize {
		return nil, errors.New("ukuran payload token tidak valid")
	}

	// Pisahkan message dan 64-byte signature
	msgLen := len(rawBody) - ed25519.SignatureSize
	payloadJSON := rawBody[:msgLen]
	sig := rawBody[msgLen:]

	header := []byte("v4.public.")
	footer := []byte("")
	implicit := []byte("")

	m2 := pae(header, payloadJSON, footer, implicit)

	// Verifikasi Ed25519
	if !ed25519.Verify(publicKey, m2, sig) {
		return nil, errors.New("tanda tangan paseto v4.public tidak valid atau telah dimodifikasi")
	}

	var payload PasetoPayload
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		return nil, errors.New("format payload token rusak")
	}

	// Validasi expiry
	if time.Now().After(payload.ExpiresAt) {
		return nil, errors.New("token paseto telah kadaluarsa")
	}

	return &payload, nil
}
