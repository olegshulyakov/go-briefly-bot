# YouTube Briefly Bot

A Telegram bot written in Python 3 that summarizes YouTube links.

## Stack

- `python-telegram-bot`
- `python-dotenv`
- `python-i18n[YAML]`
- `openai`
- `yt-dlp`

## Features

- Extracts video URLs from messages.
- Downloads subtitles/transcripts via `yt-dlp`.
- Summarizes transcript text with an OpenAI-compatible API.
- Uses localized bot messages from `locales/`.
- Applies per-user rate limiting.

## Requirements

- Python 3.11+
- `yt-dlp` available in `PATH`
- `ffmpeg` installed (recommended by `yt-dlp`)

## Setup

1. Create `.env` from `.env.example`.
2. Install dependencies:

```bash
python3 -m pip install -r requirements.txt
```

3. Run the bot:

```bash
python3 -m src.main
```

## Environment Variables

```env
TELEGRAM_BOT_TOKEN=<telegram bot token>
YT_DLP_ADDITIONAL_OPTIONS=
OPENAI_BASE_URL=https://api.openai.com/v1/
OPENAI_API_KEY=<openai compatible key>
OPENAI_MODEL=gpt-4o-mini
LOG_LEVEL=INFO
```

## Tests

```bash
python3 -m pytest
```

## Docker

Build:

```bash
./.devops/Telegram-build.sh
```

Run:

```bash
./.devops/Telegram-run.sh
```
