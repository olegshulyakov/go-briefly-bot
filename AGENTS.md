# AGENTS.md - Project Overview & Code Style

This project is a Telegram bot written in Python designed to summarize video content, primarily from YouTube and VK Video, using an OpenAI-compatible Large Language Model (LLM). It automatically extracts video links, downloads and processes subtitles using `yt-dlp`, and then generates summaries.

## Core Setup commands

- Python Version: 3.11.
- Install deps: `.venv/bin/pip install -e .[dev]`.
- Run tests: `.venv/bin/python -m pytest tests/ -v`.

## Code style & Agent Rules

- Python strict mode.
- Use object-oriented programming where possible.
- Write simple, clean and well-documented code.
- Avoid nested if statements and complex logic.
- Style Guide: Follow PEP8 conventions.
- Formatting: Use `ruff` for auto-formatting. Run with `.venv/bin/python -m ruff format`.
- Linting: Use `ruff` for linting and type checking. Run with `.venv/bin/python -m ruff check`.
- Type Hinting: All Python code MUST include type hints and return types.
- Naming Conventions: Use `snake_case` for functions/variables and `CamelCase` for classes.
- Docstrings: Use Google style docstrings for all public modules, functions, classes, and methods.
- Error Handling: Raise domain-specific exceptions with context; never swallow errors silently.

## Technologies Used

- **Core:** Python 3.11+
- **Telegram Bot Framework:** `python-telegram-bot`
- **Configuration:** `python-dotenv`
- **Internationalization:** `python-i18n[YAML]`
- **LLM Integration:** `openai` (for OpenAI-compatible API)
- **Video Processing:** `yt-dlp`
- **State & Rate Limiting:** `valkey`
- **Linting & Formatting:** `ruff`
- **Testing:** `pytest`

## Architecture Highlights

- **Video URL Extraction:** Automatically identifies video links within messages.
- **Subtitle Download:** Utilizes `yt-dlp` to download and process video subtitles.
- **AI Summarization:** Integrates with an OpenAI-compatible API to create summaries from transcripts.
- **Localization:** Supports 13 languages, with automatic detection based on user settings.
- **Rate Limiting & Caching:** Implements a pluggable `CacheProvider`. Uses `valkey` for distributed lock/rate operations and caching (summaries/transcripts), falling back to local memory if disconnected (`fail-soft`).
- **Message Chunking:** Automatically splits long responses into manageable parts for Telegram.

## Setup and Installation

### 1. Clone the Repository

```bash
git clone https://github.com/olegshulyakov/go-briefly-bot.git
cd go-briefly-bot
```

### 2. Configure Environment Variables

Copy the example environment file and fill in your details:

```bash
cp .env.example .env
```

Edit the `.env` file with your `TELEGRAM_BOT_TOKEN`, `OPENAI_API_KEY`, and `OPENAI_MODEL`.

### 3. Install Dependencies

**Using pip (recommended):**

```bash
python3 -m pip install -e ".[dev]"
```

**From `requirements.txt`:**

```bash
python3 -m pip install -r requirements.txt
```

### 4. Install Git Hooks (Optional)

To enable pre-commit hooks for testing and linting:

```bash
chmod +x install-hooks.sh
./install-hooks.sh
```

## Running the Bot

### Directly

```bash
python3 -m src.main
```

### As an Installed Package

```bash
go-briefly-bot
```

## Testing

### Run All Tests

```bash
python3 -m pytest
```

### Run Tests with Coverage

```bash
python3 -m pytest --cov=src
```

### Run Specific Test Module

```bash
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

### Using Docker Compose

```bash
docker-compose up -d
```

## Localization

The bot supports the following languages, with automatic detection from user message settings:

- English (`en`)
- Russian (`ru`)
- German (`de`)
- Spanish (`es`)
- French (`fr`)
- Italian (`it`)
- Portuguese (`pt`)
- Arabic (`ar`)
- Chinese (`zh`, `cn`)
- Japanese (`ja`)
- Korean (`ko`)
- Hindi (`hi`)

Localization files are located in the `locales/` directory.

## Development Environment (VS Code)

The project includes pre-configured VS Code settings (`.vscode/`).

### Recommended Extensions

The following extensions are recommended for optimal Python development:

- `ms-python.python`
- `ms-python.vscode-pylance`
- `charliermarsh.ruff`
- `tamasfe.even-better-toml`

### Available Tasks

Access via Command Palette (`Cmd+Shift+P` / `Ctrl+Shift+P`) → **Tasks: Run Task**:

- `📦 Install Dependencies`: Install from `pyproject.toml`.
- `▶️ Run Bot`: Run the bot application.
- `🧪 Run All Tests`: Run the `pytest` suite.
- `🧪 Run Tests with Coverage`: Run tests with coverage report.
- `🔍 Lint (Ruff Check)`: Run the linter.
- `✨ Format (Ruff Format)`: Format code.
- `🔧 Lint & Fix (Ruff)`: Auto-fix lint issues.
- `🧹 Clean Python Cache`: Remove `__pycache__` directories.
- `🔒 Install Git Hooks`: Setup git hooks.

### Debug Configurations

Press `F5` to start debugging. Available configurations:

- `🐍 Run Bot`: Debug the bot.
- `🧪 Test All`: Debug all tests.
- `🧪 Test Current File`: Debug the currently open test file.
- `🧪 Test with Coverage`: Run tests with coverage.

## Configuration (Environment Variables)

The bot's behavior can be configured using the following environment variables in the `.env` file:
| Variable | Description | Default |
| :-------------------------- | :---------------------------------------------- | :--------------------------- |
| `TELEGRAM_BOT_TOKEN` | Telegram bot token (required) | — |
| `OPENAI_API_KEY` | LLM API key (required) | — |
| `OPENAI_MODEL` | Model to use for summarization (required) | — |
| `OPENAI_BASE_URL` | Base URL for OpenAI-compatible API | `https://api.openai.com/v1/` |
| `YT_DLP_ADDITIONAL_OPTIONS` | Additional options to pass to `yt-dlp` | — |
| `VALKEY_URL` | Valkey connection URL (optional) | — |
| `CACHE_SUMMARY_TTL_SECONDS` | TTL for cached summaries | `3600` (local), `86400` (Valkey) |
| `CACHE_TRANSCRIPT_TTL_SECONDS` | TTL for cached transcripts | `3600` (local), `86400` (Valkey) |
| `CACHE_COMPRESSION_METHOD` | Compression for Valkey cache | `gzip` (none, gzip, zlib, lzma) |
| `LOG_LEVEL` | Logging level (e.g., `INFO`, `DEBUG`, `WARNING`) | `INFO` |

## License

This project is licensed under the [MIT License](LICENSE).
