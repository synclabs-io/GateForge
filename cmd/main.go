package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/synclabs-io/GateForge/internal/auth"
	"github.com/synclabs-io/GateForge/internal/cache"
	"github.com/synclabs-io/GateForge/internal/config"
	http_handlers "github.com/synclabs-io/GateForge/internal/handlers/http"
	"github.com/synclabs-io/GateForge/internal/repositories"
	redis_repository "github.com/synclabs-io/GateForge/internal/repositories/redis"
	"github.com/synclabs-io/GateForge/internal/service"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := run(logger); err != nil {
		logger.Error("server exited with error", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// --- Repositories ---
	client, err := cache.NewRedisClient(ctx, cfg)
	if err != nil {
		return err
	}
	defer client.Close()

	pool, err := pgxpool.New(ctx, cfg.Postgres.DSN())
	if err != nil {
		return err
	}
	defer pool.Close()

	userRepo := repositories.NewUsersRepository(pool)
	sessionRepo := redis_repository.NewSessionRepository(client)

	// --- Auth ---
	hasher := auth.NewPasswordHasher(64*1024, 1, 2, 32, 16)
	jwtMgr := auth.NewJWTManager(cfg.JWT.Secret, cfg.JWT.AccessExpiry)

	// --- Services ---
	usersSvc := service.NewUsersService(userRepo, sessionRepo, logger, jwtMgr, hasher, cfg.JWT.RefreshExpiry)

	// --- Handlers + routers ---
	usersHandler := http_handlers.NewUsersHandler(usersSvc, cfg)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", usersHandler.Register)
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		http_handlers.SendJSON(w, http.StatusOK, map[string]bool{"success": true})
	})

	srv := &http.Server{
		Addr:              cfg.App.HTTPPort,
		Handler:           r,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverErr := make(chan error, 1)
	go func() {
		logger.Info("http server started", slog.String("addr", cfg.App.HTTPPort))
		serverErr <- srv.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		logger.Info("shutting down...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	}
}
