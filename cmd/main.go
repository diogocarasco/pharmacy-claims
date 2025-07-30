package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/diogocarasco/go-pharmacy-service/internal/api"
	"github.com/diogocarasco/go-pharmacy-service/internal/auth"
	"github.com/diogocarasco/go-pharmacy-service/internal/config"
	"github.com/diogocarasco/go-pharmacy-service/internal/database"
	"github.com/diogocarasco/go-pharmacy-service/internal/loader"
	"github.com/diogocarasco/go-pharmacy-service/internal/logger"
	"github.com/diogocarasco/go-pharmacy-service/internal/service"
)

// @title Pharmacy Claim Service API
// @version 1.0
// @description This is a sample service for managing pharmacy claims.
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	log := logger.NewLogger()
	log.Info("Starting pharmacy service...")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading configurations: %v", err)
	}

	if _, err := os.Stat("./data"); os.IsNotExist(err) {
		err = os.Mkdir("./data", 0755)
		if err != nil {
			log.Fatal("Error creating 'data' directory: %v", err)
		}
	}

	if _, err := os.Stat(cfg.ClaimsDataPath); os.IsNotExist(err) {
		err = os.MkdirAll(cfg.ClaimsDataPath, 0755)
		if err != nil {
			log.Fatal("Error creating claims directory '%s': %v", cfg.ClaimsDataPath, err)
		}
	}
	if _, err := os.Stat(cfg.RevertsDataPath); os.IsNotExist(err) {
		err = os.MkdirAll(cfg.RevertsDataPath, 0755)
		if err != nil {
			log.Fatal("Error creating reverts directory '%s': %v", cfg.RevertsDataPath, err)
		}
	}

	dbRepo, err := database.InitDB(cfg.DatabasePath)
	if err != nil {
		log.Fatal("Error initializing database: %v", err)
	}
	defer func() {
		if err := dbRepo.Close(); err != nil {
			log.Error("Error closing database: %v", err)
		}
	}()

	sqliteRepoImpl, ok := dbRepo.(*database.SQLiteRepository)
	if !ok {
		log.Fatal("Error: DBRepository is not an expected SQLiteRepository instance for migrations.")
	}

	log.Info("Applying database migrations...")
	err = database.ApplyMigrations(sqliteRepoImpl.DB)
	if err != nil {
		log.Fatal("Error applying database migrations: %v", err)
	}
	log.Info("Database migrations applied successfully.")

	log.Info("Starting CSV pharmacies loading...")
	err = loader.LoadPharmaciesFromCSV(cfg.PharmaciesCSVPath, dbRepo)
	if err != nil {
		log.Error("Error loading pharmacies from CSV: %v", err)
	}
	log.Info("CSV pharmacies loading completed.")

	claimLoader := loader.NewClaimLoader(dbRepo)
	log.Info("Starting claims loading from directory: %s...", cfg.ClaimsDataPath)
	if err := claimLoader.LoadAndSaveClaimsFromDir(cfg.ClaimsDataPath); err != nil {
		log.Error("Error loading and saving claims: %v", err)
	}
	log.Info("Claims loading completed.")

	revertLoader := loader.NewRevertLoader(dbRepo)
	log.Info("Starting reverts loading from directory: %s...", cfg.RevertsDataPath)
	if err := revertLoader.LoadAndSaveRevertsFromDir(cfg.RevertsDataPath); err != nil {
		log.Error("Error loading and saving reverts: %v", err)
	}
	log.Info("Reverts loading completed.")

	claimService := service.NewClaimService(log, dbRepo)
	authenticator := auth.NewAuthenticator(cfg.AuthToken, log)
	handlers := api.NewHandlers(claimService, log)

	routerCfg := api.RouterConfig{
		Handlers:      handlers,
		Authenticator: authenticator,
	}
	mux := api.NewRouter(routerCfg)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	go func() {
		log.Info("HTTP server starting on port %s...", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("HTTP server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-quit

	log.Info("Shutdown signal received. Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown (timeout or error): %v", err)
	} else {
		log.Info("Server shut down gracefully.")
	}

	log.Info("Pharmacy service terminated.")
}
