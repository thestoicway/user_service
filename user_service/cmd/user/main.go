package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/thestoicway/backend/user_service/internal/api"
	"github.com/thestoicway/backend/user_service/internal/database"
	"github.com/thestoicway/backend/user_service/internal/service"
	"go.uber.org/zap"
)

func main() {
	zap, err := zap.NewProduction()

	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	sugar := zap.Sugar()
	defer sugar.Sync()

	router := httprouter.New()

	db := database.NewUserDatabase(
		sugar.Named("userDatabase"),
	)

	userService := service.NewUserService(
		sugar.Named("userService"),
		db,
	)

	userHandler := api.NewUserHandler(
		sugar.Named("userHandler"),
		userService,
	)

	userHandler.Register(router)

	log.Fatal(http.ListenAndServe(":8080", router))
}
