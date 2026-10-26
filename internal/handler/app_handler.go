package handler

import "github.com/jackc/pgx/v5/pgxpool"

type AppHandler struct {
	DB *pgxpool.Pool
}
