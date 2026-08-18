# Database and SQL

1. What is a database management system?
   - A DBMS stores, organizes, queries, secures, and manages concurrent access to persistent data.
2. What is the difference between relational and non-relational databases?
   - Relational databases use structured tables and relationships; non-relational systems use models such as documents, key-values, columns, or graphs.
3. What are tables, rows, columns, and schemas?
   - A table represents a collection, a row one record, a column one attribute, and a schema defines the database structure and rules.
4. What is a primary key? Can a table have more than one?
   - A primary key uniquely identifies each row and cannot be null; a table has one primary-key constraint, which may include multiple columns.
5. What is the difference between primary, candidate, alternate, and composite keys?
   - Candidate keys can uniquely identify rows; one becomes primary, unused candidates are alternate keys, and a composite key contains multiple columns.
6. What is a foreign key, and why is it useful?
   - A foreign key references a key in another table and helps enforce valid relationships and referential integrity.
7. What are database constraints?
   - Constraints such as `NOT NULL`, `UNIQUE`, `CHECK`, primary keys, and foreign keys enforce data-validity rules.
8. What is normalization, and why is it used?
   - Normalization separates data into related tables to reduce duplication and prevent update, insertion, and deletion anomalies.
9. Explain first, second, and third normal form.
   - 1NF requires atomic values, 2NF removes partial dependency on a composite key, and 3NF removes transitive dependency on non-key attributes.
10. When might denormalization be useful?
   - Deliberate duplication can reduce joins and improve read performance when that benefit outweighs added consistency and update complexity.
11. What is an index, and how does it improve query performance?
   - An index is an auxiliary search structure that locates rows without scanning the whole table, trading storage and write cost for faster reads.
12. What are the disadvantages of having too many indexes?
   - They consume space and slow inserts, updates, and deletes because every affected index must also be maintained.
13. What is the difference between clustered and non-clustered indexes?
   - A clustered index determines the table's physical or primary row order; a non-clustered index is a separate structure pointing to rows.
14. What is a database transaction?
   - A transaction groups operations into one logical unit that either commits successfully or rolls back.
15. What do the ACID properties mean?
   - Atomicity is all-or-nothing, consistency preserves rules, isolation controls concurrent visibility, and durability preserves committed changes.
16. What are dirty reads, non-repeatable reads, and phantom reads?
   - They mean reading uncommitted data, rereading a changed row, and rerunning a range query that returns a changed set of rows.
17. What are transaction isolation levels?
   - They define which concurrency anomalies are allowed, commonly ranging from Read Uncommitted through Serializable with increasing protection and cost.
18. What is a database deadlock, and how can it be handled?
   - It occurs when transactions wait cyclically for locks; the DB typically aborts one, while consistent lock order and short transactions reduce risk.
19. What is the difference between `DELETE`, `TRUNCATE`, and `DROP`?
   - `DELETE` removes selected rows, `TRUNCATE` quickly removes all rows, and `DROP` removes the database object itself.
20. What is the difference between `WHERE` and `HAVING`?
   - `WHERE` filters rows before grouping, while `HAVING` filters groups after aggregation.
21. What is the difference between `GROUP BY` and `ORDER BY`?
   - `GROUP BY` combines rows for aggregation; `ORDER BY` sorts the result.
22. Explain `INNER`, `LEFT`, `RIGHT`, and `FULL OUTER` joins.
   - Inner returns matches; left and right preserve all rows from their named side; full outer preserves rows from both sides, matching where possible.
23. What is a self join?
   - It joins a table to itself using aliases, often to represent relationships such as employees and their managers.
24. What is a subquery? How does it differ from a join?
   - A subquery nests one query inside another; a join directly combines rows from sources, though optimizers can sometimes produce equivalent plans.
25. What are aggregate functions?
   - Functions such as `COUNT`, `SUM`, `AVG`, `MIN`, and `MAX` calculate one result from a group of rows.
26. What is the difference between `UNION` and `UNION ALL`?
   - Both combine compatible result sets, but `UNION` removes duplicates while `UNION ALL` retains them and is usually faster.
27. What is `NULL`, and how should it be compared in SQL?
   - `NULL` represents missing or unknown data; test it with `IS NULL` or `IS NOT NULL`, not ordinary equality.
28. What is a view? What is a materialized view?
   - A view stores a query definition, while a materialized view stores its computed result and must be refreshed.
29. What is the difference between a stored procedure, function, and trigger?
   - Procedures perform callable routines, functions return values under DB-specific rules, and triggers run automatically on configured events.
30. How would you find the second-highest salary in a table?
   - Rank distinct salaries with `DENSE_RANK()` and select rank two, or select the maximum salary below the overall maximum.
31. How would you find duplicate rows in a table?
   - Group by the columns that define duplication and filter with `HAVING COUNT(*) > 1`.
32. How would you diagnose a slow SQL query?
   - Inspect its execution plan, indexes, row estimates, filters, joins, data volume, locking, and database resource metrics.

## Medium to Advanced

33. How does multi-version concurrency control work?
   - **Key note:** MVCC keeps row versions so readers can use a consistent snapshot while writers create new versions.
34. What is write skew, and which isolation level prevents it?
   - **Key note:** Concurrent transactions update different rows after reading one invariant; true serializable isolation prevents the anomaly.
35. What is the difference between optimistic and pessimistic locking?
   - **Key note:** Optimistic locking detects version conflicts at write time; pessimistic locking blocks competitors before modification.
36. How do database deadlocks occur, and how should applications handle them?
   - **Key note:** Transactions acquire locks in a cycle; keep transactions short/order locks consistently and safely retry the chosen victim.
37. How does the order of columns in a composite index affect its use?
   - **Key note:** B-tree queries generally need the leftmost prefix, with equality columns before range or sorting columns.
38. What is a covering index?
   - **Key note:** It contains all columns needed by a query, avoiding table access but increasing storage and write overhead.
39. Why might a query optimizer choose a full table scan over an index?
   - **Key note:** A scan can be cheaper for low-selectivity queries, small tables, stale statistics, or non-sargable predicates.
40. What makes a predicate sargable?
   - **Key note:** A searchable predicate lets an index directly locate a range; functions or casts on indexed columns often prevent it.
41. What is cardinality estimation, and why does it matter?
   - **Key note:** The optimizer predicts row counts to select joins and access paths; bad estimates produce poor plans.
42. Compare nested-loop, hash, and merge joins.
   - **Key note:** Nested loops suit small indexed inputs, hash joins suit equality, and merge joins exploit sorted inputs.
43. What is table partitioning, and how does partition pruning work?
   - **Key note:** Partitioning divides one logical table; pruning skips partitions proven irrelevant to the query predicate.
44. What is database sharding, and how should a shard key be chosen?
   - **Key note:** Sharding distributes rows across servers; choose a stable, evenly distributed key aligned with access patterns.
45. How do leader-follower replication and replication lag affect applications?
   - **Key note:** Followers scale reads but apply writes later, so clients may violate read-your-writes unless routed carefully.
46. What is change data capture?
   - **Key note:** CDC streams committed row changes from a transaction log for integration, replication, or derived views.
47. How would you perform a zero-downtime schema migration?
   - **Key note:** Use expand-and-contract phases, compatible application versions, batched backfill, validation, and delayed cleanup.
48. What are window functions, and how do they differ from `GROUP BY`?
   - **Key note:** Window functions calculate across related rows without collapsing each row into one grouped result.
49. What are common table expressions, and when can recursive CTEs help?
   - **Key note:** CTEs name query subresults; recursive CTEs traverse hierarchies or iteratively derive rows.
50. How do you design an idempotent database write?
   - **Key note:** Attach a stable operation key and enforce uniqueness while storing the effect and result atomically.
51. What is connection pooling, and how can incorrect pool sizing harm a database?
   - **Key note:** Too few connections queue work; too many increase contention and can overwhelm database capacity.
52. How would you archive or delete billions of old rows safely?
   - **Key note:** Prefer partition rotation or small indexed batches, throttle work, monitor replication, and verify retention requirements.
