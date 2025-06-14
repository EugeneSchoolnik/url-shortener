#!/bin/bash

# Проверка аргумента
if [ -z "$1" ]; then
  echo "Usage: $0 <commit_hash>"
  exit 1
fi

COMMIT=$1
OLD_REPO="/home/eugene/workspace/go/pet-projects/url-shortener"
TARGET_DIR="/home/eugene/workspace/go/pet-projects/url-shortener-practice"

# Проверка наличия путей
if [ ! -d "$OLD_REPO/.git" ]; then
  echo "Error: $OLD_REPO is not a valid git repository"
  exit 1
fi

if [ ! -d "$TARGET_DIR" ]; then
  echo "Error: $TARGET_DIR does not exist"
  exit 1
fi

# Перенос файлов
echo "Extracting files from commit $COMMIT into $TARGET_DIR..."
git --git-dir="$OLD_REPO/.git" --work-tree="$TARGET_DIR" checkout "$COMMIT" -- .

# Готово
echo "Done. Files from commit $COMMIT copied to $TARGET_DIR"

