# Testing — In-Depth Notes

---

## Table of Contents

1. [Unit vs Integration vs E2E Testing](#1-unit-vs-integration-vs-e2e-testing)
2. [JUnit 5, Mockito (@Mock, @Spy, @Captor)](#2-junit-5-mockito-mock-spy-captor)
3. [@SpringBootTest vs Slice Tests](#3-springboottest-vs-slice-tests)
4. [TDD Approach](#4-tdd-approach)

---

## 1. Unit vs Integration vs E2E Testing

### The Testing Pyramid

```
         /\
        /  \        ← E2E (few, slow, expensive)
       / E2E\
      /──────\
     /        \     ← Integration (some)
    /Integration\
   /────────────\
  /              \  ← Unit (many, fast, cheap)
 /   Unit Tests   \
/──────────────────\
```

### Unit Tests

Test a **single class/method in isolation** — all dependencies are mocked.

- Fast (milliseconds), no Spring context, no DB, no network
- High quantity — the foundation of the pyramid

```java
@ExtendWith(MockitoExtension.class)
class OrderServiceTest {

    @Mock  OrderRepository orderRepository;
    @Mock  PaymentService  paymentService;

    @InjectMocks
    OrderService orderService;   // dependencies auto-injected from mocks above

    @Test
    void placeOrder_shouldSaveAndChargePayment() {
        // Given
        Order order = new Order(1L, BigDecimal.valueOf(100));
        when(orderRepository.save(order)).thenReturn(order);
        when(paymentService.charge(100)).thenReturn(true);

        // When
        Order result = orderService.placeOrder(order);

        // Then
        assertThat(result).isNotNull();
        verify(orderRepository).save(order);
        verify(paymentService).charge(100);
    }
}
```

---

### Integration Tests

Test **multiple components together** — real DB, real Spring context, or real HTTP calls. Slower but catch wiring issues.

```java
@SpringBootTest
@Transactional                  // rolls back after each test
class OrderRepositoryTest {

    @Autowired OrderRepository orderRepository;

    @Test
    void findByStatus_shouldReturnMatchingOrders() {
        orderRepository.save(new Order("PENDING"));
        orderRepository.save(new Order("COMPLETED"));

        List<Order> pending = orderRepository.findByStatus("PENDING");

        assertThat(pending).hasSize(1);
    }
}
```

---

### End-to-End (E2E) Tests

Test the **entire system** from a user's perspective — real browser, real HTTP, real DB.

```java
// Using REST Assured for API-level E2E
@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
class OrderApiE2ETest {

    @LocalServerPort int port;

    @Test
    void createOrder_fullFlow() {
        given()
            .baseUri("http://localhost:" + port)
            .header("Authorization", "Bearer " + getToken())
            .contentType(ContentType.JSON)
            .body("""{"productId": 1, "quantity": 2}""")
        .when()
            .post("/api/orders")
        .then()
            .statusCode(201)
            .body("status", equalTo("PENDING"))
            .body("id", notNullValue());
    }
}
```

### Comparison

| | Unit | Integration | E2E |
|---|---|---|---|
| Scope | Single class | Multiple components | Full system |
| Speed | Very fast (ms) | Slow (seconds) | Very slow (minutes) |
| Dependencies | All mocked | Real DB, partial Spring | Real everything |
| Confidence | Low (mocks may not match) | Medium | High |
| Quantity | Many (hundreds) | Some (tens) | Few (critical paths) |
| Spring context | No | Partial/Full | Full |

---

## 2. JUnit 5, Mockito (@Mock, @Spy, @Captor)

### JUnit 5 Essentials

```java
@Test
void shouldCalculateTotal() { ... }

@Test
@DisplayName("Order total includes tax when taxable item")
void taxableItems() { ... }

// Lifecycle hooks
@BeforeAll  static void setupOnce() { ... }   // runs once before all tests in class
@AfterAll   static void teardownOnce() { ... }
@BeforeEach void setup() { ... }               // runs before each test
@AfterEach  void teardown() { ... }

// Skipping
@Test @Disabled("Fix in JIRA-456")
void brokenTest() { ... }

// Conditional
@Test @EnabledOnOs(OS.LINUX)
void linuxOnly() { ... }

@Test @EnabledIfEnvironmentVariable(named = "CI", matches = "true")
void ciOnly() { ... }
```

#### Parameterized Tests

```java
@ParameterizedTest
@ValueSource(ints = {1, 5, 10, 100})
void discountNeverExceeds50Percent(int quantity) {
    double discount = discountService.calculate(quantity);
    assertThat(discount).isLessThanOrEqualTo(0.5);
}

@ParameterizedTest
@CsvSource({
    "GOLD,   100, 10.0",
    "SILVER, 100,  5.0",
    "BRONZE, 100,  2.0"
})
void discountByTier(String tier, int amount, double expectedDiscount) {
    assertThat(discountService.calculate(tier, amount)).isEqualTo(expectedDiscount);
}

@ParameterizedTest
@MethodSource("orderProvider")
void processesVariousOrders(Order order, String expectedStatus) {
    assertThat(orderService.process(order).getStatus()).isEqualTo(expectedStatus);
}

static Stream<Arguments> orderProvider() {
    return Stream.of(
        Arguments.of(new Order("valid"),   "COMPLETED"),
        Arguments.of(new Order("invalid"), "FAILED")
    );
}
```

#### Assertions (AssertJ — recommended over JUnit assertions)

```java
// AssertJ — fluent, readable, better error messages
assertThat(result).isNotNull();
assertThat(list).hasSize(3).contains("apple").doesNotContain("grape");
assertThat(user.getName()).isEqualTo("Alice").startsWith("Al");
assertThat(price).isGreaterThan(BigDecimal.ZERO).isLessThan(BigDecimal.valueOf(1000));

// Exception assertions
assertThatThrownBy(() -> orderService.place(invalidOrder))
    .isInstanceOf(IllegalArgumentException.class)
    .hasMessageContaining("quantity must be positive");

// Soft assertions — all checked even if some fail
SoftAssertions softly = new SoftAssertions();
softly.assertThat(order.getId()).isNotNull();
softly.assertThat(order.getStatus()).isEqualTo("PENDING");
softly.assertThat(order.getTotal()).isPositive();
softly.assertAll();   // reports all failures at once
```

---

### Mockito — `@Mock`, `@Spy`, `@Captor`

#### `@Mock`

Creates a **full mock** — all methods return defaults (`null`, `0`, `false`, empty collections).  
You define behaviour explicitly via `when(...).thenReturn(...)`.

```java
@Mock
UserRepository userRepository;

@Test
void findUser() {
    User alice = new User(1L, "Alice");
    when(userRepository.findById(1L)).thenReturn(Optional.of(alice));

    Optional<User> result = userService.findUser(1L);

    assertThat(result).contains(alice);
    verify(userRepository).findById(1L);                    // verify call happened
    verify(userRepository, times(1)).findById(anyLong());   // with any arg
    verify(userRepository, never()).save(any());             // never called
}

// Stubbing exceptions
when(userRepository.findById(99L)).thenThrow(new EntityNotFoundException());

// Stubbing void methods
doNothing().when(emailService).send(any());
doThrow(new MailException()).when(emailService).send("bad@email");
```

---

#### `@Spy`

Creates a **partial mock** — wraps a **real object**. Real methods are called by default; you can override specific ones.

```java
@Spy
List<String> spyList = new ArrayList<>();

@Test
void spyCallsRealMethods() {
    spyList.add("hello");                        // real add() called
    spyList.add("world");

    assertThat(spyList).hasSize(2);              // real size() = 2

    doReturn(100).when(spyList).size();          // override just size()
    assertThat(spyList.size()).isEqualTo(100);   // stubbed

    verify(spyList, times(2)).add(anyString());  // verify real calls
}

// Real use case — spy on service with one method overridden
@Spy
OrderService orderService = new OrderService(realRepo);

// Override just the tax calculation; rest is real
doReturn(BigDecimal.ZERO).when(orderService).calculateTax(any());
```

**Mock vs Spy:**

| | `@Mock` | `@Spy` |
|---|---|---|
| Methods by default | Return defaults | Call real implementation |
| Object creation | Mockito creates | You provide real instance |
| Use when | Full isolation | Real behaviour + override one method |

---

#### `@Captor`

Captures **arguments passed to mock methods** for later assertion.

```java
@Captor
ArgumentCaptor<OrderEvent> eventCaptor;

@Test
void placeOrder_shouldPublishEvent() {
    orderService.placeOrder(new Order(1L, BigDecimal.valueOf(99.99)));

    verify(eventPublisher).publish(eventCaptor.capture());    // capture the arg

    OrderEvent publishedEvent = eventCaptor.getValue();
    assertThat(publishedEvent.getOrderId()).isEqualTo(1L);
    assertThat(publishedEvent.getAmount()).isEqualByComparingTo("99.99");
    assertThat(publishedEvent.getType()).isEqualTo("ORDER_PLACED");
}

// Multiple invocations — getAllValues()
verify(eventPublisher, times(3)).publish(eventCaptor.capture());
List<OrderEvent> allEvents = eventCaptor.getAllValues();
assertThat(allEvents).extracting(OrderEvent::getType)
    .containsExactly("ORDER_PLACED", "PAYMENT_CHARGED", "INVENTORY_RESERVED");
```

---

## 3. @SpringBootTest vs Slice Tests

### `@SpringBootTest`

Loads the **full** Spring application context — every bean, auto-configuration, embedded server.

```java
// Full context, default port (no server)
@SpringBootTest
class FullContextTest {
    @Autowired OrderService orderService;   // real bean
    @Autowired UserRepository userRepo;     // real bean

    @Test void fullIntegrationTest() { ... }
}

// Full context WITH embedded server on random port
@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
class FullServerTest {
    @LocalServerPort int port;
    @Autowired TestRestTemplate restTemplate;

    @Test
    void callRealEndpoint() {
        ResponseEntity<String> resp = restTemplate.getForEntity("/api/health", String.class);
        assertThat(resp.getStatusCode()).isEqualTo(HttpStatus.OK);
    }
}
```

**Problem**: Slow — loads everything. Use only when you need the full context.

---

### Slice Tests (Test Specific Layers)

Spring Boot provides focused slices that load only the relevant layer's beans.

#### `@WebMvcTest` — Controller Layer Only

Loads only controllers, filters, `@ControllerAdvice`, security. No `@Service` or `@Repository` beans.

```java
@WebMvcTest(OrderController.class)
class OrderControllerTest {

    @Autowired MockMvc mockMvc;

    @MockBean OrderService orderService;   // mock the service layer — not loaded by @WebMvcTest

    @Test
    void createOrder_returns201() throws Exception {
        Order saved = new Order(1L, "PENDING");
        when(orderService.create(any())).thenReturn(saved);

        mockMvc.perform(post("/api/orders")
                .contentType(MediaType.APPLICATION_JSON)
                .content("""{"productId": 1, "quantity": 2}"""))
            .andExpect(status().isCreated())
            .andExpect(jsonPath("$.id").value(1))
            .andExpect(jsonPath("$.status").value("PENDING"));
    }

    @Test
    void createOrder_withInvalidBody_returns400() throws Exception {
        mockMvc.perform(post("/api/orders")
                .contentType(MediaType.APPLICATION_JSON)
                .content("""{"quantity": -1}"""))   // invalid
            .andExpect(status().isBadRequest());
    }
}
```

---

#### `@DataJpaTest` — JPA / Repository Layer Only

Loads JPA repositories, entities, and an **in-memory database** (H2 by default). No service beans.

```java
@DataJpaTest
class OrderRepositoryTest {

    @Autowired OrderRepository orderRepository;
    @Autowired TestEntityManager em;           // JPA test helper

    @Test
    void findByStatus_shouldReturnMatchingOrders() {
        em.persist(new Order("PENDING"));
        em.persist(new Order("PENDING"));
        em.persist(new Order("COMPLETED"));
        em.flush();

        List<Order> pending = orderRepository.findByStatus("PENDING");

        assertThat(pending).hasSize(2);
    }

    @Test
    void findTop3ByOrderByCreatedAtDesc_shouldReturnLatest3() {
        for (int i = 0; i < 5; i++) em.persist(new Order("COMPLETED"));
        em.flush();

        List<Order> latest = orderRepository.findTop3ByOrderByCreatedAtDesc();

        assertThat(latest).hasSize(3);
    }
}
```

Use **real DB** with `@AutoConfigureTestDatabase(replace = NONE)`:

```java
@DataJpaTest
@AutoConfigureTestDatabase(replace = AutoConfigureTestDatabase.Replace.NONE)
@Testcontainers
class OrderRepositoryRealDbTest {

    @Container
    static PostgreSQLContainer<?> postgres = new PostgreSQLContainer<>("postgres:15");

    @DynamicPropertySource
    static void configureDb(DynamicPropertyRegistry registry) {
        registry.add("spring.datasource.url", postgres::getJdbcUrl);
        registry.add("spring.datasource.username", postgres::getUsername);
        registry.add("spring.datasource.password", postgres::getPassword);
    }

    @Autowired OrderRepository orderRepository;

    @Test
    void persistsAndFinds() { ... }
}
```

---

#### Other Slice Tests

| Annotation | What it loads |
|---|---|
| `@WebMvcTest` | Controllers, filters, security |
| `@DataJpaTest` | JPA repositories, entities, in-memory DB |
| `@DataMongoTest` | MongoDB repositories |
| `@DataRedisTest` | Redis repositories |
| `@WebFluxTest` | Reactive controllers (`WebTestClient`) |
| `@JsonTest` | JSON serialization / deserialization only |
| `@RestClientTest` | `RestTemplate`/`WebClient` with mock server |

### When to Use What

| Scenario | Use |
|---|---|
| Test controller request/response, validation, security | `@WebMvcTest` |
| Test repository queries against DB | `@DataJpaTest` |
| Test full request → service → DB flow | `@SpringBootTest` + `@Transactional` |
| Test one service class | Plain JUnit + `@ExtendWith(MockitoExtension.class)` |
| Test against real containers | `@SpringBootTest` + Testcontainers |

---

## 4. TDD Approach

### Red → Green → Refactor Cycle

```
  ┌──────────────────────────────────────────────┐
  │                                              │
  │   1. RED    Write a failing test first       │
  │      ↓      (doesn't compile or fails)       │
  │   2. GREEN  Write minimum code to pass test  │
  │      ↓      (may be ugly, just make it pass) │
  │   3. REFACTOR Clean up code while tests pass │
  │      ↓      (no new logic, only cleanup)     │
  │      └──────────────────────────────────────►│
  │              repeat                          │
  └──────────────────────────────────────────────┘
```

---

### TDD Example — Building a Shopping Cart

#### Step 1 — RED (write failing test first)

```java
@Test
void emptyCart_totalShouldBeZero() {
    ShoppingCart cart = new ShoppingCart();    // doesn't exist yet
    assertThat(cart.getTotal()).isEqualByComparingTo(BigDecimal.ZERO);
}
// ❌ Compile error — ShoppingCart doesn't exist
```

#### Step 2 — GREEN (minimal implementation)

```java
class ShoppingCart {
    public BigDecimal getTotal() {
        return BigDecimal.ZERO;
    }
}
// ✅ Test passes
```

#### Step 3 — Next RED

```java
@Test
void addItem_totalShouldInclude() {
    ShoppingCart cart = new ShoppingCart();
    cart.addItem(new Item("Apple", BigDecimal.valueOf(1.50), 3));

    assertThat(cart.getTotal()).isEqualByComparingTo("4.50");
}
// ❌ Fails — addItem not implemented
```

#### Step 4 — GREEN

```java
class ShoppingCart {
    private List<Item> items = new ArrayList<>();

    public void addItem(Item item) {
        items.add(item);
    }

    public BigDecimal getTotal() {
        return items.stream()
            .map(i -> i.getPrice().multiply(BigDecimal.valueOf(i.getQuantity())))
            .reduce(BigDecimal.ZERO, BigDecimal::add);
    }
}
// ✅ Both tests pass
```

#### Step 5 — RED (edge case)

```java
@Test
void addItem_withZeroQuantity_shouldThrow() {
    ShoppingCart cart = new ShoppingCart();
    assertThatThrownBy(() -> cart.addItem(new Item("Apple", BigDecimal.ONE, 0)))
        .isInstanceOf(IllegalArgumentException.class)
        .hasMessage("Quantity must be positive");
}
// ❌ Fails — no validation
```

#### Step 6 — GREEN

```java
public void addItem(Item item) {
    if (item.getQuantity() <= 0) {
        throw new IllegalArgumentException("Quantity must be positive");
    }
    items.add(item);
}
// ✅ All three tests pass
```

#### Step 7 — REFACTOR

```java
// Extract validation, improve readability — tests still pass
public void addItem(Item item) {
    validateItem(item);
    items.add(item);
}

private void validateItem(Item item) {
    if (item == null) throw new IllegalArgumentException("Item cannot be null");
    if (item.getQuantity() <= 0) throw new IllegalArgumentException("Quantity must be positive");
    if (item.getPrice().compareTo(BigDecimal.ZERO) < 0) throw new IllegalArgumentException("Price cannot be negative");
}
```

---

### Benefits of TDD

- **Design pressure**: hard-to-test code signals bad design (tight coupling, too many responsibilities)
- **Confidence**: refactor without fear — tests catch regressions immediately
- **Documentation**: tests describe expected behaviour better than comments
- **Fewer bugs**: edge cases considered upfront

### TDD Anti-Patterns to Avoid

```java
// ❌ Testing implementation details — brittle tests
verify(orderRepository, times(1)).save(any(Order.class));  // breaks on refactor
// ✅ Test observable behaviour instead
assertThat(orderService.getOrder(id)).isPresent();

// ❌ Over-mocking — mock everything including value objects
when(mockOrder.getId()).thenReturn(1L);
// ✅ Use real objects for simple data classes

// ❌ Too large test — tests multiple things
void testOrderService() {
    // 50 lines testing create, update, delete, and status change
}
// ✅ One concept per test — FIRST principle (Fast, Isolated, Repeatable, Self-validating, Timely)
```
