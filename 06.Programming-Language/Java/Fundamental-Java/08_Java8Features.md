# Java 8 Features

---

## Table of Contents

1. [Lambda Expressions](#1-lambda-expressions)
2. [Functional Interfaces](#2-functional-interfaces)
3. [Stream API](#3-stream-api)
4. [Method References](#4-method-references)
5. [Optional Class](#5-optional-class)
6. [Default & Static Interface Methods](#6-default--static-interface-methods)
7. [Date & Time API (`java.time`)](#7-date--time-api-javatime)

---

## 1. Lambda Expressions

A lambda is a **concise way to represent an anonymous function** — a block of code that can be passed around as data.

```java
// Syntax:  (parameters) -> expression
//          (parameters) -> { statements; }

// ── Before Java 8 — anonymous inner class ──
Runnable r1 = new Runnable() {
    @Override
    public void run() {
        System.out.println("Running...");
    }
};

// ── Java 8 lambda ──
Runnable r2 = () -> System.out.println("Running...");

// ── Single parameter (parentheses optional) ──
// Consumer<String> printer = (s) -> System.out.println(s);
Consumer<String> printer = s -> System.out.println(s);

// ── Multiple parameters ──
Comparator<String> byLength = (a, b) -> Integer.compare(a.length(), b.length());

// ── Multi-statement body ──
Comparator<String> complex = (a, b) -> {
    int lenCompare = Integer.compare(a.length(), b.length());
    return lenCompare != 0 ? lenCompare : a.compareTo(b);
};

// ── With return (single expression, no return keyword) ──
Function<Integer, Integer> square = n -> n * n;
System.out.println(square.apply(5));   // 25

// ── Capturing variables — must be effectively final ──
String prefix = "Hello";
Function<String, String> greet = name -> prefix + ", " + name;
// prefix = "Hi";   // ❌ compile error — prefix is captured, must be effectively final
```

> **Interview Q: What is a lambda expression? What problem does it solve?**  
> A lambda is an **anonymous function** — it has no name, but has parameters, a body, and a return type. It solves the **verbosity of anonymous inner classes** for single-method interfaces. Lambdas make code more concise, enable passing behavior as data, and are the foundation for functional-style programming in Java (streams, optional, method references). They can capture **effectively final** local variables from the enclosing scope.

---

## 2. Functional Interfaces

A **functional interface** has exactly **one abstract method** (SAM — Single Abstract Method). Lambdas can only be used where a functional interface is expected.

```java
// ── Built-in functional interfaces in java.util.function ──

// Supplier<T> — takes nothing, returns T
Supplier<String> greeting = () -> "Hello, World!";
System.out.println(greeting.get());   // "Hello, World!"

// Consumer<T> — takes T, returns nothing (side effects)
Consumer<String> print = s -> System.out.println(">> " + s);
print.accept("test");   // >> test

// BiConsumer<T, U> — takes T and U, returns nothing
BiConsumer<String, Integer> log = (msg, count) ->
    System.out.println(msg + " (x" + count + ")");
log.accept("error", 3);

// Function<T, R> — takes T, returns R
Function<String, Integer> length = String::length;   // method reference
Function<Integer, Integer> doubleIt = n -> n * 2;
Function<String, Integer> combined = length.andThen(doubleIt);  // compose
System.out.println(combined.apply("hello"));   // 10

// BiFunction<T, U, R>
BiFunction<String, Integer, String> repeat = (s, n) -> s.repeat(n);

// Predicate<T> — takes T, returns boolean
Predicate<String> isLong = s -> s.length() > 5;
Predicate<String> startsWithH = s -> s.startsWith("H");
Predicate<String> combined2 = isLong.and(startsWithH);   // compose
System.out.println(combined2.test("Hello World"));   // true

// UnaryOperator<T> — Function<T, T> (same type in and out)
UnaryOperator<String> shout = s -> s.toUpperCase() + "!";

// BinaryOperator<T> — BiFunction<T, T, T>
BinaryOperator<Integer> add = (a, b) -> a + b;
System.out.println(add.apply(3, 4));   // 7

// ── Custom functional interface ──
@FunctionalInterface
interface Transformer<T, R> {
    R transform(T input);
    // Can have default and static methods — still functional
    default Transformer<T, R> andLog() {
        return input -> {
            R result = this.transform(input);
            System.out.println(input + " → " + result);
            return result;
        };
    }
}

Transformer<String, Integer> parser = Integer::parseInt;
System.out.println(parser.transform("42"));   // 42
```

| Interface | Method | In | Out |
|---|---|---|---|
| `Supplier<T>` | `get()` | — | T |
| `Consumer<T>` | `accept(T)` | T | void |
| `Function<T,R>` | `apply(T)` | T | R |
| `Predicate<T>` | `test(T)` | T | boolean |
| `UnaryOperator<T>` | `apply(T)` | T | T |
| `BinaryOperator<T>` | `apply(T,T)` | T, T | T |

> **Interview Q: What is a functional interface? Can it have more than one method?**  
> A functional interface has **exactly one abstract method** (SAM). It can have any number of `default` and `static` methods — those don't count. The `@FunctionalInterface` annotation is optional but recommended (it triggers a compile error if you accidentally add a second abstract method). All lambdas and method references are typed as functional interfaces.

---

## 3. Stream API

A **Stream** is a sequence of elements supporting aggregate operations. Streams don't store data — they process it **lazily** through a pipeline.

```java
List<String> names = List.of("Alice", "Bob", "Charlie", "David", "Anna");

// ── BASIC PIPELINE: source → intermediate → terminal ──
List<String> result = names.stream()        // source: creates stream
    .filter(s -> s.startsWith("A"))         // intermediate: lazy
    .map(String::toUpperCase)               // intermediate: lazy
    .sorted()                               // intermediate: lazy
    .collect(Collectors.toList());          // terminal: triggers execution
System.out.println(result);   // [ALICE, ANNA]

// ── INTERMEDIATE OPERATIONS (lazy, return Stream) ──
names.stream()
    .filter(s -> s.length() > 3)           // keep if predicate true
    .map(String::toLowerCase)              // transform each element
    .mapToInt(String::length)              // to IntStream
    .flatMap(...)                          // flatten nested streams
    .distinct()                            // remove duplicates
    .sorted()                              // natural order
    .sorted(Comparator.reverseOrder())     // custom order
    .limit(3)                              // max 3 elements
    .skip(1)                               // skip first 1
    .peek(s -> log.debug(s));              // debug without consuming

// ── TERMINAL OPERATIONS (eager, trigger execution) ──
long count  = names.stream().filter(s -> s.length() > 3).count();
boolean any = names.stream().anyMatch(s -> s.startsWith("A"));   // true
boolean all = names.stream().allMatch(s -> s.length() > 2);      // true
boolean none= names.stream().noneMatch(s -> s.isEmpty());        // true

Optional<String> first = names.stream().filter(s -> s.startsWith("C")).findFirst();
Optional<String> any2  = names.stream().filter(s -> s.length() > 4).findAny();

// Reduce — combine elements
int totalLength = names.stream().mapToInt(String::length).sum();
Optional<Integer> product = Stream.of(1, 2, 3, 4).reduce((a, b) -> a * b);   // 24
int sum = Stream.of(1, 2, 3).reduce(0, Integer::sum);   // 6 (with identity)

// ── COLLECTORS ──
List<String>    list = names.stream().collect(Collectors.toList());
Set<String>     set  = names.stream().collect(Collectors.toSet());
String joined        = names.stream().collect(Collectors.joining(", ", "[", "]"));
// "[Alice, Bob, Charlie, David, Anna]"

Map<Integer, List<String>> byLength = names.stream()
    .collect(Collectors.groupingBy(String::length));
// {3=[Bob], 5=[Alice, David], 4=[Anna], 7=[Charlie]}

Map<Boolean, List<String>> partition = names.stream()
    .collect(Collectors.partitioningBy(s -> s.length() > 4));
// {false=[Bob, Anna], true=[Alice, Charlie, David]}

// ── PARALLEL STREAM ──
long count2 = names.parallelStream()
    .filter(s -> s.length() > 3)
    .count();
// Uses ForkJoinPool — splits work across CPU cores
// Good for CPU-intensive operations on large datasets
// Avoid for I/O-bound, stateful, or order-dependent operations
```

> **Interview Q: What is the difference between `map()` and `flatMap()` in Stream?**  
> `map()` applies a function to each element and returns a stream of the results — one-to-one transformation. `flatMap()` applies a function that returns a Stream for each element, then **flattens** all those streams into one — one-to-many transformation.  
> ```java
> // map: List<String> → Stream<Stream<char[]>>
> Stream<String[]> mapped = Stream.of("Hello", "World").map(s -> s.split(""));
>
> // flatMap: List<String> → Stream<char>  (flattened)
> Stream<String> flattened = Stream.of("Hello", "World")
>     .flatMap(s -> Arrays.stream(s.split("")));
>
> // Common use: flatten List<List<T>>
> List<List<Integer>> nested = List.of(List.of(1,2), List.of(3,4));
> List<Integer> flat = nested.stream().flatMap(Collection::stream)
>     .collect(Collectors.toList());   // [1, 2, 3, 4]
> ```

---

## 4. Method References

A shorthand syntax for lambdas that just call an existing method.

```java
// Syntax: ClassName::methodName  or  instance::methodName

List<String> names = List.of("Charlie", "Alice", "Bob");

// ── Static method reference: ClassName::staticMethod ──
// Lambda:  n -> Integer.parseInt(n)
Function<String, Integer> parser = Integer::parseInt;

// ── Instance method reference (on a specific instance): instance::method ──
// Lambda:  s -> System.out.println(s)
Consumer<String> print = System.out::println;
names.forEach(System.out::println);

// ── Instance method reference (on arbitrary instance of type): Type::instanceMethod ──
// Lambda:  s -> s.toUpperCase()
Function<String, String> upper = String::toUpperCase;
// Lambda:  (s1, s2) -> s1.compareTo(s2)
Comparator<String> comparator = String::compareTo;

names.stream().map(String::toUpperCase).forEach(System.out::println);
names.sort(String::compareTo);

// ── Constructor reference: ClassName::new ──
// Lambda:  name -> new Person(name)
Function<String, Person> factory = Person::new;
Person p = factory.apply("Alice");

// Supplier<List<String>> listFactory = ArrayList::new;
Supplier<List<String>> listFactory = ArrayList::new;
List<String> newList = listFactory.get();

// ── When to use ──
// Use method reference when lambda just calls ONE existing method with NO extra logic
// names.stream().map(s -> s.toUpperCase())   →  .map(String::toUpperCase)  ✅
// names.stream().map(s -> s.toUpperCase() + "!")  — must stay as lambda ❌
```

| Type | Syntax | Equivalent Lambda |
|---|---|---|
| Static method | `Integer::parseInt` | `s -> Integer.parseInt(s)` |
| Bound instance | `System.out::println` | `s -> System.out.println(s)` |
| Unbound instance | `String::toUpperCase` | `s -> s.toUpperCase()` |
| Constructor | `Person::new` | `name -> new Person(name)` |

> **Interview Q: What are method references? What are the four types?**  
> Method references are a compact syntax for lambdas that delegate to an existing method. The four types: (1) **Static** — `ClassName::staticMethod`; (2) **Bound instance** — `object::instanceMethod` (specific object); (3) **Unbound instance** — `ClassName::instanceMethod` (any instance of that type is the first parameter); (4) **Constructor** — `ClassName::new`. They improve readability by removing boilerplate lambda wrappers around single method calls.

---

## 5. Optional Class

`Optional<T>` is a container that may or may not hold a value — designed to **eliminate `NullPointerException`** by making nullability explicit.

```java
// ── Creating ──
Optional<String> present = Optional.of("hello");      // must not be null
Optional<String> maybe   = Optional.ofNullable(null); // null is OK
Optional<String> empty   = Optional.empty();

// ── Checking ──
present.isPresent();    // true
empty.isPresent();      // false
present.isEmpty();      // false (Java 11+)

// ── Getting value ──
present.get();                         // "hello" — throws if empty
maybe.orElse("default");               // "default" if empty
maybe.orElseGet(() -> compute());      // lazy: compute only if empty
maybe.orElseThrow();                   // throw NoSuchElementException if empty
maybe.orElseThrow(() -> new EntityNotFoundException("Not found"));

// ── Transforming (only runs if value present) ──
Optional<Integer> length = present.map(String::length);   // Optional[5]
Optional<String>  upper  = present.map(String::toUpperCase);  // Optional["HELLO"]

// flatMap — use when map returns Optional (prevents Optional<Optional<T>>)
Optional<String> result = Optional.of("user123")
    .flatMap(id -> userService.findById(id));   // findById returns Optional<User>

// filter — empty if predicate fails
Optional<String> longName = present.filter(s -> s.length() > 3); // present
Optional<String> noMatch  = present.filter(s -> s.length() > 10);// empty

// ── ifPresent — side effect ──
present.ifPresent(s -> System.out.println("Found: " + s));
present.ifPresentOrElse(
    s -> System.out.println("Found: " + s),
    () -> System.out.println("Not found")   // Java 9+
);

// ── Chaining to avoid NPE ──
// Before Optional (NPE-prone):
String city = null;
if (user != null && user.getAddress() != null) {
    city = user.getAddress().getCity();
}

// With Optional (clean):
String city2 = Optional.ofNullable(user)
    .map(User::getAddress)
    .map(Address::getCity)
    .orElse("Unknown");
```

> **Interview Q: What is `Optional`? When should you use it?**  
> `Optional<T>` is a wrapper that explicitly signals that a value **may or may not be present**, forcing callers to handle the absence case instead of risking `NullPointerException`. Use it as a **return type** for methods that might not return a value (e.g., `findById`). **Don't use it**: as a method parameter (ugly API), as instance fields (not serializable), or for collections (return empty collection instead). Never call `.get()` without first checking `.isPresent()` — use `orElse`, `orElseGet`, or `orElseThrow` instead.

---

## 6. Default & Static Interface Methods

```java
// ── DEFAULT METHODS — concrete methods in interfaces (Java 8+) ──
// Motivation: add new methods to interfaces without breaking existing implementations

interface Logger {
    void log(String message);          // abstract — implementors must provide

    default void info(String msg) {
        log("[INFO] " + msg);          // concrete — with default implementation
    }

    default void error(String msg) {
        log("[ERROR] " + msg);
    }

    static Logger console() {          // factory method
        return msg -> System.out.println(msg);
    }
}

class FileLogger implements Logger {
    private final String filename;
    FileLogger(String f) { this.filename = f; }

    @Override
    public void log(String message) {
        // write to file
        System.out.println("File[" + filename + "]: " + message);
    }
    // info() and error() inherited for free!
}

FileLogger fl = new FileLogger("app.log");
fl.info("Server started");    // File[app.log]: [INFO] Server started
fl.error("DB timeout");       // File[app.log]: [ERROR] DB timeout

Logger logger = Logger.console();   // static factory
logger.info("Hello");   // prints to console

// ── DIAMOND PROBLEM RESOLUTION with default methods ──
interface A { default void greet() { System.out.println("A"); } }
interface B { default void greet() { System.out.println("B"); } }

class C implements A, B {
    @Override
    public void greet() {
        A.super.greet();   // explicitly choose — must override to resolve conflict
    }
}
```

> **Interview Q: Why were default methods added to interfaces in Java 8?**  
> To allow **backward-compatible evolution** of interfaces. Before Java 8, adding a method to an interface would **break all existing implementations** — every class implementing that interface would fail to compile. `default` methods let the JDK add new methods (like `forEach`, `stream`, `sort` to `Collection`, `List`, `Map`) without breaking the thousands of existing implementations. They're also useful for providing utility methods related to the interface contract.

---

## 7. Date & Time API (`java.time`)

The old `java.util.Date` and `Calendar` were mutable, not thread-safe, and had confusing APIs. Java 8 introduced `java.time` — immutable, thread-safe, and intuitive.

```java
import java.time.*;
import java.time.format.DateTimeFormatter;
import java.time.temporal.ChronoUnit;

// ── CORE CLASSES ──
LocalDate date       = LocalDate.now();              // 2024-11-15
LocalTime time       = LocalTime.now();              // 14:30:45.123
LocalDateTime dtm    = LocalDateTime.now();          // 2024-11-15T14:30:45.123
ZonedDateTime zdt    = ZonedDateTime.now();          // with timezone
Instant instant      = Instant.now();                // machine time (epoch)
Duration duration    = Duration.ofHours(2);          // time-based amount
Period period        = Period.ofDays(30);            // date-based amount

// ── CREATING DATES ──
LocalDate birthday   = LocalDate.of(1995, Month.MARCH, 15);
LocalDate fromString = LocalDate.parse("2024-06-01");
LocalTime alarm      = LocalTime.of(7, 30, 0);
LocalDateTime meeting = LocalDateTime.of(2024, 6, 1, 10, 0);

// ── DATE OPERATIONS — all return NEW instances (immutable) ──
LocalDate nextWeek  = date.plusWeeks(1);
LocalDate lastMonth = date.minusMonths(1);
LocalDate firstDay  = date.withDayOfMonth(1);  // first day of current month

date.getDayOfWeek();   // FRIDAY
date.getMonth();       // NOVEMBER
date.getYear();        // 2024
date.isLeapYear();     // false
date.isAfter(birthday);   // true
date.isBefore(LocalDate.of(2025, 1, 1));  // true

// ── PERIOD & DURATION ──
Period age = Period.between(birthday, date);
System.out.println(age.getYears() + " years old");   // 29 years old

Duration between = Duration.between(LocalTime.of(9, 0), LocalTime.of(17, 30));
System.out.println(between.toHours() + " hours");    // 8 hours

long daysBetween = ChronoUnit.DAYS.between(birthday, date);

// ── FORMATTING & PARSING ──
DateTimeFormatter fmt = DateTimeFormatter.ofPattern("dd/MM/yyyy HH:mm");
String formatted = dtm.format(fmt);        // "15/11/2024 14:30"
LocalDateTime parsed = LocalDateTime.parse("15/11/2024 14:30", fmt);

// ISO format
DateTimeFormatter iso = DateTimeFormatter.ISO_LOCAL_DATE;
System.out.println(date.format(iso));      // "2024-11-15"

// ── TIMEZONE ──
ZoneId zone = ZoneId.of("America/New_York");
ZonedDateTime nyTime = ZonedDateTime.now(zone);
ZonedDateTime istTime = nyTime.withZoneSameInstant(ZoneId.of("Asia/Kolkata"));
```

| Class | Contains | Timezone |
|---|---|---|
| `LocalDate` | Date only (year/month/day) | No |
| `LocalTime` | Time only (hour/min/sec/ns) | No |
| `LocalDateTime` | Date + Time | No |
| `ZonedDateTime` | Date + Time + Zone | Yes |
| `Instant` | Machine timestamp (epoch) | UTC |
| `Duration` | Time-based amount (hours, mins) | — |
| `Period` | Date-based amount (years, months, days) | — |

> **Interview Q: What are the improvements in the Java 8 Date/Time API over `java.util.Date`?**  
> The old API problems: `Date` is **mutable** (thread-unsafe), months are 0-based (January = 0), `Calendar` is verbose and confusing, no clear separation between date-only and time-only.  
> Java 8 improvements: (1) **Immutable** — all operations return new instances; (2) **Thread-safe**; (3) **Clear separation** — `LocalDate`, `LocalTime`, `LocalDateTime`, `ZonedDateTime`; (4) **Fluent API** — easy chaining; (5) **Proper month names** — `Month.JANUARY` instead of 0; (6) **Built-in parsing/formatting** with `DateTimeFormatter`.
