package sessiondatabase

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/thestoicway/backend/user_service/internal/user/model"
	"go.uber.org/zap"
)

// SessionDatabase is an interface that declares methods for session management in a database.
type SessionDatabase interface {
	AddSession(ctx context.Context, session *model.Session) (err error)
	GetSession(ctx context.Context, jwtID string) (session *model.Session, err error)
	ReplaceSession(ctx context.Context, oldJwtID string, session *model.Session) (err error)
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

// AddSession adds a new session to the Redis store.
func (s *redisDatabase) AddSession(ctx context.Context, session *model.Session) (err error) {
	err = s.redis.Set(ctx, session.JwtID, session.RefreshToken, session.ExpirationTime).Err()

	if err != nil {
		s.logger.Errorw("failed to add session to redis", "error", err)
		return err
	}

	return nil
}

// GetSession retrieves a session from the Redis store using the JWT ID.
func (s *redisDatabase) GetSession(ctx context.Context, jwtID string) (session *model.Session, err error) {
	refreshToken, err := s.redis.Get(ctx, jwtID).Result()

	if err != nil {
		s.logger.Errorw("failed to get session from redis", "error", err)
		return nil, err
	}

	session = &model.Session{
		JwtID:        jwtID,
		RefreshToken: refreshToken,
	}

	return session, nil
}

// DeleteSession removes a session from the Redis store using the JWT ID.
func (s *redisDatabase) DeleteSession(ctx context.Context, jwtID string) (err error) {
	err = s.redis.Del(ctx, jwtID).Err()

	if err != nil {
		s.logger.Errorw("failed to delete session from redis", "error", err)
		return err
	}

	return nil
}

// ReplaceSession replaces an old session with a new one in the Redis store.
func (s *redisDatabase) ReplaceSession(ctx context.Context, oldJwtID string, session *model.Session) (err error) {
	// Start a Redis transaction
	txPipeline := s.redis.TxPipeline()

	// Queue commands in the transaction
	txPipeline.Del(ctx, oldJwtID)
	txPipeline.Set(ctx, session.JwtID, session.RefreshToken, session.ExpirationTime)

	// Execute the transaction
	_, err = txPipeline.Exec(ctx)
	if err != nil {
		s.logger.Errorw("failed to execute transaction in redis", "error", err)
		return err
	}

	return nil
}
