#!/bin/bash

image_name="go-briefly-telegram"
container_name="go-briefly-telegram"

docker build -t "$image_name" --file .build/Telegram.Dockerfile .
docker stop "$container_name" &>/dev/null || true
docker rm "$container_name" &>/dev/null || true

docker run -d \
    --env-file .env \
    --name "$container_name" \
    "$image_name"

docker image prune -f