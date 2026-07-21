# Keywords & Modifiers

---

## Table of Contents

1. [`final`](#1-final)
2. [`finally`](#2-finally)
3. [`finalize()`](#3-finalize)
4. [`static`](#4-static)
5. [`transient`](#5-transient)
6. [`volatile`](#6-volatile)
7. [`synchronized`](#7-synchronized)
8. [`native`](#8-native)
9. [`abstract`](#9-abstract)
10. [Access Modifiers](#10-access-modifiers)

---

## 1. `final`

Prevents modification — applies to variables, methods, and classes with different meanings.

```java
// ── final VARIABLE — value cannot be reassigned ──
final int MAX = 100;
// MAX = 200;   // ❌ compile error

// final reference — reference is fixed, but object's content can change
final List<String> names = new ArrayList<>();
names.add("Alice");      // ✅ — list content modified
names.add("Bob");        // ✅
// names = new ArrayList<>();  // ❌ — can't reassign the reference

// static final = constant (convention: UPPER_SNAKE_CASE)
static final double PI = 3.14159265;
static final String APP_NAME = "MyApp";

// Blank final — declared but assigned in constructor (per-instance constant)
class Circle {
    final double radius;   // blank final
    Circle(double r) {
        this.radius = r;   // must be assigned in constructor
    }
}

// ── final METHOD — cannot be overridden by subclasses ──
class Base {
    final void criticalOperation() {
        // subclass cannot change this behavior
        System.out.println("Critical logic");
    }
}
class Child extends Base {
    // @Override void criticalOperation() { }   // ❌ compile error
}

// ── final CLASS — cannot be subclassed ──
final class ImmutablePoint {
    final int x, y;
    ImmutablePoint(int x, int y) { this.x = x; this.y = y; }
}
// class ExtendedPoint extends ImmutablePoint { }  // ❌ compile error

// String, Integer, Double, etc. are all final classes
```

> **Interview Q: What are the different uses of the `final` keyword?**  
> (1) **final variable** — value cannot be reassigned after initialization (but if it's a reference, the object's internal state can still change); (2) **final method** — cannot be overridden by subclasses (used for security/template methods); (3) **final class** — cannot be extended (`String`, `Integer` are final). `final` on a field combined with no setters and a defensive copy in the constructor is how you create an **immutable class**.

---

## 2. `finally`

A block guaranteed to execute after `try`/`catch`, used for **cleanup**.

```java
public void processFile(String path) {
    FileReader fr = null;
    try {
        fr = new FileReader(path);
        // read file...
        String content = readContent(fr);
        return content;    // finally still runs before this return
    } catch (FileNotFoundException e) {
        System.err.println("File not found: " + path);
        return null;
    } finally {
        // Always executed — exception or no exception, return or no return
        if (fr != null) {
            try { fr.close(); } catch (IOException ignored) {}
            System.out.println("File closed");
        }
    }
}

// ── finally does NOT run when: ──
try {
    System.exit(0);     // JVM killed — finally does NOT run
} finally {
    System.out.println("NEVER PRINTED");
}

// ── return in finally OVERRIDES return in try ──
int weirdReturn() {
    try {
        return 1;       // would return 1...
    } finally {
        return 2;       // ...but this overrides it → returns 2
    }
}
// Avoid: never put return in finally — confusing, hides exceptions
```

> **Interview Q: What is the difference between `final`, `finally`, and `finalize()`?**  
> They are completely unrelated despite similar names:  
> - `final` is a **keyword** for constants, non-overridable methods, and non-extendable classes.  
> - `finally` is a **block** in exception handling that always executes (cleanup code).  
> - `finalize()` is a **method** in `Object` called by GC before collecting an object (deprecated, avoid).

---

## 3. `finalize()`

See [09_ObjectClassMethods.md](09_ObjectClassMethods.md#5-finalize) for full details.

```java
// Deprecated in Java 9, removed in Java 18
// DON'T USE — use AutoCloseable + try-with-resources instead

// ❌ Unreliable
protected void finalize() throws Throwable { cleanup(); }

// ✅ Correct alternative
class Resource implements AutoCloseable {
    @Override
    public void close() { cleanup(); }   // guaranteed, immediate
}
```

---

## 4. `static`

Belongs to the **class** rather than any instance.

```java
class Config {
    // ── static FIELD — one copy shared across all instances ──
    static int instanceCount = 0;
    static final String VERSION = "1.0";

    String name;

    // ── static BLOCK — runs once when class is first loaded ──
    static {
        System.out.println("Config class loading...");
        // one-time initialization, like loading a properties file
    }

    Config(String name) {
        instanceCount++;
        this.name = name;
    }

    // ── static METHOD — no 'this', cannot access instance members ──
    static int getCount() {
        return instanceCount;
        // return name;   // ❌ compile error — no instance available
    }

    // ── static NESTED CLASS — no reference to outer instance ──
    static class Builder {
        String name;
        Builder withName(String n) { this.name = n; return this; }
        Config build() { return new Config(name); }
    }
}

// Access static via class name
Config.instanceCount;
Config.getCount();
Config c = new Config.Builder().withName("prod").build();
```

> **Interview Q: Can a static method call an instance method?**  
> No — a static method belongs to the class and has no `this` reference. To call an instance method from a static method, you need an **object reference**: `MyClass obj = new MyClass(); obj.instanceMethod();`. But this is unusual. Static methods are for utility/factory functions that don't depend on any instance state.

---

## 5. `transient`

Marks a field to be **excluded from serialization**.

```java
import java.io.*;

class User implements Serializable {
    private static final long serialVersionUID = 1L;

    String username;
    String email;
    transient String password;     // excluded from serialization — security
    transient int sessionId;       // excluded — session is transient by nature
    transient DatabaseConnection conn;  // excluded — not serializable

    User(String username, String email, String password) {
        this.username = username;
        this.email = email;
        this.password = password;
        this.sessionId = (int)(Math.random() * 10000);
    }
}

// Serialization
User user = new User("alice", "alice@example.com", "secret123");
user.sessionId = 42;

try (ObjectOutputStream oos = new ObjectOutputStream(new FileOutputStream("user.ser"))) {
    oos.writeObject(user);
}

// Deserialization
try (ObjectInputStream ois = new ObjectInputStream(new FileInputStream("user.ser"))) {
    User restored = (User) ois.readObject();
    System.out.println(restored.username);    // "alice" ✅
    System.out.println(restored.password);    // null   ← transient, not saved
    System.out.println(restored.sessionId);   // 0      ← transient, not saved
}
```

> **Interview Q: What is `transient` and when do you use it?**  
> `transient` tells the Java serialization mechanism to **skip that field** when serializing. Use it for: (1) **sensitive data** (passwords, secrets) that shouldn't persist, (2) **derived/computed values** that can be recalculated after deserialization, (3) fields holding **non-serializable objects** (database connections, threads, file handles), (4) **session/runtime data** that has no meaning outside the current JVM session.

---

## 6. `volatile`

Ensures **visibility** of a variable across threads by bypassing CPU cache.

```java
class Server {
    volatile boolean running = true;    // reads/writes go to main memory

    void start() {
        new Thread(() -> {
            while (running) {     // always reads current value from main memory
                processRequest();
            }
            System.out.println("Server stopped");
        }).start();
    }

    void stop() {
        running = false;          // immediately visible to other threads
    }
}

// ── volatile guarantees ──
// 1. Visibility:  write by T1 is immediately visible to T2
// 2. Ordering:    prevents reordering across volatile access (happens-before)

// ── volatile does NOT guarantee ──
// Atomicity of compound operations:
volatile int count = 0;
count++;   // NOT atomic: read + increment + write → still a race condition

// Use AtomicInteger for atomic compound operations:
AtomicInteger atomicCount = new AtomicInteger(0);
atomicCount.incrementAndGet();   // atomic
```

> **Interview Q: What is the difference between `synchronized` and `volatile`?**  
> `volatile` ensures **visibility** — changes to a volatile variable are immediately visible to all threads. It is **not** atomic for compound operations. `synchronized` provides **mutual exclusion** AND visibility — only one thread executes the block at a time, ensuring both visibility and atomicity. Use `volatile` for simple **flags** or single-write/multiple-read patterns; use `synchronized` when you need **atomicity** across multiple operations.

---

## 7. `synchronized`

Ensures only **one thread** executes the synchronized code at a time.

```java
class Counter {
    private int count = 0;

    // ── synchronized METHOD — lock on 'this' ──
    public synchronized void increment() {
        count++;   // atomic: only one thread at a time
    }

    // ── synchronized BLOCK — finer control, specify lock object ──
    private final Object lock = new Object();

    public void safeIncrement() {
        synchronized (lock) {       // custom lock object
            count++;
        }
    }

    // ── static synchronized — lock on Class object ──
    static int globalCount = 0;
    public static synchronized void incrementGlobal() {
        globalCount++;
    }
}

// ── Reentrant: same thread can re-enter a synchronized block ──
class ReentrantExample {
    synchronized void outer() {
        System.out.println("outer");
        inner();   // ✅ same thread can call synchronized inner() — re-entrant
    }

    synchronized void inner() {
        System.out.println("inner");
    }
}
```

> **Interview Q: What is reentrant synchronization in Java?**  
> Java's intrinsic locks (monitor locks) are **reentrant** — a thread that already holds a lock **can re-acquire the same lock** without blocking itself. This allows synchronized methods to call other synchronized methods on the same object without deadlocking. Each re-entry increments the lock count; the lock is released only when the count drops to zero (after all synchronized blocks/methods exit). `ReentrantLock` from `java.util.concurrent` makes this behavior explicit.

---

## 8. `native`

Declares a method implemented in a **platform-specific language** (C/C++) via JNI (Java Native Interface).

```java
class NativeExample {
    // Declared in Java, implemented in C/C++
    public native void printHello();
    public native int multiply(int a, int b);

    static {
        System.loadLibrary("mylib");   // load the compiled native library
    }
}

// Common uses of native methods:
// - Hardware access (GPIO, graphics, audio)
// - OS-level system calls
// - Performance-critical routines
// - Wrapping legacy C/C++ code

// Examples in JDK:
// Object.hashCode()     — native
// System.arraycopy()    — native (fast memory copy)
// Thread.sleep()        — native (OS scheduler)
// Math.sqrt()           — native (hardware instruction)
```

> **Interview Q: What is the `native` keyword?**  
> `native` declares a method whose **implementation is in a native language** (C/C++) outside the JVM, accessed via the Java Native Interface (JNI). It has no method body in Java — just a declaration. Used for hardware access, OS system calls, and performance-critical operations. Many core JDK methods are native: `Object.hashCode()`, `System.arraycopy()`, `Thread.sleep()`, `Math.sqrt()`.

---

## 9. `abstract`

Declares an **incomplete type** — cannot be instantiated directly.

```java
// ── abstract CLASS ──
abstract class Shape {
    String color;

    Shape(String color) { this.color = color; }

    // Abstract method — no body, MUST be overridden
    abstract double area();
    abstract double perimeter();

    // Concrete method — shared implementation
    void describe() {
        System.out.printf("%s %s: area=%.2f%n", color, getClass().getSimpleName(), area());
    }
}

// ── abstract METHOD — forces subclass to provide implementation ──
class Circle extends Shape {
    double radius;
    Circle(String color, double r) { super(color); radius = r; }

    @Override double area() { return Math.PI * radius * radius; }
    @Override double perimeter() { return 2 * Math.PI * radius; }
}

// Cannot instantiate:
// Shape s = new Shape("red");   // ❌ compile error

// Can use as reference type (polymorphism):
Shape s = new Circle("blue", 5.0);
s.describe();   // blue Circle: area=78.54

// ── Abstract class with all concrete methods ──
// (still can't be instantiated — acts as base class with shared code)
abstract class Singleton {
    protected Singleton() {}
    abstract void initialize();
}
```

> **Interview Q: Can an abstract class have a constructor? Can it have `main()`?**  
> Yes, an abstract class **can have a constructor** — it's called by the subclass via `super()` to initialize inherited fields. It can also have a `main()` method. However, you **cannot directly instantiate** an abstract class with `new AbstractClass()`. The constructor only runs when a concrete subclass is instantiated. An abstract class with all methods concrete (no abstract methods) is valid — useful as a base class that prevents direct instantiation.

---

## 10. Access Modifiers

| Modifier | Same Class | Same Package | Subclass | Everywhere |
|---|---|---|---|---|
| `private` | ✅ | ❌ | ❌ | ❌ |
| (default/package) | ✅ | ✅ | ❌ | ❌ |
| `protected` | ✅ | ✅ | ✅ | ❌ |
| `public` | ✅ | ✅ | ✅ | ✅ |

```java
package com.example;

public class AccessDemo {
    private   int privateField   = 1;   // only within this class
              int packageField   = 2;   // within com.example package
    protected int protectedField = 3;   // package + subclasses (even in other packages)
    public    int publicField    = 4;   // everywhere

    private void privateMethod() {}
    void packageMethod() {}
    protected void protectedMethod() {}
    public void publicMethod() {}
}

// In a subclass in a DIFFERENT package:
package com.other;
class SubClass extends AccessDemo {
    void test() {
        // privateField      ❌ not accessible
        // packageField      ❌ not accessible (different package)
        protectedField;   // ✅ accessible in subclass
        publicField;      // ✅ accessible everywhere
    }
}
```

> **Interview Q: What is the default access modifier in Java?**  
> When no access modifier is specified, it's **package-private** (also called package-level or default access). The member is accessible only within the **same package**. It's less restrictive than `private` (which is class-only) but more restrictive than `protected` (which also allows subclass access from other packages). Interfaces have `public` access for their methods by default (not package-private).
