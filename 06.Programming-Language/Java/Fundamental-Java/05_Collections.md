# Collections Framework

---

## Table of Contents

1. [Collections Overview](#1-collections-overview)
2. [List — ArrayList vs LinkedList](#2-list--arraylist-vs-linkedlist)
3. [Set — HashSet vs LinkedHashSet vs TreeSet](#3-set--hashset-vs-linkedhashset-vs-treeset)
4. [Map — HashMap vs LinkedHashMap vs TreeMap](#4-map--hashmap-vs-linkedhashmap-vs-treemap)
5. [HashMap vs Hashtable vs ConcurrentHashMap](#5-hashmap-vs-hashtable-vs-concurrenthashmap)
6. [Queue & Deque](#6-queue--deque)
7. [HashMap Internals](#7-hashmap-internals)
8. [`equals()` and `hashCode()` Contract](#8-equals-and-hashcode-contract)
9. [Comparable vs Comparator](#9-comparable-vs-comparator)
10. [Iterator vs ListIterator](#10-iterator-vs-listiterator)
11. [Fail-Fast vs Fail-Safe](#11-fail-fast-vs-fail-safe)

---

## 1. Collections Overview

```
java.lang.Iterable
└── java.util.Collection
    ├── List (ordered, duplicates allowed)
    │   ├── ArrayList       — dynamic array, O(1) get
    │   ├── LinkedList      — doubly linked, O(1) insert/delete at ends
    │   └── Vector          — legacy, synchronized
    │
    ├── Set (no duplicates)
    │   ├── HashSet         — no order, O(1) ops
    │   ├── LinkedHashSet   — insertion order
    │   └── TreeSet         — sorted order, O(log n) ops
    │
    └── Queue (FIFO / priority)
        ├── LinkedList      — FIFO queue
        ├── PriorityQueue   — heap-based, sorted by priority
        └── ArrayDeque      — efficient double-ended queue

java.util.Map (key-value, NOT extends Collection)
    ├── HashMap             — no order, O(1) ops
    ├── LinkedHashMap       — insertion/access order
    ├── TreeMap             — sorted by key, O(log n) ops
    └── Hashtable           — legacy, synchronized
```

---

## 2. List — ArrayList vs LinkedList

```java
// ── ARRAYLIST — backed by Object[] array ──
List<String> arrayList = new ArrayList<>();
arrayList.add("Alice");        // O(1) amortized — may resize
arrayList.add("Bob");
arrayList.add(0, "Charlie");   // O(n) — shifts all elements right
arrayList.get(1);              // O(1) — direct index access
arrayList.remove(1);           // O(n) — shifts elements left
arrayList.contains("Alice");   // O(n) — linear scan

// Initial capacity — avoids resizing if size is known
List<Integer> large = new ArrayList<>(10_000);

// ── LINKEDLIST — doubly linked list ──
LinkedList<String> linkedList = new LinkedList<>();
linkedList.addFirst("Alice");   // O(1)
linkedList.addLast("Bob");      // O(1)
linkedList.add(1, "Charlie");   // O(n) to find position, O(1) to insert
linkedList.get(1);              // O(n) — must traverse from head
linkedList.removeFirst();       // O(1)
linkedList.removeLast();        // O(1)

// LinkedList also implements Deque — use as stack or queue
linkedList.push("first");      // stack push (addFirst)
linkedList.pop();              // stack pop (removeFirst)
linkedList.offer("last");      // queue enqueue (addLast)
linkedList.poll();             // queue dequeue (removeFirst)
```

| Operation | ArrayList | LinkedList |
|---|---|---|
| `get(i)` | **O(1)** | O(n) |
| `add` at end | O(1) amortized | **O(1)** |
| `add` at middle | O(n) | O(n) to find + O(1) to insert |
| `remove` at end | O(1) | **O(1)** |
| `remove` at middle | O(n) | O(n) to find + O(1) to remove |
| Memory | Less (compact array) | More (node + 2 pointers per element) |

> **Interview Q: When would you use LinkedList over ArrayList?**  
> Very rarely in practice. Use `LinkedList` only when you need **frequent insertions/deletions at both ends** (acts as Deque/stack/queue) and **never need random access**. `ArrayList` is almost always faster due to **CPU cache locality** — its elements are in contiguous memory, making iteration significantly faster than LinkedList's pointer-chasing. For most use cases, prefer `ArrayList` or `ArrayDeque`.

---

## 3. Set — HashSet vs LinkedHashSet vs TreeSet

```java
// ── HASHSET — no order guarantee ──
Set<String> hashSet = new HashSet<>();
hashSet.add("banana");
hashSet.add("apple");
hashSet.add("cherry");
hashSet.add("apple");         // duplicate — ignored
System.out.println(hashSet);  // [banana, cherry, apple] — unpredictable order

// ── LINKEDHASHSET — maintains INSERTION order ──
Set<String> linkedSet = new LinkedHashSet<>();
linkedSet.add("banana");
linkedSet.add("apple");
linkedSet.add("cherry");
System.out.println(linkedSet);  // [banana, apple, cherry] — insertion order

// ── TREESET — sorted (natural order or custom Comparator) ──
Set<String> treeSet = new TreeSet<>();
treeSet.add("banana");
treeSet.add("apple");
treeSet.add("cherry");
System.out.println(treeSet);    // [apple, banana, cherry] — sorted

// TreeSet extra operations:
treeSet.first();           // "apple"
treeSet.last();            // "cherry"
treeSet.headSet("cherry"); // ["apple", "banana"] — exclusive
treeSet.tailSet("banana"); // ["banana", "cherry"] — inclusive
treeSet.floor("bo");       // "banana" — greatest <= "bo"
treeSet.ceiling("bo");     // "cherry" — smallest >= "bo"

// Custom sort in TreeSet
Set<String> byLength = new TreeSet<>(Comparator.comparingInt(String::length)
                                               .thenComparing(Comparator.naturalOrder()));
byLength.addAll(List.of("fig", "apple", "kiwi", "date"));
System.out.println(byLength);  // [fig, date, kiwi, apple] — by length then alpha
```

| | HashSet | LinkedHashSet | TreeSet |
|---|---|---|---|
| Ordering | None | Insertion order | Sorted (natural/comparator) |
| `add`/`remove`/`contains` | **O(1)** | **O(1)** | O(log n) |
| Null elements | ✅ One | ✅ One | ❌ (NullPointerException) |
| Backed by | `HashMap` | `LinkedHashMap` | Red-Black Tree |

---

## 4. Map — HashMap vs LinkedHashMap vs TreeMap

```java
// ── HASHMAP — no order ──
Map<String, Integer> scores = new HashMap<>();
scores.put("Alice", 95);
scores.put("Bob", 87);
scores.put("Charlie", 92);
scores.put("Alice", 98);          // replaces — key must be unique

scores.get("Bob");                // 87
scores.getOrDefault("Dave", 0);  // 0 — safe get
scores.putIfAbsent("Eve", 70);   // adds only if key not present
scores.computeIfAbsent("Frank", k -> k.length() * 10);  // 50

// Iterate
for (Map.Entry<String, Integer> entry : scores.entrySet()) {
    System.out.println(entry.getKey() + ": " + entry.getValue());
}
scores.forEach((k, v) -> System.out.println(k + ": " + v));

// ── LINKEDHASHMAP — insertion order preserved ──
Map<String, Integer> ordered = new LinkedHashMap<>();
ordered.put("first", 1);
ordered.put("second", 2);
ordered.put("third", 3);
// Iterates in: first, second, third

// Access-ordered LinkedHashMap — useful for LRU Cache
Map<String, Integer> lru = new LinkedHashMap<>(16, 0.75f, true) {
    @Override
    protected boolean removeEldestEntry(Map.Entry<String, Integer> eldest) {
        return size() > 3;    // max 3 entries
    }
};

// ── TREEMAP — sorted by key ──
Map<String, Integer> sorted = new TreeMap<>();
sorted.put("banana", 2);
sorted.put("apple", 5);
sorted.put("cherry", 3);
// Iterates in: apple, banana, cherry (alphabetical)

sorted.firstKey();              // "apple"
sorted.lastKey();               // "cherry"
sorted.headMap("cherry");       // {apple=5, banana=2}
sorted.floorKey("bo");          // "banana"
```

| | HashMap | LinkedHashMap | TreeMap |
|---|---|---|---|
| Order | None | Insertion (or access) | Sorted by key |
| `get`/`put` | **O(1)** | **O(1)** | O(log n) |
| Null keys | ✅ One | ✅ One | ❌ |
| Use case | General purpose | LRU cache, predictable iteration | Range queries, sorted access |

---

## 5. HashMap vs Hashtable vs ConcurrentHashMap

```java
// ── HASHMAP (Java 1.2) ──
Map<String, Integer> hashMap = new HashMap<>();
// - NOT thread-safe
// - Allows ONE null key, multiple null values
// - Fast: O(1)

// ── HASHTABLE (Java 1.0 — legacy) ──
Map<String, Integer> hashtable = new Hashtable<>();
// - Thread-safe (every method synchronized on the whole object)
// - No null keys or values → NullPointerException
// - Slow in concurrent scenarios (one lock for everything)
// - AVOID — use ConcurrentHashMap instead

// ── CONCURRENTHASHMAP (Java 5) ──
Map<String, Integer> concurrent = new ConcurrentHashMap<>();
// - Thread-safe
// - No null keys or values
// - Segment/bucket-level locking (Java 8+: CAS + synchronized on bucket)
// - Better performance than Hashtable under concurrent access
// - Atomic operations: putIfAbsent, computeIfAbsent, merge

concurrent.putIfAbsent("counter", 0);
concurrent.merge("counter", 1, Integer::sum);  // atomic increment-like
concurrent.compute("counter", (k, v) -> v == null ? 1 : v + 1);
```

| | HashMap | Hashtable | ConcurrentHashMap |
|---|---|---|---|
| Thread-safe | ❌ | ✅ (full sync) | ✅ (segment/bucket sync) |
| Null key | ✅ One | ❌ | ❌ |
| Null value | ✅ | ❌ | ❌ |
| Performance | Fastest (single thread) | Slowest | Best (multi-thread) |
| Legacy | Java 1.2 | Java 1.0 (avoid) | Java 5 |

> **Interview Q: Why should you use ConcurrentHashMap instead of Hashtable?**  
> `Hashtable` synchronizes every method on the **entire object** — only one thread can operate on it at a time, even for reads on different buckets. `ConcurrentHashMap` uses **bucket-level locking** (Java 8+: CAS operations + synchronized only on the specific bucket being modified), allowing **multiple threads to read and write different buckets simultaneously** — far better throughput. Additionally, `ConcurrentHashMap` provides atomic composite operations like `putIfAbsent`, `computeIfAbsent`, and `merge`.

---

## 6. Queue & Deque

```java
// ── QUEUE (FIFO) ──
Queue<String> queue = new LinkedList<>();
queue.offer("Alice");      // enqueue — preferred over add() (no exception on full)
queue.offer("Bob");
queue.offer("Charlie");

queue.peek();              // "Alice" — see front without removing
queue.poll();              // "Alice" — remove and return front
queue.size();              // 2

// ── PRIORITYQUEUE — min-heap by default ──
PriorityQueue<Integer> pq = new PriorityQueue<>();
pq.offer(30);
pq.offer(10);
pq.offer(20);
System.out.println(pq.poll());   // 10 — smallest first

// Max-heap
PriorityQueue<Integer> maxPq = new PriorityQueue<>(Collections.reverseOrder());
maxPq.offer(30); maxPq.offer(10); maxPq.offer(20);
System.out.println(maxPq.poll()); // 30 — largest first

// ── DEQUE (double-ended queue) — stack or queue ──
Deque<String> deque = new ArrayDeque<>();
deque.addFirst("A");    // push front
deque.addLast("B");     // push back
deque.addFirst("Z");    // [Z, A, B]

deque.peekFirst();      // "Z"
deque.peekLast();       // "B"
deque.pollFirst();      // "Z" — remove front
deque.pollLast();       // "B" — remove back

// Use as STACK (LIFO)
Deque<String> stack = new ArrayDeque<>();  // PREFER over Stack class
stack.push("first");     // addFirst
stack.push("second");
System.out.println(stack.pop());   // "second" — LIFO
```

> **Interview Q: Why use `ArrayDeque` as a stack instead of `Stack`?**  
> The `Stack` class extends `Vector`, which is **synchronized on every method** (unnecessary overhead for single-threaded use). It's also part of the legacy hierarchy. `ArrayDeque` is **faster** (no synchronization), doesn't have the legacy baggage, and is the **recommended replacement** according to the Java documentation itself. For a thread-safe stack, use `Deque<T> stack = new ConcurrentLinkedDeque<>()` or use `Deque` + explicit synchronization.

---

## 7. HashMap Internals

```
HashMap<K, V> internals (Java 8+)
─────────────────────────────────────
Initial capacity: 16 buckets
Load factor: 0.75 (resize when 75% full)
After resize: doubles capacity

Bucket array: Object[] table (size = capacity)

For each key-value pair:
  1. Compute hash: hash = key.hashCode() ^ (h >>> 16)
  2. Bucket index: index = hash & (capacity - 1)
  3. Store in bucket

Each bucket:
  - Empty: null
  - 1-8 entries: LinkedList (singly linked Node<K,V>)
  - >8 entries + capacity >= 64: converted to TreeMap (Red-Black Tree)
  - TreeMap shrinks back to LinkedList when entries drop to 6
```

```java
// What happens on put(key, value):
// 1. Compute hash of key
// 2. Find bucket: index = hash % capacity
// 3. If bucket empty → insert new Node
// 4. If bucket not empty:
//    a. Check each node: if keys.equals(existingKey) → REPLACE value
//    b. If no match → APPEND to chain
// 5. If size > capacity * loadFactor → RESIZE (double capacity, rehash all)

// What happens on get(key):
// 1. Compute hash of key
// 2. Find bucket
// 3. Traverse chain/tree: compare with equals()
// 4. Return value if found, null otherwise

// Why both hashCode AND equals?
// hashCode → finds the RIGHT BUCKET (fast O(1))
// equals  → finds the RIGHT KEY within the bucket
```

> **Interview Q: What happens when two keys have the same `hashCode` in a HashMap?**  
> This is a **hash collision**. Both entries go into the **same bucket**. Before Java 8, the bucket stored a linked list of all entries with that hash — worst case O(n) for get/put. From **Java 8**, when a bucket's linked list exceeds **8 entries** (and total capacity is ≥ 64), it converts to a **Red-Black Tree**, making worst-case O(log n). When entries drop back to 6, it converts back to a linked list.

---

## 8. `equals()` and `hashCode()` Contract

```java
// THE CONTRACT:
// 1. If a.equals(b) is true → a.hashCode() == b.hashCode() MUST be true
// 2. If a.hashCode() == b.hashCode() → a.equals(b) MAY or MAY NOT be true (collision)
// 3. If !a.equals(b) → hashCode() can be equal or not (collisions OK)

// RULE: Always override BOTH together, NEVER just one

class Point {
    int x, y;

    Point(int x, int y) { this.x = x; this.y = y; }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;             // same reference
        if (!(o instanceof Point)) return false; // null-safe & type check
        Point p = (Point) o;
        return x == p.x && y == p.y;
    }

    @Override
    public int hashCode() {
        return Objects.hash(x, y);   // uses all fields in equals()
    }
}

// What breaks if you override equals() but NOT hashCode():
Set<Point> set = new HashSet<>();
set.add(new Point(1, 2));
set.contains(new Point(1, 2));   // FALSE! — different hashCode → different bucket

// What breaks if you override hashCode() but NOT equals():
// Two "equal" objects still treated as different keys
Map<Point, String> map = new HashMap<>();
Point p1 = new Point(1, 2);
Point p2 = new Point(1, 2);
map.put(p1, "origin");
map.get(p2);   // null — p1.equals(p2) is false (Object.equals checks reference)
```

> **Interview Q: What is the contract between `equals()` and `hashCode()`?**  
> If two objects are **equal** (via `equals()`), they **must** have the same `hashCode`. The reverse is not required — equal hashCodes don't mean equal objects (collision is allowed). You must override **both** together. If you override `equals()` without `hashCode()`, HashMap/HashSet will **fail silently** — two logically equal objects will land in different buckets and the map won't find them.

---

## 9. Comparable vs Comparator

```java
// ── COMPARABLE — natural ordering, intrinsic to the class ──
class Student implements Comparable<Student> {
    String name;
    double gpa;

    Student(String name, double gpa) {
        this.name = name;
        this.gpa = gpa;
    }

    @Override
    public int compareTo(Student other) {
        // Negative: this < other
        // Zero: this == other
        // Positive: this > other
        return Double.compare(other.gpa, this.gpa);   // descending by GPA
    }
}

List<Student> students = new ArrayList<>(List.of(
    new Student("Alice", 3.9),
    new Student("Bob", 3.5),
    new Student("Charlie", 3.7)
));
Collections.sort(students);   // uses Comparable.compareTo()
// Result: Alice(3.9), Charlie(3.7), Bob(3.5)

// ── COMPARATOR — external ordering, separate from class ──
// Useful when: (a) class doesn't implement Comparable, (b) multiple sort orders needed

Comparator<Student> byName = Comparator.comparing(s -> s.name);
Comparator<Student> byGpa  = Comparator.comparingDouble((Student s) -> s.gpa).reversed();
Comparator<Student> byNameThenGpa = byName.thenComparing(byGpa);

students.sort(byName);           // alphabetical
students.sort(byGpa);            // descending GPA
students.sort(byNameThenGpa);    // alpha, then descending GPA

// Inline
students.sort(Comparator.comparing(s -> s.name));

// TreeSet/TreeMap with custom ordering
TreeSet<Student> byGpaSet = new TreeSet<>(Comparator.comparingDouble(s -> s.gpa));
```

| | Comparable | Comparator |
|---|---|---|
| Package | `java.lang` | `java.util` |
| Method | `compareTo(T other)` | `compare(T o1, T o2)` |
| Defined in | The class itself | Separate class or lambda |
| Sort orders | One (natural) | Many (external) |
| Modifies class | Yes | No |
| Used by | `Collections.sort()`, `TreeSet` | `Collections.sort()`, `Arrays.sort()` |

> **Interview Q: What is the difference between Comparable and Comparator?**  
> `Comparable` defines the **natural ordering** of objects — implemented inside the class itself (`implements Comparable<T>`, overrides `compareTo`). `Comparator` is an **external ordering** — defined outside the class (separate class, anonymous class, or lambda). Use `Comparable` for the default/natural sort. Use `Comparator` when you need multiple different orderings, or when you can't modify the class.

---

## 10. Iterator vs ListIterator

```java
List<String> list = new ArrayList<>(List.of("a", "b", "c", "d"));

// ── ITERATOR — forward-only, works on any Collection ──
Iterator<String> it = list.iterator();
while (it.hasNext()) {
    String s = it.next();
    if (s.equals("b")) {
        it.remove();   // SAFE removal during iteration — no ConcurrentModificationException
    }
}
System.out.println(list);   // [a, c, d]

// ── LISTITERATOR — bidirectional, List only, supports add/set ──
ListIterator<String> lit = list.listIterator();

// Forward
while (lit.hasNext()) {
    int index = lit.nextIndex();
    String s = lit.next();
    if (s.equals("c")) {
        lit.set("C");       // replace current element
        lit.add("c+");      // insert after current
    }
}

// Backward
while (lit.hasPrevious()) {
    System.out.print(lit.previous() + " ");
}
```

| | Iterator | ListIterator |
|---|---|---|
| Direction | Forward only | Forward and backward |
| Works on | Any `Collection` | `List` only |
| Methods | `hasNext`, `next`, `remove` | + `hasPrevious`, `previous`, `add`, `set`, `nextIndex`, `previousIndex` |
| Modify during iteration | `remove()` only | `remove()`, `add()`, `set()` |

---

## 11. Fail-Fast vs Fail-Safe

```java
// ── FAIL-FAST — throws ConcurrentModificationException if modified during iteration ──
List<String> list = new ArrayList<>(List.of("a", "b", "c"));
Iterator<String> it = list.iterator();
list.add("d");              // structural modification!
it.next();                  // ❌ ConcurrentModificationException

// How it works: ArrayList maintains a 'modCount' (modification count)
// Iterator checks modCount == expectedModCount on each next()
// If they differ → CME

// Fail-fast collections: ArrayList, HashMap, HashSet, etc.
// Use iterator.remove() or removeIf() to safely modify during iteration

// ── FAIL-SAFE — works on a copy, no exception ──
CopyOnWriteArrayList<String> cowList = new CopyOnWriteArrayList<>(List.of("a", "b", "c"));
for (String s : cowList) {
    cowList.add("x");        // ✅ no exception — iterates original snapshot
}
System.out.println(cowList); // [a, b, c, x, x, x]

// ConcurrentHashMap is also fail-safe
ConcurrentHashMap<String, Integer> chm = new ConcurrentHashMap<>();
chm.put("a", 1); chm.put("b", 2);
for (String key : chm.keySet()) {
    chm.put("c", 3);   // ✅ allowed — but new entries may or may not appear in iteration
}
```

| | Fail-Fast | Fail-Safe |
|---|---|---|
| Behavior on modification | Throws `ConcurrentModificationException` | No exception |
| Iterates on | Actual collection | Copy / snapshot |
| Memory | Less (no copy) | More (copy overhead) |
| Examples | `ArrayList`, `HashMap`, `HashSet` | `CopyOnWriteArrayList`, `ConcurrentHashMap` |
| Reflects new changes | Yes | Not necessarily |

> **Interview Q: What is the difference between fail-fast and fail-safe iterators?**  
> **Fail-fast** iterators detect structural modifications (add/remove) during iteration by tracking a `modCount`, and immediately throw `ConcurrentModificationException`. They iterate directly on the collection with no copy overhead. **Fail-safe** iterators work on a **snapshot or clone** — modifications to the original don't affect the iteration, so no exception is thrown. The trade-off: fail-safe uses more memory and the iterator may not see updates made during traversal.
