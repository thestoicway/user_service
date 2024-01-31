package sessionstorage

import (
	"context"

	"github.com/thestoicway/user_service/internal/usr/model"
)

type InMemoryDatabaseParams struct {
	Sessions map[string]*model.Session
}

type inMemoryDatabase struct {
	*InMemoryDatabaseParams
}

// AddSession implements sessionstorage.
func (db *inMemoryDatabase) AddSession(ctx context.Context, session *model.Session) (err error) {
	db.Sessions[session.JwtID] = session
	return nil
}

// DeleteSession implements sessionstorage.
func (db *inMemoryDatabase) DeleteSession(ctx context.Context, jwtID string) (err error) {
	delete(db.Sessions, jwtID)
	return nil
}

// GetSession implements sessionstorage.
func (db *inMemoryDatabase) GetSession(ctx context.Context, jwtID string) (session *model.Session, err error) {
	return db.Sessions[jwtID], nil
}

// ReplaceSession implements sessionstorage.
func (db *inMemoryDatabase) ReplaceSession(ctx context.Context, oldSession *model.Session, session *model.Session) (err error) {
	delete(db.Sessions, oldSession.JwtID)
	db.Sessions[session.JwtID] = session
	return nil
}

func NewInMemoryDatabase(params *InMemoryDatabaseParams) SessionStorage {
	return &inMemoryDatabase{
		params,
	}
}
