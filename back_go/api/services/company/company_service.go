package company

import (
	"context"
	"csv-importer/api/models"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type companyService struct {
	db *sql.DB
}

func NewCompanyService(db *sql.DB) CompanyService {
	if db == nil {
		slog.Error("database connection is nil")
		os.Exit(1)
	}

	return &companyService{
		db: db,
	}
}

func (s *companyService) SearchByNaceCode(ctx context.Context, naceCode string, limit int) (*models.CompanySearchResult, error) {
	if naceCode == "" {
		return nil, fmt.Errorf("nace code cannot be empty")
	}

	if limit <= 0 || limit > 1000 {
		return nil, fmt.Errorf("invalid limit: must be between 1 and 1000")
	}

	entityNumbers, err := s.getEntityNumbersByNace(naceCode, limit)
	if err != nil {
		return nil, err
	}

	if len(entityNumbers) == 0 {
		return &models.CompanySearchResult{
			Criteria: models.CompanySearchCriteria{NaceCode: naceCode},
			Results:  []models.CompanyResult{},
			Meta:     models.Meta{Count: 0, Limit: limit},
		}, nil
	}

	companies, err := s.enrichCompanyData(entityNumbers, naceCode)
	if err != nil {
		return nil, err
	}

	return &models.CompanySearchResult{
		Criteria: models.CompanySearchCriteria{NaceCode: naceCode},
		Results:  companies,
		Meta:     models.Meta{Count: len(companies), Limit: limit},
	}, nil
}

func (s *companyService) getEntityNumbersByNace(naceCode string, limit int) ([]string, error) {
	query := `
		SELECT DISTINCT entitynumber
		FROM activity
		WHERE nacecode = $1 AND classification = 'MAIN'
		ORDER BY entitynumber
		LIMIT $2
	`

	rows, err := s.db.Query(query, naceCode, limit)
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

	return entityNumbers, nil
}

func (s *companyService) enrichCompanyData(entityNumbers []string, naceCode string) ([]models.CompanyResult, error) {
	if len(entityNumbers) == 0 {
		return []models.CompanyResult{}, nil
	}

	placeholders := make([]string, len(entityNumbers))
	args := make([]any, len(entityNumbers))
	for i, entityNumber := range entityNumbers {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = entityNumber
	}

	query := fmt.Sprintf(`
		SELECT
			e.enterprisenumber,
			COALESCE(d.denomination, '') as denomination,
			COALESCE(c.description, '') as juridical_form,
			COALESCE(e.startdate::text, '') as start_date,
			COALESCE(e.status, '') as status,
			COALESCE(addr.zipcode, '') as zipcode,
			COALESCE(addr.municipalityfr, '') as city,
			COALESCE(addr.streetfr, '') as street,
			COALESCE(addr.housenumber, '') as house_number,
			COALESCE(n.libellé_fr, '') as nace_description
		FROM enterprise e
		LEFT JOIN denomination d ON e.enterprisenumber = d.entitynumber AND d.language = '2'
		LEFT JOIN code c ON e.juridicalform = c.code AND c.category = 'JuridicalForm'
		LEFT JOIN address addr ON e.enterprisenumber = addr.entitynumber AND addr.typeofaddress = 'REGO'
		LEFT JOIN nacecode n ON '%s' = n.nacecode
		WHERE e.enterprisenumber IN (%s)
		ORDER BY e.enterprisenumber
	`, naceCode, strings.Join(placeholders, ","))

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to enrich company data: %w", err)
	}
	defer rows.Close()

	companyMap := make(map[string]*models.CompanyResult)
	for rows.Next() {
		var company models.CompanyResult
		err := rows.Scan(
			&company.EntityNumber,
			&company.Denomination,
			&company.JuridicalForm,
			&company.StartDate,
			&company.Status,
			&company.ZipCode,
			&company.City,
			&company.Street,
			&company.HouseNumber,
			&company.NaceDescription,
		)
		if err != nil {
			continue
		}
		company.NaceCode = naceCode
		companyMap[company.EntityNumber] = &company
	}

	s.enrichContactData(companyMap)

	var results []models.CompanyResult
	for _, entityNumber := range entityNumbers {
		if company, exists := companyMap[entityNumber]; exists {
			results = append(results, *company)
		}
	}

	return results, nil
}

func (s *companyService) enrichContactData(companyMap map[string]*models.CompanyResult) {
	if len(companyMap) == 0 {
		return
	}

	entityNumbers := make([]string, 0, len(companyMap))
	for entityNumber := range companyMap {
		entityNumbers = append(entityNumbers, entityNumber)
	}

	placeholders := make([]string, len(entityNumbers))
	args := make([]any, len(entityNumbers))
	for i, entityNumber := range entityNumbers {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = entityNumber
	}

	query := fmt.Sprintf(`
		SELECT entitynumber, contacttype, value
		FROM contact
		WHERE entitynumber IN (%s)
		AND contacttype IN ('EMAIL', 'WEB', 'TEL', 'FAX')
	`, strings.Join(placeholders, ","))

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var entityNumber, contactType, value string
		if err := rows.Scan(&entityNumber, &contactType, &value); err != nil {
			continue
		}

		if company, exists := companyMap[entityNumber]; exists {
			switch contactType {
			case "EMAIL":
				company.Email = value
			case "WEB":
				company.Website = value
			case "TEL":
				company.Phone = value
			case "FAX":
				company.Fax = value
			}
		}
	}
}
