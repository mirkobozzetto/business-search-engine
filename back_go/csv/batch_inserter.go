package csv

import (
	"context"
	"csv-importer/config"
	"csv-importer/database"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func InsertBatch(tableName string, headers []string, batch [][]string) error {
	if len(batch) == 0 {
		return nil
	}

	cfg := config.Load()
	conn, err := database.ConnectPgxNative(cfg)
	if err != nil {
		return fmt.Errorf("failed to connect for batch insert: %w", err)
	}
	defer conn.Close(context.Background())

	rows := make([][]any, len(batch))
	for i, record := range batch {
		row := make([]any, len(record))
		for j, v := range record {
			row[j] = v
		}
		rows[i] = row
	}

	rowsAffected, err := conn.CopyFrom(
		context.Background(),
		pgx.Identifier{tableName},
		headers,
		pgx.CopyFromSlice(len(rows), func(i int) ([]any, error) {
			return rows[i], nil
		}),
	)

	if err != nil {
		return fmt.Errorf("copy from failed: %w", err)
	}

	if rowsAffected != int64(len(batch)) {
		return fmt.Errorf("expected %d rows, got %d", len(batch), rowsAffected)
	}

	return nil
}
