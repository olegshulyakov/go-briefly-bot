from __future__ import annotations

import logging
import threading

from .bot import TelegramBrieflyBot
from .config import Settings
from .logger import configure_logging

logger = logging.getLogger(__name__)


def main() -> None:
    configure_logging()
    logger.info("Starting application")
    settings = Settings.from_env()
    telegram = TelegramBrieflyBot(settings)
    threading.Thread(target=telegram.run()).start()


if __name__ == "__main__":
    main()
