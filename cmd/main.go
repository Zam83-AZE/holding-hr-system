package main

import (
        "database/sql"
        "fmt"
        "holding-hr-system/config"
        "holding-hr-system/internal/handler"
        "holding-hr-system/internal/middleware"
        "holding-hr-system/internal/models"
        "holding-hr-system/internal/repository"
        "html/template"
        "log"
        "net/http"
        "os"
        "path/filepath"
        "sort"
        "strings"
        "time"

        _ "github.com/go-sql-driver/mysql"
        "golang.org/x/crypto/bcrypt"
)

func main() {
        // Konfiqurasiyanı yüklə
        cfg := config.Load()

        // Database bağlantısı
        db, err := connectDBWithRetry(cfg)
        if err != nil {
                log.Fatalf("Database bağlantısı alınmadı: %v", err)
        }
        defer db.Close()

        log.Println("Database bağlantısı uğurlu!")

        // Migrasiyaları işə sal
        if err := runMigrations(db); err != nil {
                log.Printf("Migrasiya xətası (davam edilir): %v", err)
        }

        // Repository-ləri yarat
        userRepo := repository.NewUserRepository(db)
        companyRepo := repository.NewCompanyRepository(db)
        employeeRepo := repository.NewEmployeeRepository(db)
        deptRepo := repository.NewDepartmentRepository(db)
        posRepo := repository.NewPositionRepository(db)

        // Demo istifadəçiləri yarat
        seedUsers(userRepo, companyRepo)

        // Template-ləri yüklə
        templates, err := loadTemplates()
        if err != nil {
                log.Fatalf("Template-lər yüklənə bilmədi: %v", err)
        }

        // Handler-ləri yarat
        authHandler := handler.NewAuthHandler(userRepo, cfg.JWTSecret, templates)
        userHandler := handler.NewUserHandler(userRepo, companyRepo, templates)
        employeeHandler := handler.NewEmployeeHandler(employeeRepo, companyRepo, deptRepo, posRepo, templates)
        companyHandler := handler.NewCompanyHandler(companyRepo, deptRepo, posRepo, templates)

        // Router
        mux := http.NewServeMux()

        // Static fayllar
        fs := http.FileServer(http.Dir("static"))
        mux.Handle("/static/", http.StripPrefix("/static/", fs))

        // Uploads
        uploadsFS := http.FileServer(http.Dir("static/uploads"))
        mux.Handle("/uploads/", http.StripPrefix("/uploads/", uploadsFS))

        // Auth route-ları
        mux.HandleFunc("/login", authHandler.ShowLogin)
        mux.HandleFunc("/auth/login", authHandler.Login)
        mux.HandleFunc("/auth/logout", authHandler.Logout)

        // Əsas route-lar (auth tələb olunur)
        mux.HandleFunc("/", middleware.AuthMiddleware(cfg.JWTSecret, employeeHandler.ShowDashboard))
        mux.HandleFunc("/employees", middleware.AuthMiddleware(cfg.JWTSecret, employeeHandler.ShowEmployees))
        mux.HandleFunc("/employee/card", middleware.AuthMiddleware(cfg.JWTSecret, employeeHandler.ShowEmployeeCard))
        mux.HandleFunc("/employee/new", middleware.AuthMiddleware(cfg.JWTSecret, employeeHandler.ShowNewEmployeeForm))
        mux.HandleFunc("/employee/create", middleware.AuthMiddleware(cfg.JWTSecret, employeeHandler.CreateEmployee))
        mux.HandleFunc("/employee/update", middleware.AuthMiddleware(cfg.JWTSecret, employeeHandler.UpdateEmployee))
        mux.HandleFunc("/employee/hire", middleware.AuthMiddleware(cfg.JWTSecret, employeeHandler.HireEmployee))
        mux.HandleFunc("/employee/terminate", middleware.AuthMiddleware(cfg.JWTSecret, employeeHandler.TerminateEmployee))
        mux.HandleFunc("/employee/reactivate", middleware.AuthMiddleware(cfg.JWTSecret, employeeHandler.ReactivateEmployee))
        mux.HandleFunc("/employee/delete", middleware.AuthMiddleware(cfg.JWTSecret, employeeHandler.DeleteEmployee))
        mux.HandleFunc("/employee/search", middleware.AuthMiddleware(cfg.JWTSecret, employeeHandler.SearchEmployees))

        // Əlavə məlumat route-ları
        mux.HandleFunc("/employee/education/add", middleware.AuthMiddleware(cfg.JWTSecret, employeeHandler.AddEducation))
        mux.HandleFunc("/employee/education/update", middleware.AuthMiddleware(cfg.JWTSecret, employeeHandler.UpdateEducation))
        mux.HandleFunc("/employee/education/delete", middleware.AuthMiddleware(cfg.JWTSecret, employeeHandler.DeleteEducation))
        mux.HandleFunc("/employee/experience/add", middleware.AuthMiddleware(cfg.JWTSecret, employeeHandler.AddExperience))
        mux.HandleFunc("/employee/experience/update", middleware.AuthMiddleware(cfg.JWTSecret, employeeHandler.UpdateExperience))
        mux.HandleFunc("/employee/experience/delete", middleware.AuthMiddleware(cfg.JWTSecret, employeeHandler.DeleteExperience))
        mux.HandleFunc("/employee/family/add", middleware.AuthMiddleware(cfg.JWTSecret, employeeHandler.AddFamily))
        mux.HandleFunc("/employee/family/update", middleware.AuthMiddleware(cfg.JWTSecret, employeeHandler.UpdateFamily))
        mux.HandleFunc("/employee/family/delete", middleware.AuthMiddleware(cfg.JWTSecret, employeeHandler.DeleteFamily))

        // API route-ları
        mux.HandleFunc("/api/departments", middleware.AuthMiddleware(cfg.JWTSecret, employeeHandler.GetDepartmentsByCompany))
        mux.HandleFunc("/api/positions", middleware.AuthMiddleware(cfg.JWTSecret, employeeHandler.GetPositionsByCompany))

        // Sertifikat route-ları
        mux.HandleFunc("/employee/certificate/add", middleware.AuthMiddleware(cfg.JWTSecret, employeeHandler.AddCertificate))
        mux.HandleFunc("/employee/certificate/update", middleware.AuthMiddleware(cfg.JWTSecret, employeeHandler.UpdateCertificate))
        mux.HandleFunc("/employee/certificate/delete", middleware.AuthMiddleware(cfg.JWTSecret, employeeHandler.DeleteCertificate))
        mux.HandleFunc("/api/work-locations", middleware.AuthMiddleware(cfg.JWTSecret, employeeHandler.GetWorkLocations))

        // Struktur route-ları
        mux.HandleFunc("/structure", middleware.AuthMiddleware(cfg.JWTSecret, companyHandler.ShowStructure))
        mux.HandleFunc("/structure/company", middleware.AuthMiddleware(cfg.JWTSecret, companyHandler.ShowCompanyStructure))
        mux.HandleFunc("/department/create", middleware.AuthMiddleware(cfg.JWTSecret, companyHandler.CreateDepartment))
        mux.HandleFunc("/position/create", middleware.AuthMiddleware(cfg.JWTSecret, companyHandler.CreatePosition))
        mux.HandleFunc("/department/delete", middleware.AuthMiddleware(cfg.JWTSecret, companyHandler.DeleteDepartment))
        mux.HandleFunc("/position/delete", middleware.AuthMiddleware(cfg.JWTSecret, companyHandler.DeletePosition))

        // Settings route-ları
        mux.HandleFunc("/settings", middleware.AuthMiddleware(cfg.JWTSecret, companyHandler.ShowSettings))
        mux.HandleFunc("/settings/users", middleware.AuthMiddleware(cfg.JWTSecret, userHandler.ShowUsers))
        mux.HandleFunc("/settings/company/create", middleware.AuthMiddleware(cfg.JWTSecret, companyHandler.CreateCompany))
        mux.HandleFunc("/settings/company/delete", middleware.AuthMiddleware(cfg.JWTSecret, companyHandler.DeleteCompany))
        mux.HandleFunc("/settings/user/create", middleware.AuthMiddleware(cfg.JWTSecret, userHandler.CreateUser))

        // Server başlat
        addr := ":" + cfg.ServerPort
        log.Printf("Server %s portunda başladılır...", addr)

        // Production üçün timeout-lu server
        server := &http.Server{
                Addr:         addr,
                Handler:      middleware.LoggingMiddleware(mux),
                ReadTimeout:  15 * time.Second,
                WriteTimeout: 15 * time.Second,
                IdleTimeout:  60 * time.Second,
        }

        if err := server.ListenAndServe(); err != nil {
                log.Fatalf("Server xətası: %v", err)
        }
}

// connectDBWithRetry - Database bağlantısını retry ilə et
func connectDBWithRetry(cfg *config.Config) (*sql.DB, error) {
        var db *sql.DB
        var err error

        maxRetries := 30
        for i := 0; i < maxRetries; i++ {
                db, err = repository.NewDB(cfg)
                if err == nil {
                        return db, nil
                }

                log.Printf("Database bağlantısı cəhdi %d/%d: %v", i+1, maxRetries, err)
                time.Sleep(2 * time.Second)
        }

        return nil, fmt.Errorf("database bağlantısı %d cəhddən sonra alınmadı: %w", maxRetries, err)
}

// loadTemplates - HTML template-ləri yüklə
func loadTemplates() (*template.Template, error) {
        // Template funksiyaları
        funcMap := template.FuncMap{
                "formatDate": func(t *time.Time) string {
                        if t == nil {
                                return ""
                        }
                        return t.Format("02.01.2006")
                },
                "formatDateInput": func(t *time.Time) string {
                        if t == nil {
                                return ""
                        }
                        return t.Format("2006-01-02")
                },
                "statusLabel": func(status string) string {
                        switch status {
                        case "CANDIDATE":
                                return "Namizəd"
                        case "ACTIVE":
                                return "Cari"
                        case "TERMINATED":
                                return "İşdən çıxan"
                        default:
                                return status
                        }
                },
                "statusBadgeClass": func(status string) string {
                        switch status {
                        case "CANDIDATE":
                                return "bg-yellow-100 text-yellow-800"
                        case "ACTIVE":
                                return "bg-green-100 text-green-800"
                        case "TERMINATED":
                                return "bg-red-100 text-red-800"
                        default:
                                return "bg-gray-100 text-gray-800"
                        }
                },
                "roleLabel": func(role string) string {
                        switch role {
                        case "ADMIN":
                                return "Sistem Admin"
                        case "HOLDING_HR":
                                return "Holding HR"
                        case "SUBSIDIARY_HR":
                                return "Alt Şirkət HR"
                        default:
                                return role
                        }
                },
                "relationLabel": func(relation string) string {
                        switch relation {
                        case "FATHER":
                                return "Ata"
                        case "MOTHER":
                                return "Ana"
                        case "SPOUSE":
                                return "Həyat yoldaşı"
                        case "CHILD":
                                return "Övlad"
                        default:
                                return relation
                        }
                },
                "degreeLabel": func(degree string) string {
                        switch degree {
                        case "SECONDARY":
                                return "Orta təhsil"
                        case "VOCATIONAL":
                                return "Peşə təhsili"
                        case "BACHELOR":
                                return "Bakalavr"
                        case "MASTER":
                                return "Magistr"
                        case "PHD":
                                return "PhD"
                        default:
                                return degree
                        }
                },
                "genderLabel": func(gender string) string {
                        switch gender {
                        case "MALE":
                                return "Kişi"
                        case "FEMALE":
                                return "Qadın"
                        default:
                                return ""
                        }
                },
                "lower": func(s string) string {
                        return strings.ToLower(s)
                },
                "add": func(a, b int) int {
                        return a + b
                },
                "now": func() time.Time {
                        return time.Now()
                },
        }

        // Template-ləri yüklə - ParseGlob ilə
        tmpl, err := template.New("").Funcs(funcMap).ParseGlob("templates/*.html")
        if err != nil {
                return nil, err
        }

        return tmpl, nil
}

// seedUsers - Demo istifadəçiləri yarat
func seedUsers(userRepo *repository.UserRepository, companyRepo *repository.CompanyRepository) {
        password := "admin123"
        hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
        if err != nil {
                log.Printf("Şifrə hash xətası: %v", err)
                return
        }

        // Demo istifadəçilər
        users := []struct {
                FullName string
                Email    string
                Role     models.Role
                CompanyID *int
        }{
                {"Sistem Admin", "admin@azmanholding.az", models.RoleAdmin, nil},
                {"Holding HR", "holding.hr@azmanholding.az", models.RoleHoldingHR, intPtr(1)},
                {"Tikinti HR", "hr@azmanconstruction.az", models.RoleSubsidiaryHR, intPtr(2)},
                {"Lojistika HR", "hr@tezlogistics.az", models.RoleSubsidiaryHR, intPtr(3)},
                {"Hotel HR", "hr@sapphirehotels.az", models.RoleSubsidiaryHR, intPtr(4)},
                {"City Service HR", "hr@cityservice.az", models.RoleSubsidiaryHR, intPtr(5)},
                {"EcoProd HR", "hr@ecoprod.az", models.RoleSubsidiaryHR, intPtr(6)},
                {"Mangal HR", "hr@mangalmmc.az", models.RoleSubsidiaryHR, intPtr(7)},
                {"Judo Club HR", "hr@judoclub.az", models.RoleSubsidiaryHR, intPtr(8)},
        }

        for _, u := range users {
                // Mövcud olub-olmadığını yoxla
                existing, _ := userRepo.GetByEmail(u.Email)
                if existing != nil {
                        continue // Artıq mövcuddur
                }

                user := &models.User{
                        FullName:     u.FullName,
                        Email:        u.Email,
                        PasswordHash: string(hash),
                        Role:         u.Role,
                        CompanyID:    u.CompanyID,
                }

                if err := userRepo.Create(user); err != nil {
                        log.Printf("İstifadəçi yaradıla bilmədi (%s): %v", u.Email, err)
                } else {
                        log.Printf("İstifadəçi yaradıldı: %s", u.Email)
                }
        }
}

func intPtr(i int) *int {
        return &i
}

// runMigrations - Migrasiya fayllarını yoxlayır və tətbiq edir
func runMigrations(db *sql.DB) error {
        // Migrasiya tracking cədvəli yarat
        _, err := db.Exec(`
                CREATE TABLE IF NOT EXISTS schema_migrations (
                        version VARCHAR(255) PRIMARY KEY,
                        applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
                ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
        `)
        if err != nil {
                return fmt.Errorf("schema_migrations cədvəli yaradıla bilmədi: %w", err)
        }

        // Migrasiya fayllarını oxu
        migrationsDir := "migrations"
        if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
                log.Println("Migrasiya qovluğu tapılmadı, gözlənilmir")
                return nil
        }

        entries, err := os.ReadDir(migrationsDir)
        if err != nil {
                return fmt.Errorf("migrasiya qovluğu oxunula bilmədi: %w", err)
        }

        // SQL fayllarını sırala
        var sqlFiles []string
        for _, entry := range entries {
                if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
                        sqlFiles = append(sqlFiles, entry.Name())
                }
        }
        sort.Strings(sqlFiles)

        // Hər bir migrasiyanı yoxla və işə sal
        for _, file := range sqlFiles {
                // Artıq tətbiq olunubmu?
                var count int
                db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = ?", file).Scan(&count)
                if count > 0 {
                        continue
                }

                // Faylı oxu
                content, err := os.ReadFile(filepath.Join(migrationsDir, file))
                if err != nil {
                        log.Printf("Migrasiya faylı oxunula bilmədi (%s): %v", file, err)
                        continue
                }

                // Tətbiq et
                log.Printf("Migrasiya tətbiq olunur: %s", file)
                _, err = db.Exec(string(content))
                if err != nil {
                        log.Printf("Migrasiya xətası (%s): %v", file, err)
                        // Davam et - digər migrasiyaları da yoxla
                        continue
                }

                // Qeyd et
                db.Exec("INSERT INTO schema_migrations (version) VALUES (?)", file)
                log.Printf("Migrasiya uğurla tətbiq edildi: %s", file)
        }

        return nil
}
