package sessionstorage

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
	customerrors "github.com/thestoicway/custom_errors"
	"github.com/thestoicway/user_service/internal/usr/model"
	"go.uber.org/zap"
)

// redisDatabase is a struct that implements the sessionstorage interface with Redis as the backend.
type redisDatabase struct {
	logger *zap.SugaredLogger // Logger for logging messages.
	redis  *redis.Client      // Redis client for interacting with the Redis store.
}

// NewRedisDatabase is a constructor for sessionstorage.
func NewRedisDatabase(logger *zap.SugaredLogger, redis *redis.Client) SessionStorage {
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

// DeleteSession deletes a session from the Redis store.
func (s *redisDatabase) DeleteSession(ctx context.Context, jwtID string) (err error) {
	err = s.redis.Del(ctx, jwtID).Err()

	if err != nil {
		s.logger.Errorw("failed to delete session from redis", "error", err)

		return err
	}

	return nil
}

// GetSession retrieves a session from the Redis store using the JWT ID.
func (s *redisDatabase) GetSession(ctx context.Context, jwtID string) (session *model.Session, err error) {
	refreshToken, err := s.redis.Get(ctx, jwtID).Result()

	if err != nil {
		if err == redis.Nil {
			return nil, customerrors.NewUnauthorizedError(
				errors.New("session not found"),
			)
		}

		return nil, err
	}

	session = &model.Session{
		JwtID:        jwtID,
		RefreshToken: refreshToken,
	}

	return session, nil
}

// ReplaceSession replaces an old session with a new one in the Redis store.
func (s *redisDatabase) ReplaceSession(ctx context.Context, oldSession *model.Session, session *model.Session) (err error) {
	// Start a Redis transaction
	txPipeline := s.redis.TxPipeline()

	// Queue commands in the transaction
	txPipeline.Del(ctx, oldSession.JwtID)
	txPipeline.Set(ctx, session.JwtID, session.RefreshToken, session.ExpirationTime)

	// Execute the transaction
	_, err = txPipeline.Exec(ctx)

	if err != nil {
		s.logger.Errorw("failed to replace session in redis", "error", err)

		if err == redis.Nil {
			return customerrors.NewUnauthorizedError(err)
		}

		return err

	}

	return nil
}
