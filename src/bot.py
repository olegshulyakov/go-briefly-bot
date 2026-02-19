from __future__ import annotations

import asyncio
import logging
import time

from telegram import Update, User
from telegram.ext import (
    Application,
    ApplicationBuilder,
    CommandHandler,
    ContextTypes,
    MessageHandler,
    filters,
)

from .config import Settings
from .load.video_loader import VideoDataLoader
from .load.video_provider import extract_urls
from .localization import translate
from .transform.summarization import OpenAISummarizer
from .utils.text import to_lexical_chunks

logger = logging.getLogger(__name__)


class UserRateLimiter:
    def __init__(self, cooldown_seconds: int) -> None:
        self.cooldown_seconds = cooldown_seconds
        self._last_request: dict[int, float] = {}
        self._lock = asyncio.Lock()

    async def is_limited(self, user_id: int) -> bool:
        async with self._lock:
            now = time.monotonic()
            last_request = self._last_request.get(user_id)
            if last_request is not None and (now - last_request) < self.cooldown_seconds:
                return True

            self._last_request[user_id] = now
            return False


class TelegramBrieflyBot:
    def __init__(self, settings: Settings) -> None:
        self.settings = settings
        self.rate_limiter = UserRateLimiter(settings.rate_limit_window_seconds)
        self.summarizer = OpenAISummarizer(settings)

        self.application: Application = ApplicationBuilder().token(self.settings.telegram_bot_token).build()
        self.application.add_handler(CommandHandler("start", self._start))
        self.application.add_handler(MessageHandler(filters.TEXT & ~filters.COMMAND, self._handle_message))
        self.application.add_error_handler(self._on_error)

    def run(self) -> None:
        logger.info("Starting Telegram bot")
        loop = asyncio.new_event_loop()
        asyncio.set_event_loop(loop)
        self.application.run_polling(drop_pending_updates=True)

    async def _start(self, update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
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
            await message.reply_text(translate("telegram.error.no_url_found", locale=language))
            return

        processing_message = await message.reply_text(translate("telegram.progress.processing", locale=language))

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

        try:
            loader = VideoDataLoader(
                video_url,
                yt_dlp_additional_options=self.settings.yt_dlp_additional_options,
            )
            await asyncio.to_thread(loader.load)
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

        if loader.transcript is None:
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

        await processing_message.edit_text(translate("telegram.progress.summarizing", locale=language))

        try:
            summary = await asyncio.to_thread(
                self.summarizer.summarize_text,
                loader.transcript.transcript,
                language,
            )
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
            title=loader.transcript.title,
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
            await message.reply_text(chunk, disable_web_page_preview=False)

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
                "Failed to summarize transcript",
                extra={
                    "userID": user.id,
                    "username": user.username,
                    "message_id": message.message_id,
                    "error": str(exc),
                },
            )

    async def _on_error(self, update: object, context: ContextTypes.DEFAULT_TYPE) -> None:
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
            await message.reply_text(translate("telegram.error.general", locale=language))
        except Exception as exc:
            logger.exception(
                "Failed to send generic error",
                extra={"error": str(exc)},
            )

    @staticmethod
    def _language(user: User | None) -> str | None:
        return user.language_code if user else None
