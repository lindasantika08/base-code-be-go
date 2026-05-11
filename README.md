#  Go Base Project

> Kerangka / template proyek Golang yang siap pakai untuk project baru.
> Tinggal clone, sesuaikan, dan langsung coding!

## Struktur Folder
```
go-base-project/
├── cmd/
│   └── api/
│       └── main.go              ← Entry point (titik mulai program)
│
├── config/
│   └── config.go                ← Baca konfigurasi dari .env
│
├── internal/                    ← Kode internal, tidak boleh diimport dari luar
│   ├── database/
│   │   └── mysql.go             ← Koneksi ke mysql
│   │
│   ├── model/
│   │   ├── user.go              ← Struct User, Request & Response DTO
│   │   └── response.go          ← Format standar API response
│   │
│   ├── repository/
│   │   └── user_repository.go   ← Query database (CRUD)
│   │
│   ├── service/
│   │   └── user_service.go      ← Logika bisnis
│   │
│   ├── handler/
│   │   └── user_handler.go      ← HTTP handler (terima request, kirim response)
│   │
│   ├── middleware/
│   │   └── middleware.go        ← Logger, JWT auth, CORS, Recovery
│   │
│   └── utils/
│       └── utils.go             ← Helper: JWT, bcrypt password
│
├── migrations/
│   └── 001_create_users.up.sql  ← File SQL untuk buat tabel
│
├── .env.example                 ← Template konfigurasi (salin jadi .env)
├── .gitignore
├── go.mod
├── Makefile                     ← Perintah-perintah berguna
└── README.md
```

## Cara Memulai

### 1. Clone & Setup
```bash
# Clone repository
git clone <url-repo>
cd go-base-project

# Salin file konfigurasi
cp .env.example .env

# Edit .env sesuai environment kamu
nano .env

### 2. Jalankan MySQL (via Docker)
```bash
make docker-up
# atau manual:
docker run --name mysql-dev \
  -e MYSQL_USER=mysql \
  -e MYSQL_PASSWORD=mysql \
  -e MYSQL_DB=mydb \
  -p 5432:5432 -d mysql:16-alpine
```

### 3. Install Dependencies
```bash
go mod tidy
# atau:
make tidy
```

### 4. Jalankan Migrasi Database
```bash
# Buat tabel di database
make migrate-up
# atau manual:
psql -U mysqls -d mydb -f migrations/001_create_users_table.up.sql
```

### 5. Jalankan Aplikasi
```bash
# Dengan hot-reload (install air dulu: go install github.com/air-verse/air@latest)
make run

# Tanpa hot-reload
make run-simple
# atau:
go run ./cmd/api/main.go
```

### 6. Test Endpoint
```bash
# Health check
curl http://localhost:8080/health

# Register user baru
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Budi","email":"budi@email.com","password":"password123"}'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"budi@email.com","password":"password123"}'

# Get all users (butuh token dari login)
curl http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer <token_dari_login>"
```

```
## Arsitektur: Clean Architecture

```
HTTP Request
    ↓
[Handler]      ← Terima request, validasi input, kirim response
    ↓
[Service]      ← Logika bisnis (business logic)
    ↓
[Repository]   ← Query database
    ↓
[Database]     ← SQL
```

**Aturan dependency:**
- Handler boleh panggil Service
- Service boleh panggil Repository
- Repository langsung akses Database
- Tidak boleh terbalik! (Repository tidak boleh panggil Service)

---

## API Endpoints

| Method | Endpoint                    | Auth | Keterangan           |
|--------|-----------------------------|------|----------------------|
| GET    | `/health`                   |    | Health check         |
| POST   | `/api/v1/auth/register`     |    | Daftar akun baru     |
| POST   | `/api/v1/auth/login`        |    | Login, dapat token   |
| GET    | `/api/v1/users`             |    | Lihat semua user     |
| GET    | `/api/v1/users/:id`         |    | Lihat user by ID     |
| PUT    | `/api/v1/users/:id`         |    | Update user          |
| DELETE | `/api/v1/users/:id`         |    | Hapus user           |

---

## Cara Menambah Fitur Baru

Misalnya mau tambah fitur **Products**:

1. **Buat model** di `internal/model/product.go`
2. **Buat migrasi** di `migrations/002_create_products_table.up.sql`
3. **Buat repository** di `internal/repository/product_repository.go`
4. **Buat service** di `internal/service/product_service.go`
5. **Buat handler** di `internal/handler/product_handler.go`
6. **Daftarkan route** di `cmd/api/main.go` (fungsi `setupRouter`)

---

## Tech Stack

| Komponen    | Library               | Keterangan                    |
|-------------|----------------------|-------------------------------|
| HTTP        | Gin                  | Framework web yang cepat      |
| Database    | SQL + sqlx           | DB relasional + query helper  |
| Auth        | JWT (golang-jwt)     | Autentikasi stateless         |
| Password    | bcrypt               | Hash password yang aman       |
| Config      | godotenv             | Baca file .env                |
| Logging     | Zap (uber-go)        | Structured logging            |

---

## Environment Variables

Lihat file `.env.example` untuk daftar lengkapnya.
