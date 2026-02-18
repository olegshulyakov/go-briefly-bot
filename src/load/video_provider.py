from __future__ import annotations

import re
from dataclasses import dataclass


@dataclass(frozen=True)
class RegexProvider:
    pattern: re.Pattern[str]
    canonical_url: str

    def is_valid_url(self, text: str) -> bool:
        return bool(self.pattern.search(text))

    def get_id(self, text: str) -> str:
        match = self.pattern.search(text)
        if not match:
            raise ValueError("no valid URL found")
        return match.group(1)

    def extract_urls(self, text: str) -> list[str]:
        return [match.group(0) for match in self.pattern.finditer(text)]

    def canonicalize(self, text: str) -> tuple[str, str]:
        video_id = self.get_id(text)
        return self.canonical_url % video_id, video_id


YOUTUBE = RegexProvider(
    pattern=re.compile(r"(?:https?://)?(?:www\.)?(?:youtube\.com/watch\?.*?v=|youtu\.be/)([a-zA-Z0-9_-]{11})"),
    canonical_url="https://www.youtube.com/watch?v=%s",
)
YOUTUBE_SHORT = RegexProvider(
    pattern=re.compile(r"(?:https?://)?(?:www\.)?youtube\.com/shorts/([A-Za-z0-9_-]{11})"),
    canonical_url="https://www.youtube.com/shorts/%s",
)
VKVIDEO = RegexProvider(
    pattern=re.compile(r"(?:https?://)?(?:www\.)?vkvideo\.ru/(video-\d+_\d+)"),
    canonical_url="https://vkvideo.ru/%s",
)

PROVIDERS: tuple[RegexProvider, ...] = (YOUTUBE, YOUTUBE_SHORT, VKVIDEO)


def extract_urls(text: str) -> list[str]:
    for provider in PROVIDERS:
        if provider.is_valid_url(text):
            return provider.extract_urls(text)
    return []


def build_video_source(url: str) -> tuple[str, str]:
    for provider in PROVIDERS:
        if provider.is_valid_url(url):
            return provider.canonicalize(url)
    raise ValueError(f"no valid URL found: {url}")
