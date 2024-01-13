package sessiondatabase

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
	"github.com/thestoicway/backend/user_service/internal/user/model"
	customerrors "github.com/thestoicway/custom_errors"
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
