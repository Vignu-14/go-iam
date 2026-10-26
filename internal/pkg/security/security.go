package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"go-iam/internal/domain"
	"time"

	"github.com/google/uuid"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func GenerateRandomSecret() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// HashBcrypt hashes user passwords or app secrets
func HashBcrypt(plain string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(bytes), err
}

// VerifyBcrypt checks if a plain password/secret matches its bcrypt hash
func VerifyBcrypt(plain, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}

// HashSHA256 is used for refresh tokens so we can look them up by index in O(log N) time
func HashSHA256(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// GenerateAccessToken signs a 15-minute JWT using the App's own JWTSecret
func GenerateAccessToken(userID, appID uuid.UUID, role, appJWTSecret string) (string, error) {
	claims := domain.AccessTokenClaims{
		UserID: userID,
		AppID:  appID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(appJWTSecret))
}
