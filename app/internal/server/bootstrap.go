package server

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/alfattd/category-service/internal/platform/config"
	"github.com/alfattd/category-service/internal/platform/logger"
	"github.com/alfattd/category-service/internal/platform/monitor"
)

func Build() (*config.Config, *http.Server, func()) {
	cfg := config.Load()

	logger.New()

	if err := cfg.Validate(); err != nil {
		slog.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

	monitor.Init()

	srv, cleanup := New(cfg)

	return cfg, srv, cleanup
}
