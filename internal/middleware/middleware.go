package middleware

import (
	"goproject/internal/domain"
	"goproject/internal/jwt"
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type MiddlewareManager struct {
	log  *slog.Logger
	repo domain.SessionRepository
}

func NewMiddlewareManager(log *slog.Logger, sessionRepo domain.SessionRepository) *MiddlewareManager {
	return &MiddlewareManager{
		log:  log,
		repo: sessionRepo,
	}
}

var sessionTimeout = 15 * time.Minute


func (m *MiddlewareManager) JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.log.Error("Authorization header is missing")
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			m.log.Error("Invalid authorization format")
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		claims, err := jwt.ValidateToken(tokenString)
		if err != nil {
			m.log.Error("Invalid token", slog.String("error", err.Error()))
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", claims.UserID)
		ctx = context.WithValue(ctx, "email", claims.Email)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}


func (m *MiddlewareManager) SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value("userID").(string)
		if !ok || userID == "" {
			m.log.Error("User not authenticated, missing or invalid userID in context")
			http.Error(w, "User not authenticated", http.StatusUnauthorized)
			return
		}

		session, err := m.repo.GetSessionByUserID(r.Context(), userID)
		if err != nil {
			m.log.Error("Failed to get session from repository",
				slog.String("userID", userID),
				slog.String("error", err.Error()),
			)
			http.Error(w, "Failed to get session", http.StatusUnauthorized)
			return
		}

		if session == nil || time.Now().UTC().After(session.ExpiresAt) {
			m.log.Error("Session expired", slog.String("userID", userID))
			http.Error(w, "Session expired, please log in again", http.StatusUnauthorized)
			return
		}

		err = m.repo.UpdateSessionExpiry(r.Context(), session.SessionID)
		if err != nil {
			m.log.Error("Failed to update session expiry",
				slog.String("sessionID", session.SessionID),
				slog.String("error", err.Error()),
			)
			http.Error(w, "Failed to update session", http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	})
}