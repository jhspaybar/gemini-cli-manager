/// Test utilities for Gemini CLI Manager tests
pub mod helpers;
pub mod theme_helpers;
pub mod mcp_fixtures;

pub use helpers::*;
pub use theme_helpers::*;
pub use mcp_fixtures::*;

// Re-export specific functions that might be needed
pub use mcp_fixtures::validate_extension_json;