package main

import (
	"context"
	"log"

	_ "github.com/lib/pq"

	"github.com/luique16/quitocoin/ent"
	"github.com/luique16/quitocoin/internal/config"
	"github.com/luique16/quitocoin/internal/domain/block"
	"github.com/luique16/quitocoin/internal/domain/transaction"
	"github.com/luique16/quitocoin/internal/domain/user"
	"github.com/luique16/quitocoin/internal/domain/utxo"
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

	blockRepo := block.NewRepository(client)
	blockService := block.NewService(3, blockRepo)

	memPool := transaction.NewMemPool()

	utxoRepo := utxo.NewRepository()
	utxoService := utxo.NewService(utxoRepo)

	userRepo := user.NewRepository(client)
	userService := user.NewService(userRepo, hasher, idGen)

	initializer := usecase.NewInitializerUseCase(blockService, memPool, utxoService)
	registerUC := usecase.NewRegisterUseCase(userService, jwtProvider)
	loginUC := usecase.NewLoginUseCase(userRepo, hasher, jwtProvider)
	getMeUC := usecase.NewGetUserUseCase(userService, utxoService)
	updateMeUC := usecase.NewUpdateUserUseCase(userService)
	updatePasswordUC := usecase.NewUpdatePasswordUseCase(userService)
	deleteMeUC := usecase.NewDeleteUserUseCase(userService)
	mineBlockUC := usecase.NewMineBlockUseCase(blockService, utxoService, memPool, 3)

	err = initializer.Execute(context.Background())

	if err != nil {
		log.Fatalf("initializer: %v", err)
	}

	router := handler.NewRouter(registerUC, loginUC, getMeUC, updateMeUC, updatePasswordUC, deleteMeUC, mineBlockUC, jwtProvider)

	log.Printf("server running on :%s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("server: %v", err)
	}
}
