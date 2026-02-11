---
name: step-02-api
description: Phase API - créer les services, handlers, routes et modèles via l'agent api-builder
next_step: steps/step-03-review.md
---

# Étape 2 : API et services

## Règles d'exécution

- Utiliser UNIQUEMENT l'agent `api-builder` pour cette phase
- Les tables créées en phase 1 sont disponibles
- Chaque sous-tâche API est exécutée séquentiellement
- Ne JAMAIS toucher aux fichiers `csv/` ou `cli/`

## Séquence

### 1. Lancer l'agent api-builder

Utiliser le Task tool avec :

```
subagent_type: "api-builder"
prompt: [description des tâches API de la phase 2]
```

Le prompt doit inclure :

- La liste exacte des sous-tâches API du plan
- Les fichiers cibles pour chaque sous-tâche
- Le schéma des tables créées en phase 1 (pour les jointures)
- Les endpoints à créer avec leurs paramètres
- Les modèles à modifier ou créer
- Le format des réponses attendues

### 2. Vérifier la compilation

Après le retour de l'agent :

```bash
cd sirene_france_backend && go build ./...
```

Si erreur : relancer l'agent avec le message d'erreur.

### 3. Vérifier le lint

```bash
cd sirene_france_backend && golangci-lint run ./...
```

Si erreurs errcheck ou autres : relancer l'agent pour corriger.

### 4. Mettre à jour le statut

```
Phase 2 : API
  ■ 2.1 : [Terminé]
  ■ 2.2 : [Terminé]
```

### 5. Point de validation

**Si `{validate_mode}` :**
Présenter un résumé avec les nouveaux endpoints et demander validation.

```
Phase 2 terminée. Changements :
  + api/services/naf/naf_service.go (nouveau)
  + api/handlers/naf_handler.go (nouveau)
  ~ api/server.go (modifié - 4 routes ajoutées)
  ~ api/models/company_models.go (modifié - 2 champs ajoutés)

Nouveaux endpoints :
  GET /api/naf/sections
  GET /api/naf/search?q={texte}
  GET /api/naf/section/:code
  GET /api/naf/code/:code

Approuver pour passer à la review ?
```

**Si `{auto_mode}` :**
Passer directement à la review.

### 6. Sauvegarder (si save_mode)

Écrire le rapport dans `.claude/output/feature-builder/{feature-id}/02-api-changes.md`.

### 7. Passer à l'étape suivante

Charger `steps/step-03-review.md`.
