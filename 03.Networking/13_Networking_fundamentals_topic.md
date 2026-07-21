# Networking Fundamentals — Interview Study Guide

A complete reference covering HTTP, TCP/IP, DNS, TLS, load balancing, caching, WebSockets, gRPC, microservices networking, security, and Go-specific networking — with practical explanations for interview prep.

---

## 1. HTTP Fundamentals ⭐⭐⭐⭐⭐

### What happens when you type google.com into a browser?
1. **URL parsing** — browser identifies scheme, host, path.
2. **DNS resolution** — hostname is resolved to an IP address (browser cache → OS cache → resolver → root/TLD/authoritative servers).
3. **TCP connection** — a three-way handshake (SYN, SYN-ACK, ACK) is established with the server IP on port 443 (or 80).
4. **TLS handshake** — if HTTPS, client and server negotiate encryption keys and verify the server's certificate.
5. **HTTP request sent** — browser sends an HTTP request (method, headers, body).
6. **Server processes request** — may hit a load balancer → reverse proxy → application server → database.
7. **HTTP response returned** — status code, headers, body (usually HTML).
8. **Browser renders the page** — parses HTML, fetches additional resources (CSS/JS/images), each of which may repeat steps 2–7.

### Explain the HTTP request lifecycle
A client sends a **request** (method + URL + headers + optional body) to a server. The server routes it to the right handler, executes business logic, and returns a **response** (status code + headers + body). The connection may be reused (`Keep-Alive`) for further requests or closed.

### HTTP vs HTTPS
- **HTTP**: plaintext, port 80. Any intermediary can read/modify traffic.
- **HTTPS**: HTTP wrapped in **TLS**, port 443. Encrypts data in transit, verifies server identity via certificates, and ensures integrity (tampering is detectable).

### HTTP methods
- **GET** — retrieve a resource, no body, safe & idempotent.
- **POST** — create a resource or trigger an action, not idempotent.
- **PUT** — replace a resource entirely, idempotent.
- **PATCH** — partially update a resource, not guaranteed idempotent.
- **DELETE** — remove a resource, idempotent.
- **HEAD** — like GET but returns only headers.
- **OPTIONS** — asks what methods/headers are allowed (used in CORS preflight).

### PUT vs PATCH
`PUT` sends the **full replacement** representation of a resource — anything you omit is effectively removed. `PATCH` sends only the **fields that changed** (a partial update / diff).

### POST vs PUT
`POST` is used to **create** a new resource or trigger a non-idempotent action — calling it twice may create two resources. `PUT` is used to **create-or-replace** a resource at a known URL — calling it twice has the same effect as calling it once (idempotent).

### 200 vs 201 vs 204
- **200 OK** — generic success, response has a body.
- **201 Created** — a new resource was created (typically with a `Location` header pointing to it).
- **204 No Content** — success, but no body is returned (common for DELETE or some PUT/PATCH operations).

### 301 vs 302
- **301 Moved Permanently** — resource permanently moved; clients/search engines should update their references and cache the redirect.
- **302 Found (Temporary Redirect)** — resource temporarily at a different URL; original URL should still be used next time.

### 401 vs 403
- **401 Unauthorized** — you are **not authenticated** (missing/invalid credentials). "Who are you?"
- **403 Forbidden** — you **are** authenticated, but not **authorized** to access this resource. "I know who you are, but you can't do this."

### 404 vs 410
- **404 Not Found** — resource doesn't exist (may never have existed, or its absence isn't specified as permanent).
- **410 Gone** — resource existed before but has been **intentionally and permanently removed**. Signals to crawlers/clients to stop looking.

### 500 vs 502 vs 503 vs 504
- **500 Internal Server Error** — generic server-side failure (unhandled exception, bug).
- **502 Bad Gateway** — a proxy/load balancer got an **invalid response** from an upstream server.
- **503 Service Unavailable** — server is temporarily overloaded or down for maintenance (often with `Retry-After`).
- **504 Gateway Timeout** — a proxy/load balancer **didn't get a response in time** from the upstream server.

### HTTP Headers
Key-value metadata sent alongside a request/response that describe the message, control caching, negotiate content, carry auth tokens, etc. They don't affect the "meaning" of the body but tell client/server how to handle it.

### Commonly used headers
`Content-Type`, `Authorization`, `Accept`, `Cache-Control`, `User-Agent`, `Host`, `Cookie`/`Set-Cookie`, `ETag`, `Origin`, `Referer`, `Content-Length`.

### Content-Type
Tells the receiver what format the body is in, e.g. `application/json`, `text/html`, `multipart/form-data`. The server uses the request's `Content-Type` to parse the body; the response's `Content-Type` tells the client how to interpret what it received.

### Accept
Sent by the **client** to tell the server what response formats it can handle (e.g. `Accept: application/json`), enabling content negotiation.

### Authorization
Carries credentials for authenticating the request — e.g. `Authorization: Bearer <JWT>` or `Basic <base64>`.

### Host
Specifies which hostname the request is for (needed since one IP/server can host multiple domains — "virtual hosting"). Mandatory in HTTP/1.1.

### User-Agent
Identifies the client software making the request (browser, OS, bot, custom app) — used for analytics, feature detection, or blocking bots.

### Origin
Sent by browsers on cross-origin requests, indicating the scheme+host+port the request originated from. Used by servers to enforce **CORS**.

### Referer
Indicates the URL of the page that linked to the resource being requested — used for analytics, but can leak sensitive URLs (misspelling "Referer" is a famous historical typo in the HTTP spec).

### Cache-Control
Directs caching behavior for both requests and responses — e.g. `no-cache`, `no-store`, `max-age=3600`, `public`/`private`.

### ETag
An opaque identifier (usually a hash) representing a specific version of a resource. Clients send it back via `If-None-Match` to let the server respond `304 Not Modified` if nothing changed, saving bandwidth.

### If-Modified-Since
A conditional request header — client sends the last-known modification timestamp; server returns `304 Not Modified` if the resource hasn't changed since, otherwise the full resource.

### Cookies vs Sessions
A **cookie** is a small piece of data stored client-side and sent with every request to the same domain. A **session** is server-side state (often keyed by a session ID stored in a cookie) representing a logged-in user's data. Cookies are the *transport*; sessions are one common *use* of that transport.

### HttpOnly
A cookie flag that prevents JavaScript (`document.cookie`) from accessing the cookie — mitigates cookie theft via XSS.

### Secure Cookie
A cookie flag that ensures the cookie is only sent over HTTPS connections, never plaintext HTTP.

### SameSite
A cookie flag controlling whether the cookie is sent on cross-site requests. `Strict` (never cross-site), `Lax` (sent on top-level navigation, default in modern browsers), `None` (always sent, requires `Secure`) — primarily a **CSRF** mitigation.

### How cookies work
Server sends `Set-Cookie: name=value; attributes` in a response. Browser stores it and automatically attaches `Cookie: name=value` on subsequent requests to matching domain/path, subject to expiry and `SameSite`/`Secure`/`HttpOnly` rules.

---

## 2. HTTP Versions ⭐⭐⭐⭐

### HTTP/1.1 vs HTTP/2
HTTP/1.1 is **text-based**, one request per TCP connection at a time (pipelining exists but is rarely used due to head-of-line blocking). HTTP/2 is **binary framed**, supports **multiplexing** many requests over a single TCP connection, header compression (HPACK), and server push.

### Why is HTTP/2 faster?
Because it multiplexes multiple streams over one connection (avoiding the overhead of opening many TCP connections), compresses headers, and prioritizes streams — reducing latency and connection overhead versus HTTP/1.1's one-request-per-connection model.

### Multiplexing
The ability to send multiple independent request/response streams **concurrently over a single connection**, interleaved as frames, rather than serializing them.

### Head-of-line blocking
When one slow/stuck request blocks others queued behind it on the same connection. HTTP/1.1 suffers this at the **application layer** (one connection = one request at a time without pipelining). HTTP/2 fixes this at the application layer via multiplexing, but still suffers head-of-line blocking at the **TCP layer** (a single lost packet stalls all streams, since TCP guarantees in-order delivery).

### HTTP/3
The next HTTP version, built on **QUIC** instead of TCP. Because QUIC runs over UDP and multiplexes streams independently, HTTP/3 eliminates TCP-level head-of-line blocking entirely.

### Why does HTTP/3 use QUIC?
QUIC integrates transport and TLS, supports independent stream multiplexing (a lost packet only stalls its own stream, not others), has faster connection establishment (0-RTT/1-RTT), and better handles connection migration (e.g. switching Wi-Fi to cellular without dropping the connection).

### TCP vs QUIC
TCP is a connection-oriented, in-order, reliable stream protocol over IP with separate TLS negotiation. QUIC is built on UDP, has reliability and multiplexing built in at the transport layer, integrates encryption (always encrypted), and avoids cross-stream head-of-line blocking.

---

## 3. TCP/IP ⭐⭐⭐⭐⭐

### TCP
**Transmission Control Protocol** — connection-oriented, reliable, ordered, byte-stream protocol. Guarantees delivery via acknowledgments and retransmission, and provides flow/congestion control.

### UDP
**User Datagram Protocol** — connectionless, unreliable, unordered datagram protocol. No handshake, no guaranteed delivery, minimal overhead — used where speed matters more than reliability (DNS, video streaming, gaming, QUIC).

### TCP vs UDP
| | TCP | UDP |
|---|---|---|
| Connection | Yes (handshake) | No |
| Reliability | Guaranteed, retransmits | Best-effort |
| Ordering | In-order | Not guaranteed |
| Overhead | Higher | Lower |
| Use cases | Web, APIs, file transfer | DNS, streaming, gaming, VoIP |

### Three-way handshake
TCP connection setup: client sends **SYN** → server responds **SYN-ACK** → client responds **ACK**. After this, both sides have agreed on initial sequence numbers and the connection is established.

### Four-way connection termination
To close a TCP connection: initiator sends **FIN** → other side **ACK**s it → other side sends its own **FIN** → initiator **ACK**s that. Each side closes its half of the connection independently.

### Why is SYN different from ACK?
**SYN** (synchronize) initiates a connection and proposes an initial sequence number. **ACK** (acknowledge) confirms receipt of data/segments up to a certain sequence number. They're different flags serving different purposes, though a segment can carry both simultaneously (SYN-ACK).

### FIN
A TCP flag signaling "I have no more data to send" — begins graceful connection termination for that direction.

### RST
A TCP flag that **abruptly** resets/aborts a connection (e.g. connecting to a closed port, or an application-level error) — no graceful close, buffers may be dropped.

### TIME_WAIT
A state the connection-closing side holds for a period (often 2×MSL) after sending the final ACK, to ensure any delayed duplicate packets from the old connection don't corrupt a new one reusing the same port pair.

### CLOSE_WAIT
The state when the local side has received a FIN from the peer but hasn't yet closed its own side — indicates the application hasn't called `close()` yet, often a sign of a resource leak if it accumulates.

### Why is TCP reliable?
Every byte is sequenced and acknowledged; unacknowledged data is retransmitted after a timeout; a checksum detects corruption; and receivers reorder out-of-sequence segments before delivering to the application.

### Sliding window
A flow-control mechanism where the sender can have multiple unacknowledged segments "in flight" up to the receiver's advertised window size, rather than waiting for each ACK before sending the next segment — improves throughput.

### Congestion control
TCP's mechanism (e.g. slow start, congestion avoidance, AIMD) for adjusting how much data it sends based on **network** conditions (packet loss/delay as signals of congestion), to avoid overwhelming the network itself.

### Flow control
TCP's mechanism for the **receiver** to tell the sender how much buffer space it has (via the advertised window), preventing the sender from overwhelming a slow receiver — distinct from congestion control, which protects the network.

### MTU
**Maximum Transmission Unit** — the largest packet size (in bytes) that can be sent over a given link layer without fragmentation (commonly 1500 bytes for Ethernet).

### MSS
**Maximum Segment Size** — the largest amount of TCP payload data (excluding headers) that can fit in a single segment, derived from the MTU minus IP/TCP header overhead.

### What happens when a packet is lost?
For TCP: the sender detects loss via a retransmission timeout or duplicate ACKs (fast retransmit), reduces its congestion window, and retransmits the missing segment; the receiver buffers out-of-order segments until the gap is filled. For UDP: nothing — the application must handle loss itself, if it cares.

### Retransmission
Resending a segment that wasn't acknowledged within an expected time, or was inferred lost (e.g. via 3 duplicate ACKs triggering fast retransmit) — a core part of TCP's reliability guarantee.

---

## 4. DNS ⭐⭐⭐⭐⭐

### What is DNS?
The **Domain Name System** — a distributed, hierarchical naming system that translates human-readable domain names (e.g. `fuad71.me`) into IP addresses.

### DNS resolution process
1. Browser/OS checks local cache.
2. If not cached, query goes to a **recursive resolver** (e.g. ISP's or `8.8.8.8`).
3. Resolver queries a **root server** → gets pointed to the right **TLD server** (e.g. `.me`).
4. TLD server points to the domain's **authoritative name server**.
5. Authoritative server returns the actual IP (A/AAAA record).
6. Resolver caches and returns the answer to the client.

### Recursive vs Iterative query
In a **recursive** query, the resolver takes full responsibility for getting the final answer, querying other servers on the client's behalf. In an **iterative** query, each server responds with either the answer or a referral to the next server to ask — the resolver itself does the walking.

### TTL
**Time To Live** — how long (in seconds) a DNS record can be cached before it must be re-queried from the authoritative source.

### DNS caching
Resolvers, OSes, and browsers cache DNS answers for the record's TTL to reduce latency and load on authoritative servers.

### A Record
Maps a hostname to an **IPv4** address.

### AAAA Record
Maps a hostname to an **IPv6** address.

### MX Record
Specifies the **mail servers** responsible for receiving email for a domain, with priority values.

### TXT Record
Holds arbitrary text data for a domain — commonly used for domain verification, SPF/DKIM/DMARC email authentication.

### CNAME
An **alias** record pointing one hostname to another hostname (which is then resolved further) — cannot coexist with other record types on the same name.

### NS Record
Specifies which **name servers** are authoritative for a domain/zone.

### PTR Record
Used for **reverse DNS** lookup — maps an IP address back to a hostname.

### What happens if DNS fails?
The client can't resolve the hostname to an IP, so the connection never even reaches the TCP handshake stage — the browser shows an error like `DNS_PROBE_FINISHED_NXDOMAIN`. Applications often mitigate this with retries, fallback resolvers, or cached/last-known-good IPs.


---

## 5. TLS / SSL ⭐⭐⭐⭐⭐

### Explain HTTPS
HTTPS is HTTP layered on top of **TLS**. It encrypts data in transit (confidentiality), verifies the server's identity via a certificate (authentication), and ensures data isn't tampered with (integrity).

### How TLS handshake works (TLS 1.2 simplified)
1. **ClientHello** — client sends supported TLS versions/cipher suites and a random value.
2. **ServerHello** — server picks a cipher suite, sends its **certificate** and a random value.
3. Client verifies the certificate against a trusted **Certificate Authority**.
4. Client and server derive a shared **session key** (via asymmetric key exchange, e.g. Diffie-Hellman).
5. Both sides switch to **symmetric encryption** using the derived session key for the rest of the connection.
(TLS 1.3 streamlines this into fewer round trips, often 1-RTT or 0-RTT for resumed sessions.)

### Symmetric vs Asymmetric encryption
**Symmetric** encryption uses the **same key** to encrypt and decrypt — fast, but the key must be securely shared beforehand. **Asymmetric** encryption uses a **public/private key pair** — anyone can encrypt with the public key, but only the private key holder can decrypt — slower, but solves the key-distribution problem.

### Why use both?
Asymmetric crypto is used briefly during the handshake to securely exchange/derive a symmetric session key (solving the key-distribution problem), then **symmetric** encryption — much faster — is used for the actual bulk data transfer for the rest of the session.

### Certificate Authority (CA)
A trusted third party that issues and digitally signs certificates, vouching that a public key belongs to the entity named in the certificate (e.g. `fuad71.me`). Browsers/OSes ship with a list of trusted root CAs.

### Certificate chain
The path from a website's certificate up through one or more **intermediate CA certificates** to a trusted **root CA certificate**. Each certificate in the chain is signed by the one above it, letting the client verify trust without directly trusting every individual site.

### Public key
The half of an asymmetric key pair that can be shared openly — used by others to encrypt data for you, or to verify your digital signatures.

### Private key
The secret half of the key pair, kept only by its owner — used to decrypt data encrypted with the matching public key, or to create digital signatures.

### Perfect Forward Secrecy (PFS)
A property where session keys are derived fresh (e.g. via ephemeral Diffie-Hellman) for each session and never stored — so if a server's long-term private key is later compromised, past recorded sessions still can't be decrypted.

### SNI (Server Name Indication)
A TLS extension where the client tells the server which hostname it's trying to connect to **during the handshake itself**, before the server sends a certificate — this lets one IP/server host multiple HTTPS domains, each with its own certificate.

### ALPN (Application-Layer Protocol Negotiation)
A TLS extension that lets client and server agree on the application protocol (e.g. HTTP/1.1 vs HTTP/2) as part of the TLS handshake, avoiding an extra round trip to negotiate it afterward.

---

## 6. Load Balancer ⭐⭐⭐⭐⭐

### What is a load balancer?
A component that distributes incoming traffic across multiple backend servers to improve availability, scalability, and fault tolerance — if one server fails, traffic is routed to healthy ones.

### Layer 4 vs Layer 7 load balancer
**Layer 4** (transport layer) balances based on IP/port info without inspecting the actual content — fast, protocol-agnostic. **Layer 7** (application layer) understands HTTP itself — can route based on URL path, headers, cookies, enabling smarter routing (e.g. `/api/*` to one service, `/static/*` to another) at the cost of more overhead.

### Round Robin
Distributes requests to backend servers in sequential, cyclic order — simple, assumes all servers have equal capacity.

### Least Connection
Routes each new request to whichever backend currently has the **fewest active connections** — better than round robin when requests have variable duration.

### Weighted Round Robin
Like round robin, but servers are assigned weights (e.g. based on capacity), so more powerful servers receive proportionally more traffic.

### Sticky Sessions
A technique (via cookie or client IP) that ensures a given client's requests always go to the **same backend server** — useful for stateful apps that store session data locally, though it works against even load distribution.

### Reverse Proxy vs Load Balancer
A **reverse proxy** sits in front of one or more servers, forwarding client requests and often adding features (caching, TLS termination, compression). A **load balancer** specifically focuses on distributing traffic across *multiple* backend instances. In practice, a reverse proxy (like Nginx) is often *used as* a load balancer — the terms overlap heavily.

### Nginx vs HAProxy
**Nginx** is a general-purpose web server / reverse proxy / load balancer, also great at serving static files and handling TLS termination. **HAProxy** is purpose-built specifically for high-performance load balancing (L4 and L7), often favored in pure load-balancing/proxying scenarios for its performance and fine-grained health-check/routing features.

### Why use a reverse proxy?
Centralizes TLS termination, load balancing, caching, compression, rate limiting, and request logging in front of application servers — decoupling these cross-cutting concerns from application code (e.g. your Go/Gin app doesn't need to handle TLS certs itself).

---

## 7. Reverse Proxy ⭐⭐⭐⭐

### What is a reverse proxy?
A server that sits between clients and backend servers, forwarding client requests to the appropriate backend and returning the backend's response to the client — the backend's existence is hidden from the client.

### Reverse proxy vs forward proxy
A **forward proxy** sits in front of *clients*, forwarding their requests to the internet on their behalf (hiding the client from the destination server — e.g. corporate proxies, VPNs). A **reverse proxy** sits in front of *servers*, forwarding client requests to the right backend (hiding backend servers from the client).

### How Nginx works
Nginx uses an event-driven, asynchronous architecture (a small number of worker processes, each handling many connections via non-blocking I/O) rather than one thread/process per connection — this lets it handle very high concurrency with low resource usage. It can terminate TLS, serve static files directly, and proxy dynamic requests to upstream application servers.

### API Gateway vs Reverse Proxy
A reverse proxy mainly does traffic forwarding, load balancing, and TLS termination. An **API Gateway** is a more feature-rich layer built for APIs specifically — adding authentication, rate limiting, request/response transformation, aggregation of multiple backend calls, and API versioning/routing logic. Every API gateway is essentially a specialized reverse proxy, but not every reverse proxy is a full gateway.

### Why put Nginx in front of Go?
Go's `net/http` can serve HTTP directly, but Nginx in front offers: TLS termination without touching Go code, serving static assets efficiently, buffering slow clients (protecting Go's goroutines from slow-client attacks), easy horizontal scaling/load balancing across multiple Go instances, and centralized logging/rate limiting/caching.

---

## 8. OSI Model ⭐⭐⭐⭐

### Explain the OSI model
A 7-layer conceptual framework describing how network communication happens:
1. **Physical** — raw bits over a medium (cables, radio).
2. **Data Link** — frames between directly connected nodes (Ethernet, MAC addresses).
3. **Network** — routing packets across networks (IP).
4. **Transport** — end-to-end delivery, reliability (TCP, UDP).
5. **Session** — managing sessions/connections between applications.
6. **Presentation** — data formatting/encryption/translation (TLS often mapped here).
7. **Application** — the protocols apps actually use (HTTP, DNS, FTP).

### Which layer does TCP belong to?
**Layer 4 — Transport.**

### Which layer is HTTP?
**Layer 7 — Application.**

### Which layer is DNS?
**Layer 7 — Application** (though it relies on transport-layer UDP/TCP underneath).

### Which layer is TLS?
Commonly mapped to **Layer 6 — Presentation** (it handles encryption/data formatting), though in practice it sits between the transport and application layers and is sometimes discussed as its own thing.

### Which layer is Ethernet?
**Layer 2 — Data Link** (with some physical signaling aspects at Layer 1).

---

## 9. IP Addressing ⭐⭐⭐⭐

### IPv4 vs IPv6
**IPv4** uses 32-bit addresses (~4.3 billion addresses, e.g. `192.168.1.1`), now largely exhausted and reliant on NAT. **IPv6** uses 128-bit addresses (a vastly larger space, e.g. `2001:db8::1`), designed to eliminate the need for NAT and simplify routing, with built-in support for things like auto-configuration.

### Private IP ranges
Reserved ranges not routable on the public internet, used within private networks: `10.0.0.0/8`, `172.16.0.0/12`, `192.168.0.0/16`.

### Public IP
An IP address that is globally unique and routable on the public internet, assigned by an ISP or cloud provider.

### NAT (Network Address Translation)
A technique that lets many devices on a private network share a single public IP address, translating private IPs to the public IP (and back) at the network boundary — conserves public IPv4 addresses.

### CIDR (Classless Inter-Domain Routing)
A notation for specifying IP ranges, e.g. `192.168.1.0/24`, where `/24` means the first 24 bits are the network portion (256 addresses in that block). Replaced the old rigid "class A/B/C" system with flexible-size blocks.

### What is subnetting?
Dividing a larger network into smaller sub-networks ("subnets") by borrowing bits from the host portion of an address — used to organize networks, limit broadcast domains, and improve security/routing efficiency.

### What is a gateway?
A device (usually a router) that connects one network to another, forwarding traffic destined outside the local network.

### Default gateway
The gateway a device sends traffic to by default when the destination isn't on its local network — typically the router.

### Loopback address
An address (`127.0.0.1` for IPv4, `::1` for IPv6) that always refers back to the local machine itself, used for local testing without touching physical network hardware.

### Localhost
The hostname that resolves to the loopback address — refers to "this machine" regardless of its actual network configuration.

### 127.0.0.1 vs 0.0.0.0
`127.0.0.1` refers specifically to the loopback interface (only reachable from the same machine). `0.0.0.0` as a **bind address** means "listen on all available network interfaces" (not just loopback) — as a destination, it's a non-routable placeholder meaning "no particular address."

---

## 10. Ports ⭐⭐⭐⭐⭐

### What is a port?
A 16-bit number (0–65535) that identifies a specific process/service on a machine, allowing multiple network applications to share a single IP address.

### Well-known ports
Ports 0–1023, reserved for standard services, e.g. `80` (HTTP), `443` (HTTPS), `22` (SSH), `53` (DNS), `25` (SMTP).

### Dynamic ports
Ports 49152–65535 (the "dynamic"/"private" range), typically used as ephemeral source ports for outgoing client connections.

### Ephemeral ports
Temporary ports assigned by the OS for the client side of a connection, released once the connection closes — lets a client make many simultaneous outgoing connections.

### Why can only one process listen on a port?
Because a listening socket is bound to a specific (IP, port, protocol) tuple exclusively — the OS uses this tuple to route incoming packets to the correct process, so allowing two listeners on the same tuple would create ambiguity about who should receive the traffic.

### Difference between listening and established connections
A **listening** socket is waiting to accept new incoming connections on a port (server-side, no specific remote peer yet). An **established** connection is an active, ongoing conversation between a specific local and remote (IP, port) pair, post-handshake.

---

## 11. Socket Programming ⭐⭐⭐⭐

### What is a socket?
An OS-level abstraction representing one endpoint of a network connection, identified by an (IP address, port, protocol) tuple — the interface applications use to send/receive data over a network.

### TCP socket
A socket using the TCP protocol — connection-oriented; requires a handshake before data can be exchanged, then provides a reliable, ordered byte stream.

### UDP socket
A socket using the UDP protocol — connectionless; you can send datagrams immediately without any handshake, with no delivery/order guarantees.

### How does net.Listen() work in Go?
`net.Listen("tcp", ":8080")` creates a listening socket bound to the given address/port, puts it into the OS's listen queue (backlog), and returns a `Listener`. Calling `.Accept()` on it blocks until a client connects, then returns a `net.Conn` representing that individual connection.

### How does net.Conn work?
`net.Conn` is Go's interface representing an established connection (TCP, UDP, etc.), exposing `Read`, `Write`, `Close`, and deadline-setting methods — abstracting over the underlying transport so code can treat different connection types uniformly.

### How many clients can connect simultaneously?
Practically bounded by system resources (file descriptor limits, memory, CPU) rather than a hard protocol limit — Go's runtime uses goroutines + the netpoller (epoll/kqueue under the hood) so it can efficiently handle tens of thousands of concurrent connections on modest hardware, since goroutines are cheap and I/O is non-blocking under the hood.


---

## 12. REST API ⭐⭐⭐⭐⭐

### REST principles
Representational State Transfer — an architectural style built on: **statelessness** (each request carries all context needed), a **uniform interface** (resources identified by URLs, manipulated via standard HTTP methods), **client-server separation**, **cacheability**, and a **layered system** (proxies/gateways can sit transparently between client and server).

### Idempotent methods
Methods where making the same request multiple times has the same effect as making it once: `GET`, `PUT`, `DELETE`, `HEAD`, `OPTIONS`. `POST` is generally **not** idempotent.

### Safe methods
Methods that don't modify server state: `GET`, `HEAD`, `OPTIONS`. Safe methods are also idempotent (but not all idempotent methods are safe — e.g. `DELETE` changes state but is idempotent).

### Stateless APIs
Each request contains all the information the server needs to process it (e.g. an auth token) — the server doesn't rely on stored session state from prior requests, which makes horizontal scaling and load balancing much simpler.

### Versioning APIs
Common strategies: URL path (`/v1/users`), a custom header (`Api-Version: 1`), query parameter (`?version=1`), or content negotiation via `Accept` header (`application/vnd.api.v1+json`) — allows evolving an API without breaking existing clients.

### Pagination
Splitting large result sets into pages to avoid huge payloads — common approaches: **offset-based** (`?page=2&limit=20`, simple but can skip/duplicate items if data changes) and **cursor-based** (`?after=<id>`, more stable under concurrent writes, common for large/real-time datasets).

### Filtering
Letting clients narrow results via query parameters, e.g. `GET /matches?sport=football&status=live` — reduces payload size and lets the server do the filtering efficiently (often via indexed DB queries) instead of the client.

### HATEOAS (Hypermedia as the Engine of Application State)
A REST constraint where responses include links to related actions/resources (like a web page with hyperlinks), so clients can discover available operations dynamically rather than hardcoding URLs — rarely fully implemented in practice, but a core original REST idea.

### Why use JSON?
Human-readable, lightweight, natively supported by JavaScript, has broad tooling/library support across virtually every language, and is simpler to work with than XML for typical API payloads.

### JSON vs Protobuf
**JSON** is text-based, human-readable, larger payload size, no strict schema by default (flexible but error-prone). **Protobuf** is binary, requires a predefined schema (`.proto` files), much smaller/faster to serialize/deserialize, and enforces strong typing — commonly used in gRPC and performance-sensitive internal service-to-service communication.

---

## 13. Caching ⭐⭐⭐⭐⭐

### Browser cache
Stores static assets (JS, CSS, images) locally in the browser per `Cache-Control`/`Expires` headers, avoiding repeat network requests entirely for unchanged resources.

### CDN cache
Caches content at geographically distributed **edge servers** close to users, reducing latency and origin server load for frequently requested (especially static) content.

### Redis cache
An in-memory key-value store commonly used as an **application-level cache** — e.g. caching expensive database query results or computed data, with configurable expiry (`TTL`), to reduce database load and latency.

### HTTP cache
Caching governed by HTTP headers (`Cache-Control`, `ETag`, `Expires`) at any point along the request path — browser, CDN, or intermediate proxy — controlling how long and under what conditions a response can be reused.

### Cache-Control
The primary HTTP header controlling caching behavior — directives like `no-store` (never cache), `no-cache` (must revalidate before use), `max-age=N` (fresh for N seconds), `private`/`public` (who can cache it).

### ETag
(See section 1) — a version fingerprint of a resource used for cache validation via conditional requests (`If-None-Match`), allowing `304 Not Modified` responses instead of resending full content.

### Cache invalidation
The process of removing or updating stale cached data when the underlying data changes — famously one of the "two hard things in computer science." Strategies include TTL expiry, explicit invalidation on write, and versioned cache keys.

### Cache stampede
When a popular cached item expires and a flood of simultaneous requests all miss the cache at once, all hammering the origin/database simultaneously to regenerate it — mitigated via locks, request coalescing, or staggered/probabilistic early expiry.

### Cache penetration
Repeated queries for data that **doesn't exist** in the cache *or* the underlying store, bypassing the cache entirely every time and hitting the database — mitigated by caching "negative results" (null markers) or using a bloom filter to short-circuit known-nonexistent keys.

### Cache avalanche
When a large number of cache entries **expire at the same time** (e.g. all set with the same TTL), causing a sudden surge of requests to the backend — mitigated by randomizing/jittering TTLs.

### Write-through cache
Writes go to the cache **and** the underlying database synchronously, together, keeping them consistent — simpler consistency model, but write latency includes both operations.

### Write-back cache
Writes go to the cache first and are asynchronously flushed to the database later — faster writes, but risks data loss if the cache fails before flushing.

### Read-through cache
Reads always go through the cache layer; on a cache miss, the cache itself (not the application) fetches from the database, stores it, and returns it — simplifies application code, since it doesn't need explicit "check cache, then check DB" logic.

---

## 14. CDN

### What is CDN?
**Content Delivery Network** — a geographically distributed network of proxy/edge servers that cache and serve content from locations physically close to end users, reducing latency and offloading traffic from the origin server.

### Why CDN?
Lower latency (content served from a nearby edge), reduced load on origin servers, better resilience against traffic spikes/DDoS (absorbed at the edge), and improved availability through redundancy across many locations.

### Edge server
A CDN server located at a point close to end users (an "edge" of the network) that caches and serves content on the origin's behalf.

### Cache hit
A request that the CDN/cache can satisfy directly from its stored copy, without contacting the origin server.

### Cache miss
A request the cache doesn't have (or has expired), requiring it to fetch fresh content from the origin server before responding — and typically caching it for next time.

### Cache invalidation
(See section 13) — for CDNs specifically, this often means "purging" content from all edge nodes so they re-fetch fresh content from origin on the next request.

---

## 15. WebSocket ⭐⭐⭐⭐⭐

### What is WebSocket?
A protocol providing a **persistent, full-duplex** connection between client and server over a single TCP connection — after an initial HTTP handshake ("upgrade"), both sides can send messages to each other at any time without the overhead of repeated HTTP requests.

### HTTP vs WebSocket
HTTP is **request-response**: the client always initiates, and the connection is typically short-lived per exchange. WebSocket is **bidirectional and persistent**: either side can push data at any time over a long-lived connection, ideal for real-time features (live scores, chat, notifications).

### WebSocket handshake
The client sends a normal HTTP GET request with an `Upgrade: websocket` and `Connection: Upgrade` header (plus a `Sec-WebSocket-Key`). If the server supports it, it responds with `101 Switching Protocols`, and the same TCP connection is repurposed for the WebSocket protocol from then on.

### Long Polling
A technique simulating real-time updates over plain HTTP: the client sends a request, the server **holds it open** until new data is available (or a timeout), then responds — the client immediately re-issues another request. Higher overhead than WebSocket but works over vanilla HTTP infrastructure.

### SSE (Server-Sent Events)
A one-way (server-to-client) streaming protocol over plain HTTP, where the server keeps a connection open and pushes text-based events to the client as they occur. Simpler than WebSocket for cases where only the server needs to push data.

### Heartbeat
Periodic small messages sent over a long-lived connection to confirm it's still alive and to keep intermediate proxies/firewalls from timing it out due to inactivity.

### Ping/Pong
A specific heartbeat mechanism built into the WebSocket protocol — one side sends a `Ping` control frame, the other must respond with a `Pong`, confirming the connection is still healthy.

### When should you use WebSocket?
When you need **low-latency, bidirectional, frequent** real-time communication — live match score updates, chat, collaborative editing, live dashboards — rather than occasional client-initiated requests, where plain HTTP (or polling) suffices.

---

## 16. gRPC ⭐⭐⭐⭐

### What is gRPC?
A high-performance RPC (Remote Procedure Call) framework from Google, built on **HTTP/2** and **Protocol Buffers**, letting you call methods on a remote service as if they were local function calls, with strongly typed request/response schemas.

### HTTP/2 dependency
gRPC relies on HTTP/2 for multiplexed streams over a single connection, binary framing, and header compression — this is what enables efficient bidirectional/streaming RPCs without the overhead HTTP/1.1 would impose.

### Protocol Buffers (Protobuf)
A binary serialization format with a schema defined in `.proto` files — compact, fast to (de)serialize, and generates strongly typed client/server code for many languages, unlike loosely-typed JSON.

### Unary RPC
A single request, single response — the classic RPC call pattern, analogous to a normal function call or a typical REST request.

### Streaming RPC
Either the client or the server sends a **stream** of multiple messages over a single call instead of just one — e.g. server-streaming (one request, many responses, like live updates) or client-streaming (many requests, one response).

### Bidirectional streaming
Both client and server send independent streams of messages over the same connection simultaneously, useful for scenarios like real-time chat or continuous data exchange where either side needs to push data at any time.

### gRPC vs REST
gRPC uses binary Protobuf (compact, fast, strongly typed, generated client code) and HTTP/2 (multiplexing, streaming) — better suited to internal service-to-service communication with performance needs. REST typically uses JSON over HTTP/1.1, is more human-readable/debuggable, broadly supported by browsers directly, and has a lower learning curve — often preferred for public-facing APIs.

### Why is gRPC faster?
Smaller binary payloads (vs verbose JSON text), HTTP/2 multiplexing avoids connection overhead for many concurrent calls, and Protobuf (de)serialization is significantly faster than JSON parsing.

---

## 17. Microservices Networking ⭐⭐⭐⭐⭐

### Service Discovery
The mechanism by which services find the network location (IP/port) of other services at runtime, since instances in microservices architectures are often dynamic (scaled up/down, rescheduled) — implemented via a registry (e.g. Consul, etcd) or DNS-based discovery (common in Kubernetes).

### API Gateway
(See section 7) — a single entry point for external clients into a microservices system, handling routing, auth, rate limiting, and aggregation across multiple backend services.

### Service Mesh
Infrastructure layer (e.g. Istio, Linkerd) that handles service-to-service communication concerns — routing, retries, load balancing, mTLS, observability — typically via lightweight proxies ("sidecars") deployed alongside each service instance, decoupling this logic from application code.

### Circuit Breaker
A pattern that stops sending requests to a failing/unresponsive downstream service after a threshold of failures, "opening the circuit" and failing fast (or falling back) instead of piling up requests that will likely fail too — protects both the failing service and the caller from cascading failures.

### Retry
Automatically re-attempting a failed request (often with backoff/jitter), useful for transient failures — but must be applied carefully (with idempotency and limits) to avoid amplifying load on an already struggling service.

### Timeout
A maximum time a caller waits for a response before giving up — essential to prevent one slow downstream service from holding resources indefinitely and cascading slowness upstream.

### Rate Limiting
Restricting how many requests a client (or the system overall) can make in a given time window — protects services from being overwhelmed and enforces fair usage/quotas.

### Bulkhead
A resilience pattern (named after ship compartments) that isolates resources (e.g. thread pools, connection pools) per dependency, so if one dependency degrades, it can't exhaust shared resources and take down unrelated parts of the system.

### Health Check
An endpoint or mechanism (e.g. `/healthz`) that reports whether a service instance is alive and ready to serve traffic — used by load balancers/orchestrators to route traffic only to healthy instances.

### Distributed Tracing
A technique for tracking a single request as it flows across multiple microservices, correlating logs/spans via a shared trace ID — essential for debugging latency and failures in complex service chains (tools: Jaeger, Zipkin, OpenTelemetry).

---

## 18. Security ⭐⭐⭐⭐⭐

### CORS (Cross-Origin Resource Sharing)
A browser-enforced security mechanism that controls whether JavaScript running on one origin (domain/scheme/port) can make requests to a different origin. Servers opt in via `Access-Control-Allow-Origin` and related headers; for non-simple requests, the browser first sends an `OPTIONS` preflight request to check permissions.

### CSRF (Cross-Site Request Forgery)
An attack where a malicious site tricks a logged-in user's browser into making an unwanted request to another site where they're authenticated (since cookies are sent automatically) — mitigated via CSRF tokens, `SameSite` cookies, and checking `Origin`/`Referer`.

### XSS (Cross-Site Scripting)
An attack where malicious JavaScript is injected into a page (via unescaped user input) and executes in other users' browsers — can steal cookies/tokens or perform actions as the victim. Mitigated via output encoding/escaping, Content Security Policy (CSP), and `HttpOnly` cookies.

### SQL Injection
An attack where untrusted input is concatenated directly into a SQL query, letting an attacker manipulate the query's logic (e.g. bypass authentication, exfiltrate data). Mitigated by always using **parameterized queries/prepared statements**, never string-concatenating user input into SQL.

### JWT (JSON Web Token)
A compact, self-contained, digitally signed token (header.payload.signature, base64-encoded) commonly used for stateless authentication — the server can verify the signature without a database lookup, though this also means a compromised signing key or a token that isn't properly revoked can be a risk (JWTs aren't inherently revocable without extra infrastructure).

### OAuth2
An authorization **framework** (not authentication itself) that lets a user grant a third-party application limited access to their resources on another service, without sharing their password — via access tokens issued through defined "grant" flows (authorization code, client credentials, etc.).

### API Key
A simple, static secret string a client includes in requests to identify/authenticate itself to an API — simpler than OAuth2 but less flexible (no fine-grained scopes/expiry by default, and must be guarded carefully since it's often long-lived).

### Mutual TLS (mTLS)
An extension of TLS where **both** the client and server present certificates to authenticate each other (not just the server, as in normal TLS) — common for securing service-to-service communication in a service mesh or zero-trust architecture.

### Rate limiting
(See section 17) — from a security angle, also protects against brute-force login attempts, credential stuffing, and API abuse, not just overload protection.

### DDoS protection
Techniques and infrastructure to withstand **Distributed Denial of Service** attacks — e.g. traffic scrubbing at the edge/CDN, rate limiting, anomaly detection, over-provisioned capacity, and upstream provider-level mitigation (e.g. Cloudflare).


---

## 19. Production Troubleshooting ⭐⭐⭐⭐⭐

### "API is slow. Where do you start?"
A structured approach: (1) check metrics/dashboards for latency percentiles (p50/p95/p99) and error rates, (2) identify whether it's the app, database, network, or an external dependency using distributed tracing, (3) check resource usage (CPU, memory, connections, disk I/O), (4) check for slow database queries or lock contention, (5) check downstream/third-party API latency, (6) reproduce and profile if needed. Narrow down layer by layer rather than guessing.

### High latency
Investigate: is it consistent or intermittent? Client-side (DNS, TLS handshake) vs server-side (processing time) vs network (packet loss, routing)? Tools: `curl -w` timing breakdown, tracing, APM tools, checking each hop (proxy → app → DB).

### Packet loss
Check with `ping` (sustained loss over time) or `mtr`/`traceroute` to identify where in the path loss occurs. Can be caused by network congestion, faulty hardware, misconfigured MTU, or an overloaded server dropping packets at the OS level.

### High CPU
Investigate with profiling tools (`pprof` in Go), check for inefficient algorithms, excessive goroutine/thread creation, GC pressure, busy-loops, or a spike in traffic. Compare against baseline to distinguish "always high" (code issue) vs "spiking under load" (capacity issue).

### High memory
Check for memory leaks (growing over time without release), large in-memory caches without eviction, goroutine leaks holding references, or genuinely needing more memory for current load. Go: use `pprof` heap profiles to find allocation sources.

### Too many TIME_WAIT sockets
Usually indicates the server is closing a very high volume of short-lived connections (e.g. not reusing HTTP client connections, or a load balancer opening a new connection per request). Mitigations: enable connection reuse/keep-alive, tune `net.ipv4.tcp_tw_reuse`, or reduce unnecessary connection churn (e.g. connection pooling).

### Database timeout
Query took longer than the configured timeout — investigate slow queries (missing indexes, table locks, large scans), connection pool exhaustion (all connections busy), or the database itself being overloaded/under-provisioned.

### Redis timeout
Similar causes: slow/blocking commands (e.g. `KEYS *` on a large dataset), network issues between app and Redis, Redis running out of memory and evicting/thrashing, or too few connections in the pool relative to concurrent demand.

### 502 Bad Gateway
The proxy/load balancer got an invalid or no response from the upstream application server — often means the backend crashed, isn't running, or returned a malformed response.

### 504 Gateway Timeout
The proxy/load balancer didn't get **any** response from the upstream within its configured timeout — the backend is alive but too slow (or hung/deadlocked).

### SSL handshake failure
Can be caused by: expired/misconfigured certificate, certificate chain issues (missing intermediate cert), mismatched hostname (SNI) vs certificate, unsupported cipher suite/TLS version between client and server, or clock skew (certificate validity checks depend on system time).

### DNS timeout
The resolver didn't get a response in time — could be an overloaded/unreachable DNS server, network issues, or misconfigured resolver settings. Mitigate with retries, fallback resolvers, and caching.

---

## 20. Golang Networking ⭐⭐⭐⭐⭐

### net/http architecture
Go's `net/http` package provides both client (`http.Client`) and server (`http.Server`) implementations built on top of `net.Listener`/`net.Conn`. The server accepts connections and dispatches each request to a `Handler` (often via a `ServeMux` or a router like Gin's engine), typically running each request in its own goroutine.

### How does http.Server work?
`http.Server` wraps a `net.Listener`, `Accept()`-ing new connections in a loop, and spawns a goroutine per connection to read requests and invoke the configured `Handler`. It manages things like read/write timeouts, TLS, and keep-alive connection reuse.

### What is http.Transport?
The underlying implementation used by `http.Client` for making outgoing HTTP requests — it manages the connection pool, TLS configuration, proxy settings, and low-level details like keep-alive and idle connection reuse for a client making requests to servers.

### Connection pooling
Reusing already-established TCP (and TLS) connections for multiple requests instead of opening a new one each time — dramatically reduces latency and resource overhead, since handshakes are expensive. `http.Transport` manages this pool automatically.

### Keep-Alive
An HTTP feature (default in HTTP/1.1) where the underlying TCP connection stays open after a response, allowing subsequent requests to reuse it instead of establishing a new connection each time.

### Idle connections
Connections in the pool that are currently open but not actively handling a request — kept around (up to configured limits) so they can be reused for the next request without a fresh handshake.

### MaxIdleConns
An `http.Transport` setting controlling the maximum total number of idle (keep-alive) connections kept open across **all** hosts — once exceeded, older idle connections are closed.

### MaxConnsPerHost
An `http.Transport` setting limiting the total number of connections (idle + active) allowed to a single host — prevents one destination from monopolizing all your outgoing connections.

### Request timeout
A deadline set on an HTTP request (via `context.WithTimeout` or `http.Client.Timeout`) after which it's aborted if not completed — essential for preventing slow/hung requests from consuming resources indefinitely.

### Context cancellation
Go's `context.Context` mechanism for propagating deadlines/cancellation signals through a call chain — e.g. if a client disconnects or a timeout fires, `ctx.Done()` fires, letting downstream operations (DB queries, external API calls) stop early instead of wasting work.

### HTTP Client reuse
Reusing a single `http.Client` instance (rather than creating new ones) is important because the client owns the connection pool (`Transport`) — reusing it lets connections actually be pooled and reused across requests.

### Why shouldn't you create a new http.Client for every request?
Because each new client gets its own fresh `Transport`/connection pool, meaning every request would need a brand-new TCP+TLS handshake — losing all the performance benefits of keep-alive/connection reuse, and potentially exhausting file descriptors/ports under load (accumulating unclosed idle connections across many discarded clients).

### How does Go scheduler handle thousands of connections?
Go uses lightweight **goroutines** (a few KB of stack each, growable) instead of OS threads, multiplexed onto a small number of OS threads by the Go runtime scheduler (M:N scheduling). Combined with the netpoller for non-blocking I/O, this lets a server handle tens of thousands of concurrent connections without the memory/context-switching overhead that one-OS-thread-per-connection would incur.

### How does the netpoller work?
Go's runtime integrates with OS-level async I/O facilities (`epoll` on Linux, `kqueue` on BSD/macOS, IOCP on Windows). When a goroutine performs a blocking-looking network call (e.g. `conn.Read()`), instead of blocking an OS thread, the runtime parks the goroutine and registers interest with the OS poller; when data becomes available, the poller wakes the goroutine back up on an available OS thread — giving blocking-style code non-blocking performance under the hood.

### Goroutines vs OS threads for networking
OS threads are expensive to create (MBs of stack, kernel-managed context switches) so one-thread-per-connection doesn't scale well past a few thousand connections. Goroutines are cheap, user-space-scheduled, and start with a tiny stack that grows as needed — combined with the netpoller, Go can handle vastly more concurrent connections per unit of memory/CPU than a thread-per-connection model.

---

## 21. Kubernetes / Cloud Networking (Senior)

### ClusterIP
The default Kubernetes Service type — exposes a service on an internal, virtual IP reachable only **within** the cluster, used for internal service-to-service communication.

### NodePort
Exposes a service on a static port on **every node's** IP in the cluster, making it reachable from outside the cluster via `<NodeIP>:<NodePort>` — a simpler but less flexible way to expose services externally than a LoadBalancer.

### LoadBalancer service
A Service type that provisions an actual external load balancer (via the cloud provider, e.g. AWS ELB/DigitalOcean LB) that routes external traffic into the cluster to the service.

### Ingress
A Kubernetes resource that manages external HTTP(S) access to services within the cluster, typically with host/path-based routing, TLS termination, and usually backed by a controller (e.g. Nginx Ingress Controller) — more flexible/economical than provisioning a separate LoadBalancer per service.

### DNS inside Kubernetes
Kubernetes runs an internal DNS service (typically CoreDNS) that gives every Service a resolvable DNS name (e.g. `my-service.my-namespace.svc.cluster.local`), letting pods discover and communicate with services by name instead of hardcoded IPs.

### Service discovery
In Kubernetes, achieved via DNS (as above) plus environment variables injected into pods — decouples pods from needing to know the actual (dynamic) IPs of other pods.

### NetworkPolicy
A Kubernetes resource that defines rules for what traffic is allowed to/from pods (like a firewall at the pod level) — e.g. restrict a database pod to only accept connections from specific application pods.

### Pod networking
Each pod in Kubernetes gets its own unique IP address (unlike Docker's default networking model), and all pods can communicate with all other pods across nodes without NAT — implemented under the hood by a CNI plugin.

### Overlay networks
A virtual network layer built on top of the underlying physical network, used to give pods a flat, cluster-wide addressable network regardless of which physical node they're on (encapsulating pod traffic to route it across nodes).

### CNI plugins
**Container Network Interface** plugins (e.g. Calico, Flannel, Cilium) implement the actual pod networking in Kubernetes — assigning IPs to pods, setting up routing/overlay networks, and (for some plugins) enforcing NetworkPolicies.

---

## 22. Practical Debugging Questions

### Why does curl work but the browser doesn't?
Common causes: **CORS** (browsers enforce it, curl doesn't), the browser is sending different headers/cookies (or missing an `Authorization` header curl had), a cached bad response/service worker in the browser, or a client-side JS error preventing the request from even being sent correctly.

### Why does Postman work but Go doesn't?
Common causes: differences in default headers (Postman might auto-set `Content-Type` or `Accept-Encoding` differently), TLS verification settings (Postman might skip cert verification by default), redirect handling differences, or the Go `http.Client`'s default timeout/transport settings behaving differently (e.g. not reusing connections properly, or missing an explicit header the Go code forgot to set).

### How would you debug a timeout?
Determine which stage timed out (DNS, TCP connect, TLS handshake, waiting for response, or reading the body) using timing instrumentation (`curl -w`, or Go's `httptrace`). Then narrow down: is it the client's timeout config, network latency, or the server being slow — check server-side logs/traces for how long it actually took to process.

### How do you identify packet loss?
Use `ping` over a sustained period to look for dropped packets, or `mtr`/`traceroute` to see where along the path loss occurs at each hop. `tcpdump`/Wireshark can confirm retransmissions at the packet level.

### Why does one API call take 20 seconds?
Investigate: is it a slow downstream call (DB query, external API) without a timeout configured? Lock contention? A cold-start/connection-establishment cost? Use distributed tracing to see exactly which span in the request's lifecycle consumed the time.

### How would you debug a TLS handshake failure?
Use `openssl s_client -connect host:443` to manually inspect the handshake, check certificate validity/expiry and chain completeness, verify the hostname matches the certificate (or SNI is being sent correctly), and check for TLS version/cipher suite mismatches between client and server.

### How would you debug DNS issues?
Use `dig`/`nslookup` to check what the domain currently resolves to, compare against expected records, check TTLs (is a stale cached record still being served somewhere), and verify authoritative name servers are responding correctly.

### Why is your application only slow in production?
Likely differences from dev/staging: real network latency (vs localhost), production data volume (bigger tables, no index yet needed at small scale), actual concurrent load, resource limits (CPU/memory caps in containers), or config differences (e.g. connection pool sizes, logging verbosity, missing caching layer that only exists in prod).

### Why is latency high only for some users?
Could be geographic (distance to nearest server/CDN edge, no regional presence near them), routing/ISP-specific network issues, those users hitting a different/degraded backend instance, or client-side factors (poor Wi-Fi, older device).

### Why is your service returning intermittent 502 errors?
Likely one or more backend instances are occasionally unhealthy/restarting/OOM-killed, a load balancer's health checks aren't catching bad instances fast enough, connection pool exhaustion causing dropped connections under load spikes, or upstream timeouts slightly shorter than what the backend actually needs.

---

## 23. Linux Networking Commands

| Command | Purpose |
|---|---|
| `ping` | Test reachability and measure round-trip latency/packet loss to a host. |
| `traceroute` / `tracepath` | Show the hop-by-hop path packets take to reach a destination. |
| `nslookup` | Query DNS records for a domain interactively. |
| `dig` | More detailed/scriptable DNS lookup tool, showing full record/response details. |
| `curl` | Make HTTP(S) (and other protocol) requests from the command line — great for inspecting headers/timing/responses. |
| `wget` | Command-line tool primarily for downloading files over HTTP(S)/FTP. |
| `netstat` (legacy) | Shows active connections, listening ports, and routing tables (largely superseded by `ss`). |
| `ss` | Modern, faster replacement for `netstat` — inspect sockets, connection states (e.g. count TIME_WAIT). |
| `lsof` | Lists open files, including network sockets — useful for finding which process holds a given port/connection. |
| `tcpdump` | Captures and inspects raw network packets from the command line. |
| `wireshark` | GUI packet analyzer — deep inspection of captured traffic, protocol decoding. |
| `ip addr` | Shows/configures network interfaces and their IP addresses. |
| `ip route` | Shows/configures the system's routing table. |
| `iptables` | Configures Linux kernel firewall rules (packet filtering, NAT). |
| `nmap` | Network scanner — discover open ports/services on hosts. |

---

## High-Priority Topics (Focus First)

If time is limited, prioritize mastery of:

- HTTP/HTTPS fundamentals and status codes
- TCP vs UDP, TCP handshake and connection lifecycle
- DNS resolution process
- TLS handshake
- Reverse proxy (Nginx) concepts
- Load balancers (L4 vs L7, algorithms)
- REST API design principles
- WebSockets
- gRPC basics
- Caching (HTTP + Redis + CDN, including invalidation strategies)
- CORS, CSRF, JWT
- Timeouts, retries, circuit breakers
- Go's `net/http`, `http.Transport`, connection pooling, `context` cancellation
- Linux networking tools (`curl`, `ss`, `dig`, `tcpdump`)
- Production debugging scenarios (latency, timeouts, DNS failures, 502/503/504 errors)
