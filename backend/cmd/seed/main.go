package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/brandon-v-seeters/go-silk-wave/internal/auth"
	"github.com/brandon-v-seeters/go-silk-wave/internal/config"
	"github.com/brandon-v-seeters/go-silk-wave/internal/database"
	"github.com/brandon-v-seeters/go-silk-wave/internal/logger"
	"github.com/brandon-v-seeters/go-silk-wave/internal/seed"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	resetSeed := flag.Bool("reset-seed", false, "delete seed-owned records before recreating them")
	namespace := flag.String("namespace", seed.DefaultNamespace, "seed namespace used to stamp and reset records")
	password := flag.String("password", seed.DefaultPassword, "password for seeded login users")
	flag.Parse()

	_ = godotenv.Load()

	cfg := config.Load()
	logger.Init(logger.GetEnv())
	defer logger.Sync()

	db, err := database.NewArangoDB(cfg)
	if err != nil {
		logger.Fatal("failed to connect to ArangoDB", zap.Error(err))
	}
	defer db.Close()

	passwords := auth.NewPasswordService(cfg.PasswordSecret)
	seeder := seed.NewSeeder(db, passwords)

	summary, err := seeder.Run(context.Background(), seed.Options{
		Namespace: *namespace,
		Password:  *password,
		ResetSeed: *resetSeed,
	})
	if err != nil {
		logger.Fatal("failed to seed database", zap.Error(err))
	}

	fmt.Println(summary.String())
}
