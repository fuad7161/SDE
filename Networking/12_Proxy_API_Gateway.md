# 🔵 Proxy & API Gateway

> **Category:** Infrastructure &nbsp;|&nbsp; **Tags:** `forward proxy` `reverse proxy` `rate limiting`

---

## Table of Contents
1. [Forward Proxy](#forward-proxy)
2. [Reverse Proxy](#reverse-proxy)
3. [Forward vs Reverse Proxy](#forward-vs-reverse-proxy)
4. [API Gateway](#api-gateway)
5. [Rate Limiting](#rate-limiting)
6. [API Gateway vs Load Balancer](#api-gateway-vs-load-balancer)
7. [Interview Questions](#interview-questions)

---

## Forward Proxy

A **forward proxy** sits between the **client and the internet**. The client sends requests to the proxy, which forwards them to the destination on the client's behalf.

```
Client ──→ [Forward Proxy] ──→ Internet / Target Server
```

**The server sees the proxy's IP, not the client's IP.**

### Use Cases
- **Anonymity:** Hide client identity (Tor, corporate proxy).
- **Content filtering:** Block certain websites (corporate/school firewalls).
- **Caching:** Cache frequently accessed content for a network.
- **Bypass geo-restrictions:** Access content available in the proxy's country.

### Example Tools
- Squid Proxy
- Nginx (forward proxy mode)
- Corporate HTTP proxy

---

## Reverse Proxy

A **reverse proxy** sits between the **internet and backend servers**. Clients send requests to the reverse proxy, which routes them to the appropriate backend server.

```
Internet / Client ──→ [Reverse Proxy] ──→ Backend Server(s)
```

**The client sees the proxy's IP, not the backend servers' IPs.**

### Use Cases
- **Load balancing:** Distribute traffic across multiple servers.
- **SSL termination:** Handle TLS at the proxy; backend runs plain HTTP.
- **Caching:** Cache backend responses (Nginx `proxy_cache`).
- **Compression:** Apply gzip before sending to clients.
- **Security:** Hide internal server structure; single entry point for WAF.
- **Authentication:** Centralize auth before routing to services.
- **URL rewriting:** Route `/api/v1` → internal service A, `/api/v2` → service B.

### Example Tools
- Nginx
- HAProxy
- AWS ALB / CloudFront
- Traefik

---

## Forward vs Reverse Proxy

| Feature | Forward Proxy | Reverse Proxy |
|---------|-------------|--------------|
| Who configures it? | Client | Server admin |
| Who does it represent? | Clients | Servers |
| What does it hide? | Client's IP from server | Server's IP from client |
| Primary use | Client anonymity, filtering | LB, SSL termination, security |
| Example tools | Squid, Tor | Nginx, HAProxy, ALB |

---

## API Gateway

An **API Gateway** is a managed reverse proxy specifically designed for **microservices and APIs**. It acts as a single entry point for all API calls and handles cross-cutting concerns.

```
Client ──→ [API Gateway] ──→ /users  → User Service
                         ──→ /orders → Order Service
                         ──→ /auth   → Auth Service
```

### API Gateway Responsibilities

| Concern | What it does |
|---------|-------------|
| **Routing** | Route requests to correct microservice based on URL, method, headers |
| **Authentication** | Validate JWT tokens, API keys, OAuth2 before forwarding |
| **Rate limiting** | Throttle requests per client/IP/API key |
| **SSL termination** | Handle HTTPS; forward HTTP internally |
| **Request/Response transformation** | Add/remove headers, transform payloads |
| **Load balancing** | Balance between instances of a service |
| **Caching** | Cache frequent responses |
| **Logging & monitoring** | Centralized access logs, metrics, tracing |
| **Circuit breaking** | Stop forwarding to failing services |
| **API versioning** | Route `/v1/` and `/v2/` to different backends |

### Example API Gateways
- AWS API Gateway
- Kong
- Nginx + Lua
- Traefik
- Apigee
- AWS AppSync (GraphQL)

---

## Rate Limiting

Rate limiting restricts the number of requests a client can make in a given time window to prevent abuse, DoS, and resource exhaustion.

### Rate Limiting Algorithms

#### Token Bucket
```
Bucket capacity: 100 tokens
Refill rate: 10 tokens/second
Each request consumes 1 token
Request allowed only if token available; else rejected (429)
```
- Allows **bursts** up to bucket capacity.
- Most common algorithm.

#### Leaky Bucket
```
Requests enter a queue (bucket); processed at a fixed rate
If bucket is full, excess requests are dropped
```
- Smooths out bursts — outputs a constant rate.
- Good for rate-limiting output (e.g., API calls to a third party).

#### Fixed Window Counter
```
Window: 1 minute
Counter: starts at 0, reset at minute boundary
Limit: 100 requests/minute
```
- Simple but vulnerable to **boundary burst**: 100 requests at 0:59 + 100 at 1:00 = 200 in 2 seconds.

#### Sliding Window Log
- Maintains a timestamp log of recent requests.
- Always checks requests within the past `window` duration.
- Accurate but memory-intensive.

#### Sliding Window Counter
- Combines fixed window counter with a weighted calculation.
- More memory-efficient approximation of sliding window.

### Rate Limit Responses
```
HTTP 429 Too Many Requests
Retry-After: 30
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 0
X-RateLimit-Reset: 1716000000
```

---

## API Gateway vs Load Balancer

| Feature | Load Balancer | API Gateway |
|---------|-------------|------------|
| Primary job | Distribute traffic | Manage API requests |
| Layer | L4 or L7 | L7 (HTTP/HTTPS) |
| Auth/AuthZ | ❌ | ✅ |
| Rate limiting | ❌ (basic) | ✅ (advanced) |
| Request transformation | ❌ | ✅ |
| API versioning | ❌ | ✅ |
| Service routing | Simple (host/URL) | Complex (method, headers, body) |
| Circuit breaking | Sometimes | ✅ |
| Cost | Lower | Higher |

**Typical architecture:** `Client → API Gateway → Load Balancer → Service instances`

---

## Interview Questions

### Q1. What is the difference between a forward proxy and a reverse proxy?

> **Answer:**
> - **Forward proxy:** Sits in front of the **client**. The client routes its requests through the proxy. The destination server sees the proxy's IP. Used for client anonymity, content filtering, caching.
> - **Reverse proxy:** Sits in front of **backend servers**. Clients think they're talking to the proxy. Used for load balancing, SSL termination, caching, security. Backend servers are hidden from the internet.

---

### Q2. What is an API Gateway and what responsibilities does it take on?

> **Answer:**  
> An API Gateway is a single entry point for all client-to-microservice communication. It handles:
> - **Routing** requests to the correct service.
> - **Authentication** — validate tokens/API keys before forwarding.
> - **Rate limiting** — throttle abusive clients.
> - **SSL termination** — handle HTTPS centrally.
> - **Request/response transformation** — modify headers, payloads.
> - **Logging & observability** — centralized access logs, metrics.
> - **Circuit breaking** — stop routing to failing services.
>
> It removes the need for each microservice to implement these concerns individually.

---

### Q3. How does the token bucket rate limiting algorithm work?

> **Answer:**  
> A **token bucket** has a maximum capacity (e.g., 100 tokens). Tokens are added at a fixed rate (e.g., 10/second). Each request consumes one token. If tokens are available, the request is allowed. If the bucket is empty, the request is rejected (HTTP 429).
>
> Key property: allows **bursts** up to the bucket size, then smooths out to the refill rate. This is more lenient than a strict fixed window and accommodates legitimate traffic spikes.

---

### Q4. What is the difference between rate limiting and throttling?

> **Answer:**
> - **Rate limiting:** Hard cap — once the limit is exceeded, requests are **rejected** (HTTP 429) until the window resets.
> - **Throttling:** Soft cap — requests are **slowed down or queued** rather than rejected outright. The server processes them at a controlled rate.
>
> In practice, the terms are often used interchangeably, but throttling implies graceful degradation while rate limiting implies hard rejection.

---

### Q5. Why would you place an API Gateway in front of microservices instead of having clients call services directly?

> **Answer:**  
> Direct client-to-service calls create several problems:
> - **Multiple round trips:** Clients make many calls to compose a single UI view.
> - **Duplication of cross-cutting concerns:** Every service would need to implement auth, rate limiting, logging, etc.
> - **Service exposure:** Internal service URLs and ports are exposed to clients, making refactoring difficult.
> - **Protocol mismatch:** Clients use HTTP/REST but some internal services may use gRPC, messaging, etc.
>
> An API Gateway centralizes all these concerns, exposes a clean stable API to clients, and decouples clients from the internal service topology.

---

<div align="center">
  <sub>← Back to <a href="networking_interview_topics.md">All Topics</a></sub>
</div>
