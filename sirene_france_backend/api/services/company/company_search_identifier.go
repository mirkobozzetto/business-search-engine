package company

import (
	"context"
	"fmt"
	"sirene-importer/api/models"
)

const (
	SIREN_LENGTH = 9
	SIRET_LENGTH = 14
)

func (s *companyService) SearchByIdentifier(ctx context.Context, identifier string) (*models.CompanySearchResult, error) {
	identifierLen := len(identifier)

	switch identifierLen {
	case SIRET_LENGTH:
		return s.searchBySiret(ctx, identifier)
	case SIREN_LENGTH:
		return s.searchBySiren(ctx, identifier)
	default:
		return nil, fmt.Errorf("identifier must be a 9-digit SIREN or 14-digit SIRET")
	}
}

func (s *companyService) searchBySiret(ctx context.Context, siret string) (*models.CompanySearchResult, error) {
	exists, err := s.verifySiretExists(ctx, siret)
	if err != nil {
		return nil, err
	}

	if !exists {
		return &models.CompanySearchResult{
			Criteria: models.CompanySearchCriteria{},
			Results:  []models.CompanyResult{},
			Meta:     models.Meta{Total: 0, Count: 0},
		}, nil
	}

	siren := siret[:SIREN_LENGTH]
	return s.searchBySiren(ctx, siren)
}

func (s *companyService) searchBySiren(ctx context.Context, siren string) (*models.CompanySearchResult, error) {
	exists, err := s.verifySirenExists(ctx, siren)
	if err != nil {
		return nil, err
	}

	if !exists {
		return &models.CompanySearchResult{
			Criteria: models.CompanySearchCriteria{},
			Results:  []models.CompanyResult{},
			Meta:     models.Meta{Total: 0, Count: 0},
		}, nil
	}

	companies := s.enrichCompanyData(ctx, []string{siren})

	if len(companies) == 0 {
		return &models.CompanySearchResult{
			Criteria: models.CompanySearchCriteria{},
			Results:  []models.CompanyResult{},
			Meta:     models.Meta{Total: 0, Count: 0},
		}, nil
	}

	result := &models.CompanySearchResult{
		Criteria: models.CompanySearchCriteria{},
		Results:  companies,
		Meta:     models.Meta{Total: 1, Count: 1},
	}

	return result, nil
}

func (s *companyService) verifySirenExists(ctx context.Context, siren string) (bool, error) {
	var count int
	query := `SELECT 1 FROM unite_legale WHERE siren = $1 LIMIT 1`
	err := s.db.QueryRowContext(ctx, query, siren).Scan(&count)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (s *companyService) verifySiretExists(ctx context.Context, siret string) (bool, error) {
	var count int
	query := `SELECT 1 FROM etablissement WHERE siret = $1 LIMIT 1`
	err := s.db.QueryRowContext(ctx, query, siret).Scan(&count)
	if err != nil {
		return false, nil
	}
	return true, nil
}
