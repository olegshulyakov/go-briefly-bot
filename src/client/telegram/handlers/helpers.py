from aiogram.types import User

from src.localization import normalize_locale


def get_language(user: User | None) -> str:
    """Retrieves the language code for a given user, if available."""
    return normalize_locale(user.language_code if user else None)
