# Computer Networks

1. What is a computer network?
   - A network is a group of connected devices that exchange data and share resources using agreed protocols.
2. What are the seven layers of the OSI model?
   - From bottom to top they are Physical, Data Link, Network, Transport, Session, Presentation, and Application.
3. How does the TCP/IP model differ from the OSI model?
   - TCP/IP is the practical Internet protocol suite with fewer layers, while OSI is a seven-layer conceptual reference model.
4. What is the difference between a hub, switch, router, and gateway?
   - A hub repeats signals, a switch forwards frames within a LAN, a router connects IP networks, and a gateway translates or mediates between systems.
5. What is an IP address?
   - It is a logical network-layer address used to identify an interface and route packets across IP networks.
6. What is the difference between IPv4 and IPv6?
   - IPv4 uses 32-bit addresses; IPv6 uses 128-bit addresses and provides a vastly larger address space plus protocol improvements.
7. What is the difference between a public and private IP address?
   - Public addresses are globally routable, while private ranges are used inside local networks and normally reach the Internet through NAT.
8. What are a subnet mask and a default gateway?
   - A subnet mask identifies the local network portion of an address; the default gateway handles traffic for destinations outside it.
9. What is a MAC address, and how does it differ from an IP address?
   - A MAC address identifies a local link-layer interface, while an IP address is a logical address used for routing between networks.
10. What is ARP?
   - ARP resolves an IPv4 address to a MAC address on the local network.
11. What is DNS, and how does name resolution work?
   - DNS maps names to records such as IP addresses through caches and a hierarchy of recursive and authoritative servers.
12. What is DHCP?
   - DHCP automatically leases network settings such as an IP address, subnet mask, gateway, and DNS servers to clients.
13. What is the difference between TCP and UDP?
   - TCP is connection-oriented, ordered, and reliable; UDP sends independent datagrams with lower overhead but no delivery guarantee.
14. Explain the TCP three-way handshake.
   - The client sends SYN, the server replies SYN-ACK, and the client sends ACK to establish synchronized connection state.
15. How does TCP close a connection?
   - Each direction is normally closed independently with FIN and ACK exchanges, followed by a temporary TIME_WAIT state on the active closer.
16. How does TCP provide reliable delivery?
   - It uses sequence numbers, acknowledgements, checksums, retransmission, ordering, flow control, and congestion control.
17. What are ports and sockets?
   - A port identifies a transport-layer service on a host; a socket is an OS communication endpoint. A connected flow is commonly identified by protocol plus local and remote address-port pairs.
18. What happens after you enter a URL in a browser?
   - The browser parses the URL, resolves DNS, establishes transport and TLS if needed, sends HTTP, receives resources, then parses and renders them.
19. What is HTTP? How is HTTPS different?
   - HTTP is an application protocol for web messages; HTTPS carries HTTP through TLS for encryption, integrity, and server authentication.
20. What are the common HTTP methods?
   - Common methods include GET, POST, PUT, PATCH, DELETE, HEAD, and OPTIONS, each expressing a different request intent.
21. What are the main groups of HTTP status codes?
   - The groups are 1xx informational, 2xx success, 3xx redirection, 4xx client error, and 5xx server error.
22. What is the difference between HTTP/1.1, HTTP/2, and HTTP/3?
   - HTTP/1.1 commonly uses multiple TCP connections, HTTP/2 multiplexes streams over one TCP connection, and HTTP/3 runs over QUIC, which normally uses UDP and avoids transport-level head-of-line blocking between streams.
23. What are TLS and an SSL certificate?
   - TLS secures data in transit; a certificate binds a public key to an identity through a trusted signature. SSL is the obsolete predecessor name.
24. What are cookies, sessions, and tokens?
   - Cookies are browser-stored values, sessions keep server-side client state, and tokens carry claims or opaque credentials used to authorize requests.
25. What is the same-origin policy, and what is CORS?
   - Same-origin policy restricts cross-origin browser access; CORS uses HTTP headers to let a server explicitly permit selected origins.
26. What is a proxy? Compare forward and reverse proxies.
   - A forward proxy represents clients to servers, while a reverse proxy represents servers to clients.
27. What is a firewall?
   - A firewall allows or blocks network traffic according to security rules based on attributes such as address, port, protocol, or application.
28. What is a VPN?
   - A VPN creates an authenticated, encrypted tunnel across an untrusted network to connect a device or network securely.
29. What is a CDN?
   - A content delivery network serves cached content from geographically distributed edge locations to reduce latency and origin load.
30. What is a load balancer?
   - It distributes requests across healthy backend instances to improve capacity, availability, and fault tolerance.

## Medium to Advanced

31. How do TCP flow control and congestion control differ?
   - **Key note:** Flow control protects the receiver; congestion control protects the network by adjusting the sender's rate.
32. What are slow start, congestion avoidance, and congestion window?
   - **Key note:** TCP grows its sending window rapidly at first, then cautiously, reducing it when congestion is inferred.
33. Why does TCP have a `TIME_WAIT` state?
   - **Key note:** It lets delayed packets expire and allows retransmission of the final ACK before the connection tuple is reused.
34. What are head-of-line blocking and multiplexing?
   - **Key note:** One delayed item blocks later work; multiplexing separates logical streams, though TCP packet loss can still block HTTP/2 streams.
35. How does QUIC differ from TCP with TLS?
   - **Key note:** QUIC integrates secure transport over UDP, supports stream-level recovery, and enables faster connection establishment.
36. How does TLS 1.3 establish a secure connection?
   - **Key note:** It negotiates parameters, authenticates the server, performs ephemeral key exchange, and derives symmetric session keys.
37. What is certificate-chain validation?
   - **Key note:** The client verifies signatures to a trusted root plus hostname, validity, usage, and revocation-related policy.
38. What is DNS caching, and how do TTL and negative caching affect changes?
   - **Key note:** Resolvers retain positive or missing answers until expiry, so updates and recoveries propagate gradually.
39. What are recursive and iterative DNS resolution?
   - **Key note:** A recursive resolver obtains the final result for a client; iterative servers return answers or referrals.
40. What is anycast routing?
   - **Key note:** Multiple sites advertise the same address and routing delivers clients to a nearby reachable site.
41. How does NAT work, and what is port-address translation?
   - **Key note:** NAT rewrites addresses; PAT also maps ports so many private connections can share one public address.
42. What is the difference between Layer 4 and Layer 7 load balancing?
   - **Key note:** L4 routes using transport data; L7 understands application protocols and can route by host, path, headers, or cookies.
43. What are sticky sessions, and what problems can they cause?
   - **Key note:** They bind clients to backends but reduce balance, complicate failure recovery, and encourage local session state.
44. How does consistent hashing help distribute network traffic?
   - **Key note:** It limits remapping when nodes change and is useful for caches or stateful routing.
45. What is connection pooling and keep-alive?
   - **Key note:** Reusing established connections avoids repeated TCP/TLS handshakes but requires limits and stale-connection handling.
46. What is the difference between a timeout and a deadline?
   - **Key note:** A timeout bounds one wait; a deadline is an absolute end time that can propagate through the whole request chain.
47. How do retries amplify load during a network or service failure?
   - **Key note:** Each layer may retry the same work, creating a retry storm; use one retry point, budgets, backoff, and jitter.
48. What is a SYN flood, and how do SYN cookies help?
   - **Key note:** Attackers fill half-open connection state; SYN cookies defer server-side allocation until the final handshake ACK.
49. What is MTU, and how can path MTU problems appear?
   - **Key note:** Oversized packets require fragmentation or rejection; blocked discovery messages can cause unexplained connection stalls.
50. How would you diagnose intermittent network latency between two services?
   - **Key note:** Correlate traces with DNS, connection setup, retransmissions, packet loss, routing, proxy queues, and endpoint saturation.
