# 🔵 Performance & Scaling

> **Category:** Performance &nbsp;|&nbsp; **Tags:** `sharding` `replication` `connection pool`

---

## Table of Contents
1. [Vertical vs Horizontal Scaling](#vertical-vs-horizontal-scaling)
2. [Database Sharding](#database-sharding)
3. [Replication](#replication)
4. [Connection Pooling](#connection-pooling)
5. [Caching Strategies](#caching-strategies)
6. [Interview Questions](#interview-questions)

---

## Vertical vs Horizontal Scaling

### Vertical Scaling (Scale Up)
Add more resources (CPU, RAM, faster disk) to the **existing single server**.

```
Before:  [DB server — 8 CPU, 32 GB RAM]
After:   [DB server — 32 CPU, 128 GB RAM]
```

- **Pros:** Simple — no application changes needed.
- **Cons:** Hardware limits; single point of failure; expensive; downtime during upgrade.
- **When:** Works well up to a point for most apps.

### Horizontal Scaling (Scale Out)
Add more servers and distribute the load.

```
Before: [DB server 1]
After:  [DB server 1] + [DB server 2] + [DB server 3]
```

- **Pros:** Near-limitless scalability; high availability.
- **Cons:** Complex — application must handle multiple nodes, data distribution, consistency.
- **When:** When vertical scaling hits limits or for global distribution.

---

## Database Sharding

**Sharding** splits a large database into smaller, independent pieces called **shards**. Each shard is a separate database holding a **subset of the data**.

```
All data       →   Shard 1: users 1–10M
                   Shard 2: users 10M–20M
                   Shard 3: users 20M–30M
```

### Sharding Strategies

#### Range-Based Sharding
```
Shard 1: user_id  1 – 10,000,000
Shard 2: user_id  10,000,001 – 20,000,000
Shard 3: user_id  20,000,001+
```
- **Pro:** Simple, good for range queries.
- **Con:** Hotspot problem — newest users cluster on the last shard.

#### Hash-Based Sharding
```
shard = hash(user_id) % num_shards
user_id = 42 → hash(42) % 3 = 0 → Shard 1
user_id = 99 → hash(99) % 3 = 2 → Shard 3
```
- **Pro:** Uniform distribution, no hotspots.
- **Con:** Range queries need all shards; resharding is painful (consistent hashing helps).

#### Directory-Based Sharding
```
Lookup service: user_id → shard_id mapping
```
- **Pro:** Flexible — can reassign shards without rehashing.
- **Con:** Lookup service is a bottleneck and single point of failure.

### Sharding Challenges
- **Cross-shard queries** — JOINs across shards are expensive or impossible.
- **Cross-shard transactions** — require distributed transactions (2PC), which are complex.
- **Resharding** — adding/removing shards requires data migration.
- **Non-uniform shard keys** — celebrity problem (one user has 10M followers → hot shard).

---

## Replication

**Replication** copies data from one database server (primary/master) to one or more other servers (replicas/slaves).

### Primary-Replica (Master-Slave) Replication

```
Writes ──→ [Primary/Master] ──replication──→ [Replica 1]
                                         ──→ [Replica 2]
Reads  ──→ [Replica 1 or 2]
```

- **Primary** accepts all writes.
- **Replicas** receive changes asynchronously (or synchronously) and serve read traffic.

### Replication Types

| Type | Behavior | Consistency | Performance |
|------|---------|-------------|-------------|
| **Synchronous** | Primary waits for replica to confirm write | Strong | Slower writes |
| **Asynchronous** | Primary confirms to client before replica syncs | Eventual | Faster writes, risk of data loss |
| **Semi-synchronous** | At least one replica confirms | Between | Balance |

### Use Cases
- **Read replicas:** Scale reads — route `SELECT` queries to replicas.
- **High availability:** If primary fails, promote a replica (failover).
- **Geographic distribution:** Replica in each region — users read locally.
- **Backups:** Take backups from replica to avoid impacting primary.

### Replication Lag
The delay between a write on the primary and its appearance on the replica.
- **Problem:** A user writes data, then immediately reads from a replica that hasn't synced yet — gets stale data.
- **Solution:** Read-your-writes consistency: after a write, direct the same user's reads to the primary for a short window.

---

## Connection Pooling

Opening a new database connection for every query is **expensive** (TCP handshake, auth, session setup — typically 20-100ms).

A **connection pool** maintains a pool of pre-opened connections and **reuses** them for incoming requests.

```
App Servers                Connection Pool             Database
  Request 1 ──────────────→ [conn 1]  ──────────────→ DB
  Request 2 ──────────────→ [conn 2]  ──────────────→ DB
  Request 3 ──────────────→ [conn 3]  ──────────────→ DB
  Request 4 ──── waits ───→ (waiting for a free conn)
```

### Common connection poolers

| Tool | Language/Layer | Notes |
|------|---------------|-------|
| **PgBouncer** | PostgreSQL (server-side) | Lightweight, transaction-mode pooling |
| **pgpool-II** | PostgreSQL | Also handles load balancing, replication |
| **HikariCP** | Java | Very fast, default in Spring Boot |
| **c3p0, DBCP** | Java | Older alternatives to HikariCP |
| **RDS Proxy** | AWS | Managed proxy for RDS/Aurora |

### Key configuration parameters
```
max_pool_size     = 20     # Maximum connections to DB
min_pool_size     = 5      # Keep this many connections open at all times
connection_timeout = 3000  # ms to wait for a connection from the pool
idle_timeout      = 600000 # ms before idle connection is closed
```

### Too many connections
Each PostgreSQL connection consumes ~5-10 MB of RAM. With 200 app server instances each opening 10 connections = 2000 connections → DB runs out of memory.

**Solution:** Use a server-side pooler (PgBouncer) — consolidate thousands of app connections into a small pool against the DB.

---

## Caching Strategies

Caching stores frequently accessed data in a faster store (Redis, Memcached) to avoid repeated DB queries.

### Cache-Aside (Lazy Loading) — Most common
```
Read:
  1. Check cache
  2. Cache miss → query DB → store in cache → return
  3. Cache hit → return from cache

Write:
  1. Update DB
  2. Invalidate (delete) the cache key
     OR update the cache key
```

```python
def get_user(user_id):
    key = f"user:{user_id}"
    cached = redis.get(key)
    if cached:
        return json.loads(cached)        # cache hit

    user = db.query("SELECT * FROM users WHERE id = ?", user_id)
    redis.setex(key, 3600, json.dumps(user))  # cache for 1 hour
    return user
```

### Write-Through
Every write goes to both cache and DB simultaneously.
- **Pro:** Cache is always fresh.
- **Con:** Every write has double overhead; cache filled with rarely-read data.

### Write-Behind (Write-Back)
Write to cache immediately; flush to DB asynchronously.
- **Pro:** Very fast writes.
- **Con:** Risk of data loss if cache crashes before DB flush.

### Read-Through
Cache sits in front of DB — application always reads from cache; cache handles DB population.
- Used by libraries/frameworks that abstract caching.

### Cache Eviction Policies

| Policy | Behavior |
|--------|---------|
| **LRU** (Least Recently Used) | Evict the least recently accessed item |
| **LFU** (Least Frequently Used) | Evict the least frequently accessed item |
| **TTL** | Evict after a time-to-live expires |
| **FIFO** | Evict the oldest item first |

---

## Interview Questions

### Q1. What is the difference between sharding and replication?

> **Answer:**
> - **Replication:** The **same data** is copied to multiple servers. Used for read scaling, high availability, and failover.
> - **Sharding:** The **data is split** across multiple servers. Each shard holds a different subset of rows. Used for write scaling and handling data too large for one server.
>
> They're complementary — a production system often uses both: data is sharded for write scale, and each shard has replicas for read scale and redundancy.

---

### Q2. What are the challenges of database sharding?

> **Answer:**
> - **Cross-shard queries:** JOINs across shards require querying multiple shards and merging in application code — expensive.
> - **Cross-shard transactions:** Require distributed transactions (2PC) which are complex and slow.
> - **Hotspots:** Uneven distribution (e.g., most writes going to the latest time-range shard). Fix: hash-based sharding.
> - **Resharding:** Adding shards requires migrating/redistributing data. Consistent hashing minimizes data movement.
> - **Schema changes:** Must be applied across all shards.

---

### Q3. What is connection pooling? Why is it important?

> **Answer:**
> Creating a new DB connection per request is expensive (20-100ms). A connection pool maintains **pre-established connections** that are reused. The application borrows a connection from the pool, uses it, and returns it.
>
> **Why important:** Without pooling, high traffic creates thousands of connections, exhausting DB memory. With pooling (e.g., PgBouncer, HikariCP), hundreds of app instances share a small pool of actual DB connections — dramatically improving throughput and stability.

---

### Q4. Explain the cache-aside pattern. What is a cache stampede?

> **Answer:**
> **Cache-aside:** On read, check cache first. On miss, fetch from DB, populate cache, return result. On write, update DB and invalidate the cache key.
>
> **Cache stampede:** When a popular cache key expires, many simultaneous requests all see a miss and all hit the DB at once — causing a traffic spike.
>
> **Prevention:**
> - **Mutex/lock:** Only one request fetches from DB; others wait.
> - **Probabilistic early expiry (XFetch):** Randomly refresh slightly before expiry.
> - **Stale-while-revalidate:** Serve stale data while one background request refreshes it.

---

### Q5. What is replication lag and how do you handle it?

> **Answer:**
> Replication lag is the delay between a write on the primary and its appearance on replicas. With async replication, replicas can be seconds behind.
>
> **Problems:** User writes data, then reads from a replica and gets stale data (a freshly created user appears to not exist).
>
> **Solutions:**
> - **Read-your-own-writes:** After a write, direct that user's reads to the primary for a short window (or until lag is resolved).
> - **Monotonic reads:** Always route a user's reads to the same replica (using sticky sessions or user-based routing).
> - **Use synchronous replication** for critical data (at cost of write latency).

---

<div align="center">
  <sub>← Back to <a href="Topic.md">All Topics</a></sub>
</div>
