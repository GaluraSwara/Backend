# 🎵 GaluraSwara Backend

Backend API untuk aplikasi pemutar playlist dan musik GaluraSwara, dibangun menggunakan **Gocroot** (Go Fiber + MongoDB).

## Tech Stack

- **Go Fiber** — Web framework
- **MongoDB** — Database
- **JWT** — Autentikasi
- **bcrypt** — Hash password
- **GitHub Actions** — CI/CD ke Alwaysdata

## Struktur Project

```
Backend/
├── .github/workflows/
│   └── deploy.yml          # CI/CD ke Alwaysdata
├── config/
│   ├── db.go               # Koneksi MongoDB
│   └── auth.go             # Konfigurasi JWT
├── controller/
│   ├── auth_controller.go  # Handler auth (login/register)
│   ├── song_controller.go  # Handler CRUD lagu
│   └── playlist_controller.go # Handler CRUD playlist
├── middleware/
│   └── auth_middleware.go  # JWT middleware
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
| Method | Endpoint             | Deskripsi         |
|--------|----------------------|-------------------|
| POST   | `/api/auth/register` | Registrasi user   |
| POST   | `/api/auth/login`    | Login user        |

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

## Setup

### 1. Clone & Install Dependencies
```bash
git clone https://github.com/GaluraSwara/Backend.git
cd Backend
go mod tidy
```

### 2. Konfigurasi Environment
```bash
cp .env.example .env
# Edit .env sesuai konfigurasi MongoDB Anda
```

### 3. Jalankan Server
```bash
go run main.go
```

Server akan berjalan di `http://localhost:8080`

## CI/CD ke Alwaysdata

Tambahkan secrets berikut di **Settings > Secrets and variables > Actions** di GitHub:

| Secret           | Deskripsi                            |
|------------------|--------------------------------------|
| `SSH_HOST`       | SSH hostname Alwaysdata              |
| `SSH_USERNAME`   | Username SSH                         |
| `SSH_PASSWORD`   | Password SSH                         |
| `SSH_PORT`       | Port SSH (default: 22)               |
| `DEPLOY_FOLDER`  | Folder tujuan deploy di server       |
| `API_KEY`        | Alwaysdata API key                   |
| `APP_ID`         | Alwaysdata App/Site ID               |

Setiap push ke branch `main` akan otomatis build dan deploy.
