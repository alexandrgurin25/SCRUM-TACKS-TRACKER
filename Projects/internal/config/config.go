package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	GRPCPort        int            `yaml:"GRPC_PORT" env:"GRPC_PORT"`
	Postgres        PostgresConfig `env:"POSTGRES"`
	AuthServiceAddr string         `env:"AUTH_SERVICE_ADDR"`
}

type PostgresConfig struct {
	Host           string `yaml:"POSTGRES_HOST" env:"POSTGRES_HOST"`
	Port           uint16 `yaml:"POSTGRES_PORT" env:"POSTGRES_PORT"`
	Username       string `yaml:"POSTGRES_USER" env:"POSTGRES_USER"`
	Password       string `yaml:"POSTGRES_PASSWORD" env:"POSTGRES_PASSWORD"`
	Database       string `yaml:"POSTGRES_DB" env:"POSTGRES_DB"`
	MinConns       int32  `yaml:"POSTGRES_MIN_CONN" env:"POSTGRES_MIN_CONN"`
	MaxConns       int32  `yaml:"POSTGRES_MAX_CONN" env:"POSTGRES_MAX_CONN"`
	MigrationsPath string `yaml:"POSTGRES_MIGRATIONS_PATH" env:"POSTGRES_MIGRATIONS_PATH"`
}

func New() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
		return nil, fmt.Errorf("config can not read env file %v", err)
	}
	return &cfg, nil
}
