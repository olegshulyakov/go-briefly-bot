#!/bin/bash

image_name="ghcr.io/olegshulyakov/go-briefly-bot"

docker build -t "$image_name" --file .devops/Telegram.Dockerfile .
