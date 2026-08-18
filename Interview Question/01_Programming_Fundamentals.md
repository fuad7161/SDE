# Programming Fundamentals

1. What is the difference between a compiler and an interpreter?
   - A compiler translates an entire program before execution, while an interpreter translates and executes instructions as the program runs.
2. What is the difference between source code, bytecode, and machine code?
   - Source code is human-readable; bytecode is portable intermediate code for a virtual machine; machine code is CPU-specific binary instructions.
3. What are variables, constants, and data types?
   - A variable stores a changeable value, a constant should not change, and a data type defines the value's form and valid operations.
4. What is the difference between primitive and reference data types?
   - A primitive directly represents a simple value, while a reference identifies an object or data stored elsewhere in memory.
5. What is type casting? What is the difference between implicit and explicit casting?
   - Casting converts one type to another; safe conversions may happen implicitly, while potentially lossy conversions require an explicit request.
6. What is the difference between static and dynamic typing?
   - Static typing checks types mainly before execution, whereas dynamic typing determines and checks them while the program runs.
7. What is the difference between pass by value and pass by reference?
   - Pass by value copies an argument's value; pass by reference lets the function access the caller's original storage. Some languages pass object references by value.
8. What is variable scope? Explain local, global, and block scope.
   - Scope determines where a name is visible: local scope is within a function, global scope spans a module or program, and block scope is limited to a block.
9. What is the lifetime of a variable?
   - A variable's lifetime is the period during execution when its storage exists and its value can be used.
10. What is the difference between stack and heap memory?
   - The stack usually stores call frames and short-lived local data, while the heap stores dynamically allocated objects with less predictable lifetimes.
11. What is a function or method, and why do we use one?
   - It is a named, reusable unit of behavior that improves decomposition, readability, testing, and reuse; a method belongs to a type or object.
12. What is recursion? What must every recursive function contain?
   - Recursion occurs when a function calls itself, and it needs a base case that stops further calls.
13. What is the difference between iteration and recursion?
   - Iteration repeats with loops, while recursion repeats through function calls; recursion can be clearer but consumes call-stack space.
14. What is a pure function, and what is a side effect?
   - A pure function gives the same output for the same input and changes no external state; modifying state or performing I/O is a side effect.
15. What is function overloading?
   - Overloading defines multiple functions with the same name but different parameter lists, letting the compiler select the appropriate version.
16. What is mutable versus immutable data?
   - Mutable data can change after creation; immutable data cannot, so an apparent modification creates a new value.
17. What is garbage collection?
   - Garbage collection automatically finds and reclaims heap memory occupied by objects that are no longer reachable.
18. What is an exception? How is it different from an error?
   - An exception represents an abnormal condition a program may handle; an error often describes a more serious fault, though exact terminology varies by language.
19. What is the purpose of `try`, `catch`, and `finally`?
   - `try` surrounds risky code, `catch` handles a matching exception, and `finally` runs cleanup whether an exception occurs or not.
20. What is the difference between syntax, runtime, and logical errors?
   - Syntax errors violate language rules, runtime errors occur during execution, and logical errors let the program run but produce incorrect results.
21. What is short-circuit evaluation?
   - A Boolean expression stops as soon as its result is known, such as skipping the right side of `false && expression`.
22. What is the difference between equality and identity/reference equality?
   - Value equality compares contents, while identity equality checks whether two references point to the exact same object.
23. What is a package, module, or namespace?
   - These organize related code, control visibility, encourage reuse, and prevent naming conflicts; their exact semantics depend on the language.
24. What are command-line arguments and environment variables?
   - Command-line arguments are values supplied when launching a program; environment variables are named configuration values inherited from its environment.
25. What are coding conventions, and why do they matter?
   - Coding conventions are shared rules for naming, formatting, and structure that make a codebase consistent and easier to maintain.
