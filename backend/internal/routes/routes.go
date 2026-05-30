package routes

import (
	"context"
	"time"

	"github.com/brandon-v-seeters/go-silk-wave/internal/auth"
	"github.com/brandon-v-seeters/go-silk-wave/internal/config"
	"github.com/brandon-v-seeters/go-silk-wave/internal/database"
	"github.com/brandon-v-seeters/go-silk-wave/internal/handlers"
	"github.com/brandon-v-seeters/go-silk-wave/internal/logger"
	"github.com/brandon-v-seeters/go-silk-wave/internal/middleware"
	"github.com/brandon-v-seeters/go-silk-wave/internal/repository"
	"github.com/brandon-v-seeters/go-silk-wave/internal/storage"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Setup(db *database.ArangoDB, cfg *config.Config) *gin.Engine {
	router := gin.Default()

	// Configure CORS
	if cfg.IsDev() {
		// In development, allow localhost origins
		router.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
		router.SetTrustedProxies([]string{"127.0.0.1", "localhost"})
	} else {
		// In production, configure allowed origins as needed
		router.SetTrustedProxies(nil)
	}

	storageClient, err := storage.NewClient(context.Background(), storage.R2Config{
		AccountID:       cfg.R2AccountID,
		AccessKeyID:     cfg.R2AccessKeyID,
		SecretAccessKey: cfg.R2SecretAccessKey,
		BucketName:      cfg.R2BucketName,
	})

	if err != nil {
		logger.Fatal("Failed to create storage client", zap.Error(err))
	}

	// Initialize services
	jwtService := auth.NewJWTService(cfg.JWTSecret)
	passwordService := auth.NewPasswordService(cfg.PasswordSecret)
	authService := auth.NewAuthService()

	userRepo := repository.NewUserRepository(db)
	releaseRepo := repository.NewReleaseRepository(db)
	artistRepo := repository.NewArtistRepository(db)
	accessRepo := repository.NewAccessRepository(db)
	userHandler := handlers.NewUserHandler(userRepo, db)
	settingsHandler := handlers.NewSettingsHandler(passwordService, db)

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler()
	authHandler := handlers.NewAuthHandler(userRepo, jwtService, passwordService, authService)
	artistHandler := handlers.NewArtistHandler(artistRepo)
	releaseHandler := handlers.NewReleaseHandler(db, accessRepo, releaseRepo, storageClient)

	// Health check
	router.GET("/health", healthHandler.Health)

	// API routes
	api := router.Group("/api")
	{
		// Public auth routes
		api.POST("/login", authHandler.Login)
		api.POST("/register", authHandler.Register)
		api.POST("/logout", authHandler.Logout)
		api.GET("/token")

		// Public routes
		api.GET("/releases", releaseHandler.GetReleases)
		api.GET("/artists/:artistSlug", artistHandler.GetArtistBySlug)
		api.GET("/artists/:artistSlug/releases/:releaseSlug", middleware.OptionalAuthMiddleware(jwtService), releaseHandler.GetReleaseByArtistAndSlug)

		// Protected routes - require authentication
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(jwtService))
		{
			// User routes
			protected.GET("/user", userHandler.GetCurrentUser)
			protected.POST("/user/mode", userHandler.UpdateMode)

			// Settings routes
			protected.POST("/settings/email", settingsHandler.UpdateEmail)
			protected.POST("/settings/delete-account", settingsHandler.DeleteAccount)
			protected.POST("/settings/password", settingsHandler.UpdatePassword)

			// Artist registration
			protected.POST("/register/artist-name", artistHandler.RegisterArtistName)
			protected.POST("/artists/:artistId/follow", artistHandler.FollowArtist)
			protected.DELETE("/artists/:artistId/follow", artistHandler.UnfollowArtist)

			// Release routes
			protected.POST("/releases/draft", releaseHandler.SaveDraftRelease)
			protected.POST("/releases/:releaseId/confirm", releaseHandler.ConfirmDraftRelease)
			protected.POST("/releases/:releaseId/publish", releaseHandler.PublishRelease)
			protected.POST("/releases/:releaseId/archive", releaseHandler.ArchiveRelease)
			protected.DELETE("/releases/:releaseId", releaseHandler.DeleteRelease)
			protected.GET("/releases/drafts", releaseHandler.GetDraftsByArtistKey)
			protected.GET("/releases/drafts/:releaseKey", releaseHandler.GetDraftByKey)
			protected.GET("/releases/:releaseId/pending/preview", releaseHandler.GetPendingPreview)
			protected.POST("/releases/:releaseId/pending", releaseHandler.StagePendingEdit)
			protected.DELETE("/releases/:releaseId/pending", releaseHandler.DiscardPendingEdit)
			protected.POST("/releases/:releaseId/pending/publish", releaseHandler.PublishPendingEdit)

			// Avatar upload
			protected.POST("/upload/avatar", releaseHandler.UploadAvatar)

			// User management (admin-like)
			// users := protected.Group("/users")
			// users.Use(middleware.AdminMiddleware())
			// {
			// 	users.GET("", userHandler.GetAllUsers)
			// 	users.GET("/:key", userHandler.GetUser)
			// 	users.PUT("/:key", userHandler.UpdateUser)
			// 	users.DELETE("/:key", userHandler.DeleteUser)
			// }
		}
	}

	return router
}
