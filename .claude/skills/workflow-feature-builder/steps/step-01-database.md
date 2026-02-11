---
name: step-01-database
description: Phase DB - tables, indexes, loaders et commandes CLI via agent Snipper
next_step: steps/step-02-api.md
---

# Etape 1 : Base de donnees

## Regles d'execution

- Utiliser un agent `Snipper` avec les contraintes DB dans le prompt
- Chaque sous-tache DB est executee sequentiellement
- Verifier la compilation apres chaque sous-tache
- Ne JAMAIS toucher aux fichiers `api/`

## Sequence

### 1. Lancer l'agent

Utiliser le Task tool avec :

```
subagent_type: "Snipper"
prompt: |
  PERIMETRE : uniquement les fichiers dans csv/, cli/, database/
  NE PAS toucher a api/ ni main.go

  Conventions du projet :
  - Colonnes PostgreSQL en TEXT, snake_case
  - Extension pg_trgm pour les recherches ILIKE
  - Pipeline pgx pour l'import parallele
  - Indexes : GIN trigram, partiels avec WHERE
  - Pas de commentaires dans le code

  Taches :
  [liste des sous-taches DB du plan]
```

### 2. Verifier la compilation

```bash
cd sirene_france_backend && go build ./...
```

Si erreur : relancer l'agent avec le message d'erreur.

### 3. Mettre a jour le statut

```
Phase 1 : Base de donnees
  ■ 1.1 : [Termine]
  ■ 1.2 : [Termine]
```

### 4. Point de validation

**Si `{validate_mode}` :**
Presenter un resume des fichiers crees/modifies et demander validation.

**Si `{auto_mode}` :**
Lancer un agent `code-reviewer` pour validation rapide.
Si APPROUVE : passer a l'etape suivante.
Si BLOQUE : relancer le `Snipper` avec les corrections.
Maximum 3 iterations.

### 5. Sauvegarder (si save_mode)

Ecrire le rapport dans `.claude/output/feature-builder/{feature-id}/01-db-changes.md`.

### 6. Passer a l'etape suivante

Charger `steps/step-02-api.md`.
