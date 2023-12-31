package sessiondatabase

import (
	"context"

	"github.com/redis/go-redis/v9"
	customerrors "github.com/thestoicway/backend/custom_errors"
	"github.com/thestoicway/backend/user_service/internal/user/model"
)

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
			return customerrors.NewUnauthorizedError("invalid session")
		}

		return err

	}

	return nil
}
