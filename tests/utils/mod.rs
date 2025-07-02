/// Test utilities for Gemini CLI Manager tests
pub mod helpers;
pub mod mcp_fixtures;
pub mod theme_helpers;

pub use helpers::*;
pub use mcp_fixtures::*;
pub use theme_helpers::*;

// Re-export specific functions that might be needed
pub use mcp_fixtures::validate_extension_json;
