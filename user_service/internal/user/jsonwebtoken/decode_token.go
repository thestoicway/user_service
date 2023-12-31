package jsonwebtoken

import (
	"github.com/golang-jwt/jwt/v5"
	customerrors "github.com/thestoicway/backend/custom_errors"
)

// DecodeToken is a method on jwtManager that takes a JWT token as a string,
// decodes it, validates it, and returns the claims in the token or an error.
func (manager *jwtManagerImpl) DecodeToken(token string) (claims *CustomClaims, err error) {
	// Parse the token with the claims
	decodedToken, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// This function returns the secret key for validating the token
		return []byte(manager.secret), nil
	})

	// If there was an error in parsing, return an UnauthorizedError
	if err != nil {
		return nil, customerrors.NewUnauthorizedError(err.Error())
	}

	// Assert that the claims in the token are of type *CustomClaims
	claims, ok := decodedToken.Claims.(*CustomClaims)

	// If the assertion was not ok, return an UnauthorizedError
	if !ok || claims == nil {
		return nil, customerrors.NewUnauthorizedError("can't get claims from token")
	}

	// Check if the duration between the issued time and the expiry time of the token is correct
	if claims.ExpiresAt.Sub(claims.IssuedAt.Time) != refreshTokenDuration {
		return nil, customerrors.NewUnauthorizedError("token duration is incorrect")
	}

	// If everything is ok, return the claims in the token
	return claims, nil
}
