---
name: step-02-api
description: Phase API - services, handlers, routes et modeles via agent Snipper
next_step: steps/step-03-review.md
---

# Etape 2 : API et services

## Regles d'execution

- Utiliser un agent `Snipper` avec les contraintes API dans le prompt
- Les tables creees en phase 1 sont disponibles
- Ne JAMAIS toucher aux fichiers `csv/` ou `cli/`

## Sequence

### 1. Lancer l'agent

Utiliser le Task tool avec :

```
subagent_type: "Snipper"
prompt: |
  PERIMETRE : uniquement les fichiers dans api/
  NE PAS toucher a csv/, cli/, main.go

  Conventions du projet :
  - Framework Gin, connexion pgx
  - Services dans api/services/{domain}/
  - Handlers dans api/handlers/
  - Routes dans api/server.go
  - Modeles dans api/models/
  - Cache Redis gzip, TTL 1h, cle sirene:v2:{type}:{params}
  - Pagination : limit (defaut 100, max 10000), offset
  - Reponses : models.Success() ou models.Error()
  - SQL : parametres positionnels ($1, $2), COALESCE sur nullable
  - Pas de commentaires dans le code

  Tables disponibles (creees en phase 1) :
  [schema des tables]

  Taches :
  [liste des sous-taches API du plan]
```

### 2. Verifier la compilation

```bash
cd sirene_france_backend && go build ./...
```

Si erreur : relancer l'agent avec le message d'erreur.

### 3. Mettre a jour le statut

```
Phase 2 : API
  ■ 2.1 : [Termine]
  ■ 2.2 : [Termine]
```

### 4. Point de validation

**Si `{validate_mode}` :**
Presenter un resume avec les nouveaux endpoints et demander validation.

**Si `{auto_mode}` :**
Passer directement a la review.

### 5. Sauvegarder (si save_mode)

Ecrire le rapport dans `.claude/output/feature-builder/{feature-id}/02-api-changes.md`.

### 6. Passer a l'etape suivante

Charger `steps/step-03-review.md`.
