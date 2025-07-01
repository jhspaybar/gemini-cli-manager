# Rust Development Guide

## Core Requirements

1. **Always use Cargo** - Never compile with rustc directly
2. **Treat warnings as errors** - Add `#![deny(warnings)]` to your crate root
3. **Format before committing** - Run `cargo fmt` on all code
4. **Test before completing** - Run `cargo test` and ensure all tests pass
5. **Check your code** - Run `cargo clippy` and fix all lints

## Development Workflow

### Before Starting Work
```bash
cargo check              # Verify project compiles
cargo test              # Ensure tests pass
cargo fmt -- --check    # Check formatting
```

### During Development
```bash
cargo check             # Fast compilation check
cargo clippy           # Catch common mistakes
cargo test <test_name> # Run specific tests
```

### Before Completing Work
```bash
cargo fmt              # Format all code
cargo test            # Run all tests
cargo clippy -- -D warnings  # Ensure no warnings
cargo build --release # Verify release build works
```

## Required Practices

### Error Handling
- Use `Result<T, E>` for fallible operations
- Prefer `?` operator over `unwrap()` in production code
- Provide context with `expect()` when panicking is acceptable
- Use `thiserror` or `anyhow` for error management

### Memory Safety
- Prefer borrowing over cloning
- Use `Arc<T>` for shared ownership, `Rc<T>` for single-threaded
- Avoid `unsafe` unless absolutely necessary and well-documented
- Always validate array/slice indices

### Code Organization
- One module per file for clarity
- Keep functions small and focused
- Use descriptive names (no abbreviations)
- Document public APIs with `///` comments

## What NOT to Do
- Never ignore compiler warnings
- Never use `unwrap()` without careful consideration
- Never skip running tests before pushing code
- Never commit code that doesn't pass `cargo fmt`
- Never use `std::process::exit()` in library code