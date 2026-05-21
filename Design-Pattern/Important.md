[50 system design pattern](https://designgurus.substack.com/p/50-system-design-patterns-every-engineer)

Conceptual questions you'll definitely face:

SOLID principles — relationship to design patterns
Singleton thread-safety — sync.Once in Go, double-checked locking in Java
Factory vs Abstract Factory vs Builder — when to use which
Decorator vs Proxy vs Adapter — they look similar, differences matter
Strategy vs State — the subtle distinction
Observer pattern — how it maps to event systems, message queues
Dependency Injection — often asked alongside patterns


🎯 Go-Specific Pattern Questions (for your interview)
Since you're prepping for a mid-senior Go role:

Functional Options pattern — idiomatic Go builder alternative
Singleton with sync.Once
Middleware chaining (Chain of Responsibility)
Interface-based Strategy — Go's duck typing makes this natural
Worker Pool pattern — goroutines + channels (Go-native)
Pipeline pattern — chained channels
Fan-Out / Fan-In — concurrency pattern often asked in Go interviews


💬 Classic Interview Questions

"Which patterns have you used in production? Why?"
"How does Go's interface system replace traditional OOP patterns?"
"Explain Observer — how would you implement it without a library?"
"When would Singleton be an anti-pattern?"
"How does the Repository Pattern relate to design patterns?"
