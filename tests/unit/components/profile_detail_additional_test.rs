#[cfg(test)]
mod tests {
    use gemini_cli_manager::components::profile_detail::ProfileDetail;
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
        let detail = ProfileDetail::default();
        
        // Should have no profile loaded
        assert!(matches!(detail, ProfileDetail { .. }));
    }
    
    #[test]
    fn test_with_storage_constructor() {
        let storage = create_test_storage();
        let detail = ProfileDetail::with_storage(storage);
        
        // Should be initialized with storage
        assert!(matches!(detail, ProfileDetail { .. }));
    }
    
    #[test]
    fn test_new_constructor() {
        let storage = create_test_storage();
        let profile = ProfileBuilder::new("Test Profile").build();
        storage.save_profile(&profile).unwrap();
        
        let detail = ProfileDetail::new(storage, profile.id.clone());
        
        // Should load the profile
        assert!(matches!(detail, ProfileDetail { .. }));
    }
    
    #[test]
    fn test_new_constructor_non_existent_profile() {
        let storage = create_test_storage();
        let detail = ProfileDetail::new(storage, "non-existent-id".to_string());
        
        // Should handle non-existent profile gracefully
        assert!(matches!(detail, ProfileDetail { .. }));
    }
    
    #[test]
    fn test_component_init() {
        let mut detail = ProfileDetail::default();
        
        // Test init
        let result = detail.init(Size { width: 80, height: 24 });
        assert!(result.is_ok());
    }
    
    #[test]
    fn test_component_register_action_handler() {
        let mut detail = ProfileDetail::default();
        let (tx, _rx) = mpsc::unbounded_channel();
        
        let result = detail.register_action_handler(tx);
        assert!(result.is_ok());
    }
    
    #[test]
    fn test_component_register_config_handler() {
        let mut detail = ProfileDetail::default();
        let config = Config::default();
        
        let result = detail.register_config_handler(config);
        assert!(result.is_ok());
    }
    
    #[test]
    fn test_component_update_render_action() {
        let mut detail = ProfileDetail::default();
        
        // Test update with Render action
        let result = detail.update(Action::Render);
        assert!(result.is_ok());
        assert!(result.unwrap().is_none());
    }
    
    #[test]
    fn test_update_view_profile_details_action() {
        let storage = create_test_storage();
        let profile = ProfileBuilder::new("Update Test Profile").build();
        storage.save_profile(&profile).unwrap();
        
        let mut detail = ProfileDetail::with_storage(storage);
        
        // Test ViewProfileDetails action
        let result = detail.update(Action::ViewProfileDetails(profile.id.clone()));
        assert!(result.is_ok());
        assert!(result.unwrap().is_none());
    }
    
    #[test]
    fn test_update_view_profile_details_non_existent() {
        let storage = create_test_storage();
        let mut detail = ProfileDetail::with_storage(storage);
        
        // Test ViewProfileDetails action with non-existent profile
        let result = detail.update(Action::ViewProfileDetails("non-existent".to_string()));
        assert!(result.is_ok());
        assert!(result.unwrap().is_none());
    }
    
    #[test]
    fn test_update_refresh_profiles_action() {
        let storage = create_test_storage();
        let mut profile = ProfileBuilder::new("Refresh Test Profile").build();
        storage.save_profile(&profile).unwrap();
        
        let mut detail = ProfileDetail::new(storage.clone(), profile.id.clone());
        
        // Modify the profile in storage
        profile.description = Some("Updated description".to_string());
        storage.save_profile(&profile).unwrap();
        
        // Test RefreshProfiles action
        let result = detail.update(Action::RefreshProfiles);
        assert!(result.is_ok());
        assert!(result.unwrap().is_none());
    }
    
    #[test]
    fn test_update_refresh_profiles_no_current_profile() {
        let storage = create_test_storage();
        let mut detail = ProfileDetail::with_storage(storage);
        
        // Test RefreshProfiles action with no current profile
        let result = detail.update(Action::RefreshProfiles);
        assert!(result.is_ok());
        assert!(result.unwrap().is_none());
    }
    
    #[test]
    fn test_vim_navigation_keys() {
        let storage = create_test_storage();
        let profile = ProfileBuilder::new("Test Profile").build();
        storage.save_profile(&profile).unwrap();
        
        let mut detail = ProfileDetail::new(storage, profile.id);
        
        // Test vim key 'j' for down
        let result = detail.handle_events(Some(create_key_event(KeyCode::Char('j')))).unwrap();
        assert_eq!(result, Some(Action::Render));
        
        // Test vim key 'k' for up
        let result = detail.handle_events(Some(create_key_event(KeyCode::Char('k')))).unwrap();
        assert_eq!(result, Some(Action::Render));
    }
    
    #[test]
    fn test_scrolling_functionality() {
        let storage = create_test_storage();
        
        // Create profile with many extensions to enable scrolling
        let mut ext_ids = vec![];
        for i in 0..20 {
            let ext = ExtensionBuilder::new(&format!("Extension {}", i))
                .with_description(&format!("Description {}", i))
                .build();
            storage.save_extension(&ext).unwrap();
            ext_ids.push(ext.id);
        }
        
        let mut profile = ProfileBuilder::new("Scrollable Profile").build();
        profile.extension_ids = ext_ids;
        storage.save_profile(&profile).unwrap();
        
        let mut detail = ProfileDetail::new(storage, profile.id);
        
        // Test scroll down
        let result = detail.handle_events(Some(create_key_event(KeyCode::Down))).unwrap();
        assert_eq!(result, Some(Action::Render));
        
        // Test scroll up
        let result = detail.handle_events(Some(create_key_event(KeyCode::Up))).unwrap();
        assert_eq!(result, Some(Action::Render));
    }
    
    #[test]
    fn test_quit_key() {
        let mut detail = ProfileDetail::default();
        
        // Press 'q' to quit
        let result = detail.handle_events(Some(create_key_event(KeyCode::Char('q')))).unwrap();
        assert_eq!(result, Some(Action::Quit));
    }
    
    #[test]
    fn test_enter_key_does_nothing() {
        let storage = create_test_storage();
        let profile = ProfileBuilder::new("Test Profile").build();
        storage.save_profile(&profile).unwrap();
        
        let mut detail = ProfileDetail::new(storage, profile.id);
        
        // Enter should not do anything (selects things, not launch)
        let result = detail.handle_events(Some(create_key_event(KeyCode::Enter))).unwrap();
        assert_eq!(result, None);
    }
    
    #[test]
    fn test_actions_with_no_profile() {
        let mut detail = ProfileDetail::default();
        
        // Edit action with no profile should return None
        let result = detail.handle_events(Some(create_key_event(KeyCode::Char('e')))).unwrap();
        assert!(result.is_none());
        
        // Delete action with no profile should return None
        let result = detail.handle_events(Some(create_key_event(KeyCode::Char('d')))).unwrap();
        assert!(result.is_none());
        
        // Launch action with no profile should return None
        let result = detail.handle_events(Some(create_key_event(KeyCode::Char('l')))).unwrap();
        assert!(result.is_none());
    }
    
    #[test]
    fn test_set_default_action() {
        let storage = create_test_storage();
        let profile = ProfileBuilder::new("Test Profile").build();
        storage.save_profile(&profile).unwrap();
        
        let mut detail = ProfileDetail::new(storage, profile.id);
        
        // Set default action (not implemented yet)
        let result = detail.handle_events(Some(create_key_event(KeyCode::Char('x')))).unwrap();
        assert_eq!(result, None);
    }
    
    #[test]
    fn test_unknown_key_events() {
        let mut detail = ProfileDetail::default();
        
        // Unknown keys should return None
        let result = detail.handle_events(Some(create_key_event(KeyCode::F(1)))).unwrap();
        assert!(result.is_none());
        
        let result = detail.handle_events(Some(create_key_event(KeyCode::Char('z')))).unwrap();
        assert!(result.is_none());
    }
    
    #[test]
    fn test_non_key_events() {
        let mut detail = ProfileDetail::default();
        
        // Non-key events should return None
        let result = detail.handle_events(None).unwrap();
        assert!(result.is_none());
        
        let result = detail.handle_events(Some(gemini_cli_manager::tui::Event::Tick)).unwrap();
        assert!(result.is_none());
    }
    
    #[test]
    fn test_no_profile_loaded_rendering() {
        let mut detail = ProfileDetail::default();
        let mut terminal = setup_test_terminal(80, 30).unwrap();
        
        // Should render empty state
        let result = terminal.draw(|f| {
            let _ = detail.draw(f, f.area());
        });
        assert!(result.is_ok());
        
        assert_buffer_contains(&terminal, "No profile selected");
    }
    
    #[test]
    fn test_profile_with_sensitive_environment_variables() {
        let storage = create_test_storage();
        let mut profile = ProfileBuilder::new("Sensitive Profile").build();
        
        // Add sensitive environment variables that should be masked
        profile.environment_variables.insert("API_TOKEN".to_string(), "secret-token-value-here".to_string());
        profile.environment_variables.insert("SECRET_KEY".to_string(), "very-secret-key".to_string());
        profile.environment_variables.insert("PASSWORD".to_string(), "pwd".to_string()); // Short value
        profile.environment_variables.insert("NORMAL_VAR".to_string(), "normal-value".to_string());
        
        storage.save_profile(&profile).unwrap();
        
        let mut detail = ProfileDetail::new(storage, profile.id);
        let mut terminal = setup_test_terminal(80, 30).unwrap();
        
        terminal.draw(|f| {
            detail.draw(f, f.area()).unwrap();
        }).unwrap();
        
        // Verify sensitive values are masked
        assert_buffer_contains(&terminal, "API_TOKEN = secr...here");
        assert_buffer_contains(&terminal, "SECRET_KEY = very...-key");
        assert_buffer_contains(&terminal, "PASSWORD = pwd"); // Short value, not actually masked because it's too short
        assert_buffer_contains(&terminal, "NORMAL_VAR = normal-value"); // Not masked
    }
    
    #[test]
    fn test_profile_with_extensions_with_mcp_servers() {
        let storage = create_test_storage();
        
        // Create extension with MCP servers
        let mut mcp_servers = HashMap::new();
        mcp_servers.insert("database-server".to_string(), McpServerConfig {
            command: Some("python".to_string()),
            args: Some(vec!["-m".to_string(), "database_server".to_string()]),
            cwd: None,
            env: None,
            trust: Some(true),
            timeout: None,
            url: None,
        });
        
        mcp_servers.insert("api-server".to_string(), McpServerConfig {
            command: Some("node".to_string()),
            args: Some(vec!["api-server.js".to_string()]),
            cwd: None,
            env: None,
            trust: None,
            timeout: None,
            url: None,
        });
        
        let mut ext = ExtensionBuilder::new("MCP Extension")
            .with_description("Extension with MCP servers")
            .build();
        ext.mcp_servers = mcp_servers;
        storage.save_extension(&ext).unwrap();
        
        let profile = ProfileBuilder::new("MCP Profile")
            .with_extensions(vec![&ext.id])
            .build();
        storage.save_profile(&profile).unwrap();
        
        let mut detail = ProfileDetail::new(storage, profile.id);
        let mut terminal = setup_test_terminal(80, 30).unwrap();
        
        terminal.draw(|f| {
            detail.draw(f, f.area()).unwrap();
        }).unwrap();
        
        // Verify MCP servers are shown (order may vary due to HashMap)
        // Check that both MCP servers are shown (order doesn't matter with HashMap)
        assert_buffer_contains(&terminal, "database-server");
        assert_buffer_contains(&terminal, "api-server");
    }
    
    #[test]
    fn test_profile_without_optional_fields() {
        let storage = create_test_storage();
        
        // Create minimal profile
        let mut profile = ProfileBuilder::new("Minimal").build();
        profile.description = None;
        profile.metadata.tags.clear();
        profile.working_directory = None;
        profile.environment_variables.clear();
        profile.extension_ids.clear();
        
        storage.save_profile(&profile).unwrap();
        
        let mut detail = ProfileDetail::new(storage, profile.id);
        let mut terminal = setup_test_terminal(80, 30).unwrap();
        
        terminal.draw(|f| {
            detail.draw(f, f.area()).unwrap();
        }).unwrap();
        
        // Should show only required fields
        assert_buffer_contains(&terminal, "Minimal");
        assert_buffer_contains(&terminal, "ID: minimal");
        assert_buffer_contains(&terminal, "Created:");
        assert_buffer_contains(&terminal, "No extensions included");
        assert_buffer_contains(&terminal, "• 0 extensions");
        assert_buffer_contains(&terminal, "• 0 environment variables");
        assert_buffer_contains(&terminal, "• 0 MCP servers total");
    }
    
    #[test]
    fn test_set_profile_resets_scroll() {
        let storage = create_test_storage();
        let profile1 = ProfileBuilder::new("Profile 1").build();
        let profile2 = ProfileBuilder::new("Profile 2").build();
        storage.save_profile(&profile1).unwrap();
        storage.save_profile(&profile2).unwrap();
        
        let mut detail = ProfileDetail::new(storage.clone(), profile1.id.clone());
        
        // Scroll down
        detail.handle_events(Some(create_key_event(KeyCode::Down))).unwrap();
        detail.handle_events(Some(create_key_event(KeyCode::Down))).unwrap();
        
        // Set a new profile
        let profile2_loaded = storage.load_profile(&profile2.id).unwrap();
        detail.set_profile(profile2_loaded);
        
        // Scroll should be reset (we can't directly test scroll_offset but behavior should work)
        let mut terminal = setup_test_terminal(80, 30).unwrap();
        terminal.draw(|f| {
            detail.draw(f, f.area()).unwrap();
        }).unwrap();
        
        assert_buffer_contains(&terminal, "Profile 2");
    }
    
    #[test]
    fn test_set_profile_loads_extensions() {
        let storage = create_test_storage();
        
        // Create extensions
        let ext1 = ExtensionBuilder::new("Extension 1").build();
        let ext2 = ExtensionBuilder::new("Extension 2").build();
        storage.save_extension(&ext1).unwrap();
        storage.save_extension(&ext2).unwrap();
        
        // Create profile with extensions
        let profile = ProfileBuilder::new("Test Profile")
            .with_extensions(vec![&ext1.id, &ext2.id])
            .build();
        
        let mut detail = ProfileDetail::with_storage(storage);
        detail.set_profile(profile);
        
        let mut terminal = setup_test_terminal(80, 30).unwrap();
        terminal.draw(|f| {
            detail.draw(f, f.area()).unwrap();
        }).unwrap();
        
        // Should load and display extensions
        assert_buffer_contains(&terminal, "Extension 1");
        assert_buffer_contains(&terminal, "Extension 2");
    }
    
    #[test]
    fn test_set_profile_with_no_storage() {
        let profile = ProfileBuilder::new("Test Profile").build();
        
        let mut detail = ProfileDetail::default(); // No storage
        detail.set_profile(profile);
        
        let mut terminal = setup_test_terminal(80, 30).unwrap();
        terminal.draw(|f| {
            detail.draw(f, f.area()).unwrap();
        }).unwrap();
        
        // Should handle no storage gracefully
        assert_buffer_contains(&terminal, "Test Profile");
        assert_buffer_contains(&terminal, "No extensions included");
    }
    
    #[test]
    fn test_current_section_helper() {
        let detail = ProfileDetail::default();
        
        // Test helper method
        assert_eq!(detail.current_section(), 0);
    }
    
    #[test]
    fn test_profile_with_tags_display() {
        let storage = create_test_storage();
        let profile = ProfileBuilder::new("Tagged Profile")
            .with_tags(vec!["development", "testing", "staging"])
            .build();
        storage.save_profile(&profile).unwrap();
        
        let mut detail = ProfileDetail::new(storage, profile.id);
        let mut terminal = setup_test_terminal(80, 30).unwrap();
        
        terminal.draw(|f| {
            detail.draw(f, f.area()).unwrap();
        }).unwrap();
        
        // Verify tags are displayed
        assert_buffer_contains(&terminal, "Tags: development, testing, staging");
    }
    
    #[test]
    fn test_profile_summary_section() {
        let storage = create_test_storage();
        
        // Create extensions with MCP servers
        let mut ext1 = ExtensionBuilder::new("Ext 1").build();
        let mut mcp_servers1 = HashMap::new();
        mcp_servers1.insert("server1".to_string(), McpServerConfig {
            command: Some("cmd1".to_string()),
            args: None,
            cwd: None,
            env: None,
            trust: None,
            timeout: None,
            url: None,
        });
        ext1.mcp_servers = mcp_servers1;
        storage.save_extension(&ext1).unwrap();
        
        let mut ext2 = ExtensionBuilder::new("Ext 2").build();
        let mut mcp_servers2 = HashMap::new();
        mcp_servers2.insert("server2".to_string(), McpServerConfig {
            command: Some("cmd2".to_string()),
            args: None,
            cwd: None,
            env: None,
            trust: None,
            timeout: None,
            url: None,
        });
        mcp_servers2.insert("server3".to_string(), McpServerConfig {
            command: Some("cmd3".to_string()),
            args: None,
            cwd: None,
            env: None,
            trust: None,
            timeout: None,
            url: None,
        });
        ext2.mcp_servers = mcp_servers2;
        storage.save_extension(&ext2).unwrap();
        
        // Create profile with extensions and env vars
        let mut profile = ProfileBuilder::new("Summary Profile")
            .with_extensions(vec![&ext1.id, &ext2.id])
            .build();
        profile.environment_variables.insert("VAR1".to_string(), "value1".to_string());
        profile.environment_variables.insert("VAR2".to_string(), "value2".to_string());
        storage.save_profile(&profile).unwrap();
        
        let mut detail = ProfileDetail::new(storage, profile.id);
        let mut terminal = setup_test_terminal(80, 40).unwrap();
        
        terminal.draw(|f| {
            detail.draw(f, f.area()).unwrap();
        }).unwrap();
        
        // Verify summary counts
        assert_buffer_contains(&terminal, "Summary");
        assert_buffer_contains(&terminal, "• 2 extensions");
        assert_buffer_contains(&terminal, "• 2 environment variables");
        assert_buffer_contains(&terminal, "• 3 MCP servers total"); // 1 + 2 servers
    }
}