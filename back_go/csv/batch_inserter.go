package csv

import (
	"database/sql"

	"github.com/lib/pq"
)

func InsertBatch(db *sql.DB, tableName string, headers []string, batch [][]string) error {
	if len(batch) == 0 {
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(pq.CopyIn(tableName, headers...))
	if err != nil {
		return err
	}

	for _, record := range batch {
		values := make([]any, len(record))
		for i, v := range record {
			values[i] = v
		}

		if _, err := stmt.Exec(values...); err != nil {
			stmt.Close()
			return err
		}
	}

	if _, err := stmt.Exec(); err != nil {
		stmt.Close()
		return err
	}

	if err := stmt.Close(); err != nil {
		return err
	}

	return tx.Commit()
}
