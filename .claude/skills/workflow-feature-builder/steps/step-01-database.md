---
name: step-01-database
description: Phase DB - créer les tables, indexes, loaders et commandes CLI via l'agent db-architect
next_step: steps/step-02-api.md
---

# Étape 1 : Base de données

## Règles d'exécution

- Utiliser UNIQUEMENT l'agent `db-architect` pour cette phase
- Chaque sous-tâche DB est exécutée séquentiellement
- Vérifier la compilation après chaque sous-tâche
- Ne JAMAIS toucher aux fichiers `api/`

## Séquence

### 1. Lancer l'agent db-architect

Utiliser le Task tool avec :

```
subagent_type: "db-architect"
prompt: [description des tâches DB de la phase 1]
```

Le prompt doit inclure :

- La liste exacte des sous-tâches DB du plan
- Les fichiers cibles pour chaque sous-tâche
- Les conventions du projet à respecter
- Le schéma des nouvelles tables si applicable

### 2. Vérifier la compilation

Après le retour de l'agent :

```bash
cd sirene_france_backend && go build ./...
```

Si erreur : relancer l'agent avec le message d'erreur.

### 3. Mettre à jour le statut

Marquer les sous-tâches DB comme terminées :

```
Phase 1 : Base de données
  ■ 1.1 : [Terminé]
  ■ 1.2 : [Terminé]
```

### 4. Point de validation

**Si `{validate_mode}` :**
Présenter un résumé des fichiers créés/modifiés et demander validation.

```
Phase 1 terminée. Fichiers modifiés :
  + csv/naf_loader.go (nouveau)
  + csv/indexer.go (modifié - 3 indexes ajoutés)
  + cli/handlers/naf_handler.go (nouveau)
  ~ cli/cli.go (modifié - 2 commandes ajoutées)

Approuver pour passer à la phase API ?
```

**Si `{auto_mode}` :**
Lancer l'agent `api-reviewer` pour une validation rapide.
Si APPROUVÉ : passer à l'étape suivante.
Si BLOQUÉ : relancer `db-architect` avec les corrections.
Maximum 3 itérations.

### 5. Sauvegarder (si save_mode)

Écrire le rapport dans `.claude/output/feature-builder/{feature-id}/01-db-changes.md`.

### 6. Passer à l'étape suivante

Charger `steps/step-02-api.md`.
