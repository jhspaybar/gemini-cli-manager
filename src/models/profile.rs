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
    
    /// Metadata
    pub metadata: ProfileMetadata,
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
    /// Create mock profiles for testing
    pub fn mock_profiles() -> Vec<Profile> {
        vec![
            Profile {
                id: "dev-fullstack".to_string(),
                name: "Full Stack Development".to_string(),
                description: Some("Complete development environment with GitHub, database, and filesystem tools".to_string()),
                extension_ids: vec![
                    "github-tools".to_string(),
                    "database-tools".to_string(),
                    "filesystem-enhanced".to_string(),
                ],
                environment_variables: {
                    let mut env = HashMap::new();
                    env.insert("GITHUB_TOKEN".to_string(), "ghp_xxx...".to_string());
                    env.insert("DATABASE_URL".to_string(), "postgresql://localhost/devdb".to_string());
                    env.insert("NODE_ENV".to_string(), "development".to_string());
                    env
                },
                working_directory: Some("~/projects".to_string()),
                metadata: ProfileMetadata {
                    created_at: Utc::now(),
                    updated_at: Utc::now(),
                    tags: vec!["development".to_string(), "fullstack".to_string()],
                    is_default: true,
                    icon: Some("🚀".to_string()),
                },
            },
            Profile {
                id: "data-science".to_string(),
                name: "Data Science Workspace".to_string(),
                description: Some("Python-based data science environment with database access".to_string()),
                extension_ids: vec![
                    "database-tools".to_string(),
                ],
                environment_variables: {
                    let mut env = HashMap::new();
                    env.insert("DATABASE_URL".to_string(), "postgresql://localhost/analytics".to_string());
                    env.insert("JUPYTER_PORT".to_string(), "8888".to_string());
                    env.insert("PYTHONPATH".to_string(), "~/notebooks/lib".to_string());
                    env
                },
                working_directory: Some("~/notebooks".to_string()),
                metadata: ProfileMetadata {
                    created_at: Utc::now(),
                    updated_at: Utc::now(),
                    tags: vec!["data-science".to_string(), "analytics".to_string(), "python".to_string()],
                    is_default: false,
                    icon: Some("📊".to_string()),
                },
            },
            Profile {
                id: "minimal".to_string(),
                name: "Minimal".to_string(),
                description: Some("Basic profile with just filesystem tools".to_string()),
                extension_ids: vec![
                    "filesystem-enhanced".to_string(),
                ],
                environment_variables: HashMap::new(),
                working_directory: None,
                metadata: ProfileMetadata {
                    created_at: Utc::now(),
                    updated_at: Utc::now(),
                    tags: vec!["minimal".to_string(), "basic".to_string()],
                    is_default: false,
                    icon: Some("📁".to_string()),
                },
            },
            Profile {
                id: "github-only".to_string(),
                name: "GitHub Focus".to_string(),
                description: Some("GitHub-only environment for code reviews and PR management".to_string()),
                extension_ids: vec![
                    "github-tools".to_string(),
                ],
                environment_variables: {
                    let mut env = HashMap::new();
                    env.insert("GITHUB_TOKEN".to_string(), "ghp_xxx...".to_string());
                    env.insert("GITHUB_DEFAULT_BRANCH".to_string(), "main".to_string());
                    env
                },
                working_directory: Some("~/github".to_string()),
                metadata: ProfileMetadata {
                    created_at: Utc::now(),
                    updated_at: Utc::now(),
                    tags: vec!["github".to_string(), "vcs".to_string(), "collaboration".to_string()],
                    is_default: false,
                    icon: Some("🐙".to_string()),
                },
            },
        ]
    }
    
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