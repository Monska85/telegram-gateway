#!/usr/bin/env bash

# Always operate from the directory where this script lives,
# so docker compose finds docker-compose.yaml (and any override) correctly.
cd "$(dirname "$(readlink -f "${0}")")" || exit 1

# Pull the latest images
echo "Pulling the latest images..."
for image in $(docker compose config | yq -r '.services | .[].image'); do
  echo "Pull image: ${image}"
  docker pull "${image}"
done
echo "Images pulled successfully."

# Restart the containers
docker compose down && docker compose up -d
