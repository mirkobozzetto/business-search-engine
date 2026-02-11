# Rapport : Changements Base de Donnees - Recherche Semantique NAF

**Commit** : `2377af9` - feat(sirene): add NAF semantic search with reference table and API endpoints

## Fichiers crees/modifies

| Fichier | Statut | Lignes |
|---------|--------|--------|
| `csv/naf_loader.go` | Cree (A) | 67 |
| `cli/handlers/naf_handler.go` | Cree (A) | 14 |
| `csv/indexer.go` | Modifie (M) | +2 lignes |
| `cli/cli.go` | Modifie (M) | +2 lignes |
| `cli/handlers/help_handler.go` | Modifie (M) | +1 ligne |

## Table naf_reference

| Colonne | Type | Contrainte |
|---------|------|------------|
| `code` | TEXT | PRIMARY KEY |
| `label` | TEXT | NOT NULL |
| `section_code` | TEXT | NOT NULL |
| `section_label` | TEXT | NOT NULL |

Creation via `CREATE TABLE IF NOT EXISTS` dans `csv/naf_loader.go`.
Donnees source : `data/naf_codes.json` (732 codes NAF, source OpenDataSoft INSEE).
Strategie de chargement : TRUNCATE puis INSERT ligne par ligne.

## Indexes ajoutes

| Nom | Table | Type | Definition |
|-----|-------|------|------------|
| `idx_naf_ref_label_trgm` | `naf_reference` | GIN trigram | `USING gin(label gin_trgm_ops)` |
| `idx_naf_ref_section` | `naf_reference` | B-tree | `(section_code)` |

Ces 2 indexes sont declares dans `csv/indexer.go` (lignes 27-28) et sont crees par la commande CLI `indexes`.

## Commande CLI ajoutee

```
sirene-api naf
```

- Enregistree dans `cli/cli.go` : `case "naf"` -> `handlers.HandleImportNaf(c.db)`
- Handler dans `cli/handlers/naf_handler.go` : appelle `csv.LoadNafCodes(db, "data/naf_codes.json")`
- Documentee dans `cli/handlers/help_handler.go` : `naf  Importer les codes NAF depuis data/naf_codes.json`

## Bilan

| Metrique | Valeur |
|----------|--------|
| Fichiers crees | 2 |
| Fichiers modifies | 3 |
| Lignes de code ajoutees | ~86 |
| Table creee | `naf_reference` (4 colonnes) |
| Indexes ajoutes | 2 (1 GIN trigram + 1 B-tree) |
| Commande CLI | `naf` |
