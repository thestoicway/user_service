package sessiondatabase

import (
	"context"

	"github.com/thestoicway/backend/user_service/internal/user/model"
	"go.uber.org/zap"
)

type SessionDatabase interface {
	AddSession(ctx context.Context, session *model.Session) (err error)
	GetSession(ctx context.Context, jwtID string) (session *model.Session, err error)
	DeleteSession(ctx context.Context, jwtID string) (err error)
}

type sessionDatabase struct {
	logger *zap.SugaredLogger
	kv     KeyValueDatabase
}

func NewSessionDatabase(logger *zap.SugaredLogger, kv KeyValueDatabase) SessionDatabase {
	return &sessionDatabase{
		logger: logger,
		kv:     kv,
	}
}

// AddSession implements SessionDatabase.
func (s *sessionDatabase) AddSession(ctx context.Context, session *model.Session) (err error) {
	err = s.kv.Set(ctx, session.JwtID, session.RefreshToken, session.ExpirationTime)

	if err != nil {
		return err
	}

	return nil
}

// DeleteSession implements SessionDatabase.
func (s *sessionDatabase) DeleteSession(ctx context.Context, jwtID string) (err error) {
	err = s.kv.Delete(ctx, jwtID)

	if err != nil {
		return err
	}

	return nil
}

// GetSession implements SessionDatabase.
func (s *sessionDatabase) GetSession(ctx context.Context, jwtID string) (session *model.Session, err error) {
	value, err := s.kv.Get(ctx, jwtID)

	if err != nil {
		return nil, err
	}

	session = &model.Session{
		JwtID:        jwtID,
		RefreshToken: value.(string),
	}

	return session, nil
}
