package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/multica-ai/team-app/internal/boot"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(
		"org_id", "",
		"workspace_id", "",
		"actor_user_id", "",
		"outcome", "",
		"duration_ms", 0,
	)
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

	if err := http.ListenAndServe(":8080", NewRouter()); err != nil {
		slog.Error("http server stopped", "error", err)
		os.Exit(1)
	}
}
