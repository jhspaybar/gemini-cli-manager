#[cfg(test)]
mod tests {
    use crate::test_utils::*;
    use crossterm::event::{KeyCode, KeyEvent, KeyEventKind};
    use gemini_cli_manager::components::Component;
    use gemini_cli_manager::components::extension_detail::ExtensionDetail;
    use gemini_cli_manager::models::extension::McpServerConfig;
    use std::collections::HashMap;

    fn create_key_event(code: KeyCode) -> gemini_cli_manager::tui::Event {
        use crossterm::event::KeyModifiers;
        gemini_cli_manager::tui::Event::Key(KeyEvent {
            code,
            modifiers: KeyModifiers::NONE,
            kind: KeyEventKind::Press,
            state: crossterm::event::KeyEventState::NONE,
        })
    }

    fn create_test_detail() -> ExtensionDetail {
        let storage = create_test_storage();

        // Create a comprehensive test extension
        let mut mcp_servers = HashMap::new();
        mcp_servers.insert(
            "echo-server".to_string(),
            McpServerConfig {
                command: Some("node".to_string()),
                args: Some(vec!["echo-server.js".to_string()]),
                cwd: None,
                env: Some(HashMap::new()),
                trust: Some(true),
                timeout: None,
                url: None,
            },
        );

        let mut ext = ExtensionBuilder::new("Test Extension")
            .with_version("2.1.0")
            .with_description("A comprehensive test extension with MCP servers")
            .with_tags(vec!["test", "mcp", "development"])
            .build();
        ext.mcp_servers = mcp_servers;
        ext.context_file_name = Some("CONTEXT.md".to_string());
        ext.context_content =
            Some("# Extension Context\n\nThis extension provides echo functionality.".to_string());

        storage.save_extension(&ext).unwrap();

        ExtensionDetail::new(storage, ext.id.clone())
    }

    #[test]
    fn test_detail_view_rendering() {
        let mut detail = create_test_detail();
        let mut terminal = setup_test_terminal(80, 30).unwrap();

        terminal
            .draw(|f| {
                detail.draw(f, f.area()).unwrap();
            })
            .unwrap();

        // Verify all extension details are displayed
        assert_buffer_contains(&terminal, "Test Extension");
        assert_buffer_contains(&terminal, "2.1.0");
        assert_buffer_contains(&terminal, "A comprehensive test extension");
        assert_buffer_contains(&terminal, "test, mcp, development");
        assert_buffer_contains(&terminal, "echo-server");
        assert_buffer_contains(&terminal, "CONTEXT.md");
    }

    // TODO: ExtensionDetail doesn't have section navigation, only scrolling
    // #[test]
    // fn test_navigation_sections() {
    //     let mut detail = create_test_detail();
    //
    //     // Should start at overview
    //     assert_eq!(detail.current_section(), 0);
    //
    //     // Navigate down through sections
    //     detail.handle_events(Some(create_key_event(KeyCode::Down))).unwrap();
    //     assert_eq!(detail.current_section(), 1); // MCP Servers
    //
    //     detail.handle_events(Some(create_key_event(KeyCode::Down))).unwrap();
    //     assert_eq!(detail.current_section(), 2); // Context
    //
    //     // Navigate back up
    //     detail.handle_events(Some(create_key_event(KeyCode::Up))).unwrap();
    //     assert_eq!(detail.current_section(), 1);
    // }

    #[test]
    fn test_back_navigation() {
        let mut detail = create_test_detail();

        // Press Escape to go back
        let result = detail
            .handle_events(Some(create_key_event(KeyCode::Esc)))
            .unwrap();

        // Should return navigation action
        assert!(result.is_some());
    }

    #[test]
    fn test_edit_action() {
        let mut detail = create_test_detail();

        // Press 'e' to edit
        let result = detail
            .handle_events(Some(create_key_event(KeyCode::Char('e'))))
            .unwrap();

        // Should return edit action
        assert!(result.is_some());
    }

    #[test]
    fn test_delete_action() {
        let mut detail = create_test_detail();

        // Press 'd' to delete
        let result = detail
            .handle_events(Some(create_key_event(KeyCode::Char('d'))))
            .unwrap();

        // Should return delete action
        assert!(result.is_some());
    }

    // TODO: ExtensionDetail doesn't have tab navigation between sections
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
    fn test_mcp_server_details() {
        let mut detail = create_test_detail();
        let mut terminal = setup_test_terminal(80, 30).unwrap();

        terminal
            .draw(|f| {
                detail.draw(f, f.area()).unwrap();
            })
            .unwrap();

        // Verify MCP server details are shown
        assert_buffer_contains(&terminal, "Type: Command");
        assert_buffer_contains(&terminal, "node");
        assert_buffer_contains(&terminal, "Args: echo-server.js");
    }

    #[test]
    fn test_context_content_display() {
        let mut detail = create_test_detail();
        let mut terminal = setup_test_terminal(80, 30).unwrap();

        terminal
            .draw(|f| {
                detail.draw(f, f.area()).unwrap();
            })
            .unwrap();

        // Verify context content is shown
        assert_buffer_contains(&terminal, "Context File: CONTEXT.md");
        assert_buffer_contains(&terminal, "Extension Context");
        assert_buffer_contains(&terminal, "echo functionality");
    }

    #[test]
    fn test_empty_extension_handling() {
        let storage = create_test_storage();

        // Create minimal extension
        let ext = ExtensionBuilder::new("Minimal Extension").build();
        storage.save_extension(&ext).unwrap();

        let mut detail = ExtensionDetail::new(storage, ext.id);
        let mut terminal = setup_test_terminal(80, 30).unwrap();

        terminal
            .draw(|f| {
                detail.draw(f, f.area()).unwrap();
            })
            .unwrap();

        // Should handle missing optional fields gracefully by not showing them
        assert_buffer_contains(&terminal, "Minimal Extension");
        assert_buffer_contains(&terminal, "ID: minimal-extension");
        // Empty fields are not displayed, so we shouldn't see "No description" etc.
    }

    #[test]
    fn test_scroll_long_content() {
        let storage = create_test_storage();

        // Create extension with long description
        let mut ext = ExtensionBuilder::new("Long Extension")
            .with_description(&"This is a very long description. ".repeat(20))
            .build();

        // Add long context content
        ext.context_content =
            Some("# Long Context\n\n".to_string() + &"Line of content\n".repeat(50));
        storage.save_extension(&ext).unwrap();

        let mut detail = ExtensionDetail::new(storage, ext.id);

        // Navigate to context and scroll
        detail
            .handle_events(Some(create_key_event(KeyCode::Down)))
            .unwrap();
        detail
            .handle_events(Some(create_key_event(KeyCode::Down)))
            .unwrap();

        // Test scrolling
        detail
            .handle_events(Some(create_key_event(KeyCode::PageDown)))
            .unwrap();
        detail
            .handle_events(Some(create_key_event(KeyCode::PageUp)))
            .unwrap();

        // Should handle scrolling without panic
    }

    #[test]
    fn test_help_text() {
        let mut detail = create_test_detail();
        let mut terminal = setup_test_terminal(80, 30).unwrap();

        terminal
            .draw(|f| {
                detail.draw(f, f.area()).unwrap();
            })
            .unwrap();

        // Verify help text is shown
        assert_buffer_contains(&terminal, "b: Back");
        assert_buffer_contains(&terminal, "e: Edit");
        assert_buffer_contains(&terminal, "d: Delete");
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

            assert!(result.is_ok(), "Failed to render at {width}x{height}");
        }
    }
}
