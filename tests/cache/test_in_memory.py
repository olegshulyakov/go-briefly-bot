import asyncio

import pytest
from src.cache import InMemoryCacheProvider


@pytest.fixture
def provider() -> InMemoryCacheProvider:
    return InMemoryCacheProvider()


@pytest.mark.asyncio
async def test_in_memory_rate_limit(provider: InMemoryCacheProvider) -> None:
    # Not limited initially
    is_limited = await provider.is_rate_limited(123, 10)
    assert not is_limited

    # Limited on next call
    is_limited = await provider.is_rate_limited(123, 10)
    assert is_limited


@pytest.mark.asyncio
async def test_in_memory_summary_multiple_languages(provider: InMemoryCacheProvider) -> None:
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
async def test_in_memory_summary_expiration(provider: InMemoryCacheProvider) -> None:
    # Set with short TTL
    await provider.put("summary:hash123:en", "expiring summary", 1)

    # Should exist initially
    assert await provider.get("summary:hash123:en") == "expiring summary"

    # Wait for expiration
    await asyncio.sleep(1.1)

    # Should be gone
    assert await provider.get("summary:hash123:en") is None


@pytest.mark.asyncio
async def test_in_memory_transcript(provider: InMemoryCacheProvider) -> None:
    data = {"text": "hello"}

    await provider.put_dict("transcript:hash123", data, 3600)

    result = await provider.get_dict("transcript:hash123")
    assert result == data


@pytest.mark.asyncio
async def test_in_memory_transcript_expiration(provider: InMemoryCacheProvider) -> None:
    data = {"text": "expiring"}

    await provider.put_dict("transcript:hash123", data, 1)

    assert await provider.get_dict("transcript:hash123") == data

    await asyncio.sleep(1.1)

    assert await provider.get_dict("transcript:hash123") is None
