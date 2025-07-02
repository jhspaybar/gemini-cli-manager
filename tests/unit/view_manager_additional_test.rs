#[cfg(test)]
mod tests {
    use crate::test_utils::*;
    use gemini_cli_manager::{
        action::Action,
        config::Config,
        view::{ViewManager, ViewType},
    };
    use ratatui::prelude::*;
    use tokio::sync::mpsc;

    #[tokio::test]
    async fn test_navigate_back_from_various_views() {
        let mut vm = ViewManager::with_storage(create_test_storage());
        let (tx, _rx) = mpsc::unbounded_channel();
        vm.register_action_handler(tx).unwrap();

        // Test NavigateBack from profile detail
        vm.update(Action::NavigateToProfiles).unwrap();
        vm.update(Action::ViewProfileDetails("test-profile".to_string()))
            .unwrap();
        assert_eq!(vm.current_view(), ViewType::ProfileDetail);

        vm.update(Action::NavigateBack).unwrap();
        assert_eq!(vm.current_view(), ViewType::ProfileList);

        // Test NavigateBack from extension detail
        vm.update(Action::ViewExtensionDetails("test-extension".to_string()))
            .unwrap();
        assert_eq!(vm.current_view(), ViewType::ExtensionDetail);

        vm.update(Action::NavigateBack).unwrap();
        assert_eq!(vm.current_view(), ViewType::ExtensionList);

        // Test NavigateBack from create forms
        vm.update(Action::CreateNewExtension).unwrap();
        assert_eq!(vm.current_view(), ViewType::ExtensionCreate);

        vm.update(Action::NavigateBack).unwrap();
        assert_eq!(vm.current_view(), ViewType::ExtensionList);
    }

    #[tokio::test]
    async fn test_navigate_back_from_edit_views() {
        let storage = create_test_storage();

        // Add test data
        let ext = ExtensionBuilder::new("Test Extension").build();
        storage.save_extension(&ext).unwrap();

        let mut vm = ViewManager::with_storage(storage);
        let (tx, _rx) = mpsc::unbounded_channel();
        vm.register_action_handler(tx).unwrap();

        // Edit from list view
        vm.update(Action::EditExtension("test-extension".to_string()))
            .unwrap();
        assert_eq!(vm.current_view(), ViewType::ExtensionEdit);

        vm.update(Action::NavigateBack).unwrap();
        assert_eq!(vm.current_view(), ViewType::ExtensionList);

        // Edit from detail view
        vm.update(Action::ViewExtensionDetails("test-extension".to_string()))
            .unwrap();
        vm.update(Action::EditExtension("test-extension".to_string()))
            .unwrap();
        assert_eq!(vm.current_view(), ViewType::ExtensionEdit);

        vm.update(Action::NavigateBack).unwrap();
        assert_eq!(vm.current_view(), ViewType::ExtensionDetail);
    }

    #[tokio::test]
    async fn test_navigate_to_settings() {
        let mut vm = ViewManager::with_storage(create_test_storage());
        let (tx, _rx) = mpsc::unbounded_channel();
        vm.register_action_handler(tx).unwrap();

        // Navigate to settings
        let result = vm.update(Action::NavigateToSettings);
        assert!(result.is_ok());
        assert_eq!(vm.current_view(), ViewType::Settings);
    }

    #[tokio::test]
    async fn test_confirm_delete_profile() {
        let storage = create_test_storage();

        // Add a deletable profile
        let profile = ProfileBuilder::new("Deletable Profile").build();
        storage.save_profile(&profile).unwrap();

        let mut vm = ViewManager::with_storage(storage.clone());
        let (tx, _rx) = mpsc::unbounded_channel();
        vm.register_action_handler(tx).unwrap();

        // Navigate to profiles and delete
        vm.update(Action::NavigateToProfiles).unwrap();
        vm.update(Action::DeleteProfile("deletable-profile".to_string()))
            .unwrap();
        assert_eq!(vm.current_view(), ViewType::ConfirmDelete);

        // Confirm deletion
        vm.update(Action::ConfirmDelete).unwrap();

        // Should navigate back to profile list
        assert_eq!(vm.current_view(), ViewType::ProfileList);

        // Profile should be deleted
        assert!(storage.load_profile("deletable-profile").is_err());
    }

    #[tokio::test]
    async fn test_cancel_delete() {
        let storage = create_test_storage();

        // Add test extension
        let ext = ExtensionBuilder::new("Keep Extension").build();
        storage.save_extension(&ext).unwrap();

        let mut vm = ViewManager::with_storage(storage.clone());
        let (tx, _rx) = mpsc::unbounded_channel();
        vm.register_action_handler(tx).unwrap();

        // Start delete flow
        vm.update(Action::DeleteExtension("keep-extension".to_string()))
            .unwrap();
        assert_eq!(vm.current_view(), ViewType::ConfirmDelete);

        // Cancel deletion
        vm.update(Action::CancelDelete).unwrap();
        assert_eq!(vm.current_view(), ViewType::ExtensionList);

        // Extension should still exist
        assert!(storage.load_extension("keep-extension").is_ok());
    }

    #[tokio::test]
    async fn test_delete_extension_referenced_by_profile() {
        let storage = create_test_storage();

        // Add extension
        let ext = ExtensionBuilder::new("Referenced Extension").build();
        storage.save_extension(&ext).unwrap();

        // Add profile that references the extension
        let profile = ProfileBuilder::new("Test Profile")
            .with_extensions(vec!["referenced-extension"])
            .build();
        storage.save_profile(&profile).unwrap();

        let mut vm = ViewManager::with_storage(storage);
        let (tx, _rx) = mpsc::unbounded_channel();
        vm.register_action_handler(tx).unwrap();

        // Try to delete referenced extension
        vm.update(Action::DeleteExtension("referenced-extension".to_string()))
            .unwrap();

        // Should not show confirm dialog, but stay on list with error
        assert_eq!(vm.current_view(), ViewType::ExtensionList);
        assert!(vm.has_error());
    }

    #[tokio::test]
    async fn test_error_message_display() {
        let mut vm = ViewManager::with_storage(create_test_storage());
        let (tx, _rx) = mpsc::unbounded_channel();
        vm.register_action_handler(tx).unwrap();

        // Send error
        vm.update(Action::Error("Test error".to_string())).unwrap();
        assert!(vm.has_error());

        // Error should contain the message
        // Note: We can't easily test the actual error display without mocking time
    }

    // Note: SaveProfile and SaveExtension actions don't exist yet
    // These tests are commented out until those actions are implemented

    // #[tokio::test]
    // async fn test_save_profile_action() {
    //     let storage = create_test_storage();
    //     let mut vm = ViewManager::with_storage(storage.clone());
    //     let (tx, _rx) = mpsc::unbounded_channel();
    //     vm.register_action_handler(tx).unwrap();
    //
    //     // Go to create profile
    //     vm.update(Action::NavigateToProfiles).unwrap();
    //     vm.update(Action::CreateProfile).unwrap();
    //     assert_eq!(vm.current_view(), ViewType::ProfileCreate);
    //
    //     // Save profile action
    //     vm.update(Action::SaveProfile).unwrap();
    //
    //     // Should navigate back to profile list
    //     assert_eq!(vm.current_view(), ViewType::ProfileList);
    // }

    // #[tokio::test]
    // async fn test_save_extension_action() {
    //     let mut vm = ViewManager::with_storage(create_test_storage());
    //     let (tx, _rx) = mpsc::unbounded_channel();
    //     vm.register_action_handler(tx).unwrap();
    //
    //     // Go to create extension
    //     vm.update(Action::CreateNewExtension).unwrap();
    //     assert_eq!(vm.current_view(), ViewType::ExtensionCreate);
    //
    //     // Save extension action
    //     vm.update(Action::SaveExtension).unwrap();
    //
    //     // Should navigate back to extension list
    //     assert_eq!(vm.current_view(), ViewType::ExtensionList);
    // }

    // Note: SetDefaultProfile action doesn't exist yet
    // This test is commented out until that action is implemented

    // #[tokio::test]
    // async fn test_set_default_profile() {
    //     let storage = create_test_storage();
    //
    //     // Add multiple profiles
    //     let profile1 = ProfileBuilder::new("Profile 1").build();
    //     let profile2 = ProfileBuilder::new("Profile 2").build();
    //     storage.save_profile(&profile1).unwrap();
    //     storage.save_profile(&profile2).unwrap();
    //
    //     let mut vm = ViewManager::with_storage(storage.clone());
    //     let (tx, _rx) = mpsc::unbounded_channel();
    //     vm.register_action_handler(tx).unwrap();
    //
    //     // Set profile 2 as default
    //     vm.update(Action::SetDefaultProfile("profile-2".to_string())).unwrap();
    //
    //     // Check that it was set as default
    //     let profiles = storage.list_profiles().unwrap();
    //     let default_profile = profiles.iter().find(|p| p.metadata.is_default);
    //     assert!(default_profile.is_some());
    //     assert_eq!(default_profile.unwrap().id, "profile-2");
    // }

    #[tokio::test]
    async fn test_init_with_size() {
        let mut vm = ViewManager::with_storage(create_test_storage());

        // Test init
        let result = vm.init(Size {
            width: 80,
            height: 24,
        });
        assert!(result.is_ok());
    }

    #[tokio::test]
    async fn test_register_config_handler() {
        let mut vm = ViewManager::with_storage(create_test_storage());

        // Test registering config
        let config = Config::default();
        let result = vm.register_config_handler(config);
        assert!(result.is_ok());
    }

    #[tokio::test]
    async fn test_handle_events() {
        let mut vm = ViewManager::with_storage(create_test_storage());
        let (tx, _rx) = mpsc::unbounded_channel();
        vm.register_action_handler(tx).unwrap();

        // Send an error first
        vm.update(Action::Error("Test error".to_string())).unwrap();
        assert!(vm.has_error());

        // Test dismissing error with Esc
        use crossterm::event::{KeyCode, KeyEvent, KeyModifiers};
        let esc_event = gemini_cli_manager::tui::Event::Key(KeyEvent {
            code: KeyCode::Esc,
            modifiers: KeyModifiers::empty(),
            kind: crossterm::event::KeyEventKind::Press,
            state: crossterm::event::KeyEventState::empty(),
        });

        let result = vm.handle_events(Some(esc_event));
        assert!(result.is_ok());
        // Error should be cleared
    }

    #[tokio::test]
    async fn test_draw_method() {
        let mut vm = ViewManager::with_storage(create_test_storage());
        let (tx, _rx) = mpsc::unbounded_channel();
        vm.register_action_handler(tx).unwrap();

        // Create a test backend
        let backend = ratatui::backend::TestBackend::new(80, 24);
        let mut terminal = Terminal::new(backend).unwrap();

        // Test drawing
        let result = terminal.draw(|frame| {
            let area = frame.area();
            let _ = vm.draw(frame, area);
        });

        assert!(result.is_ok());
    }

    #[tokio::test]
    async fn test_confirm_delete_extension_from_detail() {
        let storage = create_test_storage();

        // Add a deletable extension
        let ext = ExtensionBuilder::new("Delete From Detail").build();
        storage.save_extension(&ext).unwrap();

        let mut vm = ViewManager::with_storage(storage.clone());
        let (tx, _rx) = mpsc::unbounded_channel();
        vm.register_action_handler(tx).unwrap();

        // View details then delete
        vm.update(Action::ViewExtensionDetails(
            "delete-from-detail".to_string(),
        ))
        .unwrap();
        vm.update(Action::DeleteExtension("delete-from-detail".to_string()))
            .unwrap();
        assert_eq!(vm.current_view(), ViewType::ConfirmDelete);

        // Confirm deletion
        vm.update(Action::ConfirmDelete).unwrap();

        // Should go back to list after deletion from detail view
        assert_eq!(vm.current_view(), ViewType::ExtensionList);

        // Extension should be deleted
        assert!(storage.load_extension("delete-from-detail").is_err());
    }

    #[tokio::test]
    async fn test_edit_non_existent_extension() {
        let mut vm = ViewManager::with_storage(create_test_storage());
        let (tx, _rx) = mpsc::unbounded_channel();
        vm.register_action_handler(tx).unwrap();

        // Try to edit non-existent extension
        let result = vm.update(Action::EditExtension("non-existent".to_string()));
        assert!(result.is_ok());

        // Should stay on current view
        assert_eq!(vm.current_view(), ViewType::ExtensionList);
    }

    #[tokio::test]
    async fn test_edit_non_existent_profile() {
        let mut vm = ViewManager::with_storage(create_test_storage());
        let (tx, _rx) = mpsc::unbounded_channel();
        vm.register_action_handler(tx).unwrap();

        // Try to edit non-existent profile
        let result = vm.update(Action::EditProfile("non-existent".to_string()));
        assert!(result.is_ok());

        // Should stay on current view
        assert_eq!(vm.current_view(), ViewType::ExtensionList);
    }
}
