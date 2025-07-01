#[cfg(test)]
mod tests {
    use gemini_cli_manager::{
        launcher::Launcher,
    };
    use crate::test_utils::{
        create_temp_storage, McpFixtures, ExtensionBuilder, ProfileBuilder,
        WorkspaceVerifier,
    };
    use tempfile::TempDir;
    
    #[test]
    fn test_complete_user_journey_echo_server() {
        // Setup
        let workspace_temp = TempDir::new().unwrap();
        let (storage, _storage_temp) = create_temp_storage();
        
        // Step 1: User creates an extension with echo MCP server
        println!("Step 1: Creating extension with echo server...");
        let echo_ext = McpFixtures::echo_extension();
        storage.save_extension(&echo_ext).unwrap();
        
        // Verify extension saved correctly
        let loaded_ext = storage.load_extension(&echo_ext.id).unwrap();
        assert_eq!(loaded_ext.name, "Echo Test Extension");
        assert!(loaded_ext.mcp_servers.contains_key("echo"));
        
        // Step 2: User creates a profile using that extension
        println!("Step 2: Creating profile with echo extension...");
        let profile = ProfileBuilder::new("echo-test-profile")
            .with_description("Profile for testing echo server")
            .with_extensions(vec![&echo_ext.id])
            .with_tags(vec!["test", "echo"])
            .build();
        storage.save_profile(&profile).unwrap();
        
        // Verify profile saved correctly
        let loaded_profile = storage.load_profile(&profile.id).unwrap();
        assert_eq!(loaded_profile.extension_ids.len(), 1);
        assert_eq!(loaded_profile.extension_ids[0], echo_ext.id);
        
        // Step 3: User launches the profile
        println!("Step 3: Setting up workspace for launch...");
        let launcher = Launcher::with_storage(storage);
            
        let profile_workspace = workspace_temp.path().join(&profile.id);
        
        // Setup workspace
        launcher.setup_workspace(&profile_workspace).unwrap();
        
        // Install extensions
        launcher.install_extensions_for_profile(&profile, &profile_workspace).unwrap();
        
        // Step 4: Verify workspace is set up correctly
        println!("Step 4: Verifying workspace setup...");
        
        // Check workspace structure
        assert!(WorkspaceVerifier::verify_workspace_structure(&profile_workspace).is_ok());
        
        // Check extension installed
        assert!(WorkspaceVerifier::verify_extension_installed(
            &profile_workspace,
            &echo_ext.id
        ).is_ok());
        
        // Check gemini-extension.json
        let config_path = profile_workspace
            .join(".gemini")
            .join("extensions")
            .join(&echo_ext.id)
            .join("gemini-extension.json");
        assert!(config_path.exists());
        
        let config_content = std::fs::read_to_string(&config_path).unwrap();
        let config_json: serde_json::Value = serde_json::from_str(&config_content).unwrap();
        
        assert_eq!(config_json["name"], "Echo Test Extension");
        assert_eq!(config_json["version"], "1.0.0");
        assert!(config_json["mcpServers"]["echo"]["command"].as_str().unwrap() == "node");
        
        // Check context file
        let context_content = WorkspaceVerifier::verify_context_file(
            &profile_workspace,
            &echo_ext.id,
            "GEMINI.md"
        ).unwrap();
        assert!(context_content.contains("Echo Test Extension"));
        assert!(context_content.contains("Available Commands"));
        
        // Step 5: Verify environment setup
        println!("Step 5: Verifying environment configuration...");
        let env = launcher.prepare_environment(&profile);
        
        assert_eq!(env.get("GEMINI_PROFILE"), Some(&profile.id));
        
        println!("✅ Complete user journey test passed!");
    }
    
    #[test]
    fn test_multi_extension_profile_journey() {
        let workspace_temp = TempDir::new().unwrap();
        let (storage, _storage_temp) = create_temp_storage();
        
        // Create multiple extensions
        let echo_ext = McpFixtures::echo_extension();
        let multi_ext = McpFixtures::multi_server_extension();
        let context_ext = McpFixtures::context_only_extension();
        
        storage.save_extension(&echo_ext).unwrap();
        storage.save_extension(&multi_ext).unwrap();
        storage.save_extension(&context_ext).unwrap();
        
        // Create profile with all extensions
        let mut profile = ProfileBuilder::new("full-stack-dev")
            .with_description("Full development environment")
            .with_extensions(vec![&echo_ext.id, &multi_ext.id, &context_ext.id])
            .with_tags(vec!["development", "full-stack"])
            .as_default()
            .build();
            
        // Add environment variables
        profile.environment_variables.insert("NODE_ENV".to_string(), "development".to_string());
        profile.environment_variables.insert("API_URL".to_string(), "http://localhost:3000".to_string());
        profile.working_directory = Some("~/projects".to_string());
        
        storage.save_profile(&profile).unwrap();
        
        // Launch preparation
        let launcher = Launcher::with_storage(storage);
            
        let profile_workspace = workspace_temp.path().join(&profile.id);
        launcher.setup_workspace(&profile_workspace).unwrap();
        launcher.install_extensions_for_profile(&profile, &profile_workspace).unwrap();
        
        // Verify all extensions installed
        for ext_id in &[&echo_ext.id, &multi_ext.id, &context_ext.id] {
            assert!(
                WorkspaceVerifier::verify_extension_installed(&profile_workspace, ext_id).is_ok(),
                "Extension {} not properly installed", ext_id
            );
        }
        
        // Verify multi-server extension
        let multi_config_path = profile_workspace
            .join(".gemini")
            .join("extensions")
            .join(&multi_ext.id)
            .join("gemini-extension.json");
            
        let multi_config = std::fs::read_to_string(&multi_config_path).unwrap();
        let multi_json: serde_json::Value = serde_json::from_str(&multi_config).unwrap();
        
        assert!(multi_json["mcpServers"]["echo"].is_object());
        assert!(multi_json["mcpServers"]["python-echo"].is_object());
        assert!(multi_json["mcpServers"]["api-server"].is_object());
        
        // Verify context files with different names
        assert!(WorkspaceVerifier::verify_context_file(
            &profile_workspace,
            &multi_ext.id,
            "MULTI_SERVER.md"
        ).is_ok());
        
        assert!(WorkspaceVerifier::verify_context_file(
            &profile_workspace,
            &context_ext.id,
            "INSTRUCTIONS.md"
        ).is_ok());
        
        // Verify environment
        let env = launcher.prepare_environment(&profile);
        assert_eq!(env.get("NODE_ENV"), Some(&"development".to_string()));
        assert_eq!(env.get("API_URL"), Some(&"http://localhost:3000".to_string()));
        assert_eq!(env.get("GEMINI_PROFILE"), Some(&"full-stack-dev".to_string()));
    }
    
    #[test]
    fn test_profile_update_journey() {
        let workspace_temp = TempDir::new().unwrap();
        let (storage, _storage_temp) = create_temp_storage();
        
        // Initial setup
        let ext1 = McpFixtures::echo_extension();
        let ext2 = McpFixtures::context_only_extension();
        storage.save_extension(&ext1).unwrap();
        storage.save_extension(&ext2).unwrap();
        
        // Create initial profile
        let mut profile = ProfileBuilder::new("evolving-profile")
            .with_extensions(vec![&ext1.id])
            .build();
        storage.save_profile(&profile).unwrap();
        
        // Launch with initial configuration
        let launcher = Launcher::with_storage(storage.clone());
            
        let profile_workspace = workspace_temp.path().join(&profile.id);
        launcher.setup_workspace(&profile_workspace).unwrap();
        launcher.install_extensions_for_profile(&profile, &profile_workspace).unwrap();
        
        // Verify initial state
        assert!(WorkspaceVerifier::verify_extension_installed(
            &profile_workspace,
            &ext1.id
        ).is_ok());
        
        // User updates profile to add another extension
        profile.extension_ids.push(ext2.id.clone());
        profile.environment_variables.insert("NEW_VAR".to_string(), "new_value".to_string());
        storage.save_profile(&profile).unwrap();
        
        // Re-install with updated profile
        launcher.install_extensions_for_profile(&profile, &profile_workspace).unwrap();
        
        // Verify both extensions now installed
        assert!(WorkspaceVerifier::verify_extension_installed(
            &profile_workspace,
            &ext2.id
        ).is_ok());
        
        // Verify new environment variable
        let env = launcher.prepare_environment(&profile);
        assert_eq!(env.get("NEW_VAR"), Some(&"new_value".to_string()));
    }
    
    #[test]
    fn test_error_recovery_journey() {
        let workspace_temp = TempDir::new().unwrap();
        let (storage, _storage_temp) = create_temp_storage();
        
        // Create profile referencing non-existent extension
        let profile = ProfileBuilder::new("error-recovery-test")
            .with_extensions(vec!["non-existent-ext", "another-missing-ext"])
            .build();
        storage.save_profile(&profile).unwrap();
        
        // Also add one valid extension
        let valid_ext = McpFixtures::echo_extension();
        storage.save_extension(&valid_ext).unwrap();
        
        // Update profile to include valid extension
        let mut updated_profile = profile.clone();
        updated_profile.extension_ids.push(valid_ext.id.clone());
        storage.save_profile(&updated_profile).unwrap();
        
        // Launch should handle missing extensions gracefully
        let launcher = Launcher::with_storage(storage);
            
        let profile_workspace = workspace_temp.path().join(&updated_profile.id);
        launcher.setup_workspace(&profile_workspace).unwrap();
        
        // Should not panic when installing
        launcher.install_extensions_for_profile(&updated_profile, &profile_workspace).unwrap();
        
        // Valid extension should still be installed
        assert!(WorkspaceVerifier::verify_extension_installed(
            &profile_workspace,
            &valid_ext.id
        ).is_ok());
        
        // Workspace should be created even with errors
        assert!(WorkspaceVerifier::verify_workspace_structure(&profile_workspace).is_ok());
    }
    
    #[test]
    fn test_context_file_variations_journey() {
        let workspace_temp = TempDir::new().unwrap();
        let (storage, _storage_temp) = create_temp_storage();
        
        // Create extensions with different context file configurations
        let mut ext1 = McpFixtures::echo_extension();
        ext1.id = "ext-with-gemini-md".to_string();
        ext1.context_file_name = Some("GEMINI.md".to_string());
        
        let mut ext2 = McpFixtures::context_only_extension();
        ext2.id = "ext-with-custom-md".to_string();
        ext2.context_file_name = Some("PROJECT_CONTEXT.md".to_string());
        
        let mut ext3 = McpFixtures::full_featured_extension();
        ext3.id = "ext-with-default-md".to_string();
        ext3.context_file_name = None; // Should default to GEMINI.md
        
        let mut ext4 = ExtensionBuilder::new("ext-no-context").build();
        ext4.context_content = None; // No context file at all
        
        storage.save_extension(&ext1).unwrap();
        storage.save_extension(&ext2).unwrap();
        storage.save_extension(&ext3).unwrap();
        storage.save_extension(&ext4).unwrap();
        
        // Create profile with all extensions
        let profile = ProfileBuilder::new("context-test")
            .with_extensions(vec![&ext1.id, &ext2.id, &ext3.id, &ext4.id])
            .build();
        storage.save_profile(&profile).unwrap();
        
        // Install
        let launcher = Launcher::with_storage(storage);
            
        let profile_workspace = workspace_temp.path().join(&profile.id);
        launcher.setup_workspace(&profile_workspace).unwrap();
        launcher.install_extensions_for_profile(&profile, &profile_workspace).unwrap();
        
        // Verify each context file situation
        assert!(WorkspaceVerifier::verify_context_file(
            &profile_workspace,
            &ext1.id,
            "GEMINI.md"
        ).is_ok());
        
        assert!(WorkspaceVerifier::verify_context_file(
            &profile_workspace,
            &ext2.id,
            "PROJECT_CONTEXT.md"
        ).is_ok());
        
        assert!(WorkspaceVerifier::verify_context_file(
            &profile_workspace,
            &ext3.id,
            "GEMINI.md"
        ).is_ok());
        
        // Extension 4 should not have a context file
        let ext4_context = WorkspaceVerifier::verify_context_file(
            &profile_workspace,
            &ext4.id,
            "GEMINI.md"
        );
        assert!(ext4_context.is_err());
    }
}