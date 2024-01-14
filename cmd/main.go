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
	"github.com/thestoicway/user_service/internal/config"
	"github.com/thestoicway/user_service/internal/usr/api"
	"github.com/thestoicway/user_service/internal/usr/database"
	"github.com/thestoicway/user_service/internal/usr/jsonwebtoken"
	"github.com/thestoicway/user_service/internal/usr/service"
	sessiondatabase "github.com/thestoicway/user_service/internal/usr/session_database"
	"go.uber.org/zap"
)

func main() {
	cfg := config.NewConfig()
	logger := initializeLogger()

	db, err := database.OpenDB(cfg)
	if err != nil {
		logger.Fatalf("can't open database: %v", err)
	}

	userDb := database.NewUserDatabase(logger.Named("userDatabase"), db)
	jwtManager := jsonwebtoken.NewJwtManager(logger.Named("jwt_manager"), cfg.JwtSecret)
	redisClient := createRedisClient(cfg)
	session := sessiondatabase.NewRedisDatabase(logger.Named("sessionDatabase"), redisClient)

	userService := createUserService(logger.Named("userService"), cfg, userDb, jwtManager, session)
	userHandler := api.NewUserHandler(logger.Named("userHandler"), userService)

	router := httprouter.New()
	api.Register(router, userHandler)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerConfig.Port),
		Handler: router,
	}

	signalsChan := make(chan os.Signal, 1)
	signal.Notify(signalsChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Infof("starting server on port %d", cfg.ServerConfig.Port)
		if err := server.ListenAndServe(); err != nil {
			logger.Fatalf("can't start server: %v", err)
		}
	}()

	<-signalsChan
}

func initializeLogger() *zap.SugaredLogger {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer zapLogger.Sync()
	return zapLogger.Sugar()
}

func createRedisClient(cfg *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisConfig.RedisHost, cfg.RedisConfig.RedisPort),
		Password: cfg.RedisConfig.RedisPass,
		DB:       cfg.RedisConfig.RedisDB,
	})
}

func createUserService(logger *zap.SugaredLogger, cfg *config.Config, userDb database.UserDatabase, jwtManager jsonwebtoken.JwtManager, session sessiondatabase.SessionDatabase) service.UserService {
	userSvcParams := &service.UserServiceParams{
		Logger:     logger,
		Database:   userDb,
		JwtManager: jwtManager,
		Session:    session,
	}
	return service.NewUserService(userSvcParams)
}
