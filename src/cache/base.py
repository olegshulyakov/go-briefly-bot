"""
Base cache provider interface for Go Briefly Bot.
"""

from abc import ABC, abstractmethod


class CacheProvider(ABC):
    """Abstract base class for cache providers (cache and rate limiting)."""

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
    async def get_summary(self, video_hash: str, language_code: str | None) -> str | None:
        """Retrieves a cached summary for a video."""
        pass

    @abstractmethod
    async def set_summary(self, video_hash: str, language_code: str | None, summary: str, ttl_seconds: int) -> None:
        """Caches a summary for a video."""
        pass

    @abstractmethod
    async def get_transcript(self, video_hash: str) -> dict | None:
        """Retrieves a cached transcript dict for a video."""
        pass

    @abstractmethod
    async def set_transcript(self, video_hash: str, transcript_data: dict, ttl_seconds: int) -> None:
        """Caches a transcript dict for a video."""
        pass
