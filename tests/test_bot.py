import asyncio
from unittest.mock import MagicMock, patch

import pytest

from src.bot import TelegramBrieflyBot
from src.config import Settings
from src.cache import LocalCacheProvider
from src.rate_limiter import UserRateLimiter


@pytest.mark.asyncio
async def test_user_rate_limiter_not_limited() -> None:
    limiter = UserRateLimiter(provider=LocalCacheProvider(), cooldown_seconds=1)

    # First request should not be limited
    is_limited = await limiter.is_limited(123)
    assert is_limited is False

    # Second request immediately after should be limited
    is_limited = await limiter.is_limited(123)
    assert is_limited is True


@pytest.mark.asyncio
async def test_user_rate_limiter_different_users() -> None:
    limiter = UserRateLimiter(provider=LocalCacheProvider(), cooldown_seconds=1)

    # Different users should not affect each other
    is_limited_user1 = await limiter.is_limited(123)
    is_limited_user2 = await limiter.is_limited(456)

    assert is_limited_user1 is False
    assert is_limited_user2 is False


@pytest.mark.asyncio
async def test_user_rate_limiter_after_cooldown() -> None:
    limiter = UserRateLimiter(provider=LocalCacheProvider(), cooldown_seconds=0.1)  # Short cooldown for testing

    # First request
    is_limited = await limiter.is_limited(123)
    assert is_limited is False

    # Second request should be limited
    is_limited = await limiter.is_limited(123)
    assert is_limited is True

    # Wait for cooldown period
    await asyncio.sleep(0.2)

    # Third request after cooldown should not be limited
    is_limited = await limiter.is_limited(123)
    assert is_limited is False


def test_telegram_briefly_bot_initialization() -> None:
    mock_settings = MagicMock(spec=Settings)
    mock_settings.rate_limit_window_seconds = 10
    mock_settings.telegram_bot_token = "test_token"
    mock_settings.openai_base_url = "https://api.openai.com/v1/"
    mock_settings.openai_api_key = "test_key"
    mock_settings.openai_model = "gpt-3.5-turbo"
    mock_settings.valkey_url = None
    mock_settings.cache_compression_method = "gzip"

    with patch("src.bot.OpenAISummarizer") as mock_summarizer_class:
        mock_summarizer_instance = MagicMock()
        mock_summarizer_class.return_value = mock_summarizer_instance

        bot = TelegramBrieflyBot(mock_settings)

        # Verify initialization
        assert bot.settings == mock_settings
        assert bot.rate_limiter.cooldown_seconds == 10
        # Verify OpenAISummarizer was initialized with settings
        mock_summarizer_class.assert_called_once_with(mock_settings)
