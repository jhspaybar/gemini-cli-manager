#[cfg(test)]
mod tests {
    use gemini_cli_manager::{App, ViewType};
    use ratatui::prelude::*;
    use crate::test_utils::*;
    use crossterm::event::KeyCode;
    use insta::assert_snapshot;

    fn create_test_app() -> App {
        let storage = create_test_storage();
        
        // Add some test data
        let ext1 = create_test_extension("Extension One");
        storage.create_extension(ext1).unwrap();
        
        let profile1 = create_test_profile("Development");
        storage.create_profile(profile1).unwrap();
        
        App::new(storage)
    }

    #[test]
    fn test_tab_navigation() {
        let mut app = create_test_app();
        
        // Should start at extension list
        assert_eq!(app.current_view(), ViewType::ExtensionList);
        
        // Tab switches to profiles
        app.handle_key_event(KeyCode::Tab.into());
        assert_eq!(app.current_view(), ViewType::ProfileList);
        
        // Tab again goes back to extensions
        app.handle_key_event(KeyCode::Tab.into());
        assert_eq!(app.current_view(), ViewType::ExtensionList);
    }
    
    #[test]
    fn test_detail_view_navigation() {
        let mut app = create_test_app();
        
        // Start at extension list
        assert_eq!(app.current_view(), ViewType::ExtensionList);
        
        // Enter opens detail view
        app.handle_key_event(KeyCode::Enter.into());
        assert_eq!(app.current_view(), ViewType::ExtensionDetail);
        
        // Escape goes back to list
        app.handle_key_event(KeyCode::Esc.into());
        assert_eq!(app.current_view(), ViewType::ExtensionList);
    }
    
    #[test]
    fn test_create_form_navigation() {
        let mut app = create_test_app();
        
        // 'n' opens create form
        app.handle_key_event(KeyCode::Char('n').into());
        assert_eq!(app.current_view(), ViewType::ExtensionForm);
        
        // Escape (and confirm) returns to list
        app.handle_key_event(KeyCode::Esc.into());
        app.handle_key_event(KeyCode::Char('y').into()); // Confirm cancel
        assert_eq!(app.current_view(), ViewType::ExtensionList);
    }
    
    #[test]
    fn test_edit_navigation() {
        let mut app = create_test_app();
        
        // From list, 'e' opens edit form
        app.handle_key_event(KeyCode::Char('e').into());
        assert_eq!(app.current_view(), ViewType::ExtensionForm);
        
        // Should be in edit mode
        let mut terminal = setup_test_terminal(80, 30).unwrap();
        terminal.draw(|f| {
            app.render(f, f.area());
        }).unwrap();
        
        assert_buffer_contains(&terminal, "Edit Extension");
    }
    
    #[test]
    fn test_profile_navigation_flow() {
        let mut app = create_test_app();
        
        // Switch to profiles
        app.handle_key_event(KeyCode::Tab.into());
        assert_eq!(app.current_view(), ViewType::ProfileList);
        
        // Enter profile details
        app.handle_key_event(KeyCode::Enter.into());
        assert_eq!(app.current_view(), ViewType::ProfileDetail);
        
        // Edit profile
        app.handle_key_event(KeyCode::Char('e').into());
        assert_eq!(app.current_view(), ViewType::ProfileForm);
        
        // Cancel and go back
        app.handle_key_event(KeyCode::Esc.into());
        app.handle_key_event(KeyCode::Char('y').into());
        assert_eq!(app.current_view(), ViewType::ProfileDetail);
        
        // Back to list
        app.handle_key_event(KeyCode::Esc.into());
        assert_eq!(app.current_view(), ViewType::ProfileList);
    }
    
    #[test]
    fn test_search_mode_navigation() {
        let mut app = create_test_app();
        
        // Enter search mode
        app.handle_key_event(KeyCode::Char('/').into());
        
        // Should still be in extension list but in search mode
        assert_eq!(app.current_view(), ViewType::ExtensionList);
        assert!(app.is_search_active());
        
        // Type search query
        for ch in "test".chars() {
            app.handle_key_event(KeyCode::Char(ch).into());
        }
        
        // Escape exits search
        app.handle_key_event(KeyCode::Esc.into());
        assert!(!app.is_search_active());
    }
    
    #[test]
    fn test_navigation_history() {
        let mut app = create_test_app();
        let mut history = vec![];
        
        // Track navigation
        history.push(app.current_view());
        
        // Navigate through views
        app.handle_key_event(KeyCode::Enter.into()); // Detail
        history.push(app.current_view());
        
        app.handle_key_event(KeyCode::Char('e').into()); // Edit
        history.push(app.current_view());
        
        // Verify history
        assert_eq!(history, vec![
            ViewType::ExtensionList,
            ViewType::ExtensionDetail,
            ViewType::ExtensionForm,
        ]);
    }
    
    #[test]
    fn test_keyboard_shortcuts() {
        let mut app = create_test_app();
        
        // Test various shortcuts
        let shortcuts = vec![
            (KeyCode::Char('n'), ViewType::ExtensionForm),  // New
            (KeyCode::Esc, ViewType::ExtensionList),        // Cancel
            (KeyCode::Char('e'), ViewType::ExtensionForm),  // Edit
            (KeyCode::Esc, ViewType::ExtensionList),        // Cancel
            (KeyCode::Tab, ViewType::ProfileList),          // Switch tab
            (KeyCode::Char('n'), ViewType::ProfileForm),    // New profile
        ];
        
        for (key, expected_view) in shortcuts {
            app.handle_key_event(key.into());
            
            // Handle confirmation dialogs if needed
            if app.is_confirming() {
                app.handle_key_event(KeyCode::Char('y').into());
            }
            
            assert_eq!(
                app.current_view(), 
                expected_view, 
                "Failed for key {:?}", 
                key
            );
        }
    }
    
    #[test]
    fn test_navigation_snapshot() {
        let app = create_test_app();
        
        // Render main view
        let main_view = render_to_string(80, 30, |f| {
            app.render(f, f.area());
        }).unwrap();
        
        assert_snapshot!("navigation_main_view", main_view);
    }
    
    #[test]
    fn test_error_state_navigation() {
        let mut app = create_test_app();
        
        // Simulate an error
        app.show_error("Test error message");
        
        // Error should be displayed
        let mut terminal = setup_test_terminal(80, 30).unwrap();
        terminal.draw(|f| {
            app.render(f, f.area());
        }).unwrap();
        
        assert_buffer_contains(&terminal, "Test error message");
        
        // Any key should dismiss error
        app.handle_key_event(KeyCode::Enter.into());
        
        // Should remain in same view
        assert_eq!(app.current_view(), ViewType::ExtensionList);
    }
}

// Implement HandleKeyEvent trait for App
impl HandleKeyEvent for gemini_cli_manager::App {
    fn handle_key_event(&mut self, event: crossterm::event::KeyEvent) {
        // This would call the actual handle_key_event method on App
        // The implementation depends on the actual App API
    }
}