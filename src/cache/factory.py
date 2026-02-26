"""
Cache provider factory helpers.

Creates cache providers based on application settings.
"""

from __future__ import annotations

from threading import Lock

from ..config import Settings
from .base import CacheProvider
from .in_memory import InMemoryCacheProvider
from .valkey import ValkeyProvider

_provider_lock = Lock()


class ProviderState:
    current: CacheProvider | None = None


_state = ProviderState()


def get_cache_provider(settings: Settings) -> CacheProvider:
    """
    Build or return a cached cache provider instance.

    Args:
        settings: Application settings with cache configuration.

    Returns:
        Cache provider instance.
    """
    if _state.current is not None:
        return _state.current

    with _provider_lock:
        if _state.current is not None:
            return _state.current

        if settings.valkey_url:
            _state.current = ValkeyProvider(
                settings.valkey_url,
                compression_method=settings.cache_compression_method,
            )
        else:
            _state.current = InMemoryCacheProvider(
                compression_method=settings.cache_compression_method,
            )

    if _state.current is None:
        raise RuntimeError("Cache provider not initialized")
    return _state.current


def reset_cache_provider() -> None:
    """
    Reset the cached cache provider.

    Intended for tests to avoid cross-test state sharing.
    """
    _state.current = None
