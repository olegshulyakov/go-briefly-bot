"""
Configuration module for loading application settings.

Loads settings from environment variables with validation.
Uses dataclass for immutable, type-safe configuration.
"""

from __future__ import annotations

import os
import shlex
from dataclasses import dataclass
from typing import Any
from urllib.parse import urlparse

from dotenv import load_dotenv

DEFAULT_OPENAI_BASE_URL = "https://api.openai.com/v1/"
DEFAULT_OPENAI_TIMEOUT_SECONDS = 300
DEFAULT_OPENAI_MAX_RETRIES = 3
DEFAULT_CACHE_TTL_WITH_VALKEY = 86400
DEFAULT_CACHE_TTL_NO_VALKEY = 3600
DEFAULT_CACHE_COMPRESSION_METHOD = "gzip"
DEFAULT_RATE_LIMIT_WINDOW_SECONDS = 10
DEFAULT_MAX_TELEGRAM_MESSAGE_LENGTH = 3500


@dataclass(frozen=True)
class Settings:
    """
    Application settings loaded from environment variables.

    Immutable configuration class (frozen dataclass) that ensures
    settings cannot be modified after initialization.
    """

    telegram_bot_token: str
    telegram_proxy_url: str | None
    openai_base_url: str
    openai_api_key: str
    openai_model: str
    yt_dlp_additional_options: tuple[str, ...]
    valkey_url: str | None = None
    cache_summary_ttl_seconds: int = DEFAULT_CACHE_TTL_WITH_VALKEY
    cache_transcript_ttl_seconds: int = DEFAULT_CACHE_TTL_WITH_VALKEY
    cache_compression_method: str = DEFAULT_CACHE_COMPRESSION_METHOD
    rate_limit_window_seconds: int = DEFAULT_RATE_LIMIT_WINDOW_SECONDS
    max_telegram_message_length: int = DEFAULT_MAX_TELEGRAM_MESSAGE_LENGTH
    openai_timeout_seconds: int = DEFAULT_OPENAI_TIMEOUT_SECONDS
    openai_max_retries: int = DEFAULT_OPENAI_MAX_RETRIES

    @classmethod
    def from_env(cls) -> Settings:
        """
        Load settings from environment variables.

        Validates that all required environment variables are present.

        Returns:
            Settings instance populated with environment values.

        Raises:
            RuntimeError: If required environment variables are missing.
        """
        load_dotenv()
        env_vars = _load_env_vars()
        _validate_env_vars(env_vars)
        return cls(**env_vars)


def _load_env_vars() -> dict[str, Any]:
    """Load and parse all environment variables into a dictionary."""
    telegram_bot_token = os.getenv("TELEGRAM_BOT_TOKEN", "").strip()
    telegram_proxy_url = os.getenv("TELEGRAM_PROXY_URL", "").strip() or None
    openai_base_url = os.getenv("OPENAI_BASE_URL", DEFAULT_OPENAI_BASE_URL).strip()
    openai_api_key = os.getenv("OPENAI_API_KEY", "").strip()
    openai_model = os.getenv("OPENAI_MODEL", "").strip()
    valkey_url = os.getenv("VALKEY_URL", "").strip() or None

    default_ttl = DEFAULT_CACHE_TTL_NO_VALKEY if valkey_url is None else DEFAULT_CACHE_TTL_WITH_VALKEY

    return {
        "telegram_bot_token": telegram_bot_token,
        "telegram_proxy_url": telegram_proxy_url,
        "openai_base_url": openai_base_url,
        "openai_api_key": openai_api_key,
        "openai_model": openai_model,
        "openai_timeout_seconds": _load_int("OPENAI_TIMEOUT_SECONDS", DEFAULT_OPENAI_TIMEOUT_SECONDS),
        "openai_max_retries": _load_int("OPENAI_MAX_RETRIES", DEFAULT_OPENAI_MAX_RETRIES),
        "valkey_url": valkey_url,
        "cache_summary_ttl_seconds": _load_int("CACHE_SUMMARY_TTL_SECONDS", default_ttl),
        "cache_transcript_ttl_seconds": _load_int("CACHE_TRANSCRIPT_TTL_SECONDS", default_ttl),
        "cache_compression_method": _load_compression_method(),
        "rate_limit_window_seconds": _load_int("RATE_LIMIT_WINDOW_SECONDS", DEFAULT_RATE_LIMIT_WINDOW_SECONDS),
        "max_telegram_message_length": _load_int("MAX_TELEGRAM_MESSAGE_LENGTH", DEFAULT_MAX_TELEGRAM_MESSAGE_LENGTH),
        "yt_dlp_additional_options": tuple(shlex.split(os.getenv("YT_DLP_ADDITIONAL_OPTIONS", ""))),
    }


def _load_compression_method() -> str:
    """Load and validate cache compression method."""
    cache_compression = os.getenv("CACHE_COMPRESSION_METHOD", DEFAULT_CACHE_COMPRESSION_METHOD).strip().lower()
    valid_methods = {"none", "gzip", "zlib", "lzma"}
    return cache_compression if cache_compression in valid_methods else DEFAULT_CACHE_COMPRESSION_METHOD


def _load_int(env_var: str, default: int) -> int:
    """Load integer value from environment with fallback to default."""
    try:
        return int(os.getenv(env_var, str(default)))
    except ValueError:
        return default


def _validate_env_vars(env_vars: dict[str, Any]) -> None:
    """Validate required environment variables and proxy configuration."""
    missing = []
    if not env_vars["telegram_bot_token"]:
        missing.append("TELEGRAM_BOT_TOKEN")
    if not env_vars["openai_api_key"]:
        missing.append("OPENAI_API_KEY")
    if not env_vars["openai_model"]:
        missing.append("OPENAI_MODEL")

    if missing:
        raise RuntimeError(f"Missing required environment variables: {', '.join(missing)}")

    if env_vars["telegram_proxy_url"]:
        parsed = urlparse(env_vars["telegram_proxy_url"])
        if not parsed.scheme or not parsed.hostname or not parsed.port:
            raise RuntimeError("Invalid TELEGRAM_PROXY_URL format. Expected: protocol://user:password@host:port")
        if parsed.scheme not in {"http", "https", "socks4", "socks5"}:
            raise RuntimeError(f"Unsupported proxy protocol: {parsed.scheme}. Supported: http, https, socks4, socks5")
