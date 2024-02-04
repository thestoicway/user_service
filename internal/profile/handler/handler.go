package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	customerrors "github.com/thestoicway/custom_errors"
	"github.com/thestoicway/user_service/internal/middleware"
	"github.com/thestoicway/user_service/internal/profile/model"
	"github.com/thestoicway/user_service/internal/profile/service"
	"github.com/thestoicway/user_service/internal/usr/jsonwebtoken"
	"go.uber.org/zap"
)

type ProfileHandler interface {
	UpdateProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error
}

type profileHandlerImpl struct {
	logger  *zap.SugaredLogger
	service service.ProfileService
}

// NewProfileHandler creates a new userHandler
func NewProfileHandler(logger *zap.SugaredLogger, service service.ProfileService) ProfileHandler {
	return &profileHandlerImpl{
		logger:  logger,
		service: service,
	}
}

func Register(router *httprouter.Router, h ProfileHandler, jwt jsonwebtoken.JwtManager) {
	router.PUT("/api/v1/profile", middleware.DecodeJwtMiddleware(jwt, customerrors.HandlerWrapper(h.UpdateProfile)))
}

// UpdateProfile implements ProfileHandler.
func (h *profileHandlerImpl) UpdateProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	profile := &model.Profile{}

	err := json.NewDecoder(r.Body).Decode(profile)

	if err != nil {
		return customerrors.NewWrongInputError(err)
	}

	userId := r.Context().Value(middleware.UserIDKey)

	if userId == nil {
		// TODO: if I pass nil there is a null deference panic
		return customerrors.NewUnauthorizedError(errors.New("missing user id"))
	}

	uuid, err := uuid.Parse(userId.(string))

	if err != nil {
		return customerrors.NewWrongInputError(err)
	}

	profile.UserID = uuid

	// TODO: validate profile
	if err := h.service.UpdateProfile(profile); err != nil {
		return err
	}

	resp := customerrors.NewSuccessResponse(nil)
	w.Header().Set("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(resp)
}
