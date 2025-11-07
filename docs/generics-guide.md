# Generics in Go
*For JavaScript/TypeScript Developers*

## Table of Contents
1. [Introduction to Go Generics](#introduction-to-go-generics)
2. [Basic Generic Syntax](#basic-generic-syntax)
3. [Type Parameters and Constraints](#type-parameters-and-constraints)
4. [Generic Functions](#generic-functions)
5. [Generic Types](#generic-types)
6. [Type Inference](#type-inference)
7. [Interface-Based Constraints](#interface-based-constraints)
8. [Built-in Constraints](#built-in-constraints)
9. [Generic Data Structures](#generic-data-structures)
10. [Practical Use Cases](#practical-use-cases)
11. [Performance Considerations](#performance-considerations)
12. [Best Practices](#best-practices)

## Introduction to Go Generics

### Before Generics (Go < 1.18)
```go
// Had to use interface{} for generic behavior
func FindInSlice(slice []interface{}, target interface{}) int {
    for i, item := range slice {
        if item == target {
            return i
        }
    }
    return -1
}

// Type-specific implementations needed
func FindIntInSlice(slice []int, target int) int {
    for i, item := range slice {
        if item == target {
            return i
        }
    }
    return -1
}

func FindStringInSlice(slice []string, target string) int {
    for i, item := range slice {
        if item == target {
            return i
        }
    }
    return -1
}
```

### With Generics (Go 1.18+)
```go
// Single generic function handles all comparable types
func Find[T comparable](slice []T, target T) int {
    for i, item := range slice {
        if item == target {
            return i
        }
    }
    return -1
}

// Usage with type inference
intIndex := Find([]int{1, 2, 3}, 2)           // T inferred as int
stringIndex := Find([]string{"a", "b"}, "b")  // T inferred as string
```

### TypeScript Comparison
```typescript
// TypeScript has had generics from the beginning
function find<T>(slice: T[], target: T): number {
    for (let i = 0; i < slice.length; i++) {
        if (slice[i] === target) {
            return i;
        }
    }
    return -1;
}

// Usage
const intIndex = find([1, 2, 3], 2);        // T inferred as number
const stringIndex = find(["a", "b"], "b");  // T inferred as string
```

**Key Difference:** Go's generics are compile-time only and use constraints instead of structural typing.

## Basic Generic Syntax

### Function Generic Syntax
```go
// Basic generic function
func GenericFunction[T any](param T) T {
    return param
}

// Multiple type parameters
func Pair[T, U any](first T, second U) (T, U) {
    return first, second
}

// Type parameter with constraint
func Add[T Numeric](a, b T) T {
    return a + b
}

type Numeric interface {
    int | int8 | int16 | int32 | int64 |
    uint | uint8 | uint16 | uint32 | uint64 |
    float32 | float64
}
```

### Type Generic Syntax
```go
// Generic struct
type Container[T any] struct {
    Value T
}

// Generic interface
type Comparer[T any] interface {
    Compare(T) int
}

// Generic type alias
type Slice[T any] []T
```

### TypeScript Comparison
```typescript
// Similar syntax but different constraints
function genericFunction<T>(param: T): T {
    return param;
}

function pair<T, U>(first: T, second: U): [T, U] {
    return [first, second];
}

// TypeScript uses extends for constraints
function add<T extends number>(a: T, b: T): T {
    return a + b;
}

// Generic types
interface Container<T> {
    value: T;
}

type Slice<T> = T[];
```

## Type Parameters and Constraints

### Understanding Type Parameters
```go
// T is a type parameter
func Identity[T any](value T) T {
    return value
}

// Type parameter naming conventions
func Process[
    T any,           // Single letter for simple cases
    Key comparable,  // Descriptive for clarity
    Value any,       // Descriptive for clarity
](data map[Key]Value) []Value {
    result := make([]Value, 0, len(data))
    for _, v := range data {
        result = append(result, v)
    }
    return result
}
```

### Basic Constraints
```go
// 'any' constraint - accepts any type
func Store[T any](value T) T {
    return value
}

// 'comparable' constraint - types that can be compared with == and !=
func Contains[T comparable](slice []T, target T) bool {
    for _, item := range slice {
        if item == target {
            return true
        }
    }
    return false
}

// Custom interface constraint
type Stringer interface {
    String() string
}

func PrintAll[T Stringer](items []T) {
    for _, item := range items {
        fmt.Println(item.String())
    }
}
```

### Union Type Constraints
```go
// Union constraint using type list
type Number interface {
    int | int8 | int16 | int32 | int64 |
    uint | uint8 | uint16 | uint32 | uint64 |
    float32 | float64
}

func Max[T Number](a, b T) T {
    if a > b {
        return a
    }
    return b
}

// Combining interface methods with union types
type Ordered interface {
    int | int8 | int16 | int32 | int64 |
    uint | uint8 | uint16 | uint32 | uint64 |
    float32 | float64 | string
}

type Sortable[T Ordered] interface {
    Len() int
    Less(i, j int) bool
    Swap(i, j int)
    ~[]T  // Underlying type constraint
}
```

### TypeScript Comparison
```typescript
// TypeScript uses extends and conditional types
type Number = number; // All numbers are the same type in TS

function max<T extends number>(a: T, b: T): T {
    return a > b ? a : b;
}

// TypeScript has union types too
type StringOrNumber = string | number;

function process<T extends StringOrNumber>(value: T): T {
    return value;
}

// TypeScript conditional types (more complex than Go unions)
type IsString<T> = T extends string ? true : false;
```

## Generic Functions

### Basic Generic Functions
```go
// Simple generic utility functions
func Reverse[T any](slice []T) []T {
    result := make([]T, len(slice))
    for i, v := range slice {
        result[len(slice)-1-i] = v
    }
    return result
}

func Map[T, U any](slice []T, fn func(T) U) []U {
    result := make([]U, len(slice))
    for i, v := range slice {
        result[i] = fn(v)
    }
    return result
}

func Filter[T any](slice []T, predicate func(T) bool) []T {
    var result []T
    for _, v := range slice {
        if predicate(v) {
            result = append(result, v)
        }
    }
    return result
}

func Reduce[T, U any](slice []T, initial U, fn func(U, T) U) U {
    result := initial
    for _, v := range slice {
        result = fn(result, v)
    }
    return result
}
```

### Usage Examples
```go
func main() {
    numbers := []int{1, 2, 3, 4, 5}
    strings := []string{"hello", "world", "go"}

    // Reverse
    reversedNumbers := Reverse(numbers)  // []int{5, 4, 3, 2, 1}
    reversedStrings := Reverse(strings)  // []string{"go", "world", "hello"}

    // Map
    doubled := Map(numbers, func(n int) int { return n * 2 })
    lengths := Map(strings, func(s string) int { return len(s) })

    // Filter
    evens := Filter(numbers, func(n int) bool { return n%2 == 0 })
    longStrings := Filter(strings, func(s string) bool { return len(s) > 4 })

    // Reduce
    sum := Reduce(numbers, 0, func(acc, n int) int { return acc + n })
    concat := Reduce(strings, "", func(acc, s string) string { return acc + s })

    fmt.Printf("Sum: %d, Concat: %s\n", sum, concat)
}
```

### TypeScript Comparison
```typescript
// Similar functionality but built into arrays
const numbers = [1, 2, 3, 4, 5];
const strings = ["hello", "world", "go"];

// Built-in methods (no explicit generics needed)
const reversed = [...numbers].reverse();
const doubled = numbers.map(n => n * 2);
const evens = numbers.filter(n => n % 2 === 0);
const sum = numbers.reduce((acc, n) => acc + n, 0);

// TypeScript infers types automatically
const lengths = strings.map(s => s.length);  // number[]
```

### Advanced Generic Functions
```go
// Generic function with multiple constraints
func FindMax[T Ordered](slice []T) (T, bool) {
    if len(slice) == 0 {
        var zero T
        return zero, false
    }

    max := slice[0]
    for _, v := range slice[1:] {
        if v > max {
            max = v
        }
    }
    return max, true
}

// Generic function with method constraint
type Comparable[T any] interface {
    CompareTo(T) int
}

func Sort[T Comparable[T]](slice []T) {
    sort.Slice(slice, func(i, j int) bool {
        return slice[i].CompareTo(slice[j]) < 0
    })
}

// Generic function factory
func MakeValidator[T any](validate func(T) bool) func(T) error {
    return func(value T) error {
        if !validate(value) {
            return fmt.Errorf("validation failed for value: %v", value)
        }
        return nil
    }
}

// Usage
emailValidator := MakeValidator(func(email string) bool {
    return strings.Contains(email, "@")
})

ageValidator := MakeValidator(func(age int) bool {
    return age >= 0 && age <= 150
})
```

## Generic Types

### Generic Structs
```go
// Basic generic struct
type Box[T any] struct {
    Value T
}

func (b Box[T]) Get() T {
    return b.Value
}

func (b *Box[T]) Set(value T) {
    b.Value = value
}

// Generic struct with multiple type parameters
type Pair[T, U any] struct {
    First  T
    Second U
}

func (p Pair[T, U]) Swap() Pair[U, T] {
    return Pair[U, T]{
        First:  p.Second,
        Second: p.First,
    }
}

// Generic struct with constraints
type NumericBox[T Number] struct {
    value T
    min   T
    max   T
}

func NewNumericBox[T Number](value, min, max T) *NumericBox[T] {
    return &NumericBox[T]{
        value: value,
        min:   min,
        max:   max,
    }
}

func (nb *NumericBox[T]) Set(value T) error {
    if value < nb.min || value > nb.max {
        return fmt.Errorf("value %v out of range [%v, %v]", value, nb.min, nb.max)
    }
    nb.value = value
    return nil
}
```

### Generic Interfaces
```go
// Generic interface
type Iterator[T any] interface {
    Next() (T, bool)
    HasNext() bool
}

// Generic interface with type constraint
type Serializer[T any] interface {
    Serialize(T) ([]byte, error)
    Deserialize([]byte) (T, error)
}

// Generic interface implementation
type JSONSerializer[T any] struct{}

func (js JSONSerializer[T]) Serialize(value T) ([]byte, error) {
    return json.Marshal(value)
}

func (js JSONSerializer[T]) Deserialize(data []byte) (T, error) {
    var value T
    err := json.Unmarshal(data, &value)
    return value, err
}
```

### Generic Slice Types
```go
// Custom generic slice with methods
type List[T any] []T

func NewList[T any]() List[T] {
    return make(List[T], 0)
}

func (l *List[T]) Add(item T) {
    *l = append(*l, item)
}

func (l List[T]) Get(index int) (T, error) {
    if index < 0 || index >= len(l) {
        var zero T
        return zero, fmt.Errorf("index out of range")
    }
    return l[index], nil
}

func (l List[T]) Filter(predicate func(T) bool) List[T] {
    result := NewList[T]()
    for _, item := range l {
        if predicate(item) {
            result.Add(item)
        }
    }
    return result
}

func (l List[T]) Map[U any](fn func(T) U) List[U] {
    result := make(List[U], len(l))
    for i, item := range l {
        result[i] = fn(item)
    }
    return result
}
```

### TypeScript Comparison
```typescript
// Similar generic types in TypeScript
interface Box<T> {
    value: T;
}

class BoxImpl<T> implements Box<T> {
    constructor(public value: T) {}

    get(): T {
        return this.value;
    }

    set(value: T): void {
        this.value = value;
    }
}

// Generic class with constraints
class NumericBox<T extends number> {
    constructor(
        private value: T,
        private min: T,
        private max: T
    ) {}

    set(value: T): void {
        if (value < this.min || value > this.max) {
            throw new Error(`Value ${value} out of range [${this.min}, ${this.max}]`);
        }
        this.value = value;
    }
}
```

## Type Inference

### How Go Infers Types
```go
// Type inference from arguments
func Print[T any](value T) {
    fmt.Printf("Value: %v, Type: %T\n", value, value)
}

// Go infers T from the argument
Print("hello")    // T inferred as string
Print(42)         // T inferred as int
Print(3.14)       // T inferred as float64

// Type inference with multiple parameters
func Combine[T any](a, b T) []T {
    return []T{a, b}
}

result1 := Combine(1, 2)        // T inferred as int
result2 := Combine("a", "b")    // T inferred as string
// result3 := Combine(1, "a")   // Error: cannot infer T
```

### Explicit Type Parameters
```go
// When inference isn't possible or you want to be explicit
func MakeSlice[T any](size int) []T {
    return make([]T, size)
}

// Must specify type explicitly
intSlice := MakeSlice[int](10)
stringSlice := MakeSlice[string](5)

// Explicit types for clarity
func Convert[From, To any](value From, converter func(From) To) To {
    return converter(value)
}

// Explicit type parameters for readability
result := Convert[string, int]("42", func(s string) int {
    i, _ := strconv.Atoi(s)
    return i
})
```

### Type Inference Limitations
```go
// Cases where Go cannot infer types
func Zero[T any]() T {
    var zero T
    return zero
}

// Must be explicit
intZero := Zero[int]()
stringZero := Zero[string]()

// Complex inference scenarios
func Transform[T, U any](slice []T, fn func(T) U) []U {
    result := make([]U, len(slice))
    for i, v := range slice {
        result[i] = fn(v)
    }
    return result
}

numbers := []int{1, 2, 3}
// T inferred as int, U inferred from function return type
strings := Transform(numbers, strconv.Itoa)  // U inferred as string
```

### TypeScript Comparison
```typescript
// TypeScript has more sophisticated inference
function makeArray<T>(size: number): T[] {
    return new Array(size);
}

// TypeScript can sometimes infer from context
const numbers: number[] = makeArray(10);  // T inferred as number

// TypeScript has better contextual inference
const strings = ["1", "2", "3"].map(Number);  // infers number[]

// But sometimes needs explicit types too
function identity<T>(value: T): T {
    return value;
}

const result = identity<string>("hello");  // explicit when needed
```

## Interface-Based Constraints

### Custom Constraints
```go
// Define behavior-based constraints
type Addable interface {
    Add(Addable) Addable
}

type Numeric interface {
    int | int8 | int16 | int32 | int64 |
    uint | uint8 | uint16 | uint32 | uint64 |
    float32 | float64
}

type Stringable interface {
    String() string
}

// Constraint combining interface and type union
type OrderedStringable interface {
    Ordered
    Stringable
}

func ProcessValues[T OrderedStringable](values []T) {
    sort.Slice(values, func(i, j int) bool {
        return values[i] < values[j]
    })

    for _, v := range values {
        fmt.Println(v.String())
    }
}
```

### Real-World Constraints
```go
// Database entity constraint
type Entity interface {
    GetID() int
    SetID(int)
}

type Repository[T Entity] struct {
    data map[int]T
}

func NewRepository[T Entity]() *Repository[T] {
    return &Repository[T]{
        data: make(map[int]T),
    }
}

func (r *Repository[T]) Save(entity T) {
    if entity.GetID() == 0 {
        entity.SetID(len(r.data) + 1)
    }
    r.data[entity.GetID()] = entity
}

func (r *Repository[T]) FindByID(id int) (T, bool) {
    entity, exists := r.data[id]
    return entity, exists
}

// Implementation example
type User struct {
    ID   int
    Name string
}

func (u *User) GetID() int    { return u.ID }
func (u *User) SetID(id int)  { u.ID = id }

func main() {
    userRepo := NewRepository[*User]()
    user := &User{Name: "John"}
    userRepo.Save(user)

    found, exists := userRepo.FindByID(1)
    if exists {
        fmt.Printf("Found user: %s\n", found.Name)
    }
}
```

### Constraint Embedding
```go
// Embed constraints for composition
type Reader interface {
    Read([]byte) (int, error)
}

type Writer interface {
    Write([]byte) (int, error)
}

type ReadWriter interface {
    Reader
    Writer
}

type Closer interface {
    Close() error
}

// Combine multiple constraint interfaces
type ReadWriteCloser interface {
    ReadWriter
    Closer
}

func ProcessStream[T ReadWriteCloser](stream T, data []byte) error {
    defer stream.Close()

    _, err := stream.Write(data)
    if err != nil {
        return err
    }

    buffer := make([]byte, len(data))
    _, err = stream.Read(buffer)
    return err
}
```

## Built-in Constraints

### Using golang.org/x/exp/constraints
```go
import "golang.org/x/exp/constraints"

// Signed integer constraint
func AbsInt[T constraints.Signed](value T) T {
    if value < 0 {
        return -value
    }
    return value
}

// Unsigned integer constraint
func NextPowerOfTwo[T constraints.Unsigned](value T) T {
    if value == 0 {
        return 1
    }

    power := T(1)
    for power < value {
        power <<= 1
    }
    return power
}

// Ordered constraint (supports comparison operators)
func Clamp[T constraints.Ordered](value, min, max T) T {
    if value < min {
        return min
    }
    if value > max {
        return max
    }
    return value
}

// Float constraint
func Round[T constraints.Float](value T, decimals int) T {
    multiplier := T(1)
    for i := 0; i < decimals; i++ {
        multiplier *= 10
    }
    return T(int(value*multiplier+0.5)) / multiplier
}
```

### Custom Constraint Libraries
```go
// Create your own constraint package
package constraints

// Common constraints for your domain
type ID interface {
    int | int64 | string
}

type Timestamp interface {
    int64 | uint64
}

type Money interface {
    int64 | float64  // cents or decimal
}

// Usage in your application
func FormatID[T ID](id T) string {
    return fmt.Sprintf("ID-%v", id)
}

func FormatMoney[T Money](amount T) string {
    switch v := any(amount).(type) {
    case int64:
        return fmt.Sprintf("$%.2f", float64(v)/100)
    case float64:
        return fmt.Sprintf("$%.2f", v)
    default:
        return "$0.00"
    }
}
```

## Generic Data Structures

### Generic Stack
```go
type Stack[T any] struct {
    items []T
}

func NewStack[T any]() *Stack[T] {
    return &Stack[T]{
        items: make([]T, 0),
    }
}

func (s *Stack[T]) Push(item T) {
    s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
    if len(s.items) == 0 {
        var zero T
        return zero, false
    }

    index := len(s.items) - 1
    item := s.items[index]
    s.items = s.items[:index]
    return item, true
}

func (s *Stack[T]) Peek() (T, bool) {
    if len(s.items) == 0 {
        var zero T
        return zero, false
    }
    return s.items[len(s.items)-1], true
}

func (s *Stack[T]) Size() int {
    return len(s.items)
}

func (s *Stack[T]) IsEmpty() bool {
    return len(s.items) == 0
}
```

### Generic Binary Tree
```go
type TreeNode[T Ordered] struct {
    Value T
    Left  *TreeNode[T]
    Right *TreeNode[T]
}

type BinaryTree[T Ordered] struct {
    root *TreeNode[T]
}

func NewBinaryTree[T Ordered]() *BinaryTree[T] {
    return &BinaryTree[T]{}
}

func (bt *BinaryTree[T]) Insert(value T) {
    bt.root = bt.insertNode(bt.root, value)
}

func (bt *BinaryTree[T]) insertNode(node *TreeNode[T], value T) *TreeNode[T] {
    if node == nil {
        return &TreeNode[T]{Value: value}
    }

    if value < node.Value {
        node.Left = bt.insertNode(node.Left, value)
    } else if value > node.Value {
        node.Right = bt.insertNode(node.Right, value)
    }

    return node
}

func (bt *BinaryTree[T]) Search(value T) bool {
    return bt.searchNode(bt.root, value)
}

func (bt *BinaryTree[T]) searchNode(node *TreeNode[T], value T) bool {
    if node == nil {
        return false
    }

    if value == node.Value {
        return true
    } else if value < node.Value {
        return bt.searchNode(node.Left, value)
    } else {
        return bt.searchNode(node.Right, value)
    }
}

func (bt *BinaryTree[T]) InOrderTraversal() []T {
    var result []T
    bt.inOrder(bt.root, &result)
    return result
}

func (bt *BinaryTree[T]) inOrder(node *TreeNode[T], result *[]T) {
    if node != nil {
        bt.inOrder(node.Left, result)
        *result = append(*result, node.Value)
        bt.inOrder(node.Right, result)
    }
}
```

### Generic Cache
```go
import (
    "sync"
    "time"
)

type CacheItem[T any] struct {
    Value     T
    ExpiresAt time.Time
}

type Cache[K comparable, V any] struct {
    items map[K]CacheItem[V]
    mutex sync.RWMutex
}

func NewCache[K comparable, V any]() *Cache[K, V] {
    return &Cache[K, V]{
        items: make(map[K]CacheItem[V]),
    }
}

func (c *Cache[K, V]) Set(key K, value V, ttl time.Duration) {
    c.mutex.Lock()
    defer c.mutex.Unlock()

    expiresAt := time.Now().Add(ttl)
    c.items[key] = CacheItem[V]{
        Value:     value,
        ExpiresAt: expiresAt,
    }
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
    c.mutex.RLock()
    defer c.mutex.RUnlock()

    item, exists := c.items[key]
    if !exists {
        var zero V
        return zero, false
    }

    if time.Now().After(item.ExpiresAt) {
        go c.Delete(key) // Clean up expired item
        var zero V
        return zero, false
    }

    return item.Value, true
}

func (c *Cache[K, V]) Delete(key K) {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    delete(c.items, key)
}

func (c *Cache[K, V]) Clear() {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    c.items = make(map[K]CacheItem[V])
}

func (c *Cache[K, V]) Size() int {
    c.mutex.RLock()
    defer c.mutex.RUnlock()
    return len(c.items)
}
```

## Practical Use Cases

### Generic HTTP Response Handler
```go
// Generic API response structure
type APIResponse[T any] struct {
    Data    T      `json:"data,omitempty"`
    Error   string `json:"error,omitempty"`
    Success bool   `json:"success"`
}

// Generic response handler
func WriteJSONResponse[T any](w http.ResponseWriter, status int, data T) error {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)

    response := APIResponse[T]{
        Data:    data,
        Success: status < 400,
    }

    return json.NewEncoder(w).Encode(response)
}

func WriteErrorResponse(w http.ResponseWriter, status int, message string) error {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)

    response := APIResponse[interface{}]{
        Error:   message,
        Success: false,
    }

    return json.NewEncoder(w).Encode(response)
}

// Usage in handlers
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
    user := User{ID: 1, Name: "John", Email: "john@example.com"}
    WriteJSONResponse(w, http.StatusOK, user)
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
    users := []User{
        {ID: 1, Name: "John", Email: "john@example.com"},
        {ID: 2, Name: "Jane", Email: "jane@example.com"},
    }
    WriteJSONResponse(w, http.StatusOK, users)
}
```

### Generic Database Repository
```go
// Generic repository interface
type Repository[T any, ID comparable] interface {
    Create(T) error
    GetByID(ID) (T, error)
    Update(T) error
    Delete(ID) error
    List() ([]T, error)
}

// Generic memory repository implementation
type MemoryRepository[T Entity[ID], ID comparable] struct {
    data   map[ID]T
    nextID ID
    mutex  sync.RWMutex
}

type Entity[ID comparable] interface {
    GetID() ID
    SetID(ID)
}

func NewMemoryRepository[T Entity[ID], ID comparable]() *MemoryRepository[T, ID] {
    return &MemoryRepository[T, ID]{
        data: make(map[ID]T),
    }
}

func (r *MemoryRepository[T, ID]) Create(entity T) error {
    r.mutex.Lock()
    defer r.mutex.Unlock()

    // Simple ID generation (works for int types)
    var id ID
    switch any(id).(type) {
    case int:
        newID := any(len(r.data) + 1).(ID)
        entity.SetID(newID)
    case string:
        newID := any(fmt.Sprintf("id_%d", len(r.data)+1)).(ID)
        entity.SetID(newID)
    }

    r.data[entity.GetID()] = entity
    return nil
}

func (r *MemoryRepository[T, ID]) GetByID(id ID) (T, error) {
    r.mutex.RLock()
    defer r.mutex.RUnlock()

    entity, exists := r.data[id]
    if !exists {
        var zero T
        return zero, fmt.Errorf("entity with ID %v not found", id)
    }
    return entity, nil
}

func (r *MemoryRepository[T, ID]) List() ([]T, error) {
    r.mutex.RLock()
    defer r.mutex.RUnlock()

    entities := make([]T, 0, len(r.data))
    for _, entity := range r.data {
        entities = append(entities, entity)
    }
    return entities, nil
}

// Usage example
type Product struct {
    ID    int     `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}

func (p *Product) GetID() int      { return p.ID }
func (p *Product) SetID(id int)    { p.ID = id }

func main() {
    productRepo := NewMemoryRepository[*Product, int]()

    product := &Product{Name: "Laptop", Price: 999.99}
    productRepo.Create(product)

    found, err := productRepo.GetByID(1)
    if err == nil {
        fmt.Printf("Found product: %s - $%.2f\n", found.Name, found.Price)
    }
}
```

### Generic Validation Framework
```go
// Generic validator interface
type Validator[T any] interface {
    Validate(T) []ValidationError
}

type ValidationError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}

// Generic validation rule
type Rule[T any] struct {
    Field     string
    Predicate func(T) bool
    Message   string
}

// Generic validator implementation
type StructValidator[T any] struct {
    rules []Rule[T]
}

func NewValidator[T any]() *StructValidator[T] {
    return &StructValidator[T]{
        rules: make([]Rule[T], 0),
    }
}

func (v *StructValidator[T]) AddRule(field string, predicate func(T) bool, message string) *StructValidator[T] {
    v.rules = append(v.rules, Rule[T]{
        Field:     field,
        Predicate: predicate,
        Message:   message,
    })
    return v
}

func (v *StructValidator[T]) Validate(value T) []ValidationError {
    var errors []ValidationError

    for _, rule := range v.rules {
        if !rule.Predicate(value) {
            errors = append(errors, ValidationError{
                Field:   rule.Field,
                Message: rule.Message,
            })
        }
    }

    return errors
}

// Usage example
type CreateUserRequest struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Age      int    `json:"age"`
    Password string `json:"password"`
}

func createUserValidator() *StructValidator[CreateUserRequest] {
    return NewValidator[CreateUserRequest]().
        AddRule("name", func(req CreateUserRequest) bool {
            return len(req.Name) >= 2
        }, "Name must be at least 2 characters").
        AddRule("email", func(req CreateUserRequest) bool {
            return strings.Contains(req.Email, "@")
        }, "Email must contain @").
        AddRule("age", func(req CreateUserRequest) bool {
            return req.Age >= 0 && req.Age <= 150
        }, "Age must be between 0 and 150").
        AddRule("password", func(req CreateUserRequest) bool {
            return len(req.Password) >= 8
        }, "Password must be at least 8 characters")
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        WriteErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
        return
    }

    validator := createUserValidator()
    errors := validator.Validate(req)

    if len(errors) > 0 {
        WriteJSONResponse(w, http.StatusBadRequest, map[string]interface{}{
            "validation_errors": errors,
        })
        return
    }

    // Process valid request...
    WriteJSONResponse(w, http.StatusCreated, map[string]string{
        "message": "User created successfully",
    })
}
```

## Performance Considerations

### Compile-Time vs Runtime
```go
// Generics are compile-time only
// Each instantiation creates a separate compiled function
func Print[T any](value T) {
    fmt.Println(value)
}

// At compile time, Go generates:
// func Print_string(value string) { fmt.Println(value) }
// func Print_int(value int) { fmt.Println(value) }
// etc. for each type used

// This means no runtime overhead, but larger binary size
```

### When to Use Generics
```go
// Good use cases:
// 1. Data structures (containers, collections)
type Set[T comparable] map[T]struct{}

// 2. Algorithms that work on multiple types
func Sort[T Ordered](slice []T) { /* ... */ }

// 3. Utility functions
func Ptr[T any](value T) *T { return &value }

// 4. API helpers
func Must[T any](value T, err error) T {
    if err != nil {
        panic(err)
    }
    return value
}
```

### When NOT to Use Generics
```go
// Don't use generics for:
// 1. Simple functions with interface{} that work fine
func PrintAny(value interface{}) {
    fmt.Println(value)  // This is fine for simple cases
}

// 2. One-off functions used in single place
func processSpecificData(data []SpecificType) {
    // No need for generics if only used once
}

// 3. When interfaces provide better abstraction
func ProcessReader(r io.Reader) error {
    // io.Reader is better than generic constraint here
    return nil
}
```

### Binary Size Impact
```go
// Each generic instantiation increases binary size
func GenericFunction[T any](value T) T { return value }

// Usage creates separate compiled functions:
GenericFunction[int](42)        // Generates int version
GenericFunction[string]("hi")   // Generates string version
GenericFunction[User](user)     // Generates User version

// Balance: Use generics for reusable code, not one-offs
```

## Best Practices

### Naming Conventions
```go
// Use single letters for simple cases
func Identity[T any](value T) T { return value }
func Pair[T, U any](first T, second U) (T, U) { return first, second }

// Use descriptive names for complex cases
func Repository[Entity, Key comparable] interface {
    Get(Key) (Entity, error)
    Save(Entity) error
}

// Follow Go conventions
type Map[Key comparable, Value any] map[Key]Value
type Slice[Element any] []Element
```

### Constraint Design
```go
// Start with broader constraints, narrow as needed
func Process[T any](value T) T {  // Start with 'any'
    // Implementation
    return value
}

// Narrow down when you need specific operations
func Compare[T comparable](a, b T) bool {  // Need equality
    return a == b
}

func Add[T Numeric](a, b T) T {  // Need arithmetic
    return a + b
}

// Create meaningful constraint interfaces
type Persistable interface {
    GetID() string
    SetID(string)
    Validate() error
}

func SaveEntity[T Persistable](entity T) error {
    if err := entity.Validate(); err != nil {
        return err
    }
    // Save logic
    return nil
}
```

### Error Handling with Generics
```go
// Generic result type for better error handling
type Result[T any] struct {
    Value T
    Error error
}

func NewResult[T any](value T, err error) Result[T] {
    return Result[T]{Value: value, Error: err}
}

func (r Result[T]) IsOK() bool {
    return r.Error == nil
}

func (r Result[T]) Unwrap() T {
    if r.Error != nil {
        panic(r.Error)
    }
    return r.Value
}

func (r Result[T]) UnwrapOr(defaultValue T) T {
    if r.Error != nil {
        return defaultValue
    }
    return r.Value
}

// Usage
func FetchUser(id int) Result[User] {
    user, err := getUserFromDB(id)
    return NewResult(user, err)
}

func main() {
    result := FetchUser(123)
    if result.IsOK() {
        user := result.Unwrap()
        fmt.Printf("Found user: %s\n", user.Name)
    } else {
        fmt.Printf("Error: %v\n", result.Error)
    }
}
```

### Testing Generic Code
```go
// Test with multiple types
func TestGenericSort(t *testing.T) {
    // Test with ints
    intSlice := []int{3, 1, 4, 1, 5}
    Sort(intSlice)
    expected := []int{1, 1, 3, 4, 5}
    if !slices.Equal(intSlice, expected) {
        t.Errorf("Expected %v, got %v", expected, intSlice)
    }

    // Test with strings
    stringSlice := []string{"banana", "apple", "cherry"}
    Sort(stringSlice)
    expectedStrings := []string{"apple", "banana", "cherry"}
    if !slices.Equal(stringSlice, expectedStrings) {
        t.Errorf("Expected %v, got %v", expectedStrings, stringSlice)
    }
}

// Use table-driven tests for multiple types
func TestGenericMax(t *testing.T) {
    tests := []struct {
        name     string
        testFunc func(*testing.T)
    }{
        {"int", testMaxInt},
        {"float64", testMaxFloat64},
        {"string", testMaxString},
    }

    for _, tt := range tests {
        t.Run(tt.name, tt.testFunc)
    }
}

func testMaxInt(t *testing.T) {
    result, ok := Max([]int{1, 3, 2})
    if !ok {
        t.Fatal("Expected ok to be true")
    }
    if result != 3 {
        t.Errorf("Expected 3, got %d", result)
    }
}

func testMaxFloat64(t *testing.T) {
    result, ok := Max([]float64{1.1, 3.3, 2.2})
    if !ok {
        t.Fatal("Expected ok to be true")
    }
    if result != 3.3 {
        t.Errorf("Expected 3.3, got %f", result)
    }
}
```

### Documentation
```go
// Document type parameters and constraints clearly
// Max returns the maximum value from a slice of ordered values.
// T must be a type that supports comparison operators.
// Returns the maximum value and true if the slice is not empty,
// or the zero value and false if the slice is empty.
func Max[T Ordered](values []T) (T, bool) {
    if len(values) == 0 {
        var zero T
        return zero, false
    }

    max := values[0]
    for _, v := range values[1:] {
        if v > max {
            max = v
        }
    }
    return max, true
}

// Repository provides CRUD operations for entities.
// T must implement Entity interface with GetID/SetID methods.
// ID must be a comparable type suitable for use as map keys.
type Repository[T Entity[ID], ID comparable] interface {
    Create(T) error
    GetByID(ID) (T, error)
    Update(T) error
    Delete(ID) error
}
```

---

*Go generics provide type safety and reusability while maintaining Go's simplicity and performance. Use them judiciously to create clean, reusable APIs and data structures without sacrificing readability.*