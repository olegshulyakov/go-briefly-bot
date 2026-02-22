import asyncio
import pytest

from src.cache import LocalCacheProvider


@pytest.fixture
def provider():
    return LocalCacheProvider()


@pytest.mark.asyncio
async def test_local_rate_limit(provider):
    # Not limited initially
    is_limited = await provider.is_rate_limited(123, 10)
    assert not is_limited

    # Limited on next call
    is_limited = await provider.is_rate_limited(123, 10)
    assert is_limited


@pytest.mark.asyncio
async def test_local_summary_multiple_languages(provider):
    # Set summary in English
    await provider.set_summary("hash123", "en", "english summary", 3600)
    # Set summary in Spanish
    await provider.set_summary("hash123", "ru", "spanish summary", 3600)

    # Get English summary
    summary_en = await provider.get_summary("hash123", "en")
    assert summary_en == "english summary"

    # Get Spanish summary
    summary_ru = await provider.get_summary("hash123", "ru")
    assert summary_ru == "spanish summary"

    # Missing language
    summary_fr = await provider.get_summary("hash123", "fr")
    assert summary_fr is None


@pytest.mark.asyncio
async def test_local_summary_expiration(provider):
    # Set with short TTL
    await provider.set_summary("hash123", "en", "expiring summary", 0.1)

    # Should exist initially
    assert await provider.get_summary("hash123", "en") == "expiring summary"

    # Wait for expiration
    await asyncio.sleep(0.2)

    # Should be gone
    assert await provider.get_summary("hash123", "en") is None


@pytest.mark.asyncio
async def test_local_transcript(provider):
    transcript_data = {"text": "hello"}

    await provider.set_transcript("hash123", transcript_data, 3600)

    result = await provider.get_transcript("hash123")
    assert result == transcript_data


@pytest.mark.asyncio
async def test_local_transcript_expiration(provider):
    transcript_data = {"text": "expiring"}

    await provider.set_transcript("hash123", transcript_data, 0.1)

    assert await provider.get_transcript("hash123") == transcript_data

    await asyncio.sleep(0.2)

    assert await provider.get_transcript("hash123") is None
