package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env                string `env:"ENV" env-default:"prod"`
	Port               int    `env:"PORT" env-default:"8081"`
	AuthServiceAddr    string `env:"AUTH_SERVICE_ADDR" env-required:"true"`
	TaskServiceAddr    string `env:"TASK_SERVICE_ADDR" env-required:"true"`
	BoardServiceAddr   string `env:"BOARD_SERVICE_ADDR" env-required:"true"`
	ProjectServiceAddr string `env:"PROJECT_SERVICE_ADDR" env-required:"true"`
	CommentServiceAddr string `env:"COMMENT_SERVICE_ADDR" env-required:"true"`
}

func MustLoad(path string) *Config {
	cfg, err := Load(path)
	if err != nil {
		panic(err)
	}

	return cfg
}

func Load(path string) (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
