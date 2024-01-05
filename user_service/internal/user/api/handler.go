package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	customerrors "github.com/thestoicway/backend/custom_errors"
	"github.com/thestoicway/backend/user_service/internal/user/service"
	"go.uber.org/zap"
)

type UserHandler interface {
	SignIn(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error
	SignUp(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error
	Refresh(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error
	SignOut(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error
}

type userHandlerImpl struct {
	logger  *zap.SugaredLogger
	service service.UserService
}

// NewUserHandler creates a new userHandler
func NewUserHandler(logger *zap.SugaredLogger, service service.UserService) UserHandler {
	return &userHandlerImpl{
		logger:  logger,
		service: service,
	}
}

func Register(router *httprouter.Router, h UserHandler) {
	router.POST("/api/v1/signin", customerrors.HandlerWrapper(h.SignIn))
	router.POST("/api/v1/signup", customerrors.HandlerWrapper(h.SignUp))
	router.POST("/api/v1/refresh", customerrors.HandlerWrapper(h.Refresh))
	router.DELETE("/api/v1/signout", customerrors.HandlerWrapper(h.SignOut))
}
