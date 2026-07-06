# Finvera Backend

REST API untuk aplikasi keuangan pribadi **Finvera**, dibangun dengan **Go** dan **Gin**.

---

## Daftar Isi

- [Tentang](#tentang)
- [Tech Stack](#tech-stack)
- [Arsitektur](#arsitektur)
- [Struktur Folder](#struktur-folder)
- [Prasyarat](#prasyarat)
- [Instalasi](#instalasi)
- [Environment Variable](#environment-variable)
- [Menjalankan Project](#menjalankan-project)
- [Build](#build)
- [Database](#database)
- [API Documentation](#api-documentation)
- [Autentikasi & Keamanan](#autentikasi--keamanan)
- [Deployment](#deployment)
- [Troubleshooting](#troubleshooting)
- [Future Improvement](#future-improvement)

---

## Tentang

Finvera Backend menyediakan REST API untuk manajemen keuangan pribadi. Menangani autentikasi pengguna, manajemen rekening, transaksi, kategori, tag, dan eksekusi transaksi terjadwal otomatis via cron.

---

## Tech Stack

| Teknologi | Versi | Kegunaan |
|---|---|---|
| Go | 1.25 | Bahasa utama |
| Gin | v1.12.0 | HTTP framework |
| GORM | v1.31 | ORM |
| PostgreSQL | 15 | Database utama |
| golang-jwt/jwt | v5.3.1 | JWT authentication |
| golang.org/x/crypto | v0.53.0 | bcrypt password hashing |
| robfig/cron | v3.0.1 | Scheduler transaksi otomatis |
| swaggo/swag | v1.16.6 | Swagger doc generation |
| gin-contrib/cors | v1.7.7 | CORS middleware |

---

## Arsitektur

Project menggunakan **Layered Architecture** dengan pattern **Repository** dan **Service**:

```
HTTP Request
    └── Middleware (Auth JWT, Rate Limit, Security Headers, CORS)
            └── Handler  (validasi input, HTTP request/response)
                    └── Service  (business logic)
                            └── Repository  (data access / query)
                                        └── PostgreSQL
```

Semua Dependency Injection diwiring di satu tempat: `internal/router/router.go`.

---

## Struktur Folder

```
finvera-be/
├── cmd/server/
│   └── main.go                  # Entry point
├── internal/
│   ├── config/                  # Loader konfigurasi dari env
│   ├── cron/                    # Scheduler transaksi terjadwal otomatis
│   ├── database/                # Koneksi PostgreSQL + AutoMigrate
│   ├── dto/                     # Request & response structs
│   ├── handler/                 # HTTP handlers (7 resource)
│   ├── middleware/              # Auth JWT, rate limiter, security headers
│   ├── models/                  # GORM models + composite index
│   ├── repository/              # Data access layer (6 repository)
│   ├── router/                  # Route setup + DI wiring
│   └── service/                 # Business logic (7 service)
├── pkg/utils/                   # JWT & bcrypt helpers
├── docs/                        # Swagger docs (auto-generated via make swagger)
├── db/init.sql                  # SQL inisialisasi database
├── .env.example                 # Template environment variable
├── .air.toml                    # Konfigurasi Air (hot-reload)
├── Dockerfile                   # Multi-stage build
└── Makefile                     # Helper commands
```

---

## Prasyarat

- [Go 1.22+](https://go.dev/dl/)
- [PostgreSQL 15+](https://www.postgresql.org/)
- [Air](https://github.com/air-verse/air) — untuk hot-reload development
- [Docker & Docker Compose](https://docs.docker.com/get-docker/) — opsional

---

## Instalasi

```bash
# 1. Clone repository
git clone https://github.com/Albaihaqi354/Finvera-BE.git
cd Finvera-BE

# 2. Install dependencies
go mod download

# 3. Siapkan environment variable
cp .env.example .env
# Edit .env sesuai konfigurasi lokal
```

---

## Environment Variable

```env
PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASS=your_db_password
DB_NAME=finvera

# Generate dengan: openssl rand -hex 64
JWT_SECRET=replace_with_a_strong_random_secret_at_least_256_bits
JWT_ISSUER=finvera

# development | production
APP_ENV=development

# Comma-separated, tanpa spasi
ALLOWED_ORIGINS=http://localhost:3000
```

| Variable | Wajib | Deskripsi |
|---|---|---|
| `PORT` | Tidak (default: 8080) | Port server |
| `DB_HOST` | Ya | Host PostgreSQL |
| `DB_PORT` | Ya | Port PostgreSQL |
| `DB_USER` | Ya | Username database |
| `DB_PASS` | Ya | Password database |
| `DB_NAME` | Ya | Nama database |
| `JWT_SECRET` | Ya (min 32 char) | Secret key JWT |
| `JWT_ISSUER` | Tidak (default: finvera) | Issuer JWT |
| `APP_ENV` | Tidak | Mode aplikasi |
| `ALLOWED_ORIGINS` | Ya di production | Origin yang diizinkan CORS |

> ⚠️ Jangan pernah commit file `.env` yang berisi secret ke repository.

---

## Menjalankan Project

```bash
# Development dengan hot-reload
make dev

# Tanpa hot-reload
go run ./cmd/server/main.go

# Menggunakan Docker Compose
docker compose up --build
```

---

## Build

```bash
make build      # Build binary ke bin/finvera-be
make run        # Jalankan binary
make swagger    # Regenerate Swagger docs
make tidy       # go mod tidy
```

---

## Database

- Koneksi menggunakan environment variable `DB_*`
- **AutoMigrate** berjalan otomatis saat server start
- Connection pool: max 10 idle, max 100 open, lifetime 1 jam
- Semua tabel user-scoped memiliki **composite index `(user_id, deleted_at)`** untuk performa query optimal

> Untuk production, disarankan menggunakan migration tool seperti [golang-migrate](https://github.com/golang-migrate/migrate).

### Entity

| Model | Deskripsi |
|---|---|
| `User` | Data pengguna |
| `Account` | Rekening (kas, bank, kredit, dll.) |
| `Category` | Kategori transaksi (global + milik user) |
| `Tag` | Label transaksi |
| `Transaction` | Transaksi (income, expense, transfer) |
| `ScheduledTransaction` | Transaksi berulang otomatis |

---

## API Documentation

### Swagger UI

```
http://localhost:8080/swagger/index.html
```

### Base URL

```
http://localhost:8080/api/v1
```

### Endpoint Utama

| Method | Endpoint | Auth | Deskripsi |
|---|---|---|---|
| POST | `/auth/register` | — | Registrasi |
| POST | `/auth/login` | — | Login (username/email) |
| POST | `/auth/logout` | Bearer | Logout |
| GET | `/preset-categories` | — | Kategori preset |
| GET/PUT | `/users/me` | Bearer | Profil pengguna |
| CRUD | `/accounts` | Bearer | Manajemen rekening |
| CRUD | `/categories` | Bearer | Manajemen kategori |
| CRUD | `/tags` | Bearer | Manajemen tag |
| CRUD | `/transactions` | Bearer | Manajemen transaksi |
| CRUD | `/scheduled` | Bearer | Transaksi terjadwal |

### Format Response

```json
// Sukses
{ "success": true, "message": "...", "data": { ... } }

// Error
{ "success": false, "error": "..." }
```

---

## Autentikasi & Keamanan

| Mekanisme | Detail |
|---|---|
| JWT | HMAC-SHA256, masa berlaku 2 jam, alg:none attack prevention |
| Password | bcrypt cost factor 14 |
| Rate Limiting | Auth: 10 req/menit, API: 120 req/menit (per IP, in-memory) |
| CORS | Whitelist via `ALLOWED_ORIGINS` |
| Input Validation | Gin binding (`required`, `min`, `max`, `oneof`) |
| SQL Injection | GORM parameterized queries |
| Security Headers | X-Frame-Options, CSP, HSTS, X-Content-Type-Options, dll. |
| User Enumeration | Login selalu return `"invalid credentials"` untuk user tidak ditemukan maupun password salah |

---

## Deployment

### Docker

```bash
docker build -t finvera-be .
docker compose up --build
```

### CI/CD

Pipeline GitHub Actions berjalan otomatis saat push ke `master`/`main`:
- Build Docker image
- Push ke GitHub Container Registry (GHCR): `ghcr.io/albaihaqi354/finvera-be:master`

### Checklist Production

- [ ] Ganti `JWT_SECRET` dengan output dari `openssl rand -hex 64`
- [ ] Ganti `DB_PASS` dengan password yang kuat
- [ ] Set `APP_ENV=production`
- [ ] Set `ALLOWED_ORIGINS` ke domain yang benar
- [ ] Gunakan nginx sebagai reverse proxy (lihat root `docker-compose.prod.yml`)

---

## Troubleshooting

**"Database configuration is incomplete"**
→ Pastikan semua env var `DB_*` terisi

**"JWT_SECRET must be set" / "JWT_SECRET is too short"**
→ Generate dengan `openssl rand -hex 64`, minimal 32 karakter

**"ALLOWED_ORIGINS must be set in production"**
→ Tambahkan `ALLOWED_ORIGINS=https://yourdomain.com`

**SLOW SQL warnings di log**
→ Restart server untuk memicu AutoMigrate dan membuat composite index

---

## Future Improvement

- [ ] Refresh token (saat ini expire 2 jam)
- [ ] Redis-backed rate limiter untuk multi-instance
- [ ] Database migration tool (golang-migrate / Goose)
- [ ] Unit & integration test
- [ ] Structured logging (zerolog / zap)
- [ ] Forget password / reset password
- [ ] Export transaksi CSV / Excel
