package service

import (
	"context"

	"github.com/thestoicway/user_service/internal/usr/database"
	"github.com/thestoicway/user_service/internal/usr/jsonwebtoken"
	"github.com/thestoicway/user_service/internal/usr/model"
	sessionstorage "github.com/thestoicway/user_service/internal/usr/session_storage"
	"go.uber.org/zap"
)

type UserService interface {
	SignIn(ctx context.Context, user *model.User) (tokenPair *jsonwebtoken.TokenPair, err error)
	SignUp(ctx context.Context, user *model.User) (tokenPair *jsonwebtoken.TokenPair, err error)
	SignOut(ctx context.Context, refreshToken string) (err error)
	Refresh(ctx context.Context, refreshToken string) (tokenPair *jsonwebtoken.TokenPair, err error)
}

type userServiceImpl struct {
	*UserServiceParams
}

type UserServiceParams struct {
	Logger     *zap.SugaredLogger
	Database   database.UserDatabase
	JwtManager jsonwebtoken.JwtManager
	Session    sessionstorage.SessionStorage
}

func NewUserService(p *UserServiceParams) UserService {
	return &userServiceImpl{
		UserServiceParams: p,
	}
}
