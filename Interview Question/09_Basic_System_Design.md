# Basic System Design

1. What are functional and non-functional requirements in system design?
   - Functional requirements define capabilities; non-functional requirements define constraints and quality targets such as scale, latency, durability, and security.
2. What is scalability?
   - Scalability is a system's ability to handle growing workload by adding resources without unacceptable performance loss.
3. What is the difference between vertical and horizontal scaling?
   - Vertical scaling strengthens one machine; horizontal scaling adds machines and distributes work among them.
4. What is availability? How is it different from reliability?
   - Availability is the proportion of time a service is usable; reliability is its ability to operate correctly without failure over time.
5. What are latency and throughput?
   - Latency is the time one operation takes, while throughput is the number of operations completed per unit time.
6. What is a single point of failure?
   - It is one component whose failure can make the entire system unavailable.
7. What is load balancing, and which algorithms can a load balancer use?
   - It distributes traffic among backends using methods such as round-robin, least connections, weighted selection, or consistent hashing.
8. What is caching, and where can caching be added in a system?
   - Caching stores reusable results closer to consumers and may exist in browsers, CDNs, services, databases, and operating systems.
9. What are cache-aside, write-through, and write-back caching?
   - Cache-aside loads misses explicitly, write-through updates cache and store together, and write-back persists cached changes asynchronously.
10. What is cache invalidation?
   - It removes or refreshes stale entries when source data changes, using strategies such as TTLs, events, or explicit deletion.
11. What is database replication?
   - Replication maintains copies of data on multiple nodes for availability, read scaling, and disaster recovery.
12. What is database sharding?
   - Sharding partitions data across database nodes by a shard key to distribute storage and workload.
13. When would you choose SQL versus NoSQL?
   - Choose based on data relationships, transaction and query needs, schema flexibility, scale patterns, operational maturity, and consistency requirements.
14. What is the CAP theorem?
   - During a network partition, a distributed system cannot guarantee both complete consistency and complete availability for every request.
15. What is consistency? Compare strong and eventual consistency.
   - Strong consistency exposes the latest completed write; eventual consistency allows temporary differences but converges when updates stop.
16. What is a message queue, and why would a system use one?
   - A queue buffers messages between producers and consumers to decouple services, absorb bursts, and enable asynchronous processing.
17. What is the difference between synchronous and asynchronous communication?
   - Synchronous callers wait for a response; asynchronous callers continue while work is processed or delivered later.
18. What is the difference between a monolith and microservices?
   - A monolith deploys one cohesive application, while microservices split capabilities into independently deployable services with added distributed-system complexity.
19. What is an API gateway?
   - It is a client-facing entry point that routes requests and may centralize authentication, rate limiting, aggregation, and protocol translation.
20. What does it mean for an API operation to be idempotent?
   - Repeating the same request has the same intended effect as performing it once, which makes retries safer.
21. What are rate limiting and throttling?
   - Rate limiting enforces an allowed request quota, while throttling slows or rejects traffic when limits or capacity are reached.
22. What is a CDN, and what problem does it solve?
   - A CDN serves content from distributed edge locations, reducing user latency, origin traffic, and exposure to traffic spikes.
23. What is a reverse proxy?
   - It accepts client traffic on behalf of backend servers and can route, secure, cache, compress, or balance requests.
24. What are health checks and heartbeats?
   - Health checks actively test readiness or liveness; heartbeats are periodic signals indicating a component is still active.
25. What are retries, timeouts, and exponential backoff?
   - Timeouts bound waiting, retries repeat transient failures, and exponential backoff spaces attempts increasingly to avoid amplifying overload.
26. What is a circuit breaker?
   - It temporarily stops calls to a failing dependency, allowing recovery and preventing cascading resource exhaustion.
27. What is partitioning, and how does it differ from replication?
   - Partitioning divides different data among nodes; replication stores copies of the same data on multiple nodes.
28. What is a stateless service, and why is statelessness useful?
   - A stateless instance keeps no required client session locally, so any instance can serve a request and scaling or replacement is easier.
29. How would you remove a single point of failure?
   - Add independent redundant instances, automated failover, replicated state, health-based routing, and testing of failure scenarios.
30. What would you clarify before designing a URL shortener or chat service?
   - Clarify features, traffic and data scale, read/write patterns, latency, consistency, availability, retention, security, and cost constraints.

## Medium to Advanced

31. How would you estimate capacity for a new system from incomplete requirements?
   - **Key note:** State assumptions, derive peak QPS, storage and bandwidth from user actions, then test sensitivity and growth headroom.
32. What is consistent hashing, and when is it useful?
   - **Key note:** It distributes keys while minimizing remapping during membership changes, useful for caches and partitioned services.
33. How do replication factor and read/write quorums affect a distributed store?
   - **Key note:** Larger quorums improve overlap and consistency but increase latency and reduce availability during failures.
34. Compare leader-follower, multi-leader, and leaderless replication.
   - **Key note:** Compare write paths, conflict resolution, failover, geographic latency, and consistency guarantees.
35. What is the difference between strong consistency, linearizability, and serializability?
   - **Key note:** Linearizability concerns real-time single-operation behavior; serializability concerns transaction ordering across operations.
36. What is a distributed consensus algorithm used for?
   - **Key note:** It lets nodes agree on an ordered value or log despite limited failures, enabling leader election and replicated state.
37. What is split brain, and how can quorum prevent it?
   - **Key note:** Competing partitions both claim authority; intersecting majorities prevent two sides from committing independently.
38. How would you choose a sharding key?
   - **Key note:** Align it with primary queries while ensuring even load, stable ownership, adequate cardinality, and manageable resharding.
39. How would you reshard a live system without downtime?
   - **Key note:** Copy ranges incrementally, dual-read/write or route by version, verify consistency, cut traffic over, then clean up.
40. What is the transactional outbox pattern?
   - **Key note:** Save a domain update and pending event in one transaction, then reliably relay the event to close the dual-write gap.
41. How do Saga orchestration and choreography differ?
   - **Key note:** Orchestration centralizes workflow decisions; choreography distributes them among event-reacting services.
42. How would you make an asynchronous workflow idempotent?
   - **Key note:** Give operations stable IDs and atomically store deduplication state with each business effect.
43. What are backpressure and load shedding?
   - **Key note:** Backpressure slows producers; load shedding rejects work that cannot complete within capacity or deadlines.
44. Why do retries require exponential backoff, jitter, and a budget?
   - **Key note:** They prevent synchronized retry storms, spread recovery load, and cap traffic amplification.
45. How do bulkheads and circuit breakers prevent cascading failure?
   - **Key note:** Bulkheads isolate resource pools; circuit breakers stop calls to a failing dependency so failures do not consume everything.
46. What is a hot partition, and how can it be mitigated?
   - **Key note:** Skew directs excessive load to one owner; improve key distribution, split hot keys, cache, or isolate heavy tenants.
47. How would you design multi-region failover?
   - **Key note:** Define RTO/RPO and consistency, replicate independently, automate routing, remove regional dependencies, and test failover.
48. What is the difference between an event log and a message queue?
   - **Key note:** A log retains ordered history for independent replay; a queue primarily distributes work for consumption.
49. How would you design a globally distributed rate limiter?
   - **Key note:** Decide whether limits must be exact; partition ownership or allocate regional quotas and reconcile approximate usage.
50. How do you identify and remove a single point of failure in stateful infrastructure?
   - **Key note:** Map every dependency, add independent replicas and quorum/failover, remove shared dependencies, and exercise failure paths.
