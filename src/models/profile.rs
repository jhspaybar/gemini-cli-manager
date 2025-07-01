use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

/// A profile bundles multiple extensions with environment configuration
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Profile {
    /// Unique identifier
    pub id: String,
    
    /// Display name
    pub name: String,
    
    /// Optional description
    pub description: Option<String>,
    
    /// Extension IDs included in this profile
    pub extension_ids: Vec<String>,
    
    /// Environment variables specific to this profile
    pub environment_variables: HashMap<String, String>,
    
    /// Working directory for Gemini when launched with this profile
    pub working_directory: Option<String>,
    
    /// Launch configuration options
    #[serde(default)]
    pub launch_config: LaunchConfig,
    
    /// Metadata
    pub metadata: ProfileMetadata,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct LaunchConfig {
    /// Clean existing configuration before launching
    pub clean_launch: bool,
    
    /// Remove extensions directory after Gemini exits
    pub cleanup_on_exit: bool,
    
    /// Preserve specific extension IDs even during clean launch
    pub preserve_extensions: Vec<String>,
}

impl Default for LaunchConfig {
    fn default() -> Self {
        Self {
            clean_launch: false,
            cleanup_on_exit: true,  // Default to cleaning up after ourselves
            preserve_extensions: Vec::new(),
        }
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ProfileMetadata {
    /// When the profile was created
    pub created_at: DateTime<Utc>,
    
    /// When the profile was last modified
    pub updated_at: DateTime<Utc>,
    
    /// User-defined tags
    pub tags: Vec<String>,
    
    /// Whether this is the default profile
    pub is_default: bool,
    
    /// Icon/emoji for the profile (optional)
    pub icon: Option<String>,
}

impl Profile {
    // Mock data methods removed - profiles should be created by users
    
    /// Get the display name with optional icon
    pub fn display_name(&self) -> String {
        if let Some(icon) = &self.metadata.icon {
            format!("{} {}", icon, self.name)
        } else {
            self.name.clone()
        }
    }
    
    /// Get a summary of what's included
    pub fn summary(&self) -> String {
        let ext_count = self.extension_ids.len();
        let env_count = self.environment_variables.len();
        
        format!(
            "{} extension{}, {} env var{}",
            ext_count,
            if ext_count == 1 { "" } else { "s" },
            env_count,
            if env_count == 1 { "" } else { "s" }
        )
    }
}