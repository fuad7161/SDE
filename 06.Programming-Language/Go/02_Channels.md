# 02 — Channels

> **Topics:** buffered vs unbuffered · directional channels · select · closing

---

## Table of Contents
1. [Buffered vs Unbuffered](#1-buffered-vs-unbuffered)
2. [Sending to a Closed Channel](#2-sending-to-a-closed-channel)
3. [select Statement](#3-select-statement)
4. [chan struct{}](#4-chan-struct)
5. [Timeout with Channels](#5-timeout-with-channels)
6. [Interview Questions](#interview-questions)

---

## 1. Buffered vs Unbuffered

| | Unbuffered | Buffered |
|---|---|---|
| Created with | `make(chan T)` | `make(chan T, n)` |
| Sender blocks until | Receiver is ready | Buffer is full |
| Receiver blocks until | Sender sends | Buffer is non-empty |
| Use case | Synchronization | Decoupling producer/consumer |

```go
package main

import "fmt"

func main() {
    // --- Unbuffered ---
    uch := make(chan int)
    go func() { uch <- 42 }() // sender blocks until receiver is ready
    fmt.Println(<-uch)        // 42

    // --- Buffered ---
    bch := make(chan int, 3)
    bch <- 1 // doesn't block — buffer has room
    bch <- 2
    bch <- 3
    // bch <- 4 // would block — buffer full

    fmt.Println(<-bch) // 1
    fmt.Println(<-bch) // 2
    fmt.Println(<-bch) // 3
}
```

### Directional channels (restrict capability)

```go
// send-only channel
func producer(ch chan<- int) {
    ch <- 100
}

// receive-only channel
func consumer(ch <-chan int) {
    fmt.Println(<-ch)
}

func main() {
    ch := make(chan int, 1)
    producer(ch) // bidirectional implicitly converts to chan<-
    consumer(ch) // bidirectional implicitly converts to <-chan
}
```

---

## 2. Sending to a Closed Channel

| Operation | On closed channel | Result |
|---|---|---|
| Send (`ch <- v`) | Panics | `panic: send on closed channel` |
| Receive (`v := <-ch`) | Returns zero value | `v = zero, ok = false` |
| Range over channel | Exits loop | after last buffered value |

```go
package main

import "fmt"

func main() {
    ch := make(chan int, 3)
    ch <- 10
    ch <- 20
    close(ch)

    // Safe receive with ok idiom
    for {
        v, ok := <-ch
        if !ok {
            fmt.Println("channel closed")
            break
        }
        fmt.Println(v) // 10, then 20
    }

    // Idiomatic: range closes automatically when channel is closed
    ch2 := make(chan string, 2)
    ch2 <- "a"
    ch2 <- "b"
    close(ch2)

    for s := range ch2 {
        fmt.Println(s) // a, b
    }
}
```

> **Rule:** Only the **sender** should close a channel. Closing from the receiver side or closing twice panics.

---

## 3. select Statement

`select` waits on multiple channel operations. If **multiple cases are ready simultaneously**, Go picks one **uniformly at random**.

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    ch1 := make(chan string, 1)
    ch2 := make(chan string, 1)

    ch1 <- "one"
    ch2 <- "two"

    // Both ready — Go picks randomly
    select {
    case msg := <-ch1:
        fmt.Println("ch1:", msg)
    case msg := <-ch2:
        fmt.Println("ch2:", msg)
    }
}
```

### Non-blocking send/receive with default

```go
ch := make(chan int, 1)

select {
case ch <- 42:
    fmt.Println("sent")
default:
    fmt.Println("channel full, skipping") // executes if ch is full
}
```

### Fan-in pattern — merge two channels into one

```go
func fanIn(ch1, ch2 <-chan string) <-chan string {
    merged := make(chan string)
    go func() {
        defer close(merged)
        for {
            select {
            case v, ok := <-ch1:
                if !ok {
                    ch1 = nil // disable this case
                }
                merged <- v
            case v, ok := <-ch2:
                if !ok {
                    ch2 = nil
                }
                merged <- v
            }
            if ch1 == nil && ch2 == nil {
                return
            }
        }
    }()
    return merged
}
```

---

## 4. chan struct{}

`chan struct{}` is used as a **signal channel** — communicating an event without data. `struct{}` has zero size so it allocates no memory.

```go
package main

import (
    "fmt"
    "time"
)

func worker(done <-chan struct{}) {
    for {
        select {
        case <-done:
            fmt.Println("worker stopped")
            return
        default:
            fmt.Println("working...")
            time.Sleep(500 * time.Millisecond)
        }
    }
}

func main() {
    done := make(chan struct{})
    go worker(done)

    time.Sleep(1500 * time.Millisecond)
    close(done) // broadcast stop signal to all goroutines listening on done
    time.Sleep(100 * time.Millisecond)
}
```

Common uses:
- Goroutine cancellation signals
- Semaphores: `make(chan struct{}, N)` limits N concurrent workers
- Set membership: `map[string]struct{}`

---

## 5. Timeout with Channels

```go
package main

import (
    "fmt"
    "time"
)

func fetchData() <-chan string {
    ch := make(chan string, 1)
    go func() {
        time.Sleep(2 * time.Second) // simulate slow work
        ch <- "result"
    }()
    return ch
}

func main() {
    ch := fetchData()

    select {
    case result := <-ch:
        fmt.Println("got:", result)
    case <-time.After(1 * time.Second):
        fmt.Println("timeout — operation took too long")
    }
}
```

> In production prefer `context.WithTimeout` over `time.After` — it integrates with the whole call chain and doesn't leak the timer goroutine.

```go
// Production-grade timeout
ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
defer cancel()

select {
case result := <-fetchData():
    fmt.Println("got:", result)
case <-ctx.Done():
    fmt.Println("timeout:", ctx.Err())
}
```

---

## Interview Questions

<details>
<summary><b>What is the difference between a buffered and unbuffered channel?</b></summary>

An unbuffered channel (`make(chan T)`) synchronizes sender and receiver — the sender blocks until a receiver is ready, and vice versa. A buffered channel (`make(chan T, n)`) allows up to n sends without a receiver; the sender only blocks when the buffer is full. Use unbuffered for synchronization guarantees; use buffered to decouple producer and consumer throughput.

</details>

<details>
<summary><b>What happens if you send to a closed channel? What about receiving?</b></summary>

Sending to a closed channel causes a **panic**. Receiving from a closed channel returns immediately with the zero value and `ok = false`. If there are buffered values remaining, they are drained first before returning zero values. Only the sender should close a channel.

</details>

<details>
<summary><b>How does select work when multiple cases are ready?</b></summary>

`select` evaluates all cases simultaneously. If multiple are ready, Go picks one **uniformly at random** — this is intentional to prevent starvation. If no case is ready and there is a `default` clause, it executes immediately (non-blocking). Without `default`, `select` blocks until at least one case is ready.

</details>

<details>
<summary><b>When would you use a chan struct{} specifically?</b></summary>

When you only need to signal an event, not transmit data. `struct{}` has zero size, so it costs no memory. Common uses: cancellation/done signals, semaphores (`make(chan struct{}, N)`), and set membership in maps (`map[K]struct{}`).

</details>

<details>
<summary><b>How do you implement a timeout using channels?</b></summary>

Use `select` with `time.After(d)` for quick scripts, or `context.WithTimeout` in production. `time.After` creates a timer goroutine that leaks if the operation completes before the timeout fires. `context.WithTimeout` is preferred because cancellation propagates through the entire call chain and the timer is cleaned up via `defer cancel()`.

</details>
