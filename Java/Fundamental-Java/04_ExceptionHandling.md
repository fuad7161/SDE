# Exception Handling

---

## Table of Contents

1. [Exception Hierarchy](#1-exception-hierarchy)
2. [Checked vs Unchecked Exceptions](#2-checked-vs-unchecked-exceptions)
3. [`try`, `catch`, `finally`](#3-try-catch-finally)
4. [`throw` vs `throws`](#4-throw-vs-throws)
5. [Custom Exceptions](#5-custom-exceptions)
6. [Multiple Catch Blocks & Multi-Catch](#6-multiple-catch-blocks--multi-catch)
7. [Try-with-Resources](#7-try-with-resources)

---

## 1. Exception Hierarchy

```
java.lang.Throwable
├── java.lang.Error                    ← Don't catch! Serious JVM problems
│   ├── OutOfMemoryError               (heap full)
│   ├── StackOverflowError             (infinite recursion)
│   ├── VirtualMachineError
│   └── AssertionError
│
└── java.lang.Exception
    ├── IOException                    ← CHECKED — must handle
    │   ├── FileNotFoundException
    │   ├── SocketException
    │   └── EOFException
    ├── SQLException                   ← CHECKED
    ├── ClassNotFoundException         ← CHECKED
    ├── CloneNotSupportedException     ← CHECKED
    │
    └── RuntimeException              ← UNCHECKED — optional handling
        ├── NullPointerException       (calling method on null)
        ├── ArrayIndexOutOfBoundsException
        ├── StringIndexOutOfBoundsException
        ├── ClassCastException         (invalid type cast)
        ├── NumberFormatException      (parseInt("abc"))
        ├── IllegalArgumentException
        ├── IllegalStateException
        ├── ArithmeticException        (divide by zero)
        └── UnsupportedOperationException
```

> **Interview Q: What is the difference between `Error` and `Exception`?**  
> Both extend `Throwable`, but they serve different purposes. **`Error`** represents **serious system-level problems** that the application generally cannot recover from (JVM out of memory, stack overflow, class file corrupted). You should **not** catch `Error` in normal code. **`Exception`** represents **application-level problems** that can potentially be caught and handled (file not found, network timeout, invalid input). The `try-catch` mechanism is designed for `Exception`.

---

## 2. Checked vs Unchecked Exceptions

### Checked Exceptions

The **compiler forces** you to handle them — either with `try-catch` or declaring with `throws`.

```java
// Compile error if not handled:
public void readFile(String path) {
    FileReader fr = new FileReader(path);   // ❌ Unhandled IOException
}

// Option 1: catch it
public void readFile(String path) {
    try {
        FileReader fr = new FileReader(path);
    } catch (FileNotFoundException e) {
        System.err.println("File not found: " + e.getMessage());
    }
}

// Option 2: declare it (pass responsibility to caller)
public void readFile(String path) throws FileNotFoundException {
    FileReader fr = new FileReader(path);   // ✅ caller must handle
}
```

### Unchecked Exceptions

No compile-time enforcement — you choose whether to handle them.

```java
public void divide(int a, int b) {
    // ArithmeticException is unchecked — no forced try-catch
    int result = a / b;   // throws ArithmeticException if b == 0
    System.out.println(result);
}

// Can optionally handle:
public int safeDivide(int a, int b) {
    if (b == 0) throw new IllegalArgumentException("Divisor cannot be zero");
    return a / b;
}
```

| | Checked | Unchecked |
|---|---|---|
| Extends | `Exception` (not `RuntimeException`) | `RuntimeException` or `Error` |
| Compiler forces handling | ✅ Yes | ❌ No |
| Represents | Recoverable external conditions | Programming bugs / logic errors |
| Examples | `IOException`, `SQLException` | `NullPointerException`, `IllegalArgumentException` |
| Best practice | External resources (files, DB, network) | Validate inputs; use for contract violations |

> **Interview Q: When should you use checked vs unchecked exceptions?**  
> Use **checked exceptions** when the caller can reasonably be expected to **recover** (e.g., file not found → show dialog to user). Use **unchecked exceptions** for **programming errors** that should not occur if the code is written correctly (e.g., null pointer, out of bounds, illegal argument). Modern practice (Spring, Effective Java) leans toward **unchecked** for cleaner APIs — checked exceptions pollute method signatures and force empty catch blocks.

---

## 3. `try`, `catch`, `finally`

```java
public class ExceptionDemo {

    static int divide(int a, int b) {
        try {
            System.out.println("Trying division...");
            int result = a / b;               // may throw ArithmeticException
            System.out.println("Result: " + result);
            return result;
        } catch (ArithmeticException e) {
            System.out.println("Caught: " + e.getMessage());  // "/ by zero"
            return -1;
        } finally {
            // ALWAYS runs — even if exception occurs, even if return is executed
            System.out.println("Finally block runs");
        }
    }

    public static void main(String[] args) {
        System.out.println(divide(10, 2));    // 5
        System.out.println("---");
        System.out.println(divide(10, 0));    // -1 (caught)
    }
}

// Output:
// Trying division...
// Result: 5
// Finally block runs
// 5
// ---
// Trying division...
// Caught: / by zero
// Finally block runs
// -1
```

### When does `finally` NOT run?

```java
try {
    System.exit(0);       // JVM terminates — finally does NOT run
} finally {
    System.out.println("This will NOT print");
}

// Also skipped if JVM crashes or thread is forcibly killed
```

### `finally` for Resource Cleanup

```java
Connection conn = null;
PreparedStatement ps = null;
try {
    conn = dataSource.getConnection();
    ps = conn.prepareStatement("SELECT * FROM users WHERE id = ?");
    ps.setInt(1, userId);
    ResultSet rs = ps.executeQuery();
    // process results...
} catch (SQLException e) {
    log.error("DB error", e);
} finally {
    // Clean up in reverse order of creation
    if (ps != null) { try { ps.close(); } catch (SQLException ignored) {} }
    if (conn != null) { try { conn.close(); } catch (SQLException ignored) {} }
}
```

> **Interview Q: Will `finally` always execute?**  
> Almost always. `finally` executes even if an exception is thrown, even if there's a `return` in `try` or `catch`. The **only cases** where `finally` does NOT execute: `System.exit()` is called, the JVM crashes (`kill -9`), or the thread is killed externally. In `return` from `try`, the value is saved, `finally` runs, then the saved value is returned — `finally` cannot change the return value (but a `return` in `finally` itself overrides the `try` return).

---

## 4. `throw` vs `throws`

```java
// ─── throws ─── used in method SIGNATURE — declares possible checked exceptions ───
public void connectToDatabase(String url) throws SQLException, IOException {
    // This method MIGHT throw these exceptions
    // Caller is warned and must handle them
}

// ─── throw ─── used inside method BODY — actually throws an exception ───
public void setAge(int age) {
    if (age < 0 || age > 150) {
        throw new IllegalArgumentException("Invalid age: " + age);  // actually thrown here
    }
    this.age = age;
}

// Rethrowing — catch, do something (log), then rethrow
public void processFile(String path) throws IOException {
    try {
        readFile(path);
    } catch (IOException e) {
        log.error("Failed to process file: " + path, e);
        throw e;    // rethrow same exception
    }
}

// Wrapping — catch low-level exception, throw domain exception
public User findUser(int id) {
    try {
        return userDao.findById(id);
    } catch (SQLException e) {
        // Don't expose DB internals to caller
        throw new ServiceException("Failed to find user: " + id, e);  // chain cause!
    }
}
```

| | `throw` | `throws` |
|---|---|---|
| Location | Inside method body | Method signature |
| Purpose | Actually throws an exception | Declares possible exceptions |
| Instance | One specific exception object | One or more exception types |
| For | Checked + Unchecked | Only checked (unchecked optional) |

> **Interview Q: What is the difference between `throw` and `throws`?**  
> `throw` is a **statement** that **instantiates and throws** an exception object right there in the code (`throw new IllegalArgumentException("msg")`). `throws` is a **declaration** in the method signature that warns callers this method **might** throw the listed checked exceptions, and they must handle them. You `throw` one exception at a time; you can `throws` multiple comma-separated exception types.

---

## 5. Custom Exceptions

```java
// ─── Custom Unchecked Exception (most common in modern Java) ───
public class ProductNotFoundException extends RuntimeException {
    private final String productId;

    // Minimum: message constructor
    public ProductNotFoundException(String productId) {
        super("Product not found: " + productId);
        this.productId = productId;
    }

    // Always provide cause constructor to preserve original stack trace
    public ProductNotFoundException(String productId, Throwable cause) {
        super("Product not found: " + productId, cause);
        this.productId = productId;
    }

    // Extra data that callers might need
    public String getProductId() { return productId; }
}

// ─── Custom Checked Exception ───
public class InsufficientStockException extends Exception {
    private final int requested;
    private final int available;

    public InsufficientStockException(int requested, int available) {
        super(String.format("Requested %d units but only %d available", requested, available));
        this.requested = requested;
        this.available = available;
    }

    public int getRequested() { return requested; }
    public int getAvailable() { return available; }
}

// ─── Usage ───
public class ProductService {

    public Product findProduct(String id) {
        Product p = productRepository.findById(id);
        if (p == null) {
            throw new ProductNotFoundException(id);   // unchecked — no try-catch forced
        }
        return p;
    }

    public void placeOrder(String productId, int qty) throws InsufficientStockException {
        Product p = findProduct(productId);
        if (p.getStock() < qty) {
            throw new InsufficientStockException(qty, p.getStock());   // checked
        }
        p.setStock(p.getStock() - qty);
    }
}

// ─── Handling custom exceptions ───
try {
    service.placeOrder("PROD-001", 100);
} catch (ProductNotFoundException e) {
    System.out.println("Product ID: " + e.getProductId());
} catch (InsufficientStockException e) {
    System.out.println("Need " + e.getRequested() + ", have " + e.getAvailable());
}
```

> **Interview Q: How do you create a custom exception? What are the best practices?**  
> Extend `RuntimeException` for unchecked, `Exception` for checked. Best practices:  
> 1. Provide a **message constructor** and a **message + cause constructor** (always preserve cause)  
> 2. Include **domain-relevant fields** for structured error handling  
> 3. Name it clearly: `OrderNotFoundException` not `MyException`  
> 4. **Never swallow exceptions**: always log or rethrow  
> 5. When wrapping a lower-level exception, always pass it as `cause` to preserve the full stack trace

---

## 6. Multiple Catch Blocks & Multi-Catch

```java
// Multiple catch — most specific FIRST, most general LAST
public void processInput(String input) {
    try {
        int number = Integer.parseInt(input);   // NumberFormatException
        int[] arr = new int[number];            // NegativeArraySizeException
        arr[number - 1] = 100;                  // ArrayIndexOutOfBoundsException

    } catch (NumberFormatException e) {
        System.out.println("Not a valid number: " + input);
    } catch (NegativeArraySizeException e) {
        System.out.println("Number must be positive");
    } catch (ArrayIndexOutOfBoundsException e) {
        System.out.println("Array access error: " + e.getMessage());
    } catch (RuntimeException e) {
        System.out.println("Other runtime error: " + e.getMessage());
    } catch (Exception e) {         // must be last — most general
        System.out.println("Unexpected error: " + e.getMessage());
    }
}

// ── Multi-catch (Java 7+) — handle multiple types the same way ──
public void readData(String path) {
    try {
        // code that can throw either
    } catch (IOException | SQLException e) {
        // e is effectively final here — can't reassign
        log.error("Data access error: " + e.getMessage(), e);
        throw new DataAccessException("Failed to read data", e);
    }
}
```

**Exception chaining — always preserve cause:**

```java
try {
    // low-level operation
} catch (SQLException e) {
    // ❌ BAD — swallows original exception, loses stack trace
    throw new ServiceException("DB failed");

    // ✅ GOOD — chains cause, full stack trace preserved
    throw new ServiceException("DB failed", e);
}
```

> **Interview Q: Why must more specific exceptions be caught before general ones?**  
> Because Java evaluates `catch` blocks **top to bottom** and executes the **first matching one**. If you put `catch (Exception e)` first, it catches everything — the more specific blocks below it are **unreachable code** and the compiler will give a "catch block is unreachable" error. Always order: most specific subclass → least specific superclass.

---

## 7. Try-with-Resources

Introduced in **Java 7** — automatically closes any object that implements `AutoCloseable` when the block exits (normally or via exception).

```java
// ─── Before (Java 6) — error-prone verbose cleanup ───
BufferedReader br = null;
try {
    br = new BufferedReader(new FileReader("file.txt"));
    String line = br.readLine();
} catch (IOException e) {
    e.printStackTrace();
} finally {
    if (br != null) {
        try { br.close(); } catch (IOException e) { e.printStackTrace(); }
    }
}

// ─── After (Java 7+) — clean and safe ───
try (BufferedReader br = new BufferedReader(new FileReader("file.txt"))) {
    String line;
    while ((line = br.readLine()) != null) {
        System.out.println(line);
    }
}   // br.close() called automatically — even if exception occurs

// ─── Multiple resources ───
try (Connection conn = dataSource.getConnection();
     PreparedStatement ps = conn.prepareStatement("SELECT * FROM users");
     ResultSet rs = ps.executeQuery()) {

    while (rs.next()) {
        System.out.println(rs.getString("name"));
    }
}   // rs, ps, conn closed in REVERSE ORDER automatically

// ─── Custom AutoCloseable ───
class DatabaseSession implements AutoCloseable {
    DatabaseSession() { System.out.println("Session opened"); }

    void query(String sql) { System.out.println("Executing: " + sql); }

    @Override
    public void close() {
        System.out.println("Session closed");   // called automatically
    }
}

try (DatabaseSession session = new DatabaseSession()) {
    session.query("SELECT 1");
}
// Output:
// Session opened
// Executing: SELECT 1
// Session closed
```

**What if both `try` block and `close()` throw?**

```java
class FailingResource implements AutoCloseable {
    public void doWork() throws Exception { throw new Exception("Work failed"); }
    public void close() throws Exception { throw new Exception("Close failed"); }
}

try (FailingResource r = new FailingResource()) {
    r.doWork();
} catch (Exception e) {
    System.out.println("Primary: " + e.getMessage());           // "Work failed"
    for (Throwable s : e.getSuppressed()) {
        System.out.println("Suppressed: " + s.getMessage());   // "Close failed"
    }
}
```

> **Interview Q: What is try-with-resources? How does it differ from `finally`?**  
> Try-with-resources automatically calls `close()` on any `AutoCloseable` resource declared in its parentheses, in reverse order of declaration. Unlike `finally`, you **don't need to write explicit null checks or nested try-catch for closing**. If both the `try` block and `close()` throw exceptions, the `close()` exception is **suppressed** (attached to the primary exception via `getSuppressed()`), not swallowing the original — this was a bug-prone issue with the old `finally` approach.
