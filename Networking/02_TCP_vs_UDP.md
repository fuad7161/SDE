# 🟢 TCP vs UDP

> **Category:** Protocols &nbsp;|&nbsp; **Tags:** `3-way handshake` `reliability` `use cases`

---

## Table of Contents
1. [Overview](#overview)
2. [TCP – Transmission Control Protocol](#tcp--transmission-control-protocol)
3. [UDP – User Datagram Protocol](#udp--user-datagram-protocol)
4. [TCP vs UDP Comparison](#tcp-vs-udp-comparison)
5. [When to Use Which](#when-to-use-which)
6. [Interview Questions](#interview-questions)

---

## Overview

Both TCP and UDP are **Layer 4 (Transport)** protocols that run on top of IP. They define how data is sent between two endpoints, but differ fundamentally in their reliability guarantees.

---

## TCP – Transmission Control Protocol

TCP is a **connection-oriented**, **reliable** protocol. It guarantees delivery, ordering, and error-checking.

### 3-Way Handshake (Connection Setup)

```
Client                    Server
  |  ── SYN ──────────→   |   (I want to connect, seq=x)
  |  ←─ SYN-ACK ───────   |   (OK, seq=y, ack=x+1)
  |  ── ACK ──────────→   |   (Got it, ack=y+1)
  |                        |
  |  === Connection Established ===
```

### 4-Way Termination (Connection Teardown)

```
Client                    Server
  |  ── FIN ──────────→   |
  |  ←─ ACK ───────────   |
  |  ←─ FIN ───────────   |
  |  ── ACK ──────────→   |
  |  (waits TIME_WAIT)     |
```

### Key TCP Features
- **Acknowledgements (ACK):** Receiver confirms every segment received.
- **Sequence numbers:** Ensures packets are reassembled in order.
- **Retransmission:** Lost segments are re-sent.
- **Flow control:** Receiver advertises a window size to prevent buffer overflow.
- **Congestion control:** Algorithms (Slow Start, AIMD) prevent network congestion.

---

## UDP – User Datagram Protocol

UDP is **connectionless** and **unreliable** — it fires and forgets. No handshake, no ACKs, no retransmission.

### UDP Header (minimal, 8 bytes)

```
| Source Port (16) | Destination Port (16) |
| Length (16)      | Checksum (16)         |
| Data ...                                  |
```

TCP header is 20–60 bytes. UDP's tiny header = low overhead = **fast**.

### Key UDP Features
- No connection setup / teardown
- No guaranteed delivery or ordering
- Application-level reliability (if needed) must be implemented manually
- Supports **multicast** and **broadcast** (TCP cannot)

---

## TCP vs UDP Comparison

| Feature | TCP | UDP |
|---------|-----|-----|
| Connection | Connection-oriented | Connectionless |
| Reliability | Guaranteed delivery | Best-effort |
| Ordering | In-order delivery | No ordering |
| Error recovery | Retransmission | No retransmission |
| Speed | Slower (overhead) | Faster (low overhead) |
| Header size | 20–60 bytes | 8 bytes |
| Flow control | Yes | No |
| Congestion control | Yes | No |
| Broadcast/Multicast | No | Yes |

---

## When to Use Which

| Use TCP | Use UDP |
|---------|---------|
| HTTP/HTTPS | DNS queries |
| Email (SMTP, IMAP) | Video streaming |
| File transfer (FTP) | VoIP / online gaming |
| SSH | Live broadcasts |
| Database connections | IoT sensor data |

**Rule of thumb:** If losing a packet is **catastrophic** → TCP. If losing a packet is **tolerable** and speed matters → UDP.

---

## Interview Questions

### Q1. What is the TCP 3-way handshake and why is it needed?

> **Answer:**  
> The 3-way handshake establishes a TCP connection:
> 1. **SYN** – Client sends a segment with the SYN flag and its initial sequence number.
> 2. **SYN-ACK** – Server responds with its own sequence number and acknowledges the client's.
> 3. **ACK** – Client acknowledges the server's sequence number.
>
> It's needed to: (1) establish that both sides can send and receive, and (2) synchronize sequence numbers for reliable ordered data transfer.

---

### Q2. What happens if a TCP packet is lost?

> **Answer:**  
> TCP uses **retransmission timeouts (RTO)** and **duplicate ACKs** to detect loss:
> - If the sender doesn't receive an ACK within a timeout period, it retransmits.
> - If the receiver sends 3 duplicate ACKs (same ACK number), the sender triggers **fast retransmit** without waiting for timeout.
> - Lost segments are retransmitted and the window size may be reduced (congestion control).

---

### Q3. Why is UDP used for video streaming instead of TCP?

> **Answer:**  
> In live video/audio:
> - A slightly stale or missing frame is better than **pausing to wait** for retransmission.
> - UDP's lower overhead means lower latency.
> - TCP's retransmission of old data would cause **buffering/stuttering** which is worse than a dropped frame.
> - Applications like WebRTC and QUIC implement their own selective reliability on top of UDP.

---

### Q4. What is TCP flow control vs congestion control?

> **Answer:**
> - **Flow control:** Prevents the **sender from overwhelming the receiver**. The receiver advertises a receive window (rwnd) in ACKs — the sender can't send more data than the window allows.
> - **Congestion control:** Prevents the **sender from overwhelming the network**. Uses algorithms like Slow Start, Congestion Avoidance, and Fast Recovery. Maintains a congestion window (cwnd). The actual send rate = min(rwnd, cwnd).

---

### Q5. Can you build reliable communication on top of UDP? Give an example.

> **Answer:**  
> Yes. **QUIC** (used by HTTP/3) is built on UDP and implements:
> - Reliable ordered streams
> - TLS 1.3 encryption
> - Multiplexing without head-of-line blocking
>
> Other examples: games implementing their own ACK/retry logic, WebRTC using SRTP over UDP, DNS over QUIC.

---

### Q6. What is TIME_WAIT in TCP and why does it exist?

> **Answer:**  
> After a connection closes, the active closer enters **TIME_WAIT** state for `2 × MSL` (Maximum Segment Lifetime, typically 2 minutes). This:
> 1. Ensures the final ACK (for the server's FIN) was received — if lost, server will re-send FIN and the client can respond.
> 2. Prevents old duplicate packets from a previous connection from interfering with a new connection on the same port/address tuple.

---

<div align="center">
  <sub>← Back to <a href="networking_interview_topics.md">All Topics</a></sub>
</div>
