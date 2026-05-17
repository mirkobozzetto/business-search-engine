---
name: step-00-init
description: Initialiser le workflow - parser les flags, analyser la feature, décomposer en tâches, créer le plan
next_step: steps/step-01-database.md
---

# Étape 0 : Initialisation

## Règles d'exécution

- Analyser la feature demandée AVANT toute action
- Décomposer en tâches concrètes avec fichiers cibles
- Identifier les dépendances entre tâches
- Ne JAMAIS commencer à coder dans cette étape

## Séquence

### 1. Parser les flags

**Algorithme de parsing :**

1. Scanner l'entrée pour trouver les tokens commençant par `-`
2. Chaque token peut contenir un ou plusieurs flags courts combinés
3. Retirer tous les tokens de flags, le reste devient `{feature_description}`

**Flags combinés :**
Un token comme `-avs` se décompose lettre par lettre :
- `-avs` → `-a` + `-v` + `-s`
- `-vs` → `-v` + `-s`
- `-as` → `-a` + `-s`
- `-i` → `-i` seul

**Table de correspondance :**

```
-i / --interactive  → {interactive_mode} = true
-a / --auto         → {auto_mode} = true
-v / --validate     → {validate_mode} = true
-s / --save         → {save_mode} = true
```

**Valeurs par défaut (si aucun flag) :**

```yaml
interactive_mode: false
auto_mode: false
validate_mode: true
save_mode: false
```

**Exemples de parsing :**

```
Entrée : "-avs ajouter la recherche NAF"
  → auto_mode=true, validate_mode=true, save_mode=true
  → feature_description="ajouter la recherche NAF"

Entrée : "-vs créer la table"
  → validate_mode=true, save_mode=true
  → feature_description="créer la table"

Entrée : "-i implémenter l'export"
  → interactive_mode=true
  → feature_description="implémenter l'export"

Entrée : "ajouter un champ"
  → validate_mode=true (défaut)
  → feature_description="ajouter un champ"
```

Générer `{feature_id}` en kebab-case à partir de `{feature_description}`.

### 2. Mode interactif (si `-i`)

Si `{interactive_mode}` est activé, utiliser AskUserQuestion pour présenter les flags :

```
Question : "Quels modes activer pour ce workflow ?"
Options :
  - Auto (-a) : l'agent reviewer valide automatiquement
  - Validate (-v) : validation humaine entre chaque phase
  - Save (-s) : sauvegarder les rapports de chaque phase
multiSelect: true
```

Appliquer les choix de l'utilisateur aux variables d'état.

Si aucun mode n'est sélectionné, appliquer les défauts (validate_mode=true).

### 3. Explorer le codebase

Lancer un agent `explore-codebase` pour comprendre :

- L'architecture existante
- Les fichiers qui seront impactés
- Les patterns à suivre
- Les dépendances à respecter

### 4. Décomposer en tâches

Pour chaque phase, lister les sous-tâches concrètes :

```
Phase 1 - Base de données (agent: db-architect)
  1.1: [Description] → [fichier(s) cible(s)]
  1.2: [Description] → [fichier(s) cible(s)]

Phase 2 - API (agent: api-builder)
  2.1: [Description] → [fichier(s) cible(s)]
  2.2: [Description] → [fichier(s) cible(s)]

Phase 3 - Review (agent: api-reviewer)
  3.1: Validation complète → lecture seule
```

### 5. Vérifier les conflits de fichiers

Aucun fichier ne doit apparaître dans deux phases différentes.
Si conflit détecté : réorganiser les tâches.

### 6. Présenter le plan

Afficher le plan complet à l'utilisateur :

```
╔══════════════════════════════════════╗
║  Feature Builder : {feature_id}      ║
╠══════════════════════════════════════╣
║  Auto mode    : {auto_mode}          ║
║  Validation   : {validate_mode}      ║
║  Save mode    : {save_mode}          ║
╠══════════════════════════════════════╣
║  Phases       : N                    ║
║  Tâches       : N                    ║
║  Agents       : [liste]              ║
╚══════════════════════════════════════╝

Phase 1 : Base de données
  □ 1.1 : ...
  □ 1.2 : ...

Phase 2 : API
  □ 2.1 : ...
  □ 2.2 : ...

Phase 3 : Review
  □ 3.1 : ...
```

### 7. Attendre la validation

- Si `{validate_mode}` : demander à l'utilisateur s'il approuve le plan
- Si `{auto_mode}` : passer directement à l'étape suivante
- Sinon : demander confirmation

### 8. Créer la structure de sortie (si save_mode)

```bash
mkdir -p .claude/output/feature-builder/{feature-id}
```

Créer `00-plan.md` avec le plan détaillé.

### 9. Passer à l'étape suivante

Charger `steps/step-01-database.md`.
