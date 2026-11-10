# 🎵 GaluraSwara Backend

Backend API untuk aplikasi pemutar playlist dan musik GaluraSwara, dibangun menggunakan arsitektur **Gocroot** (Go Fiber + MongoDB + PASETO v4.public Asymmetric Cryptography).

## Tech Stack

- **Go Fiber v2** — Web framework berkinerja tinggi
- **MongoDB** — Basis data NoSQL dokumen
- **PASETO v4.public (Ed25519)** — Autentikasi asimetris modern (*Platform-Agnostic Security Tokens*)
- **bcrypt** — Hashing password
- **GitHub Actions** — CI/CD otomatis deploy ke Alwaysdata

## Keunggulan PASETO v4.public vs JWT

1. **Bebas Kerentanan Algoritma**: JWT rentan terhadap serangan manipulasi header (`"alg": "none"` atau pergantian RS256 ke HS256). PASETO v4.public secara ketat dan permanen mematok algoritma digital signature ke **EdDSA over Curve25519 (Ed25519)**.
2. **Kriptografi Asimetris**: Private key hanya disimpan di server otentikasi untuk menandatangani token, sedangkan verifikasi token dapat dilakukan secara aman menggunakan Public Key.
3. **Pre-Authentication Encoding (PAE)**: Mencegah serangan pemalsuan panjang payload (*length-extension & canonicalization attacks*).

## Struktur Project

```
Backend/
├── .github/workflows/
│   └── deploy.yml          # CI/CD ke Alwaysdata
├── config/
│   ├── db.go               # Koneksi MongoDB
│   ├── auth.go             # Implementasi PASETO v4.public (Ed25519 PAE)
│   └── auth_test.go        # Unit test otentikasi PASETO
├── controller/
│   ├── auth_controller.go  # Handler auth (login/register)
│   ├── song_controller.go  # Handler CRUD lagu
│   └── playlist_controller.go # Handler CRUD playlist
├── middleware/
│   └── auth_middleware.go  # Middleware verifikasi token PASETO
├── model/
│   ├── user.go             # Model User
│   ├── song.go             # Model Song
│   ├── playlist.go         # Model Playlist
│   └── response.go         # Model Response standar
├── routes/
│   └── routes.go           # Definisi semua routes
├── .env.example            # Template environment variables
├── .gitignore
├── go.mod
├── go.sum
├── main.go                 # Entry point
└── README.md
```

## API Endpoints

### Auth (Public)
| Method | Endpoint             | Deskripsi                                   |
|--------|----------------------|---------------------------------------------|
| POST   | `/api/auth/register` | Registrasi user & mendapatkan token PASETO  |
| POST   | `/api/auth/login`    | Login user & mendapatkan token PASETO       |

### Songs
| Method | Endpoint           | Auth | Deskripsi           |
|--------|--------------------|------|---------------------|
| GET    | `/api/songs`       | ❌   | Daftar semua lagu   |
| GET    | `/api/songs/:id`   | ❌   | Detail lagu         |
| POST   | `/api/songs`       | ✅   | Tambah lagu baru    |
| PUT    | `/api/songs/:id`   | ✅   | Update lagu         |
| DELETE | `/api/songs/:id`   | ✅   | Hapus lagu          |

**Query Parameters** untuk `GET /api/songs`:
- `?search=keyword` — cari berdasarkan judul/artis/album
- `?genre=pop` — filter berdasarkan genre
- `?artist=nama` — filter berdasarkan artis

### Playlists
| Method | Endpoint                          | Auth | Deskripsi                  |
|--------|-----------------------------------|------|----------------------------|
| GET    | `/api/playlists`                  | ❌   | Daftar playlist publik     |
| GET    | `/api/playlists/:id`              | ❌   | Detail playlist + lagu     |
| POST   | `/api/playlists`                  | ✅   | Buat playlist              |
| PUT    | `/api/playlists/:id`              | ✅   | Update playlist            |
| DELETE | `/api/playlists/:id`              | ✅   | Hapus playlist             |
| POST   | `/api/playlists/:id/songs`        | ✅   | Tambah lagu ke playlist    |
| DELETE | `/api/playlists/:id/songs/:songId`| ✅   | Hapus lagu dari playlist   |

### User (Protected)
| Method | Endpoint              | Deskripsi          |
|--------|-----------------------|--------------------|
| GET    | `/api/user/profile`   | Profil user        |
| GET    | `/api/user/playlists` | Playlist milik user|

## Menjalankan Server & Test

```bash
cd Backend

# Jalankan unit test otentikasi PASETO
go test -v ./config

# Jalankan server API
go run main.go
```
