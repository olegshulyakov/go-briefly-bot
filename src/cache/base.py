"""
Base cache provider interface for Go Briefly Bot.
"""

from abc import ABC, abstractmethod

from src.utils.compression import CompressionMethod


class CacheProvider(ABC):
    """Abstract base class for cache providers (cache and rate limiting)."""

    def _parse_compression_method(self, method: str) -> CompressionMethod:
        """Parse compression method string to enum."""
        method_map = {
            "none": CompressionMethod.NONE,
            "gzip": CompressionMethod.GZIP,
            "zlib": CompressionMethod.ZLIB,
            "lzma": CompressionMethod.LZMA,
        }
        return method_map.get(method.lower(), CompressionMethod.GZIP)

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
    async def get_dict(self, key: str) -> dict | None:
        """Retrieves a cached dict."""
        pass

    @abstractmethod
    async def put_dict(self, key: str, data: dict, ttl_seconds: int) -> None:
        """Caches a dict."""
        pass
