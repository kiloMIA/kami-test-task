package postgre

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func ConnectDB(logger *zap.Logger) *pgxpool.Pool {
	connURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)
	pool, err := pgxpool.New(context.Background(), connURL)
	if err != nil {
		logger.Error("Failed to connect to DB", zap.Error(err))
		return nil
	}
	logger.Info("Connected to DB")
	return pool
}
