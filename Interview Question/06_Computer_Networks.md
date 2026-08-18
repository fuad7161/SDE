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
   - A port identifies an application endpoint; a socket is a communication endpoint typically defined by protocol, address, and port.
18. What happens after you enter a URL in a browser?
   - The browser parses the URL, resolves DNS, establishes transport and TLS if needed, sends HTTP, receives resources, then parses and renders them.
19. What is HTTP? How is HTTPS different?
   - HTTP is an application protocol for web messages; HTTPS carries HTTP through TLS for encryption, integrity, and server authentication.
20. What are the common HTTP methods?
   - Common methods include GET, POST, PUT, PATCH, DELETE, HEAD, and OPTIONS, each expressing a different request intent.
21. What are the main groups of HTTP status codes?
   - The groups are 1xx informational, 2xx success, 3xx redirection, 4xx client error, and 5xx server error.
22. What is the difference between HTTP/1.1, HTTP/2, and HTTP/3?
   - HTTP/1.1 commonly uses multiple TCP connections, HTTP/2 multiplexes streams over TCP, and HTTP/3 uses QUIC over UDP.
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
