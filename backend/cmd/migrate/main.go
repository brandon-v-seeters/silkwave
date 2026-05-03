package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/brandon-v-seeters/go-silk-wave/internal/config"
	"github.com/brandon-v-seeters/go-silk-wave/internal/database"
	"github.com/brandon-v-seeters/go-silk-wave/internal/logger"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// Parse flags
	dryRun := flag.Bool("dry-run", false, "Show what would be done without making changes")
	flag.Parse()

	// Load .env file if it exists
	_ = godotenv.Load()

	// Load configuration
	cfg := config.Load()

	// Initialize logger
	logger.Init(logger.GetEnv())
	defer logger.Sync()

	if *dryRun {
		fmt.Println("=== DRY RUN MODE ===")
		fmt.Println("Collections and indexes that would be created/updated:")
		for _, def := range database.Schema {
			fmt.Printf("\nCollection: %s (%s)\n", def.Name, def.Type)
			for _, idx := range def.Indices {
				fmt.Printf("  - Index on %v (unique: %v, sparse: %v)\n", idx.Fields, idx.Unique, idx.Sparse)
			}
		}
		os.Exit(0)
	}

	// Initialize database connection
	db, err := database.NewArangoDB(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to ArangoDB", zap.Error(err))
	}
	defer db.Close()

	fmt.Println("Running database migrations...")

	// Run database migrations
	if err := db.Migrate(context.Background()); err != nil {
		logger.Fatal("Failed to run database migrations", zap.Error(err))
	}

	fmt.Println("✓ Migrations completed successfully!")
}
