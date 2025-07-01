#[cfg(test)]
mod tests {
    use gemini_cli_manager::{
        view::{ViewManager, ViewType},
        action::Action,
    };
    use crate::test_utils::*;
    use tokio::sync::mpsc;

    async fn create_test_view_manager() -> ViewManager {
        let storage = create_test_storage();
        
        // Add some test data
        let ext1 = ExtensionBuilder::new("Test Extension")
            .with_description("Test extension for view tests")
            .with_tags(vec!["test"])
            .build();
        storage.save_extension(&ext1).unwrap();
        
        // Add another extension not referenced by profiles (for delete tests)
        let ext2 = ExtensionBuilder::new("Deletable Extension")
            .with_description("Extension that can be deleted")
            .with_tags(vec!["test", "deletable"])
            .build();
        storage.save_extension(&ext2).unwrap();
        
        let profile1 = ProfileBuilder::new("Test Profile")
            .with_description("Test profile for view tests")
            .with_extensions(vec!["test-extension"])
            .as_default()
            .build();
        storage.save_profile(&profile1).unwrap();
        
        ViewManager::with_storage(storage)
    }

    #[tokio::test]
    async fn test_view_manager_initialization() {
        let vm = create_test_view_manager().await;
        
        // Should start at extensions list
        assert_eq!(vm.current_view(), ViewType::ExtensionList);
    }

    #[tokio::test]
    async fn test_tab_navigation() {
        let mut vm = create_test_view_manager().await;
        let (tx, _rx) = mpsc::unbounded_channel();
        vm.register_action_handler(tx).unwrap();
        
        // Start at extensions
        assert_eq!(vm.current_view(), ViewType::ExtensionList);
        
        // Navigate to profiles tab
        let result = vm.update(Action::NavigateToProfiles);
        assert!(result.is_ok());
        assert_eq!(vm.current_view(), ViewType::ProfileList);
        
        // Navigate back to extensions
        let result = vm.update(Action::NavigateToExtensions);
        assert!(result.is_ok());
        assert_eq!(vm.current_view(), ViewType::ExtensionList);
    }

    #[tokio::test]
    async fn test_create_extension_navigation() {
        let mut vm = create_test_view_manager().await;
        let (tx, _rx) = mpsc::unbounded_channel();
        vm.register_action_handler(tx).unwrap();
        
        // Navigate to create extension
        let result = vm.update(Action::CreateNewExtension);
        assert!(result.is_ok());
        assert_eq!(vm.current_view(), ViewType::ExtensionCreate);
        
        // Cancel and go back
        let result = vm.update(Action::NavigateToExtensions);
        assert!(result.is_ok());
        assert_eq!(vm.current_view(), ViewType::ExtensionList);
    }

    #[tokio::test]
    async fn test_create_profile_navigation() {
        let mut vm = create_test_view_manager().await;
        let (tx, _rx) = mpsc::unbounded_channel();
        vm.register_action_handler(tx).unwrap();
        
        // Navigate to profiles
        vm.update(Action::NavigateToProfiles).unwrap();
        assert_eq!(vm.current_view(), ViewType::ProfileList);
        
        // Navigate to create profile
        let result = vm.update(Action::CreateProfile);
        assert!(result.is_ok());
        assert_eq!(vm.current_view(), ViewType::ProfileCreate);
        
        // Cancel and go back
        let result = vm.update(Action::NavigateToProfiles);
        assert!(result.is_ok());
        assert_eq!(vm.current_view(), ViewType::ProfileList);
    }

    #[tokio::test]
    async fn test_extension_detail_navigation() {
        let mut vm = create_test_view_manager().await;
        let (tx, _rx) = mpsc::unbounded_channel();
        vm.register_action_handler(tx).unwrap();
        
        // View extension details
        let result = vm.update(Action::ViewExtensionDetails("test-extension".to_string()));
        assert!(result.is_ok());
        assert_eq!(vm.current_view(), ViewType::ExtensionDetail);
        
        // Go back
        let result = vm.update(Action::NavigateToExtensions);
        assert!(result.is_ok());
        assert_eq!(vm.current_view(), ViewType::ExtensionList);
    }

    #[tokio::test]
    async fn test_profile_detail_navigation() {
        let mut vm = create_test_view_manager().await;
        let (tx, _rx) = mpsc::unbounded_channel();
        vm.register_action_handler(tx).unwrap();
        
        // Navigate to profiles
        vm.update(Action::NavigateToProfiles).unwrap();
        
        // View profile details
        let result = vm.update(Action::ViewProfileDetails("test-profile".to_string()));
        assert!(result.is_ok());
        assert_eq!(vm.current_view(), ViewType::ProfileDetail);
        
        // Go back
        let result = vm.update(Action::NavigateToProfiles);
        assert!(result.is_ok());
        assert_eq!(vm.current_view(), ViewType::ProfileList);
    }

    #[tokio::test]
    async fn test_edit_extension_navigation() {
        let mut vm = create_test_view_manager().await;
        let (tx, _rx) = mpsc::unbounded_channel();
        vm.register_action_handler(tx).unwrap();
        
        // Edit extension
        let result = vm.update(Action::EditExtension("test-extension".to_string()));
        assert!(result.is_ok());
        assert_eq!(vm.current_view(), ViewType::ExtensionEdit);
    }

    #[tokio::test]
    async fn test_edit_profile_navigation() {
        let mut vm = create_test_view_manager().await;
        let (tx, _rx) = mpsc::unbounded_channel();
        vm.register_action_handler(tx).unwrap();
        
        // Navigate to profiles
        vm.update(Action::NavigateToProfiles).unwrap();
        
        // Edit profile
        let result = vm.update(Action::EditProfile("test-profile".to_string()));
        assert!(result.is_ok());
        assert_eq!(vm.current_view(), ViewType::ProfileEdit);
    }

    #[tokio::test]
    async fn test_delete_confirmation_flow() {
        let mut vm = create_test_view_manager().await;
        let (tx, _rx) = mpsc::unbounded_channel();
        vm.register_action_handler(tx).unwrap();
        
        // Delete extension (should show confirmation)
        let result = vm.update(Action::DeleteExtension("deletable-extension".to_string()));
        assert!(result.is_ok());
        assert_eq!(vm.current_view(), ViewType::ConfirmDelete);
        
        // Cancel deletion
        let result = vm.update(Action::CancelDelete);
        assert!(result.is_ok());
        assert_eq!(vm.current_view(), ViewType::ExtensionList);
    }

    #[tokio::test]
    async fn test_navigation_history() {
        let mut vm = create_test_view_manager().await;
        let (tx, _rx) = mpsc::unbounded_channel();
        vm.register_action_handler(tx).unwrap();
        
        // Navigate through multiple views
        assert_eq!(vm.current_view(), ViewType::ExtensionList);
        
        // Go to extension form
        vm.update(Action::CreateNewExtension).unwrap();
        assert_eq!(vm.current_view(), ViewType::ExtensionCreate);
        
        // Cancel and go back
        vm.update(Action::NavigateToExtensions).unwrap();
        assert_eq!(vm.current_view(), ViewType::ExtensionList);
        
        // Go to profiles
        vm.update(Action::NavigateToProfiles).unwrap();
        assert_eq!(vm.current_view(), ViewType::ProfileList);
        
        // Go to profile form
        vm.update(Action::CreateProfile).unwrap();
        assert_eq!(vm.current_view(), ViewType::ProfileCreate);
        
        // Cancel and verify we're back at profiles
        vm.update(Action::NavigateToProfiles).unwrap();
        assert_eq!(vm.current_view(), ViewType::ProfileList);
    }

    #[tokio::test]
    async fn test_refresh_actions() {
        let mut vm = create_test_view_manager().await;
        let (tx, _rx) = mpsc::unbounded_channel();
        vm.register_action_handler(tx).unwrap();
        
        // Test refresh extensions
        let result = vm.update(Action::RefreshExtensions);
        assert!(result.is_ok());
        
        // Navigate to profiles and test refresh
        vm.update(Action::NavigateToProfiles).unwrap();
        let result = vm.update(Action::RefreshProfiles);
        assert!(result.is_ok());
    }

    #[tokio::test]
    async fn test_error_handling() {
        let mut vm = create_test_view_manager().await;
        let (tx, _rx) = mpsc::unbounded_channel();
        vm.register_action_handler(tx).unwrap();
        
        // Send an error action
        let result = vm.update(Action::Error("Test error message".to_string()));
        assert!(result.is_ok());
        
        // Error should be displayed
        assert!(vm.has_error());
    }

    #[tokio::test]
    async fn test_launch_profile_action() {
        let mut vm = create_test_view_manager().await;
        let (tx, _rx) = mpsc::unbounded_channel();
        vm.register_action_handler(tx).unwrap();
        
        // Navigate to profiles
        vm.update(Action::NavigateToProfiles).unwrap();
        
        // Launch profile
        let result = vm.update(Action::LaunchWithProfile("test-profile".to_string()));
        assert!(result.is_ok());
    }

    // TODO: Add save action tests when SaveExtension and SaveProfile actions are implemented
    // #[tokio::test]
    // async fn test_save_extension_navigation() {
    //     let mut vm = create_test_view_manager().await;
    //     let (tx, _rx) = mpsc::unbounded_channel();
    //     vm.register_action_handler(tx).unwrap();
    //     
    //     // Go to create extension
    //     vm.update(Action::CreateNewExtension).unwrap();
    //     
    //     // Save should navigate back
    //     let result = vm.update(Action::SaveExtension);
    //     assert!(result.is_ok());
    //     assert_eq!(vm.current_view(), ViewType::ExtensionList);
    // }

    // #[tokio::test]
    // async fn test_save_profile_navigation() {
    //     let mut vm = create_test_view_manager().await;
    //     let (tx, _rx) = mpsc::unbounded_channel();
    //     vm.register_action_handler(tx).unwrap();
    //     
    //     // Navigate to profiles and create
    //     vm.update(Action::NavigateToProfiles).unwrap();
    //     vm.update(Action::CreateProfile).unwrap();
    //     
    //     // Save should navigate back
    //     let result = vm.update(Action::SaveProfile);
    //     assert!(result.is_ok());
    //     assert_eq!(vm.current_view(), ViewType::ProfileList);
    // }

    #[tokio::test]
    async fn test_delete_from_detail_view() {
        let mut vm = create_test_view_manager().await;
        let (tx, _rx) = mpsc::unbounded_channel();
        vm.register_action_handler(tx).unwrap();
        
        // View extension details
        vm.update(Action::ViewExtensionDetails("deletable-extension".to_string())).unwrap();
        assert_eq!(vm.current_view(), ViewType::ExtensionDetail);
        
        // Delete from detail view
        let result = vm.update(Action::DeleteExtension("deletable-extension".to_string()));
        assert!(result.is_ok());
        assert_eq!(vm.current_view(), ViewType::ConfirmDelete);
        
        // Confirm deletion should go back to list
        let result = vm.update(Action::ConfirmDelete);
        assert!(result.is_ok());
        assert_eq!(vm.current_view(), ViewType::ExtensionList);
    }

    #[tokio::test]
    async fn test_edit_from_detail_view() {
        let mut vm = create_test_view_manager().await;
        let (tx, _rx) = mpsc::unbounded_channel();
        vm.register_action_handler(tx).unwrap();
        
        // View extension details
        vm.update(Action::ViewExtensionDetails("test-extension".to_string())).unwrap();
        assert_eq!(vm.current_view(), ViewType::ExtensionDetail);
        
        // Edit from detail view
        let result = vm.update(Action::EditExtension("test-extension".to_string()));
        assert!(result.is_ok());
        assert_eq!(vm.current_view(), ViewType::ExtensionEdit);
        
        // Cancel should go back to detail view
        let result = vm.update(Action::NavigateToExtensions);
        assert!(result.is_ok());
        // Should return to detail view since we came from there
    }
}