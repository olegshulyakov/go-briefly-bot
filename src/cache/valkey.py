"""
Valkey cache provider with fail-soft fallback.
"""

import asyncio
import json
import logging
from typing import Any

from valkey.asyncio import Valkey
from valkey.exceptions import ConnectionError as ValkeyConnectionError

from .base import CacheProvider
from .local import LocalCacheProvider

logger = logging.getLogger(__name__)


class ValkeyProvider(CacheProvider):
    """
    Valkey-based cache provider, utilizing atomic operations and TTL.
    Implements a strict fail-soft mechanism falling back to LocalCacheProvider
    on connection drop or timeout.
    """

    def __init__(self, valkey_url: str) -> None:
        self.valkey_url = valkey_url
        self._valkey: Valkey | None = None
        self._local_fallback = LocalCacheProvider()
        self._timeout = 0.200  # Strict 200ms fail-soft timeout

    async def _get_client(self) -> Valkey:
        """Lazy initialization of the Valkey client."""
        if self._valkey is None:
            self._valkey = Valkey.from_url(self.valkey_url)
        return self._valkey

    async def _safe_execute(self, coro_fn: Any, fallback_coro_fn: Any) -> Any:
        try:
            return await asyncio.wait_for(coro_fn(), timeout=self._timeout)
        except (TimeoutError, asyncio.TimeoutError):
            logger.warning("Valkey timeout, falling back to LocalCacheProvider")
            return await fallback_coro_fn()
        except (ValkeyConnectionError, ConnectionError, OSError) as e:
            logger.warning(f"Valkey connection error: {e}, falling back to LocalCacheProvider")
            return await fallback_coro_fn()
        except Exception as e:
            logger.warning(f"Valkey operation failed: {e}, falling back to LocalCacheProvider")
            return await fallback_coro_fn()

    async def is_rate_limited(self, user_id: int, window_seconds: int) -> bool:
        """Atomic rate limiter using INCR and EXPIRE."""

        async def _valkey_rate_limit() -> bool:
            client = await self._get_client()
            key = f"user:{user_id}:limit"

            pipeline = client.pipeline()
            pipeline.incr(key)
            pipeline.expire(key, window_seconds, nx=True)
            result = await pipeline.execute()

            # result[0] is the response of INCR
            count = result[0]
            if count > 1:
                return True
            return False

        return await self._safe_execute(
            _valkey_rate_limit,
            lambda: self._local_fallback.is_rate_limited(user_id, window_seconds),
        )

    async def get_summary(self, video_hash: str, language_code: str | None) -> str | None:
        async def _valkey_get() -> str | None:
            client = await self._get_client()
            key = f"summary:{video_hash}:{language_code}"
            val = await client.get(key)
            if val is not None:
                return val.decode("utf-8")
            return None

        return await self._safe_execute(
            _valkey_get,
            lambda: self._local_fallback.get_summary(video_hash, language_code),
        )

    async def set_summary(self, video_hash: str, language_code: str | None, summary: str, ttl_seconds: int) -> None:
        async def _valkey_set() -> None:
            client = await self._get_client()
            key = f"summary:{video_hash}:{language_code}"
            await client.setex(key, ttl_seconds, summary)

        await self._safe_execute(
            _valkey_set,
            lambda: self._local_fallback.set_summary(video_hash, language_code, summary, ttl_seconds),
        )

    async def get_transcript(self, video_hash: str) -> dict | None:
        async def _valkey_get() -> dict | None:
            client = await self._get_client()
            key = f"transcript:{video_hash}"
            val = await client.get(key)
            if val is not None:
                return json.loads(val.decode("utf-8"))
            return None

        return await self._safe_execute(
            _valkey_get,
            lambda: self._local_fallback.get_transcript(video_hash),
        )

    async def set_transcript(self, video_hash: str, transcript_data: dict, ttl_seconds: int) -> None:
        async def _valkey_set() -> None:
            client = await self._get_client()
            key = f"transcript:{video_hash}"
            await client.setex(key, ttl_seconds, json.dumps(transcript_data))

        await self._safe_execute(
            _valkey_set,
            lambda: self._local_fallback.set_transcript(video_hash, transcript_data, ttl_seconds),
        )
