package server

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alfattd/category-service/internal/handler"
	"github.com/alfattd/category-service/internal/platform/config"
	"github.com/alfattd/category-service/internal/platform/database"
	"github.com/alfattd/category-service/internal/platform/monitor"
	"github.com/alfattd/category-service/internal/platform/rabbitmq"
	"github.com/alfattd/category-service/internal/repository"
	"github.com/alfattd/category-service/internal/service"
)

func New(cfg *config.Config) (*http.Server, func()) {
	mux := http.NewServeMux()

	db, err := database.NewPostgres(cfg.DBUrl())
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	publisher, err := rabbitmq.NewPublisher(cfg.RabbitMQUrl, "category_events")
	if err != nil {
		slog.Error("failed to connect to rabbitmq", "error", err)
		os.Exit(1)
	}

	cleanup := func() {
		publisher.Close()
		if err := db.Close(); err != nil {
			slog.Error("failed to close database", "error", err)
		}
	}

	categoryRepo := repository.NewPostgresCategoryRepo(db)
	categoryService := service.NewCategoryService(categoryRepo, publisher)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	mux.HandleFunc("/health", monitor.Health)
	mux.HandleFunc("/version", monitor.Version(cfg.ServiceName, cfg.ServiceVersion))
	mux.Handle("/metrics", monitor.MetricsHandler())

	mux.HandleFunc("GET /categories", categoryHandler.List)
	mux.HandleFunc("POST /categories", categoryHandler.Create)
	mux.HandleFunc("GET /categories/{id}", categoryHandler.GetByID)
	mux.HandleFunc("PUT /categories/{id}", categoryHandler.Update)
	mux.HandleFunc("DELETE /categories/{id}", categoryHandler.Delete)

	handlerWithMetrics := MetricsMiddleware(mux)

	srv := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      handlerWithMetrics,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return srv, cleanup
}

var _ *sql.DB
