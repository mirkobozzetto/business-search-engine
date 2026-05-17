# Rapport des modifications API - Recherche semantique NAF

## Resume

La feature de recherche semantique NAF ajoute un service complet de consultation des codes NAF (Nomenclature d'Activites Francaise) avec 4 endpoints dedies, et enrichit les resultats de recherche d'entreprises avec le libelle NAF via un LEFT JOIN sur la table `naf_reference`.

---

## Fichiers crees

### 1. `api/services/naf/naf_service.go` (133 lignes)

Nouveau service NAF avec 4 methodes :

| Methode | Description | Requete SQL |
|---------|-------------|-------------|
| `SearchByLabel` | Recherche par libelle (ILIKE) avec pagination et comptage parallele | `SELECT ... FROM naf_reference WHERE label ILIKE $1 ORDER BY code LIMIT $2 OFFSET $3` |
| `ListSections` | Liste les sections NAF avec comptage par section | `SELECT section_code, section_label, COUNT(*) FROM naf_reference GROUP BY section_code, section_label ORDER BY section_code` |
| `GetByCode` | Recuperation d'un code NAF specifique | `SELECT ... FROM naf_reference WHERE code = $1` |
| `GetBySection` | Liste les codes NAF d'une section donnee | `SELECT ... FROM naf_reference WHERE section_code = $1 ORDER BY code` |

**Modeles definis :**

```go
type NafCode struct {
    Code         string `json:"code"`
    Label        string `json:"label"`
    SectionCode  string `json:"section_code"`
    SectionLabel string `json:"section_label"`
}

type NafSection struct {
    Code  string `json:"code"`
    Label string `json:"label"`
    Count int    `json:"count"`
}
```

### 2. `api/services/naf/naf_handler.go` (112 lignes)

Handler HTTP Gin avec 4 methodes correspondant aux 4 endpoints. Inclut les fonctions utilitaires `parseLimit` (defaut 100, max 10000) et `parseOffset`.

---

## Fichiers modifies

### 3. `api/server.go`

**Modifications :**
- Import ajoute : `sirene-importer/api/services/naf`
- Champ `nafHandler *naf.Handler` ajoute a la struct `Server`
- Instanciation du service et handler NAF dans `NewServer()`
- 4 nouvelles routes ajoutees dans `setupRoutes()`

**Nouvelles routes :**

```go
nafGroup := api.Group("/naf")
nafGroup.GET("/search", s.nafHandler.SearchByLabel)
nafGroup.GET("/sections", s.nafHandler.ListSections)
nafGroup.GET("/code/:code", s.nafHandler.GetByCode)
nafGroup.GET("/section/:code", s.nafHandler.GetBySection)
```

### 4. `api/models/company_models.go`

**Modification :**
- Champ `NafLabel string` ajoute a la struct `CompanyResult` (ligne 25) avec le tag JSON `json:"naf_label,omitempty"`

### 5. `api/services/company/company_query.go`

**Modifications :**
- Ajout de `COALESCE(naf.label, '')` dans la constante `companySelectFields` pour recuperer le libelle NAF
- Ajout du scan de `&c.NafLabel` dans `scanCompanyRow()`
- Ajout du `LEFT JOIN naf_reference naf ON COALESCE(NULLIF(e.activite_principale_etablissement, ''), u.activite_principale_unite_legale, '') = naf.code` dans la requete `dataQuery` de `searchCompanies()`

### 6. `api/services/company/company_search_identifier.go`

**Modifications :**
- Ajout du meme `LEFT JOIN naf_reference naf` dans `lookupBySiren()` (ligne 30)
- Ajout du meme `LEFT JOIN naf_reference naf` dans `lookupBySiret()` (ligne 55)

---

## Nouveaux endpoints

| Methode | URL | Parametres | Reponse |
|---------|-----|------------|---------|
| GET | `/api/naf/search` | `q` (obligatoire), `limit` (defaut 100, max 10000), `offset` (defaut 0) | Liste de `NafCode` avec pagination (`meta.total`, `meta.count`, `meta.pages`) |
| GET | `/api/naf/sections` | Aucun | Liste de `NafSection` (code, label, count) |
| GET | `/api/naf/code/:code` | `:code` (path, ex: `62.01Z`) | Un seul `NafCode` ou 404 |
| GET | `/api/naf/section/:code` | `:code` (path, ex: `J`) | Liste de `NafCode` de la section |

### Exemples de reponse

**GET /api/naf/search?q=informatique&limit=2**
```json
{
  "success": true,
  "data": [
    {
      "code": "62.01Z",
      "label": "Programmation informatique",
      "section_code": "J",
      "section_label": "Information et communication"
    }
  ],
  "meta": {
    "total": 5,
    "count": 2,
    "limit": 2,
    "offset": 0,
    "page": 1,
    "pages": 3
  }
}
```

**GET /api/naf/sections**
```json
{
  "success": true,
  "data": [
    {
      "code": "A",
      "label": "Agriculture, sylviculture et peche",
      "count": 36
    }
  ]
}
```

---

## Modifications aux requetes existantes

Toutes les requetes de recherche d'entreprises incluent desormais un `LEFT JOIN` sur la table `naf_reference` pour enrichir les resultats avec le libelle NAF :

```sql
LEFT JOIN naf_reference naf
  ON COALESCE(
    NULLIF(e.activite_principale_etablissement, ''),
    u.activite_principale_unite_legale,
    ''
  ) = naf.code
```

**Logique de jointure :** utilise en priorite le code NAF de l'etablissement (`activite_principale_etablissement`), puis celui de l'unite legale (`activite_principale_unite_legale`) en fallback.

**Requetes impactees :**
- `searchCompanies()` dans `company_query.go` (recherche multi-criteres, par NAF, denomination, code postal, commune, etat administratif, date creation)
- `lookupBySiren()` dans `company_search_identifier.go`
- `lookupBySiret()` dans `company_search_identifier.go`

---

## Bilan quantitatif

| Fichier | Statut | Lignes |
|---------|--------|--------|
| `api/services/naf/naf_service.go` | Cree | 133 |
| `api/services/naf/naf_handler.go` | Cree | 112 |
| `api/server.go` | Modifie | ~10 lignes ajoutees |
| `api/models/company_models.go` | Modifie | 1 ligne ajoutee |
| `api/services/company/company_query.go` | Modifie | ~5 lignes ajoutees/modifiees |
| `api/services/company/company_search_identifier.go` | Modifie | ~4 lignes ajoutees |
| **Total** | | **~265 lignes de code ajoutees** |
