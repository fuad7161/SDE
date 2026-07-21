# 🟢 WebSockets & Real-time Communication

> **Category:** Protocols &nbsp;|&nbsp; **Tags:** `polling vs push` `SSE` `WS upgrade`

---

## Table of Contents
1. [The Real-time Problem](#the-real-time-problem)
2. [Polling Techniques](#polling-techniques)
3. [WebSockets](#websockets)
4. [Server-Sent Events (SSE)](#server-sent-events-sse)
5. [Comparison Table](#comparison-table)
6. [WebSocket vs HTTP/2 Push vs SSE](#websocket-vs-http2-push-vs-sse)
7. [Interview Questions](#interview-questions)

---

## The Real-time Problem

HTTP is fundamentally a **request-response** protocol — the server cannot push data to the client without a client request. For real-time applications (chat, live scores, notifications), the client needs to receive updates **as soon as they happen**.

**Approaches:**
1. Polling (client asks repeatedly)
2. Long Polling (client asks and waits)
3. Server-Sent Events (server pushes one-way)
4. WebSockets (full-duplex, bidirectional)

---

## Polling Techniques

### Short Polling
```
Client ──── GET /updates ────→ Server → {data: []}     (1s)
Client ──── GET /updates ────→ Server → {data: []}     (2s)
Client ──── GET /updates ────→ Server → {data: [msg1]} (3s)
```
- Client sends requests at fixed intervals (e.g., every 1 second).
- **Problem:** Wasteful — most requests return empty; high server load.
- **Use when:** Data doesn't need to be truly real-time, polling interval can be large.

---

### Long Polling
```
Client ──── GET /updates ────→ Server (holds connection open)
                               Server waits until data is available...
                               Server ←── New data arrives
Client ←─── Response {msg1} ──  Server  (closes connection)
Client ──── GET /updates ────→ Server (immediately re-connects)
```
- Server holds the connection open until data is available or timeout.
- Lower latency than short polling — responds immediately when data arrives.
- **Problem:** High connection overhead; not efficient for many concurrent clients.

---

## WebSockets

**WebSockets** provide a **persistent, full-duplex** (bidirectional) channel over a single TCP connection.

### WebSocket Upgrade Handshake

WebSocket starts as an HTTP request, then **upgrades** the protocol:

```
Client → Server:
GET /chat HTTP/1.1
Host: example.com
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==
Sec-WebSocket-Version: 13

Server → Client:
HTTP/1.1 101 Switching Protocols
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Accept: s3pPLMBiTxaQ9kYGzzhZRbK+xOo=
```

After `101 Switching Protocols`, the TCP connection is kept open and both sides can send **frames** at any time.

### WebSocket Frame Format
```
┌──────────┬───────────┬──────────┬───────────────────┐
│ FIN, Opcode │ Mask, Payload length │ Masking Key │ Payload │
└──────────┴───────────┴──────────┴───────────────────┘
```

- **Opcode:** Text (0x1), Binary (0x2), Close (0x8), Ping (0x9), Pong (0xA)
- **Masking:** Client→server frames must be masked (XOR with 4-byte key).

### Use Cases for WebSockets
- Chat applications
- Multiplayer games
- Collaborative editing (Google Docs-style)
- Live trading dashboards
- Real-time notifications

---

## Server-Sent Events (SSE)

**SSE** is a simpler alternative to WebSockets for **one-way** server-to-client streaming over HTTP.

### How SSE Works

Server sends a stream of `text/event-stream` content:

```
HTTP/1.1 200 OK
Content-Type: text/event-stream
Cache-Control: no-cache

data: {"price": 150.23}

data: {"price": 150.25}

event: alert
data: {"message": "Market closed"}
```

### Client-side (Browser)
```javascript
const es = new EventSource('/stock-feed');
es.onmessage = (e) => console.log(JSON.parse(e.data));
es.addEventListener('alert', (e) => console.log(e.data));
```

### SSE Features
- **Automatic reconnection** with `Last-Event-ID` header.
- **Named events** for different types of messages.
- Works over standard HTTP/1.1 and HTTP/2.
- Limited to **text** data (no binary).
- **One-way only** (server → client).

---

## Comparison Table

| Feature | Short Polling | Long Polling | SSE | WebSocket |
|---------|--------------|-------------|-----|-----------|
| Direction | Client → Server | Client → Server | Server → Client | Bidirectional |
| Connection | New per request | New per response | Persistent | Persistent |
| Latency | High | Low | Low | Lowest |
| Overhead | High | Medium | Low | Lowest |
| Protocol | HTTP | HTTP | HTTP | WS (after upgrade) |
| Binary support | ✅ (body) | ✅ | ❌ (text only) | ✅ |
| Auto-reconnect | Manual | Manual | ✅ Built-in | Manual |
| Firewall-friendly | ✅ | ✅ | ✅ | Sometimes ❌ |
| Complexity | Simple | Simple | Simple | Complex |

---

## WebSocket vs HTTP/2 Push vs SSE

| | WebSocket | HTTP/2 Server Push | SSE |
|--|-----------|-------------------|-----|
| Direction | Bidirectional | Server → Client | Server → Client |
| Use case | Chat, gaming | Preloading resources | Notifications, feeds |
| Protocol | WS (after HTTP upgrade) | HTTP/2 | HTTP/1.1 or 2 |
| Multiplexed | One WS per connection | Multiple streams | One stream |

---

## Interview Questions

### Q1. What is the difference between WebSockets and HTTP?

> **Answer:**
> - HTTP is **request-response**: the client must initiate every exchange; the server can only respond.
> - WebSockets provide a **persistent, full-duplex channel**: after a one-time HTTP upgrade handshake, both client and server can send data at any time without waiting for a request.
>
> WebSockets are ideal when the server needs to **push data** to the client (chat, live updates). HTTP is better for standard REST operations.

---

### Q2. Explain the WebSocket handshake. How does a WebSocket connection start?

> **Answer:**  
> A WebSocket connection starts as a standard HTTP GET request with special headers:
> - `Upgrade: websocket`
> - `Connection: Upgrade`
> - `Sec-WebSocket-Key: <base64-random>`
>
> If the server supports WebSockets, it responds with `101 Switching Protocols` and a `Sec-WebSocket-Accept` header (derived from the key). After this, the TCP connection stays open and both sides can exchange WebSocket frames directly.

---

### Q3. What is the difference between SSE and WebSockets? When would you choose one over the other?

> **Answer:**
> - **SSE:** Server → client only, text-based, built on HTTP (firewall-friendly), automatic reconnect.
> - **WebSockets:** Bidirectional, binary support, persistent WS connection.
>
> **Choose SSE:** For news feeds, notifications, live scores — anything where the client only reads updates. Simpler to implement and scales better (regular HTTP infrastructure).
>
> **Choose WebSockets:** For chat, collaborative editing, multiplayer games — where the client also sends data in real time.

---

### Q4. How does long polling work and what are its drawbacks?

> **Answer:**  
> In long polling, the client sends an HTTP request and the server **holds the connection open** until it has data to return. When data arrives, the server responds and the client immediately sends a new request.
>
> **Drawbacks:**
> - High resource usage on the server — each waiting client consumes a connection/thread.
> - Not efficient for many concurrent clients.
> - Latency spike when reconnecting.
> - HTTP overhead on every message exchange.

---

### Q5. How do you scale WebSocket connections in a distributed system?

> **Answer:**  
> WebSockets maintain stateful connections — a specific client is connected to a specific server. This breaks standard stateless horizontal scaling.
>
> **Solutions:**
> - **Sticky sessions at LB:** Route a client always to the same server (but limits scaling).
> - **Pub/Sub broker (Redis Pub/Sub, Kafka):** Each server subscribes to a channel. When server A needs to message a client on server B, it publishes to the channel; server B delivers to the client.
> - **Dedicated WS service:** Separate WebSocket servers (e.g., using Socket.io with Redis adapter), freeing API servers for stateless work.

---

<div align="center">
  <sub>← Back to <a href="networking_interview_topics.md">All Topics</a></sub>
</div>
