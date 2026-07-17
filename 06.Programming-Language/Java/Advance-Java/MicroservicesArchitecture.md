# Microservices & Architecture — In-Depth Notes

---

## Table of Contents

1. [Monolith vs Microservices Trade-offs](#1-monolith-vs-microservices-trade-offs)
2. [REST Best Practices, Idempotency, Versioning](#2-rest-best-practices-idempotency-versioning)
3. [Circuit Breaker (Resilience4j)](#3-circuit-breaker-resilience4j)
4. [API Gateway, Service Discovery (Eureka)](#4-api-gateway-service-discovery-eureka)
5. [Event-Driven Architecture (Kafka Basics)](#5-event-driven-architecture-kafka-basics)
6. [Saga Pattern, Outbox Pattern](#6-saga-pattern-outbox-pattern)
7. [CAP Theorem, Eventual Consistency](#7-cap-theorem-eventual-consistency)

---

## 1. Monolith vs Microservices Trade-offs

### Monolith

All functionality in a **single deployable unit** — one codebase, one database, one process.

```
┌──────────────────────────────────────────┐
│           Monolithic Application          │
│  ┌────────┐  ┌──────────┐  ┌──────────┐ │
│  │ Orders │  │ Payments │  │  Users   │ │
│  └────────┘  └──────────┘  └──────────┘ │
│  ┌────────┐  ┌──────────┐  ┌──────────┐ │
│  │ Email  │  │ Reports  │  │Inventory │ │
│  └────────┘  └──────────┘  └──────────┘ │
│               Single DB                  │
└──────────────────────────────────────────┘
```

**Pros**: Simple to develop/test/deploy initially, no network overhead, easy debugging, ACID transactions across modules.  
**Cons**: Scales as a whole unit, deployment risk (everything redeploys), tech stack lock-in, hard to scale teams.

---

### Microservices

Functionality split into **small independent services** — each owns its data, deploys independently.

```
Client → API Gateway
              │
    ┌─────────┼─────────┐
    ▼         ▼         ▼
 Order      Payment    User
 Service    Service    Service
    │          │          │
  OrderDB  PaymentDB   UserDB

Each service:
  - Has its own database
  - Communicates via REST/gRPC/events
  - Deploys independently
  - Can scale independently
```

**Pros**: Independent deployment, polyglot (different stacks per service), targeted scaling, fault isolation.  
**Cons**: Network latency, distributed transactions, complex ops (k8s, service mesh), harder debugging, data consistency challenges.

### When to Use What

| Factor | Monolith | Microservices |
|---|---|---|
| Team size | Small (< 10 devs) | Large, multiple teams |
| Stage | Early startup / MVP | Scaling product |
| Deployment frequency | Low | High (CI/CD per service) |
| Scalability needs | Uniform | Non-uniform (some services need more) |
| Operational maturity | Low | High (k8s, observability) |

> **Advice**: Start with a modular monolith. Extract services only when you feel real pain.

---

## 2. REST Best Practices, Idempotency, Versioning

### REST Design Principles

```
Resource-based URLs — nouns, not verbs:
  ✅ GET    /orders          list orders
  ✅ GET    /orders/123      get order 123
  ✅ POST   /orders          create order
  ✅ PUT    /orders/123      replace order 123 (full update)
  ✅ PATCH  /orders/123      partial update
  ✅ DELETE /orders/123      delete order 123

  ❌ GET  /getOrders
  ❌ POST /createOrder
  ❌ POST /orders/123/delete
```

**HTTP Status Codes:**

| Code | Meaning | Use |
|---|---|---|
| `200 OK` | Success with body | GET, PUT, PATCH success |
| `201 Created` | Resource created | POST success — include `Location` header |
| `204 No Content` | Success, no body | DELETE success |
| `400 Bad Request` | Invalid input | Validation errors |
| `401 Unauthorized` | Not authenticated | Missing/invalid token |
| `403 Forbidden` | Authenticated but not authorized | Insufficient permissions |
| `404 Not Found` | Resource doesn't exist | — |
| `409 Conflict` | Conflict with current state | Duplicate, optimistic lock failure |
| `422 Unprocessable Entity` | Semantically invalid | Business rule violation |
| `429 Too Many Requests` | Rate limited | Include `Retry-After` header |
| `500 Internal Server Error` | Unexpected server error | Don't expose stack traces |

---

### Idempotency

An operation is **idempotent** if calling it multiple times has the same effect as calling it once.

| Method | Idempotent | Safe (read-only) |
|---|---|---|
| `GET` | ✅ | ✅ |
| `HEAD` | ✅ | ✅ |
| `PUT` | ✅ | ❌ |
| `DELETE` | ✅ | ❌ |
| `POST` | ❌ | ❌ |
| `PATCH` | ❌ (depends on impl) | ❌ |

**Making POST idempotent** using `Idempotency-Key`:

```
Client → POST /payments
         Idempotency-Key: uuid-abc-123
         Body: {amount: 100, currency: "USD"}

Server:
  1. Check if uuid-abc-123 already processed
  2. If yes → return same response (don't charge again)
  3. If no  → process payment, store key+response
```

```java
@PostMapping("/payments")
public ResponseEntity<PaymentResponse> createPayment(
        @RequestHeader("Idempotency-Key") String idempotencyKey,
        @RequestBody PaymentRequest request) {

    // Check cache/DB for existing response
    Optional<PaymentResponse> existing = idempotencyStore.get(idempotencyKey);
    if (existing.isPresent()) {
        return ResponseEntity.ok(existing.get());
    }

    PaymentResponse response = paymentService.charge(request);
    idempotencyStore.save(idempotencyKey, response, Duration.ofHours(24));
    return ResponseEntity.status(201).body(response);
}
```

---

### API Versioning

#### URI Versioning (most common)
```
GET /api/v1/users
GET /api/v2/users   ← breaking change in a new version
```

#### Header Versioning
```
GET /api/users
Accept: application/vnd.myapp.v2+json
```

#### Query Parameter Versioning
```
GET /api/users?version=2
```

```java
// URI versioning in Spring
@RestController
@RequestMapping("/api/v1/users")
public class UserControllerV1 { ... }

@RestController
@RequestMapping("/api/v2/users")
public class UserControllerV2 { ... }
```

---

## 3. Circuit Breaker (Resilience4j)

### The Problem

Service A calls Service B. Service B is slow or down.  
Without a circuit breaker → A's threads pile up waiting → A becomes slow → cascading failure.

```
Normal:           Degraded:             Open Circuit:
  A → B ✅          A → B ⏳⏳⏳           A → fallback ✅
                    A → B ⏳⏳⏳           (B not called at all)
                    (thread pool full)    (fast fail)
```

### Circuit Breaker States

```
CLOSED ──(failure rate > threshold)──► OPEN
  ▲                                       │
  │                                       │ (wait timeout)
  │                                       ▼
  └──(half-open test succeeds)──── HALF_OPEN
                                   (allow limited calls)
```

| State | Behaviour |
|---|---|
| **CLOSED** | Normal — requests flow through, failures tracked |
| **OPEN** | Requests fail immediately (fast fail) without calling the service |
| **HALF_OPEN** | Allow N test requests; if they succeed → CLOSED; if fail → OPEN again |

---

### Resilience4j Example

```xml
<dependency>
    <groupId>io.github.resilience4j</groupId>
    <artifactId>resilience4j-spring-boot3</artifactId>
</dependency>
```

```yaml
# application.yml
resilience4j:
  circuitbreaker:
    instances:
      paymentService:
        sliding-window-type: COUNT_BASED
        sliding-window-size: 10          # track last 10 calls
        failure-rate-threshold: 50       # open if ≥ 50% fail
        wait-duration-in-open-state: 10s # stay OPEN for 10s then try HALF_OPEN
        permitted-number-of-calls-in-half-open-state: 3
        minimum-number-of-calls: 5       # don't open until at least 5 calls recorded

  retry:
    instances:
      paymentService:
        max-attempts: 3
        wait-duration: 500ms
        retry-exceptions:
          - java.io.IOException

  timelimiter:
    instances:
      paymentService:
        timeout-duration: 2s
```

```java
@Service
public class OrderService {

    @CircuitBreaker(name = "paymentService", fallbackMethod = "paymentFallback")
    @Retry(name = "paymentService")
    @TimeLimiter(name = "paymentService")
    public CompletableFuture<PaymentResult> processPayment(PaymentRequest request) {
        return CompletableFuture.supplyAsync(() -> paymentClient.charge(request));
    }

    // Fallback — called when circuit is OPEN or all retries exhausted
    public CompletableFuture<PaymentResult> paymentFallback(PaymentRequest request,
                                                             Exception ex) {
        log.warn("Payment service unavailable, queuing for later: {}", ex.getMessage());
        messageQueue.send(request);   // store for async retry
        return CompletableFuture.completedFuture(PaymentResult.queued());
    }
}
```

---

## 4. API Gateway, Service Discovery (Eureka)

### API Gateway

Single entry point for all clients — handles cross-cutting concerns:

```
Mobile App ─┐
Web App    ─┼──► API Gateway ──► Order Service
3rd Party  ─┘         │        ──► User Service
                       │        ──► Payment Service
              ┌────────┴────────┐
              │ Authentication  │
              │ Rate Limiting   │
              │ Load Balancing  │
              │ Request Routing │
              │ SSL Termination │
              │ Logging/Tracing │
              └─────────────────┘
```

**Spring Cloud Gateway:**

```yaml
spring:
  cloud:
    gateway:
      routes:
        - id: order-service
          uri: lb://ORDER-SERVICE      # lb:// = load-balanced via service discovery
          predicates:
            - Path=/api/orders/**
          filters:
            - StripPrefix=1
            - name: CircuitBreaker
              args:
                name: orderService
                fallbackUri: forward:/fallback

        - id: user-service
          uri: lb://USER-SERVICE
          predicates:
            - Path=/api/users/**
          filters:
            - AddRequestHeader=X-Internal-Token, secret123
```

---

### Service Discovery (Eureka)

Microservices register themselves; clients discover them by name instead of hardcoded IPs.

```
Service Registry (Eureka Server)
┌─────────────────────────────────┐
│  order-service  → 10.0.1.5:8080 │
│  order-service  → 10.0.1.6:8080 │
│  user-service   → 10.0.2.3:8081 │
│  payment-service→ 10.0.3.7:8082 │
└─────────────────────────────────┘
     ▲ register/heartbeat     ▲
     │                         │
Order Service              User Service
(registers on startup)   (discovers order-service by name)
```

```java
// Eureka Server
@SpringBootApplication
@EnableEurekaServer
public class DiscoveryServer { ... }

// application.yml (server)
server.port: 8761
eureka.client.register-with-eureka: false
eureka.client.fetch-registry: false
```

```java
// Eureka Client (any microservice)
@SpringBootApplication
@EnableDiscoveryClient
public class OrderServiceApp { ... }
```

```yaml
# application.yml (client)
spring.application.name: order-service
eureka:
  client:
    service-url:
      defaultZone: http://localhost:8761/eureka/
  instance:
    lease-renewal-interval-in-seconds: 10
    lease-expiration-duration-in-seconds: 30
```

```java
// Load-balanced REST call using service name
@Bean
@LoadBalanced
public RestTemplate restTemplate() { return new RestTemplate(); }

// Usage — "USER-SERVICE" is the registered service name
String user = restTemplate.getForObject("http://USER-SERVICE/api/users/1", String.class);
```

---

## 5. Event-Driven Architecture (Kafka Basics)

### Why Event-Driven?

Decouples services — producer doesn't know or wait for consumers.

```
Order Service ──publish──► [order.placed topic] ──subscribe──► Payment Service
                                                 ──subscribe──► Inventory Service
                                                 ──subscribe──► Email Service
```

### Kafka Core Concepts

```
Kafka Cluster
  ├─ Broker 1
  ├─ Broker 2
  └─ Broker 3

Topic: "order.placed"
  ├─ Partition 0: [msg0][msg1][msg4][msg7]...
  ├─ Partition 1: [msg2][msg5][msg8]...
  └─ Partition 2: [msg3][msg6][msg9]...

Consumer Group "payment-service":
  ├─ Consumer A → reads Partition 0
  └─ Consumer B → reads Partition 1, 2

  Each partition consumed by exactly one consumer in the group
  → parallel processing, ordering guaranteed within partition
```

| Concept | Description |
|---|---|
| **Topic** | Named stream of messages |
| **Partition** | Ordered, immutable log within a topic |
| **Offset** | Message position within a partition |
| **Producer** | Publishes messages to topics |
| **Consumer** | Reads messages from topics |
| **Consumer Group** | Group of consumers sharing partitions |
| **Broker** | Kafka server node |
| **Retention** | How long messages are kept (time or size based) |

---

### Spring Kafka Example

```java
// Producer
@Service
public class OrderEventPublisher {

    @Autowired
    private KafkaTemplate<String, OrderEvent> kafkaTemplate;

    public void publishOrderPlaced(Order order) {
        OrderEvent event = new OrderEvent(order.getId(), "ORDER_PLACED", order.getTotal());
        kafkaTemplate.send("order.placed", order.getId().toString(), event)
            .whenComplete((result, ex) -> {
                if (ex != null) log.error("Failed to publish event", ex);
                else log.info("Published to partition {} offset {}",
                    result.getRecordMetadata().partition(),
                    result.getRecordMetadata().offset());
            });
    }
}

// Consumer
@Component
public class PaymentEventListener {

    @KafkaListener(
        topics = "order.placed",
        groupId = "payment-service",
        concurrency = "3"          // 3 consumer threads per instance
    )
    public void handleOrderPlaced(
            @Payload OrderEvent event,
            @Header(KafkaHeaders.RECEIVED_PARTITION) int partition,
            @Header(KafkaHeaders.OFFSET) long offset) {

        log.info("Processing order {} from partition {} offset {}", event.getOrderId(), partition, offset);
        paymentService.charge(event.getOrderId(), event.getTotal());
    }
}
```

```yaml
# application.yml
spring:
  kafka:
    bootstrap-servers: localhost:9092
    producer:
      key-serializer: org.apache.kafka.common.serialization.StringSerializer
      value-serializer: org.springframework.kafka.support.serializer.JsonSerializer
      acks: all               # wait for all replicas to acknowledge
      retries: 3
    consumer:
      group-id: payment-service
      auto-offset-reset: earliest
      key-deserializer: org.apache.kafka.common.serialization.StringDeserializer
      value-deserializer: org.springframework.kafka.support.serializer.JsonDeserializer
      enable-auto-commit: false  # manual commit for at-least-once semantics
```

---

## 6. Saga Pattern, Outbox Pattern

### Saga Pattern

Manages **distributed transactions** across multiple services — replaces 2-phase commit with a sequence of local transactions and compensating transactions on failure.

#### Choreography Saga (Event-Driven)

Each service listens for events and publishes its own:

```
Order Service      Payment Service      Inventory Service
     │                    │                    │
     ├─order.placed──────►│                    │
     │                    ├─payment.completed─►│
     │                    │                    ├─inventory.reserved
     │◄───────────────────┴────────────────────┘
     │   (all success → order confirmed)

Failure path:
     │                    ├─payment.failed──────►│
     │◄──────────────────────────────────────────┘
     │   (compensate: cancel order)
```

```java
// Payment Service — listens and publishes
@KafkaListener(topics = "order.placed", groupId = "payment-saga")
public void handleOrderPlaced(OrderPlacedEvent event) {
    try {
        paymentService.charge(event.getOrderId(), event.getTotal());
        publisher.send("payment.completed", new PaymentCompletedEvent(event.getOrderId()));
    } catch (InsufficientFundsException e) {
        publisher.send("payment.failed", new PaymentFailedEvent(event.getOrderId(), e.getMessage()));
    }
}

// Order Service — listens for failure and compensates
@KafkaListener(topics = "payment.failed", groupId = "order-saga")
public void handlePaymentFailed(PaymentFailedEvent event) {
    orderService.cancelOrder(event.getOrderId());   // compensating transaction
}
```

#### Orchestration Saga (Central Coordinator)

A **Saga Orchestrator** drives the workflow and calls each service:

```
Saga Orchestrator
    │──1. charge payment──────► Payment Service
    │◄── success ──────────────────────────────
    │──2. reserve inventory──► Inventory Service
    │◄── success ──────────────────────────────
    │──3. confirm order───────► Order Service

On failure at step 2:
    │──compensate: refund─────► Payment Service
```

---

### Outbox Pattern

Solves the **dual-write problem**: atomically save to DB and publish an event.

**Problem:**
```java
// NOT atomic — what if Kafka publish fails after DB save?
orderRepository.save(order);       // saved
kafkaTemplate.send("order.placed", event);  // fails → event lost!

// Or what if DB save fails after Kafka publish?
kafkaTemplate.send("order.placed", event);  // published
orderRepository.save(order);       // fails → ghost event!
```

**Solution — Outbox table:**

```java
// 1. Write to DB AND outbox in same local transaction
@Transactional
public void placeOrder(Order order) {
    orderRepository.save(order);                  // save order

    OutboxEvent outbox = new OutboxEvent(          // save event to outbox table
        "order.placed",
        objectMapper.writeValueAsString(new OrderPlacedEvent(order))
    );
    outboxRepository.save(outbox);                 // same transaction → atomic!
}

// 2. Outbox Relay (separate process/scheduler) polls and publishes
@Scheduled(fixedDelay = 1000)
@Transactional
public void relayOutboxEvents() {
    List<OutboxEvent> pending = outboxRepository.findByPublishedFalse();
    for (OutboxEvent event : pending) {
        kafkaTemplate.send(event.getTopic(), event.getPayload());
        event.setPublished(true);         // mark as published
        outboxRepository.save(event);
    }
}
```

```sql
-- Outbox table
CREATE TABLE outbox_events (
    id          UUID PRIMARY KEY,
    topic       VARCHAR(255),
    payload     TEXT,
    published   BOOLEAN DEFAULT FALSE,
    created_at  TIMESTAMP DEFAULT NOW()
);
```

**Alternative**: Use Debezium (CDC — Change Data Capture) to stream outbox table changes directly to Kafka without a polling scheduler.

---

## 7. CAP Theorem, Eventual Consistency

### CAP Theorem

A distributed system can guarantee **at most 2 of 3** properties simultaneously:

| Property | Meaning |
|---|---|
| **C**onsistency | Every read receives the most recent write (or an error) |
| **A**vailability | Every request receives a response (not guaranteed to be latest) |
| **P**artition Tolerance | System continues despite network partition (message loss between nodes) |

**Network partitions are unavoidable** in any distributed system → you must choose between **C and A** when a partition occurs:

```
                 CAP Triangle
                     C
                    / \
              CP  /     \  CA
                /         \
               P ─────────── A
                     AP

CP systems: HBase, Zookeeper, etcd
  → Sacrifice availability for consistency under partition
  → Returns error if can't guarantee latest data

AP systems: Cassandra, CouchDB, DynamoDB (default)
  → Sacrifice consistency for availability under partition
  → Returns possibly stale data rather than an error

CA systems: Traditional RDBMS (PostgreSQL, MySQL)
  → Not partition-tolerant → single node or tight cluster
```

---

### Eventual Consistency

In an AP system, after all updates stop, all replicas will **eventually** converge to the same value — but reads may return stale data in the meantime.

```
Write "balance=100" to Node A:
  t=0  Node A: 100,  Node B: 50,  Node C: 50   ← inconsistent
  t=1  Node A: 100,  Node B: 100, Node C: 50   ← still inconsistent
  t=2  Node A: 100,  Node B: 100, Node C: 100  ← eventually consistent ✅
```

**Handling eventual consistency in practice:**

```java
// Read-your-writes: after writing, read from primary (not replica)
@Transactional(readOnly = false)
public User updateEmail(Long userId, String newEmail) {
    User user = userRepository.findById(userId).orElseThrow();
    user.setEmail(newEmail);
    return userRepository.save(user);  // write to primary
}

// For the next read in the same session, use primary
// (configure datasource routing or use sticky sessions)

// Optimistic UI: update local state immediately, sync in background
// Show user "Email updated" while async replication completes
```

**Conflict resolution strategies:**
- **Last Write Wins (LWW)**: latest timestamp wins — simple, data loss risk
- **Vector Clocks**: track causality between versions — complex but accurate
- **CRDTs** (Conflict-free Replicated Data Types): data structures that merge automatically (counters, sets)
- **Application-level merge**: business logic decides winner

### PACELC Extension

Extends CAP — even when there's **no partition (E)**, there's a trade-off between **L**atency and **C**onsistency:

```
System         Under Partition    Normal Operation
PostgreSQL     CA                 Low latency / high consistency
Cassandra      AP                 Low latency / low consistency (tunable)
DynamoDB       AP                 Low latency / eventual (tunable with strong mode)
Zookeeper      CP                 High consistency / higher latency
```
