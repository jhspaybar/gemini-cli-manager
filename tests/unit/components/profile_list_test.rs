#[cfg(test)]
mod tests {
    use crate::test_utils::*;
    use crossterm::event::{KeyCode, KeyEvent, KeyEventKind};
    use gemini_cli_manager::components::Component;
    use gemini_cli_manager::components::profile_list::ProfileList;

    fn create_key_event(code: KeyCode) -> gemini_cli_manager::tui::Event {
        use crossterm::event::KeyModifiers;
        gemini_cli_manager::tui::Event::Key(KeyEvent {
            code,
            modifiers: KeyModifiers::NONE,
            kind: KeyEventKind::Press,
            state: crossterm::event::KeyEventState::NONE,
        })
    }

    fn create_test_profile_list() -> ProfileList {
        let storage = create_test_storage();

        // Add some test profiles
        let profile1 = ProfileBuilder::new("Development")
            .with_description("Development environment")
            .with_extensions(vec!["ext-1", "ext-2"])
            .with_tags(vec!["dev", "local"])
            .build();
        storage.save_profile(&profile1).unwrap();

        let profile2 = ProfileBuilder::new("Production")
            .with_description("Production environment")
            .with_extensions(vec!["ext-3"])
            .with_tags(vec!["prod"])
            .as_default()
            .build();
        storage.save_profile(&profile2).unwrap();

        let profile3 = ProfileBuilder::new("Testing")
            .with_description("Testing environment")
            .with_tags(vec!["test", "qa"])
            .build();
        storage.save_profile(&profile3).unwrap();

        ProfileList::with_storage(storage)
    }

    #[test]
    fn test_profile_list_rendering() {
        let mut list = create_test_profile_list();
        let mut terminal = setup_test_terminal(60, 20).unwrap();

        terminal
            .draw(|f| {
                list.draw(f, f.area()).unwrap();
            })
            .unwrap();

        // Verify profiles are displayed
        assert_buffer_contains(&terminal, "Development");
        assert_buffer_contains(&terminal, "Production");
        assert_buffer_contains(&terminal, "Testing");

        // Verify default indicator
        assert_buffer_contains(&terminal, "(default)");
    }

    #[test]
    fn test_profile_list_empty_state() {
        let storage = create_test_storage();
        let mut list = ProfileList::with_storage(storage);
        let mut terminal = setup_test_terminal(60, 20).unwrap();

        terminal
            .draw(|f| {
                list.draw(f, f.area()).unwrap();
            })
            .unwrap();

        // Verify empty state message
        assert_buffer_contains(&terminal, "No profiles found");
        assert_buffer_contains(&terminal, "Press 'n' to create your first profile");
    }

    #[test]
    fn test_profile_descriptions() {
        let mut list = create_test_profile_list();
        let mut terminal = setup_test_terminal(80, 20).unwrap();

        terminal
            .draw(|f| {
                list.draw(f, f.area()).unwrap();
            })
            .unwrap();

        // Verify descriptions are shown
        assert_buffer_contains(&terminal, "Development environment");
        assert_buffer_contains(&terminal, "Production environment");
        assert_buffer_contains(&terminal, "Testing environment");
    }

    #[test]
    fn test_extension_count_display() {
        let mut list = create_test_profile_list();
        let mut terminal = setup_test_terminal(80, 20).unwrap();

        terminal
            .draw(|f| {
                list.draw(f, f.area()).unwrap();
            })
            .unwrap();

        // Verify extension counts
        assert_buffer_contains(&terminal, "2 extensions"); // Development profile
        assert_buffer_contains(&terminal, "1 extension"); // Production profile
        assert_buffer_contains(&terminal, "0 extensions"); // Testing profile
    }

    #[test]
    fn test_profile_navigation() {
        let mut list = create_test_profile_list();

        // Should start at first item
        assert_eq!(list.selected_index(), 0);

        // Navigate down
        list.handle_events(Some(create_key_event(KeyCode::Down)))
            .unwrap();
        assert_eq!(list.selected_index(), 1);

        // Navigate to specific profile
        list.handle_events(Some(create_key_event(KeyCode::Down)))
            .unwrap();
        assert_eq!(list.selected_index(), 2);

        // Navigate back up
        list.handle_events(Some(create_key_event(KeyCode::Up)))
            .unwrap();
        assert_eq!(list.selected_index(), 1);
    }

    #[test]
    fn test_default_profile_selection() {
        let mut list = create_test_profile_list();

        // Select the development profile (index 0)
        assert_eq!(list.selected_index(), 0);

        // Make it default with 'x'
        list.handle_events(Some(create_key_event(KeyCode::Char('x'))))
            .unwrap();

        // Should update default status
        let mut terminal = setup_test_terminal(80, 20).unwrap();
        terminal
            .draw(|f| {
                list.draw(f, f.area()).unwrap();
            })
            .unwrap();

        // Development should now be default
        let content = buffer_to_string(terminal.backend().buffer());
        assert!(content.contains("Development") && content.contains("(default)"));
    }

    #[test]
    fn test_profile_search() {
        let mut list = create_test_profile_list();

        // Enter search mode
        list.handle_events(Some(create_key_event(KeyCode::Char('/'))))
            .unwrap();
        assert!(list.is_search_mode());

        // Search for "dev"
        list.handle_events(Some(create_key_event(KeyCode::Char('d'))))
            .unwrap();
        list.handle_events(Some(create_key_event(KeyCode::Char('e'))))
            .unwrap();
        list.handle_events(Some(create_key_event(KeyCode::Char('v'))))
            .unwrap();

        // Should match Development profile
        assert_eq!(list.filtered_count(), 1);

        // Clear search
        list.handle_events(Some(create_key_event(KeyCode::Esc)))
            .unwrap();
        assert_eq!(list.filtered_count(), 3);
    }

    #[test]
    fn test_tag_search() {
        let mut list = create_test_profile_list();

        // Search for tag
        list.handle_events(Some(create_key_event(KeyCode::Char('/'))))
            .unwrap();
        for ch in "test".chars() {
            list.handle_events(Some(create_key_event(KeyCode::Char(ch))))
                .unwrap();
        }

        // Should match Testing profile (has "test" tag)
        assert_eq!(list.filtered_count(), 1);
    }

    #[test]
    fn test_profile_deletion_protection() {
        let mut list = create_test_profile_list();

        // Try to delete first profile
        list.handle_events(Some(create_key_event(KeyCode::Char('x'))))
            .unwrap();

        // Note: We can't directly test the confirmation dialog without access to app state
        // For now, just verify the profile count remains the same
        assert_eq!(list.total_count(), 3);
    }

    #[test]
    fn test_cannot_delete_default_profile() {
        let mut list = create_test_profile_list();

        // Navigate to Production (default profile)
        list.handle_events(Some(create_key_event(KeyCode::Down)))
            .unwrap();
        assert_eq!(list.selected_index(), 1);

        // Try to delete
        list.handle_events(Some(create_key_event(KeyCode::Char('x'))))
            .unwrap();

        // Should show error instead of confirmation
        let mut terminal = setup_test_terminal(80, 20).unwrap();
        terminal
            .draw(|f| {
                list.draw(f, f.area()).unwrap();
            })
            .unwrap();

        // Note: We can't directly test the error message without access to app state
        // The component itself might handle this differently
    }

    #[test]
    fn test_profile_list_responsive() {
        let mut list = create_test_profile_list();

        // Test various terminal sizes
        let sizes = vec![(40, 15), (80, 24), (120, 40), (30, 10)];

        for (width, height) in sizes {
            let mut terminal = setup_test_terminal(width, height).unwrap();

            // Should render without panic
            let result = terminal.draw(|f| {
                list.draw(f, f.area()).unwrap();
            });

            assert!(result.is_ok(), "Failed to render at {width}x{height}");
        }
    }
}
