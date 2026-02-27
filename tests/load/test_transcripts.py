from pathlib import Path
from unittest.mock import MagicMock, patch

from src.config import Settings
from src.load.transcripts import clean_srt
from src.load.video_loader import VideoDataLoader


def build_settings(**overrides: object) -> Settings:
    settings = MagicMock(spec=Settings)
    settings.yt_dlp_additional_options = ()
    settings.cache_transcript_ttl_seconds = 3600
    settings.valkey_url = None
    settings.cache_compression_method = "gzip"
    for key, value in overrides.items():
        setattr(settings, key, value)
    return settings


def test_clean_srt_basic() -> None:
    text = """1
00:00:00,000 --> 00:00:00,001
Hello
2
00:00:00,002 --> 00:00:00,003
World"""
    assert clean_srt(text) == "Hello World"


def test_clean_srt_deduplicates_and_strips_special_suffix() -> None:
    text = """
Some text\\h\\h

Another line\\h
Another line\\h
"""
    assert clean_srt(text) == "Some text Another line"


def test_clean_srt_removes_numeric_lines() -> None:
    text = """1
00:00:01,000 --> 00:00:04,000
First subtitle
2
00:00:05,000 --> 00:00:08,000
Second subtitle
3
00:00:09,000 --> 00:00:12,000
Third subtitle"""
    assert clean_srt(text) == "First subtitle Second subtitle Third subtitle"


def test_clean_srt_removes_timestamp_lines() -> None:
    text = """00:00:01,000 --> 00:00:04,000
Hello
00:00:05,000 --> 00:00:08,000
World"""
    assert clean_srt(text) == "Hello World"


def test_clean_srt_handles_empty_input() -> None:
    assert clean_srt("") == ""


def test_clean_srt_handles_only_whitespace() -> None:
    assert clean_srt("   \n\n\t  \n   ") == ""


def test_clean_srt_handles_duplicate_lines() -> None:
    text = """1
00:00:00,000 --> 00:00:00,001
Hello
2
00:00:00,002 --> 00:00:00,003
Hello
3
00:00:00,004 --> 00:00:00,005
World"""
    assert clean_srt(text) == "Hello World"


def test_clean_srt_removes_special_endings() -> None:
    text = """1
00:00:00,000 --> 00:00:00,001
Hello\\n
2
00:00:00,002 --> 00:00:00,003
World\\h\\c"""
    assert clean_srt(text) == "Hello World"


def test_clean_srt_preserves_punctuation() -> None:
    text = """1
00:00:00,000 --> 00:00:00,001
Hello, world!
2
00:00:00,002 --> 00:00:00,003
How are you? Fine."""
    assert clean_srt(text) == "Hello, world! How are you? Fine."


def test_clean_srt_handles_vtt() -> None:
    text = """WEBVTT
Kind: captions
Language: en

1
00:00:00.000 --> 00:00:00.001
Hello
2
00:00:00.002 --> 00:00:00.003
World"""
    assert clean_srt(text) == "Hello World"


def test_clean_srt_strips_inline_vtt_tags() -> None:
    """VTT files include per-word timing tags and <c> tags inline."""
    text = """WEBVTT

1
00:00:00.000 --> 00:00:00.500
I<00:00:00.160><c> want</c><00:00:00.320><c> to</c>
2
00:00:00.500 --> 00:00:01.000
I want to go"""
    assert clean_srt(text) == "I want to I want to go"


def test_clean_srt_handles_multiline_subtitles() -> None:
    text = """1
00:00:00,000 --> 00:00:00,001
Multi
line
subtitle
2
00:00:00,002 --> 00:00:00,003
Single line"""
    assert clean_srt(text) == "Multi line subtitle Single line"


def test_find_subtitle_file_exact_match() -> None:
    with patch("tempfile.gettempdir") as mock_temp:
        mock_temp.return_value = "/tmp"

        loader = VideoDataLoader(build_settings())

        with patch("tempfile.gettempdir", return_value="/tmp"):
            with patch.object(Path, "exists", return_value=True):
                result = loader._find_subtitle_file("test", "en")
                expected_path = Path("/tmp") / "subtitles_test.en.srt"
                assert result == expected_path


def test_find_subtitle_file_vtt_fallback() -> None:
    """When .srt does not exist, the loader should fall back to .vtt."""
    with patch("tempfile.gettempdir", return_value="/tmp"):
        loader = VideoDataLoader(build_settings())

        def exists_only_vtt(self: Path) -> bool:
            return str(self).endswith(".en.vtt")

        with patch.object(Path, "exists", exists_only_vtt):
            result = loader._find_subtitle_file("test", "en")
            assert result is not None
            assert str(result).endswith(".en.vtt")


def test_clean_srt_handles_various_special_endings() -> None:
    text = """1
00:00:00,000 --> 00:00:00,001
Text with\\h
2
00:00:00,002 --> 00:00:00,003
Text with\\n
3
00:00:00,004 --> 00:00:00,005
Text with\\c"""
    result = clean_srt(text)
    # The function removes duplicates, so we expect only one "Text with"
    assert result == "Text with"
