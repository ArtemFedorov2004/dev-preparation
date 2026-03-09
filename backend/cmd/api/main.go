// Package main DevPrep API
//
// REST API для подготовки к техническому интервью.
//
//	@title			DevPrep API
//	@version		1.0
//	@description	REST API для подготовки к техническому интервью.
//
//	@BasePath	/api/v1
//
//	@securityDefinitions.oauth2.accessCode	KeycloakAuth
//	@authorizationUrl						http://KEYCLOAK_AUTH_URL
//	@tokenUrl								http://KEYCLOAK_TOKEN_URL
//	@scope.openid
//	@scope.microprofile-jwt
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/devprep/backend/docs"
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

	oidcBase := cfg.Keycloak.RealmURL() + "/protocol/openid-connect"
	docs.SwaggerInfo.SwaggerTemplate = strings.NewReplacer(
		"http://KEYCLOAK_AUTH_URL", oidcBase+"/auth",
		"http://KEYCLOAK_TOKEN_URL", oidcBase+"/token",
	).Replace(docs.SwaggerInfo.SwaggerTemplate)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	jwks, err := keycloak.NewJWKS(ctx, keycloak.Config{
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

	activityRepo := repository.NewPGUserActivityRepo(pool)

	topicSvc := service.NewTopicService(topicRepo, questionRepo)
	questionSvc := service.NewQuestionService(questionRepo)
	activitySvc := service.NewUserActivityService(activityRepo)

	topicHandler := handler.NewTopicHandler(topicSvc)
	questionHandler := handler.NewQuestionHandler(questionSvc)
	activityHandler := handler.NewUserActivityHandler(activitySvc)

	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.CORS)

	swaggerOAuthScript := `
		window.ui.initOAuth({
				clientId: "swagger-ui",
				scopes: "openid microprofile-jwt",
				usePkceWithAuthorizationCodeGrant: true,
			});`

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.AfterScript(swaggerOAuthScript),
	))

	r.Get("/health", handler.Health)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/topics", topicHandler.ListTopics)
		r.Get("/topics/{slug}", topicHandler.GetTopicBySlug)

		r.Get("/questions", questionHandler.ListQuestions)
		r.Get("/questions/{slug}", questionHandler.GetQuestionBySlug)

		r.Get("/tags", questionHandler.ListTags)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(jwks))

			r.Post("/questions/{slug}/progress", activityHandler.UpdateProgress)
			r.Get("/questions/{slug}/progress", activityHandler.GetQuestionProgress)

			r.Post("/questions/{slug}/bookmark", activityHandler.ToggleBookmark)
			r.Get("/questions/{slug}/bookmark", activityHandler.GetBookmarkStatus)

			r.Post("/questions/{slug}/view", activityHandler.RecordView)

			r.Get("/me/progress", activityHandler.GetMyProgress)
			r.Get("/me/progress/by-topic", activityHandler.GetMyProgressByTopic)
			r.Get("/me/bookmarks", activityHandler.GetMyBookmarks)
			r.Get("/me/history", activityHandler.GetMyHistory)
		})
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
