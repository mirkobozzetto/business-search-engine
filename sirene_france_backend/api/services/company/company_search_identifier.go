package company

import (
	"context"
	"database/sql"
	"fmt"
	"sirene-importer/api/models"
)

const (
	SIREN_LENGTH = 9
	SIRET_LENGTH = 14
)

func (s *companyService) SearchByIdentifier(ctx context.Context, identifier string) (*models.CompanySearchResult, error) {
	switch len(identifier) {
	case SIRET_LENGTH:
		return s.lookupBySiret(ctx, identifier)
	case SIREN_LENGTH:
		return s.lookupBySiren(ctx, identifier)
	default:
		return nil, fmt.Errorf("identifier must be a 9-digit SIREN or 14-digit SIRET")
	}
}

func (s *companyService) lookupBySiren(ctx context.Context, siren string) (*models.CompanySearchResult, error) {
	query := fmt.Sprintf(`SELECT %s
		FROM etablissement e
		JOIN unite_legale u ON e.siren = u.siren
		LEFT JOIN naf_reference naf ON COALESCE(NULLIF(e.activite_principale_etablissement, ''), u.activite_principale_unite_legale, '') = naf.code
		WHERE e.etablissement_siege = 'true' AND u.siren = $1
		LIMIT 1`, companySelectFields)

	c, err := scanCompanyRow(s.db.QueryRowContext(ctx, query, siren))
	if err != nil {
		if err == sql.ErrNoRows {
			return &models.CompanySearchResult{
				Results: []models.CompanyResult{},
				Meta:    models.Meta{Total: 0, Count: 0},
			}, nil
		}
		return nil, fmt.Errorf("lookup failed: %w", err)
	}

	return &models.CompanySearchResult{
		Results: []models.CompanyResult{c},
		Meta:    models.Meta{Total: 1, Count: 1},
	}, nil
}

func (s *companyService) lookupBySiret(ctx context.Context, siret string) (*models.CompanySearchResult, error) {
	query := fmt.Sprintf(`SELECT %s
		FROM etablissement e
		JOIN unite_legale u ON e.siren = u.siren
		LEFT JOIN naf_reference naf ON COALESCE(NULLIF(e.activite_principale_etablissement, ''), u.activite_principale_unite_legale, '') = naf.code
		WHERE e.siret = $1
		LIMIT 1`, companySelectFields)

	c, err := scanCompanyRow(s.db.QueryRowContext(ctx, query, siret))
	if err != nil {
		if err == sql.ErrNoRows {
			return &models.CompanySearchResult{
				Results: []models.CompanyResult{},
				Meta:    models.Meta{Total: 0, Count: 0},
			}, nil
		}
		return nil, fmt.Errorf("lookup failed: %w", err)
	}

	return &models.CompanySearchResult{
		Results: []models.CompanyResult{c},
		Meta:    models.Meta{Total: 1, Count: 1},
	}, nil
}
