# Backend Interview Questions — Medium to Hard

## API Design and HTTP

1. How would you design a REST API that supports filtering, sorting, pagination, and field selection without making the interface difficult to evolve?
   - **Key note:** Use consistent query conventions, allow-listed fields, stable defaults, cursor pagination, and backward-compatible additive changes.
2. What makes an HTTP operation idempotent, and how would you make a payment-creation endpoint safely retryable?
   - **Key note:** Persist an idempotency key with the request fingerprint and original result inside the transaction boundary.
3. Compare offset-based, cursor-based, and keyset pagination. When does each approach fail?
   - **Key note:** Offset is simple but slow and unstable at scale; cursor/keyset is stable and fast but restricts navigation and ordering.
4. How would you version an API while minimizing disruption to existing clients?
   - **Key note:** Prefer compatible additions; version only breaking contracts, publish deprecation windows, and measure old-version usage.
5. What are the trade-offs between REST, GraphQL, gRPC, and asynchronous messaging for service-to-service communication?
   - **Key note:** Compare contract strength, latency, coupling, browser support, independent evolution, and delivery guarantees.
6. How would you prevent duplicate processing when a client retries after a timeout but the server completed the original request?
   - **Key note:** Combine idempotency keys, a uniqueness constraint, atomic status storage, and replay of the stored response.
7. How should an API represent validation errors, domain errors, partial success, and temporary failures?
   - **Key note:** Use stable machine-readable error codes, appropriate status codes, field details, correlation IDs, and retry guidance.
8. What is content negotiation, and when is it useful in an API?
   - **Key note:** `Accept` and `Content-Type` select representations; use it only when multiple formats are genuinely supported.
9. How would you design an API for a long-running operation such as video processing or report generation?
   - **Key note:** Return `202` with a job resource, expose status/cancellation, and notify or allow polling on completion.
10. When should an API return `200`, `201`, `202`, `204`, `409`, `422`, `429`, or `503`?
   - **Key note:** Match semantics: success, created, accepted, no body, state conflict, invalid entity, rate-limited, or temporarily unavailable.
11. How would you safely evolve a request or response schema used by independently deployed clients?
   - **Key note:** Make additive changes, tolerate unknown fields, use defaults, avoid changing meaning, and run consumer contract tests.
12. What problems can arise when a reverse proxy, CDN, or browser caches an API response incorrectly?
   - **Key note:** Private data can leak or become stale; set precise cache keys, `Vary`, authorization rules, and cache-control headers.

## Authentication, Authorization, and Security

13. What is the difference between authentication and authorization, and where should each be enforced?
   - **Key note:** Authentication proves identity; authorization checks each requested action and resource at a trusted server-side boundary.
14. Compare session-based authentication with token-based authentication.
   - **Key note:** Sessions simplify revocation but require shared state; tokens scale independently but complicate revocation and secure storage.
15. How do access tokens and refresh tokens work, and how should refresh-token rotation be implemented?
   - **Key note:** Keep access tokens short-lived; rotate refresh tokens on use and revoke the token family when reuse is detected.
16. What are the security trade-offs between opaque tokens and JWTs?
   - **Key note:** Opaque tokens enable central control; JWTs allow local validation but expose claims and make revocation harder.
17. How would you revoke a JWT before it expires?
   - **Key note:** Use short expiry plus session/version checks or a denylist for high-risk cases; rotation limits exposure.
18. What attacks become possible when JWT signature verification or algorithm validation is implemented incorrectly?
   - **Key note:** Attackers may forge identities; pin allowed algorithms, validate issuer/audience/time claims, and reject unsigned tokens.
19. Compare role-based, attribute-based, and relationship-based access control.
   - **Key note:** RBAC uses roles, ABAC evaluates attributes/policy, and ReBAC derives permission from relationships between entities.
20. How would you enforce object-level authorization and prevent insecure direct object references?
   - **Key note:** Authorize every object access against the authenticated principal; never treat an unguessable ID as permission.
21. How should passwords be stored, and why are fast hashing algorithms unsuitable?
   - **Key note:** Use a salted adaptive password hash such as Argon2id; deliberate cost makes offline guessing expensive.
22. How would you protect an API from brute-force attacks without blocking legitimate users behind a shared IP address?
   - **Key note:** Combine account, device, IP, and risk-based limits with progressive delay, MFA, alerts, and careful recovery.
23. What are CSRF, XSS, SSRF, SQL injection, and command injection, and which backend controls mitigate them?
   - **Key note:** Use CSRF defenses, output encoding/CSP, URL allow-lists, parameterized SQL, and avoid shell construction.
24. How would you store, distribute, rotate, and audit application secrets?
   - **Key note:** Use a secret manager, short-lived identities, least privilege, automated rotation, encryption, and audited access.
25. What is mutual TLS, and when would you use it between services?
   - **Key note:** Both endpoints authenticate with certificates; it provides strong workload identity but needs certificate lifecycle automation.
26. How would you design secure multi-tenant data access so one tenant can never read another tenant's data?
   - **Key note:** Derive tenant context from trusted identity and enforce it in every query, policy, cache key, and asynchronous job.

## Database Design and SQL

27. How would you choose between a relational database, document database, key-value store, wide-column store, and graph database?
   - **Key note:** Start from access patterns, relationships, transactions, consistency, scale, schema evolution, and operational capability.
28. How do B-tree and hash indexes differ, and which query patterns can each support?
   - **Key note:** B-trees support equality, ranges, and ordering; hash indexes mainly optimize exact equality lookups.
29. How does column order affect a composite index and the queries that can use it?
   - **Key note:** B-tree use generally follows the leftmost prefix; equality columns usually precede range and ordering columns.
30. What is a covering index, and what costs does it introduce?
   - **Key note:** It contains every needed column and avoids table lookup, at the cost of space and slower writes.
31. Why might a database ignore an available index?
   - **Key note:** Low selectivity, stale statistics, functions/type casts, mismatched leading columns, or a cheaper scan can be responsible.
32. How would you read and analyze an SQL execution plan?
   - **Key note:** Compare estimated versus actual rows, scan types, join order, filters, sorts, memory spills, loops, and total cost/time.
33. What causes the N+1 query problem, and how can it be fixed without loading excessive data?
   - **Key note:** Per-row lazy loading causes it; use targeted joins, batching, prefetching, or data-loader patterns with bounded results.
34. Compare optimistic and pessimistic locking. Under what contention patterns does each perform poorly?
   - **Key note:** Optimistic locking retries conflicts and suits low contention; pessimistic locking blocks and risks deadlocks under long transactions.
35. How can lost updates occur, and how would you prevent them?
   - **Key note:** Concurrent read-modify-write cycles overwrite each other; use atomic updates, version checks, or suitable isolation/locking.
36. How do MVCC databases provide isolation without locking every read?
   - **Key note:** Readers use transaction-visible row versions while writers create new versions; cleanup later reclaims obsolete data.
37. Compare Read Committed, Repeatable Read, Snapshot Isolation, and Serializable isolation.
   - **Key note:** Each blocks more anomalies; only true serializable execution guarantees equivalence to some serial transaction order.
38. What are write skew and serialization anomalies?
   - **Key note:** Transactions update different rows after reading a shared invariant; serializable isolation or explicit locking prevents it.
39. How would you investigate and reduce database deadlocks?
   - **Key note:** Inspect deadlock graphs, lock resources in consistent order, shorten transactions, index queries, and retry victims safely.
40. How would you perform a zero-downtime migration on a large, heavily used table?
   - **Key note:** Use expand/backfill/switch/contract phases, small batches, compatibility across versions, monitoring, and a rollback path.
41. How would you add a non-null column to a table containing hundreds of millions of rows?
   - **Key note:** Add it nullable or with metadata-only default, deploy dual behavior, backfill in batches, validate, then enforce non-null.
42. What are the trade-offs between UUIDs, sequential numeric IDs, and distributed time-based IDs as primary keys?
   - **Key note:** Compare predictability, locality/index fragmentation, size, coordination, global uniqueness, and timestamp leakage.
43. How would you design soft deletion while preserving uniqueness constraints and query performance?
   - **Key note:** Use a deletion timestamp, mandatory query scope, and partial or composite uniqueness indexes; plan retention and purge.
44. When should data be normalized, and when is denormalization justified?
   - **Key note:** Normalize for integrity and maintainability; denormalize measured read bottlenecks with an explicit synchronization strategy.
45. How would you model hierarchical or tree-structured data in a relational database?
   - **Key note:** Choose adjacency lists, materialized paths, nested sets, or closure tables according to read/write and subtree-query needs.
46. How would you diagnose replication lag, and how can it affect read-after-write consistency?
   - **Key note:** Measure replay delay and bottlenecks; route critical reads to the leader or use a consistency token/session stickiness.

## Transactions and Data Consistency

47. What guarantees does a local ACID transaction provide, and what does it not guarantee across services?
   - **Key note:** It protects one transactional resource; remote calls and independently committed stores are outside its atomic boundary.
48. Why is a distributed transaction difficult, and when would two-phase commit be appropriate?
   - **Key note:** Coordination reduces availability and can block; use 2PC only when participants support it and strict atomicity outweighs cost.
49. What is the Saga pattern, and how do orchestration and choreography differ?
   - **Key note:** A saga uses local transactions plus compensation; orchestration centralizes flow, while choreography reacts to events.
50. How would you design compensating actions when an operation cannot be perfectly reversed?
   - **Key note:** Model business remediation explicitly, preserve an audit trail, make compensation idempotent, and allow manual resolution.
51. What is the transactional outbox pattern, and which failure does it prevent?
   - **Key note:** Store the business change and event in one DB transaction, then relay it; this closes the database/message dual-write gap.
52. What does “exactly-once delivery” actually mean, and why is it usually achieved through processing semantics rather than transport alone?
   - **Key note:** Networks retry and duplicate; the observable effect becomes once through idempotency, deduplication, and atomic state changes.
53. How would you build an idempotent message consumer?
   - **Key note:** Persist a message/event ID with the business update atomically and return success for already processed IDs.
54. What are monotonic reads, read-your-writes, and causal consistency?
   - **Key note:** They prevent time going backward, expose a client's own writes, and preserve cause-before-effect ordering, respectively.
55. When is eventual consistency acceptable, and how should it be reflected in the user experience?
   - **Key note:** Use it where temporary staleness is tolerable; show pending state, set expectations, and provide reconciliation.
56. How would you reconcile conflicting updates made while two replicas were disconnected?
   - **Key note:** Use domain rules, versions/vector clocks, CRDTs, deterministic merge, or surface conflicts for human resolution.
57. What is a fencing token, and how does it prevent a stale lock holder from modifying data?
   - **Key note:** Each lock grant has a rising token; the resource rejects operations carrying an older token.

## Caching

58. Compare cache-aside, read-through, write-through, write-around, and write-back caching.
   - **Key note:** They differ in who loads data and when writes reach cache/storage; compare consistency, latency, and data-loss risk.
59. How would you prevent stale cache entries after a database update?
   - **Key note:** Invalidate after commit, publish change events through an outbox, use TTLs, and tolerate unavoidable race windows.
60. What are cache stampede, cache penetration, and cache avalanche, and how can each be mitigated?
   - **Key note:** Use request coalescing/jitter, negative or Bloom-filter checks, and staggered expirations with capacity protection.
61. How would you choose a cache key and expiration policy for multi-tenant data?
   - **Key note:** Include tenant, version, identity, and query dimensions; base TTL on freshness, risk, update rate, and load.
62. What problems can occur when using distributed locks to rebuild cached values?
   - **Key note:** Lock expiry, pauses, and partitions can create concurrent owners; use bounded waits, tokens, and idempotent rebuilds.
63. When would negative caching be useful, and what risks does it create?
   - **Key note:** It protects against repeated misses, but can hide newly created data; use short TTLs and invalidate on creation.
64. How would you cache paginated or filtered query results without creating an unmanageable number of keys?
   - **Key note:** Cache stable IDs or common query shapes, cap variants, use canonical keys, and avoid low-reuse combinations.
65. Compare LRU, LFU, FIFO, and TTL-based eviction policies.
   - **Key note:** They evict by recency, frequency, insertion order, or age; workload locality determines the best fit.
66. How would you measure whether adding a cache actually improves the system?
   - **Key note:** Compare end-to-end percentiles, origin load, hit quality, error rate, cost, and invalidation correctness.
67. Why can a high cache hit rate still coexist with poor latency?
   - **Key note:** Hot operations may be slow, misses may dominate tail latency, or serialization/network/contention may remain bottlenecks.

## Concurrency and Multithreading

68. What is the difference between concurrency, parallelism, asynchronous execution, and multithreading?
   - **Key note:** They mean overlapping progress, simultaneous execution, non-blocking coordination, and multiple threads, respectively.
69. What is a race condition, and why can it disappear while debugging?
   - **Key note:** Correctness depends on timing; logging and breakpoints alter scheduling and can hide the faulty interleaving.
70. Compare mutexes, semaphores, read-write locks, condition variables, and spinlocks.
   - **Key note:** Choose by ownership, permit count, read/write ratio, waiting condition, contention, and expected wait duration.
71. What conditions cause deadlock, and how would you diagnose one in production?
   - **Key note:** Look for mutual exclusion, hold-and-wait, no preemption, and circular wait using thread and lock graphs.
72. What are starvation, livelock, and priority inversion?
   - **Key note:** A task is denied progress, tasks react without progress, or a high-priority task waits behind lower-priority work.
73. What is the difference between thread safety, reentrancy, and immutability?
   - **Key note:** Thread-safe code tolerates concurrency, reentrant code tolerates overlapping invocation, and immutable state cannot change.
74. How do atomic operations and compare-and-swap work?
   - **Key note:** CAS replaces a value only if it still equals an expected value, enabling lock-free state transitions with retry loops.
75. What is the ABA problem in lock-free algorithms?
   - **Key note:** A value changes A→B→A, fooling CAS; tagged/versioned pointers help reveal the intermediate change.
76. What are memory visibility and memory ordering, and why is atomicity alone insufficient?
   - **Key note:** Threads need happens-before guarantees so writes become visible in the required order, not merely indivisible operations.
77. How would you determine an appropriate thread-pool size for CPU-bound and I/O-bound workloads?
   - **Key note:** CPU work stays near core count; I/O work can use more threads based on measured wait/compute ratio and capacity.
78. What happens when an application creates an unbounded number of threads or tasks?
   - **Key note:** Memory, scheduling, queues, and dependencies saturate, increasing latency until the service collapses.
79. How do backpressure and bounded queues protect a service under load?
   - **Key note:** They limit admitted work and push overload upstream through blocking, rejection, shedding, or demand signaling.
80. When can parallel processing make an operation slower?
   - **Key note:** Coordination, contention, serialization, context switching, small tasks, and limited downstream capacity can exceed gains.
81. How would you safely coordinate updates to shared state across multiple application instances?
   - **Key note:** Prefer database atomicity, optimistic versions, partition ownership, or queues over fragile in-memory or distributed locks.

## Messaging and Event-Driven Systems

82. When should services communicate through a message broker instead of synchronous HTTP or RPC?
   - **Key note:** Use messaging for decoupling, buffering, fan-out, and delayed work when immediate response is unnecessary.
83. Compare queues, publish-subscribe topics, event streams, and logs.
   - **Key note:** Compare consumer delivery, fan-out, retention, replay, ordering, and partitioning semantics.
84. What are at-most-once, at-least-once, and exactly-once delivery semantics?
   - **Key note:** They trade possible loss, possible duplication, and tightly scoped deduplicated effects; end-to-end guarantees matter.
85. How would you handle duplicate, delayed, lost, and out-of-order messages?
   - **Key note:** Use durable delivery, idempotency, sequence/version checks, buffering, expiry, retries, and reconciliation.
86. How can ordering be preserved for one entity while still processing many entities in parallel?
   - **Key note:** Partition by entity key and process each partition serially while running partitions concurrently.
87. What is consumer-group rebalancing, and how can it affect latency and correctness?
   - **Key note:** Partition ownership moves between consumers, pausing work and potentially replaying messages; commits must be coordinated.
88. How would you choose a partition key for an event stream?
   - **Key note:** Preserve the needed ordering boundary while distributing load evenly and allowing future scale.
89. What happens when one partition receives most of the traffic?
   - **Key note:** It becomes a hot partition that limits parallelism; redesign or salt the key if ordering permits.
90. How should retries, exponential backoff, and dead-letter queues be designed?
   - **Key note:** Retry only transient faults with jitter and limits; send poison messages to an observable, replayable quarantine.
91. When can a dead-letter queue become a place where failures are silently forgotten?
   - **Key note:** Without ownership, alerts, diagnostics, retention, and a safe replay process, DLQ data becomes permanent hidden loss.
92. What is event schema evolution, and how do you maintain backward and forward compatibility?
   - **Key note:** Prefer additive optional fields, tolerant readers, explicit versions, schema validation, and staged producer/consumer rollout.
93. What is event sourcing, and how does it differ from merely publishing domain events?
   - **Key note:** Event sourcing treats the event log as authoritative state; domain events may only describe changes to separately stored state.
94. What are the trade-offs of CQRS?
   - **Key note:** Independent read/write models improve specialization but add synchronization, eventual consistency, and operational complexity.
95. How would you replay historical events without triggering unintended external side effects?
   - **Key note:** Separate pure state projection from effects, mark replay context, and make external actions idempotent or disabled.

## Distributed Systems

96. What does the CAP theorem say, and what common interpretations of it are incorrect?
   - **Key note:** During a partition choose consistency or availability per operation; CAP is not a normal three-way database selection.
97. What is the difference between network latency, timeout, failure, and partition from a service's perspective?
   - **Key note:** A timeout reveals uncertainty, not failure; the remote operation may be slow, unreachable, failed, or already completed.
98. Why is it impossible to distinguish a failed node from a very slow node with certainty?
   - **Key note:** In an asynchronous network there is no known upper response bound, so silence is ambiguous.
99. Compare leader-follower, multi-leader, and leaderless replication.
   - **Key note:** Compare write routing, failover, conflict handling, latency, consistency, and partition behavior.
100. How is consensus different from replication?
   - **Key note:** Replication copies data; consensus makes nodes agree on one ordered value or log despite failures.
101. What problem do consensus algorithms such as Raft solve?
   - **Key note:** They elect a leader and maintain an agreed replicated log under defined crash and network assumptions.
102. What is split brain, and how can quorum-based systems reduce its risk?
   - **Key note:** Multiple groups believe they are authoritative; intersecting majorities prevent two sides from committing independently.
103. How do read and write quorums affect consistency and availability?
   - **Key note:** With `R + W > N` reads overlap writes, but larger quorums increase latency and reduce partition availability.
104. What are vector clocks or version vectors used for?
   - **Key note:** They represent causal versions and detect whether updates precede one another or conflict concurrently.
105. Why are wall clocks unsafe for ordering distributed events?
   - **Key note:** Clocks drift, jump, and synchronize imperfectly, so timestamps do not reliably establish causality.
106. What are logical clocks and hybrid logical clocks?
   - **Key note:** Logical clocks order causality; hybrid clocks combine that ordering with approximate physical time.
107. How would you generate globally unique, roughly ordered IDs without a single database sequence?
   - **Key note:** Combine time, node/region identity, and a per-tick sequence; handle clock rollback and information leakage.
108. What is consistent hashing, and why are virtual nodes useful?
   - **Key note:** It limits key movement when membership changes; virtual nodes improve balance and capacity weighting.
109. How would you design a distributed lock, and what guarantees would its users need?
   - **Key note:** Specify lease, ownership, renewal, failure, and clock assumptions; protect resources with fencing tokens.
110. What is the difference between failover, fault tolerance, and disaster recovery?
   - **Key note:** Failover switches instances, fault tolerance maintains service during faults, and disaster recovery restores after major loss.
111. How would you test a system's behavior during partial network failure?
   - **Key note:** Inject latency, loss, partitions, duplication, and asymmetric reachability while verifying invariants and recovery.

## Microservices and Service Architecture

112. When is a modular monolith a better choice than microservices?
   - **Key note:** Prefer it when one team/domain can deploy together and distributed-system independence does not justify operational cost.
113. How would you identify appropriate service boundaries?
   - **Key note:** Follow business capabilities, ownership, data/invariants, change cadence, and independent scaling—not technical layers.
114. Why is sharing one database among microservices often problematic?
   - **Key note:** It couples schemas, releases, ownership, and failure modes while allowing services to bypass each other's rules.
115. How should services share data without becoming tightly coupled?
   - **Key note:** Own data per service and expose stable APIs/events, accepting deliberate replication for independent reads.
116. What is the difference between an API gateway, reverse proxy, load balancer, and service mesh?
   - **Key note:** They focus on client API policy, server mediation, traffic distribution, and service-to-service traffic policy, respectively.
117. How does service discovery work in a dynamically scaled environment?
   - **Key note:** Instances register or are discovered through platform DNS/registry; clients or proxies route only to healthy endpoints.
118. What are the benefits and costs of sidecar proxies and service meshes?
   - **Key note:** They standardize traffic security and telemetry but add latency, resources, control-plane complexity, and harder debugging.
119. How would you prevent cascading failures across dependent services?
   - **Key note:** Apply deadlines, bounded concurrency, circuit breakers, bulkheads, shedding, graceful degradation, and controlled retries.
120. What are bulkheads, circuit breakers, deadlines, and retry budgets?
   - **Key note:** They isolate resources, stop futile calls, bound total wait, and cap retry amplification.
121. Why can retries make an overloaded system worse?
   - **Key note:** Retries multiply traffic while capacity is already exhausted; use budgets, backoff, jitter, and load shedding.
122. How would you propagate request deadlines and cancellation across services?
   - **Key note:** Carry remaining deadline in request context and cancel downstream work, queries, and queued tasks when it expires.
123. What is graceful degradation, and how would you design for it?
   - **Key note:** Preserve core behavior when dependencies fail by disabling optional features, serving safe stale data, or simplifying responses.
124. How do you maintain backward compatibility during independent service deployments?
   - **Key note:** Use additive contracts, tolerant readers, consumer tests, versioned events, and expand-contract rollouts.
125. How would you migrate functionality from a monolith to services incrementally?
   - **Key note:** Use a strangler approach: establish ownership, route one capability outward, synchronize data temporarily, then retire old paths.

## Scalability and Performance

126. How would you estimate traffic, storage, bandwidth, and memory requirements before designing a backend?
   - **Key note:** Start with users and actions, derive peak QPS and object sizes, include retention/replication, and state assumptions.
127. What is the difference between latency percentiles and average latency, and why are p95 and p99 important?
   - **Key note:** Averages hide slow tails; percentiles show the experience of the slowest meaningful fraction of requests.
128. What is Little's Law, and how can it help reason about concurrent requests?
   - **Key note:** Average concurrency equals throughput times average time in system: `L = λW` under stable conditions.
129. How would you identify whether a service is CPU-bound, memory-bound, I/O-bound, or lock-bound?
   - **Key note:** Correlate saturation metrics with profiles, allocation/GC data, I/O waits, queue depths, and lock contention.
130. What is connection pooling, and what happens when a pool is too small or too large?
   - **Key note:** Too small queues requests; too large overwhelms the dependency and increases contention and memory use.
131. How can database connection pools cause a cascading failure?
   - **Key note:** Slow queries occupy connections, callers queue and time out, retries add load, and upstream resources exhaust.
132. What is load shedding, and when should a service reject work?
   - **Key note:** Reject low-priority or excess requests before saturation when queues/deadlines show work cannot complete usefully.
133. What is the difference between horizontal partitioning, vertical partitioning, sharding, and replication?
   - **Key note:** Split rows, split columns/features, distribute partitions across nodes, or copy the same data, respectively.
134. How do hot keys and hot partitions arise, and how can they be mitigated?
   - **Key note:** Skewed access overloads one owner; cache, split/salt keys, isolate celebrities, or rebalance if ordering allows.
135. When would batch processing be preferable to processing one item at a time?
   - **Key note:** Batch when fixed per-operation overhead dominates and added latency, memory, and partial-failure handling are acceptable.
136. How would you optimize a high-throughput write-heavy system?
   - **Key note:** Batch sequential writes, partition evenly, minimize indexes, buffer asynchronously, and define durability/consistency trade-offs.
137. How would you optimize a read-heavy system with strict freshness requirements?
   - **Key note:** Optimize indexes and read models, scale authoritative reads, and use coherent invalidation or version-aware caches.
138. What measurements would you collect before and after a performance optimization?
   - **Key note:** Compare representative throughput, percentiles, errors, saturation, resources, cost, and correctness under controlled load.

## Reliability and Resilience

139. How would you define an SLI, SLO, and SLA for a backend service?
   - **Key note:** SLI is measured behavior, SLO is its target, and SLA is the external commitment and consequence.
140. What is an error budget, and how should it influence release decisions?
   - **Key note:** It is allowable unreliability under the SLO; rapid consumption should shift effort from features to reliability.
141. How do liveness, readiness, and startup health checks differ?
   - **Key note:** Liveness decides restart, readiness decides traffic eligibility, and startup allows slow initialization before liveness applies.
142. Why can poorly designed health checks cause an outage?
   - **Key note:** Shared dependency checks or tight thresholds can remove every instance or trigger restart loops during transient faults.
143. How would you perform graceful shutdown without losing in-flight requests or messages?
   - **Key note:** Stop admission, mark unready, drain with a deadline, finish/return work, commit offsets safely, then close resources.
144. What is the difference between active-active and active-passive deployment?
   - **Key note:** Active-active serves from multiple sites and handles consistency; active-passive keeps standby capacity and simpler writes.
145. How would you design for recovery from a complete region failure?
   - **Key note:** Replicate data independently, automate traffic failover, remove regional dependencies, define consistency, and test regularly.
146. What are RTO and RPO, and how do they affect disaster-recovery architecture?
   - **Key note:** RTO limits restoration time and RPO limits acceptable data loss; stricter targets demand costlier replication and automation.
147. How would you verify that backups can actually be restored?
   - **Key note:** Run scheduled isolated restores, validate integrity and application behavior, measure RTO/RPO, and record results.
148. What is chaos engineering, and what safeguards should a chaos experiment have?
   - **Key note:** Test resilience with controlled failures; set a hypothesis, small blast radius, monitoring, stop conditions, and rollback.
149. What is a brownout strategy?
   - **Key note:** Temporarily disable optional expensive features under stress so essential service remains available.

## Observability and Production Debugging

150. What is the difference between logs, metrics, traces, profiles, and audit records?
   - **Key note:** They capture events, aggregates, request paths, resource hotspots, and security-relevant accountability.
151. How would you correlate one request across multiple services and asynchronous consumers?
   - **Key note:** Propagate trace/correlation IDs and causal metadata through headers and messages while linking new asynchronous spans.
152. What makes a log structured, useful, and safe?
   - **Key note:** Use consistent fields, severity, context, and actionable events while redacting secrets and controlling volume.
153. Why should sensitive data and high-cardinality values be handled carefully in telemetry?
   - **Key note:** Sensitive values create exposure; unbounded labels make metric storage and queries prohibitively expensive.
154. What are RED and USE monitoring methods?
   - **Key note:** RED tracks rate/errors/duration for services; USE tracks utilization/saturation/errors for resources.
155. How would you investigate a sudden increase in p99 latency when average latency is unchanged?
   - **Key note:** Segment slow traces by instance, endpoint, dependency, GC, queue, lock, payload, and recent change.
156. How would you diagnose a memory leak in a long-running backend service?
   - **Key note:** Confirm post-GC growth, compare heap/allocation profiles over time, and inspect retaining references and native memory.
157. How would you investigate steadily increasing CPU use after a deployment?
   - **Key note:** Compare versions and workload, inspect CPU profiles, loops/retries, serialization, GC, logging, and dependency behavior.
158. How would you determine whether timeouts originate in the client, network, proxy, application, or database?
   - **Key note:** Align timestamps and traces across hops, inspect configured deadlines and connection events, then locate the consumed interval.
159. What signals indicate thread-pool, connection-pool, or queue exhaustion?
   - **Key note:** Look for maxed active counts, growing waits/depth, acquisition timeouts, rejected work, and rising tail latency.
160. How would you debug an intermittent production issue that cannot be reproduced locally?
   - **Key note:** Capture correlation context, compare affected cohorts, add targeted telemetry, test hypotheses, and preserve evidence before restart.
161. How would you reduce alert fatigue while still detecting meaningful failures quickly?
   - **Key note:** Alert on user-impacting symptoms and SLO burn, deduplicate, route ownership, tune windows, and remove non-actionable alerts.
162. What information should an incident runbook contain?
   - **Key note:** Include symptoms, impact checks, dashboards, safe diagnostics, mitigation/rollback, escalation, ownership, and verification.
163. What should a blameless post-incident review produce?
   - **Key note:** Produce a factual timeline, contributing system conditions, detection/response lessons, and owned prioritized actions.

## Deployment and Operations

164. Compare rolling, blue-green, and canary deployments.
   - **Key note:** Rolling gradually replaces instances, blue-green switches full environments, and canary exposes a small measured audience first.
165. How would you roll back a deployment that includes an incompatible database migration?
   - **Key note:** Avoid this trap with backward-compatible expand-contract migrations; otherwise restore compatibility forward before reverting code.
166. Why should database schema changes generally be backward compatible during deployment?
   - **Key note:** Old and new application versions coexist during rollout and rollback, so both must operate on the transitional schema.
167. What is the expand-and-contract migration pattern?
   - **Key note:** Add compatible structure, migrate reads/writes and backfill, then remove old structure only after all consumers move.
168. How should configuration differ from secrets, and how should both be delivered to an application?
   - **Key note:** Configuration may be visible; secrets require encrypted managed delivery, narrow identity access, rotation, and no logging.
169. What are the operational trade-offs between containers, virtual machines, and serverless functions?
   - **Key note:** Compare isolation, startup, control, density, autoscaling, execution limits, portability, and operational ownership.
170. How do container CPU and memory limits affect application behavior?
   - **Key note:** CPU limits can throttle; exceeding memory limits causes termination, while runtimes may need explicit limit awareness.
171. What causes a container to be terminated for out-of-memory use, and how would you investigate it?
   - **Key note:** Usage crosses its cgroup limit; inspect working set, heap/native allocation, limits, traffic, GC, and kernel events.
172. How would you safely deploy a change that modifies a widely used event schema?
   - **Key note:** Add compatible fields, update tolerant consumers first, validate registry rules, then change producers and retire old fields later.
173. What checks should a CI/CD pipeline perform before a backend release reaches production?
   - **Key note:** Compile, lint, test, scan, validate artifacts/config/migrations/contracts, deploy progressively, and verify health with rollback readiness.

## Practical Backend Scenarios

174. Design a rate limiter that works across many application instances.
   - **Key note:** Define scope and burst behavior; use token/sliding algorithms with atomic shared state or consistently partitioned ownership.
175. Design a URL-shortening service that supports custom aliases, expiration, analytics, and very high read traffic.
   - **Key note:** Focus on ID generation, unique aliases, cache-heavy redirects, expiry, asynchronous analytics, and abuse controls.
176. Design a notification service that supports email, SMS, and push delivery with retries and user preferences.
   - **Key note:** Separate orchestration from channel workers; persist intent, apply preferences, deduplicate, rate-limit, retry, and track status.
177. Design an inventory-reservation flow that prevents overselling during a flash sale.
   - **Key note:** Use atomic conditional decrement or serialized ownership, expiring reservations, idempotency, and reconciliation.
178. Design a payment workflow that safely handles duplicate callbacks, timeouts, and refunds.
   - **Key note:** Use an explicit state machine, idempotency, immutable ledger, verified callbacks, reconciliation, and auditable compensation.
179. Design a file-upload service for large files with resumable uploads and malware scanning.
   - **Key note:** Upload chunks directly to object storage, verify checksums, finalize atomically, quarantine, scan asynchronously, then publish.
180. Design a job scheduler that guarantees a task is not silently lost and tolerates worker crashes.
   - **Key note:** Persist schedules, lease jobs with heartbeats, retry expired leases, make handlers idempotent, and expose failed work.
181. Design an audit-log service whose records cannot be silently altered or deleted.
   - **Key note:** Use append-only access, immutable/WORM storage, hash chaining or signatures, strict authorization, and independent retention.
182. Design a real-time chat backend that preserves message ordering within a conversation.
   - **Key note:** Partition by conversation, assign server sequence numbers, persist before fan-out, and support reconnect/catch-up.
183. Design a webhook delivery system with signing, retries, deduplication, and observability.
   - **Key note:** Persist deliveries, sign timestamped payloads, use stable event IDs, bounded retries with jitter, DLQ, and delivery logs.
184. Design a feature-flag service that remains safe when its control plane is unavailable.
   - **Key note:** Evaluate from local cached snapshots, stream versioned updates, define fail-safe defaults, and audit changes.
185. Design a search autocomplete backend with low latency and frequent data updates.
   - **Key note:** Use prefix/trie or search indexes, popularity ranking, cached hot prefixes, incremental indexing, and typo/abuse controls.
186. A deployment doubled database CPU but request volume did not change. How would you investigate?
   - **Key note:** Compare query fingerprints/plans, call counts, indexes, connection behavior, N+1 patterns, locks, and feature/config changes.
187. An API occasionally creates duplicate orders even though the client sends one request. What failure paths would you examine?
   - **Key note:** Check proxy/client retries, message redelivery, concurrent handlers, timeout uncertainty, weak uniqueness, and non-atomic deduplication.
188. A service is healthy at low traffic but collapses suddenly beyond a threshold. What mechanisms could cause this nonlinear behavior?
   - **Key note:** Queueing, pool saturation, GC, lock contention, cache misses, retries, and downstream limits create positive feedback.
189. Users sometimes read stale data immediately after updating it. How would you locate and correct the consistency problem?
   - **Key note:** Trace write/read routes through replicas and caches; use leader reads, consistency tokens, stickiness, or correct invalidation.

## Engineering Judgment

190. How do you decide whether to build a component, use an open-source project, or buy a managed service?
   - **Key note:** Compare strategic differentiation, maturity, total cost, expertise, lock-in, security, reliability, and exit options.
191. How would you evaluate a backend design when requirements are incomplete or expected to change?
   - **Key note:** State assumptions, identify irreversible choices, preserve options at likely change points, and validate incrementally.
192. When is duplication preferable to abstraction?
   - **Key note:** Duplicate temporarily when similarity is accidental or requirements may diverge; abstract after the stable shared concept is clear.
193. How do you balance correctness, availability, latency, complexity, and cost?
   - **Key note:** Rank explicit business/SLO constraints, quantify trade-offs, choose the simplest sufficient design, and document rejected alternatives.
194. What signs indicate that a service should be split, combined, or redesigned?
   - **Key note:** Look for ownership friction, coupled changes, unclear boundaries, scaling mismatch, chatty calls, and operational overhead.
195. How would you introduce a major architectural change without stopping feature delivery?
   - **Key note:** Deliver thin reversible slices, use compatibility layers and flags, migrate measurable cohorts, and retire old paths gradually.
196. What questions would you ask before accepting a new infrastructure dependency?
   - **Key note:** Examine guarantees, limits, failure modes, operations, security, observability, cost, portability, ownership, and recovery.
197. How do you distinguish essential complexity from accidental complexity in a backend system?
   - **Key note:** Essential complexity comes from the domain and guarantees; accidental complexity comes from chosen tools, boundaries, and implementation.
