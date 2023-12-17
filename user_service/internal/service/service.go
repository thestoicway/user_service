package service

import (
	"context"

	"github.com/thestoicway/backend/user_service/internal/config"
	"github.com/thestoicway/backend/user_service/internal/database"
	jsonwebtoken "github.com/thestoicway/backend/user_service/internal/jwt"
	"github.com/thestoicway/backend/user_service/internal/model"
	"go.uber.org/zap"
)

type UserService interface {
	SignIn(ctx context.Context, user *model.User) (tokenPair *jsonwebtoken.TokenPair, err error)
	SignUp(ctx context.Context, user *model.User) (tokenPair *jsonwebtoken.TokenPair, err error)
}

type userService struct {
	logger     *zap.SugaredLogger
	config     *config.Config
	database   database.UserDatabase
	jwtManager jsonwebtoken.JwtManager
}

type UserServiceParams struct {
	Logger     *zap.SugaredLogger
	Config     *config.Config
	Database   database.UserDatabase
	JwtManager jsonwebtoken.JwtManager
}

func NewUserService(p *UserServiceParams) UserService {
	return &userService{
		logger:     p.Logger,
		database:   p.Database,
		config:     p.Config,
		jwtManager: p.JwtManager,
	}
}
