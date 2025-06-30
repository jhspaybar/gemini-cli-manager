use std::fs;
use std::path::{Path, PathBuf};

use color_eyre::{eyre::eyre, Result};
use serde::{de::DeserializeOwned, Serialize};

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

    /// Initialize storage with default data if empty
    pub fn init(&self) -> Result<()> {
        // Create subdirectories
        fs::create_dir_all(self.data_dir.join("extensions"))?;
        fs::create_dir_all(self.data_dir.join("profiles"))?;
        
        // If no extensions exist, create mock data
        if self.list_extensions()?.is_empty() {
            println!("Initializing with sample extensions...");
            for extension in Extension::mock_extensions() {
                self.save_extension(&extension)?;
            }
        }
        
        // If no profiles exist, create mock data
        if self.list_profiles()?.is_empty() {
            println!("Initializing with sample profiles...");
            for profile in Profile::mock_profiles() {
                self.save_profile(&profile)?;
            }
        }
        
        Ok(())
    }

    // Extension methods

    /// Save an extension to storage
    pub fn save_extension(&self, extension: &Extension) -> Result<()> {
        let path = self.data_dir.join("extensions").join(format!("{}.json", extension.id));
        self.save_json(&path, extension)
    }

    /// Load an extension by ID
    pub fn load_extension(&self, id: &str) -> Result<Extension> {
        let path = self.data_dir.join("extensions").join(format!("{}.json", id));
        self.load_json(&path)
    }

    /// List all extensions
    pub fn list_extensions(&self) -> Result<Vec<Extension>> {
        self.list_items("extensions")
    }

    /// Delete an extension
    #[allow(dead_code)]
    pub fn delete_extension(&self, id: &str) -> Result<()> {
        let path = self.data_dir.join("extensions").join(format!("{}.json", id));
        if path.exists() {
            fs::remove_file(path)?;
        }
        Ok(())
    }

    // Profile methods

    /// Save a profile to storage
    pub fn save_profile(&self, profile: &Profile) -> Result<()> {
        let path = self.data_dir.join("profiles").join(format!("{}.json", profile.id));
        self.save_json(&path, profile)
    }

    /// Load a profile by ID
    pub fn load_profile(&self, id: &str) -> Result<Profile> {
        let path = self.data_dir.join("profiles").join(format!("{}.json", id));
        self.load_json(&path)
    }

    /// List all profiles
    pub fn list_profiles(&self) -> Result<Vec<Profile>> {
        self.list_items("profiles")
    }

    /// Delete a profile
    pub fn delete_profile(&self, id: &str) -> Result<()> {
        let path = self.data_dir.join("profiles").join(format!("{}.json", id));
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
            for entry in fs::read_dir(dir)? {
                let entry = entry?;
                let path = entry.path();
                
                if path.extension().and_then(|s| s.to_str()) == Some("json") {
                    match self.load_json::<T>(&path) {
                        Ok(item) => items.push(item),
                        Err(e) => eprintln!("Warning: Failed to load {:?}: {}", path, e),
                    }
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
        let extensions = Extension::mock_extensions();
        let ext = &extensions[0];
        
        storage.save_extension(ext).unwrap();
        let loaded = storage.load_extension(&ext.id).unwrap();
        
        assert_eq!(loaded.id, ext.id);
        assert_eq!(loaded.name, ext.name);
    }

    #[test]
    fn test_list_extensions() {
        let (storage, _temp) = test_storage();
        let extensions = storage.list_extensions().unwrap();
        assert!(!extensions.is_empty());
    }

    #[test]
    fn test_save_and_load_profile() {
        let (storage, _temp) = test_storage();
        let profiles = Profile::mock_profiles();
        let profile = &profiles[0];
        
        storage.save_profile(profile).unwrap();
        let loaded = storage.load_profile(&profile.id).unwrap();
        
        assert_eq!(loaded.id, profile.id);
        assert_eq!(loaded.name, profile.name);
    }

    #[test]
    fn test_default_profile() {
        let (storage, _temp) = test_storage();
        
        let default = storage.get_default_profile().unwrap();
        assert!(default.is_some());
        assert!(default.unwrap().metadata.is_default);
        
        // Set a different profile as default
        storage.set_default_profile("minimal").unwrap();
        let new_default = storage.get_default_profile().unwrap().unwrap();
        assert_eq!(new_default.id, "minimal");
    }
}