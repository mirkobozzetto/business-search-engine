# Rapport SIRENE France Backend

**Date:** 2026-02-10
**Status:** Import en cours, API fonctionnelle

---

## 1. Architecture

### Structure du projet

```
sirene_france_backend/
  main.go                 -> Point d'entree, connexion DB
  config/config.go        -> Variables d'environnement
  database/               -> Connexion PostgreSQL (pgx v5)
  csv/
    zip_reader.go          -> Lecture ZIP streaming (sans extraction)
    pipeline.go            -> Import parallele (workers + channels)
    batch_inserter.go      -> pgx CopyFrom (batch 200k lignes)
    utils.go               -> Conversion camelCase -> snake_case
  cli/                     -> Commandes CLI (api, all, tables, help)
  api/
    server.go              -> Serveur Gin, routes, CORS
    cache/                 -> Redis cache avec compression gzip
    models/                -> Structures de donnees
    services/company/      -> Logique metier (recherche, enrichissement)
```

### Stack technique

| Composant   | Technologie       | Port |
|-------------|-------------------|------|
| API         | Go 1.24 + Gin     | 8081 |
| Base        | PostgreSQL 15     | 5434 |
| Cache       | Redis 7           | 6380 |
| Import      | pgx CopyFrom      | -    |

### Docker

```yaml
# Conteneurs isoles de BCE Belgium (ports 5433/6379/8080)
sirene_postgres  -> port 5434
sirene_redis     -> port 6380
```

---

## 2. Donnees SIRENE

### Source

- **Origine:** data.gouv.fr (INSEE)
- **Mise a jour:** Mensuelle
- **Format:** CSV dans ZIP (UTF-8)

### Volumes

| Fichier                        | Lignes     | Colonnes | Taille ZIP |
|-------------------------------|-----------|----------|------------|
| StockUniteLegale_utf8.zip     | 29 216 651 | 35       | 897 Mo     |
| StockEtablissement_utf8.zip   | 42 716 292 | 54       | 2.6 Go     |
| **Total**                     | **71 932 943** | -    | **3.5 Go** |

### Tables PostgreSQL

**unite_legale** (29M lignes) - Entites juridiques

Colonnes cles: `siren`, `denomination_unite_legale`, `sigle_unite_legale`,
`categorie_juridique_unite_legale`, `date_creation_unite_legale`,
`etat_administratif_unite_legale`, `tranche_effectifs_unite_legale`,
`categorie_entreprise`, `activite_principale_unite_legale`

**etablissement** (42M lignes) - Etablissements physiques

Colonnes cles: `siren`, `siret`, `etablissement_siege`,
`activite_principale_etablissement`, `code_postal_etablissement`,
`libelle_commune_etablissement`, `numero_voie_etablissement`,
`type_voie_etablissement`, `libelle_voie_etablissement`,
`enseigne1_etablissement`

### Conversion des colonnes

Les CSV SIRENE utilisent le camelCase (`denominationUniteLegale`).
Le pipeline convertit automatiquement en snake_case (`denomination_unite_legale`)
via regex `([a-z0-9])([A-Z])` -> `\1_\2`.

---

## 3. Import Pipeline

### Mecanisme ZIP Streaming

L'import lit les CSV directement depuis les ZIP sans extraction sur disque :
1. `zip.OpenReader()` ouvre le ZIP
2. Premiere passe: lit les headers, cree la table PostgreSQL
3. Seconde passe: `f.Open()` retourne un `io.ReadCloser` decompresse
4. Le reader est passe a `ProcessPipelineFromReader()`

### Pipeline parallele

```
ZIP -> CSV Reader -> Channel (10k buffer) -> 8 Workers -> pgx CopyFrom (batch 200k)
```

| Parametre        | Valeur    |
|-----------------|-----------|
| Workers         | min(CPU, 8) |
| Buffer channel  | 10 000    |
| Batch size      | 200 000   |
| Buffer CSV      | 4 Mo      |
| LazyQuotes      | true      |

### Performance observee

Import de `StockEtablissement_utf8.zip` (42.7M lignes):
- ~30M lignes en ~10 minutes
- ~50 000 lignes/seconde
- Utilisation CPU: ~270% (multi-core)
- RAM: ~1.4 Go pic

---

## 4. API Endpoints

### Routes disponibles

| Methode | Route | Parametre | Description |
|---------|-------|-----------|-------------|
| GET | `/api/health` | - | Health check |
| GET | `/api/companies/search/naf` | `?code=62.01Z&limit=100` | Recherche par code NAF |
| GET | `/api/companies/search/denomination` | `?q=google&limit=100` | Recherche par nom (multi-mots) |
| GET | `/api/companies/search/codepostal` | `?q=75001&limit=100` | Recherche par code postal |
| GET | `/api/companies/search/commune` | `?q=paris&limit=100` | Recherche par commune |
| GET | `/api/companies/search/etatadministratif` | `?q=A&limit=100` | Recherche par etat (A=actif, C=cesse) |
| GET | `/api/companies/search/datecreation` | `?from=2020-01-01&to=2024-12-31&limit=100` | Recherche par date de creation |
| GET | `/api/companies/search/multi` | `?naf=62.01Z&commune=paris&etat=A&limit=100` | Recherche multi-criteres |

### Recherche multi-mots (denomination)

La recherche par denomination supporte le multi-mots avec logique AND :
- `?q=google france` -> `denomination ILIKE '%google%' AND denomination ILIKE '%france%'`
- Equivalent du `buildMultiWordSearch()` de BCE Belgium

### Recherche multi-criteres

Fonctionne par intersection de resultats caches :
1. Chaque critere lance une recherche individuelle (avec cache Redis)
2. Les resultats sont croises (intersection sur SIREN)
3. Le plus petit dataset est utilise comme base de filtrage

### Enrichissement

Pour chaque SIREN trouve, enrichissement depuis 2 tables :

**unite_legale:** denomination, sigle, categorie juridique, date creation,
etat administratif, tranche effectifs, categorie entreprise, code NAF

**etablissement (siege):** SIRET, enseigne, adresse complete
(numero, type voie, libelle voie, code postal, commune)

Batch de 1000 SIRENs par requete avec placeholders `$1, $2, ..., $1000`.

### Cache Redis

| Parametre | Valeur |
|-----------|--------|
| TTL | 24 heures |
| Compression | gzip (seuil: 1 Ko) |
| Max decompresse | 200 Mo |
| Max compresse | 50 Mo |
| Prefixe cle | `sirene:full:{type}:{valeur}` |

---

## 5. Comparaison BCE Belgium vs SIRENE France

### Donnees source

| Aspect | BCE Belgium | SIRENE France |
|--------|------------|---------------|
| Source | BCE (SPF Economie) | INSEE (data.gouv.fr) |
| Tables | 10 (enterprise, denomination, address, contact, activity, establishment, branch, code, nacecode, status) | 2 (unite_legale, etablissement) |
| Identifiant | EntityNumber | SIREN / SIRET |
| Code activite | NACE-BEL (5 chiffres) | NAF rev2 (XX.XXZ) |
| Lignes totales | ~15M | ~72M |
| Contacts | Email, telephone, site web, fax | Aucun |
| Dirigeants | Non | Non (disponible via INPI) |
| Donnees financieres | Non | Non |

### Endpoints de recherche

| Fonctionnalite | BCE Belgium | SIRENE France |
|---------------|------------|---------------|
| Recherche par code activite | `/search/nace?code=` | `/search/naf?code=` |
| Recherche par nom | `/search/denomination?q=` | `/search/denomination?q=` |
| Recherche par code postal | `/search/zipcode?q=` | `/search/codepostal?q=` |
| Recherche par commune | Non | `/search/commune?q=` |
| Recherche par date creation | `/search/startdate?from=&to=` | `/search/datecreation?from=&to=` |
| Recherche par statut | Non (via multi) | `/search/etatadministratif?q=` |
| Multi-criteres | `/search/multi` | `/search/multi` |
| Multi-mots (nom) | Oui (AND logic) | Oui (AND logic) |
| Recherche generique | `/search/:table/:column` | Non |
| Preview table | `/data/:table/preview` | Non |
| Valeurs uniques | `/data/:table/values/:column` | Non |
| Export CSV | `/export/:table` | Non |
| Info tables | `/tables`, `/tables/structure` | CLI `tables` uniquement |

### Enrichissement

| Donnee | BCE Belgium | SIRENE France |
|--------|------------|---------------|
| Nom entreprise | denomination (multilingue FR/NL) | denomination_unite_legale |
| Adresse | address (REGO) | etablissement (siege) |
| Email | contact (EMAIL) | Non disponible |
| Telephone | contact (TEL) | Non disponible |
| Site web | contact (WEB) | Non disponible |
| Activites | activity (NACE, classification) | activite_principale |
| Etablissements | establishment | etablissement (siege uniquement) |
| Forme juridique | enterprise.juridicalform | categorie_juridique |
| Effectifs | Non | tranche_effectifs |
| Categorie | Non | categorie_entreprise (PME/ETI/GE) |

### Middleware

| Middleware | BCE Belgium | SIRENE France |
|-----------|------------|---------------|
| ValidateTableName | Oui | Non |
| ValidateColumnName | Oui | Non |
| ValidateSearchQuery | Oui | Non |
| ParseLimitParam (avec max) | Oui | Oui (max 10000) |
| ParseOffsetParam | Oui | Non |
| ParseFormatParam | Oui | Non |
| ParseSortParam | Oui | Non |
| ResponseMiddleware | Oui | Non |
| CORS | Oui | Oui |

### Recherche semantique

| Feature | BCE Belgium | SIRENE France |
|---------|------------|---------------|
| ILIKE basique | Oui | Oui |
| Multi-mots AND | Oui (buildMultiWordSearch) | Oui (denomination) |
| Full-text search (tsvector) | Non | Non |
| Trigrams (pg_trgm) | Non | Non |
| Fuzzy search | Non | Non |
| Scoring/ranking | Non | Non |

**Ni BCE ni SIRENE n'implementent de recherche semantique avancee.**
Les deux utilisent `ILIKE '%terme%'` comme base, avec du multi-mots AND pour les denominations.

---

## 6. Points forts SIRENE France

1. **ZIP streaming** - Lecture directe depuis ZIP sans extraction (economie disque)
2. **Plus de criteres de recherche** - commune, etat administratif (non disponibles dans BCE)
3. **Donnees d'effectifs** - tranche_effectifs et categorie_entreprise (PME/ETI/GE)
4. **Dates ISO** - Format YYYY-MM-DD natif (pas de conversion TO_DATE comme BCE)
5. **Volume superieur** - 72M lignes vs ~15M pour BCE

## 7. Points faibles / Ameliorations possibles

### Manquant par rapport a BCE

1. **Pas d'endpoints generiques** (`/search/:table/:column`) - BCE permet de chercher dans n'importe quelle table/colonne
2. **Pas de preview de donnees** (`/data/:table/preview`) - utile pour le debug et l'exploration
3. **Pas d'export CSV** - BCE peut exporter les resultats en CSV
4. **Pas de middleware de validation** - les noms de tables/colonnes ne sont pas valides
5. **Pas de pagination** (offset) - seulement limit
6. **Pas de contacts** - les donnees SIRENE ne contiennent pas d'emails/telephones (enrichissement externe necessaire via INPI ou Apollo.io)

### Ameliorations recommandees

1. **Index PostgreSQL** - Ajouter des index sur les colonnes de recherche apres import :
   ```sql
   CREATE INDEX idx_etab_naf ON etablissement(activite_principale_etablissement);
   CREATE INDEX idx_etab_cp ON etablissement(code_postal_etablissement);
   CREATE INDEX idx_etab_commune ON etablissement(libelle_commune_etablissement);
   CREATE INDEX idx_etab_siege ON etablissement(etablissement_siege);
   CREATE INDEX idx_etab_siren ON etablissement(siren);
   CREATE INDEX idx_ul_siren ON unite_legale(siren);
   CREATE INDEX idx_ul_denom ON unite_legale(denomination_unite_legale);
   CREATE INDEX idx_ul_etat ON unite_legale(etat_administratif_unite_legale);
   CREATE INDEX idx_ul_date ON unite_legale(date_creation_unite_legale);
   ```

2. **Full-text search** - Ajouter `tsvector` sur denomination pour recherche plus performante sur 29M lignes

3. **Enrichissement INPI** - Integrer les donnees de dirigeants depuis data.inpi.fr (SFTP gratuit)

4. **API contacts payante** - Apollo.io pour email/telephone (49-119$/mois)

---

## 8. Commandes

```bash
# Demarrer Docker
make sirene-up

# Importer les donnees (depuis sirene_data/)
make sirene-import

# Lancer l'API
make sirene-api

# Arreter Docker
make sirene-down

# Lister les tables
cd sirene_france_backend && go run . tables
```

---

## 9. Corrections appliquees

| Correction | Fichier | Impact |
|-----------|---------|--------|
| Fix imports `internal/config` -> `config` | `api/server.go` | Bloquant (ne compilait pas) |
| `rows.Err()` apres boucles de lecture | 6 fichiers search + enrichment | Fiabilite |
| Limite max 10 000 sur parametre `limit` | `company_handler.go` | Securite (DoS) |
| Redis configurable via env vars | `company_service.go` | Deploiement |
| Logging des erreurs de scan | `company_enrichment.go` | Debug |
| Buffer channel 100k -> 10k | `csv/pipeline.go` | Memoire (RAM) |
| Recherche multi-mots denomination | `company_search_denomination.go` | Fonctionnalite |
| Ajout recherche par commune | `company_search_commune.go` | Fonctionnalite |
| Ajout recherche par etat administratif | `company_search_etatadmin.go` | Fonctionnalite |
| Multi-criteres: commune + etat | `company_search_multi.go` | Fonctionnalite |
