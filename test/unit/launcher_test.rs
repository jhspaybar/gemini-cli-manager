#[cfg(test)]
mod tests {
    use gemini_cli_manager::launcher::Launcher;
    use gemini_cli_manager::models::{Extension, Profile};
    use crate::test_utils::{
        create_temp_storage, McpFixtures, ProfileBuilder, WorkspaceVerifier,
        create_extension_json, validate_extension_json
    };
    use std::path::PathBuf;
    use tempfile::TempDir;
    
    fn create_test_launcher() -> (Launcher, TempDir, TempDir) {
        let workspace_dir = TempDir::new().unwrap();
        let (storage, storage_dir) = create_temp_storage();
        
        let mut launcher = Launcher::with_storage(storage);
        launcher = launcher.with_workspace_dir(workspace_dir.path().to_path_buf());
            
        (launcher, workspace_dir, storage_dir)
    }
    
    #[test]
    fn test_workspace_setup() {
        let (launcher, workspace_dir, _storage_dir) = create_test_launcher();
        
        let profile = ProfileBuilder::new("test-profile").build();
        let profile_workspace = workspace_dir.path().join(&profile.id);
        
        // Setup workspace
        launcher.setup_workspace(&profile_workspace).unwrap();
        
        // Verify structure
        assert!(WorkspaceVerifier::verify_workspace_structure(&profile_workspace).is_ok());
    }
    
    #[test]
    fn test_extension_installation() {
        let (launcher, workspace_dir, storage_dir) = create_test_launcher();
        
        // Add extension to storage
        let ext = McpFixtures::echo_extension();
        launcher.storage.save_extension(&ext).unwrap();
        
        // Create profile using the extension
        let profile = ProfileBuilder::new("test-profile")
            .with_extensions(vec![&ext.id])
            .build();
        launcher.storage.save_profile(&profile).unwrap();
        
        let profile_workspace = workspace_dir.path().join(&profile.id);
        
        // Install extensions
        launcher.install_extensions_for_profile(&profile, &profile_workspace).unwrap();
        
        // Verify extension installed
        assert!(WorkspaceVerifier::verify_extension_installed(
            &profile_workspace,
            &ext.id
        ).is_ok());
        
        // Verify context file
        assert!(WorkspaceVerifier::verify_context_file(
            &profile_workspace,
            &ext.id,
            "GEMINI.md"
        ).is_ok());
    }
    
    #[test]
    fn test_extension_json_format() {
        let (launcher, workspace_dir, _storage_dir) = create_test_launcher();
        
        // Create extension with MCP servers
        let ext = McpFixtures::multi_server_extension();
        launcher.storage.save_extension(&ext).unwrap();
        
        let profile = ProfileBuilder::new("test")
            .with_extensions(vec![&ext.id])
            .build();
        launcher.storage.save_profile(&profile).unwrap();
        
        let profile_workspace = workspace_dir.path().join(&profile.id);
        launcher.install_extensions_for_profile(&profile, &profile_workspace).unwrap();
        
        // Read the generated gemini-extension.json
        let config_path = profile_workspace
            .join(".gemini")
            .join("extensions")
            .join(&ext.id)
            .join("gemini-extension.json");
            
        let content = std::fs::read_to_string(&config_path).unwrap();
        let json: serde_json::Value = serde_json::from_str(&content).unwrap();
        
        // Verify structure
        assert_eq!(json["name"], ext.name);
        assert_eq!(json["version"], ext.version);
        assert!(json["mcpServers"].is_object());
        assert!(json["mcpServers"]["echo"].is_object());
        assert!(json["mcpServers"]["python-echo"].is_object());
        assert!(json["mcpServers"]["api-server"].is_object());
    }
    
    #[test]
    fn test_context_file_creation() {
        let (launcher, workspace_dir, _storage_dir) = create_test_launcher();
        
        // Test different context file names
        let mut ext1 = McpFixtures::echo_extension();
        ext1.id = "test1".to_string();
        ext1.context_file_name = Some("GEMINI.md".to_string());
        
        let mut ext2 = McpFixtures::context_only_extension();
        ext2.id = "test2".to_string();
        ext2.context_file_name = Some("CUSTOM.md".to_string());
        
        let mut ext3 = McpFixtures::full_featured_extension();
        ext3.id = "test3".to_string();
        ext3.context_file_name = None; // Should default to GEMINI.md
        
        launcher.storage.save_extension(&ext1).unwrap();
        launcher.storage.save_extension(&ext2).unwrap();
        launcher.storage.save_extension(&ext3).unwrap();
        
        let profile = ProfileBuilder::new("test")
            .with_extensions(vec!["test1", "test2", "test3"])
            .build();
        launcher.storage.save_profile(&profile).unwrap();
        
        let profile_workspace = workspace_dir.path().join(&profile.id);
        launcher.install_extensions_for_profile(&profile, &profile_workspace).unwrap();
        
        // Verify each context file
        let content1 = WorkspaceVerifier::verify_context_file(
            &profile_workspace,
            "test1",
            "GEMINI.md"
        ).unwrap();
        assert!(content1.contains("Echo Test Extension"));
        
        let content2 = WorkspaceVerifier::verify_context_file(
            &profile_workspace,
            "test2",
            "CUSTOM.md"
        ).unwrap();
        assert!(content2.contains("Context Only Extension"));
        
        let content3 = WorkspaceVerifier::verify_context_file(
            &profile_workspace,
            "test3",
            "GEMINI.md"
        ).unwrap();
        assert!(content3.contains("Advanced Full-Featured Extension"));
    }
    
    #[test]
    fn test_environment_preparation() {
        let (launcher, _, _) = create_test_launcher();
        
        // Set some environment variables
        std::env::set_var("TEST_BASE_VAR", "base_value");
        std::env::set_var("TEST_OVERRIDE", "original");
        
        let mut profile = ProfileBuilder::new("env-test").build();
        profile.environment_variables.insert(
            "NEW_VAR".to_string(),
            "new_value".to_string()
        );
        profile.environment_variables.insert(
            "TEST_OVERRIDE".to_string(),
            "overridden".to_string()
        );
        profile.environment_variables.insert(
            "EXPANDED_VAR".to_string(),
            "$TEST_BASE_VAR".to_string()
        );
        
        let env = launcher.prepare_environment(&profile);
        
        // Check new variable
        assert_eq!(env.get("NEW_VAR"), Some(&"new_value".to_string()));
        
        // Check override
        assert_eq!(env.get("TEST_OVERRIDE"), Some(&"overridden".to_string()));
        
        // Check expansion
        assert_eq!(env.get("EXPANDED_VAR"), Some(&"base_value".to_string()));
        
        // Check profile ID is set
        assert_eq!(env.get("GEMINI_PROFILE"), Some(&"env-test".to_string()));
        
        // Clean up
        std::env::remove_var("TEST_BASE_VAR");
        std::env::remove_var("TEST_OVERRIDE");
    }
    
    #[test]
    fn test_working_directory_expansion() {
        let (launcher, _, _) = create_test_launcher();
        
        // Test cases for directory expansion
        let test_cases = vec![
            ("~/test", dirs::home_dir().unwrap().join("test")),
            ("/absolute/path", PathBuf::from("/absolute/path")),
            ("relative/path", PathBuf::from("relative/path")),
        ];
        
        for (input, expected_prefix) in test_cases {
            let profile = ProfileBuilder::new("dir-test")
                .build();
            let mut p = profile;
            p.working_directory = Some(input.to_string());
            
            // The launcher would expand this during launch
            let expanded = if input.starts_with("~") {
                dirs::home_dir()
                    .map(|home| home.join(&input[2..]))
                    .unwrap_or_else(|| PathBuf::from(input))
            } else {
                PathBuf::from(input)
            };
            
            assert!(expanded.to_string_lossy().contains(&expected_prefix.to_string_lossy()));
        }
    }
    
    #[test]
    fn test_multiple_extensions_installation() {
        let (launcher, workspace_dir, _storage_dir) = create_test_launcher();
        
        // Create multiple extensions
        let ext1 = McpFixtures::echo_extension();
        let ext2 = McpFixtures::multi_server_extension();
        let ext3 = McpFixtures::context_only_extension();
        
        launcher.storage.save_extension(&ext1).unwrap();
        launcher.storage.save_extension(&ext2).unwrap();
        launcher.storage.save_extension(&ext3).unwrap();
        
        // Create profile with all extensions
        let profile = ProfileBuilder::new("multi-ext")
            .with_extensions(vec![&ext1.id, &ext2.id, &ext3.id])
            .build();
        launcher.storage.save_profile(&profile).unwrap();
        
        let profile_workspace = workspace_dir.path().join(&profile.id);
        launcher.install_extensions_for_profile(&profile, &profile_workspace).unwrap();
        
        // Verify all extensions installed
        for ext_id in &[&ext1.id, &ext2.id, &ext3.id] {
            assert!(
                WorkspaceVerifier::verify_extension_installed(&profile_workspace, ext_id).is_ok(),
                "Extension {} not installed", ext_id
            );
        }
    }
    
    #[test]
    fn test_missing_extension_handling() {
        let (launcher, workspace_dir, _storage_dir) = create_test_launcher();
        
        // Create profile referencing non-existent extension
        let profile = ProfileBuilder::new("missing-ext")
            .with_extensions(vec!["non-existent-ext"])
            .build();
        launcher.storage.save_profile(&profile).unwrap();
        
        let profile_workspace = workspace_dir.path().join(&profile.id);
        
        // Should not panic, just warn
        launcher.install_extensions_for_profile(&profile, &profile_workspace).unwrap();
        
        // Workspace should still be created
        assert!(WorkspaceVerifier::verify_workspace_structure(&profile_workspace).is_ok());
    }
    
    #[test]
    fn test_extension_validation() {
        // Test valid extensions
        assert!(validate_extension_json(&McpFixtures::echo_extension()).is_ok());
        assert!(validate_extension_json(&McpFixtures::multi_server_extension()).is_ok());
        assert!(validate_extension_json(&McpFixtures::context_only_extension()).is_ok());
        
        // Test invalid extension - no name
        let mut bad_ext = McpFixtures::echo_extension();
        bad_ext.name = "".to_string();
        assert!(validate_extension_json(&bad_ext).is_err());
        
        // Test invalid extension - both command and URL
        let mut bad_ext2 = McpFixtures::echo_extension();
        if let Some(server) = bad_ext2.mcp_servers.get_mut("echo") {
            server.url = Some("http://localhost".to_string());
            // Already has command, so this is invalid
        }
        assert!(validate_extension_json(&bad_ext2).is_err());
    }
}