# JPA / Hibernate — In-Depth Notes

---

## Table of Contents

1. [N+1 Problem & Solutions](#1-n1-problem--solutions)
2. [First & Second Level Cache](#2-first--second-level-cache)
3. [LAZY vs EAGER Loading](#3-lazy-vs-eager-loading)
4. [Optimistic vs Pessimistic Locking](#4-optimistic-vs-pessimistic-locking)
5. [@Transactional with Hibernate Session Flush Modes](#5-transactional-with-hibernate-session-flush-modes)
6. [JPQL vs Criteria API vs Native Query](#6-jpql-vs-criteria-api-vs-native-query)

---

## 1. N+1 Problem & Solutions

### What Is the N+1 Problem?

When fetching a list of N entities triggers **N additional queries** to load their associations — one query per entity.

```java
// Entity setup
@Entity
class Author {
    @Id Long id;
    String name;

    @OneToMany(mappedBy = "author", fetch = FetchType.LAZY)
    List<Book> books;
}

@Entity
class Book {
    @Id Long id;
    String title;

    @ManyToOne
    Author author;
}
```

```java
// N+1 in action
List<Author> authors = em.createQuery("SELECT a FROM Author a", Author.class)
                         .getResultList();   // 1 query: SELECT * FROM author

for (Author a : authors) {
    System.out.println(a.getBooks().size()); // N queries: SELECT * FROM book WHERE author_id = ?
}

// If there are 100 authors → 1 + 100 = 101 queries!
```

---

### Solution 1 — `JOIN FETCH`

Loads the association in a single SQL JOIN query:

```java
List<Author> authors = em.createQuery(
    "SELECT DISTINCT a FROM Author a JOIN FETCH a.books", Author.class)
    .getResultList();
// Single query: SELECT a.*, b.* FROM author a JOIN book b ON b.author_id = a.id
```

```java
// Spring Data JPA
public interface AuthorRepository extends JpaRepository<Author, Long> {

    @Query("SELECT DISTINCT a FROM Author a JOIN FETCH a.books")
    List<Author> findAllWithBooks();
}
```

**Limitation**: Cannot paginate with `JOIN FETCH` on collections (HibernateException).

---

### Solution 2 — `@BatchSize`

Hibernate loads associations in batches using `IN (...)` clauses instead of individual queries:

```java
@Entity
class Author {
    @OneToMany(mappedBy = "author", fetch = FetchType.LAZY)
    @BatchSize(size = 25)           // loads books for 25 authors at a time
    List<Book> books;
}

// 100 authors → ceil(100/25) = 4 queries instead of 100
// SELECT * FROM book WHERE author_id IN (1,2,3,...,25)
// SELECT * FROM book WHERE author_id IN (26,27,...,50)
// ...
```

Global batch size in `application.properties`:
```properties
spring.jpa.properties.hibernate.default_batch_fetch_size=25
```

---

### Solution 3 — `@EntityGraph`

Defines a fetch graph at query time without changing entity annotations:

```java
@Entity
@NamedEntityGraph(
    name = "Author.withBooks",
    attributeNodes = @NamedAttributeNode("books")
)
class Author { ... }

// In repository
public interface AuthorRepository extends JpaRepository<Author, Long> {

    @EntityGraph(attributePaths = {"books"})        // inline entity graph
    List<Author> findAll();

    @EntityGraph("Author.withBooks")                // named entity graph
    Optional<Author> findById(Long id);
}
```

---

### Solution Comparison

| Solution | SQL | Pagination Safe | Best For |
|---|---|---|---|
| `JOIN FETCH` | Single JOIN query | ❌ (collection) | Small result sets |
| `@BatchSize` | Batched IN queries | ✅ | Large lists, global config |
| `@EntityGraph` | JOIN query (like FETCH) | ❌ (collection) | Per-query control |
| DTO projection | Custom SELECT | ✅ | Read-only, flat data |

---

## 2. First & Second Level Cache

### First Level Cache (Session Cache)

- **Per-session** (per `EntityManager`) — automatic, always enabled
- Stores entities loaded within the current session
- Same `id` → same object reference within session
- Cleared when session closes or `em.clear()` is called

```java
// Both queries hit DB? No — second uses 1L cache
Author a1 = em.find(Author.class, 1L);   // SELECT from DB, stored in 1st-level cache
Author a2 = em.find(Author.class, 1L);   // returned from cache — NO SQL
System.out.println(a1 == a2);            // true — same instance

em.clear();                              // clears 1st-level cache
Author a3 = em.find(Author.class, 1L);   // hits DB again
```

---

### Second Level Cache (Shared Cache)

- **Shared across sessions** (per `SessionFactory` / `EntityManagerFactory`)
- Optional — must be explicitly enabled
- Stores entity state (not objects) — serializable
- Providers: **EhCache**, **Caffeine**, **Redis**, **Infinispan**

```xml
<!-- pom.xml -->
<dependency>
    <groupId>org.hibernate.orm</groupId>
    <artifactId>hibernate-jcache</artifactId>
</dependency>
<dependency>
    <groupId>com.github.ben-manes.caffeine</groupId>
    <artifactId>caffeine</artifactId>
</dependency>
```

```properties
# application.properties
spring.jpa.properties.hibernate.cache.use_second_level_cache=true
spring.jpa.properties.hibernate.cache.region.factory_class=jcache
spring.jpa.properties.hibernate.javax.cache.provider=com.github.benmanes.caffeine.jcache.spi.CaffeineCachingProvider
```

```java
@Entity
@Cache(usage = CacheConcurrencyStrategy.READ_WRITE)   // enable 2nd level cache for this entity
class Author { ... }

// Session 1
Author a = em.find(Author.class, 1L);   // DB hit, stored in 2nd-level cache

// Session 2 (different session, same app)
Author a2 = em.find(Author.class, 1L);  // served from 2nd-level cache — no DB
```

### Cache Concurrency Strategies

| Strategy | Use case |
|---|---|
| `READ_ONLY` | Immutable data — fastest, no locking |
| `NONSTRICT_READ_WRITE` | Rarely updated — slight stale risk |
| `READ_WRITE` | Frequently updated — uses soft locks |
| `TRANSACTIONAL` | Full transactional safety — JTA required |

### Query Cache

Caches query results (the list of IDs) in addition to entities:

```java
// Enable query cache
em.createQuery("SELECT a FROM Author a WHERE a.active = true", Author.class)
  .setHint("org.hibernate.cacheable", true)
  .getResultList();
```

```properties
spring.jpa.properties.hibernate.cache.use_query_cache=true
```

---

## 3. LAZY vs EAGER Loading

### `FetchType.LAZY`

Association is loaded **on first access** — a proxy is returned until the data is needed.

```java
@Entity
class Order {
    @Id Long id;

    @ManyToOne(fetch = FetchType.LAZY)    // default for @ManyToOne is EAGER — override!
    Customer customer;

    @OneToMany(mappedBy = "order", fetch = FetchType.LAZY)  // default LAZY
    List<OrderItem> items;
}

// Usage
Order order = em.find(Order.class, 1L);  // SELECT * FROM orders WHERE id=1
                                          // customer and items NOT loaded yet
order.getCustomer().getName();            // SELECT * FROM customer WHERE id=? (triggered)
```

**LazyInitializationException** — the classic Hibernate error:

```java
@Transactional
public Order loadOrder(Long id) {
    return orderRepository.findById(id).orElseThrow();
}  // session closes here

// Later (outside transaction):
order.getItems().size();  // ❌ LazyInitializationException — session is closed!
```

**Fix**:
```java
// Option 1 — load within transaction (use JOIN FETCH)
@Query("SELECT o FROM Order o JOIN FETCH o.items WHERE o.id = :id")
Optional<Order> findWithItems(@Param("id") Long id);

// Option 2 — use DTO projection (avoids entity graph issues)
// Option 3 — Open Session in View (not recommended for production APIs)
spring.jpa.open-in-view=false   // explicitly disable OSIV
```

---

### `FetchType.EAGER`

Association is loaded **immediately** with the owning entity — always in the same query (or extra query).

```java
@Entity
class User {
    @ManyToOne(fetch = FetchType.EAGER)   // default for @ManyToOne
    Role role;                            // loaded every time User is fetched
}
```

**Default fetch types by annotation:**

| Annotation | Default FetchType |
|---|---|
| `@ManyToOne` | `EAGER` |
| `@OneToOne` | `EAGER` |
| `@OneToMany` | `LAZY` |
| `@ManyToMany` | `LAZY` |

**Best practice**: Mark everything `LAZY` explicitly, then use `JOIN FETCH` or `@EntityGraph` per query.

---

## 4. Optimistic vs Pessimistic Locking

### Optimistic Locking

Assumes conflicts are **rare** — no DB lock held. Detects conflicts at commit time using a **version field**.

```java
@Entity
class BankAccount {
    @Id Long id;
    BigDecimal balance;

    @Version                   // Hibernate manages this — increments on every update
    Integer version;
}

// Thread 1 reads account: version=1, balance=1000
// Thread 2 reads account: version=1, balance=1000

// Thread 1 updates: UPDATE bank_account SET balance=900, version=2 WHERE id=1 AND version=1  ✅
// Thread 2 updates: UPDATE bank_account SET balance=800, version=2 WHERE id=1 AND version=1  ❌
// → 0 rows affected → Hibernate throws OptimisticLockException (or StaleObjectStateException)
```

```java
// Handling OptimisticLockException
@Transactional
@Retryable(value = OptimisticLockingFailureException.class, maxAttempts = 3)
public void debit(Long accountId, BigDecimal amount) {
    BankAccount account = repository.findById(accountId).orElseThrow();
    account.setBalance(account.getBalance().subtract(amount));
    repository.save(account);  // may throw OptimisticLockException
}
```

---

### Pessimistic Locking

Assumes conflicts are **likely** — acquires a DB-level lock immediately, blocking other readers/writers.

```java
// PESSIMISTIC_READ — shared lock (others can read, not write)
BankAccount account = em.find(BankAccount.class, id,
    LockModeType.PESSIMISTIC_READ);

// PESSIMISTIC_WRITE — exclusive lock (others cannot read or write)
BankAccount account = em.find(BankAccount.class, id,
    LockModeType.PESSIMISTIC_WRITE);
// SQL: SELECT * FROM bank_account WHERE id=? FOR UPDATE

// Spring Data
public interface AccountRepository extends JpaRepository<BankAccount, Long> {
    @Lock(LockModeType.PESSIMISTIC_WRITE)
    @Query("SELECT a FROM BankAccount a WHERE a.id = :id")
    Optional<BankAccount> findByIdForUpdate(@Param("id") Long id);
}
```

### Comparison

| | Optimistic | Pessimistic |
|---|---|---|
| Lock held | No DB lock | DB row/table lock |
| Conflict detection | At commit | Upfront (blocking) |
| Performance | Better (no lock wait) | Lower (contention) |
| Failure handling | `OptimisticLockException` | Waits (or timeout) |
| Best for | Low-contention reads | High-contention writes |
| Implementation | `@Version` field | `LockModeType` |

---

## 5. @Transactional with Hibernate Session Flush Modes

### What Is Flushing?

**Flushing** synchronizes the in-memory state of the `EntityManager` (persistence context) with the database — it issues pending SQL statements.

```
EntityManager (persistence context)
  ├─ author.setName("Alice")    ← in-memory change, not yet in DB
  ├─ em.persist(newBook)        ← queued INSERT
  └─ em.remove(oldBook)         ← queued DELETE

Flush → sends all pending SQL to DB (within transaction — not yet committed)
Commit → makes changes permanent
```

---

### Flush Modes

| Mode | Flush Happens When |
|---|---|
| `AUTO` (default) | Before query execution (if pending changes affect query results) + before commit |
| `COMMIT` | Only before transaction commit |
| `ALWAYS` | Before every query execution |
| `MANUAL` | Only when `em.flush()` is called explicitly |

```java
// AUTO — default (recommended)
@Transactional
public void updateAndQuery() {
    author.setName("Updated");              // dirty entity in context
    List<Author> result = em.createQuery(   // AUTO flushes first (name change affects query)
        "SELECT a FROM Author a WHERE a.name = 'Updated'", Author.class)
        .getResultList();
    // result contains the updated author
}

// COMMIT — useful for read-heavy operations to avoid unnecessary flushes
@Transactional
public void bulkRead() {
    em.setFlushMode(FlushModeType.COMMIT);  // skip mid-query flushes
    for (int i = 0; i < 1000; i++) {
        process(repository.findByCategory(i));   // no flush before each query
    }
}

// MANUAL — for batch processing control
@Transactional
public void batchInsert(List<Entity> items) {
    for (int i = 0; i < items.size(); i++) {
        em.persist(items.get(i));
        if (i % 50 == 0) {
            em.flush();   // explicit flush every 50
            em.clear();   // clear cache to prevent OutOfMemoryError
        }
    }
}
```

```java
// Set flush mode via Spring @Transactional
@Transactional
public void process() {
    Session session = em.unwrap(Session.class);
    session.setHibernateFlushMode(FlushMode.MANUAL);
    // ...
    session.flush();
}
```

---

## 6. JPQL vs Criteria API vs Native Query

### JPQL (Java Persistence Query Language)

Object-oriented query language — operates on **entity names and fields**, not table/column names.

```java
// Basic query
List<Author> authors = em.createQuery(
    "SELECT a FROM Author a WHERE a.age > :minAge ORDER BY a.name",
    Author.class)
    .setParameter("minAge", 25)
    .getResultList();

// JOIN
List<Object[]> result = em.createQuery(
    "SELECT a.name, b.title FROM Author a JOIN a.books b WHERE b.published = true")
    .getResultList();

// Aggregate
Long count = em.createQuery(
    "SELECT COUNT(a) FROM Author a WHERE a.active = true", Long.class)
    .getSingleResult();

// UPDATE / DELETE (bulk operations — bypass persistence context)
em.createQuery("UPDATE Author a SET a.active = false WHERE a.lastLogin < :date")
  .setParameter("date", LocalDate.now().minusYears(1))
  .executeUpdate();
```

```java
// Spring Data JPA
public interface AuthorRepository extends JpaRepository<Author, Long> {

    @Query("SELECT a FROM Author a WHERE a.name LIKE :name%")
    List<Author> findByNameStartingWith(@Param("name") String name);

    // DTO projection
    @Query("SELECT new com.example.AuthorDto(a.id, a.name) FROM Author a")
    List<AuthorDto> findAllAsDto();
}
```

---

### Criteria API

Type-safe, programmatic query building — ideal for **dynamic queries** where conditions vary at runtime.

```java
CriteriaBuilder cb = em.getCriteriaBuilder();
CriteriaQuery<Author> cq = cb.createQuery(Author.class);
Root<Author> root = cq.from(Author.class);

// Build predicates dynamically
List<Predicate> predicates = new ArrayList<>();

if (name != null) {
    predicates.add(cb.like(root.get("name"), name + "%"));
}
if (minAge != null) {
    predicates.add(cb.greaterThan(root.get("age"), minAge));
}
if (active != null) {
    predicates.add(cb.equal(root.get("active"), active));
}

cq.where(predicates.toArray(new Predicate[0]))
  .orderBy(cb.asc(root.get("name")));

List<Author> authors = em.createQuery(cq).getResultList();
```

**With JPA Metamodel** (type-safe field references, generated by annotation processor):

```java
// Author_ is auto-generated
predicates.add(cb.like(root.get(Author_.name), name + "%"));    // compile-time safe
predicates.add(cb.greaterThan(root.get(Author_.age), minAge));
```

**Spring Data Specifications** (cleaner Criteria API wrapper):

```java
public class AuthorSpecs {
    public static Specification<Author> hasName(String name) {
        return (root, query, cb) -> cb.like(root.get("name"), name + "%");
    }
    public static Specification<Author> isActive() {
        return (root, query, cb) -> cb.isTrue(root.get("active"));
    }
}

// Usage
List<Author> authors = authorRepository.findAll(
    AuthorSpecs.hasName("John").and(AuthorSpecs.isActive())
);
```

---

### Native Query

Raw SQL — when JPQL/Criteria is insufficient (DB-specific functions, complex CTEs, stored procedures):

```java
// Simple native query
List<Object[]> result = em.createNativeQuery(
    "SELECT id, name, age FROM author WHERE age > ?1")
    .setParameter(1, 25)
    .getResultList();

// Map to entity
List<Author> authors = em.createNativeQuery(
    "SELECT * FROM author WHERE name ILIKE :name",  // PostgreSQL ILIKE
    Author.class)
    .setParameter("name", "%john%")
    .getResultList();

// Map to DTO via @SqlResultSetMapping
@SqlResultSetMapping(
    name = "AuthorSummaryMapping",
    classes = @ConstructorResult(
        targetClass = AuthorSummaryDto.class,
        columns = {
            @ColumnResult(name = "id",   type = Long.class),
            @ColumnResult(name = "name", type = String.class),
            @ColumnResult(name = "book_count", type = Long.class)
        }
    )
)
```

```java
// Spring Data
public interface AuthorRepository extends JpaRepository<Author, Long> {

    @Query(value = "SELECT * FROM author WHERE name ILIKE :name", nativeQuery = true)
    List<Author> findByNameNative(@Param("name") String name);
}
```

---

### Comparison

| | JPQL | Criteria API | Native Query |
|---|---|---|---|
| Type-safe | Partial (string-based) | ✅ Full (with metamodel) | ❌ |
| Dynamic queries | Awkward | ✅ Excellent | Possible (string concat) |
| DB-portable | ✅ | ✅ | ❌ DB-specific |
| Complex SQL (CTEs, window functions) | ❌ | ❌ | ✅ |
| Performance tuning | Limited | Limited | Full control |
| Readable | ✅ | Verbose | ✅ |
| Use when | Simple/medium queries | Dynamic filters | DB-specific features |
