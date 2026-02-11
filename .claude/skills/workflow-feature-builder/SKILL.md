---
name: feature-builder
description: Workflow multi-agent pour implémenter des features complètes avec décomposition en tâches, validation par étape et coordination d'agents spécialisés. Utiliser pour toute feature qui touche plusieurs couches (DB, API, modèles).
argument-hint: "[-i] [-a] [-v] [-s] [-avs] <description de la feature>"
---

<objective>
Orchestrer l'implémentation de features complètes via un système multi-agent avec décomposition en tâches, dépendances, et validation entre chaque phase. Le workflow est auto-maintenable : il peut évoluer en ajoutant de nouvelles étapes ou en modifiant les agents.
</objective>

<quick_start>

```bash
/feature-builder -i ajouter la recherche sémantique NAF

/feature-builder -avs ajouter un endpoint d'export CSV

/feature-builder -vs créer le système d'enrichissement CRM
```

**Flags combinables :**
- `-i` : mode interactif, affiche tous les flags et permet de choisir
- `-a` : mode autonome, validation automatique par l'agent reviewer
- `-v` : validation humaine obligatoire entre chaque phase
- `-s` : sauvegarder les outputs dans `.claude/output/feature-builder/`

Les flags se combinent : `-avs` = auto + validate + save, `-vs` = validate + save.

</quick_start>

<parameters>

<flags>
| Court | Long | Description |
|-------|------|-------------|
| `-i` | `--interactive` | Affiche tous les flags disponibles et permet de les activer/désactiver avant de lancer |
| `-a` | `--auto` | Mode autonome : l'agent reviewer valide automatiquement |
| `-v` | `--validate` | Validation humaine obligatoire entre chaque phase |
| `-s` | `--save` | Sauvegarder les outputs dans `.claude/output/feature-builder/` |

**Combinaison de flags :**
Les flags courts se combinent en un seul bloc :
- `-avs` équivaut à `-a -v -s`
- `-vs` équivaut à `-v -s`
- `-as` équivaut à `-a -s`
- `-i` se lance seul (les autres flags sont choisis interactivement)

**Exemples :**
```bash
/feature-builder -i ajouter un endpoint          # Interactif : choisir les flags
/feature-builder -avs implémenter la feature     # Auto + Validate + Save
/feature-builder -vs ajouter la table NAF        # Validate + Save
/feature-builder -a corriger le bug              # Auto seul
/feature-builder ajouter un champ                # Défaut : validate_mode = true
```
</flags>

<defaults>
```yaml
auto_mode: false
validate_mode: true
save_mode: false
```
</defaults>

</parameters>

<agents>

Ce workflow utilise les agents integres de Claude Code :

| Agent | Type | Role |
|-------|------|------|
| **DB** | `Snipper` | Tables, indexes, loaders, CLI (perimetre : `csv/`, `cli/`, `database/`) |
| **API** | `Snipper` | Services, handlers, routes, modeles (perimetre : `api/`) |
| **Review** | `code-reviewer` | Validation, compilation, review (lecture seule) |

Les perimetres de fichiers sont definis dans le prompt de chaque agent pour eviter les conflits.

</agents>

<workflow>

Le workflow se déroule en 5 phases séquentielles :

```
Phase 0 : Initialisation
    → Analyser la feature, décomposer en tâches
    → Identifier les agents nécessaires
    → Créer le plan avec dépendances

Phase 1 : Base de donnees (Snipper)
    → Creer les tables et indexes
    → Creer les loaders de donnees
    → Ajouter les commandes CLI
    → POINT DE VALIDATION

Phase 2 : API et services (Snipper)
    → Creer les modeles
    → Creer les services
    → Creer les handlers et routes
    → Enrichir les requetes existantes
    → POINT DE VALIDATION

Phase 3 : Review et correction (code-reviewer → Snipper si corrections)
    → Verifier compilation
    → Verifier coherence architecturale
    → Verifier securite SQL
    → Si BLOQUE : renvoyer au Snipper pour correction
    → POINT DE VALIDATION

Phase 4 : Finalisation
    → Résumé des changements
    → Liste des nouveaux endpoints
    → Commandes CLI disponibles
    → Proposition de commit
```

</workflow>

<validation_gates>

Chaque point de validation fonctionne selon le mode choisi :

**Mode `-v` (validate, par défaut) :**
1. L'agent reviewer produit son rapport
2. Le rapport est présenté à l'utilisateur
3. L'utilisateur approuve ou demande des corrections
4. Si corrections : les agents concernés sont relancés
5. Retour au point de validation

**Mode `-a` (auto) :**
1. L'agent reviewer produit son rapport
2. Si APPROUVÉ : passage automatique à la phase suivante
3. Si BLOQUÉ : les agents concernés corrigent automatiquement
4. Maximum 3 itérations de correction avant escalade à l'utilisateur

**Sans flag :**
1. L'agent reviewer produit son rapport
2. L'utilisateur est consulté pour décider de la suite

</validation_gates>

<task_decomposition>

Chaque feature est décomposée selon cette hiérarchie :

```
Feature
├── Phase 1 : Tâches DB
│   ├── Sous-tâche 1.1 : Schéma de table
│   ├── Sous-tâche 1.2 : Indexes
│   ├── Sous-tâche 1.3 : Loader de données
│   └── Sous-tâche 1.4 : Commande CLI
├── Phase 2 : Tâches API
│   ├── Sous-tâche 2.1 : Modèles
│   ├── Sous-tâche 2.2 : Services
│   ├── Sous-tâche 2.3 : Handlers
│   └── Sous-tâche 2.4 : Routes
├── Phase 3 : Validation
│   └── Sous-tâche 3.1 : Review complète
└── Phase 4 : Finalisation
    └── Sous-tâche 4.1 : Résumé et commit
```

Les sous-tâches d'une même phase peuvent être parallélisées si elles ne touchent pas les mêmes fichiers.

</task_decomposition>

<state_variables>

| Variable | Type | Description |
|----------|------|-------------|
| `{feature_description}` | string | Description de la feature |
| `{feature_id}` | string | Identifiant kebab-case |
| `{interactive_mode}` | boolean | Choix interactif des flags |
| `{auto_mode}` | boolean | Validation automatique |
| `{validate_mode}` | boolean | Validation humaine |
| `{save_mode}` | boolean | Sauvegarde des outputs |
| `{phase}` | int | Phase courante (0-4) |
| `{tasks}` | list | Liste des tâches avec statut |
| `{review_iterations}` | int | Nombre d'itérations de correction |
| `{output_dir}` | string | Répertoire de sortie |

</state_variables>

<output_structure>

Quand `{save_mode}` = true :

```
.claude/output/feature-builder/{feature-id}/
├── 00-plan.md          # Décomposition en tâches
├── 01-db-changes.md    # Rapport phase DB
├── 02-api-changes.md   # Rapport phase API
├── 03-review.md        # Rapport de review
└── 04-summary.md       # Résumé final
```

</output_structure>

<self_maintenance>

Ce workflow est conçu pour évoluer :

1. **Ajouter une phase** : creer un fichier `steps/step-0N-*.md` et l'inserer dans le workflow
3. **Modifier la validation** : ajuster les seuils dans `<validation_gates>`
4. **Ajouter des sous-tâches** : étendre `<task_decomposition>` selon les besoins

Le skill se met à jour en éditant ce fichier `SKILL.md` et les étapes dans `steps/`.

</self_maintenance>

<entry_point>

**PREMIÈRE ACTION :** Charger `steps/step-00-init.md`

</entry_point>
