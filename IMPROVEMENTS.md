# Améliorations apportées à Scotter

Ce document récapitule les améliorations apportées au projet Scotter suite aux tests d'intégration réalisés avec un projet généré par Scotter.

## Problèmes identifiés et corrections apportées

### 1. Version de Go incorrecte dans go.mod

**Problème** : Le fichier `go.mod` référençait une version de Go qui n'existe pas (`1.23.4`).

**Solution** : Mise à jour vers la version `1.21` qui est compatible avec les workflows GitHub Actions.

### 2. Gestion des tokens dans le workflow de release

**Problème** : Le workflow de release utilisait le token GitHub par défaut, qui peut avoir des limitations pour certaines actions.

**Solution** : Utilisation d'un token personnalisé `RELEASE_TOKEN` à configurer dans les secrets du dépôt GitHub.

### 3. Génération de changelog incomplète

**Problème** : Le changelog généré par GoReleaser était minimal et ne montrait pas tous les commits pertinents.

**Solution** : Ajout de l'option `use: git` dans la configuration du changelog de GoReleaser pour une génération plus complète.

## Améliorations de la génération de projets

### 1. Syntaxe des templates GoReleaser

La génération des projets doit s'assurer que les variables de template GoReleaser utilisent la syntaxe correcte :
- Utiliser `{{ .Version }}`, `{{ .Os }}`, `{{ .Arch }}`, etc. au lieu de `${VERSION}`, `${OS}`, `${ARCH}`

### 2. Fichiers nécessaires pour GoReleaser

S'assurer que les projets générés incluent tous les fichiers nécessaires pour GoReleaser :
- Présence d'un fichier `LICENSE`
- Présence d'un répertoire `docs` avec au moins un fichier

### 3. Templates de workflows GitHub Actions

Les templates de workflows doivent être générés en s'assurant que :
- Les conditionnels Go template sont correctement traités
- Aucun template non-traité n'est présent dans les fichiers YAML finaux

## Amélioration des commits conventionnels

Pour une meilleure génération automatique de changelogs, les commits doivent suivre les conventions :
- `feat:` pour les nouvelles fonctionnalités
- `fix:` pour les corrections de bugs
- `docs:` pour les mises à jour de documentation
- `perf:` pour les améliorations de performances
- etc.

## Prochaines étapes recommandées

1. Améliorer la détection et le remplacement des templates non-traités dans les fichiers générés
2. Ajouter une validation automatique des fichiers générés pour s'assurer qu'ils sont conformes
3. Automatiser la création des fichiers requis par GoReleaser si nécessaire
4. Ajouter plus de tests d'intégration pour vérifier le bon fonctionnement des workflows
