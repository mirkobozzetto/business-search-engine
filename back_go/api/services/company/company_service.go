package company

import (
	"context"
	"csv-importer/api/cache"
	"csv-importer/api/helpers/utils"
	"csv-importer/api/models"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"
)

const MAX_COMPANIES = 100000

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

	cacheKey := fmt.Sprintf("companies:full:nace:%s", naceCode)

	start := time.Now()
	var allCompanies []models.CompanyResult
	err := s.cache.Get(cacheKey, &allCompanies)
	cacheDuration := time.Since(start)

	if err != nil {
		slog.Info("Cache miss, fetching complete data from database",
			"nace_code", naceCode,
			"cache_duration_ms", cacheDuration.Milliseconds())

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

		if len(entityNumbers) > MAX_COMPANIES {
			slog.Warn("Dataset too large, truncating",
				"nace_code", naceCode,
				"original_count", len(entityNumbers),
				"truncated_to", MAX_COMPANIES)
			entityNumbers = entityNumbers[:MAX_COMPANIES]
		}

		allCompanies, err = s.enrichCompleteCompanyData(entityNumbers, naceCode)
		if err != nil {
			return nil, err
		}

		dataSize := len(allCompanies)
		estimatedSizeMB := dataSize * 2000 / 1024 / 1024

		err = s.cache.Set(cacheKey, allCompanies, 24*time.Hour)
		if err != nil {
			slog.Error("Cache write failed",
				"nace_code", naceCode,
				"companies_count", dataSize,
				"estimated_size_mb", estimatedSizeMB,
				"error", err.Error())
		} else {
			slog.Info("Cached complete company dataset",
				"nace_code", naceCode,
				"total", len(allCompanies),
				"estimated_size_mb", estimatedSizeMB)
		}
	} else {
		slog.Info("Cache hit for complete dataset",
			"nace_code", naceCode,
			"total", len(allCompanies),
			"cache_duration_ms", cacheDuration.Milliseconds())
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

func (s *companyService) enrichCompleteCompanyData(entityNumbers []string, naceCode string) ([]models.CompanyResult, error) {
	if len(entityNumbers) == 0 {
		return []models.CompanyResult{}, nil
	}

	companyMap := make(map[string]*models.CompanyResult)

	for _, entityNumber := range entityNumbers {
		companyMap[entityNumber] = &models.CompanyResult{
			EntityNumber: entityNumber,
			NaceCode:     naceCode,
		}
	}

	slog.Info("Starting complete data enrichment", "entity_count", len(entityNumbers))

	s.enrichEnterpriseData(companyMap)
	s.enrichAllDenominations(companyMap)
	s.enrichAllAddresses(companyMap)
	s.enrichAllContacts(companyMap)
	s.enrichAllActivities(companyMap)
	s.enrichAllEstablishments(companyMap)

	var results []models.CompanyResult
	for _, entityNumber := range entityNumbers {
		if company, exists := companyMap[entityNumber]; exists {
			s.setLegacyFields(company)
			results = append(results, *company)
		}
	}

	slog.Info("Complete enrichment finished", "final_count", len(results))
	return results, nil
}

func (s *companyService) enrichEnterpriseData(companyMap map[string]*models.CompanyResult) {
	if len(companyMap) == 0 {
		return
	}

	entityNumbers := make([]string, 0, len(companyMap))
	for entityNumber := range companyMap {
		entityNumbers = append(entityNumbers, entityNumber)
	}

	batchSize := 1000
	for i := 0; i < len(entityNumbers); i += batchSize {
		end := min(i+batchSize, len(entityNumbers))
		batch := entityNumbers[i:end]

		placeholders := make([]string, len(batch))
		args := make([]any, len(batch))
		for j, entityNumber := range batch {
			placeholders[j] = fmt.Sprintf("$%d", j+1)
			args[j] = entityNumber
		}

		query := fmt.Sprintf(`
			SELECT enterprisenumber, status, juridicalform, startdate
			FROM enterprise
			WHERE enterprisenumber IN (%s)
		`, strings.Join(placeholders, ","))

		rows, err := s.db.Query(query, args...)
		if err != nil {
			slog.Warn("Failed to enrich enterprise data", "batch", i, "error", err)
			continue
		}

		for rows.Next() {
			var entityNumber, status, juridicalForm, startDate string
			if err := rows.Scan(&entityNumber, &status, &juridicalForm, &startDate); err == nil {
				if company, exists := companyMap[entityNumber]; exists {
					company.Enterprise = map[string]any{
						"status":         status,
						"juridical_form": juridicalForm,
						"start_date":     startDate,
					}
					company.Status = status
					company.StartDate = startDate
				}
			}
		}
		rows.Close()
	}

	slog.Info("Enriched enterprise data", "companies", len(companyMap))
}

func (s *companyService) enrichAllDenominations(companyMap map[string]*models.CompanyResult) {
	s.enrichTableData(companyMap, "denomination",
		"SELECT entitynumber, language, denomination FROM denomination WHERE entitynumber IN (%s)",
		func(company *models.CompanyResult, row map[string]any) {
			if company.Denominations == nil {
				company.Denominations = []map[string]any{}
			}
			company.Denominations = append(company.Denominations, row)

			if row["language"] == "2" && company.Denomination == "" {
				company.Denomination = fmt.Sprintf("%v", row["denomination"])
			}
		})
}

func (s *companyService) enrichAllAddresses(companyMap map[string]*models.CompanyResult) {
	s.enrichTableData(companyMap, "address",
		`SELECT entitynumber, typeofaddress, zipcode, municipalitynl, municipalityfr,
			streetnl, streetfr, housenumber, box, extraaddressinfo
		FROM address WHERE entitynumber IN (%s)`,
		func(company *models.CompanyResult, row map[string]any) {
			if company.Addresses == nil {
				company.Addresses = []map[string]any{}
			}
			company.Addresses = append(company.Addresses, row)

			if row["typeofaddress"] == "REGO" && company.ZipCode == "" {
				company.ZipCode = fmt.Sprintf("%v", row["zipcode"])
				if row["municipalityfr"] != nil {
					company.City = fmt.Sprintf("%v", row["municipalityfr"])
				}
				if row["streetfr"] != nil {
					company.Street = fmt.Sprintf("%v", row["streetfr"])
				}
				if row["housenumber"] != nil {
					company.HouseNumber = fmt.Sprintf("%v", row["housenumber"])
				}
			}
		})
}

func (s *companyService) enrichAllContacts(companyMap map[string]*models.CompanyResult) {
	s.enrichTableData(companyMap, "contact",
		"SELECT entitynumber, contacttype, value FROM contact WHERE entitynumber IN (%s)",
		func(company *models.CompanyResult, row map[string]any) {
			if company.Contacts == nil {
				company.Contacts = []map[string]any{}
			}
			company.Contacts = append(company.Contacts, row)

			contactType := fmt.Sprintf("%v", row["contacttype"])
			value := fmt.Sprintf("%v", row["value"])
			switch contactType {
			case "EMAIL":
				if company.Email == "" {
					company.Email = value
				}
			case "WEB":
				if company.Website == "" {
					company.Website = value
				}
			case "TEL":
				if company.Phone == "" {
					company.Phone = value
				}
			case "FAX":
				if company.Fax == "" {
					company.Fax = value
				}
			}
		})
}

func (s *companyService) enrichAllActivities(companyMap map[string]*models.CompanyResult) {
	s.enrichTableData(companyMap, "activity",
		"SELECT entitynumber, activitygroup, naceversion, nacecode, classification FROM activity WHERE entitynumber IN (%s)",
		func(company *models.CompanyResult, row map[string]any) {
			if company.Activities == nil {
				company.Activities = []map[string]any{}
			}
			company.Activities = append(company.Activities, row)
		})
}

func (s *companyService) enrichAllEstablishments(companyMap map[string]*models.CompanyResult) {
	s.enrichTableData(companyMap, "establishment",
		"SELECT establishmentnumber, enterprisenumber, startdate FROM establishment WHERE enterprisenumber IN (%s)",
		func(company *models.CompanyResult, row map[string]any) {
			if company.Establishments == nil {
				company.Establishments = []map[string]any{}
			}
			company.Establishments = append(company.Establishments, row)
		})
}

func (s *companyService) enrichTableData(companyMap map[string]*models.CompanyResult, tableName, queryTemplate string, processRow func(*models.CompanyResult, map[string]any)) {
	if len(companyMap) == 0 {
		return
	}

	entityNumbers := make([]string, 0, len(companyMap))
	for entityNumber := range companyMap {
		entityNumbers = append(entityNumbers, entityNumber)
	}

	batchSize := 1000
	totalRows := 0

	for i := 0; i < len(entityNumbers); i += batchSize {
		end := min(i+batchSize, len(entityNumbers))
		batch := entityNumbers[i:end]

		placeholders := make([]string, len(batch))
		args := make([]any, len(batch))
		for j, entityNumber := range batch {
			placeholders[j] = fmt.Sprintf("$%d", j+1)
			args[j] = entityNumber
		}

		query := fmt.Sprintf(queryTemplate, strings.Join(placeholders, ","))

		rows, err := s.db.Query(query, args...)
		if err != nil {
			slog.Warn("Failed to enrich table data", "table", tableName, "batch", i, "error", err)
			continue
		}

		columns, err := rows.Columns()
		if err != nil {
			rows.Close()
			continue
		}

		data, err := utils.ScanRowsToMaps(rows, columns)
		rows.Close()

		if err != nil {
			slog.Warn("Failed to scan rows", "table", tableName, "batch", i, "error", err)
			continue
		}

		for _, row := range data {
			if entityNumber, ok := row["entitynumber"].(string); ok {
				if company, exists := companyMap[entityNumber]; exists {
					processRow(company, row)
					totalRows++
				}
			} else if enterpriseNumber, ok := row["enterprisenumber"].(string); ok {
				if company, exists := companyMap[enterpriseNumber]; exists {
					processRow(company, row)
					totalRows++
				}
			}
		}
	}

	slog.Info("Enriched table data", "table", tableName, "total_rows", totalRows)
}

func (s *companyService) setLegacyFields(company *models.CompanyResult) {
	if company.Enterprise != nil {
		if jf, ok := company.Enterprise["juridical_form"].(string); ok {
			company.JuridicalForm = jf
		}
	}
}
