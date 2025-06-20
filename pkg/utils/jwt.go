package utils

import (
	"errors"
	"fmt"
	"open-crm/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserId uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

// Generates a new JWT token
func GenerateToken(data map[string]any, expires time.Duration, audience, issuer string) (string, *Claims, error) {
	if data == nil {
		return "", nil, errors.New("provide data")
	}

	// Extrai e valida user_id do payload
	userIDRaw, ok := data["user_id"]
	if !ok {
		return "", nil, errors.New("missing user id in token data (expected key: 'user_id')")
	}

	var userID uuid.UUID
	switch v := userIDRaw.(type) {
	case uuid.UUID:
		userID = v
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			return "", nil, fmt.Errorf("invalid user id format: %w", err)
		}
		userID = parsed
	default:
		return "", nil, errors.New("invalid user id type in token data")
	}

	// Prepara os registered claims
	claims := &Claims{
		UserId: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expires)),
			Issuer:    issuer,
			Audience:  jwt.ClaimStrings{audience},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Preenche audience e issuer default se estiverem vazios
	if claims.Issuer == "" {
		claims.Issuer = config.Cfg.AUTH.JWT_ISSUER
	}
	if len(claims.Audience) == 0 || claims.Audience[0] == "" {
		claims.Audience = []string{config.Cfg.AUTH.JWT_AUDIENCE}
	}

	// Cria e assina o token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Cfg.KEYS.JWT_SECRET))
	if err != nil {
		return "", nil, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, claims, nil
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

// Extract Token
func ExtractToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{},
		func(token *jwt.Token) (any, error) {
			return []byte(secret), nil
		},
	)

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
