package model

import "time"

// Session represents a user session.
// This should be stored in K/V database like Redis.
type Session struct {
	JwtID          string
	RefreshToken   string
	ExpirationTime time.Duration
}
