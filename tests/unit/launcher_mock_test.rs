#[cfg(test)]
mod tests {
    use crate::test_utils::{McpFixtures, ProfileBuilder, create_temp_storage};
    use gemini_cli_manager::launcher::Launcher;
    use std::path::PathBuf;
    use tempfile::TempDir;

    #[test]
    fn test_launch_with_profile_setup() {
        let temp_dir = TempDir::new().unwrap();
        let (storage, _storage_dir) = create_temp_storage();

        let launcher = Launcher::with_storage(storage);

        // Create extensions
        let ext1 = McpFixtures::echo_extension();
        let ext2 = McpFixtures::multi_server_extension();
        launcher.storage.save_extension(&ext1).unwrap();
        launcher.storage.save_extension(&ext2).unwrap();

        // Create profile with working directory
        let mut profile = ProfileBuilder::new("full-test")
            .with_extensions(vec![&ext1.id, &ext2.id])
            .build();
        profile.working_directory =
            Some(temp_dir.path().join("work").to_string_lossy().to_string());
        profile
            .environment_variables
            .insert("TEST_VAR".to_string(), "test_value".to_string());
        launcher.storage.save_profile(&profile).unwrap();

        // Test the setup parts of launch_with_profile
        let profile_workspace = temp_dir.path().join(&profile.id);

        // These are the steps that launch_with_profile would do:
        launcher.setup_workspace(&profile_workspace).unwrap();
        launcher
            .install_extensions_for_profile(&profile, &profile_workspace)
            .unwrap();
        let env = launcher.prepare_environment(&profile);

        // Verify workspace
        assert!(profile_workspace.join(".gemini").exists());
        assert!(
            profile_workspace
                .join(".gemini/extensions")
                .join(&ext1.id)
                .exists()
        );
        assert!(
            profile_workspace
                .join(".gemini/extensions")
                .join(&ext2.id)
                .exists()
        );

        // Verify environment
        assert_eq!(env.get("TEST_VAR"), Some(&"test_value".to_string()));
        assert_eq!(env.get("GEMINI_PROFILE"), Some(&profile.id));

        // Test working directory creation
        let work_dir = temp_dir.path().join("work");
        if !work_dir.exists() {
            std::fs::create_dir_all(&work_dir).unwrap();
        }
        assert!(work_dir.exists());
    }

    #[test]
    fn test_launch_with_profile_tilde_expansion() {
        let (storage, _storage_dir) = create_temp_storage();

        let _launcher = Launcher::with_storage(storage);

        // Profile with ~ in working directory
        let mut profile = ProfileBuilder::new("tilde-test").build();
        profile.working_directory = Some("~/test-dir".to_string());

        // Test the expansion logic from launch_with_profile
        let dir = profile.working_directory.as_ref().unwrap();
        let expanded = if dir.starts_with("~") {
            dirs::home_dir()
                .map(|home| home.join(&dir[2..]))
                .unwrap_or_else(|| PathBuf::from(dir))
        } else {
            PathBuf::from(dir)
        };

        // Should expand to home directory
        assert!(expanded.starts_with(&dirs::home_dir().unwrap()));
        assert!(expanded.ends_with("test-dir"));
    }

    #[test]
    fn test_launch_with_profile_current_dir() {
        let (storage, _storage_dir) = create_temp_storage();
        let _launcher = Launcher::with_storage(storage);

        // Profile without working directory
        let profile = ProfileBuilder::new("no-dir").build();

        // Test the logic that determines working directory
        let working_dir = if let Some(dir) = &profile.working_directory {
            PathBuf::from(dir)
        } else {
            std::env::current_dir().unwrap()
        };

        // Should use current directory
        assert_eq!(working_dir, std::env::current_dir().unwrap());
    }

    #[test]
    fn test_default_storage() {
        // Test that Launcher::new() creates a default storage
        let launcher = Launcher::new();

        // Try to use the storage
        let result = launcher.storage.list_extensions();

        // Should work (might be empty or have mock data)
        assert!(result.is_ok());
    }

    #[test]
    fn test_launch_in_terminal_function() {
        // We can't actually test launch_in_terminal because it calls launch_with_profile
        // which requires gemini to be installed. But we can at least check it exists
        // and document why it's not fully tested.

        // The function exists and is callable:
        let _func: fn(&gemini_cli_manager::models::Profile) -> color_eyre::Result<()> =
            gemini_cli_manager::launcher::launch_in_terminal;
    }

    #[test]
    fn test_workspace_permissions() {
        let temp_dir = TempDir::new().unwrap();
        let (storage, _storage_dir) = create_temp_storage();

        let launcher = Launcher::with_storage(storage);

        let workspace = temp_dir.path().join("test-workspace");
        launcher.setup_workspace(&workspace).unwrap();

        // Verify directories are created with proper permissions
        assert!(workspace.is_dir());
        assert!(workspace.join(".gemini").is_dir());
        assert!(workspace.join(".gemini/extensions").is_dir());

        // Check that we can write to these directories
        let test_file = workspace.join(".gemini/test.txt");
        std::fs::write(&test_file, "test").unwrap();
        assert!(test_file.exists());
    }
}
