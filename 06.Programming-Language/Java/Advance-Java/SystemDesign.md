# System Design — In-Depth Notes

---

## Table of Contents

1. [Rate Limiting](#1-rate-limiting)
2. [Caching Strategies (Redis)](#2-caching-strategies-redis)
3. [Load Balancing](#3-load-balancing)
4. [Horizontal vs Vertical Scaling](#4-horizontal-vs-vertical-scaling)
5. [Database Sharding & Replication](#5-database-sharding--replication)
6. [Designing a URL Shortener](#6-designing-a-url-shortener)
7. [Designing a Notification Service](#7-designing-a-notification-service)

---

## 1. Rate Limiting

Protects services from being overwhelmed by too many requests — prevents abuse, ensures fair usage, and controls costs.

### Token Bucket Algorithm

A bucket holds tokens up to a max capacity. Each request consumes 1 token. Tokens are refilled at a fixed rate.

```
Bucket capacity: 10 tokens
Refill rate:      2 tokens/second

Time 0s:  [■■■■■■■■■■]  10 tokens
Request:  [-1]           → allowed
Request:  [-1]           → allowed
Time 1s:  [+2 refill]   → [■■■■■■■■■■] (capped at 10)

Burst:    10 requests can be processed instantly (up to bucket size)
Steady:   2 requests/second sustained
```

**Pros**: Allows bursts; smooth average rate.  
**Cons**: Slightly complex to implement accurately.

---

### Sliding Window Log

Track timestamps of all requests in a window. If count exceeds limit, reject.

```
Limit: 5 requests per 60 seconds
Window: [now - 60s ... now]

Request at t=100s:
  Log: [50, 60, 70, 80, 90, 100]  ← 6 entries → reject ❌
  Log (after cleanup): [50, 60, 70, 80, 90] within window

Request at t=115s:
  Window = [55s, 115s]
  Log after cleanup: [60, 70, 80, 90, 100, 115] → 6 → reject ❌

  At t=111s entry at 50s drops out → [60,70,80,90,100] = 5 → allow ✅
```

**Pros**: Precise, no boundary burst issues.  
**Cons**: Memory-intensive (stores all timestamps).

---

### Sliding Window Counter (Redis-based)

Compromise — divide window into small buckets, sum recent buckets.

```
Window: 60s, divided into 6 buckets of 10s each

Current time: 105s
Buckets: {50-60: 2, 60-70: 3, 70-80: 1, 80-90: 4, 90-100: 2, 100-110: 1}
Sum = 13 requests in last 60s
```

**Redis implementation:**

```java
@Service
public class RateLimiter {

    @Autowired StringRedisTemplate redis;

    private static final long LIMIT = 100;
    private static final long WINDOW_SECONDS = 60;

    public boolean isAllowed(String userId) {
        String key = "rate:" + userId;
        long now = System.currentTimeMillis();
        long windowStart = now - (WINDOW_SECONDS * 1000);

        // Sliding window using sorted set (score = timestamp)
        redis.opsForZSet().removeRangeByScore(key, 0, windowStart);   // remove old
        Long count = redis.opsForZSet().zCard(key);                    // count current

        if (count == null || count < LIMIT) {
            redis.opsForZSet().add(key, UUID.randomUUID().toString(), now);  // add request
            redis.expire(key, Duration.ofSeconds(WINDOW_SECONDS + 10));
            return true;
        }
        return false;
    }
}
```

**HTTP response when rate-limited:**

```
HTTP/1.1 429 Too Many Requests
Retry-After: 30
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 0
X-RateLimit-Reset: 1700000030
```

---

## 2. Caching Strategies (Redis)

### Why Cache?

- Database reads are slow (disk I/O, network, query execution)
- Many reads are repeated (same user profile, product catalogue, config)
- Cache stores results in fast memory (Redis: sub-millisecond)

### Cache-Aside (Lazy Loading) — Most Common

Application controls the cache. Read from cache first; on miss, load from DB and populate cache.

```
Read:
  App → Cache? HIT  → return cached value ✅
  App → Cache? MISS → App → DB → populate cache → return value

Write:
  App → DB (update record)
  App → Cache.delete(key)   ← invalidate, not update
```

```java
@Service
public class ProductService {

    @Autowired ProductRepository repo;
    @Autowired RedisTemplate<String, Product> redis;

    private static final Duration TTL = Duration.ofMinutes(30);

    public Product getProduct(Long id) {
        String key = "product:" + id;
        Product cached = redis.opsForValue().get(key);
        if (cached != null) return cached;                    // cache HIT

        Product product = repo.findById(id)                   // cache MISS
            .orElseThrow(() -> new EntityNotFoundException("Product " + id));
        redis.opsForValue().set(key, product, TTL);           // populate cache
        return product;
    }

    public Product updateProduct(Long id, ProductRequest req) {
        Product product = repo.findById(id).orElseThrow();
        product.setPrice(req.getPrice());
        product = repo.save(product);
        redis.delete("product:" + id);                        // invalidate cache
        return product;
    }
}
```

**Pros**: Only caches what's requested (no wasted memory), resilient to cache failures.  
**Cons**: Cache miss causes 3 operations (check, DB, write); initial cold cache is slow.

---

### Write-Through

Write to cache AND DB together on every write. Cache is always up-to-date.

```
Write:
  App → Cache.set(key, value)  ← write to cache
  App → DB.update(...)         ← write to DB (same transaction or right after)

Read:
  App → Cache? always HIT (if data ever written)
```

```java
public Product updateProduct(Long id, ProductRequest req) {
    Product product = repo.findById(id).orElseThrow();
    product.setPrice(req.getPrice());
    product = repo.save(product);
    redis.opsForValue().set("product:" + id, product, TTL);  // write to cache too
    return product;
}
```

**Pros**: Cache always fresh, low read latency.  
**Cons**: Write latency higher, cache filled with data that may never be read.

---

### Write-Behind (Write-Back)

Write to cache immediately, write to DB **asynchronously** later.

```
Write:
  App → Cache.set(key, value)  ← instant
  Background worker ──────────► DB.update(...)  ← async, batched

Read:
  App → Cache? always HIT
```

**Pros**: Extremely fast writes, DB protected from write spikes, batching possible.  
**Cons**: Risk of data loss if cache crashes before DB flush; complex implementation.

---

### Cache Strategies Comparison

| Strategy | Read Path | Write Path | Consistency | Use Case |
|---|---|---|---|---|
| **Cache-Aside** | Cache → DB on miss | DB → invalidate cache | Eventual | Most read-heavy apps |
| **Write-Through** | Always cache | Cache + DB together | Strong | Frequent reads after writes |
| **Write-Behind** | Always cache | Cache → async DB | Eventual | High write throughput |
| **Read-Through** | Cache handles DB fetch | — | Eventual | Framework-managed (Hibernate L2C) |

---

### Spring Cache with Redis

```java
@SpringBootApplication
@EnableCaching
public class App { ... }
```

```yaml
spring:
  data:
    redis:
      host: localhost
      port: 6379
  cache:
    type: redis
    redis:
      time-to-live: 30m
```

```java
@Service
public class ProductService {

    @Cacheable(value = "products", key = "#id")
    public Product getProduct(Long id) {
        return repo.findById(id).orElseThrow();  // only called on cache miss
    }

    @CacheEvict(value = "products", key = "#id")
    public void deleteProduct(Long id) {
        repo.deleteById(id);
    }

    @CachePut(value = "products", key = "#result.id")  // always updates cache
    public Product updateProduct(Long id, ProductRequest req) {
        ...
    }
}
```

---

### Cache Invalidation Strategies

```
1. TTL (Time-To-Live)      → expire after N seconds/minutes
2. Event-driven invalidation → delete on write (cache-aside)
3. Write-through             → always keep in sync
4. Cache versioning          → change key prefix when data schema changes
                               e.g., "products:v2:123" vs "products:v1:123"
```

**Cache stampede** — many simultaneous cache misses hitting DB:

```java
// Solution: probabilistic early expiration or distributed lock on miss
public Product getProduct(Long id) {
    String key = "product:" + id;
    Product cached = redis.opsForValue().get(key);
    if (cached != null) return cached;

    // Acquire lock so only one thread fetches from DB
    String lockKey = "lock:product:" + id;
    Boolean locked = redis.opsForValue().setIfAbsent(lockKey, "1", Duration.ofSeconds(5));
    if (Boolean.TRUE.equals(locked)) {
        try {
            Product product = repo.findById(id).orElseThrow();
            redis.opsForValue().set(key, product, TTL);
            return product;
        } finally {
            redis.delete(lockKey);
        }
    } else {
        Thread.sleep(50);     // wait and retry
        return getProduct(id);
    }
}
```

---

## 3. Load Balancing

Distributes incoming requests across multiple server instances to prevent any single instance from being overwhelmed.

```
Clients
  │
  ├──► Load Balancer
  │          │
  │    ┌─────┴──────┐
  │    ▼    ▼    ▼  ▼
  │   S1   S2   S3  S4   (server instances)
  │    │    │    │   │
  │    └────┴────┴───┘
  │          │
  │        Database
```

### Load Balancing Algorithms

#### Round Robin
Each request goes to the next server in sequence: S1, S2, S3, S1, S2, S3...

```
Request 1 → S1
Request 2 → S2
Request 3 → S3
Request 4 → S1 (wraps around)
```

**Pros**: Simple, fair distribution.  
**Cons**: Ignores server capacity or current load.

---

#### Weighted Round Robin
Servers with higher capacity get more requests.

```
S1: weight=3  (powerful server)
S2: weight=1  (weaker server)

Requests: S1, S1, S1, S2, S1, S1, S1, S2...
```

---

#### Least Connections
Route to server with fewest active connections.

```
S1: 10 active connections
S2: 3  active connections  ← next request goes here
S3: 7  active connections
```

**Best for**: Long-lived connections (WebSockets, file uploads) where sessions vary in duration.

---

#### IP Hash / Consistent Hashing
Route based on client IP — same client always hits same server (**sticky sessions**).

```
hash(clientIP) % numServers = serverIndex

Client 192.168.1.1 → always S2
Client 10.0.0.5    → always S1
```

**Problem with simple modulo**: Adding/removing servers remaps all clients.  
**Consistent Hashing** solves this — only ~1/N clients are remapped when a server is added/removed.

```
Ring-based consistent hashing:
  Servers placed at positions on a ring (hash of server name)
  Request routed to first server clockwise from hash(key)

  Adding a new server? Only the requests between new server
  and its predecessor shift — not all requests.
```

---

### Layer 4 vs Layer 7 Load Balancing

| | L4 (Transport) | L7 (Application) |
|---|---|---|
| Operates at | TCP/UDP level | HTTP/HTTPS level |
| Routing based on | IP + port | URL, headers, cookies, body |
| SSL termination | ❌ | ✅ |
| Content-based routing | ❌ | ✅ |
| Speed | Faster | Slightly slower (parses HTTP) |
| Examples | AWS NLB, HAProxy (L4 mode) | Nginx, AWS ALB, Spring Cloud Gateway |

---

## 4. Horizontal vs Vertical Scaling

### Vertical Scaling (Scale Up)

Add more resources (CPU, RAM, disk) to the **same machine**.

```
Before:        After:
┌─────────┐    ┌─────────────┐
│ 4 cores │ →  │  16 cores   │
│  8 GB   │    │   64 GB RAM │
│ 100 GB  │    │   1 TB disk │
└─────────┘    └─────────────┘
```

**Pros**: Simple — no code changes, no distributed complexity, ACID transactions easy.  
**Cons**: Hardware limits (can't scale infinitely), single point of failure, downtime to upgrade, expensive per unit.

---

### Horizontal Scaling (Scale Out)

Add **more machines** of the same size.

```
Before:          After:
┌─────────┐      ┌─────────┐  ┌─────────┐  ┌─────────┐
│ Server1 │  →   │ Server1 │  │ Server2 │  │ Server3 │
└─────────┘      └─────────┘  └─────────┘  └─────────┘
                      └──────── Load Balancer ──────────┘
```

**Pros**: Near-infinite scaling, no single point of failure, cheaper commodity hardware, zero-downtime deployments.  
**Cons**: Requires stateless services, session management complexity, distributed systems challenges (consistency, coordination).

---

### Stateless Design for Horizontal Scaling

**Session state must NOT be stored in memory** — if the load balancer routes the next request to a different server, the session is lost.

```java
// ❌ Stateful — breaks with multiple instances
@RestController
public class CartController {
    private Map<String, Cart> sessions = new HashMap<>();  // in-memory state

    @PostMapping("/cart/add")
    public void addToCart(HttpSession session, @RequestBody Item item) {
        sessions.get(session.getId()).add(item);   // only works if same server
    }
}

// ✅ Stateless — uses Redis for session storage
@SpringBootApplication
// Spring Session automatically stores sessions in Redis
public class App { ... }
```

```yaml
spring:
  session:
    store-type: redis    # sessions stored in Redis, shared across all instances
```

**JWT** is another stateless approach — auth state is in the token itself, no server-side session needed.

---

### When to Use Each

| Factor | Vertical | Horizontal |
|---|---|---|
| Application type | Stateful (legacy) | Stateless microservices |
| Traffic pattern | Predictable, moderate | Unpredictable, variable |
| Availability requirement | Low (single point of failure OK) | High (HA needed) |
| Cost | High at scale | Better cost efficiency |
| Database | Works with single DB | Needs distributed DB or read replicas |

---

## 5. Database Sharding & Replication

### Replication

**Copy data** to multiple database nodes — one primary (writes), multiple replicas (reads).

```
Primary (read + write)
    │
    ├── async replication ──► Replica 1 (read only)
    ├── async replication ──► Replica 2 (read only)
    └── async replication ──► Replica 3 (read only)
```

**Benefits:**
- **Read scalability**: distribute read queries across replicas
- **High availability**: replica can be promoted to primary on failure
- **Backup**: replicas can serve as backups

```java
// Spring Boot routing reads to replicas, writes to primary
@Configuration
public class DataSourceConfig {

    @Bean
    @Primary
    public DataSource routingDataSource(
            @Qualifier("primary") DataSource primary,
            @Qualifier("replica") DataSource replica) {

        Map<Object, Object> sources = new HashMap<>();
        sources.put("primary", primary);
        sources.put("replica", replica);

        AbstractRoutingDataSource routing = new AbstractRoutingDataSource() {
            @Override
            protected Object determineCurrentLookupKey() {
                // Route read-only transactions to replica
                return TransactionSynchronizationManager.isCurrentTransactionReadOnly()
                    ? "replica" : "primary";
            }
        };
        routing.setTargetDataSources(sources);
        routing.setDefaultTargetDataSource(primary);
        return routing;
    }
}
```

---

### Sharding

Splitting data **horizontally** across multiple database instances — each shard holds a subset of rows.

```
Without sharding:          With sharding (by user_id):
┌────────────┐             Shard 0: user_id % 3 == 0  → [0, 3, 6, 9...]
│ All users  │             ┌────────────┐
│ (100M rows)│      →      │  33M users │  DB 0
└────────────┘             └────────────┘
                           ┌────────────┐
                           │  33M users │  DB 1  (user_id % 3 == 1)
                           └────────────┘
                           ┌────────────┐
                           │  33M users │  DB 2  (user_id % 3 == 2)
                           └────────────┘
```

---

### Sharding Strategies

#### Range-Based Sharding
Shard by value ranges.

```
user_id 1-10M    → Shard A
user_id 10M-20M  → Shard B
user_id 20M+     → Shard C
```

**Pros**: Easy range queries, intuitive.  
**Cons**: Hot spots (new users all go to latest shard), uneven distribution.

#### Hash-Based Sharding
Shard by hash of key.

```
shardId = hash(userId) % numShards

hash("user:1001") % 3 = 2 → Shard 2
hash("user:1002") % 3 = 0 → Shard 0
```

**Pros**: Even distribution, no hot spots.  
**Cons**: Range queries require hitting all shards; resharding is disruptive (use consistent hashing).

#### Directory-Based Sharding
Maintain a lookup table mapping each key to its shard.

```
Lookup table:
  user 1001 → Shard B
  user 1002 → Shard A
  user 1003 → Shard C
  ...
```

**Pros**: Flexible, easy to rebalance.  
**Cons**: Lookup table is a bottleneck and single point of failure.

---

### Sharding Challenges

```
Cross-shard joins:
  SELECT u.name, o.total
  FROM users u JOIN orders o ON u.id = o.user_id
  WHERE u.region = 'EU'        ← users and orders may be on different shards

Solution: Denormalize data (store user info in orders table)
          or use application-side join (fetch from multiple shards)

Cross-shard transactions:
  Transfer from user on Shard A to user on Shard B
  → No ACID across shards → use Saga pattern

Resharding:
  Adding a new shard → need to migrate data
  → Use consistent hashing to minimize data movement
```

---

## 6. Designing a URL Shortener

### Requirements

- Given a long URL, return a short URL (e.g., `short.ly/abc123`)
- Redirect short URL to original URL
- Support ~100M URLs, ~1000 writes/sec, ~100K reads/sec

### System Components

```
Client
  │
  ├── POST /shorten {url: "https://..."}
  │         │
  │    API Gateway / Load Balancer
  │         │
  │    URL Service
  │         ├── Generate short code
  │         ├── Store in DB (shortCode → longUrl)
  │         └── Cache in Redis (shortCode → longUrl)
  │
  └── GET /abc123
            │
       URL Service (redirect)
            ├── Check Redis cache → HIT → 302 redirect
            └── Check DB → cache populate → 302 redirect
```

### Short Code Generation

#### Base62 Encoding

Encode a unique ID (from DB auto-increment or distributed ID generator) in base62.

```
Characters: 0-9 (10) + a-z (26) + A-Z (26) = 62

ID: 1000000
Base62: 4c92

6-character base62 code → 62^6 = ~56 billion unique URLs
```

```java
private static final String CHARS = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ";

public String encode(long id) {
    StringBuilder sb = new StringBuilder();
    while (id > 0) {
        sb.append(CHARS.charAt((int)(id % 62)));
        id /= 62;
    }
    return sb.reverse().toString();
}

public long decode(String code) {
    long id = 0;
    for (char c : code.toCharArray()) {
        id = id * 62 + CHARS.indexOf(c);
    }
    return id;
}
```

#### Distributed ID Generation

For multiple servers generating codes simultaneously, use a **distributed ID generator** (Snowflake ID):

```
64-bit Snowflake ID:
  [1 bit: 0] [41 bits: timestamp ms] [10 bits: machine id] [12 bits: sequence]
  → guarantees unique IDs across machines without coordination
```

---

### Database Schema

```sql
CREATE TABLE urls (
    id         BIGSERIAL PRIMARY KEY,
    short_code VARCHAR(10) UNIQUE NOT NULL,
    long_url   TEXT NOT NULL,
    user_id    BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP,
    click_count BIGINT DEFAULT 0
);

CREATE INDEX idx_short_code ON urls(short_code);  -- fast lookup
```

---

### Redirect Flow

```java
@GetMapping("/{code}")
public ResponseEntity<Void> redirect(@PathVariable String code) {
    String longUrl = cache.get(code);               // Redis first
    if (longUrl == null) {
        longUrl = urlRepository.findByShortCode(code)
            .orElseThrow(() -> new NotFoundException(code))
            .getLongUrl();
        cache.set(code, longUrl, Duration.ofHours(24));
    }

    // Async click tracking — don't slow down redirect
    eventPublisher.publish(new ClickEvent(code));

    return ResponseEntity.status(HttpStatus.FOUND)  // 302 redirect
        .location(URI.create(longUrl))
        .build();
}
```

**301 vs 302:**
- `301 Permanent` — browser caches redirect → less server traffic but can't update/track
- `302 Temporary` — browser always asks server → enables tracking and future changes (preferred)

---

## 7. Designing a Notification Service

### Requirements

- Send notifications via Email, SMS, Push
- High throughput: ~1M notifications/day
- Reliable delivery (at-least-once), idempotent consumers
- User preferences (opt-out channels)

### Architecture

```
Order Service, Payment Service, etc.
        │ (publish events)
        ▼
    [Kafka Topic: notifications.requested]
        │
    Notification Service (consumer)
        ├── Check user preferences → skip if opted out
        ├── Deduplicate (idempotency key in Redis)
        ├── Enqueue to channel-specific queue
        │
        ├── [email.queue]   → Email Worker → SendGrid / SES
        ├── [sms.queue]     → SMS Worker   → Twilio
        └── [push.queue]    → Push Worker  → FCM / APNs
```

---

### Fan-Out Pattern

One event triggers notifications to many users (e.g., system announcement to 10M users).

```
Approach 1 — Fan-out on write (push model):
  Event → Notification Service → generate 10M messages immediately
  ✅ Fast reads (pre-generated)
  ❌ Huge write amplification for large audiences

Approach 2 — Fan-out on read (pull model):
  Event → store announcement once
  When user opens app → fetch unread notifications
  ✅ No write amplification
  ❌ Slower reads (query on demand)

Approach 3 — Hybrid:
  Small audiences (< 1000) → fan-out on write
  Large audiences → fan-out on read
```

---

### Key Design Considerations

```
Reliability (at-least-once delivery):
  - Kafka consumer with manual commit
  - Only commit offset AFTER successful send
  - Retry with exponential backoff for transient failures
  - Dead-letter queue (DLQ) for permanent failures

Deduplication (avoid double notification):
  - Store idempotency key (eventId + channel + userId) in Redis
  - Check before sending: if key exists → skip
  - Set TTL matching delivery window (e.g., 24h)

User Preferences:
  - Table: user_notification_prefs(userId, channel, enabled, frequency)
  - Check before enqueuing to channel queue

Rate limiting per user:
  - Max N notifications per channel per hour
  - Sliding window counter per user per channel in Redis

Template management:
  - Store templates in DB: {id, channel, event_type, subject, body_template}
  - Use Mustache/Thymeleaf for variable substitution
  - Example: "Hello {{name}}, your order {{orderId}} has been confirmed."
```

```java
@Component
public class NotificationConsumer {

    @KafkaListener(topics = "notifications.requested", groupId = "notification-service")
    @Transactional
    public void handleNotification(NotificationRequest request,
                                   Acknowledgment ack) {
        String dedupKey = request.getEventId() + ":" + request.getChannel() + ":" + request.getUserId();

        // Deduplication check
        if (redis.hasKey(dedupKey)) {
            ack.acknowledge();   // already processed, skip
            return;
        }

        // User preference check
        if (!prefService.isEnabled(request.getUserId(), request.getChannel())) {
            ack.acknowledge();
            return;
        }

        // Rate limit check
        if (!rateLimiter.isAllowed("notif:" + request.getUserId() + ":" + request.getChannel())) {
            // Requeue with delay or drop
            return;   // don't ack → Kafka retries later
        }

        // Send notification
        channelRouter.send(request);

        // Mark as processed
        redis.set(dedupKey, "1", Duration.ofHours(24));

        ack.acknowledge();   // commit offset only after success
    }
}
```
