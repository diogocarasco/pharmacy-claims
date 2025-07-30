package auth

import (
	"net/http"
	"strings"

	"github.com/diogocarasco/go-pharmacy-service/internal/logger"
)

type Authenticator struct {
	authToken string
	logger    logger.Logger
}

func NewAuthenticator(token string, log logger.Logger) *Authenticator {
	return &Authenticator{
		authToken: token,
		logger:    log,
	}
}

// AuthMiddleware provides authentication for API requests.
func (a *Authenticator) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			a.logger.Warning("Unauthorized access attempt: Authorization header missing.") // Traduzido
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			a.logger.Warning("Unauthorized access attempt: Invalid token format.") // Traduzido
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		token := parts[1]
		if token != a.authToken {
			a.logger.Warning("Unauthorized access attempt: Invalid token.") // Traduzido
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
