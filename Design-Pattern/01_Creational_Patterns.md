# 🔵 Creational Patterns

> **Category:** Creational &nbsp;|&nbsp; **Tags:** `Singleton` `Factory Method` `Abstract Factory` `Builder` `Prototype`

Creational patterns deal with **object creation mechanisms** — controlling how objects are instantiated to improve flexibility and reuse.

---

## Table of Contents
1. [Singleton](#1-singleton)
2. [Factory Method](#2-factory-method)
3. [Abstract Factory](#3-abstract-factory)
4. [Builder](#4-builder)
5. [Prototype](#5-prototype)
6. [Interview Questions](#interview-questions)

---

## 1. Singleton

**Intent:** Ensure a class has only **one instance** and provide a global access point to it.

**When to use:**
- Shared resource (DB connection pool, config, logger, thread pool)
- Exactly one object must coordinate actions across the system

### Naive (not thread-safe)

<details>
<summary><b>Java</b></summary>

```java
public class Singleton {
    private static Singleton instance;

    private Singleton() {}   // prevent external instantiation

    public static Singleton getInstance() {
        if (instance == null) {
            instance = new Singleton();   // ⚠️ race condition in multithreading
        }
        return instance;
    }
}
```

</details>

<details>
<summary><b>Go</b></summary>

```go
package main

var instance *Singleton

type Singleton struct{}

func GetInstance() *Singleton {
    if instance == nil {
        instance = &Singleton{} // ⚠️ race condition in multithreading
    }
    return instance
}
```

</details>

### Thread-safe — Double-Checked Locking (recommended)

<details>
<summary><b>Java</b></summary>

```java
public class Singleton {
    // volatile ensures visibility across threads & prevents instruction reordering
    private static volatile Singleton instance;

    private Singleton() {}

    public static Singleton getInstance() {
        if (instance == null) {                  // 1st check — no locking
            synchronized (Singleton.class) {
                if (instance == null) {          // 2nd check — inside lock
                    instance = new Singleton();
                }
            }
        }
        return instance;
    }
}

// Usage
Singleton s = Singleton.getInstance();
```

</details>

<details>
<summary><b>Go — <code>sync.Once</code> (recommended equivalent)</b></summary>

```go
package main

import "sync"

type Singleton struct{}

var (
    instance *Singleton
    once     sync.Once
)

func GetInstance() *Singleton {
    once.Do(func() {
        instance = &Singleton{}
    })
    return instance
}

func (s *Singleton) DoWork() {
    println("Doing work...")
}

func main() {
    s := GetInstance()
    s.DoWork()
}
```

</details>

### Enum Singleton (simplest, inherently thread-safe)

<details>
<summary><b>Java</b></summary>

```java
public enum Singleton {
    INSTANCE;

    public void doWork() {
        System.out.println("Doing work...");
    }
}

// Usage
Singleton.INSTANCE.doWork();
```
> Java guarantees enum instances are created exactly once and are serialization-safe.

</details>

<details>
<summary><b>Go — package-level init (closest equivalent)</b></summary>

```go
package main

// Go has no enum; use a package-level var initialized at startup
type Singleton struct{}

var Instance = &Singleton{}  // initialized once when package loads

func (s *Singleton) DoWork() {
    println("Doing work...")
}

func main() {
    Instance.DoWork()
}
```

</details>

### Bill Pugh (Initialization-on-demand holder)

<details>
<summary><b>Java</b></summary>

```java
public class Singleton {
    private Singleton() {}

    // Inner class loaded only when getInstance() is first called
    private static class Holder {
        private static final Singleton INSTANCE = new Singleton();
    }

    public static Singleton getInstance() {
        return Holder.INSTANCE;
    }
}
```

</details>

<details>
<summary><b>Go — lazy init with <code>sync.Once</code> (same result)</b></summary>

```go
package main

import "sync"

type Singleton struct{}

var (
    instance *Singleton
    once     sync.Once
)

// GetInstance is lazy — Singleton is only created on first call
func GetInstance() *Singleton {
    once.Do(func() {
        instance = &Singleton{}
    })
    return instance
}
```

</details>

---

## 2. Factory Method

**Intent:** Define an interface for creating an object, but let **subclasses decide** which class to instantiate. Delegates instantiation to subclasses.

**When to use:**
- The exact type of object to create isn't known at compile time
- Subclasses should control what gets created
- Decouple client from concrete classes

<details>
<summary><b>Java</b></summary>

```java
// Product interface
public interface Notification {
    void send(String message);
}

// Concrete Products
public class EmailNotification implements Notification {
    @Override
    public void send(String message) {
        System.out.println("Email: " + message);
    }
}

public class SMSNotification implements Notification {
    @Override
    public void send(String message) {
        System.out.println("SMS: " + message);
    }
}

public class PushNotification implements Notification {
    @Override
    public void send(String message) {
        System.out.println("Push: " + message);
    }
}

// Creator — defines the factory method
public abstract class NotificationService {
    // Factory method — subclasses override this
    protected abstract Notification createNotification();

    public void notifyUser(String message) {
        Notification n = createNotification();
        n.send(message);
    }
}

// Concrete Creators
public class EmailService extends NotificationService {
    @Override
    protected Notification createNotification() {
        return new EmailNotification();
    }
}

public class SMSService extends NotificationService {
    @Override
    protected Notification createNotification() {
        return new SMSNotification();
    }
}

// Usage
public class Main {
    public static void main(String[] args) {
        NotificationService service = new EmailService();
        service.notifyUser("Your order has shipped!");  // Email: Your order has shipped!

        service = new SMSService();
        service.notifyUser("OTP: 123456");              // SMS: OTP: 123456
    }
}
```

</details>

<details>
<summary><b>Go</b></summary>

```go
package main

import "fmt"

type Notification interface {
    Send(message string)
}

type EmailNotification struct{}
type SMSNotification struct{}
type PushNotification struct{}

func (e *EmailNotification) Send(msg string) { fmt.Println("Email:", msg) }
func (s *SMSNotification) Send(msg string)   { fmt.Println("SMS:", msg) }
func (p *PushNotification) Send(msg string)  { fmt.Println("Push:", msg) }

// Factory Method via function field
type NotificationService struct {
    createFn func() Notification
}

func (s *NotificationService) NotifyUser(message string) {
    n := s.createFn()
    n.Send(message)
}

func NewEmailService() *NotificationService {
    return &NotificationService{
        createFn: func() Notification { return &EmailNotification{} },
    }
}

func NewSMSService() *NotificationService {
    return &NotificationService{
        createFn: func() Notification { return &SMSNotification{} },
    }
}

func main() {
    service := NewEmailService()
    service.NotifyUser("Your order has shipped!") // Email: Your order has shipped!

    service = NewSMSService()
    service.NotifyUser("OTP: 123456")             // SMS: OTP: 123456
}
```

</details>

### Simple Factory (not a GoF pattern, but commonly asked)

<details>
<summary><b>Java</b></summary>

```java
public class NotificationFactory {
    public static Notification create(String type) {
        return switch (type.toLowerCase()) {
            case "email" -> new EmailNotification();
            case "sms"   -> new SMSNotification();
            case "push"  -> new PushNotification();
            default      -> throw new IllegalArgumentException("Unknown type: " + type);
        };
    }
}

// Usage
Notification n = NotificationFactory.create("email");
n.send("Hello!");
```

</details>

<details>
<summary><b>Go</b></summary>

```go
func NewNotification(notifType string) Notification {
    switch notifType {
    case "email":
        return &EmailNotification{}
    case "sms":
        return &SMSNotification{}
    case "push":
        return &PushNotification{}
    default:
        panic("unknown notification type: " + notifType)
    }
}

// Usage
n := NewNotification("email")
n.Send("Hello!")
```

</details>

---

## 3. Abstract Factory

**Intent:** Provide an interface for creating **families of related objects** without specifying their concrete classes.

**When to use:**
- System must be independent of how products are created
- System must work with multiple families of products (e.g., UI theme: Light / Dark)
- You want to enforce that products from the same family are used together

<details>
<summary><b>Java</b></summary>

```java
// Abstract Products
public interface Button {
    void render();
    void onClick();
}

public interface Checkbox {
    void render();
}

// Concrete Products — Windows family
public class WindowsButton implements Button {
    @Override public void render() { System.out.println("[Windows Button]"); }
    @Override public void onClick() { System.out.println("Windows button clicked"); }
}

public class WindowsCheckbox implements Checkbox {
    @Override public void render() { System.out.println("[Windows Checkbox]"); }
}

// Concrete Products — Mac family
public class MacButton implements Button {
    @Override public void render() { System.out.println("[Mac Button]"); }
    @Override public void onClick() { System.out.println("Mac button clicked"); }
}

public class MacCheckbox implements Checkbox {
    @Override public void render() { System.out.println("[Mac Checkbox]"); }
}

// Abstract Factory
public interface UIFactory {
    Button createButton();
    Checkbox createCheckbox();
}

// Concrete Factories
public class WindowsFactory implements UIFactory {
    @Override public Button createButton()     { return new WindowsButton(); }
    @Override public Checkbox createCheckbox() { return new WindowsCheckbox(); }
}

public class MacFactory implements UIFactory {
    @Override public Button createButton()     { return new MacButton(); }
    @Override public Checkbox createCheckbox() { return new MacCheckbox(); }
}

// Client — works with any factory
public class Application {
    private final Button button;
    private final Checkbox checkbox;

    public Application(UIFactory factory) {
        this.button   = factory.createButton();
        this.checkbox = factory.createCheckbox();
    }

    public void render() {
        button.render();
        checkbox.render();
    }
}

// Usage
public class Main {
    public static void main(String[] args) {
        String os = System.getProperty("os.name").toLowerCase();
        UIFactory factory = os.contains("win") ? new WindowsFactory() : new MacFactory();

        Application app = new Application(factory);
        app.render();
    }
}
```

</details>

<details>
<summary><b>Go</b></summary>

```go
package main

import "fmt"

type Button interface {
    Render()
    OnClick()
}

type Checkbox interface {
    Render()
}

// Windows family
type WindowsButton struct{}
type WindowsCheckbox struct{}

func (b *WindowsButton) Render()   { fmt.Println("[Windows Button]") }
func (b *WindowsButton) OnClick()  { fmt.Println("Windows button clicked") }
func (c *WindowsCheckbox) Render() { fmt.Println("[Windows Checkbox]") }

// Mac family
type MacButton struct{}
type MacCheckbox struct{}

func (b *MacButton) Render()   { fmt.Println("[Mac Button]") }
func (b *MacButton) OnClick()  { fmt.Println("Mac button clicked") }
func (c *MacCheckbox) Render() { fmt.Println("[Mac Checkbox]") }

// Abstract Factory
type UIFactory interface {
    CreateButton() Button
    CreateCheckbox() Checkbox
}

type WindowsFactory struct{}
type MacFactory struct{}

func (f *WindowsFactory) CreateButton() Button     { return &WindowsButton{} }
func (f *WindowsFactory) CreateCheckbox() Checkbox { return &WindowsCheckbox{} }
func (f *MacFactory) CreateButton() Button         { return &MacButton{} }
func (f *MacFactory) CreateCheckbox() Checkbox     { return &MacCheckbox{} }

// Client
type Application struct {
    button   Button
    checkbox Checkbox
}

func NewApplication(factory UIFactory) *Application {
    return &Application{
        button:   factory.CreateButton(),
        checkbox: factory.CreateCheckbox(),
    }
}

func (a *Application) Render() {
    a.button.Render()
    a.checkbox.Render()
}

func main() {
    var factory UIFactory = &MacFactory{} // swap to &WindowsFactory{} as needed
    app := NewApplication(factory)
    app.Render()
}
```

</details>

### Factory Method vs Abstract Factory

| | Factory Method | Abstract Factory |
|-|---------------|-----------------|
| Creates | One product | Family of related products |
| Mechanism | Inheritance (subclass overrides) | Composition (factory passed in) |
| Extensibility | Add new subclass | Add new factory class |

---

## 4. Builder

**Intent:** Construct a **complex object step-by-step**. Separate the construction from its representation so the same construction process can create different representations.

**When to use:**
- Object requires many optional parameters (avoids telescoping constructors)
- Object construction involves multiple steps
- Need to create different representations of the same object

<details>
<summary><b>Java</b></summary>

```java
// Product
public class HttpRequest {
    private final String url;         // required
    private final String method;      // required
    private final String body;        // optional
    private final int    timeout;     // optional
    private final Map<String, String> headers;  // optional

    // Private constructor — only Builder can call it
    private HttpRequest(Builder builder) {
        this.url     = builder.url;
        this.method  = builder.method;
        this.body    = builder.body;
        this.timeout = builder.timeout;
        this.headers = Collections.unmodifiableMap(builder.headers);
    }

    @Override
    public String toString() {
        return method + " " + url + " | body=" + body + " | timeout=" + timeout;
    }

    // Static inner Builder
    public static class Builder {
        // Required
        private final String url;
        private final String method;
        // Optional with defaults
        private String body    = "";
        private int    timeout = 30;
        private Map<String, String> headers = new HashMap<>();

        public Builder(String url, String method) {
            this.url    = url;
            this.method = method;
        }

        public Builder body(String body) {
            this.body = body;
            return this;   // fluent — enables chaining
        }

        public Builder timeout(int seconds) {
            this.timeout = seconds;
            return this;
        }

        public Builder header(String key, String value) {
            this.headers.put(key, value);
            return this;
        }

        public HttpRequest build() {
            // Validate here
            if (url == null || url.isBlank()) {
                throw new IllegalStateException("URL is required");
            }
            return new HttpRequest(this);
        }
    }
}

// Usage
HttpRequest request = new HttpRequest.Builder("https://api.example.com/users", "POST")
        .body("{\"name\":\"Alice\"}")
        .timeout(60)
        .header("Content-Type", "application/json")
        .header("Authorization", "Bearer token123")
        .build();

System.out.println(request);
// POST https://api.example.com/users | body={"name":"Alice"} | timeout=60
```

</details>

<details>
<summary><b>Go</b></summary>

```go
package main

import "fmt"

type HttpRequest struct {
    url     string
    method  string
    body    string
    timeout int
    headers map[string]string
}

func (r *HttpRequest) String() string {
    return fmt.Sprintf("%s %s | body=%s | timeout=%d", r.method, r.url, r.body, r.timeout)
}

// Builder
type HttpRequestBuilder struct {
    url     string
    method  string
    body    string
    timeout int
    headers map[string]string
}

func NewRequest(url, method string) *HttpRequestBuilder {
    return &HttpRequestBuilder{
        url:     url,
        method:  method,
        timeout: 30,
        headers: make(map[string]string),
    }
}

func (b *HttpRequestBuilder) Body(body string) *HttpRequestBuilder {
    b.body = body
    return b
}

func (b *HttpRequestBuilder) Timeout(seconds int) *HttpRequestBuilder {
    b.timeout = seconds
    return b
}

func (b *HttpRequestBuilder) Header(key, value string) *HttpRequestBuilder {
    b.headers[key] = value
    return b
}

func (b *HttpRequestBuilder) Build() (*HttpRequest, error) {
    if b.url == "" {
        return nil, fmt.Errorf("URL is required")
    }
    return &HttpRequest{
        url:     b.url,
        method:  b.method,
        body:    b.body,
        timeout: b.timeout,
        headers: b.headers,
    }, nil
}

func main() {
    request, err := NewRequest("https://api.example.com/users", "POST").
        Body(`{"name":"Alice"}`).
        Timeout(60).
        Header("Content-Type", "application/json").
        Header("Authorization", "Bearer token123").
        Build()

    if err != nil {
        panic(err)
    }
    fmt.Println(request)
    // POST https://api.example.com/users | body={"name":"Alice"} | timeout=60
}
```

</details>

### Builder vs Constructor

| | Telescoping Constructor | Builder |
|--|------------------------|---------|
| Many optional params | Unreadable | Named, clear |
| Immutability | ✅ | ✅ |
| Validation | Mixed in constructor | Centralized in `build()` |
| Readability | Poor | Excellent |

> Java's `StringBuilder`, Lombok's `@Builder`, and `ProcessBuilder` all use this pattern.

---

## 5. Prototype

**Intent:** Create new objects by **cloning an existing object** (the prototype) instead of creating from scratch.

**When to use:**
- Object creation is expensive (DB query, complex calculation)
- You need many similar objects with small differences
- Avoid subclassing for object creation

<details>
<summary><b>Java</b></summary>

```java
// Prototype interface
public interface Cloneable {
    Object clone();
}

// Concrete Prototype
public class Employee implements Cloneable {
    private String name;
    private String department;
    private List<String> skills;  // mutable — must handle carefully

    public Employee(String name, String department, List<String> skills) {
        this.name       = name;
        this.department = department;
        this.skills     = skills;
    }

    // === Shallow Copy ===
    // skills list is shared — changes in copy affect original
    @Override
    public Employee clone() {
        try {
            return (Employee) super.clone();  // Object.clone() = shallow copy
        } catch (CloneNotSupportedException e) {
            throw new RuntimeException(e);
        }
    }

    // === Deep Copy ===
    // skills list is duplicated — copy is fully independent
    public Employee deepCopy() {
        return new Employee(this.name, this.department, new ArrayList<>(this.skills));
    }

    public void addSkill(String skill)    { skills.add(skill); }
    public String getName()               { return name; }
    public void setName(String name)      { this.name = name; }
    public List<String> getSkills()       { return skills; }

    @Override
    public String toString() {
        return name + " | " + department + " | skills=" + skills;
    }
}

// Usage
public class Main {
    public static void main(String[] args) {
        Employee original = new Employee("Alice", "Engineering",
                new ArrayList<>(List.of("Java", "SQL")));

        // ---- Shallow Copy ----
        Employee shallow = original.clone();
        shallow.setName("Bob");
        shallow.addSkill("Python");   // ← also modifies original.skills!

        System.out.println(original);  // Alice | Engineering | [Java, SQL, Python] ← modified!
        System.out.println(shallow);   // Bob   | Engineering | [Java, SQL, Python]

        // ---- Deep Copy ----
        Employee template = new Employee("Template", "Engineering",
                new ArrayList<>(List.of("Java", "SQL")));
        Employee deep = template.deepCopy();
        deep.setName("Carol");
        deep.addSkill("Kotlin");      // ← does NOT affect template

        System.out.println(template); // Template | Engineering | [Java, SQL]
        System.out.println(deep);     // Carol    | Engineering | [Java, SQL, Kotlin]
    }
}
```

</details>

<details>
<summary><b>Go</b></summary>

```go
package main

import "fmt"

type Employee struct {
    Name       string
    Department string
    Skills     []string
}

// ShallowCopy — Skills slice header copied, backing array shared
func (e *Employee) ShallowCopy() *Employee {
    copy := *e // struct copy = shallow
    return &copy
}

// DeepCopy — Skills slice fully duplicated, fully independent
func (e *Employee) DeepCopy() *Employee {
    skills := make([]string, len(e.Skills))
    copy(skills, e.Skills)
    return &Employee{
        Name:       e.Name,
        Department: e.Department,
        Skills:     skills,
    }
}

func main() {
    original := &Employee{"Alice", "Engineering", []string{"Java", "SQL"}}

    // ---- Shallow Copy ----
    shallow := original.ShallowCopy()
    shallow.Name = "Bob"
    shallow.Skills = append(shallow.Skills, "Python") // ← also modifies original!

    fmt.Println(original) // &{Alice Engineering [Java SQL Python]} ← modified!
    fmt.Println(shallow)  // &{Bob   Engineering [Java SQL Python]}

    // ---- Deep Copy ----
    template := &Employee{"Template", "Engineering", []string{"Java", "SQL"}}
    deep := template.DeepCopy()
    deep.Name = "Carol"
    deep.Skills = append(deep.Skills, "Kotlin") // does NOT affect template

    fmt.Println(template) // &{Template Engineering [Java SQL]}
    fmt.Println(deep)     // &{Carol    Engineering [Java SQL Kotlin]}
}
```

</details>

### Shallow vs Deep Copy

| | Shallow Copy | Deep Copy |
|--|-------------|-----------|
| Primitive fields | Copied | Copied |
| Object/Array fields | Reference copied (shared) | New copy created (independent) |
| Changes in copy affect original | ✅ (for objects) | ❌ |
| Implementation | `Object.clone()` | Manual / copy constructor / serialization |

---

## Interview Questions

### Q1. How do you make a Singleton thread-safe in Java? What are the different approaches?

> **Answer:**
> 1. **Synchronized method:** Simplest but slowest — synchronizes on every call.
> 2. **Double-checked locking + `volatile`:** Fast after initialization; `volatile` prevents instruction reordering. Pre-Java 5 this was broken without `volatile`.
> 3. **Enum Singleton:** Preferred by Joshua Bloch (Effective Java). JVM guarantees one instance, handles serialization automatically, and can't be broken by reflection.
> 4. **Bill Pugh (Holder class):** Leverages class loading guarantees — thread-safe, lazy, no synchronization overhead.
>
> For most cases, prefer **Enum** or **Bill Pugh**.

---

### Q2. What is the difference between Factory Method and Abstract Factory?

> **Answer:**
> - **Factory Method:** Defines one factory method in an abstract class; subclasses override it to create **one specific product**. Uses inheritance.
> - **Abstract Factory:** Defines an interface with multiple factory methods that produce **a family of related products**. Uses composition — you pass the factory in.
>
> Example: Factory Method → `NotificationService.createNotification()`. Abstract Factory → `UIFactory.createButton()` + `createCheckbox()` for consistent UI families.

---

### Q3. When would you use the Builder pattern over a constructor?

> **Answer:**
> Use Builder when:
> - A class has **many optional parameters** — avoids unreadable telescoping constructors like `new User(name, null, null, true, null, 30)`.
> - Object must be **immutable** but needs complex construction.
> - You want **validation** in one place (the `build()` method).
> - Self-documenting code — `builder.timeout(60).retries(3)` is clearer than positional args.
>
> Java's `StringBuilder`, `Stream.Builder`, and Lombok's `@Builder` all implement this.

---

### Q4. What is the difference between shallow copy and deep copy in Prototype pattern?

> **Answer:**
> - **Shallow copy** (`Object.clone()`): Primitive fields are copied by value. Reference-type fields (List, Map, other objects) share the same reference — modifying the copy's collection also affects the original.
> - **Deep copy**: Every field is recursively copied — collections and nested objects are duplicated. The copy is fully independent of the original.
>
> Use deep copy when the prototype has mutable reference fields that should not be shared. Implement via copy constructor, manual cloning, or serialization/deserialization.

---

### Q5. Can Singleton be broken? How do you prevent it?

> **Answer:**
> Yes, Singleton can be broken by:
> 1. **Reflection:** `Constructor.setAccessible(true)` bypasses private constructor.
> 2. **Serialization/Deserialization:** `readObject()` creates a new instance.
> 3. **Multiple ClassLoaders:** Each loader has its own class, so multiple instances can exist.
>
> Prevention:
> - Use **Enum Singleton** — immune to reflection (throws exception) and serialization (enums handle it natively).
> - For class-based Singleton, add `readResolve()` for serialization and throw exception in constructor if instance already exists.

---

### Q6. What is the difference between Abstract Factory and Dependency Injection?

> **Answer:**
> Both decouple the client from concrete classes, but differently:
> - **Abstract Factory:** The client calls `factory.createProduct()` — it actively requests objects.
> - **Dependency Injection:** The client declares what it needs; a framework (Spring, Guice) injects the dependency — the client is passive.
>
> Abstract Factory is a creational design pattern; DI is an architectural principle (often implemented using IoC containers). Spring's `@Bean` + `@Autowired` is DI; Spring's `FactoryBean` is closer to Abstract Factory.

---

<div align="center">
  <sub>← Back to <a href="Topic.md">All Patterns</a></sub>
</div>
