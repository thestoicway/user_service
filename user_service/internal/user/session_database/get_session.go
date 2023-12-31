package sessiondatabase

import (
	"context"

	"github.com/redis/go-redis/v9"
	customerrors "github.com/thestoicway/backend/custom_errors"
	"github.com/thestoicway/backend/user_service/internal/user/model"
)

// GetSession retrieves a session from the Redis store using the JWT ID.
func (s *redisDatabase) GetSession(ctx context.Context, jwtID string) (session *model.Session, err error) {
	refreshToken, err := s.redis.Get(ctx, jwtID).Result()

	if err != nil {
		s.logger.Errorw("failed to get session from redis", "error", err)

		if err == redis.Nil {
			return nil, customerrors.NewUnauthorizedError("invalid session")
		}

		return nil, err
	}

	session = &model.Session{
		JwtID:        jwtID,
		RefreshToken: refreshToken,
	}

	return session, nil
}
