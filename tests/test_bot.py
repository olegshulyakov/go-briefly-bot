import asyncio
from unittest.mock import AsyncMock, MagicMock, patch

import pytest
from src.bot import TelegramBrieflyBot
from src.cache import LocalCacheProvider
from src.config import Settings
from src.load.video_loader import VideoTranscript
from src.rate_limiter import UserRateLimiter
from telegram import Message, Update, User
from telegram.ext import ContextTypes


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
    limiter = UserRateLimiter(provider=LocalCacheProvider(), cooldown_seconds=1)  # Short cooldown for testing

    # First request
    is_limited = await limiter.is_limited(123)
    assert is_limited is False

    # Second request should be limited
    is_limited = await limiter.is_limited(123)
    assert is_limited is True

    # Wait for cooldown period
    await asyncio.sleep(1.1)

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
    mock_settings.yt_dlp_additional_options = ()

    with patch("src.bot.OpenAISummarizer") as mock_summarizer_class:
        mock_summarizer_instance = MagicMock()
        mock_summarizer_class.return_value = mock_summarizer_instance

        bot = TelegramBrieflyBot(mock_settings)

        # Verify initialization
        assert bot.settings == mock_settings
        expected_cooldown = 10
        assert bot.rate_limiter.cooldown_seconds == expected_cooldown
        # Verify OpenAISummarizer was initialized with settings
        mock_summarizer_class.assert_called_once_with(mock_settings)


@pytest.fixture
def mock_settings() -> Settings:
    settings = MagicMock(spec=Settings)
    settings.rate_limit_window_seconds = 10
    settings.telegram_bot_token = "test_token"
    settings.openai_base_url = "https://api.openai.com/v1/"
    settings.openai_api_key = "test_key"
    settings.openai_model = "gpt-3.5-turbo"
    settings.valkey_url = None
    settings.cache_compression_method = "gzip"
    settings.yt_dlp_additional_options = ()
    settings.max_telegram_message_length = 4000
    return settings


@pytest.fixture
def bot(mock_settings: Settings) -> TelegramBrieflyBot:
    with patch("src.bot.ApplicationBuilder"):
        return TelegramBrieflyBot(mock_settings)


@pytest.fixture
def mock_update() -> MagicMock:
    update = MagicMock(spec=Update)
    user = MagicMock(spec=User)
    user.id = 123
    user.username = "testuser"
    user.language_code = "en"
    user.is_bot = False
    update.effective_user = user

    message = AsyncMock(spec=Message)
    message.message_id = 1
    message.id = 1
    message.message_thread_id = None
    message.text = "https://youtube.com/watch?v=123"
    message.reply_text = AsyncMock()
    update.effective_message = message

    return update


@pytest.fixture
def mock_context() -> MagicMock:
    return MagicMock(spec=ContextTypes.DEFAULT_TYPE)


@pytest.mark.asyncio
async def test_bot_start(bot: TelegramBrieflyBot, mock_update: MagicMock, mock_context: MagicMock) -> None:
    # Need to mock translate since we aren't loading translations in unit tests by default
    with patch("src.bot.translate", return_value="Welcome"):
        await bot._start(mock_update, mock_context)
    mock_update.effective_message.reply_text.assert_called_once_with("Welcome")


@pytest.mark.asyncio
async def test_bot_handle_message_no_text(bot: TelegramBrieflyBot, mock_update: MagicMock, mock_context: MagicMock) -> None:
    mock_update.effective_message.text = None
    with (
        patch.object(bot.rate_limiter, "is_limited", return_value=False),
        patch("src.bot.translate", return_value="No URL"),
    ):
        await bot._handle_message(mock_update, mock_context)
    mock_update.effective_message.reply_text.assert_called_once()
    assert mock_update.effective_message.reply_text.call_args[0][0] == "No URL"


@pytest.mark.asyncio
async def test_bot_handle_message_rate_limited(bot: TelegramBrieflyBot, mock_update: MagicMock, mock_context: MagicMock) -> None:
    with (
        patch.object(bot.rate_limiter, "is_limited", return_value=True),
        patch("src.bot.translate", return_value="Rate Limited"),
    ):
        await bot._handle_message(mock_update, mock_context)
    mock_update.effective_message.reply_text.assert_called_once_with("Rate Limited", message_thread_id=None, reply_to_message_id=1)


@pytest.mark.asyncio
async def test_bot_handle_message_multiple_urls(bot: TelegramBrieflyBot, mock_update: MagicMock, mock_context: MagicMock) -> None:
    mock_update.effective_message.text = "url1.com url2.com"
    with (
        patch.object(bot.rate_limiter, "is_limited", return_value=False),
        patch("src.bot.extract_urls", return_value=["url1", "url2"]),
        patch("src.bot.translate", return_value="Error"),
    ):
        processing_msg_mock = AsyncMock()
        mock_update.effective_message.reply_text.return_value = processing_msg_mock
        await bot._handle_message(mock_update, mock_context)
        processing_msg_mock.edit_text.assert_called_once_with("Error")


@pytest.mark.asyncio
async def test_bot_handle_message_success(bot: TelegramBrieflyBot, mock_update: MagicMock, mock_context: MagicMock) -> None:
    transcript = VideoTranscript(id="123", language="en", uploader="test", title="Test Video", thumbnail="", transcript="Test transcript")

    processing_msg_mock = AsyncMock()
    mock_update.effective_message.reply_text.return_value = processing_msg_mock

    with (
        patch.object(bot.rate_limiter, "is_limited", return_value=False),
        patch("src.bot.extract_urls", return_value=["https://youtube.com/watch?v=123"]),
        patch.object(bot.loader, "load", return_value=transcript) as mock_load,
        patch.object(bot.summarizer, "summarize", return_value="Test summary") as mock_summarize,
        patch("src.bot.translate", side_effect=lambda key, **kw: key),
    ):
        await bot._handle_message(mock_update, mock_context)

        mock_load.assert_called_once_with("https://youtube.com/watch?v=123")
        mock_summarize.assert_called_once_with("Test transcript", "en")

        # Original message reply for processing, and second reply for final result
        expected_calls = 2
        assert mock_update.effective_message.reply_text.call_count == expected_calls
        processing_msg_mock.delete.assert_called_once()


@pytest.mark.asyncio
async def test_bot_handle_message_loader_fails(bot: TelegramBrieflyBot, mock_update: MagicMock, mock_context: MagicMock) -> None:
    processing_msg_mock = AsyncMock()
    mock_update.effective_message.reply_text.return_value = processing_msg_mock

    with (
        patch.object(bot.rate_limiter, "is_limited", return_value=False),
        patch("src.bot.extract_urls", return_value=["https://youtube.com/watch?v=123"]),
        patch.object(bot.loader, "load", side_effect=Exception("Load error")),
        patch("src.bot.translate", return_value="Fail"),
    ):
        await bot._handle_message(mock_update, mock_context)
        processing_msg_mock.edit_text.assert_called_with("Fail")


@pytest.mark.asyncio
async def test_bot_handle_message_summarizer_fails(bot: TelegramBrieflyBot, mock_update: MagicMock, mock_context: MagicMock) -> None:
    transcript = VideoTranscript(id="123", language="en", uploader="test", title="Test Video", thumbnail="", transcript="Test transcript")
    processing_msg_mock = AsyncMock()
    mock_update.effective_message.reply_text.return_value = processing_msg_mock

    with (
        patch.object(bot.rate_limiter, "is_limited", return_value=False),
        patch("src.bot.extract_urls", return_value=["https://youtube.com/watch?v=123"]),
        patch.object(bot.loader, "load", return_value=transcript),
        patch.object(bot.summarizer, "summarize", side_effect=Exception("Summarize error")),
        patch("src.bot.translate", return_value="Fail"),
    ):
        await bot._handle_message(mock_update, mock_context)
        processing_msg_mock.edit_text.assert_called_with("Fail")


@pytest.mark.asyncio
async def test_bot_on_error(bot: TelegramBrieflyBot, mock_update: MagicMock, mock_context: MagicMock) -> None:
    mock_context.error = Exception("General error")
    with patch("src.bot.translate", return_value="Error handled"):
        await bot._on_error(mock_update, mock_context)
    mock_update.effective_message.reply_text.assert_called_once_with("Error handled", message_thread_id=None, reply_to_message_id=1)
