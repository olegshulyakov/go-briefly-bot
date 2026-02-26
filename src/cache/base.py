"""
Base cache provider interface for Go Briefly Bot.
"""

import json
import logging
from abc import ABC, abstractmethod
from typing import Any

from ..utils import CompressionMethod, compress, decompress

logger = logging.getLogger(__name__)


class CacheProvider(ABC):
    """Abstract base class for cache providers (cache and rate limiting)."""

    def __init__(self, compression_method: str = "gzip") -> None:
        self._compression_method = self._parse_compression_method(compression_method)

    def _parse_compression_method(self, method: str) -> CompressionMethod:
        """Parse compression method string to enum."""
        method_map = {
            "none": CompressionMethod.NONE,
            "gzip": CompressionMethod.GZIP,
            "zlib": CompressionMethod.ZLIB,
            "lzma": CompressionMethod.LZMA,
        }
        return method_map.get(method.lower(), CompressionMethod.GZIP)

    def _encode_text(self, text: str) -> bytes:
        return compress(text.encode("utf-8"), self._compression_method)

    def _decode_text(self, data: bytes) -> str | None:
        try:
            return decompress(data).decode("utf-8")
        except Exception as e:
            logger.warning("Failed to decompress/decode cached string: %s", e)
            return None

    def _encode_dict(self, data: dict[str, Any]) -> bytes:
        return compress(json.dumps(data).encode("utf-8"), self._compression_method)

    def _decode_dict(self, data: bytes) -> dict[str, Any] | None:
        try:
            res: dict[str, Any] = json.loads(decompress(data).decode("utf-8"))
            return res
        except Exception as e:
            logger.warning("Failed to decompress/decode cached dict: %s", e)
            return None

    @abstractmethod
    async def is_rate_limited(self, user_id: int, window_seconds: int) -> bool:
        """
        Checks if a user is rate-limited. If not, records the request.

        Args:
            user_id: The ID of the user.
            window_seconds: The rate limit sliding window in seconds.

        Returns:
            True if limited, False otherwise.
        """
        pass

    @abstractmethod
    async def put(self, key: str, text: str, ttl_seconds: int) -> None:
        """Caches a text."""
        pass

    @abstractmethod
    async def get(self, key: str) -> str | None:
        """Retrieves a cached text."""
        pass

    @abstractmethod
    async def get_dict(self, key: str) -> dict[str, Any] | None:
        """Retrieves a cached dict."""
        pass

    @abstractmethod
    async def put_dict(self, key: str, data: dict[str, Any], ttl_seconds: int) -> None:
        """Caches a dict."""
        pass
