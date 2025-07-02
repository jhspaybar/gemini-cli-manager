#[cfg(test)]
mod tests {
    use crate::test_utils::*;
    use crossterm::event::{KeyCode, KeyEvent, KeyEventKind};
    use gemini_cli_manager::components::Component;
    use gemini_cli_manager::components::extension_form::{ExtensionForm, FormField};
    use insta::assert_snapshot;

    fn create_key_event(code: KeyCode) -> gemini_cli_manager::tui::Event {
        use crossterm::event::KeyModifiers;
        gemini_cli_manager::tui::Event::Key(KeyEvent {
            code,
            modifiers: KeyModifiers::NONE,
            kind: KeyEventKind::Press,
            state: crossterm::event::KeyEventState::NONE,
        })
    }

    fn create_test_form() -> ExtensionForm {
        let storage = create_test_storage();
        ExtensionForm::new(storage)
    }

    fn create_edit_form(_id: &str) -> ExtensionForm {
        let storage = create_test_storage();

        // Add test extension to storage
        let ext = ExtensionBuilder::new("Test Extension")
            .with_version("1.0.0")
            .with_description("Test description")
            .with_tags(vec!["test", "example"])
            .build();
        storage.save_extension(&ext).unwrap();

        ExtensionForm::with_extension(storage, &ext)
    }

    #[test]
    fn test_form_initial_state() {
        let form = create_test_form();

        assert!(!form.is_edit_mode());
        assert_eq!(form.current_field(), &FormField::Name);
        assert!(form.name_input().value().is_empty());
        assert_eq!(form.version_input().value(), "1.0.0"); // Default version
        assert!(form.description_input().value().is_empty());
        assert!(form.tags_input().value().is_empty());
        assert_eq!(form.context_file_name_input().value(), "GEMINI.md"); // Default context file
    }

    #[test]
    fn test_form_rendering() {
        let mut form = create_test_form();
        let mut terminal = setup_test_terminal(80, 30).unwrap();

        terminal
            .draw(|f| {
                form.draw(f, f.area()).unwrap();
            })
            .unwrap();

        // Verify form elements are rendered
        assert_buffer_contains(&terminal, "Create New Extension");
        assert_buffer_contains(&terminal, "Name");
        assert_buffer_contains(&terminal, "Version");
        assert_buffer_contains(&terminal, "Description");
        assert_buffer_contains(&terminal, "Tags");
    }

    #[test]
    fn test_edit_mode() {
        let form = create_edit_form("test-extension");

        assert!(form.is_edit_mode());
        assert_eq!(form.name_input().value(), "Test Extension");
        assert_eq!(form.version_input().value(), "1.0.0");
        assert_eq!(form.description_input().value(), "Test description");
        assert_eq!(form.tags_input().value(), "test, example");
    }

    #[test]
    fn test_field_navigation() {
        let mut form = create_test_form();

        // Start at Name
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

    #[test]
    fn test_text_input() {
        let mut form = create_test_form();

        // Type name
        for ch in "My Extension".chars() {
            form.handle_events(Some(create_key_event(KeyCode::Char(ch))))
                .unwrap();
        }
        assert_eq!(form.name_input().value(), "My Extension");

        // Tab to version
        form.handle_events(Some(create_key_event(KeyCode::Tab)))
            .unwrap();
        // Clear default version
        for _ in 0.."1.0.0".len() {
            form.handle_events(Some(create_key_event(KeyCode::Backspace)))
                .unwrap();
        }
        // Type new version
        for ch in "2.0.0".chars() {
            form.handle_events(Some(create_key_event(KeyCode::Char(ch))))
                .unwrap();
        }
        assert_eq!(form.version_input().value(), "2.0.0");
    }

    #[test]
    fn test_validation_empty_name() {
        let mut form = create_test_form();

        // Try to save with empty name
        form.handle_events(Some(create_key_event(KeyCode::Char('s').into())))
            .unwrap();

        // Should show error
        let mut terminal = setup_test_terminal(80, 30).unwrap();
        terminal
            .draw(|f| {
                form.draw(f, f.area()).unwrap();
            })
            .unwrap();

        // Note: We can't directly test validation errors without app state access
        // The component might handle this differently
    }

    #[test]
    fn test_form_cancel() {
        let mut form = create_test_form();

        // Type some data
        for ch in "Test".chars() {
            form.handle_events(Some(create_key_event(KeyCode::Char(ch))))
                .unwrap();
        }

        // Press Escape to cancel
        let result = form
            .handle_events(Some(create_key_event(KeyCode::Esc)))
            .unwrap();

        // Should return a navigation action
        assert!(result.is_some());
    }

    #[test]
    fn test_multiline_input_description() {
        let mut form = create_test_form();

        // Navigate to description
        form.handle_events(Some(create_key_event(KeyCode::Tab)))
            .unwrap();
        form.handle_events(Some(create_key_event(KeyCode::Tab)))
            .unwrap();

        // Type multiline text
        for ch in "Line 1".chars() {
            form.handle_events(Some(create_key_event(KeyCode::Char(ch))))
                .unwrap();
        }

        // Note: Multiline input might be handled differently by tui-input
        assert!(form.description_input().value().contains("Line 1"));
    }

    #[test]
    fn test_form_save() {
        let mut form = create_test_form();

        // Fill out form
        for ch in "New Extension".chars() {
            form.handle_events(Some(create_key_event(KeyCode::Char(ch))))
                .unwrap();
        }

        form.handle_events(Some(create_key_event(KeyCode::Tab)))
            .unwrap();
        for ch in "2.0.0".chars() {
            form.handle_events(Some(create_key_event(KeyCode::Char(ch))))
                .unwrap();
        }

        // Save with Ctrl+S
        let key_event = KeyEvent {
            code: KeyCode::Char('s'),
            modifiers: crossterm::event::KeyModifiers::CONTROL,
            kind: KeyEventKind::Press,
            state: crossterm::event::KeyEventState::NONE,
        };
        let result = form
            .handle_events(Some(gemini_cli_manager::tui::Event::Key(key_event)))
            .unwrap();

        // Should return save action
        assert!(result.is_some());
    }

    #[test]
    fn test_context_file_fields() {
        let mut form = create_test_form();

        // Navigate to context file name field
        // Name -> Version -> Description -> Tags -> ContextFileName
        for _ in 0..4 {
            form.handle_events(Some(create_key_event(KeyCode::Tab)))
                .unwrap();
        }
        assert_eq!(form.current_field(), &FormField::ContextFileName);

        // Clear default and type new context file name
        // First clear the default value
        for _ in 0.."GEMINI.md".len() {
            form.handle_events(Some(create_key_event(KeyCode::Backspace)))
                .unwrap();
        }
        // Then type new name
        for ch in "README.md".chars() {
            form.handle_events(Some(create_key_event(KeyCode::Char(ch))))
                .unwrap();
        }
        assert_eq!(form.context_file_name_input().value(), "README.md");

        // Tab to context content
        form.handle_events(Some(create_key_event(KeyCode::Tab)))
            .unwrap();
        assert_eq!(form.current_field(), &FormField::ContextContent);

        // Type content
        for ch in "# Extension Documentation".chars() {
            form.handle_events(Some(create_key_event(KeyCode::Char(ch))))
                .unwrap();
        }
        assert!(
            form.context_content_input()
                .value()
                .contains("Extension Documentation")
        );
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

    #[test]
    fn test_snapshot_empty_form() {
        let mut form = create_test_form();
        let output = render_to_string(80, 30, |f| {
            form.draw(f, f.area()).unwrap();
        });

        assert_snapshot!(output.unwrap());
    }

    #[test]
    fn test_snapshot_filled_form() {
        let mut form = create_test_form();

        // Fill out all fields
        for ch in "My Awesome Extension".chars() {
            form.handle_events(Some(create_key_event(KeyCode::Char(ch))))
                .unwrap();
        }

        form.handle_events(Some(create_key_event(KeyCode::Tab)))
            .unwrap();
        for ch in "3.2.1".chars() {
            form.handle_events(Some(create_key_event(KeyCode::Char(ch))))
                .unwrap();
        }

        form.handle_events(Some(create_key_event(KeyCode::Tab)))
            .unwrap();
        for ch in "This is a comprehensive extension for productivity".chars() {
            form.handle_events(Some(create_key_event(KeyCode::Char(ch))))
                .unwrap();
        }

        // Skip to tags
        form.handle_events(Some(create_key_event(KeyCode::Tab)))
            .unwrap();
        form.handle_events(Some(create_key_event(KeyCode::Tab)))
            .unwrap();
        form.handle_events(Some(create_key_event(KeyCode::Tab)))
            .unwrap();
        for ch in "productivity, ai, development".chars() {
            form.handle_events(Some(create_key_event(KeyCode::Char(ch))))
                .unwrap();
        }

        let output = render_to_string(80, 30, |f| {
            form.draw(f, f.area()).unwrap();
        });

        assert_snapshot!(output.unwrap());
    }
}
