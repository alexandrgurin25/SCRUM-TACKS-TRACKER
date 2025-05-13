package config

import (
	"errors"
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env             string      `env:"ENV" env-default:"prod"`
	Port            int         `env:"PORT" env-default:"50051"`
	AuthServiceAddr string      `env:"AUTH_SERICE_ADDR"`
	DB              DBConfig    `env:"DB"`
	Redis           RedisConfig `env:"REDIS"`
}

type DBConfig struct {
	Host           string `env:"DB_HOST" env-default:"localhost"`
	Port           int    `env:"DB_PORT" env-default:"5431"`
	User           string `env:"DB_USER" env-default:"postgres"`
	Password       string `env:"DB_PASSWORD" env-required:"true"`
	Name           string `env:"DB_NAME" env-default:"tasks"`
	MinPools       int    `env:"DB_MIN_POOLS" env-default:"3"`
	MaxPools       int    `env:"DB_MAX_POOLS" env-default:"5"`
	MigrationsPath string `env:"MIGRATIONS_PATH" env-default:"./migrations"`
}

type RedisConfig struct {
	Addr     string `env:"REDIS_ADDR" env-default:"localhost:6379"`
	User     string `env:"REDIS_USER" env-default:"root"`
	Password string `env:"REDIS_USER_PASSWORD" env-default:"root"`
	DB       int    `env:"REDIS_DB" env-default:"0"`
}

func MustLoad(path string) *Config {
	cfg, err := Load(path)
	if err != nil {
		panic(err)
	}
	return cfg
}

func Load(path string) (*Config, error) {
	if len(path) == 0 {
		return nil, errors.New("path must be not empty")
	}

	cfg := &Config{}

	err := cleanenv.ReadConfig(path, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return cfg, nil
}
