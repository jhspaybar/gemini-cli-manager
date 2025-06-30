#[cfg(test)]
mod tests {
    use gemini_cli_manager::components::{ExtensionList, ListMode};
    use gemini_cli_manager::models::Extension;
    use ratatui::prelude::*;
    use crate::test_utils::*;
    use crossterm::event::KeyCode;

    fn create_test_list() -> ExtensionList {
        let storage = create_test_storage();
        
        // Add some test extensions
        let ext1 = ExtensionBuilder::new("Extension One")
            .with_description("First test extension")
            .with_tags(vec!["productivity", "test"])
            .build();
        storage.create_extension(ext1).unwrap();
        
        let ext2 = ExtensionBuilder::new("Extension Two")
            .with_description("Second test extension")
            .with_tags(vec!["development", "test"])
            .build();
        storage.create_extension(ext2).unwrap();
        
        let ext3 = ExtensionBuilder::new("Another Extension")
            .with_description("Third test extension")
            .with_tags(vec!["utility"])
            .build();
        storage.create_extension(ext3).unwrap();
        
        ExtensionList::new(storage)
    }

    #[test]
    fn test_extension_list_rendering() {
        let list = create_test_list();
        let mut terminal = setup_test_terminal(60, 20).unwrap();
        
        terminal.draw(|f| {
            list.render(f, f.area());
        }).unwrap();
        
        // Verify extensions are displayed
        assert_buffer_contains(&terminal, "Extension One");
        assert_buffer_contains(&terminal, "Extension Two");
        assert_buffer_contains(&terminal, "Another Extension");
    }
    
    #[test]
    fn test_extension_list_empty_state() {
        let storage = create_test_storage();
        let list = ExtensionList::new(storage);
        let mut terminal = setup_test_terminal(60, 20).unwrap();
        
        terminal.draw(|f| {
            list.render(f, f.area());
        }).unwrap();
        
        // Verify empty state message
        assert_buffer_contains(&terminal, "No extensions found");
        assert_buffer_contains(&terminal, "Press 'n' to create");
    }
    
    #[test]
    fn test_keyboard_navigation() {
        let mut list = create_test_list();
        
        // Should start at first item
        assert_eq!(list.selected_index(), 0);
        
        // Move down
        list.handle_key_event(KeyCode::Down.into());
        assert_eq!(list.selected_index(), 1);
        
        // Move down again
        list.handle_key_event(KeyCode::Down.into());
        assert_eq!(list.selected_index(), 2);
        
        // At last item, down should stay
        list.handle_key_event(KeyCode::Down.into());
        assert_eq!(list.selected_index(), 2);
        
        // Move up
        list.handle_key_event(KeyCode::Up.into());
        assert_eq!(list.selected_index(), 1);
        
        // Jump to top with Home
        list.handle_key_event(KeyCode::Home.into());
        assert_eq!(list.selected_index(), 0);
        
        // Jump to bottom with End
        list.handle_key_event(KeyCode::End.into());
        assert_eq!(list.selected_index(), 2);
    }
    
    #[test]
    fn test_vim_navigation() {
        let mut list = create_test_list();
        
        // j moves down
        list.handle_key_event(KeyCode::Char('j').into());
        assert_eq!(list.selected_index(), 1);
        
        // k moves up
        list.handle_key_event(KeyCode::Char('k').into());
        assert_eq!(list.selected_index(), 0);
    }
    
    #[test]
    fn test_search_mode() {
        let mut list = create_test_list();
        
        // Should start in normal mode
        assert_eq!(list.mode(), ListMode::Normal);
        
        // '/' activates search
        list.handle_key_event(KeyCode::Char('/').into());
        assert_eq!(list.mode(), ListMode::Search);
        
        // Type search query
        list.handle_key_event(KeyCode::Char('t').into());
        list.handle_key_event(KeyCode::Char('w').into());
        list.handle_key_event(KeyCode::Char('o').into());
        
        assert_eq!(list.search_query(), "two");
        
        // Should filter to matching extensions
        assert_eq!(list.filtered_count(), 1);
        
        // Escape exits search
        list.handle_key_event(KeyCode::Esc.into());
        assert_eq!(list.mode(), ListMode::Normal);
        assert_eq!(list.search_query(), "");
    }
    
    #[test]
    fn test_search_filtering() {
        let mut list = create_test_list();
        
        // Enter search mode
        list.handle_key_event(KeyCode::Char('/').into());
        
        // Search for "extension"
        for ch in "extension".chars() {
            list.handle_key_event(KeyCode::Char(ch).into());
        }
        
        // All items should match (all have "Extension" in name)
        assert_eq!(list.filtered_count(), 3);
        
        // Clear and search for "another"
        list.clear_search();
        for ch in "another".chars() {
            list.handle_key_event(KeyCode::Char(ch).into());
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
        let result = list.handle_key_event(KeyCode::Char('d').into());
        
        // Should show confirmation dialog
        assert!(list.is_confirming_delete());
        
        // Cancel with 'n'
        list.handle_key_event(KeyCode::Char('n').into());
        assert!(!list.is_confirming_delete());
        
        // Extension should still exist
        assert_eq!(list.total_count(), 3);
    }
    
    #[test]
    fn test_minimum_size_rendering() {
        let list = create_test_list();
        
        // Test very small terminal
        let mut terminal = setup_test_terminal(20, 5).unwrap();
        
        // Should not panic
        let result = terminal.draw(|f| {
            list.render(f, f.area());
        });
        
        assert!(result.is_ok());
    }
    
    #[test]
    fn test_selection_bounds() {
        let mut list = create_test_list();
        
        // At first item, up should stay at 0
        list.handle_key_event(KeyCode::Up.into());
        assert_eq!(list.selected_index(), 0);
        
        // Jump to end
        list.handle_key_event(KeyCode::End.into());
        assert_eq!(list.selected_index(), 2);
        
        // At last item, down should stay at last
        list.handle_key_event(KeyCode::Down.into());
        assert_eq!(list.selected_index(), 2);
    }
}

// Implement the HandleKeyEvent trait for ExtensionList
impl HandleKeyEvent for gemini_cli_manager::components::ExtensionList {
    fn handle_key_event(&mut self, event: crossterm::event::KeyEvent) {
        // This would call the actual handle_key_event method on ExtensionList
        // The implementation depends on the actual ExtensionList API
    }
}