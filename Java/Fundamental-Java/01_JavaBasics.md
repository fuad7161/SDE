# Java Basics

---

## Table of Contents

1. [Java Features](#1-java-features)
2. [JDK vs JRE vs JVM](#2-jdk-vs-jre-vs-jvm)
3. [Java Program Structure](#3-java-program-structure)
4. [Data Types](#4-data-types)
5. [Variables and Scope](#5-variables-and-scope)
6. [Type Casting](#6-type-casting)
7. [Operators](#7-operators)
8. [Control Statements](#8-control-statements)
9. [Arrays](#9-arrays)
10. [Command-Line Arguments](#10-command-line-arguments)

---

## 1. Java Features

| Feature | Explanation |
|---|---|
| **Platform Independent** | Code compiled to bytecode; JVM runs it on any OS |
| **Object-Oriented** | Everything is an object (except primitives) |
| **Robust** | Strong type checking, exception handling, no pointers |
| **Multithreaded** | Built-in Thread class and concurrency utilities |
| **Secure** | No pointers, bytecode verifier, Security Manager |
| **High Performance** | JIT compiler converts hot bytecode to native code |
| **Distributed** | Built-in networking libraries (RMI, sockets) |
| **Dynamic** | Classes loaded at runtime via ClassLoader |

> **Interview Q: Why is Java platform independent?**  
> Java source code is compiled by `javac` into **bytecode** (`.class` files), not machine code. The **JVM** (which is platform-specific) reads and executes this bytecode. So the same `.class` file runs on Windows, Linux, or macOS — as long as a JVM is installed. This is the *Write Once, Run Anywhere* principle.

---

## 2. JDK vs JRE vs JVM

```
┌─────────────────────────────────────────────────────┐
│  JDK  (Java Development Kit)                         │
│  ┌───────────────────────────────────────────────┐  │
│  │  JRE  (Java Runtime Environment)               │  │
│  │  ┌─────────────────────────────────────────┐  │  │
│  │  │  JVM  (Java Virtual Machine)             │  │  │
│  │  │  ┌────────────┐  ┌───────────────────┐  │  │  │
│  │  │  │Class Loader│  │ Execution Engine  │  │  │  │
│  │  │  │            │  │  - Interpreter    │  │  │  │
│  │  │  │            │  │  - JIT Compiler   │  │  │  │
│  │  │  │            │  │  - GC             │  │  │  │
│  │  │  └────────────┘  └───────────────────┘  │  │  │
│  │  └─────────────────────────────────────────┘  │  │
│  │  + Java Standard Libraries (java.util, etc.)  │  │
│  └───────────────────────────────────────────────┘  │
│  + javac, javadoc, jar, jdb, jmap, jstack ...       │
└─────────────────────────────────────────────────────┘
```

| | JVM | JRE | JDK |
|---|---|---|---|
| **Full form** | Java Virtual Machine | Java Runtime Environment | Java Development Kit |
| **Purpose** | Executes bytecode | Provides runtime + libraries | Develop + compile + run Java |
| **Contains** | Classloader, JIT, GC, memory areas | JVM + standard libs | JRE + `javac`, `jar`, debugger, tools |
| **Can compile?** | ❌ | ❌ | ✅ |
| **Can run?** | ✅ | ✅ | ✅ |
| **Who needs it?** | — | End users | Developers |

> **Interview Q: What is the difference between JDK, JRE, and JVM?**  
> **JVM** executes bytecode and is platform-specific (different binaries for Windows/Linux). **JRE** = JVM + standard class libraries needed to run Java programs. **JDK** = JRE + development tools like the compiler (`javac`). To run a Java app you need JRE; to develop one you need JDK.

---

## 3. Java Program Structure

```java
// 1. Package declaration (optional, must be first non-comment line)
package com.company.project;

// 2. Import statements
import java.util.List;
import java.util.ArrayList;

// 3. Class declaration (filename must match public class name)
public class HelloWorld {

    // 4. Static variables — shared across all instances
    static int instanceCount = 0;

    // 5. Instance variables
    String name;

    // 6. Static initializer block — runs once when class is loaded
    static {
        System.out.println("Class loaded");
    }

    // 7. Instance initializer block — runs before every constructor
    {
        instanceCount++;
    }

    // 8. Constructor
    HelloWorld(String name) {
        this.name = name;
    }

    // 9. Methods
    void greet() {
        System.out.println("Hello from " + name);
    }

    // 10. Main method — entry point of the program
    public static void main(String[] args) {
        HelloWorld hw = new HelloWorld("World");
        hw.greet();
    }
}
```

**Execution order:**
```
1. Static blocks/fields (in order, once per class load)
2. main() is called
3. For each `new` object:
   a. Instance blocks
   b. Constructor body
```

> **Interview Q: What is the order of execution of static block, instance block, and constructor?**
> ```java
> class Demo {
>     static { System.out.println("1. Static block"); }
>     { System.out.println("2. Instance block"); }
>     Demo() { System.out.println("3. Constructor"); }
>
>     public static void main(String[] args) {
>         new Demo();   // Output: 1 → 2 → 3
>         new Demo();   // Output: 2 → 3  (static block runs only once)
>     }
> }
> ```

---

## 4. Data Types

### Primitive Types

| Type | Size | Range | Default | Example |
|---|---|---|---|---|
| `byte` | 1 byte | -128 to 127 | `0` | `byte b = 100;` |
| `short` | 2 bytes | -32,768 to 32,767 | `0` | `short s = 1000;` |
| `int` | 4 bytes | -2³¹ to 2³¹-1 (~2 billion) | `0` | `int i = 42;` |
| `long` | 8 bytes | -2⁶³ to 2⁶³-1 | `0L` | `long l = 123L;` |
| `float` | 4 bytes | ~7 decimal digits | `0.0f` | `float f = 3.14f;` |
| `double` | 8 bytes | ~15 decimal digits | `0.0d` | `double d = 3.14;` |
| `char` | 2 bytes | 0 to 65,535 (Unicode) | `'\u0000'` | `char c = 'A';` |
| `boolean` | ~1 bit | `true` / `false` | `false` | `boolean ok = true;` |

### Non-Primitive (Reference) Types

```java
// Non-primitives are objects — stored on heap, variable holds a reference
String name = "Alice";         // String object
int[] nums = {1, 2, 3};       // array object
List<String> list = new ArrayList<>();  // collection object

// Default value is null for all reference types
String s;   // s = null (uninitialized instance variable)
```

> **Interview Q: What is the difference between primitive and non-primitive types?**  
> Primitives hold the **value directly** in the variable (stored on stack or inline in object). Non-primitives (objects) hold a **reference** (memory address) to the actual data on the heap. Primitives are never `null`; non-primitives default to `null`. Primitives are faster; objects provide methods and can be used in generics.

---

## 5. Variables and Scope

```java
public class ScopeDemo {

    static int classVar = 10;     // class variable — shared, lives as long as class
    int instanceVar = 20;         // instance variable — per object, lives with object

    void method() {
        int localVar = 30;        // local variable — lives only in this method call
        System.out.println(classVar + instanceVar + localVar);

        for (int i = 0; i < 3; i++) {  // i scoped to this for loop only
            System.out.println(i);
        }
        // System.out.println(i);  // ❌ compile error — i out of scope
    }
}
```

### Variable Types Summary

| Type | Where declared | Scope | Stored | Default value |
|---|---|---|---|---|
| **Local** | Inside method/block | That block only | Stack | None (must initialize) |
| **Instance** | Inside class, outside method | Entire object | Heap | Type default |
| **Static (class)** | With `static` keyword | Entire class | Method Area | Type default |

> **Interview Q: Why must local variables be initialized before use?**  
> Instance and static variables get default values (0, null, false) automatically. But local variables are stored on the **stack**, and the JVM doesn't initialize stack memory — it's the programmer's responsibility. The compiler enforces this to prevent reading garbage values.

---

## 6. Type Casting

```java
// ── WIDENING (Implicit) ── No data loss, done automatically
byte → short → int → long → float → double

byte b = 10;
int i = b;         // automatic — byte fits in int
double d = i;      // automatic — int fits in double

// ── NARROWING (Explicit) ── May lose data, must cast manually
double pi = 3.99;
int x = (int) pi;   // x = 3 — decimal part truncated (not rounded!)

int big = 300;
byte small = (byte) big;   // small = 44 — overflow! 300 % 256 = 44

// ── OBJECT CASTING ──
Animal a = new Dog();       // upcasting — implicit, always safe
Dog d2 = (Dog) a;           // downcasting — explicit, may throw ClassCastException

// Safe downcasting with instanceof check
if (a instanceof Dog) {
    Dog d3 = (Dog) a;       // safe
}

// Java 16+ pattern matching
if (a instanceof Dog dog) {
    dog.fetch();            // no explicit cast needed
}
```

> **Interview Q: What is the difference between implicit and explicit casting?**  
> **Implicit (widening)** happens automatically when converting a smaller type to a larger type (e.g., `int` to `long`) — no data loss possible. **Explicit (narrowing)** requires a cast operator `(type)` because data loss is possible (e.g., `double` to `int` drops the decimal). For objects, upcasting (child to parent) is implicit; downcasting (parent to child) is explicit and can throw `ClassCastException`.

---

## 7. Operators

```java
// Arithmetic
int a = 10, b = 3;
a + b   // 13
a - b   // 7
a * b   // 30
a / b   // 3  (integer division — no remainder)
a % b   // 1  (remainder)

// Unary
int x = 5;
x++     // post-increment: returns 5, THEN x becomes 6
++x     // pre-increment: x becomes 7, THEN returns 7
x--     // post-decrement

// Comparison
a == b  // false
a != b  // true
a > b   // true
a >= b  // true

// Logical
true && false   // false (AND — both must be true)
true || false   // true  (OR  — at least one true)
!true           // false (NOT)

// Short-circuit evaluation
int n = null;
if (n != null && n > 0) { ... }   // safe — n > 0 not evaluated if n is null

// Bitwise
a & b    // bitwise AND
a | b    // bitwise OR
a ^ b    // bitwise XOR
~a       // bitwise NOT
a << 1   // left shift (multiply by 2)
a >> 1   // right shift (divide by 2)
a >>> 1  // unsigned right shift

// Ternary
String result = (a > b) ? "a is bigger" : "b is bigger";

// Assignment
a += 5;  a -= 3;  a *= 2;  a /= 4;  a %= 3;

// instanceof
String s = "hello";
boolean check = s instanceof String;  // true
```

> **Interview Q: What is the difference between `&` and `&&`?**  
> `&` is the **bitwise AND** operator — it always evaluates both operands. `&&` is the **logical AND** (short-circuit) — if the left side is `false`, it **skips** evaluating the right side. This makes `&&` safer: `if (obj != null && obj.value > 0)` won't NPE because `obj.value > 0` is never reached when `obj` is null.

---

## 8. Control Statements

### `if-else`

```java
int score = 75;

if (score >= 90) {
    System.out.println("A");
} else if (score >= 80) {
    System.out.println("B");
} else if (score >= 70) {
    System.out.println("C");
} else {
    System.out.println("F");
}
// Output: C
```

### `switch`

```java
// Traditional switch
String day = "MONDAY";
switch (day) {
    case "MONDAY":
    case "FRIDAY":
        System.out.println("Busy day");
        break;           // without break — falls through to next case!
    case "SATURDAY":
    case "SUNDAY":
        System.out.println("Weekend");
        break;
    default:
        System.out.println("Regular day");
}

// Switch expression (Java 14+) — cleaner, no fall-through, no break needed
String label = switch (day) {
    case "MONDAY", "FRIDAY" -> "Busy day";
    case "SATURDAY", "SUNDAY" -> "Weekend";
    default -> "Regular day";
};
```

### Loops

```java
// for loop
for (int i = 0; i < 5; i++) {
    System.out.print(i + " ");   // 0 1 2 3 4
}

// while loop — condition checked before entry
int i = 0;
while (i < 5) {
    System.out.print(i++ + " ");
}

// do-while — always executes at least once
int j = 10;
do {
    System.out.println("runs even if false: " + j);
    j++;
} while (j < 5);   // condition false, but prints once

// enhanced for (for-each) — iterates arrays/collections
int[] nums = {1, 2, 3, 4, 5};
for (int n : nums) {
    System.out.print(n + " ");
}

// break and continue
for (int k = 0; k < 10; k++) {
    if (k == 3) continue;   // skip 3
    if (k == 7) break;      // stop at 7
    System.out.print(k + " ");   // 0 1 2 4 5 6
}
```

> **Interview Q: What is the difference between `break` and `continue`?**  
> `break` **exits the loop entirely** — no more iterations. `continue` **skips the current iteration** and moves to the next one. Both can be used with a label to break/continue an **outer loop** from an inner one.
> ```java
> outer:
> for (int i = 0; i < 3; i++) {
>     for (int j = 0; j < 3; j++) {
>         if (j == 1) continue outer;  // skip to next i, not next j
>     }
> }
> ```

---

## 9. Arrays

```java
// Declaration and initialization
int[] nums = new int[5];              // default values: [0, 0, 0, 0, 0]
String[] names = {"Alice", "Bob", "Charlie"};
int[][] matrix = {{1,2},{3,4},{5,6}}; // 2D array

// Access
System.out.println(names[0]);         // "Alice"
System.out.println(names.length);     // 3  (not a method, it's a field)

// Iteration
for (int i = 0; i < nums.length; i++) { ... }
for (String name : names) { ... }     // enhanced for

// Common operations
Arrays.sort(nums);                    // sorts in place
Arrays.fill(nums, 0);                 // fill with 0
int[] copy = Arrays.copyOf(nums, nums.length);
int idx = Arrays.binarySearch(nums, 3);   // only on sorted array
System.out.println(Arrays.toString(nums)); // "[1, 2, 3, 4, 5]"

// Array to List
List<String> list = Arrays.asList(names);  // fixed size — can't add/remove
List<String> mutable = new ArrayList<>(Arrays.asList(names));  // mutable

// ArrayIndexOutOfBoundsException — most common array error
int[] arr = {1, 2, 3};
// arr[3];   // ❌ throws ArrayIndexOutOfBoundsException — valid indices: 0, 1, 2
```

> **Interview Q: What is the difference between `Array` and `ArrayList`?**  
>
> | | Array | ArrayList |
> |---|---|---|
> | Size | Fixed at creation | Dynamic (grows/shrinks) |
> | Type | Primitives + Objects | Objects only |
> | Performance | Faster (direct memory) | Slightly slower (boxing overhead) |
> | Methods | `length` field only | Rich API (add, remove, sort, etc.) |
> | Generics | ❌ | ✅ |

---

## 10. Command-Line Arguments

```java
public class Greet {
    public static void main(String[] args) {
        // args is a String array of space-separated arguments passed at runtime
        if (args.length == 0) {
            System.out.println("Usage: java Greet <name> <age>");
            return;
        }

        String name = args[0];
        int age = Integer.parseInt(args[1]);   // convert String to int

        System.out.println("Hello, " + name + "! You are " + age + " years old.");
    }
}
```

```bash
# Compile and run
javac Greet.java
java Greet Alice 30
# Output: Hello, Alice! You are 30 years old.
```

> **Interview Q: What happens if you access `args[0]` when no arguments are passed?**  
> `args` is never `null` — it's an empty array `[]`. But `args[0]` on an empty array throws `ArrayIndexOutOfBoundsException`. Always check `args.length > 0` before accessing elements. Also, all command-line arguments come in as `String` — you must parse them (`Integer.parseInt`, `Double.parseDouble`, etc.) to use as other types.
