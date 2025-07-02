pub mod e2e;
pub mod integration;
pub mod snapshots;
/// Test modules for Gemini CLI Manager
pub mod unit;
pub mod utils;

// Re-export test utilities for easy access
pub use utils::*;

// Alias for backward compatibility
pub mod test_utils {
    pub use super::utils::*;
}
