from __future__ import annotations

import logging
import tempfile
from dataclasses import dataclass
from pathlib import Path

import yt_dlp

from .transcripts import clean_srt
from .video_provider import build_video_source

logger = logging.getLogger(__name__)
max_attempts = 3


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
        logger.info(
            "Loading video info",
            extra={
                "url": self.url,
                "video_id": self.video_id,
            },
        )

        # Extract video info with retries
        last_error: Exception | None = None
        for attempt in range(max_attempts):
            try:
                ydl_opts = self._build_ydl_opts(
                    {
                        "dumpjson": True,
                    }
                )
                with yt_dlp.YoutubeDL(ydl_opts) as ydl:
                    raw_info = ydl.extract_info(self.url, download=False)
                    self.info = VideoInfo(
                        id=str(raw_info.get("id", "")),
                        language=str(raw_info.get("language", "") or ""),
                        uploader=str(raw_info.get("uploader", "") or ""),
                        title=str(raw_info.get("title", "") or ""),
                        thumbnail=str(raw_info.get("thumbnail", "") or ""),
                        subtitles=dict(raw_info.get("subtitles", {}) or {}),
                    )
                break
            except Exception as exc:
                last_error = exc
                logger.warning(
                    "Failed to load video info",
                    extra={"attempt": attempt + 1, "url": self.url, "error": str(exc)},
                )
        else:
            raise RuntimeError(f"Failed to load video info after {max_attempts} attempts: {last_error}")

        language = self._detect_language(self.info)
        logger.debug("Detected transcript language", extra={"url": self.url, "language": language})

        # Download subtitles with retries
        for attempt in range(max_attempts):
            try:
                ydl_opts = self._build_ydl_opts(
                    {
                        "no_progress": True,
                        "skip_download": True,
                        "writesubtitles": True,
                        "writeautomaticsub": True,
                        "subtitleslangs": [language, f"{language}_auto", "-live_chat"],
                        "subtitlesformat": "srt",
                        "outtmpl": self._subtitle_template_path,
                    }
                )
                with yt_dlp.YoutubeDL(ydl_opts) as ydl:
                    ydl.extract_info(self.url, download=True)
                break
            except Exception as exc:
                last_error = exc
                logger.warning(
                    "Failed to download subtitles",
                    extra={"attempt": attempt + 1, "url": self.url, "error": str(exc)},
                )
        else:
            raise RuntimeError(f"Failed to download subtitles after {max_attempts} attempts: {last_error}")

        subtitle_file = self._find_subtitle_file(language)
        if subtitle_file is None:
            logger.warning("No subtitles found", extra={"url": self.url, "language": language})
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

        logger.info("Transcript loaded successfully", extra={"url": self.url, "length": len(transcript_text)})

        self._cleanup_subtitle_files()

    def _build_ydl_opts(self, extra_options: dict | None = None) -> dict:
        """Build yt-dlp options with additional user options."""
        opts = {
            "quiet": True,
            "no_warnings": True,
            "extract_flat": False,
        }
        if extra_options:
            opts.update(extra_options)
        # Add user-provided options as keyword arguments
        i = 0
        while i < len(self.yt_dlp_additional_options):
            opt = self.yt_dlp_additional_options[i]
            if opt.startswith("--"):
                key = opt[2:].replace("-", "_")
                # Check if next option is a value (not another flag)
                if i + 1 < len(self.yt_dlp_additional_options) and not self.yt_dlp_additional_options[i + 1].startswith(
                    "--"
                ):
                    opts[key] = self.yt_dlp_additional_options[i + 1]
                    i += 1
                else:
                    opts[key] = True
            i += 1
        return opts

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
