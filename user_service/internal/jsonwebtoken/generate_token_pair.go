package jsonwebtoken

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (manager *jwtManager) GenerateTokenPair(id uuid.UUID) (tokenPair *TokenPair, info *AdditionalInfo, err error) {
	issuedAt := time.Now()
	userID := id.String()

	// Creates access and refresh tokens with claims
	aToken := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "user_service",
			Subject:   "access_token",
			ExpiresAt: jwt.NewNumericDate(issuedAt.Add(accessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
		},
	})

	accessToken, err := aToken.SignedString([]byte(manager.secret))

	if err != nil {
		return nil, nil, err
	}

	rTokenID := uuid.New().String()

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

	refreshToken, err := rToken.SignedString([]byte(manager.secret))

	if err != nil {
		return nil, nil, err
	}

	return &TokenPair{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}, &AdditionalInfo{
			RefreshTokenID:        rTokenID,
			RefreshExpirationTime: refreshTokenDuration,
		}, nil
}
