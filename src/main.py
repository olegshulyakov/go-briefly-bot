from __future__ import annotations

import logging
import os

from .bot import TelegramBrieflyBot
from .config import Settings


def configure_logging() -> None:
    level_name = os.getenv("LOG_LEVEL", "INFO").upper()
    level = getattr(logging, level_name, logging.INFO)
    logging.basicConfig(
        level=level,
        format="%(asctime)s %(levelname)s [%(name)s] %(message)s",
    )


def main() -> None:
    configure_logging()
    settings = Settings.from_env()
    bot = TelegramBrieflyBot(settings)
    bot.run()


if __name__ == "__main__":
    main()
