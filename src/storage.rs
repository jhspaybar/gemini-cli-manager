use std::fs;
use std::path::{Path, PathBuf};

use color_eyre::{Result, eyre::eyre};
use serde::{Serialize, de::DeserializeOwned};

use crate::models::{Extension, Profile};

/// Storage manager for persisting application data
#[derive(Clone)]
pub struct Storage {
    data_dir: PathBuf,
}

impl Storage {
    /// Create a new storage instance with the default data directory
    pub fn new() -> Result<Self> {
        let data_dir = Self::get_data_dir()?;
        Ok(Self { data_dir })
    }

    /// Create a storage instance with a custom data directory
    #[allow(dead_code)]
    pub fn with_data_dir(data_dir: PathBuf) -> Self {
        Self { data_dir }
    }

    /// Get the default data directory for the application
    fn get_data_dir() -> Result<PathBuf> {
        let data_dir = dirs::data_dir()
            .ok_or_else(|| eyre!("Could not determine data directory"))?
            .join("gemini-cli-manager");

        // Ensure the directory exists
        fs::create_dir_all(&data_dir)?;

        Ok(data_dir)
    }

    /// Initialize storage directories
    pub fn init(&self) -> Result<()> {
        // Create subdirectories
        fs::create_dir_all(self.data_dir.join("extensions"))?;
        fs::create_dir_all(self.data_dir.join("profiles"))?;

        // Extensions should be imported from actual extension packages
        // Profiles should be created by users

        Ok(())
    }

    // Extension methods

    /// Save an extension to storage
    pub fn save_extension(&self, extension: &Extension) -> Result<()> {
        let path = self
            .data_dir
            .join("extensions")
            .join(format!("{}.json", extension.id));
        self.save_json(&path, extension)
    }

    /// Load an extension by ID
    pub fn load_extension(&self, id: &str) -> Result<Extension> {
        let path = self.data_dir.join("extensions").join(format!("{id}.json"));
        self.load_json(&path)
    }

    /// List all extensions
    pub fn list_extensions(&self) -> Result<Vec<Extension>> {
        self.list_items("extensions")
    }

    /// Delete an extension
    #[allow(dead_code)]
    pub fn delete_extension(&self, id: &str) -> Result<()> {
        let path = self.data_dir.join("extensions").join(format!("{id}.json"));
        if path.exists() {
            fs::remove_file(path)?;
        }
        Ok(())
    }

    // Profile methods

    /// Save a profile to storage
    pub fn save_profile(&self, profile: &Profile) -> Result<()> {
        let path = self
            .data_dir
            .join("profiles")
            .join(format!("{}.json", profile.id));
        self.save_json(&path, profile)
    }

    /// Load a profile by ID
    pub fn load_profile(&self, id: &str) -> Result<Profile> {
        let path = self.data_dir.join("profiles").join(format!("{id}.json"));
        let profile: Profile = self.load_json(&path)?;

        // Ensure backward compatibility - if launch_config is missing, it will use default
        // This is handled by serde's #[serde(default)] attribute on the field

        Ok(profile)
    }

    /// List all profiles
    pub fn list_profiles(&self) -> Result<Vec<Profile>> {
        self.list_items("profiles")
    }

    /// Delete a profile
    pub fn delete_profile(&self, id: &str) -> Result<()> {
        let path = self.data_dir.join("profiles").join(format!("{id}.json"));
        if path.exists() {
            fs::remove_file(path)?;
        }
        Ok(())
    }

    /// Get the default profile
    #[allow(dead_code)]
    pub fn get_default_profile(&self) -> Result<Option<Profile>> {
        let profiles = self.list_profiles()?;
        Ok(profiles.into_iter().find(|p| p.metadata.is_default))
    }

    /// Set a profile as default
    #[allow(dead_code)]
    pub fn set_default_profile(&self, id: &str) -> Result<()> {
        let mut profiles = self.list_profiles()?;

        for profile in &mut profiles {
            profile.metadata.is_default = profile.id == id;
            self.save_profile(profile)?;
        }

        Ok(())
    }

    // Helper methods

    /// Save data as JSON
    fn save_json<T: Serialize>(&self, path: &Path, data: &T) -> Result<()> {
        let json = serde_json::to_string_pretty(data)?;
        fs::write(path, json)?;
        Ok(())
    }

    /// Load data from JSON
    fn load_json<T: DeserializeOwned>(&self, path: &Path) -> Result<T> {
        let json = fs::read_to_string(path)?;
        let data = serde_json::from_str(&json)?;
        Ok(data)
    }

    /// List all items in a subdirectory
    fn list_items<T: DeserializeOwned>(&self, subdir: &str) -> Result<Vec<T>> {
        let dir = self.data_dir.join(subdir);
        let mut items = Vec::new();

        if dir.exists() {
            // Collect all paths first to sort them
            let mut paths: Vec<PathBuf> = fs::read_dir(dir)?
                .filter_map(|entry| entry.ok())
                .map(|entry| entry.path())
                .filter(|path| path.extension().and_then(|s| s.to_str()) == Some("json"))
                .collect();
            
            // Sort paths to ensure consistent ordering
            paths.sort();

            // Load items in sorted order
            for path in paths {
                match self.load_json::<T>(&path) {
                    Ok(item) => items.push(item),
                    Err(e) => eprintln!("Warning: Failed to load {path:?}: {e}"),
                }
            }
        }

        Ok(items)
    }

    /// Get the data directory path
    #[allow(dead_code)]
    pub fn data_dir(&self) -> &Path {
        &self.data_dir
    }
}

impl Default for Storage {
    fn default() -> Self {
        Self::new().unwrap_or_else(|_| Self {
            data_dir: PathBuf::from(".gemini-cli-manager-data"),
        })
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use chrono::Utc;
    use std::collections::HashMap;
    use tempfile::TempDir;

    fn test_storage() -> (Storage, TempDir) {
        let temp_dir = TempDir::new().unwrap();
        let storage = Storage::with_data_dir(temp_dir.path().to_path_buf());
        storage.init().unwrap();
        (storage, temp_dir)
    }

    #[test]
    fn test_save_and_load_extension() {
        let (storage, _temp) = test_storage();

        // Create a test extension
        let ext = Extension {
            id: "test-ext".to_string(),
            name: "Test Extension".to_string(),
            version: "1.0.0".to_string(),
            description: Some("Test description".to_string()),
            mcp_servers: HashMap::new(),
            context_file_name: None,
            context_content: None,
            metadata: crate::models::extension::ExtensionMetadata {
                imported_at: Utc::now(),
                source_path: None,
                tags: vec!["test".to_string()],
            },
        };

        storage.save_extension(&ext).unwrap();
        let loaded = storage.load_extension(&ext.id).unwrap();

        assert_eq!(loaded.id, ext.id);
        assert_eq!(loaded.name, ext.name);
    }

    #[test]
    fn test_list_extensions() {
        let (storage, _temp) = test_storage();

        // Initially should be empty
        let extensions = storage.list_extensions().unwrap();
        assert!(extensions.is_empty());

        // Add an extension
        let ext = Extension {
            id: "test-list".to_string(),
            name: "Test List".to_string(),
            version: "1.0.0".to_string(),
            description: None,
            mcp_servers: HashMap::new(),
            context_file_name: None,
            context_content: None,
            metadata: crate::models::extension::ExtensionMetadata {
                imported_at: Utc::now(),
                source_path: None,
                tags: vec![],
            },
        };
        storage.save_extension(&ext).unwrap();

        let extensions = storage.list_extensions().unwrap();
        assert_eq!(extensions.len(), 1);
    }

    #[test]
    fn test_save_and_load_profile() {
        let (storage, _temp) = test_storage();

        // Create a test profile
        let profile = Profile {
            id: "test-profile".to_string(),
            name: "Test Profile".to_string(),
            description: Some("Test description".to_string()),
            extension_ids: vec![],
            environment_variables: HashMap::new(),
            working_directory: None,
            launch_config: crate::models::profile::LaunchConfig::default(),
            metadata: crate::models::profile::ProfileMetadata {
                created_at: Utc::now(),
                updated_at: Utc::now(),
                tags: vec!["test".to_string()],
                is_default: false,
                icon: None,
            },
        };

        storage.save_profile(&profile).unwrap();
        let loaded = storage.load_profile(&profile.id).unwrap();

        assert_eq!(loaded.id, profile.id);
        assert_eq!(loaded.name, profile.name);
    }

    #[test]
    fn test_default_profile() {
        let (storage, _temp) = test_storage();

        // Initially no default profile
        let default = storage.get_default_profile().unwrap();
        assert!(default.is_none());

        // Create profiles
        let profile1 = Profile {
            id: "profile1".to_string(),
            name: "Profile 1".to_string(),
            description: None,
            extension_ids: vec![],
            environment_variables: HashMap::new(),
            working_directory: None,
            launch_config: crate::models::profile::LaunchConfig::default(),
            metadata: crate::models::profile::ProfileMetadata {
                created_at: Utc::now(),
                updated_at: Utc::now(),
                tags: vec![],
                is_default: true,
                icon: None,
            },
        };

        let profile2 = Profile {
            id: "profile2".to_string(),
            name: "Profile 2".to_string(),
            description: None,
            extension_ids: vec![],
            environment_variables: HashMap::new(),
            working_directory: None,
            launch_config: crate::models::profile::LaunchConfig::default(),
            metadata: crate::models::profile::ProfileMetadata {
                created_at: Utc::now(),
                updated_at: Utc::now(),
                tags: vec![],
                is_default: false,
                icon: None,
            },
        };

        storage.save_profile(&profile1).unwrap();
        storage.save_profile(&profile2).unwrap();

        let default = storage.get_default_profile().unwrap();
        assert!(default.is_some());
        assert_eq!(default.unwrap().id, "profile1");

        // Set a different profile as default
        storage.set_default_profile("profile2").unwrap();
        let new_default = storage.get_default_profile().unwrap().unwrap();
        assert_eq!(new_default.id, "profile2");
    }
}
