package sessiondatabase

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/thestoicway/user_service/internal/model"
	"go.uber.org/zap"
)

// SessionDatabase is an interface that declares methods for session management in a database.
type SessionDatabase interface {
	AddSession(ctx context.Context, session *model.Session) (err error)
	GetSession(ctx context.Context, jwtID string) (session *model.Session, err error)
	ReplaceSession(ctx context.Context, oldSession *model.Session, session *model.Session) (err error)
	DeleteSession(ctx context.Context, jwtID string) (err error)
}

// redisDatabase is a struct that implements the SessionDatabase interface with Redis as the backend.
type redisDatabase struct {
	logger *zap.SugaredLogger // Logger for logging messages.
	redis  *redis.Client      // Redis client for interacting with the Redis store.
}

// NewRedisDatabase is a constructor for sessionDatabase.
func NewRedisDatabase(logger *zap.SugaredLogger, redis *redis.Client) SessionDatabase {
	return &redisDatabase{
		logger: logger,
		redis:  redis,
	}
}
