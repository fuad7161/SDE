# Multithreading & Concurrency — In-Depth Notes

---

## Table of Contents

1. [Thread Lifecycle, Runnable vs Callable](#1-thread-lifecycle-runnable-vs-callable)
2. [synchronized, volatile, Atomic Variables](#2-synchronized-volatile-atomic-variables)
3. [ReentrantLock, ReadWriteLock](#3-reentrantlock-readwritelock)
4. [ThreadLocal](#4-threadlocal)
5. [ExecutorService, ThreadPoolExecutor, Thread Pool Sizing](#5-executorservice-threadpoolexecutor-thread-pool-sizing)
6. [CompletableFuture](#6-completablefuture)
7. [Deadlock, Livelock, Starvation — Detection & Prevention](#7-deadlock-livelock-starvation--detection--prevention)
8. [Happens-Before Guarantee](#8-happens-before-guarantee)
9. [wait() / notify() vs Condition](#9-wait--notify-vs-condition)

---

## 1. Thread Lifecycle, Runnable vs Callable

### Thread Lifecycle (States)

```
               start()
NEW ──────────────────────► RUNNABLE
                               │   ▲
              scheduler picks │   │ time-slice ends
                               ▼   │
                            RUNNING
                           /   |   \
              sleep/wait/ /    |    \ run() completes
              I/O block  /     |     \
                        ▼      |      ▼
                    BLOCKED/   |    TERMINATED
                    WAITING/   |
                  TIMED_WAIT   │
                        │      │
                        └──────┘  lock acquired / notified
```

| State | Cause |
|---|---|
| `NEW` | Thread created but `start()` not called |
| `RUNNABLE` | Ready to run or currently running |
| `BLOCKED` | Waiting to acquire a `synchronized` lock |
| `WAITING` | Waiting indefinitely — `wait()`, `join()`, `park()` |
| `TIMED_WAITING` | Waiting with timeout — `sleep()`, `wait(ms)`, `join(ms)` |
| `TERMINATED` | `run()` completed or exception thrown |

```java
Thread t = new Thread(() -> System.out.println("running"));
System.out.println(t.getState()); // NEW
t.start();
System.out.println(t.getState()); // RUNNABLE
t.join();
System.out.println(t.getState()); // TERMINATED
```

---

### Runnable vs Callable

```java
// Runnable — no return value, cannot throw checked exceptions
Runnable runnable = () -> {
    System.out.println("Task running in: " + Thread.currentThread().getName());
};

Thread t = new Thread(runnable);
t.start();

// Callable — returns a value, can throw checked exceptions
Callable<Integer> callable = () -> {
    Thread.sleep(100);
    return 42;
};

ExecutorService executor = Executors.newSingleThreadExecutor();
Future<Integer> future = executor.submit(callable);
Integer result = future.get();   // blocks until result is ready
System.out.println(result);      // 42
executor.shutdown();
```

| | `Runnable` | `Callable<V>` |
|---|---|---|
| Return type | `void` | `V` (generic) |
| Checked exceptions | Cannot throw | Can throw |
| Used with | `Thread`, `ExecutorService` | `ExecutorService` only |
| Result via | — | `Future<V>` |

---

## 2. synchronized, volatile, Atomic Variables

### `synchronized`

Ensures **mutual exclusion** — only one thread can execute the block/method at a time.  
Acquiring a monitor lock also establishes a **happens-before** relationship.

```java
class Counter {
    private int count = 0;

    // synchronized method — locks on 'this'
    public synchronized void increment() {
        count++;
    }

    // synchronized block — more granular locking
    public void incrementBlock() {
        synchronized (this) {
            count++;
        }
    }

    // static synchronized — locks on Counter.class
    public static synchronized void staticMethod() { ... }

    public int getCount() { return count; }
}

Counter counter = new Counter();
Thread t1 = new Thread(() -> { for (int i = 0; i < 1000; i++) counter.increment(); });
Thread t2 = new Thread(() -> { for (int i = 0; i < 1000; i++) counter.increment(); });
t1.start(); t2.start();
t1.join();  t2.join();
System.out.println(counter.getCount()); // always 2000
```

**Without `synchronized`**, `count++` is **not atomic** — it's 3 operations:
1. Read `count`
2. Increment
3. Write back

Two threads can interleave these, causing lost updates (race condition).

---

### `volatile`

Guarantees **visibility** of changes across threads — reads/writes go directly to main memory, bypassing CPU cache.  
Does **not** guarantee atomicity for compound operations like `i++`.

```java
class FlagExample {
    private volatile boolean running = true;  // without volatile, thread may cache 'true' forever

    public void stop() {
        running = false;   // write visible to all threads immediately
    }

    public void run() {
        while (running) {  // always reads fresh value from main memory
            doWork();
        }
    }
}
```

**Double-checked locking** with `volatile` (correct singleton pattern):

```java
class Singleton {
    private static volatile Singleton instance;

    public static Singleton getInstance() {
        if (instance == null) {                    // first check (no lock)
            synchronized (Singleton.class) {
                if (instance == null) {            // second check (with lock)
                    instance = new Singleton();    // volatile prevents reordering
                }
            }
        }
        return instance;
    }
}
```

Without `volatile`, another thread could see a **partially constructed** object due to instruction reordering.

---

### Atomic Variables (`java.util.concurrent.atomic`)

Lock-free thread safety using **CAS (Compare-And-Swap)** CPU instructions.

```java
import java.util.concurrent.atomic.*;

AtomicInteger atomicCount = new AtomicInteger(0);

// All operations are atomic — no synchronized needed
atomicCount.incrementAndGet();       // ++i
atomicCount.getAndIncrement();       // i++
atomicCount.addAndGet(5);            // i += 5
atomicCount.compareAndSet(10, 20);   // if value==10, set to 20; returns boolean

// AtomicReference — atomic object reference update
AtomicReference<String> ref = new AtomicReference<>("initial");
ref.compareAndSet("initial", "updated");

// AtomicLong, AtomicBoolean follow same pattern

// LongAdder — better than AtomicLong under high contention
LongAdder adder = new LongAdder();
adder.increment();
adder.add(5);
long total = adder.sum();
```

| | `synchronized` | `volatile` | `AtomicInteger` |
|---|---|---|---|
| Mutual exclusion | Yes | No | No (CAS-based) |
| Visibility | Yes | Yes | Yes |
| Atomicity | Yes | No (for i++) | Yes |
| Performance | Lower (blocking) | High | High (lock-free) |
| Use for | Complex critical sections | Simple flags/references | Counters, single vars |

---

## 3. ReentrantLock, ReadWriteLock

### `ReentrantLock`

More flexible than `synchronized` — supports **timed**, **interruptible**, and **fair** lock acquisition.

```java
import java.util.concurrent.locks.ReentrantLock;

class SafeCounter {
    private int count = 0;
    private final ReentrantLock lock = new ReentrantLock();

    public void increment() {
        lock.lock();
        try {
            count++;
        } finally {
            lock.unlock();   // always in finally to prevent deadlock
        }
    }

    public void tryIncrement() throws InterruptedException {
        // try to acquire lock, wait up to 500ms
        if (lock.tryLock(500, TimeUnit.MILLISECONDS)) {
            try {
                count++;
            } finally {
                lock.unlock();
            }
        } else {
            System.out.println("Could not acquire lock");
        }
    }
}
```

**Reentrancy** — same thread can acquire the lock multiple times without deadlocking:

```java
ReentrantLock lock = new ReentrantLock();

void outer() {
    lock.lock();
    try {
        inner();   // same thread re-acquires lock — works fine
    } finally {
        lock.unlock();
    }
}

void inner() {
    lock.lock();   // hold count becomes 2
    try { ... }
    finally { lock.unlock(); }  // hold count back to 1
}
```

---

### `ReadWriteLock`

Allows **multiple concurrent readers** but **exclusive access for writers**.  
Ideal when reads are far more frequent than writes.

```java
import java.util.concurrent.locks.ReadWriteLock;
import java.util.concurrent.locks.ReentrantReadWriteLock;

class Cache {
    private final Map<String, String> data = new HashMap<>();
    private final ReadWriteLock rwLock = new ReentrantReadWriteLock();

    public String read(String key) {
        rwLock.readLock().lock();         // multiple threads can hold read lock
        try {
            return data.get(key);
        } finally {
            rwLock.readLock().unlock();
        }
    }

    public void write(String key, String value) {
        rwLock.writeLock().lock();        // exclusive — blocks all readers & writers
        try {
            data.put(key, value);
        } finally {
            rwLock.writeLock().unlock();
        }
    }
}
```

| Scenario | Read Lock | Write Lock |
|---|---|---|
| No locks held | Acquire ✅ | Acquire ✅ |
| Read lock(s) held | Acquire ✅ | Block ❌ |
| Write lock held | Block ❌ | Block ❌ |

### `synchronized` vs `ReentrantLock`

| Feature | `synchronized` | `ReentrantLock` |
|---|---|---|
| Syntax | Built-in keyword | Explicit lock/unlock |
| Timed lock | No | `tryLock(time, unit)` |
| Interruptible | No | `lockInterruptibly()` |
| Fairness | No | `new ReentrantLock(true)` |
| Condition variables | `wait/notify` | `newCondition()` |
| `finally` required | No | **Yes** |

---

## 4. ThreadLocal

Provides **thread-local variables** — each thread has its own independent copy.  
No synchronization needed since no sharing occurs.

```java
class UserContext {
    // Each thread gets its own copy
    private static final ThreadLocal<String> currentUser = new ThreadLocal<>();

    public static void setUser(String user) { currentUser.set(user); }
    public static String getUser()          { return currentUser.get(); }
    public static void clear()              { currentUser.remove(); }
}

// Thread 1
new Thread(() -> {
    UserContext.setUser("Alice");
    // ... do work
    System.out.println(UserContext.getUser()); // "Alice"
    UserContext.clear();
}).start();

// Thread 2 — completely independent copy
new Thread(() -> {
    UserContext.setUser("Bob");
    System.out.println(UserContext.getUser()); // "Bob"
    UserContext.clear();
}).start();
```

**With initial value:**

```java
ThreadLocal<SimpleDateFormat> dateFormat = ThreadLocal.withInitial(
    () -> new SimpleDateFormat("yyyy-MM-dd")
);
// SimpleDateFormat is not thread-safe — ThreadLocal gives each thread its own instance
String formatted = dateFormat.get().format(new Date());
```

### Common Use Cases

- Per-request user/session info in web frameworks
- Database connection / transaction context
- `SimpleDateFormat` (not thread-safe) per thread

### Memory Leak Warning

Always call `remove()` when the thread is returned to a pool (`ThreadPoolExecutor`).  
Otherwise the value persists across requests from the same thread.

```java
try {
    UserContext.setUser(resolveCurrentUser());
    handleRequest();
} finally {
    UserContext.clear();  // critical in thread pool environments
}
```

---

## 5. ExecutorService, ThreadPoolExecutor, Thread Pool Sizing

### ExecutorService

Decouples task submission from thread management:

```java
// Fixed thread pool — reuses N threads
ExecutorService fixed = Executors.newFixedThreadPool(4);

// Single thread — sequential execution
ExecutorService single = Executors.newSingleThreadExecutor();

// Cached pool — creates threads on demand, reuses idle ones (unbounded)
ExecutorService cached = Executors.newCachedThreadPool();

// Scheduled — run tasks after delay or periodically
ScheduledExecutorService scheduled = Executors.newScheduledThreadPool(2);

// Submit tasks
fixed.execute(() -> System.out.println("fire and forget"));

Future<String> future = fixed.submit(() -> "result");
String result = future.get(5, TimeUnit.SECONDS);  // timeout

// Always shut down
fixed.shutdown();                     // waits for running tasks
fixed.shutdownNow();                  // interrupts running tasks
fixed.awaitTermination(10, TimeUnit.SECONDS);
```

---

### `ThreadPoolExecutor` (Full Control)

```java
ThreadPoolExecutor executor = new ThreadPoolExecutor(
    2,                                // corePoolSize    — always-alive threads
    10,                               // maximumPoolSize — max threads under load
    60L, TimeUnit.SECONDS,            // keepAliveTime   — idle thread timeout
    new LinkedBlockingQueue<>(100),   // workQueue       — task buffer
    Executors.defaultThreadFactory(), // threadFactory
    new ThreadPoolExecutor.CallerRunsPolicy() // rejectionHandler
);
```

**Task flow:**

```
submit(task)
    │
    ├─ core threads available?  ──YES──► assign to core thread
    │
    ├─ NO → queue full?          ──NO──► enqueue task
    │
    ├─ YES → max threads reached?─NO──► create new thread (up to max)
    │
    └─ YES ──────────────────────────► RejectionHandler
```

**Rejection Policies:**

| Policy | Behavior |
|---|---|
| `AbortPolicy` (default) | Throws `RejectedExecutionException` |
| `CallerRunsPolicy` | Caller thread runs the task (natural back-pressure) |
| `DiscardPolicy` | Silently drops the task |
| `DiscardOldestPolicy` | Drops oldest queued task, retries submission |

---

### Thread Pool Sizing Rules of Thumb

```
CPU-bound tasks:
  pool size = number of CPU cores + 1
  (extra thread for occasional stalls — page fault, GC pause)

I/O-bound tasks:
  pool size = cores × (1 + wait time / compute time)
  e.g., if task is 90% waiting: pool size = 8 × (1 + 9) = 80

Mixed:
  profile first, then tune
```

```java
int cores = Runtime.getRuntime().availableProcessors();
int ioPool  = cores * 10;   // high I/O
int cpuPool = cores + 1;    // CPU-bound
```

---

## 6. CompletableFuture

Non-blocking async programming with chainable callbacks (Java 8+):

```java
// Basic async task
CompletableFuture<String> future = CompletableFuture.supplyAsync(() -> {
    // runs in ForkJoinPool.commonPool() by default
    return fetchDataFromDB();
});

// Chain transformations (non-blocking)
CompletableFuture<String> result = future
    .thenApply(data -> data.toUpperCase())          // transform result
    .thenApply(data -> "Processed: " + data);

// Chain another async task
CompletableFuture<String> chained = future
    .thenCompose(data -> CompletableFuture.supplyAsync(() -> callApi(data)));

// Combine two independent futures
CompletableFuture<String> cf1 = CompletableFuture.supplyAsync(() -> "Hello");
CompletableFuture<String> cf2 = CompletableFuture.supplyAsync(() -> "World");

CompletableFuture<String> combined = cf1.thenCombine(cf2, (a, b) -> a + " " + b);
System.out.println(combined.get()); // "Hello World"

// Wait for all to complete
CompletableFuture<Void> all = CompletableFuture.allOf(cf1, cf2);
all.join();

// First to complete wins
CompletableFuture<Object> any = CompletableFuture.anyOf(cf1, cf2);

// Exception handling
CompletableFuture<String> safe = future
    .exceptionally(ex -> "Fallback: " + ex.getMessage())
    .handle((result2, ex) -> ex != null ? "Error" : result2);  // handle both cases

// Custom thread pool
ExecutorService pool = Executors.newFixedThreadPool(4);
CompletableFuture.supplyAsync(() -> compute(), pool)
    .thenApplyAsync(r -> transform(r), pool)
    .thenAccept(System.out::println);
```

### Key Methods Summary

| Method | Purpose |
|---|---|
| `supplyAsync(Supplier)` | Start async task with return value |
| `runAsync(Runnable)` | Start async task, no return |
| `thenApply(Function)` | Transform result (sync) |
| `thenApplyAsync(Function)` | Transform result (async) |
| `thenAccept(Consumer)` | Consume result, no return |
| `thenRun(Runnable)` | Run after completion, no access to result |
| `thenCompose(Function)` | Chain another `CompletableFuture` (flatMap) |
| `thenCombine(CF, BiFunction)` | Combine two independent results |
| `allOf(CF...)` | Wait for all |
| `anyOf(CF...)` | Wait for first |
| `exceptionally(Function)` | Handle exception, provide fallback |
| `handle(BiFunction)` | Handle both result and exception |

---

## 7. Deadlock, Livelock, Starvation — Detection & Prevention

### Deadlock

Occurs when two or more threads **each hold a lock the other needs** — circular wait.

```
Thread 1 holds Lock A, waits for Lock B
Thread 2 holds Lock B, waits for Lock A
→ Neither can proceed — deadlock
```

```java
// Classic deadlock scenario
Object lockA = new Object();
Object lockB = new Object();

Thread t1 = new Thread(() -> {
    synchronized (lockA) {
        System.out.println("T1 holds A, waiting for B");
        synchronized (lockB) { System.out.println("T1 holds A and B"); }
    }
});

Thread t2 = new Thread(() -> {
    synchronized (lockB) {              // opposite order!
        System.out.println("T2 holds B, waiting for A");
        synchronized (lockA) { System.out.println("T2 holds A and B"); }
    }
});

t1.start(); t2.start();  // deadlock!
```

**Prevention — consistent lock ordering:**

```java
// Both threads acquire in same order: lockA → lockB
Thread t1 = new Thread(() -> {
    synchronized (lockA) { synchronized (lockB) { ... } }
});
Thread t2 = new Thread(() -> {
    synchronized (lockA) { synchronized (lockB) { ... } }  // same order
});
```

**Detection:** Use `jstack <pid>` or `ThreadMXBean`:

```java
ThreadMXBean bean = ManagementFactory.getThreadMXBean();
long[] deadlocked = bean.findDeadlockedThreads();
if (deadlocked != null) {
    ThreadInfo[] infos = bean.getThreadInfo(deadlocked, true, true);
    for (ThreadInfo info : infos) System.out.println(info);
}
```

---

### Livelock

Threads are **actively responding** to each other but **making no progress** — like two people in a corridor both stepping aside in the same direction.

```java
// Simplified livelock concept
class Spoon {
    private Diner owner;
    public synchronized void use() {
        // politely yield spoon if partner is hungry — both keep yielding forever
        while (owner.isHungry() && owner.getPartner().isHungry()) {
            // yield, then re-check — neither eats
            owner = owner.getPartner();
        }
    }
}
```

**Prevention**: Add randomized backoff or a priority scheme so one thread proceeds first.

---

### Starvation

A thread is **perpetually denied** CPU time because higher-priority threads always get scheduled first.

```java
// Fair lock prevents starvation — threads acquire in arrival order
ReentrantLock fairLock = new ReentrantLock(true);  // fair=true
```

**Prevention:**
- Use fair locks (`new ReentrantLock(true)`)
- Avoid thread priorities for critical correctness
- Use `ThreadPoolExecutor` with bounded queues

---

## 8. Happens-Before Guarantee

Defined by the **Java Memory Model (JMM)** — if action A *happens-before* action B, then A's effects are **visible** to B.

### Key Happens-Before Rules

```java
// 1. Thread start — all actions before start() happen-before thread's first action
int x = 10;
Thread t = new Thread(() -> System.out.println(x)); // guaranteed to see x=10
t.start();

// 2. Thread join — all actions of thread happen-before join() returns
t.join();
// everything t did is visible here

// 3. synchronized — unlock happens-before subsequent lock of same monitor
synchronized (lock) { sharedVar = 42; }
// another thread: synchronized (lock) { assert sharedVar == 42; } ← guaranteed

// 4. volatile write happens-before subsequent volatile read
volatile int flag = 0;
// Thread A: flag = 1;
// Thread B: if (flag == 1) { ... }  ← guaranteed to see flag=1 after A writes

// 5. Static initializer happens-before first use of the class
class Config {
    static final String VALUE = loadConfig(); // safely published
}
```

### Why It Matters

Without happens-before, the JIT compiler and CPU are free to **reorder** instructions for performance. Happens-before is the contract that prevents you from seeing stale or partially-written data.

```java
// Without volatile — t2 may never see ready=true (CPU cache, reordering)
boolean ready = false;
int data = 0;

// Thread 1
data = 42;
ready = true;  // may be reordered BEFORE data=42

// Thread 2
if (ready) System.out.println(data);  // could print 0!

// FIX — volatile on ready establishes happens-before
volatile boolean ready = false;
```

---

## 9. wait() / notify() vs Condition

### `wait()` / `notify()` — Classic Object Monitor

`wait()` releases the monitor lock and suspends the thread until `notify()` or `notifyAll()` is called.  
**Must** be called inside a `synchronized` block.

```java
class BlockingBuffer {
    private final Queue<Integer> queue = new LinkedList<>();
    private final int capacity;

    BlockingBuffer(int capacity) { this.capacity = capacity; }

    public synchronized void produce(int item) throws InterruptedException {
        while (queue.size() == capacity) {
            wait();             // releases lock, waits for space
        }
        queue.add(item);
        notifyAll();            // wake up consumers
    }

    public synchronized int consume() throws InterruptedException {
        while (queue.isEmpty()) {
            wait();             // releases lock, waits for item
        }
        int item = queue.poll();
        notifyAll();            // wake up producers
        return item;
    }
}
```

**Why `while` instead of `if`?** — Spurious wakeups: a thread can wake up without being notified.  
Always re-check the condition in a `while` loop.

---

### `Condition` — with `ReentrantLock`

More flexible: multiple distinct conditions per lock, supports timed waits, interruptible waits.

```java
import java.util.concurrent.locks.*;

class BetterBuffer {
    private final Queue<Integer> queue = new LinkedList<>();
    private final int capacity;
    private final ReentrantLock lock = new ReentrantLock();
    private final Condition notFull  = lock.newCondition();  // separate conditions
    private final Condition notEmpty = lock.newCondition();

    BetterBuffer(int capacity) { this.capacity = capacity; }

    public void produce(int item) throws InterruptedException {
        lock.lock();
        try {
            while (queue.size() == capacity) {
                notFull.await();          // wait for space
            }
            queue.add(item);
            notEmpty.signal();            // wake only a consumer (targeted)
        } finally {
            lock.unlock();
        }
    }

    public int consume() throws InterruptedException {
        lock.lock();
        try {
            while (queue.isEmpty()) {
                notEmpty.await();         // wait for item
            }
            int item = queue.poll();
            notFull.signal();             // wake only a producer (targeted)
            return item;
        } finally {
            lock.unlock();
        }
    }
}
```

**Advantage**: `notEmpty.signal()` wakes only a consumer — `notifyAll()` would wake all threads (producers + consumers) unnecessarily.

---

### wait()/notify() vs Condition

| Feature | `wait()` / `notify()` | `Condition` |
|---|---|---|
| Associated with | `synchronized` block (Object monitor) | `ReentrantLock` |
| Multiple wait sets | No — one per object | Yes — multiple `Condition` per lock |
| Timed wait | `wait(ms)` | `await(time, unit)` |
| Interruptible | Yes | `awaitUninterruptibly()` available |
| Targeted signal | `notify()` (random 1), `notifyAll()` | `signal()` (targeted), `signalAll()` |
| Spurious wakeup | Must guard with `while` | Must guard with `while` |
| Flexibility | Less | More |
