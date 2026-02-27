"""
SRT and WebVTT subtitle cleaning utilities.

Provides functions to parse and clean both SRT (SubRip Subtitle) and WebVTT
formats, removing timestamps, sequence numbers, WebVTT headers, inline tags,
and duplicate lines to produce clean transcript text.
"""

from __future__ import annotations

import re

_TIMELINE_RE = re.compile(r"^(?:\d{2}:)?\d{2}:\d{2}[.,]\d{3} --> (?:\d{2}:)?\d{2}:\d{2}[.,]\d{3}.*$")
"""
Regular expression to match SRT and VTT timestamp lines.

Matches format: HH:MM:SS,mmm --> HH:MM:SS,mmm or HH:MM:SS.mmm --> HH:MM:SS.mmm
Example: 00:01:23,456 --> 00:01:26,789
"""

_NUMERIC_LINE_RE = re.compile(r"^\d+$")
"""
Regular expression to match numeric lines (subtitle sequence numbers).

Matches lines containing only digits (e.g., "1", "42", "123").
"""

_SPECIAL_END_RE = re.compile(r"(\\\w)+$")
"""
Regular expression to match special SRT control sequences at line endings.

Matches patterns like \\h\\h (hard line breaks) or \\n (newlines)
that some subtitle files include for formatting.
"""

_WEBVTT_HEADER_RE = re.compile(
    r"^(?:"
    r"WEBVTT"  # File header
    r"|NOTE(?:\s|$)"  # Comment blocks
    r"|STYLE(?:\s|$)"  # Style blocks
    r"|REGION(?:\s|$)"  # Region blocks
    r"|Kind:"  # Metadata: track kind
    r"|Language:"  # Metadata: language
    r"|(?:Align|Line|Position|Size|Snap-to-lines|Vertical|Scroll):"  # Known cue settings
    r")"
)
"""
Regular expression to match WebVTT structural lines that must be excluded from
the transcript: the file header, comment/style/region blocks, metadata lines,
and known cue-setting directives (Align:, Line:, Position:, Size:, etc.).
"""

_HTML_TAG_RE = re.compile(r"<[^>]+>")
"""
Regular expression to match HTML-like tags inside subtitles.

Matches patterns like <b>, <i>, <c>, or VTT per-word timestamp tags like <00:00:00.000>.
"""


def clean_srt(text: str) -> str:
    """
    Clean SRT or WebVTT subtitle text by removing formatting artifacts.

    Handles both SRT and WebVTT formats, removing:
    - WebVTT file headers, metadata, comment/style/region blocks
    - Timestamp lines (SRT: 00:00:00,000 --> / VTT: 00:00:00.000 --> with optional cue settings)
    - Sequence numbers
    - Inline HTML-like tags (<b>, <i>, <c>, <00:00:00.000>)
    - Special SRT control characters (\\h, \\n, etc.)
    - Duplicate lines

    Args:
        text: Raw SRT or WebVTT subtitle content.

    Returns:
        Cleaned transcript as continuous text with duplicates removed.
    """

    seen: set[str] = set()
    chunks: list[str] = []

    for raw_line in text.splitlines():
        line = raw_line.strip()
        if not line:
            continue
        if _TIMELINE_RE.match(line):
            continue
        if _NUMERIC_LINE_RE.match(line):
            continue
        if _WEBVTT_HEADER_RE.match(line):
            continue

        line = _SPECIAL_END_RE.sub("", line)
        line = _HTML_TAG_RE.sub("", line).strip()

        if not line or line in seen:
            continue

        seen.add(line)
        chunks.append(line)

    return " ".join(chunks).strip()
