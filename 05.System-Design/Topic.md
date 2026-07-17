## Step 1 — Framework: How to Answer Any Design Question

1. Clarify Requirements       → Functional + Non-functional (scale, latency, availability)
2. Estimate Scale             → QPS, storage, bandwidth back-of-envelope
3. High-Level Design          → Draw the big picture (API → Service → DB)
4. Deep Dive Components       → Interviewer guides, you elaborate
5. Handle Bottlenecks         → Scaling, caching, sharding, replication
6. Discuss Tradeoffs          → Never say "this is the best" — always "it depends"

---

## Step 2 — Core Fundamentals

| Topic | Must Know |
|---|---|
| Scalability | Vertical vs Horizontal scaling |
| Load Balancing | Round robin, Least connections, Consistent hashing |
| Caching | Cache-aside, Write-through, Write-back, Eviction (LRU/LFU) |
| CAP Theorem | Can't have all 3 — tradeoffs explained |
| ACID vs BASE | When to choose which |
| SQL vs NoSQL | Decision criteria, not just definitions |
| Sharding | Range, Hash, Directory-based |
| Replication | Master-slave, Master-master, Quorum |

---

## Step 3 — Key Components

### 📦 Databases
- When to use PostgreSQL vs MongoDB vs Cassandra vs Redis
- Database indexing — B-Tree, composite index, covering index
- Read replicas and when they help
- Connection pooling — why it matters at scale

### 📨 Messaging
- Kafka vs RabbitMQ vs SQS — pull vs push, ordering, durability
- At-least-once vs exactly-once delivery
- Dead letter queues

### 🔗 APIs
- REST vs GraphQL vs gRPC — tradeoffs, when to pick each
- API Gateway — rate limiting, auth, routing
- Long polling vs WebSocket vs SSE

### ☁️ Infrastructure
- CDN — edge caching, cache invalidation
- DNS — how resolution works, GeoDNS
- Reverse Proxy — Nginx, load balancing
- Microservices vs Monolith — when to split

---

## Step 4 — Numbers Every Designer Should Know

| Metric | Value |
|---|---|
| L1 cache read | ~1 ns |
| Memory read | ~100 ns |
| SSD read | ~100 µs |
| Network round trip (same DC) | ~500 µs |
| HDD seek | ~10 ms |
| Packet: CA → Netherlands | ~150 ms |
| 1 Gbps network throughput | ~125 MB/s |

---

## Step 5 — Systems to Design (Practice)

### Tier 1 — Asked Almost Everywhere
- URL Shortener (Bitly) — hashing, redirection, analytics
- Rate Limiter — Token bucket, Leaky bucket, Sliding window
- Key-Value Store (Redis) — consistency, partitioning
- Chat System (WhatsApp) — WebSocket, message queue, presence
- News Feed / Timeline (Twitter/Facebook) — fan-out on write vs read
- Notification System — push/pull, delivery guarantees

### Tier 2 — Mid/Senior Level
- Search Autocomplete — Trie, prefix caching, ranking
- Distributed Cache (Redis Cluster) — consistent hashing, eviction
- File Storage (Google Drive / S3) — chunking, deduplication, metadata
- Video Streaming (YouTube/Netflix) — CDN, adaptive bitrate, encoding
- Web Crawler — BFS, politeness, URL frontier, deduplication
- Ride-sharing (Uber) — geospatial indexing, matching, surge pricing

### Tier 3 — Senior / Architect Level
- Distributed Message Queue (Kafka) — partitions, consumer groups, offsets
- Distributed ID Generator — Snowflake ID, UUID tradeoffs
- Distributed Locking — Redis SETNX, Redlock, Zookeeper
- Search Engine (Elasticsearch) — inverted index, sharding, relevance
- Payment System — idempotency, double-spend prevention, consistency
- Metrics & Monitoring (Prometheus/Grafana) — time-series, scraping

---

## Step 6 — Go-Specific Angles

| Topic | Go Relevance |
|---|---|
| HTTP Streaming / NDJSON | `http.Flusher`, large result sets |
| Worker Pool | Goroutines + buffered channels |
| Rate Limiter | `golang.org/x/time/rate` token bucket |
| Circuit Breaker | `sony/gobreaker` or custom |
| gRPC services | Protobuf, streaming RPCs |
| Context propagation | Cancellation, deadlines across services |
| Hexagonal Architecture | Ports & adapters in Go projects |

---

## Interview Questions

### Fundamentals
- "How does consistent hashing work and why is it used?"
- "Explain CAP theorem with a real example"
- "How would you design a system that needs high availability?"
- "What's the difference between latency and throughput?"
- "How does a CDN work?"

### Caching
- "When should you NOT use caching?"
- "How do you handle cache invalidation?"
- "What is a cache stampede and how do you prevent it?"
- "Explain the difference between write-through and write-back cache"

### Database
- "How do you handle database bottlenecks at scale?"
- "When would you shard a database? What are the downsides?"
- "How do read replicas affect consistency?"
- "Explain the N+1 query problem"

### Availability & Reliability
- "How do you achieve 99.99% uptime?"
- "What is a circuit breaker pattern?"
- "How do you handle partial failures in microservices?"
- "Explain idempotency and why it matters"

### Real Design Questions
- "Design a system that handles 1 million requests/second"
- "How would you design Twitter's trending topics?"
- "Design a distributed job scheduler"
- "How would you design Kafka from scratch?"
