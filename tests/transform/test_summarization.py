from unittest.mock import patch, MagicMock
from src.transform.summarization import OpenAISummarizer
from src.config import Settings


def test_openai_summarizer_initialization() -> None:
    mock_settings = MagicMock(spec=Settings)
    mock_settings.openai_base_url = "https://api.openai.com/v1/"
    mock_settings.openai_api_key = "test-key"
    mock_settings.openai_max_retries = 3

    with patch("src.transform.summarization.OpenAI") as mock_openai_class:
        mock_client_instance = MagicMock()
        mock_openai_class.return_value = mock_client_instance

        summarizer = OpenAISummarizer(mock_settings)

        # Verify the client was initialized with correct settings
        mock_openai_class.assert_called_once_with(
            base_url="https://api.openai.com/v1/",
            api_key="test-key",
            max_retries=3,
        )
        assert summarizer.settings == mock_settings


def test_summarizer_summarize_text_success() -> None:
    mock_settings = MagicMock(spec=Settings)
    mock_settings.openai_base_url = "https://api.openai.com/v1/"
    mock_settings.openai_api_key = "test-key"
    mock_settings.openai_model = "gpt-3.5-turbo"
    mock_settings.openai_timeout_seconds = 300
    mock_settings.openai_max_retries = 3

    with patch("src.transform.summarization.OpenAI") as mock_openai_class:
        mock_client_instance = MagicMock()
        mock_response = MagicMock()
        mock_choice = MagicMock()

        mock_choice.message.content = "This is the summary"
        mock_response.choices = [mock_choice]
        mock_client_instance.chat.completions.create.return_value = mock_response

        mock_openai_class.return_value = mock_client_instance

        # Mock the translate function to return a fixed system prompt
        with patch("src.transform.summarization.translate") as mock_translate:
            mock_translate.return_value = "Input text to summarize"

            summarizer = OpenAISummarizer(mock_settings)
            result = summarizer.summarize_text("Input text to summarize", "en")

            # Verify the API was called correctly
            mock_client_instance.chat.completions.create.assert_called_once_with(
                model="gpt-3.5-turbo",
                messages=[
                    {"role": "user", "content": "Input text to summarize"},
                ],
                timeout=300,
            )

            assert result == "This is the summary"


def test_summarizer_summarize_text_empty_response() -> None:
    mock_settings = MagicMock(spec=Settings)
    mock_settings.openai_base_url = "https://api.openai.com/v1/"
    mock_settings.openai_api_key = "test-key"
    mock_settings.openai_model = "gpt-3.5-turbo"
    mock_settings.openai_timeout_seconds = 300
    mock_settings.openai_max_retries = 3

    with patch("src.transform.summarization.OpenAI") as mock_openai_class:
        mock_client_instance = MagicMock()
        mock_response = MagicMock()
        mock_choice = MagicMock()
        mock_message = MagicMock()

        mock_message.content = None  # Empty response
        mock_choice.message = mock_message
        mock_response.choices = [mock_choice]
        # Set response attributes to avoid dynamic mock creation during logging
        mock_response.id = "test_response_id"
        mock_response.model = "gpt-3.5-turbo"
        mock_response.message = None  # Explicitly set to None to avoid mock creation
        mock_client_instance.chat.completions.create.return_value = mock_response

        mock_openai_class.return_value = mock_client_instance

        # Mock the translate function
        with patch("src.transform.summarization.translate") as mock_translate:
            mock_translate.return_value = "You are a helpful assistant."

            summarizer = OpenAISummarizer(mock_settings)

            try:
                summarizer.summarize_text("Input text to summarize", "en")
                assert False, "Expected RuntimeError for empty response"
            except RuntimeError as e:
                assert "empty OpenAI response" in str(e)


def test_summarizer_summarize_text_with_retry_success() -> None:
    mock_settings = MagicMock(spec=Settings)
    mock_settings.openai_base_url = "https://api.openai.com/v1/"
    mock_settings.openai_api_key = "test-key"
    mock_settings.openai_model = "gpt-3.5-turbo"
    mock_settings.openai_timeout_seconds = 300
    mock_settings.openai_max_retries = 3

    with patch("src.transform.summarization.OpenAI") as mock_openai_class:
        mock_client_instance = MagicMock()
        mock_response = MagicMock()
        mock_choice = MagicMock()

        mock_choice.message.content = "This is the summary"
        mock_response.choices = [mock_choice]

        # First call raises an exception, second succeeds
        mock_client_instance.chat.completions.create.side_effect = [
            Exception("Network error"),
            mock_response,
        ]

        mock_openai_class.return_value = mock_client_instance

        # Mock the translate function
        with patch("src.transform.summarization.translate") as mock_translate:
            mock_translate.return_value = "You are a helpful assistant."

            summarizer = OpenAISummarizer(mock_settings)
            result = summarizer.summarize_text("Input text to summarize", "en")

            # Should have been called twice (first failed, second succeeded)
            assert mock_client_instance.chat.completions.create.call_count == 2
            assert result == "This is the summary"


def test_summarizer_summarize_text_with_retry_failure() -> None:
    mock_settings = MagicMock(spec=Settings)
    mock_settings.openai_base_url = "https://api.openai.com/v1/"
    mock_settings.openai_api_key = "test-key"
    mock_settings.openai_model = "gpt-3.5-turbo"
    mock_settings.openai_timeout_seconds = 300
    mock_settings.openai_max_retries = 2  # Limit retries to 2

    with patch("src.transform.summarization.OpenAI") as mock_openai_class:
        mock_client_instance = MagicMock()

        # All calls raise exceptions
        mock_client_instance.chat.completions.create.side_effect = [
            Exception("Network error 1"),
            Exception("Network error 2"),
        ]

        mock_openai_class.return_value = mock_client_instance

        # Mock the translate function
        with patch("src.transform.summarization.translate") as mock_translate:
            mock_translate.return_value = "You are a helpful assistant."

            summarizer = OpenAISummarizer(mock_settings)

            try:
                summarizer.summarize_text("Input text to summarize", "en")
                assert False, "Expected RuntimeError after retries exhausted"
            except RuntimeError as e:
                assert "failed to summarize text:" in str(e)
                # Should have been called max_retries times
                assert mock_client_instance.chat.completions.create.call_count == 2
