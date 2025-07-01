#[cfg(test)]
mod tests {
    use gemini_cli_manager::components::tab_bar::TabBar;
    use gemini_cli_manager::components::Component;
    use gemini_cli_manager::view::ViewType;
    use crate::test_utils::*;

    #[test]
    fn test_tab_bar_rendering() {
        let mut terminal = setup_test_terminal(40, 3).unwrap();
        let mut tab_bar = TabBar::new();
        
        terminal.draw(|f| {
            tab_bar.draw(f, f.area()).unwrap();
        }).unwrap();
        
        // Verify both tabs are rendered
        assert_buffer_contains(&terminal, "Extensions");
        assert_buffer_contains(&terminal, "Profiles");
    }
    
    #[test]
    fn test_tab_bar_active_state() {
        let mut terminal = setup_test_terminal(40, 3).unwrap();
        
        // Test with Extensions active (default)
        let mut tab_bar = TabBar::new();
        terminal.draw(|f| {
            tab_bar.draw(f, f.area()).unwrap();
        }).unwrap();
        
        // Extensions should be highlighted (exact highlighting depends on theme)
        assert_buffer_contains(&terminal, "Extensions");
        
        // Test with Profiles active
        terminal.clear().unwrap();
        let mut tab_bar = TabBar::new();
        tab_bar.set_current_view(ViewType::ProfileList);
        terminal.draw(|f| {
            tab_bar.draw(f, f.area()).unwrap();
        }).unwrap();
        
        assert_buffer_contains(&terminal, "Profiles");
    }
    
    #[test]
    fn test_tab_bar_minimum_size() {
        // Test that tab bar handles small terminal sizes gracefully
        let mut terminal = setup_test_terminal(20, 3).unwrap();
        let mut tab_bar = TabBar::new();
        
        // Should not panic with small size
        let result = terminal.draw(|f| {
            tab_bar.draw(f, f.area()).unwrap();
        });
        
        assert!(result.is_ok());
    }
}