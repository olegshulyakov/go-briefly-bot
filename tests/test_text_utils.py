from src.text_utils import to_lexical_chunks


def test_to_lexical_chunks_keeps_word_boundaries_when_possible() -> None:
    text = "This is a very long paragraph that exceeds the chunk size and should be split."
    assert to_lexical_chunks(text, 15) == [
        "This is a very",
        "long paragraph",
        "that exceeds",
        "the chunk size",
        "and should be",
        "split.",
    ]


def test_to_lexical_chunks_handles_empty_text() -> None:
    assert to_lexical_chunks("", 10) == [""]
