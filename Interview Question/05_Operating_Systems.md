# Operating Systems

1. What is an operating system, and what are its main responsibilities?
   - An OS manages hardware resources and provides services such as process scheduling, memory management, files, devices, security, and user interfaces.
2. What is the difference between a program and a process?
   - A program is passive executable code, while a process is a running instance with memory, resources, and execution state.
3. What is the difference between a process and a thread?
   - Processes have separate address spaces; threads are execution units within a process that share its memory and resources.
4. What are the common states in a process lifecycle?
   - Common states are new, ready, running, waiting or blocked, and terminated.
5. What is a process control block?
   - A PCB is the kernel record containing a process's ID, state, registers, scheduling data, memory information, and resources.
6. What is a context switch, and why is it expensive?
   - The OS saves one task's CPU state and loads another's; this costs time and can disrupt caches and pipelines.
7. What is CPU scheduling?
   - CPU scheduling selects which ready process or thread runs next to balance responsiveness, throughput, fairness, and utilization.
8. Compare FCFS, shortest-job-first, priority, and round-robin scheduling.
   - FCFS follows arrival order, SJF favors short work, priority favors importance, and round-robin shares time in fixed quanta.
9. What is concurrency? How is it different from parallelism?
   - Concurrency manages overlapping tasks, while parallelism literally executes multiple tasks at the same instant.
10. What is a race condition?
   - A race occurs when a result depends on unpredictable timing of unsynchronized access to shared state.
11. What is a critical section?
   - It is code that accesses shared state and must be synchronized to prevent conflicting concurrent execution.
12. What is the difference between a mutex and a semaphore?
   - A mutex provides exclusive ownership of a resource; a semaphore uses a counter to permit a limited number of concurrent accesses or signal events.
13. What is a monitor in concurrency?
   - A monitor bundles shared state with mutually exclusive operations and condition variables for safe coordination.
14. What is deadlock?
   - Deadlock is a permanent wait in which tasks each hold resources needed by another task in the cycle.
15. What are the four necessary conditions for deadlock?
   - Mutual exclusion, hold and wait, no preemption, and circular wait must all be present.
16. How can deadlock be prevented, avoided, detected, or recovered from?
   - Break a necessary condition, allocate only safe states, detect wait cycles, or recover by preempting resources or terminating work.
17. What is starvation? How is it different from deadlock?
   - Starvation indefinitely denies one task resources while others progress; in deadlock, the involved tasks cannot progress at all.
18. What is memory management?
   - It tracks, allocates, protects, maps, and reclaims memory for the OS and running processes.
19. What is virtual memory?
   - Virtual memory gives each process a private logical address space mapped to physical memory and potentially secondary storage.
20. What is paging, and what is a page fault?
   - Paging divides memory into fixed-size pages and frames; a page fault occurs when a referenced page is not currently mapped in physical memory.
21. What is segmentation?
   - Segmentation divides memory into variable-size logical regions such as code, stack, and data, each with its own bounds and permissions.
22. What is the difference between internal and external fragmentation?
   - Internal fragmentation wastes space inside allocated blocks; external fragmentation leaves free space split into unusable gaps.
23. What is thrashing?
   - Thrashing occurs when excessive page faults make the system spend more time moving pages than executing useful work.
24. What is the difference between user mode and kernel mode?
   - User mode restricts privileged operations for safety; kernel mode gives the OS full hardware and memory access.
25. What is a system call?
   - A system call is the controlled interface through which user programs request kernel services.
26. What is inter-process communication? Name common IPC mechanisms.
   - IPC lets processes exchange data and coordinate through pipes, sockets, shared memory, message queues, signals, and similar mechanisms.
27. What is a file system?
   - A file system organizes persistent data into files and directories and manages names, metadata, permissions, storage, and retrieval.
28. What is the difference between a symbolic link and a hard link?
   - A symbolic link stores a path and may cross file systems; a hard link is another directory entry for the same underlying file data.

## Medium to Advanced

29. How does copy-on-write work after process creation?
   - **Key note:** Parent and child share read-only pages until one writes, when the OS copies only the modified page.
30. What is the difference between a kernel thread and a user-level thread?
   - **Key note:** The kernel schedules kernel threads; user threads are managed in a runtime and mapped onto kernel execution units.
31. What are green threads and coroutines?
   - **Key note:** They are user-space scheduled units that provide cheap concurrency, often yielding cooperatively around blocking work.
32. How does a system call transition from user mode to kernel mode?
   - **Key note:** A controlled trap switches privilege, validates arguments, runs a kernel handler, and returns the result.
33. What is the difference between blocking, non-blocking, synchronous, and asynchronous I/O?
   - **Key note:** Blocking concerns waiting of the caller; async concerns completion notification, while synchronous completes within the request flow.
34. Compare `select`, `poll`, and scalable event mechanisms such as `epoll` or `kqueue`.
   - **Key note:** `select/poll` repeatedly scan descriptors; scalable mechanisms report only ready changes and handle large sets efficiently.
35. How does virtual-to-physical address translation work?
   - **Key note:** The MMU uses page tables, accelerated by the TLB, to map virtual pages to physical frames.
36. What is a translation lookaside buffer?
   - **Key note:** A TLB caches recent page-table translations; misses require a more expensive page-table walk.
37. What is demand paging?
   - **Key note:** Pages are loaded only on first access, reducing initial memory use at the cost of possible page faults.
38. Compare page-replacement algorithms such as FIFO, LRU, and Clock.
   - **Key note:** They approximate which page is least valuable; Clock cheaply approximates LRU using reference bits.
39. What is memory-mapped I/O, and when is it useful?
   - **Key note:** Files or devices appear in virtual memory, enabling paging and simpler random access but requiring careful error handling.
40. What are zero-copy I/O techniques?
   - **Key note:** They reduce copying between kernel and user buffers through mechanisms such as `sendfile`, DMA, or buffer sharing.
41. How does a journaling file system improve crash recovery?
   - **Key note:** It records intended metadata or data changes in a log before applying them, enabling replay after failure.
42. What is CPU cache locality, and why does it affect performance?
   - **Key note:** Accessing nearby or recently used data reduces expensive memory fetches and cache misses.
43. What is false sharing?
   - **Key note:** Threads modify unrelated values on one cache line, causing costly cache-coherence invalidations.
44. What is priority inversion, and how can priority inheritance help?
   - **Key note:** A high-priority task waits on a low-priority lock holder; temporarily raising the holder's priority speeds release.
45. What is the thundering-herd problem?
   - **Key note:** Many waiters wake for one event and contend wastefully; targeted wakeups and work distribution reduce it.
46. How do namespaces and control groups provide container isolation?
   - **Key note:** Namespaces isolate resource views; cgroups account for and limit CPU, memory, and other resources.
47. What is NUMA, and why can memory placement matter?
   - **Key note:** Memory access is faster to a CPU's local node, so poor thread/data placement increases latency.
48. How would you diagnose high load average with low CPU utilization?
   - **Key note:** Check tasks blocked on disk/network I/O, uninterruptible sleep, lock contention, and resource saturation.
