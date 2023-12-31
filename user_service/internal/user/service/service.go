package service

import (
	"context"

	"github.com/thestoicway/backend/user_service/internal/config"
	"github.com/thestoicway/backend/user_service/internal/user/database"
	"github.com/thestoicway/backend/user_service/internal/user/jsonwebtoken"
	"github.com/thestoicway/backend/user_service/internal/user/model"
	sessiondatabase "github.com/thestoicway/backend/user_service/internal/user/session_database"
	"go.uber.org/zap"
)

type UserService interface {
	SignIn(ctx context.Context, user *model.User) (tokenPair *jsonwebtoken.TokenPair, err error)
	SignUp(ctx context.Context, user *model.User) (tokenPair *jsonwebtoken.TokenPair, err error)
	Refresh(ctx context.Context, refreshToken string) (tokenPair *jsonwebtoken.TokenPair, err error)
}

type userServiceImpl struct {
	*UserServiceParams
}

type UserServiceParams struct {
	Logger     *zap.SugaredLogger
	Config     *config.Config
	Database   database.UserDatabase
	JwtManager jsonwebtoken.JwtManager
	Session    sessiondatabase.SessionDatabase
}

func NewUserService(p *UserServiceParams) UserService {
	return &userServiceImpl{
		UserServiceParams: p,
	}
}
