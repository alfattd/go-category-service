package server

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/alfattd/crud/internal/platform/config"
	"github.com/alfattd/crud/internal/platform/logger"
	"github.com/alfattd/crud/internal/platform/monitor"
)

func Run() (*config.Config, *http.Server) {
	cfg := config.Load()

	logger.New()

	if err := cfg.Validate(); err != nil {
		slog.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

	monitor.Init()

	srv := New(cfg)

	return cfg, srv
}
