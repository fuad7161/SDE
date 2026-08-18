# Software Engineering and Testing

1. What is the software development lifecycle?
   - The SDLC is the structured process used to plan, build, test, release, operate, and eventually retire software.
2. What are the typical phases of software development?
   - Typical phases are requirements, planning and design, implementation, testing, deployment, operation, and maintenance.
3. What is the difference between functional and non-functional requirements?
   - Functional requirements describe what a system does; non-functional requirements describe qualities and constraints such as speed, security, and availability.
4. Compare Agile and Waterfall development.
   - Waterfall follows largely sequential phases, while Agile delivers in short feedback-driven iterations and adapts requirements as learning occurs.
5. What are Scrum roles, events, and artifacts?
   - Scrum defines the Product Owner, Scrum Master, and Developers; the Sprint plus Planning, Daily Scrum, Review, and Retrospective events; and Product Backlog, Sprint Backlog, and Increment artifacts.
6. What are cohesion and coupling?
   - Cohesion measures how closely a module's responsibilities belong together; coupling measures its dependency on other modules.
7. What does DRY mean?
   - “Don't Repeat Yourself” means keeping each piece of knowledge in one authoritative place, while avoiding premature or misleading abstractions.
8. What do KISS and YAGNI mean?
   - KISS favors simple solutions, and YAGNI advises against building speculative features before they are actually required.
9. What is technical debt?
   - Technical debt is the future cost created by expedient design or implementation choices that make later changes harder.
10. What is refactoring, and how is it different from rewriting?
   - Refactoring incrementally improves internal structure without changing behavior; rewriting replaces substantial implementation, usually with greater risk.
11. What is clean code?
   - Clean code is easy to understand, test, change, and operate, with clear names, focused units, and minimal accidental complexity.
12. Why are code reviews useful?
   - Reviews find defects, improve design and consistency, spread knowledge, and provide shared ownership of code.
13. What is an API contract?
   - It defines an API's inputs, outputs, behavior, errors, and compatibility expectations on which clients may rely.
14. What is the difference between a library, framework, and platform?
   - Code calls a library, a framework typically controls flow and calls application code, and a platform supplies a broader runtime and service environment.
15. What is unit testing?
   - Unit testing checks a small unit of behavior with fast, deterministic, focused tests; isolation may come from design or test doubles, not necessarily from mocking every collaborator.
16. What is the difference between unit, integration, system, and acceptance testing?
   - They test individual units, component interactions, the complete system, and whether the product satisfies user or business requirements.
17. What are black-box and white-box testing?
   - Black-box testing uses externally visible behavior; white-box testing designs checks with knowledge of internal structure.
18. What is regression testing?
   - It reruns tests after changes to ensure previously working behavior has not broken.
19. What is test-driven development?
   - TDD cycles through writing a failing test, implementing the simplest passing code, and refactoring safely.
20. What is the Arrange-Act-Assert test pattern?
   - Arrange prepares data and dependencies, Act performs the behavior, and Assert verifies the result.
21. What are mocks, stubs, spies, and fakes?
   - Stubs return canned data, mocks verify expected interactions, spies record calls, and fakes provide lightweight working implementations.
22. What makes a test reliable and maintainable?
   - It is deterministic, isolated, focused on behavior, clearly named, fast enough for its level, and free from unnecessary implementation coupling.
23. What is code coverage, and why is 100% coverage not sufficient?
   - Coverage reports which code tests execute, but it cannot prove the assertions are meaningful or that all behaviors and edge cases are correct.
24. What is continuous integration?
   - CI frequently combines changes and automatically builds and tests them to reveal integration problems early.
25. What is the difference between continuous delivery and continuous deployment?
   - Delivery keeps every change releasable with a manual production decision; deployment automatically releases every passing change.
26. What is semantic versioning?
   - SemVer uses `MAJOR.MINOR.PATCH`: breaking changes increment major, compatible features minor, and compatible fixes patch.
27. What are logs, metrics, and traces?
   - Logs record events, metrics aggregate numerical measurements, and traces follow a request across components.
28. How would you debug an issue that cannot be reproduced locally?
   - Gather exact context and observability data, compare environments, correlate recent changes, form hypotheses, and add safe targeted instrumentation.

## Medium to Advanced

29. What is the testing pyramid, and when does it become misleading?
   - **Key note:** Favor many fast focused tests and fewer broad tests, but choose layers according to actual architecture and risk.
30. What is the difference between state-based and interaction-based testing?
   - **Key note:** State tests verify outcomes; interaction tests verify collaborator calls and can couple tests more tightly to implementation.
31. What are contract tests, and which integration failures do they catch?
   - **Key note:** They verify provider/consumer interface expectations without full end-to-end environments, catching incompatible independent changes.
32. What are property-based tests?
   - **Key note:** They generate many inputs to verify general invariants and shrink failures to small counterexamples.
33. What are mutation tests?
   - **Key note:** They deliberately alter production code; surviving mutations reveal tests that execute code without effectively checking behavior.
34. How do flaky tests harm a delivery pipeline?
   - **Key note:** They destroy trust, hide real failures, slow feedback, and encourage unsafe reruns or ignored results.
35. How would you test concurrent code deterministically?
   - **Key note:** Control scheduling and clocks, synchronize known points, assert invariants, and repeat stress/model checks without arbitrary sleeps.
36. How should external dependencies be tested?
   - **Key note:** Combine unit fakes, contract tests, targeted integration tests, and a small number of end-to-end checks.
37. What is consumer-driven contract testing?
   - **Key note:** Consumers publish required interactions and providers verify them, supporting safe independent deployment.
38. What is the difference between load, stress, spike, soak, and capacity testing?
   - **Key note:** They test expected traffic, breaking limits, sudden bursts, long-duration stability, and maximum sustainable throughput.
39. How do you choose useful test data?
   - **Key note:** Cover boundaries, invalid states, realistic distributions, security cases, and production-like scale without exposing real personal data.
40. What is hexagonal or ports-and-adapters architecture?
   - **Key note:** Core policy exposes ports while replaceable adapters connect databases, frameworks, queues, and user interfaces.
41. Compare layered, clean, hexagonal, and event-driven architectures.
   - **Key note:** Compare dependency direction, boundary enforcement, coupling, flow, testability, and operational complexity—not labels alone.
42. What is a bounded context?
   - **Key note:** It is a boundary within which one domain model and vocabulary have precise, consistent meaning.
43. What is an architecture decision record?
   - **Key note:** An ADR captures context, chosen decision, alternatives, and consequences so future teams understand the trade-off.
44. What are feature flags, and what operational debt can they create?
   - **Key note:** Flags decouple release from deployment but add state combinations and dead paths; assign ownership and removal dates.
45. How do blue-green and canary releases reduce deployment risk?
   - **Key note:** Blue-green enables environment switching; canary limits initial exposure and evaluates real production signals.
46. What is a blameless incident review?
   - **Key note:** It studies system conditions and decisions without personal blame, producing owned actions that reduce recurrence and impact.
47. What are SLI, SLO, SLA, and error budget?
   - **Key note:** They define measurement, reliability target, external promise, and permitted unreliability for balancing change and stability.
48. How do trunk-based development and long-lived feature branches differ?
   - **Key note:** Trunk-based work integrates small changes frequently; long-lived branches defer integration and increase merge and release risk.
