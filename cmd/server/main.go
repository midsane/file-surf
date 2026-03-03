package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/midsane/file-surf/internal/config"
	"github.com/midsane/file-surf/internal/database"
	"github.com/midsane/file-surf/internal/document"
	"github.com/midsane/file-surf/internal/server"
	"github.com/midsane/file-surf/internal/storage"

	// "github.com/midsane/file-surf/internal/storage"
	"github.com/midsane/file-surf/internal/tenant"
	"github.com/midsane/file-surf/internal/user"
)

func main() {
	// Load config
	godotenv.Load()
	cfg := config.Load()
	ctx := context.Background()

	dynamoClient, err := database.NewDynamoClient(ctx, cfg.AWSRegion)
	if err != nil {
		log.Fatalf("failed to init dynamodb: %v", err)
	}

	s3Client, err := storage.NewS3Client(ctx, cfg.AWSRegion)
	if err != nil {
		log.Fatalf("faild to init s3 %v", err)
	}

	s3Storage := storage.NewS3Storage(s3Client, cfg.S3Bucket)

	// Create router
	router := gin.New()
	router.Use(server.RequestLogger())
	router.Use(gin.Recovery())

	tenantRepo := tenant.NewRepository(dynamoClient, cfg.TenantTable)
	tenantService := tenant.NewService(tenantRepo)
	tenantHandler := tenant.NewHandler(tenantService)
	tenantHandler.RegisterRoutes(router)

	userRepo := user.NewRepository(dynamoClient, cfg.TenantTable)
	userService := user.NewService(userRepo, tenantRepo)
	userHandler := user.NewHandler(userService)
	userHandler.RegisterRoutes(router)

	docRepo := document.NewRepository(dynamoClient, cfg.TenantTable)
	docService := document.NewService(docRepo, tenantRepo, s3Storage)
	docHandler := document.NewHandler(docService)
	docHandler.RegisterRoutes(router)

	// Health route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Start server in goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Println("Server started on port", cfg.Port)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %s\n", err)
	}

	log.Println("Server exited properly")
}
