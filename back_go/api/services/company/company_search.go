package company

import (
	"fmt"
	"log/slog"
)

func (s *companyService) getAllEntityNumbersByNace(naceCode string) ([]string, error) {
	query := `
		SELECT DISTINCT entitynumber
		FROM activity
		WHERE nacecode = $1 AND classification = 'MAIN'
		ORDER BY entitynumber
	`

	rows, err := s.db.Query(query, naceCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get entity numbers: %w", err)
	}
	defer rows.Close()

	var entityNumbers []string
	for rows.Next() {
		var entityNumber string
		if err := rows.Scan(&entityNumber); err == nil {
			entityNumbers = append(entityNumbers, entityNumber)
		}
	}

	slog.Info("Found entity numbers", "nace_code", naceCode, "count", len(entityNumbers))
	return entityNumbers, nil
}

func (s *companyService) getAllEntityNumbersByDenomination(query string) ([]string, error) {
	querySQL := `
		SELECT DISTINCT entitynumber
		FROM denomination
		WHERE denomination ILIKE $1
		ORDER BY entitynumber
	`

	rows, err := s.db.Query(querySQL, "%"+query+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to get entity numbers by denomination: %w", err)
	}
	defer rows.Close()

	var entityNumbers []string
	for rows.Next() {
		var entityNumber string
		if err := rows.Scan(&entityNumber); err == nil {
			entityNumbers = append(entityNumbers, entityNumber)
		}
	}

	slog.Info("Found entity numbers by denomination", "query", query, "count", len(entityNumbers))
	return entityNumbers, nil
}

func (s *companyService) getAllEntityNumbersByZipcode(zipcode string) ([]string, error) {
	query := `
		SELECT DISTINCT entitynumber
		FROM address
		WHERE zipcode = $1 AND typeofaddress = 'REGO'
		ORDER BY entitynumber
	`

	rows, err := s.db.Query(query, zipcode)
	if err != nil {
		return nil, fmt.Errorf("failed to get entity numbers by zipcode: %w", err)
	}
	defer rows.Close()

	var entityNumbers []string
	for rows.Next() {
		var entityNumber string
		if err := rows.Scan(&entityNumber); err == nil {
			entityNumbers = append(entityNumbers, entityNumber)
		}
	}

	slog.Info("Found entity numbers by zipcode", "zipcode", zipcode, "count", len(entityNumbers))
	return entityNumbers, nil
}
