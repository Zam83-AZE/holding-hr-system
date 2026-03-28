package middleware

import (
	"context"
	"holding-hr-system/internal/models"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserKey contextKey = "user"

func AuthMiddleware(jwtSecret string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Cookie-dən token al
		cookie, err := r.Cookie("token")
		if err != nil {
			if r.Header.Get("HX-Request") == "true" {
				w.Header().Set("HX-Redirect", "/login")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		tokenString := cookie.Value

		// Token parse et
		claims := &models.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			if r.Header.Get("HX-Request") == "true" {
				w.Header().Set("HX-Redirect", "/login")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Context-ə user məlumatını əlavə et
		ctx := context.WithValue(r.Context(), UserKey, claims)
		next(w, r.WithContext(ctx))
	}
}

func GetCurrentUser(r *http.Request) *models.Claims {
	if user, ok := r.Context().Value(UserKey).(*models.Claims); ok {
		return user
	}
	return nil
}

func GenerateToken(user *models.User, jwtSecret string) (string, error) {
	claims := &models.Claims{
		UserID:    user.ID,
		CompanyID: user.CompanyID,
		Email:     user.Email,
		Role:      user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func SetAuthCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Production-da true edilmelidir
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400 * 7, // 7 gun
	})
}

func ClearAuthCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
}

// RequireRole - müəyyən rol tələb et
func RequireRole(roles ...models.Role) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			user := GetCurrentUser(r)
			if user == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			for _, role := range roles {
				if user.Role == role {
					next(w, r)
					return
				}
			}

			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	}
}

// CanAccessCompany - şirkətə çıxış icazəsi
func CanAccessCompany(user *models.Claims, companyID int) bool {
	// Admin və Holding HR bütün şirkətləri görə bilər
	if user.Role == models.RoleAdmin || user.Role == models.RoleHoldingHR {
		return true
	}

	// Alt şirkət HR yalnız öz şirkətini görə bilər
	if user.CompanyID != nil && *user.CompanyID == companyID {
		return true
	}

	return false
}

// GetCompanyFilter - user-in görə biləcəyi şirkət filteri
func GetCompanyFilter(user *models.Claims) *int {
	if user.Role == models.RoleAdmin || user.Role == models.RoleHoldingHR {
		return nil // Bütün şirkətlər
	}
	return user.CompanyID // Yalnız öz şirkəti
}

// IsHTMX - HTMX request-i yoxla
func IsHTMX(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

// GetPartial - partial template render etmək lazımdır mı?
func GetPartial(r *http.Request) bool {
	return strings.HasPrefix(r.URL.Path, "/partial/")
}

// CSRF middleware (sadə implementation)
func CSRFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" || r.Method == "HEAD" || r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}

		// POST, PUT, DELETE üçün CSRF token yoxla
		referer := r.Header.Get("Referer")
		if referer == "" {
			http.Error(w, "CSRF token missing", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Logging middleware
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		// Log request
		// log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
		_ = start
	})
}
