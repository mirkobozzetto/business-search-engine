package models

type CompanySearchCriteria struct {
	Siren              string `json:"siren,omitempty"`
	Siret              string `json:"siret,omitempty"`
	NafCode            string `json:"naf_code,omitempty"`
	Denomination       string `json:"denomination,omitempty"`
	CodePostal         string `json:"code_postal,omitempty"`
	Commune            string `json:"commune,omitempty"`
	EtatAdministratif  string `json:"etat_administratif,omitempty"`
	DateCreationFrom   string `json:"date_creation_from,omitempty"`
	DateCreationTo     string `json:"date_creation_to,omitempty"`
	CategorieJuridique string `json:"categorie_juridique,omitempty"`
	TrancheEffectifs   string `json:"tranche_effectifs,omitempty"`
}

type CompanyResult struct {
	Siren               string           `json:"siren"`
	Denomination        string           `json:"denomination,omitempty"`
	Sigle               string           `json:"sigle,omitempty"`
	CategorieJuridique  string           `json:"categorie_juridique,omitempty"`
	DateCreation        string           `json:"date_creation,omitempty"`
	EtatAdministratif   string           `json:"etat_administratif,omitempty"`
	TrancheEffectifs    string           `json:"tranche_effectifs,omitempty"`
	CategorieEntreprise string           `json:"categorie_entreprise,omitempty"`
	NafCode             string           `json:"naf_code,omitempty"`
	NafLabel            string           `json:"naf_label,omitempty"`
	Siret               string           `json:"siret,omitempty"`
	Enseigne            string           `json:"enseigne,omitempty"`
	NumeroVoie          string           `json:"numero_voie,omitempty"`
	TypeVoie            string           `json:"type_voie,omitempty"`
	LibelleVoie         string           `json:"libelle_voie,omitempty"`
	CodePostal          string           `json:"code_postal,omitempty"`
	LibelleCommune      string           `json:"libelle_commune,omitempty"`
	Email               string           `json:"email,omitempty"`
	Telephone           string           `json:"telephone,omitempty"`
	Website             string           `json:"website,omitempty"`
	UniteLegale         map[string]any   `json:"unite_legale,omitempty"`
	Etablissements      []map[string]any `json:"etablissements,omitempty"`
}

type CompanySearchResult struct {
	Criteria CompanySearchCriteria `json:"criteria"`
	Results  []CompanyResult       `json:"results"`
	Meta     Meta                  `json:"meta"`
}
