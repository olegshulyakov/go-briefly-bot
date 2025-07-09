#!/bin/bash

image_name="go-briefly-telegram"

docker build -t "$image_name" --file .build/Telegram.Dockerfile .
