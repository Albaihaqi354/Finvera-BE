# Finvera Backend

REST API untuk aplikasi keuangan pribadi **Finvera**, dibangun dengan **Go** dan **Gin**.

> 🚀 Production: [finvera-be-production.up.railway.app](https://finvera-be-production.up.railway.app)  
> 🔗 Frontend Repository: [Finvera-App](https://github.com/Albaihaqi354/Finvera-app)

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
- [Database & Migrasi](#database--migrasi)
- [API Documentation](#api-documentation)
- [Autentikasi & Keamanan](#autentikasi--keamanan)
- [Deployment](#deployment)
- [Troubleshooting](#troubleshooting)

---

## Tentang

Finvera Backend menyediakan REST API untuk manajemen keuangan pribadi. Menangani autentikasi pengguna, manajemen rekening, transaksi, kategori, tag, dan eksekusi transaksi terjadwal otomatis via cron job.

---

## Tech Stack

| Teknologi | Versi | Kegunaan |
|---|---|---|
| Go | 1.25 | Bahasa utama |
| Gin | v1.12.0 | HTTP framework |
| GORM | v1.31 | ORM |
| PostgreSQL (Neon) | 15 | Database production |
| golang-jwt/jwt | v5.3.1 | JWT authentication |
| golang.org/x/crypto | v0.53.0 | bcrypt password hashing |
| robfig/cron | v3.0.1 | Scheduler transaksi otomatis |
| golang-migrate | v4.19.1 | Database migration tool |
| swaggo/swag | v1.16.6 | Swagger doc generation |
| gin-contrib/cors | v1.7.7 | CORS middleware |
| google/uuid | v1.6.0 | UUID primary key |

---

## Arsitektur

Project menggunakan **Layered Architecture** dengan pattern **Repository** dan **Service**:

```
HTTP Request
    └── Middleware (Auth JWT, Rate Limit, Security Headers, CORS)
            └── Handler  (validasi input, HTTP request/response)
                    └── Service  (business logic)
                            └── Repository  (data access / query)
                                        └── PostgreSQL (Neon)
```

Semua Dependency Injection diwiring di satu tempat: `internal/router/router.go`.

---

## Struktur Folder

```
finvera-be/
├── cmd/server/
│   └── main.go                  # Entry point + Graceful Shutdown
├── internal/
│   ├── config/                  # Loader konfigurasi dari env
│   ├── cron/                    # Scheduler transaksi terjadwal (deadlock-safe)
│   ├── database/                # Koneksi PostgreSQL + Connection Pool
│   ├── dto/                     # Request & response structs
│   ├── handler/                 # HTTP handlers (7 resource)
│   ├── middleware/              # Auth JWT, rate limiter, security headers
│   ├── models/                  # GORM models + composite index
│   ├── repository/              # Data access layer (6 repository)
│   ├── router/                  # Route setup + DI wiring
│   └── service/                 # Business logic (7 service)
├── pkg/
│   ├── blacklist/               # JWT blacklist (persistent di DB)
│   └── utils/                   # JWT & bcrypt helpers
├── db/
│   ├── migrations/              # Versioned SQL migrations (golang-migrate)
│   │   ├── 000001_init_schema.up.sql
│   │   ├── 000001_init_schema.down.sql
│   │   ├── 000002_add_missing_columns.up.sql
│   │   ├── 000002_add_missing_columns.down.sql
│   │   ├── 000003_add_tag_groups.up.sql
│   │   └── 000003_add_tag_groups.down.sql
│   └── init.sql                 # DDL referensi (legacy)
├── docs/                        # Swagger docs (auto-generated)
├── .env.example                 # Template environment variable
├── .air.toml                    # Konfigurasi Air (hot-reload)
├── Dockerfile                   # Multi-stage build
└── Makefile                     # Helper commands
```

---

## Prasyarat

- [Go 1.22+](https://go.dev/dl/)
- [PostgreSQL 15+](https://www.postgresql.org/) atau akun [Neon](https://neon.tech) untuk production
- [Air](https://github.com/air-verse/air) — untuk hot-reload development (opsional)
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

# Set ke false di production — gunakan golang-migrate
AUTO_MIGRATE=false

# Comma-separated, tanpa spasi, tanpa trailing slash
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
| `JWT_SECRET` | Ya (min 32 char) | Secret key JWT — generate dengan `openssl rand -hex 64` |
| `JWT_ISSUER` | Tidak (default: finvera) | Issuer JWT |
| `APP_ENV` | Tidak (default: development) | Mode aplikasi |
| `AUTO_MIGRATE` | Tidak (default: false) | Jalankan GORM AutoMigrate saat start |
| `ALLOWED_ORIGINS` | Ya di production | Origin yang diizinkan CORS (pisahkan dengan koma) |

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

## Database & Migrasi

Project menggunakan **[golang-migrate](https://github.com/golang-migrate/migrate)** sebagai migration tool untuk versioned schema management.

### Menjalankan Migrasi

```bash
# Install migrate CLI (sekali saja)
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Apply semua migrasi ke database
migrate -path db/migrations -database "postgresql://user:pass@host/dbname?sslmode=require" up

# Rollback 1 step
migrate -path db/migrations -database "postgresql://..." down 1

# Cek versi migrasi saat ini
migrate -path db/migrations -database "postgresql://..." version
```

### Daftar Migrasi

| Versi | Nama | Deskripsi |
|---|---|---|
| 000001 | init_schema | Tabel dasar: users, accounts, categories, tags, transactions, scheduled_transactions + preset categories |
| 000002 | add_missing_columns | Tambah kolom currency, geo_lat, geo_lng, geo_name di transactions; currency, last_run di scheduled_transactions |
| 000003 | add_tag_groups | Buat tabel tag_groups; tambah group_id, sort_order di tags |

### Entity

| Model | Deskripsi |
|---|---|
| `User` | Data pengguna |
| `Account` | Rekening (kas, bank, kredit, dll.) |
| `Category` | Kategori transaksi (global/preset + milik user) |
| `TagGroup` | Grup untuk mengelompokkan tag |
| `Tag` | Label transaksi |
| `Transaction` | Transaksi (income, expense, transfer) dengan geolokasi & currency |
| `ScheduledTransaction` | Transaksi berulang otomatis dengan last_run tracking |
| `JWTBlacklist` | Token JWT yang sudah di-revoke (persistent) |

### Connection Pool

- Max 10 idle connections
- Max 100 open connections
- Lifetime 1 jam
- Semua tabel user-scoped memiliki **composite index `(user_id, deleted_at)`** untuk performa query optimal

---

## API Documentation

### Swagger UI

Hanya tersedia di mode **development**:

```
http://localhost:8080/swagger/index.html
```

> Swagger dinonaktifkan di production untuk alasan keamanan.

### Base URL

```
# Development
http://localhost:8080/api/v1

# Production
https://finvera-be-production.up.railway.app/api/v1
```

### Endpoint Utama

| Method | Endpoint | Auth | Deskripsi |
|---|---|---|---|
| GET | `/ping` | — | Health check (cek koneksi DB) |
| POST | `/auth/register` | — | Registrasi |
| POST | `/auth/login` | — | Login (username/email) |
| POST | `/auth/logout` | Bearer | Logout (revoke token) |
| GET | `/preset-categories` | — | Kategori preset sistem |
| GET/PUT | `/users/me` | Bearer | Profil pengguna |
| PATCH | `/users/me/settings` | Bearer | Update pengaturan |
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
| JWT Blacklist | Persistent di database — logout benar-benar merevoke token meski server restart |
| Password | bcrypt cost factor 14 |
| Rate Limiting | Auth: 10 req/menit, API: 120 req/menit (per IP, in-memory) |
| CORS | Whitelist via `ALLOWED_ORIGINS` env var |
| Input Validation | Gin binding (`required`, `min`, `max`, `oneof`) |
| SQL Injection | GORM parameterized queries |
| Security Headers | X-Frame-Options, CSP, HSTS, X-Content-Type-Options, dll. |
| User Enumeration | Login selalu return `"invalid credentials"` untuk user tidak ditemukan maupun password salah |
| Graceful Shutdown | Handle SIGTERM/SIGINT — request aktif diselesaikan sebelum shutdown |
| Swagger | Dinonaktifkan di production |
| Pagination | Max limit 100 rows per request |

---

## Deployment

### Railway (Production)

Backend di-deploy ke **[Railway](https://railway.app)**. Deploy otomatis berjalan setiap kali push ke branch `master`.

**Environment Variables yang harus diset di Railway:**

| Variable | Contoh Value |
|---|---|
| `DB_HOST` | `ep-xxxxx.ap-southeast-1.aws.neon.tech` |
| `DB_PORT` | `5432` |
| `DB_USER` | `neondb_owner` |
| `DB_PASS` | `your_neon_password` |
| `DB_NAME` | `neondb` |
| `JWT_SECRET` | Output dari `openssl rand -hex 64` |
| `JWT_ISSUER` | `finvera` |
| `APP_ENV` | `production` |
| `AUTO_MIGRATE` | `false` |
| `ALLOWED_ORIGINS` | `https://finvera-app.vercel.app` |

> Gunakan `DB_HOST` dari bagian **PGHOST_UNPOOLED** (tanpa `-pooler`) agar GORM bisa menggunakan connection pool bawaan Go secara optimal.

### Menjalankan Migrasi ke Production

```bash
# Dari local — arahkan ke Neon production database
migrate \
  -path db/migrations \
  -database "postgresql://neondb_owner:password@ep-xxx.neon.tech/neondb?sslmode=require" \
  up
```

### Docker

```bash
docker build -t finvera-be .
docker compose up --build
```

---

## Troubleshooting

**"Database configuration is incomplete"**  
→ Pastikan semua env var `DB_*` terisi

**"JWT_SECRET must be set" / "JWT_SECRET is too short"**  
→ Generate dengan `openssl rand -hex 64`, minimal 32 karakter

**"ALLOWED_ORIGINS must be set in production"**  
→ Tambahkan `ALLOWED_ORIGINS=https://finvera-app.vercel.app` (tanpa trailing slash)

**500 Error — column "xxx" does not exist**  
→ Ada kolom di GORM model yang belum ada di database. Jalankan `migrate up` ke production database.

**CORS error dari frontend**  
→ Pastikan `ALLOWED_ORIGINS` di Railway persis sama dengan origin frontend (cek URL di pesan error browser, tidak ada trailing slash)

**Scheduled transaction tidak berjalan**  
→ Cek Railway logs untuk error di cron service. Pastikan semua foreign key (account_id, category_id) valid.
