package company

import (
	"context"
	"csv-importer/api/cache"
	"csv-importer/api/models"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"
)

type companyService struct {
	db    *sql.DB
	cache *cache.RedisCache
}

func NewCompanyService(db *sql.DB) CompanyService {
	if db == nil {
		slog.Error("database connection is nil")
		os.Exit(1)
	}

	redisCache := cache.NewRedisCache(cache.CacheConfig{
		Host:     "localhost",
		Port:     "6379",
		Password: "",
		DB:       0,
	})

	return &companyService{
		db:    db,
		cache: redisCache,
	}
}

func (s *companyService) SearchByNaceCode(ctx context.Context, naceCode string, limit int) (*models.CompanySearchResult, error) {
	if naceCode == "" {
		return nil, fmt.Errorf("nace code cannot be empty")
	}

	if limit <= 0 {
		limit = 50
	}

	cacheKey := fmt.Sprintf("companies:nace:%s", naceCode)

	var allCompanies []models.CompanyResult
	err := s.cache.Get(cacheKey, &allCompanies)

	if err != nil {
		slog.Info("Cache miss, fetching from database", "nace_code", naceCode)

		entityNumbers, err := s.getAllEntityNumbersByNace(naceCode)
		if err != nil {
			return nil, err
		}

		if len(entityNumbers) == 0 {
			return &models.CompanySearchResult{
				Criteria: models.CompanySearchCriteria{NaceCode: naceCode},
				Results:  []models.CompanyResult{},
				Meta:     models.Meta{Count: 0, Total: 0, Limit: limit},
			}, nil
		}

		allCompanies, err = s.enrichCompanyData(entityNumbers, naceCode)
		if err != nil {
			return nil, err
		}

		err = s.cache.Set(cacheKey, allCompanies, 24*time.Hour)
		if err != nil {
			slog.Warn("Failed to cache results", "error", err)
		}

		slog.Info("Cached all companies", "nace_code", naceCode, "total", len(allCompanies))
	} else {
		slog.Info("Cache hit", "nace_code", naceCode, "total", len(allCompanies))
	}

	total := len(allCompanies)
	var results []models.CompanyResult

	if limit > total {
		results = allCompanies
	} else {
		results = allCompanies[:limit]
	}

	return &models.CompanySearchResult{
		Criteria: models.CompanySearchCriteria{NaceCode: naceCode},
		Results:  results,
		Meta:     models.Meta{Count: len(results), Total: total, Limit: limit},
	}, nil
}

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

func (s *companyService) enrichCompanyData(entityNumbers []string, naceCode string) ([]models.CompanyResult, error) {
	if len(entityNumbers) == 0 {
		return []models.CompanyResult{}, nil
	}

	batchSize := 1000
	var allCompanies []models.CompanyResult

	for i := 0; i < len(entityNumbers); i += batchSize {
		end := min(i+batchSize, len(entityNumbers))

		batch := entityNumbers[i:end]
		companies, err := s.enrichBatch(batch, naceCode)
		if err != nil {
			slog.Warn("Failed to enrich batch", "start", i, "end", end, "error", err)
			continue
		}

		allCompanies = append(allCompanies, companies...)
		slog.Info("Enriched batch", "start", i, "end", end, "total_so_far", len(allCompanies))
	}

	s.enrichContactDataForAll(allCompanies)

	return allCompanies, nil
}

func (s *companyService) enrichBatch(entityNumbers []string, naceCode string) ([]models.CompanyResult, error) {
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

	var companies []models.CompanyResult
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
		companies = append(companies, company)
	}

	return companies, nil
}

func (s *companyService) enrichContactDataForAll(companies []models.CompanyResult) {
	if len(companies) == 0 {
		return
	}

	entityNumbers := make([]string, len(companies))
	companyMap := make(map[string]*models.CompanyResult)

	for i := range companies {
		entityNumbers[i] = companies[i].EntityNumber
		companyMap[companies[i].EntityNumber] = &companies[i]
	}

	batchSize := 1000
	for i := 0; i < len(entityNumbers); i += batchSize {
		end := min(i+batchSize, len(entityNumbers))

		batch := entityNumbers[i:end]
		s.enrichContactBatch(batch, companyMap)
	}
}

func (s *companyService) enrichContactBatch(entityNumbers []string, companyMap map[string]*models.CompanyResult) {
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
