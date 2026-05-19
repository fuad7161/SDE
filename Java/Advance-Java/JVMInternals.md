# JVM Internals — In-Depth Notes

---

## Table of Contents

1. [Class Loading Mechanism](#1-class-loading-mechanism)
2. [JIT Compilation](#2-jit-compilation)
3. [GC Tuning Flags](#3-gc-tuning-flags)
4. [Memory Leak Identification & Heap Dump Analysis](#4-memory-leak-identification--heap-dump-analysis)
5. [StackOverflowError vs OutOfMemoryError](#5-stackoverflowerror-vs-outofmemoryerror)

---

## 1. Class Loading Mechanism

### Overview

When the JVM needs a class for the first time, the **ClassLoader subsystem** loads, links, and initializes it.

```
Source (.java)
    │ javac
    ▼
Bytecode (.class)
    │ ClassLoader
    ▼
┌──────────────────────────────────────┐
│             ClassLoader              │
│  1. Loading   → read .class bytes    │
│  2. Linking                          │
│     a. Verification  → bytecode valid│
│     b. Preparation   → static fields │
│     c. Resolution    → symbolic refs │
│  3. Initialization → static blocks   │
└──────────────────────────────────────┘
    │
    ▼
Method Area (Metaspace) — class metadata stored
```

---

### The Three Built-in ClassLoaders

```
Bootstrap ClassLoader         (native C++ — no Java class)
    │   loads: rt.jar / java.base module (java.lang, java.util, ...)
    │
    ▼
Extension / Platform ClassLoader    (java.lang.ClassLoader)
    │   loads: lib/ext/, javax.*, java.sql, ...
    │
    ▼
Application / System ClassLoader
        loads: classpath (-cp), your own classes
```

```java
// Observing classloaders
System.out.println(String.class.getClassLoader());         // null → Bootstrap
System.out.println(com.sun.nio.fs.LinuxFileSystem.class
                       .getClassLoader());                 // PlatformClassLoader
System.out.println(MyApp.class.getClassLoader());          // AppClassLoader
```

---

### Delegation Model (Parent-First)

Before loading a class, a ClassLoader **delegates to its parent** first.  
Only if the parent cannot find it does the child attempt to load it.

```
AppClassLoader.loadClass("com.example.Foo")
  → delegates to ExtClassLoader
      → delegates to Bootstrap
          → not found in rt.jar
      ← not found in ext
  → AppClassLoader loads from classpath ✅
```

This prevents user code from replacing core classes like `java.lang.String`.

```java
// Custom ClassLoader example
class MyClassLoader extends ClassLoader {
    @Override
    protected Class<?> findClass(String name) throws ClassNotFoundException {
        byte[] bytes = loadClassBytes(name);   // read .class file from custom source
        return defineClass(name, bytes, 0, bytes.length);
    }

    private byte[] loadClassBytes(String name) {
        String path = name.replace('.', '/') + ".class";
        try (InputStream is = getClass().getResourceAsStream("/" + path)) {
            return is.readAllBytes();
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }
}

// Usage
MyClassLoader loader = new MyClassLoader();
Class<?> clazz = loader.loadClass("com.example.MyPlugin");
Object instance = clazz.getDeclaredConstructor().newInstance();
```

---

### Class Initialization Order

```java
class Parent {
    static int x = 10;
    static { System.out.println("Parent static block, x=" + x); }
    int y = 20;
    { System.out.println("Parent instance block"); }
    Parent() { System.out.println("Parent constructor"); }
}

class Child extends Parent {
    static int a = 30;
    static { System.out.println("Child static block, a=" + a); }
    int b = 40;
    { System.out.println("Child instance block"); }
    Child() { System.out.println("Child constructor"); }
}

new Child();
// Output order:
// 1. Parent static block, x=10    ← static init, parent first
// 2. Child static block, a=30     ← static init, child
// 3. Parent instance block        ← instance init, parent
// 4. Parent constructor
// 5. Child instance block         ← instance init, child
// 6. Child constructor
```

---

## 2. JIT Compilation

### Interpretation vs Compilation

```
Java Bytecode (.class)
        │
        ▼
   JVM Interpreter          ← executes bytecode line by line — slow start, no warmup
        │
        │ (hot code detected — called frequently)
        ▼
    JIT Compiler            ← compiles hot bytecode → native machine code
        │
        ▼
   Native Machine Code      ← cached in Code Cache, executed directly — very fast
```

The JVM starts **interpreting** bytecode immediately (fast startup), then identifies **hot spots** — methods called frequently — and compiles them to native code at runtime.

---

### Tiered Compilation (Java 7+, default since Java 8)

```
Level 0 — Interpreter
Level 1 — C1 (client compiler) — simple fast compilation, basic optimizations
Level 2 — C1 with limited profiling
Level 3 — C1 with full profiling (counts method calls, branch frequencies)
Level 4 — C2 (server compiler) — deep aggressive optimization using profiling data
```

A method typically goes: `0 → 3 → 4` under load.

```bash
# View JIT decisions
java -XX:+PrintCompilation MyApp

# Output format:
# timestamp  compile_id  flags  method_name  (size)
#   127       23         %      com.example.HotMethod::compute @ 12 (45 bytes)
```

---

### Key JIT Optimizations

**Inlining** — replaces a method call with the method body:
```java
// Source
int result = add(a, b);
int add(int x, int y) { return x + y; }

// After JIT inlining — no method call overhead
int result = a + b;
```

**Escape Analysis** — if an object never escapes the method, allocate it on the **stack** (or eliminate it entirely):
```java
void process() {
    Point p = new Point(1, 2);  // p never returned or stored externally
    int sum = p.x + p.y;        // JIT may eliminate the allocation entirely
}
```

**Loop Unrolling** — reduces loop overhead:
```java
// Source
for (int i = 0; i < 4; i++) arr[i] = i;

// After unrolling
arr[0]=0; arr[1]=1; arr[2]=2; arr[3]=3;  // no loop counter, no branch
```

**Dead Code Elimination**, **Null Check Elimination**, **Devirtualization** (convert virtual call to direct call when only one implementation).

---

### Monitoring JIT

```bash
# Print compiled methods
-XX:+PrintCompilation

# Print inlining decisions
-XX:+UnlockDiagnosticVMOptions -XX:+PrintInlining

# Disable JIT (interpreter only — for debugging)
-Xint

# Compile everything immediately (no interpretation — for benchmarking)
-Xcomp

# Code Cache size (where native code is stored)
-XX:ReservedCodeCacheSize=256m
```

---

## 3. GC Tuning Flags

### Heap Sizing

```bash
-Xms512m                  # Initial heap size (ms = memory start)
-Xmx4g                    # Maximum heap size (mx = memory max)
                          # Best practice: set Xms == Xmx in production
                          # (avoids heap resize pauses)

-Xss512k                  # Per-thread stack size (ss = stack size)
                          # Reduce if creating many threads
```

### Choosing a GC

```bash
-XX:+UseSerialGC          # Single-threaded GC — only for small apps / embedded
-XX:+UseParallelGC        # Throughput-focused, multi-threaded — batch workloads
-XX:+UseG1GC              # Balanced latency+throughput — default Java 9+
-XX:+UseZGC               # Ultra-low latency < 1ms — large heaps (Java 15+ prod)
-XX:+UseShenandoahGC      # Low-pause, concurrent — RedHat's alternative to ZGC
```

### G1GC Tuning (most common in production)

```bash
-XX:+UseG1GC
-XX:MaxGCPauseMillis=200        # Target max pause (G1 tries to meet this — soft goal)
-XX:G1HeapRegionSize=16m        # Region size: 1,2,4,8,16,32 MB
                                # Set higher for large heaps (> 32GB → use 32m)
-XX:G1NewSizePercent=20         # Min % of heap for young generation
-XX:G1MaxNewSizePercent=60      # Max % of heap for young generation
-XX:InitiatingHeapOccupancyPercent=45  # Start concurrent GC cycle at 45% heap use
-XX:G1ReservePercent=10         # Reserve 10% headroom to avoid evacuation failures
```

### Generation Sizing

```bash
-XX:NewRatio=2            # Old:Young = 2:1 → young gen = 1/3 of heap
-XX:SurvivorRatio=8       # Eden:Survivor = 8:1 (each survivor = 1/10 of young)
-XX:MaxTenuringThreshold=15  # Promotions to old gen after 15 minor GCs
```

### Metaspace

```bash
-XX:MetaspaceSize=128m         # Initial metaspace commit (triggers GC when reached)
-XX:MaxMetaspaceSize=512m      # Cap metaspace (unlimited by default — risky)
```

### GC Logging (Java 9+ unified logging)

```bash
-Xlog:gc                           # Basic GC events
-Xlog:gc*                          # All GC details
-Xlog:gc*:file=gc.log:time,uptime  # Write to file with timestamps
-Xlog:gc+heap=debug                # Heap size before/after each GC
```

### Diagnostic & Crash

```bash
-XX:+HeapDumpOnOutOfMemoryError          # Auto-dump heap on OOM
-XX:HeapDumpPath=/tmp/heap-dump.hprof    # Dump file location
-XX:+ExitOnOutOfMemoryError              # Crash instead of limping (for containers)
-XX:+CrashOnOutOfMemoryError             # Force crash with full core dump
-XX:OnOutOfMemoryError="kill -9 %p"      # Run shell command on OOM
```

### Quick Reference Card

```bash
# Typical production JVM flags for a Spring Boot service
java \
  -server \
  -Xms2g -Xmx2g \
  -XX:+UseG1GC \
  -XX:MaxGCPauseMillis=200 \
  -XX:+HeapDumpOnOutOfMemoryError \
  -XX:HeapDumpPath=/var/log/app/heap.hprof \
  -Xlog:gc*:file=/var/log/app/gc.log:time,uptime:filecount=5,filesize=20m \
  -jar app.jar
```

---

## 4. Memory Leak Identification & Heap Dump Analysis

### What Is a Memory Leak in Java?

Objects that are **no longer needed** but still held by a reference chain from GC roots — the GC cannot collect them.

```
GC Root (static field)
    └─► Cache (static Map)
            └─► Entry (old data, never removed)
                    └─► Large object graph  ← never collected = memory leak
```

---

### Common Memory Leak Patterns

#### 1. Static Collections That Grow Forever

```java
// LEAK — static cache never evicted
class Registry {
    private static final Map<String, byte[]> CACHE = new HashMap<>();

    public static void register(String key, byte[] data) {
        CACHE.put(key, data);   // entries added but never removed
    }
}

// FIX — use a bounded cache or WeakReference values
private static final Map<String, WeakReference<byte[]>> CACHE = new HashMap<>();
// or use a proper LRU cache: new LinkedHashMap<>(100, 0.75f, true) { removeEldestEntry }
```

#### 2. Unclosed Resources

```java
// LEAK — streams / connections left open
public String readFile(String path) throws IOException {
    InputStream is = new FileInputStream(path);  // never closed if exception thrown
    return new String(is.readAllBytes());
}

// FIX — try-with-resources
public String readFile(String path) throws IOException {
    try (InputStream is = new FileInputStream(path)) {
        return new String(is.readAllBytes());
    }
}
```

#### 3. Listeners / Callbacks Never Removed

```java
// LEAK — registered listener holds reference to subscriber object
eventBus.register(myListener);    // myListener kept alive by eventBus
// ... myListener goes out of scope but eventBus still references it

// FIX — always deregister
eventBus.unregister(myListener);
```

#### 4. ThreadLocal Not Cleared in Thread Pools

```java
// LEAK — thread pool reuses threads; ThreadLocal value persists between requests
ThreadLocal<UserSession> SESSION = new ThreadLocal<>();

void handleRequest() {
    SESSION.set(new UserSession(...));
    processRequest();
    // forgot SESSION.remove() → session leaks to next request on this thread
}

// FIX
try {
    SESSION.set(new UserSession(...));
    processRequest();
} finally {
    SESSION.remove();   // always clean up
}
```

---

### Detecting Memory Leaks

#### Step 1 — Monitor Heap Growth

```bash
# Command-line heap stats every 1s (GC events)
jstat -gcutil <pid> 1000

# Output columns:
# S0    S1    E      O      M     YGC  YGCT  FGC  FGCT   GCT
# 0.00  0.00  72.49  94.12  ...   12   0.3   3    2.1    2.4
# Old gen (O) keeps growing after full GCs → likely leak
```

#### Step 2 — Take a Heap Dump

```bash
# Via jmap
jmap -dump:format=b,file=heap.hprof <pid>

# Via jcmd (preferred, safer)
jcmd <pid> GC.heap_dump /tmp/heap.hprof

# Automatically on OOM (add to JVM flags)
-XX:+HeapDumpOnOutOfMemoryError -XX:HeapDumpPath=/tmp/heap.hprof
```

#### Step 3 — Analyze with Eclipse MAT or VisualVM

Key views in **Eclipse Memory Analyzer (MAT)**:
- **Dominator Tree** — which objects retain the most heap (biggest suspects)
- **Leak Suspects Report** — auto-detects likely leaks
- **Object Query Language (OQL)** — query like SQL: `SELECT * FROM java.util.HashMap`
- **Retained Heap** — memory freed if this object were collected

```
Leak Suspects output example:
  Problem 1: 847 instances of com.example.UserSession
  occupy 238 MB (73% of heap)
  retained by: static field Registry.CACHE
  → never evicted from cache
```

#### Useful JVM Tools

```bash
jps                          # list running JVM processes
jstack <pid>                 # thread dump (deadlock detection)
jmap -histo <pid>            # histogram of objects in heap (quick check)
jcmd <pid> VM.native_memory  # native memory usage
jcmd <pid> GC.run            # trigger GC manually
```

---

## 5. StackOverflowError vs OutOfMemoryError

Both are `Error` subclasses (not `Exception`) — signal **JVM-level** resource exhaustion.

---

### `StackOverflowError`

Thrown when the **thread stack** runs out of space — typically caused by deep or infinite recursion.

Each method call pushes a **stack frame** (local vars, operand stack, return address).  
Stack size is fixed per thread (`-Xss`).

```java
// Classic cause — unbounded recursion
public int factorial(int n) {
    return n * factorial(n - 1);   // no base case → StackOverflowError
}

// Also triggers with very deep (but finite) call chains on small stacks
public void a() { b(); }
public void b() { c(); }
// ... 10,000 levels deep
```

```
Exception in thread "main" java.lang.StackOverflowError
    at com.example.App.factorial(App.java:5)
    at com.example.App.factorial(App.java:5)   ← repeated hundreds of times
    ...
```

**Fix:**
```java
// 1. Add base case
public long factorial(int n) {
    if (n <= 1) return 1;          // base case
    return n * factorial(n - 1);
}

// 2. Convert recursion to iteration
public long factorial(int n) {
    long result = 1;
    for (int i = 2; i <= n; i++) result *= i;
    return result;
}

// 3. Increase stack size (last resort)
// -Xss2m   (default is usually 512k–1m)
// Or per-thread: new Thread(null, runnable, "name", 2 * 1024 * 1024).start();
```

---

### `OutOfMemoryError`

Thrown when the JVM **cannot allocate an object** because the memory region is exhausted.  
Different subtypes point to different memory areas:

#### `java.lang.OutOfMemoryError: Java heap space`
Heap is full — GC cannot reclaim enough memory.

```java
// Cause: allocating more than -Xmx allows
List<byte[]> list = new ArrayList<>();
while (true) {
    list.add(new byte[1024 * 1024]);   // keep allocating 1MB chunks
}
// → OutOfMemoryError: Java heap space
```

```bash
# Diagnosis
-XX:+HeapDumpOnOutOfMemoryError
# Then analyze heap.hprof with MAT → find what's holding memory
```

#### `java.lang.OutOfMemoryError: Metaspace`
Class metadata area is full — too many classes loaded (e.g., dynamic proxy generation, hot deployment).

```java
// Cause: generating classes at runtime without bound
while (true) {
    new ByteBuddy().subclass(Object.class).make()
        .load(getClass().getClassLoader())
        .getLoaded();  // creates a new class each iteration
}
// → OutOfMemoryError: Metaspace
```

```bash
# Fix: cap and tune Metaspace
-XX:MaxMetaspaceSize=256m
```

#### `java.lang.OutOfMemoryError: GC overhead limit exceeded`
JVM is spending **> 98% of time in GC** but recovering **< 2% of heap**.  
Signals that the heap is effectively full and GC is making no progress.

```bash
# Disable this check (not recommended — hides the real problem)
-XX:-UseGCOverheadLimit
```

#### `java.lang.OutOfMemoryError: unable to create native thread`
OS cannot create more threads — either the process hit its thread limit or OS-level limits.

```java
// Cause: creating threads without bound
while (true) {
    new Thread(() -> { try { Thread.sleep(Long.MAX_VALUE); } catch ... }).start();
}
// → OutOfMemoryError: unable to create native thread
```

```bash
# Linux fix: increase per-process thread limit
ulimit -u 65535
# Or reduce -Xss to allow more threads within the same virtual address space
```

---

### Side-by-Side Comparison

| | `StackOverflowError` | `OutOfMemoryError` |
|---|---|---|
| Memory area | Thread Stack | Heap / Metaspace / Native |
| Common cause | Infinite / deep recursion | Memory leak, too-small heap |
| Catchable | Yes (but rarely useful) | Yes (but usually fatal) |
| JVM flag | `-Xss` (stack size) | `-Xmx`, `-XX:MaxMetaspaceSize` |
| Recoverable | Sometimes (per-thread) | Rarely — usually restart |
| Detection | Stack trace repeats same method | Heap dump, `jstat -gcutil` |

```java
// Catching errors — only do this if you have a very specific recovery plan
try {
    deepRecursion();
} catch (StackOverflowError e) {
    System.err.println("Stack overflow — returning fallback result");
    return fallback;
}

// For OOM — almost never appropriate to catch; log and let the process die
try {
    loadHugeData();
} catch (OutOfMemoryError e) {
    log.error("OOM — initiating graceful shutdown", e);
    System.exit(1);
}
```
