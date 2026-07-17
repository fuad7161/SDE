# 🟡 Behavioral Patterns

> **Category:** Behavioral &nbsp;|&nbsp; **Tags:** `Observer` `Strategy` `Command` `Template Method` `Chain of Responsibility` `State` `Iterator`

Behavioral patterns deal with **algorithms and communication between objects** — how objects interact and distribute responsibility.

---

## Table of Contents
1. [Observer](#1-observer)
2. [Strategy](#2-strategy)
3. [Command](#3-command)
4. [Template Method](#4-template-method)
5. [Chain of Responsibility](#5-chain-of-responsibility)
6. [State](#6-state)
7. [Iterator](#7-iterator)
8. [Interview Questions](#interview-questions)

---

## 1. Observer

**Intent:** Define a **one-to-many dependency** between objects. When one object (subject) changes state, all its dependents (observers) are notified and updated automatically.

**Analogy:** Subscribing to a YouTube channel — when the creator uploads, all subscribers are notified.

**When to use:**
- An event in one object should trigger updates in multiple others
- Loose coupling between publisher and subscribers
- Event-driven systems, UI frameworks, messaging

<details>
<summary><b>Java</b></summary>

```java
import java.util.*;

// Observer interface
public interface Observer {
    void update(String event, Object data);
}

// Subject interface
public interface Subject {
    void subscribe(String event, Observer observer);
    void unsubscribe(String event, Observer observer);
    void publish(String event, Object data);
}

// Concrete Subject — Event Bus
public class EventBus implements Subject {
    private final Map<String, List<Observer>> listeners = new HashMap<>();

    @Override
    public void subscribe(String event, Observer observer) {
        listeners.computeIfAbsent(event, k -> new ArrayList<>()).add(observer);
    }

    @Override
    public void unsubscribe(String event, Observer observer) {
        List<Observer> obs = listeners.get(event);
        if (obs != null) obs.remove(observer);
    }

    @Override
    public void publish(String event, Object data) {
        List<Observer> obs = listeners.getOrDefault(event, List.of());
        obs.forEach(o -> o.update(event, data));
    }
}

// Concrete Observers
public class EmailAlert implements Observer {
    @Override
    public void update(String event, Object data) {
        System.out.println("EMAIL ALERT [" + event + "]: " + data);
    }
}

public class SMSAlert implements Observer {
    @Override
    public void update(String event, Object data) {
        System.out.println("SMS ALERT [" + event + "]: " + data);
    }
}

public class AuditLog implements Observer {
    @Override
    public void update(String event, Object data) {
        System.out.println("AUDIT LOG [" + event + "]: " + data + " at " + new Date());
    }
}

// Usage
public class Main {
    public static void main(String[] args) {
        EventBus bus = new EventBus();

        bus.subscribe("ORDER_PLACED", new EmailAlert());
        bus.subscribe("ORDER_PLACED", new SMSAlert());
        bus.subscribe("ORDER_PLACED", new AuditLog());
        bus.subscribe("PAYMENT_FAILED", new EmailAlert());

        bus.publish("ORDER_PLACED", "Order #1001");
        // EMAIL ALERT [ORDER_PLACED]: Order #1001
        // SMS ALERT [ORDER_PLACED]: Order #1001
        // AUDIT LOG [ORDER_PLACED]: Order #1001 at ...

        bus.publish("PAYMENT_FAILED", "Order #1002 — card declined");
        // EMAIL ALERT [PAYMENT_FAILED]: Order #1002 — card declined
    }
}
```

</details>

<details>
<summary><b>Go</b></summary>

```go
package main

import (
    "fmt"
    "time"
)

type Observer interface {
    Update(event string, data interface{})
}

type EventBus struct {
    listeners map[string][]Observer
}

func NewEventBus() *EventBus {
    return &EventBus{listeners: make(map[string][]Observer)}
}

func (b *EventBus) Subscribe(event string, observer Observer) {
    b.listeners[event] = append(b.listeners[event], observer)
}

func (b *EventBus) Unsubscribe(event string, observer Observer) {
    list := b.listeners[event]
    for i, o := range list {
        if o == observer {
            b.listeners[event] = append(list[:i], list[i+1:]...)
            return
        }
    }
}

func (b *EventBus) Publish(event string, data interface{}) {
    for _, o := range b.listeners[event] {
        o.Update(event, data)
    }
}

type EmailAlert struct{}
type SMSAlert   struct{}
type AuditLog   struct{}

func (e *EmailAlert) Update(event string, data interface{}) {
    fmt.Printf("EMAIL ALERT [%s]: %v\n", event, data)
}
func (s *SMSAlert) Update(event string, data interface{}) {
    fmt.Printf("SMS ALERT [%s]: %v\n", event, data)
}
func (a *AuditLog) Update(event string, data interface{}) {
    fmt.Printf("AUDIT LOG [%s]: %v at %s\n", event, data, time.Now().Format(time.RFC3339))
}

func main() {
    bus := NewEventBus()
    bus.Subscribe("ORDER_PLACED", &EmailAlert{})
    bus.Subscribe("ORDER_PLACED", &SMSAlert{})
    bus.Subscribe("ORDER_PLACED", &AuditLog{})
    bus.Subscribe("PAYMENT_FAILED", &EmailAlert{})

    bus.Publish("ORDER_PLACED", "Order #1001")
    bus.Publish("PAYMENT_FAILED", "Order #1002 — card declined")
}
```

</details>

### Push vs Pull Model

| | Push | Pull |
|--|------|------|
| **How data is sent** | Subject pushes all data in notification | Observer pulls only what it needs |
| **Coupling** | Subject must know what observers need | Observer queries subject when notified |
| **Example** | `observer.update(event, data)` | `observer.update(subject)` then `subject.getState()` |
| **Best when** | Data is small, always needed | Data is large or conditionally needed |

---

## 2. Strategy

**Intent:** Define a family of algorithms, encapsulate each one, and make them **interchangeable**. Strategy lets the algorithm vary independently from the clients that use it.

**When to use:**
- Multiple variants of an algorithm (sort, payment, compression)
- Replace conditionals (`if/switch`) with polymorphism
- Swap behavior at runtime

<details>
<summary><b>Java</b></summary>

```java
// Strategy interface
public interface SortStrategy {
    void sort(int[] array);
}

// Concrete Strategies
public class BubbleSort implements SortStrategy {
    @Override
    public void sort(int[] array) {
        // simplified
        System.out.println("BubbleSort applied");
        Arrays.sort(array);  // using built-in for brevity
    }
}

public class QuickSort implements SortStrategy {
    @Override
    public void sort(int[] array) {
        System.out.println("QuickSort applied");
        Arrays.sort(array);
    }
}

public class MergeSort implements SortStrategy {
    @Override
    public void sort(int[] array) {
        System.out.println("MergeSort applied");
        Arrays.sort(array);
    }
}

// Context — uses a strategy, can swap at runtime
public class Sorter {
    private SortStrategy strategy;

    public Sorter(SortStrategy strategy) {
        this.strategy = strategy;
    }

    public void setStrategy(SortStrategy strategy) {
        this.strategy = strategy;
    }

    public void sort(int[] array) {
        strategy.sort(array);
    }
}

// Usage
int[] data = {5, 2, 8, 1, 9};
Sorter sorter = new Sorter(new QuickSort());
sorter.sort(data);                    // QuickSort applied

sorter.setStrategy(new MergeSort());
sorter.sort(data);                    // MergeSort applied

// Lambda strategy (Java 8+)
sorter.setStrategy(arr -> {
    System.out.println("Lambda sort");
    Arrays.sort(arr);
});
sorter.sort(data);
```

</details>

<details>
<summary><b>Go</b></summary>

```go
package main

import (
    "fmt"
    "sort"
)

type SortStrategy interface {
    Sort(data []int)
}

type BubbleSort struct{}
type QuickSort  struct{}
type MergeSort  struct{}

func (b *BubbleSort) Sort(data []int) { fmt.Println("BubbleSort applied"); sort.Ints(data) }
func (q *QuickSort) Sort(data []int)  { fmt.Println("QuickSort applied");  sort.Ints(data) }
func (m *MergeSort) Sort(data []int)  { fmt.Println("MergeSort applied");  sort.Ints(data) }

type Sorter struct{ strategy SortStrategy }

func (s *Sorter) SetStrategy(strategy SortStrategy) { s.strategy = strategy }
func (s *Sorter) Sort(data []int)                   { s.strategy.Sort(data) }

// SortFunc lets a plain function implement SortStrategy (like Java lambdas)
type SortFunc func([]int)

func (f SortFunc) Sort(data []int) { f(data) }

func main() {
    data := []int{5, 2, 8, 1, 9}
    sorter := &Sorter{strategy: &QuickSort{}}
    sorter.Sort(data) // QuickSort applied

    sorter.SetStrategy(&MergeSort{})
    sorter.Sort(data) // MergeSort applied

    // Function-based strategy
    sorter.SetStrategy(SortFunc(func(d []int) {
        fmt.Println("Custom sort")
        sort.Ints(d)
    }))
    sorter.Sort(data)
}
```

</details>

### Real-world: Payment Strategy

<details>
<summary><b>Java</b></summary>

```java
public interface PaymentStrategy {
    void pay(double amount);
}

public class CreditCardPayment implements PaymentStrategy {
    private final String cardNumber;
    public CreditCardPayment(String cardNumber) { this.cardNumber = cardNumber; }
    @Override public void pay(double amount) {
        System.out.printf("Paid $%.2f with credit card ending %s%n",
                          amount, cardNumber.substring(cardNumber.length() - 4));
    }
}

public class PayPalPayment implements PaymentStrategy {
    private final String email;
    public PayPalPayment(String email) { this.email = email; }
    @Override public void pay(double amount) {
        System.out.printf("Paid $%.2f via PayPal (%s)%n", amount, email);
    }
}

public class ShoppingCart {
    private PaymentStrategy paymentStrategy;

    public void setPaymentStrategy(PaymentStrategy strategy) {
        this.paymentStrategy = strategy;
    }

    public void checkout(double amount) {
        paymentStrategy.pay(amount);
    }
}
```

</details>

<details>
<summary><b>Go</b></summary>

```go
package main

import "fmt"

type PaymentStrategy interface {
    Pay(amount float64)
}

type CreditCardPayment struct{ cardNumber string }
type PayPalPayment     struct{ email string }

func (c *CreditCardPayment) Pay(amount float64) {
    last4 := c.cardNumber[len(c.cardNumber)-4:]
    fmt.Printf("Paid $%.2f with credit card ending %s\n", amount, last4)
}
func (p *PayPalPayment) Pay(amount float64) {
    fmt.Printf("Paid $%.2f via PayPal (%s)\n", amount, p.email)
}

type ShoppingCart struct{ strategy PaymentStrategy }

func (c *ShoppingCart) SetStrategy(s PaymentStrategy) { c.strategy = s }
func (c *ShoppingCart) Checkout(amount float64)       { c.strategy.Pay(amount) }

func main() {
    cart := &ShoppingCart{}
    cart.SetStrategy(&CreditCardPayment{"4111111111111234"})
    cart.Checkout(99.99) // Paid $99.99 with credit card ending 1234

    cart.SetStrategy(&PayPalPayment{"alice@example.com"})
    cart.Checkout(49.99) // Paid $49.99 via PayPal (alice@example.com)
}
```

</details>

### Strategy vs State

| | Strategy | State |
|--|----------|-------|
| **Who changes behavior** | Client (explicitly swaps strategy) | Object itself (transitions based on state) |
| **Strategies aware of each other** | No | Yes (can trigger transitions) |
| **Use when** | Interchangeable algorithms | Object behavior changes with its internal state |

---

## 3. Command

**Intent:** Encapsulate a request as an object, letting you **parameterize, queue, log, and undo operations**.

**When to use:**
- Undo/redo functionality
- Queueing or scheduling operations
- Transactional behavior (rollback)
- Macro recording

<details>
<summary><b>Java</b></summary>

```java
// Command interface
public interface Command {
    void execute();
    void undo();
}

// Receiver — the object that does the actual work
public class TextEditor {
    private final StringBuilder text = new StringBuilder();

    public void insertText(String text, int pos) {
        this.text.insert(pos, text);
        System.out.println("Text: " + this.text);
    }

    public void deleteText(int pos, int length) {
        this.text.delete(pos, pos + length);
        System.out.println("Text: " + this.text);
    }

    public String getText() { return text.toString(); }
}

// Concrete Command — InsertText
public class InsertCommand implements Command {
    private final TextEditor editor;
    private final String text;
    private final int position;

    public InsertCommand(TextEditor editor, String text, int position) {
        this.editor   = editor;
        this.text     = text;
        this.position = position;
    }

    @Override
    public void execute() { editor.insertText(text, position); }

    @Override
    public void undo() { editor.deleteText(position, text.length()); }
}

// Invoker — stores and executes commands
public class CommandManager {
    private final Deque<Command> history = new ArrayDeque<>();
    private final Deque<Command> redoStack = new ArrayDeque<>();

    public void execute(Command command) {
        command.execute();
        history.push(command);
        redoStack.clear();   // new action clears redo stack
    }

    public void undo() {
        if (!history.isEmpty()) {
            Command cmd = history.pop();
            cmd.undo();
            redoStack.push(cmd);
        }
    }

    public void redo() {
        if (!redoStack.isEmpty()) {
            Command cmd = redoStack.pop();
            cmd.execute();
            history.push(cmd);
        }
    }
}

// Usage
public class Main {
    public static void main(String[] args) {
        TextEditor editor = new TextEditor();
        CommandManager mgr = new CommandManager();

        mgr.execute(new InsertCommand(editor, "Hello", 0));  // Text: Hello
        mgr.execute(new InsertCommand(editor, " World", 5)); // Text: Hello World
        mgr.undo();                                          // Text: Hello
        mgr.redo();                                          // Text: Hello World
        mgr.undo();                                          // Text: Hello
        mgr.undo();                                          // Text:
    }
}
```

</details>

<details>
<summary><b>Go</b></summary>

```go
package main

import "fmt"

type Command interface {
    Execute()
    Undo()
}

// Receiver
type TextEditor struct{ text []rune }

func (e *TextEditor) InsertText(text string, pos int) {
    runes := []rune(text)
    e.text = append(e.text[:pos], append(runes, e.text[pos:]...)...)
    fmt.Println("Text:", string(e.text))
}

func (e *TextEditor) DeleteText(pos, length int) {
    e.text = append(e.text[:pos], e.text[pos+length:]...)
    fmt.Println("Text:", string(e.text))
}

// Concrete Command
type InsertCommand struct {
    editor   *TextEditor
    text     string
    position int
}

func (c *InsertCommand) Execute() { c.editor.InsertText(c.text, c.position) }
func (c *InsertCommand) Undo()    { c.editor.DeleteText(c.position, len([]rune(c.text))) }

// Invoker
type CommandManager struct {
    history   []Command
    redoStack []Command
}

func (m *CommandManager) Execute(cmd Command) {
    cmd.Execute()
    m.history = append(m.history, cmd)
    m.redoStack = nil
}

func (m *CommandManager) Undo() {
    if len(m.history) == 0 {
        return
    }
    cmd := m.history[len(m.history)-1]
    m.history = m.history[:len(m.history)-1]
    cmd.Undo()
    m.redoStack = append(m.redoStack, cmd)
}

func (m *CommandManager) Redo() {
    if len(m.redoStack) == 0 {
        return
    }
    cmd := m.redoStack[len(m.redoStack)-1]
    m.redoStack = m.redoStack[:len(m.redoStack)-1]
    cmd.Execute()
    m.history = append(m.history, cmd)
}

func main() {
    editor := &TextEditor{}
    mgr := &CommandManager{}

    mgr.Execute(&InsertCommand{editor, "Hello", 0})  // Text: Hello
    mgr.Execute(&InsertCommand{editor, " World", 5}) // Text: Hello World
    mgr.Undo()                                        // Text: Hello
    mgr.Redo()                                        // Text: Hello World
    mgr.Undo()                                        // Text: Hello
    mgr.Undo()                                        // Text:
}
```

</details>

---

## 4. Template Method

**Intent:** Define the **skeleton of an algorithm** in a base class, deferring some steps to subclasses. Subclasses override specific steps without changing the algorithm's overall structure.

**When to use:**
- Same algorithm structure, but some steps differ per subclass
- Avoid code duplication across similar classes
- Framework hooks — call user code at defined points

<details>
<summary><b>Java</b></summary>

```java
// Abstract class — defines the template
public abstract class DataProcessor {

    // Template method — fixed algorithm skeleton (final = can't be overridden)
    public final void process() {
        readData();
        processData();
        writeResult();
        sendReport();      // hook — default implementation, can be overridden
    }

    protected abstract void readData();
    protected abstract void processData();
    protected abstract void writeResult();

    // Hook — optional step with a default implementation
    protected void sendReport() {
        System.out.println("Sending default email report");
    }
}

// Concrete class — CSV processing
public class CSVDataProcessor extends DataProcessor {
    @Override
    protected void readData()     { System.out.println("Reading CSV file"); }

    @Override
    protected void processData()  { System.out.println("Parsing CSV rows"); }

    @Override
    protected void writeResult()  { System.out.println("Writing to database"); }

    @Override
    protected void sendReport()   { System.out.println("Sending Slack notification"); }
}

// Concrete class — JSON processing
public class JSONDataProcessor extends DataProcessor {
    @Override
    protected void readData()     { System.out.println("Fetching JSON from API"); }

    @Override
    protected void processData()  { System.out.println("Deserializing JSON"); }

    @Override
    protected void writeResult()  { System.out.println("Saving to S3"); }
    // sendReport() not overridden — uses default
}

// Usage
DataProcessor csv  = new CSVDataProcessor();
csv.process();   // always: readData → processData → writeResult → sendReport

DataProcessor json = new JSONDataProcessor();
json.process();
```

</details>

<details>
<summary><b>Go</b></summary>

```go
package main

import "fmt"

// In Go, Template Method uses an interface + a runner function
type DataProcessor interface {
    ReadData()
    ProcessData()
    WriteResult()
    SendReport() // hook — override or use default
}

// Default hook implementation via embedding
type DefaultReporter struct{}

func (d *DefaultReporter) SendReport() {
    fmt.Println("Sending default email report")
}

// Template runner — defines the fixed sequence
func Process(p DataProcessor) {
    p.ReadData()
    p.ProcessData()
    p.WriteResult()
    p.SendReport()
}

// CSV Processor
type CSVDataProcessor struct{ DefaultReporter }

func (p *CSVDataProcessor) ReadData()    { fmt.Println("Reading CSV file") }
func (p *CSVDataProcessor) ProcessData() { fmt.Println("Parsing CSV rows") }
func (p *CSVDataProcessor) WriteResult() { fmt.Println("Writing to database") }
func (p *CSVDataProcessor) SendReport()  { fmt.Println("Sending Slack notification") }

// JSON Processor
type JSONDataProcessor struct{ DefaultReporter }

func (p *JSONDataProcessor) ReadData()    { fmt.Println("Fetching JSON from API") }
func (p *JSONDataProcessor) ProcessData() { fmt.Println("Deserializing JSON") }
func (p *JSONDataProcessor) WriteResult() { fmt.Println("Saving to S3") }
// SendReport() not overridden — uses DefaultReporter.SendReport

func main() {
    Process(&CSVDataProcessor{})  // ReadData → ProcessData → WriteResult → Slack
    fmt.Println()
    Process(&JSONDataProcessor{}) // ReadData → ProcessData → WriteResult → default email
}
```

</details>

> JUnit's test lifecycle is Template Method: `@BeforeAll`, `@BeforeEach`, `@Test`, `@AfterEach`, `@AfterAll`.

### Template Method vs Strategy

| | Template Method | Strategy |
|--|----------------|---------|
| **Varies by** | Overriding steps in subclass | Delegating to a separate strategy object |
| **Mechanism** | Inheritance | Composition |
| **Algorithm structure** | Fixed in base class | Fully swappable |
| **Best for** | Steps vary, overall flow is same | Entire algorithm swappable |

---

## 5. Chain of Responsibility

**Intent:** Pass a request along a **chain of handlers**. Each handler decides to process the request or pass it to the next handler in the chain.

**Analogy:** HTTP middleware pipeline — request passes through authentication, logging, rate limiting, then reaches the controller.

**When to use:**
- Multiple handlers for a request, decoupled from each other
- The handler is determined at runtime
- Middleware, filters, event propagation

<details>
<summary><b>Java</b></summary>

```java
// Handler interface
public abstract class RequestHandler {
    protected RequestHandler next;

    public RequestHandler setNext(RequestHandler next) {
        this.next = next;
        return next;  // enables fluent chaining
    }

    public abstract void handle(HttpRequest request);

    protected void passToNext(HttpRequest request) {
        if (next != null) next.handle(request);
        else System.out.println("No handler processed the request");
    }
}

// Simulated HttpRequest
public class HttpRequest {
    public final String token;
    public final int    requestsPerMinute;
    public final String body;

    public HttpRequest(String token, int rpm, String body) {
        this.token              = token;
        this.requestsPerMinute  = rpm;
        this.body               = body;
    }
}

// Concrete Handlers
public class AuthenticationHandler extends RequestHandler {
    @Override
    public void handle(HttpRequest request) {
        if (request.token == null || request.token.isBlank()) {
            System.out.println("AUTH: Rejected — missing token");
            return;
        }
        System.out.println("AUTH: Token valid");
        passToNext(request);
    }
}

public class RateLimitHandler extends RequestHandler {
    private static final int MAX_RPM = 100;

    @Override
    public void handle(HttpRequest request) {
        if (request.requestsPerMinute > MAX_RPM) {
            System.out.println("RATE LIMIT: Rejected — " + request.requestsPerMinute + " rpm");
            return;
        }
        System.out.println("RATE LIMIT: OK");
        passToNext(request);
    }
}

public class LoggingHandler extends RequestHandler {
    @Override
    public void handle(HttpRequest request) {
        System.out.println("LOG: Request received, body length = " + request.body.length());
        passToNext(request);
    }
}

public class BusinessLogicHandler extends RequestHandler {
    @Override
    public void handle(HttpRequest request) {
        System.out.println("BUSINESS: Processing — " + request.body);
    }
}

// Build and run the chain
public class Main {
    public static void main(String[] args) {
        RequestHandler auth = new AuthenticationHandler();

        // Build the chain with fluent linking
        auth.setNext(new RateLimitHandler())
            .setNext(new LoggingHandler())
            .setNext(new BusinessLogicHandler());

        System.out.println("--- Valid request ---");
        auth.handle(new HttpRequest("token123", 50, "Create user"));
        // AUTH: Token valid  / RATE LIMIT: OK / LOG: ... / BUSINESS: ...

        System.out.println("\n--- Missing token ---");
        auth.handle(new HttpRequest(null, 50, "Create user"));
        // AUTH: Rejected — missing token

        System.out.println("\n--- Rate limit exceeded ---");
        auth.handle(new HttpRequest("token123", 200, "Create user"));
        // AUTH: Token valid  /  RATE LIMIT: Rejected — 200 rpm
    }
}
```

</details>

<details>
<summary><b>Go</b></summary>

```go
package main

import "fmt"

type HttpRequest struct {
    Token             string
    RequestsPerMinute int
    Body              string
}

type Handler interface {
    Handle(req *HttpRequest)
    SetNext(h Handler) Handler
}

type BaseHandler struct{ next Handler }

func (b *BaseHandler) SetNext(h Handler) Handler {
    b.next = h
    return h // enables fluent chaining
}

func (b *BaseHandler) PassToNext(req *HttpRequest) {
    if b.next != nil {
        b.next.Handle(req)
    } else {
        fmt.Println("No handler processed the request")
    }
}

type AuthHandler      struct{ BaseHandler }
type RateLimitHandler struct{ BaseHandler }
type LoggingHandler   struct{ BaseHandler }
type BusinessHandler  struct{ BaseHandler }

func (h *AuthHandler) Handle(req *HttpRequest) {
    if req.Token == "" {
        fmt.Println("AUTH: Rejected — missing token")
        return
    }
    fmt.Println("AUTH: Token valid")
    h.PassToNext(req)
}

func (h *RateLimitHandler) Handle(req *HttpRequest) {
    if req.RequestsPerMinute > 100 {
        fmt.Printf("RATE LIMIT: Rejected — %d rpm\n", req.RequestsPerMinute)
        return
    }
    fmt.Println("RATE LIMIT: OK")
    h.PassToNext(req)
}

func (h *LoggingHandler) Handle(req *HttpRequest) {
    fmt.Printf("LOG: Request received, body length = %d\n", len(req.Body))
    h.PassToNext(req)
}

func (h *BusinessHandler) Handle(req *HttpRequest) {
    fmt.Println("BUSINESS: Processing —", req.Body)
}

func main() {
    auth := &AuthHandler{}
    auth.SetNext(&RateLimitHandler{}).
        SetNext(&LoggingHandler{}).
        SetNext(&BusinessHandler{})

    fmt.Println("--- Valid request ---")
    auth.Handle(&HttpRequest{"token123", 50, "Create user"})

    fmt.Println("\n--- Missing token ---")
    auth.Handle(&HttpRequest{"", 50, "Create user"})

    fmt.Println("\n--- Rate limit exceeded ---")
    auth.Handle(&HttpRequest{"token123", 200, "Create user"})
}
```

</details>

---

## 6. State

**Intent:** Allow an object to **alter its behavior when its internal state changes**. The object will appear to change its class.

**Analogy:** A vending machine — it behaves differently depending on whether it has items, whether money was inserted, etc.

**When to use:**
- Object behavior is state-dependent
- Complex conditionals based on object's state → replace with state objects
- States and transitions are numerous

<details>
<summary><b>Java</b></summary>

```java
// State interface
public interface OrderState {
    void next(Order order);
    void cancel(Order order);
    String getStatus();
}

// Concrete States
public class PendingState implements OrderState {
    @Override public void next(Order order)   { order.setState(new ProcessingState()); }
    @Override public void cancel(Order order) { order.setState(new CancelledState()); }
    @Override public String getStatus()       { return "PENDING"; }
}

public class ProcessingState implements OrderState {
    @Override public void next(Order order)   { order.setState(new ShippedState()); }
    @Override public void cancel(Order order) { order.setState(new CancelledState()); }
    @Override public String getStatus()       { return "PROCESSING"; }
}

public class ShippedState implements OrderState {
    @Override public void next(Order order)   { order.setState(new DeliveredState()); }
    @Override public void cancel(Order order) {
        throw new IllegalStateException("Cannot cancel a shipped order");
    }
    @Override public String getStatus()       { return "SHIPPED"; }
}

public class DeliveredState implements OrderState {
    @Override public void next(Order order)   {
        throw new IllegalStateException("Order already delivered");
    }
    @Override public void cancel(Order order) {
        throw new IllegalStateException("Cannot cancel a delivered order");
    }
    @Override public String getStatus()       { return "DELIVERED"; }
}

public class CancelledState implements OrderState {
    @Override public void next(Order order)   {
        throw new IllegalStateException("Order is cancelled");
    }
    @Override public void cancel(Order order) {
        System.out.println("Order already cancelled");
    }
    @Override public String getStatus()       { return "CANCELLED"; }
}

// Context
public class Order {
    private OrderState state = new PendingState();
    private final String id;

    public Order(String id) { this.id = id; }

    public void setState(OrderState state) { this.state = state; }

    public void advance() {
        state.next(this);
        System.out.println("Order " + id + " → " + state.getStatus());
    }

    public void cancel() {
        state.cancel(this);
        System.out.println("Order " + id + " → " + state.getStatus());
    }
}

// Usage
Order order = new Order("ORD-001");
order.advance();   // Order ORD-001 → PROCESSING
order.advance();   // Order ORD-001 → SHIPPED
order.advance();   // Order ORD-001 → DELIVERED
// order.cancel(); // IllegalStateException: Cannot cancel a delivered order
```

</details>

<details>
<summary><b>Go</b></summary>

```go
package main

import "fmt"

type OrderState interface {
    Next(order *Order)
    Cancel(order *Order)
    GetStatus() string
}

type Order struct {
    id    string
    state OrderState
}

func NewOrder(id string) *Order       { return &Order{id: id, state: &PendingState{}} }
func (o *Order) SetState(s OrderState) { o.state = s }
func (o *Order) Advance() {
    o.state.Next(o)
    fmt.Printf("Order %s → %s\n", o.id, o.state.GetStatus())
}
func (o *Order) Cancel() {
    o.state.Cancel(o)
    fmt.Printf("Order %s → %s\n", o.id, o.state.GetStatus())
}

type PendingState    struct{}
type ProcessingState struct{}
type ShippedState    struct{}
type DeliveredState  struct{}
type CancelledState  struct{}

func (s *PendingState) Next(o *Order)    { o.SetState(&ProcessingState{}) }
func (s *PendingState) Cancel(o *Order)  { o.SetState(&CancelledState{}) }
func (s *PendingState) GetStatus() string { return "PENDING" }

func (s *ProcessingState) Next(o *Order)    { o.SetState(&ShippedState{}) }
func (s *ProcessingState) Cancel(o *Order)  { o.SetState(&CancelledState{}) }
func (s *ProcessingState) GetStatus() string { return "PROCESSING" }

func (s *ShippedState) Next(o *Order)    { o.SetState(&DeliveredState{}) }
func (s *ShippedState) Cancel(*Order)    { panic("cannot cancel a shipped order") }
func (s *ShippedState) GetStatus() string { return "SHIPPED" }

func (s *DeliveredState) Next(*Order)    { panic("order already delivered") }
func (s *DeliveredState) Cancel(*Order)  { panic("cannot cancel a delivered order") }
func (s *DeliveredState) GetStatus() string { return "DELIVERED" }

func (s *CancelledState) Next(*Order)    { panic("order is cancelled") }
func (s *CancelledState) Cancel(*Order)  { fmt.Println("Order already cancelled") }
func (s *CancelledState) GetStatus() string { return "CANCELLED" }

func main() {
    order := NewOrder("ORD-001")
    order.Advance() // Order ORD-001 → PROCESSING
    order.Advance() // Order ORD-001 → SHIPPED
    order.Advance() // Order ORD-001 → DELIVERED
}
```

</details>

---

## 7. Iterator

**Intent:** Provide a way to **sequentially access elements** of a collection without exposing its underlying structure.

**When to use:**
- Traverse different collection types with a uniform interface
- Hide internal structure (array, tree, graph, linked list)
- Multiple simultaneous traversals

<details>
<summary><b>Java</b></summary>

```java
// Custom collection — a binary search tree that can be iterated in-order
public class BinarySearchTree implements Iterable<Integer> {

    private Node root;

    private static class Node {
        int value;
        Node left, right;
        Node(int value) { this.value = value; }
    }

    public void insert(int value) {
        root = insert(root, value);
    }

    private Node insert(Node node, int value) {
        if (node == null) return new Node(value);
        if (value < node.value) node.left  = insert(node.left, value);
        else if (value > node.value) node.right = insert(node.right, value);
        return node;
    }

    // Return an in-order iterator
    @Override
    public Iterator<Integer> iterator() {
        return new InOrderIterator(root);
    }

    // Concrete Iterator — in-order traversal using a stack
    private static class InOrderIterator implements Iterator<Integer> {
        private final Deque<Node> stack = new ArrayDeque<>();

        InOrderIterator(Node root) {
            pushLeft(root);
        }

        private void pushLeft(Node node) {
            while (node != null) {
                stack.push(node);
                node = node.left;
            }
        }

        @Override
        public boolean hasNext() { return !stack.isEmpty(); }

        @Override
        public Integer next() {
            if (!hasNext()) throw new NoSuchElementException();
            Node node = stack.pop();
            pushLeft(node.right);
            return node.value;
        }
    }
}

// Usage
BinarySearchTree bst = new BinarySearchTree();
bst.insert(5); bst.insert(3); bst.insert(7);
bst.insert(1); bst.insert(4); bst.insert(6);

for (int val : bst) {
    System.out.print(val + " ");   // 1 3 4 5 6 7
}
```

</details>

<details>
<summary><b>Go</b></summary>

```go
package main

import "fmt"

type BSTNode struct {
    value       int
    left, right *BSTNode
}

type BinarySearchTree struct{ root *BSTNode }

func (t *BinarySearchTree) Insert(value int) {
    t.root = insertNode(t.root, value)
}

func insertNode(node *BSTNode, value int) *BSTNode {
    if node == nil {
        return &BSTNode{value: value}
    }
    if value < node.value {
        node.left = insertNode(node.left, value)
    } else if value > node.value {
        node.right = insertNode(node.right, value)
    }
    return node
}

// Iterator — in-order traversal using a stack
type InOrderIterator struct{ stack []*BSTNode }

func NewInOrderIterator(root *BSTNode) *InOrderIterator {
    it := &InOrderIterator{}
    it.pushLeft(root)
    return it
}

func (it *InOrderIterator) pushLeft(node *BSTNode) {
    for node != nil {
        it.stack = append(it.stack, node)
        node = node.left
    }
}

func (it *InOrderIterator) HasNext() bool { return len(it.stack) > 0 }

func (it *InOrderIterator) Next() int {
    node := it.stack[len(it.stack)-1]
    it.stack = it.stack[:len(it.stack)-1]
    it.pushLeft(node.right)
    return node.value
}

func main() {
    bst := &BinarySearchTree{}
    for _, v := range []int{5, 3, 7, 1, 4, 6} {
        bst.Insert(v)
    }

    it := NewInOrderIterator(bst.root)
    for it.HasNext() {
        fmt.Print(it.Next(), " ") // 1 3 4 5 6 7
    }
    fmt.Println()
}
```

</details>

> Java's `java.util.Iterator` and `Iterable` interfaces formalize this pattern. All `Collection` implementations use it.

---

## Interview Questions

### Q1. Explain the Observer pattern. What is the difference between push and pull models?

> **Answer:**
> Observer defines a one-to-many relationship: when the subject changes state, all registered observers are notified.
>
> - **Push:** Subject sends all relevant data in the notification: `observer.update(event, data)`. Simple but may send unnecessary data.
> - **Pull:** Subject sends only a reference to itself: `observer.update(subject)`. Observer then calls `subject.getState()` to fetch only what it needs. More flexible but adds coupling back to subject.
>
> Real-world: Java's `PropertyChangeListener`, Android's `LiveData`, `EventBus`, Spring's `ApplicationEvent`.

---

### Q2. What is the Strategy pattern? How does it replace if/else chains?

> **Answer:**
> Strategy encapsulates interchangeable algorithms behind a common interface. Instead of:
> ```java
> if (type.equals("CREDIT")) { // credit payment }
> else if (type.equals("PAYPAL")) { // paypal payment }
> ```
> You inject a `PaymentStrategy` and call `strategy.pay(amount)`. Adding a new payment type means creating a new class — no modification to existing code (Open/Closed Principle).
>
> Used in: Java's `Comparator` (passed to `Collections.sort()`), `java.util.Comparator`, Spring's `ResourceLoader`, servlet filters.

---

### Q3. How does the Command pattern enable Undo/Redo?

> **Answer:**
> Each `Command` object implements both `execute()` and `undo()`. The Invoker maintains a history stack:
> - `execute()`: run the command, push to history, clear redo stack.
> - `undo()`: pop from history, call `undo()`, push to redo stack.
> - `redo()`: pop from redo stack, call `execute()`, push to history.
>
> The trick is each command stores enough state to reverse itself (e.g., `InsertCommand` stores the text and position, so `undo()` can delete exactly that text).

---

### Q4. What is the difference between Strategy and State patterns?

> **Answer:**
> Both encapsulate behavior and use composition, but differ in intent:
> - **Strategy:** Algorithms are interchangeable. The **client** actively swaps strategies. Strategies are independent and don't know about each other.
> - **State:** The object transitions between states based on its own logic. **States know about each other** (transitions). The client calls the same methods; behavior changes automatically based on state.
>
> Strategy = "which algorithm to use" (client choice). State = "what this object can do right now" (auto-managed).

---

### Q5. What is the Chain of Responsibility pattern? Give a real-world example.

> **Answer:**
> Chain of Responsibility passes a request through a chain of handlers. Each handler either processes the request or forwards it to the next.
>
> Real-world examples:
> - **Java Servlet Filters:** `Filter.doFilter()` chains security → logging → compression → controller.
> - **Spring Security filter chain:** Authentication → authorization → CSRF → CORS filters.
> - **Exception handling:** `catch` blocks form a chain — most specific first, then more general.
> - **Approval workflows:** Manager → Director → VP based on expense amount.
>
> Key advantage: new handlers can be added or reordered without changing the client or other handlers.

---

### Q6. When would you use Template Method vs Strategy?

> **Answer:**
> - **Template Method:** When the **overall algorithm is fixed**, but specific steps differ. Use when the variation is in steps, not the entire algorithm. Relies on inheritance — subclass overrides specific steps.
> - **Strategy:** When the **entire algorithm can be swapped**. Use when different clients need completely different algorithms. Relies on composition — algorithm is injected.
>
> If you're subclassing just to override a few steps, Template Method is natural. If you want runtime swappability without inheritance, use Strategy.

---

### Q7. How does Iterator improve over direct collection access?

> **Answer:**
> - **Encapsulation:** Client doesn't know if the collection is an array, linked list, tree, or graph — it just calls `next()`.
> - **Single Responsibility:** Traversal logic lives in the iterator, not the collection.
> - **Multiple traversals:** Each call to `iterator()` creates an independent cursor — multiple traversals can run simultaneously.
> - **Lazy evaluation:** Iterators can generate elements on demand (e.g., infinite sequences, streams).
>
> Java's `Iterator` + `Iterable` + `for-each` loop is the language-level implementation. Java `Stream` is a functional evolution of this pattern.

---

<div align="center">
  <sub>← Back to <a href="Topic.md">All Patterns</a></sub>
</div>
