package sessiondatabase

import (
	"context"
)

// DeleteSession removes a session from the Redis store using the JWT ID.
func (s *redisDatabase) DeleteSession(ctx context.Context, jwtID string) (err error) {
	err = s.redis.Del(ctx, jwtID).Err()

	if err != nil {
		s.logger.Errorw("failed to delete session from redis", "error", err)

		return err
	}

	return nil
}
