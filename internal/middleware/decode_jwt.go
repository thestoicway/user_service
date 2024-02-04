package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-errors/errors"

	"github.com/julienschmidt/httprouter"
	customerrors "github.com/thestoicway/custom_errors"
	"github.com/thestoicway/user_service/internal/usr/jsonwebtoken"
)

// TODO: make this accept *CustomError
func handleErrorResponse(w http.ResponseWriter, customError error) {
	w.Header().Set("Content-Type", "application/json")
	err := customError.(*customerrors.CustomError)
	w.WriteHeader(err.StatusCode())
	resp := customerrors.NewErrorResponse(err)
	jsonEncoder := json.NewEncoder(w)
	jsonEncoder.Encode(resp)
}

type UserId string

const UserIDKey UserId = "user_id"

// This middleware decodes the JWT token from "Authorization" header
// and sets the user id in the request context.
// TODO: THINK ABOUT COMPATIBILITY WITH OTHER MIDDLEWARES
func DecodeJwtMiddleware(jwtManager jsonwebtoken.JwtManager, h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		token := r.Header.Get("Authorization")

		if token == "" {
			resp := customerrors.NewUnauthorizedError(
				errors.New("missing token"),
			)
			handleErrorResponse(w, resp)
			return
		}

		claims, err := jwtManager.DecodeToken(token)

		if err != nil {
			resp := customerrors.NewUnauthorizedError(err)
			handleErrorResponse(w, resp)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
		print("claims.UserID: ", claims.UserID)
		r = r.WithContext(ctx)

		h(w, r, ps)
	}
}
