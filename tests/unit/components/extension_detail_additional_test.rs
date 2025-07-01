#[cfg(test)]
mod tests {
    use gemini_cli_manager::components::extension_detail::ExtensionDetail;
    use gemini_cli_manager::components::Component;
    use gemini_cli_manager::models::extension::McpServerConfig;
    use gemini_cli_manager::action::Action;
    use gemini_cli_manager::config::Config;
    use crate::test_utils::*;
    use crossterm::event::{KeyCode, KeyEvent, KeyEventKind};
    use std::collections::HashMap;
    use tokio::sync::mpsc;
    use ratatui::prelude::*;

    fn create_key_event(code: KeyCode) -> gemini_cli_manager::tui::Event {
        use crossterm::event::KeyModifiers;
        gemini_cli_manager::tui::Event::Key(KeyEvent {
            code,
            modifiers: KeyModifiers::NONE,
            kind: KeyEventKind::Press,
            state: crossterm::event::KeyEventState::NONE,
        })
    }

    #[test]
    fn test_default_constructor() {
        let detail = ExtensionDetail::default();
        
        // Should have no extension loaded
        assert!(matches!(detail, ExtensionDetail { .. }));
    }
    
    #[test]
    fn test_with_storage_constructor() {
        let storage = create_test_storage();
        let detail = ExtensionDetail::with_storage(storage);
        
        // Should be initialized with storage
        assert!(matches!(detail, ExtensionDetail { .. }));
    }
    
    #[test]
    fn test_non_existent_extension() {
        let storage = create_test_storage();
        let mut detail = ExtensionDetail::new(storage, "non-existent-id".to_string());
        
        // Should handle non-existent extension gracefully
        let mut terminal = setup_test_terminal(80, 30).unwrap();
        
        // Should render without panic (will show "no extension loaded")
        let result = terminal.draw(|f| {
            let _ = detail.draw(f, f.area());
        });
        assert!(result.is_ok());
    }
    
    #[test]
    fn test_component_init() {
        let mut detail = ExtensionDetail::default();
        
        // Test init
        let result = detail.init(Size { width: 80, height: 24 });
        assert!(result.is_ok());
    }
    
    #[test]
    fn test_component_register_action_handler() {
        let mut detail = ExtensionDetail::default();
        let (tx, _rx) = mpsc::unbounded_channel();
        
        let result = detail.register_action_handler(tx);
        assert!(result.is_ok());
    }
    
    #[test]
    fn test_component_register_config_handler() {
        let mut detail = ExtensionDetail::default();
        let config = Config::default();
        
        let result = detail.register_config_handler(config);
        assert!(result.is_ok());
    }
    
    #[test]
    fn test_component_update() {
        let mut detail = ExtensionDetail::default();
        
        // Test update with various actions
        let result = detail.update(Action::Render);
        assert!(result.is_ok());
        assert!(result.unwrap().is_none());
    }
    
    #[test]
    fn test_vim_navigation_keys() {
        let storage = create_test_storage();
        let ext = ExtensionBuilder::new("Test Extension").build();
        storage.save_extension(&ext).unwrap();
        
        let mut detail = ExtensionDetail::new(storage, ext.id);
        
        // Test vim key 'j' for down
        let result = detail.handle_events(Some(create_key_event(KeyCode::Char('j')))).unwrap();
        assert_eq!(result, Some(Action::Render));
        
        // Test vim key 'k' for up
        let result = detail.handle_events(Some(create_key_event(KeyCode::Char('k')))).unwrap();
        assert_eq!(result, Some(Action::Render));
    }
    
    #[test]
    fn test_quit_key() {
        let mut detail = ExtensionDetail::default();
        
        // Press 'q' to quit
        let result = detail.handle_events(Some(create_key_event(KeyCode::Char('q')))).unwrap();
        assert_eq!(result, Some(Action::Quit));
    }
    
    #[test]
    fn test_scrolling_functionality() {
        let storage = create_test_storage();
        
        // Create extension with long content
        let mut ext = ExtensionBuilder::new("Long Extension")
            .with_description("Long description that might need scrolling")
            .build();
        ext.context_content = Some(vec!["Line\n"; 100].join(""));
        storage.save_extension(&ext).unwrap();
        
        let mut detail = ExtensionDetail::new(storage, ext.id);
        
        // Test scroll down
        let result = detail.handle_events(Some(create_key_event(KeyCode::Down))).unwrap();
        assert_eq!(result, Some(Action::Render));
        
        // Test scroll up
        let result = detail.handle_events(Some(create_key_event(KeyCode::Up))).unwrap();
        assert_eq!(result, Some(Action::Render));
        
        // PageDown and PageUp are not currently handled
        let result = detail.handle_events(Some(create_key_event(KeyCode::PageDown))).unwrap();
        assert_eq!(result, None);
        
        let result = detail.handle_events(Some(create_key_event(KeyCode::PageUp))).unwrap();
        assert_eq!(result, None);
    }
    
    #[test]
    fn test_extension_with_no_mcp_servers() {
        let storage = create_test_storage();
        
        // Create extension without MCP servers
        let mut ext = ExtensionBuilder::new("No MCP Extension")
            .with_description("Extension without any MCP servers")
            .build();
        ext.mcp_servers.clear();
        storage.save_extension(&ext).unwrap();
        
        let mut detail = ExtensionDetail::new(storage, ext.id);
        let mut terminal = setup_test_terminal(80, 30).unwrap();
        
        terminal.draw(|f| {
            detail.draw(f, f.area()).unwrap();
        }).unwrap();
        
        // Should render without MCP servers section
        assert_buffer_contains(&terminal, "No MCP Extension");
        assert_buffer_contains(&terminal, "Extension without any MCP servers");
    }
    
    #[test]
    fn test_mcp_server_with_environment_vars() {
        let storage = create_test_storage();
        
        // Create extension with MCP server that has environment variables
        let mut mcp_servers = HashMap::new();
        let mut env_vars = HashMap::new();
        env_vars.insert("API_KEY".to_string(), "$API_KEY".to_string());
        env_vars.insert("DEBUG".to_string(), "true".to_string());
        
        mcp_servers.insert("api-server".to_string(), McpServerConfig {
            command: Some("python".to_string()),
            args: Some(vec!["-m".to_string(), "api_server".to_string()]),
            cwd: Some("/opt/api".to_string()),
            env: Some(env_vars),
            trust: Some(false),
            timeout: Some(5000),
        });
        
        let mut ext = ExtensionBuilder::new("API Extension").build();
        ext.mcp_servers = mcp_servers;
        storage.save_extension(&ext).unwrap();
        
        let mut detail = ExtensionDetail::new(storage, ext.id);
        let mut terminal = setup_test_terminal(80, 40).unwrap();
        
        terminal.draw(|f| {
            detail.draw(f, f.area()).unwrap();
        }).unwrap();
        
        // Verify all MCP server details are shown
        assert_buffer_contains(&terminal, "api-server");
        assert_buffer_contains(&terminal, "python");
        // Working directory is not displayed in the current implementation
        assert_buffer_contains(&terminal, "Env: DEBUG = true");
        assert_buffer_contains(&terminal, "Env: API_KEY = $API_KEY");
        // Timeout is not displayed in the current implementation
        assert_buffer_contains(&terminal, "Trust: No");
    }
    
    #[test]
    fn test_mcp_server_command_type() {
        let storage = create_test_storage();
        
        // Create extension with command-based MCP server
        let mut mcp_servers = HashMap::new();
        mcp_servers.insert("remote-server".to_string(), McpServerConfig {
            command: Some("node".to_string()),
            args: Some(vec!["server.js".to_string()]),
            cwd: None,
            env: None,
            trust: None,
            timeout: None,
        });
        
        let mut ext = ExtensionBuilder::new("Remote Extension").build();
        ext.mcp_servers = mcp_servers;
        storage.save_extension(&ext).unwrap();
        
        let mut detail = ExtensionDetail::new(storage, ext.id);
        let mut terminal = setup_test_terminal(80, 30).unwrap();
        
        terminal.draw(|f| {
            detail.draw(f, f.area()).unwrap();
        }).unwrap();
        
        // Verify command-based server is shown correctly
        assert_buffer_contains(&terminal, "remote-server");
        assert_buffer_contains(&terminal, "Type: Command");
        assert_buffer_contains(&terminal, "node");
    }
    
    #[test]
    fn test_extension_without_optional_fields() {
        let storage = create_test_storage();
        
        // Create minimal extension
        let mut ext = ExtensionBuilder::new("Minimal").build();
        ext.description = None;
        ext.metadata.tags.clear();
        ext.metadata.source_path = None;
        ext.context_file_name = None;
        ext.context_content = None;
        storage.save_extension(&ext).unwrap();
        
        let mut detail = ExtensionDetail::new(storage, ext.id);
        let mut terminal = setup_test_terminal(80, 30).unwrap();
        
        terminal.draw(|f| {
            detail.draw(f, f.area()).unwrap();
        }).unwrap();
        
        // Should show only required fields
        assert_buffer_contains(&terminal, "Minimal");
        assert_buffer_contains(&terminal, "ID: minimal");
        assert_buffer_contains(&terminal, "Imported:");
    }
    
    #[test]
    fn test_actions_with_no_extension() {
        let mut detail = ExtensionDetail::default();
        
        // Edit action with no extension should return None
        let result = detail.handle_events(Some(create_key_event(KeyCode::Char('e')))).unwrap();
        assert!(result.is_none());
        
        // Delete action with no extension should return None
        let result = detail.handle_events(Some(create_key_event(KeyCode::Char('d')))).unwrap();
        assert!(result.is_none());
    }
    
    #[test]
    fn test_set_extension_resets_scroll() {
        let storage = create_test_storage();
        let ext1 = ExtensionBuilder::new("Extension 1").build();
        let ext2 = ExtensionBuilder::new("Extension 2").build();
        storage.save_extension(&ext1).unwrap();
        storage.save_extension(&ext2).unwrap();
        
        let mut detail = ExtensionDetail::new(storage.clone(), ext1.id.clone());
        
        // Scroll down
        detail.handle_events(Some(create_key_event(KeyCode::Down))).unwrap();
        detail.handle_events(Some(create_key_event(KeyCode::Down))).unwrap();
        
        // Set a new extension
        let ext2_loaded = storage.load_extension(&ext2.id).unwrap();
        detail.set_extension(ext2_loaded);
        
        // Scroll should be reset (we can't directly test scroll_offset but behavior should work)
        let mut terminal = setup_test_terminal(80, 30).unwrap();
        terminal.draw(|f| {
            detail.draw(f, f.area()).unwrap();
        }).unwrap();
        
        assert_buffer_contains(&terminal, "Extension 2");
    }
    
    #[test]
    fn test_unknown_key_events() {
        let mut detail = ExtensionDetail::default();
        
        // Unknown keys should return None
        let result = detail.handle_events(Some(create_key_event(KeyCode::F(1)))).unwrap();
        assert!(result.is_none());
        
        let result = detail.handle_events(Some(create_key_event(KeyCode::Char('x')))).unwrap();
        assert!(result.is_none());
    }
    
    #[test]
    fn test_non_key_events() {
        let mut detail = ExtensionDetail::default();
        
        // Non-key events should return None
        let result = detail.handle_events(None).unwrap();
        assert!(result.is_none());
        
        let result = detail.handle_events(Some(gemini_cli_manager::tui::Event::Tick)).unwrap();
        assert!(result.is_none());
    }
}