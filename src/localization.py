from __future__ import annotations

from pathlib import Path

import i18n

DEFAULT_LOCALE = "en"

_initialized = False


def _setup_i18n() -> None:
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


def _normalize_locale(locale: str | None) -> str:
    if not locale:
        return DEFAULT_LOCALE

    normalized = locale.strip().lower().replace("_", "-")
    if not normalized:
        return DEFAULT_LOCALE

    # Telegram can send language tags like "en-US". Our locale files use base language only.
    return normalized.split("-", maxsplit=1)[0]


def translate(key: str, locale: str | None = None, **kwargs: object) -> str:
    _setup_i18n()

    lang = _normalize_locale(locale)
    for candidate in (lang, DEFAULT_LOCALE):
        try:
            value = i18n.t(key, locale=candidate, **kwargs)
        except Exception:
            continue
        if value != key:
            return str(value)

    return key
