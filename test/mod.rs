/// Test modules for Gemini CLI Manager
pub mod unit;
pub mod integration;
pub mod e2e;
pub mod snapshots;
pub mod utils;

// Re-export test utilities for easy access
pub use utils::*;