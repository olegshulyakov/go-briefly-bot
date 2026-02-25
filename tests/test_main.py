import logging
import os
from unittest.mock import MagicMock, patch

from src.main import configure_logging, main


def test_configure_logging_default() -> None:
    # Ensure no LOG_LEVEL is set
    with (
        patch.dict(os.environ, {}, clear=True),
        patch("src.logger.logging.getLogger") as mock_get_logger,
        patch("src.logger.logging.StreamHandler") as mock_handler_class,
        patch("src.logger.os.getenv", return_value="INFO"),
    ):
        mock_root_logger = MagicMock()
        mock_httpx_logger = MagicMock()
        mock_handler = MagicMock()

        # Return different loggers based on name
        def get_logger_side_effect(name=None):
            if name == "httpx":
                return mock_httpx_logger
            return mock_root_logger

        mock_get_logger.side_effect = get_logger_side_effect
        mock_handler_class.return_value = mock_handler
        mock_root_logger.handlers = []

        configure_logging()

        # Verify root logger level was set to INFO (20)
        mock_root_logger.setLevel.assert_called_with(logging.INFO)
        # Verify handler was added
        mock_root_logger.addHandler.assert_called_once_with(mock_handler)
        # Verify handler level was set
        mock_handler.setLevel.assert_called_with(logging.INFO)
        # Verify formatter was set
        mock_handler.setFormatter.assert_called_once()
        # Verify httpx logger level was set to WARNING
        mock_httpx_logger.setLevel.assert_called_with(logging.WARNING)


def test_configure_logging_custom_level() -> None:
    # Set a custom log level
    with (
        patch.dict(os.environ, {"LOG_LEVEL": "DEBUG"}),
        patch("src.logger.logging.getLogger") as mock_get_logger,
        patch("src.logger.logging.StreamHandler") as mock_handler_class,
        patch("src.logger.os.getenv", return_value="DEBUG"),
    ):
        mock_root_logger = MagicMock()
        mock_httpx_logger = MagicMock()
        mock_handler = MagicMock()

        # Return different loggers based on name
        def get_logger_side_effect(name=None):
            if name == "httpx":
                return mock_httpx_logger
            return mock_root_logger

        mock_get_logger.side_effect = get_logger_side_effect
        mock_handler_class.return_value = mock_handler
        mock_root_logger.handlers = []

        configure_logging()

        # Verify root logger level was set to DEBUG (10)
        mock_root_logger.setLevel.assert_called_with(logging.DEBUG)
        # Verify handler was added
        mock_root_logger.addHandler.assert_called_once_with(mock_handler)
        # Verify handler level was set
        mock_handler.setLevel.assert_called_with(logging.DEBUG)


def test_configure_logging_invalid_level() -> None:
    # Set an invalid log level - should default to INFO
    with (
        patch.dict(os.environ, {"LOG_LEVEL": "INVALID_LEVEL"}),
        patch("src.logger.logging.getLogger") as mock_get_logger,
        patch("src.logger.logging.StreamHandler") as mock_handler_class,
        patch("src.logger.os.getenv", return_value="INVALID_LEVEL"),
    ):
        mock_root_logger = MagicMock()
        mock_httpx_logger = MagicMock()
        mock_handler = MagicMock()

        # Return different loggers based on name
        def get_logger_side_effect(name=None):
            if name == "httpx":
                return mock_httpx_logger
            return mock_root_logger

        mock_get_logger.side_effect = get_logger_side_effect
        mock_handler_class.return_value = mock_handler
        mock_root_logger.handlers = []

        configure_logging()

        # Verify root logger level was set to INFO (default for invalid level)
        mock_root_logger.setLevel.assert_called_with(logging.INFO)
        # Verify handler was added
        mock_root_logger.addHandler.assert_called_once_with(mock_handler)


def test_configure_logging_removes_existing_handlers() -> None:
    # Test that existing handlers are removed before adding new one
    with (
        patch("src.logger.logging.getLogger") as mock_get_logger,
        patch("src.logger.logging.StreamHandler") as mock_handler_class,
        patch("src.logger.os.getenv", return_value="INFO"),
    ):
        mock_root_logger = MagicMock()
        mock_httpx_logger = MagicMock()
        mock_existing_handler = MagicMock()
        mock_handler = MagicMock()

        def get_logger_side_effect(name=None):
            if name == "httpx":
                return mock_httpx_logger
            return mock_root_logger

        mock_get_logger.side_effect = get_logger_side_effect
        mock_handler_class.return_value = mock_handler
        mock_root_logger.handlers = [mock_existing_handler]

        configure_logging()

        # Verify existing handler was removed
        mock_root_logger.removeHandler.assert_called_with(mock_existing_handler)
        # Verify new handler was added
        mock_root_logger.addHandler.assert_called_with(mock_handler)


def test_configure_logging_uses_custom_formatter() -> None:
    # Test that CustomFormatter is used
    with (
        patch("src.logger.logging.getLogger") as mock_get_logger,
        patch("src.logger.logging.StreamHandler") as mock_handler_class,
        patch("src.logger.CustomFormatter") as mock_formatter_class,
        patch("src.logger.os.getenv", return_value="INFO"),
    ):
        mock_root_logger = MagicMock()
        mock_httpx_logger = MagicMock()
        mock_handler = MagicMock()
        mock_formatter = MagicMock()

        def get_logger_side_effect(name=None):
            if name == "httpx":
                return mock_httpx_logger
            return mock_root_logger

        mock_get_logger.side_effect = get_logger_side_effect
        mock_handler_class.return_value = mock_handler
        mock_formatter_class.return_value = mock_formatter
        mock_root_logger.handlers = []

        configure_logging()

        # Verify CustomFormatter was created with correct format
        mock_formatter_class.assert_called_once()
        # Verify formatter was set on handler
        mock_handler.setFormatter.assert_called_once_with(mock_formatter)


def test_configure_logging_httpx_warning() -> None:
    # Test that httpx logger is set to WARNING
    with (
        patch("src.logger.logging.getLogger") as mock_get_logger,
        patch("src.logger.logging.StreamHandler") as mock_handler_class,
    ):
        mock_root_logger = MagicMock()
        mock_httpx_logger = MagicMock()
        mock_handler = MagicMock()

        def get_logger_side_effect(name=None):
            if name == "httpx":
                return mock_httpx_logger
            return mock_root_logger

        mock_get_logger.side_effect = get_logger_side_effect
        mock_handler_class.return_value = mock_handler
        mock_root_logger.handlers = []

        configure_logging()

        # Verify httpx logger level was set to WARNING
        mock_httpx_logger.setLevel.assert_called_with(logging.WARNING)


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
