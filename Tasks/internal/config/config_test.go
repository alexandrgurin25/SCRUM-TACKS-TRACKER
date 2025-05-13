package config

import "testing"

func TestLoad(t *testing.T) {
	path := "test.env"
	target := &Config{
		Env:             "dev",
		Port:            50050,
		AuthServiceAddr: "localhost:50052",
		DB: DBConfig{
			Host:           "localhost",
			Port:           5432,
			User:           "postgres",
			Password:       "root",
			Name:           "tasks",
			MinPools:       3,
			MaxPools:       5,
			MigrationsPath: "./migrations",
		},
		Redis: RedisConfig{
			Addr:     "localhost:6379",
			User:     "root",
			Password: "root",
			DB:       0,
		},
	}
	cfg, err := Load(path)
	if err != nil || *cfg != *target {
		t.Log(*cfg)
		t.Fatalf("failed to load config: %v", err)
	}
}
