package server

import (
	"net/http"

	"github.com/alfattd/category-service/internal/handler"
	"github.com/alfattd/category-service/internal/platform/config"
	"github.com/alfattd/category-service/internal/platform/database"
	"github.com/alfattd/category-service/internal/platform/monitor"
	"github.com/alfattd/category-service/internal/platform/rabbitmq"
	"github.com/alfattd/category-service/internal/repository"
	"github.com/alfattd/category-service/internal/repository/memory"
	"github.com/alfattd/category-service/internal/repository/postgres"
	"github.com/alfattd/category-service/internal/service"
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

	publisher, err := rabbitmq.NewPublisher(cfg.RabbitMQUrl, "category_events")
	if err != nil {
		panic(err)
	}

	categoryHandler := handler.NewCategoryHandler(categoryService, publisher)

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
