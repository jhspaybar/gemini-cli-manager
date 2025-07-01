#[cfg(test)]
mod tests {
    use chrono::Utc;
    use std::collections::HashMap;
    use gemini_cli_manager::models::profile::LaunchConfig;
    
    // Since we can't directly test main(), we'll test the list_storage_contents functionality
    // by extracting its logic into a testable function
    
    #[test]
    fn test_cli_list_storage_output() {
        // This tests that the --list-storage flag would produce output
        // We can't easily test the actual main() function, but we can verify
        // the storage listing logic works
        
        // Create a temporary storage location
        let temp_dir = tempfile::tempdir().unwrap();
        
        // Create storage with custom directory
        let storage = gemini_cli_manager::storage::Storage::with_data_dir(temp_dir.path().to_path_buf());
        storage.init().unwrap();
        
        // Clear any mock data
        for ext in storage.list_extensions().unwrap() {
            storage.delete_extension(&ext.id).unwrap();
        }
        for profile in storage.list_profiles().unwrap() {
            storage.delete_profile(&profile.id).unwrap();
        }
        
        // Add test data
        let extension = gemini_cli_manager::models::Extension {
            id: "test-ext".to_string(),
            name: "Test Extension".to_string(),
            version: "1.0.0".to_string(),
            description: Some("A test extension".to_string()),
            mcp_servers: HashMap::new(),
            context_file_name: None,
            context_content: None,
            metadata: gemini_cli_manager::models::extension::ExtensionMetadata {
                imported_at: Utc::now(),
                source_path: None,
                tags: vec![],
            },
        };
        storage.save_extension(&extension).unwrap();
        
        let profile = gemini_cli_manager::models::Profile {
            id: "test-profile".to_string(),
            name: "Test Profile".to_string(),
            description: Some("A test profile".to_string()),
            extension_ids: vec!["test-ext".to_string()],
            environment_variables: HashMap::new(),
            working_directory: None,
            launch_config: LaunchConfig::default(),
            metadata: gemini_cli_manager::models::profile::ProfileMetadata {
                created_at: Utc::now(),
                updated_at: Utc::now(),
                tags: vec!["test".to_string()],
                is_default: false,
                icon: None,
            },
        };
        storage.save_profile(&profile).unwrap();
        
        // Verify we can list the data
        let extensions = storage.list_extensions().unwrap();
        assert_eq!(extensions.len(), 1);
        assert_eq!(extensions[0].name, "Test Extension");
        
        let profiles = storage.list_profiles().unwrap();
        assert_eq!(profiles.len(), 1);
        assert_eq!(profiles[0].name, "Test Profile");
    }
    
    #[test]
    fn test_theme_initialization() {
        // Test that theme is set to Mocha
        use gemini_cli_manager::theme::ThemeFlavour;
        
        // Main sets the theme to Mocha
        gemini_cli_manager::theme::set_flavour(ThemeFlavour::Mocha);
        // We can't test get_flavour because it's not exposed, but set_flavour doesn't panic
    }
    
    
    #[test]
    fn test_list_storage_with_empty_storage() {
        // Test listing empty storage
        let temp_dir = tempfile::tempdir().unwrap();
        
        let storage = gemini_cli_manager::storage::Storage::with_data_dir(temp_dir.path().to_path_buf());
        // Don't call init() to avoid mock data
        
        let extensions = storage.list_extensions().unwrap_or_default();
        assert_eq!(extensions.len(), 0);
        
        let profiles = storage.list_profiles().unwrap_or_default();
        assert_eq!(profiles.len(), 0);
    }
    
    #[test]
    fn test_list_storage_with_multiple_items() {
        // Test with multiple extensions and profiles
        let temp_dir = tempfile::tempdir().unwrap();
        
        let storage = gemini_cli_manager::storage::Storage::with_data_dir(temp_dir.path().to_path_buf());
        storage.init().unwrap();
        
        // Clear any mock data
        for ext in storage.list_extensions().unwrap() {
            storage.delete_extension(&ext.id).unwrap();
        }
        for profile in storage.list_profiles().unwrap() {
            storage.delete_profile(&profile.id).unwrap();
        }
        
        // Add multiple extensions
        for i in 1..=3 {
            let extension = gemini_cli_manager::models::Extension {
                id: format!("ext-{}", i),
                name: format!("Extension {}", i),
                version: "1.0.0".to_string(),
                description: if i == 2 { None } else { Some(format!("Description {}", i)) },
                mcp_servers: HashMap::new(),
                context_file_name: None,
                context_content: None,
                metadata: gemini_cli_manager::models::extension::ExtensionMetadata {
                    imported_at: Utc::now(),
                    source_path: None,
                    tags: vec![],
                },
            };
            storage.save_extension(&extension).unwrap();
        }
        
        // Add multiple profiles
        for i in 1..=2 {
            let profile = gemini_cli_manager::models::Profile {
                id: format!("profile-{}", i),
                name: format!("Profile {}", i),
                description: if i == 1 { Some("First profile".to_string()) } else { None },
                extension_ids: vec![format!("ext-{}", i)],
                environment_variables: HashMap::new(),
                working_directory: None,
                launch_config: LaunchConfig::default(),
                metadata: gemini_cli_manager::models::profile::ProfileMetadata {
                    created_at: Utc::now(),
                    updated_at: Utc::now(),
                    tags: if i == 1 { vec!["tag1".to_string(), "tag2".to_string()] } else { vec![] },
                    is_default: false,
                    icon: None,
                },
            };
            storage.save_profile(&profile).unwrap();
        }
        
        let extensions = storage.list_extensions().unwrap();
        assert_eq!(extensions.len(), 3);
        
        let profiles = storage.list_profiles().unwrap();
        assert_eq!(profiles.len(), 2);
        
        // Verify the profile with tags
        let profile_with_tags = profiles.iter().find(|p| p.id == "profile-1").unwrap();
        assert_eq!(profile_with_tags.metadata.tags.len(), 2);
    }
}