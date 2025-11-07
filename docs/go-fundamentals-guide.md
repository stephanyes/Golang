# Go Fundamentals Guide
*For JavaScript/TypeScript Developers*

## Table of Contents
1. [Variable Declarations](#variable-declarations)
2. [Constants](#constants)
3. [Data Types](#data-types)
4. [Type Inference](#type-inference)
5. [Zero Values](#zero-values)
6. [Quick Reference](#quick-reference)

## Variable Declarations

### The `var` Keyword

**Go:**
```go
// Explicit type declaration
var name string = "John"
var age int = 25

// Type inference (Go infers the type)
var name = "John"
var age = 25

// Declaration without initialization (gets zero value)
var name string  // ""
var age int      // 0
```

**JavaScript Equivalent:**
```javascript
// JavaScript (no type annotations)
let name = "John";
let age = 25;

// TypeScript equivalent
let name: string = "John";
let age: number = 25;
```

### The `:=` Short Declaration (Go's Special Feature)

**Go:**
```go
// Short variable declaration (only inside functions)
name := "John"        // equivalent to: var name = "John"
age := 25            // equivalent to: var age = 25
isActive := true     // equivalent to: var isActive = true

// Multiple assignments
name, age := "John", 25
```

**JavaScript Equivalent:**
```javascript
// JavaScript doesn't have := but similar concept
const name = "John";
const age = 25;
const isActive = true;

// Destructuring assignment (similar to multiple assignment)
const [name, age] = ["John", 25];
```

**Important:** `:=` can only be used inside functions, not at package level.

## Constants

### Basic Constants

**Go:**
```go
// Typed constants
const PI float64 = 3.14159
const MaxUsers int = 100

// Untyped constants (more flexible)
const PI = 3.14159
const Greeting = "Hello"
const IsDebug = true
```

**JavaScript Equivalent:**
```javascript
// JavaScript
const PI = 3.14159;
const MAX_USERS = 100;
const GREETING = "Hello";
const IS_DEBUG = true;

// TypeScript
const PI: number = 3.14159;
const MAX_USERS: number = 100;
```

### Constant Blocks

**Go:**
```go
const (
    StatusPending = "pending"
    StatusActive  = "active"
    StatusClosed  = "closed"
)

// With iota (auto-incrementing)
const (
    Red = iota    // 0
    Green         // 1
    Blue          // 2
)
```

**JavaScript Equivalent:**
```javascript
// JavaScript - using object or enum
const Status = {
    PENDING: "pending",
    ACTIVE: "active",
    CLOSED: "closed"
};

// TypeScript enum
enum Status {
    Pending = "pending",
    Active = "active",
    Closed = "closed"
}

enum Color {
    Red,    // 0
    Green,  // 1
    Blue    // 2
}
```

## Data Types

### Numeric Types

**Go:**
```go
// Integers
var smallInt int8 = 127          // -128 to 127
var regularInt int = 42          // Platform dependent (32 or 64 bit)
var bigInt int64 = 9223372036854775807

// Unsigned integers
var smallUint uint8 = 255        // 0 to 255
var regularUint uint = 42        // Platform dependent
var bigUint uint64 = 18446744073709551615

// Floating point
var smallFloat float32 = 3.14
var bigFloat float64 = 3.141592653589793

// Complex numbers (JavaScript doesn't have these!)
var complexNum complex64 = 1 + 2i
var bigComplex complex128 = 1.5 + 2.5i
```

**JavaScript Equivalent:**
```javascript
// JavaScript - only has 'number' type (64-bit float)
let anyNumber = 42;
let floatNumber = 3.14;
let bigNumber = 9223372036854775807;

// TypeScript - same as JavaScript for numbers
let anyNumber: number = 42;

// For very large integers, JavaScript has BigInt
let veryBigInt: bigint = 9223372036854775807n;
```

### String and Character Types

**Go:**
```go
// Strings (UTF-8 encoded)
var message string = "Hello, 世界"
var multiline string = `This is a
multi-line string
with backticks`

// Individual characters (runes)
var char rune = 'A'              // Single quotes for runes
var unicodeChar rune = '世'       // Unicode support

// Byte (uint8)
var byteVal byte = 65            // ASCII 'A'
```

**JavaScript Equivalent:**
```javascript
// JavaScript
let message = "Hello, 世界";
let multiline = `This is a
multi-line string
with template literals`;

let char = 'A';                  // No distinction between char and string
let unicodeChar = '世';

// No direct byte type, but you can use:
let byteVal = 65;
```

### Boolean Type

**Go:**
```go
var isActive bool = true
var isComplete bool = false

// Only true/false, no truthy/falsy like JavaScript
```

**JavaScript Equivalent:**
```javascript
// JavaScript (truthy/falsy values exist)
let isActive = true;
let isComplete = false;

// These are falsy in JavaScript but would be errors in Go:
// false, 0, "", null, undefined, NaN
```

### Arrays and Slices

**Go:**
```go
// Arrays (fixed size)
var numbers [5]int = [5]int{1, 2, 3, 4, 5}
var fruits [3]string = [3]string{"apple", "banana", "orange"}

// Slices (dynamic arrays)
var dynamicNumbers []int = []int{1, 2, 3}
var dynamicFruits []string = []string{"apple", "banana"}

// Using make for slices
var slice []int = make([]int, 5)    // length 5, all zeros
```

**JavaScript Equivalent:**
```javascript
// JavaScript - all arrays are dynamic
let numbers = [1, 2, 3, 4, 5];
let fruits = ["apple", "banana", "orange"];

// TypeScript with type annotations
let numbers: number[] = [1, 2, 3, 4, 5];
let fruits: string[] = ["apple", "banana"];
```

### Maps (Objects)

**Go:**
```go
// Maps (like JavaScript objects)
var userAges map[string]int = map[string]int{
    "John":  25,
    "Alice": 30,
}

// Using make
var scores map[string]int = make(map[string]int)
scores["player1"] = 100
```

**JavaScript Equivalent:**
```javascript
// JavaScript objects
let userAges = {
    "John": 25,
    "Alice": 30
};

// Map object (closer to Go maps)
let scores = new Map();
scores.set("player1", 100);

// TypeScript
let userAges: { [key: string]: number } = {
    "John": 25,
    "Alice": 30
};
```

## Type Inference

**Go:**
```go
// Go infers types when using := or var without explicit type
name := "John"          // string
age := 25              // int
price := 19.99         // float64
isActive := true       // bool
```

**JavaScript/TypeScript:**
```javascript
// JavaScript - dynamic typing
let name = "John";      // can change type later
name = 42;             // allowed

// TypeScript - type inference
let name = "John";      // inferred as string
// name = 42;          // Error: Type 'number' not assignable to 'string'
```

## Zero Values

**Go Concept:** Variables declared without initialization get "zero values"

**Go:**
```go
var text string      // ""
var number int       // 0
var flag bool        // false
var pointer *int     // nil
var slice []int      // nil
var mapVar map[string]int  // nil
```

**JavaScript Equivalent:**
```javascript
// JavaScript - uninitialized variables are 'undefined'
let text;               // undefined
let number;             // undefined

// You'd need to explicitly set default values
let text = "";          // empty string
let number = 0;         // zero
let flag = false;       // false
```

## Quick Reference

| Concept | Go | JavaScript/TypeScript |
|---------|----|-----------------------|
| Variable Declaration | `var name string` | `let name: string` |
| Short Declaration | `name := "value"` | `const name = "value"` |
| Constant | `const PI = 3.14` | `const PI = 3.14` |
| Type Inference | `age := 25` | `let age = 25` |
| Arrays | `[5]int{1,2,3,4,5}` | `[1,2,3,4,5]` |
| Dynamic Arrays | `[]int{1,2,3}` | `[1,2,3]` |
| Objects/Maps | `map[string]int` | `{[key: string]: number}` |
| Zero Values | Automatic | `undefined` (manual defaults) |

## Key Takeaways for JS Developers

1. **Static Typing:** Go requires you to think about types upfront
2. **Zero Values:** Go initializes variables automatically (no undefined!)
3. **`:=` is Special:** Only use inside functions, great for quick declarations
4. **No Truthy/Falsy:** Only `true` and `false` are boolean values
5. **Explicit Error Handling:** No exceptions, functions return errors
6. **Simpler but Explicit:** Less "magic" than JavaScript, more explicit code

---

*This guide covers the fundamentals. As you progress, we'll add more advanced topics like pointers, interfaces, goroutines, and channels!*