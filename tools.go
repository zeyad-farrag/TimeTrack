//go:build tools

package tools

import (
	_ "github.com/go-chi/chi/v5"
	_ "github.com/google/uuid"
	_ "github.com/gorilla/websocket"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/prometheus/client_golang/prometheus"
	_ "github.com/sqlc-dev/sqlc/cmd/sqlc"
	_ "github.com/teambition/rrule-go"
)
