package server

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/alfattd/category-service/internal/pkg/config"
	"github.com/alfattd/category-service/internal/pkg/logger"
	"github.com/alfattd/category-service/internal/pkg/monitor"
)

func Build() (*config.Config, *http.Server, func(), *slog.Logger) {
	log := logger.New()

	cfg := config.Load()

	if err := cfg.Validate(); err != nil {
		log.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

	monitor.Init()

	srv, cleanup := New(cfg, log)

	return cfg, srv, cleanup, log
}
