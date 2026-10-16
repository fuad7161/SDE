# Backend Engineering Interview Questions — Categorized

## Table of Contents

- [1. Idempotency & Duplicate Processing](#1-idempotency--duplicate-processing)
- [2. Transactions & Spring `@Transactional`](#2-transactions--spring-transactional)
- [3. Database Indexing & Query Optimization](#3-database-indexing--query-optimization)
- [4. Concurrency Control & Lost Updates](#4-concurrency-control--lost-updates)
- [5. Caching & Cache Stampede](#5-caching--cache-stampede)
- [6. Rate Limiting](#6-rate-limiting)
- [7. API Design (Errors, Pagination, Consistency)](#7-api-design-errors-pagination-consistency)
- [8. Microservices & Latency Optimization](#8-microservices--latency-optimization)
- [9. Kafka & Event Streaming](#9-kafka--event-streaming)
- [10. Design Patterns & OOP / Spring DI](#10-design-patterns--oop--spring-di)
- [11. System Design (Misc.)](#11-system-design-misc)
- [12. Data Structures & Algorithms](#12-data-structures--algorithms)

---

## 1. Idempotency & Duplicate Processing

- Kafka consumer processes a payment event, updates DB, crashes before committing offset — how to prevent processing twice?
- Go/Java backend receives a create-order request; client may retry on slow network — how to prevent duplicate orders?
- Designing `POST /orders` to be idempotent when clients retry on timeout.
- Kafka consumer processes `OrderPlaced` events, crashes after DB commit but before offset ack — how to make processing safe?
- Kafka consumer processes `PaymentCompleted` events, crashes after DB commit but before offset ack/commit.

## 2. Transactions & Spring `@Transactional`

- `@Transactional` public method called via `this.saveOrder()` from same class — why doesn't transaction start, and how to fix?
- `OrderService.processOrders()` calls `this.updateOrder()` (annotated `@Transactional`) — does it work? Why not? Fix?
- Spring service saves an order then calls external payment API — should DB roll back automatically on failure? How to design safely?
- `transfer()` method is `@Transactional` and calls an external HTTP notification inside — problems and redesign.

## 3. Database Indexing & Query Optimization

- `orders` query: `WHERE customer_id = ? AND status = 'PAID' ORDER BY created_at DESC LIMIT 20` — which index and why?
- Same query shape (`customer_id`, `status`, `created_at DESC`, `LIMIT 20`) — index choice and reasoning.
- `orders(id, customer_id, status, created_at, total_amount)` query with `status = 'COMPLETED'` + `ORDER BY created_at DESC LIMIT 20` — index + considerations for 100M+ rows.

## 4. Concurrency Control & Lost Updates

- Two users fetch same customer record; User B overwrites A's changes with stale data — how to prevent lost update via REST API?
- Two users update the same order nearly simultaneously — prevent silent overwrite in Spring Boot.
- Concurrent updates to a user's account balance — prevent race condition where both read old balance and overwrite.

## 5. Caching & Cache Stampede

- `GET /products/{id}` cached in Redis; popular product expires → many simultaneous DB hits — prevent cache stampede.
- Redis-cached product details with 10-min TTL; popular item expires, thousands of concurrent misses hit DB — prevent stampede + trade-offs.

## 6. Rate Limiting

- Login API rate limiter: max 5 attempts per user per 10 minutes — implement with Redis + behavior on limit exceeded.
- Rate limiter for `POST /login`: max 5 requests per user per minute — multi-instance backend using Redis.
- REST API: max 100 requests per user per minute across 5 instances behind load balancer — algorithm, Redis usage, simultaneous requests.

## 7. API Design (Errors, Pagination, Consistency)

- Spring Boot REST API returns inconsistent errors (plain text, stack traces, varied validation formats) — design consistent error response.
- `GET /orders` on large table — pagination design that stays fast and consistent during inserts.
- `GET /api/transactions?page=10000&size=50` slow on millions of rows with continuous inserts — redesign pagination to avoid skipped/duplicated records.

## 8. Microservices & Latency Optimization

- Spring Boot endpoint calls 3 independent external services sequentially — reduce latency, handle partial failures.

## 9. Kafka & Event Streaming

- Kafka topic with millions of user accounts — per-user ordering, cross-user parallelism; choose partition key and consumer design.

## 10. Design Patterns & OOP / Spring DI

- Notification service with several `if/else` for email — add SMS, push, future channels without changing main service.
- `PaymentService` directly instantiates `StripeClient`, `PayPalClient`, `BankTransferClient` — redesign so new methods don't modify core logic.
- Spring Boot: three `PaymentGateway` implementations (Bkash, Card, Bank) — select by `PaymentMethod` without large `if/else`/`switch` using DI.
- `Bird`/`Sparrow`/`Penguin` where `Penguin.fly()` throws `UnsupportedOperationException` — why is this an OOP problem, and how to redesign?

## 11. System Design (Misc.)

- URL-shortening service — DB schema and constraints to guarantee unique short codes + fast lookups.

## 12. Data Structures & Algorithms

- Design an LRU cache with `get(key)` / `put(key, value)` in **O(1)** average time; evict least recently used on overflow.
- Longest contiguous subarray with at most **two distinct values**.
- Millions of API request logs with endpoint paths — find **top K most frequent endpoints**; approach, complexity, tie handling.
