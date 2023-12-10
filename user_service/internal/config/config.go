package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	JwtSecret    string
	ServerConfig ServerConfig
}

type ServerConfig struct {
	Port int
}

// NewConfig returns a new Config struct from ENV variables
func NewConfig() *Config {
	jwtSecret := os.Getenv("JWT_SECRET")

	if jwtSecret == "" {
		log.Fatal("JWT_SECRET must be set")
	}

	serverConfig := newServerConfig()

	return &Config{
		JwtSecret:    jwtSecret,
		ServerConfig: serverConfig,
	}
}

func newServerConfig() ServerConfig {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	intPort, err := strconv.Atoi(port)

	if err != nil {
		log.Fatalf("can't convert port to int: %v", err)
	}

	return ServerConfig{
		Port: intPort,
	}
}
