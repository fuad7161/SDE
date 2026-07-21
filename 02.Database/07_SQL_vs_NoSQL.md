# 🟢 SQL vs NoSQL

> **Category:** Queries &nbsp;|&nbsp; **Tags:** `MongoDB` `Redis` `Cassandra` `eventual consistency`

---

## Table of Contents
1. [SQL Databases](#sql-databases)
2. [NoSQL Databases](#nosql-databases)
3. [NoSQL Types](#nosql-types)
4. [SQL vs NoSQL Comparison](#sql-vs-nosql-comparison)
5. [Eventual Consistency](#eventual-consistency)
6. [Schema-on-Read vs Schema-on-Write](#schema-on-read-vs-schema-on-write)
7. [Interview Questions](#interview-questions)

---

## SQL Databases

**Relational databases** store data in structured tables with a predefined schema. They use SQL (Structured Query Language) and enforce ACID properties.

### Characteristics
- Structured data with predefined schema
- Strong relationships via foreign keys
- ACID transactions
- Vertical scaling (traditionally)
- Mature — decades of tooling, optimization

### Popular SQL databases

| Database | Best for |
|---------|---------|
| PostgreSQL | Complex queries, JSONB, full-featured |
| MySQL / MariaDB | Web apps, e-commerce, high read traffic |
| SQLite | Embedded, local storage, testing |
| SQL Server | Enterprise, Microsoft stack |
| Oracle | Enterprise, complex analytics |

---

## NoSQL Databases

**NoSQL** ("Not Only SQL") databases break from the relational model to provide **flexibility, scalability, and performance** for specific use cases.

### Why NoSQL emerged
- Web-scale data (billions of records)
- Unstructured or semi-structured data
- Need for horizontal scaling beyond single-node SQL
- High write throughput requirements
- Schema evolution without downtime

---

## NoSQL Types

### 1. Document Store
Stores **JSON/BSON documents** — schemaless, self-describing.

```json
// MongoDB document
{
  "_id": "user_001",
  "name": "Alice",
  "email": "alice@example.com",
  "address": {
    "city": "New York",
    "zip": "10001"
  },
  "orders": [
    {"id": "ord_1", "total": 99.99},
    {"id": "ord_2", "total": 49.99}
  ]
}
```

```javascript
// MongoDB query
db.users.find({ "address.city": "New York", "orders.total": { $gt: 50 } });
```

**Best for:** Content management, user profiles, catalogs, any data with variable structure.
**Examples:** MongoDB, CouchDB, Firestore, DynamoDB (can be used as document store).

---

### 2. Key-Value Store
Simple **key → value** lookup. Extremely fast O(1) operations.

```
SET user:session:abc123  {"user_id": 42, "expires": "2024-12-01"}  EX 3600
GET user:session:abc123  → {"user_id": 42, "expires": "2024-12-01"}
DEL user:session:abc123
```

**Best for:** Sessions, caching, rate limiting, leaderboards, pub/sub, queues.
**Examples:** Redis, Memcached, DynamoDB (key-value mode), etcd.

---

### 3. Wide-Column Store (Column-Family)
Data organized in **rows and columns**, but columns are grouped into **column families** and can vary per row. Designed for massive write throughput.

```
Row key: user_001
  Profile CF:  { name: "Alice", email: "alice@example.com" }
  Activity CF: { last_login: "2024-05-01", login_count: 42 }
  Prefs CF:    { theme: "dark", lang: "en" }
```

**CQL (Cassandra Query Language):**
```sql
CREATE TABLE user_activity (
    user_id  UUID,
    day      DATE,
    event    TEXT,
    payload  TEXT,
    PRIMARY KEY ((user_id), day, event)   -- partition key + clustering key
) WITH CLUSTERING ORDER BY (day DESC);

SELECT * FROM user_activity WHERE user_id = ? AND day > '2024-01-01';
```

**Best for:** Time-series data, IoT sensor data, activity feeds, write-heavy workloads at scale.
**Examples:** Apache Cassandra, Apache HBase, Google Bigtable, Amazon Keyspaces.

---

### 4. Graph Database
Data stored as **nodes** (entities) and **edges** (relationships), with properties on both.

```
(Alice)──[FRIENDS_WITH]──(Bob)
(Alice)──[LIKES]──(Post:123)
(Bob)──[WROTE]──(Post:123)
```

```cypher
// Cypher (Neo4j) — find friends of friends
MATCH (alice:User {name: "Alice"})-[:FRIENDS_WITH*2]-(fof)
WHERE fof.name <> "Alice"
RETURN DISTINCT fof.name;
```

**Best for:** Social networks, recommendation engines, fraud detection, knowledge graphs.
**Examples:** Neo4j, Amazon Neptune, ArangoDB.

---

## SQL vs NoSQL Comparison

| Feature | SQL | NoSQL |
|---------|-----|-------|
| Schema | Fixed, predefined | Flexible / schemaless |
| Data model | Tables + rows | Document/KV/Column/Graph |
| Relationships | FK + JOINs (native) | Application-level (mostly) |
| ACID | Full ACID | Varies (BASE typical) |
| Scaling | Vertical (primarily) | Horizontal (designed for it) |
| Query language | SQL (standardized) | API/query varies per DB |
| Consistency | Strong (default) | Eventual (often default) |
| Write throughput | Moderate | Very high (e.g., Cassandra) |
| Joins | First-class | Expensive/unsupported |
| Best for | Complex queries, financial data | Scale, unstructured, flexible |

### When to use SQL
- Financial transactions, banking
- Complex reporting with ad-hoc queries
- Strong consistency required
- Well-defined schema, normalized data
- You need full ACID and complex JOINs

### When to use NoSQL
- Massive scale (billions of records)
- Variable or evolving schema (user-generated content)
- High write throughput (IoT, logging, metrics)
- Simple access patterns (lookup by key)
- Graph/relationship-heavy data

---

## Eventual Consistency

**BASE** (alternative to ACID in distributed NoSQL):
- **B**asically Available — system always responds
- **S**oft state — state may change over time even without new input
- **E**ventually Consistent — all nodes will converge to same state *eventually*

```
Write to Cassandra Node A:  user.email = "new@example.com"
                                ↓ (async replication)
Read from Node B (immediately): user.email = "old@example.com"  ← stale!
Read from Node B (1 second later): user.email = "new@example.com" ← consistent
```

### Conflict resolution
When two nodes accept conflicting writes, they must be resolved:
- **Last Write Wins (LWW):** Highest timestamp wins.
- **Version vectors:** Track causality across nodes.
- **CRDTs:** Conflict-free Replicated Data Types — mathematically merge conflicts.
- **Application-level merge:** Application decides (e.g., shopping cart merge).

### Tunable consistency (Cassandra)
```
Write quorum:  W = majority of N replicas must confirm write
Read quorum:   R = majority of N replicas must confirm read
Strong consistency: W + R > N
```

---

## Schema-on-Read vs Schema-on-Write

| | Schema-on-Write | Schema-on-Read |
|--|----------------|----------------|
| **When schema applied** | At write time (enforced by DB) | At read time (by application) |
| **Example** | PostgreSQL: `ALTER TABLE` before inserting | MongoDB: insert any JSON; parse when reading |
| **Flexibility** | Low — schema must be updated before new fields | High — add fields anytime |
| **Consistency** | Strong — DB guarantees structure | Weak — application may see different structures |
| **Best for** | Known, stable domain model | Evolving data, multiple formats, prototyping |

---

## Interview Questions

### Q1. What is the difference between SQL and NoSQL databases?

> **Answer:**
> - **SQL:** Relational, fixed schema, ACID transactions, powerful JOINs, vertical scaling. Best for structured data with complex relationships.
> - **NoSQL:** Schema-flexible, horizontal scaling, BASE consistency, various data models (document, KV, column, graph). Best for massive scale, unstructured/semi-structured data, or specialized access patterns.
>
> Neither is better — choose based on data model, consistency requirements, and scale needs.

---

### Q2. What is eventual consistency? How does it differ from strong consistency?

> **Answer:**
> - **Strong consistency:** After a write, all subsequent reads return the new value. System pauses or errors if it can't guarantee this.
> - **Eventual consistency:** After a write, reads may temporarily return stale values, but all nodes will converge to the same value *eventually* (no timeline guarantee).
>
> Eventual consistency allows higher availability and write throughput at the cost of potentially stale reads. Used by Cassandra, DynamoDB (default), DNS.

---

### Q3. When would you choose MongoDB over PostgreSQL?

> **Answer:**
> Choose **MongoDB** when:
> - Schema is highly variable or evolving (product catalog with different attributes per category).
> - Data is naturally document-shaped and deeply nested (no JOINs needed).
> - Horizontal scaling is required from the start.
> - Rapid iteration/prototyping.
>
> Choose **PostgreSQL** when:
> - You need ACID transactions across multiple entities.
> - Complex reporting queries with JOINs.
> - Strong data integrity (FK constraints, NOT NULL, unique).
> - PostgreSQL also supports JSONB — bridging the gap for semi-structured data.

---

### Q4. What is Cassandra's data model? Why is it good for high write throughput?

> **Answer:**
> Cassandra uses a **wide-column store** with a distributed, masterless architecture. Data is partitioned by a **partition key** and sorted by a **clustering key** within each partition.
>
> High write throughput because:
> - Writes are appended to an **in-memory structure (MemTable)** and a sequential **commit log** — no random disk writes.
> - No single master — all nodes accept writes (writes go to `W` replicas in parallel).
> - No UPDATE in place — writes are always appends; **compaction** merges later.
> - Optimized for the write path at the expense of flexible querying (queries must match the primary key structure).

---

### Q5. What is Redis and what are its common use cases?

> **Answer:**
> Redis is an **in-memory key-value data structure store**. It supports strings, lists, sets, sorted sets, hashes, bitmaps, streams, and more.
>
> Common use cases:
> - **Caching:** Store DB query results, HTML fragments.
> - **Sessions:** Fast session storage and retrieval.
> - **Rate limiting:** Increment counters with TTL.
> - **Leaderboards:** Sorted sets for score-based ranking.
> - **Pub/Sub messaging:** Real-time event broadcasting.
> - **Distributed locks:** `SET key value NX PX timeout`.
> - **Job queues:** Lists as FIFO queues.

---

<div align="center">
  <sub>← Back to <a href="Topic.md">All Topics</a></sub>
</div>
