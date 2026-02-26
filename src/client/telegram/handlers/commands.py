import logging

from aiogram import Router
from aiogram.filters import Command
from aiogram.types import Message

from src.client.telegram.handlers.helpers import get_language
from src.localization import translate

logger = logging.getLogger(__name__)

start_router = Router()


@start_router.message(Command("start"))
async def start_command(message: Message) -> None:
    """Handles the /start command, sending a welcome message to the user."""
    user = message.from_user
    language = get_language(user)
    logger.info(
        "User started bot",
        extra={
            "userID": user.id if user else None,
            "username": user.username if user else None,
            "language": language,
        },
    )
    await message.reply(translate("telegram.welcome.message", locale=language))
