#!/bin/bash

image_name="telegram-youtube-briefly"
container_name="telegram-youtube-briefly"

docker build --tag "$image_name" --file Dockerfile.Telegram .
docker stop "$container_name" &>/dev/null || true
docker rm "$container_name" &>/dev/null || true

docker run -d \
    --env-file .env \
    --name "$container_name" \
    "$image_name"

docker image prune -f