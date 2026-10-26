package handler

import "github.com/jackc/pgx/v5/pgxpool"

type AuthHandler struct {
	DB *pgxpool.Pool
}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Role         string `json:"role"`
	ExpiresIn    int    `json:"expires_in"`
}
