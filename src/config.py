from __future__ import annotations

import os
import shlex
from dataclasses import dataclass

from dotenv import load_dotenv


@dataclass(frozen=True)
class Settings:
    telegram_bot_token: str
    openai_base_url: str
    openai_api_key: str
    openai_model: str
    yt_dlp_additional_options: tuple[str, ...]
    rate_limit_window_seconds: int = 10
    max_telegram_message_length: int = 3500
    openai_timeout_seconds: int = 300
    openai_max_retries: int = 3

    @classmethod
    def from_env(cls) -> "Settings":
        load_dotenv()

        telegram_bot_token = os.getenv("TELEGRAM_BOT_TOKEN", "").strip()
        openai_base_url = os.getenv("OPENAI_BASE_URL", "https://api.openai.com/v1/").strip()
        openai_api_key = os.getenv("OPENAI_API_KEY", "").strip()
        openai_model = os.getenv("OPENAI_MODEL", "").strip()
        yt_dlp_additional_options = tuple(shlex.split(os.getenv("YT_DLP_ADDITIONAL_OPTIONS", "")))

        missing = []
        if not telegram_bot_token:
            missing.append("TELEGRAM_BOT_TOKEN")
        if not openai_api_key:
            missing.append("OPENAI_API_KEY")
        if not openai_model:
            missing.append("OPENAI_MODEL")

        if missing:
            raise RuntimeError(f"Missing required environment variables: {', '.join(missing)}")

        return cls(
            telegram_bot_token=telegram_bot_token,
            openai_base_url=openai_base_url,
            openai_api_key=openai_api_key,
            openai_model=openai_model,
            yt_dlp_additional_options=yt_dlp_additional_options,
        )
