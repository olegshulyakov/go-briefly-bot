# YouTube Briefly Bot

A Telegram bot written in Python that summarizes Video content using LLM.

## Stack

- `python-telegram-bot`
- `python-dotenv`
- `python-i18n[YAML]`
- `openai`
- `yt-dlp`

## Features

- 🎬 **Video URL Extraction** — automatically finds Video links in messages
- 📝 **Subtitle Download** — downloads and processes subtitles via `yt-dlp`
- 🤖 **AI Summarization** — creates transcript summaries via OpenAI-compatible API
- 🌍 **Localization** — supports 13 languages (en, ru, de, es, fr, it, pt, ar, zh, cn, ja, ko, hi)
- ⏱️ **Rate Limiting** — abuse protection with per-user cooldown
- 📊 **Message Chunking** — automatic splitting of long responses into parts

## Supported Platforms

- YouTube (regular videos and Shorts)
- VK Video

## Requirements

- Python 3.11+
- `ffmpeg` (recommended by yt-dlp)

## Quick Start

### 1. Clone and Setup

```bash
git clone https://github.com/olegshulyakov/go-briefly-bot.git
cd go-briefly-bot
cp .env.example .env
```

### 2. Configure `.env`

```env
TELEGRAM_BOT_TOKEN=your_bot_token_here
OPENAI_API_KEY=your_openai_api_key_here
OPENAI_MODEL=gpt-4o-mini
OPENAI_BASE_URL=https://api.openai.com/v1/
LOG_LEVEL=INFO
```

### 3. Install Dependencies

**Option A: Using pip (recommended)**

```bash
python3 -m pip install -e ".[dev]"
```

**Option B: From requirements.txt**

```bash
python3 -m pip install -r requirements.txt
```

### 4. Install Git Hooks (optional)

```bash
chmod +x install-hooks.sh
./install-hooks.sh
```

This enables pre-commit hooks that run tests and linting before each commit.

### 5. Run the Bot

```bash
python3 -m src.main
```

Or if installed as a package:

```bash
go-briefly-bot
```

## VS Code Setup

This project includes pre-configured VS Code settings for optimal Python development.

### Recommended Extensions

Install the recommended extensions:

```bash
# Open the Command Palette (Cmd+Shift+P) and run:
# "Extensions: Show Recommended Extensions"
```

Or install via CLI:

```bash
code --install-extension ms-python.python
code --install-extension ms-python.vscode-pylance
code --install-extension charliermarsh.ruff
code --install-extension tamasfe.even-better-toml
```

### Available Tasks

Open the Command Palette (`Cmd+Shift+P` or `Ctrl+Shift+P`) → **Tasks: Run Task**:

| Task | Description |
|------|-------------|
| `📦 Install Dependencies` | Install from pyproject.toml |
| `▶️ Run Bot` | Run the bot application |
| `🧪 Run All Tests` | Run pytest suite |
| `🧪 Run Tests with Coverage` | Run tests with coverage report |
| `🔍 Lint (Ruff Check)` | Run linter |
| `✨ Format (Ruff Format)` | Format code |
| `🔧 Lint & Fix (Ruff)` | Auto-fix lint issues |
| `🧹 Clean Python Cache` | Remove __pycache__ |
| `🔒 Install Git Hooks` | Setup git hooks |

### Debug Configurations

Press `F5` to start debugging. Available configurations:

| Configuration | Description |
|---------------|-------------|
| `🐍 Run Bot` | Debug the bot |
| `🧪 Test All` | Debug all tests |
| `🧪 Test Current File` | Debug current test file |
| `🧪 Test with Coverage` | Run tests with coverage |

### Keyboard Shortcuts

- **Run Bot**: `Cmd+Shift+B` (macOS) / `Ctrl+Shift+B` (Windows/Linux)
- **Run Tests**: Open Command Palette → Tasks → Run Task → Tests
- **Format Code**: `Shift+Alt+F` (default)
- **Quick Fix**: `Cmd+.` (macOS) / `Ctrl+.` (Windows/Linux)

## Environment Variables

| Variable                    | Description                        | Default                      |
| --------------------------- | ---------------------------------- | ---------------------------- |
| `TELEGRAM_BOT_TOKEN`        | Telegram bot token (required)      | —                            |
| `OPENAI_API_KEY`            | LLM API key (required)             | —                            |
| `OPENAI_MODEL`              | Model for summarization (required) | —                            |
| `OPENAI_BASE_URL`           | OpenAI-compatible API base URL     | `https://api.openai.com/v1/` |
| `YT_DLP_ADDITIONAL_OPTIONS` | Additional yt-dlp options          | —                            |
| `LOG_LEVEL`                 | Logging level                      | `INFO`                       |

## Tests

```bash
# Run all tests
python3 -m pytest

# Run with coverage
python3 -m pytest --cov=src

# Run specific module
python3 -m pytest tests/test_bot.py -v
```

## Docker

### Build Image

```bash
./.devops/Telegram-build.sh
```

### Run Container

```bash
./.devops/Telegram-run.sh
```

### Docker Compose

```bash
docker-compose up -d
```

## Localization

The bot supports the following languages:

- 🇬🇧 English (`en`)
- 🇷🇺 Русский (`ru`)
- 🇩🇪 Deutsch (`de`)
- 🇪🇸 Español (`es`)
- 🇫🇷 Français (`fr`)
- 🇮🇹 Italiano (`it`)
- 🇵🇹 Português (`pt`)
- 🇸🇦 العربية (`ar`)
- 🇨🇳 中文 (`zh`, `cn`)
- 🇯🇵 日本語 (`ja`)
- 🇰🇷 한국어 (`ko`)
- 🇮🇳 हिन्दी (`hi`)

Language is automatically detected from the user's message settings.

## License

[MIT License](LICENSE)
