from src.video_provider import build_video_source, extract_urls


def test_extract_urls_youtube() -> None:
    text = "text https://youtu.be/validID1234 and https://www.youtube.com/watch?v=anotherID12"
    assert extract_urls(text) == [
        "https://youtu.be/validID1234",
        "https://www.youtube.com/watch?v=anotherID12",
    ]


def test_build_video_source_canonicalizes_short_link() -> None:
    canonical, video_id = build_video_source("youtu.be/abcdefghijk")
    assert video_id == "abcdefghijk"
    assert canonical == "https://www.youtube.com/watch?v=abcdefghijk"
