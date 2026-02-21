"""
Configuration module for loading application settings.

Loads settings from environment variables with validation.
Uses dataclass for immutable, type-safe configuration.
"""

from __future__ import annotations

import os
import shlex
from dataclasses import dataclass

from dotenv import load_dotenv


@dataclass(frozen=True)
class Settings:
    """
    Application settings loaded from environment variables.

    Immutable configuration class (frozen dataclass) that ensures
    settings cannot be modified after initialization.

    Attributes:
        telegram_bot_token: Telegram Bot API token.
        openai_base_url: Base URL for OpenAI-compatible API.
        openai_api_key: API key for LLM service.
        openai_model: Model name to use for summarization.
        valkey_url: Optional connection string for Valkey (replaces in-memory state).
        yt_dlp_additional_options: Additional yt-dlp CLI options.
        rate_limit_window_seconds: Cooldown between user requests.
        cache_summary_ttl_seconds: TTL for cached video summaries (in seconds).
        cache_transcript_ttl_seconds: TTL for cached video transcripts (in seconds).
        max_telegram_message_length: Maximum message length before chunking.
        openai_timeout_seconds: Timeout for LLM API requests.
        openai_max_retries: Maximum retry attempts for LLM API.
    """

    telegram_bot_token: str
    openai_base_url: str
    openai_api_key: str
    openai_model: str
    yt_dlp_additional_options: tuple[str, ...]
    valkey_url: str | None = None
    cache_summary_ttl_seconds: int = 86400
    cache_transcript_ttl_seconds: int = 86400
    rate_limit_window_seconds: int = 10
    max_telegram_message_length: int = 3500
    openai_timeout_seconds: int = 300
    openai_max_retries: int = 3

    @classmethod
    def from_env(cls) -> "Settings":
        """
        Load settings from environment variables.

        Validates that all required environment variables are present.

        Returns:
            Settings instance populated with environment values.

        Raises:
            RuntimeError: If required environment variables are missing.
        """

        load_dotenv()

        telegram_bot_token = os.getenv("TELEGRAM_BOT_TOKEN", "").strip()
        openai_base_url = os.getenv("OPENAI_BASE_URL", "https://api.openai.com/v1/").strip()
        openai_api_key = os.getenv("OPENAI_API_KEY", "").strip()
        openai_model = os.getenv("OPENAI_MODEL", "").strip()
        valkey_url = os.getenv("VALKEY_URL", "").strip() or None

        try:
            cache_summary_ttl = int(os.getenv("CACHE_SUMMARY_TTL_SECONDS", "86400"))
        except ValueError:
            cache_summary_ttl = 86400

        try:
            cache_transcript_ttl = int(os.getenv("CACHE_TRANSCRIPT_TTL_SECONDS", "86400"))
        except ValueError:
            cache_transcript_ttl = 86400

        try:
            rate_limit_window = int(os.getenv("RATE_LIMIT_WINDOW_SECONDS", "10"))
        except ValueError:
            rate_limit_window = 10

        yt_dlp_additional_options = tuple(shlex.split(os.getenv("YT_DLP_ADDITIONAL_OPTIONS", "")))

        missing = []
        if not telegram_bot_token:
            missing.append("TELEGRAM_BOT_TOKEN")
        if not openai_api_key:
            missing.append("OPENAI_API_KEY")
        if not openai_model:
            missing.append("OPENAI_MODEL")

        if missing:
            raise RuntimeError(f"Missing required environment variables: {', '.join(missing)}")

        return cls(
            telegram_bot_token=telegram_bot_token,
            openai_base_url=openai_base_url,
            openai_api_key=openai_api_key,
            openai_model=openai_model,
            valkey_url=valkey_url,
            cache_summary_ttl_seconds=cache_summary_ttl,
            cache_transcript_ttl_seconds=cache_transcript_ttl,
            rate_limit_window_seconds=rate_limit_window,
            yt_dlp_additional_options=yt_dlp_additional_options,
        )
