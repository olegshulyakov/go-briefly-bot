"""
SRT subtitle cleaning utilities.

Provides functions to parse and clean SRT (SubRip Subtitle) format,
removing timestamps, sequence numbers, and duplicate lines.
"""

from __future__ import annotations

import re

_TIMELINE_RE = re.compile(r"^\d{2}:\d{2}:\d{2},\d{3} --> \d{2}:\d{2}:\d{2},\d{3}$")
"""
Regular expression to match SRT timestamp lines.

Matches format: HH:MM:SS,mmm --> HH:MM:SS,mmm
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


def clean_srt(text: str) -> str:
    """
    Clean SRT subtitle text by removing formatting artifacts.

    Removes:
    - Timestamp lines (00:00:00,000 --> 00:00:00,000)
    - Sequence numbers
    - Special control characters (\\h, \\n, etc.)
    - Duplicate lines

    Args:
        text: Raw SRT subtitle content.

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

        line = _SPECIAL_END_RE.sub("", line).strip()
        if not line or line in seen:
            continue

        seen.add(line)
        chunks.append(line)

    return " ".join(chunks).strip()
