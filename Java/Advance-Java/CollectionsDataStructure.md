# Collections & Data Structures — In-Depth Notes

---

## Table of Contents

1. [HashMap Internals](#1-hashmap-internals)
2. [ConcurrentHashMap vs Collections.synchronizedMap](#2-concurrenthashmap-vs-collectionssynchronizedmap)
3. [LinkedHashMap, TreeMap, LinkedList vs ArrayList](#3-linkedhashmap-treemap-linkedlist-vs-arraylist)
4. [Comparable vs Comparator](#4-comparable-vs-comparator)
5. [Fail-fast vs Fail-safe Iterators](#5-fail-fast-vs-fail-safe-iterators)

---

## 1. HashMap Internals

### Core Structure

A `HashMap` is backed by an **array of buckets** (`Node<K,V>[] table`).  
Each bucket holds a **linked list** of entries; from Java 8 onwards it converts to a **Red-Black Tree** when a bucket gets too large.

```
table[]
 [0] → null
 [1] → Node("apple", 1) → Node("grape", 4) → null   (collision)
 [2] → Node("banana", 2)
 [3] → null
 [4] → Node("cherry", 3)
 ...
```

### How `put(K key, V value)` Works

```java
Map<String, Integer> map = new HashMap<>();
map.put("apple", 1);
```

**Step-by-step:**

1. **Hash the key** — `hash = key.hashCode()` then spread bits: `(h ^ (h >>> 16))`
2. **Find bucket index** — `index = hash & (capacity - 1)`
3. **Check bucket:**
   - Empty → create new `Node` and place it
   - Non-empty → walk the list/tree, comparing hash + `equals()`
     - Match found → **update** value
     - No match → **append** new node

```java
// Internally (simplified):
int hash = hash(key.hashCode());          // spread hash
int index = hash & (table.length - 1);   // bucket index
Node<K,V> existing = table[index];

if (existing == null) {
    table[index] = new Node<>(hash, key, value, null);
} else {
    // traverse list / tree, check equals()
}
```

### Treeification

When a single bucket's linked list length reaches **≥ 8** AND the table size is **≥ 64**, the list is converted to a **Red-Black Tree** (TreeNode).  
This improves worst-case lookup from **O(n)** to **O(log n)**.

```
Linked list (length < 8):         Red-Black Tree (length ≥ 8):
  A → B → C → D → E → F → G → H     [D]
                                    /     \
                                  [B]     [F]
                                 / \     / \
                               [A] [C] [E] [G]
                                           \
                                           [H]
```

If the bucket size drops below **6** (after removals), the tree is converted back to a linked list.

### Resize (Rehashing)

- Default initial capacity: **16**; load factor: **0.75**
- When `size > capacity × loadFactor` → resize: **double the table**, rehash all entries

```java
// Default construction
Map<String, Integer> map = new HashMap<>();          // capacity=16, LF=0.75
// Resizes at 16 × 0.75 = 12 entries

// Custom
Map<String, Integer> map = new HashMap<>(32, 0.5f); // capacity=32, LF=0.5
```

Resize is expensive — if you know approximate size, pre-size to avoid it:
```java
// Pre-size to hold 100 elements without resizing
Map<String, Integer> map = new HashMap<>(128);  // 100 / 0.75 ≈ 134, next power of 2 = 128
```

### Key Points Summary

| Property | Value |
|---|---|
| Default capacity | 16 |
| Default load factor | 0.75 |
| Treeification threshold | 8 nodes in bucket + table size ≥ 64 |
| Untreeify threshold | 6 nodes |
| Null keys allowed | Yes (1 null key) |
| Thread-safe | No |
| Order | Not guaranteed |

---

## 2. ConcurrentHashMap vs Collections.synchronizedMap

### `Collections.synchronizedMap`

Wraps any map with a **single mutex lock** on the entire map. Every operation acquires a lock on the wrapper object.

```java
Map<String, Integer> syncMap = Collections.synchronizedMap(new HashMap<>());

syncMap.put("a", 1);   // locks entire map
syncMap.get("a");      // locks entire map

// IMPORTANT: iteration must be manually synchronized
synchronized (syncMap) {
    for (Map.Entry<String, Integer> e : syncMap.entrySet()) {
        System.out.println(e.getKey() + "=" + e.getValue());
    }
}
```

**Problem**: Only one thread can read or write at a time → poor throughput under contention.

---

### `ConcurrentHashMap`

Uses **segment-level / bucket-level locking** (Java 8+: CAS + `synchronized` on individual bucket heads).  
Multiple threads can read and write to **different buckets simultaneously**.

```java
ConcurrentHashMap<String, Integer> map = new ConcurrentHashMap<>();

// Thread-safe put / get — no external locking needed
map.put("a", 1);
map.get("a");

// Atomic compound operations
map.putIfAbsent("b", 2);
map.computeIfAbsent("c", k -> k.length());
map.merge("a", 1, Integer::sum);   // atomic increment

// Safe iteration — no ConcurrentModificationException
for (Map.Entry<String, Integer> e : map.entrySet()) {
    System.out.println(e.getKey());
}
```

### Head-to-Head Comparison

| Feature | `synchronizedMap` | `ConcurrentHashMap` |
|---|---|---|
| Locking granularity | Entire map | Per-bucket (fine-grained) |
| Read performance | Locked | Lock-free (volatile reads) |
| Write performance | Single lock | Concurrent on different buckets |
| Null keys/values | Depends on wrapped map | **Not allowed** |
| Iteration | Must sync manually | Weakly consistent (no CME) |
| Atomic ops (`putIfAbsent`) | No | Yes |
| Recommended | Legacy / simple cases | High-concurrency production code |

```java
// ConcurrentHashMap — null throws NullPointerException
ConcurrentHashMap<String, String> map = new ConcurrentHashMap<>();
map.put(null, "v");   // throws NullPointerException
map.put("k", null);   // throws NullPointerException
```

---

## 3. LinkedHashMap, TreeMap, LinkedList vs ArrayList

### LinkedHashMap

Extends `HashMap` but maintains a **doubly-linked list** across all entries to preserve **insertion order** (or access order).

```java
// Insertion-order (default)
Map<String, Integer> linked = new LinkedHashMap<>();
linked.put("banana", 2);
linked.put("apple", 1);
linked.put("cherry", 3);
System.out.println(linked.keySet()); // [banana, apple, cherry]

// Access-order — use as LRU cache
Map<String, Integer> lru = new LinkedHashMap<>(16, 0.75f, true) {
    @Override
    protected boolean removeEldestEntry(Map.Entry<String, Integer> eldest) {
        return size() > 3;  // max 3 entries
    }
};
lru.put("a", 1);
lru.put("b", 2);
lru.put("c", 3);
lru.get("a");          // "a" becomes most recently used
lru.put("d", 4);       // evicts "b" (least recently used)
System.out.println(lru.keySet()); // [c, a, d]
```

---

### TreeMap

Implements `NavigableMap` backed by a **Red-Black Tree**. Keys are always **sorted** (natural order or custom `Comparator`).

```java
TreeMap<String, Integer> tree = new TreeMap<>();
tree.put("banana", 2);
tree.put("apple", 1);
tree.put("cherry", 3);
System.out.println(tree.keySet()); // [apple, banana, cherry] — sorted

// Navigation methods
tree.firstKey();            // "apple"
tree.lastKey();             // "cherry"
tree.floorKey("b");         // "banana" — greatest key ≤ "b"
tree.ceilingKey("b");       // "banana" — smallest key ≥ "b"
tree.headMap("cherry");     // {apple=1, banana=2}
tree.tailMap("banana");     // {banana=2, cherry=3}
tree.subMap("apple", "cherry"); // {apple=1, banana=2}
```

| | HashMap | LinkedHashMap | TreeMap |
|---|---|---|---|
| Order | None | Insertion / access | Sorted (natural/custom) |
| `null` key | 1 allowed | 1 allowed | Not allowed (needs compare) |
| `get`/`put` time | O(1) avg | O(1) avg | O(log n) |
| Implements | `Map` | `Map` | `NavigableMap`, `SortedMap` |

---

### ArrayList vs LinkedList

```
ArrayList (dynamic array):
 [0][1][2][3][4][ ][ ][ ]   ← contiguous memory
  ↑ fast random access

LinkedList (doubly-linked list):
 head ↔ [A] ↔ [B] ↔ [C] ↔ [D] ↔ tail
         ↑ fast insert/delete at ends
```

```java
// ArrayList
List<String> arrayList = new ArrayList<>();
arrayList.add("a");
arrayList.get(2);        // O(1) — index directly
arrayList.add(0, "z");   // O(n) — shifts all elements right
arrayList.remove(0);     // O(n) — shifts all elements left

// LinkedList
LinkedList<String> linkedList = new LinkedList<>();
linkedList.add("a");
linkedList.get(2);           // O(n) — must traverse
linkedList.addFirst("z");    // O(1) — pointer update only
linkedList.removeLast();     // O(1) — pointer update only

// LinkedList as Deque
linkedList.offerFirst("x");  // push to front
linkedList.pollLast();       // pop from back
```

| Operation | ArrayList | LinkedList |
|---|---|---|
| Random access (`get(i)`) | **O(1)** | O(n) |
| Add/remove at end | O(1) amortized | **O(1)** |
| Add/remove at front/middle | O(n) | **O(1)** (with iterator) |
| Memory | Less (array) | More (node pointers) |
| Iteration | **Faster** (cache-friendly) | Slower (pointer chasing) |
| Use when | Random access, iteration | Frequent insert/delete at ends |

---

## 4. Comparable vs Comparator

### `Comparable` — Natural Ordering

Implement on the class itself. Defines the **default/natural** sort order.  
Method: `int compareTo(T other)`

```java
class Student implements Comparable<Student> {
    String name;
    int grade;

    @Override
    public int compareTo(Student other) {
        return Integer.compare(this.grade, other.grade); // sort by grade ascending
    }
}

List<Student> students = Arrays.asList(
    new Student("Alice", 85),
    new Student("Bob", 72),
    new Student("Charlie", 90)
);

Collections.sort(students);  // uses compareTo — natural order
// Result: Bob(72), Alice(85), Charlie(90)
```

Return contract:
- Negative → `this` comes **before** `other`
- Zero → equal
- Positive → `this` comes **after** `other`

---

### `Comparator` — Custom Ordering

External comparator — defines ordering **outside** the class. Use when:
- You don't own the class
- You need **multiple different** sort orders

```java
// Sort by name
Comparator<Student> byName = Comparator.comparing(s -> s.name);

// Sort by grade descending
Comparator<Student> byGradeDesc = Comparator.comparingInt(Student::getGrade).reversed();

// Chained: sort by grade desc, then name asc
Comparator<Student> combined = Comparator
    .comparingInt(Student::getGrade).reversed()
    .thenComparing(Student::getName);

students.sort(combined);

// With TreeMap
Map<Student, String> map = new TreeMap<>(byName);

// With Stream
students.stream()
    .sorted(byGradeDesc)
    .forEach(System.out::println);
```

### Comparison

| | `Comparable` | `Comparator` |
|---|---|---|
| Package | `java.lang` | `java.util` |
| Method | `compareTo(T o)` | `compare(T o1, T o2)` |
| Defined in | The class itself | Separate class / lambda |
| Sort orders | One (natural) | Many |
| Modifying class | Required | Not required |
| Used by | `Collections.sort()`, `TreeMap` | `Collections.sort()`, `Stream.sorted()` |

---

## 5. Fail-fast vs Fail-safe Iterators

### Fail-fast Iterators

Throw `ConcurrentModificationException` immediately if the collection is **structurally modified** while iterating (outside of iterator's own `remove()`).

Detected via an internal **`modCount`** counter — incremented on every structural change.

```java
List<String> list = new ArrayList<>(Arrays.asList("a", "b", "c", "d"));
Iterator<String> it = list.iterator();

while (it.hasNext()) {
    String val = it.next();
    if (val.equals("b")) {
        list.remove(val);  // ❌ ConcurrentModificationException!
    }
}

// ✅ Correct way — use iterator's own remove()
Iterator<String> it2 = list.iterator();
while (it2.hasNext()) {
    if (it2.next().equals("b")) {
        it2.remove();  // safe — updates modCount internally
    }
}

// ✅ Or use removeIf (Java 8+)
list.removeIf(s -> s.equals("b"));
```

**Collections with fail-fast iterators**: `ArrayList`, `HashMap`, `HashSet`, `LinkedList`, `TreeMap`

---

### Fail-safe Iterators

Operate on a **snapshot / copy** of the collection — structural changes to the original do **not** affect ongoing iteration.

```java
// CopyOnWriteArrayList — creates a new array on every write
List<String> cowList = new CopyOnWriteArrayList<>(Arrays.asList("a", "b", "c"));
for (String s : cowList) {
    if (s.equals("b")) {
        cowList.remove(s);  // ✅ no exception — iterates original snapshot
    }
}
System.out.println(cowList); // [a, c]

// ConcurrentHashMap — weakly consistent iterator
ConcurrentHashMap<String, Integer> map = new ConcurrentHashMap<>();
map.put("x", 1); map.put("y", 2); map.put("z", 3);

for (Map.Entry<String, Integer> e : map.entrySet()) {
    map.put("w", 4);  // ✅ no exception — may or may not see "w" during iteration
}
```

**Collections with fail-safe iterators**: `CopyOnWriteArrayList`, `CopyOnWriteArraySet`, `ConcurrentHashMap`

---

### Head-to-Head Comparison

| Feature | Fail-fast | Fail-safe |
|---|---|---|
| Exception on modification | `ConcurrentModificationException` | None |
| Works on | Original collection | Snapshot / copy |
| Memory overhead | Low | Higher (copy) |
| Reflects latest changes | Yes | Not guaranteed |
| Thread safety | No | Yes |
| Examples | `ArrayList`, `HashMap` | `CopyOnWriteArrayList`, `ConcurrentHashMap` |

### Common Trap — Enhanced For Loop

```java
// This also uses iterator internally — same fail-fast behavior
for (String s : list) {
    list.remove(s);  // ❌ ConcurrentModificationException
}
```

### Best Practices

```java
// 1. Use iterator.remove()
Iterator<String> it = list.iterator();
while (it.hasNext()) {
    if (condition(it.next())) it.remove();
}

// 2. Use removeIf (cleanest, Java 8+)
list.removeIf(s -> condition(s));

// 3. Collect to remove list, then remove after iteration
List<String> toRemove = new ArrayList<>();
for (String s : list) {
    if (condition(s)) toRemove.add(s);
}
list.removeAll(toRemove);

// 4. Use CopyOnWriteArrayList for concurrent scenarios
List<String> safeList = new CopyOnWriteArrayList<>(list);
```
