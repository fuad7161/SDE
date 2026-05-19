# 🟢 DNS – Domain Name System

> **Category:** Protocols &nbsp;|&nbsp; **Tags:** `resolution flow` `record types` `caching / TTL`

---

## Table of Contents
1. [What is DNS?](#what-is-dns)
2. [DNS Resolution Flow](#dns-resolution-flow)
3. [DNS Record Types](#dns-record-types)
4. [Caching & TTL](#caching--ttl)
5. [DNS Security](#dns-security)
6. [Interview Questions](#interview-questions)

---

## What is DNS?

**DNS (Domain Name System)** is the internet's phone book — it translates human-readable domain names (e.g., `example.com`) into IP addresses (e.g., `93.184.216.34`) that machines use to communicate.

Without DNS, you would have to memorize IP addresses for every website.

### DNS Hierarchy

```
Root (.)
  └── Top-Level Domains: .com, .org, .net, .io, .uk
        └── Second-Level Domains: example.com, google.com
              └── Subdomains: www.example.com, api.example.com
```

---

## DNS Resolution Flow

When you type `www.example.com` in a browser:

```
Browser → OS Cache → Recursive Resolver (ISP/Google 8.8.8.8)
                           ↓
                      Root Name Server (.)
                           ↓ (returns .com NS)
                      TLD Name Server (.com)
                           ↓ (returns example.com NS)
                      Authoritative Name Server (example.com)
                           ↓ (returns IP: 93.184.216.34)
                      Recursive Resolver caches & returns IP
                           ↓
                         Browser connects to 93.184.216.34
```

### Key Players

| Component | Role |
|-----------|------|
| **DNS Stub Resolver** | On your OS — checks local cache first |
| **Recursive Resolver** | ISP or public resolver (8.8.8.8, 1.1.1.1) — does the work |
| **Root Name Server** | 13 root servers (a.root-servers.net – m.root-servers.net) |
| **TLD Name Server** | Manages `.com`, `.org`, etc. |
| **Authoritative Name Server** | Holds the actual DNS records for the domain |

### Resolution Steps (detailed)
1. Browser checks its **own cache**.
2. OS checks `/etc/hosts` and **OS DNS cache**.
3. Query sent to **Recursive Resolver** (configured in network settings).
4. Resolver checks its cache. If miss → queries **Root NS**.
5. Root NS returns address of the **.com TLD NS**.
6. Resolver queries **TLD NS** → returns address of **example.com's authoritative NS**.
7. Resolver queries **authoritative NS** → returns the IP address.
8. Resolver **caches** the result for TTL duration and returns it to the client.

---

## DNS Record Types

| Record | Purpose | Example |
|--------|---------|---------|
| **A** | Maps domain to IPv4 address | `example.com → 93.184.216.34` |
| **AAAA** | Maps domain to IPv6 address | `example.com → 2001:db8::1` |
| **CNAME** | Alias to another domain | `www.example.com → example.com` |
| **MX** | Mail exchange server | `example.com → mail.example.com` |
| **TXT** | Arbitrary text (SPF, DKIM, verification) | `"v=spf1 include:..."` |
| **NS** | Authoritative name servers for a zone | `example.com → ns1.example.com` |
| **SOA** | Start of Authority — zone metadata | Serial, refresh, retry intervals |
| **PTR** | Reverse DNS (IP → domain) | `34.216.184.93.in-addr.arpa → example.com` |
| **SRV** | Service discovery (host + port) | `_http._tcp.example.com` |
| **CAA** | Which CAs can issue certs for the domain | `example.com CAA letsencrypt.org` |

**Important:** A `CNAME` cannot coexist with other records at the same name. The root domain (`example.com`) cannot use a CNAME — use `ALIAS` or `ANAME` (provider-specific flattening).

---

## Caching & TTL

**TTL (Time To Live):** A value (in seconds) set on each DNS record. It tells resolvers how long to cache the record before re-querying.

```
example.com.  300  IN  A  93.184.216.34
              ↑
              TTL = 300 seconds (5 minutes)
```

### Caching layers (in order)
1. Browser DNS cache
2. OS DNS cache
3. Recursive resolver cache
4. TLD/root servers (rarely change)

### Practical implications
- **Low TTL (60–300s):** Changes propagate quickly. Useful before migrations.
- **High TTL (3600–86400s):** Fewer DNS queries, faster resolution, but changes are slow to propagate.
- **Before a migration:** Lower TTL to 60s a day in advance, then switch the A record.

---

## DNS Security

### DNS Spoofing / Cache Poisoning
An attacker injects a fake DNS record into a resolver's cache, redirecting users to a malicious IP.

**Defense:** DNSSEC (DNS Security Extensions) — digitally signs DNS responses so resolvers can verify they haven't been tampered with.

### DNS over HTTPS (DoH)
Encrypts DNS queries using HTTPS, preventing ISPs and attackers from eavesdropping on which domains you're querying.

### DNS over TLS (DoT)
Similar to DoH but uses a dedicated TLS connection on port 853.

---

## Interview Questions

### Q1. Walk me through what happens when you type a URL in the browser and press Enter — just the DNS part.

> **Answer:**
> 1. Browser checks its local DNS cache.
> 2. If miss, asks the OS resolver — which checks `/etc/hosts` and the system cache.
> 3. If still a miss, the OS queries the **recursive resolver** (e.g., 8.8.8.8).
> 4. The recursive resolver queries the **root name server** → gets the TLD name server address.
> 5. Queries the **TLD name server** → gets the authoritative name server address.
> 6. Queries the **authoritative name server** → gets the actual IP address.
> 7. The IP is cached at each level (for the record's TTL) and returned to the browser.

---

### Q2. What is a CNAME record and when should you NOT use one?

> **Answer:**  
> A CNAME (Canonical Name) is an alias that maps one domain to another. When resolved, the resolver follows the chain to the final A record.
>
> **Don't use CNAME at the zone apex (root domain).** `example.com` cannot be a CNAME because RFC mandates the root domain must have an SOA and NS record, which can't coexist with CNAME. Use provider-specific ALIAS/ANAME records or point the A record directly.

---

### Q3. What is DNS TTL and what are the tradeoffs of high vs low values?

> **Answer:**  
> TTL tells caching resolvers how long to keep a DNS record before re-querying.
> - **High TTL:** Fewer DNS lookups (faster resolution, less load on DNS servers). Downside: changes take longer to propagate.
> - **Low TTL:** Changes propagate quickly. Downside: more DNS queries, higher load, slightly slower resolution.
>
> Best practice: Lower TTL before a planned IP change, then restore it afterward.

---

### Q4. What is the difference between authoritative and recursive DNS servers?

> **Answer:**
> - **Authoritative DNS server:** Holds the actual DNS records for a zone. It gives a definitive answer for that domain. (e.g., the nameservers you configure at your registrar.)
> - **Recursive resolver:** Does the work of looking up DNS records on your behalf by querying the hierarchy. It caches results. (e.g., 8.8.8.8, 1.1.1.1, or your ISP's resolver.)

---

### Q5. What is DNS cache poisoning? How is it prevented?

> **Answer:**  
> Cache poisoning is when an attacker injects a forged DNS response into a resolver's cache. The next client querying that resolver gets the malicious IP (e.g., redirected to a phishing site).
>
> **Prevention:**
> - **DNSSEC:** Adds digital signatures to DNS records, so resolvers can verify authenticity.
> - **Source port randomization:** Makes it harder for attackers to guess the transaction ID.
> - **DNS over HTTPS/TLS:** Encrypts the DNS transaction entirely.

---

<div align="center">
  <sub>← Back to <a href="networking_interview_topics.md">All Topics</a></sub>
</div>
