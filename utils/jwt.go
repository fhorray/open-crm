package utils

import (
	"errors"
	"fmt"
	"open-crm/pkg/config"
	"time"

	"maps"

	"github.com/golang-jwt/jwt/v5"
)

// Generate
func GenerateJWT(data map[string]any, expires time.Duration, audience, issuer string) (tokenString string, err error) {
	if data == nil {
		return "", errors.New("provide data")
	}

	// Inicializa os claims
	claims := jwt.MapClaims{}

	// Copia os dados fornecidos
	maps.Copy(claims, data)

	// Define a expiração padrão
	claims["exp"] = time.Now().Add(expires).Unix()

	if audience != "" {
		claims["aud"] = audience
	} else {
		claims["aud"] = config.Cfg.AUTH.JWT_AUDIENCE
	}

	if issuer != "" {
		claims["iss"] = issuer
	} else {
		claims["iss"] = config.Cfg.AUTH.JWT_ISSUER
	}

	// Cria e assina o token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err = token.SignedString([]byte(config.Cfg.KEYS.JWT_SECRET))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Parse Token
func ParseJWT(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		// Verify algorithym
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return []byte(config.Cfg.KEYS.JWT_SECRET), nil
		// ADD AUDIENCE & ISSUER
	}, jwt.WithAudience(config.Cfg.AUTH.JWT_AUDIENCE), jwt.WithIssuer((config.Cfg.AUTH.JWT_ISSUER)))
}
