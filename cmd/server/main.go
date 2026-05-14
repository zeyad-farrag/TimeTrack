package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/multica-ai/team-app/internal/boot"
)

// shutdownTimeout caps how long the server will wait for in-flight
// requests to drain on SIGTERM/SIGINT before forcibly closing.
const shutdownTimeout = 30 * time.Second

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

	server := &http.Server{
		Addr:    ":8080",
		Handler: NewRouter(),
		// TIME-RULE: keep direct time.Duration constants in boot wiring only; inject clocks in service code.
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	// Run ListenAndServe in a goroutine so we can react to SIGTERM/SIGINT
	// on the main goroutine. listenErr surfaces non-shutdown errors.
	listenErr := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			listenErr <- err
			return
		}
		listenErr <- nil
	}()

	// Block until either ListenAndServe returns an error or a termination
	// signal arrives.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	exitCode := 0
	select {
	case err := <-listenErr:
		if err != nil {
			slog.Error("http server stopped", "error", err)
			exitCode = 1
		}
	case sig := <-stop:
		slog.Info("shutdown signal received", "signal", sig.String())
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			slog.Error("http server shutdown", "error", err)
			exitCode = 1
		}
		// Drain ListenAndServe's exit so we don't leak the goroutine.
		if err := <-listenErr; err != nil {
			slog.Error("http server post-shutdown", "error", err)
			exitCode = 1
		}
	}

	// Close the DB pool AFTER the server has drained, so in-flight handlers
	// finish their queries before connections are torn down.
	pool.Close()

	if exitCode != 0 {
		os.Exit(exitCode)
	}
}
