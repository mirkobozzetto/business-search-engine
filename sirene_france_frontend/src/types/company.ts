export interface CompanyResult {
  siren: string;
  denomination?: string;
  sigle?: string;
  categorie_juridique?: string;
  date_creation?: string;
  etat_administratif?: string;
  tranche_effectifs?: string;
  categorie_entreprise?: string;
  naf_code?: string;
  naf_label?: string;
  siret?: string;
  enseigne?: string;
  numero_voie?: string;
  type_voie?: string;
  libelle_voie?: string;
  code_postal?: string;
  libelle_commune?: string;
  email?: string;
  telephone?: string;
  website?: string;
  unite_legale?: Record<string, unknown>;
  etablissements?: Record<string, unknown>[];
}

export interface CompanySearchCriteria {
  naf_code?: string;
  denomination?: string;
  code_postal?: string;
  commune?: string;
  etat_administratif?: string;
  date_creation_from?: string;
  date_creation_to?: string;
  categorie_juridique?: string;
  tranche_effectifs?: string;
}
