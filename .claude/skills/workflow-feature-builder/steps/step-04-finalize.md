---
name: step-04-finalize
description: Finalisation - résumé des changements, proposition de commit, documentation des endpoints
---

# Étape 4 : Finalisation

## Règles d'exécution

- Ne plus modifier de fichier
- Produire un résumé clair et complet
- Proposer un message de commit cohérent
- Attendre la confirmation de l'utilisateur avant de commit

## Séquence

### 1. Inventaire des changements

Lister tous les fichiers créés et modifiés avec `git status` et `git diff --stat`.

### 2. Résumé de la feature

Présenter un résumé structuré :

```
╔══════════════════════════════════════════╗
║  Feature terminée : {feature_id}         ║
╠══════════════════════════════════════════╣
║  Fichiers créés   : N                    ║
║  Fichiers modifiés : N                   ║
║  Itérations review : {review_iterations} ║
╚══════════════════════════════════════════╝

### Nouveaux fichiers
  + chemin/fichier.go : description

### Fichiers modifiés
  ~ chemin/fichier.go : description du changement

### Nouveaux endpoints API
  GET /api/...  → description
  GET /api/...  → description

### Nouvelles commandes CLI
  sirene-api commande  → description

### Tables/indexes ajoutés
  TABLE nom_table (N colonnes)
  INDEX idx_nom ON table(colonnes)
```

### 3. Proposition de commit

Proposer un message de commit au format conventionnel :

```
feat(sirene): {description courte de la feature}
```

### 4. Attendre la décision de l'utilisateur

- Commit maintenant ?
- Modifier quelque chose avant ?
- Lancer l'API pour tester ?

Ne JAMAIS commit sans validation explicite de l'utilisateur.

### 5. Sauvegarder (si save_mode)

Écrire le résumé dans `.claude/output/feature-builder/{feature-id}/04-summary.md`.

### 6. Fin du workflow

Le workflow est terminé. Afficher :

```
Workflow feature-builder terminé pour : {feature_id}
```
