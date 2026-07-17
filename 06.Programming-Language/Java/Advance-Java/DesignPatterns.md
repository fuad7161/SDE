# Design Patterns — In-Depth Notes

---

## Table of Contents

1. [Singleton — Thread-Safe Variants](#1-singleton--thread-safe-variants)
2. [Factory, Abstract Factory, Builder](#2-factory-abstract-factory-builder)
3. [Proxy, Decorator, Adapter](#3-proxy-decorator-adapter)
4. [Observer, Strategy, Command](#4-observer-strategy-command)
5. [Template Method](#5-template-method)

---

## 1. Singleton — Thread-Safe Variants

Ensures a class has **exactly one instance** and provides a global access point.

---

### ❌ Naive Singleton (Not Thread-Safe)

```java
public class Singleton {
    private static Singleton instance;

    private Singleton() {}

    public static Singleton getInstance() {
        if (instance == null) {             // two threads can both see null
            instance = new Singleton();     // and both create an instance
        }
        return instance;
    }
}
```

---

### ✅ Variant 1 — Synchronized Method (Simple, Slow)

```java
public class Singleton {
    private static Singleton instance;

    private Singleton() {}

    public static synchronized Singleton getInstance() {
        if (instance == null) {
            instance = new Singleton();
        }
        return instance;
    }
}
// Safe, but every call acquires a lock — unnecessary after first initialization
```

---

### ✅ Variant 2 — Double-Checked Locking (Fast + Safe)

```java
public class Singleton {
    private static volatile Singleton instance;  // volatile is mandatory

    private Singleton() {}

    public static Singleton getInstance() {
        if (instance == null) {                     // first check — no lock
            synchronized (Singleton.class) {
                if (instance == null) {             // second check — with lock
                    instance = new Singleton();     // volatile prevents reordering
                }
            }
        }
        return instance;
    }
}
// volatile prevents the JIT from exposing a partially constructed object
// to the first check in another thread
```

---

### ✅ Variant 3 — Initialization-on-Demand Holder (Best Practice)

```java
public class Singleton {
    private Singleton() {}

    // Inner class is not loaded until getInstance() is called
    private static class Holder {
        private static final Singleton INSTANCE = new Singleton();
    }

    public static Singleton getInstance() {
        return Holder.INSTANCE;
    }
}
// Lazy — loaded only on first call
// Thread-safe — guaranteed by class loading mechanism (static initializer)
// No synchronization overhead
```

---

### ✅ Variant 4 — Enum Singleton (Most Robust)

```java
public enum Singleton {
    INSTANCE;

    public void doWork() { ... }
}

// Usage
Singleton.INSTANCE.doWork();
```

- Guaranteed single instance by the JVM
- Serialization-safe (enum deserialize always returns the same instance)
- Reflection-safe (cannot call private constructor via reflection)
- Recommended by Joshua Bloch (*Effective Java*)

---

### Comparison

| Variant | Thread-Safe | Lazy | Reflection-Safe | Serialization-Safe |
|---|---|---|---|---|
| Naive | ❌ | ✅ | ❌ | ❌ |
| Synchronized method | ✅ | ✅ | ❌ | ❌ |
| Double-checked locking | ✅ | ✅ | ❌ | ❌ |
| Holder pattern | ✅ | ✅ | ❌ | ❌ |
| Enum | ✅ | ❌ | ✅ | ✅ |

---

## 2. Factory, Abstract Factory, Builder

### Factory Method

Defines an interface for creating an object but lets **subclasses decide which class to instantiate**.

```
Creator
  └─ factoryMethod() → Product
        ▲                  ▲
ConcreteCreator    ConcreteProduct
```

```java
// Product interface
interface Notification {
    void send(String message);
}

// Concrete products
class EmailNotification implements Notification {
    public void send(String message) {
        System.out.println("Email: " + message);
    }
}
class SMSNotification implements Notification {
    public void send(String message) {
        System.out.println("SMS: " + message);
    }
}
class PushNotification implements Notification {
    public void send(String message) {
        System.out.println("Push: " + message);
    }
}

// Factory
class NotificationFactory {
    public static Notification create(String type) {
        return switch (type.toUpperCase()) {
            case "EMAIL" -> new EmailNotification();
            case "SMS"   -> new SMSNotification();
            case "PUSH"  -> new PushNotification();
            default      -> throw new IllegalArgumentException("Unknown type: " + type);
        };
    }
}

// Client — decoupled from concrete classes
Notification n = NotificationFactory.create("EMAIL");
n.send("Welcome!");
```

---

### Abstract Factory

Creates **families of related objects** without specifying their concrete classes.  
Think of it as a factory of factories.

```
AbstractFactory
  ├─ createButton()  → Button
  └─ createCheckbox()→ Checkbox
        ▲                  ▲
WindowsFactory          MacFactory
  ├─ WindowsButton      MacButton
  └─ WindowsCheckbox    MacCheckbox
```

```java
// Abstract products
interface Button   { void render(); }
interface Checkbox { void render(); }

// Concrete products — Windows family
class WindowsButton   implements Button   { public void render() { System.out.println("Windows Button");   } }
class WindowsCheckbox implements Checkbox { public void render() { System.out.println("Windows Checkbox"); } }

// Concrete products — Mac family
class MacButton   implements Button   { public void render() { System.out.println("Mac Button");   } }
class MacCheckbox implements Checkbox { public void render() { System.out.println("Mac Checkbox"); } }

// Abstract factory
interface GUIFactory {
    Button   createButton();
    Checkbox createCheckbox();
}

// Concrete factories
class WindowsFactory implements GUIFactory {
    public Button   createButton()   { return new WindowsButton();   }
    public Checkbox createCheckbox() { return new WindowsCheckbox(); }
}
class MacFactory implements GUIFactory {
    public Button   createButton()   { return new MacButton();   }
    public Checkbox createCheckbox() { return new MacCheckbox(); }
}

// Client — works with any factory/family
class Application {
    private final Button button;
    private final Checkbox checkbox;

    Application(GUIFactory factory) {
        this.button   = factory.createButton();
        this.checkbox = factory.createCheckbox();
    }

    void render() { button.render(); checkbox.render(); }
}

// Usage
GUIFactory factory = isWindows() ? new WindowsFactory() : new MacFactory();
new Application(factory).render();
```

---

### Builder

Constructs **complex objects step by step**, separating construction from representation.  
Avoids telescoping constructors.

```java
// Without Builder — telescoping constructor problem
new Pizza("large", true, false, true, false, "thin", "tomato"); // unreadable

// With Builder
class Pizza {
    private final String size;
    private final boolean cheese;
    private final boolean pepperoni;
    private final boolean mushrooms;
    private final String crust;

    private Pizza(Builder builder) {
        this.size      = builder.size;
        this.cheese    = builder.cheese;
        this.pepperoni = builder.pepperoni;
        this.mushrooms = builder.mushrooms;
        this.crust     = builder.crust;
    }

    public static class Builder {
        private final String size;   // required
        private boolean cheese    = false;
        private boolean pepperoni = false;
        private boolean mushrooms = false;
        private String  crust     = "regular";

        public Builder(String size) { this.size = size; }

        public Builder cheese()         { this.cheese    = true;   return this; }
        public Builder pepperoni()      { this.pepperoni = true;   return this; }
        public Builder mushrooms()      { this.mushrooms = true;   return this; }
        public Builder crust(String c)  { this.crust     = c;      return this; }

        public Pizza build() {
            // validation here if needed
            return new Pizza(this);
        }
    }
}

// Usage — fluent, self-documenting
Pizza pizza = new Pizza.Builder("large")
    .cheese()
    .pepperoni()
    .crust("thin")
    .build();
```

**Java's built-in builders**: `StringBuilder`, `Stream.Builder`, `Locale.Builder`, Lombok's `@Builder`.

---

### Factory vs Abstract Factory vs Builder

| Pattern | Intent | Creates |
|---|---|---|
| Factory Method | Delegate instantiation to subclass | One product type |
| Abstract Factory | Create families of related products | Multiple related products |
| Builder | Construct complex object step-by-step | One complex object |

---

## 3. Proxy, Decorator, Adapter

### Proxy

Provides a **surrogate or placeholder** to control access to another object — same interface as the real object.

**Use cases**: lazy initialization, access control, logging, caching, remote proxies.

```java
interface DatabaseService {
    String query(String sql);
}

class RealDatabaseService implements DatabaseService {
    public String query(String sql) {
        System.out.println("Executing: " + sql);
        return "result";
    }
}

// Proxy — adds logging + access control without changing real service
class DatabaseServiceProxy implements DatabaseService {
    private final RealDatabaseService real = new RealDatabaseService();
    private final String currentUser;

    DatabaseServiceProxy(String user) { this.currentUser = user; }

    public String query(String sql) {
        if (!hasPermission(currentUser)) {
            throw new SecurityException("Access denied for: " + currentUser);
        }
        System.out.println("[LOG] " + currentUser + " querying: " + sql);
        long start = System.currentTimeMillis();
        String result = real.query(sql);
        System.out.println("[LOG] Took " + (System.currentTimeMillis() - start) + "ms");
        return result;
    }

    private boolean hasPermission(String user) { return !user.equals("guest"); }
}

// Client sees the same interface
DatabaseService db = new DatabaseServiceProxy("alice");
db.query("SELECT * FROM users");
```

**JDK Dynamic Proxy** (used by Spring AOP):

```java
DatabaseService proxy = (DatabaseService) Proxy.newProxyInstance(
    DatabaseService.class.getClassLoader(),
    new Class[]{DatabaseService.class},
    (proxyObj, method, args) -> {
        System.out.println("Before: " + method.getName());
        Object result = method.invoke(new RealDatabaseService(), args);
        System.out.println("After: " + method.getName());
        return result;
    }
);
```

---

### Decorator

**Wraps** an object to add new behaviour dynamically — same interface, but behaviour is composable.  
Unlike inheritance, decorators can be stacked at runtime.

```java
interface TextProcessor {
    String process(String text);
}

// Base
class PlainTextProcessor implements TextProcessor {
    public String process(String text) { return text; }
}

// Decorators — each wraps another TextProcessor
class TrimDecorator implements TextProcessor {
    private final TextProcessor wrapped;
    TrimDecorator(TextProcessor t) { this.wrapped = t; }
    public String process(String text) { return wrapped.process(text).trim(); }
}

class UpperCaseDecorator implements TextProcessor {
    private final TextProcessor wrapped;
    UpperCaseDecorator(TextProcessor t) { this.wrapped = t; }
    public String process(String text) { return wrapped.process(text).toUpperCase(); }
}

class HtmlEscapeDecorator implements TextProcessor {
    private final TextProcessor wrapped;
    HtmlEscapeDecorator(TextProcessor t) { this.wrapped = t; }
    public String process(String text) {
        return wrapped.process(text)
            .replace("&", "&amp;").replace("<", "&lt;").replace(">", "&gt;");
    }
}

// Compose at runtime — stacking behaviours like layers
TextProcessor processor = new HtmlEscapeDecorator(
                            new UpperCaseDecorator(
                              new TrimDecorator(
                                new PlainTextProcessor())));

System.out.println(processor.process("  hello <world>  "));
// → "HELLO &lt;WORLD&gt;"
```

**Java's built-in decorators**: `BufferedReader(new FileReader(...))`, `Collections.unmodifiableList(list)`, `Collections.synchronizedList(list)`.

---

### Adapter

**Converts** the interface of a class into another interface the client expects.  
Bridges incompatible interfaces — like a plug adapter.

```java
// Existing class with incompatible interface
class LegacyPaymentGateway {
    public void makePayment(double amount, String currency) {
        System.out.println("Legacy paying " + amount + " " + currency);
    }
}

// Target interface our application expects
interface PaymentProcessor {
    void pay(int amountInCents, String currencyCode);
}

// Adapter — wraps Legacy, exposes modern interface
class PaymentAdapter implements PaymentProcessor {
    private final LegacyPaymentGateway legacy;

    PaymentAdapter(LegacyPaymentGateway gateway) { this.legacy = gateway; }

    @Override
    public void pay(int amountInCents, String currencyCode) {
        double amount = amountInCents / 100.0;   // convert cents → decimal
        legacy.makePayment(amount, currencyCode);
    }
}

// Client works with PaymentProcessor — unaware of legacy system
PaymentProcessor processor = new PaymentAdapter(new LegacyPaymentGateway());
processor.pay(4999, "USD");  // → "Legacy paying 49.99 USD"
```

---

### Proxy vs Decorator vs Adapter

| Pattern | Intent | Interface | Adds behaviour |
|---|---|---|---|
| Proxy | Control access to the same object | Same as subject | Access control, lazy load, logging |
| Decorator | Add behaviour dynamically, stackable | Same as component | Yes — new features |
| Adapter | Convert incompatible interfaces | Different from adaptee | No — translation only |

---

## 4. Observer, Strategy, Command

### Observer

Defines a **one-to-many dependency** — when one object changes state, all dependants are notified automatically.  
Also called Publish-Subscribe.

```java
// Observer interface
interface EventListener {
    void onEvent(String eventType, Object data);
}

// Subject (Observable)
class EventManager {
    private final Map<String, List<EventListener>> listeners = new HashMap<>();

    public void subscribe(String eventType, EventListener listener) {
        listeners.computeIfAbsent(eventType, k -> new ArrayList<>()).add(listener);
    }

    public void unsubscribe(String eventType, EventListener listener) {
        listeners.getOrDefault(eventType, List.of()).remove(listener);
    }

    public void notify(String eventType, Object data) {
        listeners.getOrDefault(eventType, List.of())
                 .forEach(l -> l.onEvent(eventType, data));
    }
}

// Concrete observers
class EmailAlertListener implements EventListener {
    public void onEvent(String type, Object data) {
        System.out.println("Email alert — " + type + ": " + data);
    }
}

class AuditLogListener implements EventListener {
    public void onEvent(String type, Object data) {
        System.out.println("Audit log — " + type + ": " + data);
    }
}

// Usage
EventManager events = new EventManager();
events.subscribe("USER_LOGIN",  new EmailAlertListener());
events.subscribe("USER_LOGIN",  new AuditLogListener());
events.subscribe("FILE_UPLOAD", new AuditLogListener());

events.notify("USER_LOGIN", "alice");
// → Email alert — USER_LOGIN: alice
// → Audit log  — USER_LOGIN: alice
```

---

### Strategy

Defines a family of algorithms, **encapsulates each one**, and makes them **interchangeable** at runtime.  
Replaces conditionals with polymorphism.

```java
// Strategy interface
interface SortStrategy {
    void sort(int[] array);
}

// Concrete strategies
class BubbleSort implements SortStrategy {
    public void sort(int[] arr) {
        System.out.println("Bubble sorting " + arr.length + " elements");
        // bubble sort impl
    }
}
class QuickSort implements SortStrategy {
    public void sort(int[] arr) {
        System.out.println("Quick sorting " + arr.length + " elements");
        // quicksort impl
    }
}
class MergeSort implements SortStrategy {
    public void sort(int[] arr) {
        System.out.println("Merge sorting " + arr.length + " elements");
        // mergesort impl
    }
}

// Context — uses a strategy
class Sorter {
    private SortStrategy strategy;

    public void setStrategy(SortStrategy strategy) { this.strategy = strategy; }

    public void sort(int[] array) {
        strategy.sort(array);
    }
}

// Switch strategy at runtime
Sorter sorter = new Sorter();
int[] data = {5, 3, 8, 1};

sorter.setStrategy(new QuickSort());
sorter.sort(data);                    // Quick sorting...

sorter.setStrategy(new MergeSort());  // swap strategy without changing Sorter
sorter.sort(data);                    // Merge sorting...

// With lambdas (functional approach)
sorter.setStrategy(arr -> Arrays.sort(arr));   // JDK built-in as strategy
```

---

### Command

**Encapsulates a request as an object**, allowing parameterization, queuing, logging, and undo/redo.

```java
// Command interface
interface Command {
    void execute();
    void undo();
}

// Receiver — the object that does the actual work
class TextEditor {
    private StringBuilder text = new StringBuilder();

    public void insertText(String t) { text.append(t); }
    public void deleteText(int len)  { text.delete(text.length() - len, text.length()); }
    public String getText()          { return text.toString(); }
}

// Concrete commands
class InsertCommand implements Command {
    private final TextEditor editor;
    private final String text;

    InsertCommand(TextEditor editor, String text) {
        this.editor = editor;
        this.text   = text;
    }

    public void execute() { editor.insertText(text); }
    public void undo()    { editor.deleteText(text.length()); }
}

// Invoker — holds and executes commands, supports undo stack
class CommandHistory {
    private final Deque<Command> history = new ArrayDeque<>();

    public void execute(Command cmd) {
        cmd.execute();
        history.push(cmd);
    }

    public void undo() {
        if (!history.isEmpty()) history.pop().undo();
    }
}

// Usage
TextEditor editor = new TextEditor();
CommandHistory history = new CommandHistory();

history.execute(new InsertCommand(editor, "Hello"));
history.execute(new InsertCommand(editor, " World"));
System.out.println(editor.getText());   // "Hello World"

history.undo();
System.out.println(editor.getText());   // "Hello"

history.undo();
System.out.println(editor.getText());   // ""
```

---

### Observer vs Strategy vs Command

| Pattern | Intent | Key Benefit |
|---|---|---|
| Observer | Notify many objects on state change | Decoupled event broadcasting |
| Strategy | Swap algorithms at runtime | Eliminate conditionals, open/closed |
| Command | Encapsulate request as object | Undo/redo, queuing, logging |

---

## 5. Template Method

Defines the **skeleton of an algorithm** in a base class, deferring specific steps to subclasses.  
The overall structure (template) is fixed; individual steps can be customised.

```
AbstractClass
  └─ templateMethod()          ← final — defines the algorithm skeleton
       ├─ step1()              ← concrete — shared implementation
       ├─ step2()              ← abstract — subclass must implement
       ├─ step3()              ← concrete — shared implementation
       └─ hook()               ← hook — optional override (default no-op)
```

```java
// Abstract class with the template
abstract class DataProcessor {

    // Template method — final prevents subclasses from changing the algorithm
    public final void process() {
        readData();
        processData();
        if (shouldValidate()) {    // hook — subclass can override
            validateData();
        }
        writeData();
    }

    private void readData() {
        System.out.println("Reading data from source");
    }

    protected abstract void processData();   // subclass defines this step

    private void validateData() {
        System.out.println("Validating data");
    }

    private void writeData() {
        System.out.println("Writing data to destination");
    }

    // Hook method — optional override, default does nothing
    protected boolean shouldValidate() { return true; }
}

// CSV implementation
class CsvDataProcessor extends DataProcessor {
    @Override
    protected void processData() {
        System.out.println("Parsing CSV rows and mapping columns");
    }
}

// JSON implementation — skips validation
class JsonDataProcessor extends DataProcessor {
    @Override
    protected void processData() {
        System.out.println("Deserializing JSON and transforming fields");
    }

    @Override
    protected boolean shouldValidate() { return false; }  // override hook
}

// Usage
System.out.println("=== CSV ===");
new CsvDataProcessor().process();
// Reading data from source
// Parsing CSV rows and mapping columns
// Validating data
// Writing data to destination

System.out.println("=== JSON ===");
new JsonDataProcessor().process();
// Reading data from source
// Deserializing JSON and transforming fields
// Writing data to destination  ← validation skipped
```

---

### Real-World Uses of Template Method

| Framework | Template Class | Hook/Step |
|---|---|---|
| Spring | `AbstractController` | `handleRequestInternal()` |
| Spring | `JdbcTemplate` | `query()` with `RowMapper` callback |
| Spring | `AbstractBatchConfiguration` | `reader()`, `processor()`, `writer()` |
| Java | `AbstractList` | `get(int index)`, `size()` |
| Java | `HttpServlet` | `doGet()`, `doPost()` |

---

### Template Method vs Strategy

Both allow varying part of an algorithm — key difference is **mechanism**:

| | Template Method | Strategy |
|---|---|---|
| Mechanism | Inheritance (subclass overrides steps) | Composition (inject strategy object) |
| Variation scope | Part of algorithm in subclass | Entire algorithm swapped |
| Coupling | Tighter (extends abstract class) | Looser (depends on interface) |
| Runtime swap | ❌ fixed at compile-time | ✅ swap strategy at runtime |
| Prefer when | Steps share lots of common code | Algorithm families need runtime switching |
