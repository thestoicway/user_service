package model

import "time"

// Session represents a user session.
// This should be stored in K/V database like Redis.
type Session struct {
	// JwtID is the ID of the refresh token.
	JwtID string
	// RefreshToken is the refresh token.
	RefreshToken string
	// ExpirationTime is the expiration time of the refresh token.
	ExpirationTime time.Duration
}
