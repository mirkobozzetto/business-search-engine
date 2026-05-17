package helpers

import (
	"csv-importer/api/helpers/utils"
	"database/sql"
	"fmt"
	"regexp"
	"slices"
)

var ValidTables = []string{
	"activity", "address", "branch", "code", "contact",
	"denomination", "enterprise", "establishment", "meta", "nacecode",
}

func ValidateTableName(tableName string) error {
	if slices.Contains(ValidTables, tableName) {
		return nil
	}
	return fmt.Errorf("invalid table name: %s", tableName)
}

func ValidateColumnExists(db *sql.DB, tableName, columnName string) error {
	var count int
	err := db.QueryRow(`
		SELECT count(*)
		FROM information_schema.columns
		WHERE table_name = $1 AND column_name = $2
	`, tableName, columnName).Scan(&count)

	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("column %s not found in table %s", columnName, tableName)
	}
	return nil
}

func ValidateIdentifier(name string) error {
	matched, _ := regexp.MatchString(`^[a-zA-Z_][a-zA-Z0-9_]*$`, name)
	if !matched {
		return fmt.Errorf("invalid identifier: %s", name)
	}
	return nil
}

func SafeQuery(db *sql.DB, tableName string, columns []string) (*sql.Rows, error) {
	if err := ValidateTableName(tableName); err != nil {
		return nil, err
	}

	for _, col := range columns {
		if err := ValidateIdentifier(col); err != nil {
			return nil, err
		}
		if err := ValidateColumnExists(db, tableName, col); err != nil {
			return nil, err
		}
	}

	query := fmt.Sprintf("SELECT %s FROM %s", utils.JoinColumns(columns), tableName)
	return db.Query(query)
}

func SafeQueryWithBuilder(db *sql.DB, tableName string, columns []string, builder *utils.QueryBuilder) (*sql.Rows, error) {
	if err := ValidateTableName(tableName); err != nil {
		return nil, err
	}

	for _, col := range columns {
		if err := ValidateIdentifier(col); err != nil {
			return nil, err
		}
		if err := ValidateColumnExists(db, tableName, col); err != nil {
			return nil, err
		}
	}

	query, args := utils.BuildSafeQuery(tableName, columns, builder)
	return utils.ExecuteSafeQuery(db, query, args)
}
