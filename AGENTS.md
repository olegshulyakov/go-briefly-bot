# Briefly Bot

This project is a Telegram bot written in Python designed to summarize video content, primarily from YouTube and VK Video, using an OpenAI-compatible Large Language Model (LLM). It automatically extracts video links, downloads and processes subtitles using `yt-dlp`, and then generates summaries.

## Core Setup commands

- Install deps: `.venv/bin/pip install -e .[dev]`.
- Run tests: `.venv/bin/python -m pytest tests/ -v`.

## Code style & Agent Rules

- **Python strict mode.**
- **Use object-oriented programming** where possible.
- **Style Guide:** Follow PEP8 conventions and use `ruff` for formatting/linting.
- **Type Hinting:** All code MUST include type hints and return types.
- **Docstrings:** Use Google style docstrings for all public modules and classes.
- **Error Handling:** Raise domain-specific exceptions; never swallow errors silently.

### 🏗️ Mandatory Architecture Rules (Zero Tolerance)

1.  **`__init__.py` is an Entry Point:** It MUST ONLY contain re-exports or package wiring. No business logic or utility functions allowed here.
2.  **No Catch-All Files:** Names like `utils.py`, `helpers.py`, or `services.py` are BANNED. Use domain-specific names (e.g., `yt_dlp_handler.py`, `date_formatter.py`).
3.  **Single Responsibility Principle:** Every file MUST have one clear, nameable purpose. If you can't describe it in one phrase, split it.
4.  **300 LOC Hard Limit:** Any file > 300 lines of logic (excluding docstrings/prompts) is a code smell and MUST be refactored immediately.
5.  **Refactor-First:** If a task requires touching a file that violates these rules, you MUST refactor it BEFORE adding new logic.
