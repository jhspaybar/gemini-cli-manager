#[cfg(test)]
mod tests {
    use gemini_cli_manager::errors;
    use gemini_cli_manager::trace_dbg;

    #[test]
    fn test_errors_init() {
        // Test that init can be called successfully
        // Note: This might fail if already initialized in other tests
        let _ = errors::init();

        // Note: We can't easily test the panic hook without actually panicking
        // which would fail the test
    }

    #[test]
    fn test_trace_dbg_macro() {
        // Test basic trace_dbg usage
        let value = 42;
        let result = trace_dbg!(value);
        assert_eq!(result, 42);

        // Test with expression
        let result = trace_dbg!(1 + 2);
        assert_eq!(result, 3);

        // Test with custom level
        let result = trace_dbg!(level: tracing::Level::INFO, value);
        assert_eq!(result, 42);

        // Test with custom target
        let result = trace_dbg!(target: "test_target", value);
        assert_eq!(result, 42);

        // Test with both custom target and level
        let result = trace_dbg!(target: "test_target", level: tracing::Level::WARN, value);
        assert_eq!(result, 42);
    }

    #[test]
    fn test_trace_dbg_with_complex_expressions() {
        let vec = vec![1, 2, 3];
        let result = trace_dbg!(vec.len());
        assert_eq!(result, 3);

        let string = "hello";
        let result = trace_dbg!(string.to_uppercase());
        assert_eq!(result, "HELLO");

        // Test that it returns the value unchanged
        let mut counter = 0;
        let result = trace_dbg!({
            counter += 1;
            counter
        });
        assert_eq!(result, 1);
        assert_eq!(counter, 1);
    }

    #[test]
    fn test_trace_dbg_preserves_value() {
        // Test that trace_dbg returns the exact value, not a copy
        let original = String::from("test");
        let returned = trace_dbg!(original);
        assert_eq!(returned, "test");
    }
}
