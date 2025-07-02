#[cfg(test)]
mod tests {
    use crate::test_utils::*;
    use crossterm::event::KeyCode;

    // NOTE: The App struct doesn't expose public methods for testing navigation
    // These tests would need to be rewritten as component-level tests or
    // the App struct would need to expose test-friendly APIs

    #[test]
    fn test_extension_list_navigation() {
        // Test navigation within extension list component
        let storage = create_test_storage();

        // Add test data
        let ext1 = create_test_extension("Extension One");
        storage.save_extension(&ext1).unwrap();

        let ext2 = create_test_extension("Extension Two");
        storage.save_extension(&ext2).unwrap();

        use gemini_cli_manager::components::Component;
        use gemini_cli_manager::components::extension_list::ExtensionList;

        let mut list = ExtensionList::with_storage(storage);

        // Test navigation
        assert_eq!(list.selected_index(), 0);

        // Move down
        list.handle_events(Some(create_key_event(KeyCode::Down)))
            .unwrap();
        assert_eq!(list.selected_index(), 1);

        // Move up
        list.handle_events(Some(create_key_event(KeyCode::Up)))
            .unwrap();
        assert_eq!(list.selected_index(), 0);
    }

    #[test]
    fn test_profile_list_navigation() {
        // Test navigation within profile list component
        let storage = create_test_storage();

        // Add test data
        let profile1 = create_test_profile("Development");
        storage.save_profile(&profile1).unwrap();

        let profile2 = create_test_profile("Production");
        storage.save_profile(&profile2).unwrap();

        use gemini_cli_manager::components::Component;
        use gemini_cli_manager::components::profile_list::ProfileList;

        let mut list = ProfileList::with_storage(storage);

        // Test navigation
        assert_eq!(list.selected_index(), 0);

        // Move down
        list.handle_events(Some(create_key_event(KeyCode::Down)))
            .unwrap();
        assert_eq!(list.selected_index(), 1);

        // Move up
        list.handle_events(Some(create_key_event(KeyCode::Up)))
            .unwrap();
        assert_eq!(list.selected_index(), 0);
    }

    #[test]
    fn test_search_mode_activation() {
        let storage = create_test_storage();

        // Add test data
        let ext1 = create_test_extension("Test Extension");
        storage.save_extension(&ext1).unwrap();

        use gemini_cli_manager::components::Component;
        use gemini_cli_manager::components::extension_list::ExtensionList;

        let mut list = ExtensionList::with_storage(storage);

        // Not in search mode initially
        assert!(!list.is_search_mode());

        // Press '/' to enter search mode
        list.handle_events(Some(create_key_event(KeyCode::Char('/'))))
            .unwrap();
        assert!(list.is_search_mode());

        // Press Esc to exit search mode
        list.handle_events(Some(create_key_event(KeyCode::Esc)))
            .unwrap();
        assert!(!list.is_search_mode());
    }

    #[test]
    fn test_form_field_navigation() {
        let storage = create_test_storage();

        use gemini_cli_manager::components::Component;
        use gemini_cli_manager::components::extension_form::{ExtensionForm, FormField};

        let mut form = ExtensionForm::new(storage);

        // Should start at Name field
        assert_eq!(form.current_field(), &FormField::Name);

        // Tab to Version
        form.handle_events(Some(create_key_event(KeyCode::Tab)))
            .unwrap();
        assert_eq!(form.current_field(), &FormField::Version);

        // Tab to Description
        form.handle_events(Some(create_key_event(KeyCode::Tab)))
            .unwrap();
        assert_eq!(form.current_field(), &FormField::Description);

        // Shift+Tab back to Version
        form.handle_events(Some(create_key_event(KeyCode::BackTab)))
            .unwrap();
        assert_eq!(form.current_field(), &FormField::Version);
    }

    fn create_key_event(code: KeyCode) -> gemini_cli_manager::tui::Event {
        use crossterm::event::{KeyEvent, KeyEventKind, KeyModifiers};
        gemini_cli_manager::tui::Event::Key(KeyEvent {
            code,
            modifiers: KeyModifiers::NONE,
            kind: KeyEventKind::Press,
            state: crossterm::event::KeyEventState::NONE,
        })
    }
}
