#[cfg(test)]
mod tests {
    use crate::test_utils::{
        ExtensionBuilder, McpFixtures, ProfileBuilder, create_test_storage, render_to_string,
    };
    use gemini_cli_manager::{
        components::{
            Component, extension_form::ExtensionForm, extension_list::ExtensionList,
            profile_list::ProfileList, tab_bar::TabBar,
        },
        theme::{self, ThemeFlavour},
        view::ViewType,
    };
    use insta::assert_snapshot;
    use ratatui::prelude::*;
    use ratatui::widgets::*;

    fn test_all_themes<F>(name: &str, render_fn: F)
    where
        F: Fn() -> Result<String, Box<dyn std::error::Error>>,
    {
        let themes = vec![
            ("mocha", ThemeFlavour::Mocha),
            ("macchiato", ThemeFlavour::Macchiato),
            ("frappe", ThemeFlavour::Frappe),
            ("latte", ThemeFlavour::Latte),
        ];

        for (theme_name, flavour) in themes {
            theme::set_flavour(flavour);
            let output = render_fn().unwrap_or_else(|e| format!("Error: {e}"));
            assert_snapshot!(format!("{}_{}", name, theme_name), output);
        }
    }

    #[test]
    fn test_extension_list_empty_snapshot() {
        test_all_themes("extension_list_empty", || {
            let storage = create_test_storage();
            let mut list = ExtensionList::with_storage(storage);

            render_to_string(80, 24, |f| {
                list.draw(f, f.area()).unwrap();
            })
        });
    }

    #[test]
    fn test_extension_list_populated_snapshot() {
        test_all_themes("extension_list_populated", || {
            let storage = create_test_storage();

            // Add test extensions
            let ext1 = McpFixtures::echo_extension();
            storage.save_extension(&ext1).unwrap();

            let ext2 = ExtensionBuilder::new("AI Assistant")
                .with_version("2.0.0")
                .with_description("Advanced AI-powered development assistant")
                .with_tags(vec!["ai", "productivity", "coding"])
                .build();
            storage.save_extension(&ext2).unwrap();

            let ext3 = ExtensionBuilder::new("Database Tools")
                .with_version("1.5.3")
                .with_description("Comprehensive database management and migration tools")
                .with_tags(vec!["database", "sql", "migration"])
                .build();
            storage.save_extension(&ext3).unwrap();

            let mut list = ExtensionList::with_storage(storage);

            render_to_string(80, 24, |f| {
                list.draw(f, f.area()).unwrap();
            })
        });
    }

    #[test]
    fn test_profile_list_empty_snapshot() {
        test_all_themes("profile_list_empty", || {
            let storage = create_test_storage();
            let mut list = ProfileList::with_storage(storage);

            render_to_string(80, 24, |f| {
                list.draw(f, f.area()).unwrap();
            })
        });
    }

    #[test]
    fn test_profile_list_populated_snapshot() {
        test_all_themes("profile_list_populated", || {
            let storage = create_test_storage();

            // Add test profiles
            let profile1 = ProfileBuilder::new("Development")
                .with_description("Local development environment")
                .with_extensions(vec!["echo-server", "ai-assistant"])
                .with_tags(vec!["dev", "local"])
                .build();
            storage.save_profile(&profile1).unwrap();

            let profile2 = ProfileBuilder::new("Production")
                .with_description("Production environment with minimal extensions")
                .with_extensions(vec!["database-tools"])
                .with_tags(vec!["prod", "minimal"])
                .as_default()
                .build();
            storage.save_profile(&profile2).unwrap();

            let profile3 = ProfileBuilder::new("Testing")
                .with_description("QA and testing environment")
                .with_tags(vec!["test", "qa"])
                .build();
            storage.save_profile(&profile3).unwrap();

            let mut list = ProfileList::with_storage(storage);

            render_to_string(80, 24, |f| {
                list.draw(f, f.area()).unwrap();
            })
        });
    }

    #[test]
    fn test_extension_form_empty_snapshot() {
        test_all_themes("extension_form_empty", || {
            let storage = create_test_storage();
            let mut form = ExtensionForm::new(storage);

            render_to_string(80, 30, |f| {
                form.draw(f, f.area()).unwrap();
            })
        });
    }

    #[test]
    fn test_extension_form_filled_snapshot() {
        test_all_themes("extension_form_filled", || {
            let storage = create_test_storage();
            let ext = ExtensionBuilder::new("Test Extension")
                .with_version("1.0.0")
                .with_description("This is a test extension for snapshot testing")
                .with_tags(vec!["test", "snapshot"])
                .build();

            let mut form = ExtensionForm::with_extension(storage, &ext);

            render_to_string(80, 30, |f| {
                form.draw(f, f.area()).unwrap();
            })
        });
    }

    #[test]
    fn test_tab_bar_extensions_active_snapshot() {
        test_all_themes("tab_bar_extensions", || {
            let mut tab_bar = TabBar::new();
            tab_bar.set_current_view(ViewType::ExtensionList);

            render_to_string(80, 3, |f| {
                tab_bar.draw(f, f.area()).unwrap();
            })
        });
    }

    #[test]
    fn test_tab_bar_profiles_active_snapshot() {
        test_all_themes("tab_bar_profiles", || {
            let mut tab_bar = TabBar::new();
            tab_bar.set_current_view(ViewType::ProfileList);

            render_to_string(80, 3, |f| {
                tab_bar.draw(f, f.area()).unwrap();
            })
        });
    }

    #[test]
    fn test_complex_layout_snapshot() {
        test_all_themes("complex_layout", || {
            let storage = create_test_storage();

            // Add some data
            let ext = McpFixtures::echo_extension();
            storage.save_extension(&ext).unwrap();

            let mut tab_bar = TabBar::new();
            let mut ext_list = ExtensionList::with_storage(storage);

            render_to_string(80, 24, |f| {
                let chunks = Layout::default()
                    .direction(Direction::Vertical)
                    .constraints([
                        Constraint::Length(3), // Tab bar
                        Constraint::Min(10),   // Content
                    ])
                    .split(f.area());

                tab_bar.draw(f, chunks[0]).unwrap();
                ext_list.draw(f, chunks[1]).unwrap();
            })
        });
    }

    #[test]
    fn test_error_state_snapshot() {
        test_all_themes("error_state", || {
            render_to_string(60, 15, |f| {
                let block = Block::default()
                    .borders(Borders::ALL)
                    .border_style(Style::default().fg(theme::error()))
                    .title(" Error ");

                let inner = block.inner(f.area());
                f.render_widget(block, f.area());

                let error_style = Style::default()
                    .fg(theme::error())
                    .add_modifier(Modifier::BOLD);

                let error_text =
                    Paragraph::new("Failed to delete extension:\nExtension is in use by profiles")
                        .style(error_style)
                        .wrap(Wrap { trim: true });

                f.render_widget(error_text, inner);
            })
        });
    }
}
