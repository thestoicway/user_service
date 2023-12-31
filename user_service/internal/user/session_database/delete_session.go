package sessiondatabase

import (
	"context"

	"github.com/redis/go-redis/v9"
	customerrors "github.com/thestoicway/backend/custom_errors"
)

// DeleteSession removes a session from the Redis store using the JWT ID.
func (s *redisDatabase) DeleteSession(ctx context.Context, jwtID string) (err error) {
	err = s.redis.Del(ctx, jwtID).Err()

	if err != nil {
		s.logger.Errorw("failed to delete session from redis", "error", err)

		if err == redis.Nil {
			return customerrors.NewUnauthorizedError("invalid session")
		}

		return err
	}

	return nil
}
