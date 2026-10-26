package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type App struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	AppSecretHash string    `json:"-"` // Never expose hashes in JSON
	JWTSecret     string    `json:"-"` // Used to sign JWTs for this app
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
}

type Role struct {
	ID        uuid.UUID `json:"id"`
	AppID     uuid.UUID `json:"app_id"`
	Name      string    `json:"name"`
	IsDefault bool      `json:"is_default"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	ID           uuid.UUID `json:"id"`
	AppID        uuid.UUID `json:"app_id"`
	RoleID       uuid.UUID `json:"role_id"`
	RoleName     string    `json:"role"` // Joined from roles.name
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
}

// The exact payload inside the JWT token sent to the Bookstore app
type AccessTokenClaims struct {
	UserID uuid.UUID `json:"user_id"`
	AppID  uuid.UUID `json:"app_id"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}
