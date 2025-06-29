#!/bin/bash

image_name="telegram-youtube-briefly"

docker build --tag "$image_name" --file Dockerfile.Telegram .
