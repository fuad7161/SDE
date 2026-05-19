# 🔵 CDN & Caching

> **Category:** Infrastructure &nbsp;|&nbsp; **Tags:** `edge servers` `cache-control` `invalidation`

---

## Table of Contents
1. [What is a CDN?](#what-is-a-cdn)
2. [How a CDN Works](#how-a-cdn-works)
3. [Cache-Control Headers](#cache-control-headers)
4. [Cache Invalidation](#cache-invalidation)
5. [Caching Layers](#caching-layers)
6. [Cache Strategies](#cache-strategies)
7. [Interview Questions](#interview-questions)

---

## What is a CDN?

A **Content Delivery Network (CDN)** is a geographically distributed network of **edge servers** (also called Points of Presence — PoPs) that cache and serve content from locations close to the user.

**Goal:** Reduce latency by serving content from a server that is geographically close to the client, rather than the origin server.

```
Without CDN:
User (Tokyo) ──────── 150ms ────────→ Origin (New York)

With CDN:
User (Tokyo) ── 5ms ──→ Edge (Tokyo) ── (cache hit) ──→ Response
```

**Popular CDNs:** Cloudflare, AWS CloudFront, Akamai, Fastly, Google Cloud CDN.

---

## How a CDN Works

```
1. Client requests https://example.com/logo.png
2. DNS resolves to the nearest CDN edge server (via anycast or geo-DNS)
3. Edge server checks its cache:
   → Cache HIT:  serve cached content immediately
   → Cache MISS: fetch from origin, cache it, then serve
4. Subsequent requests served from edge cache
```

### Cache Hit Ratio
`Cache Hit Ratio = (Cache Hits) / (Total Requests) × 100%`

A higher ratio = less load on origin = better performance.

---

## Cache-Control Headers

HTTP cache behavior is controlled via response headers.

### `Cache-Control` directives

| Directive | Meaning |
|-----------|---------|
| `max-age=3600` | Cache for 3600 seconds (1 hour) |
| `s-maxage=3600` | Shared cache (CDN) max-age (overrides `max-age` for CDN) |
| `no-store` | Never cache (sensitive data) |
| `no-cache` | Cache but must revalidate with server before use |
| `private` | Browser can cache, CDN must not |
| `public` | Any cache (browser + CDN) can cache |
| `immutable` | Content will never change — don't revalidate |
| `stale-while-revalidate=60` | Serve stale while fetching fresh in background |
| `must-revalidate` | Once stale, must revalidate before serving |

### Example for a static asset (image)

```
Cache-Control: public, max-age=31536000, immutable
```
→ Cache for 1 year, never revalidate (fingerprint the filename instead).

### Example for an API response

```
Cache-Control: private, no-cache
```
→ Browser can cache but must check with server; CDN must not cache.

---

### ETag and Last-Modified (Conditional Requests)

When a cached resource expires, the browser sends a **conditional request** to check if the content changed:

```
# Browser sends:
If-None-Match: "abc123"          (ETag from previous response)
If-Modified-Since: Mon, 1 Jan 2024 00:00:00 GMT

# Server responds:
304 Not Modified  (if content unchanged — browser uses cached copy)
200 OK + new body (if content changed)
```

---

## Cache Invalidation

> "There are only two hard things in CS: cache invalidation and naming things." — Phil Karlton

### Strategies

| Method | How | When to use |
|--------|-----|------------|
| **TTL expiry** | Cache expires naturally after `max-age` | Acceptable staleness |
| **URL fingerprinting** | Include content hash in filename: `main.abc123.js` | Static assets |
| **Purge API** | CDN API call to delete specific cached objects | Urgent content updates |
| **Cache-busting** | Append query string: `logo.png?v=2` | Simple versioning |
| **Surrogate-Key/Tag** | Tag related objects, purge by tag | Structured content invalidation |

### Best practice for static assets:
- Set `max-age=31536000, immutable` (1 year).
- Change the filename when content changes (webpack fingerprinting).
- No need for manual invalidation — the URL changes automatically.

---

## Caching Layers

| Layer | Location | What it caches |
|-------|----------|---------------|
| **Browser cache** | Client | HTTP responses per origin |
| **Service Worker** | Client (JS) | Programmable cache, offline support |
| **CDN / Edge cache** | PoP servers | Static assets, API responses |
| **Reverse proxy cache** | Nginx/Varnish in front of app | HTML, API responses |
| **Application cache** | In-process (in-memory map) | DB query results |
| **Distributed cache** | Redis, Memcached | Shared session data, DB results |
| **Database cache** | Buffer pool, query cache | Row and page-level caching |

---

## Cache Strategies

### Cache-Aside (Lazy Loading)
```
Read:  check cache → miss → query DB → populate cache → return
Write: update DB → invalidate or update cache
```
Most common. Cache only what's actually read.

### Write-Through
```
Write: update cache AND DB synchronously
Read:  always hits cache
```
Consistent, but writes are slower.

### Write-Behind (Write-Back)
```
Write: update cache immediately, async write to DB later
Read:  fast
```
Risk: data loss if cache crashes before DB write.

### Read-Through
Cache sits in front of DB — application always reads from cache, which populates itself on miss.

---

## Interview Questions

### Q1. What is a CDN and why would you use one?

> **Answer:**  
> A CDN is a network of edge servers distributed globally that cache and serve content close to users. Benefits:
> - **Reduced latency:** Users get content from a nearby edge server instead of a distant origin.
> - **Reduced origin load:** Most requests served from cache, saving bandwidth and compute.
> - **Better availability:** CDN absorbs DDoS traffic, provides redundancy.
> - **Improved throughput:** Edge servers are optimized for high-volume static content delivery.

---

### Q2. What is the difference between `Cache-Control: no-cache` and `no-store`?

> **Answer:**
> - **`no-cache`:** The resource **can be cached**, but the cache must **revalidate** with the server (using ETag or Last-Modified) before serving it. If unchanged (304), the cached copy is used.
> - **`no-store`:** The resource **must never be cached** at all — every request goes to the origin. Used for sensitive data (banking, private user data).

---

### Q3. How would you handle cache invalidation for a CDN?

> **Answer:**  
> **For static assets (JS, CSS, images):** Use **content-addressed filenames** (hash in the filename, e.g., `app.a1b2c3.js`). Set `Cache-Control: immutable, max-age=31536000`. When the content changes, the filename changes, so no invalidation needed.
>
> **For dynamic content:** Use the CDN's **purge API** to explicitly remove specific URLs or cache tags. Alternatively, use short TTLs and `stale-while-revalidate` for freshness with performance.

---

### Q4. What is cache stampede (thundering herd) and how do you prevent it?

> **Answer:**  
> When a cached item expires, multiple requests arrive simultaneously — all see a cache miss and all go to the origin/DB at once, causing a spike.
>
> **Prevention strategies:**
> - **Mutex/lock:** Only one request populates the cache; others wait.
> - **Probabilistic early expiration:** Before the TTL expires, randomly start refreshing (XFetch algorithm).
> - **Stale-while-revalidate:** Serve the stale value while refreshing in the background.
> - **Pre-warming:** Proactively refresh cache before expiry for hot keys.

---

### Q5. What is the difference between an edge cache and a CDN origin?

> **Answer:**
> - **Origin:** The authoritative source of content — your web server, S3 bucket, or API. It holds the source of truth.
> - **Edge cache (CDN PoP):** A geographically distributed cache that stores copies of origin content close to users. It serves cached content directly, only going to origin on a cache miss or expired TTL.

---

<div align="center">
  <sub>← Back to <a href="networking_interview_topics.md">All Topics</a></sub>
</div>
