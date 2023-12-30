package jsonwebtoken

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	accessTokenDuration  = time.Minute * 30
	refreshTokenDuration = time.Hour * 24 * 30
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AdditionalInfo struct {
	RefreshTokenID        string        `json:"refresh_token_id"`
	RefreshExpirationTime time.Duration `json:"refresh_expiration_time"`
}

type JwtManager interface {
	GenerateTokenPair(uuid uuid.UUID) (tokenPair *TokenPair, info *AdditionalInfo, err error)
	DecodeToken(token string) (claims *CustomClaims, err error)
}

type jwtManager struct {
	logger *zap.SugaredLogger
	secret string
}

func NewJwtManager(logger *zap.SugaredLogger, secret string) JwtManager {
	return &jwtManager{
		logger: logger,
		secret: secret,
	}
}

type CustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}
