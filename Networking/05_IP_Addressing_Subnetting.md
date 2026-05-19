# 🟣 IP Addressing & Subnetting

> **Category:** Fundamentals &nbsp;|&nbsp; **Tags:** `IPv4 vs IPv6` `CIDR` `NAT`

---

## Table of Contents
1. [IPv4 Addressing](#ipv4-addressing)
2. [IPv6 Addressing](#ipv6-addressing)
3. [CIDR – Classless Inter-Domain Routing](#cidr--classless-inter-domain-routing)
4. [Subnetting](#subnetting)
5. [NAT – Network Address Translation](#nat--network-address-translation)
6. [Special IP Ranges](#special-ip-ranges)
7. [Interview Questions](#interview-questions)

---

## IPv4 Addressing

An **IPv4 address** is a 32-bit number represented as four 8-bit octets in dotted-decimal notation.

```
192  .  168  .   1  .   1
 ↑        ↑       ↑      ↑
8 bits  8 bits  8 bits  8 bits  =  32 bits total
```

- Total address space: **2³² = ~4.3 billion** addresses.
- Addresses are split into a **network portion** and a **host portion** by the subnet mask.

### Original IP Classes (classful, now obsolete)

| Class | Range | Default Mask | Hosts |
|-------|-------|--------------|-------|
| A | 0.0.0.0 – 127.255.255.255 | /8 | 16M |
| B | 128.0.0.0 – 191.255.255.255 | /16 | 65K |
| C | 192.0.0.0 – 223.255.255.255 | /24 | 254 |
| D | 224.0.0.0 – 239.255.255.255 | Multicast | — |
| E | 240.0.0.0 – 255.255.255.255 | Reserved | — |

---

## IPv6 Addressing

IPv6 uses **128-bit** addresses, written as 8 groups of 4 hex digits:

```
2001:0db8:85a3:0000:0000:8a2e:0370:7334
```

**Shortening rules:**
- Leading zeros in each group can be omitted: `0db8` → `db8`
- One consecutive group of all-zeros can be replaced by `::`: `2001:db8::8a2e:370:7334`

- Total address space: **2¹²⁸ ≈ 3.4 × 10³⁸** — effectively unlimited.
- No need for NAT — every device can have a globally unique address.
- Built-in **IPSec** support.
- Replaced broadcasts with **multicast** and **anycast**.

---

## CIDR – Classless Inter-Domain Routing

**CIDR** (introduced 1993) replaced classful addressing. It uses a **prefix length** (the `/n` notation) to specify how many bits are the network portion.

```
192.168.1.0/24
             ↑
             24 bits for network, 8 bits for hosts
             → 2^8 - 2 = 254 usable host addresses
```

### CIDR to Subnet Mask

| CIDR | Subnet Mask | # Hosts |
|------|-------------|---------|
| /8 | 255.0.0.0 | 16,777,214 |
| /16 | 255.255.0.0 | 65,534 |
| /24 | 255.255.255.0 | 254 |
| /25 | 255.255.255.128 | 126 |
| /26 | 255.255.255.192 | 62 |
| /27 | 255.255.255.224 | 30 |
| /28 | 255.255.255.240 | 14 |
| /30 | 255.255.255.252 | 2 |
| /32 | 255.255.255.255 | 1 (single host) |

**Usable hosts = 2^(32-prefix) - 2** (subtract network address + broadcast)

---

## Subnetting

Subnetting divides a network into smaller **subnetworks** for better organization, security, and traffic management.

### Example: Split 192.168.1.0/24 into 4 subnets

- Need 4 subnets → need 2 extra bits → /26
- Each subnet has 2^6 - 2 = **62 usable hosts**

| Subnet | Network | Usable Range | Broadcast |
|--------|---------|--------------|-----------|
| 1 | 192.168.1.0/26 | .1 – .62 | .63 |
| 2 | 192.168.1.64/26 | .65 – .126 | .127 |
| 3 | 192.168.1.128/26 | .129 – .190 | .191 |
| 4 | 192.168.1.192/26 | .193 – .254 | .255 |

---

## NAT – Network Address Translation

**NAT** allows multiple devices on a private network to share a single public IP address. The router maintains a translation table mapping `(private IP : port) ↔ (public IP : port)`.

```
[Device A: 192.168.1.10:5000] ──→ Router NAT ──→ [Public IP: 203.0.113.1:40001] ──→ Internet
[Device B: 192.168.1.11:5000] ──→ Router NAT ──→ [Public IP: 203.0.113.1:40002] ──→ Internet
```

### Types of NAT

| Type | Description |
|------|-------------|
| **SNAT** (Source NAT) | Replaces source IP — outbound connections (most common) |
| **DNAT** (Destination NAT) | Replaces destination IP — port forwarding |
| **PAT / Masquerade** | Many-to-one (uses port numbers to distinguish connections) |

**Why NAT exists:** IPv4 exhaustion. ~4B addresses aren't enough for billions of devices.  
**NAT disadvantages:** Breaks end-to-end connectivity, complicates peer-to-peer protocols.

---

## Special IP Ranges

| Range | Purpose |
|-------|---------|
| `10.0.0.0/8` | Private (RFC 1918) |
| `172.16.0.0/12` | Private (RFC 1918) |
| `192.168.0.0/16` | Private (RFC 1918) |
| `127.0.0.0/8` | Loopback (`127.0.0.1` = localhost) |
| `169.254.0.0/16` | Link-local / APIPA (no DHCP found) |
| `0.0.0.0/0` | Default route (all traffic) |
| `255.255.255.255/32` | Limited broadcast |
| `224.0.0.0/4` | Multicast |

---

## Interview Questions

### Q1. What is the difference between a public and private IP address?

> **Answer:**
> - **Private IPs** (RFC 1918): `10.x.x.x`, `172.16-31.x.x`, `192.168.x.x` — used within private networks, not routable on the public internet.
> - **Public IPs:** Globally unique, routable on the internet, assigned by ISPs and IANA.
>
> NAT allows devices with private IPs to communicate on the internet via a shared public IP.

---

### Q2. How do you calculate the number of usable hosts in a subnet?

> **Answer:**  
> **Formula:** `2^(32 - prefix_length) - 2`
>
> The `-2` removes the **network address** (all host bits = 0) and **broadcast address** (all host bits = 1).
>
> Example: `/26` → `2^(32-26) - 2 = 2^6 - 2 = 64 - 2 = 62` usable hosts.

---

### Q3. What is CIDR and why was it introduced?

> **Answer:**  
> CIDR (Classless Inter-Domain Routing) allows flexible allocation of IP blocks using prefix lengths (e.g., `/22`), instead of the rigid Class A/B/C system. It was introduced to:
> 1. **Slow IPv4 exhaustion** by allowing allocations of any size.
> 2. **Reduce routing table size** through supernetting/route aggregation (many small networks summarized as one prefix).

---

### Q4. What is NAT and what problem does it solve?

> **Answer:**  
> NAT (Network Address Translation) maps multiple private IPs to a single public IP by using different port numbers. It solves **IPv4 address exhaustion** — a home router with one public IP can serve hundreds of devices.
>
> The router maintains a NAT table: `(privateIP:port) ↔ (publicIP:port)`. For outbound traffic, source IP is rewritten; inbound responses are translated back.

---

### Q5. What is the difference between IPv4 and IPv6?

> **Answer:**
>
> | | IPv4 | IPv6 |
> |--|------|------|
> | Address size | 32-bit | 128-bit |
> | Format | Dotted decimal | Hex groups |
> | Address space | ~4.3 billion | ~3.4 × 10³⁸ |
> | NAT required | Often yes | No (abundant addresses) |
> | Header | Variable (20-60B) | Fixed 40B |
> | Fragmentation | Router or host | Source only |
> | IPSec | Optional | Built-in |
> | Broadcast | Yes | Replaced by multicast |

---

### Q6. A server is at 10.0.0.5/28. What is its subnet, and how many hosts are in that subnet?

> **Answer:**  
> - `/28` → subnet mask `255.255.255.240` → 4 host bits
> - Block size = `2^4 = 16`. Subnets are at multiples of 16: `0, 16, 32...`
> - `10.0.0.5` falls in the `10.0.0.0` – `10.0.0.15` block.
> - **Network:** `10.0.0.0`, **Broadcast:** `10.0.0.15`, **Usable:** `10.0.0.1` – `10.0.0.14` → **14 hosts**.

---

<div align="center">
  <sub>← Back to <a href="networking_interview_topics.md">All Topics</a></sub>
</div>
