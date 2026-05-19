# Important Concepts

---

## Table of Contents

1. [Immutable Class](#1-immutable-class)
2. [Singleton Pattern](#2-singleton-pattern)
3. [Marker Interface](#3-marker-interface)
4. [Wrapper Classes](#4-wrapper-classes)
5. [Autoboxing & Unboxing](#5-autoboxing--unboxing)
6. [Pass by Value in Java](#6-pass-by-value-in-java)
7. [Shallow Copy vs Deep Copy](#7-shallow-copy-vs-deep-copy)
8. [Reflection API](#8-reflection-api)

---

## 1. Immutable Class

An immutable class's **state cannot change after construction**. `String`, `Integer`, `LocalDate` are all immutable.

```java
// ── Rules to create an immutable class ──

public final class Money {                   // 1. Declare class final (prevent subclassing)
    private final double amount;             // 2. All fields private final
    private final String currency;           // 3. No setters
    private final List<String> tags;         // 4. Mutable fields need defensive copy

    public Money(double amount, String currency, List<String> tags) {
        if (amount < 0) throw new IllegalArgumentException("Amount cannot be negative");
        this.amount = amount;
        this.currency = currency;
        // 5. Defensive copy of mutable parameter
        this.tags = Collections.unmodifiableList(new ArrayList<>(tags));
    }

    public double getAmount() { return amount; }
    public String getCurrency() { return currency; }

    // 5. Return defensive copy (or unmodifiable view) for mutable fields
    public List<String> getTags() {
        return tags;   // already unmodifiable — safe to return directly
    }

    // "Modification" returns a NEW object
    public Money add(Money other) {
        if (!this.currency.equals(other.currency))
            throw new IllegalArgumentException("Currency mismatch");
        List<String> combined = new ArrayList<>(this.tags);
        combined.addAll(other.tags);
        return new Money(this.amount + other.amount, this.currency, combined);
    }

    @Override
    public boolean equals(Object o) {
        if (!(o instanceof Money)) return false;
        Money m = (Money) o;
        return Double.compare(amount, m.amount) == 0 && currency.equals(m.currency);
    }

    @Override
    public int hashCode() { return Objects.hash(amount, currency); }

    @Override
    public String toString() { return amount + " " + currency; }
}

Money price = new Money(10.00, "USD", List.of("sale"));
Money tax   = new Money(1.50, "USD", List.of());
Money total = price.add(tax);   // price and tax unchanged
System.out.println(price);      // 10.0 USD — unchanged
System.out.println(total);      // 11.5 USD
```

**Why create immutable classes?**
- Thread-safe by design — no synchronization needed
- Safe to use as HashMap/HashSet keys (hashCode never changes)
- Easier to reason about — no unexpected state changes
- Defensive programming — can share references freely

> **Interview Q: What are the steps to create an immutable class?**  
> 1. Make the class `final` (prevent subclassing and overriding)  
> 2. Make all fields `private final`  
> 3. No setter methods  
> 4. Initialize all fields via constructor  
> 5. For any **mutable field** (List, array, custom object), make a **defensive copy** in the constructor and return a defensive copy (or unmodifiable view) from getters

---

## 2. Singleton Pattern

Ensures **only one instance** of a class exists across the entire application.

```java
// ── Approach 1: Eager Initialization ──
// Instance created at class load time (thread-safe, simple)
class EagerSingleton {
    private static final EagerSingleton INSTANCE = new EagerSingleton();
    private EagerSingleton() {}   // private constructor
    public static EagerSingleton getInstance() { return INSTANCE; }
}

// ── Approach 2: Lazy (Double-Checked Locking) ── (pre-Java 5 had issues)
class LazySingleton {
    private static volatile LazySingleton instance;   // volatile is essential!
    private LazySingleton() {}

    public static LazySingleton getInstance() {
        if (instance == null) {                // first check (no lock)
            synchronized (LazySingleton.class) {
                if (instance == null) {        // second check (with lock)
                    instance = new LazySingleton();
                }
            }
        }
        return instance;
    }
}

// ── Approach 3: Initialization-on-Demand Holder (BEST) ──
// Thread-safe by class loading guarantees — no synchronization needed
class HolderSingleton {
    private HolderSingleton() {}

    private static class Holder {
        static final HolderSingleton INSTANCE = new HolderSingleton();
    }

    public static HolderSingleton getInstance() {
        return Holder.INSTANCE;   // inner class loaded lazily on first access
    }
}

// ── Approach 4: Enum (Josh Bloch's recommendation) ──
// Serialization-safe, reflection-attack-proof, thread-safe
enum EnumSingleton {
    INSTANCE;

    public void doSomething() {
        System.out.println("Singleton behavior");
    }
}
EnumSingleton.INSTANCE.doSomething();

// ── Why volatile is needed in double-checked locking ──
// instance = new Singleton() is NOT atomic:
// 1. Allocate memory
// 2. Initialize object
// 3. Assign reference to 'instance'
// Steps 2 and 3 can be REORDERED without volatile
// Thread B might see non-null 'instance' but uninitialized object
```

> **Interview Q: What is the Singleton pattern? What is the best implementation?**  
> Singleton ensures a class has only one instance and provides global access to it. The **best implementation** for most cases is the **Initialization-on-Demand Holder** (static inner class) — lazy (created on first use), thread-safe (class loading is thread-safe by JVM spec), no synchronization overhead. For simplicity or when serialization safety is required, use an **enum**. Double-checked locking requires `volatile` and is prone to subtle bugs.

---

## 3. Marker Interface

An interface with **no methods** — it simply "marks" a class to signal some property to the JVM or frameworks.

```java
// java.io.Serializable — marks class as serializable
class User implements Serializable {
    String name;
    int age;
    // No methods to implement — JVM checks 'instanceof Serializable' before serializing
}

// java.lang.Cloneable — marks class as supporting clone()
class Point implements Cloneable {
    int x, y;
    @Override
    public Object clone() throws CloneNotSupportedException {
        return super.clone();   // if Cloneable not implemented → CloneNotSupportedException
    }
}

// java.io.Externalizable — technically has methods, but similar concept
// java.util.RandomAccess  — marks List as supporting O(1) random access (e.g., ArrayList)

// ── Custom Marker Interface ──
interface AuditLog { }   // mark classes that should be audit-logged

class PaymentTransaction implements AuditLog {
    double amount;
}

// Framework checks marker:
public void process(Object obj) {
    if (obj instanceof AuditLog) {
        auditLogger.log(obj);   // only log marked classes
    }
}
```

**Modern alternative: Annotations**

```java
// Annotations are generally preferred over marker interfaces now
@Retention(RetentionPolicy.RUNTIME)
@Target(ElementType.TYPE)
@interface AuditLog { }

@AuditLog
class PaymentTransaction { }

// Check at runtime:
if (obj.getClass().isAnnotationPresent(AuditLog.class)) {
    auditLogger.log(obj);
}
```

> **Interview Q: What is a marker interface? Can you give examples?**  
> A marker interface has **no methods** — it's an empty interface used to **tag** classes with some property. Examples: `Serializable` (allows serialization), `Cloneable` (allows `clone()`), `RandomAccess` (signals O(1) access). The JVM or framework checks `instanceof MarkerInterface` to decide behavior. Modern Java prefers **annotations** over marker interfaces because annotations can carry metadata and be inspected at compile time. However, marker interfaces have one advantage: the compiler can enforce them at the variable reference level (`Serializable s = obj` vs annotation check at runtime only).

---

## 4. Wrapper Classes

Each primitive type has a corresponding **wrapper class** in `java.lang`.

| Primitive | Wrapper | Default value | Size |
|---|---|---|---|
| `byte` | `Byte` | 0 | 8-bit |
| `short` | `Short` | 0 | 16-bit |
| `int` | `Integer` | 0 | 32-bit |
| `long` | `Long` | 0L | 64-bit |
| `float` | `Float` | 0.0f | 32-bit |
| `double` | `Double` | 0.0d | 64-bit |
| `char` | `Character` | '\u0000' | 16-bit |
| `boolean` | `Boolean` | false | 1-bit |

```java
// ── Why wrapper classes? ──
// 1. Use primitives in Collections (generics require objects)
List<Integer> numbers = new ArrayList<>();   // can't use List<int>
numbers.add(42);   // autoboxing: int → Integer

// 2. Null representation (primitives can't be null)
Integer score = null;   // valid — "unknown score"
int primitive = null;   // ❌ compile error

// 3. Utility methods
Integer.parseInt("42");                  // String → int
Integer.toBinaryString(42);             // "101010"
Integer.toHexString(255);               // "ff"
Integer.max(10, 20);                    // 20
Integer.compare(5, 10);                 // negative (5 < 10)

Double.parseDouble("3.14");
Boolean.parseBoolean("true");           // true
Character.isDigit('5');                 // true
Character.isLetter('a');                // true
Character.toUpperCase('a');             // 'A'

// 4. Integer cache (-128 to 127)
Integer a = 100;
Integer b = 100;
System.out.println(a == b);    // true  — cached, same object

Integer c = 200;
Integer d = 200;
System.out.println(c == d);    // false — outside cache, different objects
System.out.println(c.equals(d)); // true — same value
```

> **Interview Q: What is Integer caching? What is the range?**  
> Java caches `Integer` objects for values **-128 to 127** (configurable with `-XX:AutoBoxCacheMax`). When you autobox or use `Integer.valueOf(n)` for values in this range, you get the **same cached object** back — so `==` returns true. Outside this range, a new `Integer` object is created each time, so `==` returns false. This is why you should always use `equals()` to compare `Integer` values, not `==`.

---

## 5. Autoboxing & Unboxing

**Autoboxing** — automatic conversion from primitive to wrapper (Java 5+).  
**Unboxing** — automatic conversion from wrapper to primitive.

```java
// ── AUTOBOXING ──
int i = 42;
Integer wrapped = i;        // autoboxing: int → Integer
// Equivalent to: Integer wrapped = Integer.valueOf(i);

List<Integer> list = new ArrayList<>();
list.add(10);               // autoboxed: int 10 → Integer(10)
list.add(20);

// ── UNBOXING ──
Integer obj = 100;
int primitive = obj;        // unboxing: Integer → int
// Equivalent to: int primitive = obj.intValue();

// Math operations trigger unboxing
Integer x = 5, y = 3;
int sum = x + y;            // both unboxed for arithmetic

// ── PERFORMANCE WARNING ──
// Repeated boxing in loops creates many objects
Long total = 0L;
for (long i = 0; i < 1_000_000; i++) {
    total += i;   // ❌ unbox total, add i, rebox result — 1 million objects!
}

long total2 = 0L;
for (long i = 0; i < 1_000_000; i++) {
    total2 += i;  // ✅ no boxing — use primitive
}

// ── NULL TRAP — NullPointerException on unboxing ──
Integer score = null;
int points = score;   // ❌ NullPointerException — unboxing null!

// Fix: null check before unboxing
int points = (score != null) ? score : 0;
```

> **Interview Q: What is autoboxing? What are its pitfalls?**  
> Autoboxing is Java's automatic conversion between primitive types and their wrapper counterparts. Pitfalls: (1) **Performance** — creates many wrapper objects in tight loops; use primitives where possible; (2) **NullPointerException** — unboxing a `null` wrapper throws NPE; (3) **Unexpected `==` comparisons** — two `Integer` values outside the -128 to 127 cache will be `!=` even if equal; (4) **Overloading ambiguity** — autoboxing can affect method resolution unexpectedly.

---

## 6. Pass by Value in Java

Java is **always pass-by-value** — but the "value" of an object reference is the memory address.

```java
// ── PRIMITIVE — pass-by-value (copy of value) ──
void doubleIt(int x) {
    x = x * 2;   // modifies local copy only
    System.out.println("Inside: " + x);  // 20
}

int num = 10;
doubleIt(num);
System.out.println("Outside: " + num);   // 10 — unchanged

// ── OBJECT — pass reference by value (copy of reference) ──
void rename(Person p) {
    p.name = "Bob";   // modifies the OBJECT the reference points to
    // Both 'p' (local) and 'alice' (caller) point to same object
}

void reassign(Person p) {
    p = new Person("Charlie");  // reassigns LOCAL copy of reference
    // 'alice' in caller still points to original object
}

Person alice = new Person("Alice");
rename(alice);
System.out.println(alice.name);    // "Bob" — object was modified

reassign(alice);
System.out.println(alice.name);    // "Bob" — NOT "Charlie", reassignment only local

// ── String "gotcha" — strings are immutable ──
void change(String s) {
    s = s.toUpperCase();   // creates NEW String, s is local reference
}

String name = "alice";
change(name);
System.out.println(name);   // "alice" — unchanged (immutable + by-value)
```

**Memory diagram:**

```
STACK (caller)         STACK (called)         HEAP
──────────────         ──────────────         ──────────────────
alice → ─────────────────────────────────────► Person{name="Alice"}
                       p    → ──────────────► Person{name="Alice"}
                                               (same object!)
After p.name = "Bob":
alice → ─────────────────────────────────────► Person{name="Bob"}
```

> **Interview Q: Is Java pass-by-value or pass-by-reference?**  
> Java is **always pass-by-value**. For primitives, the value itself is copied. For objects, the **reference value** (memory address) is copied — both the original and the parameter point to the same object, so modifications to the object's state are visible to the caller. But **reassigning** the parameter (`p = new Person()`) only changes the local copy — the caller's reference is unaffected. Java never passes a reference to a reference.

---

## 7. Shallow Copy vs Deep Copy

```java
class Address {
    String city;
    Address(String city) { this.city = city; }
}

class Employee {
    String name;
    int salary;
    Address address;   // mutable nested object

    Employee(String name, int salary, Address address) {
        this.name = name;
        this.salary = salary;
        this.address = address;
    }

    // ── SHALLOW COPY — copies reference, not object ──
    Employee shallowCopy() {
        return new Employee(this.name, this.salary, this.address);
        // address is SHARED — both copies point to same Address
    }

    // ── DEEP COPY — copies all nested objects too ──
    Employee deepCopy() {
        return new Employee(this.name, this.salary, new Address(this.address.city));
        // address is a NEW independent object
    }
}

Employee original = new Employee("Alice", 75000, new Address("NYC"));

Employee shallow = original.shallowCopy();
shallow.address.city = "LA";
System.out.println(original.address.city);  // "LA" — SHARED, both affected!

Employee deep = original.deepCopy();
deep.address.city = "Chicago";
System.out.println(original.address.city);  // "LA" — NOT affected, independent!
```

**Ways to deep copy:**

```java
// 1. Manual copy constructor (most common, clearest)
Employee copy = new Employee(original);

// 2. Override clone() + clone nested objects
Employee copy = original.deepClone();

// 3. Serialization (works for entire object graph)
ByteArrayOutputStream bos = new ByteArrayOutputStream();
new ObjectOutputStream(bos).writeObject(original);
Employee copy = (Employee) new ObjectInputStream(
    new ByteArrayInputStream(bos.toByteArray())).readObject();

// 4. Libraries: Apache Commons, Kryo, Gson round-trip
```

> **Interview Q: When is a shallow copy sufficient and when do you need a deep copy?**  
> **Shallow copy** is sufficient when all fields are **primitive or immutable** (String, Integer, LocalDate) — there's no shared mutable state to worry about. **Deep copy** is needed when the object contains **mutable nested objects** — without it, both copies share the same nested object, so modifying one affects the other unexpectedly. Rule of thumb: if your class has any mutable reference fields, a deep copy is needed for true independence.

---

## 8. Reflection API

Reflection allows inspecting and manipulating classes, methods, and fields **at runtime**, even if they're private.

```java
import java.lang.reflect.*;

// ── Getting Class object ──
Class<?> clazz = String.class;
Class<?> clazz2 = "hello".getClass();
Class<?> clazz3 = Class.forName("java.lang.String");  // by fully qualified name

// ── Inspect class info ──
System.out.println(clazz.getName());          // "java.lang.String"
System.out.println(clazz.getSimpleName());    // "String"
System.out.println(clazz.isInterface());      // false
System.out.println(clazz.getSuperclass());    // class java.lang.Object

// ── Inspect fields ──
Field[] fields = clazz.getDeclaredFields();   // all fields (including private)
for (Field f : fields) {
    System.out.println(f.getName() + " : " + f.getType());
}

// ── Inspect methods ──
Method[] methods = clazz.getMethods();        // public methods (inherited too)
Method lengthMethod = clazz.getMethod("length");
System.out.println(lengthMethod.getReturnType());  // int

// ── Invoke methods ──
Method upperCase = String.class.getMethod("toUpperCase");
String result = (String) upperCase.invoke("hello");  // "HELLO"
System.out.println(result);

// ── Access private fields (with care!) ──
class Secret {
    private String code = "XYZ";
}

Secret s = new Secret();
Field field = Secret.class.getDeclaredField("code");
field.setAccessible(true);                    // bypass private access
System.out.println(field.get(s));             // "XYZ"
field.set(s, "ABC");                          // modify private field
System.out.println(field.get(s));             // "ABC"

// ── Create objects dynamically ──
Constructor<StringBuilder> ctor = StringBuilder.class.getConstructor(String.class);
StringBuilder sb = ctor.newInstance("Hello");
System.out.println(sb.toString());   // "Hello"

// ── Annotations via reflection (common in frameworks) ──
@Retention(RetentionPolicy.RUNTIME)
@interface MyAnnotation { String value(); }

@MyAnnotation("test")
class MyClass { }

MyAnnotation ann = MyClass.class.getAnnotation(MyAnnotation.class);
System.out.println(ann.value());   // "test"
```

**Common uses of Reflection:**
- **Frameworks**: Spring, Hibernate, JUnit use reflection to wire beans, map entities, and discover test methods
- **Serialization libraries**: Jackson uses reflection to discover fields
- **Dependency injection**: constructor/field injection
- **Testing**: Mockito uses reflection to intercept method calls

> **Interview Q: What is the Reflection API? What are its drawbacks?**  
> Reflection allows **runtime inspection and manipulation** of classes, methods, and fields — even private ones. It's heavily used by frameworks (Spring, Hibernate, JUnit). Drawbacks: (1) **Performance** — reflection is 10-50x slower than direct calls; (2) **Type safety** — errors are discovered at runtime, not compile time; (3) **Security** — `setAccessible(true)` bypasses access modifiers; (4) **Maintainability** — refactoring tools may not update reflective code. Use it only when necessary (framework code, not application code).
