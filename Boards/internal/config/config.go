package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	GRPCPort        int            `env:"GRPC_PORT"`
	Postgres        PostgresConfig `env:"POSTGRES"`
	AuthServiceAddr string         `env:"AUTH_SERVICE_ADDR"`
}

type PostgresConfig struct {
	Host           string `env:"POSTGRES_HOST"`
	Port           uint16 `env:"POSTGRES_PORT"`
	Username       string `env:"POSTGRES_USER"`
	Password       string `env:"POSTGRES_PASSWORD"`
	Database       string `env:"POSTGRES_DB"`
	MinConns       int32  `env:"POSTGRES_MIN_CONN"`
	MaxConns       int32  `env:"POSTGRES_MAX_CONN"`
	MigrationsPath string `env:"POSTGRES_MIGRATIONS_PATH"`
}

func New() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
		return nil, fmt.Errorf("config can not read env file %v", err)
	}
	return &cfg, nil
}
