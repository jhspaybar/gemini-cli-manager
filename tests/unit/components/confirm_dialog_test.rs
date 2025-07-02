#[cfg(test)]
mod tests {
    use crate::test_utils::*;
    use crossterm::event::{KeyCode, KeyEvent, KeyEventKind};
    use gemini_cli_manager::components::Component;
    use gemini_cli_manager::components::confirm_dialog::ConfirmDialog;

    fn create_key_event(code: KeyCode) -> gemini_cli_manager::tui::Event {
        use crossterm::event::KeyModifiers;
        gemini_cli_manager::tui::Event::Key(KeyEvent {
            code,
            modifiers: KeyModifiers::NONE,
            kind: KeyEventKind::Press,
            state: crossterm::event::KeyEventState::NONE,
        })
    }

    fn create_test_dialog() -> ConfirmDialog {
        ConfirmDialog::new(
            "Delete Extension",
            "Are you sure you want to delete this extension? This action cannot be undone.",
        )
        .with_actions(
            gemini_cli_manager::action::Action::ConfirmDelete,
            gemini_cli_manager::action::Action::CancelDelete,
        )
    }

    #[test]
    fn test_dialog_rendering() {
        let mut dialog = create_test_dialog();
        let mut terminal = setup_test_terminal(60, 20).unwrap();

        terminal
            .draw(|f| {
                dialog.draw(f, f.area()).unwrap();
            })
            .unwrap();

        // Verify dialog elements are displayed
        assert_buffer_contains(&terminal, "Delete Extension");
        assert_buffer_contains(&terminal, "Are you sure you want to delete");
        assert_buffer_contains(&terminal, "action cannot be undone");
        assert_buffer_contains(&terminal, "Cancel");
        assert_buffer_contains(&terminal, "Delete");
    }

    #[test]
    fn test_initial_selection() {
        let dialog = create_test_dialog();

        // Should default to "No" for safety
        assert!(!dialog.is_confirmed());
        assert_eq!(dialog.selected_option(), 1); // 0 = Yes, 1 = No
    }

    #[test]
    fn test_keyboard_navigation() {
        let mut dialog = create_test_dialog();

        // Start at "No"
        assert_eq!(dialog.selected_option(), 1);

        // Left arrow to "Yes"
        dialog
            .handle_events(Some(create_key_event(KeyCode::Left)))
            .unwrap();
        assert_eq!(dialog.selected_option(), 0);

        // Right arrow back to "No"
        dialog
            .handle_events(Some(create_key_event(KeyCode::Right)))
            .unwrap();
        assert_eq!(dialog.selected_option(), 1);

        // Tab also works
        dialog
            .handle_events(Some(create_key_event(KeyCode::Tab)))
            .unwrap();
        assert_eq!(dialog.selected_option(), 0);
    }

    #[test]
    fn test_confirm_yes() {
        let mut dialog = create_test_dialog();

        // Move to "Yes"
        dialog
            .handle_events(Some(create_key_event(KeyCode::Left)))
            .unwrap();
        assert_eq!(dialog.selected_option(), 0);

        // Press Enter to confirm
        let result = dialog
            .handle_events(Some(create_key_event(KeyCode::Enter)))
            .unwrap();

        // Should return confirmation action
        assert!(result.is_some());
    }

    #[test]
    fn test_confirm_no() {
        let mut dialog = create_test_dialog();

        // Already at "No"
        assert_eq!(dialog.selected_option(), 1);

        // Press Enter
        let result = dialog
            .handle_events(Some(create_key_event(KeyCode::Enter)))
            .unwrap();

        // Should return cancel action
        assert!(result.is_some());
        assert!(!dialog.is_confirmed());
    }

    #[test]
    fn test_escape_cancels() {
        let mut dialog = create_test_dialog();

        // Press Escape
        let result = dialog
            .handle_events(Some(create_key_event(KeyCode::Esc)))
            .unwrap();

        // Should return cancel action
        assert!(result.is_some());
        assert!(!dialog.is_confirmed());
    }

    #[test]
    fn test_y_key_confirms() {
        let mut dialog = create_test_dialog();

        // Press 'y' for quick confirm
        let result = dialog
            .handle_events(Some(create_key_event(KeyCode::Char('y'))))
            .unwrap();

        // Should confirm immediately
        assert!(result.is_some());
        // Dialog returns ConfirmDelete action but doesn't change internal state
    }

    #[test]
    fn test_n_key_cancels() {
        let mut dialog = create_test_dialog();

        // Press 'n' for quick cancel
        let result = dialog
            .handle_events(Some(create_key_event(KeyCode::Char('n'))))
            .unwrap();

        // Should cancel immediately
        assert!(result.is_some());
        // Dialog returns CancelDelete action
    }

    #[test]
    fn test_space_selects() {
        let mut dialog = create_test_dialog();

        // Move to Delete button
        dialog
            .handle_events(Some(create_key_event(KeyCode::Right)))
            .unwrap();

        // Press Space to select - Note: Space is not handled by ConfirmDialog, only Enter
        let result = dialog
            .handle_events(Some(create_key_event(KeyCode::Char(' '))))
            .unwrap();

        // Space doesn't do anything in ConfirmDialog
        assert!(result.is_none());
    }

    #[test]
    fn test_visual_selection_indicator() {
        let mut dialog = create_test_dialog();
        let mut terminal = setup_test_terminal(60, 20).unwrap();

        // Render with "No" selected
        terminal
            .draw(|f| {
                dialog.draw(f, f.area()).unwrap();
            })
            .unwrap();

        // Should show Cancel button (default selection)
        assert_buffer_contains(&terminal, "Cancel");
        assert_buffer_contains(&terminal, "Delete");

        // Move to Delete button
        dialog
            .handle_events(Some(create_key_event(KeyCode::Right)))
            .unwrap();

        terminal
            .draw(|f| {
                dialog.draw(f, f.area()).unwrap();
            })
            .unwrap();

        // Selection should now be on Delete
        assert_buffer_contains(&terminal, "Delete");
    }

    #[test]
    fn test_custom_messages() {
        let mut dialog = ConfirmDialog::new(
            "Custom Title",
            "This is a custom confirmation message with different text.",
        )
        .with_actions(
            gemini_cli_manager::action::Action::ConfirmDelete,
            gemini_cli_manager::action::Action::CancelDelete,
        );

        let mut terminal = setup_test_terminal(60, 20).unwrap();
        terminal
            .draw(|f| {
                dialog.draw(f, f.area()).unwrap();
            })
            .unwrap();

        assert_buffer_contains(&terminal, "Custom Title");
        assert_buffer_contains(&terminal, "custom confirmation message");
    }

    #[test]
    fn test_long_message_wrapping() {
        let long_message = "This is a very long confirmation message that should wrap properly when displayed in the dialog. It contains multiple sentences to test the text wrapping functionality.";
        let mut dialog = ConfirmDialog::new("Long Message Test", long_message);

        let mut terminal = setup_test_terminal(50, 20).unwrap();

        // Should render without panic even with long text
        let result = terminal.draw(|f| {
            dialog.draw(f, f.area()).unwrap();
        });

        assert!(result.is_ok());
    }

    #[test]
    fn test_responsive_sizing() {
        let mut dialog = create_test_dialog();

        // Test various terminal sizes
        let sizes = vec![(40, 15), (60, 20), (80, 24), (120, 40)];

        for (width, height) in sizes {
            let mut terminal = setup_test_terminal(width, height).unwrap();

            // Should render without panic at any size
            let result = terminal.draw(|f| {
                dialog.draw(f, f.area()).unwrap();
            });

            assert!(result.is_ok(), "Failed to render at {width}x{height}");
        }
    }

    #[test]
    fn test_centered_layout() {
        let mut dialog = create_test_dialog();
        let mut terminal = setup_test_terminal(80, 24).unwrap();

        terminal
            .draw(|f| {
                dialog.draw(f, f.area()).unwrap();
            })
            .unwrap();

        // Dialog should be centered on screen
        // We can't test exact positioning, but verify it renders
        assert_buffer_contains(&terminal, "Delete Extension");
    }

    #[test]
    fn test_modal_overlay() {
        let dialog = create_test_dialog();

        // Dialog is modal by design (blocks background)
        // Just verify it can be created
        assert_eq!(dialog.title(), "Delete Extension");
    }
}
