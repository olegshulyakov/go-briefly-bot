from __future__ import annotations

import logging
import time

from openai import OpenAI

from ..config import Settings
from ..localization import translate

logger = logging.getLogger(__name__)


class OpenAISummarizer:
    def __init__(self, settings: Settings) -> None:
        self.settings = settings
        self.client = OpenAI(
            base_url=settings.openai_base_url,
            api_key=settings.openai_api_key,
            max_retries=settings.openai_max_retries,
        )

    def summarize_text(self, text: str, locale: str | None) -> str:
        logger.info(
            "Summarizing text",
            extra={"locale": locale, "text_length": len(text), "model": self.settings.openai_model},
        )
        system_prompt = translate("llm.system", locale=locale)

        last_error: Exception | None = None
        for attempt in range(self.settings.openai_max_retries):
            try:
                start_time = time.monotonic()
                response = self.client.chat.completions.create(
                    model=self.settings.openai_model,
                    messages=[
                        {"role": "system", "content": system_prompt},
                        {"role": "user", "content": text},
                    ],
                    timeout=self.settings.openai_timeout_seconds,
                )
                elapsed = time.monotonic() - start_time

                if not response or not response.choices:
                    logger.warning(
                        "OpenAI returned no response",
                        extra={"response_id": getattr(response, "id", None), "model": getattr(response, "model", None)},
                    )
                    raise RuntimeError("no OpenAI response")

                content = response.choices[0].message.content
                if not content:
                    logger.warning(
                        "OpenAI returned empty response",
                        extra={
                            "response_id": getattr(response, "id", None),
                            "model": getattr(response, "model", None),
                            "choices": getattr(response, "choices", None),
                        },
                    )
                    raise RuntimeError("empty OpenAI response")

                logger.info(
                    "Summary received",
                    extra={
                        "locale": locale,
                        "model": self.settings.openai_model,
                        "elapsed_ms": int(elapsed * 1000),
                        "content_length": len(content),
                    },
                )
                return content
            except Exception as exc:  # pragma: no cover - upstream SDK/runtime errors
                last_error = exc
                logger.warning(
                    "OpenAI summarization attempt failed",
                    extra={"attempt": attempt + 1, "error": str(exc)},
                )
                time.sleep(1)

        raise RuntimeError(f"failed to summarize text: {last_error}")
