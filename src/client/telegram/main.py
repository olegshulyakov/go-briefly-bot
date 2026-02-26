from __future__ import annotations

import asyncio
import logging

from aiogram import Bot, Dispatcher
from aiogram.client.default import DefaultBotProperties
from aiogram.enums import ParseMode

from src.cache.base import CacheProvider
from src.cache.factory import get_cache_provider
from src.client.telegram.handlers import commands_router, errors_router, messages_router
from src.config import Settings
from src.load.video_loader import VideoDataLoader
from src.logger import configure_logging
from src.rate_limiter import UserRateLimiter
from src.transform.summarization import OpenAISummarizer

logger = logging.getLogger(__name__)


async def main() -> None:
    """
    Main entry point of the application.
    Configures logging, loads settings, initializes and runs the Telegram bot.
    """
    settings = Settings.from_env()
    provider: CacheProvider = get_cache_provider(settings)

    rate_limiter = UserRateLimiter(provider, settings.rate_limit_window_seconds)
    loader = VideoDataLoader(settings)
    summarizer = OpenAISummarizer(settings)

    # aiogram setup
    dp = Dispatcher()
    dp.include_routers(commands_router, messages_router, errors_router)

    bot = Bot(
        token=settings.telegram_bot_token,
        default=DefaultBotProperties(parse_mode=ParseMode.HTML),
    )

    await dp.start_polling(
        bot,
        settings=settings,
        rate_limiter=rate_limiter,
        loader=loader,
        summarizer=summarizer,
    )


if __name__ == "__main__":
    configure_logging()
    logger.info("Starting Telegram bot")
    try:
        asyncio.run(main())
    except (KeyboardInterrupt, SystemExit):
        logger.info("Bot stopped")
