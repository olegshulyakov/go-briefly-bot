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
from .localization import translate
from .summarization import OpenAISummarizer
from .text_utils import to_lexical_chunks
from .video_loader import VideoDataLoader
from .video_provider import extract_urls

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

    def run(self) -> None:
        application: Application = ApplicationBuilder().token(self.settings.telegram_bot_token).build()
        application.add_handler(CommandHandler("start", self._start))
        application.add_handler(MessageHandler(filters.TEXT & ~filters.COMMAND, self._handle_message))
        application.add_error_handler(self._on_error)
        application.run_polling(drop_pending_updates=True)

    async def _start(self, update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
        del context

        message = update.effective_message
        if message is None:
            return

        language = self._language(update.effective_user)
        await message.reply_text(translate("telegram.welcome.message", locale=language))

    async def _handle_message(self, update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
        del context

        user = update.effective_user
        if user is None:
            return

        language = self._language(user)
        if user.is_bot:
            logger.warning("Ignored bot message", extra={"user_id": user.id, "language": language})
            return

        message = update.effective_message
        if message is None:
            return

        if await self.rate_limiter.is_limited(user.id):
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
            await message.reply_text(translate("telegram.error.no_url_found", locale=language))
            return

        processing_message = await message.reply_text(translate("telegram.progress.processing", locale=language))

        if len(urls) > 1:
            await processing_message.edit_text(translate("telegram.error.multiple_urls", locale=language))

        video_url = urls[0]

        await processing_message.edit_text(translate("telegram.progress.fetching_info", locale=language))

        try:
            loader = VideoDataLoader(
                video_url,
                yt_dlp_additional_options=self.settings.yt_dlp_additional_options,
            )
            await asyncio.to_thread(loader.load)
        except Exception as exc:
            logger.exception("Failed to load transcript", extra={"url": video_url, "error": str(exc)})
            await processing_message.edit_text(translate("telegram.error.transcript_failed", locale=language))
            return

        if loader.transcript is None:
            await processing_message.edit_text(translate("telegram.error.transcript_failed", locale=language))
            return

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
                extra={"url": video_url, "error": str(exc)},
            )
            await processing_message.edit_text(translate("telegram.error.summary_failed", locale=language))
            return

        title = translate(
            "telegram.response.title",
            locale=language,
            title=loader.transcript.title,
            url=video_url,
        )
        response = f"{title}\n{summary}".strip()

        for chunk in to_lexical_chunks(response, self.settings.max_telegram_message_length):
            await message.reply_text(chunk, disable_web_page_preview=False)

        try:
            await processing_message.delete()
        except Exception:
            logger.debug("Failed to delete processing message", exc_info=True)

    async def _on_error(self, update: object, context: ContextTypes.DEFAULT_TYPE) -> None:
        logger.exception("Telegram update handling failed", exc_info=context.error)

        if not isinstance(update, Update):
            return

        message = update.effective_message
        if message is None:
            return

        language = self._language(update.effective_user)
        try:
            await message.reply_text(translate("telegram.error.general", locale=language))
        except Exception:
            logger.debug("Failed to send generic error", exc_info=True)

    @staticmethod
    def _language(user: User | None) -> str | None:
        return user.language_code if user else None
