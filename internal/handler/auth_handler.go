package handler

import "github.com/jackc/pgx/v5/pgxpool"

type AuthHandler struct {
	DB *pgxpool.Pool
}

type AuthRequest struct{
	Email string `json:"email"`
	Password string `json:"password"`
}