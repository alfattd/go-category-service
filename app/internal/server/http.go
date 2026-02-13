package server

import (
	"net/http"

	"github.com/alfattd/crud/internal/handler"
	"github.com/alfattd/crud/internal/platform/config"
	"github.com/alfattd/crud/internal/platform/database"
	"github.com/alfattd/crud/internal/platform/monitor"
	"github.com/alfattd/crud/internal/repository"
	"github.com/alfattd/crud/internal/repository/memory"
	"github.com/alfattd/crud/internal/repository/postgres"
	"github.com/alfattd/crud/internal/service"
)

func New(cfg *config.Config) *http.Server {

	mux := http.NewServeMux()

	var categoryRepo repository.CategoryRepository

	if cfg.ServiceVersion == "dev" {
		categoryRepo = memory.NewInMemoryCategoryRepo()
	} else {
		db := database.NewPostgres(cfg.DBUrl())
		categoryRepo = postgres.NewPostgresCategoryRepo(db)
	}

	categoryService := service.NewCategoryService(categoryRepo)

	categoryHandler := handler.NewCategoryHandler(categoryService)

	mux.HandleFunc("/health", monitor.Health)
	mux.HandleFunc("/version", monitor.Version(cfg.ServiceName, cfg.ServiceVersion))
	mux.Handle("/metrics", monitor.MetricsHandler())

	mux.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		case http.MethodPost:
			categoryHandler.CreateCategory(w, r)

		case http.MethodGet:
			categoryHandler.ListCategory(w, r)

		case http.MethodPut:
			categoryHandler.UpdateCategory(w, r)

		case http.MethodDelete:
			categoryHandler.DeleteCategory(w, r)

		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/categories/detail", categoryHandler.GetCategoryByID)

	handlerWithMetrics := MetricsMiddleware(mux)

	return &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: handlerWithMetrics,
	}
}
