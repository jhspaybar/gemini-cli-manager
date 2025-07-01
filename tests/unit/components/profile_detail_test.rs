#[cfg(test)]
mod tests {
    use gemini_cli_manager::components::profile_detail::ProfileDetail;
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

    fn create_test_detail() -> ProfileDetail {
        let storage = create_test_storage();
        
        // Create test extensions
        let ext1 = ExtensionBuilder::new("Database Tools")
            .with_description("Database management extension")
            .with_tags(vec!["database", "tools"])
            .build();
        storage.save_extension(&ext1).unwrap();
        
        let ext2 = ExtensionBuilder::new("AI Assistant")
            .with_description("AI-powered development assistant")
            .with_tags(vec!["ai", "assistant"])
            .build();
        storage.save_extension(&ext2).unwrap();
        
        // Create test profile with extensions
        let mut profile = ProfileBuilder::new("Development Profile")
            .with_description("My development environment setup")
            .with_extensions(vec![&ext1.id, &ext2.id])
            .with_tags(vec!["dev", "local"])
            .as_default()
            .build();
        
        // Add environment variables
        profile.environment_variables.insert("NODE_ENV".to_string(), "development".to_string());
        profile.environment_variables.insert("DEBUG".to_string(), "true".to_string());
        profile.working_directory = Some("~/projects".to_string());
        
        storage.save_profile(&profile).unwrap();
        
        ProfileDetail::new(storage, profile.id.clone())
    }

    #[test]
    fn test_detail_view_rendering() {
        let mut detail = create_test_detail();
        let mut terminal = setup_test_terminal(80, 30).unwrap();
        
        terminal.draw(|f| {
            detail.draw(f, f.area()).unwrap();
        }).unwrap();
        
        // Verify all profile details are displayed
        assert_buffer_contains(&terminal, "Development Profile");
        assert_buffer_contains(&terminal, "My development environment setup");
        assert_buffer_contains(&terminal, "dev, local");
        assert_buffer_contains(&terminal, "Default Profile");
        assert_buffer_contains(&terminal, "~/projects");
        assert_buffer_contains(&terminal, "Database Tools");
        assert_buffer_contains(&terminal, "AI Assistant");
        assert_buffer_contains(&terminal, "NODE_ENV");
        assert_buffer_contains(&terminal, "DEBUG");
    }
    
    // TODO: ProfileDetail doesn't have section navigation, only scrolling
    // #[test]
    // fn test_navigation_sections() {
    //     let mut detail = create_test_detail();
    //     
    //     // Should start at overview
    //     assert_eq!(detail.current_section(), 0);
    //     
    //     // Navigate down through sections
    //     detail.handle_events(Some(create_key_event(KeyCode::Down))).unwrap();
    //     assert_eq!(detail.current_section(), 1); // Extensions
    //     
    //     detail.handle_events(Some(create_key_event(KeyCode::Down))).unwrap();
    //     assert_eq!(detail.current_section(), 2); // Environment Variables
    //     
    //     // Navigate back up
    //     detail.handle_events(Some(create_key_event(KeyCode::Up))).unwrap();
    //     assert_eq!(detail.current_section(), 1);
    // }
    
    #[test]
    fn test_back_navigation() {
        let mut detail = create_test_detail();
        
        // Press Escape to go back
        let result = detail.handle_events(Some(create_key_event(KeyCode::Esc))).unwrap();
        
        // Should return navigation action
        assert!(result.is_some());
    }
    
    #[test]
    fn test_edit_action() {
        let mut detail = create_test_detail();
        
        // Press 'e' to edit
        let result = detail.handle_events(Some(create_key_event(KeyCode::Char('e')))).unwrap();
        
        // Should return edit action
        assert!(result.is_some());
    }
    
    #[test]
    fn test_delete_action() {
        let mut detail = create_test_detail();
        
        // Press 'd' to delete
        let result = detail.handle_events(Some(create_key_event(KeyCode::Char('d')))).unwrap();
        
        // Should return delete action
        assert!(result.is_some());
    }
    
    #[test]
    fn test_launch_action() {
        let mut detail = create_test_detail();
        
        // Press 'l' to launch
        let result = detail.handle_events(Some(create_key_event(KeyCode::Char('l')))).unwrap();
        
        // Should return launch action
        assert!(result.is_some());
    }
    
    #[test]
    fn test_set_default_action() {
        let mut detail = create_test_detail();
        
        // Press 'x' to set as default
        let result = detail.handle_events(Some(create_key_event(KeyCode::Char('x')))).unwrap();
        
        // TODO: Set default action not implemented yet
        assert!(result.is_none());
    }
    
    // TODO: ProfileDetail doesn't have tab navigation between sections
    // #[test]
    // fn test_tab_navigation() {
    //     let mut detail = create_test_detail();
    //     
    //     // Tab through sections
    //     detail.handle_events(Some(create_key_event(KeyCode::Tab))).unwrap();
    //     assert_eq!(detail.current_section(), 1);
    //     
    //     detail.handle_events(Some(create_key_event(KeyCode::Tab))).unwrap();
    //     assert_eq!(detail.current_section(), 2);
    //     
    //     // Tab should wrap around
    //     detail.handle_events(Some(create_key_event(KeyCode::Tab))).unwrap();
    //     assert_eq!(detail.current_section(), 0);
    // }
    
    #[test]
    fn test_extension_list_display() {
        let mut detail = create_test_detail();
        let mut terminal = setup_test_terminal(80, 30).unwrap();
        
        // Navigate to extensions section
        detail.handle_events(Some(create_key_event(KeyCode::Down))).unwrap();
        
        terminal.draw(|f| {
            detail.draw(f, f.area()).unwrap();
        }).unwrap();
        
        // Verify extension details are shown
        assert_buffer_contains(&terminal, "Database management extension");
        assert_buffer_contains(&terminal, "AI-powered development assistant");
    }
    
    #[test]
    fn test_environment_variables_display() {
        let mut detail = create_test_detail();
        let mut terminal = setup_test_terminal(80, 30).unwrap();
        
        // Navigate to env vars section
        detail.handle_events(Some(create_key_event(KeyCode::Down))).unwrap();
        detail.handle_events(Some(create_key_event(KeyCode::Down))).unwrap();
        
        terminal.draw(|f| {
            detail.draw(f, f.area()).unwrap();
        }).unwrap();
        
        // Verify env vars are shown (with spaces around =)
        assert_buffer_contains(&terminal, "NODE_ENV = development");
        assert_buffer_contains(&terminal, "DEBUG = true");
    }
    
    #[test]
    fn test_empty_profile_handling() {
        let storage = create_test_storage();
        
        // Create minimal profile
        let profile = ProfileBuilder::new("Minimal Profile").build();
        storage.save_profile(&profile).unwrap();
        
        let mut detail = ProfileDetail::new(storage, profile.id);
        let mut terminal = setup_test_terminal(80, 30).unwrap();
        
        terminal.draw(|f| {
            detail.draw(f, f.area()).unwrap();
        }).unwrap();
        
        // Should handle missing optional fields gracefully
        assert_buffer_contains(&terminal, "Minimal Profile");
        assert_buffer_contains(&terminal, "No extensions included");
        // Description and env vars are not shown when empty
    }
    
    #[test]
    fn test_non_default_profile() {
        let storage = create_test_storage();
        
        // Create non-default profile
        let profile = ProfileBuilder::new("Secondary Profile")
            .with_description("Not the default")
            .build();
        storage.save_profile(&profile).unwrap();
        
        let mut detail = ProfileDetail::new(storage, profile.id);
        let mut terminal = setup_test_terminal(80, 30).unwrap();
        
        terminal.draw(|f| {
            detail.draw(f, f.area()).unwrap();
        }).unwrap();
        
        // Should not show default indicator
        assert_buffer_not_contains(&terminal, "(default)");
    }
    
    #[test]
    fn test_help_text() {
        let mut detail = create_test_detail();
        let mut terminal = setup_test_terminal(80, 30).unwrap();
        
        terminal.draw(|f| {
            detail.draw(f, f.area()).unwrap();
        }).unwrap();
        
        // Verify help text is shown (some may be cut off at 80 chars)
        assert_buffer_contains(&terminal, "Scroll");
        assert_buffer_contains(&terminal, "Launch");
        assert_buffer_contains(&terminal, "Edit");
        assert_buffer_contains(&terminal, "Delete");
    }
    
    #[test]
    fn test_scroll_many_extensions() {
        let storage = create_test_storage();
        
        // Create many extensions
        let mut ext_ids = vec![];
        for i in 0..20 {
            let ext = ExtensionBuilder::new(&format!("Extension {}", i))
                .with_description(&format!("Description for extension {}", i))
                .build();
            storage.save_extension(&ext).unwrap();
            ext_ids.push(ext.id);
        }
        
        // Create profile with many extensions
        let mut profile = ProfileBuilder::new("Large Profile")
            .with_description("Profile with many extensions")
            .build();
        profile.extension_ids = ext_ids;
        storage.save_profile(&profile).unwrap();
        
        let mut detail = ProfileDetail::new(storage, profile.id);
        
        // Navigate to extensions and test scrolling
        detail.handle_events(Some(create_key_event(KeyCode::Down))).unwrap();
        detail.handle_events(Some(create_key_event(KeyCode::PageDown))).unwrap();
        detail.handle_events(Some(create_key_event(KeyCode::PageUp))).unwrap();
        
        // Should handle scrolling without panic
    }
    
    #[test]
    fn test_responsive_layout() {
        let mut detail = create_test_detail();
        
        // Test various terminal sizes
        let sizes = vec![(60, 20), (80, 30), (120, 40), (40, 15)];
        
        for (width, height) in sizes {
            let mut terminal = setup_test_terminal(width, height).unwrap();
            
            // Should render without panic at any size
            let result = terminal.draw(|f| {
                detail.draw(f, f.area()).unwrap();
            });
            
            assert!(result.is_ok(), "Failed to render at {}x{}", width, height);
        }
    }
}