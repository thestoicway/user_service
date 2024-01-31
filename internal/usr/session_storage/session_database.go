package sessionstorage

import (
	"context"

	"github.com/thestoicway/user_service/internal/usr/model"
)

// SessionStorage is an interface that declares methods for session management in a database.
type SessionStorage interface {
	// AddSession adds a new session to the storage.
	AddSession(ctx context.Context, session *model.Session) (err error)
	// GetSession retrieves a session from the storage using the JWT ID.
	GetSession(ctx context.Context, jwtID string) (session *model.Session, err error)
	// ReplaceSession replaces an old session with a new session in the storage.
	ReplaceSession(ctx context.Context, oldSession *model.Session, session *model.Session) (err error)
	// DeleteSession deletes a session from the storage.
	DeleteSession(ctx context.Context, jwtID string) (err error)
}
