from unittest.mock import patch, MagicMock
from pathlib import Path
from src.load.video_loader import VideoDataLoader, VideoInfo, VideoTranscript


def test_video_info_creation() -> None:
    info = VideoInfo(
        id="test_id",
        language="en",
        uploader="test_uploader",
        title="test_title",
        thumbnail="test_thumbnail",
        subtitles={"en": [{"url": "sub_url"}]},
    )

    assert info.id == "test_id"
    assert info.language == "en"
    assert info.uploader == "test_uploader"
    assert info.title == "test_title"
    assert info.thumbnail == "test_thumbnail"
    assert info.subtitles == {"en": [{"url": "sub_url"}]}


def test_video_transcript_creation() -> None:
    transcript = VideoTranscript(
        id="test_id",
        language="en",
        uploader="test_uploader",
        title="test_title",
        thumbnail="test_thumbnail",
        transcript="test transcript",
    )

    assert transcript.id == "test_id"
    assert transcript.language == "en"
    assert transcript.uploader == "test_uploader"
    assert transcript.title == "test_title"
    assert transcript.thumbnail == "test_thumbnail"
    assert transcript.transcript == "test transcript"


def test_video_data_loader_initialization() -> None:
    with patch("src.load.video_loader.build_video_source") as mock_build:
        mock_build.return_value = ("https://www.youtube.com/watch?v=test", "test")

        loader = VideoDataLoader("https://youtu.be/test")

        mock_build.assert_called_once_with("https://youtu.be/test")
        assert loader.url == "https://www.youtube.com/watch?v=test"
        assert loader.video_id == "test"
        assert loader.info is None
        assert loader.transcript is None


def test_video_data_loader_initialization_with_options() -> None:
    with patch("src.load.video_loader.build_video_source") as mock_build:
        mock_build.return_value = ("https://www.youtube.com/watch?v=test", "test")

        loader = VideoDataLoader("https://youtu.be/test", yt_dlp_additional_options=("--format", "mp4"))

        assert loader.yt_dlp_additional_options == ("--format", "mp4")


def test_video_data_loader_subtitle_template_path() -> None:
    with patch("src.load.video_loader.build_video_source") as mock_build:
        mock_build.return_value = ("https://www.youtube.com/watch?v=test", "test")
        with patch("tempfile.gettempdir") as mock_temp:
            mock_temp.return_value = "/tmp"

            loader = VideoDataLoader("https://youtu.be/test")
            expected_path = "/tmp/subtitles_test.%(ext)s"

            assert loader._subtitle_template_path == expected_path


def test_video_data_loader_subtitle_prefix() -> None:
    with patch("src.load.video_loader.build_video_source") as mock_build:
        mock_build.return_value = ("https://www.youtube.com/watch?v=test", "test")
        with patch("tempfile.gettempdir") as mock_temp:
            mock_temp.return_value = "/tmp"

            loader = VideoDataLoader("https://youtu.be/test")
            expected_prefix = Path("/tmp/subtitles_test")

            assert loader._subtitle_prefix == expected_prefix


def test_detect_language_with_video_info_language() -> None:
    with patch("src.load.video_loader.build_video_source") as mock_build:
        mock_build.return_value = ("https://www.youtube.com/watch?v=test", "test")

        loader = VideoDataLoader("https://youtu.be/test")

        info = VideoInfo(
            id="test_id",
            language="es-ES",  # Full language code with region
            uploader="test_uploader",
            title="test_title",
            thumbnail="test_thumbnail",
            subtitles={},
        )

        # Should return base language without region
        assert loader._detect_language(info) == "es"


def test_detect_language_with_subtitles_en_priority() -> None:
    with patch("src.load.video_loader.build_video_source") as mock_build:
        mock_build.return_value = ("https://www.youtube.com/watch?v=test", "test")

        loader = VideoDataLoader("https://youtu.be/test")

        info = VideoInfo(
            id="test_id",
            language="",  # No language specified
            uploader="test_uploader",
            title="test_title",
            thumbnail="test_thumbnail",
            subtitles={"en": [], "es": []},  # English available
        )

        # Should prefer English when available
        assert loader._detect_language(info) == "en"


def test_detect_language_with_subtitles_first_available() -> None:
    with patch("src.load.video_loader.build_video_source") as mock_build:
        mock_build.return_value = ("https://www.youtube.com/watch?v=test", "test")

        loader = VideoDataLoader("https://youtu.be/test")

        info = VideoInfo(
            id="test_id",
            language="",  # No language specified
            uploader="test_uploader",
            title="test_title",
            thumbnail="test_thumbnail",
            subtitles={"fr-FR": [], "de": []},  # No English available
        )

        # Should return first available language (fr)
        assert loader._detect_language(info) == "fr"


def test_detect_language_no_subtitles_or_language() -> None:
    with patch("src.load.video_loader.build_video_source") as mock_build:
        mock_build.return_value = ("https://www.youtube.com/watch?v=test", "test")

        loader = VideoDataLoader("https://youtu.be/test")

        info = VideoInfo(
            id="test_id",
            language="",  # No language specified
            uploader="test_uploader",
            title="test_title",
            thumbnail="test_thumbnail",
            subtitles={},  # No subtitles available
        )

        # Should raise RuntimeError when no subtitles available
        try:
            loader._detect_language(info)
            assert False, "Expected RuntimeError"
        except RuntimeError as e:
            assert "no subtitles available" in str(e)


def test_find_subtitle_file_exact_match() -> None:
    with patch("src.load.video_loader.build_video_source") as mock_build:
        mock_build.return_value = ("https://www.youtube.com/watch?v=test", "test")
        with patch("tempfile.gettempdir") as mock_temp:
            mock_temp.return_value = "/tmp"

            loader = VideoDataLoader("https://youtu.be/test")

            with patch("tempfile.gettempdir", return_value="/tmp"):
                with patch.object(Path, "exists", return_value=True):
                    result = loader._find_subtitle_file("en")
                    expected_path = Path("/tmp") / "subtitles_test.en.srt"
                    assert result == expected_path


def test_build_ydl_opts_base_options() -> None:
    with patch("src.load.video_loader.build_video_source") as mock_build:
        mock_build.return_value = ("https://www.youtube.com/watch?v=test", "test")

        loader = VideoDataLoader("https://youtu.be/test")
        opts = loader._build_ydl_opts()

        assert opts["quiet"] is True
        assert opts["no_warnings"] is True
        assert opts["extract_flat"] is False


def test_build_ydl_opts_with_extra_options() -> None:
    with patch("src.load.video_loader.build_video_source") as mock_build:
        mock_build.return_value = ("https://www.youtube.com/watch?v=test", "test")

        loader = VideoDataLoader("https://youtu.be/test")
        opts = loader._build_ydl_opts({"dumpjson": True, "skip_download": True})

        assert opts["dumpjson"] is True
        assert opts["skip_download"] is True


def test_build_ydl_opts_with_user_options() -> None:
    with patch("src.load.video_loader.build_video_source") as mock_build:
        mock_build.return_value = ("https://www.youtube.com/watch?v=test", "test")

        loader = VideoDataLoader(
            "https://youtu.be/test", yt_dlp_additional_options=("--format", "mp4", "--cookies", "cookies.txt")
        )
        opts = loader._build_ydl_opts()

        assert opts["format"] == "mp4"
        assert opts["cookies"] == "cookies.txt"


@patch("yt_dlp.YoutubeDL")
def test_load_success(mock_youtube_dl_class) -> None:
    with patch("src.load.video_loader.build_video_source") as mock_build:
        mock_build.return_value = ("https://www.youtube.com/watch?v=test", "test")

        # Mock the context manager
        mock_ydl = MagicMock()
        mock_ydl.__enter__ = MagicMock(return_value=mock_ydl)
        mock_ydl.__exit__ = MagicMock(return_value=False)
        mock_youtube_dl_class.return_value = mock_ydl

        # First call: extract_info for video info
        # Second call: extract_info for subtitles
        mock_ydl.extract_info.side_effect = [
            {
                "id": "test_id",
                "language": "en",
                "uploader": "test_uploader",
                "title": "test_title",
                "thumbnail": "test_thumbnail",
                "subtitles": {"en": []},
            },
            None,  # Download returns None
        ]

        # Mock subtitle file
        mock_subtitle_file = MagicMock()
        mock_subtitle_file.read_text.return_value = "1\n00:00:00,000 --> 00:00:01,000\nTest subtitle"
        mock_subtitle_file.exists.return_value = True

        with patch.object(VideoDataLoader, "_find_subtitle_file", return_value=mock_subtitle_file):
            with patch("src.load.video_loader.clean_srt", return_value="Test subtitle"):
                loader = VideoDataLoader("https://youtu.be/test")
                loader.load()

                assert loader.info is not None
                assert loader.info.id == "test_id"
                assert loader.transcript is not None
                assert loader.transcript.transcript == "Test subtitle"


@patch("yt_dlp.YoutubeDL")
def test_load_retry_on_info_failure(mock_youtube_dl_class) -> None:
    with patch("src.load.video_loader.build_video_source") as mock_build:
        mock_build.return_value = ("https://www.youtube.com/watch?v=test", "test")

        # Mock the context manager
        mock_ydl = MagicMock()
        mock_ydl.__enter__ = MagicMock(return_value=mock_ydl)
        mock_ydl.__exit__ = MagicMock(return_value=False)
        mock_youtube_dl_class.return_value = mock_ydl

        # Fail twice, succeed on third
        mock_ydl.extract_info.side_effect = [
            Exception("Network error"),
            Exception("Network error"),
            {
                "id": "test_id",
                "language": "en",
                "uploader": "test_uploader",
                "title": "test_title",
                "thumbnail": "test_thumbnail",
                "subtitles": {"en": []},
            },
            None,  # Download returns None
        ]

        # Mock subtitle file
        mock_subtitle_file = MagicMock()
        mock_subtitle_file.read_text.return_value = "1\n00:00:00,000 --> 00:00:01,000\nTest subtitle"
        mock_subtitle_file.exists.return_value = True

        with patch.object(VideoDataLoader, "_find_subtitle_file", return_value=mock_subtitle_file):
            with patch("src.load.video_loader.clean_srt", return_value="Test subtitle"):
                loader = VideoDataLoader("https://youtu.be/test")
                loader.load()

                # Should have been called 4 times (2 info failures + 1 info success + 1 subtitle download)
                assert mock_ydl.extract_info.call_count == 4


@patch("yt_dlp.YoutubeDL")
def test_load_failure_all_attempts(mock_youtube_dl_class) -> None:
    with patch("src.load.video_loader.build_video_source") as mock_build:
        mock_build.return_value = ("https://www.youtube.com/watch?v=test", "test")

        # Mock the context manager
        mock_ydl = MagicMock()
        mock_ydl.__enter__ = MagicMock(return_value=mock_ydl)
        mock_ydl.__exit__ = MagicMock(return_value=False)
        mock_youtube_dl_class.return_value = mock_ydl

        # Always fail
        mock_ydl.extract_info.side_effect = Exception("Network error")

        loader = VideoDataLoader("https://youtu.be/test")
        try:
            loader.load()
            assert False, "Expected RuntimeError"
        except RuntimeError as e:
            assert "Failed to load video info after 3 attempts" in str(e)


@patch("yt_dlp.YoutubeDL")
def test_load_no_subtitles_found(mock_youtube_dl_class) -> None:
    with patch("src.load.video_loader.build_video_source") as mock_build:
        mock_build.return_value = ("https://www.youtube.com/watch?v=test", "test")

        # Mock the context manager
        mock_ydl = MagicMock()
        mock_ydl.__enter__ = MagicMock(return_value=mock_ydl)
        mock_ydl.__exit__ = MagicMock(return_value=False)
        mock_youtube_dl_class.return_value = mock_ydl

        # Info succeeds, download succeeds but no file found
        mock_ydl.extract_info.side_effect = [
            {
                "id": "test_id",
                "language": "en",
                "uploader": "test_uploader",
                "title": "test_title",
                "thumbnail": "test_thumbnail",
                "subtitles": {"en": []},
            },
            None,
        ]

        with patch.object(VideoDataLoader, "_find_subtitle_file", return_value=None):
            loader = VideoDataLoader("https://youtu.be/test")
            try:
                loader.load()
                assert False, "Expected FileNotFoundError"
            except FileNotFoundError as e:
                assert "no subtitles found" in str(e)
