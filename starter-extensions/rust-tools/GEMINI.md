# GEMINI.md - Rust Development Guide

This comprehensive guide provides patterns, best practices, and conventions for Rust development. It is designed to help you write safe, performant, and idiomatic Rust code.

## Table of Contents
1. [Cargo Best Practices](#cargo-best-practices)
2. [Error Handling](#error-handling)
3. [Memory Management](#memory-management)
4. [Async Programming](#async-programming)
5. [Testing Strategy](#testing-strategy)
6. [Performance Guidelines](#performance-guidelines)
7. [Security Considerations](#security-considerations)
8. [Common Patterns](#common-patterns)
9. [Development Workflow](#development-workflow)
10. [Debugging Tips](#debugging-tips)

## Cargo Best Practices

### Essential Cargo Commands

```bash
# Create a new project
cargo new my_project --bin  # Binary project
cargo new my_lib --lib      # Library project

# Build with optimizations
cargo build --release

# Run with optimizations
cargo run --release

# Check code without building
cargo check

# Run tests
cargo test
cargo test -- --nocapture  # Show println! output

# Generate and open documentation
cargo doc --open

# Update dependencies
cargo update

# Audit dependencies for security vulnerabilities
cargo audit
```

### Cargo.toml Configuration

```toml
[package]
name = "my_project"
version = "0.1.0"
edition = "2021"  # Always use the latest edition

[dependencies]
# Specify exact versions for reproducible builds
serde = { version = "1.0.193", features = ["derive"] }
tokio = { version = "1.35.0", features = ["full"] }

[dev-dependencies]
# Development-only dependencies
pretty_assertions = "1.4.0"
criterion = "0.5.1"

[profile.release]
# Optimize for size
opt-level = "z"
lto = true
codegen-units = 1
strip = true

[profile.dev]
# Faster compilation for development
opt-level = 0
debug = true
```

### Make Warnings Errors

**Always treat warnings as errors in CI/production:**

```rust
// In lib.rs or main.rs
#![deny(warnings)]
#![deny(clippy::all)]
#![deny(clippy::pedantic)]
#![deny(clippy::nursery)]
#![deny(missing_docs)]
#![deny(missing_debug_implementations)]
```

Or in Cargo.toml:

```toml
[lints.rust]
warnings = "deny"
unsafe_code = "deny"
missing_docs = "deny"

[lints.clippy]
all = "deny"
pedantic = "deny"
nursery = "deny"
```

### Dependency Management

```bash
# Add dependencies properly
cargo add serde --features derive
cargo add tokio --features full
cargo add --dev pretty_assertions

# Check for outdated dependencies
cargo install cargo-outdated
cargo outdated

# Check for security vulnerabilities
cargo install cargo-audit
cargo audit

# Minimize dependencies
cargo tree  # View dependency tree
cargo tree --duplicates  # Find duplicate dependencies
```

## Error Handling

### Use Result<T, E> Everywhere

```rust
use thiserror::Error;

#[derive(Error, Debug)]
pub enum AppError {
    #[error("IO error: {0}")]
    Io(#[from] std::io::Error),
    
    #[error("Parse error: {0}")]
    Parse(#[from] std::num::ParseIntError),
    
    #[error("Custom error: {message}")]
    Custom { message: String },
}

// Always return Result
pub fn risky_operation() -> Result<String, AppError> {
    let data = std::fs::read_to_string("config.toml")?;
    let value: i32 = data.parse()?;
    
    if value < 0 {
        return Err(AppError::Custom {
            message: "Value must be positive".to_string(),
        });
    }
    
    Ok(format!("Success: {}", value))
}
```

### Error Context with anyhow/eyre

```rust
use color_eyre::Result;
use color_eyre::eyre::{eyre, WrapErr};

fn load_config(path: &str) -> Result<Config> {
    let content = std::fs::read_to_string(path)
        .wrap_err_with(|| format!("Failed to read config file: {}", path))?;
    
    let config: Config = toml::from_str(&content)
        .wrap_err("Failed to parse config file")?;
    
    if config.version < MIN_VERSION {
        return Err(eyre!("Config version {} is too old", config.version));
    }
    
    Ok(config)
}
```

### Custom Error Types

```rust
// Define domain-specific errors
#[derive(Debug, thiserror::Error)]
pub enum ValidationError {
    #[error("Value {0} is out of range [{1}, {2}]")]
    OutOfRange(i32, i32, i32),
    
    #[error("Invalid format: expected {expected}, got {actual}")]
    InvalidFormat { expected: String, actual: String },
}

// Implement conversions
impl From<ValidationError> for AppError {
    fn from(err: ValidationError) -> Self {
        AppError::Validation(err)
    }
}
```

## Memory Management

### Ownership Rules

```rust
// 1. Each value has a single owner
let s1 = String::from("hello");
let s2 = s1;  // s1 is moved, no longer valid
// println!("{}", s1);  // ERROR: use of moved value

// 2. References must not outlive their data
fn invalid_reference() -> &String {  // ERROR: missing lifetime
    let s = String::from("hello");
    &s  // s is dropped here
}

// 3. Cannot have mutable and immutable references simultaneously
let mut s = String::from("hello");
let r1 = &s;  // Immutable borrow
let r2 = &s;  // Another immutable borrow - OK
// let r3 = &mut s;  // ERROR: cannot borrow as mutable
```

### Smart Pointers

```rust
use std::rc::Rc;
use std::cell::RefCell;
use std::sync::Arc;
use std::sync::Mutex;

// Single-threaded reference counting
let rc = Rc::new(vec![1, 2, 3]);
let rc_clone = Rc::clone(&rc);

// Interior mutability for single-threaded code
let cell = RefCell::new(5);
*cell.borrow_mut() += 1;

// Thread-safe reference counting
let arc = Arc::new(Mutex::new(0));
let arc_clone = Arc::clone(&arc);

// Spawn thread with shared data
std::thread::spawn(move || {
    let mut data = arc_clone.lock().unwrap();
    *data += 1;
});
```

### Avoiding Memory Leaks

```rust
// Use weak references to break cycles
use std::rc::{Rc, Weak};

struct Node {
    value: i32,
    parent: RefCell<Weak<Node>>,
    children: RefCell<Vec<Rc<Node>>>,
}

// RAII pattern - automatic cleanup
struct TempFile {
    path: PathBuf,
}

impl Drop for TempFile {
    fn drop(&mut self) {
        let _ = std::fs::remove_file(&self.path);
    }
}
```

## Async Programming

### Tokio Best Practices

```rust
use tokio::time::{sleep, Duration};
use tokio::sync::mpsc;

#[tokio::main]
async fn main() -> Result<()> {
    // Spawn concurrent tasks
    let handle1 = tokio::spawn(async {
        sleep(Duration::from_secs(1)).await;
        "Task 1 complete"
    });
    
    let handle2 = tokio::spawn(async {
        sleep(Duration::from_secs(2)).await;
        "Task 2 complete"
    });
    
    // Wait for both tasks
    let (result1, result2) = tokio::join!(handle1, handle2);
    println!("{}, {}", result1?, result2?);
    
    Ok(())
}

// Channels for communication
async fn channel_example() {
    let (tx, mut rx) = mpsc::channel(100);
    
    tokio::spawn(async move {
        for i in 0..10 {
            tx.send(i).await.unwrap();
        }
    });
    
    while let Some(value) = rx.recv().await {
        println!("Received: {}", value);
    }
}
```

### Async Error Handling

```rust
use tokio::time::timeout;

async fn with_timeout() -> Result<String> {
    match timeout(Duration::from_secs(5), slow_operation()).await {
        Ok(result) => result,
        Err(_) => Err(eyre!("Operation timed out")),
    }
}

// Retry logic
async fn with_retry<F, T>(mut f: F) -> Result<T>
where
    F: FnMut() -> futures::future::BoxFuture<'static, Result<T>>,
{
    let mut retries = 3;
    loop {
        match f().await {
            Ok(value) => return Ok(value),
            Err(e) if retries > 0 => {
                retries -= 1;
                sleep(Duration::from_secs(1)).await;
                continue;
            }
            Err(e) => return Err(e),
        }
    }
}
```

## Testing Strategy

### Unit Tests

```rust
#[cfg(test)]
mod tests {
    use super::*;
    use pretty_assertions::assert_eq;

    #[test]
    fn test_basic_functionality() {
        let result = calculate(2, 3);
        assert_eq!(result, 5);
    }

    #[test]
    #[should_panic(expected = "divide by zero")]
    fn test_panic() {
        divide(10, 0);
    }

    #[test]
    fn test_result() -> Result<()> {
        let value = risky_operation()?;
        assert!(value > 0);
        Ok(())
    }
}
```

### Integration Tests

```rust
// tests/integration_test.rs
use my_crate::Client;

#[tokio::test]
async fn test_client_connection() {
    let client = Client::new("localhost:8080");
    let response = client.get("/health").await.unwrap();
    assert_eq!(response.status(), 200);
}

// Test with fixtures
#[test]
fn test_with_temp_dir() {
    let temp_dir = tempfile::tempdir().unwrap();
    let file_path = temp_dir.path().join("test.txt");
    
    std::fs::write(&file_path, "test content").unwrap();
    
    // Test code that uses the file
    let content = std::fs::read_to_string(&file_path).unwrap();
    assert_eq!(content, "test content");
    
    // temp_dir is automatically cleaned up
}
```

### Property-Based Testing

```rust
use proptest::prelude::*;

proptest! {
    #[test]
    fn test_parse_and_format(s in "[0-9]{1,10}") {
        let n: u64 = s.parse().unwrap();
        let formatted = format!("{}", n);
        prop_assert_eq!(s, formatted);
    }
}
```

## Performance Guidelines

### Benchmarking

```rust
use criterion::{black_box, criterion_group, criterion_main, Criterion};

fn fibonacci(n: u64) -> u64 {
    match n {
        0 => 1,
        1 => 1,
        n => fibonacci(n-1) + fibonacci(n-2),
    }
}

fn criterion_benchmark(c: &mut Criterion) {
    c.bench_function("fib 20", |b| {
        b.iter(|| fibonacci(black_box(20)))
    });
}

criterion_group!(benches, criterion_benchmark);
criterion_main!(benches);
```

### Performance Tips

```rust
// 1. Use iterators instead of indexing
let sum: i32 = vec.iter().sum();  // Good
let mut sum = 0;
for i in 0..vec.len() {  // Less efficient
    sum += vec[i];
}

// 2. Avoid unnecessary allocations
fn process(data: &[u8]) -> Vec<u8> {
    let mut result = Vec::with_capacity(data.len());  // Pre-allocate
    for &byte in data {
        if byte != 0 {
            result.push(byte);
        }
    }
    result
}

// 3. Use const generics for compile-time optimization
fn dot_product<const N: usize>(a: &[f64; N], b: &[f64; N]) -> f64 {
    a.iter().zip(b.iter()).map(|(x, y)| x * y).sum()
}

// 4. Enable link-time optimization
// In Cargo.toml:
// [profile.release]
// lto = true
```

## Security Considerations

### Input Validation

```rust
use validator::Validate;

#[derive(Debug, Validate)]
struct UserInput {
    #[validate(length(min = 1, max = 100))]
    username: String,
    
    #[validate(email)]
    email: String,
    
    #[validate(range(min = 18, max = 150))]
    age: u8,
}

fn process_input(input: UserInput) -> Result<()> {
    input.validate()?;
    // Process validated input
    Ok(())
}
```

### Secure Random

```rust
use rand::Rng;
use rand::distributions::Alphanumeric;

fn generate_token() -> String {
    rand::thread_rng()
        .sample_iter(&Alphanumeric)
        .take(32)
        .map(char::from)
        .collect()
}
```

### Avoiding Common Vulnerabilities

```rust
// SQL injection prevention
use sqlx::query;

async fn get_user(id: i32) -> Result<User> {
    // Use parameterized queries
    let user = query_as!(User, "SELECT * FROM users WHERE id = $1", id)
        .fetch_one(&pool)
        .await?;
    Ok(user)
}

// Path traversal prevention
use std::path::{Path, PathBuf};

fn safe_path_join(base: &Path, user_input: &str) -> Result<PathBuf> {
    let path = base.join(user_input);
    
    // Ensure the path is within the base directory
    if !path.starts_with(base) {
        return Err(eyre!("Path traversal attempt detected"));
    }
    
    Ok(path)
}
```

## Common Patterns

### Builder Pattern

```rust
#[derive(Default)]
struct ServerBuilder {
    host: Option<String>,
    port: Option<u16>,
    threads: Option<usize>,
}

impl ServerBuilder {
    fn new() -> Self {
        Self::default()
    }
    
    fn host(mut self, host: impl Into<String>) -> Self {
        self.host = Some(host.into());
        self
    }
    
    fn port(mut self, port: u16) -> Self {
        self.port = Some(port);
        self
    }
    
    fn threads(mut self, threads: usize) -> Self {
        self.threads = Some(threads);
        self
    }
    
    fn build(self) -> Result<Server> {
        Ok(Server {
            host: self.host.unwrap_or_else(|| "localhost".to_string()),
            port: self.port.unwrap_or(8080),
            threads: self.threads.unwrap_or(4),
        })
    }
}

// Usage
let server = ServerBuilder::new()
    .host("0.0.0.0")
    .port(3000)
    .threads(8)
    .build()?;
```

### Type State Pattern

```rust
struct Connection<State> {
    _state: PhantomData<State>,
}

struct Disconnected;
struct Connected;

impl Connection<Disconnected> {
    fn connect(self) -> Result<Connection<Connected>> {
        // Connection logic
        Ok(Connection { _state: PhantomData })
    }
}

impl Connection<Connected> {
    fn send(&self, data: &[u8]) -> Result<()> {
        // Can only send when connected
        Ok(())
    }
    
    fn disconnect(self) -> Connection<Disconnected> {
        Connection { _state: PhantomData }
    }
}
```

### Newtype Pattern

```rust
#[derive(Debug, Clone, PartialEq, Eq)]
struct UserId(u64);

#[derive(Debug, Clone, PartialEq, Eq)]
struct PostId(u64);

impl UserId {
    fn new(id: u64) -> Self {
        Self(id)
    }
}

// Prevents mixing up user IDs and post IDs
fn get_posts_by_user(user_id: UserId) -> Vec<Post> {
    // Implementation
}
```

## Development Workflow

### Code Organization

```
my_project/
├── Cargo.toml
├── src/
│   ├── main.rs          # Binary entry point
│   ├── lib.rs           # Library root
│   ├── config.rs        # Configuration
│   ├── errors.rs        # Error types
│   ├── models/          # Data structures
│   │   ├── mod.rs
│   │   ├── user.rs
│   │   └── post.rs
│   ├── handlers/        # Business logic
│   │   ├── mod.rs
│   │   └── api.rs
│   └── utils/           # Utilities
│       ├── mod.rs
│       └── validation.rs
├── tests/
│   └── integration_test.rs
├── benches/
│   └── benchmark.rs
└── examples/
    └── example.rs
```

### CI/CD Configuration

```yaml
# .github/workflows/rust.yml
name: Rust

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - uses: dtolnay/rust-toolchain@stable
      with:
        components: rustfmt, clippy
    
    - name: Check formatting
      run: cargo fmt -- --check
    
    - name: Clippy
      run: cargo clippy -- -D warnings
    
    - name: Test
      run: cargo test
    
    - name: Audit
      run: cargo audit
```

### Pre-commit Hooks

```bash
#!/bin/bash
# .git/hooks/pre-commit

# Format code
cargo fmt

# Run clippy
cargo clippy -- -D warnings || exit 1

# Run tests
cargo test || exit 1

# Check for TODOs
if grep -r "TODO\|FIXME\|XXX" src/; then
    echo "Found TODOs in code"
    exit 1
fi
```

## Debugging Tips

### Using dbg! Macro

```rust
let x = 5;
let y = dbg!(x * 2) + 1;  // Prints: [src/main.rs:2] x * 2 = 10
dbg!(&y);  // Prints: [src/main.rs:3] &y = 11
```

### Logging with tracing

```rust
use tracing::{debug, info, warn, error, instrument};

#[instrument]
fn process_request(id: u64) -> Result<Response> {
    debug!("Processing request");
    
    let data = fetch_data(id)?;
    info!(%id, data_len = data.len(), "Data fetched successfully");
    
    if data.is_empty() {
        warn!("Empty data returned");
    }
    
    Ok(Response::new(data))
}

// Initialize tracing
fn main() {
    tracing_subscriber::fmt()
        .with_max_level(tracing::Level::DEBUG)
        .init();
    
    // Your application code
}
```

### Using rust-gdb/rust-lldb

```bash
# Compile with debug symbols
cargo build

# Debug with gdb
rust-gdb target/debug/my_program

# Or with lldb
rust-lldb target/debug/my_program

# Common gdb commands:
# break main              # Set breakpoint at main
# run                     # Run the program
# next                    # Step over
# step                    # Step into
# print variable_name     # Print variable value
# backtrace              # Show stack trace
```

Remember: 
- Always handle errors explicitly
- Use `cargo clippy` and `cargo fmt` regularly
- Write tests for all public APIs
- Document your code with examples
- Keep dependencies minimal and up to date
- Profile before optimizing
- Prefer safe code, use `unsafe` only when necessary with proper documentation