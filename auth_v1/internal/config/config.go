package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		DB          DbConfig `yaml:"postgres"`
		HttpServer        HTTPConfig `yaml:"http"`
		Auth        AuthConfig `yaml:"auth"`
	}

	DbConfig struct {
		Name     string `yaml:"dbName"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		MaxConn  int32  `yaml:"maxConn" env-default:"10"`
	}

	AuthConfig struct {
		AccessTokenTTL  time.Duration `yaml:"acesssTokenTTL" env-default:"120m"`
		RefreshTokenTTL time.Duration `yaml:"refreshTokenTTL" env-default:"43200m"`
		// SigningKey
	}

	HTTPConfig struct {
		Host               string        `yaml:"host" env-default:"127.0.0.1"`
		Port               string        `yaml:"port" env-default:"8080"`
		// Address     string        `yaml:"address" env-default:"8080"`
		ReadTimeout     time.Duration `yaml:"readTimeout" env-default:"10s"`
		WriteTimeout time.Duration `yaml:"writeTimeout" env-default:"10s"`
		MaxHeaderMegabytes int `yaml:"maxHeaderMegabytes" env-default:"1"`
	}
)

func LoadConfig() (*Config, error) {

	var cfg Config
	if err := cleanenv.ReadConfig("config/main.yaml", &cfg); err != nil {
		return nil, fmt.Errorf("failed to read config: %s", err)
	}

	return &cfg, nil
}
