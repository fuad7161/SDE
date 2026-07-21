# Core Java — In-Depth Notes

---

## Table of Contents

1. [OOP Principles (SOLID, DRY, KISS)](#1-oop-principles-solid-dry-kiss)
2. [`final`, `finally`, `finalize` Differences](#2-final-finally-finalize-differences)
3. [String Immutability, String Pool, `intern()`](#3-string-immutability-string-pool-intern)
4. [Exception Hierarchy (Checked vs Unchecked)](#4-exception-hierarchy-checked-vs-unchecked-custom-exceptions)
5. [`equals()` / `hashCode()` Contract](#5-equals--hashcode-contract)
6. [Generics (Wildcards, Bounded Types, Type Erasure)](#6-generics-wildcards-bounded-types-type-erasure)
7. [Java 8–21 Features](#7-java-821-features)
8. [Reflection & Annotations](#8-reflection--annotations)
9. [Java Memory Model (Heap, Stack, Metaspace)](#9-java-memory-model-heap-stack-metaspace-gc-roots)
10. [Garbage Collection Algorithms & Tuning](#10-garbage-collection-algorithms-g1-zgc-cms--tuning)

---

## 1. OOP Principles (SOLID, DRY, KISS)

### Four Pillars of OOP

| Pillar | Description |
|---|---|
| **Encapsulation** | Bundling data + behavior; hide internals via access modifiers |
| **Abstraction** | Expose only what is necessary (interfaces, abstract classes) |
| **Inheritance** | Child class reuses/extends parent behavior (`extends`) |
| **Polymorphism** | One interface, many forms — compile-time (overloading) & runtime (overriding) |

---

### SOLID Principles

#### S — Single Responsibility Principle (SRP)
> A class should have **one, and only one, reason to change.**

```java
// BAD — Report class handles both data and printing
class Report {
    String getContent() { ... }
    void printReport() { ... }   // second responsibility
}

// GOOD — separate concerns
class Report { String getContent() { ... } }
class ReportPrinter { void print(Report r) { ... } }
```

#### O — Open/Closed Principle (OCP)
> Classes should be **open for extension, closed for modification.**

```java
// Use abstraction so new shapes need no changes to existing code
interface Shape { double area(); }
class Circle  implements Shape { public double area() { return Math.PI * r * r; } }
class Square  implements Shape { public double area() { return side * side; } }
// Adding Triangle never touches existing classes
```

#### L — Liskov Substitution Principle (LSP)
> Subtypes must be **substitutable** for their base types without altering correctness.

```java
// VIOLATION — Square extends Rectangle but breaks setWidth/setHeight contract
class Rectangle { void setWidth(int w); void setHeight(int h); }
class Square extends Rectangle { /* must keep w==h — breaks LSP */ }
```

#### I — Interface Segregation Principle (ISP)
> Clients should not be forced to depend on **interfaces they do not use.**

```java
// BAD — fat interface
interface Worker { void work(); void eat(); }

// GOOD — segregated
interface Workable { void work(); }
interface Eatable  { void eat();  }
```

#### D — Dependency Inversion Principle (DIP)
> High-level modules should not depend on low-level modules. **Both should depend on abstractions.**

```java
// BAD
class OrderService { MySQLDatabase db = new MySQLDatabase(); }

// GOOD
class OrderService {
    private final Database db;  // abstraction
    OrderService(Database db) { this.db = db; }
}
```

---

### DRY — Don't Repeat Yourself
Extract duplicated logic into a single method/class/constant. Every piece of knowledge must have a **single, authoritative representation** in the system.

### KISS — Keep It Simple, Stupid
Prefer the simplest solution that works. Avoid unnecessary complexity, over-engineering, or premature abstraction.

---

## 2. `final`, `finally`, `finalize` Differences

### `final`
A **modifier** keyword with different effects:

```java
final int x = 10;          // constant — cannot reassign
final class MyClass {}     // cannot be subclassed (e.g., String)
final void method() {}     // cannot be overridden in subclasses

final List<String> list = new ArrayList<>();
list.add("ok");            // allowed — reference is final, not the object
// list = new ArrayList<>(); // ERROR
```

### `finally`
A block in try-catch that **always executes**, used for cleanup:

```java
try {
    riskyOperation();
} catch (IOException e) {
    handle(e);
} finally {
    connection.close();    // always runs, even if exception is thrown
}
```

**Exception**: `finally` does NOT run if:
- `System.exit()` is called
- JVM crashes
- The thread is killed

### `finalize()`
A method called by the **GC before collecting an object** — deprecated since Java 9, removed in Java 18:

```java
@Override
protected void finalize() throws Throwable {
    // cleanup — unreliable, unpredictable timing
    super.finalize();
}
```

**Prefer** `Closeable` / `AutoCloseable` with try-with-resources instead.

| | `final` | `finally` | `finalize()` |
|---|---|---|---|
| Type | Modifier | Block | Method |
| Purpose | Prevent change/override | Guaranteed cleanup | Pre-GC hook (deprecated) |

---

## 3. String Immutability, String Pool, `intern()`

### Why Strings are Immutable
- `String` class is `final`; char array (`value[]`) is `private final`
- **Security**: Safe for class loading, network connections, file paths
- **Thread safety**: Immutable objects are inherently thread-safe
- **Caching**: `hashCode()` can be cached safely

### String Pool (String Intern Pool)
- A special memory region in the **Heap** (Metaspace pre-Java 7)
- String **literals** are automatically interned
- Avoids creating duplicate string objects

```java
String a = "hello";          // goes to pool
String b = "hello";          // reuses same pool object
String c = new String("hello"); // new object on heap, NOT in pool

System.out.println(a == b);  // true  — same pool reference
System.out.println(a == c);  // false — different objects
System.out.println(a.equals(c)); // true — same content
```

### `intern()`
Forces a String to be placed in the pool (or returns the existing pool reference):

```java
String c = new String("hello").intern();
System.out.println(a == c);  // true — now same pool object
```

### StringBuilder vs StringBuffer vs String

| | String | StringBuilder | StringBuffer |
|---|---|---|---|
| Mutable | No | Yes | Yes |
| Thread-safe | Yes (immutable) | No | Yes (synchronized) |
| Performance | Slow for concat | Fastest | Slower (sync overhead) |

---

## 4. Exception Hierarchy (Checked vs Unchecked), Custom Exceptions

### Hierarchy

```
Throwable
├── Error                     (serious JVM problems — don't catch)
│   ├── OutOfMemoryError
│   ├── StackOverflowError
│   └── VirtualMachineError
└── Exception
    ├── RuntimeException      (UNCHECKED — compiler does not enforce)
    │   ├── NullPointerException
    │   ├── IllegalArgumentException
    │   ├── ArrayIndexOutOfBoundsException
    │   └── ClassCastException
    └── IOException           (CHECKED — must declare or handle)
        ├── FileNotFoundException
        └── SQLException
```

### Checked vs Unchecked

| | Checked | Unchecked |
|---|---|---|
| Extends | `Exception` (not RuntimeException) | `RuntimeException` or `Error` |
| Compiler | Forces `try-catch` or `throws` | No enforcement |
| Use case | Recoverable (file not found, network) | Programming bugs (null, bad args) |

### Custom Exceptions

```java
// Checked custom exception
public class InsufficientFundsException extends Exception {
    private final double amount;

    public InsufficientFundsException(double amount) {
        super("Insufficient funds: need " + amount + " more");
        this.amount = amount;
    }

    public double getAmount() { return amount; }
}

// Unchecked custom exception
public class InvalidOrderStateException extends RuntimeException {
    public InvalidOrderStateException(String message) {
        super(message);
    }

    public InvalidOrderStateException(String message, Throwable cause) {
        super(message, cause);   // preserve original stack trace
    }
}
```

### Best Practices
- Always include the **cause** when wrapping exceptions: `new MyException("msg", cause)`
- Prefer **unchecked** for programming errors, **checked** for recoverable conditions
- Never `catch (Exception e) {}` silently — at minimum, log it
- Use **try-with-resources** for `AutoCloseable` resources

```java
try (Connection conn = dataSource.getConnection();
     PreparedStatement ps = conn.prepareStatement(sql)) {
    // both are auto-closed
}
```

---

## 5. `equals()` / `hashCode()` Contract

### The Contract
1. If `a.equals(b)` is `true`, then `a.hashCode() == b.hashCode()` **must** be true
2. If `a.hashCode() == b.hashCode()`, `a.equals(b)` **may or may not** be true (collision allowed)
3. `equals()` must be: **reflexive**, **symmetric**, **transitive**, **consistent**, and `x.equals(null) == false`

### Why This Matters for Collections
`HashMap` / `HashSet` use `hashCode()` to find the bucket, then `equals()` to find the key:
- If you override `equals()` but not `hashCode()`, two "equal" objects may land in **different buckets** → `HashMap` treats them as different keys.

```java
class Point {
    int x, y;

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (!(o instanceof Point)) return false;
        Point p = (Point) o;
        return x == p.x && y == p.y;
    }

    @Override
    public int hashCode() {
        return Objects.hash(x, y);  // consistent with equals
    }
}
```

### Common Mistakes
- Using mutable fields in `hashCode()` — hash changes after object is stored in a map
- Not handling `null` in `equals()`
- Forgetting to check `instanceof` before casting

---

## 6. Generics (Wildcards, Bounded Types, Type Erasure)

### Why Generics?
Compile-time type safety + code reuse without casting:

```java
List<String> list = new ArrayList<>();
list.add("hello");
String s = list.get(0);  // no cast needed
```

### Bounded Type Parameters

```java
// Upper bound — T must be Number or subtype
<T extends Number>
List<? extends Number>  // can read as Number, cannot add

// Lower bound — T must be Integer or supertype
<T super Integer>
List<? super Integer>   // can add Integer, reading requires cast

// PECS — Producer Extends, Consumer Super
void copy(List<? extends T> src, List<? super T> dst)
```

### Wildcards

```java
// Unbounded — accepts any type
void print(List<?> list)

// Upper bounded — read-only, safe to read as Animal
void process(List<? extends Animal> animals)

// Lower bounded — write-safe, can add Dog
void addDogs(List<? super Dog> list)
```

### Type Erasure
Generics are a **compile-time** feature. At runtime, all generic type info is **erased**:

```java
List<String> strings = new ArrayList<>();
List<Integer> ints   = new ArrayList<>();

// Both have the same runtime type:
System.out.println(strings.getClass() == ints.getClass()); // true

// Cannot do at runtime:
// if (list instanceof List<String>)  // compile error
// new T()                            // compile error
// T[] array = new T[10]              // compile error
```

**Bridge methods** are generated by the compiler to maintain polymorphism after erasure.

---

## 7. Java 8–21 Features

### Java 8

#### Lambda Expressions
```java
// Before
Comparator<String> c = new Comparator<String>() {
    public int compare(String a, String b) { return a.compareTo(b); }
};

// After
Comparator<String> c = (a, b) -> a.compareTo(b);
```

#### Streams API
```java
List<String> names = people.stream()
    .filter(p -> p.getAge() > 18)
    .map(Person::getName)
    .sorted()
    .collect(Collectors.toList());

// Terminal operations: collect, forEach, reduce, count, findFirst, anyMatch
// Intermediate: filter, map, flatMap, distinct, sorted, limit, skip
```

#### Optional
```java
Optional<String> opt = Optional.ofNullable(getValue());
String result = opt
    .filter(s -> !s.isEmpty())
    .map(String::toUpperCase)
    .orElse("DEFAULT");

// Avoid: opt.get() without isPresent() check
// Prefer: orElse, orElseGet, orElseThrow, ifPresent, map
```

#### Default & Static Interface Methods
```java
interface Greeter {
    void greet(String name);                          // abstract
    default void greetAll(List<String> names) {       // default impl
        names.forEach(this::greet);
    }
    static Greeter formal() { return n -> "Dear " + n; } // static factory
}
```

---

### Java 9
- **Module System (JPMS)**: `module-info.java` with `requires`, `exports`
- **`List.of()`, `Map.of()`, `Set.of()`** — immutable factory methods
- **`Stream.takeWhile()`, `dropWhile()`, `iterate()` with predicate**

---

### Java 10
- **`var` (local variable type inference)**:
```java
var list = new ArrayList<String>();  // inferred as ArrayList<String>
var map  = Map.of("a", 1);
```

---

### Java 14–16 — Records
Concise data carriers — automatically generate constructor, getters, `equals`, `hashCode`, `toString`:

```java
record Point(int x, int y) {}

// Equivalent to a class with:
// - final fields x and y
// - canonical constructor
// - x(), y() accessors
// - equals/hashCode/toString

Point p = new Point(1, 2);
System.out.println(p.x()); // 1
```

Custom compact constructor:
```java
record Range(int min, int max) {
    Range {  // compact constructor
        if (min > max) throw new IllegalArgumentException();
    }
}
```

---

### Java 17 — Sealed Classes
Restrict which classes can extend/implement a type:

```java
public sealed interface Shape permits Circle, Rectangle, Triangle {}

public record Circle(double radius) implements Shape {}
public record Rectangle(double w, double h) implements Shape {}
public final class Triangle implements Shape { ... }
```

---

### Java 16–21 — Pattern Matching

#### `instanceof` Pattern Matching (Java 16)
```java
// Before
if (obj instanceof String) {
    String s = (String) obj;
    System.out.println(s.length());
}

// After
if (obj instanceof String s) {
    System.out.println(s.length());
}
```

#### Switch Expressions (Java 14+)
```java
String result = switch (day) {
    case MONDAY, FRIDAY -> "Weekday";
    case SATURDAY, SUNDAY -> "Weekend";
    default -> "Midweek";
};
```

#### Pattern Matching for Switch (Java 21)
```java
String format(Object obj) {
    return switch (obj) {
        case Integer i -> "int: " + i;
        case String s  -> "str: " + s;
        case null      -> "null";
        default        -> "other: " + obj;
    };
}
```

Works powerfully with sealed classes:
```java
double area(Shape s) {
    return switch (s) {
        case Circle c      -> Math.PI * c.radius() * c.radius();
        case Rectangle r   -> r.w() * r.h();
        case Triangle t    -> t.base() * t.height() / 2;
        // no default needed — compiler knows all permits
    };
}
```

---

### Java 21 — Virtual Threads (Project Loom)
Lightweight threads managed by the JVM, not OS:

```java
// Old OS thread
Thread t = new Thread(() -> handleRequest());

// Virtual thread — millions can exist simultaneously
Thread vt = Thread.ofVirtual().start(() -> handleRequest());

// With ExecutorService
try (var executor = Executors.newVirtualThreadPerTaskExecutor()) {
    executor.submit(() -> processRequest(req));
}
```

**Key benefit**: Blocking I/O on a virtual thread doesn't block the OS thread — JVM parks the virtual thread and reuses the carrier thread.

---

### Java Feature Timeline Summary

| Version | Key Feature |
|---|---|
| 8 | Lambdas, Streams, Optional, Default methods, `java.time` |
| 9 | Modules, `List.of()`, `Stream` enhancements |
| 10 | `var` |
| 11 | `String` methods, `var` in lambdas, HTTP Client |
| 14 | Switch expressions (standard), Records (preview) |
| 15 | Sealed classes (preview), Text blocks |
| 16 | Records (standard), `instanceof` pattern matching |
| 17 | Sealed classes (standard), LTS |
| 21 | Pattern matching for switch, Virtual threads, Sequenced Collections, LTS |

---

## 8. Reflection & Annotations

### Reflection
Allows inspection and manipulation of classes, methods, and fields **at runtime**:

```java
Class<?> clazz = Class.forName("com.example.MyClass");

// Get methods
Method[] methods = clazz.getDeclaredMethods();

// Invoke a method dynamically
Method method = clazz.getMethod("greet", String.class);
Object instance = clazz.getDeclaredConstructor().newInstance();
method.invoke(instance, "World");

// Access private field
Field field = clazz.getDeclaredField("name");
field.setAccessible(true);
field.set(instance, "Alice");
```

**Use cases**: Frameworks (Spring, Hibernate), serialization, testing (Mockito), DI containers.
**Drawbacks**: Performance overhead, breaks encapsulation, disables compiler checks.

---

### Annotations
Metadata attached to code elements, processed at compile-time or runtime:

```java
// Defining a custom annotation
@Retention(RetentionPolicy.RUNTIME)  // available at runtime
@Target(ElementType.METHOD)          // applies to methods
public @interface Loggable {
    String level() default "INFO";
}

// Using it
@Loggable(level = "DEBUG")
public void process() { ... }

// Processing at runtime via reflection
Method m = obj.getClass().getMethod("process");
if (m.isAnnotationPresent(Loggable.class)) {
    Loggable ann = m.getAnnotation(Loggable.class);
    System.out.println("Level: " + ann.level());
}
```

### Retention Policies

| Policy | When Available |
|---|---|
| `SOURCE` | Compile-time only (e.g., `@Override`) |
| `CLASS` | In `.class` file, not at runtime (default) |
| `RUNTIME` | Available via reflection at runtime |

---

## 9. Java Memory Model (Heap, Stack, Metaspace, GC Roots)

```
┌─────────────────────────────────────────────────┐
│                      JVM Memory                  │
│  ┌──────────┐  ┌────────────────┐  ┌──────────┐ │
│  │  Stack   │  │     Heap       │  │Metaspace │ │
│  │ (per     │  │ ┌────────────┐ │  │(native   │ │
│  │ thread)  │  │ │  Young Gen │ │  │ memory)  │ │
│  │          │  │ │ Eden|S0|S1 │ │  └──────────┘ │
│  │ frames   │  │ ├────────────┤ │               │
│  │ locals   │  │ │  Old Gen   │ │               │
│  │ refs     │  │ └────────────┘ │               │
│  └──────────┘  └────────────────┘               │
└─────────────────────────────────────────────────┘
```

### Stack
- One **per thread**; stores stack frames (method calls)
- Each frame holds: local variables, operand stack, reference to constant pool
- **LIFO** structure; frame pushed on call, popped on return
- Stores primitives and **object references** (not objects themselves)
- Fixed size → `StackOverflowError` on deep/infinite recursion

### Heap
- **Shared** across all threads; all objects live here
- Divided into:
  - **Young Generation**: Eden + Survivor spaces (S0, S1) — new objects allocated here
  - **Old (Tenured) Generation**: long-lived objects promoted from Young
- Subject to Garbage Collection

### Metaspace (Java 8+, replaced PermGen)
- Stores **class metadata** (bytecode, method info, static fields, interned strings in Java 7+)
- Lives in **native memory** (not heap) — grows dynamically
- Controlled by `-XX:MaxMetaspaceSize`

### GC Roots
Objects reachable from GC roots are considered **live** and will not be collected:
- Local variables in active stack frames
- Static fields
- Active Java threads
- JNI references

---

## 10. Garbage Collection Algorithms (G1, ZGC, CMS) & Tuning

### Minor GC vs Major GC vs Full GC

| GC Type | Area Collected | Frequency |
|---|---|---|
| Minor GC | Young Generation | Frequent, fast |
| Major GC | Old Generation | Less frequent, slower |
| Full GC | Entire Heap + Metaspace | Slowest, avoid in production |

---

### CMS (Concurrent Mark Sweep) — deprecated in Java 14
- Mostly concurrent, minimizes pause times for Old Gen
- Phases: Initial Mark (STW) → Concurrent Mark → Remark (STW) → Concurrent Sweep
- **Cons**: Fragmentation (no compaction), CPU intensive, deprecated

### G1GC (Garbage First) — default since Java 9
- Divides heap into equal-sized **regions** (~1–32 MB each)
- Prioritizes collecting regions with the most garbage first (hence "Garbage First")
- Handles both Young and Old generation
- Compacts heap → no fragmentation
- Predictable pause times via `-XX:MaxGCPauseMillis`

```
Heap divided into regions:
[E][E][S][O][O][H][E][O][S][E]
 E=Eden  S=Survivor  O=Old  H=Humongous
```

### ZGC (Z Garbage Collector) — production since Java 15
- Ultra-low latency: pauses < 1ms regardless of heap size
- Uses **colored pointers** and **load barriers** for concurrent work
- Supports heaps from MBs to TBs

### Key GC Tuning Flags

```bash
-Xms512m                    # Initial heap size
-Xmx4g                      # Max heap size
-XX:+UseG1GC                # Enable G1GC
-XX:MaxGCPauseMillis=200    # Target max pause time (G1)
-XX:NewRatio=3              # Old:Young ratio
-XX:SurvivorRatio=8         # Eden:Survivor ratio
-XX:+PrintGCDetails         # GC logging
-Xlog:gc*                   # Modern GC logging (Java 9+)
-XX:MaxMetaspaceSize=256m   # Cap metaspace
```

---
