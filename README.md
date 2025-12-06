# Traspac Backend Application

## Deskripsi Proyek
Aplikasi backend untuk kompetisi Traspac yang dibangun menggunakan framework Gin dalam bahasa Go. Aplikasi ini menggunakan raw SQL queries untuk interaksi database dan dikontainerisasi menggunakan Docker. Aplikasi dirancang untuk mengelola fungsionalitas yang berkaitan dengan manajemen user dan fitur-fitur terkait kompetisi.

## Struktur Proyek
```
traspac-backend/
├── cmd/
│   └── server/
│       └── main.go              # Entry point aplikasi
├── internal/
│   ├── config/
│   │   └── config.go           # Konfigurasi dan environment variables
│   ├── controllers/
│   │   └── user_controller.go  # Handler untuk HTTP requests user
│   ├── database/
│   │   ├── connection.go       # Manajemen koneksi database
│   │   └── migrations/
│   │       └── 001_create_users_table.sql # SQL untuk membuat tabel users
│   ├── models/
│   │   └── user.go             # Definisi model User
│   ├── repositories/
│   │   └── user_repository.go  # Interaksi dengan tabel users
│   ├── routes/
│   │   └── routes.go           # Setup routing aplikasi
│   └── services/
│       └── user_service.go     # Business logic untuk users
├── pkg/
│   └── utils/
│       └── response.go         # Utility functions untuk HTTP responses
├── docker/
│   └── Dockerfile              # Instruksi untuk build Docker image
├── scripts/
│   └── migrate.sh              # Script untuk menjalankan database migrations
│   └── init.sql                # Script inisialisasi database
├── .env.example                # Template untuk environment variables
├── docker-compose.yml          # Konfigurasi services Docker
├── go.mod                      # Module dan dependencies Go
├── go.sum                      # Checksums untuk dependencies
└── README.md                   # Dokumentasi proyek
```

## Teknologi yang Digunakan
- **Backend Framework**: Gin (Go)
- **Database**: PostgreSQL 15
- **Cache**: Redis 7
- **Containerization**: Docker & Docker Compose
- **Database Admin**: PgAdmin 4
- **API Documentation**: Swagger (akan tersedia di `/swagger/index.html`)

## Persyaratan Sistem
- Docker Desktop atau Docker Engine
- Docker Compose
- Git
- Port yang tersedia: 8080, 5432, 6379, 5050

## Panduan Setup untuk Juri

### 1. Clone Repository
```bash
git clone <repository-url>
cd traspac-backend
```

### 2. Konfigurasi Environment Variables
```bash
# Salin file .env.example ke .env
cp .env.example .env

# Edit file .env sesuai kebutuhan (opsional, default sudah bisa digunakan)
nano .env
```

### 3. Jalankan Aplikasi dengan Docker Compose
```bash
# Jalankan semua services (Database, Redis, Backend App, PgAdmin)
docker-compose up -d

# Atau tanpa detached mode untuk melihat logs
docker-compose up
```

### 4. Verifikasi Services Berjalan
```bash
# Cek status containers
docker-compose ps

# Cek logs aplikasi
docker-compose logs app

# Cek logs database
docker-compose logs postgres
```

### 5. Akses Aplikasi

| Service | URL | Keterangan |
|---------|-----|------------|
| **Backend API** | `http://localhost:8080` | Main application |
| **API Documentation** | `http://localhost:8080/swagger/index.html` | Swagger UI |
| **Health Check** | `http://localhost:8080/health` | Status aplikasi |
| **PgAdmin** | `http://localhost:5050` | Database management |

#### Kredensial Default:
- **PgAdmin**: 
  - Email: `admin@example.com`
  - Password: `your_pgadmin_password_here`
- **PostgreSQL**:
  - Host: `localhost:5432`
  - Database: `traspac_db`
  - Username: `admin`
  - Password: `your_postgres_password_here`

### 6. Testing API Endpoints
```bash
# Test health endpoint
curl http://localhost:8080/health

# Test dengan tools lain
# - Postman: Import collection dari dokumentasi Swagger
# - Insomnia: Gunakan base URL http://localhost:8080
```

## Panduan Pengembangan

### Menjalankan Migrations
```bash
# Masuk ke container aplikasi
docker-compose exec app ./scripts/migrate.sh

# Atau jalankan manual
docker-compose exec postgres psql -U admin -d traspac_db -f /docker-entrypoint-initdb.d/init.sql
```

### Monitoring Logs
```bash
# Semua services
docker-compose logs -f

# Service tertentu
docker-compose logs -f app
docker-compose logs -f postgres
docker-compose logs -f redis
```

### Menghentikan Aplikasi
```bash
# Stop semua containers
docker-compose down

# Stop dan hapus volumes (data akan hilang)
docker-compose down -v

# Rebuild dan restart
docker-compose down && docker-compose up --build
```

## Troubleshooting

### Port Sudah Digunakan
Jika ada error port sudah digunakan, ubah port di file `.env`:
```env
PORT=8081  # Ubah dari 8080
DB_PORT=5433  # Ubah dari 5432
REDIS_PORT=6380  # Ubah dari 6379
```

### Container Tidak Bisa Connect ke Database
1. Pastikan semua containers running: `docker-compose ps`
2. Cek health check: `docker-compose logs postgres`
3. Restart services: `docker-compose restart`

### Reset Database
```bash
# Hapus volume database dan restart
docker-compose down -v
docker-compose up -d
```

## Fitur Utama
- ✅ RESTful API dengan Gin framework
- ✅ PostgreSQL dengan raw SQL queries
- ✅ Redis untuk caching
- ✅ Docker containerization
- ✅ Health check endpoints
- ✅ Database migrations
- ✅ API documentation dengan Swagger
- ✅ Structured logging
- ✅ Environment-based configuration

## API Documentation
Setelah aplikasi berjalan, akses dokumentasi Swagger di:
`http://localhost:8080/swagger/index.html`

Dokumentasi ini berisi:
- List semua endpoints
- Request/response schemas
- Try-it-out functionality
- Authentication requirements

## Kontak & Support
Untuk pertanyaan teknis atau issue, silakan buat issue di repository ini atau hubungi tim developer.

---
**Catatan untuk Juri**: Aplikasi ini sudah dikonfigurasi untuk berjalan out-of-the-box dengan Docker. Cukup jalankan `docker-compose up` dan semua services akan tersedia otomatis.