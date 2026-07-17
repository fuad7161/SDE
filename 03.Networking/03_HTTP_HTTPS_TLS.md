# 🟢 HTTP / HTTPS / TLS

> **Category:** Protocols &nbsp;|&nbsp; **Tags:** `HTTP/1.1 vs 2 vs 3` `TLS handshake` `certs`

---

## Table of Contents
1. [HTTP Basics](#http-basics)
2. [HTTP/1.1 vs HTTP/2 vs HTTP/3](#http11-vs-http2-vs-http3)
3. [HTTPS & TLS](#https--tls)
4. [TLS Handshake](#tls-handshake)
5. [Certificates](#certificates)
6. [Interview Questions](#interview-questions)

---

## HTTP Basics

**HTTP (HyperText Transfer Protocol)** is a stateless, application-layer protocol for transferring data on the web. Every HTTP exchange follows a **request–response** model.

### Request Structure
```
GET /index.html HTTP/1.1
Host: example.com
Accept: text/html
User-Agent: Mozilla/5.0
```

### Response Structure
```
HTTP/1.1 200 OK
Content-Type: text/html
Content-Length: 1234

<html>...</html>
```

### Common HTTP Methods

| Method | Purpose | Idempotent | Safe |
|--------|---------|------------|------|
| GET | Retrieve resource | ✅ | ✅ |
| POST | Create resource | ❌ | ❌ |
| PUT | Replace resource | ✅ | ❌ |
| PATCH | Partially update | ❌ | ❌ |
| DELETE | Remove resource | ✅ | ❌ |
| HEAD | Headers only (no body) | ✅ | ✅ |
| OPTIONS | Describe allowed methods | ✅ | ✅ |

### Common Status Codes
| Code | Meaning |
|------|---------|
| 200 | OK |
| 201 | Created |
| 204 | No Content |
| 301 | Moved Permanently |
| 302 | Found (Temporary Redirect) |
| 304 | Not Modified (cached) |
| 400 | Bad Request |
| 401 | Unauthorized |
| 403 | Forbidden |
| 404 | Not Found |
| 429 | Too Many Requests |
| 500 | Internal Server Error |
| 502 | Bad Gateway |
| 503 | Service Unavailable |

---

## HTTP/1.1 vs HTTP/2 vs HTTP/3

| Feature | HTTP/1.1 | HTTP/2 | HTTP/3 |
|---------|----------|--------|--------|
| Year | 1997 | 2015 | 2022 |
| Underlying protocol | TCP | TCP | QUIC (UDP) |
| Multiplexing | ❌ (one req/TCP conn) | ✅ (streams) | ✅ (streams) |
| Head-of-line blocking | At request level | At TCP level | ❌ (solved) |
| Header compression | ❌ | ✅ HPACK | ✅ QPACK |
| Server push | ❌ | ✅ | ✅ |
| Binary protocol | ❌ (text) | ✅ | ✅ |
| TLS | Optional | Optional (but de facto required) | Required |

### HTTP/1.1 Problems
- **One request per TCP connection** by default.
- Workaround: browsers open up to 6 parallel TCP connections per domain.
- **Head-of-line blocking:** Later requests must wait for earlier ones to complete.

### HTTP/2 Improvements
- **Multiplexing:** Multiple requests over one TCP connection via streams.
- **Binary framing layer:** More efficient parsing.
- **HPACK header compression:** Repeated headers sent as short indices.
- **Server Push:** Server can proactively send resources (e.g., CSS before browser asks).

### HTTP/3 Improvements (QUIC)
- Runs over **QUIC** (UDP-based), eliminating TCP head-of-line blocking.
- **0-RTT connection resumption:** Faster reconnections to known servers.
- Built-in TLS 1.3.

---

## HTTPS & TLS

**HTTPS = HTTP + TLS**. TLS (Transport Layer Security) provides:
- **Confidentiality** – Data is encrypted.
- **Integrity** – Data cannot be tampered with (MAC/HMAC).
- **Authentication** – Server identity verified via certificates.

TLS replaced SSL (Secure Sockets Layer). Current version: **TLS 1.3** (2018).

---

## TLS Handshake

### TLS 1.2 Handshake (2-RTT)
```
Client                            Server
  |  ── ClientHello ───────────→  |  (TLS version, cipher suites, random)
  |  ←─ ServerHello ───────────   |  (chosen cipher, random, certificate)
  |  ←─ Certificate ───────────   |
  |  ←─ ServerHelloDone ────────  |
  |  ── ClientKeyExchange ──────→ |  (pre-master secret, encrypted with server pubkey)
  |  ── ChangeCipherSpec ───────→ |
  |  ── Finished ───────────────→ |
  |  ←─ ChangeCipherSpec ────────  |
  |  ←─ Finished ────────────────  |
  |  === Encrypted Communication ===
```

### TLS 1.3 Handshake (1-RTT, or 0-RTT resumption)
- Removed weak cipher suites and RSA key exchange.
- Client sends key share in the first message — server can respond immediately.
- Supports **0-RTT** for session resumption (sends data with first packet, risk of replay attacks).

---

## Certificates

An **X.509 certificate** contains:
- **Subject** – Who the cert belongs to (domain name).
- **Issuer** – The Certificate Authority (CA) that signed it.
- **Public key** – Used to establish encrypted communication.
- **Validity period** – Not before / not after dates.
- **Signature** – CA's digital signature proving authenticity.

### Certificate Chain
```
Root CA (self-signed, trusted by OS/browser)
  └── Intermediate CA (signed by Root CA)
        └── End-entity cert (signed by Intermediate CA, used by the server)
```

Browsers ship with a list of **trusted Root CAs**. If the chain leads to a trusted root, the certificate is valid.

### Certificate Types
| Type | Validates |
|------|----------|
| DV (Domain Validated) | Domain ownership only |
| OV (Org Validated) | Domain + organization identity |
| EV (Extended Validation) | Full legal entity verification |
| Wildcard (`*.example.com`) | All subdomains |
| SAN (Subject Alt Name) | Multiple domains in one cert |

---

## Interview Questions

### Q1. What is the difference between HTTP and HTTPS?

> **Answer:**  
> HTTP transmits data in **plaintext** — anyone who intercepts the traffic can read it. HTTPS wraps HTTP inside a **TLS tunnel**, providing encryption (so data is unreadable to eavesdroppers), integrity (data cannot be modified in transit), and authentication (proves you're talking to the real server, not an impostor).

---

### Q2. Explain the TLS handshake process.

> **Answer:**  
> In TLS 1.3:
> 1. Client sends **ClientHello** with supported TLS versions, cipher suites, and a key share.
> 2. Server responds with **ServerHello**, its certificate, and its key share — both sides now derive the session key.
> 3. Client verifies the server's certificate against trusted CAs.
> 4. Client sends **Finished** — encrypted communication begins.
>
> This is **1-RTT**. For returning clients, TLS 1.3 supports **0-RTT** session resumption.

---

### Q3. What is the difference between HTTP/1.1, HTTP/2, and HTTP/3?

> **Answer:**
> - **HTTP/1.1:** Text-based, one request per connection (or pipelining with HOL blocking).
> - **HTTP/2:** Binary, multiplexes many requests over one TCP connection, adds HPACK header compression and server push. Still suffers from TCP-level HOL blocking.
> - **HTTP/3:** Uses QUIC over UDP, eliminating TCP HOL blocking entirely. Built-in TLS 1.3, faster connection setup (1-RTT or 0-RTT).

---

### Q4. What is a Certificate Authority and why do we trust it?

> **Answer:**  
> A CA is a trusted third party that issues digital certificates. Operating systems and browsers ship with a pre-installed list of **trusted Root CAs** (e.g., DigiCert, Let's Encrypt). When a server presents a certificate, the browser verifies its chain of trust leads back to a trusted root. If it does, the server's identity is considered verified.

---

### Q5. What is the difference between 301 and 302 redirects? When would you use each?

> **Answer:**
> - **301 Moved Permanently:** The resource has permanently moved. Browsers and search engines cache this redirect and update bookmarks/index.
> - **302 Found (Temporary Redirect):** The resource is temporarily at a different location. Not cached — browser re-checks each time.
>
> Use 301 for permanent URL changes (e.g., HTTP → HTTPS migration). Use 302 for temporary redirects (e.g., A/B testing, maintenance pages).

---

### Q6. What is HSTS?

> **Answer:**  
> **HTTP Strict Transport Security** — a response header (`Strict-Transport-Security: max-age=31536000`) that tells browsers to **always use HTTPS** for this domain for the specified duration. Prevents SSL stripping attacks. Once received, the browser refuses to connect over plain HTTP, even if the user types `http://`.

---

<div align="center">
  <sub>← Back to <a href="networking_interview_topics.md">All Topics</a></sub>
</div>
