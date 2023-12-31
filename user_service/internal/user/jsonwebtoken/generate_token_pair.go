package jsonwebtoken

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

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
