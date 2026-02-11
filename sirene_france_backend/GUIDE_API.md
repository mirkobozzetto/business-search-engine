# SIRENE France API - Guide d'utilisation

API de recherche d'entreprises francaises basee sur la base SIRENE de l'INSEE (72 millions de lignes).
Toutes les donnees sont en local, aucune API externe, aucune cle, aucun quota.

---

## Demarrage rapide

```bash
# 1. Lancer les services (PostgreSQL + Redis)
make sirene-up

# 2. Lancer l'API (dans un autre terminal)
make sirene-api

# 3. Verifier que ca fonctionne
curl -s "localhost:8081/api/health" | jq .
```

L'API tourne sur le port **8081**.

---

## Recherche par SIREN ou SIRET

Recherche directe d'une entreprise par son identifiant.

```bash
# Par SIRET (14 chiffres)
curl -s "localhost:8081/api/companies/lookup/97994855100010" | jq .

# Par SIREN (9 chiffres)
curl -s "localhost:8081/api/companies/lookup/979948551" | jq .
```

Exemple de resultat :

```json
{
  "siren": "979948551",
  "denomination": "CREACH AGENCY",
  "categorie_juridique": "5710",
  "date_creation": "2023-09-19",
  "etat_administratif": "A",
  "tranche_effectifs": "NN",
  "categorie_entreprise": "PME",
  "naf_code": "62.01Z",
  "siret": "97994855100010",
  "numero_voie": "60",
  "type_voie": "RUE",
  "libelle_voie": "FRANCOIS IER",
  "code_postal": "75008",
  "libelle_commune": "PARIS"
}
```

---

## Recherche multi-criteres

C'est l'endpoint le plus puissant. Combine autant de criteres que tu veux.

```
GET /api/companies/search/multi?critere1=valeur1&critere2=valeur2&limit=N&offset=N
```

### Tous les criteres disponibles

| Parametre             | Description                   | Exemple      | Type                             |
| --------------------- | ----------------------------- | ------------ | -------------------------------- |
| `naf`                 | Code d'activite NAF           | `62.01Z`     | Exact                            |
| `denomination`        | Nom de l'entreprise           | `creach`     | Contient (insensible a la casse) |
| `codepostal`          | Code postal                   | `75008`      | Exact                            |
| `commune`             | Nom de la commune             | `paris`      | Contient (insensible a la casse) |
| `etat`                | Etat administratif            | `A` ou `C`   | Exact                            |
| `from`                | Date de creation (debut)      | `2025-01-01` | >= date                          |
| `to`                  | Date de creation (fin)        | `2025-12-31` | <= date                          |
| `categorie_juridique` | Forme juridique               | `5710`       | Exact                            |
| `tranche_effectifs`   | Tranche d'effectifs           | `03`         | Exact                            |
| `limit`               | Nombre de resultats par page  | `10`         | Pagination                       |
| `offset`              | Decalage (pour page suivante) | `10`         | Pagination                       |

### Exemples concrets

```bash
# Boites de dev actives a Paris creees en 2025
curl -s "localhost:8081/api/companies/search/multi?naf=62.01Z&commune=paris&etat=A&from=2025-01-01&to=2025-12-31&limit=10" | jq .

# SAS de consulting IT a Lyon
curl -s "localhost:8081/api/companies/search/multi?naf=62.02A&commune=lyon&etat=A&categorie_juridique=5710&limit=10" | jq .

# Restaurants crees en 2024 a Marseille
curl -s "localhost:8081/api/companies/search/multi?naf=56.10A&commune=marseille&from=2024-01-01&to=2024-12-31&limit=10" | jq .

# Nouvelles entreprises informatiques a Paris en 2026
curl -s "localhost:8081/api/companies/search/multi?naf=62.01Z&commune=paris&etat=A&from=2026-01-01&to=2026-02-28&limit=10" | jq .

# Agences de pub actives a Bordeaux
curl -s "localhost:8081/api/companies/search/multi?naf=73.11Z&commune=bordeaux&etat=A&limit=10" | jq .
```

---

## Recherches simples (un seul critere)

```bash
# Par code NAF
curl -s "localhost:8081/api/companies/search/naf?code=62.01Z&limit=5" | jq .

# Par nom d'entreprise
curl -s "localhost:8081/api/companies/search/denomination?q=creach&limit=5" | jq .

# Par code postal
curl -s "localhost:8081/api/companies/search/codepostal?q=75008&limit=5" | jq .

# Par commune
curl -s "localhost:8081/api/companies/search/commune?q=paris&limit=5" | jq .

# Par etat administratif
curl -s "localhost:8081/api/companies/search/etatadministratif?q=A&limit=5" | jq .

# Par date de creation
curl -s "localhost:8081/api/companies/search/datecreation?from=2025-01-01&to=2025-12-31&limit=5" | jq .
```

---

## Pagination

Tous les endpoints supportent `limit` et `offset`.

```bash
# Page 1 (resultats 1 a 10)
curl -s "localhost:8081/api/companies/search/multi?naf=62.01Z&commune=paris&etat=A&limit=10&offset=0" | jq .

# Page 2 (resultats 11 a 20)
curl -s "localhost:8081/api/companies/search/multi?naf=62.01Z&commune=paris&etat=A&limit=10&offset=10" | jq .

# Page 3 (resultats 21 a 30)
curl -s "localhost:8081/api/companies/search/multi?naf=62.01Z&commune=paris&etat=A&limit=10&offset=20" | jq .
```

La reponse contient les infos de pagination dans `meta` :

```json
"meta": {
  "count": 10,
  "total": 3790,
  "limit": 10,
  "offset": 0,
  "page": 1,
  "pages": 379
}
```

---

## Comprendre les champs de reponse

| Champ                  | Signification                                | Exemple          |
| ---------------------- | -------------------------------------------- | ---------------- |
| `siren`                | Identifiant entreprise (9 chiffres)          | `979948551`      |
| `siret`                | Identifiant etablissement (14 chiffres)      | `97994855100010` |
| `denomination`         | Nom de la societe                            | `CREACH AGENCY`  |
| `categorie_juridique`  | Forme juridique (voir table ci-dessous)      | `5710`           |
| `date_creation`        | Date de creation (AAAA-MM-JJ)                | `2023-09-19`     |
| `etat_administratif`   | `A` = Active, `C` = Cessation                | `A`              |
| `tranche_effectifs`    | Tranche d'effectifs (voir table ci-dessous)  | `NN`             |
| `categorie_entreprise` | Taille : PME, ETI, GE                        | `PME`            |
| `naf_code`             | Code d'activite NAF                          | `62.01Z`         |
| `enseigne`             | Nom commercial                               |                  |
| `numero_voie`          | Numero de rue                                | `60`             |
| `type_voie`            | Type (RUE, AVENUE, BOULEVARD...)             | `RUE`            |
| `libelle_voie`         | Nom de la voie                               | `FRANCOIS IER`   |
| `code_postal`          | Code postal                                  | `75008`          |
| `libelle_commune`      | Ville                                        | `PARIS`          |
| `[ND]`                 | Non Diffusible (auto-entrepreneurs proteges) |                  |

---

## Codes NAF les plus utiles

### Informatique et digital

| Code NAF | Activite                                                       |
| -------- | -------------------------------------------------------------- |
| `62.01Z` | Programmation informatique                                     |
| `62.02A` | Conseil en systemes et logiciels informatiques                 |
| `62.02B` | Tierce maintenance de systemes et d'applications informatiques |
| `62.03Z` | Gestion d'installations informatiques                          |
| `62.09Z` | Autres activites informatiques                                 |
| `63.11Z` | Traitement de donnees, hebergement et activites connexes       |
| `63.12Z` | Portails internet                                              |

### Conseil et services aux entreprises

| Code NAF | Activite                                                |
| -------- | ------------------------------------------------------- |
| `70.10Z` | Activites des sieges sociaux                            |
| `70.21Z` | Conseil en relations publiques et communication         |
| `70.22Z` | Conseil pour les affaires et autres conseils de gestion |
| `69.10Z` | Activites juridiques                                    |
| `69.20Z` | Activites comptables                                    |
| `73.11Z` | Activites des agences de publicite                      |
| `73.12Z` | Regie publicitaire de medias                            |
| `74.10Z` | Activites specialisees de design                        |

### Commerce et restauration

| Code NAF | Activite                                  |
| -------- | ----------------------------------------- |
| `47.11B` | Commerce d'alimentation generale          |
| `47.71Z` | Commerce de detail d'habillement          |
| `47.91A` | Vente a distance sur catalogue general    |
| `47.91B` | Vente a distance sur catalogue specialise |
| `56.10A` | Restauration traditionnelle               |
| `56.10B` | Cafeterias et autres libres-services      |
| `56.10C` | Restauration de type rapide               |
| `56.30Z` | Debits de boissons                        |

### Construction et immobilier

| Code NAF | Activite                                           |
| -------- | -------------------------------------------------- |
| `41.20A` | Construction de maisons individuelles              |
| `43.21A` | Travaux d'installation electrique                  |
| `43.22A` | Travaux d'installation d'eau et de gaz             |
| `43.34Z` | Travaux de peinture et vitrerie                    |
| `68.20A` | Location de logements                              |
| `68.20B` | Location de terrains et d'autres biens immobiliers |
| `68.31Z` | Agences immobilieres                               |

### Sante

| Code NAF | Activite                                         |
| -------- | ------------------------------------------------ |
| `86.21Z` | Activite des medecins generalistes               |
| `86.22C` | Autres activites des medecins specialistes       |
| `86.23Z` | Pratique dentaire                                |
| `86.90D` | Activites des infirmiers et des sages-femmes     |
| `86.90E` | Activites des professionnels de la reeducation   |
| `86.90F` | Activites de sante humaine non classees ailleurs |

### Transport

| Code NAF | Activite                                 |
| -------- | ---------------------------------------- |
| `49.32Z` | Transports de voyageurs par taxis        |
| `49.41A` | Transports routiers de fret interurbains |
| `53.20Z` | Autres activites de poste et de courrier |

### Education et formation

| Code NAF | Activite                                                        |
| -------- | --------------------------------------------------------------- |
| `85.51Z` | Enseignement de disciplines sportives et d'activites de loisirs |
| `85.52Z` | Enseignement culturel                                           |
| `85.59A` | Formation continue d'adultes                                    |
| `85.59B` | Autres enseignements                                            |

---

## Formes juridiques (categorie_juridique)

| Code   | Forme juridique                                               |
| ------ | ------------------------------------------------------------- |
| `1000` | Entrepreneur individuel (auto-entrepreneur, micro-entreprise) |
| `2110` | Indivision                                                    |
| `3220` | Societe en nom collectif (SNC)                                |
| `5202` | Societe en nom collectif cooperative                          |
| `5485` | Societe europeenne (SE)                                       |
| `5499` | Societe par actions simplifiee (SAS autre)                    |
| `5499` | SA a conseil d'administration                                 |
| `5599` | SA a directoire                                               |
| `5710` | Societe par actions simplifiee (SAS)                          |
| `6540` | Societe civile immobiliere (SCI)                              |
| `6541` | SCI de construction vente                                     |
| `6598` | Societe civile de moyens (SCM)                                |
| `6599` | Autres societes civiles                                       |
| `9110` | Syndicat de proprieteires                                     |
| `9220` | Association declaree                                          |

---

## Tranches d'effectifs salaries

| Code | Tranche                 |
| ---- | ----------------------- |
| `NN` | Non renseigne           |
| `00` | 0 salarie               |
| `01` | 1 ou 2 salaries         |
| `02` | 3 a 5 salaries          |
| `03` | 6 a 9 salaries          |
| `11` | 10 a 19 salaries        |
| `12` | 20 a 49 salaries        |
| `21` | 50 a 99 salaries        |
| `22` | 100 a 199 salaries      |
| `31` | 200 a 249 salaries      |
| `32` | 250 a 499 salaries      |
| `41` | 500 a 999 salaries      |
| `42` | 1 000 a 1 999 salaries  |
| `51` | 2 000 a 4 999 salaries  |
| `52` | 5 000 a 9 999 salaries  |
| `53` | 10 000 salaries et plus |

---

## Etat administratif

| Code | Signification      |
| ---- | ------------------ |
| `A`  | Active             |
| `C`  | Cessation (fermee) |

---

## Sources de donnees et telechargements

Toutes les donnees sont **gratuites et open data**.

### Base SIRENE (deja integree)

| Fichier                       | URL                                                                                                  | Format      |
| ----------------------------- | ---------------------------------------------------------------------------------------------------- | ----------- |
| `StockUniteLegale_utf8.zip`   | https://www.data.gouv.fr/datasets/base-sirene-des-entreprises-et-de-leurs-etablissements-siren-siret | CSV (~2 Go) |
| `StockEtablissement_utf8.zip` | Meme page                                                                                            | CSV (~5 Go) |

Pour reimporter les donnees a jour :

```bash
make sirene-reimport
```

### Codes NAF (a integrer)

| Source                         | URL                                                                                       | Format  |
| ------------------------------ | ----------------------------------------------------------------------------------------- | ------- |
| NAF rev.2 complete (732 codes) | https://public.opendatasoft.com/explore/dataset/naf2008 -> Export -> CSV                  | CSV     |
| NAF officielle INSEE           | https://www.insee.fr/fr/information/2406147 -> "Telecharger les fichiers de la NAF rev.2" | CSV/PDF |

### Dirigeants et comptes annuels - INPI (a integrer)

Gratuit, necessite un compte (validation instantanee).

1. Creer un compte : https://data.inpi.fr/register
2. Activer l'acces SFTP dans "Mes acces API / SFTP"
3. Se connecter au serveur SFTP

| Dossier SFTP        | Contenu                                                              | Format          |
| ------------------- | -------------------------------------------------------------------- | --------------- |
| `/Stock/`           | 5M+ societes avec dirigeants, representants, beneficiaires effectifs | JSON (.json.gz) |
| `/Flux/`            | Mises a jour quotidiennes                                            | JSON            |
| `/Comptes annuels/` | Bilans, comptes de resultat depuis 2017                              | JSON + PDF      |
| `/Actes/`           | Statuts, PV d'AG depuis 1993                                         | PDF             |

Documentation technique : https://data.inpi.fr/content/editorial/Acces_API_Entreprises

---

## Architecture technique

```
curl (client) --> API Go (Gin) :8081 --> Redis :6380 (cache 24h)
                                    --> PostgreSQL :5434 (72M lignes)
```

- **Go 1.24** avec framework Gin
- **PostgreSQL 15** : 2 tables (unite_legale 29M lignes + etablissement 42M lignes)
- **Redis 7** : cache avec compression automatique pour les gros jeux de donnees
- Les requetes multi-criteres font une seule requete SQL avec JOIN
- Les recherches simples sont cachees 24h dans Redis

---

## Commandes disponibles

```bash
make sirene-up           # Demarre PostgreSQL + Redis
make sirene-down         # Arrete les services
make sirene-api          # Lance l'API sur le port 8081
make sirene-import       # Importe les fichiers SIRENE (premiere fois)
make sirene-reimport     # Reimporte les donnees (mise a jour)
make sirene-indexes      # Cree les index PostgreSQL
make help                # Affiche l'aide
```
