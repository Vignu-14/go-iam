package middleware

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)


type contextKey string

const AppIDContextKey contextKey = "app_id"

func RequireAppID(db *pgxpool.Pool) func(http.Handler) http.Handler{
	return func(next http.Handler)	http.Handler{
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			
		})
	}
}