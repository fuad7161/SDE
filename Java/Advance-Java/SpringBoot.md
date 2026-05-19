# Spring & Spring Boot — In-Depth Notes

---

## Table of Contents

1. [IoC Container & Dependency Injection Internals](#1-ioc-container--dependency-injection-internals)
2. [Bean Lifecycle & Scopes](#2-bean-lifecycle--scopes)
3. [@Transactional — Propagation, Isolation, Rollback](#3-transactional--propagation-isolation-rollback)
4. [AOP — Proxy Mechanism (JDK vs CGLIB)](#4-aop--proxy-mechanism-jdk-vs-cglib)
5. [Spring Security (JWT, OAuth2, Filter Chain)](#5-spring-security-jwt-oauth2-filter-chain)
6. [Spring Boot Auto-Configuration Internals](#6-spring-boot-auto-configuration-internals)
7. [ApplicationContext vs BeanFactory](#7-applicationcontext-vs-beanfactory)

---

## 1. IoC Container & Dependency Injection Internals

### Inversion of Control (IoC)

Traditional code: **you** create and wire dependencies.  
IoC: the **container** creates and injects dependencies — control is inverted.

```java
// Without IoC — tight coupling
class OrderService {
    private PaymentService payment = new StripePaymentService(); // hardcoded
    private EmailService   email   = new SmtpEmailService();     // hardcoded
}

// With IoC — loose coupling
class OrderService {
    private final PaymentService payment;
    private final EmailService   email;

    // Spring injects the implementations
    public OrderService(PaymentService payment, EmailService email) {
        this.payment = payment;
        this.email   = email;
    }
}
```

---

### Dependency Injection Types

#### Constructor Injection (Recommended)
```java
@Service
public class OrderService {
    private final PaymentService paymentService;
    private final InventoryService inventoryService;

    // @Autowired is optional on single constructor (Spring 4.3+)
    public OrderService(PaymentService paymentService,
                        InventoryService inventoryService) {
        this.paymentService   = paymentService;
        this.inventoryService = inventoryService;
    }
}
// Advantages: immutable fields, all deps visible, easy to unit test (no Spring needed)
```

#### Setter Injection (Optional dependencies)
```java
@Service
public class ReportService {
    private CacheService cacheService;

    @Autowired(required = false)   // optional — works without cache too
    public void setCacheService(CacheService cacheService) {
        this.cacheService = cacheService;
    }
}
```

#### Field Injection (Avoid in production)
```java
@Service
public class UserService {
    @Autowired                 // hidden dependency — bad for testing
    private UserRepository repo;
}
```

---

### How the IoC Container Works Internally

```
Startup:
  1. Scan classpath for @Component, @Service, @Repository, @Controller
  2. Read @Configuration classes and @Bean methods
  3. Build BeanDefinition registry (metadata: class, scope, init method, ...)
  4. Resolve dependencies (topological sort)
  5. Instantiate beans (constructor / factory method)
  6. Inject dependencies (@Autowired, @Value, @Resource)
  7. Run BeanPostProcessors (AOP proxies created here)
  8. Call @PostConstruct / InitializingBean.afterPropertiesSet()
  9. ApplicationContext is ready
```

```java
// @Configuration — explicit bean wiring
@Configuration
public class AppConfig {

    @Bean
    public DataSource dataSource() {
        return DataSourceBuilder.create()
            .url("jdbc:postgresql://localhost/mydb")
            .username("user")
            .password("pass")
            .build();
    }

    @Bean
    public UserRepository userRepository(DataSource ds) {
        return new JdbcUserRepository(ds);  // Spring injects DataSource bean
    }
}
```

---

### `@Qualifier` and `@Primary`

When multiple beans match a type, Spring needs a hint:

```java
interface PaymentService { void pay(double amount); }

@Service @Primary           // default choice when no qualifier specified
class StripePaymentService implements PaymentService { ... }

@Service
class PayPalPaymentService implements PaymentService { ... }

// Injection
@Autowired
private PaymentService payment;            // gets StripePaymentService (@Primary)

@Autowired @Qualifier("payPalPaymentService")
private PaymentService payment;            // gets PayPalPaymentService
```

---

## 2. Bean Lifecycle & Scopes

### Bean Lifecycle

```
Container starts
    │
    ▼
1. Instantiation        new MyBean()
    │
    ▼
2. Dependency Injection  @Autowired fields/setters populated
    │
    ▼
3. BeanNameAware         setBeanName(name)
   BeanFactoryAware      setBeanFactory(factory)
   ApplicationContextAware setApplicationContext(ctx)
    │
    ▼
4. BeanPostProcessor     postProcessBeforeInitialization()   ← AOP proxy created here
    │
    ▼
5. @PostConstruct        init method (custom setup)
   InitializingBean      afterPropertiesSet()
    │
    ▼
6. Bean is READY — used by application
    │
    ▼  (container shuts down)
7. @PreDestroy           cleanup method
   DisposableBean        destroy()
    │
    ▼
Bean destroyed
```

```java
@Component
public class DatabaseConnection {
    private Connection connection;

    @PostConstruct
    public void init() {
        connection = openConnection();   // runs after injection, before use
        System.out.println("Connection opened");
    }

    @PreDestroy
    public void cleanup() {
        closeConnection(connection);     // runs before bean is destroyed
        System.out.println("Connection closed");
    }
}
```

---

### Bean Scopes

| Scope | Instance | Use case |
|---|---|---|
| `singleton` (default) | One per ApplicationContext | Stateless services, repositories |
| `prototype` | New instance per injection/request | Stateful beans, non-thread-safe objects |
| `request` | One per HTTP request | Web: request-scoped data |
| `session` | One per HTTP session | Web: user session data |
| `application` | One per ServletContext | Web: app-wide shared state |

```java
@Component
@Scope("singleton")           // default — one shared instance
public class UserService { ... }

@Component
@Scope("prototype")           // new instance every time it's injected
public class ShoppingCart { ... }

// Web scopes
@Component
@Scope(value = WebApplicationContext.SCOPE_REQUEST, proxyMode = ScopedProxyMode.TARGET_CLASS)
public class RequestContext { ... }
```

**Injecting prototype into singleton** (scoped proxy needed):
```java
@Component
@Scope(value = "prototype", proxyMode = ScopedProxyMode.TARGET_CLASS)
public class TaskProcessor { ... }

@Service
public class TaskService {
    @Autowired
    private TaskProcessor processor;  // Spring injects a proxy; each call gets a new instance
}
```

---

## 3. @Transactional — Propagation, Isolation, Rollback

### How `@Transactional` Works

Spring wraps the bean in a **proxy**. The proxy intercepts calls to `@Transactional` methods and manages the transaction:

```
Caller → [Proxy] → begin transaction → [Your method] → commit/rollback
```

```java
@Service
public class OrderService {

    @Transactional
    public void placeOrder(Order order) {
        orderRepository.save(order);
        paymentService.charge(order.getTotal()); // if this throws → rollback
        inventoryService.reserve(order.getItems());
    }
}
```

---

### Propagation Levels

Controls what happens when a `@Transactional` method is called from within an existing transaction:

| Propagation | Behaviour |
|---|---|
| `REQUIRED` (default) | Join existing tx; create new one if none |
| `REQUIRES_NEW` | Always suspend existing tx and start a new one |
| `NESTED` | Create a savepoint; inner rolls back to savepoint, outer can still commit |
| `SUPPORTS` | Join existing tx if present; run non-transactionally if none |
| `NOT_SUPPORTED` | Suspend existing tx; run non-transactionally |
| `MANDATORY` | Must run inside existing tx; throw if none |
| `NEVER` | Must NOT run inside a tx; throw if one exists |

```java
@Service
public class AuditService {

    // Always runs in its own transaction — audit log persisted even if outer tx rolls back
    @Transactional(propagation = Propagation.REQUIRES_NEW)
    public void logAction(String action) {
        auditRepository.save(new AuditLog(action));
    }
}

@Service
public class OrderService {
    @Autowired AuditService auditService;

    @Transactional
    public void placeOrder(Order order) {
        orderRepository.save(order);
        auditService.logAction("ORDER_PLACED");   // runs in separate tx
        throw new RuntimeException("test");       // outer tx rolls back, but audit log persists
    }
}
```

---

### Isolation Levels

Controls **visibility of concurrent transactions** to each other:

| Isolation | Dirty Read | Non-Repeatable Read | Phantom Read |
|---|---|---|---|
| `READ_UNCOMMITTED` | ✅ possible | ✅ possible | ✅ possible |
| `READ_COMMITTED` (most DBs default) | ❌ prevented | ✅ possible | ✅ possible |
| `REPEATABLE_READ` (MySQL default) | ❌ prevented | ❌ prevented | ✅ possible |
| `SERIALIZABLE` | ❌ prevented | ❌ prevented | ❌ prevented |

```java
@Transactional(isolation = Isolation.REPEATABLE_READ)
public BigDecimal getAccountBalance(Long accountId) {
    // Same query returns same result within this transaction
    return accountRepository.findBalanceById(accountId);
}

@Transactional(isolation = Isolation.SERIALIZABLE)
public void transferFunds(Long from, Long to, BigDecimal amount) {
    // Fully isolated — no concurrent interference
}
```

---

### Rollback Rules

Default: rollback on **unchecked** exceptions (`RuntimeException`, `Error`).  
Checked exceptions do **not** trigger rollback by default.

```java
// Custom rollback rules
@Transactional(
    rollbackFor    = {IOException.class, CustomException.class},  // also rollback for these
    noRollbackFor  = {OptimisticLockException.class}              // don't rollback for this
)
public void processFile(String path) throws IOException { ... }
```

### Common `@Transactional` Pitfalls

```java
// ❌ PITFALL 1 — self-invocation bypasses proxy
@Service
public class UserService {
    @Transactional
    public void doA() { doB(); }  // calls doB directly — no proxy — no transaction!

    @Transactional
    public void doB() { ... }
}
// FIX: inject self or use AspectJ weaving

// ❌ PITFALL 2 — private methods are not proxied
@Transactional
private void save() { ... }  // Spring ignores @Transactional on private methods

// ❌ PITFALL 3 — @Transactional on @Component without interface can fail with JDK proxy
// FIX: use CGLIB (spring.aop.proxy-target-class=true, default in Spring Boot)
```

---

## 4. AOP — Proxy Mechanism (JDK vs CGLIB)

### What is AOP?

**Aspect-Oriented Programming** separates cross-cutting concerns (logging, security, transactions) from business logic.

```
Without AOP — repeated boilerplate everywhere:
  UserService.save()     → log → check security → begin tx → business logic → end tx
  OrderService.place()   → log → check security → begin tx → business logic → end tx

With AOP — define once, apply everywhere:
  @Before("@annotation(Transactional)") → begin transaction
  @After("@annotation(Transactional)")  → commit/rollback
```

### Key AOP Concepts

| Term | Meaning |
|---|---|
| **Aspect** | Class with cross-cutting logic (`@Aspect`) |
| **Advice** | The action to run (`@Before`, `@After`, `@Around`) |
| **Pointcut** | Expression defining which methods to intercept |
| **Join point** | The actual method execution being intercepted |
| **Weaving** | Linking aspects with target objects |

```java
@Aspect
@Component
public class LoggingAspect {

    // Pointcut — all methods in service package
    @Pointcut("execution(* com.example.service.*.*(..))")
    public void serviceMethods() {}

    @Before("serviceMethods()")
    public void logBefore(JoinPoint jp) {
        System.out.println("Calling: " + jp.getSignature().getName());
    }

    @AfterReturning(pointcut = "serviceMethods()", returning = "result")
    public void logAfter(Object result) {
        System.out.println("Returned: " + result);
    }

    @AfterThrowing(pointcut = "serviceMethods()", throwing = "ex")
    public void logException(Exception ex) {
        System.out.println("Exception: " + ex.getMessage());
    }

    // Around — full control (used by @Transactional internally)
    @Around("@annotation(com.example.Timed)")
    public Object measureTime(ProceedingJoinPoint pjp) throws Throwable {
        long start = System.currentTimeMillis();
        Object result = pjp.proceed();           // invoke the actual method
        long elapsed = System.currentTimeMillis() - start;
        System.out.println(pjp.getSignature() + " took " + elapsed + "ms");
        return result;
    }
}
```

---

### JDK Dynamic Proxy vs CGLIB

Spring AOP creates a **proxy object** that wraps the bean. Two mechanisms:

#### JDK Dynamic Proxy
- Requires the bean to implement **at least one interface**
- Proxy implements the same interface(s)
- Uses `java.lang.reflect.Proxy`

```java
interface UserService { User findById(Long id); }

@Service
class UserServiceImpl implements UserService {
    public User findById(Long id) { return repo.findById(id); }
}

// Spring creates: Proxy implements UserService
//                 delegates to UserServiceImpl + applies advice
UserService proxy = (UserService) Proxy.newProxyInstance(...);
```

#### CGLIB Proxy
- Works on **concrete classes** (no interface needed)
- Creates a **subclass** of the target at runtime using bytecode generation
- Cannot proxy `final` classes or `final` methods

```java
@Service
class OrderService {           // no interface
    public void place() { ... }
}
// Spring (CGLIB) creates: class OrderService$$SpringCGLIB$$0 extends OrderService
```

```yaml
# Spring Boot default (CGLIB for all, including interface beans)
spring:
  aop:
    proxy-target-class: true   # default since Spring Boot 2.x
```

| | JDK Dynamic Proxy | CGLIB |
|---|---|---|
| Requires interface | Yes | No |
| Mechanism | `Proxy.newProxyInstance` | Subclass generation |
| `final` classes | Can proxy | ❌ Cannot |
| `final` methods | Proxied via interface | ❌ Not intercepted |
| Performance | Slightly slower (reflection) | Faster |
| Spring Boot default | Overridden by CGLIB | ✅ Default |

---

## 5. Spring Security (JWT, OAuth2, Filter Chain)

### Security Filter Chain

Every HTTP request passes through a chain of `Filter` objects before reaching the controller:

```
HTTP Request
    │
    ▼
SecurityContextPersistenceFilter   — restore SecurityContext from session
    │
    ▼
UsernamePasswordAuthenticationFilter  — process login form POST
    │
    ▼
JwtAuthenticationFilter  (custom)  — validate JWT, set Authentication
    │
    ▼
ExceptionTranslationFilter         — handle AuthenticationException, AccessDeniedException
    │
    ▼
FilterSecurityInterceptor          — enforce @PreAuthorize, URL access rules
    │
    ▼
DispatcherServlet → Controller
```

```java
@Configuration
@EnableWebSecurity
public class SecurityConfig {

    @Bean
    public SecurityFilterChain filterChain(HttpSecurity http) throws Exception {
        return http
            .csrf(csrf -> csrf.disable())
            .sessionManagement(sm -> sm.sessionCreationPolicy(SessionCreationPolicy.STATELESS))
            .authorizeHttpRequests(auth -> auth
                .requestMatchers("/api/auth/**").permitAll()
                .requestMatchers("/api/admin/**").hasRole("ADMIN")
                .anyRequest().authenticated()
            )
            .addFilterBefore(jwtAuthFilter, UsernamePasswordAuthenticationFilter.class)
            .build();
    }
}
```

---

### JWT Authentication Flow

```
Login:
  POST /auth/login {username, password}
      → AuthenticationManager.authenticate()
      → UserDetailsService.loadUserByUsername()
      → BCryptPasswordEncoder.matches()
      → generate JWT (header.payload.signature)
      → return JWT to client

Subsequent Requests:
  GET /api/orders
  Authorization: Bearer <JWT>
      → JwtAuthFilter intercepts
      → validate signature + expiry
      → extract username from payload
      → load UserDetails, set Authentication in SecurityContext
      → proceed to controller
```

```java
@Component
public class JwtAuthFilter extends OncePerRequestFilter {

    @Autowired private JwtService jwtService;
    @Autowired private UserDetailsService userDetailsService;

    @Override
    protected void doFilterInternal(HttpServletRequest request,
                                    HttpServletResponse response,
                                    FilterChain chain) throws ServletException, IOException {
        String authHeader = request.getHeader("Authorization");
        if (authHeader == null || !authHeader.startsWith("Bearer ")) {
            chain.doFilter(request, response);
            return;
        }

        String token    = authHeader.substring(7);
        String username = jwtService.extractUsername(token);

        if (username != null && SecurityContextHolder.getContext().getAuthentication() == null) {
            UserDetails user = userDetailsService.loadUserByUsername(username);
            if (jwtService.isTokenValid(token, user)) {
                UsernamePasswordAuthenticationToken auth =
                    new UsernamePasswordAuthenticationToken(user, null, user.getAuthorities());
                auth.setDetails(new WebAuthenticationDetailsSource().buildDetails(request));
                SecurityContextHolder.getContext().setAuthentication(auth);
            }
        }
        chain.doFilter(request, response);
    }
}
```

---

### OAuth2 (Brief)

```
Resource Owner (User)
    │  1. Login / consent
    ▼
Authorization Server (Google, GitHub, Keycloak)
    │  2. Issues access_token + refresh_token
    ▼
Client App (your Spring Boot service)
    │  3. Send access_token
    ▼
Resource Server (your API)
    │  4. Validate token, serve resource
```

```yaml
# application.yml — OAuth2 client (social login)
spring:
  security:
    oauth2:
      client:
        registration:
          google:
            client-id: YOUR_CLIENT_ID
            client-secret: YOUR_SECRET
            scope: openid,email,profile
```

---

## 6. Spring Boot Auto-Configuration Internals

Spring Boot eliminates boilerplate XML/Java config by auto-configuring beans based on **what's on the classpath**.

### How It Works

```
@SpringBootApplication
  = @Configuration
  + @ComponentScan
  + @EnableAutoConfiguration   ← the magic

@EnableAutoConfiguration triggers:
  1. Reads META-INF/spring/org.springframework.boot.autoconfigure.AutoConfiguration.imports
  2. Finds all @AutoConfiguration classes (e.g., DataSourceAutoConfiguration)
  3. Each has @ConditionalOn... guards — only activates if conditions are met
  4. Registers beans that are not already defined by the user
```

```java
// Spring Boot's DataSourceAutoConfiguration (simplified)
@AutoConfiguration
@ConditionalOnClass(DataSource.class)             // only if DataSource class is on classpath
@ConditionalOnMissingBean(DataSource.class)       // only if user hasn't defined their own
@EnableConfigurationProperties(DataSourceProperties.class)
public class DataSourceAutoConfiguration {

    @Bean
    public DataSource dataSource(DataSourceProperties props) {
        return props.initializeDataSourceBuilder().build();
    }
}
```

### Key `@Conditional` Annotations

| Annotation | Activates when |
|---|---|
| `@ConditionalOnClass` | Class is present on classpath |
| `@ConditionalOnMissingClass` | Class is absent |
| `@ConditionalOnBean` | Bean of type exists |
| `@ConditionalOnMissingBean` | Bean of type does NOT exist |
| `@ConditionalOnProperty` | Property has specific value |
| `@ConditionalOnWebApplication` | Running as web app |

### Custom Auto-Configuration

```java
@AutoConfiguration
@ConditionalOnProperty(prefix = "mylib", name = "enabled", havingValue = "true")
@ConditionalOnMissingBean(MyService.class)
public class MyLibAutoConfiguration {

    @Bean
    public MyService myService() {
        return new DefaultMyService();
    }
}
```

```
# META-INF/spring/org.springframework.boot.autoconfigure.AutoConfiguration.imports
com.example.MyLibAutoConfiguration
```

### Debugging Auto-Configuration

```bash
# See which auto-configurations were applied and which were skipped (and why)
java -jar app.jar --debug

# Or in application.properties
logging.level.org.springframework.boot.autoconfigure=DEBUG
```

---

## 7. ApplicationContext vs BeanFactory

### `BeanFactory`

The root interface for the Spring IoC container — basic bean instantiation and wiring.

- Lazy initialization by default (beans created only on first `getBean()`)
- No built-in support for AOP, events, internationalization

```java
BeanFactory factory = new XmlBeanFactory(new ClassPathResource("beans.xml")); // legacy
MyBean bean = (MyBean) factory.getBean("myBean");
```

---

### `ApplicationContext`

Extends `BeanFactory` with enterprise features — this is what Spring Boot uses.

- **Eager initialization** of singleton beans at startup
- `ApplicationEvent` publishing/listening
- `MessageSource` (i18n)
- `Environment` / `@Value` property resolution
- `BeanFactoryPostProcessor` and `BeanPostProcessor` auto-detection
- Integration with AOP, transactions, security

```java
ApplicationContext ctx = new AnnotationConfigApplicationContext(AppConfig.class);
MyService service = ctx.getBean(MyService.class);

// Publishing custom events
ctx.publishEvent(new UserRegisteredEvent(this, user));

// Listening to events
@EventListener
public void onUserRegistered(UserRegisteredEvent event) {
    emailService.sendWelcome(event.getUser());
}
```

### `ApplicationContext` Implementations

| Class | Use case |
|---|---|
| `AnnotationConfigApplicationContext` | Java config / annotations (non-web) |
| `AnnotationConfigServletWebServerApplicationContext` | Spring Boot web apps (auto-used) |
| `ClassPathXmlApplicationContext` | Legacy XML config |
| `FileSystemXmlApplicationContext` | XML config from filesystem path |

### Comparison

| Feature | `BeanFactory` | `ApplicationContext` |
|---|---|---|
| Bean instantiation | Lazy (on demand) | Eager (at startup) |
| AOP support | No | Yes |
| Event publishing | No | Yes (`ApplicationEvent`) |
| i18n (`MessageSource`) | No | Yes |
| `@PostConstruct` / `@PreDestroy` | No | Yes |
| `Environment` / properties | No | Yes |
| Used in production | Rarely | Always |
