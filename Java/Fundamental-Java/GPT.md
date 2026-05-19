# Java Interview Preparation — Study Guide

A structured roadmap covering all core Java topics commonly asked in technical interviews.

---

## Table of Contents

1. [Java Basics](#1-java-basics)
2. [Object-Oriented Programming (OOP)](#2-object-oriented-programming-oop)
3. [String Handling](#3-string-handling)
4. [Exception Handling](#4-exception-handling)
5. [Collections Framework](#5-collections-framework)
6. [Java Memory Management](#6-java-memory-management)
7. [Multithreading & Concurrency](#7-multithreading--concurrency)
8. [Java 8 Features](#8-java-8-features-very-important)
9. [Object Class Methods](#9-object-class-methods)
10. [Keywords & Modifiers](#10-keywords--modifiers)
11. [File Handling & I/O](#11-file-handling--io)
12. [Important Interview Comparisons](#12-important-interview-comparisons)
13. [Important Java Concepts](#13-important-java-concepts)

---

## 1. Java Basics

- Java features — platform independent, OOP, robust, multithreaded
- JDK vs JRE vs JVM
- Java program structure
- Data types — primitive & non-primitive
- Variables and scope
- Operators
- Type casting
- Control statements — `if`, `switch`, loops
- Arrays
- Command-line arguments

---

## 2. Object-Oriented Programming (OOP)

- Class & Object
- Constructor — default, parameterized, constructor chaining
- `this` keyword
- `static` keyword
- Inheritance
- Polymorphism — compile-time & runtime
- Method overloading vs overriding
- Encapsulation
- Abstraction
- Interface vs Abstract Class
- `super` keyword
- Association, Aggregation, Composition

---

## 3. String Handling

- String immutability
- String pool
- `String` vs `StringBuilder` vs `StringBuffer`
- `equals()` vs `==`
- Common String methods
- Mutable vs immutable objects

---

## 4. Exception Handling

- `try`, `catch`, `finally`
- `throw` vs `throws`
- Checked vs unchecked exceptions
- Custom exceptions
- Exception hierarchy
- Multiple catch blocks

---

## 5. Collections Framework

- `List`, `Set`, `Map`, `Queue`
- Key comparisons:
  - `ArrayList` vs `LinkedList`
  - `HashMap` vs `Hashtable`
  - `HashMap` vs `ConcurrentHashMap`
  - `HashSet` vs `TreeSet`
- Internal working of `HashMap`
- `equals()` and `hashCode()`
- `Comparable` vs `Comparator`
- `Iterator` vs `ListIterator`
- Fail-fast vs Fail-safe iterator

---

## 6. Java Memory Management

- Heap vs Stack memory
- Garbage Collection
- Memory leaks
- Object lifecycle
- JVM architecture
- Class loader

---

## 7. Multithreading & Concurrency

- Thread lifecycle
- Creating threads — `Thread` class vs `Runnable`
- Synchronization
- `volatile` keyword
- Deadlock
- Race condition
- Inter-thread communication
- Executor Framework
- `Callable` & `Future`
- Concurrent collections

---

## 8. Java 8 Features *(Very Important)*

- Lambda expressions
- Functional interfaces
- Stream API
- Method references
- `Optional` class
- Default & static methods in interfaces
- Date & Time API (`java.time`)

---

## 9. Object Class Methods

Key methods inherited by every Java class from `java.lang.Object`:

| Method | Purpose |
|---|---|
| `toString()` | String representation of the object |
| `equals()` | Logical equality comparison |
| `hashCode()` | Hash code for use in hash-based collections |
| `clone()` | Creates a shallow copy |
| `finalize()` | Called by GC before object is collected (deprecated) |
| `wait()` | Causes thread to wait until notified |
| `notify()` | Wakes up a single waiting thread |
| `notifyAll()` | Wakes up all waiting threads |

---

## 10. Keywords & Modifiers

| Keyword | Description |
|---|---|
| `final` | Prevents modification of variables, methods, or classes |
| `finally` | Block that always executes after try-catch |
| `finalize()` | Pre-GC cleanup method (deprecated since Java 9) |
| `static` | Belongs to the class, not an instance |
| `transient` | Excludes field from serialization |
| `volatile` | Ensures visibility of changes across threads |
| `synchronized` | Restricts access to one thread at a time |
| `native` | Method implemented in native (non-Java) code |
| `abstract` | Declares a method or class without implementation |

---

## 11. File Handling & I/O

- Byte stream vs Character stream
- `File` class
- Serialization & Deserialization
- `BufferedReader` / `BufferedWriter`
- `Scanner` class

---

## 12. Important Interview Comparisons

| Topic | Compare |
|---|---|
| Equality | `==` vs `equals()` |
| Abstraction | Abstract class vs Interface |
| Lists | `ArrayList` vs `LinkedList` |
| Maps | `HashMap` vs `Hashtable` |
| Strings | `StringBuffer` vs `StringBuilder` |
| Polymorphism | Overloading vs Overriding |
| Concurrency | Process vs Thread |
| Error handling | Exception vs Error |

---

## 13. Important Java Concepts

- Immutable class
- Singleton class
- Marker interface
- Wrapper classes
- Autoboxing & Unboxing
- Pass by value in Java
- Shallow copy vs Deep copy
- Reflection API
