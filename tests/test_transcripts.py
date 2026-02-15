from src.transcripts import clean_srt


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
