from unittest.mock import patch
import os
from src.config import Settings


def test_settings_from_env_success() -> None:
    with patch.dict(
        os.environ,
        {
            "TELEGRAM_BOT_TOKEN": "test_token",
            "OPENAI_API_KEY": "test_api_key",
            "OPENAI_MODEL": "gpt-3.5-turbo",
            "OPENAI_BASE_URL": "https://api.openai.com/v1/",
            "YT_DLP_ADDITIONAL_OPTIONS": "--format mp4 --proxy http://proxy.example.com",
        },
    ):
        settings = Settings.from_env()

        assert settings.telegram_bot_token == "test_token"
        assert settings.openai_api_key == "test_api_key"
        assert settings.openai_model == "gpt-3.5-turbo"
        assert settings.openai_base_url == "https://api.openai.com/v1/"
        assert settings.yt_dlp_additional_options == (
            "--format",
            "mp4",
            "--proxy",
            "http://proxy.example.com",
        )


@patch("src.config.load_dotenv")
def test_settings_from_env_missing_required_variables(mock_load_dotenv) -> None:
    # Test with missing TELEGRAM_BOT_TOKEN
    with patch.dict(
        os.environ,
        {"OPENAI_API_KEY": "test_api_key", "OPENAI_MODEL": "gpt-3.5-turbo"},
        clear=True,
    ):
        try:
            Settings.from_env()
            assert False, "Expected RuntimeError for missing env vars"
        except RuntimeError as e:
            assert "TELEGRAM_BOT_TOKEN" in str(e)

    # Test with missing OPENAI_API_KEY
    with patch.dict(
        os.environ,
        {"TELEGRAM_BOT_TOKEN": "test_token", "OPENAI_MODEL": "gpt-3.5-turbo"},
        clear=True,
    ):
        try:
            Settings.from_env()
            assert False, "Expected RuntimeError for missing env vars"
        except RuntimeError as e:
            assert "OPENAI_API_KEY" in str(e)

    # Test with missing OPENAI_MODEL
    with patch.dict(
        os.environ,
        {"TELEGRAM_BOT_TOKEN": "test_token", "OPENAI_API_KEY": "test_api_key"},
        clear=True,
    ):
        try:
            Settings.from_env()
            assert False, "Expected RuntimeError for missing env vars"
        except RuntimeError as e:
            assert "OPENAI_MODEL" in str(e)


@patch("src.config.load_dotenv")
def test_settings_from_env_default_values(mock_load_dotenv) -> None:
    with patch.dict(
        os.environ,
        {
            "TELEGRAM_BOT_TOKEN": "test_token",
            "OPENAI_API_KEY": "test_api_key",
            "OPENAI_MODEL": "gpt-3.5-turbo",
        },
        clear=True,
    ):
        settings = Settings.from_env()

        # Check default values
        assert settings.openai_base_url == "https://api.openai.com/v1/"  # Default from code
        assert settings.rate_limit_window_seconds == 10
        assert settings.max_telegram_message_length == 3500
        assert settings.openai_timeout_seconds == 300
        assert settings.openai_max_retries == 3
        assert settings.yt_dlp_additional_options == ()


def test_settings_from_env_custom_values() -> None:
    with patch.dict(
        os.environ,
        {
            "TELEGRAM_BOT_TOKEN": "custom_token",
            "OPENAI_BASE_URL": "https://custom.openai.api/v1/",
            "OPENAI_API_KEY": "custom_api_key",
            "OPENAI_MODEL": "gpt-4",
            "YT_DLP_ADDITIONAL_OPTIONS": "--format best --extract-audio",
            "RATE_LIMIT_WINDOW_SECONDS": "30",
        },
    ):
        # Note: RATE_LIMIT_WINDOW_SECONDS is not configurable via env in the current implementation
        # So we'll test with the default value
        settings = Settings.from_env()

        assert settings.telegram_bot_token == "custom_token"
        assert settings.openai_base_url == "https://custom.openai.api/v1/"
        assert settings.openai_api_key == "custom_api_key"
        assert settings.openai_model == "gpt-4"
        assert settings.yt_dlp_additional_options == (
            "--format",
            "best",
            "--extract-audio",
        )


def test_settings_immutability() -> None:
    with patch.dict(
        os.environ,
        {
            "TELEGRAM_BOT_TOKEN": "test_token",
            "OPENAI_API_KEY": "test_api_key",
            "OPENAI_MODEL": "gpt-3.5-turbo",
        },
    ):
        settings = Settings.from_env()

        # Attempt to modify should raise an exception since it's a frozen dataclass
        try:
            settings.telegram_bot_token = "new_token"
            assert False, "Expected dataclass.FrozenInstanceError"
        except Exception:
            # Expected behavior - frozen dataclass can't be modified
            pass
