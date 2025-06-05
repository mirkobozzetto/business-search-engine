# API BCE - Guide Utilisateur

## Vue d'ensemble

API de recherche d'entreprises belges avec 47M+ lignes réparties sur 10 tables. Système basé sur `entitynumber` comme clé universelle.

## Tables Principales

```
entitynumber (clé universelle)
├── activity (36M lignes) - Codes NACE, classifications
├── address (2.8M lignes) - Adresses complètes
├── contact (683k lignes) - Email, téléphone, web, fax
├── denomination (3.3M lignes) - Noms FR/NL/DE
├── enterprise (1.9M lignes) - Statut, forme juridique, dates
└── establishment (1.6M lignes) - Unités locales
```

## Comment ça Marche

### Principe Simple

1. **Point d'entrée** → Récupère les entitynumbers
2. **Enrichissement automatique** → Joint les 6 tables
3. **Cache Redis** → Stocke le tout
4. **Multi-critères** → Intersection des caches

### Prérequis Important

Pour utiliser `/search/multi`, il faut OBLIGATOIREMENT créer les caches individuels d'abord.

## Endpoints Disponibles

```bash
# Recherche par code NACE
curl "localhost:8080/api/companies/search/nace?code=62010" | jq

# Recherche par nom d'entreprise
curl "localhost:8080/api/companies/search/denomination?q=informatique" | jq

# Recherche par code postal
curl "localhost:8080/api/companies/search/zipcode?q=1000" | jq

# Recherche par date de création
curl "localhost:8080/api/companies/search/startdate?from=01-01-2025" | jq

# Multi-critères (après avoir créé les caches)
curl "localhost:8080/api/companies/search/multi?nace=62010&zipcode=1000" | jq
```

## Workflow Basique

### Exemple : Coaches de santé créés en 2025

**Étape 1 : Cache NACE**

```bash
curl "localhost:8080/api/companies/search/nace?code=56309" | jq
# Crée le cache pour toutes les entreprises NACE 56309
```

**Étape 2 : Cache Date**

```bash
curl "localhost:8080/api/companies/search/startdate?from=01-01-2025" | jq
# Crée le cache pour toutes les entreprises créées en 2025+
```

**Étape 3 : Intersection**

```bash
curl "localhost:8080/api/companies/search/multi?nace=56309&startdate_from=01-01-2025" | jq
# Intersection des 2 caches = entreprises NACE 56309 créées en 2025
```

## Erreurs Fréquentes

### Cache manquant

```bash
# ❌ Erreur
curl "localhost:8080/api/companies/search/multi?nace=62010&zipcode=1000" | jq
# "NACE cache not found: 62010. Please search by NACE first"

# ✅ Solution
curl "localhost:8080/api/companies/search/nace?code=62010" | jq
curl "localhost:8080/api/companies/search/zipcode?q=1000" | jq
curl "localhost:8080/api/companies/search/multi?nace=62010&zipcode=1000" | jq
```

### Paramètre incorrect

```bash
# ❌ Erreur - Paramètre "nace="
curl "localhost:8080/api/companies/search/nace?nace=62010" | jq

# ✅ Correct - Paramètre "code="
curl "localhost:8080/api/companies/search/nace?code=62010" | jq
```

## Structure de Réponse

Chaque entreprise retournée contient :

- **Données de base** : entitynumber, denomination, zipcode, city
- **Contacts** : email, web, tel, fax
- **Détails** : status, juridical_form, start_date
- **Données complètes** : denominations[], addresses[], contacts[], activities[], etc.

## Cas d'Usage Simples

### Toutes les entreprises d'un secteur

```bash
curl "localhost:8080/api/companies/search/nace?code=62010" | jq
```

### Entreprises d'une ville

```bash
curl "localhost:8080/api/companies/search/zipcode?q=1000" | jq
```

### Croisement secteur + ville

```bash
curl "localhost:8080/api/companies/search/nace?code=62010" | jq
curl "localhost:8080/api/companies/search/zipcode?q=1000" | jq
curl "localhost:8080/api/companies/search/multi?nace=62010&zipcode=1000" | jq
```

### Nouvelles entreprises d'un secteur

```bash
curl "localhost:8080/api/companies/search/nace?code=62010" | jq
curl "localhost:8080/api/companies/search/startdate?from=01-01-2024" | jq
curl "localhost:8080/api/companies/search/multi?nace=62010&startdate_from=01-01-2024" | jq
```
