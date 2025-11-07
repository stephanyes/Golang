# Arrays, Loops, and Collections in Go
*For JavaScript/TypeScript Developers*

## Table of Contents
1. [Memory Layout Fundamentals](#memory-layout-fundamentals)
2. [Arrays (Fixed Size)](#arrays-fixed-size)
3. [Slices (Dynamic Arrays)](#slices-dynamic-arrays)
4. [For Loops](#for-loops)
5. [While Loops (Go Style)](#while-loops-go-style)
6. [Maps (Hash Tables)](#maps-hash-tables)
7. [Memory Performance Comparison](#memory-performance-comparison)
8. [Practical Use Cases](#practical-use-cases)
9. [Best Practices](#best-practices)

## Memory Layout Fundamentals

### What "Contiguous in Memory" Means

**Concept:**
```
JavaScript Array (not truly contiguous):
Memory: [ptr] -> [obj1] [obj2] [obj3] (scattered in heap)
        0x100   0x500  0x800  0x300

Go Array (contiguous):
Memory: [val1][val2][val3][val4] (sequential memory addresses)
        0x100 0x104 0x108 0x10C
```

**JavaScript:**
```javascript
// JavaScript arrays are actually objects with numeric keys
let jsArray = [1, 2, 3, 4];
// Internally: { "0": 1, "1": 2, "2": 3, "3": 4 }
// Elements can be scattered throughout memory
// Each access might require pointer dereferencing

// Mixed types are allowed (different memory layouts)
let mixedArray = [1, "hello", true, {name: "John"}];
```

**Go:**
```go
// Go arrays store values directly in contiguous memory
var goArray [4]int = [4]int{1, 2, 3, 4}
// Memory layout: [1][2][3][4] - each int takes 4 or 8 bytes
// Direct memory access, no pointer dereferencing needed

// All elements MUST be the same type (same memory size)
// var mixedArray [4]interface{} = [4]interface{}{1, "hello", true, map[string]string{}} // possible but not typical
```

**Performance Implications:**
- **Cache Efficiency:** Contiguous memory = better CPU cache utilization
- **Predictable Access:** O(1) access time with simple arithmetic
- **Memory Efficiency:** No pointer overhead per element

## Arrays (Fixed Size)

### Declaration and Initialization

**Go:**
```go
// Declaration with size - size is part of the type!
var numbers [5]int                    // [0, 0, 0, 0, 0]
var names [3]string                   // ["", "", ""]

// Initialization
var fruits [3]string = [3]string{"apple", "banana", "orange"}

// Shorthand initialization
colors := [4]string{"red", "green", "blue", "yellow"}

// Let Go count the elements
autoSize := [...]int{1, 2, 3, 4, 5}  // becomes [5]int

// Specific index initialization
sparse := [5]int{0: 10, 2: 20, 4: 40} // [10, 0, 20, 0, 40]
```

**JavaScript Equivalent:**
```javascript
// JavaScript doesn't have fixed-size arrays
let numbers = new Array(5).fill(0);     // [0, 0, 0, 0, 0]
let names = new Array(3).fill("");      // ["", "", ""]

let fruits = ["apple", "banana", "orange"];
let colors = ["red", "green", "blue", "yellow"];

// No direct equivalent to sparse initialization
let sparse = [10, 0, 20, 0, 40];
```

### Array Properties

**Go:**
```go
arr := [5]int{1, 2, 3, 4, 5}

// Length is fixed and part of the type
fmt.Println(len(arr))     // 5
fmt.Println(cap(arr))     // 5 (capacity same as length for arrays)

// Arrays are values, not references!
arr2 := arr              // Copies entire array
arr2[0] = 100
fmt.Println(arr[0])      // Still 1 (original unchanged)
fmt.Println(arr2[0])     // 100

// Passing to function copies the array
func modifyArray(a [5]int) {
    a[0] = 999  // Only modifies the copy
}

modifyArray(arr)
fmt.Println(arr[0])      // Still 1
```

**JavaScript Equivalent:**
```javascript
let arr = [1, 2, 3, 4, 5];

console.log(arr.length);     // 5

// Arrays are references in JavaScript
let arr2 = arr;              // Copies reference, not array
arr2[0] = 100;
console.log(arr[0]);         // 100 (original changed!)

// Passing to function passes reference
function modifyArray(a) {
    a[0] = 999;  // Modifies original array
}

modifyArray(arr);
console.log(arr[0]);         // 999
```

### Memory Layout Example

**Go Array Memory Layout:**
```go
// [4]int32 array
arr := [4]int32{10, 20, 30, 40}

/*
Memory visualization (each int32 = 4 bytes):
Address:  0x1000  0x1004  0x1008  0x100C
Value:    [  10  ][  20  ][  30  ][  40  ]

Accessing arr[2]:
- Calculate: base_address + (index * element_size)
- 0x1000 + (2 * 4) = 0x1008
- Direct memory access, no pointer following
*/
```

## Slices (Dynamic Arrays)

### Important: Arrays vs Slices - Different Types!

**Common Misconception:** "Slices are arrays with additional functionality"

**Reality:** Arrays and slices are completely different types in Go!

```go
// These are DIFFERENT types - cannot be assigned to each other
var array [5]int    // Array: fixed size, value type
var slice []int     // Slice: dynamic, reference type

// You cannot assign one to the other directly
// slice = array  // ** Compiler error!
```

**Think of it this way:**
- **Array** = The actual data storage (like a parking garage)
- **Slice** = A "view" or "window" into an array (like a camera pointing at part of the garage)

**Slice Internal Structure:**
```go
// A slice is actually this struct:
type slice struct {
    ptr *T    // Pointer to underlying array
    len int   // Current length
    cap int   // Capacity (size of underlying array)
}
```

**Example Showing the Relationship:**
```go
// Create an array (the actual storage)
arr := [6]int{1, 2, 3, 4, 5, 6}

// Create slices that "view" parts of the array
slice1 := arr[1:4]   // Views elements [2, 3, 4]
slice2 := arr[2:5]   // Views elements [3, 4, 5]

// Modifying through slice affects the underlying array
slice1[0] = 999      // Changes arr[1] to 999
fmt.Println(arr)     // [1, 999, 3, 4, 5, 6]
fmt.Println(slice2)  // [999, 4, 5] (same underlying data!)
```

**JavaScript Comparison:**
In JavaScript, there's only one "array" type that's dynamic. Go separates the concepts:
- **Go arrays** = Fixed-size, like `new Int32Array(5)` in JavaScript
- **Go slices** = Dynamic, like regular JavaScript arrays `[]`

### Understanding Slices vs Arrays

**Go:**
```go
// Slice: dynamic, reference type
var numbers []int                    // nil slice
numbers = []int{1, 2, 3, 4, 5}      // slice literal

// Slice from array
arr := [5]int{1, 2, 3, 4, 5}
slice := arr[1:4]                   // [2, 3, 4] - references arr

// Make function
made := make([]int, 5)              // length 5, capacity 5, [0,0,0,0,0]
withCap := make([]int, 3, 10)       // length 3, capacity 10, [0,0,0]
```

**JavaScript Equivalent:**
```javascript
// JavaScript arrays are always dynamic (like Go slices)
let numbers = [];                    // empty array
numbers = [1, 2, 3, 4, 5];         // array literal

// Slicing creates new array
let arr = [1, 2, 3, 4, 5];
let slice = arr.slice(1, 4);        // [2, 3, 4] - copies elements

// No direct equivalent to make(), but similar:
let made = new Array(5).fill(0);    // [0,0,0,0,0]
```

### Slice Internal Structure

**Go Slice Structure:**
```go
/*
A slice is actually a struct:
type slice struct {
    ptr    *T      // pointer to underlying array
    len    int     // current length
    cap    int     // capacity
}
*/

slice := []int{1, 2, 3, 4, 5}

fmt.Println(len(slice))    // 5
fmt.Println(cap(slice))    // 5

// Slicing doesn't copy, it creates new slice header
subSlice := slice[1:4]     // [2, 3, 4]
fmt.Println(len(subSlice)) // 3
fmt.Println(cap(subSlice)) // 4 (capacity from original)

// Modifying subSlice affects original!
subSlice[0] = 999
fmt.Println(slice)         // [1, 999, 3, 4, 5]
```

### Memory Layout: Slice vs Array

**Memory Visualization:**
```go
// Array: values stored directly
arr := [4]int{1, 2, 3, 4}
/*
Stack/Memory:
arr: [1][2][3][4]  (16 bytes on stack for 64-bit ints)
*/

// Slice: header + underlying array
slice := []int{1, 2, 3, 4}
/*
Stack:
slice: [ptr][len][cap]  (24 bytes: 8+8+8 on 64-bit)
        ↓
Heap:
       [1][2][3][4]     (underlying array)
*/
```

### Growing Slices (append)

**Go:**
```go
slice := []int{1, 2, 3}
fmt.Printf("len=%d, cap=%d\n", len(slice), cap(slice))  // len=3, cap=3

// Append within capacity
slice = append(slice, 4)
fmt.Printf("len=%d, cap=%d\n", len(slice), cap(slice))  // len=4, cap=6

// Append beyond capacity - triggers reallocation
slice = append(slice, 5, 6, 7)
fmt.Printf("len=%d, cap=%d\n", len(slice), cap(slice))  // len=7, cap=12

/*
Capacity growth strategy (approximately):
- When cap < 256: new_cap = old_cap * 2
- When cap >= 256: new_cap = old_cap * 1.25
- Plus some rounding for memory alignment
*/
```

**JavaScript Equivalent:**
```javascript
let arr = [1, 2, 3];
console.log(arr.length);    // 3 (no capacity concept)

arr.push(4);
console.log(arr.length);    // 4

arr.push(5, 6, 7);
console.log(arr.length);    // 7

// JavaScript engines handle growth automatically
// V8 uses similar doubling strategies internally
```

## For Loops

### Basic For Loop

**Go:**
```go
// C-style for loop
for i := 0; i < 10; i++ {
    fmt.Printf("%d ", i)
}

// Infinite loop
for {
    // break or return to exit
    if condition {
        break
    }
}

// Condition-only loop (while-style)
i := 0
for i < 10 {
    fmt.Printf("%d ", i)
    i++
}
```

**JavaScript Equivalent:**
```javascript
// Similar C-style for loop
for (let i = 0; i < 10; i++) {
    console.log(i);
}

// Infinite loop
for (;;) {
    if (condition) {
        break;
    }
}

// While loop (Go doesn't have separate while keyword)
let i = 0;
while (i < 10) {
    console.log(i);
    i++;
}
```

### Range Loops (Go's foreach)

**Go:**
```go
// Array/slice iteration
numbers := []int{10, 20, 30, 40, 50}

// Index and value
for index, value := range numbers {
    fmt.Printf("Index: %d, Value: %d\n", index, value)
}

// Index only
for index := range numbers {
    fmt.Printf("Index: %d\n", index)
}

// Value only (blank identifier for index)
for _, value := range numbers {
    fmt.Printf("Value: %d\n", value)
}

// String iteration (runes, not bytes!)
str := "Hello 世界"
for index, runeValue := range str {
    fmt.Printf("Index: %d, Rune: %c (Unicode: %d)\n", index, runeValue, runeValue)
}
```

**JavaScript Equivalent:**
```javascript
const numbers = [10, 20, 30, 40, 50];

// Index and value
numbers.forEach((value, index) => {
    console.log(`Index: ${index}, Value: ${value}`);
});

// Or with for...of
for (const [index, value] of numbers.entries()) {
    console.log(`Index: ${index}, Value: ${value}`);
}

// Index only
for (const index in numbers) {
    console.log(`Index: ${index}`);
}

// Value only
for (const value of numbers) {
    console.log(`Value: ${value}`);
}

// String iteration (characters)
const str = "Hello 世界";
for (const [index, char] of [...str].entries()) {
    console.log(`Index: ${index}, Char: ${char}`);
}
```

### Common Pitfall: Format Verbs in Printf

**The Problem - Type Mismatch:**
```go
myMap := map[string]int{
    "John":  25,
    "Alice": 30,
}

// ** This will cause a runtime error!
for name, age := range myMap {
    fmt.Printf("Name: %s, Age: %s\n", name, age)
    //                        ^^^ %s expects string, but age is int!
}
```

**Solutions:**
```go
// ** Solution 1: Use correct format verbs
for name, age := range myMap {
    fmt.Printf("Name: %s, Age: %d\n", name, age)
    //                        ^^^ %d for integers
}

// ** Solution 2: Use universal %v
for name, age := range myMap {
    fmt.Printf("Name: %v, Age: %v\n", name, age)
    // %v works with any type (less specific)
}

// ** Solution 3: Convert explicitly
for name, age := range myMap {
    fmt.Printf("Name: %s, Age: %s\n", name, strconv.Itoa(age))
    //                                       ^^^ Convert int to string
}
```

**Format Verb Reference:**
| Verb | Type Expected | Example Output |
|------|---------------|----------------|
| `%s` | string only | `"John"` |
| `%d` | integer only | `25` |
| `%f` | float only | `3.14` |
| `%t` | boolean only | `true` |
| `%v` | **any type** | Works with everything |

**JavaScript Comparison:**
```javascript
// JavaScript automatically converts types
for (const [name, age] of Object.entries(myMap)) {
    console.log(`Name: ${name}, Age: ${age}`); // No type errors
}
```

JavaScript's template literals automatically convert everything to strings, but Go requires explicit type matching for safety.

**Best Practice:** Use specific format verbs (`%s`, `%d`, `%f`) when you know the types - it's more explicit and catches type errors at runtime!

### Performance Considerations

**Go Range Performance:**
```go
// For large slices, range can be slower due to value copying
type LargeStruct struct {
    data [1000]int
}

largeSlice := make([]LargeStruct, 1000)

// Slower: copies each LargeStruct
for _, item := range largeSlice {
    processLargeStruct(item)  // item is a copy!
}

// Faster: use index to avoid copying
for i := range largeSlice {
    processLargeStruct(largeSlice[i])  // direct access
}

// Or range with pointer
for i := range largeSlice {
    processLargeStructPtr(&largeSlice[i])  // pass address
}
```

## While Loops (Go Style)

**Go doesn't have a `while` keyword - use `for`:**
```go
// While-style loop
count := 0
for count < 5 {
    fmt.Printf("Count: %d\n", count)
    count++
}

// Do-while style (condition at end)
count = 0
for {
    fmt.Printf("Count: %d\n", count)
    count++
    if count >= 5 {
        break
    }
}

// Reading until EOF
scanner := bufio.NewScanner(os.Stdin)
for scanner.Scan() {
    line := scanner.Text()
    if line == "quit" {
        break
    }
    fmt.Printf("You said: %s\n", line)
}
```

**JavaScript Equivalent:**
```javascript
// While loop
let count = 0;
while (count < 5) {
    console.log(`Count: ${count}`);
    count++;
}

// Do-while loop
count = 0;
do {
    console.log(`Count: ${count}`);
    count++;
} while (count < 5);
```

## Maps (Hash Tables)

### Declaration and Initialization

**Go:**
```go
// Declaration
var ages map[string]int              // nil map (can't write to it!)

// Make a map
ages = make(map[string]int)

// Map literal
ages = map[string]int{
    "Alice": 30,
    "Bob":   25,
    "Carol": 35,
}

// Shorthand
scores := map[string]int{
    "player1": 100,
    "player2": 85,
}
```

**JavaScript Equivalent:**
```javascript
// Object literal (similar to Go map)
let ages = {
    "Alice": 30,
    "Bob": 25,
    "Carol": 35
};

// Map object (closer to Go maps)
let scores = new Map([
    ["player1", 100],
    ["player2", 85]
]);

// TypeScript with type annotation
let ages: { [key: string]: number } = {
    "Alice": 30,
    "Bob": 25
};
```

### Map Operations

**Go:**
```go
ages := map[string]int{
    "Alice": 30,
    "Bob":   25,
}

// Set value
ages["Carol"] = 35

// Get value (returns zero value if key doesn't exist)
age := ages["Alice"]        // 30
missing := ages["David"]    // 0 (zero value for int)

// Check if key exists (comma ok idiom)
age, exists := ages["Alice"]
if exists {
    fmt.Printf("Alice is %d years old\n", age)
}

// Delete
delete(ages, "Bob")

// Check length
fmt.Printf("Map has %d entries\n", len(ages))

// Iterate
for name, age := range ages {
    fmt.Printf("%s: %d\n", name, age)
}
```

**JavaScript Equivalent:**
```javascript
let ages = {
    "Alice": 30,
    "Bob": 25
};

// Set value
ages["Carol"] = 35;

// Get value (returns undefined if key doesn't exist)
let age = ages["Alice"];        // 30
let missing = ages["David"];    // undefined

// Check if key exists
if ("Alice" in ages) {
    console.log(`Alice is ${ages["Alice"]} years old`);
}

// Delete
delete ages["Bob"];

// Check length
console.log(`Object has ${Object.keys(ages).length} entries`);

// Iterate
for (const [name, age] of Object.entries(ages)) {
    console.log(`${name}: ${age}`);
}

// Using Map object
let ageMap = new Map([["Alice", 30], ["Bob", 25]]);
ageMap.set("Carol", 35);
console.log(ageMap.get("Alice"));  // 30
console.log(ageMap.has("Alice"));  // true
ageMap.delete("Bob");
```

### Map Memory Layout

**Go Map Internals:**
```go
/*
Go maps use hash tables with buckets:

map[string]int structure:
┌─────────────────────────────┐
│ Hash Table Header           │
│ - bucket array pointer      │
│ - count, flags, etc.        │
└─────────────────────────────┘
              │
              ▼
┌─────────────────────────────┐
│ Bucket Array                │
│ [bucket0][bucket1][bucket2] │
└─────────────────────────────┘
              │
              ▼
┌─────────────────────────────┐
│ Bucket (holds ~8 key-value) │
│ [hash][hash][hash]...       │
│ [key ][key ][key ]...       │
│ [val ][val ][val ]...       │
│ [overflow ptr]              │
└─────────────────────────────┘
*/

// Maps are reference types
m1 := map[string]int{"a": 1}
m2 := m1                    // m2 points to same map
m2["a"] = 2
fmt.Println(m1["a"])        // 2 (same underlying map)
```

### Zero Values and Nil Maps

**Go:**
```go
// Nil map - cannot write to it!
var nilMap map[string]int
fmt.Println(nilMap == nil)      // true

// Reading from nil map returns zero value
value := nilMap["key"]          // 0 (doesn't panic)

// Writing to nil map panics!
// nilMap["key"] = 1            // panic: assignment to entry in nil map

// Must initialize before writing
nilMap = make(map[string]int)
nilMap["key"] = 1               // OK now

// Check for nil
if nilMap != nil {
    nilMap["key"] = 1
}
```

**JavaScript Equivalent:**
```javascript
// JavaScript objects are never "nil" in the Go sense
let obj;                        // undefined
console.log(obj === undefined); // true

// Accessing undefined object throws error
// console.log(obj["key"]);     // TypeError: Cannot read property

// Must initialize
obj = {};
obj["key"] = 1;                 // OK

// Check for undefined
if (obj !== undefined) {
    obj["key"] = 1;
}
```

## Memory Performance Comparison

### Cache Performance

**Go (Contiguous Memory):**
```go
// Array: excellent cache performance
arr := [1000000]int{}
for i := range arr {
    arr[i] = i * 2          // Sequential memory access
}

// Slice: good cache performance
slice := make([]int, 1000000)
for i := range slice {
    slice[i] = i * 2        // Sequential access to underlying array
}
```

**JavaScript (Non-Contiguous):**
```javascript
// JavaScript arrays: variable cache performance
let arr = new Array(1000000);
for (let i = 0; i < arr.length; i++) {
    arr[i] = i * 2;         // May require pointer dereferencing
}

// Modern JavaScript engines optimize for numeric arrays
// but still not as predictable as Go
```

### Memory Usage Comparison

**Go:**
```go
// Array: exact memory usage
arr := [1000]int64{}        // Exactly 8000 bytes (1000 * 8)

// Slice: overhead + underlying array
slice := make([]int64, 1000) // 24 bytes (header) + 8000 bytes (data)

// Map: significant overhead
m := make(map[int]int64, 1000) // Much more than 8000 bytes due to hash table structure
```

**JavaScript:**
```javascript
// Array: engine-dependent overhead
let arr = new Array(1000);   // Varies by engine, typically more overhead than Go

// Object: high overhead per property
let obj = {};
for (let i = 0; i < 1000; i++) {
    obj[i] = i;             // Each property has overhead (property descriptors, etc.)
}
```

## Practical Use Cases

### When to Use Arrays - Real Examples

**Use Case 1: Fixed Configuration Data**
```go
// RGB color values - always 3 components
type Color [3]uint8  // Red, Green, Blue

func NewColor(r, g, b uint8) Color {
    return Color{r, g, b}  // Compile-time guarantee of size
}

// 3D coordinates - always x, y, z
type Point3D [3]float64

func Distance(p1, p2 Point3D) float64 {
    // We know exactly how to iterate - no bounds checking needed
    dx := p1[0] - p2[0]
    dy := p1[1] - p2[1]
    dz := p1[2] - p2[2]
    return math.Sqrt(dx*dx + dy*dy + dz*dz)
}
```

**Use Case 2: Small Lookup Tables**
```go
// Days of the week - fixed size, known at compile time
var daysOfWeek = [7]string{
    "Monday", "Tuesday", "Wednesday", "Thursday",
    "Friday", "Saturday", "Sunday",
}

func getDayName(dayIndex int) string {
    if dayIndex < 0 || dayIndex >= 7 {
        return "Invalid"
    }
    return daysOfWeek[dayIndex]  // No bounds checking overhead
}
```

### When to Use Slices - Real Examples

**Use Case 1: Dynamic Lists (Most Common)**
```go
// User management - unknown number of users
type User struct {
    ID   int
    Name string
    Email string
}

type UserService struct {
    users []User  // Dynamic list - grows as needed
}

func (s *UserService) AddUser(name, email string) {
    newUser := User{
        ID:    len(s.users) + 1,
        Name:  name,
        Email: email,
    }
    s.users = append(s.users, newUser)  // Dynamic growth
}

func (s *UserService) GetActiveUsers() []User {
    var active []User
    for _, user := range s.users {
        if user.isActive() {
            active = append(active, user)  // Build result dynamically
        }
    }
    return active
}
```

**Use Case 2: File Processing**
```go
// Reading lines from a file - unknown number of lines
func ReadFileLines(filename string) ([]string, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var lines []string  // Start empty, grow as needed
    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        lines = append(lines, scanner.Text())  // Dynamic growth
    }

    return lines, scanner.Err()
}

// Process large files efficiently
func ProcessLargeFile(filename string) error {
    lines, err := ReadFileLines(filename)
    if err != nil {
        return err
    }

    // Pre-allocate result slice for better performance
    results := make([]ProcessedLine, 0, len(lines))

    for _, line := range lines {
        if processed := processLine(line); processed.IsValid() {
            results = append(results, processed)
        }
    }

    return saveResults(results)
}
```

**Use Case 3: API Responses**
```go
// HTTP API that returns variable number of items
type APIResponse struct {
    Items []Product `json:"items"`
    Total int       `json:"total"`
}

func GetProducts(category string, limit int) ([]Product, error) {
    // Query database - unknown result size
    rows, err := db.Query("SELECT * FROM products WHERE category = ? LIMIT ?", category, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var products []Product  // Dynamic based on query results
    for rows.Next() {
        var p Product
        if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
            return nil, err
        }
        products = append(products, p)
    }

    return products, nil
}
```

### When to Use Maps - Real Examples

**Use Case 1: Caching/Lookups**
```go
// User session cache
type SessionManager struct {
    sessions map[string]*Session  // sessionID -> Session
    mutex    sync.RWMutex
}

func (sm *SessionManager) GetSession(sessionID string) (*Session, bool) {
    sm.mutex.RLock()
    defer sm.mutex.RUnlock()

    session, exists := sm.sessions[sessionID]
    return session, exists
}

func (sm *SessionManager) CreateSession(userID int) string {
    sessionID := generateSessionID()
    session := &Session{
        ID:        sessionID,
        UserID:    userID,
        CreatedAt: time.Now(),
        ExpiresAt: time.Now().Add(24 * time.Hour),
    }

    sm.mutex.Lock()
    sm.sessions[sessionID] = session
    sm.mutex.Unlock()

    return sessionID
}
```

**Use Case 2: Configuration Management**
```go
// Application configuration
type Config struct {
    settings map[string]string  // key-value configuration
}

func LoadConfig(filename string) (*Config, error) {
    config := &Config{
        settings: make(map[string]string),
    }

    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if parts := strings.SplitN(line, "=", 2); len(parts) == 2 {
            key := strings.TrimSpace(parts[0])
            value := strings.TrimSpace(parts[1])
            config.settings[key] = value  // Dynamic key-value storage
        }
    }

    return config, nil
}

func (c *Config) Get(key string) (string, bool) {
    value, exists := c.settings[key]
    return value, exists
}

func (c *Config) GetWithDefault(key, defaultValue string) string {
    if value, exists := c.settings[key]; exists {
        return value
    }
    return defaultValue
}
```

**Use Case 3: Counting/Aggregation**
```go
// Log analysis - count occurrences
func AnalyzeLogs(logFile string) (map[string]int, error) {
    file, err := os.Open(logFile)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    errorCounts := make(map[string]int)  // error type -> count
    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        line := scanner.Text()
        if strings.Contains(line, "ERROR") {
            // Extract error type
            errorType := extractErrorType(line)
            errorCounts[errorType]++  // Increment counter
        }
    }

    return errorCounts, nil
}

// Find most common errors
func GetTopErrors(errorCounts map[string]int, limit int) []ErrorStat {
    // Convert map to slice for sorting
    var stats []ErrorStat
    for errorType, count := range errorCounts {
        stats = append(stats, ErrorStat{
            Type:  errorType,
            Count: count,
        })
    }

    // Sort by count (descending)
    sort.Slice(stats, func(i, j int) bool {
        return stats[i].Count > stats[j].Count
    })

    if len(stats) > limit {
        stats = stats[:limit]  // Slice to limit results
    }

    return stats
}
```

### Real-World Example: Building a Simple Web Server Cache

**Combining All Collection Types:**
```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "sync"
    "time"
)

// Array: Fixed status codes
var httpStatusNames = [6]string{
    200: "OK",
    404: "Not Found",
    500: "Internal Server Error",
    // ... more status codes
}

// Cache entry with expiration
type CacheEntry struct {
    Data      []byte    // Slice: variable size data
    ExpiresAt time.Time
}

// Web cache using all collection types
type WebCache struct {
    entries map[string]*CacheEntry  // Map: key-value cache
    stats   map[string]int          // Map: hit/miss statistics
    mutex   sync.RWMutex
}

func NewWebCache() *WebCache {
    return &WebCache{
        entries: make(map[string]*CacheEntry),
        stats:   make(map[string]int),
    }
}

func (wc *WebCache) Get(key string) ([]byte, bool) {
    wc.mutex.RLock()
    defer wc.mutex.RUnlock()

    entry, exists := wc.entries[key]
    if !exists || time.Now().After(entry.ExpiresAt) {
        wc.stats["misses"]++
        return nil, false
    }

    wc.stats["hits"]++
    return entry.Data, true
}

func (wc *WebCache) Set(key string, data []byte, ttl time.Duration) {
    wc.mutex.Lock()
    defer wc.mutex.Unlock()

    wc.entries[key] = &CacheEntry{
        Data:      data,  // Slice: stores variable-length data
        ExpiresAt: time.Now().Add(ttl),
    }
}

func (wc *WebCache) GetStats() map[string]int {
    wc.mutex.RLock()
    defer wc.mutex.RUnlock()

    // Return copy of stats map
    stats := make(map[string]int)
    for k, v := range wc.stats {
        stats[k] = v
    }
    return stats
}

// HTTP handler using our cache
func (wc *WebCache) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    cacheKey := r.URL.Path

    // Try to get from cache
    if data, found := wc.Get(cacheKey); found {
        w.Header().Set("X-Cache", "HIT")
        w.Write(data)
        return
    }

    // Generate response (simulate expensive operation)
    response := map[string]interface{}{
        "path":      r.URL.Path,
        "timestamp": time.Now(),
        "message":   "Hello from cache!",
    }

    data, _ := json.Marshal(response)

    // Store in cache
    wc.Set(cacheKey, data, 5*time.Minute)

    w.Header().Set("X-Cache", "MISS")
    w.Write(data)
}
```

**JavaScript Equivalent (for comparison):**
```javascript
// JavaScript version - less type safety, more dynamic
class WebCache {
    constructor() {
        this.entries = new Map();  // Similar to Go map
        this.stats = { hits: 0, misses: 0 };
    }

    get(key) {
        const entry = this.entries.get(key);
        if (!entry || Date.now() > entry.expiresAt) {
            this.stats.misses++;
            return null;
        }

        this.stats.hits++;
        return entry.data;
    }

    set(key, data, ttl) {
        this.entries.set(key, {
            data: data,  // No type constraints
            expiresAt: Date.now() + ttl
        });
    }
}
```

## Best Practices

### Choosing the Right Collection Type

**Go:**
```go
// Use arrays when:
// - Size is known at compile time
// - Small, fixed collections
// - Need guarantee of exact memory layout
var coordinates [3]float64

// Use slices when:
// - Dynamic size needed
// - Most common choice for collections
// - Need to pass to functions efficiently
var items []Item

// Pre-allocate slice capacity if size is known
items := make([]Item, 0, expectedSize)

// Use maps when:
// - Key-value lookups needed
// - Keys are not sequential integers
var cache map[string][]byte
```

### Performance Tips

**Go:**
```go
// 1. Pre-allocate slices when possible
// Bad: grows multiple times
var items []int
for i := 0; i < 1000; i++ {
    items = append(items, i)    // Multiple reallocations
}

// Good: allocate once
items := make([]int, 0, 1000)
for i := 0; i < 1000; i++ {
    items = append(items, i)    // No reallocations
}

// 2. Use range carefully with large structs
// Bad: copies large structs
for _, item := range largeStructSlice {
    process(item)               // item is a copy!
}

// Good: use index
for i := range largeStructSlice {
    process(largeStructSlice[i]) // direct access
}

// 3. Initialize maps with capacity hint
m := make(map[string]int, expectedSize)
```

---

*Understanding memory layout is crucial for Go performance. The contiguous memory model makes Go collections much more predictable and cache-friendly than JavaScript's dynamic object system!*