package middleware

import (
	"context"
	"go-iam/internal/pkg/security"
	"net/http"
	"uuid"

	"github.com/jackc/pgx/v5/pgxpool"
)

type contextKey string

const AppIDContextKey contextKey = "app_id"

func RequireAppID(db *pgxpool.Pool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			appIDStr := r.Header.Get("X-App-ID")
			appID, err := uuid.Parse(appIDStr)
			if err != nil {
				http.Error(w, `{"error":"invalid or missing X-App-ID header"}`, http.StatusBadRequest)
				return
			}

			var isActive bool
			err = db.QueryRow(r.Context(), `SELECT is_active FROM apps WHERE id = $1`, appID).Scan(&isActive)
			if err != nil || !isActive {
				http.Error(w, `{"error":"app not found or inactive"}`, http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), AppIDContextKey, appID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireAppOwner(db *pgxpool.Pool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			appID, err := uuid.Parse(r.Header.Get("X-App-ID"))
			appSecret := r.Header.Get("X-app-Secret")
			if err != nil || appSecret == "" {
				http.Error(w, `{"error":"missing X-App-ID or X-App-Secret"}`, http.StatusUnauthorized)
				return
			}
			var secretHash string
			var isActive bool

			err = db.QueryRow(r.Context(),
				`SELECT app_secret_hash, is_active FROM apps WHERE id = $1`, appID,
			).Scan(&secretHash, &isActive)

			if err != nil || !isActive || !security.VerifyBcrypt(appSecret, secretHash) {
				http.Error(w, `{"error":"invalid app credentials"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), AppIDContextKey, appID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
