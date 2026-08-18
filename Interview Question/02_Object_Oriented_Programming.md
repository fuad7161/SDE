# Object-Oriented Programming

1. What is object-oriented programming?
   - OOP models software as objects that combine state and behavior and collaborate through well-defined interfaces.
2. What is the difference between a class and an object?
   - A class is a blueprint defining data and behavior; an object is a concrete instance of that class.
3. What are the four pillars of OOP?
   - They are encapsulation, abstraction, inheritance, and polymorphism.
4. What is encapsulation, and why is it useful?
   - Encapsulation hides internal state behind controlled operations, protecting invariants and reducing dependence on implementation details.
5. What is abstraction? Give a practical example.
   - Abstraction exposes essential behavior while hiding complexity; for example, a `pay()` method can hide payment-provider details.
6. What is inheritance, and when should it be avoided?
   - Inheritance creates a subtype from an existing type; avoid it when the relationship is not truly “is-a” or when it creates tight coupling.
7. What is polymorphism?
   - Polymorphism lets one interface represent multiple concrete types, each supplying its own implementation of an operation.
8. What is the difference between compile-time and runtime polymorphism?
   - Compile-time polymorphism resolves overloaded calls during compilation; runtime polymorphism selects overridden behavior from the object's actual type.
9. What is method overriding?
   - Overriding occurs when a subclass provides its own implementation of an inherited method with a compatible signature.
10. What is the difference between overloading and overriding?
   - Overloading uses the same name with different parameters; overriding replaces inherited behavior in a subtype.
11. What is a constructor? Can a constructor be overloaded?
   - A constructor initializes a new object, and many languages allow multiple constructors with different parameter lists.
12. What is the purpose of the `this` or `self` reference?
   - It refers to the current instance, allowing its fields and methods to be accessed and distinguished from local names.
13. What is the difference between an interface and an abstract class?
   - An interface primarily defines a contract, while an abstract class can also provide shared state and implementation; language rules vary.
14. What is the difference between association, aggregation, and composition?
   - Association is a general relationship, aggregation is weak whole-part ownership, and composition is strong ownership where a part's lifecycle depends on the whole.
15. Why is composition often preferred over inheritance?
   - Composition assembles behavior from replaceable collaborators, usually providing more flexibility and less coupling than a fixed inheritance hierarchy.
16. What are access modifiers?
   - Access modifiers such as public, protected, and private control which code may use a class or member.
17. What are static members, and how do they differ from instance members?
   - Static members belong to the class itself and are shared, while instance members belong to individual objects.
18. What is an immutable class, and how would you design one?
   - Its state cannot change after construction; use private final fields, controlled construction, no mutators, and defensive copies of mutable values.
19. What is dependency injection?
   - Dependency injection supplies an object's collaborators from outside instead of having the object construct them, improving flexibility and testability.
20. What do high cohesion and low coupling mean?
   - High cohesion keeps related responsibilities together, while low coupling minimizes dependencies between components.
21. What are the SOLID principles?
   - SOLID is a set of five design guidelines intended to create understandable, flexible, and maintainable object-oriented code.
22. What is the Single Responsibility Principle?
   - A module should have one focused responsibility, commonly expressed as having one main reason to change.
23. What is the Open/Closed Principle?
   - Software entities should be open to extension but closed to modification, so new behavior can be added without repeatedly changing stable code.
24. What is the Liskov Substitution Principle?
   - A subtype must be safely usable wherever its base type is expected without breaking the caller's valid assumptions.
25. What are the Interface Segregation and Dependency Inversion principles?
   - ISP favors small client-specific interfaces; DIP says high-level policy should depend on abstractions rather than concrete low-level details.
