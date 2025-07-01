#[cfg(test)]
mod tests {
    use gemini_cli_manager::cli::{Cli, version};
    use clap::Parser;

    #[test]
    fn test_cli_default_values() {
        let cli = Cli::parse_from(&["gemini-cli-manager"]);
        
        assert_eq!(cli.list_storage, false);
    }

    #[test]
    fn test_cli_list_storage_flag() {
        let cli = Cli::parse_from(&["gemini-cli-manager", "--list-storage"]);
        
        assert_eq!(cli.list_storage, true);
    }

    #[test]
    fn test_version_function() {
        let version_str = version();
        
        // Version should contain certain elements
        assert!(version_str.contains("Authors:"));
        assert!(version_str.contains("Config directory:"));
        assert!(version_str.contains("Data directory:"));
        
        // Should not be empty
        assert!(!version_str.is_empty());
    }

    #[test]
    fn test_cli_help() {
        // Test that help doesn't panic
        let result = Cli::try_parse_from(&["gemini-cli-manager", "--help"]);
        
        // Help should return an error (clap's way of handling help)
        assert!(result.is_err());
        
        // But the error should be a help error
        let err = result.unwrap_err();
        assert_eq!(err.kind(), clap::error::ErrorKind::DisplayHelp);
    }

    #[test]
    fn test_cli_version() {
        // Test that version doesn't panic
        let result = Cli::try_parse_from(&["gemini-cli-manager", "--version"]);
        
        // Version should return an error (clap's way of handling version)
        assert!(result.is_err());
        
        // But the error should be a version error
        let err = result.unwrap_err();
        assert_eq!(err.kind(), clap::error::ErrorKind::DisplayVersion);
    }
}