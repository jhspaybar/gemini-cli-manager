use crate::components::settings_view::UserSettings;
use crossterm::event::KeyEvent;
use std::sync::{Arc, RwLock};

pub struct KeybindingManager {
    settings: Arc<RwLock<UserSettings>>,
}

impl KeybindingManager {
    pub fn new(settings: Arc<RwLock<UserSettings>>) -> Self {
        Self { settings }
    }

    /// Check if a key event matches a specific action
    pub fn matches(&self, key: &KeyEvent, action: &str) -> bool {
        if let Ok(settings_lock) = self.settings.read() {
            let configured_keys = settings_lock.keybindings.get_keys_for_action(action);
            let key_str = format_key_event(key);
            configured_keys.contains(&key_str)
        } else {
            false
        }
    }

    /// Build help text with current keybindings
    pub fn build_help_text(&self, items: &[(&str, &str)]) -> String {
        if let Ok(settings_lock) = self.settings.read() {
            let keybindings = &settings_lock.keybindings;

            items
                .iter()
                .filter_map(|(action, description)| {
                    let keys = keybindings.get_keys_for_action(action);
                    if !keys.is_empty() {
                        let key_str = keys.join(", ");
                        Some(format!("{key_str}: {description}"))
                    } else {
                        None
                    }
                })
                .collect::<Vec<_>>()
                .join(" | ")
        } else {
            // Fallback if settings can't be read
            items
                .iter()
                .map(|(key, desc)| format!("{key}: {desc}"))
                .collect::<Vec<_>>()
                .join(" | ")
        }
    }
}

fn format_key_event(key: &KeyEvent) -> String {
    use crossterm::event::{KeyCode, KeyModifiers};

    let mut parts = Vec::new();

    // Add modifiers
    if key.modifiers.contains(KeyModifiers::CONTROL) {
        parts.push("Ctrl");
    }
    if key.modifiers.contains(KeyModifiers::ALT) {
        parts.push("Alt");
    }
    if key.modifiers.contains(KeyModifiers::SHIFT) {
        parts.push("Shift");
    }

    // Add the key itself
    let key_str = match key.code {
        KeyCode::Char(c) => match c {
            ' ' => "Space".to_string(),
            c => {
                if key.modifiers.contains(KeyModifiers::SHIFT) {
                    c.to_uppercase().to_string()
                } else {
                    c.to_string()
                }
            }
        },
        KeyCode::F(n) => format!("F{n}"),
        KeyCode::Up => "Up".to_string(),
        KeyCode::Down => "Down".to_string(),
        KeyCode::Left => "Left".to_string(),
        KeyCode::Right => "Right".to_string(),
        KeyCode::Enter => "Enter".to_string(),
        KeyCode::Tab => "Tab".to_string(),
        KeyCode::BackTab => "BackTab".to_string(),
        KeyCode::Backspace => "Backspace".to_string(),
        KeyCode::Delete => "Delete".to_string(),
        KeyCode::Insert => "Insert".to_string(),
        KeyCode::Home => "Home".to_string(),
        KeyCode::End => "End".to_string(),
        KeyCode::PageUp => "PageUp".to_string(),
        KeyCode::PageDown => "PageDown".to_string(),
        KeyCode::Esc => "Esc".to_string(),
        _ => return "Unknown".to_string(),
    };

    parts.push(&key_str);

    if parts.len() > 1 && !parts[0].starts_with('F') {
        parts.join("+")
    } else {
        key_str
    }
}
