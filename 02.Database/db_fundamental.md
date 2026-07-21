1. Database Fundamentals
Q1. What is ACID?
Atomicity
Consistency
Isolation
Durability

Follow-up:

Why do we need ACID?
Which databases fully support ACID?
Does MongoDB support ACID?
Q2. Explain database normalization.

Know

1NF
2NF
3NF
BCNF

Follow-up

Why normalize?
When should you denormalize?
Q3. Normalization vs Denormalization

Expected discussion

Advantages

Less redundancy
Easier updates

Disadvantages

More joins

When denormalization is better

Reporting
Analytics
Read-heavy systems
Q4. Primary Key vs Unique Key

Expected

NULL behavior
Number of keys allowed
Index creation
Q5. Clustered vs Non-clustered Index

Very common.

Need to explain

Physical storage
B+ Tree
Leaf nodes
Performance differences
2. Indexing (Very Frequently Asked)
Q6. What is an Index?

Explain

B+ Tree
Faster lookup
Extra storage
Slower writes
Q7. How does a B+ Tree work?

Senior interviews almost always ask this.

Need to explain

Internal nodes
Leaf nodes
Balanced tree
O(log n)
Q8. Why don't databases use Binary Search Trees?

Need

Unbalanced
Disk I/O
Fan-out
Tree height
Q9. Hash Index vs B-tree Index

Need comparison

Equality search

vs

Range search

Q10. Composite Index

Example

CREATE INDEX idx_user_age
ON users(country, age);

Questions

Can this query use the index?

WHERE country='BD'

YES

How about

WHERE age=25

NO

Need leftmost prefix rule.

Q11. Covering Index

What is it?

Why does it improve performance?

Q12. Why are too many indexes bad?

Expected

Insert slower
Update slower
Delete slower
More storage
3. SQL Questions
Q13. Difference between WHERE and HAVING

Classic interview question.

Q14. INNER JOIN vs LEFT JOIN vs RIGHT JOIN vs FULL JOIN

Know with examples.

Q15. Explain different JOIN algorithms.

Senior interviews

Nested Loop Join
Hash Join
Merge Join
Q16. GROUP BY

How does it work internally?

Q17. DISTINCT vs GROUP BY
Q18. UNION vs UNION ALL
Q19. EXISTS vs IN

Performance discussion.

Q20. DELETE vs TRUNCATE vs DROP

Very common.

4. Transactions
Q21. What is a transaction?
Q22. Transaction Lifecycle
BEGIN

SQL

COMMIT

ROLLBACK
Q23. What is Isolation Level?

Need all four

Read Uncommitted

Read Committed

Repeatable Read

Serializable

Q24. Explain Dirty Read
Q25. Explain Non-repeatable Read
Q26. Explain Phantom Read
Q27. Which isolation level does PostgreSQL use?

Expected

Read Committed

Repeatable Read

Serializable

5. Locks
Q28. Optimistic Locking
Q29. Pessimistic Locking
Q30. Row Lock vs Table Lock
Q31. What causes deadlock?

Very common.

Need example.

Q32. How do you avoid deadlocks?
6. Query Optimization
Q33. Why is my query slow?

Possible answers

Missing index

Too many joins

Sequential scan

Large result

Bad execution plan

Q34. How do you analyze a slow query?

Postgres

EXPLAIN

or

EXPLAIN ANALYZE
Q35. What is Query Execution Plan?
Q36. What is Sequential Scan?
Q37. Index Scan vs Bitmap Scan

Senior-level topic.

7. PostgreSQL (Very Common for Go)
Q38. MVCC

Must know.

Explain

Snapshot
Versioning
No read lock
Q39. VACUUM

Why needed?

Q40. Autovacuum
Q41. WAL (Write Ahead Log)

Very common.

Q42. Checkpoint
Q43. HOT Update
Q44. TOAST
8. Replication
Q45. Primary-Replica Replication
Q46. Synchronous vs Asynchronous Replication
Q47. Replication Lag
Q48. Read Replica

When should you use it?

9. Scaling
Q49. Vertical vs Horizontal Scaling
Q50. Sharding

Need

Pros

Cons

Shard key

Q51. Partitioning

Difference between

Partitioning

vs

Sharding

10. Database Design
Q52. Design a URL Shortener Database
Q53. Design an E-commerce Database

Tables

Users
Orders
Products
Inventory
Q54. Design Chat Application Database
Q55. How would you model many-to-many relationships?
Q56. Soft Delete vs Hard Delete
11. Concurrency
Q57. Two users update same row.

How do you prevent data loss?

Expected

Version column
Optimistic Lock
Q58. Bank Transfer Problem

Transfer

100$

Account A

↓

Account B

How ensure consistency?

12. Golang + Database
Q59. Why use connection pooling?
Q60. How does database/sql manage connections?
Q61. What happens if rows aren't closed?

Need

defer rows.Close()
Q62. Context timeout in SQL queries
ctx, cancel := context.WithTimeout(...)

Why important?

Q63. SQL Injection

How to prevent?

Parameterized queries.

Q64. Prepared Statement

Benefits

Q65. ORM vs Raw SQL

Trade-offs.

13. Redis + Database
Q66. Why cache database queries?
Q67. Cache Aside Pattern
Q68. Write Through Cache
Q69. Cache Invalidation

Classic interview topic.

Q70. Prevent Cache Stampede
14. Scenario-Based Questions
Q71. Your API suddenly became slow.

How do you investigate?

Q72. Database CPU reaches 100%.

What do you do?

Q73. A query takes 10 seconds.

How do you optimize it?

Q74. Millions of inserts per second.

Database design?

Q75. Read-heavy system.

How would you scale?

Q76. Write-heavy system.

How would you scale?

Q77. User reports missing data after deployment.

How do you debug?

Q78. Duplicate orders are being created.

How would you fix it?

Q79. How would you migrate a database with zero downtime?
Q80. Your production database is full.

What do you do?

15. Advanced Topics (Senior)
Q81. CAP Theorem
Q82. Eventual Consistency
Q83. CQRS
Q84. Database Failover
Q85. Leader Election
Q86. Idempotency in Database Operations
Q87. Outbox Pattern

Very common in microservices.

Q88. Saga Pattern
Q89. Why PostgreSQL over MySQL?
Q90. Explain MVCC Internally
Top 20 "Must Master" Questions

If you're interviewing for a mid/senior Go backend role, prioritize these:

ACID properties
Isolation levels
Deadlocks and prevention
MVCC
WAL (Write-Ahead Logging)
B+ Tree indexes
Composite indexes and the leftmost prefix rule
EXPLAIN ANALYZE
Query optimization techniques
Transactions in PostgreSQL
Optimistic vs pessimistic locking
Database connection pooling in Go
Prepared statements and SQL injection prevention
Read replicas and replication lag
Sharding vs partitioning
Redis caching patterns (Cache Aside, Write Through)
Soft delete vs hard delete
Zero-downtime database migrations
Idempotency for APIs and database writes
Outbox pattern for reliable event publishing
