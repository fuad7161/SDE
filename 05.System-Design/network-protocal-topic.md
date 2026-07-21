1. HTTP/HTTPS (Most Important)

These are asked in almost every backend interview.

Basic
Explain the HTTP request lifecycle.
What are the different HTTP methods?
Difference between GET and POST?
Difference between PUT and PATCH?
DELETE vs POST for deletion?
Safe methods vs Idempotent methods?
What is an idempotent API?
What are HTTP status codes?
200
201
202
204
301
302
307
308
400
401
403
404
409
429
500
502
503
504
Headers
What are HTTP headers?
Common headers?
Authorization
Content-Type
Accept
Cache-Control
ETag
If-None-Match
Cookie
Set-Cookie
Origin
Referer
Connection
What is Keep-Alive?
What is Connection: close?
Persistent connection?
What is chunked transfer encoding?
What is Content-Length?
Caching
Browser cache
CDN cache
Reverse proxy cache
Cache-Control
ETag
Last-Modified
Cache invalidation
HTTP Versions
HTTP/1.1
HTTP/2
HTTP/3

Questions:

Why was HTTP/2 introduced?
Multiplexing?
Header compression?
Server Push?
Why HTTP/3 uses UDP?
2. HTTPS / TLS

Very common.

Questions:

Difference between HTTP and HTTPS.
Explain TLS handshake.
What is a certificate?
Public key vs private key.
Why symmetric encryption after handshake?
Certificate Authority.
What is mTLS?
Self-signed certificate?
TLS versions.
3. TCP/IP

Extremely common.

Questions:

Explain TCP.
TCP three-way handshake.
TCP four-way termination.
Why is TCP reliable?
Sliding window.
Flow control.
Congestion control.
Retransmission.
Sequence number.
ACK.
SYN flood attack.
TIME_WAIT.
Why does TIME_WAIT exist?
What is socket?
What is port?
4. UDP

Questions:

UDP vs TCP
When use UDP?
Is UDP reliable?
How can you make UDP reliable?

Examples:

DNS
Gaming
VoIP
Video Streaming
5. DNS

Questions:

Explain DNS lookup.

What is:

A record
AAAA
MX
TXT
CNAME
NS
PTR

DNS cache?

Recursive resolver?

Authoritative server?

6. WebSocket

Very common.

Questions:

What is WebSocket?

Difference:

HTTP vs WebSocket

How handshake works?

When use WebSocket?

How is it different from polling?

Long polling vs SSE vs WebSocket?

7. gRPC

Very common for Go.

Questions:

What is gRPC?

Advantages over REST?

What is Protocol Buffers?

Unary RPC

Streaming RPC

Bidirectional streaming

HTTP/2 dependency

Deadlines

Metadata

Interceptors

8. REST API

Questions:

REST principles

Stateless

Resource

URI design

Versioning

Pagination

Filtering

Sorting

HATEOAS

9. GraphQL

Questions:

REST vs GraphQL

N+1 problem

Resolver

Schema

Mutation

Subscription

Advantages

Disadvantages

10. Authentication Protocols

Questions:

OAuth2

OpenID Connect

JWT

Session

Bearer token

Refresh token

Access token

SAML

Kerberos

LDAP

11. Message Brokers (VERY IMPORTANT)

Most companies ask this.

RabbitMQ

Questions:

Exchange types

Direct
Topic
Fanout
Headers

Queue

Binding

Routing key

Publisher confirm

Consumer ACK

Dead Letter Queue

Retry

TTL

Priority Queue

Prefetch

Kafka

Questions:

Broker

Producer

Consumer

Partition

Offset

Replication

Leader

ISR

Consumer Group

Exactly once

At least once

At most once

Retention

Compaction

Ordering

NATS

Questions:

Core NATS

JetStream

Subject

Queue Group

Durable Consumer

Streaming

Redis Pub/Sub

Questions:

Pub/Sub

Streams

Consumer Group

When use Redis instead of Kafka?

12. Event Streaming

Questions:

Difference:

Queue vs Stream

Event sourcing

CQRS

Replay

Immutable event log

13. MQTT

Often IoT companies ask.

Questions:

Publish

Subscribe

QoS

Retained messages

Will message

Broker

14. AMQP

Questions:

What is AMQP?

RabbitMQ uses AMQP?

AMQP vs MQTT

15. FTP

Questions:

FTP

SFTP

FTPS

Difference?

16. SSH

Questions:

SSH handshake

Public key authentication

SSH tunnel

Port forwarding

17. SMTP / POP3 / IMAP

Questions:

SMTP

IMAP

POP3

Which sends mail?

Which receives mail?

Difference?

18. Web Security Protocols

Questions:

CORS

CSRF

XSS

CSP

HSTS

SameSite Cookie

Secure Cookie

HttpOnly

19. Load Balancer Protocols

Questions:

Layer 4

Layer 7

TCP Load Balancer

HTTP Load Balancer

Sticky Session

Health Check

Reverse Proxy

20. Service Discovery

Questions:

Consul

etcd

ZooKeeper

DNS-based discovery

Client-side discovery

Server-side discovery

21. Serialization Protocols

Questions:

JSON

XML

MessagePack

Protocol Buffers

Avro

Thrift

When use each?

22. RPC

Questions:

RPC

gRPC

JSON-RPC

Apache Thrift

REST vs RPC

23. API Gateway

Questions:

Why API Gateway?

Authentication

Rate limiting

Circuit breaker

Request aggregation

Caching

24. Reverse Proxy

Questions:

Nginx

HAProxy

Envoy

Traefik

Forward proxy vs Reverse proxy

25. Real-World Design Questions

These often test practical decision-making rather than definitions.

When would you choose Kafka over RabbitMQ?
When is REST better than gRPC?
When should you use WebSockets instead of HTTP polling?
How would you notify 1 million users in real time?
Which protocol would you use for a live chat application, and why?
How would you design communication between microservices?
How would you guarantee that a message is processed exactly once?
How would you implement retries without creating duplicate side effects?
How would you handle backpressure when consumers are slower than producers?
How would you build a resilient event-driven system using Go?
High-Priority Topics to Master

If you're preparing for mid/senior Golang backend interviews, focus on these first:

HTTP/1.1, HTTP/2, HTTP/3
HTTPS and the TLS handshake
TCP/IP fundamentals (handshake, retransmission, flow/congestion control)
REST and gRPC
WebSockets and Server-Sent Events (SSE)
Kafka (partitions, offsets, consumer groups, delivery semantics)
RabbitMQ (exchanges, acknowledgments, dead-letter queues, retries)
Redis Pub/Sub vs Redis Streams
OAuth2, JWT, and session-based authentication
Reverse proxies, load balancers, and caching (browser, CDN, reverse proxy)
