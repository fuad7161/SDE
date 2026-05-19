# 🟢 Structural Patterns

> **Category:** Structural &nbsp;|&nbsp; **Tags:** `Adapter` `Decorator` `Facade` `Proxy` `Composite` `Bridge`

Structural patterns deal with **object composition** — how classes and objects are combined to form larger structures while keeping them flexible and efficient.

---

## Table of Contents
1. [Adapter](#1-adapter)
2. [Decorator](#2-decorator)
3. [Facade](#3-facade)
4. [Proxy](#4-proxy)
5. [Composite](#5-composite)
6. [Bridge](#6-bridge)
7. [Interview Questions](#interview-questions)

---

## 1. Adapter

**Intent:** Convert the interface of a class into another interface that clients expect. Allows classes with **incompatible interfaces** to work together.

**Analogy:** A power plug adapter that converts a US plug to a European socket.

**When to use:**
- Integrating a third-party library with an incompatible interface
- Reusing existing classes that can't be changed
- Legacy system integration

<details>
<summary><b>Java</b></summary>

```java
// Existing interface the client depends on
public interface MediaPlayer {
    void play(String fileName);
}

// New third-party library with a different interface
public class VLCPlayer {
    public void playVLC(String fileName) {
        System.out.println("VLC playing: " + fileName);
    }
}

public class MP4Player {
    public void playMP4(String fileName) {
        System.out.println("MP4 player playing: " + fileName);
    }
}

// Adapter — wraps the incompatible class, exposes MediaPlayer interface
public class MediaAdapter implements MediaPlayer {
    private final VLCPlayer vlcPlayer;
    private final MP4Player mp4Player;

    public MediaAdapter() {
        this.vlcPlayer = new VLCPlayer();
        this.mp4Player = new MP4Player();
    }

    @Override
    public void play(String fileName) {
        if (fileName.endsWith(".vlc")) {
            vlcPlayer.playVLC(fileName);
        } else if (fileName.endsWith(".mp4")) {
            mp4Player.playMP4(fileName);
        } else {
            throw new UnsupportedOperationException("Format not supported: " + fileName);
        }
    }
}

// Default media player (plays mp3 natively)
public class AudioPlayer implements MediaPlayer {
    private final MediaAdapter adapter = new MediaAdapter();

    @Override
    public void play(String fileName) {
        if (fileName.endsWith(".mp3")) {
            System.out.println("Native MP3: " + fileName);
        } else {
            adapter.play(fileName);    // delegate to adapter
        }
    }
}

// Client
public class Main {
    public static void main(String[] args) {
        MediaPlayer player = new AudioPlayer();
        player.play("song.mp3");       // Native MP3: song.mp3
        player.play("video.vlc");      // VLC playing: video.vlc
        player.play("movie.mp4");      // MP4 player playing: movie.mp4
    }
}
```

</details>

<details>
<summary><b>Go</b></summary>

```go
package main

import (
    "fmt"
    "strings"
)

// Target interface
type MediaPlayer interface {
    Play(fileName string)
}

// Adaptees (third-party — can't change)
type VLCPlayer struct{}
type MP4Player struct{}

func (v *VLCPlayer) PlayVLC(fileName string) { fmt.Println("VLC playing:", fileName) }
func (m *MP4Player) PlayMP4(fileName string) { fmt.Println("MP4 player playing:", fileName) }

// Adapter — wraps adaptees, implements MediaPlayer
type MediaAdapter struct {
    vlc *VLCPlayer
    mp4 *MP4Player
}

func NewMediaAdapter() *MediaAdapter {
    return &MediaAdapter{vlc: &VLCPlayer{}, mp4: &MP4Player{}}
}

func (a *MediaAdapter) Play(fileName string) {
    switch {
    case strings.HasSuffix(fileName, ".vlc"):
        a.vlc.PlayVLC(fileName)
    case strings.HasSuffix(fileName, ".mp4"):
        a.mp4.PlayMP4(fileName)
    default:
        panic("format not supported: " + fileName)
    }
}

type AudioPlayer struct {
    adapter *MediaAdapter
}

func NewAudioPlayer() *AudioPlayer {
    return &AudioPlayer{adapter: NewMediaAdapter()}
}

func (p *AudioPlayer) Play(fileName string) {
    if strings.HasSuffix(fileName, ".mp3") {
        fmt.Println("Native MP3:", fileName)
    } else {
        p.adapter.Play(fileName)
    }
}

func main() {
    player := NewAudioPlayer()
    player.Play("song.mp3")  // Native MP3: song.mp3
    player.Play("video.vlc") // VLC playing: video.vlc
    player.Play("movie.mp4") // MP4 player playing: movie.mp4
}
```

</details>

### Class Adapter vs Object Adapter

| | Class Adapter | Object Adapter |
|--|--------------|---------------|
| Mechanism | Extends adaptee (inheritance) | Wraps adaptee (composition) |
| Flexibility | Less — tied to one class | More — can wrap any subtype |
| Java support | Multiple inheritance needed | ✅ Preferred in Java |

---

## 2. Decorator

**Intent:** Attach additional responsibilities to an object **dynamically**. Decorators provide a flexible alternative to subclassing for extending functionality.

**Analogy:** A coffee order — start with plain coffee, then "decorate" with milk, sugar, whip cream.

**When to use:**
- Add behavior without changing the original class
- Combine behaviors flexibly at runtime
- Inheritance leads to class explosion

<details>
<summary><b>Java</b></summary>

```java
// Component interface
public interface Coffee {
    String getDescription();
    double getCost();
}

// Concrete Component
public class SimpleCoffee implements Coffee {
    @Override public String getDescription() { return "Simple coffee"; }
    @Override public double getCost()        { return 1.00; }
}

// Base Decorator — implements the same interface, wraps a component
public abstract class CoffeeDecorator implements Coffee {
    protected final Coffee coffee;   // wrapped component

    public CoffeeDecorator(Coffee coffee) {
        this.coffee = coffee;
    }

    @Override public String getDescription() { return coffee.getDescription(); }
    @Override public double getCost()        { return coffee.getCost(); }
}

// Concrete Decorators
public class MilkDecorator extends CoffeeDecorator {
    public MilkDecorator(Coffee coffee) { super(coffee); }

    @Override public String getDescription() { return coffee.getDescription() + ", Milk"; }
    @Override public double getCost()        { return coffee.getCost() + 0.25; }
}

public class SugarDecorator extends CoffeeDecorator {
    public SugarDecorator(Coffee coffee) { super(coffee); }

    @Override public String getDescription() { return coffee.getDescription() + ", Sugar"; }
    @Override public double getCost()        { return coffee.getCost() + 0.10; }
}

public class WhipDecorator extends CoffeeDecorator {
    public WhipDecorator(Coffee coffee) { super(coffee); }

    @Override public String getDescription() { return coffee.getDescription() + ", Whip"; }
    @Override public double getCost()        { return coffee.getCost() + 0.50; }
}

// Client
public class Main {
    public static void main(String[] args) {
        Coffee order = new SimpleCoffee();
        order = new MilkDecorator(order);
        order = new SugarDecorator(order);
        order = new WhipDecorator(order);

        System.out.println(order.getDescription());  // Simple coffee, Milk, Sugar, Whip
        System.out.println("$" + order.getCost());   // $1.85
    }
}
```

</details>

<details>
<summary><b>Go</b></summary>

```go
package main

import "fmt"

type Coffee interface {
    GetDescription() string
    GetCost() float64
}

type SimpleCoffee struct{}

func (c *SimpleCoffee) GetDescription() string { return "Simple coffee" }
func (c *SimpleCoffee) GetCost() float64       { return 1.00 }

// Base decorator
type CoffeeDecorator struct{ coffee Coffee }

func (d *CoffeeDecorator) GetDescription() string { return d.coffee.GetDescription() }
func (d *CoffeeDecorator) GetCost() float64       { return d.coffee.GetCost() }

// Concrete decorators
type MilkDecorator  struct{ CoffeeDecorator }
type SugarDecorator struct{ CoffeeDecorator }
type WhipDecorator  struct{ CoffeeDecorator }

func NewMilk(c Coffee) *MilkDecorator   { return &MilkDecorator{CoffeeDecorator{c}} }
func NewSugar(c Coffee) *SugarDecorator { return &SugarDecorator{CoffeeDecorator{c}} }
func NewWhip(c Coffee) *WhipDecorator   { return &WhipDecorator{CoffeeDecorator{c}} }

func (m *MilkDecorator) GetDescription() string  { return m.coffee.GetDescription() + ", Milk" }
func (m *MilkDecorator) GetCost() float64        { return m.coffee.GetCost() + 0.25 }
func (s *SugarDecorator) GetDescription() string { return s.coffee.GetDescription() + ", Sugar" }
func (s *SugarDecorator) GetCost() float64       { return s.coffee.GetCost() + 0.10 }
func (w *WhipDecorator) GetDescription() string  { return w.coffee.GetDescription() + ", Whip" }
func (w *WhipDecorator) GetCost() float64        { return w.coffee.GetCost() + 0.50 }

func main() {
    var order Coffee = &SimpleCoffee{}
    order = NewMilk(order)
    order = NewSugar(order)
    order = NewWhip(order)

    fmt.Println(order.GetDescription())       // Simple coffee, Milk, Sugar, Whip
    fmt.Printf("$%.2f\n", order.GetCost())   // $1.85
}
```

</details>

> Java I/O is built on Decorator: `new BufferedReader(new InputStreamReader(new FileInputStream("file.txt")))`.

### Decorator vs Inheritance

| | Inheritance | Decorator |
|--|------------|-----------|
| Behavior added at | Compile time | Runtime |
| Combinations possible | N subclasses for N combos | Infinite combos from few classes |
| Open/Closed principle | Violates (modify class) | Follows (wrap without changing) |

---

## 3. Facade

**Intent:** Provide a **simplified interface** to a complex subsystem. Doesn't hide the subsystem — just provides a convenient front door.

**Analogy:** A home theater remote — one button starts the projector, dims lights, turns on sound, and starts the movie.

**When to use:**
- Simplify complex library or framework usage
- Provide a clean API over a legacy system
- Reduce client dependencies on internal subsystems

<details>
<summary><b>Java</b></summary>

```java
// Complex subsystem classes
public class Amplifier {
    public void on()               { System.out.println("Amplifier on"); }
    public void setVolume(int v)   { System.out.println("Volume: " + v); }
    public void off()              { System.out.println("Amplifier off"); }
}

public class Projector {
    public void on()               { System.out.println("Projector on"); }
    public void setInput(String s) { System.out.println("Input: " + s); }
    public void off()              { System.out.println("Projector off"); }
}

public class Lights {
    public void dim(int level)     { System.out.println("Lights dimmed to " + level + "%"); }
    public void on()               { System.out.println("Lights on"); }
}

public class DVDPlayer {
    public void on()               { System.out.println("DVD on"); }
    public void play(String movie) { System.out.println("Playing: " + movie); }
    public void stop()             { System.out.println("DVD stopped"); }
    public void off()              { System.out.println("DVD off"); }
}

// Facade — simplified interface over all the subsystems
public class HomeTheaterFacade {
    private final Amplifier amp;
    private final Projector projector;
    private final Lights lights;
    private final DVDPlayer dvd;

    public HomeTheaterFacade(Amplifier amp, Projector projector,
                              Lights lights, DVDPlayer dvd) {
        this.amp       = amp;
        this.projector = projector;
        this.lights    = lights;
        this.dvd       = dvd;
    }

    public void watchMovie(String movie) {
        System.out.println("--- Watch Movie ---");
        lights.dim(10);
        amp.on();
        amp.setVolume(8);
        projector.on();
        projector.setInput("DVD");
        dvd.on();
        dvd.play(movie);
    }

    public void endMovie() {
        System.out.println("--- Shutting Down ---");
        dvd.stop();
        dvd.off();
        projector.off();
        amp.off();
        lights.on();
    }
}

// Client — only interacts with the Facade
public class Main {
    public static void main(String[] args) {
        HomeTheaterFacade theater = new HomeTheaterFacade(
            new Amplifier(), new Projector(), new Lights(), new DVDPlayer()
        );
        theater.watchMovie("Inception");
        // ...
        theater.endMovie();
    }
}
```

</details>

<details>
<summary><b>Go</b></summary>

```go
package main

import "fmt"

type Amplifier struct{}
type Projector struct{}
type Lights    struct{}
type DVDPlayer struct{}

func (a *Amplifier) On()               { fmt.Println("Amplifier on") }
func (a *Amplifier) SetVolume(v int)   { fmt.Println("Volume:", v) }
func (a *Amplifier) Off()              { fmt.Println("Amplifier off") }
func (p *Projector) On()               { fmt.Println("Projector on") }
func (p *Projector) SetInput(s string) { fmt.Println("Input:", s) }
func (p *Projector) Off()              { fmt.Println("Projector off") }
func (l *Lights) Dim(level int)        { fmt.Printf("Lights dimmed to %d%%\n", level) }
func (l *Lights) On()                  { fmt.Println("Lights on") }
func (d *DVDPlayer) On()               { fmt.Println("DVD on") }
func (d *DVDPlayer) Play(m string)     { fmt.Println("Playing:", m) }
func (d *DVDPlayer) Stop()             { fmt.Println("DVD stopped") }
func (d *DVDPlayer) Off()              { fmt.Println("DVD off") }

// Facade
type HomeTheaterFacade struct {
    amp       *Amplifier
    projector *Projector
    lights    *Lights
    dvd       *DVDPlayer
}

func NewHomeTheater() *HomeTheaterFacade {
    return &HomeTheaterFacade{
        amp:       &Amplifier{},
        projector: &Projector{},
        lights:    &Lights{},
        dvd:       &DVDPlayer{},
    }
}

func (h *HomeTheaterFacade) WatchMovie(movie string) {
    fmt.Println("--- Watch Movie ---")
    h.lights.Dim(10)
    h.amp.On()
    h.amp.SetVolume(8)
    h.projector.On()
    h.projector.SetInput("DVD")
    h.dvd.On()
    h.dvd.Play(movie)
}

func (h *HomeTheaterFacade) EndMovie() {
    fmt.Println("--- Shutting Down ---")
    h.dvd.Stop()
    h.dvd.Off()
    h.projector.Off()
    h.amp.Off()
    h.lights.On()
}

func main() {
    theater := NewHomeTheater()
    theater.WatchMovie("Inception")
    theater.EndMovie()
}
```

</details>

---

## 4. Proxy

**Intent:** Provide a **surrogate or placeholder** for another object to control access to it.

**Three main types:**
- **Virtual Proxy:** Defers expensive object creation until needed (lazy init)
- **Protection Proxy:** Controls access based on permissions
- **Remote Proxy:** Represents a remote object (e.g., RMI, REST)

### Virtual Proxy — Lazy Image Loading

<details>
<summary><b>Java</b></summary>

```java
public interface Image {
    void display();
}

// Real subject — expensive to create (loads from disk)
public class RealImage implements Image {
    private final String filename;

    public RealImage(String filename) {
        this.filename = filename;
        loadFromDisk();
    }

    private void loadFromDisk() {
        System.out.println("Loading image: " + filename);  // expensive!
    }

    @Override
    public void display() {
        System.out.println("Displaying: " + filename);
    }
}

// Proxy — defers creation of RealImage until display() is called
public class ImageProxy implements Image {
    private final String filename;
    private RealImage realImage;   // null until first use

    public ImageProxy(String filename) {
        this.filename = filename;
    }

    @Override
    public void display() {
        if (realImage == null) {
            realImage = new RealImage(filename);  // create on first call
        }
        realImage.display();
    }
}

// Usage
Image image = new ImageProxy("photo.jpg");
// Image NOT loaded yet
image.display();  // Loading image: photo.jpg  → Displaying: photo.jpg
image.display();  // Displaying: photo.jpg  (no reload)
```

</details>

<details>
<summary><b>Go</b></summary>

```go
package main

import "fmt"

type Image interface {
    Display()
}

type RealImage struct{ filename string }

func NewRealImage(filename string) *RealImage {
    img := &RealImage{filename: filename}
    fmt.Println("Loading image:", filename) // expensive!
    return img
}

func (r *RealImage) Display() { fmt.Println("Displaying:", r.filename) }

// Proxy — lazy loading
type ImageProxy struct {
    filename  string
    realImage *RealImage
}

func NewImageProxy(filename string) *ImageProxy {
    return &ImageProxy{filename: filename}
}

func (p *ImageProxy) Display() {
    if p.realImage == nil {
        p.realImage = NewRealImage(p.filename) // create on first call
    }
    p.realImage.Display()
}

func main() {
    image := NewImageProxy("photo.jpg")
    // Image NOT loaded yet
    image.Display() // Loading image: photo.jpg → Displaying: photo.jpg
    image.Display() // Displaying: photo.jpg (no reload)
}
```

</details>

### Protection Proxy — Role-based Access

<details>
<summary><b>Java</b></summary>

```java
public interface DocumentService {
    String read(String docId);
    void   write(String docId, String content);
}

public class RealDocumentService implements DocumentService {
    @Override public String read(String docId)                   { return "Content of " + docId; }
    @Override public void write(String docId, String content)    { System.out.println("Saved: " + docId); }
}

public class DocumentProxy implements DocumentService {
    private final RealDocumentService service = new RealDocumentService();
    private final String userRole;

    public DocumentProxy(String userRole) { this.userRole = userRole; }

    @Override
    public String read(String docId) {
        return service.read(docId);  // everyone can read
    }

    @Override
    public void write(String docId, String content) {
        if (!"ADMIN".equals(userRole)) {
            throw new SecurityException("Write access denied for role: " + userRole);
        }
        service.write(docId, content);
    }
}

// Usage
DocumentService service = new DocumentProxy("USER");
service.read("doc1");           // OK
service.write("doc1", "data");  // throws SecurityException
```

</details>

<details>
<summary><b>Go</b></summary>

```go
package main

import (
    "errors"
    "fmt"
)

type DocumentService interface {
    Read(docID string) string
    Write(docID, content string) error
}

type RealDocumentService struct{}

func (s *RealDocumentService) Read(docID string) string {
    return "Content of " + docID
}
func (s *RealDocumentService) Write(docID, content string) error {
    fmt.Println("Saved:", docID)
    return nil
}

type DocumentProxy struct {
    service  *RealDocumentService
    userRole string
}

func NewDocumentProxy(role string) *DocumentProxy {
    return &DocumentProxy{service: &RealDocumentService{}, userRole: role}
}

func (p *DocumentProxy) Read(docID string) string { return p.service.Read(docID) }

func (p *DocumentProxy) Write(docID, content string) error {
    if p.userRole != "ADMIN" {
        return errors.New("write access denied for role: " + p.userRole)
    }
    return p.service.Write(docID, content)
}

func main() {
    service := NewDocumentProxy("USER")
    fmt.Println(service.Read("doc1"))           // Content of doc1
    if err := service.Write("doc1", "data"); err != nil {
        fmt.Println("Error:", err)              // Error: write access denied for role: USER
    }
}
```

</details>

---

## 5. Composite

**Intent:** Compose objects into **tree structures** to represent part-whole hierarchies. Lets clients treat individual objects (leaves) and compositions (nodes) uniformly.

**Analogy:** File system — files and directories are both "file system items". A directory contains items; a file is a leaf item.

**When to use:**
- Tree structures (file system, org chart, UI components, HTML DOM)
- Clients should ignore the difference between single objects and compositions

<details>
<summary><b>Java</b></summary>

```java
// Component — common interface for leaves and composites
public abstract class FileSystemItem {
    protected final String name;

    public FileSystemItem(String name) { this.name = name; }
    public String getName()            { return name; }

    public abstract long getSize();
    public abstract void print(String indent);

    // Optional — overridden only by Directory
    public void add(FileSystemItem item)    { throw new UnsupportedOperationException(); }
    public void remove(FileSystemItem item) { throw new UnsupportedOperationException(); }
}

// Leaf
public class File extends FileSystemItem {
    private final long size;

    public File(String name, long size) {
        super(name);
        this.size = size;
    }

    @Override public long getSize() { return size; }
    @Override public void print(String indent) {
        System.out.println(indent + "📄 " + name + " (" + size + " bytes)");
    }
}

// Composite
public class Directory extends FileSystemItem {
    private final List<FileSystemItem> children = new ArrayList<>();

    public Directory(String name) { super(name); }

    @Override public void add(FileSystemItem item)    { children.add(item); }
    @Override public void remove(FileSystemItem item) { children.remove(item); }

    @Override
    public long getSize() {
        return children.stream().mapToLong(FileSystemItem::getSize).sum();
    }

    @Override
    public void print(String indent) {
        System.out.println(indent + "📁 " + name + " (" + getSize() + " bytes)");
        children.forEach(child -> child.print(indent + "  "));
    }
}

// Client
public class Main {
    public static void main(String[] args) {
        Directory root = new Directory("root");

        Directory src = new Directory("src");
        src.add(new File("Main.java", 2048));
        src.add(new File("App.java",  1024));

        Directory test = new Directory("test");
        test.add(new File("MainTest.java", 512));

        root.add(src);
        root.add(test);
        root.add(new File("README.md", 256));

        root.print("");
        // 📁 root (3840 bytes)
        //   📁 src (3072 bytes)
        //     📄 Main.java (2048 bytes)
        //     📄 App.java (1024 bytes)
        //   📁 test (512 bytes)
        //     📄 MainTest.java (512 bytes)
        //   📄 README.md (256 bytes)
    }
}
```

</details>

<details>
<summary><b>Go</b></summary>

```go
package main

import "fmt"

type FileSystemItem interface {
    GetName() string
    GetSize() int64
    Print(indent string)
}

// Leaf
type File struct {
    name string
    size int64
}

func (f *File) GetName() string     { return f.name }
func (f *File) GetSize() int64      { return f.size }
func (f *File) Print(indent string) {
    fmt.Printf("%s📄 %s (%d bytes)\n", indent, f.name, f.size)
}

// Composite
type Directory struct {
    name     string
    children []FileSystemItem
}

func NewDirectory(name string) *Directory { return &Directory{name: name} }

func (d *Directory) Add(item FileSystemItem) { d.children = append(d.children, item) }
func (d *Directory) GetName() string        { return d.name }
func (d *Directory) GetSize() int64 {
    var total int64
    for _, child := range d.children {
        total += child.GetSize()
    }
    return total
}
func (d *Directory) Print(indent string) {
    fmt.Printf("%s📁 %s (%d bytes)\n", indent, d.name, d.GetSize())
    for _, child := range d.children {
        child.Print(indent + "  ")
    }
}

func main() {
    root := NewDirectory("root")
    src := NewDirectory("src")
    src.Add(&File{"Main.java", 2048})
    src.Add(&File{"App.java", 1024})
    test := NewDirectory("test")
    test.Add(&File{"MainTest.java", 512})
    root.Add(src)
    root.Add(test)
    root.Add(&File{"README.md", 256})
    root.Print("")
}
```

</details>

---

## 6. Bridge

**Intent:** **Decouple an abstraction from its implementation** so the two can vary independently. Uses composition over inheritance.

**When to use:**
- Avoid permanent binding between abstraction and implementation
- Both abstraction and implementation should be extensible via subclassing
- Implementation details should be hidden from the client

<details>
<summary><b>Java</b></summary>

```java
// Implementor interface
public interface MessageSender {
    void sendMessage(String to, String body);
}

// Concrete Implementors
public class EmailSender implements MessageSender {
    @Override
    public void sendMessage(String to, String body) {
        System.out.println("EMAIL to " + to + ": " + body);
    }
}

public class SMSSender implements MessageSender {
    @Override
    public void sendMessage(String to, String body) {
        System.out.println("SMS to " + to + ": " + body);
    }
}

// Abstraction — holds reference to implementor (the "bridge")
public abstract class Notification {
    protected final MessageSender sender;   // bridge to implementation

    public Notification(MessageSender sender) {
        this.sender = sender;
    }

    public abstract void send(String recipient, String message);
}

// Refined Abstractions
public class UrgentNotification extends Notification {
    public UrgentNotification(MessageSender sender) { super(sender); }

    @Override
    public void send(String recipient, String message) {
        sender.sendMessage(recipient, "[URGENT] " + message);
    }
}

public class ScheduledNotification extends Notification {
    private final String sendTime;

    public ScheduledNotification(MessageSender sender, String sendTime) {
        super(sender);
        this.sendTime = sendTime;
    }

    @Override
    public void send(String recipient, String message) {
        sender.sendMessage(recipient, "[Scheduled " + sendTime + "] " + message);
    }
}

// Client — mix and match abstraction + implementation independently
public class Main {
    public static void main(String[] args) {
        Notification urgent = new UrgentNotification(new SMSSender());
        urgent.send("+1234567890", "Server is down!");
        // SMS to +1234567890: [URGENT] Server is down!

        Notification scheduled = new ScheduledNotification(new EmailSender(), "09:00");
        scheduled.send("alice@example.com", "Weekly report ready");
        // EMAIL to alice@example.com: [Scheduled 09:00] Weekly report ready
    }
}
```

</details>

<details>
<summary><b>Go</b></summary>

```go
package main

import "fmt"

// Implementor
type MessageSender interface {
    SendMessage(to, body string)
}

type EmailSender struct{}
type SMSSender  struct{}

func (e *EmailSender) SendMessage(to, body string) {
    fmt.Printf("EMAIL to %s: %s\n", to, body)
}
func (s *SMSSender) SendMessage(to, body string) {
    fmt.Printf("SMS to %s: %s\n", to, body)
}

// Abstraction (bridge holds the implementor)
type NotificationBase struct {
    sender MessageSender
}

// Refined Abstractions
type UrgentNotification struct{ NotificationBase }
type ScheduledNotification struct {
    NotificationBase
    sendTime string
}

func (u *UrgentNotification) Send(recipient, message string) {
    u.sender.SendMessage(recipient, "[URGENT] "+message)
}

func (s *ScheduledNotification) Send(recipient, message string) {
    s.sender.SendMessage(recipient, "[Scheduled "+s.sendTime+"] "+message)
}

func main() {
    urgent := &UrgentNotification{NotificationBase{sender: &SMSSender{}}}
    urgent.Send("+1234567890", "Server is down!")
    // SMS to +1234567890: [URGENT] Server is down!

    scheduled := &ScheduledNotification{
        NotificationBase: NotificationBase{sender: &EmailSender{}},
        sendTime:         "09:00",
    }
    scheduled.Send("alice@example.com", "Weekly report ready")
    // EMAIL to alice@example.com: [Scheduled 09:00] Weekly report ready
}
```

</details>

### Bridge vs Adapter

| | Bridge | Adapter |
|--|--------|---------|
| Intent | Design from the start to allow variation | Fix incompatibility after the fact |
| Used when | New system, planned extension | Integrating existing / legacy code |
| Direction | Proactive (planned) | Reactive (retrofit) |

---

## Interview Questions

### Q1. What is the Adapter pattern? Give a real-world Java example.

> **Answer:**
> Adapter converts one interface into another expected by the client — like a plug adapter. In Java:
> - `Arrays.asList()` adapts an array to the `List` interface.
> - `InputStreamReader` adapts a byte-stream (`InputStream`) to a char-stream (`Reader`).
> - `Collections.enumeration(list)` adapts a `List` to the legacy `Enumeration` interface.
>
> Prefer **object adapter** (composition) over class adapter (inheritance) in Java since Java doesn't support multiple inheritance.

---

### Q2. How does the Decorator pattern differ from inheritance?

> **Answer:**
> - **Inheritance:** Behavior added at compile time; for N combinations of N behaviors you need N subclasses.
> - **Decorator:** Behavior added at runtime by wrapping; N decorators produce unlimited combinations with N classes. Follows the Open/Closed principle — add behavior by writing new decorators, not modifying existing code.
>
> Java I/O is the classic example: `BufferedInputStream`, `DataInputStream`, `GZIPInputStream` are all decorators over `InputStream`.

---

### Q3. What is the Proxy pattern and what are its types?

> **Answer:**
> A Proxy is a placeholder that controls access to another object.
>
> - **Virtual Proxy:** Defers expensive creation until first use (lazy initialization). e.g., Hibernate's lazy-loaded entities.
> - **Protection Proxy:** Checks permissions before delegating to the real subject.
> - **Remote Proxy:** Represents an object in another process or machine (RMI, gRPC stub).
> - **Caching Proxy:** Caches results of expensive operations.
> - **Logging Proxy:** Logs calls before/after delegating.
>
> Spring AOP creates proxies around beans to implement `@Transactional`, `@Cacheable`, etc.

---

### Q4. When would you use Facade vs Adapter?

> **Answer:**
> - **Facade:** Simplifies a complex subsystem by providing a high-level interface. Doesn't change interfaces — just provides a convenient wrapper over multiple classes.
> - **Adapter:** Makes an incompatible interface compatible with what the client expects. Changes how an interface looks.
>
> Use Facade to simplify; use Adapter to make incompatible things work together.

---

### Q5. Explain the Composite pattern with a real-world example.

> **Answer:**
> Composite lets you treat single objects and groups of objects uniformly using a common interface.
>
> Classic examples:
> - **File system:** Both `File` and `Directory` implement `FileSystemItem`. You can call `getSize()` on either — directory sums its children recursively.
> - **UI components:** A `Panel` and a `Button` both implement `render()`. A panel renders by recursively rendering its children.
> - **Org chart:** Both `Employee` and `Manager` implement `OrgNode`. `Manager.getSalary()` returns the sum of all reports.
>
> Key: the **client never checks** if it's dealing with a leaf or a composite — it just calls the interface method.

---

### Q6. What is the difference between Bridge and Adapter patterns?

> **Answer:**
> - **Bridge:** Designed **upfront** to let abstraction and implementation vary independently. Both sides are new; the pattern prevents class explosion (e.g., `Shape × Color` = Bridge instead of `RedCircle`, `BlueCircle`, `RedSquare`...).
> - **Adapter:** Applied **after the fact** to make an existing interface compatible with another. One side already exists and can't be changed.
>
> Mnemonic: Bridge is about design-time flexibility; Adapter is a retrofit.

---

<div align="center">
  <sub>← Back to <a href="Topic.md">All Patterns</a></sub>
</div>
