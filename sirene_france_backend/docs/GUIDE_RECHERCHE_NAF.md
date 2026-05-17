# Guide de la recherche sémantique NAF

## Vue d'ensemble

L'API Sirene France permet de rechercher des entreprises par activité via la nomenclature NAF (Nomenclature des Activités Françaises). Les 732 codes NAF rev.2 sont stockés dans une table de référence `naf_reference` et enrichissent automatiquement les résultats d'entreprises avec le libellé d'activité en clair.

**Base de données** : 732 codes NAF, 21 sections, ~36M établissements, ~24M unités légales.

---

## Démarrage rapide

```bash
# 1. Importer les codes NAF
./sirene-api naf

# 2. Créer les indexes (si pas déjà fait)
./sirene-api indexes

# 3. Lancer l'API
./sirene-api api
```

---

## Endpoints disponibles

### 1. Recherche par label (recherche sémantique)

```
GET /api/naf/search?q={texte}&limit={n}&offset={n}
```

Recherche les codes NAF dont le libellé contient le texte donné (insensible à la casse, via trigram GIN).

**Exemples et résultats réels :**

```bash
curl http://localhost:8081/api/naf/search?q=informatique
```

| Code | Label | Section |
|------|-------|---------|
| 46.51Z | Commerce de gros d'ordinateurs, d'équipements informatiques périphériques et de logiciels | G |
| 62.01Z | Programmation informatique | J |
| 62.02A | Conseil en systèmes et logiciels informatiques | J |
| 62.02B | Tierce maintenance de systèmes et d'applications informatiques | J |
| 62.03Z | Gestion d'installations informatiques | J |
| 62.09Z | Autres activités informatiques | J |
| 77.33Z | Location et location-bail de machines de bureau et de matériel informatique | N |

**7 résultats** pour "informatique".

```bash
curl http://localhost:8081/api/naf/search?q=boulangerie
```

| Code | Label | Section |
|------|-------|---------|
| 10.71B | Cuisson de produits de boulangerie | C |
| 10.71C | Boulangerie et boulangerie-pâtisserie | C |

**2 résultats** pour "boulangerie".

```bash
curl http://localhost:8081/api/naf/search?q=transport
```

**25 résultats** couvrant : transport d'électricité, ferroviaire, routier, maritime, fluvial, aérien, spatial, par conduites, services auxiliaires, location de matériels, affrètement.

```bash
curl http://localhost:8081/api/naf/search?q=conseil&limit=5
```

| Code | Label | Section |
|------|-------|---------|
| 62.02A | Conseil en systèmes et logiciels informatiques | J |
| 70.21Z | Conseil en relations publiques et communication | M |
| 70.22Z | Conseil pour les affaires et autres conseils de gestion | M |

**3 résultats** pour "conseil" (limit de 5, 3 trouvés).

### 2. Lister les sections

```
GET /api/naf/sections
```

Retourne les 21 sections de la nomenclature avec le nombre de codes par section.

**Extrait des résultats :**

| Section | Label | Codes |
|---------|-------|-------|
| A | Agriculture, sylviculture et pêche | 39 |
| C | Industrie manufacturière | 259 |
| G | Commerce, réparation d'automobiles et de motocycles | 94 |
| J | Information et communication | 33 |
| M | Activités spécialisées, scientifiques et techniques | 44 |
| U | Activités extra-territoriales | 1 |

### 3. Détail d'un code NAF

```
GET /api/naf/code/:code
```

```bash
curl http://localhost:8081/api/naf/code/62.01Z
```

```json
{
  "success": true,
  "data": {
    "code": "62.01Z",
    "label": "Programmation informatique",
    "section_code": "J",
    "section_label": "Information et communication"
  }
}
```

```bash
curl http://localhost:8081/api/naf/code/56.10A
```

```json
{
  "success": true,
  "data": {
    "code": "56.10A",
    "label": "Restauration traditionnelle",
    "section_code": "I",
    "section_label": "Hébergement et restauration"
  }
}
```

### 4. Codes d'une section

```
GET /api/naf/section/:code
```

```bash
curl http://localhost:8081/api/naf/section/J
```

Retourne les **33 codes** de la section "Information et communication" : édition, production audiovisuelle, télécommunications, informatique, traitement de données, agences de presse.

```bash
curl http://localhost:8081/api/naf/section/C
```

Retourne les **259 codes** de la section "Industrie manufacturière" (la plus grande).

---

## Enrichissement des résultats d'entreprises

Toutes les recherches d'entreprises retournent désormais le champ `naf_label` en plus du `naf_code`. Le LEFT JOIN sur `naf_reference` enrichit automatiquement les résultats.

### Exemples réels

**Recherche par code NAF :**

```bash
curl "http://localhost:8081/api/companies/search/naf?code=62.01Z&limit=3"
```

```
227 644 entreprises avec le code 62.01Z "Programmation informatique"
```

| Entreprise | Ville | naf_code | naf_label |
|------------|-------|----------|-----------|
| AMPLYHUB | Paris 75015 | 62.01Z | Programmation informatique |
| ... | Castelnau-le-Lez | 62.01Z | Programmation informatique |
| ... | Nice 06300 | 62.01Z | Programmation informatique |

**Recherche par dénomination :**

```bash
curl "http://localhost:8081/api/companies/search/denomination?q=google&limit=3"
```

| Entreprise | SIREN | naf_label |
|------------|-------|-----------|
| GOOGLE PAYMENT IRELAND LIMITED | — | Autres intermédiations monétaires |
| GOOGLE CLOUD FRANCE | 910738392 | Activités des syndicats de salariés |
| GOOGLE CLOUD FRANCE | 881721583 | Commerce de gros d'ordinateurs... |

**Recherche multi-critères (NAF + localisation) :**

```bash
curl "http://localhost:8081/api/companies/search/multi?naf=56.10A&commune=paris&limit=3"
```

```
28 454 restaurants traditionnels à Paris
```

| Entreprise | Code postal | naf_label |
|------------|-------------|-----------|
| PANTHERA | 75015 | Restauration traditionnelle |
| B.O & FILS / LA CAVE BAJLA | 75018 | Restauration traditionnelle |
| SAVEURS DE MO | 75011 | Restauration traditionnelle |

```bash
curl "http://localhost:8081/api/companies/search/multi?naf=62.01Z&codepostal=75008&limit=3"
```

```
5 850 entreprises IT dans le 8ème arrondissement
```

| Entreprise | Adresse | naf_label |
|------------|---------|-----------|
| CORELYNK | 66 Avenue des Champs Elysées | Programmation informatique |

**Lookup par SIREN :**

```bash
curl http://localhost:8081/api/companies/lookup/356000000
```

```json
{
  "siren": "356000000",
  "denomination": "LA POSTE",
  "naf_code": "53.10Z",
  "naf_label": "Activités de poste dans le cadre d'une obligation de service universel",
  "code_postal": "75015",
  "libelle_commune": "PARIS 15",
  "categorie_entreprise": "GE",
  "date_creation": "1991-01-01"
}
```

---

## Conseils pour les meilleurs résultats

### Recherche par label NAF

| Requête | Résultats | Qualité |
|---------|-----------|---------|
| `q=informatique` | 7 | Excellent |
| `q=boulangerie` | 2 | Excellent |
| `q=transport` | 25 | Excellent |
| `q=conseil` | 3 | Bon |
| `q=restaurant` | 0 | Mauvais |

**Bonnes pratiques :**
- Utiliser le mot exact tel qu'il apparaît dans les libellés NAF
- Chercher "restauration" plutôt que "restaurant"
- Chercher "programmation" plutôt que "développement"
- En cas de doute, chercher un mot partiel : "boul" trouvera "boulangerie"
- Utiliser `/api/naf/sections` pour naviguer par section d'abord

### Recherche multi-critères d'entreprises

Les critères combinables pour `/api/companies/search/multi` :

| Paramètre | Description | Exemple |
|-----------|-------------|---------|
| `naf` | Code NAF exact | `62.01Z` |
| `denomination` | Nom d'entreprise (ILIKE) | `google` |
| `codepostal` | Code postal exact | `75008` |
| `commune` | Commune (ILIKE) | `paris` |
| `etat` | État administratif | `A` (actif) |
| `from` / `to` | Plage de date de création | `2020-01-01` |
| `categorie_juridique` | Catégorie juridique | `5710` (SAS) |
| `tranche_effectifs` | Tranche d'effectifs | `12` (20-49) |
| `limit` | Nombre max de résultats | `100` (défaut) |
| `offset` | Décalage pour pagination | `0` (défaut) |

**Scénario type : trouver des ESN à Paris créées après 2020 :**

```bash
curl "http://localhost:8081/api/companies/search/multi?naf=62.01Z&commune=paris&from=2020-01-01&etat=A&limit=50"
```

---

## Limitation actuelle et améliorations prévues

### Recherche textuelle vs sémantique

La recherche actuelle utilise `ILIKE '%mot%'` avec un index trigram GIN. C'est rapide et efficace pour les correspondances exactes, mais ne gère pas :

- **Le stemming** : "restaurant" ne trouve pas "restauration"
- **Les synonymes** : "développeur" ne trouve pas "programmation"
- **Les fautes de frappe** : "informatque" ne trouvera rien

**Amélioration prévue** : migration vers PostgreSQL full-text search (`tsvector` / `tsquery`) avec un dictionnaire français pour le stemming automatique. Cela permettrait :

```sql
-- Recherche actuelle (ILIKE)
WHERE label ILIKE '%restaurant%'  -- ne trouve PAS "restauration"

-- Recherche améliorée (full-text)
WHERE to_tsvector('french', label) @@ to_tsquery('french', 'restaurant')
-- trouvera "restauration", "restaurants", "restaurateur"
```

### Enrichissement CRM (futur)

Tables prévues pour enrichir les données d'entreprises :

| Table | Champs | Source |
|-------|--------|--------|
| `company_enrichment` | website, phone, linkedin, founder_name | Scraping, APIs tierces |
| `company_contacts` | contacts JSONB (nom, email, poste, linkedin) | LinkedIn, Societe.com |

Ces tables seront liées par SIREN et alimenteront de nouveaux champs dans `CompanyResult` via LEFT JOIN, exactement comme `naf_reference` enrichit `naf_label` aujourd'hui.

### Autres améliorations possibles

- **Cache Redis** sur les endpoints NAF pour les requêtes fréquentes
- **Autocomplétion** : endpoint `/api/naf/autocomplete?q=info` retournant des suggestions au fil de la frappe
- **Recherche par code partiel** : `62.*` pour tous les codes de la division 62
- **Export CSV** : `/api/companies/export?naf=62.01Z&commune=paris&format=csv`
- **Statistiques** : `/api/naf/stats/62.01Z` retournant le nombre d'entreprises, la répartition géographique, l'évolution temporelle

---

## Architecture technique

```
data/naf_codes.json (732 codes, 21 sections)
         │
         ▼
   csv/naf_loader.go ──→ TABLE naf_reference
         │                  code TEXT (PK)
         │                  label TEXT
         │                  section_code TEXT
         │                  section_label TEXT
         │
         ▼
   Indexes :
     idx_naf_ref_label_trgm (GIN trigram sur label)
     idx_naf_ref_section (B-tree sur section_code)
         │
         ▼
   api/services/naf/ ──→ 4 endpoints REST
         │
         ▼
   api/services/company/ ──→ LEFT JOIN naf_reference
                              enrichit naf_label sur TOUTES
                              les recherches d'entreprises
```

---

## Résultats des tests

**Date** : 11 février 2026

### Tests endpoints NAF (10 tests)

| # | Endpoint | Résultat | Détails |
|---|----------|----------|---------|
| 1 | /api/naf/sections | OK | 21 sections |
| 2 | /api/naf/search?q=informatique | OK | 7 résultats |
| 3 | /api/naf/search?q=boulangerie | OK | 2 résultats |
| 4 | /api/naf/search?q=transport | OK | 25 résultats |
| 5 | /api/naf/search?q=restaurant | ECHEC | 0 résultat (limitation ILIKE) |
| 6 | /api/naf/code/62.01Z | OK | Programmation informatique |
| 7 | /api/naf/code/56.10A | OK | Restauration traditionnelle |
| 8 | /api/naf/section/J | OK | 33 codes |
| 9 | /api/naf/section/C | OK | 259 codes |
| 10 | /api/naf/search?q=conseil&limit=5 | OK | 3 résultats |

### Tests enrichissement entreprises (6 tests)

| # | Endpoint | Résultat | naf_label |
|---|----------|----------|-----------|
| 1 | search/naf?code=62.01Z | OK | Programmation informatique (227 644 entreprises) |
| 2 | search/denomination?q=google | OK | 3 labels différents |
| 3 | multi?naf=56.10A&commune=paris | OK | Restauration traditionnelle (28 454 entreprises) |
| 4 | multi?naf=62.01Z&codepostal=75008 | OK | Programmation informatique (5 850 entreprises) |
| 5 | lookup/356000000 | OK | LA POSTE - Activités de poste... |
| 6 | search/commune?q=lyon | OK | 3 labels différents (347 987 entreprises) |

**Bilan : 15/16 tests OK. 1 limitation identifiée (stemming).**
