---
name: step-03-review
description: Phase review - validation par agent code-reviewer avec boucle de correction
next_step: steps/step-04-finalize.md
---

# Etape 3 : Review et correction

## Regles d'execution

- L'agent `code-reviewer` ne modifie JAMAIS de fichier
- Si des corrections sont necessaires, relancer un `Snipper`
- Maximum 3 iterations de correction avant escalade humaine
- La compilation et le lint DOIVENT passer

## Sequence

### 1. Lancer la review

Utiliser le Task tool avec :

```
subagent_type: "code-reviewer"
prompt: |
  Valider l'ensemble des changements pour la feature {feature_id}.

  Fichiers crees/modifies :
  [liste des fichiers des phases 1 et 2]

  Verifier :
  1. go build ./... compile sans erreur
  2. go vet ./... sans warning
  3. Pas d'injection SQL (parametres positionnels $1, $2)
  4. Erreurs gerees ou explicitement ignorees (_ = ...)
  5. Cache Redis coherent (cles, TTL)
  6. Patterns existants respectes
  7. Pas de commentaires dans le code

  Produire un rapport avec verdict APPROUVE ou BLOQUE.
```

### 2. Analyser le verdict

**Si APPROUVE :**
Passer a l'etape 4 (finalisation)

**Si BLOQUE :**
Analyser les actions requises, relancer un `Snipper` avec les corrections.

### 3. Boucle de correction

```
Tant que verdict == BLOQUE et review_iterations < 3 :
    1. Identifier les fichiers a corriger
    2. Si fichier dans csv/ ou cli/ → Snipper avec perimetre DB
    3. Si fichier dans api/ → Snipper avec perimetre API
    4. Recompiler : go build ./...
    5. Relancer code-reviewer
    6. review_iterations++

Si review_iterations >= 3 et toujours BLOQUE :
    Presenter le rapport a l'utilisateur
    Demander comment proceder
```

### 4. Point de validation finale

**Si `{validate_mode}` :**
Presenter le rapport de review a l'utilisateur.

**Si `{auto_mode}` :**
Si APPROUVE, passer directement a la finalisation.

### 5. Sauvegarder (si save_mode)

Ecrire le rapport dans `.claude/output/feature-builder/{feature-id}/03-review.md`.

### 6. Passer a l'etape suivante

Charger `steps/step-04-finalize.md`.
