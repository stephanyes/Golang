# Structs and Interfaces in Go
*For JavaScript/TypeScript Developers*

## Table of Contents
1. [Object-Oriented Programming: Go vs JavaScript](#object-oriented-programming-go-vs-javascript)
2. [Structs - Go's Objects](#structs---gos-objects)
3. [Methods and Receivers](#methods-and-receivers)
4. [Interfaces - Implicit Contracts](#interfaces---implicit-contracts)
5. [Composition over Inheritance](#composition-over-inheritance)
6. [Empty Interface and Type Assertions](#empty-interface-and-type-assertions)
7. [Interface Segregation](#interface-segregation)
8. [Practical Design Patterns](#practical-design-patterns)
9. [Memory Layout and Performance](#memory-layout-and-performance)
10. [Best Practices](#best-practices)

## Object-Oriented Programming: Go vs JavaScript

### JavaScript's Approach
```javascript
// JavaScript: Classes and prototypes
class User {
    constructor(name, email) {
        this.name = name;
        this.email = email;
        this.createdAt = new Date();
    }

    // Methods
    greet() {
        return `Hello, I'm ${this.name}`;
    }

    // Inheritance
    static fromJSON(json) {
        const data = JSON.parse(json);
        return new User(data.name, data.email);
    }
}

// Inheritance
class AdminUser extends User {
    constructor(name, email, permissions) {
        super(name, email);
        this.permissions = permissions;
    }

    // Override method
    greet() {
        return `Hello, I'm ${this.name} (Admin)`;
    }

    // Additional method
    grantPermission(permission) {
        this.permissions.push(permission);
    }
}

// Usage
const user = new User("John", "john@example.com");
console.log(user.greet());  // "Hello, I'm John"

const admin = new AdminUser("Alice", "alice@example.com", ["read", "write"]);
console.log(admin.greet()); // "Hello, I'm Alice (Admin)"
```

### Go's Approach (No Classes!)
```go
// Go: Structs + Methods + Interfaces
import (
    "encoding/json"
    "fmt"
    "time"
)

// Struct definition (like a class without methods)
type User struct {
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}

// Methods are defined separately with receivers
func (u User) Greet() string {
    return fmt.Sprintf("Hello, I'm %s", u.Name)
}

// Constructor function (by convention)
func NewUser(name, email string) User {
    return User{
        Name:      name,
        Email:     email,
        CreatedAt: time.Now(),
    }
}

// "Static" method (package-level function)
func UserFromJSON(jsonData []byte) (User, error) {
    var user User
    err := json.Unmarshal(jsonData, &user)
    return user, err
}

// Composition instead of inheritance
type AdminUser struct {
    User                    // Embedded struct (composition)
    Permissions []string    `json:"permissions"`
}

// Methods for AdminUser
func (au AdminUser) Greet() string {
    return fmt.Sprintf("Hello, I'm %s (Admin)", au.Name)
}

func (au *AdminUser) GrantPermission(permission string) {
    au.Permissions = append(au.Permissions, permission)
}

// Usage
user := NewUser("John", "john@example.com")
fmt.Println(user.Greet())  // "Hello, I'm John"

admin := AdminUser{
    User:        NewUser("Alice", "alice@example.com"),
    Permissions: []string{"read", "write"},
}
fmt.Println(admin.Greet()) // "Hello, I'm Alice (Admin)"
admin.GrantPermission("delete")
```

**Key Difference:** Go uses composition and interfaces instead of classes and inheritance.

## Structs - Go's Objects

### Basic Struct Definition
```go
// Simple struct
type Point struct {
    X float64
    Y float64
}

// Struct with different field types
type Product struct {
    ID          int       `json:"id"`
    Name        string    `json:"name"`
    Price       float64   `json:"price"`
    InStock     bool      `json:"in_stock"`
    Categories  []string  `json:"categories"`
    CreatedAt   time.Time `json:"created_at"`
}

// Nested structs
type Address struct {
    Street   string `json:"street"`
    City     string `json:"city"`
    ZipCode  string `json:"zip_code"`
    Country  string `json:"country"`
}

type Customer struct {
    ID       int     `json:"id"`
    Name     string  `json:"name"`
    Email    string  `json:"email"`
    Address  Address `json:"address"`  // Nested struct
}
```

**JavaScript Equivalent:**
```javascript
// JavaScript objects are more dynamic
const point = { x: 10.5, y: 20.3 };

const product = {
    id: 1,
    name: "Laptop",
    price: 999.99,
    inStock: true,
    categories: ["electronics", "computers"],
    createdAt: new Date()
};

const customer = {
    id: 1,
    name: "John Doe",
    email: "john@example.com",
    address: {
        street: "123 Main St",
        city: "Anytown",
        zipCode: "12345",
        country: "USA"
    }
};
```

### Struct Initialization
```go
// Different ways to initialize structs
product := Product{
    ID:         1,
    Name:       "Laptop",
    Price:      999.99,
    InStock:    true,
    Categories: []string{"electronics", "computers"},
    CreatedAt:  time.Now(),
}

// Positional initialization (not recommended for readability)
point := Point{10.5, 20.3}

// Zero value initialization
var emptyProduct Product  // All fields get zero values
fmt.Printf("%+v\n", emptyProduct)  // {ID:0 Name: Price:0 InStock:false Categories:[] CreatedAt:0001-01-01...}

// Partial initialization
partialProduct := Product{
    Name:  "Mouse",
    Price: 29.99,
    // Other fields get zero values
}

// Using new (returns pointer)
productPtr := &Product{
    Name:  "Keyboard",
    Price: 79.99,
}
```

### Struct Tags and JSON Marshaling
```go
type User struct {
    ID        int       `json:"id" db:"user_id" validate:"required"`
    FirstName string    `json:"first_name" db:"first_name"`
    LastName  string    `json:"last_name" db:"last_name"`
    Email     string    `json:"email" db:"email" validate:"email"`
    Password  string    `json:"-" db:"password_hash"`  // "-" means omit from JSON
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`  // omit if zero value
}

// JSON marshaling
user := User{
    ID:        1,
    FirstName: "John",
    LastName:  "Doe",
    Email:     "john@example.com",
    Password:  "secret123",
    CreatedAt: time.Now(),
}

jsonData, _ := json.Marshal(user)
fmt.Println(string(jsonData))
// Output: {"id":1,"first_name":"John","last_name":"Doe","email":"john@example.com","created_at":"2023..."}
// Note: Password is omitted due to json:"-" tag
```

**JavaScript Equivalent:**
```javascript
// JavaScript doesn't have struct tags, but you can control JSON serialization
class User {
    constructor(id, firstName, lastName, email, password) {
        this.id = id;
        this.firstName = firstName;
        this.lastName = lastName;
        this.email = email;
        this.password = password;  // Would be included in JSON.stringify()
        this.createdAt = new Date();
    }

    // Custom JSON serialization
    toJSON() {
        return {
            id: this.id,
            first_name: this.firstName,
            last_name: this.lastName,
            email: this.email,
            created_at: this.createdAt
            // password intentionally omitted
        };
    }
}
```

## Methods and Receivers

### Value vs Pointer Receivers
```go
type Counter struct {
    Value int
}

// Value receiver - receives a copy
func (c Counter) GetValue() int {
    return c.Value  // Can read
}

func (c Counter) TryIncrement() {
    c.Value++  // ** This only modifies the copy!
}

// Pointer receiver - receives a pointer
func (c *Counter) Increment() {
    c.Value++  // ** This modifies the original
}

func (c *Counter) Reset() {
    c.Value = 0  // ** This modifies the original
}

// Usage
counter := Counter{Value: 5}

fmt.Println(counter.GetValue())  // 5
counter.TryIncrement()
fmt.Println(counter.GetValue())  // Still 5! (copy was modified)

counter.Increment()
fmt.Println(counter.GetValue())  // 6 (original was modified)
```

### When to Use Pointer vs Value Receivers
```go
type SmallStruct struct {
    ID   int
    Name string
}

type LargeStruct struct {
    Data [1000]int
    // ... many fields
}

// Use value receivers for:
// - Small structs (cheap to copy)
// - Immutable operations
// - When you don't need to modify the struct
func (s SmallStruct) String() string {
    return fmt.Sprintf("ID: %d, Name: %s", s.ID, s.Name)
}

// Use pointer receivers for:
// - Large structs (expensive to copy)
// - When you need to modify the struct
// - For consistency (if any method uses pointer, use pointer for all)
func (ls *LargeStruct) ProcessData() {
    for i := range ls.Data {
        ls.Data[i] = i * 2
    }
}

func (ls *LargeStruct) GetSum() int {  // Pointer even for read-only to avoid copying
    sum := 0
    for _, v := range ls.Data {
        sum += v
    }
    return sum
}
```

**JavaScript Comparison:**
```javascript
// JavaScript objects are always passed by reference
class Counter {
    constructor(value = 0) {
        this.value = value;
    }

    getValue() {
        return this.value;  // Always reads from original
    }

    increment() {
        this.value++;  // Always modifies original
    }
}

const counter = new Counter(5);
counter.increment();  // Always modifies the original object
```

### Method Chaining
```go
type StringBuilder struct {
    parts []string
}

// Methods that return *StringBuilder for chaining
func (sb *StringBuilder) Add(text string) *StringBuilder {
    sb.parts = append(sb.parts, text)
    return sb
}

func (sb *StringBuilder) AddLine(text string) *StringBuilder {
    sb.parts = append(sb.parts, text, "\n")
    return sb
}

func (sb *StringBuilder) String() string {
    return strings.Join(sb.parts, "")
}

// Usage with method chaining
result := &StringBuilder{}
output := result.
    Add("Hello ").
    Add("World").
    AddLine("!").
    Add("How are you?").
    String()

fmt.Println(output)  // "Hello World!\nHow are you?"
```

**JavaScript Equivalent:**
```javascript
class StringBuilder {
    constructor() {
        this.parts = [];
    }

    add(text) {
        this.parts.push(text);
        return this;  // Return 'this' for chaining
    }

    addLine(text) {
        this.parts.push(text, '\n');
        return this;
    }

    toString() {
        return this.parts.join('');
    }
}

// Method chaining works the same way
const output = new StringBuilder()
    .add("Hello ")
    .add("World")
    .addLine("!")
    .add("How are you?")
    .toString();
```

## Interfaces - Implicit Contracts

### Basic Interface Concepts
```go
// Interface definition
type Writer interface {
    Write([]byte) (int, error)
}

type Reader interface {
    Read([]byte) (int, error)
}

// Combining interfaces
type ReadWriter interface {
    Reader
    Writer
}

// Any type that implements these methods automatically satisfies the interface
type FileHandler struct {
    filename string
}

// Implementing Writer interface (implicitly)
func (fh FileHandler) Write(data []byte) (int, error) {
    // Implementation here
    return len(data), nil
}

// Implementing Reader interface (implicitly)
func (fh FileHandler) Read(data []byte) (int, error) {
    // Implementation here
    return len(data), nil
}

// FileHandler automatically satisfies ReadWriter interface!
// No explicit "implements" declaration needed

// Function that accepts any Writer
func WriteData(w Writer, data []byte) error {
    _, err := w.Write(data)
    return err
}

// Usage
fh := FileHandler{filename: "test.txt"}
WriteData(fh, []byte("Hello, World!"))  // FileHandler satisfies Writer
```

**JavaScript/TypeScript Comparison:**
```javascript
// JavaScript doesn't have built-in interfaces, but TypeScript does:

// TypeScript interface
interface Writer {
    write(data: Uint8Array): Promise<number>;
}

interface Reader {
    read(data: Uint8Array): Promise<number>;
}

interface ReadWriter extends Reader, Writer {}

class FileHandler implements ReadWriter {
    constructor(private filename: string) {}

    async write(data: Uint8Array): Promise<number> {
        // Implementation
        return data.length;
    }

    async read(data: Uint8Array): Promise<number> {
        // Implementation
        return data.length;
    }
}

// Function that accepts any Writer
async function writeData(writer: Writer, data: Uint8Array): Promise<void> {
    await writer.write(data);
}

// Usage
const fh = new FileHandler("test.txt");
await writeData(fh, new Uint8Array([1, 2, 3]));
```

### Common Built-in Interfaces
```go
import (
    "fmt"
    "io"
    "sort"
    "strings"
)

// fmt.Stringer interface
type Stringer interface {
    String() string
}

type Person struct {
    FirstName string
    LastName  string
    Age       int
}

// Implementing Stringer
func (p Person) String() string {
    return fmt.Sprintf("%s %s (%d years old)", p.FirstName, p.LastName, p.Age)
}

// Now Person can be printed directly
person := Person{"John", "Doe", 30}
fmt.Println(person)  // "John Doe (30 years old)"

// sort.Interface for custom sorting
type People []Person

func (p People) Len() int           { return len(p) }
func (p People) Less(i, j int) bool { return p[i].Age < p[j].Age }
func (p People) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Now People can be sorted
people := People{
    {"Alice", "Smith", 25},
    {"Bob", "Jones", 30},
    {"Charlie", "Brown", 20},
}

sort.Sort(people)  // Sorts by age
fmt.Println(people)  // Sorted by age: Charlie (20), Alice (25), Bob (30)
```

### Interface Design Patterns
```go
// Small, focused interfaces (Interface Segregation Principle)
type Validator interface {
    Validate() error
}

type Saver interface {
    Save() error
}

type Loader interface {
    Load() error
}

// Compose larger interfaces from smaller ones
type Persistable interface {
    Saver
    Loader
    Validator
}

// Real-world example: HTTP handlers
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}

// Custom handler
type APIHandler struct {
    database DB
}

func (h APIHandler) ServeHTTP(w ResponseWriter, r *Request) {
    // Handle HTTP request
    switch r.Method {
    case "GET":
        h.handleGet(w, r)
    case "POST":
        h.handlePost(w, r)
    default:
        w.WriteHeader(405)
    }
}

// Function that works with any handler
func RegisterHandler(pattern string, handler Handler) {
    http.Handle(pattern, handler)
}
```

## Composition over Inheritance

### Embedding Structs
```go
// Base functionality
type Logger struct {
    prefix string
}

func (l Logger) Log(message string) {
    fmt.Printf("[%s] %s\n", l.prefix, message)
}

// Service with embedded logger
type UserService struct {
    Logger  // Embedded struct - promotes Logger's methods
    db      Database
}

func NewUserService(db Database) UserService {
    return UserService{
        Logger: Logger{prefix: "UserService"},
        db:     db,
    }
}

func (us UserService) CreateUser(name, email string) error {
    us.Log("Creating user: " + name)  // Can call embedded method directly

    user := User{Name: name, Email: email}
    if err := us.db.Save(user); err != nil {
        us.Log("Failed to create user: " + err.Error())
        return err
    }

    us.Log("User created successfully")
    return nil
}

// Multiple embeddings
type AuditedUserService struct {
    UserService  // Embedded
    Auditor     // Another embedded struct
}

func (aus AuditedUserService) CreateUser(name, email string) error {
    // Call the original method
    err := aus.UserService.CreateUser(name, email)

    // Add auditing
    aus.Auditor.LogAction("CreateUser", name)

    return err
}
```

**JavaScript Equivalent (Composition):**
```javascript
// JavaScript composition pattern
class Logger {
    constructor(prefix) {
        this.prefix = prefix;
    }

    log(message) {
        console.log(`[${this.prefix}] ${message}`);
    }
}

class UserService {
    constructor(database) {
        this.logger = new Logger("UserService");  // Composition
        this.db = database;
    }

    createUser(name, email) {
        this.logger.log(`Creating user: ${name}`);

        const user = { name, email };
        try {
            this.db.save(user);
            this.logger.log("User created successfully");
        } catch (error) {
            this.logger.log(`Failed to create user: ${error.message}`);
            throw error;
        }
    }
}

// Extending with more composition
class AuditedUserService {
    constructor(database, auditor) {
        this.userService = new UserService(database);  // Composition
        this.auditor = auditor;
    }

    createUser(name, email) {
        const result = this.userService.createUser(name, email);
        this.auditor.logAction("CreateUser", name);
        return result;
    }
}
```

### Interface Composition
```go
// Small, focused interfaces
type Reader interface {
    Read([]byte) (int, error)
}

type Writer interface {
    Write([]byte) (int, error)
}

type Closer interface {
    Close() error
}

// Composed interfaces
type ReadCloser interface {
    Reader
    Closer
}

type WriteCloser interface {
    Writer
    Closer
}

type ReadWriteCloser interface {
    Reader
    Writer
    Closer
}

// File implements all these interfaces
type File struct {
    name string
    data []byte
    pos  int
}

func (f *File) Read(p []byte) (int, error) {
    // Implementation
    return 0, nil
}

func (f *File) Write(p []byte) (int, error) {
    // Implementation
    return len(p), nil
}

func (f *File) Close() error {
    // Implementation
    return nil
}

// Functions can accept specific interface combinations
func CopyAndClose(dst WriteCloser, src ReadCloser) error {
    defer src.Close()
    defer dst.Close()

    _, err := io.Copy(dst, src)
    return err
}

// File can be used in any context requiring these interfaces
file1 := &File{name: "source.txt"}
file2 := &File{name: "dest.txt"}
CopyAndClose(file2, file1)  // Both satisfy the required interfaces
```

## Empty Interface and Type Assertions

### The Empty Interface `interface{}`
```go
// interface{} can hold any value (like 'any' in TypeScript)
func printAnything(value interface{}) {
    fmt.Printf("Value: %v, Type: %T\n", value, value)
}

// Usage
printAnything(42)          // Value: 42, Type: int
printAnything("hello")     // Value: hello, Type: string
printAnything([]int{1,2,3}) // Value: [1 2 3], Type: []int

// Storing different types in a slice
var mixed []interface{}
mixed = append(mixed, 42)
mixed = append(mixed, "hello")
mixed = append(mixed, true)

fmt.Println(mixed)  // [42 hello true]
```

### Type Assertions
```go
// Type assertion syntax: value.(Type)
func processValue(value interface{}) {
    // Safe type assertion with ok check
    if str, ok := value.(string); ok {
        fmt.Printf("String value: %s (length: %d)\n", str, len(str))
        return
    }

    if num, ok := value.(int); ok {
        fmt.Printf("Integer value: %d (double: %d)\n", num, num*2)
        return
    }

    fmt.Printf("Unknown type: %T\n", value)
}

// Type switch for multiple types
func handleValue(value interface{}) {
    switch v := value.(type) {
    case string:
        fmt.Printf("String: %s\n", v)
    case int:
        fmt.Printf("Integer: %d\n", v)
    case bool:
        fmt.Printf("Boolean: %t\n", v)
    case []int:
        fmt.Printf("Integer slice: %v\n", v)
    default:
        fmt.Printf("Unknown type: %T\n", v)
    }
}
```

**JavaScript Equivalent:**
```javascript
// JavaScript is dynamically typed, so this is natural
function processValue(value) {
    if (typeof value === 'string') {
        console.log(`String value: ${value} (length: ${value.length})`);
        return;
    }

    if (typeof value === 'number') {
        console.log(`Number value: ${value} (double: ${value * 2})`);
        return;
    }

    console.log(`Unknown type: ${typeof value}`);
}

// TypeScript equivalent with union types
function handleValue(value: string | number | boolean | number[]): void {
    switch (typeof value) {
        case 'string':
            console.log(`String: ${value}`);
            break;
        case 'number':
            console.log(`Number: ${value}`);
            break;
        case 'boolean':
            console.log(`Boolean: ${value}`);
            break;
        default:
            if (Array.isArray(value)) {
                console.log(`Number array: ${value}`);
            }
    }
}
```

### Modern Go: Using `any` (Go 1.18+)
```go
// Go 1.18+ introduced 'any' as an alias for interface{}
type any = interface{}

// More readable than interface{}
func processAny(value any) {
    switch v := value.(type) {
    case string:
        fmt.Println("String:", v)
    case int:
        fmt.Println("Int:", v)
    default:
        fmt.Printf("Other: %T\n", v)
    }
}
```

## Interface Segregation

### Bad Interface Design (Too Big)
```go
// ** BAD: Interface does too many things
type UserManager interface {
    CreateUser(User) error
    UpdateUser(User) error
    DeleteUser(int) error
    GetUser(int) (User, error)
    ListUsers() ([]User, error)
    AuthenticateUser(string, string) (User, error)
    SendWelcomeEmail(User) error
    GenerateReport() (Report, error)
    BackupUsers() error
}

// Any implementation must implement ALL methods, even if they only need some
type SimpleUserStorage struct{}

// Must implement all methods, even if not needed
func (s SimpleUserStorage) CreateUser(User) error         { return nil }
func (s SimpleUserStorage) UpdateUser(User) error         { return nil }
func (s SimpleUserStorage) DeleteUser(int) error          { return nil }
func (s SimpleUserStorage) GetUser(int) (User, error)     { return User{}, nil }
func (s SimpleUserStorage) ListUsers() ([]User, error)    { return nil, nil }
func (s SimpleUserStorage) AuthenticateUser(string, string) (User, error) { return User{}, nil }
func (s SimpleUserStorage) SendWelcomeEmail(User) error   { return nil }  // ** Doesn't belong here!
func (s SimpleUserStorage) GenerateReport() (Report, error) { return Report{}, nil }  // ** Wrong responsibility!
func (s SimpleUserStorage) BackupUsers() error            { return nil }  // ** Different concern!
```

### Good Interface Design (Small and Focused)
```go
// ** GOOD: Small, focused interfaces
type UserReader interface {
    GetUser(int) (User, error)
    ListUsers() ([]User, error)
}

type UserWriter interface {
    CreateUser(User) error
    UpdateUser(User) error
    DeleteUser(int) error
}

type UserAuthenticator interface {
    AuthenticateUser(username, password string) (User, error)
}

type EmailSender interface {
    SendWelcomeEmail(User) error
}

type ReportGenerator interface {
    GenerateReport() (Report, error)
}

type BackupService interface {
    BackupUsers() error
}

// Compose interfaces when needed
type UserRepository interface {
    UserReader
    UserWriter
}

// Implementations only need to implement what they actually do
type DatabaseUserStorage struct {
    db Database
}

// Only implements storage-related methods
func (dus DatabaseUserStorage) GetUser(id int) (User, error)     { /* db implementation */ return User{}, nil }
func (dus DatabaseUserStorage) ListUsers() ([]User, error)       { /* db implementation */ return nil, nil }
func (dus DatabaseUserStorage) CreateUser(user User) error       { /* db implementation */ return nil }
func (dus DatabaseUserStorage) UpdateUser(user User) error       { /* db implementation */ return nil }
func (dus DatabaseUserStorage) DeleteUser(id int) error          { /* db implementation */ return nil }

type SMTPEmailService struct {
    host string
    port int
}

// Only implements email-related methods
func (ses SMTPEmailService) SendWelcomeEmail(user User) error { /* email implementation */ return nil }

type FileBackupService struct {
    backupPath string
}

// Only implements backup-related methods
func (fbs FileBackupService) BackupUsers() error { /* backup implementation */ return nil }

// Service that composes different implementations
type UserService struct {
    repository UserRepository
    emailer    EmailSender
    backup     BackupService
}

func (us UserService) RegisterUser(user User) error {
    // Create user
    if err := us.repository.CreateUser(user); err != nil {
        return err
    }

    // Send welcome email
    if err := us.emailer.SendWelcomeEmail(user); err != nil {
        // Log error but don't fail registration
        log.Printf("Failed to send welcome email: %v", err)
    }

    return nil
}
```

## Practical Design Patterns

### Repository Pattern
```go
// Domain model
type User struct {
    ID       int
    Username string
    Email    string
    CreatedAt time.Time
}

// Repository interface
type UserRepository interface {
    Save(user *User) error
    FindByID(id int) (*User, error)
    FindByEmail(email string) (*User, error)
    Delete(id int) error
}

// Database implementation
type SQLUserRepository struct {
    db *sql.DB
}

func (r *SQLUserRepository) Save(user *User) error {
    query := `INSERT INTO users (username, email, created_at) VALUES (?, ?, ?)`
    result, err := r.db.Exec(query, user.Username, user.Email, user.CreatedAt)
    if err != nil {
        return err
    }

    id, err := result.LastInsertId()
    if err != nil {
        return err
    }

    user.ID = int(id)
    return nil
}

func (r *SQLUserRepository) FindByID(id int) (*User, error) {
    user := &User{}
    query := `SELECT id, username, email, created_at FROM users WHERE id = ?`
    err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
    if err != nil {
        return nil, err
    }
    return user, nil
}

// In-memory implementation (for testing)
type MemoryUserRepository struct {
    users  map[int]*User
    nextID int
}

func NewMemoryUserRepository() *MemoryUserRepository {
    return &MemoryUserRepository{
        users:  make(map[int]*User),
        nextID: 1,
    }
}

func (r *MemoryUserRepository) Save(user *User) error {
    if user.ID == 0 {
        user.ID = r.nextID
        r.nextID++
    }
    r.users[user.ID] = user
    return nil
}

func (r *MemoryUserRepository) FindByID(id int) (*User, error) {
    user, exists := r.users[id]
    if !exists {
        return nil, errors.New("user not found")
    }
    return user, nil
}

// Service using repository
type UserService struct {
    repo UserRepository  // Interface, not concrete type
}

func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) CreateUser(username, email string) (*User, error) {
    user := &User{
        Username:  username,
        Email:     email,
        CreatedAt: time.Now(),
    }

    return user, s.repo.Save(user)
}

// Usage - can swap implementations easily
func main() {
    // Production: use database
    db, _ := sql.Open("mysql", "connection-string")
    userService := NewUserService(&SQLUserRepository{db: db})

    // Testing: use in-memory
    // userService := NewUserService(NewMemoryUserRepository())

    user, err := userService.CreateUser("john", "john@example.com")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Created user: %+v\n", user)
}
```

### Strategy Pattern
```go
// Strategy interface
type PaymentProcessor interface {
    ProcessPayment(amount float64) error
}

// Concrete strategies
type CreditCardProcessor struct {
    cardNumber string
}

func (ccp CreditCardProcessor) ProcessPayment(amount float64) error {
    fmt.Printf("Processing $%.2f via credit card ending in %s\n",
        amount, ccp.cardNumber[len(ccp.cardNumber)-4:])
    return nil
}

type PayPalProcessor struct {
    email string
}

func (pp PayPalProcessor) ProcessPayment(amount float64) error {
    fmt.Printf("Processing $%.2f via PayPal account %s\n", amount, pp.email)
    return nil
}

type BankTransferProcessor struct {
    accountNumber string
}

func (btp BankTransferProcessor) ProcessPayment(amount float64) error {
    fmt.Printf("Processing $%.2f via bank transfer to account %s\n", amount, btp.accountNumber)
    return nil
}

// Context that uses strategy
type OrderProcessor struct {
    paymentProcessor PaymentProcessor
}

func (op *OrderProcessor) SetPaymentProcessor(processor PaymentProcessor) {
    op.paymentProcessor = processor
}

func (op *OrderProcessor) ProcessOrder(orderID string, amount float64) error {
    fmt.Printf("Processing order %s for $%.2f\n", orderID, amount)
    return op.paymentProcessor.ProcessPayment(amount)
}

// Usage
processor := &OrderProcessor{}

// Switch between different payment methods
processor.SetPaymentProcessor(CreditCardProcessor{cardNumber: "1234567890123456"})
processor.ProcessOrder("ORD-001", 99.99)

processor.SetPaymentProcessor(PayPalProcessor{email: "user@example.com"})
processor.ProcessOrder("ORD-002", 149.99)
```

## Memory Layout and Performance

### Struct Memory Layout
```go
// Memory layout matters for performance
type BadStruct struct {
    A bool    // 1 byte
    B int64   // 8 bytes (but needs 8-byte alignment, so 7 bytes padding)
    C bool    // 1 byte
    D int64   // 8 bytes (but needs 8-byte alignment, so 7 bytes padding)
}
// Total: 1 + 7 + 8 + 1 + 7 + 8 = 32 bytes

type GoodStruct struct {
    B int64   // 8 bytes
    D int64   // 8 bytes
    A bool    // 1 byte
    C bool    // 1 byte (+ 6 bytes padding at end for struct alignment)
}
// Total: 8 + 8 + 1 + 1 + 6 = 24 bytes (25% smaller!)

// Check struct size
fmt.Printf("BadStruct size: %d bytes\n", unsafe.Sizeof(BadStruct{}))   // 32
fmt.Printf("GoodStruct size: %d bytes\n", unsafe.Sizeof(GoodStruct{})) // 24
```

### Interface Performance
```go
// Interface values have overhead
type Animal interface {
    Speak() string
}

type Dog struct {
    Name string
}

func (d Dog) Speak() string {
    return "Woof!"
}

// Direct call - faster
dog := Dog{Name: "Buddy"}
sound := dog.Speak()  // Direct method call

// Interface call - slower (but more flexible)
var animal Animal = dog  // Interface conversion
sound2 := animal.Speak()  // Dynamic dispatch

// Benchmark shows interface calls are ~2-3x slower
// But the flexibility is often worth it
```

### Value vs Pointer Performance
```go
type SmallStruct struct {
    X, Y int
}

type LargeStruct struct {
    Data [1000]int
}

// Small structs: value receivers often better (no pointer indirection)
func (s SmallStruct) Process() int {
    return s.X + s.Y  // Fast, direct access
}

// Large structs: pointer receivers almost always better
func (s *LargeStruct) Process() int {  // Avoids copying 8KB
    sum := 0
    for _, v := range s.Data {
        sum += v
    }
    return sum
}
```

## Best Practices

### Interface Design Guidelines

**1. Keep interfaces small:**
```go
// ** Good: Small, focused interface
type Writer interface {
    Write([]byte) (int, error)
}

// ** Bad: Too many responsibilities
type FileManager interface {
    Write([]byte) (int, error)
    Read([]byte) (int, error)
    Backup() error
    Compress() error
    Encrypt() error
}
```

**2. Accept interfaces, return structs:**
```go
// ** Good: Accept interface for flexibility
func ProcessData(w Writer, data []byte) error {
    _, err := w.Write(data)
    return err
}

// ** Good: Return concrete type for clarity
func NewFileWriter(filename string) *FileWriter {
    return &FileWriter{filename: filename}
}

// ** Bad: Return interface unnecessarily
func NewWriter(filename string) Writer {  // Limits flexibility
    return &FileWriter{filename: filename}
}
```

**3. Use composition over embedding when relationships aren't "is-a":**
```go
// ** Good: User HAS a logger, not IS a logger
type UserService struct {
    logger Logger  // Composition
    db     Database
}

// ** Bad: User IS a logger? (doesn't make sense)
type UserService struct {
    Logger  // Embedding suggests "is-a" relationship
    db Database
}
```

### Struct Design Guidelines

**1. Group related fields:**
```go
// ** Good: Logical grouping
type User struct {
    // Identity
    ID       int
    Username string
    Email    string

    // Metadata
    CreatedAt time.Time
    UpdatedAt time.Time
    LastLogin time.Time

    // Preferences
    Theme    string
    Language string
}
```

**2. Use tags appropriately:**
```go
type User struct {
    ID       int    `json:"id" db:"user_id" validate:"required"`
    Username string `json:"username" db:"username" validate:"required,min=3,max=20"`
    Email    string `json:"email" db:"email" validate:"required,email"`
    Password string `json:"-" db:"password_hash"`  // Never serialize password
}
```

**3. Provide constructors for complex structs:**
```go
type DatabaseConfig struct {
    Host     string
    Port     int
    Database string
    Username string
    Password string
    MaxConns int
    Timeout  time.Duration
}

// Constructor with sensible defaults
func NewDatabaseConfig(host, database, username, password string) *DatabaseConfig {
    return &DatabaseConfig{
        Host:     host,
        Port:     5432,  // Default PostgreSQL port
        Database: database,
        Username: username,
        Password: password,
        MaxConns: 10,    // Reasonable default
        Timeout:  30 * time.Second,
    }
}

// Fluent API for optional configuration
func (dc *DatabaseConfig) WithPort(port int) *DatabaseConfig {
    dc.Port = port
    return dc
}

func (dc *DatabaseConfig) WithMaxConnections(maxConns int) *DatabaseConfig {
    dc.MaxConns = maxConns
    return dc
}

// Usage
config := NewDatabaseConfig("localhost", "myapp", "user", "pass").
    WithPort(5433).
    WithMaxConnections(20)
```

### Testing with Interfaces
```go
// Interface makes testing easy
type EmailSender interface {
    SendEmail(to, subject, body string) error
}

type UserService struct {
    emailSender EmailSender
}

func (us *UserService) RegisterUser(user User) error {
    // ... registration logic ...

    return us.emailSender.SendEmail(user.Email, "Welcome!", "Welcome to our service!")
}

// Mock for testing
type MockEmailSender struct {
    SentEmails []Email
}

type Email struct {
    To      string
    Subject string
    Body    string
}

func (m *MockEmailSender) SendEmail(to, subject, body string) error {
    m.SentEmails = append(m.SentEmails, Email{
        To:      to,
        Subject: subject,
        Body:    body,
    })
    return nil
}

// Test
func TestUserRegistration(t *testing.T) {
    mockEmailer := &MockEmailSender{}
    userService := &UserService{emailSender: mockEmailer}

    user := User{Email: "test@example.com"}
    err := userService.RegisterUser(user)

    assert.NoError(t, err)
    assert.Len(t, mockEmailer.SentEmails, 1)
    assert.Equal(t, "test@example.com", mockEmailer.SentEmails[0].To)
    assert.Equal(t, "Welcome!", mockEmailer.SentEmails[0].Subject)
}
```

---

*Remember: Go's approach to OOP emphasizes composition, interfaces, and explicit design over inheritance and complex class hierarchies. This leads to more maintainable and testable code, especially for large applications and APIs!*