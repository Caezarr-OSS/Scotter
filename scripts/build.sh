#!/bin/bash

# Script pour compiler Scotter avec les informations de version appropriées
# Règles de versionnement:
# 1. Si un tag Git existe, utiliser ce tag comme version
# 2. Si on est sur une branche autre que main/master, utiliser le nom de la branche
# 3. Pour les builds locaux sans tag, utiliser un timestamp

# Répertoire du script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# Répertoire racine du projet (un niveau au-dessus du script)
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

# Se déplacer dans le répertoire du projet
cd "$PROJECT_DIR" || exit 1

# Récupérer le commit SHA
COMMIT_SHA=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Récupérer la date de build au format ISO
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Déterminer la version
# 1. Vérifier si nous sommes sur un tag
GIT_TAG=$(git describe --tags --exact-match 2>/dev/null)
if [ -n "$GIT_TAG" ]; then
    # Nous sommes sur un tag, utiliser ce tag comme version
    VERSION=${GIT_TAG#v}  # Enlever le 'v' préfixe si présent pour les bibliothèques
    echo "Building tagged version: $VERSION"
else
    # 2. Vérifier si nous sommes sur une branche
    GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null)
    if [ -n "$GIT_BRANCH" ] && [ "$GIT_BRANCH" != "HEAD" ]; then
        # Nous sommes sur une branche, utiliser son nom comme version
        VERSION="${GIT_BRANCH}-${COMMIT_SHA}"
        echo "Building branch version: $VERSION"
    else
        # 3. Build local, utiliser un timestamp
        TIMESTAMP=$(date -u +"%Y%m%d%H%M%S")
        VERSION="dev-${TIMESTAMP}"
        echo "Building local version: $VERSION"
    fi
fi

# Compiler le binaire avec les informations de version
echo "Compiling Scotter..."
go build -ldflags "-X main.Version=${VERSION} -X main.CommitSHA=${COMMIT_SHA} -X main.BuildDate=${BUILD_DATE}" -o scotter .

echo "Build completed: ./scotter"
echo "Version: ${VERSION}"
echo "Commit: ${COMMIT_SHA}"
echo "Build Date: ${BUILD_DATE}"
