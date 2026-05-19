# 🟣 OSI Model

> **Category:** Fundamentals &nbsp;|&nbsp; **Tags:** `7 layers` `encapsulation` `PDUs`

---

## Table of Contents
1. [What is the OSI Model?](#what-is-the-osi-model)
2. [The 7 Layers](#the-7-layers)
3. [Encapsulation & PDUs](#encapsulation--pdus)
4. [OSI vs TCP/IP Model](#osi-vs-tcpip-model)
5. [Interview Questions](#interview-questions)

---

## What is the OSI Model?

The **Open Systems Interconnection (OSI)** model is a conceptual framework that standardizes how different network systems communicate. It divides network communication into **7 distinct layers**, each with a specific responsibility.

It was introduced by ISO (International Organization for Standardization) in 1984. It is not a real-world implementation — it's a **reference model** used to understand and troubleshoot networking.

---

## The 7 Layers

| # | Layer | Role | Protocol Examples | PDU |
|---|-------|------|-------------------|-----|
| 7 | **Application** | User-facing interfaces, network services | HTTP, FTP, SMTP, DNS | Data |
| 6 | **Presentation** | Data formatting, encryption, compression | SSL/TLS, JPEG, ASCII | Data |
| 5 | **Session** | Session establishment, maintenance, termination | NetBIOS, RPC | Data |
| 4 | **Transport** | End-to-end delivery, error recovery, flow control | TCP, UDP | Segment |
| 3 | **Network** | Logical addressing, routing | IP, ICMP, OSPF | Packet |
| 2 | **Data Link** | Physical addressing (MAC), error detection | Ethernet, Wi-Fi (802.11) | Frame |
| 1 | **Physical** | Bits over physical medium (cables, radio) | USB, Fiber, Coax | Bit |

**Memory trick:** _"All People Seem To Need Data Processing"_ (Application → Physical)  
Or bottom-up: _"Please Do Not Throw Sausage Pizza Away"_

---

## Encapsulation & PDUs

When data travels **down** the OSI stack (sender side), each layer **wraps** the data with its own header (and sometimes trailer). This is called **encapsulation**.

```
Application Data
    ↓ + HTTP header         → Data
    ↓ + TCP header          → Segment
    ↓ + IP header           → Packet
    ↓ + Ethernet header/FCS → Frame
    ↓ converted to bits     → Bits
```

On the **receiver side**, each layer **strips** its header and passes the data up — called **de-encapsulation**.

**PDU (Protocol Data Unit):** The name given to data at each layer:
- Layer 4 → **Segment**
- Layer 3 → **Packet**
- Layer 2 → **Frame**
- Layer 1 → **Bit**

---

## OSI vs TCP/IP Model

| OSI Layer | TCP/IP Layer |
|-----------|-------------|
| Application (7) | Application |
| Presentation (6) | Application |
| Session (5) | Application |
| Transport (4) | Transport |
| Network (3) | Internet |
| Data Link (2) | Network Access |
| Physical (1) | Network Access |

The **TCP/IP model** is the practical model used in the real internet — it collapses the top 3 OSI layers into one "Application" layer.

---

## Interview Questions

### Q1. What are the 7 layers of the OSI model and what does each do?

> **Answer:**
> 1. **Physical** – Transmits raw bits over a physical medium (cables, radio waves).
> 2. **Data Link** – Handles MAC addresses, frames data, error detection (CRC).
> 3. **Network** – Logical addressing (IP), routing between networks.
> 4. **Transport** – Reliable (TCP) or unreliable (UDP) end-to-end delivery.
> 5. **Session** – Manages sessions/connections between applications.
> 6. **Presentation** – Data serialization, encryption, compression.
> 7. **Application** – Protocols used by applications (HTTP, FTP, DNS).

---

### Q2. At which OSI layer does a router operate? What about a switch?

> **Answer:**
> - **Router** → Layer 3 (Network) — routes packets based on IP addresses.
> - **Switch** → Layer 2 (Data Link) — forwards frames based on MAC addresses.
> - A **Layer 3 switch** can also route at the Network layer.

---

### Q3. What is encapsulation in the OSI model?

> **Answer:**  
> Encapsulation is the process of each layer adding its own header (and sometimes trailer) to the data passed down from the layer above. It allows each layer to operate independently, only concerned with its own header. The reverse process (de-encapsulation) happens at the receiver.

---

### Q4. What is the difference between the OSI model and the TCP/IP model?

> **Answer:**
> - The OSI model has 7 layers; TCP/IP has 4 layers.
> - TCP/IP merges OSI's Application + Presentation + Session into one Application layer.
> - TCP/IP merges OSI's Data Link + Physical into one Network Access layer.
> - OSI is a **conceptual reference model**; TCP/IP is the **actual protocol suite** used on the internet.

---

### Q5. Where does TLS/SSL operate in the OSI model?

> **Answer:**  
> TLS operates primarily at **Layer 6 (Presentation)** — it handles encryption, decryption, and certificate verification. Some refer to it as straddling Layers 4–7 since it uses TCP (Layer 4) and is used by application protocols (Layer 7), but the encryption logic is at the Presentation layer.

---

<div align="center">
  <sub>← Back to <a href="networking_interview_topics.md">All Topics</a></sub>
</div>
