package security

import (
	"crypto/rand"
	"encoding/hex"
)


func GenerateRandomSecret()(string ,error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil{
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}