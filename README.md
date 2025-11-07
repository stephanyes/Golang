# Go Learning Code

This directory contains my hands-on Go programming exercises and projects as I learn the language coming from a Node.js background.

## Structure

### `/cmd/tutorial_*`
Progressive exercises covering Go fundamentals:
- `tutorial_1` - Basic variables, types, and syntax
- `tutorial_2` - Functions, control flow, and basic structures
- `tutorial_3` - More advanced concepts as I progress

Each tutorial builds on the previous one, with code that actually runs and demonstrates specific concepts.

### `/docs`
Comprehensive guides written for JavaScript developers learning Go. These serve as my reference material and document the mental model shifts required.

## Go Module
- **Module name**: `github.com/stephanyes/Golang`
- Run any tutorial with: `go run cmd/tutorial_X/main.go`

## Learning Notes

Key differences from Node.js development that I'm adapting to:
- Explicit error handling instead of try/catch
- Static typing with compile-time checks
- Memory management concepts (pointers, value vs reference)
- Concurrency through goroutines instead of async/await
- Composition over inheritance

## Compilation

Go compiles to native binaries. Build any tutorial:
```bash
go build cmd/tutorial_X/main.go
./main
```

This directory grows as I learn more Go concepts and build actual projects.