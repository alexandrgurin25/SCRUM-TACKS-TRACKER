package config

import "testing"

func TestLoad(t *testing.T) {
	path := "test.env"
	target := &Config{
		Env:             "local",
		Port:            8000,
		AuthServiceAddr: "1234",
		DB: DBConfig{
			Host:           "1234",
			Port:           1234,
			User:           "1234",
			Password:       "1234",
			Name:           "1234",
			MinPools:       1234,
			MaxPools:       4321,
			MigrationsPath: "./migrations",
		},
		Redis: RedisConfig{
			Addr:     "1234",
			User:     "1234",
			Password: "1234",
			DB:       0,
		},
	}

	cfg, err := Load(path)
	if err != nil || *cfg != *target {
		t.Log(*cfg)
		t.Fatalf("failed to load config: %v", err)
	}
}
