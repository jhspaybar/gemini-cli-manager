#[cfg(test)]
mod tests {
    use gemini_cli_manager::storage::Storage;
    use crate::test_utils::{create_temp_storage, McpFixtures, ExtensionBuilder, ProfileBuilder};
    use chrono::Utc;
    
    #[test]
    fn test_extension_crud_operations() {
        let (storage, _temp) = create_temp_storage();
        
        // Create
        let ext = McpFixtures::echo_extension();
        storage.save_extension(&ext).expect("Failed to save extension");
        
        // Read
        let loaded = storage.load_extension(&ext.id).expect("Failed to load extension");
        assert_eq!(loaded.id, ext.id);
        assert_eq!(loaded.name, ext.name);
        assert_eq!(loaded.version, ext.version);
        assert_eq!(loaded.description, ext.description);
        
        // Update
        let mut updated = loaded.clone();
        updated.version = "2.0.0".to_string();
        updated.description = Some("Updated description".to_string());
        storage.save_extension(&updated).expect("Failed to update extension");
        
        let reloaded = storage.load_extension(&ext.id).expect("Failed to reload extension");
        assert_eq!(reloaded.version, "2.0.0");
        assert_eq!(reloaded.description, Some("Updated description".to_string()));
        
        // Delete
        storage.delete_extension(&ext.id).expect("Failed to delete extension");
        assert!(storage.load_extension(&ext.id).is_err());
    }
    
    #[test]
    fn test_extension_list_operations() {
        let (storage, _temp) = create_temp_storage();
        
        // Start empty
        assert!(storage.list_extensions().unwrap().is_empty());
        
        // Add multiple extensions
        let ext1 = McpFixtures::echo_extension();
        let ext2 = McpFixtures::multi_server_extension();
        let ext3 = McpFixtures::context_only_extension();
        
        storage.save_extension(&ext1).unwrap();
        storage.save_extension(&ext2).unwrap();
        storage.save_extension(&ext3).unwrap();
        
        // List all
        let extensions = storage.list_extensions().unwrap();
        assert_eq!(extensions.len(), 3);
        
        // Verify all IDs present
        let ids: Vec<String> = extensions.iter().map(|e| e.id.clone()).collect();
        assert!(ids.contains(&ext1.id));
        assert!(ids.contains(&ext2.id));
        assert!(ids.contains(&ext3.id));
    }
    
    #[test]
    fn test_extension_context_persistence() {
        let (storage, _temp) = create_temp_storage();
        
        // Extension with context file
        let ext = McpFixtures::full_featured_extension();
        storage.save_extension(&ext).unwrap();
        
        let loaded = storage.load_extension(&ext.id).unwrap();
        assert_eq!(loaded.context_file_name, ext.context_file_name);
        assert_eq!(loaded.context_content, ext.context_content);
        assert!(loaded.context_content.unwrap().contains("Advanced Full-Featured Extension"));
    }
    
    #[test]
    fn test_mcp_server_persistence() {
        let (storage, _temp) = create_temp_storage();
        
        let ext = McpFixtures::multi_server_extension();
        storage.save_extension(&ext).unwrap();
        
        let loaded = storage.load_extension(&ext.id).unwrap();
        assert_eq!(loaded.mcp_servers.len(), 3);
        
        // Verify echo server
        let echo_server = loaded.mcp_servers.get("echo").unwrap();
        assert_eq!(echo_server.command, Some("node".to_string()));
        assert!(echo_server.args.is_some());
        
        // Verify python server with environment
        let python_server = loaded.mcp_servers.get("python-echo").unwrap();
        assert_eq!(python_server.command, Some("python".to_string()));
        assert_eq!(python_server.cwd, Some("./servers".to_string()));
        assert!(python_server.env.is_some());
        let env = python_server.env.as_ref().unwrap();
        assert_eq!(env.get("ECHO_PREFIX"), Some(&"[ECHO]".to_string()));
        
        // Verify env-configured server
        let api_server = loaded.mcp_servers.get("api-server").unwrap();
        assert!(api_server.command.is_some());
        assert_eq!(api_server.timeout, Some(10000));
    }
    
    #[test]
    fn test_profile_crud_operations() {
        let (storage, _temp) = create_temp_storage();
        
        // Create profile
        let profile = ProfileBuilder::new("Test Profile")
            .with_description("Test description")
            .with_extensions(vec!["ext1", "ext2"])
            .with_tags(vec!["test", "development"])
            .build();
            
        storage.save_profile(&profile).unwrap();
        
        // Read
        let loaded = storage.load_profile(&profile.id).unwrap();
        assert_eq!(loaded.name, "Test Profile");
        assert_eq!(loaded.extension_ids.len(), 2);
        assert_eq!(loaded.metadata.tags.len(), 2);
        
        // Update
        let mut updated = loaded.clone();
        updated.extension_ids.push("ext3".to_string());
        storage.save_profile(&updated).unwrap();
        
        let reloaded = storage.load_profile(&profile.id).unwrap();
        assert_eq!(reloaded.extension_ids.len(), 3);
        
        // Delete
        storage.delete_profile(&profile.id).unwrap();
        assert!(storage.load_profile(&profile.id).is_err());
    }
    
    #[test]
    fn test_profile_default_handling() {
        let (storage, _temp) = create_temp_storage();
        
        // Create profiles
        let profile1 = ProfileBuilder::new("Profile 1").build();
        let profile2 = ProfileBuilder::new("Profile 2").as_default().build();
        let profile3 = ProfileBuilder::new("Profile 3").build();
        
        storage.save_profile(&profile1).unwrap();
        storage.save_profile(&profile2).unwrap();
        storage.save_profile(&profile3).unwrap();
        
        // Check default
        let default = storage.get_default_profile().unwrap().unwrap();
        assert_eq!(default.id, profile2.id);
        
        // Change default
        storage.set_default_profile(&profile3.id).unwrap();
        
        // Verify change
        let new_default = storage.get_default_profile().unwrap().unwrap();
        assert_eq!(new_default.id, profile3.id);
        
        // Old default should not be default anymore
        let old_default = storage.load_profile(&profile2.id).unwrap();
        assert!(!old_default.metadata.is_default);
    }
    
    #[test]
    fn test_profile_environment_variables() {
        let (storage, _temp) = create_temp_storage();
        
        let mut profile = ProfileBuilder::new("Env Test").build();
        profile.environment_variables.insert("API_KEY".to_string(), "secret123".to_string());
        profile.environment_variables.insert("DATABASE_URL".to_string(), "postgres://localhost".to_string());
        profile.working_directory = Some("~/projects/test".to_string());
        
        storage.save_profile(&profile).unwrap();
        
        let loaded = storage.load_profile(&profile.id).unwrap();
        assert_eq!(loaded.environment_variables.len(), 2);
        assert_eq!(loaded.environment_variables.get("API_KEY"), Some(&"secret123".to_string()));
        assert_eq!(loaded.working_directory, Some("~/projects/test".to_string()));
    }
    
    #[test]
    fn test_corrupted_json_handling() {
        let (storage, temp) = create_temp_storage();
        
        // Save valid extension
        let ext = McpFixtures::echo_extension();
        storage.save_extension(&ext).unwrap();
        
        // Corrupt the JSON file
        let json_path = temp.path()
            .join("extensions")
            .join(format!("{}.json", ext.id));
        std::fs::write(&json_path, "{ invalid json").unwrap();
        
        // Should handle error gracefully
        assert!(storage.load_extension(&ext.id).is_err());
        
        // List should skip corrupted file
        let extensions = storage.list_extensions().unwrap();
        assert_eq!(extensions.len(), 0);
    }
    
    #[test]
    fn test_extension_metadata_persistence() {
        let (storage, _temp) = create_temp_storage();
        
        let ext = ExtensionBuilder::new("Metadata Test")
            .with_version("1.2.3")
            .with_tags(vec!["tag1", "tag2", "tag3"])
            .build();
            
        storage.save_extension(&ext).unwrap();
        
        let loaded = storage.load_extension(&ext.id).unwrap();
        assert_eq!(loaded.metadata.tags.len(), 3);
        assert!(loaded.metadata.imported_at <= Utc::now());
    }
    
    #[test]
    fn test_profile_metadata_persistence() {
        let (storage, _temp) = create_temp_storage();
        
        let profile = ProfileBuilder::new("Metadata Test")
            .with_tags(vec!["prod", "main"])
            .build();
            
        storage.save_profile(&profile).unwrap();
        
        let loaded = storage.load_profile(&profile.id).unwrap();
        assert_eq!(loaded.metadata.tags, vec!["prod", "main"]);
        assert!(loaded.metadata.created_at <= Utc::now());
        assert!(loaded.metadata.updated_at <= Utc::now());
        assert_eq!(loaded.metadata.icon, None);
    }
    
    #[test]
    fn test_extension_id_generation() {
        let ext1 = ExtensionBuilder::new("Test Extension").build();
        assert_eq!(ext1.id, "test-extension");
        
        let ext2 = ExtensionBuilder::new("Another Test Extension").build();
        assert_eq!(ext2.id, "another-test-extension");
        
        let ext3 = ExtensionBuilder::new("Test_With_Underscores").build();
        assert_eq!(ext3.id, "test-with-underscores");
    }
    
    #[test]
    fn test_storage_persistence_across_instances() {
        let temp_dir = tempfile::TempDir::new().unwrap();
        let data_dir = temp_dir.path().to_path_buf();
        
        // First instance - save data
        {
            let storage1 = Storage::with_data_dir(data_dir.clone());
            // Create directories without initializing mock data
            std::fs::create_dir_all(data_dir.join("extensions")).unwrap();
            std::fs::create_dir_all(data_dir.join("profiles")).unwrap();
            
            let ext = McpFixtures::echo_extension();
            storage1.save_extension(&ext).unwrap();
            
            let profile = ProfileBuilder::new("Test").build();
            storage1.save_profile(&profile).unwrap();
        }
        
        // Second instance - load data
        {
            let storage2 = Storage::with_data_dir(data_dir);
            
            let extensions = storage2.list_extensions().unwrap();
            assert_eq!(extensions.len(), 1);
            assert_eq!(extensions[0].id, "echo-test");
            
            let profiles = storage2.list_profiles().unwrap();
            assert_eq!(profiles.len(), 1);
            assert_eq!(profiles[0].id, "test");
        }
    }
}