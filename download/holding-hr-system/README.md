# Holding HR System

Vahid HR idarəetmə platforması - Holding və alt şirkətlər üçün.

## 🚀 Xüsusiyyətlər

- **Multi-tenancy**: Hər alt şirkət yalnız öz məlumatlarını görür
- **3 Əsas Qrup**: Namizədlər, Cari İşçilər, İşdən Çıxanlar
- **İşçi Kartoçkası**: Şəxsi məlumatlar, təhsil, iş təcrübəsi, ailə məlumatları
- **Status Keçidləri**: Namizəd → Cari → Arxiv
- **Responsive Dizayn**: Mobil və desktop uyğunluğu

## 🛠️ Texnologiyalar

- **Backend**: Golang (net/http, html/template)
- **Database**: MariaDB 10.11
- **Frontend**: Tailwind CSS (CDN)
- **Container**: Docker & Docker Compose

## 📦 Quraşdırma

### Docker ilə (Tövsiyə olunur)

```bash
# Proyekti klonla
git clone <repository-url>
cd holding-hr-system

# Docker compose ilə başlat
docker-compose up -d

# Brauzerdə aç
http://localhost:8080
```

### Əl ilə

```bash
# Go modulları yüklə
go mod tidy

# MariaDB-ni quraşdır və migrate et
mysql -u root -p < migrations/01_schema.sql
mysql -u root -p < migrations/02_seed.sql

# Environment dəyişənlərini təyin et
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=hruser
export DB_PASSWORD=hrpass123
export DB_NAME=holding_hr
export JWT_SECRET=your-secret-key

# Serveri başlat
go run cmd/main.go
```

## 👤 Demo Hesabları

| Email | Şifrə | Rol |
|-------|-------|-----|
| admin@abcholding.az | admin123 | Sistem Admin |
| holding.hr@abcholding.az | admin123 | Holding HR |
| hr@abctekstil.az | admin123 | Alt Şirkət HR |
| hr@abclogistika.az | admin123 | Alt Şirkət HR |

## 🌐 Linode Serverinə Daşınma

### 1. Server Hazırlığı

```bash
# Serverə qoşul
ssh root@LINODE_IP

# Docker quraşdır
curl -fsSL https://get.docker.com | sh
systemctl start docker
systemctl enable docker

# Docker Compose quraşdır
apt install docker-compose-plugin
```

### 2. Proyekti Daşı

```bash
# GitHub-dan klonla (tövsiyə olunur)
git clone https://github.com/USERNAME/holding-hr-system.git
cd holding-hr-system

# Və ya SCP ilə köçür
scp -r ./holding-hr-system root@LINODE_IP:/opt/
```

### 3. Başlat

```bash
# Production üçün environment təyin et
export JWT_SECRET="production-secret-key-change-this"

# Docker compose ilə başlat
docker compose up -d

# Logları izlə
docker compose logs -f
```

### 4. Nginx Reverse Proxy (Opsional)

```nginx
server {
    listen 80;
    server_name hr.yourcompany.az;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 📁 Proyekt Strukturu

```
holding-hr-system/
├── cmd/
│   └── main.go              # Entry point
├── config/
│   └── config.go            # Konfiqurasiya
├── internal/
│   ├── handler/             # HTTP handler-lər
│   ├── middleware/          # Auth middleware
│   ├── models/              # Data modelləri
│   └── repository/          # Database əməliyyatları
├── migrations/              # SQL migration faylları
├── static/                  # Static fayllar (CSS, JS, uploads)
├── templates/               # HTML template-lər
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── README.md
```

## 🔒 Təhlükəsizlik

- JWT ilə authentication
- Bcrypt ilə şifrə hash
- Row-Level Security (company_id filter)
- CSRF qorunması

## 📝 Lisenziya

MIT License
