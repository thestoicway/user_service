package jsonwebtoken

import (
	"github.com/golang-jwt/jwt/v5"
	customerrors "github.com/thestoicway/backend/custom_errors/custom_errors"
)

func (manager *jwtManager) DecodeToken(token string) (claims *CustomClaims, err error) {
	// decode token
	decodedToken, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(manager.secret), nil
	})

	if err != nil {
		return nil, customerrors.NewUnauthorizedError(err.Error())
	}

	exp, err := decodedToken.Claims.GetExpirationTime()

	if err != nil {
		return nil, customerrors.NewUnauthorizedError(
			"can't get expiration time from claims: " + err.Error(),
		)
	}

	iat, err := decodedToken.Claims.GetIssuedAt()

	if err != nil {
		return nil, customerrors.NewUnauthorizedError(
			"can't get issued at time from claims: " + err.Error(),
		)
	}

	// check if the duration is correct
	if exp.Sub(iat.Time) != refreshTokenDuration {
		return nil, customerrors.NewUnauthorizedError(
			"token duration is incorrect",
		)
	}

	if claims, ok := decodedToken.Claims.(*CustomClaims); ok {
		return claims, nil
	} else {
		return nil, customerrors.NewUnauthorizedError(
			"can't get claims from token",
		)
	}
}
