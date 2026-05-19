# Interview Comparisons

---

## Table of Contents

1. [`==` vs `equals()`](#1--vs-equals)
2. [Abstract Class vs Interface](#2-abstract-class-vs-interface)
3. [ArrayList vs LinkedList](#3-arraylist-vs-linkedlist)
4. [HashMap vs Hashtable](#4-hashmap-vs-hashtable)
5. [StringBuffer vs StringBuilder](#5-stringbuffer-vs-stringbuilder)
6. [Overloading vs Overriding](#6-overloading-vs-overriding)
7. [Process vs Thread](#7-process-vs-thread)
8. [Exception vs Error](#8-exception-vs-error)
9. [Checked vs Unchecked Exception](#9-checked-vs-unchecked-exception)
10. [Iterator vs ListIterator](#10-iterator-vs-listiterator)
11. [Comparable vs Comparator](#11-comparable-vs-comparator)
12. [sleep() vs wait()](#12-sleep-vs-wait)
13. [Stack vs Heap](#13-stack-vs-heap)

---

## 1. `==` vs `equals()`

| | `==` | `equals()` |
|---|---|---|
| Type | Operator | Method (from `Object`) |
| Compares | References (memory address) | Content (logical value) |
| For primitives | Value comparison | N/A (can't call on primitives) |
| For objects | Reference equality | Depends on override |
| Null-safe | Yes (`null == null` is true) | NPE if called on null |

```java
String a = "hello";
String b = "hello";
String c = new String("hello");

System.out.println(a == b);          // true  (same pool object)
System.out.println(a == c);          // false (different heap objects)
System.out.println(a.equals(c));     // true  (same content)

// Rule: always use equals() to compare String/Object content
// Exception: checking for null: obj == null is correct and idiomatic
```

> **Key rule**: Use `==` only for **null checks** and **comparing primitives**. Use `equals()` for **all object content comparison**. For null-safe comparison: `Objects.equals(a, b)`.

---

## 2. Abstract Class vs Interface

| | Abstract Class | Interface |
|---|---|---|
| Instance fields | ✅ Any type | ❌ Only `public static final` |
| Constructor | ✅ | ❌ |
| Concrete methods | ✅ | ✅ `default`/`static` (Java 8+) |
| Multiple inheritance | ❌ (only one `extends`) | ✅ (many `implements`) |
| Access modifiers | Any | `public` by default |
| `abstract` keyword | Required for abstract methods | Not required (implicitly abstract) |
| Extends/Implements | `extends` | `implements` |

```java
// Use Abstract Class when:
// - You want to share STATE (instance variables) among related classes
// - You need a constructor for common initialization
// - You have a "is-a" relationship with common behavior

// Use Interface when:
// - Defining a CAPABILITY or CONTRACT (Serializable, Comparable, AutoCloseable)
// - You need multiple inheritance
// - Unrelated classes share the same behavior

abstract class Animal {          // shared state + partial implementation
    String name;
    abstract void speak();
    void breathe() { System.out.println("breathing"); }
}

interface Swimmable {            // pure capability
    void swim();
}

interface Flyable {              // pure capability
    void fly();
}

class Duck extends Animal implements Swimmable, Flyable {
    @Override public void speak() { System.out.println("Quack"); }
    @Override public void swim()  { System.out.println("Swimming"); }
    @Override public void fly()   { System.out.println("Flying"); }
}
```

> **When to choose**: Abstract class when you have IS-A with shared code/state. Interface when defining CAN-DO capabilities across unrelated classes.

---

## 3. ArrayList vs LinkedList

| Operation | ArrayList | LinkedList |
|---|---|---|
| `get(i)` | **O(1)** — direct index | O(n) — traverse from head |
| `add` at end | O(1) amortized | **O(1)** |
| `add` at index | O(n) — shift right | O(n) to find + O(1) to insert |
| `remove` at end | O(1) | **O(1)** |
| `remove` at index | O(n) — shift left | O(n) to find + O(1) to remove |
| Memory | Less (array) | More (Node + 2 pointers each) |
| Cache performance | **Excellent** (contiguous) | Poor (scattered pointers) |
| Implements Deque | ❌ | ✅ (use as stack/queue) |

```java
// In practice: almost always use ArrayList
// Use LinkedList ONLY when:
// - Frequent add/remove at BOTH ends (Deque operations)
// - Never need random access
// For queue/stack, ArrayDeque is even better than LinkedList
```

> **Rule**: Default to `ArrayList`. Only switch to `LinkedList` if you need it as a Deque with heavy end operations and no random access.

---

## 4. HashMap vs Hashtable

| | HashMap | Hashtable |
|---|---|---|
| Thread-safe | ❌ No | ✅ Yes (full synchronized) |
| Null key | ✅ One allowed | ❌ NullPointerException |
| Null value | ✅ Multiple | ❌ NullPointerException |
| Performance | **Faster** (no sync) | Slower |
| Iterator | Fail-fast | Fail-fast |
| Introduced | Java 1.2 | Java 1.0 (legacy) |
| Ordered | ❌ | ❌ |

```java
// For thread safety, use ConcurrentHashMap — not Hashtable
ConcurrentHashMap<String, Integer> chm = new ConcurrentHashMap<>();
// - Segment/bucket-level locking (not whole-object)
// - No null keys or values
// - Atomic ops: putIfAbsent, computeIfAbsent, merge
```

> **Rule**: Use `HashMap` (single-thread) or `ConcurrentHashMap` (multi-thread). **Never use Hashtable** in new code.

---

## 5. StringBuffer vs StringBuilder

| | StringBuffer | StringBuilder |
|---|---|---|
| Mutable | ✅ | ✅ |
| Thread-safe | ✅ (synchronized methods) | ❌ |
| Performance | Slower (sync overhead) | **Faster** |
| When to use | Shared across multiple threads | Single-threaded use (99% of cases) |
| Introduced | Java 1.0 | Java 5 |
| API | Identical | Identical |

```java
// Single-threaded: ALWAYS use StringBuilder
StringBuilder sb = new StringBuilder("Hello");
sb.append(" World").append("!");
System.out.println(sb.toString());

// Multi-threaded shared string building: StringBuffer
// But in practice, refactor to avoid sharing mutable state
```

> **Rule**: Use `StringBuilder` in 99% of cases. `StringBuffer` only if the same instance is genuinely accessed by multiple threads simultaneously.

---

## 6. Overloading vs Overriding

| | Overloading | Overriding |
|---|---|---|
| Where | Same class | Subclass |
| Method name | Same | Same |
| Parameters | Must differ | Must be identical |
| Return type | Can differ | Same (or covariant subtype) |
| Access modifier | No restriction | Cannot be more restrictive |
| `static`/`private`/`final` | Can overload | Cannot override |
| Resolution | **Compile-time** | **Runtime** (dynamic dispatch) |
| Polymorphism type | Compile-time (ad-hoc) | Runtime (true polymorphism) |

```java
class Math {
    int add(int a, int b)         { return a + b; }      // overloading
    double add(double a, double b){ return a + b; }      // different param types
    int add(int a, int b, int c)  { return a + b + c; }  // different param count
}

class Shape { double area() { return 0; } }
class Circle extends Shape {
    @Override
    double area() { return Math.PI * r * r; }            // overriding
}
```

---

## 7. Process vs Thread

| | Process | Thread |
|---|---|---|
| Definition | Independent program instance | Lightweight unit of execution within a process |
| Memory | Separate address space | Shared memory within same process |
| Creation cost | Expensive | Lightweight |
| Communication | IPC (pipes, sockets, shared memory) | Direct (shared variables) |
| Isolation | High — crash doesn't affect others | Low — thread crash can kill process |
| Context switch | Expensive | Cheaper |
| Example | Two running Java programs | Two threads in same JVM |

```java
// Threads share heap memory within a process
class SharedState {
    static int counter = 0;  // shared across all threads in the JVM
}

// Threads can access and modify shared variables — needs synchronization
Thread t1 = new Thread(() -> { SharedState.counter++; });
Thread t2 = new Thread(() -> { SharedState.counter++; });
t1.start(); t2.start();   // race condition!
```

---

## 8. Exception vs Error

| | Exception | Error |
|---|---|---|
| Superclass | `java.lang.Throwable` | `java.lang.Throwable` |
| Represents | Application-level problems | System-level JVM problems |
| Recoverable | ✅ Usually | ❌ Generally not |
| Should catch | ✅ Yes | ❌ No (don't catch Error) |
| Examples | `IOException`, `NullPointerException` | `OutOfMemoryError`, `StackOverflowError` |

```java
// Exception — can and should handle
try {
    String s = null;
    s.length();
} catch (NullPointerException e) {
    System.out.println("Handled: " + e.getMessage());
}

// Error — JVM problem; catching is generally wrong
try {
    int[] arr = new int[Integer.MAX_VALUE];  // OutOfMemoryError
} catch (OutOfMemoryError e) {
    // Rarely useful — JVM state is unreliable at this point
}
```

---

## 9. Checked vs Unchecked Exception

| | Checked | Unchecked |
|---|---|---|
| Extends | `Exception` (not `RuntimeException`) | `RuntimeException` or `Error` |
| Compiler | Forces handle or declare | No requirement |
| For | External conditions (file, network, DB) | Programming bugs |
| Examples | `IOException`, `SQLException` | `NullPointerException`, `IllegalArgumentException` |

```java
// Checked — must handle
public void readFile() throws IOException { ... }

// Unchecked — optional
public void setAge(int age) {
    if (age < 0) throw new IllegalArgumentException("Negative age");  // no throws needed
}
```

> **Modern trend**: Prefer unchecked exceptions for cleaner APIs (Spring, Effective Java Item 71). Checked exceptions are best when the caller can actually recover.

---

## 10. Iterator vs ListIterator

| | Iterator | ListIterator |
|---|---|---|
| Works on | Any `Collection` | `List` only |
| Direction | Forward only | Forward and backward |
| `add` | ❌ | ✅ |
| `set` | ❌ | ✅ |
| `remove` | ✅ | ✅ |
| Index access | ❌ | `nextIndex()`, `previousIndex()` |

```java
List<String> list = new ArrayList<>(List.of("a", "b", "c"));

// Iterator — safe remove during iteration
Iterator<String> it = list.iterator();
while (it.hasNext()) {
    if (it.next().equals("b")) it.remove();  // safe
}

// ListIterator — bidirectional + modify
ListIterator<String> lit = list.listIterator(list.size());  // start at end
while (lit.hasPrevious()) {
    System.out.print(lit.previous() + " ");  // c a
}
```

---

## 11. Comparable vs Comparator

| | Comparable | Comparator |
|---|---|---|
| Package | `java.lang` | `java.util` |
| Method | `compareTo(T o)` | `compare(T o1, T o2)` |
| Defined in | Inside the class | External (separate / lambda) |
| Sort orders | One (natural) | Many (external) |
| Modifies class | Yes | No |

```java
// Comparable — class defines its own order
class Student implements Comparable<Student> {
    String name; int gpa;
    @Override public int compareTo(Student o) { return this.name.compareTo(o.name); }
}
Collections.sort(students);  // uses compareTo

// Comparator — external ordering
students.sort(Comparator.comparingInt(s -> s.gpa));       // by GPA
students.sort(Comparator.comparing(s -> s.name).reversed()); // by name desc
```

---

## 12. `sleep()` vs `wait()`

| | `Thread.sleep()` | `Object.wait()` |
|---|---|---|
| Class | `Thread` (static) | `Object` (instance) |
| Lock released | ❌ Keeps the lock | ✅ Releases the lock |
| `synchronized` required | ❌ No | ✅ Yes (IllegalMonitorStateException if not) |
| Woken by | Timeout or interrupt | `notify()`/`notifyAll()` or timeout |
| Purpose | Pause execution for time | Wait for condition + signal |

```java
// sleep: pause without releasing lock
synchronized void method() {
    Thread.sleep(1000);   // sleeps 1s but STILL HOLDS the lock
}

// wait: release lock and wait for signal
synchronized void method() throws InterruptedException {
    while (!condition) {
        wait();   // releases lock, waits for notify()
    }
}
```

> **Key difference**: `sleep()` keeps the lock — other threads can't enter the synchronized block. `wait()` releases the lock — other threads can proceed.

---

## 13. Stack vs Heap

| | Stack | Heap |
|---|---|---|
| Per-thread | ✅ Each thread has its own | ❌ Shared across all threads |
| Stores | Local vars, method params, return addresses | Objects, arrays |
| Allocation | Automatic (method call/return) | `new` keyword |
| Deallocation | Automatic (method returns) | Garbage Collector |
| Size | Small (limited, ~512KB-1MB) | Large (configured by -Xmx) |
| Speed | Faster | Slower |
| Error | `StackOverflowError` | `OutOfMemoryError` |
| Lifetime | Method scope | Until no more references |

```java
void method() {
    int x = 10;              // x on STACK
    String s = "hello";      // s (reference) on STACK, "hello" in Heap (pool)
    Person p = new Person(); // p (reference) on STACK, Person object on HEAP
}   // x, s, p removed from stack; Person object eligible for GC
```
