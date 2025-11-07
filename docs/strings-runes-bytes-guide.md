# Strings, Runes, and Bytes in Go
*For JavaScript/TypeScript Developers*

## Table of Contents
1. [The String Universe in Go vs JavaScript](#the-string-universe-in-go-vs-javascript)
2. [Understanding UTF-8 and Unicode](#understanding-utf-8-and-unicode)
3. [Strings in Go](#strings-in-go)
4. [Runes (Unicode Code Points)](#runes-unicode-code-points)
5. [Bytes (Raw Data)](#bytes-raw-data)
6. [String Iteration Mysteries](#string-iteration-mysteries)
7. [Conversion Between Types](#conversion-between-types)
8. [Practical Use Cases](#practical-use-cases)
9. [Performance Considerations](#performance-considerations)
10. [Best Practices](#best-practices)

## The String Universe in Go vs JavaScript

### JavaScript's Simple World
```javascript
// JavaScript: One string type, everything "just works"
let str = "Hello 世界! 🌍";
console.log(str.length);        // 9 (characters)
console.log(str[0]);           // "H"
console.log(str[6]);           // "世"

// Iteration is straightforward
for (let char of str) {
    console.log(char);  // H, e, l, l, o, space, 世, 界, !, space, 🌍
}

// String manipulation is simple
let upper = str.toUpperCase();
let slice = str.slice(0, 5);   // "Hello"
```

### Go's Complex but Powerful World
```go
// Go: Three related but different types!
str := "Hello 世界! 🌍"

// 1. String (UTF-8 encoded bytes)
fmt.Println(len(str))                    // 16 (BYTES, not characters!)
fmt.Println(str[0])                      // 72 (byte value of 'H')
// fmt.Println(str[6])                   // 228 (first byte of '世', not the character!)

// 2. Runes (Unicode code points)
runes := []rune(str)
fmt.Println(len(runes))                  // 9 (actual characters)
fmt.Println(runes[0])                    // 72 ('H' as Unicode code point)
fmt.Println(runes[6])                    // 19990 ('世' as Unicode code point)

// 3. Bytes (raw data)
bytes := []byte(str)
fmt.Println(len(bytes))                  // 16 (raw bytes)
fmt.Println(bytes[0])                    // 72 (byte value of 'H')
```

**Key Insight:** JavaScript hides the complexity, Go exposes it for performance and control.

## Understanding UTF-8 and Unicode

### Unicode Basics (What JavaScript Hides)
```
Unicode assigns numbers to characters:
'A' = U+0041 (65)
'世' = U+4E16 (19990)
'🌍' = U+1F30D (127757)

UTF-8 encoding (how these numbers are stored as bytes):
'A' = 1 byte:  [65]
'世' = 3 bytes: [228, 184, 150]
'🌍' = 4 bytes: [240, 159, 140, 141]
```

### Memory Layout Visualization
```go
str := "A世🌍"

// String view: "A世🌍"
// Rune view:   ['A', '世', '🌍']  (3 characters)
// Byte view:   [65, 228, 184, 150, 240, 159, 140, 141]  (8 bytes)

/*
Memory layout:
Address:  0x100  0x101  0x102  0x103  0x104  0x105  0x106  0x107
Byte:     [ 65 ] [228 ] [184 ] [150 ] [240 ] [159 ] [140 ] [141 ]
Char:     [ A  ] [    世    ] [         🌍         ]
*/
```

**JavaScript Equivalent:**
```javascript
// JavaScript abstracts this away
let str = "A世🌍";
console.log(str.length);  // 4 (but this is wrong for emoji!)

// Modern JavaScript has some UTF-16 quirks with emoji
console.log([...str].length);  // 3 (correct character count)
```

## Strings in Go

### String Properties
```go
// Strings are immutable byte slices with UTF-8 encoding
str := "Hello, 世界"

// Properties
fmt.Println(len(str))           // 13 (bytes)
fmt.Println(utf8.RuneCountInString(str))  // 9 (characters/runes)

// Strings are immutable
// str[0] = 'h'  // ** Compile error!

// Indexing gives bytes, not characters
fmt.Printf("str[0] = %d ('%c')\n", str[0], str[0])     // 72 ('H')
fmt.Printf("str[7] = %d\n", str[7])                    // 228 (first byte of '世')

// Slicing works on byte boundaries
fmt.Println(str[0:5])           // "Hello" (safe - ASCII)
// fmt.Println(str[0:8])        // "Hello, ä¸" (broken! cuts through '世')
```

**JavaScript Equivalent:**
```javascript
let str = "Hello, 世界";
console.log(str.length);        // 9 (characters, not bytes)
console.log(str[0]);           // "H" (character, not byte value)
console.log(str.slice(0, 5));  // "Hello" (character-based slicing)

// JavaScript strings are mutable in a sense
str = str.replace('H', 'h');   // Creates new string
```

### String Literals and Raw Strings
```go
// Regular string literals (interpret escapes)
str1 := "Hello\nWorld\t🌍"
str2 := "Path: C:\\Users\\name"  // Need to escape backslashes

// Raw string literals (backticks - no escapes)
str3 := `Hello
World	🌍`  // Literal newline and tab

str4 := `Path: C:\Users\name`    // No need to escape

// Multiline strings
sql := `
SELECT id, name, email
FROM users
WHERE active = true
ORDER BY name
`
```

**JavaScript Equivalent:**
```javascript
// Regular strings
let str1 = "Hello\nWorld\t🌍";
let str2 = "Path: C:\\Users\\name";

// Template literals (similar to Go raw strings)
let str3 = `Hello
World	🌍`;

let str4 = `Path: C:\Users\name`;

// Multiline strings
let sql = `
SELECT id, name, email
FROM users
WHERE active = true
ORDER BY name
`;
```

## Runes (Unicode Code Points)

### What Are Runes?
```go
// rune is an alias for int32
type rune = int32

// Rune literals use single quotes
var r1 rune = 'A'          // 65
var r2 rune = '世'         // 19990
var r3 rune = '🌍'         // 127757

// Unicode escape sequences
var r4 rune = '\u4E16'     // 世 (Unicode U+4E16)
var r5 rune = '\U0001F30D' // 🌍 (Unicode U+1F30D)

fmt.Printf("'A' = %d\n", r1)       // 65
fmt.Printf("'世' = %d\n", r2)      // 19990
fmt.Printf("'🌍' = %d\n", r3)      // 127757
```

### Working with Runes
```go
str := "Hello, 世界! 🌍"

// Convert string to rune slice for character-level operations
runes := []rune(str)
fmt.Println(len(runes))            // 11 (actual characters)

// Safe character access
fmt.Printf("First char: %c (%d)\n", runes[0], runes[0])    // H (72)
fmt.Printf("7th char: %c (%d)\n", runes[7], runes[7])      // 世 (19990)

// Character-level slicing
fmt.Println(string(runes[0:5]))    // "Hello"
fmt.Println(string(runes[7:9]))    // "世界"

// Modify characters (runes are mutable in slices)
runes[0] = 'h'
modified := string(runes)          // "hello, 世界! 🌍"
fmt.Println(modified)
```

**JavaScript Equivalent:**
```javascript
let str = "Hello, 世界! 🌍";

// JavaScript strings give you character access directly
console.log(str.length);           // 11 (characters)
console.log(str[0]);              // "H"
console.log(str[7]);              // "世"

// Character-level slicing works naturally
console.log(str.slice(0, 5));     // "Hello"
console.log(str.slice(7, 9));     // "世界"

// No direct character modification (strings are immutable)
let modified = 'h' + str.slice(1); // "hello, 世界! 🌍"
```

### Unicode Categories and Properties
```go
import "unicode"

func analyzeRune(r rune) {
    fmt.Printf("Rune: %c (%d)\n", r, r)
    fmt.Printf("  IsLetter: %v\n", unicode.IsLetter(r))
    fmt.Printf("  IsDigit: %v\n", unicode.IsDigit(r))
    fmt.Printf("  IsSpace: %v\n", unicode.IsSpace(r))
    fmt.Printf("  IsUpper: %v\n", unicode.IsUpper(r))
    fmt.Printf("  ToUpper: %c\n", unicode.ToUpper(r))
    fmt.Printf("  ToLower: %c\n", unicode.ToLower(r))
}

// Examples
analyzeRune('A')    // Letter, Upper
analyzeRune('世')   // Letter, neither upper nor lower
analyzeRune('5')    // Digit
analyzeRune(' ')    // Space
```

## Bytes (Raw Data)

### Understanding Bytes
```go
str := "Hello, 世界"

// Convert to bytes
bytes := []byte(str)
fmt.Println(len(bytes))            // 13 (raw bytes)
fmt.Println(bytes)                 // [72 101 108 108 111 44 32 228 184 150 231 149 140]

// Byte-level access
fmt.Printf("First byte: %d ('%c')\n", bytes[0], bytes[0])  // 72 ('H')
fmt.Printf("Byte 7: %d\n", bytes[7])                       // 228 (part of '世')

// Modify bytes (bytes are mutable)
bytes[0] = 104  // Change 'H' to 'h'
modified := string(bytes)
fmt.Println(modified)              // "hello, 世界"
```

### Binary Data and File Operations
```go
// Reading binary files
func readBinaryFile(filename string) ([]byte, error) {
    return os.ReadFile(filename)
}

// Writing binary data
func writeBinaryFile(filename string, data []byte) error {
    return os.WriteFile(filename, data, 0644)
}

// Network operations often use bytes
func handleHTTPBody(body io.Reader) ([]byte, error) {
    return io.ReadAll(body)
}

// Base64 encoding/decoding
text := "Hello, 世界"
encoded := base64.StdEncoding.EncodeToString([]byte(text))
fmt.Println(encoded)  // SGVsbG8sIOS4lueVjA==

decoded, _ := base64.StdEncoding.DecodeString(encoded)
fmt.Println(string(decoded))  // "Hello, 世界"
```

**JavaScript Equivalent:**
```javascript
// JavaScript has ArrayBuffer/Uint8Array for binary data
let str = "Hello, 世界";
let encoder = new TextEncoder();
let bytes = encoder.encode(str);
console.log(bytes);  // Uint8Array([72, 101, 108, 108, 111, 44, 32, 228, 184, 150, 231, 149, 140])

let decoder = new TextDecoder();
let decoded = decoder.decode(bytes);
console.log(decoded);  // "Hello, 世界"

// Base64 encoding
let encoded = btoa(str);  // May not work correctly with Unicode
// Use: btoa(unescape(encodeURIComponent(str))) for Unicode
```

## String Iteration Mysteries

### The Range Loop Revelation
```go
str := "Hello, 世界! 🌍"

// Iteration gives different results depending on method!

// Method 1: Index-based (byte iteration) - Usually WRONG
fmt.Println("Byte iteration:")
for i := 0; i < len(str); i++ {
    fmt.Printf("Index %d: byte %d\n", i, str[i])
}
// Output: 16 iterations (bytes), broken characters

// Method 2: Range loop (rune iteration) - Usually CORRECT
fmt.Println("Rune iteration:")
for i, r := range str {
    fmt.Printf("Index %d: rune %c (%d)\n", i, r, r)
}
// Output: 11 iterations (characters), proper Unicode

// Method 3: Manual rune conversion
fmt.Println("Manual rune iteration:")
runes := []rune(str)
for i, r := range runes {
    fmt.Printf("Index %d: rune %c (%d)\n", i, r, r)
}
```

**JavaScript Equivalent:**
```javascript
let str = "Hello, 世界! 🌍";

// JavaScript for...of gives you characters (mostly)
for (let char of str) {
    console.log(char);  // H, e, l, l, o, ,, space, 世, 界, !, space, 🌍
}

// Array.from handles Unicode properly
Array.from(str).forEach((char, index) => {
    console.log(`Index ${index}: ${char}`);
});
```

### When Each Method Is Useful
```go
// Byte iteration: When working with ASCII or binary data
func countASCIILetters(s string) int {
    count := 0
    for i := 0; i < len(s); i++ {
        if s[i] >= 'A' && s[i] <= 'Z' || s[i] >= 'a' && s[i] <= 'z' {
            count++
        }
    }
    return count
}

// Rune iteration: When working with Unicode text (most common)
func countUnicodeLetters(s string) int {
    count := 0
    for _, r := range s {
        if unicode.IsLetter(r) {
            count++
        }
    }
    return count
}

// Mixed: When you need both byte position and character
func findUnicodeChar(s string, target rune) int {
    for byteIndex, r := range s {
        if r == target {
            return byteIndex  // Returns byte position, not character position
        }
    }
    return -1
}
```

## Conversion Between Types

### Safe Conversions
```go
// String ↔ Runes (safe, handles Unicode)
str := "Hello, 世界! 🌍"
runes := []rune(str)              // String to runes
back := string(runes)             // Runes to string
fmt.Println(str == back)          // true

// String ↔ Bytes (safe, preserves UTF-8)
bytes := []byte(str)              // String to bytes
back2 := string(bytes)            // Bytes to string
fmt.Println(str == back2)         // true

// Runes ↔ Bytes (through string)
runeBytes := []byte(string(runes))  // Runes → String → Bytes
```

### Type Conversion Examples
```go
// Converting individual runes
r := '世'
fmt.Printf("Rune: %c (%d)\n", r, r)                    // 世 (19990)
fmt.Printf("As string: %s\n", string(r))               // "世"
fmt.Printf("As bytes: %v\n", []byte(string(r)))        // [228 184 150]

// Converting numbers to strings
num := 42
strNum := strconv.Itoa(num)                             // "42"
backNum, _ := strconv.Atoi(strNum)                      // 42

// Converting runes to numbers
digitRune := '5'
digitValue := digitRune - '0'                           // 5 (int32)
fmt.Printf("'5' as number: %d\n", digitValue)          // 5
```

**JavaScript Equivalent:**
```javascript
// JavaScript conversions are simpler but less explicit
let num = 42;
let strNum = num.toString();        // "42"
let backNum = parseInt(strNum);     // 42

// Character to number
let digitChar = '5';
let digitValue = parseInt(digitChar);  // 5
// or: digitChar.charCodeAt(0) - '0'.charCodeAt(0)
```

## Practical Use Cases

### Use Case 1: Text Processing and Validation
```go
// Email validation with Unicode support
func isValidEmail(email string) bool {
    // Convert to runes for proper Unicode handling
    runes := []rune(email)

    hasAt := false
    for _, r := range runes {
        if r == '@' {
            if hasAt {
                return false  // Multiple @
            }
            hasAt = true
        } else if !unicode.IsLetter(r) && !unicode.IsDigit(r) &&
                  r != '.' && r != '-' && r != '_' {
            return false  // Invalid character
        }
    }

    return hasAt && len(runes) > 3
}

// Text cleaning
func cleanText(input string) string {
    var result []rune

    for _, r := range input {
        if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
            result = append(result, unicode.ToLower(r))
        }
    }

    return string(result)
}
```

### Use Case 2: File and Network Operations
```go
// Reading different encodings
func readFileAsUTF8(filename string) (string, error) {
    bytes, err := os.ReadFile(filename)
    if err != nil {
        return "", err
    }

    // Validate UTF-8
    if !utf8.Valid(bytes) {
        return "", fmt.Errorf("file contains invalid UTF-8")
    }

    return string(bytes), nil
}

// HTTP response processing
func processHTTPResponse(resp *http.Response) (string, error) {
    defer resp.Body.Close()

    // Read as bytes first
    bodyBytes, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    // Check content encoding, convert if needed
    contentType := resp.Header.Get("Content-Type")
    if strings.Contains(contentType, "charset=utf-8") {
        return string(bodyBytes), nil
    }

    // Handle other encodings...
    return string(bodyBytes), nil
}
```

### Use Case 3: String Manipulation
```go
// Safe string truncation (character-aware)
func truncateText(text string, maxChars int) string {
    runes := []rune(text)
    if len(runes) <= maxChars {
        return text
    }

    return string(runes[:maxChars]) + "..."
}

// Word count (Unicode-aware)
func wordCount(text string) int {
    count := 0
    inWord := false

    for _, r := range text {
        if unicode.IsLetter(r) || unicode.IsDigit(r) {
            if !inWord {
                count++
                inWord = true
            }
        } else {
            inWord = false
        }
    }

    return count
}

// Case conversion (Unicode-aware)
func toTitleCase(text string) string {
    var result []rune
    capitalizeNext := true

    for _, r := range text {
        if unicode.IsLetter(r) {
            if capitalizeNext {
                result = append(result, unicode.ToUpper(r))
                capitalizeNext = false
            } else {
                result = append(result, unicode.ToLower(r))
            }
        } else {
            result = append(result, r)
            capitalizeNext = unicode.IsSpace(r)
        }
    }

    return string(result)
}
```

**JavaScript Equivalents:**
```javascript
// JavaScript text processing
function isValidEmail(email) {
    const hasAt = email.includes('@');
    const validChars = /^[a-zA-Z0-9@._-]+$/.test(email);
    return hasAt && validChars && email.length > 3;
}

function truncateText(text, maxChars) {
    if (text.length <= maxChars) return text;
    return text.slice(0, maxChars) + "...";
}

function wordCount(text) {
    return text.trim().split(/\s+/).filter(word => word.length > 0).length;
}

function toTitleCase(text) {
    return text.toLowerCase().replace(/\b\w/g, char => char.toUpperCase());
}
```

## Performance Considerations

### String Building Performance
```go
// ** Inefficient: String concatenation creates new strings each time
func buildStringBad(words []string) string {
    result := ""
    for _, word := range words {
        result += word + " "  // Creates new string each time!
    }
    return result
}

// ** Efficient: Use strings.Builder for multiple concatenations
func buildStringGood(words []string) string {
    var builder strings.Builder

    // Pre-allocate capacity if you know the approximate size
    builder.Grow(len(words) * 10)  // Estimate

    for _, word := range words {
        builder.WriteString(word)
        builder.WriteByte(' ')
    }

    return builder.String()
}

// ** Most efficient for simple cases: strings.Join
func buildStringBest(words []string) string {
    return strings.Join(words, " ")
}
```

### Iteration Performance
```go
// Benchmark different iteration methods
func benchmarkIterations(s string) {
    // Byte iteration: Fastest for ASCII
    start := time.Now()
    for i := 0; i < len(s); i++ {
        _ = s[i]
    }
    fmt.Println("Byte iteration:", time.Since(start))

    // Range iteration: Good balance
    start = time.Now()
    for _, r := range s {
        _ = r
    }
    fmt.Println("Range iteration:", time.Since(start))

    // Rune slice: Slowest but most flexible
    start = time.Now()
    runes := []rune(s)
    for i := 0; i < len(runes); i++ {
        _ = runes[i]
    }
    fmt.Println("Rune slice:", time.Since(start))
}
```

**JavaScript Performance:**
```javascript
// JavaScript string operations are generally optimized by the engine
function buildStringJS(words) {
    // Modern JavaScript engines optimize string concatenation
    return words.join(' ');

    // Template literals are also efficient
    // return words.map(word => `${word} `).join('');
}
```

## Best Practices

### When to Use Each Type

**Use `string` when:**
```go
// ** Storing and passing text data
func processUserName(name string) error {
    if name == "" {
        return errors.New("name cannot be empty")
    }
    return saveUser(name)
}

// ** String literals and constants
const welcomeMessage = "Welcome to our application!"

// ** Simple ASCII operations
func isASCIIAlphanumeric(s string) bool {
    for i := 0; i < len(s); i++ {
        b := s[i]
        if !((b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') || (b >= '0' && b <= '9')) {
            return false
        }
    }
    return true
}
```

**Use `[]rune` when:**
```go
// ** Character-level manipulation
func reverseText(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}

// ** Safe character counting and indexing
func getCharAt(s string, index int) (rune, bool) {
    runes := []rune(s)
    if index < 0 || index >= len(runes) {
        return 0, false
    }
    return runes[index], true
}

// ** Unicode-aware text processing
func removeAccents(s string) string {
    runes := []rune(s)
    var result []rune

    for _, r := range runes {
        normalized := norm.NFD.String(string(r))
        for _, nr := range normalized {
            if !unicode.Is(unicode.Mn, nr) {  // Skip combining marks
                result = append(result, nr)
            }
        }
    }

    return string(result)
}
```

**Use `[]byte` when:**
```go
// ** Binary data and file operations
func processImageFile(filename string) error {
    data, err := os.ReadFile(filename)  // []byte
    if err != nil {
        return err
    }

    // Check file signature
    if len(data) >= 4 && string(data[:4]) == "\x89PNG" {
        return processPNG(data)
    }

    return errors.New("unsupported format")
}

// ** Network protocols
func parseHTTPHeader(data []byte) map[string]string {
    headers := make(map[string]string)

    lines := bytes.Split(data, []byte("\r\n"))
    for _, line := range lines {
        parts := bytes.SplitN(line, []byte(":"), 2)
        if len(parts) == 2 {
            key := string(bytes.TrimSpace(parts[0]))
            value := string(bytes.TrimSpace(parts[1]))
            headers[key] = value
        }
    }

    return headers
}

// ** Performance-critical operations
func fastContains(haystack, needle []byte) bool {
    return bytes.Contains(haystack, needle)  // Optimized implementation
}
```

### Common Pitfalls and How to Avoid Them

** Don't slice strings arbitrarily:**
```go
// BAD: Can break Unicode characters
func badTruncate(s string, n int) string {
    if len(s) <= n {
        return s
    }
    return s[:n]  // Might cut through Unicode character!
}

// GOOD: Use rune-aware truncation
func goodTruncate(s string, n int) string {
    runes := []rune(s)
    if len(runes) <= n {
        return s
    }
    return string(runes[:n])
}
```

** Don't assume len(string) = character count:**
```go
// BAD: Assumes bytes = characters
func badCharCount(s string) int {
    return len(s)  // Wrong for Unicode!
}

// GOOD: Count actual characters
func goodCharCount(s string) int {
    return utf8.RuneCountInString(s)
    // or: return len([]rune(s))
}
```

** Don't build strings inefficiently:**
```go
// BAD: Creates many temporary strings
func badBuild(items []string) string {
    result := ""
    for _, item := range items {
        result = result + item + ", "  // Very inefficient!
    }
    return result
}

// GOOD: Use strings.Builder or strings.Join
func goodBuild(items []string) string {
    return strings.Join(items, ", ")
}
```

### Testing String Operations
```go
func TestStringOperations(t *testing.T) {
    testCases := []struct {
        input    string
        expected int
    }{
        {"Hello", 5},
        {"世界", 2},
        {"Hello, 世界! 🌍", 11},
        {"", 0},
    }

    for _, tc := range testCases {
        got := utf8.RuneCountInString(tc.input)
        if got != tc.expected {
            t.Errorf("RuneCountInString(%q) = %d, want %d",
                tc.input, got, tc.expected)
        }
    }
}
```

---

*Remember: JavaScript hides string complexity for ease of use, while Go exposes it for performance and control. Understanding the difference between bytes, runes, and strings is crucial for correct Unicode handling in Go!*