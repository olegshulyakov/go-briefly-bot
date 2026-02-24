"""
User rate limiting helpers.

Provides a per-user rate limiter backed by the cache provider.
"""

from __future__ import annotations

from .cache import CacheProvider


class UserRateLimiter:
    """Limits the rate of requests for individual users."""

    def __init__(self, provider: CacheProvider, cooldown_seconds: int) -> None:
        """
        Initializes the UserRateLimiter.

        Args:
            provider: The cache provider for state management.
            cooldown_seconds: The cooldown period in seconds for each user.
        """

        self.provider = provider
        self.cooldown_seconds = cooldown_seconds

    async def is_limited(self, user_id: int) -> bool:
        """
        Checks if a user is rate-limited. If not, records the current request time.

        Args:
            user_id: The ID of the user to check.

        Returns:
            True if the user is rate-limited, False otherwise.
        """
        return await self.provider.is_rate_limited(user_id, self.cooldown_seconds)
