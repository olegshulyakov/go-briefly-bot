## Valkey Integration & State Distribution

### 1. Falsifiable Thesis (The Spine)

Integrating Valkey as a pluggable state provider will reduce LLM API overhead by a measurable percentage for duplicate video requests and enable -node horizontal scaling without violating per-user rate limits. If the latency overhead of Valkey connection timeouts exceeds the time saved by cache hits, the integration is functionally net-negative.

### 2. Implementation Timeline

| Milestone                  | Date       | Deliverable                                                |
| -------------------------- | ---------- | ---------------------------------------------------------- |
| **Phase 1: Discovery**     | 2026-02-22 | Finalize `StorageInterface` abstraction.                   |
| **Phase 2: Core Cache**    | 2026-02-24 | Implement `ValkeyCache` with `SETEX` logic for summaries.  |
| **Phase 3: Rate Limiting** | 2026-02-26 | Refactor `RateLimiter` to use atomic increments in Valkey. |
| **Phase 4: Resilience**    | 2026-02-28 | Integration tests for connection fail-soft behavior.       |

---

### 3. Technical Hurdles & Solutions

#### A. Abstraction of Storage Layer

The bot currently couples rate-limiting logic to in-memory dictionaries. To support optionality, we must introduce a `StorageProvider` interface.

- **Hurdle:** Avoiding code duplication between `LocalProvider` and `ValkeyProvider`.
- **Solution:** Use a Strategy Pattern where the `RateLimiter` and `SummaryCache` inject a provider based on the presence of `VALKEY_URL`.

#### B. The Fail-Soft Mechanism

**Fail-soft** is a design principle where a system maintains core functionality even when a secondary component (like a cache) is unavailable.

- **Hurdle:** Sync/Async connection timeouts in Python can hang the main event loop if not handled.
- **Solution:** Wrap Valkey calls in a `try-except` block with a strict 200ms timeout. If a `ConnectionError` occurs, the bot must log a `WARNING` and default to the `LocalProvider` logic for that specific request lifecycle.

#### C. Atomic Rate Limiting

- **Hurdle:** In a distributed environment, "read-then-write" operations on rate limits lead to race conditions.
- **Solution:** Utilize the Valkey `INCR` and `EXPIRE` commands in a single atomic pipeline or Lua script to ensure consistency across multiple bot instances.

---

### 4. Technical Mechanisms

#### The Caching Logic

When a video URL is processed, the bot generates a unique key: `summary:{sha256(video_url)}`.

- **Mechanism:** `SETEX key TTL value`
- **TTL Calculation:** Defaults to seconds (24 hours).
- **Logic:** If `GET key` returns data, the LLM call is bypassed entirely.

#### Rate Limit Windowing

The rate limit is governed by the formula:

- **Mechanism:** Upon message receipt, the bot executes `INCR user:{id}:limit`. If the result is , it immediately calls `EXPIRE user:{id}:limit RATE_LIMIT_WINDOW_SECONDS`.

---

### 5. Explicit Trade-offs

1. **Latency vs. Scalability:** We accept a increase in message processing time (network round-trip to Valkey) to gain the ability to run multiple bot instances.
2. **Memory Complexity:** We trade the simplicity of Python dictionaries for the operational overhead of managing a separate Valkey container.
3. **Consistency vs. Availability:** By choosing a "fail-soft" approach, we prioritize **Availability**. If Valkey goes down, the bot reverts to local memory; this may allow users to bypass global rate limits temporarily, but keeps the service running.

---

### 6. Configuration Schema

| Variable                    | Type    | Default | Description                                         |
| --------------------------- | ------- | ------- | --------------------------------------------------- |
| `VALKEY_URL`                | String  | `None`  | `valkey://[[username]:[password]]@host[:port][/db]` |
| `CACHE_TTL_SECONDS`         | Integer | `86400` | Expiration for cached video summaries.              |
| `RATE_LIMIT_WINDOW_SECONDS` | Integer | `60`    | Duration of the rate limit sliding window.          |
