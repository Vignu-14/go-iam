package handler

import (
	"encoding/json"
	"go-iam/internal/pkg/security"
	"net/http"
	"uuid"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AppHandler struct {
	DB *pgxpool.Pool
}

type RegisterAppRequest struct {
	Name string `json:"name"`
}

type RegisterAppResponse struct {
	AppID     uuid.UUID `json:"app_id"`
	Name      string    `json:"name"`
	AppSecret string    `json:"app_secret"`
	JWTSecret string    `json:"jwt_secret"`
}

func (h *AppHandler) RegisterApp(w http.ResponseWriter, r *http.Request) {
	var req RegisterAppRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		http.Error(w, `{"error":"app name is required"}`, http.StatusBadRequest)
		return
	}

	rawAppSecret, err := security.GenerateRandomSecret()
	if err != nil {
		http.Error(w, `{"error":"failed to generate secret"}`, http.StatusInternalServerError)
		return
	}

	jwtSecret, err := security.GenerateRandomSecret()
	if err != nil {
		http.Error(w, `{"error":"failed to generate jwt secret"}`, http.StatusInternalServerError)
		return
	}

	appSecretHash, err := security.HashBcrypt(rawAppSecret)
	if err != nil {
		http.Error(w, `{"error":"failed to hash secret"}`, http.StatusInternalServerError)
		return
	}

	// Start PostgreSQL Transaction
	tx, err := h.DB.Begin(r.Context())
	if err != nil {
		http.Error(w, `{"error":"db error"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(r.Context())

	// 1. Insert the App
	var appID uuid.UUID
	err = tx.QueryRow(r.Context(),
		`INSERT INTO apps (name, app_secret_hash, jwt_secret)
		 VALUES ($1, $2, $3) RETURNING id`,
		req.Name, appSecretHash, jwtSecret,
	).Scan(&appID)
	if err != nil {
		http.Error(w, `{"error":"failed to create app"}`, http.StatusInternalServerError)
		return
	}

	// 2. Insert starter roles ('user' as default, 'admin' as non-default)
	_, err = tx.Exec(r.Context(),
		`INSERT INTO roles (app_id, name, is_default)
		 VALUES ($1, 'user', TRUE), ($1, 'admin', FALSE)`,
		appID,
	)
	if err != nil {
		http.Error(w, `{"error":"failed to create default roles"}`, http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(r.Context()); err != nil {
		http.Error(w, `{"error":"failed to commit transaction"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(RegisterAppResponse{
		AppID:     appID,
		Name:      req.Name,
		AppSecret: rawAppSecret,
		JWTSecret: jwtSecret,
	})
}
