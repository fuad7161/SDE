# 🔵 Load Balancing

> **Category:** Infrastructure &nbsp;|&nbsp; **Tags:** `L4 vs L7` `round-robin` `sticky sessions`

---

## Table of Contents
1. [What is Load Balancing?](#what-is-load-balancing)
2. [L4 vs L7 Load Balancing](#l4-vs-l7-load-balancing)
3. [Load Balancing Algorithms](#load-balancing-algorithms)
4. [Sticky Sessions](#sticky-sessions)
5. [Health Checks](#health-checks)
6. [High Availability Patterns](#high-availability-patterns)
7. [Interview Questions](#interview-questions)

---

## What is Load Balancing?

A **load balancer** distributes incoming network traffic across multiple backend servers to:
- Prevent any single server from being overwhelmed
- Improve availability and fault tolerance
- Scale horizontally

```
                     ┌─── Server A
Client ──→ Load Balancer ──→ Server B
                     └─── Server C
```

---

## L4 vs L7 Load Balancing

### L4 – Transport Layer Load Balancing

Operates on **TCP/UDP** — routes traffic based on IP and port without inspecting the content.

- Faster — no need to decrypt or parse application data.
- Cannot make routing decisions based on URL, headers, or cookies.
- Works with any TCP/UDP protocol.

```
Client → LB sees: srcIP:srcPort → dstIP:443
LB routes based on IP hash or round-robin → forwards TCP stream to chosen server
```

**Examples:** AWS Network Load Balancer (NLB), HAProxy (TCP mode).

---

### L7 – Application Layer Load Balancing

Operates on **HTTP/HTTPS** — can inspect headers, URLs, cookies, and request body.

- Can route `/api/*` to API servers and `/static/*` to CDN/file servers.
- Can route based on hostname (virtual hosting).
- Supports SSL termination (decrypt once at LB, forward HTTP internally).
- Can add/remove headers, rewrite URLs.

```
GET /api/users → routed to API server cluster
GET /images/logo.png → routed to static file server
```

**Examples:** AWS Application Load Balancer (ALB), Nginx, HAProxy (HTTP mode).

---

### Comparison

| Feature | L4 | L7 |
|---------|----|----|
| OSI Layer | Transport (4) | Application (7) |
| Protocol aware | TCP/UDP only | HTTP/HTTPS |
| Routing basis | IP + Port | URL, headers, cookies |
| SSL termination | No (pass-through) | Yes |
| Performance | Faster | Slightly more overhead |
| Content-based routing | ❌ | ✅ |

---

## Load Balancing Algorithms

| Algorithm | Description | Best For |
|-----------|-------------|---------|
| **Round Robin** | Requests distributed evenly in sequence | Equal-capacity servers |
| **Weighted Round Robin** | Like round-robin but proportional to server weight | Different-capacity servers |
| **Least Connections** | Route to server with fewest active connections | Varying request durations |
| **Least Response Time** | Route to server with lowest latency + fewest connections | Latency-sensitive apps |
| **IP Hash** | Hash source IP → always same server | Sticky sessions without cookies |
| **Random** | Random server selection | Simple, low overhead |
| **Resource Based** | Route based on actual CPU/memory usage | Dynamic workloads |

---

## Sticky Sessions

**Sticky sessions (session persistence)** ensure that all requests from a specific client are routed to the **same backend server**.

**Why needed:** Some applications store session state (e.g., shopping cart) in server memory. If the client is routed to a different server, the session is lost.

### Methods

| Method | How |
|--------|-----|
| **Cookie-based** | LB inserts a cookie with server ID (`AWSALB`, `JSESSIONID`) |
| **IP-based** | Hash the client's IP to consistently pick the same server |
| **URL parameter** | Embed server ID in URL (rarely used) |

**Drawbacks of sticky sessions:**
- Uneven load distribution if some users have long sessions.
- If the assigned server fails, the session is lost anyway.
- Hinders auto-scaling.

**Better alternative:** Store session state in a shared store (Redis, Memcached) so any server can handle the request.

---

## Health Checks

Load balancers periodically probe backend servers to detect failures.

### Types
- **Active (ping-based):** LB sends periodic HTTP/TCP requests and checks response.
- **Passive:** LB monitors actual traffic — marks server unhealthy if it returns 5xx or times out.

### Common configuration
```yaml
healthCheck:
  path: /health
  interval: 10s
  timeout: 5s
  unhealthyThreshold: 3   # fail 3 checks → mark down
  healthyThreshold: 2     # pass 2 checks → mark up
```

---

## High Availability Patterns

### Active-Passive LB Pair
Two load balancers: one active, one standby. If the active one fails, a virtual IP (VIP) floats to the passive via **VRRP/HSRP**.

### Active-Active LB Pair
Both load balancers handle traffic simultaneously. DNS round-robin points to both IPs.

### Global Load Balancing
DNS-based load balancing across regions (e.g., AWS Route 53 latency routing). Directs users to the nearest region.

---

## Interview Questions

### Q1. What is the difference between L4 and L7 load balancing?

> **Answer:**
> - **L4 (Transport):** Routes based on IP and port only. Doesn't inspect request content. Fast, works with any TCP/UDP protocol. Cannot do URL-based routing or SSL termination.
> - **L7 (Application):** Inspects HTTP headers, URLs, and cookies. Can route `/api/*` differently from `/static/*`, do SSL termination, host-based routing, and header manipulation.
>
> L7 is more flexible and powerful; L4 is faster and protocol-agnostic.

---

### Q2. What is a sticky session? When would you avoid using it?

> **Answer:**  
> Sticky sessions ensure a client always hits the same backend server, typically via a cookie containing the server ID.
>
> **Avoid it when:**
> - Servers hold no local state (stateless apps) — stickiness provides no benefit.
> - You need reliable failover — if the sticky server goes down, the session is lost regardless.
> - You're auto-scaling — new servers get no traffic until stickiness is broken.
>
> **Better alternative:** Externalize session state to Redis/Memcached and make backends stateless.

---

### Q3. What happens when a backend server behind a load balancer fails?

> **Answer:**  
> The load balancer detects the failure via health checks (active or passive). Once the server is marked unhealthy:
> - New requests are not routed to it.
> - In-flight requests may fail (returned as errors to clients).
> - If sticky sessions were used, those clients may need to re-establish session state.
>
> When the server recovers and passes health checks, it is re-added to the rotation.

---

### Q4. What is SSL termination at the load balancer? What are the pros and cons?

> **Answer:**  
> **SSL termination** means the LB decrypts HTTPS traffic and forwards unencrypted HTTP to backend servers.
>
> **Pros:**
> - Offloads CPU-intensive TLS from backend servers.
> - LB can inspect and route based on HTTP content (L7 decisions).
> - Centralized certificate management.
>
> **Cons:**
> - Traffic between LB and backend is unencrypted (mitigated if they're on a private VPC/LAN).
> - End-to-end encryption broken.
>
> **Alternative:** SSL pass-through (L4) or re-encryption (LB decrypts and re-encrypts to backend).

---

### Q5. How does the least connections algorithm work and when is it better than round-robin?

> **Answer:**  
> Least connections routes each new request to the server with the **fewest active connections** at that moment.
>
> **Better than round-robin when:**
> - Requests have **varying durations** (some quick, some long-running). Round-robin would pile up slow requests on a server while fast ones free up connections quickly.
> - Example: A mix of file uploads (slow, holds connections) and API pings (fast). Least connections keeps load more balanced.

---

<div align="center">
  <sub>← Back to <a href="networking_interview_topics.md">All Topics</a></sub>
</div>
