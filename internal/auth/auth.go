package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"

	"github.com/geneowak/go-expense-tracker/internal/hashing"
)

func HashPassword(password string) (string, error) {
	// in prod env we should probably use hashing.ProdParams
	return hashing.CreateHash(password, hashing.DefaultParams)
}

func CheckPasswordHash(password, hash string) (bool, error) {
	return hashing.ComparePasswordAndHash(password, hash)
}

func GetBearerToken(headers http.Header) (string, error) {
	header := headers.Get("Authorization")
	if header == "" {
		return "", errors.New("No Authorization header included in request")
	}
	token, ok := strings.CutPrefix(header, "Bearer ")
	if !ok {
		return "", errors.New("Malformed Authorization header")
	}

	return token, nil
}

func MakeRefreshToken() string {
	key := make([]byte, 32)
	rand.Read(key)

	return hex.EncodeToString(key)
}

func GetAPIKey(headers http.Header) (string, error) {
	header := headers.Get("Authorization")
	if header == "" {
		return "", errors.New("No Authorization header included in request")
	}
	key, ok := strings.CutPrefix(header, "ApiKey ")
	if !ok {
		return "", errors.New("Malformed Authorization header")
	}

	return key, nil
}
