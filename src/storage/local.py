"""
Local in-memory storage provider implementation.
"""

import asyncio
import time

from .base import StorageProvider


class LocalProvider(StorageProvider):
    """
    In-memory storage provider using dictionaries.
    Used when Valkey is not configured or unavailable.
    """

    def __init__(self) -> None:
        self._rate_limits: dict[int, float] = {}
        self._rate_limit_lock = asyncio.Lock()

        # Caches
        self._summaries: dict[str, tuple[str, float]] = {}
        self._transcripts: dict[str, tuple[dict, float]] = {}
        self._cache_lock = asyncio.Lock()

    async def is_rate_limited(self, user_id: int, window_seconds: int) -> bool:
        async with self._rate_limit_lock:
            now = time.monotonic()
            last_request = self._rate_limits.get(user_id)
            if last_request is not None and (now - last_request) < window_seconds:
                return True

            self._rate_limits[user_id] = now
            return False

    async def get_summary(self, video_hash: str, language_code: str | None) -> str | None:
        key = f"{video_hash}:{language_code}"
        async with self._cache_lock:
            cached = self._summaries.get(key)
            if cached is None:
                return None

            summary, expires_at = cached
            if time.monotonic() > expires_at:
                del self._summaries[key]
                return None

            return summary

    async def set_summary(self, video_hash: str, language_code: str | None, summary: str, ttl_seconds: int) -> None:
        key = f"{video_hash}:{language_code}"
        async with self._cache_lock:
            self._summaries[key] = (summary, time.monotonic() + ttl_seconds)

    async def get_transcript(self, video_hash: str) -> dict | None:
        async with self._cache_lock:
            cached = self._transcripts.get(video_hash)
            if cached is None:
                return None

            transcript, expires_at = cached
            if time.monotonic() > expires_at:
                del self._transcripts[video_hash]
                return None

            return transcript

    async def set_transcript(self, video_hash: str, transcript_data: dict, ttl_seconds: int) -> None:
        async with self._cache_lock:
            self._transcripts[video_hash] = (transcript_data, time.monotonic() + ttl_seconds)
