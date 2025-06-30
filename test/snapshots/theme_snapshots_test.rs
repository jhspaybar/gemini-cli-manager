#[cfg(test)]
mod tests {
    use gemini_cli_manager::{
        theme::{self, ThemeFlavour},
        components::{ExtensionList, ProfileList, ExtensionForm, ProfileForm, TabBar},
        view::ViewType,
        models::{Extension, Profile},
    };
    use crate::test_utils::{
        render_to_string, create_temp_storage, McpFixtures, 
        ExtensionBuilder, ProfileBuilder,
    };
    use insta::assert_snapshot;
    use ratatui::prelude::*;
    
    fn test_all_themes<F>(name: &str, render_fn: F)
    where
        F: Fn() -> String,
    {
        let themes = vec![
            ("mocha", ThemeFlavour::Mocha),
            ("macchiato", ThemeFlavour::Macchiato),
            ("frappe", ThemeFlavour::Frappe),
            ("latte", ThemeFlavour::Latte),
        ];
        
        for (theme_name, flavour) in themes {
            theme::set_flavour(flavour);
            let output = render_fn();
            assert_snapshot!(format!("{}_{}", name, theme_name), output);
        }
    }
    
    #[test]
    fn test_extension_list_empty_snapshot() {
        test_all_themes("extension_list_empty", || {
            let storage = create_temp_storage().0;
            let list = ExtensionList::new(storage);
            
            render_to_string(80, 24, |f| {
                list.render(f, f.area());
            }).unwrap()
        });
    }
    
    #[test]
    fn test_extension_list_populated_snapshot() {
        test_all_themes("extension_list_populated", || {
            let storage = create_temp_storage().0;
            
            // Add test extensions
            let ext1 = McpFixtures::echo_extension();
            let ext2 = McpFixtures::multi_server_extension();
            let ext3 = McpFixtures::context_only_extension();
            
            storage.save_extension(&ext1).unwrap();
            storage.save_extension(&ext2).unwrap();
            storage.save_extension(&ext3).unwrap();
            
            let list = ExtensionList::new(storage);
            
            render_to_string(80, 24, |f| {
                list.render(f, f.area());
            }).unwrap()
        });
    }
    
    #[test]
    fn test_profile_list_snapshot() {
        test_all_themes("profile_list", || {
            let storage = create_temp_storage().0;
            
            // Add test profiles
            let profile1 = ProfileBuilder::new("Development")
                .with_description("Dev environment with all tools")
                .with_extensions(vec!["ext1", "ext2", "ext3"])
                .with_tags(vec!["dev", "local"])
                .as_default()
                .build();
                
            let profile2 = ProfileBuilder::new("Production")
                .with_description("Minimal production environment")
                .with_extensions(vec!["ext1"])
                .with_tags(vec!["prod"])
                .build();
                
            storage.save_profile(&profile1).unwrap();
            storage.save_profile(&profile2).unwrap();
            
            let list = ProfileList::new(storage);
            
            render_to_string(80, 24, |f| {
                list.render(f, f.area());
            }).unwrap()
        });
    }
    
    #[test]
    fn test_extension_form_create_snapshot() {
        test_all_themes("extension_form_create", || {
            let storage = create_temp_storage().0;
            let form = ExtensionForm::new(storage);
            
            render_to_string(80, 35, |f| {
                form.render(f, f.area());
            }).unwrap()
        });
    }
    
    #[test]
    fn test_extension_form_filled_snapshot() {
        test_all_themes("extension_form_filled", || {
            let storage = create_temp_storage().0;
            let mut form = ExtensionForm::new(storage);
            
            // Simulate filling out the form
            form.set_name("Test Extension");
            form.set_version("1.0.0");
            form.set_description("A test extension for snapshot testing");
            form.set_tags("test, snapshot, example");
            
            render_to_string(80, 35, |f| {
                form.render(f, f.area());
            }).unwrap()
        });
    }
    
    #[test]
    fn test_profile_form_snapshot() {
        test_all_themes("profile_form", || {
            let storage = create_temp_storage().0;
            
            // Add some extensions for selection
            let ext1 = ExtensionBuilder::new("Extension One").build();
            let ext2 = ExtensionBuilder::new("Extension Two").build();
            storage.save_extension(&ext1).unwrap();
            storage.save_extension(&ext2).unwrap();
            
            let mut form = ProfileForm::new(storage);
            form.set_name("Test Profile");
            form.set_description("Profile for testing");
            
            render_to_string(80, 30, |f| {
                form.render(f, f.area());
            }).unwrap()
        });
    }
    
    #[test]
    fn test_tab_bar_extensions_active_snapshot() {
        test_all_themes("tab_bar_extensions", || {
            let tab_bar = TabBar::new(ViewType::ExtensionList);
            
            render_to_string(80, 3, |f| {
                tab_bar.render(f, f.area()).unwrap();
            }).unwrap()
        });
    }
    
    #[test]
    fn test_tab_bar_profiles_active_snapshot() {
        test_all_themes("tab_bar_profiles", || {
            let tab_bar = TabBar::new(ViewType::ProfileList);
            
            render_to_string(80, 3, |f| {
                tab_bar.render(f, f.area()).unwrap();
            }).unwrap()
        });
    }
    
    #[test]
    fn test_error_dialog_snapshot() {
        test_all_themes("error_dialog", || {
            // Simulate an error dialog overlay
            render_to_string(60, 20, |f| {
                let area = f.area();
                
                // Render background
                let block = Block::default()
                    .borders(Borders::ALL)
                    .title("Background View");
                f.render_widget(block, area);
                
                // Calculate centered error dialog area
                let dialog_width = 40;
                let dialog_height = 8;
                let x = (area.width.saturating_sub(dialog_width)) / 2;
                let y = (area.height.saturating_sub(dialog_height)) / 2;
                
                let dialog_area = Rect {
                    x: area.x + x,
                    y: area.y + y,
                    width: dialog_width.min(area.width),
                    height: dialog_height.min(area.height),
                };
                
                // Render error dialog
                let error_block = Block::default()
                    .borders(Borders::ALL)
                    .border_style(Style::default().fg(theme::error()))
                    .title("Error")
                    .title_style(Style::default().fg(theme::error()).bold());
                    
                let error_text = Paragraph::new("Failed to delete extension:\nExtension is in use by profiles")
                    .style(Style::default().fg(theme::text_primary()))
                    .alignment(Alignment::Center)
                    .wrap(Wrap { trim: true });
                    
                f.render_widget(Clear, dialog_area);
                f.render_widget(error_block.clone(), dialog_area);
                
                let inner = error_block.inner(dialog_area);
                f.render_widget(error_text, inner);
            }).unwrap()
        });
    }
    
    #[test]
    fn test_confirmation_dialog_snapshot() {
        test_all_themes("confirmation_dialog", || {
            render_to_string(60, 20, |f| {
                let area = f.area();
                
                // Background
                let block = Block::default()
                    .borders(Borders::ALL)
                    .title("Extension List");
                f.render_widget(block, area);
                
                // Confirmation dialog
                let dialog_width = 45;
                let dialog_height = 10;
                let x = (area.width.saturating_sub(dialog_width)) / 2;
                let y = (area.height.saturating_sub(dialog_height)) / 2;
                
                let dialog_area = Rect {
                    x: area.x + x,
                    y: area.y + y,
                    width: dialog_width.min(area.width),
                    height: dialog_height.min(area.height),
                };
                
                let confirm_block = Block::default()
                    .borders(Borders::ALL)
                    .border_style(Style::default().fg(theme::warning()))
                    .title("Confirm Deletion")
                    .title_style(Style::default().fg(theme::warning()).bold());
                    
                let confirm_text = vec![
                    Line::from("Are you sure you want to delete this extension?"),
                    Line::from(""),
                    Line::from(Span::styled("Extension: ", Style::default().fg(theme::text_secondary())))
                        .push_span(Span::styled("echo-test", Style::default().fg(theme::text_primary()).bold())),
                    Line::from(""),
                    Line::from(""),
                    Line::from(vec![
                        Span::styled("[Y]es  ", Style::default().fg(theme::success()).bold()),
                        Span::styled("[N]o", Style::default().fg(theme::error()).bold()),
                    ]),
                ];
                
                let paragraph = Paragraph::new(confirm_text)
                    .alignment(Alignment::Center)
                    .wrap(Wrap { trim: true });
                    
                f.render_widget(Clear, dialog_area);
                f.render_widget(confirm_block.clone(), dialog_area);
                
                let inner = confirm_block.inner(dialog_area);
                f.render_widget(paragraph, inner);
            }).unwrap()
        });
    }
    
    #[test]
    fn test_search_mode_snapshot() {
        test_all_themes("search_mode", || {
            let storage = create_temp_storage().0;
            
            // Add extensions to search through
            for i in 1..=5 {
                let ext = ExtensionBuilder::new(&format!("Extension {}", i))
                    .with_tags(vec!["test", &format!("tag{}", i)])
                    .build();
                storage.save_extension(&ext).unwrap();
            }
            
            render_to_string(80, 24, |f| {
                let area = f.area();
                
                // Render list with search bar
                let list_area = Rect {
                    x: area.x,
                    y: area.y + 3, // Leave space for search
                    width: area.width,
                    height: area.height.saturating_sub(3),
                };
                
                // Search bar
                let search_block = Block::default()
                    .borders(Borders::ALL)
                    .border_style(Style::default().fg(theme::highlight()))
                    .title("Search");
                    
                let search_text = Paragraph::new("tag2")
                    .style(Style::default().fg(theme::text_primary()));
                    
                let search_area = Rect {
                    x: area.x,
                    y: area.y,
                    width: area.width,
                    height: 3,
                };
                
                f.render_widget(search_block.clone(), search_area);
                let inner = search_block.inner(search_area);
                f.render_widget(search_text, inner);
                
                // Filtered results
                let results_block = Block::default()
                    .borders(Borders::ALL)
                    .title("Extensions (1 match)");
                f.render_widget(results_block, list_area);
            }).unwrap()
        });
    }
    
    #[test]
    fn test_minimal_terminal_size_snapshot() {
        // Test rendering at minimum supported size
        test_all_themes("minimal_size", || {
            let storage = create_temp_storage().0;
            let list = ExtensionList::new(storage);
            
            render_to_string(40, 12, |f| {
                list.render(f, f.area());
            }).unwrap()
        });
    }
}