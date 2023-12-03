package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/thestoicway/backend/user_service/internal/service"
	"go.uber.org/zap"
)

type UserHandler interface {
	SignIn(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
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
	router.GET("/signin", h.SignIn)
}
