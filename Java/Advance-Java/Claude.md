# Java Interview Topics

## Core Java

- OOP principles (SOLID, DRY, KISS)
- Java Memory Model (Heap, Stack, Metaspace, GC roots)
- Garbage Collection algorithms (G1, ZGC, CMS) & tuning
- String immutability, String Pool, `intern()`
- `equals()` / `hashCode()` contract
- Generics (wildcards, bounded types, type erasure)
- Reflection & Annotations
- `final`, `finally`, `finalize` differences
- Exception hierarchy (checked vs unchecked), custom exceptions
- Java 8–21 features (Streams, Optional, Records, Sealed classes, Pattern Matching)

## Collections & Data Structures

- `HashMap` internals (hashing, buckets, treeification, resize)
- `ConcurrentHashMap` vs `Collections.synchronizedMap`
- `LinkedHashMap`, `TreeMap`, `LinkedList` vs `ArrayList`
- `Comparable` vs `Comparator`
- fail-fast vs fail-safe iterators

## Multithreading & Concurrency

- Thread lifecycle, `Runnable` vs `Callable`
- `synchronized`, `volatile`, atomic variables
- `ReentrantLock`, `ReadWriteLock`
- `ThreadLocal`
- `ExecutorService`, `ThreadPoolExecutor`, thread pool sizing
- `CompletableFuture`
- Deadlock, livelock, starvation — detection & prevention
- happens-before guarantee
- `wait()` / `notify()` vs `Condition`

## JVM Internals

- Class loading mechanism (Bootstrap, Extension, Application)
- JIT compilation
- GC tuning flags (`-Xms`, `-Xmx`, `-XX:+UseG1GC`)
- Memory leak identification & heap dump analysis
- Stack overflow vs `OutOfMemoryError`

## Design Patterns

- Singleton (thread-safe variants)
- Factory, Abstract Factory, Builder
- Proxy, Decorator, Adapter
- Observer, Strategy, Command
- Template Method

## Spring & Spring Boot

- IoC container & Dependency Injection internals
- Bean lifecycle & scopes
- `@Transactional` — propagation, isolation levels, rollback rules
- AOP — proxy mechanism (JDK vs CGLIB)
- Spring Security (JWT, OAuth2, filter chain)
- Spring Boot auto-configuration internals
- `ApplicationContext` vs `BeanFactory`

## JPA / Hibernate

- N+1 problem & solutions (`JOIN FETCH`, `@BatchSize`, `EntityGraph`)
- First & second level cache
- `LAZY` vs `EAGER` loading
- Optimistic vs Pessimistic locking
- `@Transactional` with Hibernate session flush modes
- JPQL vs Criteria API vs Native Query

## Database & SQL

- Indexes (B-Tree, composite, covering index)
- Transaction isolation levels (dirty read, phantom read, etc.)
- `EXPLAIN` / `EXPLAIN ANALYZE`
- Connection pooling (HikariCP)
- Deadlocks in DB

## Microservices & Architecture

- Monolith vs Microservices trade-offs
- REST best practices, idempotency, versioning
- Circuit Breaker (Resilience4j)
- API Gateway, Service Discovery (Eureka)
- Event-driven architecture (Kafka basics)
- Saga pattern, Outbox pattern
- CAP theorem, eventual consistency

## Testing

- Unit vs Integration vs E2E testing
- JUnit 5, Mockito (`@Mock`, `@Spy`, `@Captor`)
- `@SpringBootTest` vs slice tests (`@WebMvcTest`, `@DataJpaTest`)
- TDD approach

## System Design

- Rate limiting, caching strategies (Redis)
- Load balancing
- Horizontal vs vertical scaling
- Database sharding & replication
- Designing systems like URL shortener, notification service, etc.

[InterviewBit Java question](https://www.interviewbit.com/java-interview-questions/)
