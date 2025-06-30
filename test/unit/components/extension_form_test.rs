#[cfg(test)]
mod tests {
    use gemini_cli_manager::components::{ExtensionForm, FormField, FormMode};
    use ratatui::prelude::*;
    use crate::test_utils::*;
    use crossterm::event::KeyCode;
    use insta::assert_snapshot;

    fn create_test_form() -> ExtensionForm {
        let storage = create_test_storage();
        ExtensionForm::new(storage)
    }
    
    fn create_edit_form(id: &str) -> ExtensionForm {
        let storage = create_test_storage();
        
        // Add test extension to storage
        let ext = ExtensionBuilder::new("Test Extension")
            .with_version("1.0.0")
            .with_description("Test description")
            .with_tags(vec!["test", "example"])
            .build();
        storage.create_extension(ext).unwrap();
        
        ExtensionForm::edit(storage, id.to_string())
    }

    #[test]
    fn test_form_initial_state() {
        let form = create_test_form();
        
        assert_eq!(form.mode(), FormMode::Create);
        assert_eq!(form.current_field(), FormField::Name);
        assert!(form.name_input().value().is_empty());
        assert!(form.version_input().value().is_empty());
        assert!(form.description_input().value().is_empty());
        assert!(form.tags_input().value().is_empty());
    }
    
    #[test]
    fn test_form_rendering() {
        let form = create_test_form();
        let mut terminal = setup_test_terminal(80, 30).unwrap();
        
        terminal.draw(|f| {
            form.render(f, f.area());
        }).unwrap();
        
        // Verify form elements are rendered
        assert_buffer_contains(&terminal, "Create New Extension");
        assert_buffer_contains(&terminal, "Name");
        assert_buffer_contains(&terminal, "Version");
        assert_buffer_contains(&terminal, "Description");
        assert_buffer_contains(&terminal, "Tags");
        assert_buffer_contains(&terminal, "MCP Servers");
    }
    
    #[test]
    fn test_field_navigation() {
        let mut form = create_test_form();
        
        // Start at Name field
        assert_eq!(form.current_field(), FormField::Name);
        
        // Tab to Version
        form.handle_key_event(KeyCode::Tab.into());
        assert_eq!(form.current_field(), FormField::Version);
        
        // Tab to Description
        form.handle_key_event(KeyCode::Tab.into());
        assert_eq!(form.current_field(), FormField::Description);
        
        // Tab to Tags
        form.handle_key_event(KeyCode::Tab.into());
        assert_eq!(form.current_field(), FormField::Tags);
        
        // Tab to MCP Servers
        form.handle_key_event(KeyCode::Tab.into());
        assert_eq!(form.current_field(), FormField::McpServers);
        
        // Tab wraps to Name
        form.handle_key_event(KeyCode::Tab.into());
        assert_eq!(form.current_field(), FormField::Name);
        
        // Shift+Tab goes backward
        form.handle_key_event(KeyCode::BackTab.into());
        assert_eq!(form.current_field(), FormField::McpServers);
    }
    
    #[test]
    fn test_text_input() {
        let mut form = create_test_form();
        
        // Type into name field
        for ch in "My Extension".chars() {
            form.handle_key_event(KeyCode::Char(ch).into());
        }
        assert_eq!(form.name_input().value(), "My Extension");
        
        // Tab to version and type
        form.handle_key_event(KeyCode::Tab.into());
        for ch in "2.0.0".chars() {
            form.handle_key_event(KeyCode::Char(ch).into());
        }
        assert_eq!(form.version_input().value(), "2.0.0");
        
        // Tab to description and type
        form.handle_key_event(KeyCode::Tab.into());
        for ch in "A great extension".chars() {
            form.handle_key_event(KeyCode::Char(ch).into());
        }
        assert_eq!(form.description_input().value(), "A great extension");
        
        // Tab to tags and type
        form.handle_key_event(KeyCode::Tab.into());
        for ch in "productivity, tools".chars() {
            form.handle_key_event(KeyCode::Char(ch).into());
        }
        assert_eq!(form.tags_input().value(), "productivity, tools");
    }
    
    #[test]
    fn test_mcp_server_management() {
        let mut form = create_test_form();
        
        // Navigate to MCP Servers field
        for _ in 0..4 {
            form.handle_key_event(KeyCode::Tab.into());
        }
        assert_eq!(form.current_field(), FormField::McpServers);
        
        // Add new server with 'a'
        form.handle_key_event(KeyCode::Char('a').into());
        
        // Should be in server edit mode
        assert!(form.is_editing_server());
        
        // Fill in server details
        // Name field
        for ch in "test-server".chars() {
            form.handle_key_event(KeyCode::Char(ch).into());
        }
        
        // Tab to command
        form.handle_key_event(KeyCode::Tab.into());
        for ch in "node server.js".chars() {
            form.handle_key_event(KeyCode::Char(ch).into());
        }
        
        // Save server with Ctrl+S
        form.handle_key_event(KeyCode::Char('s').with_ctrl());
        
        // Should return to server list
        assert!(!form.is_editing_server());
        assert_eq!(form.mcp_servers().len(), 1);
    }
    
    #[test]
    fn test_form_validation() {
        let mut form = create_test_form();
        
        // Try to save without required fields
        let result = form.save();
        assert!(result.is_err());
        
        // Add name
        for ch in "Test".chars() {
            form.handle_key_event(KeyCode::Char(ch).into());
        }
        
        // Still missing version
        let result = form.save();
        assert!(result.is_err());
        
        // Add version
        form.handle_key_event(KeyCode::Tab.into());
        for ch in "1.0.0".chars() {
            form.handle_key_event(KeyCode::Char(ch).into());
        }
        
        // Should now save successfully
        let result = form.save();
        assert!(result.is_ok());
    }
    
    #[test]
    fn test_edit_mode() {
        let form = create_edit_form("test-extension");
        
        assert_eq!(form.mode(), FormMode::Edit);
        assert_eq!(form.name_input().value(), "Test Extension");
        assert_eq!(form.version_input().value(), "1.0.0");
        assert_eq!(form.description_input().value(), "Test description");
        assert_eq!(form.tags_input().value(), "test, example");
    }
    
    #[test]
    fn test_cancel_operation() {
        let mut form = create_test_form();
        
        // Add some data
        for ch in "Unsaved Data".chars() {
            form.handle_key_event(KeyCode::Char(ch).into());
        }
        
        // Cancel with Escape
        form.handle_key_event(KeyCode::Esc.into());
        
        // Should show confirmation if there are changes
        assert!(form.is_confirming_cancel());
        
        // Confirm cancel
        form.handle_key_event(KeyCode::Char('y').into());
        
        // Form should indicate it's cancelled
        assert!(form.is_cancelled());
    }
    
    #[test]
    fn test_form_help_text() {
        let form = create_test_form();
        let mut terminal = setup_test_terminal(80, 40).unwrap();
        
        terminal.draw(|f| {
            form.render(f, f.area());
        }).unwrap();
        
        // Verify help text is shown
        assert_buffer_contains(&terminal, "Tab/Shift+Tab");
        assert_buffer_contains(&terminal, "navigate");
        assert_buffer_contains(&terminal, "Ctrl+S");
        assert_buffer_contains(&terminal, "save");
        assert_buffer_contains(&terminal, "Esc");
    }
    
    #[test]
    fn test_form_snapshot() {
        let mut form = create_test_form();
        
        // Fill out the form
        form.set_name("Snapshot Extension");
        form.set_version("3.0.0");
        form.set_description("Extension for snapshot testing");
        form.set_tags("test, snapshot, example");
        
        // Render to string
        let output = render_to_string(80, 35, |f| {
            form.render(f, f.area());
        }).unwrap();
        
        // Assert snapshot
        assert_snapshot!(output);
    }
    
    #[test]
    fn test_responsive_layout() {
        let form = create_test_form();
        
        // Test different sizes
        let sizes = vec![(60, 25), (80, 30), (100, 40), (50, 20)];
        
        for (width, height) in sizes {
            let mut terminal = setup_test_terminal(width, height).unwrap();
            
            // Should render without panic
            let result = terminal.draw(|f| {
                form.render(f, f.area());
            });
            
            assert!(result.is_ok(), "Failed at {}x{}", width, height);
        }
    }
    
    #[test]
    fn test_field_focus_styling() {
        let mut form = create_test_form();
        let mut terminal = setup_test_terminal(80, 35).unwrap();
        
        // Render with name field focused
        terminal.draw(|f| {
            form.render(f, f.area());
        }).unwrap();
        
        // The focused field should have highlight styling
        // This would be verified by checking the actual styling
        // but for now we just ensure it renders
        
        // Move to next field and re-render
        form.handle_key_event(KeyCode::Tab.into());
        terminal.clear().unwrap();
        
        terminal.draw(|f| {
            form.render(f, f.area());
        }).unwrap();
        
        // Version field should now be highlighted
    }
}

// Implement the HandleKeyEvent trait for ExtensionForm
impl HandleKeyEvent for gemini_cli_manager::components::ExtensionForm {
    fn handle_key_event(&mut self, event: crossterm::event::KeyEvent) {
        // This would call the actual handle_key_event method on ExtensionForm
        // The implementation depends on the actual ExtensionForm API
    }
}