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
