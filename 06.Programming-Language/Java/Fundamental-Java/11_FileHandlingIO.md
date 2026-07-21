# File Handling & I/O

---

## Table of Contents

1. [Byte Stream vs Character Stream](#1-byte-stream-vs-character-stream)
2. [File Class](#2-file-class)
3. [Reading & Writing Files](#3-reading--writing-files)
4. [BufferedReader & BufferedWriter](#4-bufferedreader--bufferedwriter)
5. [Scanner Class](#5-scanner-class)
6. [Serialization & Deserialization](#6-serialization--deserialization)
7. [Java NIO (New I/O)](#7-java-nio-new-io)

---

## 1. Byte Stream vs Character Stream

```
I/O Stream Hierarchy:
─────────────────────

Byte Streams (8-bit)                 Character Streams (16-bit Unicode)
────────────────────                 ─────────────────────────────────
InputStream (abstract)               Reader (abstract)
├── FileInputStream                  ├── FileReader
├── BufferedInputStream              ├── BufferedReader
├── ByteArrayInputStream             ├── InputStreamReader
├── ObjectInputStream                └── StringReader
└── DataInputStream

OutputStream (abstract)              Writer (abstract)
├── FileOutputStream                 ├── FileWriter
├── BufferedOutputStream             ├── BufferedWriter
├── ByteArrayOutputStream            ├── OutputStreamWriter
├── ObjectOutputStream               └── PrintWriter
└── DataOutputStream
```

```java
// ── BYTE STREAMS — read raw bytes (images, audio, binary data) ──
try (FileInputStream fis = new FileInputStream("image.png");
     FileOutputStream fos = new FileOutputStream("copy.png")) {

    byte[] buffer = new byte[4096];
    int bytesRead;
    while ((bytesRead = fis.read(buffer)) != -1) {
        fos.write(buffer, 0, bytesRead);   // write exactly bytesRead bytes
    }
    System.out.println("File copied");
}

// ── CHARACTER STREAMS — read/write text with encoding ──
try (FileReader fr = new FileReader("notes.txt");
     FileWriter fw = new FileWriter("output.txt")) {

    char[] buffer = new char[1024];
    int charsRead;
    while ((charsRead = fr.read(buffer)) != -1) {
        fw.write(buffer, 0, charsRead);
    }
}

// ── Bridging: InputStreamReader / OutputStreamWriter ──
// Convert byte stream to character stream with explicit encoding
try (InputStreamReader isr = new InputStreamReader(
        new FileInputStream("data.txt"), StandardCharsets.UTF_8)) {
    // read as UTF-8 characters
}
```

| | Byte Stream | Character Stream |
|---|---|---|
| Unit | 8-bit byte | 16-bit char (Unicode) |
| For | Binary data, images, audio | Text files |
| Base classes | `InputStream`/`OutputStream` | `Reader`/`Writer` |
| Encoding | No (raw bytes) | Yes (handles charset) |
| Default encoding | N/A | Platform default (risky — use explicit) |

> **Interview Q: What is the difference between byte stream and character stream?**  
> **Byte streams** handle raw 8-bit data — use them for binary files (images, audio, PDFs). **Character streams** handle 16-bit Unicode characters — use them for text files; they handle character encoding (UTF-8, UTF-16) transparently. For text, always use character streams with an explicit charset (`StandardCharsets.UTF_8`) to avoid platform-dependent encoding bugs.

---

## 2. File Class

```java
import java.io.File;

// ── Create File object (does NOT create the actual file) ──
File file = new File("/home/user/docs/notes.txt");
File dir  = new File("/home/user/docs");
File rel  = new File("notes.txt");   // relative to working directory

// ── File info ──
file.getName();           // "notes.txt"
file.getParent();         // "/home/user/docs"
file.getAbsolutePath();   // "/home/user/docs/notes.txt"
file.length();            // size in bytes
file.lastModified();      // timestamp in ms

// ── Checks ──
file.exists();        // does file exist?
file.isFile();        // is it a regular file?
file.isDirectory();   // is it a directory?
file.canRead();
file.canWrite();
file.canExecute();

// ── Creating files and directories ──
File newFile = new File("/tmp/test.txt");
newFile.createNewFile();         // creates file, returns false if already exists

File newDir = new File("/tmp/mydir");
newDir.mkdir();                  // creates single directory
newDir.mkdirs();                 // creates directory + all parents

// ── Listing directory contents ──
File docs = new File("/home/user/docs");
String[] names = docs.list();                        // names only
File[]   files = docs.listFiles();                   // File objects
File[]   txts  = docs.listFiles((d, name) -> name.endsWith(".txt"));  // filtered

// ── Operations ──
file.renameTo(new File("/home/user/docs/renamed.txt"));
file.delete();
file.deleteOnExit();   // deleted when JVM exits

// ── java.nio.file.Path (modern, preferred) ──
Path path = Paths.get("/home/user/docs/notes.txt");
path.getFileName();     // notes.txt
path.getParent();       // /home/user/docs
Files.exists(path);
Files.isDirectory(path);
Files.size(path);
```

> **Interview Q: What is the difference between `File` and `Path`/`Files` in Java?**  
> `java.io.File` is the legacy class — works but has limitations (poor error messages, unreliable return values from `delete()`/`mkdir()`). `java.nio.file.Path` (Java 7+) is the modern replacement — better API, throws exceptions instead of returning false, supports symbolic links, file attributes, and watch services. Use `Files.copy()`, `Files.move()`, `Files.delete()` etc. In new code, prefer `Path`/`Files`.

---

## 3. Reading & Writing Files

```java
import java.nio.file.*;
import java.nio.charset.StandardCharsets;

// ── Modern: Files utility class (Java 7+) ──

// Write all lines at once
List<String> lines = List.of("Line 1", "Line 2", "Line 3");
Files.write(Paths.get("output.txt"), lines, StandardCharsets.UTF_8);

// Append to file
Files.write(Paths.get("output.txt"), List.of("New line"),
    StandardCharsets.UTF_8,
    StandardOpenOption.APPEND);

// Read all lines at once
List<String> readLines = Files.readAllLines(Paths.get("input.txt"), StandardCharsets.UTF_8);
readLines.forEach(System.out::println);

// Read as single string
String content = Files.readString(Paths.get("input.txt"));  // Java 11+

// Write string
Files.writeString(Paths.get("output.txt"), "Hello, World!", StandardCharsets.UTF_8);  // Java 11+

// Stream lines (memory-efficient for large files)
try (Stream<String> stream = Files.lines(Paths.get("large.txt"), StandardCharsets.UTF_8)) {
    stream.filter(line -> line.contains("ERROR"))
          .forEach(System.out::println);
}

// ── Copy and Move ──
Files.copy(Paths.get("src.txt"), Paths.get("dst.txt"),
    StandardCopyOption.REPLACE_EXISTING);
Files.move(Paths.get("old.txt"), Paths.get("new.txt"),
    StandardCopyOption.REPLACE_EXISTING);
Files.delete(Paths.get("toDelete.txt"));           // throws if not found
Files.deleteIfExists(Paths.get("mayNotExist.txt")); // safe

// ── Walk directory tree ──
Files.walk(Paths.get("/home/user/docs"))
    .filter(p -> p.toString().endsWith(".java"))
    .forEach(System.out::println);
```

---

## 4. BufferedReader & BufferedWriter

Wraps a reader/writer with an **in-memory buffer** to reduce the number of I/O operations.

```java
import java.io.*;
import java.nio.charset.StandardCharsets;

// ── BUFFEREDREADER ──
try (BufferedReader br = new BufferedReader(
        new InputStreamReader(new FileInputStream("data.txt"), StandardCharsets.UTF_8))) {

    String line;
    int lineNum = 0;
    while ((line = br.readLine()) != null) {   // returns null at EOF
        lineNum++;
        System.out.println(lineNum + ": " + line);
    }
}

// Cleaner with Files.newBufferedReader():
try (BufferedReader br = Files.newBufferedReader(Paths.get("data.txt"), StandardCharsets.UTF_8)) {
    br.lines()                        // Java 8: Stream<String>
      .filter(l -> !l.isBlank())
      .forEach(System.out::println);
}

// ── BUFFEREDWRITER ──
try (BufferedWriter bw = Files.newBufferedWriter(Paths.get("output.txt"),
        StandardCharsets.UTF_8,
        StandardOpenOption.CREATE,
        StandardOpenOption.TRUNCATE_EXISTING)) {

    bw.write("First line");
    bw.newLine();                 // OS-appropriate line separator
    bw.write("Second line");
    bw.newLine();
    bw.flush();                   // optional when using try-with-resources (auto-closed)
}

// ── PrintWriter — adds print()/println()/printf() convenience ──
try (PrintWriter pw = new PrintWriter(new BufferedWriter(new FileWriter("report.txt")))) {
    pw.println("=== Report ===");
    pw.printf("%-20s %5d%n", "Alice", 95);
    pw.printf("%-20s %5d%n", "Bob", 87);
}
```

**Why buffering matters:**

```
Without buffer:
  1 write call → 1 OS disk write → 1 context switch
  1000 writes  → 1000 OS calls (very slow)

With buffer:
  1000 writes → accumulate in 8KB buffer → 1-2 OS disk writes → very fast

Default buffer size: 8192 characters for BufferedReader/BufferedWriter
```

> **Interview Q: Why do we use `BufferedReader` instead of `FileReader` directly?**  
> `FileReader` reads one character at a time, making one system call per character — very slow. `BufferedReader` reads a large chunk (default 8KB) at once into memory, then serves individual characters from that buffer. This drastically reduces system calls and I/O overhead. Also, `BufferedReader.readLine()` is convenient for line-by-line text processing. Always wrap `FileReader` in `BufferedReader` for text file reading.

---

## 5. Scanner Class

```java
import java.util.Scanner;

// ── Reading from CONSOLE ──
Scanner scanner = new Scanner(System.in);

System.out.print("Enter name: ");
String name = scanner.nextLine();    // reads entire line including spaces

System.out.print("Enter age: ");
int age = scanner.nextInt();
scanner.nextLine();                  // consume the leftover newline after nextInt()

System.out.print("Enter salary: ");
double salary = scanner.nextDouble();
scanner.nextLine();

System.out.printf("Hello, %s! Age: %d, Salary: %.2f%n", name, age, salary);
scanner.close();   // close when done with System.in? — debatable (closes System.in)

// ── Reading from FILE ──
try (Scanner fileScan = new Scanner(new File("data.txt"), StandardCharsets.UTF_8)) {
    while (fileScan.hasNextLine()) {
        String line = fileScan.nextLine();
        System.out.println(line);
    }
}

// ── Reading from STRING ──
Scanner strScan = new Scanner("42 3.14 hello true");
System.out.println(strScan.nextInt());      // 42
System.out.println(strScan.nextDouble());   // 3.14
System.out.println(strScan.next());         // "hello"
System.out.println(strScan.nextBoolean());  // true
strScan.close();

// ── Custom delimiter ──
Scanner csv = new Scanner("Alice,25,Engineer").useDelimiter(",");
System.out.println(csv.next());   // "Alice"
System.out.println(csv.nextInt());// 25
System.out.println(csv.next());   // "Engineer"
```

**Scanner vs BufferedReader:**

| | Scanner | BufferedReader |
|---|---|---|
| Tokenizing | ✅ nextInt, nextDouble, next | ❌ (read as String, parse manually) |
| Performance | Slower (regex tokenizing) | Faster (raw reads) |
| For | User input, small structured data | Large file reading |
| Thread-safe | ❌ | ✅ |

> **Interview Q: What is the common pitfall when mixing `nextInt()` and `nextLine()` in Scanner?**  
> After `nextInt()` (or `nextDouble()`, `nextBoolean()`), the Scanner's position is just after the number — the **newline character** from pressing Enter is still in the buffer. The next `nextLine()` immediately returns an empty string (consuming that leftover newline) instead of reading the next actual line. **Fix**: call `scanner.nextLine()` immediately after `nextInt()` to consume the newline, before reading the next full line.

---

## 6. Serialization & Deserialization

Serialization converts an object to a **byte stream** so it can be saved to disk or sent over a network. Deserialization reconstructs the object.

```java
import java.io.*;

// ── Class must implement Serializable ──
class Employee implements Serializable {
    private static final long serialVersionUID = 1L;  // version control

    String name;
    int salary;
    transient String password;         // excluded from serialization
    static int totalEmployees = 0;     // static — excluded (belongs to class)

    Employee(String name, int salary, String password) {
        this.name = name;
        this.salary = salary;
        this.password = password;
        totalEmployees++;
    }
}

// ── SERIALIZE: write object to file ──
Employee emp = new Employee("Alice", 75000, "secret123");
try (ObjectOutputStream oos = new ObjectOutputStream(
        new BufferedOutputStream(new FileOutputStream("employee.ser")))) {
    oos.writeObject(emp);
    System.out.println("Serialized");
}

// ── DESERIALIZE: read object from file ──
try (ObjectInputStream ois = new ObjectInputStream(
        new BufferedInputStream(new FileInputStream("employee.ser")))) {
    Employee restored = (Employee) ois.readObject();
    System.out.println(restored.name);      // "Alice"
    System.out.println(restored.salary);    // 75000
    System.out.println(restored.password);  // null (transient)
}

// ── serialVersionUID ──
// If you modify the class without updating serialVersionUID,
// deserializing old data will throw InvalidClassException
// Always declare it explicitly for production classes
```

**Custom serialization:**

```java
class SecureUser implements Serializable {
    private static final long serialVersionUID = 1L;
    String username;
    private String encryptedPassword;

    // Custom read — runs during deserialization instead of default
    private void readObject(ObjectInputStream ois) throws IOException, ClassNotFoundException {
        ois.defaultReadObject();   // read all non-transient fields
        // post-processing after deserialization
        validate();
    }

    // Custom write — runs during serialization instead of default
    private void writeObject(ObjectOutputStream oos) throws IOException {
        oos.defaultWriteObject();  // write all non-transient fields
        // extra data
    }
}
```

> **Interview Q: What is `serialVersionUID` and why is it important?**  
> `serialVersionUID` is a **version identifier** for a serializable class. Java uses it to verify that a serialized object is compatible with the current class definition when deserializing. If you modify the class (add/remove fields) without updating `serialVersionUID`, Java throws `InvalidClassException` when trying to deserialize old data. By declaring it explicitly (e.g., `1L`), you control compatibility — you can choose when breaking changes require a new UID.

---

## 7. Java NIO (New I/O)

Java NIO (`java.nio`) provides non-blocking, channel-based I/O — more efficient for high-throughput and scalable applications.

```java
import java.nio.*;
import java.nio.file.*;
import java.nio.channels.*;

// ── FILES UTILITY (NIO convenience methods, covered above) ──
// Files.readAllLines(), Files.write(), Files.copy(), etc.

// ── CHANNEL + BUFFER — low-level, efficient ──
// Write to file using channel
try (FileChannel fc = FileChannel.open(Paths.get("out.bin"),
        StandardOpenOption.CREATE, StandardOpenOption.WRITE)) {

    ByteBuffer buffer = ByteBuffer.allocate(64);
    buffer.put("Hello NIO".getBytes(StandardCharsets.UTF_8));
    buffer.flip();          // switch from write mode to read mode
    fc.write(buffer);
}

// Read from file using channel
try (FileChannel fc = FileChannel.open(Paths.get("out.bin"), StandardOpenOption.READ)) {
    ByteBuffer buffer = ByteBuffer.allocate(64);
    int bytesRead = fc.read(buffer);
    buffer.flip();
    byte[] data = new byte[bytesRead];
    buffer.get(data);
    System.out.println(new String(data, StandardCharsets.UTF_8));  // "Hello NIO"
}

// ── MEMORY-MAPPED FILE — fastest for large files ──
try (FileChannel fc = FileChannel.open(Paths.get("largefile.dat"), StandardOpenOption.READ);
     MappedByteBuffer mbb = fc.map(FileChannel.MapMode.READ_ONLY, 0, fc.size())) {

    // File is mapped directly to memory — OS handles paging
    // Access as if it's a byte array in memory
    byte firstByte = mbb.get(0);
}

// ── WatchService — monitor directory for changes ──
WatchService watcher = FileSystems.getDefault().newWatchService();
Path dir = Paths.get("/tmp/watched");
dir.register(watcher, StandardWatchEventKinds.ENTRY_CREATE,
                      StandardWatchEventKinds.ENTRY_MODIFY,
                      StandardWatchEventKinds.ENTRY_DELETE);

WatchKey key = watcher.take();   // blocks until event
for (WatchEvent<?> event : key.pollEvents()) {
    System.out.println(event.kind() + ": " + event.context());
}
key.reset();
```

> **Interview Q: What are the main differences between Java I/O and Java NIO?**  
> Java I/O is **stream-based** (one-directional), **blocking** (thread waits for data), and uses `InputStream`/`OutputStream`. Java NIO is **buffer-based** (bidirectional `ByteBuffer`), can be **non-blocking** (one thread manages multiple channels), and uses `Channel`/`Buffer`. NIO is better for **high-concurrency servers** (handling thousands of connections without a thread per connection). For simple file operations, NIO's `Files` utility class is cleaner and preferred over the old `java.io.File` API.
