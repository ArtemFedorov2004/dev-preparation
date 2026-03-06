// Package main DevPrep API
//
// REST API для подготовки к техническому интервью.
//
//	@title			DevPrep API
//	@version		1.0
//	@description	REST API для подготовки к техническому интервью.
//
//	@BasePath	/api/v1
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/devprep/backend/docs"
	"github.com/devprep/backend/internal/config"
	"github.com/devprep/backend/internal/database"
	"github.com/devprep/backend/internal/handler"
	"github.com/devprep/backend/internal/keycloak"
	"github.com/devprep/backend/internal/middleware"
	"github.com/devprep/backend/internal/redis"
	"github.com/devprep/backend/internal/repository"
	"github.com/devprep/backend/internal/service"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	_, err = keycloak.NewJWKS(ctx, keycloak.Config{
		RealmURL:        cfg.Keycloak.RealmURL(),
		RefreshInterval: cfg.Keycloak.RefreshInterval,
	})
	if err != nil {
		slog.Error("failed to initialize keycloak jwks", "error", err)
		os.Exit(1)
	}
	slog.Info("keycloak jwks initialized", "realm_url", cfg.Keycloak.RealmURL())

	migrateDSN := "pgx5://" + cfg.Database.DSN()[len("postgres://"):]
	if err := database.RunMigrations(migrateDSN); err != nil {
		slog.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	pool, err := database.Connect(ctx, cfg.Database)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()
	slog.Info("connected to database")

	rdb, err := redis.Connect(ctx, cfg.Redis)
	if err != nil {
		slog.Warn("redis unavailable, caching disabled", "error", err)
	}

	var topicRepo repository.TopicRepository = repository.NewPGTopicRepo(pool)
	var questionRepo repository.QuestionRepository = repository.NewPGQuestionRepo(pool)

	if rdb != nil {
		topicRepo = repository.NewCachedTopicRepo(topicRepo, rdb, cfg.Redis.TopicTTL)
		questionRepo = repository.NewCachedQuestionRepo(questionRepo, rdb, cfg.Redis.QuestionTTL)
		slog.Info("redis cache enabled")
	}

	topicSvc := service.NewTopicService(topicRepo, questionRepo)
	questionSvc := service.NewQuestionService(questionRepo)

	topicHandler := handler.NewTopicHandler(topicSvc)
	questionHandler := handler.NewQuestionHandler(questionSvc)

	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.CORS)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	r.Get("/health", handler.Health)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/topics", topicHandler.ListTopics)
		r.Get("/topics/{slug}", topicHandler.GetTopicBySlug)

		r.Get("/questions", questionHandler.ListQuestions)
		r.Get("/questions/{slug}", questionHandler.GetQuestionBySlug)

		r.Get("/tags", questionHandler.ListTags)
	})

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		slog.Info("server starting", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	}

	slog.Info("server stopped")
}
