// Package jsonwebtoken provides utilities for managing JSON Web Tokens.
package jsonwebtoken

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Constants for access and refresh token durations.
const (
	// accessTokenDuration specifies the duration that an access token is valid for.
	accessTokenDuration = time.Minute * 30
	// refreshTokenDuration specifies the duration that a refresh token is valid for.
	refreshTokenDuration = time.Hour * 24 * 30
)

// TokenPair represents a pair of access and refresh tokens.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// AdditionalInfo represents additional information related to a token.
type AdditionalInfo struct {
	RefreshTokenID        string        `json:"refresh_token_id"`
	RefreshExpirationTime time.Duration `json:"refresh_expiration_time"`
}

// JwtManager is an interface for managing JSON Web Tokens.
type JwtManager interface {
	// GenerateTokenPair generates a new pair of access and refresh tokens.
	GenerateTokenPair(uuid uuid.UUID) (tokenPair *TokenPair, info *AdditionalInfo, err error)
	// DecodeToken decodes a token and returns the claims.
	DecodeToken(token string) (claims *CustomClaims, err error)
}

// jwtManager is an implementation of the JwtManager interface.
type jwtManager struct {
	logger *zap.SugaredLogger
	secret string
}

// NewJwtManager creates a new instance of JwtManager.
func NewJwtManager(logger *zap.SugaredLogger, secret string) JwtManager {
	return &jwtManager{
		logger: logger,
		secret: secret,
	}
}

// CustomClaims represents the custom claims in a JSON Web Token.
type CustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}
