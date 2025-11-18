package main

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"os"
)

func main() {
	cfg := config{
		addr: ":8080",
		db: dbConfig{
			dsn: "host=localhost user=root password=root dbname=ecom sslmode=disable",
		},
	}

	// Logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	slog.SetDefault(logger)

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, cfg.db.dsn)
	if err != nil {
		slog.Error("failed to connect to database")
		os.Exit(1)
	}
	defer conn.Close(ctx)

	slog.Info("connected to database")

	api := application{
		config: cfg,
		db:     conn,
	}

	handler := api.mount()
	err = api.run(handler)
	if err != nil {
		slog.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}
