use gemini_cli_manager::theme::{self, ThemeFlavour};
use ratatui::style::{Color, Style};

/// Test helpers for theme validation and contrast checking
pub struct ThemeTestHelper;

impl ThemeTestHelper {
    /// Get luminance of a color (approximation for terminal colors)
    pub fn get_luminance(color: Color) -> f32 {
        match color {
            Color::Rgb(r, g, b) => {
                // Use relative luminance formula
                let r = r as f32 / 255.0;
                let g = g as f32 / 255.0;
                let b = b as f32 / 255.0;
                0.2126 * r + 0.7152 * g + 0.0722 * b
            }
            // For indexed colors, use approximations
            Color::Black => 0.0,
            Color::White => 1.0,
            Color::Gray | Color::DarkGray => 0.5,
            Color::Red | Color::LightRed => 0.3,
            Color::Green | Color::LightGreen => 0.4,
            Color::Blue | Color::LightBlue => 0.2,
            Color::Yellow | Color::LightYellow => 0.7,
            Color::Magenta | Color::LightMagenta => 0.3,
            Color::Cyan | Color::LightCyan => 0.5,
            _ => 0.5, // Default for unknown colors
        }
    }

    /// Calculate contrast ratio between two colors
    pub fn contrast_ratio(fg: Color, bg: Color) -> f32 {
        let l1 = Self::get_luminance(fg);
        let l2 = Self::get_luminance(bg);

        let lighter = l1.max(l2);
        let darker = l1.min(l2);

        (lighter + 0.05) / (darker + 0.05)
    }

    /// Check if contrast meets WCAG AA standards (4.5:1 for normal text)
    pub fn meets_wcag_aa(fg: Color, bg: Color) -> bool {
        Self::contrast_ratio(fg, bg) >= 4.5
    }

    /// Check if contrast meets WCAG AAA standards (7:1 for normal text)
    pub fn meets_wcag_aaa(fg: Color, bg: Color) -> bool {
        Self::contrast_ratio(fg, bg) >= 7.0
    }

    /// Test all text color combinations for a theme
    pub fn test_theme_contrast(flavour: ThemeFlavour) -> Vec<ContrastTestResult> {
        theme::set_flavour(flavour);

        vec![
            // Test primary text colors
            ContrastTestResult {
                name: "Primary text on background".to_string(),
                fg: theme::text_primary(),
                bg: theme::background(),
                ratio: Self::contrast_ratio(theme::text_primary(), theme::background()),
                meets_aa: Self::meets_wcag_aa(theme::text_primary(), theme::background()),
            },
            ContrastTestResult {
                name: "Secondary text on background".to_string(),
                fg: theme::text_secondary(),
                bg: theme::background(),
                ratio: Self::contrast_ratio(theme::text_secondary(), theme::background()),
                meets_aa: Self::meets_wcag_aa(theme::text_secondary(), theme::background()),
            },
            ContrastTestResult {
                name: "Muted text on background".to_string(),
                fg: theme::text_muted(),
                bg: theme::background(),
                ratio: Self::contrast_ratio(theme::text_muted(), theme::background()),
                meets_aa: Self::meets_wcag_aa(theme::text_muted(), theme::background()),
            },
            // Test selection colors
            ContrastTestResult {
                name: "Text on selection".to_string(),
                fg: theme::text_primary(),
                bg: theme::selection(),
                ratio: Self::contrast_ratio(theme::text_primary(), theme::selection()),
                meets_aa: Self::meets_wcag_aa(theme::text_primary(), theme::selection()),
            },
            // Test accent colors
            ContrastTestResult {
                name: "Highlight on background".to_string(),
                fg: theme::highlight(),
                bg: theme::background(),
                ratio: Self::contrast_ratio(theme::highlight(), theme::background()),
                meets_aa: Self::meets_wcag_aa(theme::highlight(), theme::background()),
            },
            // Test semantic colors
            ContrastTestResult {
                name: "Error text on background".to_string(),
                fg: theme::error(),
                bg: theme::background(),
                ratio: Self::contrast_ratio(theme::error(), theme::background()),
                meets_aa: Self::meets_wcag_aa(theme::error(), theme::background()),
            },
            ContrastTestResult {
                name: "Success text on background".to_string(),
                fg: theme::success(),
                bg: theme::background(),
                ratio: Self::contrast_ratio(theme::success(), theme::background()),
                meets_aa: Self::meets_wcag_aa(theme::success(), theme::background()),
            },
        ]
    }

    /// Verify that focus indicators are visible
    pub fn test_focus_visibility(flavour: ThemeFlavour) -> bool {
        theme::set_flavour(flavour);

        // Border focused should be visible against background
        // Using a lower threshold since focus indicators often use color rather than just brightness
        Self::contrast_ratio(theme::border_focused(), theme::background()) >= 2.0
    }
}

#[derive(Debug)]
pub struct ContrastTestResult {
    pub name: String,
    pub fg: Color,
    pub bg: Color,
    pub ratio: f32,
    pub meets_aa: bool,
}

impl ContrastTestResult {
    pub fn display(&self) -> String {
        format!(
            "{}: {:.2}:1 {}",
            self.name,
            self.ratio,
            if self.meets_aa { "✓" } else { "✗" }
        )
    }
}

/// Test helper to verify component styling
pub struct StyleVerifier;

impl StyleVerifier {
    /// Check if a style uses theme colors (not hardcoded)
    pub fn uses_theme_colors(style: &Style) -> bool {
        // This is a simplified check - in practice we'd need more sophisticated verification
        if let Some(fg) = style.fg {
            // Check if the color matches any theme color
            let theme_colors = [
                theme::text_primary(),
                theme::text_secondary(),
                theme::text_muted(),
                theme::primary(),
                theme::secondary(),
                theme::accent(),
                theme::highlight(),
                theme::error(),
                theme::success(),
                theme::warning(),
                theme::info(),
            ];

            theme_colors.contains(&fg)
        } else {
            // No foreground color is OK
            true
        }
    }

    /// Verify all text in a rendered component uses proper colors
    pub fn verify_buffer_colors(buffer: &ratatui::buffer::Buffer) -> Vec<ColorIssue> {
        let mut issues = vec![];

        for y in 0..buffer.area.height {
            for x in 0..buffer.area.width {
                let cell = &buffer[(x, y)];

                // Check for black text on dark backgrounds
                if let Some(fg) = cell.style().fg {
                    if fg == Color::Black || fg == Color::Rgb(0, 0, 0) {
                        if let Some(bg) = cell.style().bg {
                            // Check if background is dark
                            let bg_luminance = ThemeTestHelper::get_luminance(bg);
                            if bg_luminance < 0.5 {
                                issues.push(ColorIssue {
                                    position: (x, y),
                                    description: "Black text on dark background".to_string(),
                                    fg,
                                    bg,
                                });
                            }
                        }
                    }
                }

                // Check for invisible text (same fg and bg)
                if let (Some(fg), Some(bg)) = (cell.style().fg, cell.style().bg) {
                    if fg == bg {
                        issues.push(ColorIssue {
                            position: (x, y),
                            description: "Text same color as background".to_string(),
                            fg,
                            bg,
                        });
                    }
                }
            }
        }

        issues
    }
}

#[derive(Debug)]
pub struct ColorIssue {
    pub position: (u16, u16),
    pub description: String,
    pub fg: Color,
    pub bg: Color,
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_contrast_calculation() {
        // White on black should have high contrast
        let ratio = ThemeTestHelper::contrast_ratio(Color::White, Color::Black);
        assert!(ratio > 20.0);

        // Same color should have ratio of 1
        let ratio = ThemeTestHelper::contrast_ratio(Color::Gray, Color::Gray);
        assert!((ratio - 1.0).abs() < 0.1);
    }

    #[test]
    fn test_wcag_compliance() {
        // White on black should meet both AA and AAA
        assert!(ThemeTestHelper::meets_wcag_aa(Color::White, Color::Black));
        assert!(ThemeTestHelper::meets_wcag_aaa(Color::White, Color::Black));

        // Very dark gray (RGB 40,40,40) on black should fail AA (contrast ~2.5:1)
        assert!(!ThemeTestHelper::meets_wcag_aa(
            Color::Rgb(40, 40, 40),
            Color::Black
        ));
    }
}
