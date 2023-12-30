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
	PostgresPort int    `env:"POSTGRES_PORT" env-default:"5432"`
	PostgresHost string `env:"POSTGRES_HOST" env-default:"localhost"`
	PostgresUser string `env:"POSTGRES_USER" env-default:"postgres"`
	PostgresPass string `env:"POSTGRES_PASS" env-default:"postgres"`
	PostgresDB   string `env:"POSTGRES_DB" env-default:"postgres"`
}

// NewConfig returns a new Config struct from ENV variables
func NewConfig() *Config {
	cfg := &Config{}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		log.Fatalf("can't read config: %v", err)
	}

	return cfg
}
