#!/bin/bash

image_name="ghcr.io/olegshulyakov/go-briefly-bot"
container_name="go-briefly-telegram"

docker build -t "$image_name" --file .devops/Telegram.Dockerfile .
docker stop "$container_name" &>/dev/null || true
docker rm "$container_name" &>/dev/null || true

docker run -d \
    --env-file .env \
    --volume "./data:/app/data" \
    --name "$container_name" \
    "$image_name"

docker image prune -f
