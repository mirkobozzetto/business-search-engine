---
name: step-03-review
description: Phase review - validation complète par l'agent api-reviewer avec boucle de correction
next_step: steps/step-04-finalize.md
---

# Étape 3 : Review et correction

## Règles d'exécution

- L'agent `api-reviewer` ne modifie JAMAIS de fichier
- Si des corrections sont nécessaires, relancer l'agent approprié
- Maximum 3 itérations de correction avant escalade humaine
- La compilation et le lint DOIVENT passer

## Séquence

### 1. Lancer l'agent api-reviewer

Utiliser le Task tool avec :

```
subagent_type: "api-reviewer"
prompt: |
  Valider l'ensemble des changements effectués pour la feature {feature_id}.

  Fichiers créés/modifiés :
  [liste des fichiers des phases 1 et 2]

  Exécuter :
  1. go build ./...
  2. go vet ./...
  3. Review de chaque fichier modifié

  Produire un rapport avec verdict APPROUVÉ ou BLOQUÉ.
```

### 2. Analyser le verdict

**Si APPROUVÉ :**
→ Passer à l'étape 4 (finalisation)

**Si BLOQUÉ :**
→ Analyser les actions requises
→ Identifier l'agent responsable (db-architect ou api-builder)
→ Relancer l'agent avec les corrections demandées
→ Incrémenter `{review_iterations}`
→ Relancer le reviewer

### 3. Boucle de correction

```
Tant que verdict == BLOQUÉ et review_iterations < 3 :
    1. Identifier les fichiers à corriger
    2. Si fichier dans csv/ ou cli/ → relancer db-architect
    3. Si fichier dans api/ → relancer api-builder
    4. Recompiler : go build ./...
    5. Relancer api-reviewer
    6. review_iterations++

Si review_iterations >= 3 et toujours BLOQUÉ :
    → Présenter le rapport à l'utilisateur
    → Demander comment procéder
```

### 4. Point de validation finale

**Si `{validate_mode}` :**
Présenter le rapport de review à l'utilisateur.

```
Review terminée après {review_iterations} itération(s).

Verdict : APPROUVÉ

Compilation : OK
Architecture : OK
Sécurité SQL : OK
Gestion erreurs : OK
Cache : OK
Conventions : OK

Passer à la finalisation ?
```

**Si `{auto_mode}` :**
Si APPROUVÉ, passer directement à la finalisation.

### 5. Sauvegarder (si save_mode)

Écrire le rapport dans `.claude/output/feature-builder/{feature-id}/03-review.md`.

### 6. Passer à l'étape suivante

Charger `steps/step-04-finalize.md`.
