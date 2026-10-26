package handler

import (
	"uuid"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AppHandler struct {
	DB *pgxpool.Pool
}

type RegisterAppRequest struct{
	Name string `json:"name"`
}

type RegisterAppResponse struct{
	AppID uuid.UUID `json:"app_id"`
	Name string `json:"name"`
	AppSecret string `json:"app_secret"`
	JWTSecret string `json:"jwt_secret"`
}