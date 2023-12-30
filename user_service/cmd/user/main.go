package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/julienschmidt/httprouter"
	"github.com/thestoicway/backend/user_service/internal/api"
	"github.com/thestoicway/backend/user_service/internal/config"
	"github.com/thestoicway/backend/user_service/internal/database"
	"github.com/thestoicway/backend/user_service/internal/jsonwebtoken"
	"github.com/thestoicway/backend/user_service/internal/service"
	"go.uber.org/zap"
)

func main() {
	zap, err := zap.NewProduction()
	cfg := config.NewConfig()

	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	sugar := zap.Sugar()
	defer sugar.Sync()

	db, err := database.OpenDB(cfg)

	if err != nil {
		sugar.Fatalf("can't open database: %v", err)
	}

	userDb := database.NewUserDatabase(
		sugar.Named("userDatabase"),
		db,
	)

	jwtManager := jsonwebtoken.NewJwtManager(sugar.Named("jwt_manager"), cfg.JwtSecret)

	userSvcParams := &service.UserServiceParams{
		Logger:     sugar.Named("userService"),
		Config:     cfg,
		Database:   userDb,
		JwtManager: jwtManager,
	}

	userService := service.NewUserService(userSvcParams)

	userHandler := api.NewUserHandler(
		sugar.Named("userHandler"),
		userService,
	)

	router := httprouter.New()

	userHandler.Register(router)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerConfig.Port),
		Handler: router,
	}

	signalsChan := make(chan os.Signal, 1)

	signal.Notify(signalsChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("starting server on port %d", cfg.ServerConfig.Port)
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("can't start server: %v", err)
		}
	}()

	<-signalsChan
}
