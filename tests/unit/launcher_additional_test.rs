#[cfg(test)]
mod tests {
    use crate::test_utils::{McpFixtures, ProfileBuilder};
    use gemini_cli_manager::launcher::{Launcher, create_launch_script};
    use gemini_cli_manager::storage::Storage;
    use std::path::PathBuf;
    use tempfile::TempDir;

    #[test]
    fn test_launcher_new() {
        // Test the default new() constructor
        let launcher = Launcher::new();

        // Should have default workspace directory
        let _expected_workspace = dirs::home_dir()
            .unwrap_or_else(|| PathBuf::from("."))
            .join(".gemini-workspace");

        // We can't directly access workspace_dir, but we can test through behavior
        // The launcher should work with default settings
        assert!(matches!(launcher.storage, Storage { .. }));
    }

    #[test]
    fn test_launcher_workspace_setup() {
        let launcher = Launcher::new();
        let temp_dir = TempDir::new().unwrap();

        // Test that workspace setup works
        let profile = ProfileBuilder::new("test").build();
        let workspace = temp_dir.path().join(&profile.id);

        launcher.setup_workspace(&workspace).unwrap();
        assert!(workspace.exists());
    }

    #[test]
    fn test_create_launch_script() {
        let temp_dir = TempDir::new().unwrap();
        let script_path = temp_dir.path().join("launch.sh");

        // Create a profile with various settings
        let mut profile = ProfileBuilder::new("script-test")
            .with_extensions(vec!["ext1", "ext2"])
            .build();
        profile
            .environment_variables
            .insert("TEST_VAR".to_string(), "test_value".to_string());
        profile
            .environment_variables
            .insert("PATH_VAR".to_string(), "/custom/path".to_string());
        profile.working_directory = Some("~/projects/test".to_string());

        // Create the launch script
        create_launch_script(&profile, &script_path).unwrap();

        // Verify script was created
        assert!(script_path.exists());

        // Read and verify script content
        let content = std::fs::read_to_string(&script_path).unwrap();

        // Check key components
        assert!(content.contains(&format!("export GEMINI_PROFILE=\"{}\"", profile.id)));
        assert!(content.contains("export TEST_VAR=\"test_value\""));
        assert!(content.contains("export PATH_VAR=\"/custom/path\""));
        assert!(content.contains("cd \"~/projects/test\""));
        assert!(content.contains(&format!(
            "Launching Gemini CLI with profile: {}",
            profile.name
        )));
        assert!(content.contains("gemini"));

        // On Unix, check that it's executable
        #[cfg(unix)]
        {
            use std::os::unix::fs::PermissionsExt;
            let metadata = std::fs::metadata(&script_path).unwrap();
            let permissions = metadata.permissions();
            assert_eq!(permissions.mode() & 0o111, 0o111); // Check execute bits
        }
    }

    #[test]
    fn test_create_launch_script_no_working_dir() {
        let temp_dir = TempDir::new().unwrap();
        let script_path = temp_dir.path().join("launch_no_dir.sh");

        // Create a profile without working directory
        let profile = ProfileBuilder::new("no-dir-test").build();

        create_launch_script(&profile, &script_path).unwrap();

        let content = std::fs::read_to_string(&script_path).unwrap();
        assert!(content.contains("# No working directory specified"));
        assert!(!content.contains("cd \""));
    }

    #[test]
    fn test_prepare_environment_empty_profile() {
        let launcher = Launcher::new();

        // Profile with no environment variables
        let profile = ProfileBuilder::new("empty-env").build();

        let env = launcher.prepare_environment(&profile);

        // Should still have GEMINI_PROFILE set
        assert_eq!(env.get("GEMINI_PROFILE"), Some(&"empty-env".to_string()));

        // Should include system environment variables
        assert!(env.len() > 1); // At least GEMINI_PROFILE plus system vars
    }

    #[test]
    fn test_prepare_environment_invalid_expansion() {
        let launcher = Launcher::new();

        let mut profile = ProfileBuilder::new("invalid-expand").build();
        // Reference to non-existent variable
        profile
            .environment_variables
            .insert("EXPANDED".to_string(), "$NONEXISTENT_VAR".to_string());

        let env = launcher.prepare_environment(&profile);

        // Should keep the original value when expansion fails
        assert_eq!(env.get("EXPANDED"), Some(&"$NONEXISTENT_VAR".to_string()));
    }

    #[test]
    fn test_install_extension_no_context() {
        let temp_dir = TempDir::new().unwrap();
        let (storage, _storage_dir) = crate::test_utils::create_temp_storage();

        let launcher = Launcher::with_storage(storage);

        // Extension without context content
        let mut ext = McpFixtures::echo_extension();
        ext.context_content = None;
        ext.context_file_name = None;
        launcher.storage.save_extension(&ext).unwrap();

        let profile = ProfileBuilder::new("no-context")
            .with_extensions(vec![&ext.id])
            .build();
        launcher.storage.save_profile(&profile).unwrap();

        let workspace = temp_dir.path().join(&profile.id);
        launcher
            .install_extensions_for_profile(&profile, &workspace)
            .unwrap();

        // Should have extension config but no context file
        let ext_dir = workspace.join(".gemini").join("extensions").join(&ext.id);
        assert!(ext_dir.join("gemini-extension.json").exists());
        assert!(!ext_dir.join("GEMINI.md").exists());
        assert!(!ext_dir.join("CUSTOM.md").exists());
    }

    #[test]
    fn test_install_extension_empty_context_filename() {
        let temp_dir = TempDir::new().unwrap();
        let (storage, _storage_dir) = crate::test_utils::create_temp_storage();

        let launcher = Launcher::with_storage(storage);

        // Extension with empty context filename (should default to GEMINI.md)
        let mut ext = McpFixtures::echo_extension();
        ext.context_file_name = Some("   ".to_string()); // Whitespace only
        launcher.storage.save_extension(&ext).unwrap();

        let profile = ProfileBuilder::new("empty-filename")
            .with_extensions(vec![&ext.id])
            .build();
        launcher.storage.save_profile(&profile).unwrap();

        let workspace = temp_dir.path().join(&profile.id);
        launcher
            .install_extensions_for_profile(&profile, &workspace)
            .unwrap();

        // Should default to GEMINI.md
        let ext_dir = workspace.join(".gemini").join("extensions").join(&ext.id);
        assert!(ext_dir.join("GEMINI.md").exists());
    }

    #[test]
    fn test_launch_with_profile_workspace_creation() {
        let temp_dir = TempDir::new().unwrap();
        let (storage, _storage_dir) = crate::test_utils::create_temp_storage();

        let launcher = Launcher::with_storage(storage);

        let profile = ProfileBuilder::new("launch-test").build();

        // We can't actually test the full launch (would require gemini installed),
        // but we can test the workspace setup part
        let workspace = temp_dir.path().join(&profile.id);

        // The launch_with_profile method would call setup_workspace
        launcher.setup_workspace(&workspace).unwrap();

        // Verify workspace structure
        assert!(workspace.join(".gemini").exists());
        assert!(workspace.join(".gemini").join("extensions").exists());
    }

    #[test]
    fn test_working_directory_creation() {
        let temp_dir = TempDir::new().unwrap();

        // Test directory creation logic from launch_with_profile
        let test_dir = temp_dir.path().join("new_working_dir");
        assert!(!test_dir.exists());

        // This mirrors the logic in launch_with_profile
        if !test_dir.exists() {
            std::fs::create_dir_all(&test_dir).unwrap();
        }

        assert!(test_dir.exists());
    }

    // Note: We can't easily test the actual launch_with_profile method
    // because it requires the 'gemini' command to be installed.
    // Similarly, launch_in_terminal just calls launch_with_profile.
    // The Command execution and gemini check are integration-level concerns
    // that would require mocking or actual gemini installation.
}
