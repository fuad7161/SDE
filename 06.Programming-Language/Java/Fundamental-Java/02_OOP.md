# Object-Oriented Programming (OOP)

---

## Table of Contents

1. [Four Pillars of OOP](#1-four-pillars-of-oop)
2. [Class & Object](#2-class--object)
3. [Constructors](#3-constructors)
4. [`this` Keyword](#4-this-keyword)
5. [`static` Keyword](#5-static-keyword)
6. [Inheritance](#6-inheritance)
7. [Polymorphism](#7-polymorphism)
8. [Encapsulation](#8-encapsulation)
9. [Abstraction — Interface vs Abstract Class](#9-abstraction--interface-vs-abstract-class)
10. [`super` Keyword](#10-super-keyword)
11. [Association, Aggregation, Composition](#11-association-aggregation-composition)

---

## 1. Four Pillars of OOP

```
╔══════════════╗  ╔══════════════╗  ╔══════════════╗  ╔══════════════╗
║ Encapsulation║  ║  Abstraction ║  ║  Inheritance ║  ║Polymorphism  ║
║  hide data   ║  ║ hide details ║  ║  reuse code  ║  ║ many forms   ║
╚══════════════╝  ╚══════════════╝  ╚══════════════╝  ╚══════════════╝
```

| Pillar | Definition | Achieved via |
|---|---|---|
| **Encapsulation** | Bundle data + behavior; restrict direct access | `private` fields + getters/setters |
| **Abstraction** | Expose what, hide how | Interfaces, abstract classes |
| **Inheritance** | Child reuses parent fields and methods | `extends`, `implements` |
| **Polymorphism** | Same method name, different behavior | Overloading (compile-time), Overriding (runtime) |

> **Interview Q: What are the four pillars of OOP? Give an example of each.**  
> See sections below for full examples.

---

## 2. Class & Object

```java
// Class — blueprint/template
class Car {
    // Fields (state)
    String brand;
    int speed;

    // Methods (behavior)
    void accelerate(int amount) {
        speed += amount;
        System.out.println(brand + " speed: " + speed);
    }

    void brake() {
        speed = 0;
        System.out.println(brand + " stopped");
    }
}

// Object — instance of the class (created on heap)
Car tesla = new Car();     // 'new' allocates memory on heap
tesla.brand = "Tesla";
tesla.speed = 0;
tesla.accelerate(60);      // Tesla speed: 60

Car bmw = new Car();       // second independent object
bmw.brand = "BMW";
bmw.accelerate(80);        // BMW speed: 80

// Each object has its own copy of instance fields
// but shares the same method code
```

> **Interview Q: What is the difference between a class and an object?**  
> A **class** is a blueprint — it defines structure (fields) and behavior (methods) but occupies no memory at runtime on its own. An **object** is a concrete **instance** of a class — it exists in memory (heap) and has its own copy of instance variables. You can create many objects from one class, each with independent state.

---

## 3. Constructors

```java
class Student {
    String name;
    int age;
    String course;

    // Default constructor — compiler adds this only if you define NO constructor
    // Student() { }

    // Parameterized constructor
    Student(String name, int age) {
        this.name = name;
        this.age = age;
        this.course = "General";      // default value for course
    }

    // Overloaded constructor
    Student(String name, int age, String course) {
        this(name, age);              // constructor chaining — calls above constructor
        this.course = course;         // only sets the extra field
    }

    // Copy constructor — creates a new object from an existing one
    Student(Student other) {
        this.name = other.name;
        this.age = other.age;
        this.course = other.course;
    }
}

Student s1 = new Student("Alice", 20);
Student s2 = new Student("Bob", 22, "CS");
Student s3 = new Student(s1);         // copy of s1
```

**Key rules:**
- Constructor name must **match the class name exactly**
- No return type (not even `void`)
- `this()` call must be the **first statement**
- If you define any constructor, the compiler will **NOT** add a default one

> **Interview Q: What is constructor chaining? Why use it?**  
> Constructor chaining calls one constructor from another using `this()` (same class) or `super()` (parent class). It avoids **code duplication** — common initialization logic is written once in one constructor, and other constructors delegate to it. `this()` must always be the very first statement in a constructor.

---

## 4. `this` Keyword

```java
class Employee {
    String name;
    int salary;

    // Use 1: Resolve field vs parameter name conflict
    Employee(String name, int salary) {
        this.name = name;       // 'this.name' = field, 'name' = parameter
        this.salary = salary;
    }

    // Use 2: Call another constructor in same class
    Employee(String name) {
        this(name, 50000);      // calls Employee(String, int)
    }

    // Use 3: Pass current object as argument
    void register(EmployeeRegistry registry) {
        registry.add(this);     // pass the current Employee object
    }

    // Use 4: Return current object (method chaining / Builder pattern)
    Employee withSalary(int salary) {
        this.salary = salary;
        return this;            // enables chaining: emp.withSalary(60000).register(reg)
    }
}
```

> **Interview Q: What are the uses of the `this` keyword?**  
> 1. Distinguish **instance variable from parameter** with the same name (`this.name = name`)  
> 2. Call **another constructor** in the same class (`this(args)`) — must be first statement  
> 3. Pass the **current object** as an argument to a method  
> 4. Return the **current object** from a method (enables fluent/builder-style chaining)

---

## 5. `static` Keyword

```java
class Counter {
    // Static field — shared across ALL instances, not per object
    static int totalCount = 0;
    int id;

    // Static block — runs once when class is first loaded
    static {
        System.out.println("Counter class loaded");
        totalCount = 0;
    }

    Counter() {
        totalCount++;
        this.id = totalCount;
    }

    // Static method — can call without creating an object
    static int getTotal() {
        return totalCount;
        // Cannot access 'this' or instance fields/methods here
    }
}

Counter c1 = new Counter();
Counter c2 = new Counter();
System.out.println(Counter.getTotal());   // 2 — accessed via class name
System.out.println(Counter.totalCount);   // 2

// Static inner class — does not hold reference to outer class instance
class Outer {
    static class StaticNested {
        void show() { System.out.println("Static nested class"); }
    }
}
Outer.StaticNested obj = new Outer.StaticNested();  // no Outer instance needed
```

> **Interview Q: Can a `static` method access instance (non-static) variables?**  
> No. A static method belongs to the **class**, not to any object. Instance variables only exist when an object is created. Since static methods can be called without creating any object, there's no `this` reference and no instance variables to access. You'd get a compile error: "non-static variable cannot be referenced from a static context."

---

## 6. Inheritance

```java
// Parent class (superclass)
class Animal {
    String name;
    int age;

    Animal(String name, int age) {
        this.name = name;
        this.age = age;
    }

    void eat() {
        System.out.println(name + " is eating");
    }

    void sleep() {
        System.out.println(name + " is sleeping");
    }
}

// Child class (subclass) inherits Animal's fields and methods
class Dog extends Animal {
    String breed;

    Dog(String name, int age, String breed) {
        super(name, age);          // call parent constructor — must be FIRST
        this.breed = breed;
    }

    // Override parent method
    @Override
    void eat() {
        System.out.println(name + " (a " + breed + ") is eating dog food");
    }

    // New method specific to Dog
    void fetch() {
        System.out.println(name + " is fetching the ball!");
    }
}

Dog dog = new Dog("Rex", 3, "Labrador");
dog.eat();    // uses Dog's overridden version
dog.sleep();  // inherited from Animal — not redefined
dog.fetch();  // Dog-specific method

// Upcasting — Dog IS-A Animal
Animal a = new Dog("Buddy", 2, "Poodle");
a.eat();      // still calls Dog's eat() — runtime polymorphism
// a.fetch(); // ❌ compile error — Animal reference can't see Dog-specific methods
```

**Types of inheritance in Java:**

| Type | Supported? | Example |
|---|---|---|
| Single | ✅ | `Dog extends Animal` |
| Multi-level | ✅ | `GuideDog extends Dog extends Animal` |
| Hierarchical | ✅ | `Dog extends Animal`, `Cat extends Animal` |
| Multiple (class) | ❌ | `class C extends A, B` — not allowed (diamond problem) |
| Multiple (interface) | ✅ | `class C implements A, B` — allowed |

> **Interview Q: Why doesn't Java support multiple inheritance of classes?**  
> To avoid the **Diamond Problem** — if class `C` extends both `A` and `B`, and both have a method `display()`, which one does `C` inherit? Java avoids this ambiguity by restricting to single class inheritance. Multiple **interface** inheritance is allowed because interfaces only define contracts (not state), and conflicts are resolved explicitly.

---

## 7. Polymorphism

### Compile-time Polymorphism (Method Overloading)

```java
class Printer {
    // Same name, different parameter lists — resolved at compile time
    void print(int n) {
        System.out.println("Printing int: " + n);
    }

    void print(String s) {
        System.out.println("Printing String: " + s);
    }

    void print(int n, int m) {
        System.out.println("Printing two ints: " + n + ", " + m);
    }

    // ❌ NOT overloading — return type alone cannot differentiate
    // double print(int n) { return n; }  // compile error: duplicate method
}

Printer p = new Printer();
p.print(42);           // calls print(int)
p.print("hello");      // calls print(String)
p.print(1, 2);         // calls print(int, int)
```

### Runtime Polymorphism (Method Overriding)

```java
class Shape {
    double area() { return 0; }
    void describe() { System.out.println("I am a shape, area = " + area()); }
}

class Circle extends Shape {
    double radius;
    Circle(double r) { this.radius = r; }

    @Override
    double area() { return Math.PI * radius * radius; }
}

class Rectangle extends Shape {
    double w, h;
    Rectangle(double w, double h) { this.w = w; this.h = h; }

    @Override
    double area() { return w * h; }
}

// Polymorphic array — one reference type, many object types
Shape[] shapes = { new Circle(5), new Rectangle(4, 6), new Circle(3) };

for (Shape s : shapes) {
    s.describe();    // correct area() called based on actual runtime type
}
// I am a shape, area = 78.539...
// I am a shape, area = 24.0
// I am a shape, area = 28.274...
```

> **Interview Q: What is the difference between method overloading and overriding?**
>
> | | Overloading | Overriding |
> |---|---|---|
> | Where | Same class | Subclass |
> | Signature | Must differ (params) | Must be identical |
> | Return type | Can differ | Same (or covariant) |
> | Access modifier | Any | Cannot restrict |
> | `static`/`private`/`final` | Can overload | Cannot override |
> | Resolution | Compile-time | Runtime (dynamic dispatch) |

---

## 8. Encapsulation

```java
class BankAccount {
    // Private fields — hidden from outside world
    private String accountNumber;
    private double balance;
    private String owner;

    public BankAccount(String accountNumber, String owner, double initialBalance) {
        this.accountNumber = accountNumber;
        this.owner = owner;
        if (initialBalance < 0) throw new IllegalArgumentException("Initial balance cannot be negative");
        this.balance = initialBalance;
    }

    // Controlled read access
    public double getBalance() { return balance; }
    public String getOwner() { return owner; }
    public String getAccountNumber() { return accountNumber; }

    // Controlled write — business rules enforced in setters
    public void deposit(double amount) {
        if (amount <= 0) throw new IllegalArgumentException("Deposit must be positive");
        balance += amount;
    }

    public void withdraw(double amount) {
        if (amount <= 0) throw new IllegalArgumentException("Withdrawal must be positive");
        if (amount > balance) throw new IllegalStateException("Insufficient funds");
        balance -= amount;
    }

    // No setter for accountNumber or balance — direct modification not allowed
}

BankAccount acc = new BankAccount("ACC001", "Alice", 1000);
acc.deposit(500);
acc.withdraw(200);
System.out.println(acc.getBalance());   // 1300.0
// acc.balance = 1000000;  // ❌ compile error — private field
```

> **Interview Q: What is encapsulation and why is it important?**  
> Encapsulation means binding data and methods together and **restricting direct access** to the data using access modifiers. It is important because:
> 1. **Data protection** — no one can set `balance = -999` directly
> 2. **Validation** — you can enforce rules in setters/methods
> 3. **Flexibility** — internal implementation can change without affecting external code
> 4. **Maintainability** — changes are localized

---

## 9. Abstraction — Interface vs Abstract Class

```java
// Abstract Class — partial abstraction (can have state + concrete methods)
abstract class Vehicle {
    String brand;
    int year;

    Vehicle(String brand, int year) {
        this.brand = brand;
        this.year = year;
    }

    abstract void start();     // subclass MUST implement
    abstract void stop();

    // Concrete method — shared implementation
    void displayInfo() {
        System.out.println(brand + " (" + year + ")");
    }
}

// Interface — pure abstraction (contract, no state)
interface Electric {
    int MAX_CHARGE = 100;           // implicitly public static final
    int getChargeLevel();           // implicitly public abstract
    default void plugIn() {         // concrete, but can be overridden
        System.out.println("Plugging in to charge...");
    }
}

interface GPS {
    String getLocation();
}

// A class can extend ONE abstract class and implement MANY interfaces
class Tesla extends Vehicle implements Electric, GPS {
    private int charge;

    Tesla(int charge) {
        super("Tesla", 2024);
        this.charge = charge;
    }

    @Override public void start() { System.out.println("Tesla starting silently"); }
    @Override public void stop()  { System.out.println("Tesla stopped"); }
    @Override public int getChargeLevel() { return charge; }
    @Override public String getLocation() { return "37.7749° N, 122.4194° W"; }
}
```

| | Abstract Class | Interface |
|---|---|---|
| State (fields) | ✅ (any type) | ❌ (only `public static final`) |
| Constructor | ✅ | ❌ |
| Concrete methods | ✅ | Only `default`/`static` (Java 8+) |
| Multiple inheritance | ❌ (single extend) | ✅ (multiple implement) |
| Access modifiers | Any | `public` by default |
| **When to use** | Shared base with code/state | Define a contract |

> **Interview Q: When would you use an abstract class instead of an interface?**  
> Use an **abstract class** when: (1) you want to share **code/state** among related classes (e.g., all vehicles share `brand` and `displayInfo()`), (2) you need a **constructor**, or (3) you want non-public methods. Use an **interface** when: (1) defining a **capability** that unrelated classes can share (e.g., `Serializable`, `Comparable`), (2) you want **multiple inheritance**, or (3) you're defining a pure contract with no state.

---

## 10. `super` Keyword

```java
class Person {
    String name;
    int age;

    Person(String name, int age) {
        this.name = name;
        this.age = age;
    }

    void introduce() {
        System.out.println("I'm " + name + ", age " + age);
    }
}

class Employee extends Person {
    String company;
    double salary;

    Employee(String name, int age, String company, double salary) {
        super(name, age);          // Use 1: call parent constructor — MUST be first
        this.company = company;
        this.salary = salary;
    }

    @Override
    void introduce() {
        super.introduce();         // Use 2: call parent's method
        System.out.println("I work at " + company + " earning $" + salary);
    }

    void showParentName() {
        System.out.println(super.name);   // Use 3: access parent's field (if not overridden)
    }
}

Employee emp = new Employee("Alice", 30, "Google", 120000);
emp.introduce();
// I'm Alice, age 30
// I work at Google earning $120000.0
```

> **Interview Q: What is the difference between `this()` and `super()`?**  
> `this()` calls another constructor **in the same class** — used for constructor chaining. `super()` calls the **parent class constructor** — used to initialize the inherited part. Both must be the **first statement** in a constructor, so they **cannot both appear** in the same constructor. If neither is written, Java inserts `super()` (no-arg parent constructor) automatically.

---

## 11. Association, Aggregation, Composition

```java
// ─── ASSOCIATION — "uses-a" — no ownership, both live independently ───
class Teacher {
    String name;
    Teacher(String name) { this.name = name; }
    void teach(Student s) {
        System.out.println(name + " is teaching " + s.name);
    }
}
class Student {
    String name;
    Student(String name) { this.name = name; }
}
// Teacher uses Student but doesn't own it

// ─── AGGREGATION — "has-a" — weak ownership, child exists without parent ───
class Department {
    String name;
    List<Employee> employees;     // Department has Employees

    Department(String name, List<Employee> employees) {
        this.name = name;
        this.employees = employees;  // employees passed in — they exist outside
    }
}
class Employee {
    String name;
    Employee(String name) { this.name = name; }
}
// If Department is deleted, employees still exist

// ─── COMPOSITION — "part-of" — strong ownership, part can't exist alone ───
class House {
    private final List<Room> rooms;  // House OWNS rooms

    House(int numberOfRooms) {
        rooms = new ArrayList<>();
        for (int i = 1; i <= numberOfRooms; i++) {
            rooms.add(new Room("Room " + i));  // Rooms created BY House
        }
    }
    // If House is destroyed, Rooms are destroyed too
}
class Room {
    String name;
    Room(String name) { this.name = name; }
}
```

| Relationship | Symbol | Ownership | Lifecycle dependency | Example |
|---|---|---|---|---|
| **Association** | uses-a | None | Independent | Teacher ↔ Student |
| **Aggregation** | has-a (weak) | Partial | Child can exist alone | Department → Employee |
| **Composition** | part-of (strong) | Full | Child dies with parent | House → Room |

> **Interview Q: What is the difference between aggregation and composition?**  
> Both are "has-a" relationships, but they differ in **ownership and lifecycle**. In **aggregation**, the child can exist independently of the parent (e.g., employees exist even if the department is deleted). In **composition**, the child's lifecycle is **tied to the parent** — if the parent is destroyed, so are the children (e.g., a `Room` only makes sense inside a `House`; it's created and destroyed with it). Composition represents a **stronger** dependency.
