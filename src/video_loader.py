from __future__ import annotations

import json
import logging
import subprocess
import tempfile
from dataclasses import dataclass
from pathlib import Path

from .transcripts import clean_srt
from .video_provider import build_video_source

logger = logging.getLogger(__name__)


@dataclass(frozen=True)
class VideoInfo:
    id: str
    language: str
    uploader: str
    title: str
    thumbnail: str
    subtitles: dict[str, list[dict[str, str]]]


@dataclass(frozen=True)
class VideoTranscript:
    id: str
    language: str
    uploader: str
    title: str
    thumbnail: str
    transcript: str


class VideoDataLoader:
    def __init__(self, url: str, yt_dlp_additional_options: tuple[str, ...] = ()) -> None:
        self.url, self.video_id = build_video_source(url)
        self.yt_dlp_additional_options = yt_dlp_additional_options
        self.info: VideoInfo | None = None
        self.transcript: VideoTranscript | None = None

    def load(self) -> None:
        logger.debug("Loading video info", extra={"url": self.url})
        dump_output = self._exec(["--dump-json"], self.url)
        payload = self._extract_json_payload(dump_output)

        raw_info = json.loads(payload)
        self.info = VideoInfo(
            id=str(raw_info.get("id", "")),
            language=str(raw_info.get("language", "") or ""),
            uploader=str(raw_info.get("uploader", "") or ""),
            title=str(raw_info.get("title", "") or ""),
            thumbnail=str(raw_info.get("thumbnail", "") or ""),
            subtitles=dict(raw_info.get("subtitles", {}) or {}),
        )

        language = self._detect_language(self.info)
        logger.debug("Downloading transcript", extra={"url": self.url, "language": language})

        self._exec(
            [
                "--no-progress",
                "--skip-download",
                "--write-subs",
                "--write-auto-subs",
                "--convert-subs",
                "srt",
                "--sub-lang",
                f"{language},{language}_auto,-live_chat",
                "--output",
                self._subtitle_template_path,
            ],
            self.url,
        )

        subtitle_file = self._find_subtitle_file(language)
        if subtitle_file is None:
            raise FileNotFoundError("no subtitles found")

        raw_transcript = subtitle_file.read_text(encoding="utf-8", errors="ignore")
        transcript_text = clean_srt(raw_transcript)

        self.transcript = VideoTranscript(
            id=self.info.id,
            language=language,
            uploader=self.info.uploader,
            title=self.info.title,
            thumbnail=self.info.thumbnail,
            transcript=transcript_text,
        )

        self._cleanup_subtitle_files()

    @property
    def _subtitle_template_path(self) -> str:
        temp_dir = Path(tempfile.gettempdir())
        return str(temp_dir / f"subtitles_{self.video_id}.%(ext)s")

    @property
    def _subtitle_prefix(self) -> Path:
        return Path(tempfile.gettempdir()) / f"subtitles_{self.video_id}"

    def _detect_language(self, info: VideoInfo) -> str:
        if info.language:
            return info.language.split("-", maxsplit=1)[0]

        if info.subtitles:
            if "en" in info.subtitles:
                return "en"
            return next(iter(info.subtitles.keys())).split("-", maxsplit=1)[0]

        raise RuntimeError("no subtitles available")

    def _find_subtitle_file(self, language: str) -> Path | None:
        exact = self._subtitle_prefix.with_suffix(f".{language}.srt")
        if exact.exists():
            return exact

        auto = self._subtitle_prefix.with_suffix(f".{language}_auto.srt")
        if auto.exists():
            return auto

        candidates = sorted(Path(tempfile.gettempdir()).glob(f"subtitles_{self.video_id}*.srt"))
        return candidates[0] if candidates else None

    def _cleanup_subtitle_files(self) -> None:
        for path in Path(tempfile.gettempdir()).glob(f"subtitles_{self.video_id}*"):
            try:
                path.unlink(missing_ok=True)
            except OSError:
                logger.warning("Failed to cleanup temp subtitle file", extra={"path": str(path)})

    def _exec(self, arguments: list[str], url: str) -> str:
        max_attempts = 3
        args = [*self.yt_dlp_additional_options, *arguments, url]

        last_error: Exception | None = None
        last_output = ""

        for attempt in range(max_attempts):
            try:
                result = subprocess.run(
                    ["yt-dlp", *args],
                    check=True,
                    capture_output=True,
                    text=True,
                )
                return (result.stdout or result.stderr or "").strip()
            except subprocess.CalledProcessError as exc:
                last_error = exc
                last_output = f"{exc.stdout or ''}\n{exc.stderr or ''}".strip()
                logger.warning(
                    "yt-dlp failed",
                    extra={
                        "attempt": attempt + 1,
                        "url": url,
                        "returncode": exc.returncode,
                    },
                )

        raise RuntimeError(f"yt-dlp failed after {max_attempts} attempts: {last_error}\n{last_output}")

    @staticmethod
    def _extract_json_payload(output: str) -> str:
        lines = [line.strip() for line in output.splitlines() if line.strip()]
        if not lines:
            raise ValueError("empty yt-dlp metadata output")

        for line in reversed(lines):
            if line.startswith("{") and line.endswith("}"):
                return line

        return lines[-1]
