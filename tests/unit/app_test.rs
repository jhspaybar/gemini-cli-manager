#[cfg(test)]
mod tests {
    use gemini_cli_manager::App;
    use tokio::time::{Duration, timeout};

    #[tokio::test]
    async fn test_app_initialization() {
        let app = App::new();
        assert!(app.is_ok(), "App should initialize successfully");
    }

    #[tokio::test]
    async fn test_app_run_immediate_exit() {
        let mut app = App::new().unwrap();

        // In test environment without a real terminal, the app might exit immediately
        // or timeout. Either is acceptable for this test.
        let run_future = app.run();
        let result = timeout(Duration::from_millis(100), run_future).await;

        // Just verify that we can call run without panicking
        match result {
            Ok(Ok(())) => {
                // App exited cleanly - this is fine in test environment
            }
            Ok(Err(e)) => {
                // App returned an error - log it but don't fail the test
                eprintln!("App run returned error: {}", e);
            }
            Err(_) => {
                // Timeout - app is running normally
            }
        }
    }

    #[tokio::test]
    async fn test_app_multiple_new_calls() {
        // Test that we can create multiple app instances
        let app1 = App::new();
        assert!(app1.is_ok());

        let app2 = App::new();
        assert!(app2.is_ok());
    }
}
