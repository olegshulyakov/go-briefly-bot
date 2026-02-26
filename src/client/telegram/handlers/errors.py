import logging

from aiogram import Router
from aiogram.types import ErrorEvent

from src.client.telegram.handlers.helpers import get_language
from src.localization import translate

logger = logging.getLogger(__name__)

error_router = Router()


@error_router.errors()
async def error_handler(event: ErrorEvent) -> None:
    """Handles errors during Telegram update processing, logging the exception and notifying the user."""

    logger.exception(
        "Telegram update handling failed",
        extra={"error": str(event.exception)},
    )

    update = event.update
    if update.message is None:
        return

    message = update.message
    language = get_language(message.from_user)
    try:
        await message.reply(translate("telegram.error.general", locale=language))
    except Exception as exc:
        logger.exception(
            "Failed to send generic error",
            extra={"error": str(exc)},
        )
