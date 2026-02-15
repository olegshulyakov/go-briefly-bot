from unittest.mock import patch, MagicMock
import os
import logging
from src.main import configure_logging, main


def test_configure_logging_default() -> None:
    # Ensure no LOG_LEVEL is set
    if "LOG_LEVEL" in os.environ:
        del os.environ["LOG_LEVEL"]

    # Capture the basicConfig call
    with patch("logging.basicConfig") as mock_basic_config:
        configure_logging()

        # Verify basicConfig was called with default INFO level
        mock_basic_config.assert_called_once()
        args, kwargs = mock_basic_config.call_args
        assert kwargs["level"] == logging.INFO


def test_configure_logging_custom_level() -> None:
    # Set a custom log level
    with patch.dict(os.environ, {"LOG_LEVEL": "DEBUG"}):
        with patch("logging.basicConfig") as mock_basic_config:
            configure_logging()

            # Verify basicConfig was called with DEBUG level
            mock_basic_config.assert_called_once()
            args, kwargs = mock_basic_config.call_args
            assert kwargs["level"] == logging.DEBUG


def test_configure_logging_invalid_level() -> None:
    # Set an invalid log level - should default to INFO
    with patch.dict(os.environ, {"LOG_LEVEL": "INVALID_LEVEL"}):
        with patch("logging.basicConfig") as mock_basic_config:
            configure_logging()

            # Verify basicConfig was called with default INFO level
            mock_basic_config.assert_called_once()
            args, kwargs = mock_basic_config.call_args
            assert kwargs["level"] == logging.INFO


def test_main_function_calls_configure_and_runs() -> None:
    # Mock all dependencies to test the main flow
    with (
        patch("src.main.configure_logging") as mock_configure_logging,
        patch("src.main.Settings") as mock_settings,
        patch("src.main.TelegramBrieflyBot") as mock_bot_class,
    ):
        # Create mock objects
        mock_settings_obj = MagicMock()
        mock_bot_obj = MagicMock()

        mock_settings.from_env.return_value = mock_settings_obj
        mock_bot_class.return_value = mock_bot_obj

        # Call main function
        main()

        # Verify all steps were called in sequence
        mock_configure_logging.assert_called_once()
        mock_settings.from_env.assert_called_once()
        mock_bot_class.assert_called_once_with(mock_settings_obj)
        mock_bot_obj.run.assert_called_once()


def test_main_function_structure() -> None:
    # Just test that the main function exists and has the expected structure
    # by checking it has the expected attributes
    assert callable(main)
    assert hasattr(main, "__name__")
    assert main.__name__ == "main"
