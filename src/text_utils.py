from __future__ import annotations


def to_lexical_chunks(text: str, chunk_size: int) -> list[str]:
    if chunk_size <= 0 or len(text) <= chunk_size:
        return [text.strip()]

    chunks: list[str] = []
    left = 0
    text_length = len(text)

    while left < text_length:
        right = min(left + chunk_size, text_length)
        if right == text_length:
            chunk = text[left:right].strip()
            if chunk:
                chunks.append(chunk)
            break

        right = _find_natural_breakpoint(text, left, right)
        right = min(right, text_length)

        chunk = text[left:right].strip()
        if chunk:
            chunks.append(chunk)

        left = right
        while left < text_length and text[left] in {" ", "\n"}:
            left += 1

    return chunks or [""]


def _find_natural_breakpoint(text: str, left: int, right: int) -> int:
    if right < len(text) and text[right] == "\n":
        return right

    piece = text[left:right]

    paragraph_idx = piece.rfind("\n")
    if paragraph_idx > 0:
        return left + paragraph_idx + 1

    sentence_idx = max(piece.rfind("."), piece.rfind("!"), piece.rfind("?"))
    if sentence_idx > 0:
        return left + sentence_idx + 1

    if right < len(text) and text[right] != " ":
        word_idx = piece.rfind(" ")
        if word_idx > 0:
            return left + word_idx + 1

    return right
