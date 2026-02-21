import asyncio
from unittest.mock import AsyncMock, MagicMock, patch

import pytest
from valkey.exceptions import ConnectionError as ValkeyConnectionError

from src.storage import ValkeyProvider


@pytest.fixture
def provider():
    return ValkeyProvider("valkey://localhost:6379")


@pytest.mark.asyncio
async def test_valkey_provider_rate_limit_success(provider):
    with patch("src.storage.valkey_provider.Valkey") as mock_valkey:
        mock_client = AsyncMock()
        mock_pipeline = MagicMock()
        mock_pipeline.execute = AsyncMock(return_value=[1])  # INCR result
        mock_client.pipeline = MagicMock(return_value=mock_pipeline)
        mock_valkey.from_url.return_value = mock_client

        is_limited = await provider.is_rate_limited(123, 10)

        assert not is_limited
        mock_client.pipeline.assert_called_once()
        mock_pipeline.incr.assert_called_with("user:123:limit")
        mock_pipeline.expire.assert_called_with("user:123:limit", 10, nx=True)


@pytest.mark.asyncio
async def test_valkey_provider_rate_limit_exceeded(provider):
    with patch("src.storage.valkey_provider.Valkey") as mock_valkey:
        mock_client = AsyncMock()
        mock_pipeline = MagicMock()
        mock_pipeline.execute = AsyncMock(return_value=[2])  # INCR result > 1
        mock_client.pipeline = MagicMock(return_value=mock_pipeline)
        mock_valkey.from_url.return_value = mock_client

        is_limited = await provider.is_rate_limited(123, 10)

        assert is_limited


@pytest.mark.asyncio
async def test_valkey_provider_fail_soft_timeout(provider):
    # Reduce timeout for faster test
    provider._timeout = 0.01

    with patch("src.storage.valkey_provider.Valkey") as mock_valkey:
        mock_client = AsyncMock()

        async def slow_execute():
            await asyncio.sleep(0.05)
            return [1]

        mock_pipeline = MagicMock()
        mock_pipeline.execute = slow_execute
        mock_client.pipeline = MagicMock(return_value=mock_pipeline)
        mock_valkey.from_url.return_value = mock_client

        # Should fall back to local provider and return False (not limited)
        is_limited = await provider.is_rate_limited(123, 10)

        assert not is_limited
        # Verify it fallback logic works by checking if it cached in local
        assert 123 in provider._local_fallback._rate_limits


@pytest.mark.asyncio
async def test_valkey_provider_fail_soft_connection_error(provider):
    with patch("src.storage.valkey_provider.Valkey") as mock_valkey:
        mock_client = AsyncMock()
        mock_pipeline = MagicMock()
        mock_pipeline.execute = AsyncMock(side_effect=ValkeyConnectionError("Connection refused"))
        mock_client.pipeline = MagicMock(return_value=mock_pipeline)
        mock_valkey.from_url.return_value = mock_client

        # Should fall back to local provider without raising an error
        is_limited = await provider.is_rate_limited(456, 10)

        assert not is_limited
        assert 456 in provider._local_fallback._rate_limits


@pytest.mark.asyncio
async def test_valkey_provider_set_get_summary(provider):
    with patch("src.storage.valkey_provider.Valkey") as mock_valkey:
        mock_client = AsyncMock()
        mock_client.get.return_value = b"test summary"
        mock_valkey.from_url.return_value = mock_client

        # Test set
        await provider.set_summary("hash123", "test summary", 3600)
        mock_client.setex.assert_called_with("summary:hash123", 3600, "test summary")

        # Test get
        summary = await provider.get_summary("hash123")
        assert summary == "test summary"
        mock_client.get.assert_called_with("summary:hash123")


@pytest.mark.asyncio
async def test_valkey_provider_set_get_transcript(provider):
    with patch("src.storage.valkey_provider.Valkey") as mock_valkey:
        mock_client = AsyncMock()
        mock_client.get.return_value = b'{"text": "hello"}'
        mock_valkey.from_url.return_value = mock_client

        # Test set
        await provider.set_transcript("hash123", {"text": "hello"}, 3600)
        mock_client.setex.assert_called_with("transcript:hash123", 3600, '{"text": "hello"}')

        # Test get
        transcript = await provider.get_transcript("hash123")
        assert transcript == {"text": "hello"}
        mock_client.get.assert_called_with("transcript:hash123")
