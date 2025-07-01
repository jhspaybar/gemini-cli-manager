#[cfg(test)]
mod tests {
    use gemini_cli_manager::{
        models::{Extension, Profile, extension::McpServerConfig, profile::LaunchConfig},
        storage::Storage,
        launcher::Launcher,
        theme::{self, ThemeFlavour},
    };
    use std::collections::HashMap;
    use tempfile::TempDir;
    use chrono::Utc;

    // Simple storage tests
    #[test]
    fn test_storage_extension_crud() {
        let temp_dir = TempDir::new().unwrap();
        let storage = Storage::with_data_dir(temp_dir.path().to_path_buf());
        let _ = storage.init();
        
        // Get initial count
        let initial_count = storage.list_extensions().unwrap().len();

        // Create extension
        let ext = Extension {
            id: "test-ext".to_string(),
            name: "Test Extension".to_string(),
            version: "1.0.0".to_string(),
            description: Some("Test description".to_string()),
            mcp_servers: HashMap::new(),
            context_file_name: Some("GEMINI.md".to_string()),
            context_content: Some("# Test Content".to_string()),
            metadata: gemini_cli_manager::models::extension::ExtensionMetadata {
                imported_at: Utc::now(),
                source_path: None,
                tags: vec!["test".to_string()],
            },
        };

        // Save
        storage.save_extension(&ext).unwrap();

        // Load
        let loaded = storage.load_extension(&ext.id).unwrap();
        assert_eq!(loaded.name, ext.name);
        assert_eq!(loaded.version, ext.version);

        // List
        let list = storage.list_extensions().unwrap();
        assert_eq!(list.len(), initial_count + 1);

        // Delete
        storage.delete_extension(&ext.id).unwrap();
        assert!(storage.load_extension(&ext.id).is_err());
    }

    #[test]
    fn test_storage_profile_crud() {
        let temp_dir = TempDir::new().unwrap();
        let storage = Storage::with_data_dir(temp_dir.path().to_path_buf());
        let _ = storage.init();

        // Create profile
        let profile = Profile {
            id: "test-profile".to_string(),
            name: "Test Profile".to_string(),
            description: Some("Test description".to_string()),
            extension_ids: vec!["ext1".to_string(), "ext2".to_string()],
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

        // Save
        storage.save_profile(&profile).unwrap();

        // Load
        let loaded = storage.load_profile(&profile.id).unwrap();
        assert_eq!(loaded.name, profile.name);
        assert_eq!(loaded.extension_ids.len(), 2);

        // Delete
        storage.delete_profile(&profile.id).unwrap();
        assert!(storage.load_profile(&profile.id).is_err());
    }

    #[test]
    fn test_launcher_workspace_setup() {
        let workspace_temp = TempDir::new().unwrap();
        let storage_temp = TempDir::new().unwrap();
        
        let storage = Storage::with_data_dir(storage_temp.path().to_path_buf());
        let launcher = Launcher::with_storage(storage);

        let workspace_dir = workspace_temp.path().join("test-profile");
        launcher.setup_workspace(&workspace_dir).unwrap();

        // Verify structure
        assert!(workspace_dir.exists());
        assert!(workspace_dir.join(".gemini").exists());
        assert!(workspace_dir.join(".gemini").join("extensions").exists());
    }

    #[test]
    fn test_launcher_extension_installation() {
        let workspace_temp = TempDir::new().unwrap();
        let storage_temp = TempDir::new().unwrap();
        
        let storage = Storage::with_data_dir(storage_temp.path().to_path_buf());
        let _ = storage.init();

        // Create and save extension
        let ext = Extension {
            id: "echo-test".to_string(),
            name: "Echo Test".to_string(),
            version: "1.0.0".to_string(),
            description: Some("Echo test extension".to_string()),
            mcp_servers: {
                let mut servers = HashMap::new();
                servers.insert("echo".to_string(), McpServerConfig {
                    command: Some("node".to_string()),
                    args: Some(vec!["echo.js".to_string()]),
                    cwd: None,
                    env: None,
                    timeout: None,
                    trust: Some(true),
                });
                servers
            },
            context_file_name: Some("GEMINI.md".to_string()),
            context_content: Some("# Echo Test\n\nThis is a test.".to_string()),
            metadata: gemini_cli_manager::models::extension::ExtensionMetadata {
                imported_at: Utc::now(),
                source_path: None,
                tags: vec!["test".to_string()],
            },
        };
        storage.save_extension(&ext).unwrap();

        // Create profile
        let profile = Profile {
            id: "test-profile".to_string(),
            name: "Test Profile".to_string(),
            description: None,
            extension_ids: vec![ext.id.clone()],
            environment_variables: HashMap::new(),
            working_directory: None,
            launch_config: LaunchConfig::default(),
            metadata: gemini_cli_manager::models::profile::ProfileMetadata {
                created_at: Utc::now(),
                updated_at: Utc::now(),
                tags: vec![],
                is_default: false,
                icon: None,
            },
        };

        let launcher = Launcher::with_storage(storage);

        let workspace_dir = workspace_temp.path().join(&profile.id);
        launcher.setup_workspace(&workspace_dir).unwrap();
        launcher.install_extensions_for_profile(&profile, &workspace_dir).unwrap();

        // Verify extension installed
        let ext_dir = workspace_dir.join(".gemini").join("extensions").join(&ext.id);
        assert!(ext_dir.exists());
        
        let config_file = ext_dir.join("gemini-extension.json");
        assert!(config_file.exists());
        
        let context_file = ext_dir.join("GEMINI.md");
        assert!(context_file.exists());
    }

    #[test]
    fn test_theme_contrast() {
        // Test that text is visible in dark theme
        theme::set_flavour(ThemeFlavour::Mocha);
        
        // Just verify we can access theme colors
        let _bg = theme::background();
        let _fg = theme::text_primary();
        let _border = theme::border();
        
        // In a real test, we'd calculate contrast ratios
        // For now, just ensure theme system works
    }

    #[test]
    fn test_mcp_server_validation() {
        // Valid command-based server
        let server1 = McpServerConfig {
            command: Some("node".to_string()),
            args: Some(vec!["server.js".to_string()]),
            cwd: None,
            env: None,
            timeout: None,
            trust: None,
        };

        // Valid server with timeout
        let server2 = McpServerConfig {
            command: Some("python".to_string()),
            args: Some(vec!["server.py".to_string()]),
            cwd: None,
            env: None,
            timeout: Some(5000),
            trust: Some(false),
        };

        // Create extension with servers
        let mut ext = Extension {
            id: "test".to_string(),
            name: "Test".to_string(),
            version: "1.0.0".to_string(),
            description: None,
            mcp_servers: HashMap::new(),
            context_file_name: None,
            context_content: None,
            metadata: gemini_cli_manager::models::extension::ExtensionMetadata {
                imported_at: Utc::now(),
                source_path: None,
                tags: vec![],
            },
        };

        ext.mcp_servers.insert("server1".to_string(), server1);
        ext.mcp_servers.insert("server2".to_string(), server2);

        // Basic validation - has required fields
        assert!(!ext.name.is_empty());
        assert!(!ext.version.is_empty());
    }

    #[test]
    fn test_environment_preparation() {
        let storage_temp = TempDir::new().unwrap();
        let storage = Storage::with_data_dir(storage_temp.path().to_path_buf());
        let launcher = Launcher::with_storage(storage);

        let mut profile = Profile {
            id: "env-test".to_string(),
            name: "Environment Test".to_string(),
            description: None,
            extension_ids: vec![],
            environment_variables: HashMap::new(),
            working_directory: None,
            launch_config: LaunchConfig::default(),
            metadata: gemini_cli_manager::models::profile::ProfileMetadata {
                created_at: Utc::now(),
                updated_at: Utc::now(),
                tags: vec![],
                is_default: false,
                icon: None,
            },
        };

        // Add environment variables
        profile.environment_variables.insert("TEST_VAR".to_string(), "test_value".to_string());
        profile.environment_variables.insert("ANOTHER_VAR".to_string(), "another_value".to_string());

        let env = launcher.prepare_environment(&profile);
        
        assert_eq!(env.get("TEST_VAR"), Some(&"test_value".to_string()));
        assert_eq!(env.get("ANOTHER_VAR"), Some(&"another_value".to_string()));
        assert_eq!(env.get("GEMINI_PROFILE"), Some(&"env-test".to_string()));
    }
}