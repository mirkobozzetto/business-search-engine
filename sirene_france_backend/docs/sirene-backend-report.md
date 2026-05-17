# SIRENE France Backend - Rapport de verification

## Resume

Backend Go pour l'exploitation de la base SIRENE (INSEE) des entreprises francaises.
Import de 71,932,943 lignes (42.7M etablissements + 29.2M unites legales) via streaming ZIP.
7 endpoints de recherche testes et operationnels.

---

## Donnees importees

| Table         | Lignes     | Colonnes |
| ------------- | ---------- | -------- |
| etablissement | 42,716,292 | 54       |
| unite_legale  | 29,216,651 | 35       |

### Statistiques

| Metrique                | Valeur     |
| ----------------------- | ---------- |
| Sieges sociaux          | 29,194,868 |
| Codes NAF distincts     | 2,132      |
| Codes postaux distincts | 16,313     |
| Communes distinctes     | 41,196     |

---

## Endpoints API

Base URL: `http://localhost:8081/api`

| Endpoint                              | Methode | Params                                                                        | Description                                 | Statut |
| ------------------------------------- | ------- | ----------------------------------------------------------------------------- | ------------------------------------------- | ------ |
| `/health`                             | GET     | -                                                                             | Health check                                | OK     |
| `/companies/search/naf`               | GET     | `code`, `limit`                                                               | Recherche par code NAF                      | OK     |
| `/companies/search/denomination`      | GET     | `q`, `limit`                                                                  | Recherche par nom (multi-mots AND)          | OK     |
| `/companies/search/codepostal`        | GET     | `q`, `limit`                                                                  | Recherche par code postal (sieges)          | OK     |
| `/companies/search/commune`           | GET     | `q`, `limit`                                                                  | Recherche par commune (ILIKE, sieges)       | OK     |
| `/companies/search/etatadministratif` | GET     | `q`, `limit`                                                                  | Recherche par etat (A=actif, C=cesse)       | OK     |
| `/companies/search/datecreation`      | GET     | `from`, `to`, `limit`                                                         | Recherche par date de creation (YYYY-MM-DD) | OK     |
| `/companies/search/multi`             | GET     | `naf`, `denomination`, `codepostal`, `commune`, `etat`, `from`, `to`, `limit` | Multi-criteres (intersection)               | OK     |

### Exemples de requetes

```
GET /api/companies/search/naf?code=62.01Z&limit=10
GET /api/companies/search/denomination?q=google+france&limit=10
GET /api/companies/search/codepostal?q=75009&limit=10
GET /api/companies/search/commune?q=LYON&limit=10
GET /api/companies/search/etatadministratif?q=A&limit=10
GET /api/companies/search/datecreation?from=2024-01-01&to=2024-12-31&limit=10
GET /api/companies/search/multi?naf=62.01Z&codepostal=75009&limit=10
```

---

## Comparaison BCE Belgique vs SIRENE France

### Fonctionnalites communes

| Feature                     | BCE Belgique                           | SIRENE France                                    |
| --------------------------- | -------------------------------------- | ------------------------------------------------ |
| Recherche par code activite | NACE code (activity table)             | NAF code (etablissement table)                   |
| Recherche par denomination  | ILIKE simple                           | ILIKE multi-mots AND                             |
| Recherche par code postal   | zipcode (address table, REGO)          | code_postal (etablissement, siege)               |
| Recherche par date creation | startdate (DD-MM-YYYY, TO_DATE)        | date_creation (YYYY-MM-DD natif)                 |
| Multi-criteres intersection | nace + denom + zipcode + status + date | naf + denom + codepostal + commune + etat + date |
| Cache Redis                 | gz: compression, 24h TTL               | gz: compression, 24h TTL                         |
| Import bulk                 | CSV direct, pgx CopyFrom               | ZIP streaming, pgx CopyFrom                      |

### Fonctionnalites supplementaires dans SIRENE France

| Feature                          | Detail                                              |
| -------------------------------- | --------------------------------------------------- |
| Recherche par commune            | ILIKE sur libelle_commune_etablissement             |
| Recherche par etat administratif | A (actif) ou C (cesse)                              |
| Denomination multi-mots          | "google france" cherche google AND france           |
| ZIP streaming                    | Pas d'extraction disque, lecture directe depuis ZIP |

### Fonctionnalites presentes dans BCE mais absentes de SIRENE

| Feature BCE                                  | Raison d'absence                                 | Portabilite                                        |
| -------------------------------------------- | ------------------------------------------------ | -------------------------------------------------- |
| Recherche generique `/search/:table/:column` | Non implemente                                   | Facile a porter                                    |
| Middleware ValidateTableName/ColumnName      | Non implemente                                   | Recommande                                         |
| Recherche NACE multi-mots intelligente       | NAF search est exact match seulement             | A ameliorer                                        |
| Contacts (email, tel, web)                   | Absents des donnees SIRENE                       | Necessite enrichissement externe (Apollo.io, INPI) |
| Multi-langue (FR/NL)                         | Donnees SIRENE en francais uniquement            | N/A                                                |
| Table establishment separee                  | Tout est dans etablissement                      | Architecture differente                            |
| Preview/export generique                     | Non implemente                                   | Facile a porter                                    |
| Middleware ParseLimit/Offset/Sort            | parseLimit inline dans handler                   | A refactorer                                       |
| Enrichissement 6 tables                      | 2 tables seulement (unite_legale, etablissement) | Architecture differente                            |

### Donnees disponibles mais non exploitees dans SIRENE

| Champ                   | Colonne                                  | Utilite potentielle             |
| ----------------------- | ---------------------------------------- | ------------------------------- |
| Categorie juridique     | categorie_juridique_unite_legale         | Filtrer SAS/SARL/SA/etc         |
| Tranche effectifs       | tranche_effectifs_unite_legale           | Filtrer par taille              |
| Categorie entreprise    | categorie_entreprise                     | PME/ETI/GE                      |
| Economie sociale        | economie_sociale_solidaire_unite_legale  | Filtre ESS                      |
| Caractere employeur     | caractere_employeur_unite_legale         | Entreprises qui embauchent      |
| Code commune            | code_commune_etablissement               | Code INSEE (plus precis que CP) |
| NAF 2025                | activite_principale_n_a_f25_unite_legale | Nouvelle nomenclature           |
| Coordonnees Lambert     | coordonnee*lambert*\*                    | Geolocalisation                 |
| Date dernier traitement | date*dernier_traitement*\*               | Fraicheur des donnees           |

---

## Architecture technique

### Pipeline d'import

```
sirene_data/*.zip
    |
    v
zip.OpenReader (streaming, pas d'extraction)
    |
    v
csv.Reader (LazyQuotes, buffer 4MB)
    |
    v
camelCase -> snake_case (headers)
    |
    v
DROP TABLE IF EXISTS / CREATE TABLE (TEXT)
    |
    v
8 workers paralleles
    |
    v
Batch 200k lignes -> pgx CopyFrom
    |
    v
PostgreSQL (synchronous_commit=OFF)
```

### Performance d'import mesuree

| Table         | Lignes | Duree estimee |
| ------------- | ------ | ------------- |
| etablissement | 42.7M  | ~6 min        |
| unite_legale  | 29.2M  | ~4 min        |
| Total         | 71.9M  | ~10 min       |

### Stack

| Composant  | Version | Port |
| ---------- | ------- | ---- |
| Go         | 1.24.2  | -    |
| PostgreSQL | 15      | 5434 |
| Redis      | 7       | 6380 |
| Gin        | 1.10.1  | 8081 |
| pgx        | 5.7.5   | -    |

---

## Corrections appliquees

### Bugs corriges

| Fichier                        | Correction                                   |
| ------------------------------ | -------------------------------------------- |
| api/server.go                  | Fix imports `internal/config` -> `config`    |
| company_enrichment.go          | Ajout `rows.Err()` apres boucles de lecture  |
| company_enrichment.go          | Logging des erreurs de scan                  |
| company_handler.go             | Limite max 10,000 sur le parametre `limit`   |
| company_service.go             | Redis configurable via REDIS_HOST/REDIS_PORT |
| csv/pipeline.go                | Buffer channel reduit 100k -> 10k (RAM)      |
| company_search_naf.go          | Ajout `rows.Err()`                           |
| company_search_denomination.go | Ajout `rows.Err()` + recherche multi-mots    |
| company_search_codepostal.go   | Ajout `rows.Err()`                           |
| company_search_datecreation.go | Ajout `rows.Err()`                           |

### Nouvelles fonctionnalites

| Fichier                     | Ajout                                                   |
| --------------------------- | ------------------------------------------------------- |
| company_search_commune.go   | Recherche par commune (ILIKE)                           |
| company_search_etatadmin.go | Recherche par etat administratif                        |
| company_models.go           | Champs Commune et EtatAdministratif dans les criteres   |
| company_handler.go          | Handlers commune et etat administratif                  |
| company_search_multi.go     | Criteres commune et etat dans multi-criteres            |
| server.go                   | Routes `/search/commune` et `/search/etatadministratif` |

---

## Recommandations

### Priorite haute

1. **Ajouter des index PostgreSQL** sur les colonnes de recherche pour ameliorer les performances

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

2. **Middleware de validation** des parametres d'entree (format NAF, code postal, dates)

3. **Reutilisation des connexions** dans le batch inserter (une connexion par worker au lieu d'une par batch)

### Priorite moyenne

4. **Recherche NAF multi-mots** (actuellement exact match seulement)
5. **Recherche par categorie juridique** (SAS, SARL, SA, etc.)
6. **Recherche par tranche effectifs** (filtrer par taille d'entreprise)
7. **Pagination** (ajout d'offset en plus de limit)

### Priorite basse

8. **Endpoint d'export** CSV/JSON
9. **Endpoint de statistiques** (nombre d'entreprises par NAF, par commune, etc.)
10. **Enrichissement INPI** (dirigeants, via data.inpi.fr)
11. **Enrichissement contacts** (Apollo.io pour email/tel/web)
