# Multithreading & Concurrency

---

## Table of Contents

1. [Thread Lifecycle](#1-thread-lifecycle)
2. [Thread vs Runnable](#2-thread-vs-runnable)
3. [Synchronization](#3-synchronization)
4. [`volatile` Keyword](#4-volatile-keyword)
5. [Deadlock & Race Condition](#5-deadlock--race-condition)
6. [Inter-Thread Communication](#6-inter-thread-communication)
7. [Executor Framework](#7-executor-framework)
8. [Callable & Future](#8-callable--future)
9. [Concurrent Collections](#9-concurrent-collections)

---

## 1. Thread Lifecycle

```
                   start()
NEW ──────────────────────────► RUNNABLE
                                  │
                                  │ CPU schedules thread
                                  ▼
                              RUNNING
                              /   |   \
         wait()/sleep()/     /    |    \  run() completes
         block on I/O      ▼     |     ▼
                        WAITING  |   TERMINATED (Dead)
                           │     |
              notify()/    │     │ synchronized block available
              notifyAll()  │     │ sleep ends / I/O done
                           ▼     ▼
                          RUNNABLE (back in queue)
```

| State | Description |
|---|---|
| **NEW** | Thread created but `start()` not yet called |
| **RUNNABLE** | Ready to run; waiting for CPU time |
| **RUNNING** | Actually executing on CPU |
| **BLOCKED** | Waiting to acquire a monitor lock (synchronized) |
| **WAITING** | Waiting indefinitely (`wait()`, `join()`, `park()`) |
| **TIMED_WAITING** | Waiting with timeout (`sleep(n)`, `wait(n)`, `join(n)`) |
| **TERMINATED** | `run()` has finished or exception thrown |

```java
Thread t = new Thread(() -> System.out.println("running"));
System.out.println(t.getState());   // NEW
t.start();
System.out.println(t.getState());   // RUNNABLE or RUNNING
t.join();
System.out.println(t.getState());   // TERMINATED
```

> **Interview Q: What are the different states of a thread?**  
> NEW → RUNNABLE → RUNNING → (BLOCKED/WAITING/TIMED_WAITING) → TERMINATED. A thread enters BLOCKED when it tries to acquire a monitor lock held by another thread. It enters WAITING when it calls `wait()` or `join()` without timeout. TIMED_WAITING is like WAITING but with a timeout (`sleep(ms)`, `wait(ms)`). Once `run()` completes or an uncaught exception occurs, the thread goes to TERMINATED and cannot be restarted.

---

## 2. Thread vs Runnable

```java
// ── METHOD 1: Extend Thread ──
class PrintTask extends Thread {
    private final String message;

    PrintTask(String message) {
        this.message = message;
    }

    @Override
    public void run() {
        for (int i = 0; i < 3; i++) {
            System.out.println(Thread.currentThread().getName() + ": " + message);
        }
    }
}

PrintTask t1 = new PrintTask("Hello");
PrintTask t2 = new PrintTask("World");
t1.start();    // starts new OS thread — calls run() internally
t2.start();
// t1.run();   // ❌ does NOT start new thread — runs on current thread

// ── METHOD 2: Implement Runnable (PREFERRED) ──
class CounterTask implements Runnable {
    private final String name;

    CounterTask(String name) { this.name = name; }

    @Override
    public void run() {
        for (int i = 1; i <= 3; i++) {
            System.out.println(name + ": " + i);
        }
    }
}

Thread thread = new Thread(new CounterTask("Counter"));
thread.start();

// With lambda (Java 8+)
Thread t = new Thread(() -> System.out.println("Lambda thread!"));
t.start();

// ── THREAD METHODS ──
Thread t = new Thread(() -> doWork());
t.setName("WorkerThread");
t.setDaemon(true);          // daemon threads die when all non-daemon threads finish
t.setPriority(Thread.MAX_PRIORITY);  // 1-10, default 5

t.start();
t.join();           // wait for t to finish
t.join(1000);       // wait max 1000ms
Thread.sleep(500);  // pause current thread 500ms (throws InterruptedException)
Thread.currentThread().getName();   // get current thread name
```

| | `Thread` | `Runnable` |
|---|---|---|
| Coupling | Tight (IS-A Thread) | Loose (HAS-A task) |
| Multiple inheritance | ❌ (can't extend anything else) | ✅ (can implement other interfaces) |
| Reuse | Task tied to thread | Same task submitted to any executor |
| Modern usage | Rarely | Preferred (+ `Callable`) |

> **Interview Q: What is the difference between extending `Thread` and implementing `Runnable`?**  
> Extending `Thread` is limiting because Java doesn't support multiple inheritance — your class can't extend anything else. Implementing `Runnable` **separates the task** (what to do) from the thread (how to execute it), allowing the same task to run in a thread pool, scheduled executor, or new thread without changes. `Runnable` is always preferred; `Thread` extension is only useful if you need to override Thread's behavior itself.

---

## 3. Synchronization

```java
// ── RACE CONDITION WITHOUT SYNCHRONIZATION ──
class UnsafeCounter {
    int count = 0;

    void increment() {
        count++;   // NOT atomic: read count → add 1 → write back
        // Two threads can both read "5", both write "6" → lost update
    }
}

// ── SYNCHRONIZED METHOD ──
class SafeCounter {
    private int count = 0;

    synchronized void increment() {    // acquires lock on 'this'
        count++;
    }

    synchronized int getCount() {
        return count;
    }
}

// ── SYNCHRONIZED BLOCK — finer granularity ──
class BankAccount {
    private double balance;
    private final Object lock = new Object();   // dedicated lock object

    void transfer(BankAccount target, double amount) {
        synchronized (this.lock) {          // minimal critical section
            if (this.balance < amount) throw new IllegalStateException("Insufficient funds");
            this.balance -= amount;
        }
        synchronized (target.lock) {
            target.balance += amount;
        }
    }

    // Static synchronized — locks on Class object
    static synchronized void printTotal() {
        // only one thread can call this across ALL instances
    }
}

// ── REENTRANT LOCK — explicit, more flexible ──
import java.util.concurrent.locks.ReentrantLock;

class FlexibleCounter {
    private int count = 0;
    private final ReentrantLock lock = new ReentrantLock(true); // fair lock

    void increment() {
        lock.lock();
        try {
            count++;
        } finally {
            lock.unlock();   // ALWAYS unlock in finally
        }
    }

    boolean tryIncrement() {
        if (lock.tryLock()) {     // non-blocking attempt
            try {
                count++;
                return true;
            } finally {
                lock.unlock();
            }
        }
        return false;   // couldn't get lock, did nothing
    }
}
```

> **Interview Q: What is synchronization? What is the difference between synchronized method and synchronized block?**  
> Synchronization ensures only one thread executes a **critical section** at a time by acquiring a **monitor lock** (intrinsic lock). A **synchronized method** acquires the lock on `this` (or the Class object for static methods) — the entire method is the critical section. A **synchronized block** lets you specify exactly which object to lock on and which code to protect, giving **finer granularity** — less contention and better performance. Prefer synchronized blocks to minimize the locked region.

---

## 4. `volatile` Keyword

```java
// Problem: without volatile, each thread may cache variable in CPU register
// Changes by one thread invisible to other threads (visibility problem)

class StopFlag {
    // ❌ WITHOUT volatile — thread may never see flag = true
    // boolean running = true;

    // ✅ WITH volatile — all reads/writes go directly to main memory
    volatile boolean running = true;

    void stop() {
        running = false;    // write to main memory immediately
    }

    void run() {
        while (running) {   // always reads from main memory
            // do work
        }
        System.out.println("Stopped");
    }
}

// volatile guarantees:
// 1. VISIBILITY  — write by one thread is immediately visible to all threads
// 2. ORDERING    — no reordering across a volatile access (happens-before)

// volatile does NOT guarantee:
// - Atomicity of compound operations (count++ is NOT safe with volatile alone)
class UnsafeWithVolatile {
    volatile int count = 0;

    void increment() {
        count++;   // still not atomic: read-modify-write (3 operations)
    }
}

// For atomic compound operations, use AtomicInteger:
import java.util.concurrent.atomic.AtomicInteger;

class SafeCounter {
    AtomicInteger count = new AtomicInteger(0);

    void increment() {
        count.incrementAndGet();   // atomic CAS operation
    }
}
```

| | `volatile` | `synchronized` |
|---|---|---|
| Guarantees | Visibility + ordering | Visibility + atomicity + mutual exclusion |
| Compound ops | ❌ Not atomic | ✅ Atomic |
| Performance | Faster (no lock) | Slower (lock acquisition) |
| Use for | Simple flags, status variables | Critical sections with multiple ops |

> **Interview Q: What is the `volatile` keyword and when should you use it?**  
> `volatile` tells the JVM to always read a variable from **main memory** and write directly to main memory, bypassing CPU-level caching. It guarantees **visibility** (writes by one thread are immediately visible to all others) and prevents instruction **reordering** across that variable. Use it for simple **flag variables** (like `running = false` to stop a thread). Do NOT use it when you need atomicity for compound operations (like `count++`) — use `AtomicInteger` or `synchronized` instead.

---

## 5. Deadlock & Race Condition

### Deadlock

```java
// Deadlock — two threads each waiting for a lock the other holds
class DeadlockExample {
    static Object lock1 = new Object();
    static Object lock2 = new Object();

    static void thread1() {
        synchronized (lock1) {
            System.out.println("T1: got lock1, waiting for lock2");
            try { Thread.sleep(100); } catch (InterruptedException e) { }
            synchronized (lock2) {    // BLOCKED — T2 holds lock2
                System.out.println("T1: got both locks");
            }
        }
    }

    static void thread2() {
        synchronized (lock2) {
            System.out.println("T2: got lock2, waiting for lock1");
            synchronized (lock1) {    // BLOCKED — T1 holds lock1
                System.out.println("T2: got both locks");
            }
        }
    }
}
// Both threads wait forever — DEADLOCK

// ── Prevention: always acquire locks in the SAME ORDER ──
static void thread1Safe() {
    synchronized (lock1) {
        synchronized (lock2) { /* ... */ }   // same order
    }
}
static void thread2Safe() {
    synchronized (lock1) {                   // same order as thread1
        synchronized (lock2) { /* ... */ }
    }
}

// Detection: jstack <pid> shows "Found one Java-level deadlock"
```

### Race Condition

```java
// Race condition — outcome depends on thread scheduling order
class TicketSystem {
    private int availableTickets = 1;

    void bookTicket(String customer) {
        // Two threads can both see availableTickets == 1
        if (availableTickets > 0) {
            // context switch here → both proceed
            availableTickets--;
            System.out.println(customer + " booked last ticket");
        }
    }
}

// Fix: synchronize the check-then-act
synchronized void bookTicket(String customer) {
    if (availableTickets > 0) {
        availableTickets--;
        System.out.println(customer + " booked last ticket");
    } else {
        System.out.println(customer + ": no tickets available");
    }
}
```

> **Interview Q: What are the four conditions for deadlock?**  
> **Coffman conditions** — all four must hold simultaneously:  
> 1. **Mutual Exclusion** — resource held by only one thread at a time  
> 2. **Hold and Wait** — thread holds resource while waiting for another  
> 3. **No Preemption** — resource can only be released voluntarily  
> 4. **Circular Wait** — thread A waits for B, B waits for A (cycle)  
> **Prevention**: break any one condition — most practical: enforce **lock ordering** (break circular wait) or use `tryLock()` with timeout (break hold-and-wait).

---

## 6. Inter-Thread Communication

```java
// wait() and notify() — must be called inside synchronized block
class SharedBuffer {
    private int data;
    private boolean hasData = false;

    // Producer — sets data and notifies consumer
    synchronized void produce(int value) throws InterruptedException {
        while (hasData) {
            wait();          // releases lock, suspends until notify()
        }
        data = value;
        hasData = true;
        System.out.println("Produced: " + value);
        notify();            // wakes up ONE waiting thread
        // notifyAll();     // wakes up ALL waiting threads
    }

    // Consumer — waits for data then reads it
    synchronized int consume() throws InterruptedException {
        while (!hasData) {
            wait();          // releases lock, suspends until notify()
        }
        hasData = false;
        System.out.println("Consumed: " + data);
        notify();            // notify producer
        return data;
    }
}

// Usage
SharedBuffer buffer = new SharedBuffer();

Thread producer = new Thread(() -> {
    for (int i = 1; i <= 5; i++) {
        try { buffer.produce(i); } catch (InterruptedException e) { Thread.currentThread().interrupt(); }
    }
});

Thread consumer = new Thread(() -> {
    for (int i = 0; i < 5; i++) {
        try { buffer.consume(); } catch (InterruptedException e) { Thread.currentThread().interrupt(); }
    }
});

producer.start();
consumer.start();
```

> **Interview Q: Why must `wait()` and `notify()` be called inside a `synchronized` block?**  
> Because they operate on the **object's monitor lock**. `wait()` releases the lock and suspends the thread; `notify()` signals a waiting thread to reacquire the lock. If called outside `synchronized`, the thread doesn't hold the monitor, so the JVM throws `IllegalMonitorStateException`. Also, the check-then-wait pattern (`while (!condition) { wait(); }`) must be atomic — without `synchronized`, there's a race between checking the condition and calling `wait()`.

---

## 7. Executor Framework

```java
import java.util.concurrent.*;

// ── THREAD POOL TYPES ──

// Fixed pool — N threads, reused across tasks
ExecutorService fixed = Executors.newFixedThreadPool(4);

// Single thread — tasks execute sequentially
ExecutorService single = Executors.newSingleThreadExecutor();

// Cached pool — creates threads as needed, reuses idle ones (no limit!)
ExecutorService cached = Executors.newCachedThreadPool();

// Scheduled — run after delay or periodically
ScheduledExecutorService scheduled = Executors.newScheduledThreadPool(2);

// ── SUBMITTING TASKS ──
ExecutorService pool = Executors.newFixedThreadPool(4);

// submit Runnable — no result
pool.execute(() -> System.out.println("Task 1: " + Thread.currentThread().getName()));
pool.submit(() -> System.out.println("Task 2"));    // returns Future<?>

// ── SHUTDOWN ──
pool.shutdown();          // stop accepting new tasks, wait for running tasks to finish
pool.shutdownNow();       // interrupt running tasks, return list of pending tasks
pool.awaitTermination(10, TimeUnit.SECONDS);   // wait up to 10s for completion

// ── SCHEDULED EXECUTOR ──
scheduled.schedule(() -> System.out.println("Delayed by 2s"), 2, TimeUnit.SECONDS);
scheduled.scheduleAtFixedRate(() -> System.out.println("Every 1s"), 0, 1, TimeUnit.SECONDS);
scheduled.scheduleWithFixedDelay(() -> System.out.println("Delay after completion"), 0, 1, TimeUnit.SECONDS);

// ── CUSTOM THREAD POOL ──
ThreadPoolExecutor executor = new ThreadPoolExecutor(
    2,                          // corePoolSize
    10,                         // maximumPoolSize
    60, TimeUnit.SECONDS,       // keepAliveTime
    new LinkedBlockingQueue<>(100),  // work queue
    new ThreadFactory() {            // custom thread factory
        int n = 0;
        @Override
        public Thread newThread(Runnable r) {
            return new Thread(r, "Worker-" + n++);
        }
    },
    new ThreadPoolExecutor.CallerRunsPolicy()  // rejection policy
);
```

> **Interview Q: What are the different types of thread pools provided by `Executors`?**  
> `newFixedThreadPool(n)` — fixed N threads, tasks queue up when all busy; good for CPU-bound work. `newCachedThreadPool()` — grows unboundedly; threads idle for 60s then die; good for short-lived tasks but can overwhelm memory. `newSingleThreadExecutor()` — one thread, serial execution, tasks never overlap. `newScheduledThreadPool(n)` — supports delayed and periodic execution. **In production**, prefer `ThreadPoolExecutor` directly with explicit queue bounds and rejection policies to avoid unbounded queues.

---

## 8. Callable & Future

```java
// Runnable — no return value, can't throw checked exceptions
// Callable — has return value, can throw checked exceptions

ExecutorService pool = Executors.newFixedThreadPool(3);

// ── CALLABLE — returns a result ──
Callable<Integer> task = () -> {
    Thread.sleep(1000);   // simulate computation
    return 42;
};

Future<Integer> future = pool.submit(task);

// Do other work while task runs...
System.out.println("Task submitted, doing other work...");

// Get result — BLOCKS until result is ready
Integer result = future.get();           // waits indefinitely
// future.get(5, TimeUnit.SECONDS);      // timeout version
System.out.println("Result: " + result); // 42

// ── FUTURE METHODS ──
future.isDone();           // true if completed (normally or with exception)
future.isCancelled();      // true if cancelled
future.cancel(true);       // attempt to cancel; true = interrupt if running

// ── COMPLETABLEFUTURE (Java 8) — async pipeline ──
CompletableFuture<String> cf = CompletableFuture
    .supplyAsync(() -> fetchUser(1))                   // async: fetch user
    .thenApply(user -> user.getName().toUpperCase())   // transform result
    .thenApply(name -> "Hello, " + name)              // chain transform
    .exceptionally(ex -> "Error: " + ex.getMessage()); // error handling

cf.thenAccept(System.out::println);   // consume result
String result2 = cf.get();            // block and get

// ── COMBINING FUTURES ──
CompletableFuture<String> f1 = CompletableFuture.supplyAsync(() -> "Hello");
CompletableFuture<String> f2 = CompletableFuture.supplyAsync(() -> "World");

CompletableFuture<String> combined = f1.thenCombine(f2, (s1, s2) -> s1 + " " + s2);
System.out.println(combined.get());   // "Hello World"

// Run all in parallel, wait for all
CompletableFuture<Void> all = CompletableFuture.allOf(f1, f2);
all.get();   // blocks until both done
```

> **Interview Q: What is the difference between `Runnable` and `Callable`?**  
> `Runnable.run()` returns `void` and cannot throw checked exceptions. `Callable.call()` returns a typed result (`T`) and can throw checked exceptions. `Callable` is submitted to an `ExecutorService` via `submit()`, which returns a `Future<T>` that you can use to retrieve the result, check completion status, or cancel the task. Use `Callable` whenever you need a result from a concurrent task.

---

## 9. Concurrent Collections

```java
// java.util.concurrent — thread-safe without explicit synchronization

// ── ConcurrentHashMap ── (prefer over Hashtable/synchronized Map)
ConcurrentHashMap<String, Integer> chm = new ConcurrentHashMap<>();
chm.put("a", 1);
chm.putIfAbsent("b", 2);
chm.computeIfAbsent("c", k -> k.length());   // atomic: compute if absent
chm.merge("a", 1, Integer::sum);             // atomic: merge with existing

// ── CopyOnWriteArrayList ── (reads > writes; e.g., event listeners)
CopyOnWriteArrayList<String> cowList = new CopyOnWriteArrayList<>();
cowList.add("listener1");
// Iteration never throws ConcurrentModificationException
for (String s : cowList) {
    cowList.add("new");   // safe: iterates original snapshot
}

// ── BlockingQueue — producer-consumer ──
BlockingQueue<Integer> queue = new LinkedBlockingQueue<>(10);  // bounded

// Producer thread
new Thread(() -> {
    for (int i = 1; i <= 20; i++) {
        try {
            queue.put(i);   // blocks when queue full
            System.out.println("Produced: " + i);
        } catch (InterruptedException e) { Thread.currentThread().interrupt(); }
    }
}).start();

// Consumer thread
new Thread(() -> {
    while (true) {
        try {
            int item = queue.take();   // blocks when queue empty
            System.out.println("Consumed: " + item);
        } catch (InterruptedException e) { Thread.currentThread().interrupt(); break; }
    }
}).start();

// ── CountDownLatch — wait for N tasks to complete ──
CountDownLatch latch = new CountDownLatch(3);
for (int i = 0; i < 3; i++) {
    new Thread(() -> {
        doWork();
        latch.countDown();   // decrement count
    }).start();
}
latch.await();   // blocks until count reaches 0

// ── CyclicBarrier — synchronize N threads at a checkpoint ──
CyclicBarrier barrier = new CyclicBarrier(3, () -> System.out.println("All arrived!"));
for (int i = 0; i < 3; i++) {
    new Thread(() -> {
        doPhase1();
        barrier.await();   // wait for all 3 threads to reach this point
        doPhase2();
    }).start();
}
```

| Collection | Thread-safe | Use case |
|---|---|---|
| `ConcurrentHashMap` | ✅ bucket-level | General purpose map |
| `CopyOnWriteArrayList` | ✅ copy on write | Read-heavy, few writes |
| `LinkedBlockingQueue` | ✅ two locks | Producer-consumer |
| `PriorityBlockingQueue` | ✅ | Priority producer-consumer |
| `ConcurrentLinkedQueue` | ✅ CAS | Non-blocking queue |

> **Interview Q: What is `CountDownLatch` and how is it different from `CyclicBarrier`?**  
> `CountDownLatch` is initialized with a count N. Threads call `countDown()` when done; a waiting thread calls `await()` which blocks until the count reaches 0. It is **one-shot** — cannot be reset. `CyclicBarrier` synchronizes N threads at a **common barrier point** — all N threads call `await()`, and all are released together when the last one arrives. It can be **reused** (cyclic) for the next phase, and an optional `Runnable` action runs when all threads reach the barrier. Use `CountDownLatch` for "wait for N tasks to finish"; use `CyclicBarrier` for "synchronize N threads at phase boundaries".
