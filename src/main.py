from __future__ import annotations

import logging

from .bot import TelegramBrieflyBot
from .config import Settings
from .logger import configure_logging

logger = logging.getLogger(__name__)


def main() -> None:
    """
    Main entry point of the application.
    Configures logging, loads settings, initializes and runs the Telegram bot.
    """
    configure_logging()
    logger.info("Starting application")
    settings = Settings.from_env()
    telegram = TelegramBrieflyBot(settings)
    telegram.run()


if __name__ == "__main__":
    main()
