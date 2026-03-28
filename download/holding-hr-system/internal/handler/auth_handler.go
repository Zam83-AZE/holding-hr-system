package handler

import (
        "holding-hr-system/internal/middleware"
        "holding-hr-system/internal/models"
        "holding-hr-system/internal/repository"
        "html/template"
        "net/http"

        "golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
        userRepo   *repository.UserRepository
        jwtSecret  string
        templates  *template.Template
}

func NewAuthHandler(userRepo *repository.UserRepository, jwtSecret string, templates *template.Template) *AuthHandler {
        return &AuthHandler{
                userRepo:  userRepo,
                jwtSecret: jwtSecret,
                templates: templates,
        }
}

// ShowLogin - Login səhifəsini göstər
func (h *AuthHandler) ShowLogin(w http.ResponseWriter, r *http.Request) {
        data := PageData{
                Title: "Giriş",
        }

        // Əgər artıq login olubsa, dashboard-a yönləndir
        if cookie, err := r.Cookie("token"); err == nil && cookie.Value != "" {
                http.Redirect(w, r, "/", http.StatusSeeOther)
                return
        }

        h.templates.ExecuteTemplate(w, "login.html", data)
}

// Login - Login prosesi
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
        email := r.FormValue("email")
        password := r.FormValue("password")

        if email == "" || password == "" {
                h.renderLoginError(w, "Email və şifrə daxil edilməlidir")
                return
        }

        // İstifadəçini tap
        user, err := h.userRepo.GetByEmail(email)
        if err != nil {
                h.renderLoginError(w, "Email və ya şifrə yanlışdır")
                return
        }

        // Şifrəni yoxla
        if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
                h.renderLoginError(w, "Email və ya şifrə yanlışdır")
                return
        }

        // Token yarat
        token, err := middleware.GenerateToken(user, h.jwtSecret)
        if err != nil {
                h.renderLoginError(w, "Sistem xətası")
                return
        }

        // Cookie set et
        middleware.SetAuthCookie(w, token)

        // Redirect
        http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout - Çıxış
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
        middleware.ClearAuthCookie(w)
        http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *AuthHandler) renderLoginError(w http.ResponseWriter, message string) {
        data := PageData{
                Title: "Giriş",
                Error: message,
        }
        h.templates.ExecuteTemplate(w, "login.html", data)
}

// ========== USER MANAGEMENT ==========

type UserHandler struct {
        userRepo    *repository.UserRepository
        companyRepo *repository.CompanyRepository
        templates   *template.Template
}

func NewUserHandler(userRepo *repository.UserRepository, companyRepo *repository.CompanyRepository, templates *template.Template) *UserHandler {
        return &UserHandler{
                userRepo:    userRepo,
                companyRepo: companyRepo,
                templates:   templates,
        }
}

// ShowUsers - İstifadəçilər siyahısı
func (h *UserHandler) ShowUsers(w http.ResponseWriter, r *http.Request) {
        user := middleware.GetCurrentUser(r)

        users, err := h.userRepo.GetAll()
        if err != nil {
                http.Error(w, "Xəta baş verdi", http.StatusInternalServerError)
                return
        }

        companies, _ := h.companyRepo.GetAll()

        data := PageData{
                Title:     "İstifadəçilər",
                User:      user,
                Users:     users,
                Companies: companies,
        }

        h.templates.ExecuteTemplate(w, "users.html", data)
}

// CreateUser - Yeni istifadəçi yarat
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
        currentUser := middleware.GetCurrentUser(r)

        // Yalnız Admin yarada bilər
        if currentUser.Role != models.RoleAdmin {
                http.Error(w, "Səlahiyyətsiz", http.StatusForbidden)
                return
        }

        fullName := r.FormValue("full_name")
        email := r.FormValue("email")
        password := r.FormValue("password")
        role := r.FormValue("role")
        companyIDStr := r.FormValue("company_id")

        // Şifrəni hash et
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
        if err != nil {
                http.Error(w, "Şifrə hash xətası", http.StatusInternalServerError)
                return
        }

        newUser := &models.User{
                FullName:     fullName,
                Email:        email,
                PasswordHash: string(hashedPassword),
                Role:         models.Role(role),
        }

        if companyIDStr != "" && companyIDStr != "0" {
                companyID := parseInt(companyIDStr)
                if companyID > 0 {
                        newUser.CompanyID = &companyID
                }
        }

        if err := h.userRepo.Create(newUser); err != nil {
                http.Error(w, "İstifadəçi yaradıla bilmədi", http.StatusInternalServerError)
                return
        }

        http.Redirect(w, r, "/settings/users", http.StatusSeeOther)
}

func parseInt(s string) int {
        var result int
        for _, c := range s {
                if c >= '0' && c <= '9' {
                        result = result*10 + int(c-'0')
                }
        }
        return result
}
