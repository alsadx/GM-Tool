package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Env     string `yaml:"env" env-default:"local"`
		DB DBConfig `yaml:"postgres" env-required:"true"`
		GRPC GrpcConfig `yaml:"grpc"`
		Auth AuthConfig `yaml:"auth"`
	}

	DBConfig struct { 
		Host     string `yaml:"host"`
		Port     string    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
		MaxConn  int32    `yaml:"max_conn"`
	}

	GrpcConfig struct {
		Port     int    `yaml:"port"`
		Timeout  time.Duration `yaml:"timeout"`
	}

	AuthConfig struct {
		AccessTokenTTL  time.Duration `yaml:"accessTokenTTL"`
		RefreshTokenTTL time.Duration `yaml:"refreshTokenTTL"`
	}
)

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}
	
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist")
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}

func MustLoadByPath(path string) *Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist " + path)
	}

	var cfg Config
	
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}