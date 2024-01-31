// Package jsonwebtoken provides utilities for managing JSON Web Tokens.
package jsonwebtoken

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	customerrors "github.com/thestoicway/custom_errors"
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

// CustomClaims represents the custom claims in a JSON Web Token.
type CustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// JwtManager is an interface for managing JSON Web Tokens.
type JwtManager interface {
	// GenerateTokenPair generates a new pair of access and refresh tokens.
	GenerateTokenPair(uuid uuid.UUID) (tokenPair *TokenPair, info *AdditionalInfo, err error)
	// DecodeToken decodes a token and returns the claims.
	DecodeToken(token string) (claims *CustomClaims, err error)
}

// jwtManagerImpl is an implementation of the JwtManager interface.
type jwtManagerImpl struct {
	logger *zap.SugaredLogger
	secret string
}

// NewJwtManager creates a new instance of JwtManager.
func NewJwtManager(logger *zap.SugaredLogger, secret string) JwtManager {
	return &jwtManagerImpl{
		logger: logger,
		secret: secret,
	}
}

// GenerateTokenPair is a method on jwtManager that generates a pair of JWT tokens (access and refresh)
// for a given user ID. It returns the token pair, additional info, and any error that occurred.
func (manager *jwtManagerImpl) GenerateTokenPair(id uuid.UUID) (tokenPair *TokenPair, info *AdditionalInfo, err error) {
	// Get the current time, which will be used as the issued time of the tokens
	issuedAt := time.Now()
	userID := id.String()

	// Create an access token with claims
	aToken := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "user_service",
			Subject:   "access_token",
			ExpiresAt: jwt.NewNumericDate(issuedAt.Add(accessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
		},
	})

	// Sign the access token with the secret key
	accessToken, err := aToken.SignedString([]byte(manager.secret))
	if err != nil {
		return nil, nil, err
	}

	// Generate a new UUID for the refresh token
	rTokenID := uuid.New().String()

	// Create a refresh token with claims
	rToken := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "user_service",
			Subject:   "refresh_token",
			ExpiresAt: jwt.NewNumericDate(issuedAt.Add(refreshTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ID:        rTokenID,
		},
	})

	// Sign the refresh token with the secret key
	refreshToken, err := rToken.SignedString([]byte(manager.secret))
	if err != nil {
		return nil, nil, err
	}

	// Return the pair of tokens and the additional info
	return &TokenPair{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}, &AdditionalInfo{
			RefreshTokenID:        rTokenID,
			RefreshExpirationTime: refreshTokenDuration,
		}, nil
}

// DecodeToken is a method on jwtManager that takes a JWT token as a string,
// decodes it, validates it, and returns the claims in the token or an error.
func (manager *jwtManagerImpl) DecodeToken(token string) (claims *CustomClaims, err error) {
	// Parse the token with the claims
	decodedToken, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// This function returns the secret key for validating the token
		return []byte(manager.secret), nil
	})

	// If there was an error in parsing, return an UnauthorizedError
	if err != nil {
		return nil, customerrors.NewUnauthorizedError(err)
	}

	// Assert that the claims in the token are of type *CustomClaims
	claims, ok := decodedToken.Claims.(*CustomClaims)

	// If the assertion was not ok, return an UnauthorizedError
	if !ok || claims == nil {
		return nil, customerrors.NewUnauthorizedError(errors.New("token claims are not of type *CustomClaims"))
	}

	// Check if the duration between the issued time and the expiry time of the token is correct
	if claims.ExpiresAt.Sub(claims.IssuedAt.Time) != refreshTokenDuration {
		return nil, customerrors.NewUnauthorizedError(errors.New("token duration is not correct"))
	}

	// If everything is ok, return the claims in the token
	return claims, nil
}
