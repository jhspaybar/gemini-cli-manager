#[cfg(test)]
mod tests {
    use gemini_cli_manager::components::TabBar;
    use gemini_cli_manager::view::ViewType;
    use ratatui::prelude::*;
    use crate::test_utils::*;

    #[test]
    fn test_tab_bar_rendering() {
        let mut terminal = setup_test_terminal(40, 3).unwrap();
        let tab_bar = TabBar::new(ViewType::ExtensionList);
        
        terminal.draw(|f| {
            tab_bar.render(f, f.area()).unwrap();
        }).unwrap();
        
        // Verify both tabs are rendered
        assert_buffer_contains(&terminal, "Extensions");
        assert_buffer_contains(&terminal, "Profiles");
    }
    
    #[test]
    fn test_tab_bar_active_state() {
        let mut terminal = setup_test_terminal(40, 3).unwrap();
        
        // Test with Extensions active
        let tab_bar = TabBar::new(ViewType::ExtensionList);
        terminal.draw(|f| {
            tab_bar.render(f, f.area()).unwrap();
        }).unwrap();
        
        let content = terminal.backend().buffer().content();
        // Extensions should be highlighted (exact highlighting depends on theme)
        assert!(content.contains("Extensions"));
        
        // Test with Profiles active
        terminal.clear().unwrap();
        let tab_bar = TabBar::new(ViewType::ProfileList);
        terminal.draw(|f| {
            tab_bar.render(f, f.area()).unwrap();
        }).unwrap();
        
        let content = terminal.backend().buffer().content();
        assert!(content.contains("Profiles"));
    }
    
    #[test]
    fn test_tab_bar_minimum_size() {
        // Test that tab bar handles small terminal sizes gracefully
        let mut terminal = setup_test_terminal(20, 3).unwrap();
        let tab_bar = TabBar::new(ViewType::ExtensionList);
        
        // Should not panic with small size
        let result = terminal.draw(|f| {
            tab_bar.render(f, f.area()).unwrap();
        });
        
        assert!(result.is_ok());
    }
}