package sessiondatabase

import (
	"context"

	"github.com/thestoicway/user_service/internal/model"
)

// AddSession adds a new session to the Redis store.
func (s *redisDatabase) AddSession(ctx context.Context, session *model.Session) (err error) {
	err = s.redis.Set(ctx, session.JwtID, session.RefreshToken, session.ExpirationTime).Err()

	if err != nil {
		s.logger.Errorw("failed to add session to redis", "error", err)
		return err
	}

	return nil
}
