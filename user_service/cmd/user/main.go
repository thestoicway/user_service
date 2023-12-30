package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/julienschmidt/httprouter"
	"github.com/redis/go-redis/v9"
	"github.com/thestoicway/backend/user_service/internal/config"
	"github.com/thestoicway/backend/user_service/internal/jsonwebtoken"
	"github.com/thestoicway/backend/user_service/internal/user/api"
	"github.com/thestoicway/backend/user_service/internal/user/database"
	"github.com/thestoicway/backend/user_service/internal/user/service"
	sessiondatabase "github.com/thestoicway/backend/user_service/internal/user/session_database"
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

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisConfig.RedisHost, cfg.RedisConfig.RedisPort),
		Password: cfg.RedisConfig.RedisPass,
		DB:       cfg.RedisConfig.RedisDB,
	})

	session := sessiondatabase.NewRedisDatabase(sugar.Named("sessionDatabase"), redisClient)

	userSvcParams := &service.UserServiceParams{
		Logger:     sugar.Named("userService"),
		Config:     cfg,
		Database:   userDb,
		JwtManager: jwtManager,
		Session:    session,
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
