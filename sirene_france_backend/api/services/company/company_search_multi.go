package company

import (
	"context"
	"fmt"
	"sirene-importer/api/models"
)

func (s *companyService) SearchMultiCriteria(ctx context.Context, criteria models.CompanySearchCriteria, limit int, offset int) (*models.CompanySearchResult, error) {
	conditions := []string{"e.etablissement_siege = 'true'"}
	var args []any
	argN := 1

	if criteria.Siren != "" {
		conditions = append(conditions, fmt.Sprintf("u.siren = $%d", argN))
		args = append(args, criteria.Siren)
		argN++
	}

	if criteria.Siret != "" {
		conditions = append(conditions, fmt.Sprintf("e.siret = $%d", argN))
		args = append(args, criteria.Siret)
		argN++
	}

	if criteria.NafCode != "" {
		conditions = append(conditions, fmt.Sprintf("e.activite_principale_etablissement = $%d", argN))
		args = append(args, criteria.NafCode)
		argN++
	}

	if criteria.Denomination != "" {
		conditions = append(conditions, fmt.Sprintf("unaccent(u.denomination_unite_legale) ILIKE unaccent($%d)", argN))
		args = append(args, "%"+criteria.Denomination+"%")
		argN++
	}

	if criteria.CodePostal != "" {
		conditions = append(conditions, fmt.Sprintf("e.code_postal_etablissement = $%d", argN))
		args = append(args, criteria.CodePostal)
		argN++
	}

	if criteria.Commune != "" {
		conditions = append(conditions, fmt.Sprintf("unaccent(e.libelle_commune_etablissement) ILIKE unaccent($%d)", argN))
		args = append(args, "%"+criteria.Commune+"%")
		argN++
	}

	if criteria.EtatAdministratif != "" {
		conditions = append(conditions, fmt.Sprintf("u.etat_administratif_unite_legale = $%d", argN))
		args = append(args, criteria.EtatAdministratif)
		argN++
	}

	if criteria.DateCreationFrom != "" {
		conditions = append(conditions, fmt.Sprintf("u.date_creation_unite_legale >= $%d", argN))
		args = append(args, criteria.DateCreationFrom)
		argN++
	}

	if criteria.DateCreationTo != "" {
		conditions = append(conditions, fmt.Sprintf("u.date_creation_unite_legale <= $%d", argN))
		args = append(args, criteria.DateCreationTo)
		argN++
	}

	if criteria.CategorieJuridique != "" {
		conditions = append(conditions, fmt.Sprintf("u.categorie_juridique_unite_legale = $%d", argN))
		args = append(args, criteria.CategorieJuridique)
		argN++
	}

	if criteria.TrancheEffectifs != "" {
		conditions = append(conditions, fmt.Sprintf("u.tranche_effectifs_unite_legale = $%d", argN))
		args = append(args, criteria.TrancheEffectifs)
	}

	if len(conditions) == 1 {
		return &models.CompanySearchResult{
			Criteria: criteria,
			Results:  []models.CompanyResult{},
			Meta:     models.Meta{Total: 0, Count: 0},
		}, nil
	}

	cacheKey := fmt.Sprintf("sirene:v2:multi:%s:%s:%s:%s:%s:%s:%s:%s:%s:%s:%s",
		criteria.Siren, criteria.Siret,
		criteria.NafCode, criteria.Denomination, criteria.CodePostal, criteria.Commune,
		criteria.EtatAdministratif, criteria.DateCreationFrom, criteria.DateCreationTo,
		criteria.CategorieJuridique, criteria.TrancheEffectifs)

	return s.searchCompanies(ctx, conditions, args, limit, offset, cacheKey, criteria)
}
