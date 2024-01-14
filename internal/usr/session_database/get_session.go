package sessiondatabase

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
	customerrors "github.com/thestoicway/custom_errors"
	"github.com/thestoicway/user_service/internal/usr/model"
)

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
