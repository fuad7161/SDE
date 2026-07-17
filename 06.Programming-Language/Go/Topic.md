# Go Interview Topics

> Mid-Senior level reference. Each topic has a dedicated file with explanations, examples, and interview Q&A.

---

## Topics

| # | Topic | Key Subtopics |
|---|---|---|
| 01 | [Goroutines & Concurrency](./01_Goroutines_Concurrency.md) | goroutine lifecycle, GMP scheduler, goroutine leaks |
| 02 | [Channels](./02_Channels.md) | buffered vs unbuffered, directional, select, closing |
| 03 | [Interfaces & Type System](./03_Interfaces_TypeSystem.md) | implicit implementation, empty interface, type assertion, embedding |
| 04 | [Error Handling](./04_Error_Handling.md) | error interface, errors.Is/As, wrapping, sentinel errors |
| 05 | [Memory Management & Pointers](./05_Memory_Pointers.md) | stack vs heap, escape analysis, GC |
| 06 | [Slices & Maps Internals](./06_Slices_Maps.md) | slice header, append, map internals, nil maps |
| 07 | [defer, panic, recover](./07_Defer_Panic_Recover.md) | execution order, named returns, recover pattern |
| 08 | [Sync Primitives](./08_Sync_Primitives.md) | Mutex, WaitGroup, Once, Pool, atomic |
| 09 | [Context Package](./09_Context.md) | cancellation, timeout, deadline, value propagation |
| 10 | [Generics](./10_Generics.md) | type parameters, constraints, comparable |
| 11 | [Testing](./11_Testing.md) | table-driven tests, benchmarks, mocking, t.Parallel |
| 12 | [Package & Module System](./12_Packages_Modules.md) | go.mod, go.sum, build tags, init() order |
| 13 | [Input Stream Buffering](./13_Input_Stream_Buffering.md) | bufio, Scanner, Reader, performance, io.Reader |

---

## Quick Reference — Depth Expected at Mid-Level

| Area | Depth Expected |
|---|---|
| Goroutines + Channels | Deep — concurrency bugs, patterns |
| Interface design | Medium — nil traps, composition |
| Error wrapping | Medium — `%w`, `errors.Is/As` |
| Context propagation | Deep — cancellation chains |
| Slice/Map internals | Medium — gotchas, performance |
| Testing | Medium — table-driven, benchmarks |
| Generics | Basic to Medium |
