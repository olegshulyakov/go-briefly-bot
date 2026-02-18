from __future__ import annotations

import re

_TIMELINE_RE = re.compile(r"^\d{2}:\d{2}:\d{2},\d{3} --> \d{2}:\d{2}:\d{2},\d{3}$")
_NUMERIC_LINE_RE = re.compile(r"^\d+$")
_SPECIAL_END_RE = re.compile(r"(\\\w)+$")


def clean_srt(text: str) -> str:
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
