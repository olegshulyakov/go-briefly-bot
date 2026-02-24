"""
Localization module for internationalization (i18n) support.

Provides translation functionality using the python-i18n library.
Supports multiple locales with automatic fallback to English.
"""

from __future__ import annotations

from pathlib import Path

import i18n

DEFAULT_LOCALE = "en"

_initialized = False


def _setup_i18n() -> None:
    """
    Initialize i18n configuration (thread-safe).

    Sets up:
    - Locale file path
    - File format (YAML)
    - Fallback locale
    - Memoization for performance
    """
    global _initialized
    if _initialized:
        return

    locales_path = Path(__file__).resolve().parents[1] / "locales"
    i18n.load_path.append(str(locales_path))
    i18n.set("filename_format", "locale.{locale}.{format}")
    i18n.set("file_format", "yml")
    i18n.set("fallback", DEFAULT_LOCALE)
    i18n.set("skip_locale_root_data", True)
    i18n.set("locale", DEFAULT_LOCALE)
    i18n.set("enable_memoization", True)

    _initialized = True


def normalize_locale(locale: str | None) -> str:
    """
    Normalize a locale code to base language.

    Converts locale tags like 'en-US' to 'en' for matching
    against available locale files.

    Args:
        locale: Raw locale string (e.g., 'en-US', 'pt_BR').

    Returns:
        Normalized base language code (e.g., 'en', 'pt').
    """
    if not locale:
        return DEFAULT_LOCALE

    normalized = locale.strip().lower().replace("_", "-")
    if not normalized:
        return DEFAULT_LOCALE

    # Telegram can send language tags like "en-US". Our locale files use base language only.
    return normalized.split("-", maxsplit=1)[0]


def translate(key: str, locale: str | None = None, **kwargs: object) -> str:
    """
    Translate a localization key to the specified language.

    Attempts translation in this order:
    1. Requested locale
    2. Default locale (English)
    3. Return key as fallback

    Args:
        key: Localization key (e.g., 'telegram.welcome.message').
        locale: Target locale code (e.g., 'ru', 'es').
        **kwargs: Variables to interpolate in the translation.

    Returns:
        Translated string or the original key if translation fails.
    """
    _setup_i18n()

    lang = normalize_locale(locale)
    for candidate in (lang, DEFAULT_LOCALE):
        try:
            value = i18n.t(key, locale=candidate, **kwargs)
        except Exception:
            continue
        if value != key:
            return str(value)

    return key
