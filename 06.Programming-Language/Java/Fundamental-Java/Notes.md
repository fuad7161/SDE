# Java Interview — In-Depth Notes with Examples

---

## Table of Contents

1. [Java Basics](#1-java-basics)
2. [Object-Oriented Programming (OOP)](#2-object-oriented-programming-oop)
3. [String Handling](#3-string-handling)
4. [Exception Handling](#4-exception-handling)
5. [Collections Framework](#5-collections-framework)
6. [Java Memory Management](#6-java-memory-management)
7. [Multithreading & Concurrency](#7-multithreading--concurrency)
8. [Java 8 Features](#8-java-8-features)
9. [Object Class Methods](#9-object-class-methods)
10. [Keywords & Modifiers](#10-keywords--modifiers)
11. [File Handling & I/O](#11-file-handling--io)
12. [Important Interview Comparisons](#12-important-interview-comparisons)
13. [Important Java Concepts](#13-important-java-concepts)

---

## 1. Java Basics

### JDK vs JRE vs JVM

```
JDK (Java Development Kit)
└── JRE (Java Runtime Environment)
    └── JVM (Java Virtual Machine)
        ├── Class Loader
        ├── Bytecode Verifier
        └── Execution Engine (JIT Compiler)
```

| | JVM | JRE | JDK |
|---|---|---|---|
| Purpose | Runs bytecode | Provides runtime libs + JVM | Full dev kit (compiler, debugger, JRE) |
| Contains | Execution engine | JVM + standard libraries | JRE + `javac`, `javadoc`, tools |
| Use | Running apps | Running apps | Developing apps |

> **Interview**: "JVM is platform-specific (different JVM for Windows/Linux), but bytecode is platform-independent — that's how Java achieves Write Once, Run Anywhere."

---

### Data Types

| Type | Size | Range | Default |
|---|---|---|---|
| `byte` | 1 byte | -128 to 127 | 0 |
| `short` | 2 bytes | -32,768 to 32,767 | 0 |
| `int` | 4 bytes | -2^31 to 2^31-1 | 0 |
| `long` | 8 bytes | -2^63 to 2^63-1 | 0L |
| `float` | 4 bytes | ~3.4e38 | 0.0f |
| `double` | 8 bytes | ~1.8e308 | 0.0d |
| `char` | 2 bytes | 0 to 65,535 (Unicode) | '\u0000' |
| `boolean` | 1 bit | true / false | false |

---

### Type Casting

```java
// Widening (implicit) — no data loss
int i = 100;
long l = i;       // int → long, automatic
double d = l;     // long → double, automatic

// Narrowing (explicit) — possible data loss
double pi = 3.14;
int x = (int) pi;  // x = 3 — decimal part lost

// Common interview trap
byte b = (byte) 130;  // overflow! → b = -126
// 130 in binary: 10000010 → as signed byte = -126
```

---

### Java Program Structure

```java
package com.example;           // 1. package declaration (optional)

import java.util.List;         // 2. import statements

public class HelloWorld {      // 3. class declaration

    static int count = 0;      // 4. static/instance fields

    public static void main(String[] args) {   // 5. entry point
        System.out.println("Hello, World!");
    }
}
```

> **Interview**: "What is the order of execution in Java?" → Static blocks → Static fields → `main()` → Instance blocks → Constructors (for each object).

---

## 2. Object-Oriented Programming (OOP)

### Four Pillars

```
Encapsulation  → hide data, expose via methods (getters/setters)
Abstraction    → hide implementation, expose interface
Inheritance    → reuse parent behavior (extends / implements)
Polymorphism   → one name, many behaviors (overloading/overriding)
```

---

### Constructor Chaining

```java
class Employee {
    String name;
    int age;
    String dept;

    Employee() {
        this("Unknown", 0);              // calls Employee(String, int)
    }

    Employee(String name, int age) {
        this(name, age, "General");      // calls Employee(String, int, String)
    }

    Employee(String name, int age, String dept) {
        this.name = name;
        this.age = age;
        this.dept = dept;
    }
}
```

> **Interview**: "`this()` must be the first statement in a constructor. You can't call both `this()` and `super()` — only one first statement is allowed."

---

### Inheritance

```java
class Animal {
    String name;

    Animal(String name) { this.name = name; }

    void speak() { System.out.println(name + " makes a sound"); }
}

class Dog extends Animal {
    Dog(String name) {
        super(name);    // must call parent constructor first
    }

    @Override
    void speak() { System.out.println(name + " barks"); }   // runtime polymorphism
}

Animal a = new Dog("Rex");  // upcasting
a.speak();                  // prints "Rex barks" — dynamic dispatch
```

---

### Method Overloading vs Overriding

| | Overloading | Overriding |
|---|---|---|
| Where | Same class | Subclass |
| Signature | Different (param type/count) | Must be identical |
| Return type | Can differ | Must be same (or covariant) |
| Access | Any | Cannot be more restrictive |
| Binding | Compile-time (static) | Runtime (dynamic) |
| `static`/`private` | Can overload | Cannot override (hidden) |

```java
// Overloading — compile-time polymorphism
class Calculator {
    int add(int a, int b) { return a + b; }
    double add(double a, double b) { return a + b; }      // different params
    int add(int a, int b, int c) { return a + b + c; }   // different count
}

// Overriding — runtime polymorphism
class Shape { double area() { return 0; } }
class Circle extends Shape {
    double radius;
    Circle(double r) { this.radius = r; }
    @Override double area() { return Math.PI * radius * radius; }
}
```

---

### Interface vs Abstract Class

| | Abstract Class | Interface |
|---|---|---|
| Keyword | `abstract class` | `interface` |
| Instantiation | ❌ | ❌ |
| Inheritance | Single (`extends`) | Multiple (`implements`) |
| Constructor | ✅ | ❌ |
| Fields | Any type | `public static final` only |
| Methods | Abstract + concrete | Abstract + `default` + `static` (Java 8+) |
| Access modifiers | Any | `public` by default |
| Use case | Shared base with state | Contract / capability |

```java
abstract class Vehicle {
    String brand;                          // state is allowed
    Vehicle(String brand) { this.brand = brand; }  // constructor allowed
    abstract void start();                 // subclass must implement
    void stop() { System.out.println("Stopped"); }  // concrete method
}

interface Electric {
    int getRange();                        // abstract
    default void charge() {               // default — Java 8+
        System.out.println("Charging...");
    }
}

class Tesla extends Vehicle implements Electric {
    Tesla() { super("Tesla"); }
    @Override public void start() { System.out.println("Silently starting"); }
    @Override public int getRange() { return 500; }
}
```

---

### Association, Aggregation, Composition

```
Association  → "uses-a"       — Teacher ↔ Student (loose, both exist independently)
Aggregation  → "has-a"        — Department has Employees (employees exist without department)
Composition  → "part-of"      — House has Rooms (rooms cannot exist without house)
```

```java
// Composition — room cannot exist without house
class Room {
    String type;
    Room(String type) { this.type = type; }
}

class House {
    private final List<Room> rooms;   // House owns Rooms
    House() {
        rooms = new ArrayList<>();
        rooms.add(new Room("Bedroom"));  // created inside House
        rooms.add(new Room("Kitchen"));
    }
    // if House is destroyed, Rooms are too
}

// Aggregation — employee can exist without department
class Employee { String name; }
class Department {
    private List<Employee> employees;  // employees passed in, not created here
    Department(List<Employee> employees) { this.employees = employees; }
}
```

---

## 3. String Handling

### String Immutability

```java
String s = "hello";
s = s + " world";   // does NOT modify "hello"
                    // creates a NEW String "hello world"
                    // "hello" may be GC'd if no other reference
```

Why immutable?
- **Security** — file paths, class names, network URLs cannot be tampered with mid-use
- **Thread safety** — safe to share across threads with no synchronization
- **Caching** — `hashCode` computed once and cached; critical for `HashMap` keys

---

### String Pool

```java
String a = "java";            // stored in String pool
String b = "java";            // same pool object reused
String c = new String("java"); // new object on heap, outside pool

System.out.println(a == b);           // true  — same reference
System.out.println(a == c);           // false — different object
System.out.println(a.equals(c));      // true  — same content

String d = c.intern();                // moves c to pool (or returns existing)
System.out.println(a == d);           // true
```

---

### String vs StringBuilder vs StringBuffer

```java
// String — immutable, poor for repeated concatenation in loops
String result = "";
for (int i = 0; i < 1000; i++) result += i;  // creates 1000 objects! ❌

// StringBuilder — mutable, fast, NOT thread-safe
StringBuilder sb = new StringBuilder();
for (int i = 0; i < 1000; i++) sb.append(i);  // single object ✅
String result = sb.toString();

// StringBuffer — mutable, thread-safe (synchronized methods)
StringBuffer sbuf = new StringBuffer();
sbuf.append("thread").append("-").append("safe");
```

> **Interview**: "When should you use `StringBuffer`?" → Only when the same `StringBuilder` is modified by multiple threads simultaneously — which is rare. For single-thread, always prefer `StringBuilder`.

---

### Common String Methods

```java
String s = "  Hello, World!  ";

s.length()              // 17
s.trim()                // "Hello, World!"
s.toLowerCase()         // "  hello, world!  "
s.toUpperCase()         // "  HELLO, WORLD!  "
s.contains("World")     // true
s.startsWith("  He")    // true
s.replace("World", "Java")  // "  Hello, Java!  "
s.split(", ")           // ["  Hello", "World!  "]
s.substring(2, 7)       // "Hello"
s.charAt(2)             // 'H'
s.indexOf("o")          // 4
s.isEmpty()             // false
s.isBlank()             // false (Java 11+) — also checks whitespace-only
s.strip()               // "Hello, World!" (Java 11+) — handles Unicode whitespace
```

---

## 4. Exception Handling

### Hierarchy

```
Throwable
├── Error              ← Don't catch — JVM problems
│   ├── OutOfMemoryError
│   └── StackOverflowError
└── Exception
    ├── IOException    ← CHECKED — must handle
    │   └── FileNotFoundException
    ├── SQLException   ← CHECKED
    └── RuntimeException  ← UNCHECKED — no forced handling
        ├── NullPointerException
        ├── ArrayIndexOutOfBoundsException
        ├── ClassCastException
        └── IllegalArgumentException
```

---

### `throw` vs `throws`

```java
// throws — declares that a method MIGHT throw (checked exceptions)
public void readFile(String path) throws IOException {
    // ...
}

// throw — actually throws an exception
public void setAge(int age) {
    if (age < 0) throw new IllegalArgumentException("Age cannot be negative: " + age);
    this.age = age;
}
```

---

### `try-catch-finally`

```java
public int divide(int a, int b) {
    try {
        return a / b;                      // may throw ArithmeticException
    } catch (ArithmeticException e) {
        System.out.println("Division by zero: " + e.getMessage());
        return -1;
    } finally {
        System.out.println("Always runs");  // runs even with return in try/catch
    }
}

// Multiple catch blocks — more specific first
try {
    String s = null;
    s.length();
} catch (NullPointerException e) {
    System.out.println("Null: " + e.getMessage());
} catch (RuntimeException e) {          // catches other runtime exceptions
    System.out.println("Runtime: " + e.getMessage());
} catch (Exception e) {                 // most general — must be last
    System.out.println("General: " + e.getMessage());
}

// Multi-catch (Java 7+)
} catch (IOException | SQLException e) { ... }
```

---

### Custom Exception

```java
// Best practice: provide message + cause constructors
public class OrderNotFoundException extends RuntimeException {
    private final long orderId;

    public OrderNotFoundException(long orderId) {
        super("Order not found: " + orderId);
        this.orderId = orderId;
    }

    // Always provide cause constructor to preserve stack trace
    public OrderNotFoundException(long orderId, Throwable cause) {
        super("Order not found: " + orderId, cause);
        this.orderId = orderId;
    }

    public long getOrderId() { return orderId; }
}

// Usage
public Order findOrder(long id) {
    return orderRepository.findById(id)
        .orElseThrow(() -> new OrderNotFoundException(id));
}
```

---

### Try-with-Resources (Java 7+)

```java
// Automatically closes resources — no need for finally
try (FileReader fr = new FileReader("file.txt");
     BufferedReader br = new BufferedReader(fr)) {

    String line;
    while ((line = br.readLine()) != null) {
        System.out.println(line);
    }
}
// br and fr are closed automatically even if exception occurs
```

> **Interview**: "What if both `try` block and `close()` throw exceptions?" → The exception from `close()` is *suppressed*. You can retrieve it via `e.getSuppressed()`.

---

## 5. Collections Framework

### Hierarchy

```
Iterable
└── Collection
    ├── List         (ordered, duplicates allowed)
    │   ├── ArrayList
    │   └── LinkedList
    ├── Set          (no duplicates)
    │   ├── HashSet
    │   ├── LinkedHashSet  (insertion order)
    │   └── TreeSet        (sorted order)
    └── Queue
        ├── LinkedList
        └── PriorityQueue

Map (not a Collection)
    ├── HashMap
    ├── LinkedHashMap
    ├── TreeMap
    └── Hashtable (legacy)
```

---

### HashMap Internal Working

```
Key → hashCode() → index → bucket (array slot)
Each bucket is a linked list (or red-black tree for ≥8 entries)

put("Alice", 30):
  1. hash = "Alice".hashCode() → modified hash
  2. index = hash & (capacity - 1)   e.g. index = 3
  3. bucket[3].add(Entry("Alice", 30))
  4. On collision (same index, different key) → add to linked list

get("Alice"):
  1. Compute index (same as above)
  2. Traverse linked list at that index
  3. Find node where node.key.equals("Alice") → return value
```

Key facts:
- Default capacity: **16**, load factor: **0.75** → resizes at 12 entries
- Resize: doubles capacity, **rehashes** all entries
- Java 8+: linked list converts to **red-black tree** when bucket size ≥ 8 (O(log n) lookup)
- **NOT thread-safe** — use `ConcurrentHashMap` for concurrent access
- Allows **one `null` key**, multiple `null` values

```java
Map<String, Integer> map = new HashMap<>();
map.put("Alice", 30);
map.put("Bob", 25);
map.put(null, 0);       // null key allowed

// Safe access patterns
map.getOrDefault("Charlie", -1);                    // -1 if not found
map.putIfAbsent("Alice", 99);                       // won't overwrite
map.computeIfAbsent("Dave", k -> k.length());       // compute and put
map.merge("Alice", 1, Integer::sum);                // 30 + 1 = 31
```

---

### ArrayList vs LinkedList

| | ArrayList | LinkedList |
|---|---|---|
| Structure | Dynamic array | Doubly-linked list |
| Random access `get(i)` | O(1) | O(n) |
| Insert/delete at end | O(1) amortized | O(1) |
| Insert/delete at middle | O(n) — shifts elements | O(1) once at position (traversal is O(n)) |
| Memory | Less — only data | More — data + prev + next pointers |
| Cache performance | Better (contiguous) | Worse (scattered) |
| Use case | Read-heavy | Frequent insert/delete at ends |

```java
// ArrayList — backed by Object[]
List<String> list = new ArrayList<>(16);  // initial capacity hint
list.add("a"); list.add("b"); list.add("c");
list.get(1);  // O(1) → "b"

// LinkedList — also implements Deque (double-ended queue)
Deque<String> deque = new LinkedList<>();
deque.addFirst("front");
deque.addLast("back");
deque.pollFirst();  // removes and returns "front"
```

---

### Comparable vs Comparator

```java
// Comparable — natural ordering, implemented in the class itself
class Student implements Comparable<Student> {
    String name;
    int marks;

    @Override
    public int compareTo(Student other) {
        return Integer.compare(this.marks, other.marks);  // ascending by marks
    }
}

List<Student> students = new ArrayList<>();
Collections.sort(students);  // uses compareTo

// Comparator — external ordering, flexible
Comparator<Student> byName = Comparator.comparing(s -> s.name);
Comparator<Student> byMarksDesc = Comparator.comparingInt((Student s) -> s.marks).reversed();
Comparator<Student> complex = Comparator.comparing((Student s) -> s.name)
                                         .thenComparingInt(s -> s.marks);

students.sort(byName);
students.sort(byMarksDesc);
```

> **Interview**: "When to use Comparable vs Comparator?" → `Comparable` for the default/natural ordering of a class (e.g., `String` sorts alphabetically). `Comparator` for multiple orderings or when you can't modify the class.

---

### Fail-Fast vs Fail-Safe Iterators

```java
List<String> list = new ArrayList<>(Arrays.asList("a", "b", "c"));

// Fail-fast — throws ConcurrentModificationException if modified during iteration
Iterator<String> it = list.iterator();
while (it.hasNext()) {
    String s = it.next();
    if (s.equals("b")) list.remove(s);   // ❌ throws ConcurrentModificationException
}

// Correct way — use iterator's own remove
Iterator<String> it2 = list.iterator();
while (it2.hasNext()) {
    if (it2.next().equals("b")) it2.remove();   // ✅ safe
}

// Fail-safe — CopyOnWriteArrayList iterates a snapshot copy
List<String> cowList = new CopyOnWriteArrayList<>(Arrays.asList("a", "b", "c"));
for (String s : cowList) {
    if (s.equals("b")) cowList.remove(s);   // ✅ safe, iterates original snapshot
}
```

---

## 6. Java Memory Management

### Heap vs Stack

| | Stack | Heap |
|---|---|---|
| Contents | Method frames, local vars, references | All objects, instance variables |
| Shared? | Per-thread (private) | All threads (shared) |
| Size | Small (fixed) | Large (dynamic) |
| Management | Auto (LIFO push/pop) | GC managed |
| Error | `StackOverflowError` | `OutOfMemoryError` |
| Speed | Fast | Slower |

```java
void method() {
    int x = 10;               // x stored on STACK
    String s = "hello";       // reference 's' on STACK, "hello" object on HEAP
    Person p = new Person();  // reference 'p' on STACK, Person object on HEAP
}
// when method returns: x and s and p are popped from stack
// Person object on heap becomes eligible for GC (if no other references)
```

---

### Garbage Collection

```
Young Generation (Eden + S0 + S1):
  New objects → Eden
  Minor GC runs → live objects copied to Survivor (S0 or S1)
  After N minor GCs → promoted to Old Generation

Old Generation:
  Long-lived objects
  Major GC (slower)

Full GC:
  Entire heap + Metaspace — avoid in production (STW pause)
```

```java
// Object eligible for GC when no reachable references
String s = new String("hello");  // object on heap
s = null;                         // now eligible for GC

// Suggest GC (no guarantee)
System.gc();   // not recommended in production

// Finalization — deprecated, unreliable
// Use try-with-resources or Cleaner (Java 9+) instead
```

---

### Memory Leaks (Common Causes in Java)

```java
// 1. Unclosed resources
Connection conn = dataSource.getConnection();
// ... forgot to close → connection leak

// 2. Static collections holding references
static Map<String, Object> cache = new HashMap<>();
cache.put("key", new HeavyObject());  // never removed → never GC'd

// 3. Inner class holding outer class reference
class Outer {
    int data = 100;
    class Inner { void use() { System.out.println(data); } }
    // Inner holds implicit reference to Outer — Outer can't be GC'd while Inner is alive
}

// Fix — use static nested class or WeakReference
static class StaticInner { }  // no implicit outer reference
WeakReference<Object> ref = new WeakReference<>(new HeavyObject());
// GC can collect the object, ref.get() returns null after collection
```

---

### JVM Architecture

```
Source (.java)
    │ javac
    ▼
Bytecode (.class)
    │
    ▼
JVM
├── Class Loader Subsystem
│     ├── Bootstrap CL   (loads rt.jar / JDK core classes)
│     ├── Extension CL   (loads ext/*.jar)
│     └── Application CL (loads classpath classes)
│
├── Runtime Data Areas
│     ├── Method Area    (class metadata, static fields) → shared
│     ├── Heap           (objects) → shared
│     ├── Stack          (frames) → per thread
│     ├── PC Register    (current instruction) → per thread
│     └── Native Stack   → per thread
│
└── Execution Engine
      ├── Interpreter    (executes bytecode line by line — slow)
      ├── JIT Compiler   (compiles hot paths to native → fast)
      └── GC
```

---

## 7. Multithreading & Concurrency

### Thread Lifecycle

```
NEW ──start()──► RUNNABLE ──scheduler──► RUNNING
                     ▲                      │
                     │               sleep/wait/IO
                     │                      │
                     └──────────────── WAITING/BLOCKED
                                            │
                               TERMINATED ◄─┘ (run() completes)
```

---

### Creating Threads

```java
// 1. Extend Thread
class MyThread extends Thread {
    @Override
    public void run() {
        System.out.println("Thread: " + Thread.currentThread().getName());
    }
}
new MyThread().start();

// 2. Implement Runnable (preferred — allows extending other classes)
Runnable task = () -> System.out.println("Runnable running");
new Thread(task).start();

// 3. Callable — returns a result and can throw checked exceptions
Callable<Integer> calc = () -> {
    Thread.sleep(1000);
    return 42;
};
ExecutorService exec = Executors.newSingleThreadExecutor();
Future<Integer> future = exec.submit(calc);
Integer result = future.get();   // blocks until complete
exec.shutdown();
```

> **Interview**: "What is the difference between `Runnable` and `Callable`?"  
> `Runnable.run()` returns void and cannot throw checked exceptions.  
> `Callable.call()` returns a value (`Future<T>`) and can throw checked exceptions.

---

### Synchronization

```java
class Counter {
    private int count = 0;

    // synchronized method — only one thread can execute at a time
    public synchronized void increment() {
        count++;
    }

    // synchronized block — finer-grained (better performance)
    public void incrementBlock() {
        synchronized (this) {
            count++;
        }
    }

    public int getCount() { return count; }
}

// Race condition without synchronization:
// Thread A reads count=5, Thread B reads count=5
// Thread A writes 6, Thread B writes 6 — both incremented but count is 6, not 7!
```

---

### `volatile` Keyword

```java
class FlagExample {
    private volatile boolean running = true;  // visible to all threads immediately

    public void stop() { running = false; }   // write visible to other threads

    public void run() {
        while (running) {   // reads fresh value from main memory, not CPU cache
            doWork();
        }
    }
}
```

> **Interview**: "`volatile` guarantees visibility but NOT atomicity. `count++` on a volatile variable is still NOT thread-safe (it's read + increment + write — 3 operations). Use `AtomicInteger` for that."

```java
AtomicInteger atomicCount = new AtomicInteger(0);
atomicCount.incrementAndGet();   // atomic, thread-safe
atomicCount.compareAndSet(5, 10); // CAS — sets to 10 only if current value is 5
```

---

### Deadlock

```java
// Classic deadlock — two threads, two locks, opposite order
Object lockA = new Object();
Object lockB = new Object();

Thread t1 = new Thread(() -> {
    synchronized (lockA) {
        System.out.println("T1 holds A, waiting for B");
        synchronized (lockB) { System.out.println("T1 holds both"); }
    }
});

Thread t2 = new Thread(() -> {
    synchronized (lockB) {
        System.out.println("T2 holds B, waiting for A");
        synchronized (lockA) { System.out.println("T2 holds both"); }
    }
});

t1.start(); t2.start();   // may deadlock!

// Prevention: always acquire locks in the same order
Thread t2Fixed = new Thread(() -> {
    synchronized (lockA) {   // same order as T1: A then B
        synchronized (lockB) { System.out.println("T2 fixed"); }
    }
});
```

---

### Executor Framework

```java
// Fixed thread pool — reuses N threads
ExecutorService pool = Executors.newFixedThreadPool(4);

for (int i = 0; i < 10; i++) {
    int taskId = i;
    pool.submit(() -> {
        System.out.println("Task " + taskId + " on " + Thread.currentThread().getName());
    });
}

pool.shutdown();                         // no new tasks
pool.awaitTermination(10, TimeUnit.SECONDS);  // wait for completion

// Types of executors
Executors.newSingleThreadExecutor();     // 1 thread, sequential
Executors.newCachedThreadPool();         // unlimited threads, reuses idle ones
Executors.newScheduledThreadPool(2);     // run after delay or periodically
```

---

### Inter-Thread Communication

```java
class Buffer {
    private int data;
    private boolean hasData = false;

    // Producer
    public synchronized void produce(int value) throws InterruptedException {
        while (hasData) wait();     // wait if buffer full
        data = value;
        hasData = true;
        System.out.println("Produced: " + value);
        notify();                   // wake up consumer
    }

    // Consumer
    public synchronized int consume() throws InterruptedException {
        while (!hasData) wait();    // wait if buffer empty
        hasData = false;
        notify();                   // wake up producer
        return data;
    }
}
```

> **Interview**: "Why use `while` instead of `if` with `wait()`?" → **Spurious wakeups** — a thread can wake up without being notified. The `while` loop re-checks the condition after waking.

---

## 8. Java 8 Features

### Lambda Expressions

```java
// Functional interface — exactly ONE abstract method
@FunctionalInterface
interface Greeting {
    String greet(String name);
}

// Lambda = anonymous implementation of functional interface
Greeting formal = name -> "Dear " + name;
Greeting casual = name -> "Hey " + name + "!";

System.out.println(formal.greet("Alice"));  // "Dear Alice"

// Common built-in functional interfaces
Predicate<String> isLong = s -> s.length() > 5;       // boolean test
Function<String, Integer> length = String::length;     // transform
Consumer<String> print = System.out::println;          // consume, no return
Supplier<List<String>> listFactory = ArrayList::new;   // produce value
```

---

### Stream API

```java
List<String> names = Arrays.asList("Alice", "Bob", "Charlie", "Dave", "Anna");

// Filter + map + collect
List<String> result = names.stream()
    .filter(n -> n.startsWith("A"))    // ["Alice", "Anna"]
    .map(String::toUpperCase)          // ["ALICE", "ANNA"]
    .sorted()                          // ["ALICE", "ANNA"]
    .collect(Collectors.toList());

// Reduce — aggregate to single value
int totalLength = names.stream()
    .mapToInt(String::length)
    .sum();   // 5+3+7+4+4 = 23

// groupingBy
Map<Integer, List<String>> byLength = names.stream()
    .collect(Collectors.groupingBy(String::length));
// {3=[Bob], 4=[Dave, Anna], 5=[Alice], 7=[Charlie]}

// findFirst vs findAny (findAny is faster in parallel streams)
Optional<String> first = names.stream().filter(n -> n.length() > 4).findFirst();

// flatMap — flatten nested lists
List<List<Integer>> nested = Arrays.asList(Arrays.asList(1, 2), Arrays.asList(3, 4));
List<Integer> flat = nested.stream()
    .flatMap(Collection::stream)
    .collect(Collectors.toList());  // [1, 2, 3, 4]
```

> **Interview**: "What is the difference between `map` and `flatMap`?" → `map` transforms each element 1-to-1. `flatMap` transforms each element into a stream and merges all streams into one (1-to-many, then flatten).

---

### Optional

```java
// Avoid null checks with Optional
Optional<String> opt = Optional.ofNullable(getUsername());  // may be null

// BAD — defeats the purpose
if (opt.isPresent()) opt.get();

// GOOD — chain operations
String name = opt
    .filter(s -> s.length() > 3)
    .map(String::toUpperCase)
    .orElse("ANONYMOUS");           // default if empty

opt.ifPresent(n -> System.out.println("Hello " + n));  // only runs if present
opt.orElseThrow(() -> new UserNotFoundException());     // throw if empty
opt.orElseGet(() -> db.getDefaultUsername());           // lazy default (call only if empty)
```

---

### Method References

```java
// Instance method reference
names.forEach(System.out::println);       // same as s -> System.out.println(s)

// Static method reference
List<String> sorted = names.stream()
    .sorted(String::compareTo)            // same as (a,b) -> a.compareTo(b)
    .collect(Collectors.toList());

// Constructor reference
Supplier<ArrayList<String>> factory = ArrayList::new;  // same as () -> new ArrayList<>()
```

---

### Default & Static Interface Methods

```java
interface Validator<T> {
    boolean validate(T value);   // abstract

    // Default — can be overridden
    default Validator<T> and(Validator<T> other) {
        return value -> this.validate(value) && other.validate(value);
    }

    // Static factory
    static Validator<String> nonEmpty() {
        return s -> s != null && !s.isEmpty();
    }
}

Validator<String> nonEmpty = Validator.nonEmpty();
Validator<String> longEnough = s -> s.length() >= 8;
Validator<String> combined = nonEmpty.and(longEnough);

System.out.println(combined.validate("hello123"));  // true
System.out.println(combined.validate("hi"));        // false
```

---

## 9. Object Class Methods

Every Java class implicitly extends `java.lang.Object`.

### `toString()`

```java
class Person {
    String name; int age;
    Person(String name, int age) { this.name = name; this.age = age; }

    @Override
    public String toString() {
        return "Person{name='" + name + "', age=" + age + "}";
    }
}

Person p = new Person("Alice", 30);
System.out.println(p);           // calls toString() → "Person{name='Alice', age=30}"
// Without override → "com.example.Person@1b6d3586" (class@hashcode in hex)
```

---

### `equals()` and `hashCode()`

```java
class Point {
    int x, y;
    Point(int x, int y) { this.x = x; this.y = y; }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;             // same reference
        if (!(o instanceof Point)) return false; // null or different type
        Point p = (Point) o;
        return x == p.x && y == p.y;
    }

    @Override
    public int hashCode() {
        return Objects.hash(x, y);  // must be consistent with equals
    }
}

Set<Point> set = new HashSet<>();
set.add(new Point(1, 2));
System.out.println(set.contains(new Point(1, 2)));  // true (with correct equals+hashCode)
                                                     // false (without override)
```

---

### `clone()`

```java
class Config implements Cloneable {
    String host;
    int port;
    List<String> tags;

    // Shallow clone — new Config object but same 'tags' list reference!
    @Override
    public Config clone() throws CloneNotSupportedException {
        return (Config) super.clone();
    }

    // Deep clone — also clone mutable fields
    public Config deepClone() throws CloneNotSupportedException {
        Config copy = (Config) super.clone();
        copy.tags = new ArrayList<>(this.tags);  // new list, same string contents
        return copy;
    }
}
```

> **Interview**: "Shallow vs Deep copy?" → Shallow copy creates a new object but copies references (nested objects are shared). Deep copy creates new copies of all nested objects too.

---

### `wait()`, `notify()`, `notifyAll()`

```java
// Must be called inside synchronized block — else IllegalMonitorStateException
Object lock = new Object();

Thread waiter = new Thread(() -> {
    synchronized (lock) {
        try {
            System.out.println("Waiting...");
            lock.wait();           // releases lock and waits
            System.out.println("Resumed!");
        } catch (InterruptedException e) { Thread.currentThread().interrupt(); }
    }
});

Thread notifier = new Thread(() -> {
    synchronized (lock) {
        System.out.println("Notifying...");
        lock.notify();             // wakes ONE waiting thread
        // lock.notifyAll();       // wakes ALL waiting threads
    }
});

waiter.start();
Thread.sleep(100);
notifier.start();
```

---

## 10. Keywords & Modifiers

### `static`

```java
class MathUtils {
    static final double PI = 3.14159;   // class-level constant

    static int add(int a, int b) { return a + b; }  // no instance needed

    static {
        System.out.println("Static initializer — runs once when class is loaded");
    }
}

// Access without creating object
MathUtils.add(3, 4);
double pi = MathUtils.PI;
```

---

### `final`

```java
final int MAX = 100;                // cannot reassign
final class ImmutableClass {}       // cannot subclass
final void lockedMethod() {}        // cannot override

// Effectively final — variables used in lambdas must be final or effectively final
int multiplier = 3;   // effectively final (never reassigned)
Function<Integer, Integer> triple = n -> n * multiplier;   // ✅
```

---

### `transient`

```java
import java.io.*;

class User implements Serializable {
    String username;
    transient String password;   // excluded from serialization — won't be written to file
    transient int sessionToken;  // sensitive / temporary — don't persist
}

// After deserializing: username is restored, password and sessionToken are null/0
```

---

### `volatile`

```java
class Singleton {
    // Double-checked locking — volatile needed to prevent instruction reordering
    private static volatile Singleton instance;

    public static Singleton getInstance() {
        if (instance == null) {
            synchronized (Singleton.class) {
                if (instance == null) {
                    instance = new Singleton();  // without volatile, other threads
                    // might see partially constructed object
                }
            }
        }
        return instance;
    }
}
```

---

### `synchronized`

```java
class BankAccount {
    private double balance;

    public BankAccount(double initialBalance) { this.balance = initialBalance; }

    // Method-level lock — locks on 'this'
    public synchronized void deposit(double amount) {
        balance += amount;
    }

    // Block-level lock — same effect, but more granular
    public void withdraw(double amount) {
        synchronized (this) {
            if (balance >= amount) balance -= amount;
        }
    }

    // Static synchronized — locks on Class object (not instance)
    public static synchronized void staticMethod() { ... }
}
```

---

## 11. File Handling & I/O

### Byte Stream vs Character Stream

| | Byte Stream | Character Stream |
|---|---|---|
| Classes | `InputStream` / `OutputStream` | `Reader` / `Writer` |
| Unit | 1 byte at a time | 1 char (2 bytes, Unicode) at a time |
| Use for | Binary data (images, audio) | Text data |
| Examples | `FileInputStream`, `FileOutputStream` | `FileReader`, `FileWriter`, `BufferedReader` |

---

### Reading a File

```java
// BufferedReader — efficient line-by-line reading
try (BufferedReader br = new BufferedReader(new FileReader("data.txt"))) {
    String line;
    while ((line = br.readLine()) != null) {
        System.out.println(line);
    }
}

// Java 8+ — Files.lines() with Stream API
import java.nio.file.*;

try (Stream<String> lines = Files.lines(Paths.get("data.txt"))) {
    lines.filter(l -> l.contains("error"))
         .forEach(System.out::println);
}

// Read all at once (small files)
String content = Files.readString(Paths.get("data.txt"));   // Java 11+
List<String> allLines = Files.readAllLines(Paths.get("data.txt"));
```

---

### Writing a File

```java
// BufferedWriter
try (BufferedWriter bw = new BufferedWriter(new FileWriter("output.txt"))) {
    bw.write("Hello, World!");
    bw.newLine();
    bw.write("Second line");
}

// Files API — Java 7+
Files.writeString(Paths.get("output.txt"), "content");  // Java 11+
Files.write(Paths.get("output.txt"), List.of("line1", "line2"));
```

---

### Serialization

```java
// Serialize (Object → File)
class Student implements Serializable {
    private static final long serialVersionUID = 1L;  // version control for deserialization
    String name;
    int marks;
    transient String tempPassword;  // excluded
}

// Write to file
try (ObjectOutputStream oos = new ObjectOutputStream(new FileOutputStream("student.ser"))) {
    oos.writeObject(new Student("Alice", 95));
}

// Read from file (Deserialize)
try (ObjectInputStream ois = new ObjectInputStream(new FileInputStream("student.ser"))) {
    Student s = (Student) ois.readObject();
    System.out.println(s.name + ": " + s.marks);  // Alice: 95
    // s.tempPassword = null (transient)
}
```

> **Interview**: "What is `serialVersionUID`?" → A version identifier for serialized classes. If you change the class and the UID doesn't match the serialized file's UID, deserialization throws `InvalidClassException`. Always declare it explicitly to control versioning.

---

## 12. Important Interview Comparisons

### `==` vs `equals()`

```java
String a = "hello";
String b = "hello";
String c = new String("hello");

a == b        // true  — same pool reference
a == c        // false — different heap objects
a.equals(c)   // true  — same character content

Integer x = 127;
Integer y = 127;
x == y        // true  — Integer cache (-128 to 127)

Integer p = 200;
Integer q = 200;
p == q        // false — outside cache range, different objects
p.equals(q)   // true
```

---

### Abstract Class vs Interface

```java
// Abstract class — when sharing CODE and STATE
abstract class Logger {
    private String format;    // shared state
    Logger(String format) { this.format = format; }  // constructor

    abstract void log(String msg);   // subclass defines WHERE to log

    void logInfo(String msg) { log("[INFO] " + format + " " + msg); }  // shared behavior
}

// Interface — when defining a CONTRACT across unrelated classes
interface Auditable {
    void onCreated(Object entity);
    void onModified(Object entity);
    default void onDeleted(Object entity) {  // optional default
        System.out.println("Deleted: " + entity);
    }
}

// A class CAN extend one abstract class AND implement many interfaces
class AuditedService extends Logger implements Auditable, Serializable { ... }
```

---

### `HashMap` vs `Hashtable`

| | `HashMap` | `Hashtable` |
|---|---|---|
| Thread safety | Not thread-safe | Thread-safe (synchronized) |
| `null` keys | 1 null key allowed | ❌ `NullPointerException` |
| Performance | Faster | Slower (synchronization overhead) |
| Iteration order | Unordered | Unordered |
| Java version | Java 2 | Java 1 (legacy) |
| Preferred alternative | `ConcurrentHashMap` (thread-safe) | `ConcurrentHashMap` |

---

### `StringBuffer` vs `StringBuilder`

```java
// StringBuilder — fast, single-threaded
StringBuilder sb = new StringBuilder("Hello");
sb.append(", ").append("World").insert(5, "!"); // Hello!, World
sb.reverse();  // dlroW ,!olleH
sb.delete(0, 6);

// StringBuffer — thread-safe, use only when shared across threads
StringBuffer sbuf = new StringBuffer();
// Both have identical API; StringBuffer methods are synchronized
```

---

### Process vs Thread

| | Process | Thread |
|---|---|---|
| Definition | Independent program in execution | Smallest unit of execution within a process |
| Memory | Separate memory space | Shared memory within same process |
| Communication | IPC (expensive) | Direct (via shared memory) |
| Overhead | High (creation, context switch) | Low |
| Crash isolation | Crash doesn't affect other processes | Crash can kill all threads in process |

---

### Exception vs Error

```java
// Error — serious JVM/system problems; don't catch
try {
    recursiveMethod();
} catch (StackOverflowError e) {  // ← should NOT catch in production
    // can't reliably recover
}

// Exception — application-level problems; handle gracefully
try {
    int[] arr = new int[Integer.MAX_VALUE];
} catch (OutOfMemoryError e) {   // ← usually unrecoverable, log and exit
    log.error("OOM", e);
    System.exit(1);
}

// Checked Exception — must handle
try {
    Files.readAllLines(Paths.get("missing.txt"));  // throws IOException
} catch (IOException e) {
    System.out.println("File not found: " + e.getMessage());
}
```

---

## 13. Important Java Concepts

### Immutable Class

```java
// Rules for immutable class:
// 1. Class is final
// 2. All fields are private final
// 3. No setters
// 4. Deep copy mutable fields in constructor and getters

public final class ImmutablePerson {
    private final String name;
    private final int age;
    private final List<String> hobbies;   // mutable — needs defensive copy

    public ImmutablePerson(String name, int age, List<String> hobbies) {
        this.name = name;
        this.age = age;
        this.hobbies = List.copyOf(hobbies);  // defensive copy — immutable copy
    }

    public String getName() { return name; }
    public int getAge() { return age; }
    public List<String> getHobbies() { return hobbies; }  // already unmodifiable
}
```

---

### Singleton Class

```java
// Thread-safe Singleton — Bill Pugh / Initialization-on-demand Holder
public class Singleton {
    private Singleton() {}   // prevent instantiation

    private static class Holder {
        // Loaded only when getInstance() is first called
        private static final Singleton INSTANCE = new Singleton();
    }

    public static Singleton getInstance() {
        return Holder.INSTANCE;  // thread-safe without synchronization
    }
}

// Enum Singleton — simplest, handles serialization and reflection attacks
public enum ConfigManager {
    INSTANCE;

    private String dbUrl;
    public void setDbUrl(String url) { this.dbUrl = url; }
    public String getDbUrl() { return dbUrl; }
}

ConfigManager.INSTANCE.setDbUrl("jdbc:postgresql://localhost/mydb");
```

---

### Marker Interface

```java
// Marker interface — no methods, just "marks" a class for special behavior
// Examples: Serializable, Cloneable, RandomAccess

// JVM check:
if (obj instanceof Serializable) { /* serialize it */ }

// Custom marker interface
interface Auditable {}   // marks entities that should be audit-logged

class Order implements Auditable {
    long id; String status;
}

// Framework checks at runtime:
void save(Object entity) {
    repository.save(entity);
    if (entity instanceof Auditable) {
        auditLog.record(entity);   // only audit Auditable entities
    }
}
```

---

### Autoboxing & Unboxing

```java
// Autoboxing — primitive to wrapper (automatic)
Integer a = 5;           // int 5 → Integer.valueOf(5)
List<Integer> list = new ArrayList<>();
list.add(10);            // int 10 → Integer.valueOf(10)

// Unboxing — wrapper to primitive (automatic)
int b = a;               // Integer.intValue()
int sum = list.get(0) + list.get(1);  // unboxed for + operation

// Pitfall — NullPointerException from unboxing null
Integer x = null;
int y = x;               // NullPointerException! — tries x.intValue() on null

// Pitfall — Integer cache
Integer p = 127, q = 127;
System.out.println(p == q);   // true  (cached -128 to 127)
Integer r = 128, s = 128;
System.out.println(r == s);   // false (new objects)
System.out.println(r.equals(s)); // true
```

---

### Pass by Value in Java

```java
// Java is ALWAYS pass by value — for objects, it passes the VALUE of the reference

void changeValue(int x) { x = 99; }  // only changes local copy

int a = 5;
changeValue(a);
System.out.println(a);  // still 5 — primitive copied

void changeName(Person p) { p.name = "Bob"; }    // modifies object via reference
void reassign(Person p)   { p = new Person(); }  // only changes local reference copy

Person person = new Person("Alice");
changeName(person);
System.out.println(person.name);  // "Bob" — same object was modified

reassign(person);
System.out.println(person.name);  // still "Bob" — original reference unchanged
```

---

### Reflection API

```java
// Inspect class structure at runtime
Class<?> clazz = Class.forName("java.util.ArrayList");

System.out.println(clazz.getName());         // java.util.ArrayList
System.out.println(clazz.getSuperclass());   // java.util.AbstractList

// List all methods
for (Method m : clazz.getDeclaredMethods()) {
    System.out.println(m.getName());
}

// Create instance dynamically
Object obj = clazz.getDeclaredConstructor().newInstance();

// Access private field
class Secret { private String code = "abc123"; }
Field field = Secret.class.getDeclaredField("code");
field.setAccessible(true);   // bypasses private — security risk in prod!
String value = (String) field.get(new Secret());  // "abc123"

// Use cases: Spring DI, Mockito, JPA, Jackson serialization
// Avoid in performance-critical paths — reflection is 10-50x slower
```
