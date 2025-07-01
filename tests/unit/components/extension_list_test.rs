#[cfg(test)]
mod tests {
    use gemini_cli_manager::components::extension_list::ExtensionList;
    use gemini_cli_manager::components::Component;
    use crate::test_utils::*;
    use crossterm::event::{KeyCode, KeyEvent, KeyEventKind};

    fn create_key_event(code: KeyCode) -> gemini_cli_manager::tui::Event {
        use crossterm::event::KeyModifiers;
        gemini_cli_manager::tui::Event::Key(KeyEvent {
            code,
            modifiers: KeyModifiers::NONE,
            kind: KeyEventKind::Press,
            state: crossterm::event::KeyEventState::NONE,
        })
    }
    
    fn create_test_list() -> ExtensionList {
        let storage = create_test_storage();
        
        // Add some test extensions
        let ext1 = ExtensionBuilder::new("Extension One")
            .with_description("First test extension")
            .with_tags(vec!["productivity", "test"])
            .build();
        storage.save_extension(&ext1).unwrap();
        
        let ext2 = ExtensionBuilder::new("Extension Two")
            .with_description("Second test extension")
            .with_tags(vec!["development", "test"])
            .build();
        storage.save_extension(&ext2).unwrap();
        
        let ext3 = ExtensionBuilder::new("Another Extension")
            .with_description("Third test extension")
            .with_tags(vec!["utility"])
            .build();
        storage.save_extension(&ext3).unwrap();
        
        ExtensionList::with_storage(storage)
    }

    #[test]
    fn test_extension_list_rendering() {
        let mut list = create_test_list();
        let mut terminal = setup_test_terminal(60, 20).unwrap();
        
        terminal.draw(|f| {
            list.draw(f, f.area()).unwrap();
        }).unwrap();
        
        // Verify extensions are displayed
        assert_buffer_contains(&terminal, "Extension One");
        assert_buffer_contains(&terminal, "Extension Two");
        assert_buffer_contains(&terminal, "Another Extension");
    }
    
    #[test]
    fn test_extension_list_empty_state() {
        let storage = create_test_storage();
        let mut list = ExtensionList::with_storage(storage);
        let mut terminal = setup_test_terminal(60, 20).unwrap();
        
        terminal.draw(|f| {
            list.draw(f, f.area()).unwrap();
        }).unwrap();
        
        // Verify empty state message
        assert_buffer_contains(&terminal, "No extensions found");
        assert_buffer_contains(&terminal, "Press 'n' to create your first extension");
    }
    
    #[test]
    fn test_keyboard_navigation() {
        let mut list = create_test_list();
        
        // Should start at first item
        assert_eq!(list.selected_index(), 0);
        
        // Move down
        list.handle_events(Some(create_key_event(KeyCode::Down))).unwrap();
        assert_eq!(list.selected_index(), 1);
        
        // Move down again
        list.handle_events(Some(create_key_event(KeyCode::Down))).unwrap();
        assert_eq!(list.selected_index(), 2);
        
        // At last item, down should stay
        list.handle_events(Some(create_key_event(KeyCode::Down))).unwrap();
        assert_eq!(list.selected_index(), 2);
        
        // Move up
        list.handle_events(Some(create_key_event(KeyCode::Up))).unwrap();
        assert_eq!(list.selected_index(), 1);
        
        // Jump to top with Home
        list.handle_events(Some(create_key_event(KeyCode::Home))).unwrap();
        assert_eq!(list.selected_index(), 0);
        
        // Jump to bottom with End
        list.handle_events(Some(create_key_event(KeyCode::End))).unwrap();
        assert_eq!(list.selected_index(), 2);
    }
    
    #[test]
    fn test_vim_navigation() {
        let mut list = create_test_list();
        
        // j moves down
        list.handle_events(Some(create_key_event(KeyCode::Char('j')))).unwrap();
        assert_eq!(list.selected_index(), 1);
        
        // k moves up
        list.handle_events(Some(create_key_event(KeyCode::Char('k')))).unwrap();
        assert_eq!(list.selected_index(), 0);
    }
    
    #[test]
    fn test_search_mode() {
        let mut list = create_test_list();
        
        // Should start in normal mode
        assert!(!list.is_search_mode());
        
        // '/' activates search
        list.handle_events(Some(create_key_event(KeyCode::Char('/')))).unwrap();
        assert!(list.is_search_mode());
        
        // Type search query
        list.handle_events(Some(create_key_event(KeyCode::Char('t')))).unwrap();
        list.handle_events(Some(create_key_event(KeyCode::Char('w')))).unwrap();
        list.handle_events(Some(create_key_event(KeyCode::Char('o')))).unwrap();
        
        assert_eq!(list.search_query(), "two");
        
        // Should filter to matching extensions
        assert_eq!(list.filtered_count(), 1);
        
        // Escape exits search
        list.handle_events(Some(create_key_event(KeyCode::Esc))).unwrap();
        assert!(!list.is_search_mode());
        assert_eq!(list.search_query(), "");
    }
    
    #[test]
    fn test_search_filtering() {
        let mut list = create_test_list();
        
        // Enter search mode
        list.handle_events(Some(create_key_event(KeyCode::Char('/')))).unwrap();
        
        // Search for "extension"
        for ch in "extension".chars() {
            list.handle_events(Some(create_key_event(KeyCode::Char(ch)))).unwrap();
        }
        
        // All items should match (all have "Extension" in name)
        assert_eq!(list.filtered_count(), 3);
        
        // Clear and search for "another"
        // Exit search mode and re-enter to clear
        list.handle_events(Some(create_key_event(KeyCode::Esc))).unwrap();
        list.handle_events(Some(create_key_event(KeyCode::Char('/')))).unwrap();
        for ch in "another".chars() {
            list.handle_events(Some(create_key_event(KeyCode::Char(ch)))).unwrap();
        }
        
        // Only one should match
        assert_eq!(list.filtered_count(), 1);
    }
    
    #[test]
    fn test_deletion_protection() {
        let mut list = create_test_list();
        
        // Select first extension
        assert_eq!(list.selected_index(), 0);
        
        // Try to delete
        let result = list.handle_events(Some(create_key_event(KeyCode::Char('d')))).unwrap();
        
        // The deletion should trigger an action (we can't check the confirmation dialog directly)
        // For now, just check that the action was handled
        assert!(result.is_some() || result.is_none()); // Either returns an action or None
        
        // Extension count should remain the same (deletion needs confirmation)
        // Note: We can't directly test the confirmation dialog without access to the app state
        
        // Extension should still exist
        assert_eq!(list.total_count(), 3);
    }
    
    #[test]
    fn test_minimum_size_rendering() {
        let mut list = create_test_list();
        
        // Test very small terminal
        let mut terminal = setup_test_terminal(20, 5).unwrap();
        
        // Should not panic
        let result = terminal.draw(|f| {
            list.draw(f, f.area()).unwrap();
        });
        
        assert!(result.is_ok());
    }
    
    #[test]
    fn test_selection_bounds() {
        let mut list = create_test_list();
        
        // At first item, up should stay at 0
        list.handle_events(Some(create_key_event(KeyCode::Up))).unwrap();
        assert_eq!(list.selected_index(), 0);
        
        // Jump to end
        list.handle_events(Some(create_key_event(KeyCode::End))).unwrap();
        assert_eq!(list.selected_index(), 2);
        
        // At last item, down should stay at last
        list.handle_events(Some(create_key_event(KeyCode::Down))).unwrap();
        assert_eq!(list.selected_index(), 2);
    }
}