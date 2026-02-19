"""
Logging configuration module.

Sets up application-wide logging with configurable log level.
Suppresses verbose HTTP client logs to reduce noise.
"""

import logging
import os


def configure_logging() -> None:
    """
    Configure application logging.

    Sets up:
    - Root logger with specified level from LOG_LEVEL env var
    - Custom format with timestamp, level, logger name, and message
    - Suppresses httpx library logs (too verbose)
    """
    level_name = os.getenv("LOG_LEVEL", "INFO").upper()
    level = getattr(logging, level_name, logging.INFO)
    logging.basicConfig(
        level=level,
        format="%(asctime)s %(levelname)s [%(name)s] %(message)s",
    )

    logging.getLogger("httpx").setLevel(logging.WARNING)
