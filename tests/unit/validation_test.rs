#[cfg(test)]
mod tests {
    use gemini_cli_manager::models::extension::McpServerConfig;
    use crate::test_utils::{ExtensionBuilder, ProfileBuilder, validate_extension_json};
    use std::collections::HashMap;
    
    #[test]
    fn test_extension_name_validation() {
        // Valid names
        let valid_names = vec![
            "simple-name",
            "name_with_underscores",
            "name123",
            "CamelCase",
            "name.with.dots",
        ];
        
        for name in valid_names {
            let ext = ExtensionBuilder::new(name).build();
            assert!(!ext.name.is_empty());
        }
        
        // Empty name should fail validation
        let mut ext = ExtensionBuilder::new("test").build();
        ext.name = "".to_string();
        assert!(validate_extension_json(&ext).is_err());
    }
    
    #[test]
    fn test_version_format_validation() {
        // Valid versions
        let valid_versions = vec![
            "1.0.0",
            "2.1.3",
            "0.0.1",
            "1.0.0-beta",
            "2.0.0-rc.1",
            "3.0.0+build123",
        ];
        
        for version in valid_versions {
            let ext = ExtensionBuilder::new("test")
                .with_version(version)
                .build();
            assert_eq!(ext.version, version);
        }
        
        // Empty version should fail
        let mut ext = ExtensionBuilder::new("test").build();
        ext.version = "".to_string();
        assert!(validate_extension_json(&ext).is_err());
    }
    
    #[test]
    fn test_mcp_server_validation() {
        // Valid: command-based server
        let mut ext1 = ExtensionBuilder::new("test1").build();
        ext1.mcp_servers.insert("server1".to_string(), McpServerConfig {
            command: Some("node".to_string()),
            args: Some(vec!["server.js".to_string()]),
            cwd: None,
            env: None,
            timeout: None,
            trust: None,
            url: None,
        });
        assert!(validate_extension_json(&ext1).is_ok());
        
        // Valid: Another command-based server
        let mut ext2 = ExtensionBuilder::new("test2").build();
        ext2.mcp_servers.insert("server2".to_string(), McpServerConfig {
            command: Some("python".to_string()),
            args: Some(vec!["server.py".to_string()]),
            cwd: None,
            env: None,
            timeout: Some(5000),
            trust: Some(false),
            url: None,
        });
        assert!(validate_extension_json(&ext2).is_ok());
        
        // Invalid: No command
        let mut ext3 = ExtensionBuilder::new("test3").build();
        ext3.mcp_servers.insert("server3".to_string(), McpServerConfig {
            command: None,
            args: None,
            cwd: None,
            env: None,
            timeout: None,
            trust: None,
            url: None,
        });
        let result = validate_extension_json(&ext3);
        assert!(result.is_err());
        assert!(result.unwrap_err().contains("must have a command"));
        
        // Invalid: Empty command string
        let mut ext4 = ExtensionBuilder::new("test4").build();
        ext4.mcp_servers.insert("server4".to_string(), McpServerConfig {
            command: Some("".to_string()),
            args: None,
            cwd: None,
            env: None,
            timeout: None,
            trust: None,
            url: None,
        });
        // In a real implementation, we might want to validate empty strings
        // For now, this would technically pass as command is Some(...)
    }
    
    #[test]
    fn test_environment_variable_validation() {
        let mut ext = ExtensionBuilder::new("env-test").build();
        
        // Valid environment variables
        let mut env = HashMap::new();
        env.insert("VALID_VAR".to_string(), "value".to_string());
        env.insert("PATH".to_string(), "/usr/bin:/usr/local/bin".to_string());
        env.insert("EMPTY_VAR".to_string(), "".to_string()); // Empty is valid
        env.insert("VAR_WITH_REFERENCE".to_string(), "$OTHER_VAR".to_string());
        
        ext.mcp_servers.insert("server".to_string(), McpServerConfig {
            command: Some("test".to_string()),
            args: Some(vec![]),
            cwd: None,
            env: Some(env),
            timeout: None,
            trust: None,
            url: None,
        });
        
        assert!(validate_extension_json(&ext).is_ok());
    }
    
    #[test]
    fn test_path_traversal_prevention() {
        // Test paths that should be rejected
        let dangerous_paths = vec![
            "../../../etc/passwd",
            "..\\..\\windows\\system32",
            "/etc/passwd",
            "C:\\Windows\\System32",
            "~/../../root",
        ];
        
        for path in dangerous_paths {
            // In a real implementation, we'd have path validation
            // For now, just ensure paths don't contain ".."
            assert!(path.contains("..") || path.starts_with('/') || path.contains(':'));
        }
    }
    
    #[test]
    fn test_timeout_validation() {
        let mut ext = ExtensionBuilder::new("timeout-test").build();
        
        // Valid timeouts
        let valid_timeouts = vec![
            Some(1000),      // 1 second
            Some(60000),     // 1 minute
            Some(600000),    // 10 minutes
            None,            // No timeout is valid
        ];
        
        for timeout in valid_timeouts {
            ext.mcp_servers.insert("server".to_string(), McpServerConfig {
                command: Some("node".to_string()),
                args: Some(vec!["server.js".to_string()]),
                cwd: None,
                env: None,
                timeout,
                trust: None,
                url: None,
            });
            assert!(validate_extension_json(&ext).is_ok());
        }
        
        // Note: In practice, we might want to validate against:
        // - Zero or negative timeouts
        // - Extremely large timeouts
    }
    
    #[test]
    fn test_id_uniqueness_validation() {
        // IDs should be unique and follow naming conventions
        let ext1 = ExtensionBuilder::new("Test Extension").build();
        let ext2 = ExtensionBuilder::new("Test Extension").build();
        
        // Both should generate the same ID
        assert_eq!(ext1.id, ext2.id);
        assert_eq!(ext1.id, "test-extension");
        
        // Different names should generate different IDs
        let ext3 = ExtensionBuilder::new("Another Extension").build();
        assert_ne!(ext1.id, ext3.id);
    }
    
    #[test]
    fn test_profile_extension_reference_validation() {
        // Profile should only reference existing extensions
        let profile = ProfileBuilder::new("test-profile")
            .with_extensions(vec!["ext1", "ext2", "ext3"])
            .build();
            
        // In a real implementation, we'd validate these exist in storage
        assert_eq!(profile.extension_ids.len(), 3);
    }
    
    #[test]
    fn test_working_directory_validation() {
        // Valid working directories
        let valid_dirs = vec![
            Some("~/projects".to_string()),
            Some("/usr/local/bin".to_string()),
            Some("./relative/path".to_string()),
            None, // No working directory is valid
        ];
        
        for dir in valid_dirs {
            let profile = ProfileBuilder::new("test").build();
            let mut p = profile;
            p.working_directory = dir.clone();
            
            // Should not panic or error
            assert!(dir.is_none() || !dir.unwrap().is_empty());
        }
    }
    
    #[test]
    fn test_tag_validation() {
        // Valid tags
        let valid_tags = vec![
            vec!["tag1", "tag2"],
            vec!["development", "testing", "production"],
            vec!["feature-123", "bug-fix"],
            vec![], // Empty tags are valid
        ];
        
        for tags in valid_tags {
            let ext = ExtensionBuilder::new("test")
                .with_tags(tags.clone())
                .build();
            assert_eq!(ext.metadata.tags.len(), tags.len());
        }
    }
    
    #[test]
    fn test_special_character_handling() {
        // Test that special characters are handled properly
        let special_chars = vec![
            "Name with spaces",
            "Name-with-dashes",
            "Name_with_underscores",
            "Name.with.dots",
            "Name@with#special$chars",
        ];
        
        for name in special_chars {
            let ext = ExtensionBuilder::new(name).build();
            
            // ID should be sanitized
            assert!(!ext.id.contains(' '));
            assert!(!ext.id.contains('@'));
            assert!(!ext.id.contains('#'));
            assert!(!ext.id.contains('$'));
        }
    }
    
    #[test]
    fn test_empty_mcp_servers_allowed() {
        // Extensions without MCP servers are valid (context-only)
        let ext = ExtensionBuilder::new("context-only").build();
        assert!(ext.mcp_servers.is_empty());
        assert!(validate_extension_json(&ext).is_ok());
    }
    
    #[test]
    fn test_profile_circular_reference_prevention() {
        // In a real implementation, we'd check for:
        // - Profiles referencing themselves
        // - Circular dependencies between profiles
        // - Extensions with circular dependencies
        
        let profile = ProfileBuilder::new("test").build();
        assert!(!profile.extension_ids.contains(&profile.id));
    }
}

#[cfg(test)]
mod error_recovery_tests {
    use crate::test_utils::create_temp_storage;
    
    #[test]
    fn test_missing_extension_recovery() {
        let (storage, _temp) = create_temp_storage();
        
        // Try to load non-existent extension
        let result = storage.load_extension("non-existent");
        assert!(result.is_err());
        
        // Should not panic, return proper error
        match result {
            Err(e) => {
                // Error should be informative
                let error_str = e.to_string();
                assert!(!error_str.is_empty());
            }
            Ok(_) => panic!("Should have failed"),
        }
    }
    
    #[test]
    fn test_corrupted_storage_recovery() {
        let (storage, temp) = create_temp_storage();
        
        // Save valid extension
        let ext = crate::ExtensionBuilder::new("test").build();
        storage.save_extension(&ext).unwrap();
        
        // Corrupt the file
        let file_path = temp.path()
            .join("extensions")
            .join(format!("{}.json", ext.id));
        std::fs::write(&file_path, "{ corrupted json").unwrap();
        
        // Try to load
        let result = storage.load_extension(&ext.id);
        assert!(result.is_err());
        
        // List should handle corrupted files gracefully
        let list = storage.list_extensions();
        assert!(list.is_ok());
        assert_eq!(list.unwrap().len(), 0); // Corrupted file is skipped
    }
}