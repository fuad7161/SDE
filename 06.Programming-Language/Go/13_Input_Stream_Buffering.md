# 13. Input Stream Buffering

## Overview
In Go, reading from I/O streams (like files, network connections, or `os.Stdin`) directly can be inefficient due to the overhead of frequent system calls for every small read. Input stream buffering solves this by reading larger blocks of data into a memory buffer all at once, and then serving subsequent smaller reads from that buffer. The standard library provides the `bufio` package to wrap an `io.Reader` or `io.Writer` and provide buffering.

## Key Subtopics

### 1. The `bufio` Package
The `bufio` package implements buffered I/O. It wraps an `io.Reader` or `io.Writer` object, creating another object (`Reader` or `Writer`) that also implements the interface but provides buffering and some help for textual I/O.

### 2. `bufio.Reader`
A `bufio.Reader` wraps an `io.Reader` and provides buffered operations.
- **Initialization**: `reader := bufio.NewReader(os.Stdin)` or `bufio.NewReaderSize(reader, size)`. Default size is usually 4KB.
- **Reading**: Offers methods like `Read`, `ReadByte`, `ReadBytes`, `ReadLine`, and `ReadString`.
- **Peeking**: `Peek(n int)` lets you look ahead at the next `n` bytes without advancing the reader.

### 3. `bufio.Scanner`
`bufio.Scanner` provides a convenient interface for reading data such as a file of newline-delimited lines of text. Successive calls to the `Scan` method will step through the 'tokens' of a file, skipping the bytes between the tokens.
- **Default behavior**: It splits the input into lines (`bufio.ScanLines`).
- **Custom Splitters**: You can customize it to split by words (`bufio.ScanWords`), runes (`bufio.ScanRunes`), bytes, or a custom split function.
- **Key Methods**: `Scan() bool`, `Text() string`, `Bytes() []byte`, `Err() error`.

### 4. Buffer Size and `ErrTooLong`
When using a `Scanner`, the default max token size is 64KB (`bufio.MaxScanTokenSize`). If a line (or token) exceeds this limit, the scanner will stop and return an `bufio.ErrTooLong` error.
- **Solution**: To handle larger lines, you need to increase the buffer capacity using `Scanner.Buffer(buf []byte, max int)`.

## Code Examples

### Example 1: `bufio.Scanner`
```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

func main() {
    input := "line 1\nline 2\nline 3"
    scanner := bufio.NewScanner(strings.NewReader(input))
    
    for scanner.Scan() {
        fmt.Println(scanner.Text())
    }
    
    if err := scanner.Err(); err != nil {
        fmt.Fprintln(os.Stderr, "reading standard input:", err)
    }
}
```

### Example 2: Handling Large Lines with `Scanner`
```go
package main

import (
    "bufio"
    "fmt"
    "strings"
)

func main() {
    // A long line exceeding default token size
    longLine := strings.Repeat("A", 100000) 
    scanner := bufio.NewScanner(strings.NewReader(longLine))
    
    // Create a larger buffer
    buf := make([]byte, 0, 64*1024)
    scanner.Buffer(buf, 1024*1024) // 1MB max token size
    
    for scanner.Scan() {
        fmt.Println("Scanned length:", len(scanner.Text()))
    }
}
```

## Potential Interview Questions

<details>
<summary><b>Q1: Why use `bufio` to read a file instead of raw `os.File` reads?</b></summary>

**Answer**: `os.File` reads directly from the OS file system. Calling it for each byte or small chunk forces a context switch and system call (`read`), which is slow. `bufio.Reader` mitigates this by making large chunk reads from the OS into a memory buffer, and then small reads are served instantly from memory, significantly enhancing performance.

</details>

<details>
<summary><b>Q2: What is the difference between `bufio.Reader` and `bufio.Scanner`?</b></summary>

**Answer**: 
- `bufio.Reader` is a lower-level construct that adds buffering to an `io.Reader`. It operates on streams of data and provides methods to read specific delimiters, bytes, or chunks. It is more flexible but can be more complex (e.g., handling `isPrefix` with `ReadLine`).
- `bufio.Scanner` is a higher-level abstraction designed specifically for parsing streams into distinct tokens (like lines or words). It is much easier to use but has limitations, such as keeping the entire token in memory (which can fail with `ErrTooLong` if the token exceeds max capacity).

</details>

<details>
<summary><b>Q3: You get a `bufio.Scanner: token too long` error. What happened and how do you fix it?</b></summary>

**Answer**: This happens when the scanner encounters a token (usually a line) that exceeds its maximum buffer size (default 64KB). It can be fixed by allocating a larger initial buffer and setting a larger capacity limit using the `scanner.Buffer(customBuf, maxCapacity)` method before invoking `scanner.Scan()`.

</details>

<details>
<summary><b>Q4: Should I use `io.ReadAll` or `bufio` for a large file?</b></summary>

**Answer**: If a file is extremely large, `io.ReadAll` could cause memory exhaustion or an Out of Memory (OOM) panic because it loads the entire file into memory at once. For large files, stream processing with `bufio.Scanner` or `bufio.Reader` is the strictly preferred approach to maintain a stable, low memory footprint.

</details>

<details>
<summary><b>Q5: The Core Trap: What is the danger of using `.Bytes()` inside a scanner loop?</b></summary>

**Answer**: `scanner.Bytes()` returns a slice pointing directly to the scanner's internal buffer. The contents of this slice will be overwritten on the next call to `Scan()`. If you need to keep the data across loop iterations or pass it to another goroutine, you must copy the slice's contents. Note that `scanner.Text()` is safe because it allocates a new string (strings are immutable in Go), effectively copying the underlying bytes.

</details>

<details>
<summary><b>Q6: Architectural Trade-off: When should you use `bufio.Scanner` vs. `bufio.Reader`?</b></summary>

**Answer**: Use `bufio.Scanner` for simple, token-based reading (like reading line-by-line or word-by-word) where tokens are reasonably sized. It provides a cleaner and simpler API. Use `bufio.Reader` when you need fine-grained control, are dealing with mixed data formats (binary and text), need to `Peek()` ahead without advancing the reader, or when handling arbitrarily long lines where you prefer granular chunked reading (via `ReadLine` or `ReadSlice`) rather than dealing with `Scanner` buffer resizing.

</details>

<details>
<summary><b>Q7: Deep Optimization: Why is buffered I/O faster than direct unbuffered I/O?</b></summary>

**Answer**: Direct unbuffered I/O (e.g., calling `os.File.Read` for single bytes) involves an expensive context switch from user space to kernel space for every single call. Buffered I/O performs a single system call to read a large contiguous chunk (e.g., 4KB) of data into a user-space memory buffer. Subsequent small reads are then fulfilled almost instantaneously from this local memory buffer without traversing the kernel boundary, drastically reducing system call overhead.

</details>
