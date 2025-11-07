# Go Error Handling Guide
*For JavaScript/TypeScript Developers*

## Table of Contents
1. [Philosophy Differences](#philosophy-differences)
2. [Basic Error Handling](#basic-error-handling)
3. [Creating Errors](#creating-errors)
4. [Custom Error Types](#custom-error-types)
5. [Error Wrapping](#error-wrapping)
6. [Best Practices](#best-practices)
7. [Common Patterns](#common-patterns)
8. [Quick Reference](#quick-reference)

## Philosophy Differences

### JavaScript/TypeScript Approach
```javascript
// JavaScript uses try/catch with exceptions
try {
    const result = riskyOperation();
    console.log(result);
} catch (error) {
    console.error("Something went wrong:", error.message);
} finally {
    // cleanup code
}

// Async operations with promises
async function fetchData() {
    try {
        const response = await fetch('/api/data');
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        return await response.json();
    } catch (error) {
        console.error('Fetch failed:', error);
        throw error; // re-throw
    }
}
```

### Go Approach
```go
// Go uses explicit error returns (no exceptions!)
result, err := riskyOperation()
if err != nil {
    fmt.Printf("Something went wrong: %v\n", err)
    return // or handle the error
}
fmt.Println(result)

// Functions return errors as second value
func fetchData() ([]byte, error) {
    response, err := http.Get("/api/data")
    if err != nil {
        return nil, fmt.Errorf("fetch failed: %w", err)
    }
    defer response.Body.Close() // cleanup (like finally)

    if response.StatusCode != 200 {
        return nil, fmt.Errorf("HTTP error: %d", response.StatusCode)
    }

    return io.ReadAll(response.Body)
}
```

**Key Difference:** Go treats errors as values, not exceptions!

## Basic Error Handling

### The Error Interface

**Go:**
```go
// The built-in error interface
type error interface {
    Error() string
}

// Functions that can fail return (result, error)
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

// Always check errors!
result, err := divide(10, 0)
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}
fmt.Printf("Result: %f\n", result)
```

**JavaScript Equivalent:**
```javascript
// JavaScript would throw an exception
function divide(a, b) {
    if (b === 0) {
        throw new Error("Division by zero");
    }
    return a / b;
}

try {
    const result = divide(10, 0);
    console.log(`Result: ${result}`);
} catch (error) {
    console.error(`Error: ${error.message}`);
}

// TypeScript with explicit error returns (Go-style)
function divideGoStyle(a: number, b: number): [number, Error | null] {
    if (b === 0) {
        return [0, new Error("Division by zero")];
    }
    return [a / b, null];
}

const [result, err] = divideGoStyle(10, 0);
if (err) {
    console.error(`Error: ${err.message}`);
} else {
    console.log(`Result: ${result}`);
}
```

### Multiple Return Values

**Go:**
```go
// Multiple return values are common in Go
func processFile(filename string) (string, int, error) {
    content, err := os.ReadFile(filename)
    if err != nil {
        return "", 0, fmt.Errorf("failed to read file: %w", err)
    }

    lines := strings.Split(string(content), "\n")
    return string(content), len(lines), nil
}

// Usage
content, lineCount, err := processFile("data.txt")
if err != nil {
    log.Printf("Error processing file: %v", err)
    return
}
fmt.Printf("File has %d lines\n", lineCount)
```

**JavaScript Equivalent:**
```javascript
// JavaScript would typically use objects or arrays
async function processFile(filename) {
    try {
        const content = await fs.readFile(filename, 'utf8');
        const lines = content.split('\n');
        return {
            content: content,
            lineCount: lines.length,
            error: null
        };
    } catch (error) {
        return {
            content: "",
            lineCount: 0,
            error: error
        };
    }
}

// Usage
const result = await processFile("data.txt");
if (result.error) {
    console.error("Error processing file:", result.error.message);
    return;
}
console.log(`File has ${result.lineCount} lines`);
```

## Creating Errors

### Basic Error Creation

**Go:**
```go
import (
    "errors"
    "fmt"
)

// Simple error
err1 := errors.New("something went wrong")

// Formatted error
err2 := fmt.Errorf("user %s not found", username)

// Error with formatting and wrapping
err3 := fmt.Errorf("failed to connect to database: %w", originalErr)
```

**JavaScript Equivalent:**
```javascript
// Simple error
const err1 = new Error("Something went wrong");

// Formatted error (using template literals)
const err2 = new Error(`User ${username} not found`);

// Error with cause (modern JavaScript/Node.js)
const err3 = new Error("Failed to connect to database", { cause: originalErr });
```

### Custom Error Types

**Go:**
```go
// Custom error type
type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation failed for field '%s': %s", e.Field, e.Message)
}

// Usage
func validateEmail(email string) error {
    if !strings.Contains(email, "@") {
        return ValidationError{
            Field:   "email",
            Message: "must contain @ symbol",
        }
    }
    return nil
}

// Type assertion to access custom fields
err := validateEmail("invalid-email")
if err != nil {
    if validationErr, ok := err.(ValidationError); ok {
        fmt.Printf("Field: %s, Message: %s\n", validationErr.Field, validationErr.Message)
    }
}
```

**JavaScript Equivalent:**
```javascript
// Custom error class
class ValidationError extends Error {
    constructor(field, message) {
        super(`Validation failed for field '${field}': ${message}`);
        this.name = 'ValidationError';
        this.field = field;
        this.validationMessage = message;
    }
}

// Usage
function validateEmail(email) {
    if (!email.includes('@')) {
        throw new ValidationError('email', 'must contain @ symbol');
    }
}

try {
    validateEmail('invalid-email');
} catch (error) {
    if (error instanceof ValidationError) {
        console.log(`Field: ${error.field}, Message: ${error.validationMessage}`);
    }
}
```

## Error Wrapping

### Go 1.13+ Error Wrapping

**Go:**
```go
import (
    "errors"
    "fmt"
)

func readConfig() error {
    _, err := os.Open("config.json")
    if err != nil {
        // Wrap the error with context
        return fmt.Errorf("failed to read config: %w", err)
    }
    return nil
}

func startServer() error {
    err := readConfig()
    if err != nil {
        // Wrap again with more context
        return fmt.Errorf("server startup failed: %w", err)
    }
    return nil
}

// Check for specific errors in the chain
err := startServer()
if err != nil {
    // Check if the error chain contains a specific error
    if errors.Is(err, os.ErrNotExist) {
        fmt.Println("Config file doesn't exist")
    }

    // Unwrap to get the original error
    var pathErr *os.PathError
    if errors.As(err, &pathErr) {
        fmt.Printf("Path error: %s\n", pathErr.Path)
    }

    fmt.Printf("Full error: %v\n", err)
}
```

**JavaScript Equivalent:**
```javascript
// JavaScript error chaining (newer Node.js/browsers)
function readConfig() {
    try {
        // Simulate file reading
        throw new Error("ENOENT: no such file or directory");
    } catch (error) {
        throw new Error("Failed to read config", { cause: error });
    }
}

function startServer() {
    try {
        readConfig();
    } catch (error) {
        throw new Error("Server startup failed", { cause: error });
    }
}

// Usage
try {
    startServer();
} catch (error) {
    console.log("Full error:", error.message);

    // Walk the error chain
    let currentError = error;
    while (currentError.cause) {
        currentError = currentError.cause;
        console.log("Caused by:", currentError.message);
    }
}
```

## Best Practices

### Go Error Handling Patterns

**Go:**
```go
// 1. Always check errors immediately
result, err := someOperation()
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}

// 2. Don't ignore errors
_, _ = someOperation() // BAD: ignoring error

// 3. Handle errors at the right level
func processData() error {
    data, err := fetchData()
    if err != nil {
        // Log and wrap, don't just return
        log.Printf("Failed to fetch data: %v", err)
        return fmt.Errorf("data processing failed: %w", err)
    }

    return processRawData(data)
}

// 4. Use defer for cleanup
func processFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close() // Always cleanup

    return processFileContent(file)
}

// 5. Sentinel errors for expected conditions
var ErrNotFound = errors.New("item not found")

func findUser(id int) (*User, error) {
    // ... search logic
    if userNotFound {
        return nil, ErrNotFound
    }
    return user, nil
}

// Usage
user, err := findUser(123)
if errors.Is(err, ErrNotFound) {
    // Handle expected "not found" case
    return createDefaultUser()
}
```

**JavaScript Equivalent:**
```javascript
// 1. Always handle promise rejections
try {
    const result = await someOperation();
    // process result
} catch (error) {
    throw new Error(`Operation failed: ${error.message}`);
}

// 2. Don't ignore errors
someOperation().catch(() => {}); // BAD: ignoring error

// 3. Handle errors at the right level
async function processData() {
    try {
        const data = await fetchData();
        return await processRawData(data);
    } catch (error) {
        console.error("Failed to fetch data:", error);
        throw new Error(`Data processing failed: ${error.message}`);
    }
}

// 4. Use finally for cleanup
async function processFile(filename) {
    let file;
    try {
        file = await fs.open(filename);
        return await processFileContent(file);
    } catch (error) {
        throw new Error(`Failed to process file: ${error.message}`);
    } finally {
        if (file) {
            await file.close();
        }
    }
}

// 5. Custom error types for expected conditions
class NotFoundError extends Error {
    constructor(message) {
        super(message);
        this.name = 'NotFoundError';
    }
}

async function findUser(id) {
    // ... search logic
    if (userNotFound) {
        throw new NotFoundError("User not found");
    }
    return user;
}

// Usage
try {
    const user = await findUser(123);
} catch (error) {
    if (error instanceof NotFoundError) {
        return createDefaultUser();
    }
    throw error;
}
```

## Common Patterns

### Early Return Pattern

**Go:**
```go
func validateAndProcess(data string) error {
    if data == "" {
        return errors.New("data cannot be empty")
    }

    if len(data) > 100 {
        return errors.New("data too long")
    }

    processed, err := processData(data)
    if err != nil {
        return fmt.Errorf("processing failed: %w", err)
    }

    err = saveData(processed)
    if err != nil {
        return fmt.Errorf("saving failed: %w", err)
    }

    return nil
}
```

**JavaScript Equivalent:**
```javascript
async function validateAndProcess(data) {
    if (!data) {
        throw new Error("Data cannot be empty");
    }

    if (data.length > 100) {
        throw new Error("Data too long");
    }

    try {
        const processed = await processData(data);
        await saveData(processed);
    } catch (error) {
        throw new Error(`Processing failed: ${error.message}`);
    }
}
```

### Error Aggregation

**Go:**
```go
func processMultipleFiles(filenames []string) error {
    var errors []error

    for _, filename := range filenames {
        err := processFile(filename)
        if err != nil {
            errors = append(errors, fmt.Errorf("file %s: %w", filename, err))
        }
    }

    if len(errors) > 0 {
        return fmt.Errorf("processing failed for %d files: %v", len(errors), errors)
    }

    return nil
}
```

**JavaScript Equivalent:**
```javascript
async function processMultipleFiles(filenames) {
    const errors = [];

    for (const filename of filenames) {
        try {
            await processFile(filename);
        } catch (error) {
            errors.push(new Error(`File ${filename}: ${error.message}`));
        }
    }

    if (errors.length > 0) {
        const combinedMessage = errors.map(e => e.message).join(', ');
        throw new Error(`Processing failed for ${errors.length} files: ${combinedMessage}`);
    }
}
```

## Quick Reference

| Concept | Go | JavaScript/TypeScript |
|---------|----|-----------------------|
| Error Creation | `errors.New("message")` | `new Error("message")` |
| Formatted Error | `fmt.Errorf("user %s not found", name)` | `new Error(\`user ${name} not found\`)` |
| Error Checking | `if err != nil { ... }` | `try/catch` or `if (error) { ... }` |
| Error Wrapping | `fmt.Errorf("context: %w", err)` | `new Error("context", { cause: err })` |
| Custom Errors | Implement `Error() string` | Extend `Error` class |
| No Error | `return result, nil` | `return result` (no throw) |
| Multiple Returns | `func() (T, error)` | `async function(): Promise<T>` |
| Error Chain Check | `errors.Is(err, target)` | `error instanceof ErrorType` |

## Key Takeaways for JS Developers

1. **Errors are Values:** Not exceptions! Functions return them explicitly
2. **Always Check:** Unlike exceptions that you can ignore, Go errors must be handled
3. **Multiple Returns:** Go functions commonly return `(result, error)`
4. **No Try/Catch:** Use `if err != nil` instead
5. **Error Wrapping:** Use `%w` verb to maintain error chains
6. **Explicit Handling:** Makes error paths visible in the code
7. **defer for Cleanup:** Like `finally` but more explicit and powerful

## Common Mistakes JS Developers Make

```go
// ** Ignoring errors (don't do this!)
result, _ := someOperation()

// ** Not wrapping errors
if err != nil {
    return err  // loses context
}

// ** Proper error handling
result, err := someOperation()
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}

// ** Early return pattern
if err != nil {
    return fmt.Errorf("validation failed: %w", err)
}
// continue with happy path
```

---

*Remember: Go's explicit error handling might feel verbose at first, but it makes your code more reliable and error conditions more visible!*