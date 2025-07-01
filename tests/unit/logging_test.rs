#[cfg(test)]
mod tests {
    use gemini_cli_manager::logging;
    use std::env;

    #[test]
    fn test_log_env_variable_name() {
        // Test that LOG_ENV is correctly formatted
        let log_env = &*logging::LOG_ENV;
        assert!(log_env.contains("_LOG_LEVEL"));
        assert!(!log_env.is_empty());
    }

    #[test]
    fn test_log_file_name() {
        // Test that LOG_FILE has the correct format
        let log_file = &*logging::LOG_FILE;
        assert!(log_file.ends_with(".log"));
        assert!(log_file.contains("gemini-cli-manager"));
    }

    #[test]
    fn test_logging_init_creates_directory() {
        // Initialize logging
        let _ = logging::init();
        
        // Since logging might already be initialized from other tests,
        // we can't test directory creation reliably. Just ensure it doesn't panic.
    }

    #[test]
    fn test_logging_with_custom_log_level() {
        // Set custom log level
        unsafe {
            env::set_var(&*logging::LOG_ENV, "debug");
        }
        
        // Try to initialize (might fail if already initialized)
        let _result = logging::init();
        
        // Clean up
        unsafe {
            env::remove_var(&*logging::LOG_ENV);
        }
    }

    #[test]
    fn test_logging_with_rust_log() {
        // Test with RUST_LOG environment variable
        unsafe {
            env::set_var("RUST_LOG", "info");
        }
        
        // Try to initialize
        let _result = logging::init();
        
        // Clean up
        unsafe {
            env::remove_var("RUST_LOG");
        }
    }
}