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
