# Goroutines and Concurrency in Go
*For JavaScript/TypeScript Developers*

## Table of Contents
1. [Concurrency: Go vs JavaScript](#concurrency-go-vs-javascript)
2. [What Are Goroutines?](#what-are-goroutines)
3. [Basic Goroutine Usage](#basic-goroutine-usage)
4. [Channels - Communication Between Goroutines](#channels---communication-between-goroutines)
5. [Channel Patterns](#channel-patterns)
6. [Select Statement](#select-statement)
7. [Synchronization and WaitGroups](#synchronization-and-waitgroups)
8. [Context Package for Cancellation](#context-package-for-cancellation)
9. [Common Concurrency Patterns](#common-concurrency-patterns)
10. [Race Conditions and Mutexes](#race-conditions-and-mutexes)
11. [Practical Use Cases](#practical-use-cases)
12. [Performance and Best Practices](#performance-and-best-practices)

## Concurrency: Go vs JavaScript

### JavaScript's Event Loop Model
```javascript
// JavaScript: Single-threaded with event loop
console.log("1. Start");

setTimeout(() => {
    console.log("3. Timeout callback");
}, 0);

Promise.resolve().then(() => {
    console.log("2. Promise callback");
});

console.log("4. End");

// Output: 1. Start, 4. End, 2. Promise callback, 3. Timeout callback

// Async/await for non-blocking operations
async function fetchData() {
    console.log("Fetching...");
    const response = await fetch("/api/data");  // Non-blocking
    const data = await response.json();
    console.log("Data received:", data);
}

// Multiple async operations
async function fetchMultipleData() {
    const [users, posts, comments] = await Promise.all([
        fetch("/api/users").then(r => r.json()),
        fetch("/api/posts").then(r => r.json()),
        fetch("/api/comments").then(r => r.json())
    ]);

    return { users, posts, comments };
}
```

### Go's Goroutine Model
```go
// Go: True parallelism with goroutines
import (
    "fmt"
    "time"
)

func main() {
    fmt.Println("1. Start")

    // Launch goroutine (concurrent execution)
    go func() {
        fmt.Println("3. Goroutine")
    }()

    time.Sleep(10 * time.Millisecond)  // Give goroutine time to run
    fmt.Println("2. Main thread")
    fmt.Println("4. End")
}

// Equivalent to JavaScript's Promise.all
func fetchMultipleData() ([]User, []Post, []Comment, error) {
    type result struct {
        users    []User
        posts    []Post
        comments []Comment
        err      error
    }

    resultChan := make(chan result)

    // Launch concurrent fetches
    go func() {
        users, err := fetchUsers()
        posts, err2 := fetchPosts()
        comments, err3 := fetchComments()

        if err != nil {
            resultChan <- result{err: err}
            return
        }
        if err2 != nil {
            resultChan <- result{err: err2}
            return
        }
        if err3 != nil {
            resultChan <- result{err: err3}
            return
        }

        resultChan <- result{users: users, posts: posts, comments: comments}
    }()

    res := <-resultChan
    return res.users, res.posts, res.comments, res.err
}
```

**Key Differences:**
- **JavaScript**: Single-threaded, event-driven, async/await for I/O
- **Go**: Multi-threaded, goroutines for true parallelism, channels for communication

## What Are Goroutines?

### Conceptual Understanding
```
Traditional Threads vs Goroutines:

OS Threads:
- Heavy (8MB stack space each)
- Expensive to create/destroy
- Limited by OS (typically 1000s)
- Scheduled by OS

Goroutines:
- Lightweight (2KB initial stack)
- Cheap to create/destroy
- Can have millions
- Scheduled by Go runtime
```

### Basic Goroutine Creation
```go
import (
    "fmt"
    "time"
)

func sayHello(name string) {
    for i := 0; i < 3; i++ {
        fmt.Printf("Hello, %s! (%d)\n", name, i)
        time.Sleep(100 * time.Millisecond)
    }
}

func main() {
    // Regular function call (synchronous)
    sayHello("Alice")

    fmt.Println("---")

    // Goroutine (asynchronous)
    go sayHello("Bob")     // Runs concurrently
    go sayHello("Charlie") // Runs concurrently

    // Main function needs to wait for goroutines
    time.Sleep(500 * time.Millisecond)
    fmt.Println("Main function ending")
}

// Output (order may vary):
// Hello, Alice! (0)
// Hello, Alice! (1)
// Hello, Alice! (2)
// ---
// Hello, Bob! (0)
// Hello, Charlie! (0)
// Hello, Bob! (1)
// Hello, Charlie! (1)
// Hello, Bob! (2)
// Hello, Charlie! (2)
// Main function ending
```

### Anonymous Goroutines
```go
func main() {
    // Anonymous function goroutine
    go func() {
        fmt.Println("Anonymous goroutine 1")
    }()

    // Anonymous function with parameters
    go func(message string) {
        fmt.Printf("Anonymous goroutine 2: %s\n", message)
    }("Hello from goroutine!")

    // Capture variables from closure
    name := "World"
    go func() {
        fmt.Printf("Anonymous goroutine 3: Hello, %s!\n", name)
    }()

    time.Sleep(100 * time.Millisecond)
}
```

**JavaScript Equivalent:**
```javascript
// JavaScript async functions (closest equivalent)
async function sayHello(name) {
    for (let i = 0; i < 3; i++) {
        console.log(`Hello, ${name}! (${i})`);
        await new Promise(resolve => setTimeout(resolve, 100));
    }
}

async function main() {
    // Sequential (like regular function call)
    await sayHello("Alice");

    console.log("---");

    // Concurrent (like goroutines)
    Promise.all([
        sayHello("Bob"),
        sayHello("Charlie")
    ]);

    // Wait for completion
    await new Promise(resolve => setTimeout(resolve, 500));
    console.log("Main function ending");
}
```

## Channels - Communication Between Goroutines

### Basic Channel Operations
```go
// Channel creation
ch := make(chan int)        // Unbuffered channel
buffered := make(chan int, 5) // Buffered channel (capacity 5)

// Send and receive operations
func main() {
    messages := make(chan string)

    // Send value to channel (in goroutine to avoid deadlock)
    go func() {
        messages <- "Hello"     // Send operation
        messages <- "World"
    }()

    // Receive values from channel
    msg1 := <-messages          // Receive operation
    msg2 := <-messages

    fmt.Println(msg1, msg2)     // "Hello World"
}
```

### Channel Types and Directions
```go
// Bidirectional channel
var ch chan int = make(chan int)

// Send-only channel
var sendOnly chan<- int = ch

// Receive-only channel
var receiveOnly <-chan int = ch

// Function parameters with channel directions
func sender(ch chan<- string) {
    ch <- "message"
    // Can only send, not receive
}

func receiver(ch <-chan string) {
    msg := <-ch
    fmt.Println(msg)
    // Can only receive, not send
}

func main() {
    ch := make(chan string, 1)

    sender(ch)    // Send to channel
    receiver(ch)  // Receive from channel
}
```

### Buffered vs Unbuffered Channels
```go
func demonstrateChannels() {
    // Unbuffered channel - synchronous communication
    unbuffered := make(chan int)

    go func() {
        fmt.Println("Sending to unbuffered channel")
        unbuffered <- 42  // Blocks until someone receives
        fmt.Println("Sent to unbuffered channel")
    }()

    time.Sleep(100 * time.Millisecond)
    fmt.Println("Receiving from unbuffered channel")
    value := <-unbuffered
    fmt.Printf("Received: %d\n", value)

    fmt.Println("---")

    // Buffered channel - asynchronous communication
    buffered := make(chan int, 2)

    fmt.Println("Sending to buffered channel")
    buffered <- 1  // Doesn't block (buffer has space)
    buffered <- 2  // Doesn't block (buffer has space)
    fmt.Println("Sent to buffered channel")

    fmt.Printf("Received: %d\n", <-buffered)  // 1
    fmt.Printf("Received: %d\n", <-buffered)  // 2
}
```

**JavaScript Comparison:**
```javascript
// JavaScript doesn't have direct channel equivalent
// Closest would be async generators or custom event emitters

// Custom channel-like implementation
class Channel {
    constructor(bufferSize = 0) {
        this.buffer = [];
        this.bufferSize = bufferSize;
        this.waitingReceivers = [];
        this.waitingSenders = [];
    }

    async send(value) {
        if (this.waitingReceivers.length > 0) {
            const receiver = this.waitingReceivers.shift();
            receiver.resolve(value);
            return;
        }

        if (this.buffer.length < this.bufferSize) {
            this.buffer.push(value);
            return;
        }

        // Block until space available
        return new Promise(resolve => {
            this.waitingSenders.push({ value, resolve });
        });
    }

    async receive() {
        if (this.buffer.length > 0) {
            const value = this.buffer.shift();
            if (this.waitingSenders.length > 0) {
                const sender = this.waitingSenders.shift();
                this.buffer.push(sender.value);
                sender.resolve();
            }
            return value;
        }

        // Block until value available
        return new Promise(resolve => {
            this.waitingReceivers.push({ resolve });
        });
    }
}

// Usage
const ch = new Channel(1);
ch.send("Hello").then(() => console.log("Sent"));
ch.receive().then(msg => console.log("Received:", msg));
```

## Channel Patterns

### Worker Pool Pattern
```go
func workerPool() {
    const numWorkers = 3
    const numJobs = 10

    jobs := make(chan int, numJobs)
    results := make(chan int, numJobs)

    // Start workers
    for w := 1; w <= numWorkers; w++ {
        go worker(w, jobs, results)
    }

    // Send jobs
    for j := 1; j <= numJobs; j++ {
        jobs <- j
    }
    close(jobs)

    // Collect results
    for r := 1; r <= numJobs; r++ {
        result := <-results
        fmt.Printf("Result: %d\n", result)
    }
}

func worker(id int, jobs <-chan int, results chan<- int) {
    for job := range jobs {
        fmt.Printf("Worker %d processing job %d\n", id, job)
        time.Sleep(time.Second) // Simulate work
        results <- job * 2
    }
}
```

### Pipeline Pattern
```go
func pipeline() {
    // Stage 1: Generate numbers
    numbers := make(chan int)
    go func() {
        for i := 1; i <= 10; i++ {
            numbers <- i
        }
        close(numbers)
    }()

    // Stage 2: Square numbers
    squares := make(chan int)
    go func() {
        for num := range numbers {
            squares <- num * num
        }
        close(squares)
    }()

    // Stage 3: Filter even squares
    evens := make(chan int)
    go func() {
        for square := range squares {
            if square%2 == 0 {
                evens <- square
            }
        }
        close(evens)
    }()

    // Final stage: Print results
    for even := range evens {
        fmt.Printf("Even square: %d\n", even)
    }
}
```

### Fan-Out/Fan-In Pattern
```go
// Fan-out: Distribute work to multiple goroutines
func fanOut(input <-chan int) (<-chan int, <-chan int) {
    out1 := make(chan int)
    out2 := make(chan int)

    go func() {
        for val := range input {
            select {
            case out1 <- val:
            case out2 <- val:
            }
        }
        close(out1)
        close(out2)
    }()

    return out1, out2
}

// Fan-in: Combine results from multiple goroutines
func fanIn(input1, input2 <-chan int) <-chan int {
    output := make(chan int)

    go func() {
        for {
            select {
            case val, ok := <-input1:
                if !ok {
                    input1 = nil
                } else {
                    output <- val
                }
            case val, ok := <-input2:
                if !ok {
                    input2 = nil
                } else {
                    output <- val
                }
            }

            if input1 == nil && input2 == nil {
                break
            }
        }
        close(output)
    }()

    return output
}
```

### Channel Closing and Range
```go
func channelClosing() {
    ch := make(chan int)

    // Producer
    go func() {
        for i := 1; i <= 5; i++ {
            ch <- i
        }
        close(ch)  // Signal that no more values will be sent
    }()

    // Consumer using range (automatically stops when channel is closed)
    for value := range ch {
        fmt.Printf("Received: %d\n", value)
    }

    // Alternative: manual check for closed channel
    ch2 := make(chan int)
    go func() {
        ch2 <- 42
        close(ch2)
    }()

    value, ok := <-ch2
    if ok {
        fmt.Printf("Received: %d\n", value)
    } else {
        fmt.Println("Channel closed")
    }
}
```

## Select Statement

### Basic Select Usage
```go
func selectExample() {
    ch1 := make(chan string)
    ch2 := make(chan string)

    go func() {
        time.Sleep(1 * time.Second)
        ch1 <- "Message from ch1"
    }()

    go func() {
        time.Sleep(2 * time.Second)
        ch2 <- "Message from ch2"
    }()

    // Select waits for the first available channel operation
    select {
    case msg1 := <-ch1:
        fmt.Println("Received from ch1:", msg1)
    case msg2 := <-ch2:
        fmt.Println("Received from ch2:", msg2)
    }
}
```

### Select with Default Case (Non-blocking)
```go
func nonBlockingSelect() {
    ch := make(chan string, 1)

    // Non-blocking send
    select {
    case ch <- "Hello":
        fmt.Println("Sent message")
    default:
        fmt.Println("Channel full, couldn't send")
    }

    // Non-blocking receive
    select {
    case msg := <-ch:
        fmt.Println("Received:", msg)
    default:
        fmt.Println("No message available")
    }
}
```

### Select with Timeout
```go
func selectWithTimeout() {
    ch := make(chan string)

    go func() {
        time.Sleep(2 * time.Second)
        ch <- "Delayed message"
    }()

    select {
    case msg := <-ch:
        fmt.Println("Received:", msg)
    case <-time.After(1 * time.Second):
        fmt.Println("Timeout: no message received")
    }
}
```

**JavaScript Equivalent:**
```javascript
// JavaScript Promise.race is similar to select
async function selectExample() {
    const promise1 = new Promise(resolve =>
        setTimeout(() => resolve("Message from promise1"), 1000)
    );

    const promise2 = new Promise(resolve =>
        setTimeout(() => resolve("Message from promise2"), 2000)
    );

    // Race: first promise to resolve wins
    const result = await Promise.race([promise1, promise2]);
    console.log("First result:", result);  // "Message from promise1"

    // With timeout
    const timeoutPromise = new Promise((_, reject) =>
        setTimeout(() => reject(new Error("Timeout")), 1000)
    );

    try {
        const result = await Promise.race([promise2, timeoutPromise]);
        console.log("Result:", result);
    } catch (error) {
        console.log("Timeout occurred");
    }
}
```

## Synchronization and WaitGroups

### WaitGroup for Goroutine Synchronization
```go
import (
    "sync"
    "time"
)

func waitGroupExample() {
    var wg sync.WaitGroup

    // Launch multiple goroutines
    for i := 1; i <= 3; i++ {
        wg.Add(1)  // Increment counter

        go func(id int) {
            defer wg.Done()  // Decrement counter when done

            fmt.Printf("Worker %d starting\n", id)
            time.Sleep(time.Duration(id) * time.Second)
            fmt.Printf("Worker %d done\n", id)
        }(i)
    }

    fmt.Println("Waiting for all workers to complete...")
    wg.Wait()  // Block until counter reaches 0
    fmt.Println("All workers completed!")
}
```

### WaitGroup with Error Handling
```go
func waitGroupWithErrors() {
    var wg sync.WaitGroup
    errors := make(chan error, 3)  // Buffered channel for errors

    for i := 1; i <= 3; i++ {
        wg.Add(1)

        go func(id int) {
            defer wg.Done()

            // Simulate some work that might fail
            if id == 2 {
                errors <- fmt.Errorf("worker %d failed", id)
                return
            }

            fmt.Printf("Worker %d completed successfully\n", id)
            errors <- nil  // No error
        }(i)
    }

    // Close errors channel when all goroutines are done
    go func() {
        wg.Wait()
        close(errors)
    }()

    // Collect errors
    var workerErrors []error
    for err := range errors {
        if err != nil {
            workerErrors = append(workerErrors, err)
        }
    }

    if len(workerErrors) > 0 {
        fmt.Printf("Errors occurred: %v\n", workerErrors)
    } else {
        fmt.Println("All workers completed successfully!")
    }
}
```

**JavaScript Equivalent:**
```javascript
// JavaScript Promise.all and Promise.allSettled
async function waitGroupExample() {
    const workers = [];

    for (let i = 1; i <= 3; i++) {
        const worker = new Promise((resolve) => {
            console.log(`Worker ${i} starting`);
            setTimeout(() => {
                console.log(`Worker ${i} done`);
                resolve(i);
            }, i * 1000);
        });
        workers.push(worker);
    }

    console.log("Waiting for all workers to complete...");
    await Promise.all(workers);
    console.log("All workers completed!");
}

// With error handling (like Go's error channel)
async function waitGroupWithErrors() {
    const workers = [];

    for (let i = 1; i <= 3; i++) {
        const worker = new Promise((resolve, reject) => {
            if (i === 2) {
                reject(new Error(`Worker ${i} failed`));
                return;
            }
            console.log(`Worker ${i} completed successfully`);
            resolve(i);
        });
        workers.push(worker);
    }

    // Promise.allSettled waits for all, regardless of success/failure
    const results = await Promise.allSettled(workers);

    const errors = results
        .filter(result => result.status === 'rejected')
        .map(result => result.reason);

    if (errors.length > 0) {
        console.log("Errors occurred:", errors);
    } else {
        console.log("All workers completed successfully!");
    }
}
```

## Context Package for Cancellation

### Basic Context Usage
```go
import (
    "context"
    "time"
)

func contextExample() {
    // Create context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()  // Always call cancel to free resources

    // Start work that can be cancelled
    result := make(chan string)
    go doWork(ctx, result)

    select {
    case res := <-result:
        fmt.Println("Work completed:", res)
    case <-ctx.Done():
        fmt.Println("Work cancelled:", ctx.Err())
    }
}

func doWork(ctx context.Context, result chan<- string) {
    for i := 0; i < 5; i++ {
        select {
        case <-ctx.Done():
            fmt.Println("Work cancelled early")
            return
        default:
            // Do some work
            time.Sleep(500 * time.Millisecond)
            fmt.Printf("Work step %d completed\n", i+1)
        }
    }

    result <- "All work completed"
}
```

### Context with Values
```go
type contextKey string

const userIDKey contextKey = "userID"

func contextWithValues() {
    // Create context with value
    ctx := context.WithValue(context.Background(), userIDKey, 12345)

    processRequest(ctx)
}

func processRequest(ctx context.Context) {
    // Extract value from context
    if userID, ok := ctx.Value(userIDKey).(int); ok {
        fmt.Printf("Processing request for user %d\n", userID)
    } else {
        fmt.Println("No user ID in context")
    }

    // Pass context to other functions
    authenticateUser(ctx)
}

func authenticateUser(ctx context.Context) {
    if userID, ok := ctx.Value(userIDKey).(int); ok {
        fmt.Printf("Authenticating user %d\n", userID)
    }
}
```

### Context Cancellation Propagation
```go
func contextCancellation() {
    // Parent context
    parentCtx, parentCancel := context.WithCancel(context.Background())

    // Child context
    childCtx, childCancel := context.WithTimeout(parentCtx, 5*time.Second)
    defer childCancel()

    // Start work with child context
    go func() {
        select {
        case <-time.After(3 * time.Second):
            fmt.Println("Work completed normally")
        case <-childCtx.Done():
            fmt.Println("Child context cancelled:", childCtx.Err())
        }
    }()

    // Cancel parent after 1 second (cancels child too)
    time.Sleep(1 * time.Second)
    parentCancel()

    time.Sleep(2 * time.Second)
}
```

## Race Conditions and Mutexes

### Race Condition Example
```go
import "sync"

var counter int

func raceCondition() {
    var wg sync.WaitGroup

    // Launch 1000 goroutines that increment counter
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter++  // ** Race condition! Multiple goroutines accessing same variable
        }()
    }

    wg.Wait()
    fmt.Printf("Final counter value: %d\n", counter)  // Unpredictable result!
}
```

### Fixing with Mutex
```go
var (
    counter int
    mutex   sync.Mutex
)

func fixedWithMutex() {
    var wg sync.WaitGroup

    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()

            mutex.Lock()    // ** Acquire lock
            counter++       // Safe access
            mutex.Unlock()  // ** Release lock
        }()
    }

    wg.Wait()
    fmt.Printf("Final counter value: %d\n", counter)  // Always 1000
}
```

### RWMutex for Read/Write Operations
```go
type SafeMap struct {
    data map[string]int
    mu   sync.RWMutex
}

func NewSafeMap() *SafeMap {
    return &SafeMap{
        data: make(map[string]int),
    }
}

func (sm *SafeMap) Set(key string, value int) {
    sm.mu.Lock()         // Write lock (exclusive)
    sm.data[key] = value
    sm.mu.Unlock()
}

func (sm *SafeMap) Get(key string) (int, bool) {
    sm.mu.RLock()        // Read lock (shared)
    value, ok := sm.data[key]
    sm.mu.RUnlock()
    return value, ok
}

func (sm *SafeMap) GetAll() map[string]int {
    sm.mu.RLock()
    result := make(map[string]int)
    for k, v := range sm.data {
        result[k] = v    // Copy to avoid sharing internal map
    }
    sm.mu.RUnlock()
    return result
}
```

### Atomic Operations
```go
import (
    "sync/atomic"
)

var atomicCounter int64

func atomicExample() {
    var wg sync.WaitGroup

    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            atomic.AddInt64(&atomicCounter, 1)  // ** Atomic increment
        }()
    }

    wg.Wait()
    fmt.Printf("Atomic counter: %d\n", atomic.LoadInt64(&atomicCounter))  // Always 1000
}
```

## Common Concurrency Patterns

### Rate Limiting
```go
func rateLimitedWork() {
    // Create rate limiter (1 operation per 100ms)
    rateLimiter := time.Tick(100 * time.Millisecond)

    for i := 0; i < 5; i++ {
        <-rateLimiter  // Wait for rate limiter
        fmt.Printf("Performing work %d at %s\n", i, time.Now().Format("15:04:05.000"))
    }
}

// Advanced rate limiting with burst capacity
func burstRateLimiter() {
    requests := make(chan int, 5)
    limiter := make(chan time.Time, 3)  // Burst capacity of 3

    // Fill up the bucket
    for i := 0; i < 3; i++ {
        limiter <- time.Now()
    }

    // Refill bucket every 200ms
    go func() {
        for t := range time.Tick(200 * time.Millisecond) {
            select {
            case limiter <- t:
            default:
                // Bucket is full, drop the tick
            }
        }
    }()

    // Simulate 5 requests
    for i := 1; i <= 5; i++ {
        requests <- i
    }
    close(requests)

    // Process requests with rate limiting
    for req := range requests {
        <-limiter  // Wait for available token
        fmt.Printf("Processing request %d at %s\n", req, time.Now().Format("15:04:05.000"))
    }
}
```

### Timeout Pattern
```go
func timeoutPattern() {
    result := make(chan string)

    // Start long-running operation
    go func() {
        time.Sleep(2 * time.Second)
        result <- "Operation completed"
    }()

    // Wait with timeout
    select {
    case res := <-result:
        fmt.Println("Success:", res)
    case <-time.After(1 * time.Second):
        fmt.Println("Operation timed out")
    }
}
```

### Graceful Shutdown Pattern
```go
func gracefulShutdown() {
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt)

    done := make(chan bool)

    // Simulate work
    go func() {
        for {
            select {
            case <-quit:
                fmt.Println("Received interrupt signal, shutting down...")

                // Perform cleanup
                time.Sleep(1 * time.Second)
                fmt.Println("Cleanup completed")

                done <- true
                return
            default:
                fmt.Println("Working...")
                time.Sleep(500 * time.Millisecond)
            }
        }
    }()

    <-done
    fmt.Println("Server stopped gracefully")
}
```

## Practical Use Cases

### Web Server with Concurrent Request Handling
```go
import (
    "encoding/json"
    "net/http"
    "strconv"
    "sync"
    "time"
)

type Server struct {
    users map[int]*User
    mu    sync.RWMutex
}

type User struct {
    ID       int    `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Created  time.Time `json:"created"`
}

func NewServer() *Server {
    return &Server{
        users: make(map[int]*User),
    }
}

func (s *Server) GetUserHandler(w http.ResponseWriter, r *http.Request) {
    // Each request runs in its own goroutine automatically
    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    s.mu.RLock()  // Multiple concurrent reads are safe
    user, exists := s.users[id]
    s.mu.RUnlock()

    if !exists {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

func (s *Server) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    user.Created = time.Now()

    s.mu.Lock()  // Write operation needs exclusive lock
    user.ID = len(s.users) + 1
    s.users[user.ID] = &user
    s.mu.Unlock()

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

// Concurrent request processing
func (s *Server) ProcessBulkUsers(users []User) {
    var wg sync.WaitGroup

    for _, user := range users {
        wg.Add(1)
        go func(u User) {
            defer wg.Done()

            // Simulate processing (API calls, database operations, etc.)
            time.Sleep(100 * time.Millisecond)

            s.mu.Lock()
            u.ID = len(s.users) + 1
            s.users[u.ID] = &u
            s.mu.Unlock()

            fmt.Printf("Processed user: %s\n", u.Name)
        }(user)
    }

    wg.Wait()
    fmt.Println("All users processed")
}
```

### Database Connection Pool
```go
type ConnectionPool struct {
    connections chan *sql.DB
    factory     func() (*sql.DB, error)
    closed      bool
    mu          sync.Mutex
}

func NewConnectionPool(maxConnections int, factory func() (*sql.DB, error)) *ConnectionPool {
    pool := &ConnectionPool{
        connections: make(chan *sql.DB, maxConnections),
        factory:     factory,
    }

    // Pre-fill pool with connections
    for i := 0; i < maxConnections; i++ {
        conn, err := factory()
        if err != nil {
            panic(err)
        }
        pool.connections <- conn
    }

    return pool
}

func (p *ConnectionPool) Get() (*sql.DB, error) {
    p.mu.Lock()
    defer p.mu.Unlock()

    if p.closed {
        return nil, errors.New("pool is closed")
    }

    select {
    case conn := <-p.connections:
        return conn, nil
    default:
        // No available connections, create new one
        return p.factory()
    }
}

func (p *ConnectionPool) Put(conn *sql.DB) {
    p.mu.Lock()
    defer p.mu.Unlock()

    if p.closed {
        conn.Close()
        return
    }

    select {
    case p.connections <- conn:
        // Successfully returned to pool
    default:
        // Pool is full, close the connection
        conn.Close()
    }
}

func (p *ConnectionPool) Close() {
    p.mu.Lock()
    defer p.mu.Unlock()

    if p.closed {
        return
    }

    p.closed = true
    close(p.connections)

    // Close all connections in pool
    for conn := range p.connections {
        conn.Close()
    }
}

// Usage example
func databaseExample() {
    factory := func() (*sql.DB, error) {
        return sql.Open("mysql", "user:password@tcp(localhost:3306)/database")
    }

    pool := NewConnectionPool(10, factory)
    defer pool.Close()

    // Simulate concurrent database operations
    var wg sync.WaitGroup
    for i := 0; i < 50; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()

            // Get connection from pool
            conn, err := pool.Get()
            if err != nil {
                fmt.Printf("Failed to get connection: %v\n", err)
                return
            }
            defer pool.Put(conn)  // Return to pool

            // Use connection
            _, err = conn.Exec("SELECT * FROM users WHERE id = ?", id)
            if err != nil {
                fmt.Printf("Query failed: %v\n", err)
            } else {
                fmt.Printf("Query %d completed\n", id)
            }
        }(i)
    }

    wg.Wait()
}
```

### Real-time Chat Server
```go
type ChatServer struct {
    clients    map[string]chan string
    joining    chan *Client
    leaving    chan *Client
    messages   chan Message
    mu         sync.RWMutex
}

type Client struct {
    ID       string
    Username string
    Messages chan string
}

type Message struct {
    From    string `json:"from"`
    To      string `json:"to,omitempty"`  // Empty for broadcast
    Content string `json:"content"`
    Type    string `json:"type"`  // "message", "join", "leave"
}

func NewChatServer() *ChatServer {
    return &ChatServer{
        clients:  make(map[string]chan string),
        joining:  make(chan *Client),
        leaving:  make(chan *Client),
        messages: make(chan Message),
    }
}

func (cs *ChatServer) Run() {
    for {
        select {
        case client := <-cs.joining:
            cs.mu.Lock()
            cs.clients[client.ID] = client.Messages
            cs.mu.Unlock()

            // Notify all clients about new user
            cs.broadcast(Message{
                From:    "System",
                Content: fmt.Sprintf("%s joined the chat", client.Username),
                Type:    "join",
            })

            fmt.Printf("Client %s joined\n", client.Username)

        case client := <-cs.leaving:
            cs.mu.Lock()
            if _, exists := cs.clients[client.ID]; exists {
                close(cs.clients[client.ID])
                delete(cs.clients, client.ID)
            }
            cs.mu.Unlock()

            // Notify all clients about user leaving
            cs.broadcast(Message{
                From:    "System",
                Content: fmt.Sprintf("%s left the chat", client.Username),
                Type:    "leave",
            })

            fmt.Printf("Client %s left\n", client.Username)

        case msg := <-cs.messages:
            if msg.To == "" {
                // Broadcast to all clients
                cs.broadcast(msg)
            } else {
                // Send to specific client
                cs.sendToClient(msg.To, msg)
            }
        }
    }
}

func (cs *ChatServer) broadcast(msg Message) {
    cs.mu.RLock()
    for _, client := range cs.clients {
        select {
        case client <- formatMessage(msg):
        default:
            // Client's channel is full, skip
        }
    }
    cs.mu.RUnlock()
}

func (cs *ChatServer) sendToClient(clientID string, msg Message) {
    cs.mu.RLock()
    client, exists := cs.clients[clientID]
    cs.mu.RUnlock()

    if exists {
        select {
        case client <- formatMessage(msg):
        default:
            // Client's channel is full
        }
    }
}

func formatMessage(msg Message) string {
    data, _ := json.Marshal(msg)
    return string(data)
}

// Client handler for WebSocket connections
func (cs *ChatServer) HandleClient(conn *websocket.Conn, userID, username string) {
    client := &Client{
        ID:       userID,
        Username: username,
        Messages: make(chan string, 10),  // Buffered channel
    }

    // Add client to server
    cs.joining <- client

    // Cleanup when client disconnects
    defer func() {
        cs.leaving <- client
        conn.Close()
    }()

    // Start goroutine to send messages to client
    go func() {
        for msg := range client.Messages {
            if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
                break
            }
        }
    }()

    // Read messages from client
    for {
        _, data, err := conn.ReadMessage()
        if err != nil {
            break
        }

        var msg Message
        if err := json.Unmarshal(data, &msg); err != nil {
            continue
        }

        msg.From = username
        cs.messages <- msg
    }
}
```

## Performance and Best Practices

### Goroutine Pool vs Unlimited Goroutines
```go
// ** Bad: Unlimited goroutines (can overwhelm system)
func badConcurrency(tasks []Task) {
    var wg sync.WaitGroup

    for _, task := range tasks {  // Could be millions of tasks!
        wg.Add(1)
        go func(t Task) {
            defer wg.Done()
            processTask(t)
        }(task)
    }

    wg.Wait()
}

// ** Good: Limited worker pool
func goodConcurrency(tasks []Task) {
    const numWorkers = 10
    taskChan := make(chan Task, 100)  // Buffered channel

    var wg sync.WaitGroup

    // Start fixed number of workers
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for task := range taskChan {
                processTask(task)
            }
        }()
    }

    // Send tasks to workers
    for _, task := range tasks {
        taskChan <- task
    }
    close(taskChan)

    wg.Wait()
}
```

### Channel Buffer Sizing
```go
// Choose buffer size based on use case:

// Unbuffered: Synchronous communication (sender blocks until receiver ready)
ch1 := make(chan int)

// Small buffer: Slight decoupling, good for bursty workloads
ch2 := make(chan int, 10)

// Large buffer: High throughput, but uses more memory
ch3 := make(chan int, 1000)

// Rule of thumb:
// - Start with unbuffered
// - Add buffer if you observe blocking
// - Size buffer based on expected burst capacity
```

### Memory Management
```go
// ** Memory leak: Goroutine never exits
func memoryLeak() {
    ch := make(chan int)

    go func() {
        for {
            select {
            case val := <-ch:
                fmt.Println(val)
            // No case to exit the goroutine!
            }
        }
    }()

    // Channel is never used, goroutine leaks
}

// ** Proper cleanup: Always provide exit condition
func properCleanup() {
    ch := make(chan int)
    done := make(chan bool)

    go func() {
        for {
            select {
            case val := <-ch:
                fmt.Println(val)
            case <-done:
                return  // Exit goroutine
            }
        }
    }()

    // Later... signal to stop
    close(done)
}
```

### Best Practices Summary

1. **Don't create too many goroutines** - Use worker pools for large numbers of tasks
2. **Always provide exit conditions** - Prevent goroutine leaks
3. **Use buffered channels wisely** - Size based on expected throughput
4. **Prefer channels over shared memory** - "Don't communicate by sharing memory; share memory by communicating"
5. **Use context for cancellation** - Propagate cancellation signals properly
6. **Protect shared data** - Use mutexes or channels to avoid race conditions
7. **Test for race conditions** - Use `go run -race` to detect race conditions
8. **Monitor goroutine count** - Use `runtime.NumGoroutine()` to detect leaks

---

*Goroutines make Go incredibly powerful for building concurrent applications like web servers, APIs, and microservices. They allow you to handle thousands of concurrent requests efficiently while keeping your code readable and maintainable!*