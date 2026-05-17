# Guide du workflow Feature Builder

## Qu'est-ce que c'est ?

Feature Builder est un système multi-agent qui automatise l'implémentation de features complètes sur le projet Business Search Engine. Il coordonne trois agents spécialisés qui travaillent chacun sur leur périmètre, avec des points de validation entre chaque phase pour garantir la qualité.

---

## Prérequis

### 1. Mode multi-agent activé

Vérifier que le setting global est en place :

```json
// ~/.claude/settings.json
{
  "env": {
    "CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS": "1"
  }
}
```

### 2. Fichiers du workflow

Tout est dans le répertoire du projet :

```
.claude/
├── agents/
│   ├── db-architect.md       # Agent base de données
│   ├── api-builder.md        # Agent API
│   └── api-reviewer.md       # Agent de review
│
└── skills/
    └── workflow-feature-builder/
        ├── SKILL.md               # Définition du skill
        └── steps/
            ├── step-00-init.md        # Initialisation
            ├── step-01-database.md    # Phase DB
            ├── step-02-api.md         # Phase API
            ├── step-03-review.md      # Phase review
            └── step-04-finalize.md    # Finalisation
```

### 3. Backend fonctionnel

- Go installé avec `go build ./...` qui compile
- PostgreSQL accessible
- Redis accessible (pour le cache)

---

## Lancer le workflow

### Commande de base

```bash
/feature-builder ajouter la recherche sémantique NAF
```

Par défaut, le workflow demande une validation humaine entre chaque phase. Claude analyse la feature, crée un plan, puis exécute phase par phase en te demandant ton accord avant de continuer.

### Avec les flags

```bash
# Mode interactif : choisir les flags avant de lancer
/feature-builder -i ajouter la recherche sémantique NAF

# Mode autonome : Claude valide tout seul via l'agent reviewer
/feature-builder -a ajouter la recherche sémantique NAF

# Mode validation stricte : tu dois approuver chaque phase
/feature-builder -v ajouter la recherche sémantique NAF

# Mode avec sauvegarde des rapports
/feature-builder -s ajouter la recherche sémantique NAF

# Combiner les flags en un seul bloc
/feature-builder -avs ajouter la recherche sémantique NAF
/feature-builder -vs ajouter la recherche sémantique NAF
/feature-builder -as ajouter la recherche sémantique NAF
```

| Flag | Effet |
| ---- | ----- |
| `-i` | Mode interactif : affiche tous les flags et te laisse choisir avant de lancer |
| `-a` | L'agent reviewer valide automatiquement. Si un problème est détecté, les agents corrigent sans intervention (max 3 tentatives avant de te demander) |
| `-v` | Tu dois approuver explicitement chaque phase. C'est le mode par défaut |
| `-s` | Les rapports de chaque phase sont sauvegardés dans `.claude/output/feature-builder/` |

Les flags courts se combinent : `-avs` = auto + validate + save, `-vs` = validate + save, `-as` = auto + save. Le flag `-i` se lance seul car les autres sont choisis interactivement.

---

## Ce qui se passe quand tu lances le workflow

### Phase 0 : Initialisation

Claude analyse ta demande et crée un plan détaillé :

```
╔══════════════════════════════════════╗
║  Feature : recherche-semantique-naf  ║
╠══════════════════════════════════════╣
║  Auto mode    : false                ║
║  Validation   : true                 ║
║  Save mode    : false                ║
╠══════════════════════════════════════╣
║  Phases       : 4                    ║
║  Tâches       : 8                    ║
║  Agents       : 3                    ║
╚══════════════════════════════════════╝

Phase 1 : Base de données
  □ 1.1 : Créer la table naf_reference
  □ 1.2 : Ajouter les indexes
  □ 1.3 : Créer le loader JSON
  □ 1.4 : Ajouter la commande CLI naf-load

Phase 2 : API
  □ 2.1 : Créer le service NafService
  □ 2.2 : Créer le handler NafHandler
  □ 2.3 : Ajouter les routes
  □ 2.4 : Enrichir CompanyResult avec le label NAF
```

Tu approuves ou tu modifies le plan.

### Phase 1 : Base de données

L'agent **db-architect** prend le relais. Il crée les tables, les indexes, les loaders de données et les commandes CLI.

Périmètre de fichiers : `csv/`, `cli/`, `database/`

A la fin de cette phase, un point de validation te montre ce qui a été créé :

```
Phase 1 terminée. Fichiers modifiés :
  + csv/naf_loader.go (nouveau)
  ~ csv/indexer.go (3 indexes ajoutés)
  + cli/handlers/naf_handler.go (nouveau)
  ~ cli/cli.go (2 commandes ajoutées)

Approuver pour passer à la phase API ?
```

### Phase 2 : API et services

L'agent **api-builder** crée les services, handlers, routes et modèles.

Périmètre de fichiers : `api/`

A la fin :

```
Phase 2 terminée. Changements :
  + api/services/naf/naf_service.go (nouveau)
  + api/handlers/naf_handler.go (nouveau)
  ~ api/server.go (4 routes ajoutées)
  ~ api/models/company_models.go (2 champs ajoutés)

Nouveaux endpoints :
  GET /api/naf/sections
  GET /api/naf/search?q={texte}
  GET /api/naf/section/:code
  GET /api/naf/code/:code

Approuver pour passer à la review ?
```

### Phase 3 : Review

L'agent **api-reviewer** vérifie tout :

- Compilation (`go build ./...`)
- Vétérinaire (`go vet ./...`)
- Cohérence architecturale
- Sécurité SQL (pas d'injection)
- Gestion des erreurs
- Conventions du projet

Il produit un rapport avec un verdict :

```
## Rapport de validation

Compilation    : OK
Architecture   : OK
Sécurité SQL   : OK
Erreurs        : OK
Cache          : OK
Conventions    : OK

### Verdict
APPROUVÉ
```

Si le verdict est **BLOQUÉ**, les agents concernés sont relancés pour corriger les problèmes. Maximum 3 itérations avant que le workflow te demande d'intervenir.

### Phase 4 : Finalisation

Un résumé complet est affiché avec tous les changements, les nouveaux endpoints, et une proposition de message de commit. Le commit n'est jamais fait automatiquement : tu dois donner ton accord.

---

## Les agents utilises

Le workflow utilise les agents integres de Claude Code, pas d'agents custom. Les contraintes de perimetre sont passees dans le prompt du Task tool.

| Phase | Agent integre | Perimetre |
| ----- | ------------- | --------- |
| DB | `Snipper` | `csv/`, `cli/`, `database/` |
| API | `Snipper` | `api/` |
| Review | `code-reviewer` | Lecture seule |

Les perimetres sont isoles pour eviter que deux agents ecrivent dans le meme fichier.

---

## Modes de validation

### Mode par défaut (aucun flag, ou `-v`)

```
Phase 0 → tu approuves le plan
Phase 1 → Snipper (DB) travaille → tu approuves
Phase 2 → Snipper (API) travaille → tu approuves
Phase 3 → code-reviewer verifie → tu approuves
Phase 4 → tu décides de commit ou pas
```

Tu gardes le contrôle total. A chaque étape, tu peux :

- Approuver et passer à la suite
- Demander des modifications
- Arrêter le workflow

### Mode autonome (`-a`)

```
Phase 0 → plan créé automatiquement
Phase 1 → Snipper (DB) travaille → code-reviewer valide automatiquement
Phase 2 → Snipper (API) travaille → code-reviewer valide automatiquement
Phase 3 → code-reviewer verifie → approuve automatiquement si OK
Phase 4 → tu décides de commit ou pas (toujours manuel)
```

Les agents se corrigent entre eux. Si après 3 tentatives un problème persiste, le workflow t'escalade la situation.

Le commit reste toujours à ta discrétion.

---

## Sauvegardes (`-s`)

Quand le flag `-s` est activé, chaque phase produit un rapport sauvegardé :

```
.claude/output/feature-builder/recherche-semantique-naf/
├── 00-plan.md          # Plan avec toutes les tâches
├── 01-db-changes.md    # Rapport phase DB
├── 02-api-changes.md   # Rapport phase API
├── 03-review.md        # Rapport de review
└── 04-summary.md       # Résumé final
```

Ces fichiers servent de documentation et de référence pour les features suivantes.

---

## Exemples d'utilisation concrets

### Recherche sémantique NAF

```bash
/feature-builder -v -s ajouter la recherche sémantique NAF avec table de référence et enrichissement des résultats
```

Ce que le workflow va faire :

1. Créer la table `naf_reference` avec les 732 codes
2. Créer un loader qui lit `data/naf_codes.json`
3. Ajouter les indexes de recherche (trigram pour ILIKE)
4. Créer les endpoints `/api/naf/search`, `/api/naf/sections`, etc.
5. Enrichir `CompanyResult` avec le libellé NAF
6. Valider, résumer, proposer le commit

### Table d'enrichissement CRM

```bash
/feature-builder -v créer la table company_enrichment pour le CRM avec champs website, phone, linkedin, contacts JSONB
```

### Export CSV

```bash
/feature-builder -a ajouter un endpoint d'export CSV pour les résultats de recherche multi-critères
```

---

## Faire evoluer le workflow

### Ajouter une phase

1. Créer le fichier `steps/step-0N-ma-phase.md` avec :

```yaml
---
name: step-0N-ma-phase
description: Ce que fait cette phase
next_step: steps/step-0M-suivante.md
---
```

2. Mettre à jour le `next_step` de l'étape précédente pour pointer vers la nouvelle

### Modifier les seuils de validation

Dans `steps/step-03-review.md`, la variable `review_iterations` contrôle le nombre maximum de corrections automatiques avant escalade humaine. Par défaut : 3.

---

## Architecture technique rappel

```
sirene_france_backend/
├── api/
│   ├── cache/          # Client Redis + compression
│   ├── handlers/       # Handlers HTTP (Gin)
│   ├── models/         # Structs de données
│   ├── services/       # Logique métier
│   │   ├── company/    # Recherches entreprises
│   │   └── naf/        # Recherches NAF (à créer)
│   └── server.go       # Routes + middleware
├── cli/
│   ├── cli.go          # Dispatch des commandes
│   └── handlers/       # Logique CLI
├── csv/
│   ├── indexer.go      # Création des indexes
│   ├── pipeline.go     # Import parallèle
│   └── zip_reader.go   # Lecture des ZIP INSEE
├── config/             # Variables d'environnement
├── data/
│   └── naf_codes.json  # 732 codes NAF (déjà en place)
├── database/           # Connexion PostgreSQL
└── main.go             # Point d'entrée
```

### Tables existantes

| Table           | Lignes | Usage                       |
| --------------- | ------ | --------------------------- |
| `etablissement` | ~36M   | Données des établissements  |
| `unite_legale`  | ~24M   | Données des entités légales |

### Tables à créer

| Table                | Lignes     | Usage                                                 |
| -------------------- | ---------- | ----------------------------------------------------- |
| `naf_reference`      | 732        | Référentiel des codes NAF                             |
| `company_enrichment` | 0 (future) | Données CRM : site web, téléphone, LinkedIn, contacts |

### Indexes existants

14 indexes dont :

- `idx_etab_naf` sur `activite_principale_etablissement`
- `idx_etab_commune_trgm` et `idx_ul_denom_trgm` en GIN trigram
- 3 indexes partiels avec `WHERE etablissement_siege = 'true'`

---

## Résolution de problèmes

### Le skill ne se lance pas

Vérifier que le fichier `SKILL.md` est bien dans `.claude/skills/workflow-feature-builder/SKILL.md` et que la propriété `name` correspond.

### Un agent échoue à compiler

Le workflow relance l'agent avec le message d'erreur. Si le problème persiste après 3 tentatives, il te demande d'intervenir. Tu peux alors corriger manuellement et relancer la phase suivante.

### Conflit de fichiers entre agents

Impossible par design : chaque agent a son périmètre de fichiers défini dans son `.md`. Si un fichier doit être touché par deux agents, il faut le déplacer dans le périmètre d'un seul agent et adapter l'autre.

### Le mode autonome boucle

Après 3 itérations de correction sans résolution, le workflow s'arrête et te donne la main. Le rapport du reviewer indique exactement ce qui bloque.
