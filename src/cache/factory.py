"""
Cache provider factory helpers.

Creates cache providers based on application settings.
"""

from __future__ import annotations

from threading import Lock

from ..config import Settings
from .base import CacheProvider
from .local import LocalCacheProvider
from .valkey import ValkeyProvider

_provider_lock = Lock()
_provider: CacheProvider | None = None


def get_cache_provider(settings: Settings) -> CacheProvider:
    """
    Build or return a cached cache provider instance.

    Args:
        settings: Application settings with cache configuration.

    Returns:
        Cache provider instance.
    """
    global _provider
    if _provider is not None:
        return _provider

    with _provider_lock:
        if _provider is not None:
            return _provider

        if settings.valkey_url:
            _provider = ValkeyProvider(
                settings.valkey_url,
                compression_method=settings.cache_compression_method,
            )
        else:
            _provider = LocalCacheProvider()

    return _provider


def reset_cache_provider() -> None:
    """
    Reset the cached cache provider.

    Intended for tests to avoid cross-test state sharing.
    """
    global _provider
    _provider = None
