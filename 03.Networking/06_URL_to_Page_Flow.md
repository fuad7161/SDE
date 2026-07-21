# 🟣 URL → Page: The Full Request Flow

> **Category:** Fundamentals &nbsp;|&nbsp; **Tags:** `DNS → TCP → TLS` `HTTP → render`

---

## Table of Contents
1. [Overview](#overview)
2. [Step-by-Step Flow](#step-by-step-flow)
3. [The Full Diagram](#the-full-diagram)
4. [Browser Rendering](#browser-rendering)
5. [Interview Questions](#interview-questions)

---

## Overview

When a user types `https://www.example.com/index.html` and presses Enter, a complex sequence of events unfolds — involving DNS, TCP, TLS, HTTP, and browser rendering — before the page appears.

This is one of the **most common system design interview questions**.

---

## Step-by-Step Flow

### Step 1: URL Parsing

The browser parses the URL:
```
https://www.example.com/index.html?q=hello#section
  ↑          ↑                ↑        ↑        ↑
scheme     host             path    query   fragment
```
- **Scheme** → use HTTPS (port 443 default)
- **Host** → need to resolve `www.example.com` to an IP

---

### Step 2: DNS Resolution

1. Browser checks its **DNS cache**.
2. OS checks `/etc/hosts` and its **system cache**.
3. Query sent to **recursive resolver** (e.g., 8.8.8.8).
4. Resolver walks the DNS hierarchy:
   - Root NS → TLD NS (`.com`) → Authoritative NS → returns IP (e.g., `93.184.216.34`)
5. IP is cached per TTL; returned to browser.

---

### Step 3: TCP Connection (3-Way Handshake)

Browser opens a TCP connection to `93.184.216.34:443`:

```
Client  ──── SYN ────────→  Server
Client  ←─── SYN-ACK ─────  Server
Client  ──── ACK ────────→  Server
```

---

### Step 4: TLS Handshake (HTTPS)

On top of the TCP connection, TLS is negotiated (TLS 1.3, 1-RTT):

```
Client  ──── ClientHello (key share) ────→  Server
Client  ←─── ServerHello + Certificate  ──  Server
             ← both derive session key →
Client  ──── Finished (encrypted) ───────→  Server
```

The browser verifies the server's certificate against trusted CAs.

---

### Step 5: HTTP Request

Browser sends the HTTP request over the encrypted TLS connection:

```
GET /index.html?q=hello HTTP/1.1
Host: www.example.com
Accept: text/html,application/xhtml+xml
Accept-Language: en-US
Cookie: session=abc123
```

---

### Step 6: Server Processing

1. Request hits a **load balancer** (if present) → routed to an app server.
2. App server processes the request:
   - Authentication / authorization check.
   - Business logic execution.
   - Database query (if dynamic content).
   - CDN cache check (if static content).
3. Server generates the HTTP response.

---

### Step 7: HTTP Response

```
HTTP/1.1 200 OK
Content-Type: text/html; charset=UTF-8
Content-Encoding: gzip
Cache-Control: max-age=3600
Set-Cookie: session=abc123; Secure; HttpOnly

<!DOCTYPE html>
<html>...
```

---

### Step 8: Browser Rendering

1. **HTML Parsing** → Builds the **DOM** (Document Object Model).
2. **CSS Parsing** → Builds the **CSSOM** (CSS Object Model).
3. **Render Tree** = DOM + CSSOM (only visible elements).
4. **Layout (Reflow)** → Calculates position and size of every element.
5. **Painting** → Fills in pixels.
6. **Compositing** → Layers are merged and displayed on screen.

Additional requests triggered:
- `<link rel="stylesheet">` → CSS files (render-blocking)
- `<script src="...">` → JS files (parser-blocking unless `async`/`defer`)
- `<img src="...">` → Images
- Fonts, API calls, etc.

---

## The Full Diagram

```
User types URL
      ↓
  URL Parsing
      ↓
  DNS Lookup ──→ [Cache Hit] ──→ IP returned
      ↓ (cache miss)
  Root NS → TLD NS → Auth NS → IP
      ↓
  TCP 3-Way Handshake (SYN/SYN-ACK/ACK)
      ↓
  TLS 1.3 Handshake (1-RTT)
      ↓
  HTTP GET Request sent
      ↓
  Load Balancer → App Server → DB/Cache
      ↓
  HTTP Response (HTML)
      ↓
  Browser: Parse HTML → DOM
                      → Fetch CSS/JS/Images
                      → CSSOM → Render Tree → Layout → Paint
      ↓
  Page visible to user
```

---

## Browser Rendering

### Critical Rendering Path

The sequence from receiving HTML to displaying pixels:

1. **Parse HTML** → DOM Tree
2. **Parse CSS** → CSSOM Tree
3. **Execute JS** (may modify DOM/CSSOM)
4. **Render Tree** (DOM + CSSOM, visible nodes only)
5. **Layout** (position + size of each element)
6. **Paint** (rasterize each node to pixels)
7. **Composite** (GPU layers merged, display)

### Render-blocking resources
- CSS is **render-blocking** — browser won't paint until all CSS is loaded.
- JS is **parser-blocking** by default — use `async` or `defer`:
  - `async` — download in parallel, execute immediately when ready.
  - `defer` — download in parallel, execute after HTML parsing completes.

---

## Interview Questions

### Q1. Walk me through what happens when you type a URL and press Enter.

> **Answer** (summary):
> 1. **URL parsing** — identify scheme, host, path.
> 2. **DNS resolution** — resolve hostname to IP (browser cache → OS → recursive resolver → DNS hierarchy).
> 3. **TCP handshake** — 3-way SYN/SYN-ACK/ACK to establish connection.
> 4. **TLS handshake** — negotiate encryption, verify certificate.
> 5. **HTTP request** — send GET request with headers/cookies.
> 6. **Server processing** — load balancer → app server → database → response.
> 7. **HTTP response** — receive HTML with status 200.
> 8. **Browser rendering** — parse HTML/CSS, execute JS, layout, paint, display.

---

### Q2. How does HTTPS differ from HTTP in this flow?

> **Answer:**  
> After the TCP handshake, an additional **TLS handshake** occurs before any HTTP data is sent. This negotiates a symmetric encryption key, verifies the server's identity via its certificate, and establishes an encrypted tunnel. All HTTP traffic (request + response headers + body) is then encrypted inside TLS.

---

### Q3. What is the Critical Rendering Path and how would you optimize it?

> **Answer:**  
> The CRP is the sequence of steps the browser takes to convert HTML/CSS/JS into pixels. Optimizations:
> - **Minimize render-blocking CSS:** Inline critical CSS, defer non-critical stylesheets.
> - **Defer/async JS:** Move `<script>` to end of body or use `defer`/`async`.
> - **Reduce resource size:** Minify, compress (gzip/brotli), optimize images.
> - **Use HTTP/2:** Multiplexing for concurrent resource loading.
> - **Preconnect/preload:** `<link rel="preconnect">` and `<link rel="preload">` hints.
> - **CDN:** Serve static assets from edge servers close to the user.

---

### Q4. What could cause a page to be slow to load? How would you diagnose it?

> **Answer:**  
> Potential bottlenecks at each stage:
> - **DNS:** High TTL means no caching issue; low TTL or slow resolver adds latency. Fix: use fast resolver, prefetch DNS.
> - **TCP/TLS:** High latency to server (geographic distance). Fix: CDN, edge servers.
> - **Server processing:** Slow DB queries, no caching. Fix: query optimization, Redis caching.
> - **Response size:** Large uncompressed HTML/JS/CSS. Fix: gzip, minification, code splitting.
> - **Rendering:** Render-blocking resources. Fix: async/defer JS, inline critical CSS.
>
> Diagnosis: Chrome DevTools Network tab, Lighthouse, WebPageTest.

---

### Q5. What is the difference between `async` and `defer` on a script tag?

> **Answer:**
> - **`async`:** Script is downloaded in parallel with HTML parsing. It **executes immediately** when downloaded, which may interrupt parsing. Order of execution is not guaranteed.
> - **`defer`:** Script is downloaded in parallel. It **executes after HTML parsing** is complete, in document order.
>
> Use `defer` for scripts that need the DOM. Use `async` for independent scripts (e.g., analytics).

---

<div align="center">
  <sub>← Back to <a href="networking_interview_topics.md">All Topics</a></sub>
</div>
