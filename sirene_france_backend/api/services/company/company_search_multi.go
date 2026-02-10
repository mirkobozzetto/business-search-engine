package company

import (
	"context"
	"fmt"
	"sirene-importer/api/models"
)

func (s *companyService) SearchMultiCriteria(ctx context.Context, criteria models.CompanySearchCriteria, limit int) (*models.CompanySearchResult, error) {
	var datasets [][]models.CompanyResult

	if criteria.NafCode != "" {
		result, err := s.SearchByNafCode(ctx, criteria.NafCode, 0)
		if err != nil {
			return nil, fmt.Errorf("naf search failed: %w", err)
		}
		datasets = append(datasets, result.Results)
	}

	if criteria.Denomination != "" {
		result, err := s.SearchByDenomination(ctx, criteria.Denomination, 0)
		if err != nil {
			return nil, fmt.Errorf("denomination search failed: %w", err)
		}
		datasets = append(datasets, result.Results)
	}

	if criteria.CodePostal != "" {
		result, err := s.SearchByCodePostal(ctx, criteria.CodePostal, 0)
		if err != nil {
			return nil, fmt.Errorf("codepostal search failed: %w", err)
		}
		datasets = append(datasets, result.Results)
	}

	if criteria.DateCreationFrom != "" {
		result, err := s.SearchByDateCreation(ctx, criteria.DateCreationFrom, criteria.DateCreationTo, 0)
		if err != nil {
			return nil, fmt.Errorf("datecreation search failed: %w", err)
		}
		datasets = append(datasets, result.Results)
	}

	if criteria.Commune != "" {
		result, err := s.SearchByCommune(ctx, criteria.Commune, 0)
		if err != nil {
			return nil, fmt.Errorf("commune search failed: %w", err)
		}
		datasets = append(datasets, result.Results)
	}

	if criteria.EtatAdministratif != "" {
		result, err := s.SearchByEtatAdministratif(ctx, criteria.EtatAdministratif, 0)
		if err != nil {
			return nil, fmt.Errorf("etatadministratif search failed: %w", err)
		}
		datasets = append(datasets, result.Results)
	}

	if len(datasets) == 0 {
		return &models.CompanySearchResult{
			Criteria: criteria,
			Results:  []models.CompanyResult{},
			Meta:     models.Meta{Total: 0, Count: 0},
		}, nil
	}

	results := intersectResults(datasets)

	searchResult := &models.CompanySearchResult{
		Criteria: criteria,
		Results:  results,
		Meta:     models.Meta{Total: len(results)},
	}

	return buildSearchResult(searchResult, limit), nil
}

func intersectResults(datasets [][]models.CompanyResult) []models.CompanyResult {
	if len(datasets) == 1 {
		return datasets[0]
	}

	smallest := 0
	for i, ds := range datasets {
		if len(ds) < len(datasets[smallest]) {
			smallest = i
		}
	}

	sirenMap := make(map[string]models.CompanyResult)
	for _, company := range datasets[smallest] {
		sirenMap[company.Siren] = company
	}

	for i, ds := range datasets {
		if i == smallest {
			continue
		}
		otherSirens := make(map[string]bool)
		for _, company := range ds {
			otherSirens[company.Siren] = true
		}
		for siren := range sirenMap {
			if !otherSirens[siren] {
				delete(sirenMap, siren)
			}
		}
	}

	results := make([]models.CompanyResult, 0, len(sirenMap))
	for _, company := range sirenMap {
		results = append(results, company)
	}
	return results
}
