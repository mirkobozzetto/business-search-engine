package database

import (
	"context"
	"csv-importer/config"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func ConnectPgxNative(cfg *config.Config) (*pgx.Conn, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
	cfg.DBUser,
	cfg.DBPassword,
	cfg.DBHost,
	cfg.DBPort,
	cfg.DBName,
	)

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect with pgx native: %w", err)
	}

	return conn, nil
}

func CopyFromSlice(conn *pgx.Conn, tableName string, columnNames []string, rows [][]any) (int64, error) {
	return conn.CopyFrom(
		context.Background(),
		pgx.Identifier{tableName},
		columnNames,
		pgx.CopyFromSlice(len(rows), func(i int) ([]any, error) {
			return rows[i], nil
		}),
	)
}
