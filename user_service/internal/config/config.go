package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	JwtSecret        string `env:"JWT_SECRET" env-required:"true"`
	ServerConfig     ServerConfig
	PostgresDatabase PostgresDatabase
}

type ServerConfig struct {
	Port int `env:"PORT" env-default:"8080"`
}

type PostgresDatabase struct {
	PostgresURL string `env:"POSTGRES_URL" env-required:"true"`
}

// NewConfig returns a new Config struct from ENV variables
func NewConfig() *Config {
	cfg := &Config{}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		log.Fatalf("can't read config: %v", err)
	}

	return cfg
}
