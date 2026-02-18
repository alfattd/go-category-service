package server

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alfattd/category-service/internal/handler"
	"github.com/alfattd/category-service/internal/pkg/config"
	"github.com/alfattd/category-service/internal/pkg/database"
	"github.com/alfattd/category-service/internal/pkg/monitor"
	"github.com/alfattd/category-service/internal/pkg/rabbitmq"
	"github.com/alfattd/category-service/internal/pkg/requestid"
	"github.com/alfattd/category-service/internal/repository"
	"github.com/alfattd/category-service/internal/service"
)

func New(cfg *config.Config, log *slog.Logger) (*http.Server, func()) {
	mux := http.NewServeMux()

	db, err := database.NewPostgres(cfg.DBUrl())
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	publisher, err := rabbitmq.NewPublisher(cfg.RabbitMQUrl, "category_events")
	if err != nil {
		log.Error("failed to connect to rabbitmq", "error", err)
		os.Exit(1)
	}

	cleanup := func() {
		publisher.Close()
		if err := db.Close(); err != nil {
			log.Error("failed to close database", "error", err)
		}
	}

	categoryRepo := repository.NewPostgresCategoryRepo(db)
	categoryService := service.NewCategoryService(categoryRepo, publisher, log)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	mux.HandleFunc("/health", monitor.Health)
	mux.HandleFunc("/version", monitor.Version(cfg.ServiceName, cfg.ServiceVersion))
	mux.Handle("/metrics", monitor.MetricsHandler())

	mux.HandleFunc("GET /categories", categoryHandler.List)
	mux.HandleFunc("POST /categories", categoryHandler.Create)
	mux.HandleFunc("GET /categories/{id}", categoryHandler.GetByID)
	mux.HandleFunc("PUT /categories/{id}", categoryHandler.Update)
	mux.HandleFunc("DELETE /categories/{id}", categoryHandler.Delete)

	handler := MetricsMiddleware(
		LoggingMiddleware(
			requestid.Middleware(mux),
		),
	)

	srv := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return srv, cleanup
}
