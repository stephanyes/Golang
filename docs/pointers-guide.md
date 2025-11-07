# Pointers in Go
*For JavaScript/TypeScript Developers*

## Table of Contents
1. [Memory Fundamentals: Go vs JavaScript](#memory-fundamentals-go-vs-javascript)
2. [What Are Pointers?](#what-are-pointers)
3. [Pointer Syntax and Operations](#pointer-syntax-and-operations)
4. [Value vs Reference Semantics](#value-vs-reference-semantics)
5. [Pointers with Functions](#pointers-with-functions)
6. [Pointers with Structs](#pointers-with-structs)
7. [Arrays, Slices, and Maps](#arrays-slices-and-maps)
8. [Common Pointer Patterns](#common-pointer-patterns)
9. [Memory Management](#memory-management)
10. [Practical Use Cases](#practical-use-cases)
11. [Performance Considerations](#performance-considerations)
12. [Best Practices and Pitfalls](#best-practices-and-pitfalls)

## Memory Fundamentals: Go vs JavaScript

### JavaScript's Hidden Memory Management
```javascript
// JavaScript: Everything is managed automatically
let obj1 = { name: "John", age: 30 };
let obj2 = obj1;  // obj2 points to same object

obj2.age = 31;
console.log(obj1.age);  // 31 - both reference same object

// Primitives are copied
let a = 5;
let b = a;  // b gets copy of a's value
b = 10;
console.log(a);  // 5 - a is unchanged

// You never see memory addresses
console.log(obj1);  // { name: "John", age: 31 }
// No way to get the actual memory address
```

### Go's Explicit Memory Control
```go
// Go: Explicit control over values vs addresses
var obj1 = struct{ name string; age int }{"John", 30}
var obj2 = obj1  // obj2 gets a COPY of obj1

obj2.age = 31
fmt.Printf("obj1.age: %d\n", obj1.age)  // 30 - obj1 unchanged!

// To share the same memory, use pointers
var ptr1 = &obj1  // ptr1 points to obj1's memory address
var ptr2 = ptr1   // ptr2 points to same address

ptr2.age = 31
fmt.Printf("obj1.age: %d\n", obj1.age)  // 31 - obj1 changed!

// You can see memory addresses
fmt.Printf("Address of obj1: %p\n", &obj1)     // 0xc000010040
fmt.Printf("Value of ptr1: %p\n", ptr1)        // 0xc000010040 (same address)
fmt.Printf("Address of ptr1: %p\n", &ptr1)     // 0xc000010048 (different address)
```

**Key Difference:** JavaScript hides memory management, Go makes it explicit and controllable.

## What Are Pointers?

### Conceptual Understanding
```
Memory visualization:

Address:    Value:
0x1000      42        <- Variable 'x' lives here
0x1004      0x1000    <- Pointer 'ptr' lives here, contains address of x

Go code:
x := 42
ptr := &x

JavaScript equivalent concept:
let x = 42;
let ptr = /* no direct equivalent - JavaScript hides this */
```

### Basic Pointer Concepts
```go
// Declare an integer
x := 42

// Get pointer to x (address-of operator &)
ptr := &x
fmt.Printf("x = %d\n", x)           // 42
fmt.Printf("&x = %p\n", &x)         // 0xc000010040 (memory address)
fmt.Printf("ptr = %p\n", ptr)       // 0xc000010040 (same address)
fmt.Printf("*ptr = %d\n", *ptr)     // 42 (dereference - get value at address)

// Modify value through pointer
*ptr = 100
fmt.Printf("x = %d\n", x)           // 100 (x changed!)

// Zero value of pointer is nil
var uninitializedPtr *int
fmt.Printf("uninitializedPtr = %v\n", uninitializedPtr)  // <nil>

// Check for nil before dereferencing
if uninitializedPtr != nil {
    fmt.Printf("*uninitializedPtr = %d\n", *uninitializedPtr)
} else {
    fmt.Println("Pointer is nil!")
}
```

### Pointer Types
```go
// Pointers are typed - pointer to int, pointer to string, etc.
var intPtr *int
var stringPtr *string
var boolPtr *bool

// Can't assign wrong type
// intPtr = &"hello"  // ** Compile error!

// Custom types
type Person struct {
    Name string
    Age  int
}

var personPtr *Person
person := Person{Name: "Alice", Age: 25}
personPtr = &person

fmt.Printf("(*personPtr).Name = %s\n", (*personPtr).Name)  // Alice
fmt.Printf("personPtr.Name = %s\n", personPtr.Name)        // Alice (syntactic sugar)
```

**JavaScript Comparison:**
```javascript
// JavaScript doesn't have explicit pointer types
// But conceptually similar with object references
let person = { name: "Alice", age: 25 };
let personRef = person;  // Reference, not pointer

// No way to get actual memory address
// No way to create pointer to primitive
// No nil checking needed (has undefined/null instead)
```

## Pointer Syntax and Operations

### The Three Pointer Operations
```go
// 1. Address-of operator (&) - get pointer to variable
name := "John"
namePtr := &name
fmt.Printf("Address of name: %p\n", namePtr)

// 2. Dereference operator (*) - get value at pointer
fmt.Printf("Value at namePtr: %s\n", *namePtr)

// 3. Pointer declaration (*Type) - declare pointer variable
var agePtr *int  // Pointer to int, initially nil

age := 30
agePtr = &age    // Assign address of age
fmt.Printf("Value through agePtr: %d\n", *agePtr)
```

### Pointer Arithmetic (Limited in Go)
```go
// Go doesn't allow arbitrary pointer arithmetic (unlike C/C++)
// This is for safety and garbage collector compatibility

slice := []int{10, 20, 30, 40, 50}
ptr := &slice[2]  // Pointer to third element

fmt.Printf("*ptr = %d\n", *ptr)  // 30

// ** These don't work in Go:
// ptr++          // No pointer arithmetic
// ptr + 2        // No pointer arithmetic
// ptr - 1        // No pointer arithmetic

// ** Use slice operations instead:
fmt.Printf("Next element: %d\n", slice[3])    // 40
fmt.Printf("Previous element: %d\n", slice[1]) // 20
```

### Pointer to Pointer
```go
x := 42
ptr := &x        // Pointer to int
ptrPtr := &ptr   // Pointer to pointer to int

fmt.Printf("x = %d\n", x)                    // 42
fmt.Printf("*ptr = %d\n", *ptr)             // 42
fmt.Printf("**ptrPtr = %d\n", **ptrPtr)     // 42

// Modify through double pointer
**ptrPtr = 100
fmt.Printf("x = %d\n", x)                   // 100

// Type of ptrPtr
fmt.Printf("Type of ptrPtr: %T\n", ptrPtr)  // **int
```

## Value vs Reference Semantics

### Pass by Value (Default in Go)
```go
func modifyValue(x int) {
    x = 100  // Only modifies the copy
    fmt.Printf("Inside function: x = %d\n", x)  // 100
}

func main() {
    original := 42
    modifyValue(original)
    fmt.Printf("After function: original = %d\n", original)  // 42 (unchanged)
}
```

### Pass by Pointer (Explicit Reference)
```go
func modifyValueByPointer(ptr *int) {
    *ptr = 100  // Modifies the original value
    fmt.Printf("Inside function: *ptr = %d\n", *ptr)  // 100
}

func main() {
    original := 42
    modifyValueByPointer(&original)  // Pass address
    fmt.Printf("After function: original = %d\n", original)  // 100 (changed!)
}
```

**JavaScript Comparison:**
```javascript
// JavaScript: Objects passed by reference, primitives by value
function modifyObject(obj) {
    obj.value = 100;  // Modifies original object
}

function modifyPrimitive(x) {
    x = 100;  // Only modifies the copy
}

let obj = { value: 42 };
let primitive = 42;

modifyObject(obj);        // obj.value becomes 100
modifyPrimitive(primitive); // primitive remains 42

console.log(obj.value);   // 100 (changed)
console.log(primitive);   // 42 (unchanged)

// Go makes this explicit with pointers
// JavaScript does this automatically based on type
```

### Struct Value vs Pointer Semantics
```go
type Point struct {
    X, Y float64
}

// Value receiver - operates on copy
func (p Point) MoveValue(dx, dy float64) {
    p.X += dx  // Only modifies the copy
    p.Y += dy
}

// Pointer receiver - operates on original
func (p *Point) MovePointer(dx, dy float64) {
    p.X += dx  // Modifies the original
    p.Y += dy
}

func main() {
    point := Point{X: 1, Y: 2}

    point.MoveValue(5, 5)
    fmt.Printf("After MoveValue: %+v\n", point)    // {X:1 Y:2} (unchanged)

    point.MovePointer(5, 5)
    fmt.Printf("After MovePointer: %+v\n", point)  // {X:6 Y:7} (changed)
}
```

## Pointers with Functions

### Returning Pointers
```go
// ** Safe: Go handles memory management
func createPoint(x, y float64) *Point {
    point := Point{X: x, Y: y}  // Local variable
    return &point               // Return pointer to local variable
    // Go's escape analysis moves this to heap automatically
}

// ** Also safe: Direct creation
func createPointDirect(x, y float64) *Point {
    return &Point{X: x, Y: y}
}

func main() {
    p1 := createPoint(1, 2)
    p2 := createPointDirect(3, 4)

    fmt.Printf("p1: %+v\n", *p1)  // {X:1 Y:2}
    fmt.Printf("p2: %+v\n", *p2)  // {X:3 Y:4}
}
```

### Multiple Return Values with Pointers
```go
// Common pattern: return value and error
func findUser(id int) (*User, error) {
    if id <= 0 {
        return nil, errors.New("invalid user ID")
    }

    // Simulate database lookup
    user := &User{
        ID:   id,
        Name: "User " + strconv.Itoa(id),
    }

    return user, nil
}

func main() {
    user, err := findUser(123)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    if user != nil {
        fmt.Printf("Found user: %+v\n", *user)
    }
}
```

### Nil Pointer Handling
```go
func processUser(user *User) {
    // Always check for nil before dereferencing
    if user == nil {
        fmt.Println("User is nil")
        return
    }

    fmt.Printf("Processing user: %s\n", user.Name)
}

func safeGetUserName(user *User) string {
    if user == nil {
        return "Unknown"
    }
    return user.Name
}

// Nil-safe method pattern
func (u *User) SafeName() string {
    if u == nil {
        return "Unknown"
    }
    return u.Name
}

func main() {
    var user *User  // nil pointer

    processUser(user)           // "User is nil"
    name := safeGetUserName(user) // "Unknown"

    // Method call on nil pointer (safe if method handles it)
    fmt.Println(user.SafeName())  // "Unknown"

    // ** This would panic:
    // fmt.Println(user.Name)     // panic: runtime error: invalid memory address
}
```

**JavaScript Comparison:**
```javascript
// JavaScript uses null/undefined instead of nil
function processUser(user) {
    if (user == null) {  // Checks both null and undefined
        console.log("User is null/undefined");
        return;
    }

    console.log(`Processing user: ${user.name}`);
}

function safeGetUserName(user) {
    return user?.name || "Unknown";  // Optional chaining (ES2020)
}

let user = null;
processUser(user);           // "User is null/undefined"
let name = safeGetUserName(user); // "Unknown"

// ** This would throw:
// console.log(user.name);   // TypeError: Cannot read property 'name' of null
```

## Pointers with Structs

### Creating and Using Struct Pointers
```go
type User struct {
    ID    int
    Name  string
    Email string
}

func main() {
    // Method 1: Create value, then get pointer
    user1 := User{ID: 1, Name: "John", Email: "john@example.com"}
    userPtr1 := &user1

    // Method 2: Create pointer directly
    userPtr2 := &User{ID: 2, Name: "Alice", Email: "alice@example.com"}

    // Method 3: Using new (returns pointer to zero value)
    userPtr3 := new(User)
    userPtr3.ID = 3
    userPtr3.Name = "Bob"
    userPtr3.Email = "bob@example.com"

    // Access fields (Go provides syntactic sugar)
    fmt.Printf("Method 1 - ID: %d, Name: %s\n", userPtr1.ID, userPtr1.Name)
    fmt.Printf("Method 2 - ID: %d, Name: %s\n", userPtr2.ID, userPtr2.Name)
    fmt.Printf("Method 3 - ID: %d, Name: %s\n", userPtr3.ID, userPtr3.Name)

    // Explicit dereferencing (equivalent to above)
    fmt.Printf("Explicit - ID: %d, Name: %s\n", (*userPtr1).ID, (*userPtr1).Name)
}
```

### Constructor Functions with Pointers
```go
// Constructor that returns pointer (common pattern)
func NewUser(name, email string) *User {
    return &User{
        ID:    generateID(),  // Assume this function exists
        Name:  name,
        Email: email,
    }
}

// Constructor with validation
func NewUserWithValidation(name, email string) (*User, error) {
    if name == "" {
        return nil, errors.New("name cannot be empty")
    }
    if !isValidEmail(email) {
        return nil, errors.New("invalid email format")
    }

    return &User{
        ID:    generateID(),
        Name:  name,
        Email: email,
    }, nil
}

// Builder pattern using pointers
type UserBuilder struct {
    user *User
}

func NewUserBuilder() *UserBuilder {
    return &UserBuilder{
        user: &User{},
    }
}

func (ub *UserBuilder) WithName(name string) *UserBuilder {
    ub.user.Name = name
    return ub
}

func (ub *UserBuilder) WithEmail(email string) *UserBuilder {
    ub.user.Email = email
    return ub
}

func (ub *UserBuilder) Build() *User {
    ub.user.ID = generateID()
    return ub.user
}

// Usage
user := NewUserBuilder().
    WithName("Charlie").
    WithEmail("charlie@example.com").
    Build()
```

### Modifying Structs Through Pointers
```go
type BankAccount struct {
    AccountNumber string
    Balance       float64
    Owner         string
}

// Methods that modify the struct need pointer receivers
func (ba *BankAccount) Deposit(amount float64) error {
    if amount <= 0 {
        return errors.New("deposit amount must be positive")
    }
    ba.Balance += amount
    return nil
}

func (ba *BankAccount) Withdraw(amount float64) error {
    if amount <= 0 {
        return errors.New("withdrawal amount must be positive")
    }
    if amount > ba.Balance {
        return errors.New("insufficient funds")
    }
    ba.Balance -= amount
    return nil
}

// Read-only methods can use value receivers (but pointer is often preferred for consistency)
func (ba BankAccount) GetBalance() float64 {
    return ba.Balance
}

func main() {
    account := &BankAccount{
        AccountNumber: "12345",
        Balance:       1000.0,
        Owner:         "John Doe",
    }

    fmt.Printf("Initial balance: $%.2f\n", account.GetBalance())  // $1000.00

    account.Deposit(500.0)
    fmt.Printf("After deposit: $%.2f\n", account.GetBalance())   // $1500.00

    err := account.Withdraw(200.0)
    if err != nil {
        fmt.Printf("Withdrawal error: %v\n", err)
    } else {
        fmt.Printf("After withdrawal: $%.2f\n", account.GetBalance())  // $1300.00
    }
}
```

## Arrays, Slices, and Maps

### Arrays and Pointers
```go
// Arrays are values - copying creates independent copies
arr1 := [3]int{1, 2, 3}
arr2 := arr1  // Copies entire array

arr2[0] = 100
fmt.Printf("arr1: %v\n", arr1)  // [1 2 3] (unchanged)
fmt.Printf("arr2: %v\n", arr2)  // [100 2 3]

// Pointer to array
arrPtr := &arr1
(*arrPtr)[0] = 999
fmt.Printf("arr1: %v\n", arr1)  // [999 2 3] (changed through pointer)

// Function with array pointer
func modifyArray(arr *[3]int) {
    arr[1] = 888  // Syntactic sugar for (*arr)[1] = 888
}

modifyArray(&arr1)
fmt.Printf("arr1: %v\n", arr1)  // [999 888 3]
```

### Slices and Pointers
```go
// Slices are reference types (contain pointer to underlying array)
slice1 := []int{1, 2, 3}
slice2 := slice1  // Both point to same underlying array

slice2[0] = 100
fmt.Printf("slice1: %v\n", slice1)  // [100 2 3] (changed!)
fmt.Printf("slice2: %v\n", slice2)  // [100 2 3]

// Pointer to slice
slicePtr := &slice1
(*slicePtr)[1] = 200
fmt.Printf("slice1: %v\n", slice1)  // [100 200 3]

// Appending to slice might change underlying array
*slicePtr = append(*slicePtr, 4)
fmt.Printf("slice1: %v\n", slice1)  // [100 200 3 4]

// Slice of pointers
var numbers []int = []int{10, 20, 30}
var pointerSlice []*int

for i := range numbers {
    pointerSlice = append(pointerSlice, &numbers[i])
}

// Modify through pointers
*pointerSlice[0] = 999
fmt.Printf("numbers: %v\n", numbers)  // [999 20 30]
```

### Maps and Pointers
```go
// Maps are reference types
map1 := map[string]int{"a": 1, "b": 2}
map2 := map1  // Both reference same underlying map

map2["c"] = 3
fmt.Printf("map1: %v\n", map1)  // map[a:1 b:2 c:3] (changed!)

// Pointer to map
mapPtr := &map1
(*mapPtr)["d"] = 4
fmt.Printf("map1: %v\n", map1)  // map[a:1 b:2 c:3 d:4]

// Map of pointers
type Person struct {
    Name string
    Age  int
}

people := map[string]*Person{
    "john":  {Name: "John", Age: 30},
    "alice": {Name: "Alice", Age: 25},
}

// Modify through pointer
people["john"].Age = 31
fmt.Printf("John's age: %d\n", people["john"].Age)  // 31

// Safe access with nil check
if person, exists := people["bob"]; exists && person != nil {
    fmt.Printf("Bob's age: %d\n", person.Age)
} else {
    fmt.Println("Bob not found or nil")
}
```

**JavaScript Comparison:**
```javascript
// JavaScript arrays and objects are always references
let arr1 = [1, 2, 3];
let arr2 = arr1;  // Reference, not copy

arr2[0] = 100;
console.log(arr1);  // [100, 2, 3] (changed)

let map1 = { a: 1, b: 2 };
let map2 = map1;  // Reference, not copy

map2.c = 3;
console.log(map1);  // { a: 1, b: 2, c: 3 } (changed)

// To copy, you need to be explicit:
let arr3 = [...arr1];    // Shallow copy array
let map3 = {...map1};    // Shallow copy object
```

## Common Pointer Patterns

### Optional Values (Pointer to Indicate Presence)
```go
type UserProfile struct {
    Name     string
    Email    string
    Age      *int     // Optional - nil means not provided
    Bio      *string  // Optional - nil means not provided
}

func NewUserProfile(name, email string) *UserProfile {
    return &UserProfile{
        Name:  name,
        Email: email,
        // Age and Bio are nil (not provided)
    }
}

func (up *UserProfile) SetAge(age int) {
    up.Age = &age  // Set pointer to age value
}

func (up *UserProfile) SetBio(bio string) {
    up.Bio = &bio  // Set pointer to bio value
}

func (up *UserProfile) GetAge() (int, bool) {
    if up.Age == nil {
        return 0, false  // Age not set
    }
    return *up.Age, true
}

func (up *UserProfile) GetBio() (string, bool) {
    if up.Bio == nil {
        return "", false  // Bio not set
    }
    return *up.Bio, true
}

// Usage
profile := NewUserProfile("John", "john@example.com")

// Check optional fields
if age, hasAge := profile.GetAge(); hasAge {
    fmt.Printf("Age: %d\n", age)
} else {
    fmt.Println("Age not provided")
}

// Set optional fields
profile.SetAge(30)
profile.SetBio("Software developer")

// JSON marshaling respects nil pointers
jsonData, _ := json.Marshal(profile)
fmt.Println(string(jsonData))  // {"Name":"John","Email":"john@example.com","Age":30,"Bio":"Software developer"}
```

### Linked Data Structures
```go
// Linked List Node
type ListNode struct {
    Value int
    Next  *ListNode  // Pointer to next node
}

// Binary Tree Node
type TreeNode struct {
    Value int
    Left  *TreeNode  // Pointer to left child
    Right *TreeNode  // Pointer to right child
}

// Create a simple linked list
func createLinkedList(values []int) *ListNode {
    if len(values) == 0 {
        return nil
    }

    head := &ListNode{Value: values[0]}
    current := head

    for i := 1; i < len(values); i++ {
        current.Next = &ListNode{Value: values[i]}
        current = current.Next
    }

    return head
}

// Traverse linked list
func printList(head *ListNode) {
    current := head
    for current != nil {
        fmt.Printf("%d -> ", current.Value)
        current = current.Next
    }
    fmt.Println("nil")
}

// Usage
list := createLinkedList([]int{1, 2, 3, 4, 5})
printList(list)  // 1 -> 2 -> 3 -> 4 -> 5 -> nil
```

### Circular References and Cleanup
```go
type Parent struct {
    Name     string
    Children []*Child
}

type Child struct {
    Name   string
    Parent *Parent  // Back reference
}

func NewParent(name string) *Parent {
    return &Parent{
        Name:     name,
        Children: make([]*Child, 0),
    }
}

func (p *Parent) AddChild(name string) *Child {
    child := &Child{
        Name:   name,
        Parent: p,
    }
    p.Children = append(p.Children, child)
    return child
}

// Cleanup method to break circular references if needed
func (p *Parent) Cleanup() {
    for _, child := range p.Children {
        child.Parent = nil  // Break circular reference
    }
    p.Children = nil
}

// Usage
parent := NewParent("John")
child1 := parent.AddChild("Alice")
child2 := parent.AddChild("Bob")

fmt.Printf("Parent: %s\n", parent.Name)
fmt.Printf("Child1's parent: %s\n", child1.Parent.Name)
fmt.Printf("Child2's parent: %s\n", child2.Parent.Name)

// Cleanup when done (not always necessary due to GC)
defer parent.Cleanup()
```

### Method Chaining with Pointers
```go
type QueryBuilder struct {
    table   string
    fields  []string
    where   []string
    orderBy string
    limit   int
}

func NewQueryBuilder() *QueryBuilder {
    return &QueryBuilder{
        fields: make([]string, 0),
        where:  make([]string, 0),
    }
}

func (qb *QueryBuilder) Select(fields ...string) *QueryBuilder {
    qb.fields = append(qb.fields, fields...)
    return qb
}

func (qb *QueryBuilder) From(table string) *QueryBuilder {
    qb.table = table
    return qb
}

func (qb *QueryBuilder) Where(condition string) *QueryBuilder {
    qb.where = append(qb.where, condition)
    return qb
}

func (qb *QueryBuilder) OrderBy(field string) *QueryBuilder {
    qb.orderBy = field
    return qb
}

func (qb *QueryBuilder) Limit(n int) *QueryBuilder {
    qb.limit = n
    return qb
}

func (qb *QueryBuilder) Build() string {
    query := "SELECT " + strings.Join(qb.fields, ", ")
    query += " FROM " + qb.table

    if len(qb.where) > 0 {
        query += " WHERE " + strings.Join(qb.where, " AND ")
    }

    if qb.orderBy != "" {
        query += " ORDER BY " + qb.orderBy
    }

    if qb.limit > 0 {
        query += " LIMIT " + strconv.Itoa(qb.limit)
    }

    return query
}

// Usage
query := NewQueryBuilder().
    Select("id", "name", "email").
    From("users").
    Where("age > 18").
    Where("active = true").
    OrderBy("name").
    Limit(10).
    Build()

fmt.Println(query)
// SELECT id, name, email FROM users WHERE age > 18 AND active = true ORDER BY name LIMIT 10
```

## Memory Management

### Stack vs Heap Allocation
```go
// Go's escape analysis automatically decides stack vs heap

func stackExample() {
    // Local variable - likely on stack
    x := 42
    fmt.Println(x)
    // x is destroyed when function returns
}

func heapExample() *int {
    // Local variable that escapes - moved to heap
    x := 42
    return &x  // Returning pointer forces heap allocation
}

// Large structs often go to heap
func createLargeStruct() LargeStruct {
    return LargeStruct{
        Data: [1000000]int{},  // Might be allocated on heap
    }
}

// Explicit heap allocation
func useNew() {
    ptr := new(int)    // Always allocates on heap
    *ptr = 42
    fmt.Printf("Value: %d, Address: %p\n", *ptr, ptr)
}
```

### Memory Leaks and Prevention
```go
// Potential memory leak: retaining references
type Cache struct {
    data map[string]*Data
}

type Data struct {
    Value     string
    Timestamp time.Time
}

func (c *Cache) Set(key string, value string) {
    c.data[key] = &Data{
        Value:     value,
        Timestamp: time.Now(),
    }
}

// ** Without cleanup, cache grows forever
func (c *Cache) Get(key string) *Data {
    return c.data[key]
}

// ** With cleanup
func (c *Cache) Cleanup(maxAge time.Duration) {
    cutoff := time.Now().Add(-maxAge)
    for key, data := range c.data {
        if data.Timestamp.Before(cutoff) {
            delete(c.data, key)  // Remove old entries
        }
    }
}

// ** With size limit
func (c *Cache) SetWithLimit(key string, value string, maxSize int) {
    if len(c.data) >= maxSize {
        // Remove oldest entry
        var oldestKey string
        var oldestTime time.Time
        for k, data := range c.data {
            if oldestKey == "" || data.Timestamp.Before(oldestTime) {
                oldestKey = k
                oldestTime = data.Timestamp
            }
        }
        delete(c.data, oldestKey)
    }

    c.Set(key, value)
}
```

### Finalizers (Advanced)
```go
import (
    "runtime"
    "unsafe"
)

// Resource that needs cleanup
type FileResource struct {
    filename string
    handle   uintptr  // Simulated file handle
}

func NewFileResource(filename string) *FileResource {
    fr := &FileResource{
        filename: filename,
        handle:   openFile(filename),  // Simulated native call
    }

    // Set finalizer to ensure cleanup
    runtime.SetFinalizer(fr, (*FileResource).cleanup)
    return fr
}

func (fr *FileResource) cleanup() {
    if fr.handle != 0 {
        closeFile(fr.handle)  // Simulated native call
        fr.handle = 0
    }
}

// Explicit close (preferred)
func (fr *FileResource) Close() {
    fr.cleanup()
    runtime.SetFinalizer(fr, nil)  // Remove finalizer
}

// Simulated native functions
func openFile(filename string) uintptr {
    return uintptr(unsafe.Pointer(&filename))  // Fake handle
}

func closeFile(handle uintptr) {
    // Cleanup native resource
}
```

## Practical Use Cases

### Database Connections and Connection Pooling
```go
type Database struct {
    host     string
    port     int
    username string
    password string
    conn     *sql.DB  // Pointer to connection
}

func NewDatabase(host string, port int, username, password string) (*Database, error) {
    db := &Database{
        host:     host,
        port:     port,
        username: username,
        password: password,
    }

    if err := db.Connect(); err != nil {
        return nil, err
    }

    return db, nil
}

func (db *Database) Connect() error {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", db.username, db.password, db.host, db.port)

    conn, err := sql.Open("mysql", dsn)
    if err != nil {
        return err
    }

    db.conn = conn
    return nil
}

func (db *Database) Close() error {
    if db.conn != nil {
        err := db.conn.Close()
        db.conn = nil
        return err
    }
    return nil
}

func (db *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
    if db.conn == nil {
        return nil, errors.New("database connection is nil")
    }
    return db.conn.Query(query, args...)
}

// Connection pool
type ConnectionPool struct {
    connections []*Database
    available   chan *Database
    maxSize     int
}

func NewConnectionPool(maxSize int, host string, port int, username, password string) (*ConnectionPool, error) {
    pool := &ConnectionPool{
        connections: make([]*Database, 0, maxSize),
        available:   make(chan *Database, maxSize),
        maxSize:     maxSize,
    }

    for i := 0; i < maxSize; i++ {
        db, err := NewDatabase(host, port, username, password)
        if err != nil {
            return nil, err
        }
        pool.connections = append(pool.connections, db)
        pool.available <- db
    }

    return pool, nil
}

func (cp *ConnectionPool) GetConnection() *Database {
    return <-cp.available
}

func (cp *ConnectionPool) ReleaseConnection(db *Database) {
    cp.available <- db
}

func (cp *ConnectionPool) Close() {
    close(cp.available)
    for _, db := range cp.connections {
        db.Close()
    }
}
```

### HTTP Request/Response Processing
```go
type APIServer struct {
    database *Database
    cache    *Cache
    logger   *Logger
}

type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"-"`  // Never serialize
}

type APIResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   *string     `json:"error,omitempty"`  // Pointer - nil means no error
}

func (s *APIServer) GetUserHandler(w http.ResponseWriter, r *http.Request) {
    userID, err := strconv.Atoi(r.URL.Query().Get("id"))
    if err != nil {
        s.sendErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    user, err := s.getUser(userID)
    if err != nil {
        s.sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if user == nil {
        s.sendErrorResponse(w, "User not found", http.StatusNotFound)
        return
    }

    s.sendSuccessResponse(w, user)
}

func (s *APIServer) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        s.sendErrorResponse(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    if err := s.validateUser(&user); err != nil {
        s.sendErrorResponse(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := s.createUser(&user); err != nil {
        s.sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
        return
    }

    s.sendSuccessResponse(w, user)
}

func (s *APIServer) getUser(userID int) (*User, error) {
    // Try cache first
    if cached := s.cache.Get(fmt.Sprintf("user:%d", userID)); cached != nil {
        if user, ok := cached.(*User); ok {
            return user, nil
        }
    }

    // Query database
    user := &User{}
    query := "SELECT id, username, email FROM users WHERE id = ?"
    err := s.database.QueryRow(query, userID).Scan(&user.ID, &user.Username, &user.Email)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil  // User not found
        }
        return nil, err
    }

    // Cache result
    s.cache.Set(fmt.Sprintf("user:%d", userID), user, 5*time.Minute)

    return user, nil
}

func (s *APIServer) createUser(user *User) error {
    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    user.Password = string(hashedPassword)

    // Insert into database
    query := "INSERT INTO users (username, email, password) VALUES (?, ?, ?)"
    result, err := s.database.Exec(query, user.Username, user.Email, user.Password)
    if err != nil {
        return err
    }

    id, err := result.LastInsertId()
    if err != nil {
        return err
    }

    user.ID = int(id)
    user.Password = ""  // Clear password from response

    return nil
}

func (s *APIServer) sendSuccessResponse(w http.ResponseWriter, data interface{}) {
    response := APIResponse{
        Success: true,
        Data:    data,
    }
    s.sendJSONResponse(w, response, http.StatusOK)
}

func (s *APIServer) sendErrorResponse(w http.ResponseWriter, message string, status int) {
    response := APIResponse{
        Success: false,
        Error:   &message,  // Pointer to string
    }
    s.sendJSONResponse(w, response, status)
}

func (s *APIServer) sendJSONResponse(w http.ResponseWriter, data interface{}, status int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}

func (s *APIServer) validateUser(user *User) error {
    if user.Username == "" {
        return errors.New("username is required")
    }
    if user.Email == "" {
        return errors.New("email is required")
    }
    if user.Password == "" {
        return errors.New("password is required")
    }
    return nil
}
```

### Configuration Management
```go
type DatabaseConfig struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    Username string `json:"username"`
    Password string `json:"password"`
    Database string `json:"database"`
}

type ServerConfig struct {
    Port int `json:"port"`
    Host string `json:"host"`
}

type Config struct {
    Database *DatabaseConfig `json:"database,omitempty"`
    Server   *ServerConfig   `json:"server,omitempty"`
    Debug    bool           `json:"debug"`
}

// Default configuration
func DefaultConfig() *Config {
    return &Config{
        Database: &DatabaseConfig{
            Host:     "localhost",
            Port:     3306,
            Username: "root",
            Password: "",
            Database: "app",
        },
        Server: &ServerConfig{
            Port: 8080,
            Host: "0.0.0.0",
        },
        Debug: false,
    }
}

func LoadConfig(filename string) (*Config, error) {
    config := DefaultConfig()

    file, err := os.Open(filename)
    if err != nil {
        if os.IsNotExist(err) {
            // Return default config if file doesn't exist
            return config, nil
        }
        return nil, err
    }
    defer file.Close()

    if err := json.NewDecoder(file).Decode(config); err != nil {
        return nil, err
    }

    return config, nil
}

func (c *Config) SaveConfig(filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    return encoder.Encode(c)
}

// Merge configurations (useful for environment overrides)
func (c *Config) Merge(other *Config) {
    if other.Database != nil {
        if c.Database == nil {
            c.Database = &DatabaseConfig{}
        }
        if other.Database.Host != "" {
            c.Database.Host = other.Database.Host
        }
        if other.Database.Port != 0 {
            c.Database.Port = other.Database.Port
        }
        // ... other fields
    }

    if other.Server != nil {
        if c.Server == nil {
            c.Server = &ServerConfig{}
        }
        if other.Server.Host != "" {
            c.Server.Host = other.Server.Host
        }
        if other.Server.Port != 0 {
            c.Server.Port = other.Server.Port
        }
    }
}

// Usage
config, err := LoadConfig("config.json")
if err != nil {
    log.Fatal(err)
}

// Override with environment-specific config
envConfig, _ := LoadConfig("config.production.json")
if envConfig != nil {
    config.Merge(envConfig)
}

fmt.Printf("Database: %s:%d\n", config.Database.Host, config.Database.Port)
fmt.Printf("Server: %s:%d\n", config.Server.Host, config.Server.Port)
```

## Performance Considerations

### Pointer vs Value Performance
```go
type SmallStruct struct {
    A, B int
}

type LargeStruct struct {
    Data [1000]int
}

// Benchmark: Small structs - value vs pointer
func BenchmarkSmallStructValue(b *testing.B) {
    s := SmallStruct{A: 1, B: 2}
    for i := 0; i < b.N; i++ {
        processSmallStructByValue(s)
    }
}

func BenchmarkSmallStructPointer(b *testing.B) {
    s := SmallStruct{A: 1, B: 2}
    for i := 0; i < b.N; i++ {
        processSmallStructByPointer(&s)
    }
}

func processSmallStructByValue(s SmallStruct) int {
    return s.A + s.B
}

func processSmallStructByPointer(s *SmallStruct) int {
    return s.A + s.B
}

// Benchmark: Large structs - value vs pointer
func BenchmarkLargeStructValue(b *testing.B) {
    s := LargeStruct{}
    for i := 0; i < b.N; i++ {
        processLargeStructByValue(s)  // Copies 8KB each time!
    }
}

func BenchmarkLargeStructPointer(b *testing.B) {
    s := LargeStruct{}
    for i := 0; i < b.N; i++ {
        processLargeStructByPointer(&s)  // Copies only 8 bytes (pointer)
    }
}

func processLargeStructByValue(s LargeStruct) int {
    return len(s.Data)
}

func processLargeStructByPointer(s *LargeStruct) int {
    return len(s.Data)
}

/*
Typical results:
BenchmarkSmallStructValue     1000000000    0.3 ns/op
BenchmarkSmallStructPointer   1000000000    0.5 ns/op  (slightly slower due to indirection)

BenchmarkLargeStructValue     10000000      150 ns/op  (slow due to copying)
BenchmarkLargeStructPointer   1000000000    0.5 ns/op  (fast)
*/
```

### Memory Allocation Patterns
```go
// ** Inefficient: Many small allocations
func buildStringBad(items []string) string {
    var result string
    for _, item := range items {
        ptr := &item  // Unnecessary allocation
        result += *ptr + " "
    }
    return result
}

// ** Efficient: Minimize allocations
func buildStringGood(items []string) string {
    if len(items) == 0 {
        return ""
    }

    // Pre-allocate builder with estimated capacity
    var builder strings.Builder
    builder.Grow(len(items) * 10)  // Estimate

    for _, item := range items {
        builder.WriteString(item)
        builder.WriteByte(' ')
    }

    return builder.String()
}

// Slice pre-allocation
func processItemsBad() []*Item {
    var items []*Item
    for i := 0; i < 1000; i++ {
        item := &Item{ID: i}
        items = append(items, item)  // Multiple reallocations
    }
    return items
}

func processItemsGood() []*Item {
    items := make([]*Item, 0, 1000)  // Pre-allocate capacity
    for i := 0; i < 1000; i++ {
        item := &Item{ID: i}
        items = append(items, item)  // No reallocations
    }
    return items
}
```

## Best Practices and Pitfalls

### When to Use Pointers
```go
// ** Use pointers when:

// 1. You need to modify the value
func (u *User) SetName(name string) {
    u.Name = name
}

// 2. Struct is large (avoid copying)
func ProcessLargeData(data *LargeStruct) {
    // Process without copying 8KB
}

// 3. You need optional values
type Config struct {
    Database *DatabaseConfig  // nil means no database config
}

// 4. Implementing interfaces with pointer receivers
func (u *User) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Interface implementation
}

// 5. Recursive data structures
type Node struct {
    Value int
    Next  *Node
}
```

### When NOT to Use Pointers
```go
// ** Don't use pointers when:

// 1. Struct is small and rarely modified
type Point struct {
    X, Y float64
}

func (p Point) Distance() float64 {  // Value receiver is fine
    return math.Sqrt(p.X*p.X + p.Y*p.Y)
}

// 2. Function doesn't need to modify input
func FormatUser(u User) string {  // Value parameter is fine
    return fmt.Sprintf("%s <%s>", u.Name, u.Email)
}

// 3. Maps and slices (already reference types)
func ProcessUsers(users []User) {  // Don't need *[]User
    for _, user := range users {
        fmt.Println(user.Name)
    }
}
```

### Common Pitfalls

**1. Nil Pointer Dereference:**
```go
// ** Dangerous
func printUser(user *User) {
    fmt.Println(user.Name)  // Panic if user is nil
}

// ** Safe
func printUserSafe(user *User) {
    if user == nil {
        fmt.Println("User is nil")
        return
    }
    fmt.Println(user.Name)
}
```

**2. Taking Address of Loop Variable:**
```go
// ** Bug: All pointers point to same variable
func createUserPointers(names []string) []*string {
    var pointers []*string
    for _, name := range names {
        pointers = append(pointers, &name)  // ** All point to loop variable!
    }
    return pointers
}

// ** Correct: Create copy of loop variable
func createUserPointersCorrect(names []string) []*string {
    var pointers []*string
    for _, name := range names {
        nameCopy := name  // Create copy
        pointers = append(pointers, &nameCopy)
    }
    return pointers
}
```

**3. Pointer to Interface:**
```go
// ** Usually wrong
var writer *io.Writer  // Pointer to interface

// ** Usually correct
var writer io.Writer   // Interface value
```

**4. Returning Pointer to Local Slice Element:**
```go
// ** Dangerous if slice reallocates
func getFirstElement(items []int) *int {
    if len(items) == 0 {
        return nil
    }
    return &items[0]  // Dangerous if slice grows
}

// ** Safer approach
func getFirstElementSafe(items []int) (int, bool) {
    if len(items) == 0 {
        return 0, false
    }
    return items[0], true  // Return value, not pointer
}
```

### Testing with Pointers
```go
func TestUserService(t *testing.T) {
    // Setup
    user := &User{ID: 1, Name: "John", Email: "john@example.com"}

    // Test modification
    originalName := user.Name
    service := &UserService{}
    service.UpdateUserName(user, "Jane")

    // Verify original was modified
    if user.Name == originalName {
        t.Error("User name was not updated")
    }
    if user.Name != "Jane" {
        t.Errorf("Expected 'Jane', got %s", user.Name)
    }

    // Test nil safety
    err := service.UpdateUserName(nil, "Test")
    if err == nil {
        t.Error("Expected error for nil user")
    }
}

// Helper for creating test users
func newTestUser() *User {
    return &User{
        ID:    1,
        Name:  "Test User",
        Email: "test@example.com",
    }
}

// Deep copy for testing (when you need independent copies)
func (u *User) Copy() *User {
    if u == nil {
        return nil
    }
    return &User{
        ID:    u.ID,
        Name:  u.Name,
        Email: u.Email,
    }
}
```

---

*Remember: Pointers in Go provide explicit control over memory and sharing semantics. Unlike JavaScript's automatic reference handling, Go makes you decide when to share vs copy data. This explicit control leads to more predictable performance and clearer code intent!*