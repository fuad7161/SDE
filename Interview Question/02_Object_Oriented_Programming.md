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
   - Its observable state cannot change after construction; prevent uncontrolled mutation, initialize all state during construction, and defensively copy mutable values where needed.
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

## Medium to Advanced

26. What is the difference between subtype polymorphism and parametric polymorphism?
   - **Key note:** Subtype polymorphism dispatches through inheritance/interfaces; parametric polymorphism uses generics uniformly across types.
27. What is the fragile base-class problem?
   - **Key note:** A base-class change can silently break subclasses that depend on its internal behavior, showing inheritance's tight coupling.
28. What is the diamond inheritance problem, and how do languages address it?
   - **Key note:** Multiple paths inherit the same base member; languages use restrictions, explicit resolution, virtual inheritance, or linearization.
29. What is the difference between delegation and inheritance?
   - **Key note:** Delegation forwards work to a collaborator, while inheritance obtains behavior through an “is-a” relationship.
30. What is the Law of Demeter?
   - **Key note:** An object should talk mainly to close collaborators, avoiding long navigation chains and structural coupling.
31. What is tell-don't-ask?
   - **Key note:** Ask an object to perform behavior instead of extracting its state and making decisions elsewhere.
32. What is an anemic domain model, and when might it be acceptable?
   - **Key note:** It separates data from behavior; it may suit simple CRUD but loses encapsulated domain rules in complex systems.
33. What are entities, value objects, and aggregates in domain-driven design?
   - **Key note:** Entities have identity, value objects are defined by attributes, and aggregates protect consistency behind one root.
34. How do you preserve class invariants?
   - **Key note:** Validate construction and every state transition while preventing uncontrolled mutation of internal state.
35. What is double dispatch, and how does the Visitor pattern use it?
   - **Key note:** Behavior is selected from two runtime types; Visitor combines element dispatch with an overloaded visitor method.
36. What is the difference between inheritance and interface implementation?
   - **Key note:** Inheritance may reuse state and implementation; an interface establishes a behavioral contract without class ancestry.
37. What is dependency inversion at an architectural boundary?
   - **Key note:** Domain policy owns abstractions, while external details implement adapters that point inward toward the domain.
38. What is constructor injection, and why is it usually preferred?
   - **Key note:** Required dependencies are supplied at creation, making valid construction, immutability, and testing straightforward.
39. When is a service locator considered an anti-pattern?
   - **Key note:** It hides dependencies behind global lookup, making contracts unclear and tests/order of initialization fragile.
40. What is object identity, and how does it affect equality and hashing?
   - **Key note:** Identity distinguishes instances; equality and hash codes must remain consistent, especially when objects are map keys.
41. Why can mutable objects be unsafe as hash-map keys?
   - **Key note:** Changing fields used in hashing can move the logical bucket, making the stored key impossible to find.
42. What is the difference between a domain model and a data-transfer object?
   - **Key note:** A domain model contains rules and behavior; a DTO carries data across a boundary and should not own domain logic.
43. How would you break a circular dependency between classes or modules?
   - **Key note:** Extract a stable abstraction, move shared responsibility, introduce events, or redraw ownership boundaries.
44. What are temporal coupling and hidden state?
   - **Key note:** Temporal coupling requires calls in a specific order; hidden state makes behavior depend on non-obvious history.
45. When should a class be sealed or final?
   - **Key note:** Prevent extension when invariants, security, immutability, or an intentionally closed hierarchy require controlled subtypes.
