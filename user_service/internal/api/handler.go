package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	customerrors "github.com/thestoicway/backend/custom_errors/custom_errors"
	"github.com/thestoicway/backend/user_service/internal/service"
	"go.uber.org/zap"
)

type UserHandler interface {
	SignIn(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error
	SignUp(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error
	Refresh(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error
	SignOut(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error
	Register(router *httprouter.Router)
}

type userHandler struct {
	logger  *zap.SugaredLogger
	service service.UserService
}

// NewUserHandler creates a new userHandler
func NewUserHandler(logger *zap.SugaredLogger, service service.UserService) UserHandler {
	return &userHandler{
		logger:  logger,
		service: service,
	}
}

func (h *userHandler) Register(router *httprouter.Router) {
	router.POST("/signin", customerrors.HandlerWrapper(h.SignIn))
	router.POST("/signup", customerrors.HandlerWrapper(h.SignUp))
	router.POST("/refresh", customerrors.HandlerWrapper(h.Refresh))
	router.DELETE("/signout", customerrors.HandlerWrapper(h.SignOut))
}
