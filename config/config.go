package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	APP  `yaml:"app"`
	HTTP `yaml:"http"`
	PG   `yaml:"pg"`
	Log  `yaml:"logger"`
}

type APP struct {
	MainStorage string `env-required:"false" yaml:"main_storage" env:"MAIN_STORAGE"`
}

type HTTP struct {
	Port string `env-required:"false" yaml:"port" env:"PORT"`
}

type PG struct {
	PoolMax    int    `env-required:"false" yaml:"pool_max" env:"PG_POOL_MAX"`
	URL        string `env-required:"false" yaml:"pg_url"   env:"PG_URL"`
	Schema_Url string `env-required:"false" yaml:"schema_url"   env:"SCHEMA_URL"`
}

type Log struct {
	Level string `env-required:"false" yaml:"log_level" env:"LOG_LEVEL"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yaml", cfg)
	if err != nil {
		return nil, fmt.Errorf("read config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
