# Object Class Methods

---

## Table of Contents

1. [`toString()`](#1-tostring)
2. [`equals()`](#2-equals)
3. [`hashCode()`](#3-hashcode)
4. [`clone()`](#4-clone)
5. [`finalize()`](#5-finalize)
6. [`wait()`, `notify()`, `notifyAll()`](#6-wait-notify-notifyall)
7. [`getClass()`](#7-getclass)

---

## 1. `toString()`

Called whenever an object is **printed or concatenated with a String**. Default implementation returns `ClassName@hashCode` in hexadecimal — almost always useless.

```java
class Product {
    String name;
    double price;
    int quantity;

    Product(String name, double price, int quantity) {
        this.name = name;
        this.price = price;
        this.quantity = quantity;
    }

    // Without override:
    // System.out.println(p) → "Product@6d06d69c" — meaningless

    @Override
    public String toString() {
        return String.format("Product{name='%s', price=%.2f, quantity=%d}",
                             name, price, quantity);
    }
}

Product p = new Product("Laptop", 999.99, 5);
System.out.println(p);                   // Product{name='Laptop', price=999.99, quantity=5}
System.out.println("Item: " + p);        // "Item: Product{name='Laptop'...}" — auto calls toString()
log.info("Created: {}", p);              // same auto-call in logging frameworks
String s = p.toString();                 // explicit call
```

> **Interview Q: Why should you override `toString()`?**  
> The default `Object.toString()` returns `ClassName@hexHashCode` which is useless for debugging and logging. Overriding it provides a **human-readable representation** of the object's state. It's called automatically when an object is printed, logged, or concatenated with a String. Always override `toString()` for domain objects — it dramatically improves debuggability. IDEs (IntelliJ, Eclipse) can auto-generate it.

---

## 2. `equals()`

Default `Object.equals()` uses **reference equality** (`==`). Override to define logical equality based on field values.

```java
class Point {
    int x, y;

    Point(int x, int y) {
        this.x = x;
        this.y = y;
    }

    @Override
    public boolean equals(Object o) {
        // 1. Reflexivity: x.equals(x) must be true
        if (this == o) return true;

        // 2. Null check + type check (instanceof handles null safely)
        if (!(o instanceof Point)) return false;

        // 3. Cast and compare all fields used in equals
        Point other = (Point) o;
        return this.x == other.x && this.y == other.y;
    }
}

Point p1 = new Point(3, 4);
Point p2 = new Point(3, 4);
Point p3 = p1;

System.out.println(p1 == p2);         // false — different objects
System.out.println(p1.equals(p2));    // true  — same coordinates
System.out.println(p1 == p3);         // true  — same reference
System.out.println(p1.equals(p3));    // true
```

**The `equals()` contract (from `Object` Javadoc):**

```
1. Reflexive:   x.equals(x) == true
2. Symmetric:   x.equals(y) == y.equals(x)
3. Transitive:  x.equals(y) && y.equals(z) → x.equals(z)
4. Consistent:  x.equals(y) returns same result on repeated calls
5. Null-safe:   x.equals(null) == false (never true, never NPE)
```

> **Interview Q: What is the `equals()` contract and what happens if you violate it?**  
> The contract requires equals to be reflexive, symmetric, transitive, consistent, and return false for null. Violations can cause **broken behavior in collections** — HashMap, HashSet, and `Collections.sort()` all depend on a correct `equals()`. For example, breaking symmetry where `A.equals(B)` is true but `B.equals(A)` is false can cause elements to appear "lost" in a `Set`, or objects to be found with one key type but not another.

---

## 3. `hashCode()`

Must be consistent with `equals()` — **always override both together**.

```java
class Point {
    int x, y;

    Point(int x, int y) { this.x = x; this.y = y; }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (!(o instanceof Point)) return false;
        Point p = (Point) o;
        return x == p.x && y == p.y;
    }

    @Override
    public int hashCode() {
        // Use Objects.hash() — handles null, combines fields
        return Objects.hash(x, y);
    }

    // Manual (same result):
    // public int hashCode() {
    //     int result = 31 * x + y;
    //     return result;
    // }
}

// Why 31? It's prime, odd, and JVM optimizes 31*x to (x<<5) - x

// Demonstration of the contract:
Point a = new Point(1, 2);
Point b = new Point(1, 2);

System.out.println(a.equals(b));            // true
System.out.println(a.hashCode() == b.hashCode()); // true ← required by contract

// Without hashCode override:
Set<Point> set = new HashSet<>();
set.add(new Point(1, 2));
System.out.println(set.contains(new Point(1, 2)));  // false! (without hashCode)
// Both equal but different hashCodes → different buckets → not found
```

> **Interview Q: What happens if you override `equals()` but not `hashCode()`?**  
> The `equals()`-`hashCode()` contract is broken. Objects that are logically equal (via `equals()`) may have **different hashCodes** (the default hashCode uses memory address). This means `HashMap` and `HashSet` will **fail silently** — two equal objects land in different buckets, so `contains()` returns false, `get()` returns null, and `put()` creates duplicate keys. Always override both, using the **same fields** in both methods.

---

## 4. `clone()`

Creates a copy of an object. `Object.clone()` performs a **shallow copy** by default.

```java
class Address {
    String city;
    Address(String city) { this.city = city; }
}

class Employee implements Cloneable {
    String name;
    int salary;
    Address address;   // mutable field — shallow copy shares this reference

    Employee(String name, int salary, Address address) {
        this.name = name;
        this.salary = salary;
        this.address = address;
    }

    // ── SHALLOW COPY — default clone() behavior ──
    @Override
    public Employee clone() throws CloneNotSupportedException {
        return (Employee) super.clone();
        // Primitive/immutable fields (name, salary) are copied by value
        // address reference is COPIED — both original and clone share same Address object
    }

    // ── DEEP COPY — manually clone mutable nested objects ──
    public Employee deepClone() throws CloneNotSupportedException {
        Employee copy = (Employee) super.clone();
        copy.address = new Address(this.address.city);  // new Address object
        return copy;
    }
}

Employee original = new Employee("Alice", 50000, new Address("NYC"));
Employee shallow = original.clone();
Employee deep = original.deepClone();

shallow.address.city = "LA";
System.out.println(original.address.city);   // "LA"  — shared! shallow copy

deep.address.city = "Chicago";
System.out.println(original.address.city);   // "LA"  — NOT affected, deep copy
```

**Alternatives to `clone()`:**

```java
// Copy constructor (preferred)
class Employee {
    Employee(Employee other) {
        this.name = other.name;
        this.salary = other.salary;
        this.address = new Address(other.address.city);
    }
}

// Serialization-based deep copy (all fields, including nested)
Employee deepCopy = SerializationUtils.clone(original);  // Apache Commons
```

> **Interview Q: What is the difference between shallow copy and deep copy?**  
> **Shallow copy** copies the object and all its field values — for primitive fields, this means a copy of the value; for reference fields, it copies the **reference** (both original and copy point to the same nested object). **Deep copy** recursively copies all objects in the graph — the copy is completely independent. `Object.clone()` gives shallow copy; for deep copy you must override `clone()` to copy nested objects, use a copy constructor, or use serialization. In practice, copy constructors are preferred over `clone()`.

---

## 5. `finalize()`

Called by the GC **before** collecting an object. **Deprecated in Java 9, removed in Java 18** — do not use.

```java
class Resource {
    String name;
    boolean closed = false;

    Resource(String name) { this.name = name; }

    void close() {
        closed = true;
        System.out.println(name + " closed properly");
    }

    @Override
    @Deprecated
    protected void finalize() throws Throwable {
        // ❌ Problems with finalize():
        // 1. No guarantee when (or if) it runs
        // 2. Can resurrect objects (finalize can re-register the object)
        // 3. GC is delayed for objects with finalize — extra GC cycle needed
        // 4. Exceptions thrown here are ignored (swallowed)
        // 5. Can cause OutOfMemoryError if finalizable objects pile up
        if (!closed) {
            System.out.println("Warning: " + name + " not closed (finalize)");
        }
        super.finalize();
    }
}

// ✅ CORRECT: use AutoCloseable + try-with-resources instead
class Resource implements AutoCloseable {
    String name;
    Resource(String name) { this.name = name; }

    @Override
    public void close() {
        System.out.println(name + " closed");
    }
}

try (Resource r = new Resource("DB Connection")) {
    r.doWork();
}   // close() called automatically — guaranteed, immediate, no GC needed
```

> **Interview Q: Why was `finalize()` deprecated?**  
> `finalize()` has severe problems: (1) **No timing guarantee** — GC may run it after long delays or never; (2) **Revives objects** — a finalizer can re-register the object, keeping it alive; (3) **Performance** — objects with finalizers require two GC cycles to collect; (4) **Security** — subclasses can override and delay collection (finalizer attack); (5) **Exceptions ignored** — uncaught exceptions in `finalize()` are silently swallowed. Use `AutoCloseable` with try-with-resources for deterministic resource cleanup.

---

## 6. `wait()`, `notify()`, `notifyAll()`

Used for **inter-thread communication** on a shared monitor lock. Must be called inside a `synchronized` block.

```java
class MessageQueue {
    private final Queue<String> queue = new LinkedList<>();
    private final int MAX_SIZE = 5;

    // Producer: add message, notify consumer
    public synchronized void send(String message) throws InterruptedException {
        while (queue.size() == MAX_SIZE) {
            System.out.println("Queue full, producer waiting...");
            wait();   // releases lock, thread suspended
        }
        queue.add(message);
        System.out.println("Sent: " + message);
        notifyAll();  // wake up all waiting threads (consumer might be waiting)
    }

    // Consumer: wait for message, notify producer
    public synchronized String receive() throws InterruptedException {
        while (queue.isEmpty()) {
            System.out.println("Queue empty, consumer waiting...");
            wait();   // releases lock, thread suspended
        }
        String message = queue.poll();
        System.out.println("Received: " + message);
        notifyAll();  // wake up all waiting threads (producer might be waiting)
        return message;
    }
}
```

**Why `while` instead of `if` before `wait()`:**

```java
// ❌ WRONG — use if
synchronized void consume() throws InterruptedException {
    if (queue.isEmpty()) {
        wait();
    }
    // Dangerous: spurious wakeup can occur
    // Thread wakes up but queue might still be empty!
    queue.poll();   // could throw NoSuchElementException
}

// ✅ CORRECT — use while (handles spurious wakeups)
synchronized void consume() throws InterruptedException {
    while (queue.isEmpty()) {
        wait();    // loops back to check condition after waking up
    }
    queue.poll();   // safe — condition verified
}
```

| Method | What it does | Wakes up |
|---|---|---|
| `wait()` | Releases lock; suspends thread until `notify()` | N/A |
| `wait(long ms)` | Same but with timeout | After ms or notify |
| `notify()` | Wakes up **one** arbitrary waiting thread | One thread |
| `notifyAll()` | Wakes up **all** waiting threads | All threads |

> **Interview Q: What is the difference between `notify()` and `notifyAll()`?**  
> `notify()` wakes up **one** arbitrary thread waiting on the object's monitor. `notifyAll()` wakes up **all** threads waiting on that monitor. Use `notify()` only when all waiting threads are waiting for the **same condition** and any one of them can proceed (e.g., producers all wait for the same "queue not full" condition). Use `notifyAll()` when different threads wait for different conditions or when correctness requires all to re-check. When in doubt, `notifyAll()` is safer — extra wakeups just result in threads re-entering `wait()`.

---

## 7. `getClass()`

Returns the **runtime class** of the object — the actual type it was created as, not the reference type.

```java
class Animal { }
class Dog extends Animal { }

Animal a = new Dog();           // reference is Animal, object is Dog

System.out.println(a.getClass());               // class Dog
System.out.println(a.getClass().getName());     // "Dog" (simple name)
System.out.println(a.getClass().getSimpleName()); // "Dog"
System.out.println(a.getClass().getSuperclass()); // class Animal

// getClass() vs instanceof:
System.out.println(a instanceof Animal);   // true — IS-A (includes subclasses)
System.out.println(a instanceof Dog);      // true — IS-A
System.out.println(a.getClass() == Dog.class);   // true — exact type
System.out.println(a.getClass() == Animal.class); // false — not exactly Animal

// Why getClass() vs instanceof matters for equals():
// The standard pattern:
@Override
public boolean equals(Object o) {
    if (!(o instanceof MyClass)) return false;   // true for subclasses too
    // vs.
    if (o == null || getClass() != o.getClass()) return false; // exact type only
}
// Using instanceof: allows Dog to equal Animal (may break symmetry)
// Using getClass(): strict — only same runtime type (safer for equals contract)
```

> **Interview Q: What is the difference between `instanceof` and `getClass()` in `equals()`?**  
> `instanceof` returns true if the object is an instance of the given type **or any subclass**. `getClass()` checks the **exact runtime type**. In `equals()`, using `instanceof` can break the **symmetry** requirement — `animal.equals(dog)` might be true but `dog.equals(animal)` might be false if they check different fields. Using `getClass()` ensures only objects of the exact same class are compared, preserving all contract properties. Effective Java recommends `instanceof` combined with the subclass not being able to add equality-relevant state.
