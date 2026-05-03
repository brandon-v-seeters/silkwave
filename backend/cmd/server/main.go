package main

import (
	"github.com/brandon-v-seeters/go-silk-wave/internal/config"
	"github.com/brandon-v-seeters/go-silk-wave/internal/database"
	"github.com/brandon-v-seeters/go-silk-wave/internal/logger"
	"github.com/brandon-v-seeters/go-silk-wave/internal/routes"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)


func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// Can't use logger yet, it's not initialized
	}

	// Load configuration
	cfg := config.Load()

	// Initialize logger
	logger.Init(logger.GetEnv())
	defer logger.Sync()

	// Initialize database connection
	db, err := database.NewArangoDB(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to ArangoDB", zap.Error(err))
	}
	defer db.Close()

	// Setup router and routes
	router := routes.Setup(db, cfg)

	// Start server
	logger.Info("Server starting", zap.String("port", cfg.ServerPort))
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
