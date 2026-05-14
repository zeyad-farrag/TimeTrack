package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/multica-ai/team-app/internal/boot"
)

func main() {
	// Lean base logger. Per-request fields (org_id, workspace_id, actor_user_id,
	// outcome, duration_ms) are attached by request middleware via logger.With(...)
	// once the values are actually known — pre-allocating zero/empty values here
	// would emit them on every log line and drown out real signal.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := boot.ValidateRequiredEnv(); err != nil {
		slog.Error("missing required env var", "missing_env_var", boot.MissingEnvVar(err), "error", err)
		os.Exit(1)
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		slog.Error("connect postgres", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	server := &http.Server{
		Addr:    ":8080",
		Handler: NewRouter(),
		// TIME-RULE: keep direct time.Duration constants in boot wiring only; inject clocks in service code.
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("http server stopped", "error", err)
		os.Exit(1)
	}
}
