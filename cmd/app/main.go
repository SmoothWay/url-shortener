package main

import (
	"os"

	"github.com/SmoothWay/url-shortener/internal/config"
	mw "github.com/SmoothWay/url-shortener/internal/http-server/middleware"
	"github.com/SmoothWay/url-shortener/internal/lib/logger/handlers/slogpretty"
	"github.com/SmoothWay/url-shortener/internal/lib/logger/sl"
	"github.com/SmoothWay/url-shortener/internal/storage/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("Starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("Debug message")

	storage, err := sqlite.New(cfg.Storage)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		return
	}

	_ = storage

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(mw.New(log))
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	// middleware

	// TODO: run server
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	default: // If env config is invalid, set prod settings by default due to security
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
