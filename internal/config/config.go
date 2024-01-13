package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	JwtSecret        string `env:"JWT_SECRET" env-required:"true"`
	ServerConfig     ServerConfig
	PostgresDatabase PostgresDatabase
	RedisConfig      RedisDatabase
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

type RedisDatabase struct {
	RedisPort int    `env:"REDIS_PORT" env-default:"6379"`
	RedisHost string `env:"REDIS_HOST" env-default:"localhost"`
	RedisPass string `env:"REDIS_PASS" env-default:""`
	RedisDB   int    `env:"REDIS_DB" env-default:"0"`
}

// NewConfig returns a new Config struct from ENV variables
func NewConfig() *Config {
	cfg := &Config{}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		log.Fatalf("can't read config: %v", err)
	}

	return cfg
}
