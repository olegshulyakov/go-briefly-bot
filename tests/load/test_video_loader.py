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


def test_extract_json_payload_simple() -> None:
    output = '{"id": "test", "title": "Test Video"}'

    with patch("src.load.video_loader.build_video_source") as mock_build:
        mock_build.return_value = ("https://www.youtube.com/watch?v=test", "test")

        loader = VideoDataLoader("https://youtu.be/test")
        result = loader._extract_json_payload(output)

        assert result == '{"id": "test", "title": "Test Video"}'


def test_extract_json_payload_with_extra_lines() -> None:
    output = """[info] Some info message
{"id": "test", "title": "Test Video"}
[download] Completed download"""

    with patch("src.load.video_loader.build_video_source") as mock_build:
        mock_build.return_value = ("https://www.youtube.com/watch?v=test", "test")

        loader = VideoDataLoader("https://youtu.be/test")
        result = loader._extract_json_payload(output)

        assert result == '{"id": "test", "title": "Test Video"}'


def test_extract_json_payload_empty_output() -> None:
    with patch("src.load.video_loader.build_video_source") as mock_build:
        mock_build.return_value = ("https://www.youtube.com/watch?v=test", "test")

        loader = VideoDataLoader("https://youtu.be/test")
        try:
            loader._extract_json_payload("")
            assert False, "Expected ValueError for empty output"
        except ValueError as e:
            assert "empty yt-dlp metadata output" in str(e)


def test_extract_json_payload_last_json_object() -> None:
    output = """[info] Processing...
{"version": "1.0"}
{"id": "test", "title": "Test Video"}"""

    with patch("src.load.video_loader.build_video_source") as mock_build:
        mock_build.return_value = ("https://www.youtube.com/watch?v=test", "test")

        loader = VideoDataLoader("https://youtu.be/test")
        result = loader._extract_json_payload(output)

        assert result == '{"id": "test", "title": "Test Video"}'


@patch("subprocess.run")
def test_exec_success_single_attempt(mock_run) -> None:
    mock_result = MagicMock()
    mock_result.stdout = '{"id": "test"}'
    mock_result.stderr = ""
    mock_result.returncode = 0
    mock_run.return_value = mock_result

    with patch("src.load.video_loader.build_video_source") as mock_build:
        mock_build.return_value = ("https://www.youtube.com/watch?v=test", "test")

        loader = VideoDataLoader("https://youtu.be/test")
        result = loader._exec(["--dump-json"], "https://youtu.be/test")

        mock_run.assert_called_once()
        args, kwargs = mock_run.call_args
        cmd_list = args[0]  # First argument is the command list
        assert "yt-dlp" in cmd_list[0]  # First element should be the command
        assert "--dump-json" in cmd_list  # Option should be in the command list
        assert result == '{"id": "test"}'


@patch("subprocess.run")
def test_exec_success_with_additional_options(mock_run) -> None:
    mock_result = MagicMock()
    mock_result.stdout = '{"id": "test"}'
    mock_result.stderr = ""
    mock_run.return_value = mock_result

    with patch("src.load.video_loader.build_video_source") as mock_build:
        mock_build.return_value = ("https://www.youtube.com/watch?v=test", "test")

        loader = VideoDataLoader("https://youtu.be/test", yt_dlp_additional_options=("--format", "mp4"))
        result = loader._exec(["--dump-json"], "https://youtu.be/test")

        mock_run.assert_called_once()
        args, kwargs = mock_run.call_args
        # Should include additional options
        full_cmd = args[0]
        assert "--format" in full_cmd
        assert "mp4" in full_cmd
        assert "--dump-json" in full_cmd
        assert "https://youtu.be/test" in full_cmd
        assert result == '{"id": "test"}'


@patch("subprocess.run")
def test_exec_failure_eventually_succeeds(mock_run) -> None:
    # Create a mock that raises CalledProcessError on first call, succeeds on second
    def run_side_effect(*args, **kwargs):
        if run_side_effect.call_count == 1:
            run_side_effect.call_count += 1
            # Raise CalledProcessError for first call
            from subprocess import CalledProcessError

            raise CalledProcessError(returncode=1, cmd=args[0])
        else:
            # Return success result for subsequent calls
            result = MagicMock()
            result.stdout = '{"id": "test"}'
            result.stderr = ""
            result.returncode = 0
            result.check_returncode.return_value = None
            return result

    run_side_effect.call_count = 1
    mock_run.side_effect = run_side_effect

    with patch("src.load.video_loader.build_video_source") as mock_build:
        mock_build.return_value = ("https://www.youtube.com/watch?v=test", "test")

        loader = VideoDataLoader("https://youtu.be/test")
        result = loader._exec(["--dump-json"], "https://youtu.be/test")

        assert mock_run.call_count == 2  # Called twice
        assert result == '{"id": "test"}'


@patch("subprocess.run")
def test_exec_failure_all_attempts(mock_run) -> None:
    # All calls should raise CalledProcessError
    def run_side_effect(*args, **kwargs):
        from subprocess import CalledProcessError

        raise CalledProcessError(returncode=1, cmd=args[0])

    mock_run.side_effect = run_side_effect

    with patch("src.load.video_loader.build_video_source") as mock_build:
        mock_build.return_value = ("https://www.youtube.com/watch?v=test", "test")

        loader = VideoDataLoader("https://youtu.be/test")
        try:
            loader._exec(["--dump-json"], "https://youtu.be/test")
            assert False, "Expected RuntimeError after all attempts fail"
        except RuntimeError as e:
            assert "yt-dlp failed after 3 attempts" in str(e)
            assert mock_run.call_count == 3
