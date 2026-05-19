# Memory Management

---

## Table of Contents

1. [JVM Architecture](#1-jvm-architecture)
2. [Heap vs Stack](#2-heap-vs-stack)
3. [Garbage Collection](#3-garbage-collection)
4. [GC Algorithms](#4-gc-algorithms)
5. [Memory Leaks in Java](#5-memory-leaks-in-java)
6. [Object Lifecycle](#6-object-lifecycle)
7. [ClassLoader Subsystem](#7-classloader-subsystem)

---

## 1. JVM Architecture

```
┌─────────────────────────────────────────────────────┐
│                      JVM                            │
│                                                     │
│  ┌─────────────────────────────────────────────┐   │
│  │           Class Loader Subsystem            │   │
│  │   Loading → Linking → Initialization        │   │
│  └─────────────────────────────────────────────┘   │
│                        │                           │
│  ┌─────────────────────────────────────────────┐   │
│  │              Runtime Data Areas             │   │
│  │                                             │   │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐ │   │
│  │  │  Method  │  │   Heap   │  │  Stack   │ │   │
│  │  │  Area    │  │          │  │(per thread│ │   │
│  │  │(Metaspace│  │Young Gen │  │          │ │   │
│  │  │in Java 8)│  │Old Gen   │  │ Frames   │ │   │
│  │  └──────────┘  └──────────┘  └──────────┘ │   │
│  │                                             │   │
│  │  ┌──────────┐  ┌──────────────────────┐   │   │
│  │  │  PC Reg  │  │  Native Method Stack │   │   │
│  │  └──────────┘  └──────────────────────┘   │   │
│  └─────────────────────────────────────────────┘   │
│                                                     │
│  ┌──────────────────┐   ┌───────────────────────┐ │
│  │ Execution Engine │   │ Native Method Interface│ │
│  │  JIT Compiler    │   │       (JNI)            │ │
│  │  GC              │   └───────────────────────┘ │
│  └──────────────────┘                             │
└─────────────────────────────────────────────────────┘
```

| Component | Purpose |
|---|---|
| **Method Area (Metaspace)** | Class metadata, static variables, constant pool |
| **Heap** | All objects and arrays; shared across threads |
| **Stack** | Each thread has one; stores frames (local vars, operand stack, return addr) |
| **PC Register** | Points to current instruction being executed (per thread) |
| **Native Method Stack** | Supports native (C/C++) method calls |
| **JIT Compiler** | Compiles hot bytecode to native machine code at runtime |

> **Interview Q: What is the difference between JDK, JRE, and JVM?**  
> **JVM** (Java Virtual Machine) — executes bytecode; provides memory management and garbage collection; platform-specific (different JVM for Windows/Linux/Mac).  
> **JRE** (Java Runtime Environment) — JVM + standard libraries (java.lang, java.util, etc.); everything needed to **run** Java programs.  
> **JDK** (Java Development Kit) — JRE + development tools (compiler `javac`, debugger `jdb`, `jar`, `javap`); everything needed to **develop and compile** Java programs.  
> Relationship: JDK ⊃ JRE ⊃ JVM

---

## 2. Heap vs Stack

```
STACK (per thread)                HEAP (shared across all threads)
──────────────────                ────────────────────────────────
• Each thread has its own         • All objects created with 'new'
• Stores stack frames             • All arrays
  - Local variables               • String pool (Java 7+)
  - Method parameters             • Static variables' object values
  - Return addresses              • Shared by all threads
• LIFO — frame pushed on call,    • Managed by Garbage Collector
  popped on return                • Larger, slower, flexible size
• Faster, limited size            • Can grow/shrink dynamically
• Automatic cleanup on return
```

```java
public class MemoryDemo {
    static int staticVar = 42;    // value in Method Area

    public static void main(String[] args) {
        int x = 10;               // x lives on STACK
        String name = "Alice";    // 'name' reference on STACK, "Alice" in String Pool (HEAP)

        Person p = new Person("Bob", 25);  // p reference on STACK, Person object on HEAP
        calculate(x, p);
    }

    static void calculate(int num, Person person) {
        // New stack frame pushed:
        // - 'num' (copy of x = 10) on stack
        // - 'person' reference (copy of p) on stack
        // - 'result' on stack
        int result = num * 2;     // stack
        person.setAge(30);        // modifies the HEAP object that person points to
        // Frame popped when method returns — result, num, person reference gone
    }
}
```

**What causes StackOverflowError?**

```java
// Infinite recursion — each call pushes a new frame onto the stack
public int factorial(int n) {
    return n * factorial(n - 1);   // missing base case → StackOverflowError
}

// Fix:
public int factorial(int n) {
    if (n <= 1) return 1;           // base case
    return n * factorial(n - 1);
}
```

> **Interview Q: What is the difference between Heap and Stack memory in Java?**  
> **Stack** is per-thread, stores method frames (local variables, parameters, return addresses), follows LIFO order, automatically freed when a method returns. **Heap** is shared by all threads, stores all object instances and arrays, managed by the Garbage Collector. Stack is faster but limited in size — deep recursion causes `StackOverflowError`. Heap is larger but slower — too many live objects causes `OutOfMemoryError`.

---

## 3. Garbage Collection

GC automatically reclaims memory from objects that are no longer **reachable** (no live reference pointing to them).

```
Heap Generation Layout (classic):
──────────────────────────────────
Young Generation (new objects)
├── Eden Space      ← new objects created here
├── Survivor 0 (S0)
└── Survivor 1 (S1)

Old Generation (long-lived objects)

Metaspace (Java 8+) ← class metadata (was PermGen before Java 8)
```

**GC Lifecycle:**

```
1. Object created → Eden Space

2. Eden fills up → Minor GC triggered
   - Live objects from Eden → Survivor space (S0 or S1)
   - Objects alternate between S0/S1 on each Minor GC
   - Each surviving GC increments object's "age"

3. Object age reaches threshold (default: 15) → promoted to Old Generation
   Also promoted if: Survivor space full

4. Old Generation fills → Major/Full GC triggered
   - Scans entire heap (expensive — causes "stop-the-world" pause)

5. Unreachable objects → garbage collected, memory reclaimed
```

```java
// Making objects eligible for GC:
Person p = new Person("Alice");
p = null;           // original Person object now unreachable — eligible for GC

// Setting to null is one way, but usually not necessary —
// most objects become eligible when they go out of scope:
{
    Person local = new Person("Bob");
    // local goes out of scope here
}   // Bob is now eligible for GC

// System.gc() is just a HINT — JVM may or may not run GC
System.gc();   // not recommended in production code
```

> **Interview Q: What triggers Garbage Collection in Java?**  
> GC is triggered **automatically** by the JVM when: (1) the heap space (Eden, Survivor, or Old Gen) becomes full, (2) explicit `System.gc()` is called (but the JVM can ignore it), or (3) the JVM decides based on GC algorithm heuristics. An object becomes **eligible** for GC when no live reference path exists from any GC root (thread stacks, static variables, JNI references). You cannot force GC or predict exactly when it runs.

---

## 4. GC Algorithms

| Algorithm | Java version | Pause | Throughput | Best for |
|---|---|---|---|---|
| **Serial GC** | All | High (stop-the-world) | Low | Single-core, small heaps |
| **Parallel GC** | Default pre-Java 9 | Medium | High | Batch/throughput workloads |
| **CMS (Concurrent Mark Sweep)** | Deprecated Java 9 | Low | Medium | Low-latency apps |
| **G1 GC** | Default Java 9+ | Low-medium | High | Large heaps, general purpose |
| **ZGC** | Java 15+ (production) | Ultra-low (<1ms) | Good | Large heaps, low latency |
| **Shenandoah** | Java 12+ | Ultra-low | Good | Low-latency requirement |

```
G1 GC (Garbage First) — Java 9+ default:
─────────────────────────────────────────
• Divides heap into equal-sized regions (not contiguous Young/Old areas)
• Prioritizes collecting regions with most garbage first ("Garbage First")
• Concurrent marking runs alongside application
• Predictable pause time goals (-XX:MaxGCPauseMillis=200)
• Good balance of latency and throughput
```

**Key JVM flags:**

```bash
# Heap size
-Xms512m              # initial heap size
-Xmx2g                # max heap size

# GC selection
-XX:+UseG1GC          # use G1 (default Java 9+)
-XX:+UseZGC           # use ZGC (Java 15+)

# GC logging
-Xlog:gc*             # GC logs (Java 9+)

# Tune G1
-XX:MaxGCPauseMillis=200     # target max pause time
-XX:G1HeapRegionSize=16m     # region size
```

---

## 5. Memory Leaks in Java

Java has GC, but **memory leaks are still possible** when objects are kept alive unintentionally.

```java
// ── LEAK 1: Static collection growing forever ──
class Registry {
    private static List<Object> cache = new ArrayList<>();  // never cleared

    public static void register(Object obj) {
        cache.add(obj);   // objects accumulate — never become unreachable
    }
}

// Fix: use WeakReference or evict old entries

// ── LEAK 2: Listener/Callback not removed ──
class EventSource {
    private List<EventListener> listeners = new ArrayList<>();

    public void addListener(EventListener l) { listeners.add(l); }
    // Forgot to provide removeListener() — listener objects held forever
}

// Fix: always provide removeListener() and call it when done

// ── LEAK 3: ThreadLocal not cleaned ──
ThreadLocal<HeavyObject> threadLocal = new ThreadLocal<>();

// In a thread pool, threads are reused — ThreadLocal persists across tasks!
threadLocal.set(new HeavyObject());  // set
// ... use it ...
// threadLocal.remove();  // ← MUST call this at end of task!

// ── LEAK 4: Unclosed resources ──
public void loadData() throws IOException {
    FileInputStream fis = new FileInputStream("data.txt");
    // process...
    // forgot fis.close() — OS file descriptor leaked
}

// Fix: use try-with-resources
try (FileInputStream fis = new FileInputStream("data.txt")) {
    // process...
}   // automatically closed

// ── LEAK 5: Inner class holding outer class reference ──
class Outer {
    private byte[] data = new byte[1024 * 1024];   // 1MB

    class Inner {   // non-static inner class implicitly holds reference to Outer
        void doWork() { /* uses Outer.this */ }
    }
}

// If Inner instance lives longer than Outer should,
// the 1MB data cannot be collected
// Fix: use static nested class instead
```

> **Interview Q: Can memory leaks occur in Java? Give an example.**  
> Yes. Even with GC, memory leaks occur when **objects are kept reachable but are never used again**. Common causes: (1) **static collections** that grow indefinitely, (2) **event listeners/observers** never unregistered, (3) **`ThreadLocal` values** not removed at end of tasks (critical in thread pools), (4) **unclosed resources** (streams, connections), (5) **non-static inner classes** holding implicit reference to outer class. Tools: JVisualVM, JProfiler, Eclipse MAT to analyze heap dumps.

---

## 6. Object Lifecycle

```
1. CLASS LOADING (once per class)
   ClassLoader reads .class file
   JVM creates Class object in Metaspace

2. OBJECT CREATION
   a. Memory allocated on Heap (Eden space)
   b. Instance fields set to defaults (0, false, null)
   c. Static initializer blocks run (once, at class load)
   d. Constructor executes:
      - super() called first (recursively up to Object)
      - Instance initializer blocks run (in order)
      - Constructor body runs

3. IN USE
   Object is reachable via references
   Methods called, fields read/modified

4. UNREACHABLE
   No more references from any GC root
   Object eligible for GC

5. FINALIZATION (optional, deprecated Java 9)
   GC may call finalize() before reclaiming
   Unreliable — don't depend on it

6. GARBAGE COLLECTED
   Memory reclaimed
```

```java
class LifecycleExample {
    static {
        System.out.println("1. Static block (runs once at class load)");
    }

    {
        System.out.println("3. Instance initializer block (runs before constructor body)");
    }

    int value;

    LifecycleExample(int v) {
        System.out.println("4. Constructor body");
        this.value = v;
    }

    public static void main(String[] args) {
        System.out.println("2. main() starts");
        LifecycleExample obj = new LifecycleExample(42);
        // Output:
        // 1. Static block (runs once at class load)
        // 2. main() starts
        // 3. Instance initializer block (runs before constructor body)
        // 4. Constructor body
    }
}
```

> **Interview Q: What is the order of execution when an object is created?**  
> 1. **Parent class static blocks** (once at class load, top-to-bottom)  
> 2. **Subclass static blocks** (once at class load)  
> 3. **Parent class instance initializer blocks** (every instantiation)  
> 4. **Parent class constructor body**  
> 5. **Subclass instance initializer blocks**  
> 6. **Subclass constructor body**

---

## 7. ClassLoader Subsystem

```
Delegation Model (Parent-First):
─────────────────────────────────
Bootstrap ClassLoader      ← loads java.lang.*, java.util.*, etc. (JDK core)
       ↑ parent of
Extension ClassLoader      ← loads jdk/lib/ext/ (JDK extensions)
       ↑ parent of
Application ClassLoader    ← loads your app's classpath (your .class files)
       ↑ parent of
Custom ClassLoader         ← user-defined (hot reload, plugin systems)
```

```java
// ClassLoader delegation:
// When A.class is requested:
// 1. Check Application ClassLoader's cache
// 2. If not found → delegate to parent (Extension)
// 3. Extension delegates to Bootstrap
// 4. Bootstrap tries to load → if found, returns
// 5. If not found, Extension tries → if found, returns
// 6. If not found, Application ClassLoader loads it
// This prevents user code from overriding core Java classes

// Viewing classloaders
System.out.println(String.class.getClassLoader());       // null = Bootstrap
System.out.println(Main.class.getClassLoader());         // ApplicationClassLoader
System.out.println(ClassLoader.getSystemClassLoader());  // ApplicationClassLoader
```

> **Interview Q: What is the parent delegation model in Java ClassLoader?**  
> When a ClassLoader is asked to load a class, it first **delegates the request to its parent** ClassLoader. Only if the parent cannot find the class does the child attempt to load it itself. This ensures that **core Java classes** (like `java.lang.String`) are always loaded by the Bootstrap ClassLoader — preventing user-defined classes from maliciously replacing core classes. The chain: Bootstrap → Extension → Application → Custom.
