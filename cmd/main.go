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
	profileDatabase "github.com/thestoicway/user_service/internal/profile/database"
	profileHandler "github.com/thestoicway/user_service/internal/profile/handler"
	profileService "github.com/thestoicway/user_service/internal/profile/service"
	"github.com/thestoicway/user_service/internal/usr/api"
	"github.com/thestoicway/user_service/internal/usr/database"
	"github.com/thestoicway/user_service/internal/usr/jsonwebtoken"
	"github.com/thestoicway/user_service/internal/usr/service"
	sessionstorage "github.com/thestoicway/user_service/internal/usr/session_storage"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.NewConfig()
	logger := initializeLogger()

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatalf("can't open database: %v", err)
	}

	userDb := database.NewUserDatabase(logger.Named("userDatabase"), db)
	jwtManager := jsonwebtoken.NewJwtManager(logger.Named("jwt_manager"), cfg.JwtSecret)
	redisClient := createRedisClient(cfg)
	session := sessionstorage.NewRedisDatabase(logger.Named("sessionstorage"), redisClient)

	userService := createUserService(logger.Named("userService"), cfg, userDb, jwtManager, session)
	userHandler := api.NewUserHandler(logger.Named("userHandler"), userService)

	profileDB := profileDatabase.NewProfileDatabase(logger.Named("profileDatabase"), db)
	profileService := createProfileService(logger.Named("profileService"), cfg, profileDB)
	pHandler := profileHandler.NewProfileHandler(logger.Named("profileHandler"), profileService)

	router := httprouter.New()
	api.Register(router, userHandler)
	profileHandler.Register(router, pHandler, jwtManager)

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

func createProfileService(logger *zap.SugaredLogger, cfg *config.Config, profileDb profileDatabase.ProfileDatabase) profileService.ProfileService {
	profileSvcParams := &profileService.ProfileServiceParams{
		Logger:   logger,
		Database: profileDb,
	}
	return profileService.NewProfileService(profileSvcParams)
}

func createRedisClient(cfg *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisConfig.RedisHost, cfg.RedisConfig.RedisPort),
		Password: cfg.RedisConfig.RedisPass,
		DB:       cfg.RedisConfig.RedisDB,
	})
}

func createUserService(logger *zap.SugaredLogger, cfg *config.Config, userDb database.UserDatabase, jwtManager jsonwebtoken.JwtManager, session sessionstorage.SessionStorage) service.UserService {
	userSvcParams := &service.UserServiceParams{
		Logger:     logger,
		Database:   userDb,
		JwtManager: jwtManager,
		Session:    session,
	}
	return service.NewUserService(userSvcParams)
}

func openDB(config *config.Config) (*gorm.DB, error) {
	// "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable"
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		config.PostgresDatabase.PostgresHost,
		config.PostgresDatabase.PostgresUser,
		config.PostgresDatabase.PostgresPass,
		config.PostgresDatabase.PostgresDB,
		config.PostgresDatabase.PostgresPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})

	if err != nil {
		return nil, err
	}

	return db, err
}
