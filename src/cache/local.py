"""
Local in-memory cache provider implementation.
"""

import asyncio
import time

from .base import CacheProvider


class LocalCacheProvider(CacheProvider):
    """
    In-memory cache provider using dictionaries.
    Used when Valkey is not configured or unavailable.
    """

    def __init__(self) -> None:
        self._rate_limits: dict[int, float] = {}
        self._rate_limit_lock = asyncio.Lock()

        # Caches
        self._cache: dict[str, tuple[str, float]] = {}
        self._cache_dict: dict[str, tuple[dict, float]] = {}
        self._cache_lock = asyncio.Lock()

    async def is_rate_limited(self, user_id: int, window_seconds: int) -> bool:
        async with self._rate_limit_lock:
            now = time.monotonic()
            last_request = self._rate_limits.get(user_id)
            if last_request is not None and (now - last_request) < window_seconds:
                return True

            self._rate_limits[user_id] = now
            return False

    async def get(self, key: str) -> str | None:
        async with self._cache_lock:
            cached = self._cache.get(key)
            if cached is None:
                return None

            text, expires_at = cached
            if time.monotonic() > expires_at:
                del self._cache[key]
                return None

            return text

    async def put(self, key: str, text: str, ttl_seconds: int) -> None:
        async with self._cache_lock:
            self._cache[key] = (text, time.monotonic() + ttl_seconds)

    async def get_dict(self, key: str) -> dict | None:
        async with self._cache_lock:
            cached = self._cache_dict.get(key)
            if cached is None:
                return None

            data, expires_at = cached
            if time.monotonic() > expires_at:
                del self._cache_dict[key]
                return None

            return data

    async def put_dict(self, key: str, data: dict, ttl_seconds: int) -> None:
        async with self._cache_lock:
            self._cache_dict[key] = (data, time.monotonic() + ttl_seconds)
