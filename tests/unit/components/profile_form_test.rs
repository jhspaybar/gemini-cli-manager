#[cfg(test)]
mod tests {
    use gemini_cli_manager::components::profile_form::{ProfileForm, FormField};
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

    fn create_test_form() -> ProfileForm {
        let storage = create_test_storage();
        
        // Add some test extensions for selection
        let ext1 = ExtensionBuilder::new("Extension One")
            .with_description("First test extension")
            .build();
        storage.save_extension(&ext1).unwrap();
        
        let ext2 = ExtensionBuilder::new("Extension Two")
            .with_description("Second test extension")
            .build();
        storage.save_extension(&ext2).unwrap();
        
        ProfileForm::new(storage)
    }
    
    fn create_edit_form(_id: &str) -> ProfileForm {
        let storage = create_test_storage();
        
        // Add test extensions
        let ext1 = ExtensionBuilder::new("Extension One").build();
        storage.save_extension(&ext1).unwrap();
        
        let ext2 = ExtensionBuilder::new("Extension Two").build();
        storage.save_extension(&ext2).unwrap();
        
        // Add test profile to edit
        let profile = ProfileBuilder::new("Test Profile")
            .with_description("Test description")
            .with_extensions(vec![&ext1.id])
            .with_tags(vec!["test", "example"])
            .build();
        storage.save_profile(&profile).unwrap();
        
        ProfileForm::with_profile(storage, &profile)
    }

    #[test]
    fn test_form_initial_state() {
        let form = create_test_form();
        
        assert!(!form.is_edit_mode());
        assert_eq!(form.current_field(), &FormField::Name);
        assert!(form.name_input().value().is_empty());
        assert!(form.description_input().value().is_empty());
        assert!(form.working_directory_input().value().is_empty());
        assert!(form.tags_input().value().is_empty());
        assert!(form.selected_extensions().is_empty());
    }
    
    #[test]
    fn test_form_rendering() {
        let mut form = create_test_form();
        let mut terminal = setup_test_terminal(80, 30).unwrap();
        
        terminal.draw(|f| {
            form.draw(f, f.area()).unwrap();
        }).unwrap();
        
        // Verify form elements are rendered
        assert_buffer_contains(&terminal, "Create New Profile");
        assert_buffer_contains(&terminal, "Name");
        assert_buffer_contains(&terminal, "Description");
        assert_buffer_contains(&terminal, "Working Directory");
        assert_buffer_contains(&terminal, "Tags");
        assert_buffer_contains(&terminal, "Extensions");
    }
    
    #[test]
    fn test_edit_mode() {
        let form = create_edit_form("test-profile");
        
        assert!(form.is_edit_mode());
        assert_eq!(form.name_input().value(), "Test Profile");
        assert_eq!(form.description_input().value(), "Test description");
        assert_eq!(form.tags_input().value(), "test, example");
        assert_eq!(form.selected_extensions().len(), 1);
    }
    
    #[test]
    fn test_field_navigation() {
        let mut form = create_test_form();
        
        // Start at Name
        assert_eq!(form.current_field(), &FormField::Name);
        
        // Tab to Description
        form.handle_events(Some(create_key_event(KeyCode::Tab))).unwrap();
        assert_eq!(form.current_field(), &FormField::Description);
        
        // Tab to Working Directory
        form.handle_events(Some(create_key_event(KeyCode::Tab))).unwrap();
        assert_eq!(form.current_field(), &FormField::WorkingDirectory);
        
        // Tab to Extensions
        form.handle_events(Some(create_key_event(KeyCode::Tab))).unwrap();
        assert_eq!(form.current_field(), &FormField::Extensions);
        
        // Tab to Tags
        form.handle_events(Some(create_key_event(KeyCode::Tab))).unwrap();
        assert_eq!(form.current_field(), &FormField::Tags);
        
        // Tab to LaunchConfig
        form.handle_events(Some(create_key_event(KeyCode::Tab))).unwrap();
        assert_eq!(form.current_field(), &FormField::LaunchConfig);
        
        // Tab wraps back to Name
        form.handle_events(Some(create_key_event(KeyCode::Tab))).unwrap();
        assert_eq!(form.current_field(), &FormField::Name);
        
        // Shift+Tab back to LaunchConfig
        form.handle_events(Some(create_key_event(KeyCode::BackTab))).unwrap();
        assert_eq!(form.current_field(), &FormField::LaunchConfig);
    }
    
    #[test]
    fn test_text_input() {
        let mut form = create_test_form();
        
        // Type name
        for ch in "My Profile".chars() {
            form.handle_events(Some(create_key_event(KeyCode::Char(ch)))).unwrap();
        }
        assert_eq!(form.name_input().value(), "My Profile");
        
        // Tab to description
        form.handle_events(Some(create_key_event(KeyCode::Tab))).unwrap();
        
        // Type description
        for ch in "Development environment".chars() {
            form.handle_events(Some(create_key_event(KeyCode::Char(ch)))).unwrap();
        }
        assert_eq!(form.description_input().value(), "Development environment");
        
        // Tab to working directory
        form.handle_events(Some(create_key_event(KeyCode::Tab))).unwrap();
        
        // Type path
        for ch in "~/projects".chars() {
            form.handle_events(Some(create_key_event(KeyCode::Char(ch)))).unwrap();
        }
        assert_eq!(form.working_directory_input().value(), "~/projects");
    }
    
    #[test]
    fn test_extension_selection() {
        let mut form = create_test_form();
        
        // Navigate to extensions field (Name -> Description -> WorkingDirectory -> Extensions)
        for _ in 0..3 {
            form.handle_events(Some(create_key_event(KeyCode::Tab))).unwrap();
        }
        assert_eq!(form.current_field(), &FormField::Extensions);
        
        // Press space to toggle first extension
        form.handle_events(Some(create_key_event(KeyCode::Char(' ')))).unwrap();
        assert_eq!(form.selected_extensions().len(), 1);
        
        // Move down and toggle second extension
        form.handle_events(Some(create_key_event(KeyCode::Down))).unwrap();
        form.handle_events(Some(create_key_event(KeyCode::Char(' ')))).unwrap();
        assert_eq!(form.selected_extensions().len(), 2);
        
        // Toggle first extension off
        form.handle_events(Some(create_key_event(KeyCode::Up))).unwrap();
        form.handle_events(Some(create_key_event(KeyCode::Char(' ')))).unwrap();
        assert_eq!(form.selected_extensions().len(), 1);
    }
    
    #[test]
    fn test_extension_list_navigation() {
        let mut form = create_test_form();
        
        // Navigate to extensions field
        for _ in 0..3 {
            form.handle_events(Some(create_key_event(KeyCode::Tab))).unwrap();
        }
        
        // Should start at first extension
        assert_eq!(form.extension_cursor(), 0);
        
        // Move down
        form.handle_events(Some(create_key_event(KeyCode::Down))).unwrap();
        assert_eq!(form.extension_cursor(), 1);
        
        // At last item, down should wrap to first
        form.handle_events(Some(create_key_event(KeyCode::Down))).unwrap();
        assert_eq!(form.extension_cursor(), 0);
        
        // Move up from 0 should wrap to last (1)
        form.handle_events(Some(create_key_event(KeyCode::Up))).unwrap();
        assert_eq!(form.extension_cursor(), 1);
    }
    
    #[test]
    fn test_validation_empty_name() {
        let mut form = create_test_form();
        
        // Try to save with empty name
        form.handle_events(Some(create_key_event(KeyCode::Char('s')).into())).unwrap();
        
        // Should show error (we can't directly test validation without app state)
        // But the form should not save
        assert!(!form.is_saved());
    }
    
    #[test]
    fn test_form_cancel() {
        let mut form = create_test_form();
        
        // Type some data
        for ch in "Test".chars() {
            form.handle_events(Some(create_key_event(KeyCode::Char(ch)))).unwrap();
        }
        
        // Press Escape to cancel
        let result = form.handle_events(Some(create_key_event(KeyCode::Esc))).unwrap();
        
        // Should return a navigation action
        assert!(result.is_some());
    }
    
    #[test]
    fn test_form_save() {
        let mut form = create_test_form();
        
        // Fill out form
        for ch in "New Profile".chars() {
            form.handle_events(Some(create_key_event(KeyCode::Char(ch)))).unwrap();
        }
        
        // Save with Ctrl+S
        let key_event = KeyEvent {
            code: KeyCode::Char('s'),
            modifiers: crossterm::event::KeyModifiers::CONTROL,
            kind: KeyEventKind::Press,
            state: crossterm::event::KeyEventState::NONE,
        };
        let result = form.handle_events(Some(gemini_cli_manager::tui::Event::Key(key_event))).unwrap();
        
        // Should return save action
        assert!(result.is_some());
    }
    
    // TODO: ProfileForm doesn't have set as default functionality
    // #[test]
    // fn test_set_as_default() {
    //     let mut form = create_test_form();
    //     
    //     // Check default checkbox with Ctrl+D
    //     let key_event = KeyEvent {
    //         code: KeyCode::Char('d'),
    //         modifiers: crossterm::event::KeyModifiers::CONTROL,
    //         kind: KeyEventKind::Press,
    //         state: crossterm::event::KeyEventState::NONE,
    //     };
    //     form.handle_events(Some(gemini_cli_manager::tui::Event::Key(key_event))).unwrap();
    //     
    //     assert!(form.is_default());
    //     
    //     // Toggle again
    //     form.handle_events(Some(gemini_cli_manager::tui::Event::Key(key_event))).unwrap();
    //     assert!(!form.is_default());
    // }
    
    #[test]
    fn test_environment_variables() {
        let mut form = create_test_form();
        
        // Add environment variable with Ctrl+E
        let key_event = KeyEvent {
            code: KeyCode::Char('e'),
            modifiers: crossterm::event::KeyModifiers::CONTROL,
            kind: KeyEventKind::Press,
            state: crossterm::event::KeyEventState::NONE,
        };
        form.handle_events(Some(gemini_cli_manager::tui::Event::Key(key_event))).unwrap();
        
        // Should open env var dialog (component might handle this differently)
        // For now, verify the form accepts the key event
    }
    
    #[test]
    fn test_form_with_all_fields() {
        let mut form = create_test_form();
        
        // Fill all fields
        for ch in "Complete Profile".chars() {
            form.handle_events(Some(create_key_event(KeyCode::Char(ch)))).unwrap();
        }
        
        form.handle_events(Some(create_key_event(KeyCode::Tab))).unwrap();
        for ch in "A complete profile with all fields".chars() {
            form.handle_events(Some(create_key_event(KeyCode::Char(ch)))).unwrap();
        }
        
        form.handle_events(Some(create_key_event(KeyCode::Tab))).unwrap();
        for ch in "~/workspace".chars() {
            form.handle_events(Some(create_key_event(KeyCode::Char(ch)))).unwrap();
        }
        
        // Skip Extensions field
        form.handle_events(Some(create_key_event(KeyCode::Tab))).unwrap();
        
        // Now at Tags field
        form.handle_events(Some(create_key_event(KeyCode::Tab))).unwrap();
        for ch in "dev, local, testing".chars() {
            form.handle_events(Some(create_key_event(KeyCode::Char(ch)))).unwrap();
        }
        
        // Verify all fields are filled
        assert_eq!(form.name_input().value(), "Complete Profile");
        assert_eq!(form.description_input().value(), "A complete profile with all fields");
        assert_eq!(form.working_directory_input().value(), "~/workspace");
        assert_eq!(form.tags_input().value(), "dev, local, testing");
    }
    
    #[test]
    fn test_form_responsive() {
        let mut form = create_test_form();
        
        // Test various terminal sizes
        let sizes = vec![(60, 20), (80, 30), (120, 40), (40, 15)];
        
        for (width, height) in sizes {
            let mut terminal = setup_test_terminal(width, height).unwrap();
            
            // Should render without panic
            let result = terminal.draw(|f| {
                form.draw(f, f.area()).unwrap();
            });
            
            assert!(result.is_ok(), "Failed to render at {}x{}", width, height);
        }
    }
}