# String Handling

---

## Table of Contents

1. [String Immutability](#1-string-immutability)
2. [String Pool](#2-string-pool)
3. [String vs StringBuilder vs StringBuffer](#3-string-vs-stringbuilder-vs-stringbuffer)
4. [`equals()` vs `==`](#4-equals-vs-)
5. [Common String Methods](#5-common-string-methods)
6. [Mutable vs Immutable Objects](#6-mutable-vs-immutable-objects)

---

## 1. String Immutability

Once a `String` object is created, **its value can never change**. Any operation that appears to "modify" a String actually creates a **new String object**.

```java
String s1 = "Hello";
String s2 = s1.concat(" World");   // creates a NEW String object
String s3 = s1.toUpperCase();      // creates another NEW String object

System.out.println(s1);   // "Hello" — unchanged
System.out.println(s2);   // "Hello World"
System.out.println(s3);   // "HELLO"
```

### Why is String immutable?

```
1. SECURITY
   Strings are used for file paths, network URLs, class names, DB passwords.
   If mutable, another thread could change "admin" to "guest" mid-authentication.

2. THREAD SAFETY
   Immutable objects can be shared freely across threads — no synchronization needed.
   String literals in the pool are shared by all threads safely.

3. HASHCODE CACHING
   String caches its hashCode after first computation.
   HashMap/HashSet rely on this — key's hash must not change after insertion.

4. STRING POOL
   Immutability makes it safe to reuse the same literal across many references.
```

```java
// How Java achieves immutability
public final class String {           // final — cannot be subclassed
    private final byte[] value;       // final byte array — cannot be reassigned
    private int hashCode;             // cached hash

    // No setter methods
    // All "modifying" methods return NEW String instances
}

// Performance trap — never concatenate in a loop with +
String result = "";
for (int i = 0; i < 100; i++) {
    result += i;   // creates 100 new String objects — very inefficient!
}

// Correct — use StringBuilder
StringBuilder sb = new StringBuilder();
for (int i = 0; i < 100; i++) {
    sb.append(i);
}
String result = sb.toString();   // one final conversion
```

> **Interview Q: Why is String immutable in Java?**  
> String is immutable for four reasons: **security** (strings are used as keys for file paths, DB passwords, and class loading — mutability would be a security risk), **thread safety** (immutable objects can be shared without synchronization), **caching** (the `hashCode` is cached after first call — safe because value never changes), and **String pool efficiency** (same literal can be safely shared by multiple references).

---

## 2. String Pool

The **String pool** (intern pool) is a special area in the **Heap** (moved from PermGen to Heap in Java 7) where Java stores unique String literals to avoid duplication.

```java
String a = "hello";           // "hello" created in pool
String b = "hello";           // reuses same pool object — NO new object
String c = new String("hello"); // FORCES a new object on heap, bypasses pool
String d = "hel" + "lo";     // compile-time constant → stored in pool → same as a

System.out.println(a == b);          // true  — same pool reference
System.out.println(a == c);          // false — c is a different heap object
System.out.println(a == d);          // true  — compile-time concatenation goes to pool
System.out.println(a.equals(c));     // true  — same content

// intern() — force a heap string into the pool
String e = c.intern();
System.out.println(a == e);          // true — e is now the pool reference
```

### String Pool Diagram

```
STACK                    HEAP
─────                    ─────────────────────────────────────
a ──────────────────────► ┌──────────────────────────────┐
b ──────────────────────► │   String Pool                │
d ──────────────────────► │   ┌──────────────────┐       │
                          │   │  "hello"          │◄──────│── a, b, d, e all point here
e ──────────────────────► │   └──────────────────┘       │
                          │                              │
c ─────────────────────────────────────────────────────► │
                          │   Regular Heap               │
                          │   ┌──────────────────┐       │
                          │   │ new String("hello")│◄─── │── c points here
                          │   └──────────────────┘       │
                          └──────────────────────────────┘
```

> **Interview Q: What is the String constant pool? When is a String added to it?**  
> The String constant pool is a **cache of String literals** inside the Heap. A String goes into the pool when:  
> 1. You write a **string literal**: `String s = "hello";`  
> 2. You call **`intern()`**: `new String("hello").intern()`  
> 3. The result of a **compile-time constant expression**: `"he" + "llo"`  
> Using `new String("hello")` **bypasses** the pool — always creates a new heap object. The pool saves memory when the same string value is used many times.

---

## 3. String vs StringBuilder vs StringBuffer

| Feature | `String` | `StringBuilder` | `StringBuffer` |
|---|---|---|---|
| Mutability | Immutable | Mutable | Mutable |
| Thread-safe | Yes (immutable) | ❌ No | ✅ Yes (synchronized) |
| Performance | Slow (creates new objects) | **Fastest** | Slower (sync overhead) |
| Storage | String pool or heap | Heap | Heap |
| Java version | Java 1 | Java 5 | Java 1 |

```java
// ── String — immutable, bad for repeated changes ──
String s = "Hello";
s += " World";      // creates a new String "Hello World", s now points to it
s += "!";           // another new String — 3 objects total for 2 appends

// ── StringBuilder — mutable, single-threaded ──
StringBuilder sb = new StringBuilder("Hello");
sb.append(" World");        // modifies the same object
sb.append("!");
sb.insert(5, ",");          // "Hello, World!"
sb.delete(5, 6);            // "Hello World!"
sb.reverse();               // "!dlroW olleH"
sb.replace(0, 1, "?");      // "?dlroW olleH"

System.out.println(sb.toString());
System.out.println(sb.length());    // current length
System.out.println(sb.charAt(0));   // '?'

// ── StringBuffer — thread-safe, use ONLY when shared between threads ──
StringBuffer sbuf = new StringBuffer();
sbuf.append("thread").append("-").append("safe");
// Methods identical to StringBuilder, but synchronized
```

### Performance Comparison

```java
long start, end;

// String concatenation in loop — creates N objects
String result = "";
start = System.nanoTime();
for (int i = 0; i < 10_000; i++) result += i;
end = System.nanoTime();
System.out.println("String: " + (end - start) / 1_000_000 + "ms");     // ~300ms

// StringBuilder — single mutable buffer
StringBuilder sbResult = new StringBuilder();
start = System.nanoTime();
for (int i = 0; i < 10_000; i++) sbResult.append(i);
end = System.nanoTime();
System.out.println("StringBuilder: " + (end - start) / 1_000_000 + "ms"); // ~1ms
```

> **Interview Q: When would you choose `StringBuffer` over `StringBuilder`?**  
> Only when the **same `StringBuffer` object** is being modified by **multiple threads simultaneously**. In practice this is extremely rare — you'd usually prefer a thread-safe design that doesn't share mutable state at all. For any single-threaded code (which is almost everything), always use `StringBuilder` — it has the same API but without the synchronization overhead. **99% of the time, `StringBuilder` is the right choice.**

---

## 4. `equals()` vs `==`

```java
// == compares REFERENCES (memory addresses)
// equals() compares CONTENT (character values)

String a = "hello";
String b = "hello";
String c = new String("hello");

System.out.println(a == b);          // true  — same pool object
System.out.println(a == c);          // false — different objects in memory
System.out.println(a.equals(c));     // true  — same content "hello"
System.out.println(a.equals(b));     // true  — same content

// Common mistake: using == for object comparison
String input = new String("yes");    // simulating user input
if (input == "yes") {                // ❌ WRONG — always false for new String
    System.out.println("Never prints");
}
if ("yes".equals(input)) {           // ✅ CORRECT — use equals()
    System.out.println("Correct!");  // prints
}

// Null-safe equals — put the literal first to avoid NullPointerException
String userInput = null;
// if (userInput.equals("hello")) { }   // ❌ NullPointerException
if ("hello".equals(userInput)) { }     // ✅ safe — returns false
if (Objects.equals(userInput, "hello")) { }  // ✅ also safe — Java 7+
```

> **Interview Q: What is the difference between `==` and `equals()` for Strings?**  
> `==` is the **reference equality operator** — it checks if two variables point to the **same object in memory**. `equals()` checks **content equality** — whether two strings have the same sequence of characters.  
> For String literals, `==` may return `true` due to the String pool, but this is an **implementation detail** you shouldn't rely on. Always use `equals()` to compare String content. Also, put the known non-null value first: `"expected".equals(variable)` to avoid `NullPointerException`.

---

## 5. Common String Methods

```java
String s = "  Hello, Java World!  ";

// ── LENGTH & ACCESS ──
s.length()                  // 22 — number of chars
s.charAt(7)                 // 'J' — char at index 7
s.indexOf("Java")           // 8  — first occurrence (-1 if not found)
s.lastIndexOf("o")          // 17 — last occurrence
s.isEmpty()                 // false — length == 0?
s.isBlank()                 // false — only whitespace? (Java 11+)

// ── COMPARISON ──
s.equals("  Hello, Java World!  ")          // true
s.equalsIgnoreCase("  hello, java world!  ") // true
s.compareTo("abc")          // negative/zero/positive (lexicographic)
s.startsWith("  Hello")     // true
s.endsWith("!  ")           // true
s.contains("Java")          // true
s.matches(".*Java.*")       // true (regex)

// ── TRANSFORMATION ──
s.trim()                    // "Hello, Java World!" — removes leading/trailing whitespace
s.strip()                   // "Hello, Java World!" — Java 11+, handles Unicode spaces
s.toLowerCase()             // "  hello, java world!  "
s.toUpperCase()             // "  HELLO, JAVA WORLD!  "
s.replace("Java", "Python") // "  Hello, Python World!  "
s.replaceAll("\\s+", "_")   // replaces all whitespace runs with _
s.substring(8, 12)          // "Java" — from index 8 (inclusive) to 12 (exclusive)
s.substring(8)              // "Java World!  " — from index 8 to end

// ── SPLITTING & JOINING ──
"a,b,c,d".split(",")        // ["a", "b", "c", "d"]
"a,b,c".split(",", 2)       // ["a", "b,c"] — limit to 2 parts
String.join("-", "a", "b", "c")  // "a-b-c"
String.join(", ", List.of("x", "y", "z"))  // "x, y, z"

// ── CONVERSION ──
String.valueOf(42)          // "42" — converts any type to String
String.valueOf(3.14)        // "3.14"
Integer.parseInt("42")      // 42 — String to int
Double.parseDouble("3.14")  // 3.14

// ── FORMATTING ──
String.format("Name: %s, Age: %d, Score: %.2f", "Alice", 25, 98.567);
// "Name: Alice, Age: 25, Score: 98.57"

// ── JAVA 11+ NEW METHODS ──
"  ".isBlank()              // true
"hello\nworld".lines()      // Stream of ["hello", "world"]
"  hi  ".strip()            // "hi" (Unicode-aware trim)
"ha".repeat(3)              // "hahaha"
```

> **Interview Q: How do you reverse a String in Java?**
> ```java
> // Method 1: StringBuilder
> String reversed = new StringBuilder("hello").reverse().toString();  // "olleh"
>
> // Method 2: char array
> String original = "hello";
> char[] chars = original.toCharArray();
> int left = 0, right = chars.length - 1;
> while (left < right) {
>     char temp = chars[left];
>     chars[left++] = chars[right];
>     chars[right--] = temp;
> }
> String reversed2 = new String(chars);  // "olleh"
>
> // Method 3: Stream (Java 8+)
> String reversed3 = "hello".chars()
>     .mapToObj(c -> String.valueOf((char) c))
>     .collect(Collectors.collectingAndThen(Collectors.toList(),
>         list -> { Collections.reverse(list); return String.join("", list); }));
> ```

---

## 6. Mutable vs Immutable Objects

```java
// ── IMMUTABLE OBJECT — state cannot change after construction ──
// Examples: String, Integer, Double, LocalDate, BigDecimal

String s = "hello";
// All methods return NEW objects — original never changes
String upper = s.toUpperCase();  // new object
System.out.println(s);           // still "hello"

Integer x = 5;
Integer y = x + 3;   // x is still 5, y is a new Integer(8)

// ── MUTABLE OBJECT — state can change ──
// Examples: StringBuilder, ArrayList, HashMap, most custom classes

StringBuilder sb = new StringBuilder("hello");
sb.append(" world");   // SAME object modified
System.out.println(sb); // "hello world"

List<String> list = new ArrayList<>();
list.add("a");   // same object modified

// ── Creating an immutable class ──
public final class Money {              // 1. Make class final
    private final double amount;        // 2. All fields private final
    private final String currency;      // 3. No setters

    public Money(double amount, String currency) {
        this.amount = amount;
        this.currency = currency;
    }

    public double getAmount() { return amount; }
    public String getCurrency() { return currency; }

    // "Operations" return NEW Money objects
    public Money add(Money other) {
        if (!this.currency.equals(other.currency))
            throw new IllegalArgumentException("Currency mismatch");
        return new Money(this.amount + other.amount, this.currency);
    }

    @Override public String toString() {
        return amount + " " + currency;
    }
}

Money price = new Money(10.00, "USD");
Money tax   = new Money(1.50, "USD");
Money total = price.add(tax);    // price and tax unchanged, total is new
System.out.println(price);       // 10.0 USD — unchanged
System.out.println(total);       // 11.5 USD
```

| | Immutable | Mutable |
|---|---|---|
| State after creation | Cannot change | Can change |
| Thread safety | Inherently safe | Needs synchronization |
| `hashCode` | Can cache | Must recompute (risky in maps) |
| Defensive copy | Not needed | Required to protect state |
| Examples | `String`, `Integer`, `LocalDate` | `StringBuilder`, `ArrayList` |

> **Interview Q: How do you create an immutable class in Java?**  
> 1. Declare the class as **`final`** (prevents subclassing and overriding)  
> 2. Make all fields **`private final`** (no reassignment)  
> 3. Provide **no setters**  
> 4. Initialize all fields **in the constructor**  
> 5. If any field holds a **mutable object** (like a `List`), store a **defensive copy** in the constructor and return a copy (or unmodifiable view) from the getter — otherwise a caller could get the reference and mutate it externally.
