package jsonwebtoken

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type JwtManager interface {
	GenerateTokenPair(uuid uuid.UUID) (tokenPair *TokenPair, err error)
}

type jwtManager struct {
	logger *zap.SugaredLogger
	secret string
}

func (manager *jwtManager) GenerateTokenPair(uuid uuid.UUID) (tokenPair *TokenPair, err error) {
	iat := time.Now()

	id := uuid.String()

	// Creates access and refresh tokens with claims
	aToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"iat": iat.Unix(),
		"exp": iat.Add(time.Minute * 30).Unix(),
	})

	accessToken, err := aToken.SignedString([]byte(manager.secret))

	if err != nil {
		return nil, err
	}

	rToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"iat": iat.Unix(),
		"exp": iat.Add(time.Hour * 24 * 30).Unix(),
	})

	refreshToken, err := rToken.SignedString([]byte(manager.secret))

	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func NewJwtManager(logger *zap.SugaredLogger, secret string) JwtManager {
	return &jwtManager{
		logger: logger,
		secret: secret,
	}
}
