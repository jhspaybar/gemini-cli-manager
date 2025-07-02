use crate::components::settings_view::KeybindingConfig;

/// Builds dynamic help text based on current keybindings
pub struct HelpTextBuilder<'a> {
    keybindings: &'a KeybindingConfig,
}

impl<'a> HelpTextBuilder<'a> {
    pub fn new(keybindings: &'a KeybindingConfig) -> Self {
        Self { keybindings }
    }

    /// Get all keys for an action (for display in help text)
    fn get_keys(&self, action: &str) -> String {
        let keys = self.keybindings.get_keys_for_action(action);
        if keys.is_empty() {
            "?".to_string()
        } else {
            keys.join(", ")
        }
    }

    /// Build help text for a list of actions
    pub fn build(&self, actions: &[(&str, &str)]) -> String {
        let parts: Vec<String> = actions
            .iter()
            .map(|(action, label)| format!("{}: {}", self.get_keys(action), label))
            .collect();

        if parts.is_empty() {
            String::new()
        } else {
            format!(" {} ", parts.join(" | "))
        }
    }

    /// Common help text patterns
    #[allow(dead_code)]
    pub fn navigation_help(&self) -> String {
        self.build(&[
            ("up", "Up"),
            ("down", "Down"),
            ("select", "Select"),
            ("back", "Back"),
            ("quit", "Quit"),
        ])
    }

    #[allow(dead_code)]
    pub fn list_help(&self) -> String {
        self.build(&[
            ("up", "Navigate"),
            ("select", "View"),
            ("edit", "Edit"),
            ("create", "New"),
            ("delete", "Delete"),
            ("search", "Search"),
            ("quit", "Quit"),
        ])
    }

    #[allow(dead_code)]
    pub fn form_help(&self) -> String {
        self.build(&[
            ("up", "Previous field"),
            ("down", "Next field"),
            ("select", "Submit"),
            ("back", "Cancel"),
            ("quit", "Quit"),
        ])
    }
}

/// Global function to get current keybindings
pub fn get_current_keybindings() -> KeybindingConfig {
    // Try to load from settings manager
    if let Ok(manager) = crate::components::settings_view::SettingsManager::new() {
        manager.get_settings().keybindings.clone()
    } else {
        // Fall back to defaults if settings can't be loaded
        KeybindingConfig::default()
    }
}

/// Convenience function to build help text with current keybindings
pub fn build_help_text(actions: &[(&str, &str)]) -> String {
    let keybindings = get_current_keybindings();
    HelpTextBuilder::new(&keybindings).build(actions)
}
