package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Postgres struct {
	Host     string `yaml:"POSTGRES_HOST" env:"POSTGRES_HOST"`
	Port     uint16 `yaml:"POSTGRES_PORT" env:"POSTGRES_PORT"`
	Username string `yaml:"POSTGRES_USER" env:"POSTGRES_USER"`
	Password string `yaml:"POSTGRES_PASSWORD" env:"POSTGRES_PASSWORD"`
	Database string `yaml:"POSTGRES_DB" env:"POSTGRES_DB"`

	MinConns int32 `yaml:"POSTGRES_MIN_CONN" env:"POSTGRES_MIN_CONN"`
	MaxConns int32 `yaml:"POSTGRES_MAX_CONN" env:"POSTGRES_MAX_CONN"`
}

type Auth struct {
	AccessTokenSecret  string        `env:"ACCESS_TOKEN_SECRET,required"`
	RefreshTokenSecret string        `env:"REFRESH_TOKEN_SECRET,required"`
	AccessTokenTTL     time.Duration `env:"ACCESS_TOKEN_TTL" env-default:"15m"`
	RefreshTokenTTL    time.Duration `env:"REFRESH_TOKEN_TTL" env-default:"168h"` // 7 days
}

type Grpc struct {
	GRPCPort int `env:"GPRC_PORT"`
}
type Config struct {
	Postgres
	Auth
	Grpc
}

func New() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig("./config/.env", &cfg); err != nil {
		return nil, fmt.Errorf("config can not read env file %v", err)
	}
	return &cfg, nil
}

func NewTest(path string) (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		if errNoRead := cleanenv.ReadEnv(&cfg); errNoRead != nil {
			return nil, errNoRead
		}
		return nil, err

	}

	return &cfg, nil
}
