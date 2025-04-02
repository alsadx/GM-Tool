package main

import (
	"fmt"
	"time"

	"github.com/alsadx/GM-Tool/internal/auth"
	"github.com/alsadx/GM-Tool/internal/config"
	"github.com/alsadx/GM-Tool/internal/hash"
	"github.com/alsadx/GM-Tool/internal/repositories/postgres"
	"github.com/alsadx/GM-Tool/internal/repository"
	"github.com/alsadx/GM-Tool/internal/service"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", cfg)

	dbPool, err := postgres.InitDbPool(&cfg.DB)
	if err != nil {
		panic(err)
	}

	hasher := hash.NewHasher()
	tokenManager, err := auth.NewManager("123")

	repositories := repository.NewRepositories(dbPool)

	serviceDeps := service.ServisesDeps{
		Repos:           repositories,
		TokenManager:    tokenManager,
		Hasher:          *hasher,
		AccessTokenTTL:     cfg.Auth.JWT.AccessTokenTTL,
		RefreshTokenTTL:    cfg.Auth.JWT.RefreshTokenTTL,
		AccessTokenTTL:  time.Minute * 15,
		RefreshTokenTTL: time.Hour * 24 * 7,
	}

}
