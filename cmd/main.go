package main

import (
	"context"
	"log"

	_ "github.com/lib/pq"

	"github.com/luique16/quitocoin/ent"
	"github.com/luique16/quitocoin/internal/config"
	"github.com/luique16/quitocoin/internal/domain/user"
	"github.com/luique16/quitocoin/internal/handler"
	"github.com/luique16/quitocoin/internal/provider"
	"github.com/luique16/quitocoin/internal/usecase"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	client, err := ent.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer client.Close()

	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("migration: %v", err)
	}

	hasher := provider.NewPasswordHasher()
	idGen := provider.NewIdGenerator()
	jwtProvider := provider.NewJWTProvider(cfg.JWTSecret)

	repo := user.NewRepository(client)
	userService := user.NewService(repo, hasher, idGen)

	registerUC := usecase.NewRegisterUseCase(userService, jwtProvider)
	loginUC := usecase.NewLoginUseCase(repo, hasher, jwtProvider)
	getMeUC := usecase.NewGetUserUseCase(userService)
	updateMeUC := usecase.NewUpdateUserUseCase(userService)
	deleteMeUC := usecase.NewDeleteUserUseCase(userService)

	router := handler.NewRouter(registerUC, loginUC, getMeUC, updateMeUC, deleteMeUC, jwtProvider)

	log.Printf("server running on :%s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("server: %v", err)
	}
}
