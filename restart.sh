#!/usr/bin/env bash

# Pull the latest images
echo "Pulling the latest images..."
for image in $(yq -r '.services | .[].image' docker-compose.yml); do
  echo "Pull image: ${image}"
  docker pull "${image}"
done
echo "Images pulled successfully."

# Restart the containers
docker compose down && docker compose up -d
