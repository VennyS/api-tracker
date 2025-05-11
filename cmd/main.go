package main

import (
	"api-tracker/internal/config"
	"api-tracker/internal/handlers/loghandler"
	"api-tracker/internal/lib/logger/handlers/slogpretty"
	"api-tracker/internal/service/logservice"
	"api-tracker/internal/storage/clickhouse"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	router := chi.NewRouter()

	clickStorage, err := clickhouse.New(cfg.ClickHouse, log)
	if err != nil {
		log.Error("cannot init storage", "error", err)
		os.Exit(1)
	}
	defer clickStorage.Close()

	logSrv := logservice.New(clickStorage, log)
	logHandler := loghandler.New(logSrv, log)

	router.Post("/", logHandler.PostLog())

	log.Info("starting server", "port", 8080)
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Error("server failed", "error", err)
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
