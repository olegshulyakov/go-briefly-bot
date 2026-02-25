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
    await provider.put("summary:hash123:en", "english summary", 3600)
    # Set summary in Spanish
    await provider.put("summary:hash123:es", "spanish summary", 3600)

    # Get English summary
    text_en = await provider.get("summary:hash123:en")
    assert text_en == "english summary"

    # Get Spanish summary
    text_es = await provider.get("summary:hash123:es")
    assert text_es == "spanish summary"

    # Missing language
    text_fr = await provider.get("summary:hash123:fr")
    assert text_fr is None


@pytest.mark.asyncio
async def test_local_summary_expiration(provider):
    # Set with short TTL
    await provider.put("summary:hash123:en", "expiring summary", 0.1)

    # Should exist initially
    assert await provider.get("summary:hash123:en") == "expiring summary"

    # Wait for expiration
    await asyncio.sleep(0.2)

    # Should be gone
    assert await provider.get("summary:hash123:en") is None


@pytest.mark.asyncio
async def test_local_transcript(provider):
    data = {"text": "hello"}

    await provider.put_dict("transcript:hash123", data, 3600)

    result = await provider.get_dict("transcript:hash123")
    assert result == data


@pytest.mark.asyncio
async def test_local_transcript_expiration(provider):
    data = {"text": "expiring"}

    await provider.put_dict("transcript:hash123", data, 0.1)

    assert await provider.get_dict("transcript:hash123") == data

    await asyncio.sleep(0.2)

    assert await provider.get_dict("transcript:hash123") is None
