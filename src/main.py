from __future__ import annotations

import logging
import os
import threading

from .bot import TelegramBrieflyBot
from .config import Settings

logger = logging.getLogger(__name__)


def configure_logging() -> None:
    level_name = os.getenv("LOG_LEVEL", "INFO").upper()
    level = getattr(logging, level_name, logging.INFO)
    logging.basicConfig(
        level=level,
        format="%(asctime)s %(levelname)s [%(name)s] %(message)s",
    )


def main() -> None:
    configure_logging()
    logger.info("Starting application")
    settings = Settings.from_env()
    telegram = TelegramBrieflyBot(settings)
    telegram_thread = threading.Thread(target=telegram.run())
    telegram_thread.start()


if __name__ == "__main__":
    main()
