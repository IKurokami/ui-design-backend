package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"ui-design-backend/internal/auth"
	"ui-design-backend/internal/config"
	"ui-design-backend/internal/domain"
	"ui-design-backend/internal/mailer"
	"ui-design-backend/internal/middleware"
	"ui-design-backend/internal/repository"
	"ui-design-backend/internal/service"
	transporthttp "ui-design-backend/internal/transport/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("[INIT] Starting ui-design auth backend in %s mode...", cfg.Env)

	// 1. Connect to database
	db, err := repository.NewDB(cfg.DBDriver, cfg.DBDSN)
	if err != nil {
		log.Fatalf("Failed to connect database (%s): %v", cfg.DBDriver, err)
	}
	defer db.Close()
	log.Printf("[DB] Connected to %s database successfully", cfg.DBDriver)

	// 2. Initialize Repositories
	userRepo := repository.NewUserRepo(db)
	if cfg.AdminEmail != "" {
		adminUser, findErr := userRepo.FindByEmail(context.Background(), strings.ToLower(strings.TrimSpace(cfg.AdminEmail)))
		if findErr == nil {
			findErr = userRepo.UpdateRole(context.Background(), adminUser.ID, domain.RoleAdmin)
		}
		if findErr != nil {
			log.Printf("[ADMIN] Could not promote ADMIN_EMAIL=%s: %v", cfg.AdminEmail, findErr)
		} else {
			log.Printf("[ADMIN] Promoted %s to admin", cfg.AdminEmail)
		}
	}
	sessionRepo := repository.NewSessionRepo(db)
	resetRepo := repository.NewPasswordResetRepo(db)
	profileRepo := repository.NewProfileRepo(db)
	storyRepo := repository.NewStoryRepo(db)
	ebookRepo := repository.NewEbookRepo(db)

	// 3. Initialize Services
	tokenService := auth.NewTokenService(cfg.JWTSecret, cfg.AccessTokenTTL)
	mailService := mailer.NewMailer(cfg)
	authService := service.NewAuthService(cfg, userRepo, sessionRepo, resetRepo, tokenService, mailService)
	profileService := service.NewProfileService(userRepo, profileRepo)
	storyService := service.NewStoryService(storyRepo, userRepo)
	ebookService := service.NewEbookService(ebookRepo)
	managementService := service.NewManagementService(userRepo, storyService)

	// 4. Initialize Middleware & Handlers
	authMiddleware := middleware.NewAuthMiddleware(cfg, tokenService, authService, userRepo)
	authHandler := transporthttp.NewAuthHandler(authService, authMiddleware, transporthttp.NewAdminHandler(managementService))
	profileHandler := transporthttp.NewProfileHandler(profileService, "./uploads")
	storyHandler := transporthttp.NewStoryHandler(storyService)
	ebookHandler := transporthttp.NewEbookHandler(ebookService)

	// 5. Router
	router := transporthttp.NewRouter(cfg, authHandler, profileHandler, storyHandler, authMiddleware, ebookHandler)

	// 6. Server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Graceful shutdown channel
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("[SERVER] Server listening on http://localhost:%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	<-shutdownChan
	log.Println("[SERVER] Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("[SERVER] Server stopped.")
}
