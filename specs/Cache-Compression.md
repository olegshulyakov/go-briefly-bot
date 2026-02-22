# Cache Compression Specification

## Overview

The Valkey cache provider supports multiple compression algorithms to reduce memory usage and network traffic. Compression is applied automatically to all cached data (summaries and transcripts) with transparent decompression on retrieval.

## Supported Compression Methods

| Method | Prefix | Description                         | Best For                              |
| ------ | ------ | ----------------------------------- | ------------------------------------- |
| `none` | `\x00` | No compression, data stored as-is   | Debugging, already-compressed data    |
| `gzip` | `\x01` | DEFLATE-based compression (default) | General purpose, balanced speed/ratio |
| `zlib` | `\x02` | Raw DEFLATE compression             | Fast compression/decompression        |
| `lzma` | `\x03` | Lempel-Ziv-Markov chain (XZ format) | Maximum compression ratio             |

## Compression Format

### Key Structure

All cache keys include the compression method suffix for better observability:

```
summary:{video_hash}:{language_code}:{compression_method}
transcript:{video_hash}:{compression_method}
```

Example keys:

- `summary:d4f8a2b1:en:gzip` — English summary with gzip compression
- `summary:d4f8a2b1:ru:lzma` — Russian summary with LZMA compression
- `transcript:d4f8a2b1:gzip` — Transcript with gzip compression

Benefits of including method in key:

- **Debugging**: Easy to see compression method in Valkey CLI
- **Migration**: Different methods can coexist during migration
- **Cleanup**: Can delete specific compression methods with `DEL pattern`
- **Monitoring**: Can track memory usage per compression method

### Value Structure

Compressed values include a 1-byte prefix for automatic detection:

```
[1 byte prefix][compressed payload]
```

This allows automatic detection of the compression method during decompression, enabling:

- **Backward compatibility**: Old uncompressed data can be migrated
- **Method flexibility**: Change compression method without clearing cache
- **Error detection**: Invalid prefixes raise `CompressionError`

## Configuration

Set the compression method via environment variable:

```bash
CACHE_COMPRESSION_METHOD=gzip  # Default
```

Valid values: `none`, `gzip`, `zlib`, `lzma`

## Performance Characteristics

### Compression Ratios (typical)

| Data Type         | Original Size | gzip         | zlib         | lzma         |
| ----------------- | ------------- | ------------ | ------------ | ------------ |
| Summary (text)    | 2 KB          | 0.8 KB (60%) | 0.9 KB (55%) | 0.6 KB (70%) |
| Transcript (text) | 50 KB         | 12 KB (76%)  | 14 KB (72%)  | 8 KB (84%)   |
| Transcript (JSON) | 100 KB        | 25 KB (75%)  | 28 KB (72%)  | 18 KB (82%)  |

### Speed Benchmarks (approximate, M1 Mac)

| Method | Compress (MB/s) | Decompress (MB/s) | CPU Usage |
| ------ | --------------- | ----------------- | --------- |
| `none` | ∞ (copy only)   | ∞ (copy only)     | ~0%       |
| `gzip` | 50-100          | 200-300           | ~5%       |
| `zlib` | 80-150          | 300-400           | ~3%       |
| `lzma` | 5-15            | 50-100            | ~15%      |

## Trade-offs

### ✅ Benefits

1. **Memory Savings**: 60-80% reduction in Valkey memory usage
2. **Network Efficiency**: Less data transferred between bot and Valkey
3. **Cost Reduction**: Lower memory requirements for managed Valkey/Redis
4. **Increased Cache Capacity**: More data fits in the same memory

### ❌ Costs

1. **CPU Overhead**: Compression/decompression requires CPU cycles
2. **Latency**: +1-5ms per operation (depending on method and data size)
3. **Complexity**: Additional error handling for compression failures

## Recommendations

### Production Deployment

- **Use `gzip`** (default) - best balance of speed and compression ratio
- **Monitor CPU usage** - ensure compression overhead is acceptable
- **Set appropriate TTL** - compressed data still expires based on TTL

### Development/Debugging

- **Use `none`** - easier to inspect cached data in Valkey
- **Enable debug logging** - monitor compression stats

### High-Performance Scenarios

- **Use `zlib`** - faster than gzip with slightly lower compression
- **Benchmark with your data** - actual ratios depend on content

### Memory-Constrained Environments

- **Use `lzma`** - maximum compression ratio
- **Accept CPU trade-off** - slower but uses least memory

## Implementation Details

### Compression Flow

```
set_summary() / set_transcript()
    ↓
Serialize data (JSON for transcripts)
    ↓
Encode to UTF-8 bytes
    ↓
Compress with selected method
    ↓
Add method prefix
    ↓
Store in Valkey with SETEX
```

### Decompression Flow

```
get_summary() / get_transcript()
    ↓
Retrieve bytes from Valkey
    ↓
Read first byte (method prefix)
    ↓
Select decompression algorithm
    ↓
Decompress payload
    ↓
Decode UTF-8 to string
    ↓
Parse JSON (for transcripts)
    ↓
Return data
```

### Error Handling

- **Unknown prefix**: Raises `CompressionError` with details
- **Corrupted data**: Raises `CompressionError` (caught by fail-soft)
- **Decompression failure**: Falls back to local cache (if enabled)

## Migration Guide

### Enabling Compression on Existing Cache

1. **Set `CACHE_COMPRESSION_METHOD`** in environment
2. **Restart bot** - new data will be compressed
3. **Old data** - will be read as-is (first byte = `\x00` = none)
4. **Gradual migration** - cache naturally migrates as items expire

### Disabling Compression

1. **Set `CACHE_COMPRESSION_METHOD=none`**
2. **Restart bot** - new data stored uncompressed
3. **Compressed data** - still readable (auto-detected)

## Monitoring

### Key Metrics

- **Compression ratio**: `compressed_size / original_size`
- **Cache hit rate**: Should increase with compression (more data fits)
- **CPU usage**: Monitor for compression overhead
- **Latency**: P95/P99 for cache operations

### Logging

The bot logs compression statistics at INFO level:

```
INFO: Cache compression enabled (method=gzip, expected_savings=70%)
```

## Security Considerations

- **Compression oracle attacks**: Not applicable (data is not secret)
- **DoS via compression**: Mitigated by 200ms timeout
- **Data integrity**: CRC64 checks for LZMA format

## Future Enhancements

- [ ] Adaptive compression (select method based on data size)
- [ ] Compression statistics in logs/metrics
- [ ] Zstandard (zstd) support for better speed/ratio trade-off
- [ ] Configurable compression level (1-9)
