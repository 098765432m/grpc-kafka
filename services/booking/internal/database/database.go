package booking_database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Connect() (*pgxpool.Pool, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	// Create dsn string
	dsn := viper.GetString("DB_URL")
	if dsn == "" {
		zap.S().Fatal("DB_URL is not set in env")
	}

	// Connect database
	conn, err := pgxpool.New(ctx, dsn)
	if err != nil {
		zap.S().Fatal("Failed to connect to database: ", err)
	}

	zap.S().Info("Connected to database successfully")

	return conn, nil
}
