# 🟠 Firewalls & VPN

> **Category:** Security &nbsp;|&nbsp; **Tags:** `stateful/stateless` `tunneling` `IPSec`

---

## Table of Contents
1. [Firewalls](#firewalls)
2. [Stateless vs Stateful Firewalls](#stateless-vs-stateful-firewalls)
3. [Firewall Types](#firewall-types)
4. [VPN – Virtual Private Network](#vpn--virtual-private-network)
5. [VPN Protocols](#vpn-protocols)
6. [IPSec](#ipsec)
7. [Interview Questions](#interview-questions)

---

## Firewalls

A **firewall** is a network security device (hardware or software) that monitors and controls incoming and outgoing traffic based on defined security rules.

**Core function:** Filter traffic to allow legitimate connections and block malicious or unauthorized ones.

### What firewalls filter on:
- IP source/destination address
- Port numbers (source and destination)
- Protocol (TCP, UDP, ICMP)
- Connection state (for stateful firewalls)
- Application-layer content (for next-gen firewalls)

---

## Stateless vs Stateful Firewalls

### Stateless Firewall
- Evaluates **each packet independently** against a fixed rule set.
- Does not track connection state.
- Fast, simple — but easily fooled.

```
Rule: Allow TCP from 10.0.0.0/8 on port 443
Every packet evaluated independently — no context of whether a connection exists.
```

**Weakness:** An attacker can craft a packet with ACK flag set (looks like part of an established connection) and it will pass through.

---

### Stateful Firewall
- Maintains a **state table** of active connections.
- Tracks TCP handshakes — only allows packets that belong to a known, established connection.
- Blocks unsolicited inbound packets even if they match IP/port rules.

```
State table entry:
  SrcIP: 10.0.0.5  SrcPort: 52341
  DstIP: 93.184.216.34  DstPort: 443
  State: ESTABLISHED
```

**Advantages:** More secure, aware of connection context.  
**Disadvantages:** Higher memory/CPU usage, can be overwhelmed by SYN flood attacks.

---

## Firewall Types

| Type | Layer | Capability |
|------|-------|-----------|
| **Packet filter** | L3/L4 | IP, port, protocol rules |
| **Stateful inspection** | L3/L4 | + Connection state tracking |
| **Application-layer (proxy)** | L7 | Inspects HTTP, FTP, DNS content |
| **Next-Generation Firewall (NGFW)** | L3–L7 | + IDS/IPS, deep packet inspection, app awareness |
| **WAF (Web Application Firewall)** | L7 | HTTP-specific — blocks SQLi, XSS, CSRF |

### Firewall Placement

```
Internet → [Edge Firewall] → DMZ (public-facing servers) → [Internal Firewall] → Internal Network
```

- **DMZ (Demilitarized Zone):** Segment for internet-facing servers (web, email, DNS) isolated from internal network.

---

## VPN – Virtual Private Network

A **VPN** creates an encrypted tunnel over a public network (internet), allowing remote users or networks to communicate as if they were on the same private network.

### Use Cases
- **Remote access VPN:** Employee connects to corporate network from home.
- **Site-to-site VPN:** Connect two office networks over the internet.
- **Consumer VPN:** Mask IP address, bypass geo-restrictions.

### How Tunneling Works

```
[Client] ─── Encrypted Tunnel ──────────────────→ [VPN Gateway] ──→ [Private Network]
  Outer packet: Client IP → VPN Server IP
  Inner packet: Client IP → Internal server IP (encrypted)
```

The VPN gateway **decapsulates** the outer packet and forwards the inner packet to the internal network.

---

## VPN Protocols

| Protocol | Description | Port |
|----------|-------------|------|
| **IPSec** | Industry-standard, operates at Network layer | UDP 500/4500 |
| **OpenVPN** | Open-source, TLS-based, highly configurable | UDP/TCP 1194 |
| **WireGuard** | Modern, fast, minimal codebase (~4K lines) | UDP 51820 |
| **L2TP/IPSec** | L2TP for tunneling + IPSec for encryption | UDP 1701 |
| **SSTP** | Microsoft, uses HTTPS (port 443), bypasses firewalls | TCP 443 |
| **IKEv2** | Fast reconnection (MOBIKE), mobile-friendly | UDP 500/4500 |

---

## IPSec

**IPSec (Internet Protocol Security)** is a suite of protocols for securing IP communication by authenticating and encrypting each IP packet.

### Two Main Protocols

| Protocol | Provides | Header Overhead |
|----------|----------|-----------------|
| **AH (Authentication Header)** | Integrity + authentication (no encryption) | 24 bytes |
| **ESP (Encapsulating Security Payload)** | Encryption + integrity + authentication | 24+ bytes |

In practice, **ESP** is almost always used (provides all three).

### Two Modes

| Mode | What's Encrypted | Use Case |
|------|-----------------|---------|
| **Transport Mode** | Payload only (IP header unchanged) | Host-to-host |
| **Tunnel Mode** | Entire original IP packet (new outer IP header added) | Gateway-to-gateway VPN |

### IKE (Internet Key Exchange)
IPSec uses **IKE** (usually IKEv2) to negotiate security associations (SAs) and exchange cryptographic keys:
1. **Phase 1:** Establish a secure, authenticated channel (IKE SA).
2. **Phase 2:** Negotiate IPSec SA — agree on encryption algorithm, keys, lifetime.

---

## Interview Questions

### Q1. What is the difference between a stateless and stateful firewall?

> **Answer:**
> - **Stateless:** Examines each packet independently using static rules (IP, port, protocol). Fast but easily bypassed — it can't distinguish a legitimate response packet from a crafted attack packet.
> - **Stateful:** Maintains a connection state table. Only allows packets that belong to a known established connection. Can detect and block SYN-only packets, ACK without SYN, etc.
>
> Stateful is more secure for general use. Stateless is used in high-speed scenarios like ACLs on core routers.

---

### Q2. What is a VPN and how does tunneling work?

> **Answer:**  
> A VPN creates an **encrypted tunnel** over an untrusted network. Tunneling works by **encapsulating** the original (inner) packet inside a new (outer) packet. The outer packet is routed to the VPN gateway over the internet. The gateway decapsulates it, decrypts the inner packet, and forwards it to the private network.
>
> This makes remote traffic appear as if it originates from inside the private network.

---

### Q3. What is the difference between IPSec Transport mode and Tunnel mode?

> **Answer:**
> - **Transport mode:** Only the **payload** is encrypted. The original IP header is preserved. Used between two hosts directly.
> - **Tunnel mode:** The **entire original IP packet** (header + payload) is encrypted and encapsulated in a new IP packet. Used for site-to-site VPNs between gateways.
>
> Tunnel mode is more common because it hides internal network structure.

---

### Q4. What is a DMZ in network security?

> **Answer:**  
> A **DMZ (Demilitarized Zone)** is a network segment that sits between the public internet and the private internal network. Internet-facing servers (web, email, DNS) are placed here, so if they are compromised, the attacker cannot directly reach the internal corporate network. Typically enforced by two firewalls:
> - External firewall: internet → DMZ
> - Internal firewall: DMZ → internal network

---

### Q5. What is a WAF and how is it different from a regular firewall?

> **Answer:**
> - A **regular firewall** operates at L3/L4 — it filters based on IP addresses, ports, and protocols. It cannot inspect HTTP request content.
> - A **WAF (Web Application Firewall)** operates at L7 — it understands HTTP and inspects request/response bodies, headers, cookies, and URLs. It can detect and block **SQL injection, XSS, CSRF, path traversal**, and other application-layer attacks.
>
> A WAF is not a replacement for a firewall — they are complementary.

---

<div align="center">
  <sub>← Back to <a href="networking_interview_topics.md">All Topics</a></sub>
</div>
