"""
Telegram bot implementation for video summarization.

Handles:
- User rate limiting
- URL extraction from messages
- Video transcript loading and summarization
- Response formatting and delivery
"""

from __future__ import annotations

import asyncio
import hashlib
import logging
from dataclasses import asdict

from telegram import LinkPreviewOptions, Update, User
from telegram.constants import ParseMode
from telegram.ext import (
    Application,
    ApplicationBuilder,
    CommandHandler,
    ContextTypes,
    MessageHandler,
    filters,
)

from src.utils.markdown import markdown_to_telegram_html

from .config import Settings
from .load.video_loader import VideoDataLoader, VideoTranscript
from .load.video_provider import extract_urls
from .localization import translate
from .storage import LocalProvider, StorageProvider, ValkeyProvider
from .transform.summarization import OpenAISummarizer
from .utils.text import to_lexical_chunks

logger = logging.getLogger(__name__)


class UserRateLimiter:
    """Limits the rate of requests for individual users."""

    def __init__(self, provider: StorageProvider, cooldown_seconds: int) -> None:
        """
        Initializes the UserRateLimiter.

        Args:
            provider: The storage provider for state management.
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


class TelegramBrieflyBot:
    """A Telegram bot that summarizes video content from provided URLs."""

    def __init__(self, settings: Settings) -> None:
        """
        Initializes the TelegramBrieflyBot.

        Args:
            settings: Application settings.
        """

        self.settings = settings
        if settings.valkey_url:
            self.provider: StorageProvider = ValkeyProvider(settings.valkey_url)
        else:
            self.provider: StorageProvider = LocalProvider()

        self.rate_limiter = UserRateLimiter(self.provider, settings.rate_limit_window_seconds)
        self.summarizer = OpenAISummarizer(settings)

        self.application: Application = ApplicationBuilder().token(self.settings.telegram_bot_token).build()
        self.application.add_handler(CommandHandler("start", self._start))
        self.application.add_handler(MessageHandler(filters.TEXT & ~filters.COMMAND, self._handle_message))
        self.application.add_error_handler(self._on_error)

    def run(self) -> None:
        """Starts the bot and runs it indefinitely using polling."""

        logger.info("Starting Telegram bot")
        loop = asyncio.new_event_loop()
        asyncio.set_event_loop(loop)
        self.application.run_polling(drop_pending_updates=True)

    async def _start(self, update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
        """Handles the /start command, sending a welcome message to the user."""

        del context

        message = update.effective_message
        if message is None:
            return

        user = update.effective_user
        language = self._language(user)
        logger.info(
            "User started bot",
            extra={
                "userID": user.id if user else None,
                "username": user.username if user else None,
                "language": language,
            },
        )
        await message.reply_text(translate("telegram.welcome.message", locale=language))

    async def _handle_message(self, update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
        """Handles incoming text messages, extracts URLs, loads video transcripts, summarizes them, and sends the summary back to the user."""

        del context

        user = update.effective_user
        if user is None:
            return

        language = self._language(user)
        if user.is_bot:
            logger.warning("Ignored bot message", extra={"userID": user.id})
            return

        message = update.effective_message
        if message is None:
            logger.warning("Got no message from ", extra={"userID": user.id})
            return

        logger.info(
            "Processing message",
            extra={
                "userID": user.id,
                "username": user.username,
                "language": language,
                "message_id": message.message_id,
            },
        )

        if await self.rate_limiter.is_limited(user.id):
            logger.warning(
                "Rate Limit exceeded",
                extra={
                    "userID": user.id,
                    "username": user.username,
                    "language": language,
                },
            )
            await message.reply_text(
                translate(
                    "telegram.error.rate_limited",
                    locale=language,
                    rateLimitWindow=self.settings.rate_limit_window_seconds,
                ),
                message_thread_id=message.message_thread_id,
                reply_to_message_id=message.id,
            )
            return

        text = message.text or ""
        urls = extract_urls(text)
        if not urls:
            logger.info(
                "No URL found in message",
                extra={
                    "userID": user.id,
                    "username": user.username,
                    "message_id": message.message_id,
                    "text": message,
                },
            )
            await message.reply_text(
                translate("telegram.error.no_url_found", locale=language),
                message_thread_id=message.message_thread_id,
                reply_to_message_id=message.id,
            )
            return

        processing_message = await message.reply_text(
            translate("telegram.progress.processing", locale=language),
            message_thread_id=message.message_thread_id,
            reply_to_message_id=message.id,
        )

        if len(urls) > 1:
            logger.info(
                "Multiple URLs detected",
                extra={
                    "userID": user.id,
                    "username": user.username,
                    "message_id": message.message_id,
                    "url_count": len(urls),
                },
            )
            await processing_message.edit_text(translate("telegram.error.multiple_urls", locale=language))
            return

        video_url = urls[0]
        logger.info(
            "Processing video URL",
            extra={
                "userID": user.id,
                "username": user.username,
                "message_id": message.message_id,
                "url": video_url,
            },
        )

        await processing_message.edit_text(translate("telegram.progress.fetching_info", locale=language))

        video_hash = hashlib.sha256(video_url.encode("utf-8")).hexdigest()
        transcript = None

        cached_transcript_data = await self.provider.get_transcript(video_hash)
        if cached_transcript_data:
            transcript = VideoTranscript(**cached_transcript_data)
            logger.info("Transcript loaded from cache", extra={"userID": user.id, "url": video_url})
        else:
            try:
                loader = VideoDataLoader(
                    video_url,
                    yt_dlp_additional_options=self.settings.yt_dlp_additional_options,
                )
                await asyncio.to_thread(loader.load)
                transcript = loader.transcript

                if transcript:
                    await self.provider.set_transcript(
                        video_hash, asdict(transcript), self.settings.cache_transcript_ttl_seconds
                    )
            except Exception as exc:
                logger.exception(
                    "Failed to load transcript",
                    extra={
                        "userID": user.id,
                        "username": user.username,
                        "message_id": message.message_id,
                        "url": video_url,
                        "error": str(exc),
                    },
                )
                await processing_message.edit_text(translate("telegram.error.transcript_failed", locale=language))
                return

        if transcript is None:
            logger.warning(
                "Transcript is None",
                extra={
                    "userID": user.id,
                    "username": user.username,
                    "message_id": message.message_id,
                    "url": video_url,
                },
            )
            await processing_message.edit_text(translate("telegram.error.transcript_failed", locale=language))
            return

        logger.info(
            "Transcript loaded",
            extra={
                "userID": user.id,
                "username": user.username,
                "message_id": message.message_id,
            },
        )

        cached_summary = await self.provider.get_summary(video_hash)
        if cached_summary:
            summary = cached_summary
            logger.info("Summary loaded from cache", extra={"userID": user.id, "url": video_url})
        else:
            await processing_message.edit_text(translate("telegram.progress.summarizing", locale=language))

            try:
                summary = await asyncio.to_thread(
                    self.summarizer.summarize_text,
                    transcript.transcript,
                    language,
                )
                await self.provider.set_summary(video_hash, summary, self.settings.cache_summary_ttl_seconds)
            except Exception as exc:
                logger.exception(
                    "Failed to summarize transcript",
                    extra={
                        "userID": user.id,
                        "username": user.username,
                        "message_id": message.message_id,
                        "error": str(exc),
                    },
                )
                await processing_message.edit_text(translate("telegram.error.summary_failed", locale=language))
                return

        logger.info(
            "Summary generated",
            extra={
                "userID": user.id,
                "username": user.username,
                "message_id": message.message_id,
                "summary_length": len(summary),
            },
        )

        title = translate(
            "telegram.response.title",
            locale=language,
            title=transcript.title,
            url=video_url,
        )
        response = f"{title}\n{summary}".strip()
        chunks = to_lexical_chunks(response, self.settings.max_telegram_message_length)
        for i, chunk in enumerate(chunks):
            logger.debug(
                "Sending response chunk",
                extra={
                    "userID": user.id,
                    "username": user.username,
                    "message_id": message.message_id,
                    "chunk_index": i,
                },
            )
            await message.reply_text(
                parse_mode=ParseMode.HTML,
                text=markdown_to_telegram_html(chunk),
                message_thread_id=message.message_thread_id,
                reply_to_message_id=message.id,
                link_preview_options=LinkPreviewOptions(is_disabled=False, url=video_url, show_above_text=True),
            )

        logger.info(
            "Response sent",
            extra={
                "userID": user.id,
                "username": user.username,
                "message_id": message.message_id,
                "url": video_url,
            },
        )

        try:
            await processing_message.delete()
        except Exception as exc:
            logger.exception(
                "Failed to delete processing message",
                extra={
                    "userID": user.id,
                    "username": user.username,
                    "message_id": message.message_id,
                    "error": str(exc),
                },
            )

    async def _on_error(self, update: object, context: ContextTypes.DEFAULT_TYPE) -> None:
        """Handles errors during Telegram update processing, logging the exception and notifying the user."""

        logger.exception(
            "Telegram update handling failed",
            extra={"error": str(context.error)},
        )

        if not isinstance(update, Update):
            return

        message = update.effective_message
        if message is None:
            return

        language = self._language(update.effective_user)
        try:
            await message.reply_text(
                translate("telegram.error.general", locale=language),
                message_thread_id=message.message_thread_id,
                reply_to_message_id=message.id,
            )
        except Exception as exc:
            logger.exception(
                "Failed to send generic error",
                extra={"error": str(exc)},
            )

    @staticmethod
    def _language(user: User | None) -> str | None:
        """Retrieves the language code for a given user, if available."""

        return user.language_code if user else None
