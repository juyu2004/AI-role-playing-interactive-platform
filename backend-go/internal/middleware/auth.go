package middleware

import (
	"context"
	"net/http"
	"strings"

	"ai-role-playing-platform/backend-go/internal/config"
	"ai-role-playing-platform/backend-go/internal/security"
)

func JWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") || len(auth) <= 7 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		cfg := config.Load()
		claims, err := security.VerifyHS256(cfg.JWTSecret, token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "user", claims.Sub)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
