#![allow(dead_code)]

use catppuccin::{PALETTE, Flavor, FlavorColors};
use ratatui::style::Color;
use serde::{Deserialize, Serialize};
use std::sync::Mutex;
use lazy_static::lazy_static;

/// Available theme flavours from Catppuccin
#[derive(Debug, Clone, Copy, Serialize, Deserialize, PartialEq)]
pub enum ThemeFlavour {
    Latte,     // Light theme
    Frappe,    // Light-dark theme  
    Macchiato, // Dark theme
    Mocha,     // Darker theme
}

impl ThemeFlavour {
    fn to_flavor(&self) -> &'static Flavor {
        match self {
            ThemeFlavour::Latte => &PALETTE.latte,
            ThemeFlavour::Frappe => &PALETTE.frappe,
            ThemeFlavour::Macchiato => &PALETTE.macchiato,
            ThemeFlavour::Mocha => &PALETTE.mocha,
        }
    }
}

/// A theme defines the colors used throughout the application
#[derive(Debug, Clone)]
pub struct Theme {
    pub name: String,
    flavor: &'static Flavor,
    colors: &'static FlavorColors,
}

impl Theme {
    /// Create a theme with the specified flavour
    pub fn new(flavour: ThemeFlavour) -> Self {
        let flavor = flavour.to_flavor();
        Self {
            name: format!("Catppuccin {}", flavor.name),
            flavor,
            colors: &flavor.colors,
        }
    }

    // Base colors
    pub fn background(&self) -> Color {
        self.colors.base.into()
    }
    
    pub fn surface(&self) -> Color {
        self.colors.surface0.into()
    }
    
    pub fn overlay(&self) -> Color {
        self.colors.surface1.into()
    }

    // Text colors
    pub fn text_primary(&self) -> Color {
        self.colors.text.into()
    }
    
    pub fn text_secondary(&self) -> Color {
        self.colors.subtext1.into()
    }
    
    pub fn text_muted(&self) -> Color {
        self.colors.subtext0.into()
    }
    
    pub fn text_disabled(&self) -> Color {
        self.colors.overlay0.into()
    }

    // Accent colors
    pub fn primary(&self) -> Color {
        self.colors.blue.into()
    }
    
    pub fn secondary(&self) -> Color {
        self.colors.mauve.into()
    }
    
    pub fn accent(&self) -> Color {
        self.colors.pink.into()
    }
    
    pub fn highlight(&self) -> Color {
        self.colors.yellow.into()
    }

    // Semantic colors
    pub fn success(&self) -> Color {
        self.colors.green.into()
    }
    
    pub fn warning(&self) -> Color {
        self.colors.peach.into()
    }
    
    pub fn error(&self) -> Color {
        self.colors.red.into()
    }
    
    pub fn info(&self) -> Color {
        self.colors.sky.into()
    }

    // UI element colors
    pub fn border(&self) -> Color {
        self.colors.surface2.into()
    }
    
    pub fn border_focused(&self) -> Color {
        self.colors.blue.into()
    }
    
    pub fn selection(&self) -> Color {
        self.colors.surface2.into()
    }
    
    pub fn cursor(&self) -> Color {
        self.colors.rosewater.into()
    }
}

impl Default for Theme {
    fn default() -> Self {
        // Default to Mocha (dark theme) for good contrast
        Self::new(ThemeFlavour::Mocha)
    }
}

lazy_static! {
    /// Global theme instance using Mutex for thread safety
    static ref CURRENT_THEME: Mutex<Theme> = Mutex::new(Theme::default());
}

/// Get the current theme and apply a function to it
fn with_theme<F, R>(f: F) -> R
where
    F: FnOnce(&Theme) -> R,
{
    let theme = CURRENT_THEME.lock().unwrap();
    f(&*theme)
}

/// Set the current theme
pub fn set_theme(theme: Theme) {
    let mut current = CURRENT_THEME.lock().unwrap();
    *current = theme;
}

/// Set the theme by flavour
pub fn set_flavour(flavour: ThemeFlavour) {
    set_theme(Theme::new(flavour));
}

/// Helper functions for common color needs
pub fn background() -> Color { with_theme(|t| t.background()) }
pub fn surface() -> Color { with_theme(|t| t.surface()) }
pub fn overlay() -> Color { with_theme(|t| t.overlay()) }
pub fn text_primary() -> Color { with_theme(|t| t.text_primary()) }
pub fn text_secondary() -> Color { with_theme(|t| t.text_secondary()) }
pub fn text_muted() -> Color { with_theme(|t| t.text_muted()) }
pub fn text_disabled() -> Color { with_theme(|t| t.text_disabled()) }
pub fn primary() -> Color { with_theme(|t| t.primary()) }
pub fn secondary() -> Color { with_theme(|t| t.secondary()) }
pub fn accent() -> Color { with_theme(|t| t.accent()) }
pub fn highlight() -> Color { with_theme(|t| t.highlight()) }
pub fn border() -> Color { with_theme(|t| t.border()) }
pub fn border_focused() -> Color { with_theme(|t| t.border_focused()) }
pub fn selection() -> Color { with_theme(|t| t.selection()) }
pub fn cursor() -> Color { with_theme(|t| t.cursor()) }
pub fn success() -> Color { with_theme(|t| t.success()) }
pub fn error() -> Color { with_theme(|t| t.error()) }
pub fn warning() -> Color { with_theme(|t| t.warning()) }
pub fn info() -> Color { with_theme(|t| t.info()) }