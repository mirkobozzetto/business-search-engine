package company

import (
	"context"
	"fmt"
	"log/slog"
	"sirene-importer/api/models"
	"strings"
	"sync"
	"time"
)

const companySelectFields = `
	COALESCE(u.siren, ''),
	COALESCE(u.denomination_unite_legale, ''),
	COALESCE(u.sigle_unite_legale, ''),
	COALESCE(u.categorie_juridique_unite_legale, ''),
	COALESCE(u.date_creation_unite_legale, ''),
	COALESCE(u.etat_administratif_unite_legale, ''),
	COALESCE(u.tranche_effectifs_unite_legale, ''),
	COALESCE(u.categorie_entreprise, ''),
	COALESCE(e.siret, ''),
	COALESCE(e.enseigne1_etablissement, ''),
	COALESCE(e.numero_voie_etablissement, ''),
	COALESCE(e.type_voie_etablissement, ''),
	COALESCE(e.libelle_voie_etablissement, ''),
	COALESCE(e.code_postal_etablissement, ''),
	COALESCE(e.libelle_commune_etablissement, ''),
	COALESCE(NULLIF(e.activite_principale_etablissement, ''), u.activite_principale_unite_legale, '')`

const companySelectFieldsNoCount = `
	COALESCE(u.siren, ''),
	COALESCE(u.denomination_unite_legale, ''),
	COALESCE(u.sigle_unite_legale, ''),
	COALESCE(u.categorie_juridique_unite_legale, ''),
	COALESCE(u.date_creation_unite_legale, ''),
	COALESCE(u.etat_administratif_unite_legale, ''),
	COALESCE(u.tranche_effectifs_unite_legale, ''),
	COALESCE(u.categorie_entreprise, ''),
	COALESCE(e.siret, ''),
	COALESCE(e.enseigne1_etablissement, ''),
	COALESCE(e.numero_voie_etablissement, ''),
	COALESCE(e.type_voie_etablissement, ''),
	COALESCE(e.libelle_voie_etablissement, ''),
	COALESCE(e.code_postal_etablissement, ''),
	COALESCE(e.libelle_commune_etablissement, ''),
	COALESCE(NULLIF(e.activite_principale_etablissement, ''), u.activite_principale_unite_legale, '')`

func scanCompanyRow(scanner interface{ Scan(...any) error }) (models.CompanyResult, error) {
	var c models.CompanyResult
	err := scanner.Scan(
		&c.Siren, &c.Denomination, &c.Sigle, &c.CategorieJuridique,
		&c.DateCreation, &c.EtatAdministratif, &c.TrancheEffectifs,
		&c.CategorieEntreprise,
		&c.Siret, &c.Enseigne, &c.NumeroVoie, &c.TypeVoie,
		&c.LibelleVoie, &c.CodePostal, &c.LibelleCommune,
		&c.NafCode,
	)
	return c, err
}

func (s *companyService) searchCompanies(ctx context.Context, conditions []string, args []any, limit, offset int, cacheKey string, criteria models.CompanySearchCriteria) (*models.CompanySearchResult, error) {
	pageCacheKey := fmt.Sprintf("%s:l%d:o%d", cacheKey, limit, offset)
	var cached models.CompanySearchResult
	if err := s.cache.Get(pageCacheKey, &cached); err == nil {
		return &cached, nil
	}

	where := "WHERE " + strings.Join(conditions, " AND ")
	argN := len(args) + 1

	dataQuery := fmt.Sprintf(`SELECT %s
		FROM etablissement e
		JOIN unite_legale u ON e.siren = u.siren
		%s
		ORDER BY u.date_creation_unite_legale DESC
		LIMIT $%d OFFSET $%d`, companySelectFields, where, argN, argN+1)

	dataArgs := make([]any, len(args)+2)
	copy(dataArgs, args)
	dataArgs[len(args)] = limit
	dataArgs[len(args)+1] = offset

	countCacheKey := cacheKey + ":count"
	var totalCount int
	var countCached bool

	var cachedCount int
	if err := s.cache.Get(countCacheKey, &cachedCount); err == nil {
		totalCount = cachedCount
		countCached = true
	}

	var wg sync.WaitGroup
	var countErr error

	if !countCached {
		countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM etablissement e
			JOIN unite_legale u ON e.siren = u.siren %s`, where)
		countArgs := make([]any, len(args))
		copy(countArgs, args)

		wg.Add(1)
		go func() {
			defer wg.Done()
			countErr = s.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&totalCount)
			if countErr == nil {
				_ = s.cache.Set(countCacheKey, totalCount, 1*time.Hour)
			}
		}()
	}

	rows, err := s.db.QueryContext(ctx, dataQuery, dataArgs...)
	if err != nil {
		return nil, fmt.Errorf("search query failed: %w", err)
	}
	defer func() { _ = rows.Close() }()

	companies := make([]models.CompanyResult, 0, limit)
	for rows.Next() {
		c, err := scanCompanyRow(rows)
		if err != nil {
			continue
		}
		companies = append(companies, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	if !countCached {
		wg.Wait()
		if countErr != nil {
			slog.Warn("Count query failed, using result length", "error", countErr, "key", cacheKey)
			totalCount = len(companies)
		}
	}

	page := 1
	if limit > 0 {
		page = (offset / limit) + 1
	}
	pages := 0
	if limit > 0 && totalCount > 0 {
		pages = (totalCount + limit - 1) / limit
	}

	result := &models.CompanySearchResult{
		Criteria: criteria,
		Results:  companies,
		Meta: models.Meta{
			Total:  totalCount,
			Count:  len(companies),
			Limit:  limit,
			Offset: offset,
			Page:   page,
			Pages:  pages,
		},
	}

	_ = s.cache.Set(pageCacheKey, result, 1*time.Hour)
	return result, nil
}
