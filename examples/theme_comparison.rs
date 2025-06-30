use gemini_cli_manager::theme::{self, ThemeFlavour};
use ratatui::{prelude::*, widgets::*};

fn main() -> color_eyre::Result<()> {
    // Initialize color_eyre for error handling
    color_eyre::install()?;
    
    // Create a simple terminal backend for demo
    let mut stdout = std::io::stdout();
    let backend = ratatui::backend::CrosstermBackend::new(&mut stdout);
    let mut terminal = Terminal::new(backend)?;
    
    // Clear the terminal
    terminal.clear()?;
    
    // Test each theme
    let themes = [
        ("Mocha (Dark)", ThemeFlavour::Mocha),
        ("Macchiato", ThemeFlavour::Macchiato),
        ("Frappe", ThemeFlavour::Frappe),
        ("Latte (Light)", ThemeFlavour::Latte),
    ];
    
    for (name, flavour) in themes {
        theme::set_flavour(flavour);
        
        terminal.draw(|frame| {
            let area = frame.area();
            
            // Fill background
            frame.render_widget(
                Block::default()
                    .style(Style::default().bg(theme::background())),
                area,
            );
            
            // Create a demo UI
            let chunks = Layout::vertical([
                Constraint::Length(3),  // Title
                Constraint::Min(10),    // Content
                Constraint::Length(3),  // Status
            ])
            .split(area);
            
            // Title bar
            let title = Paragraph::new(format!(" Gemini CLI Manager - {} Theme ", name))
                .style(Style::default()
                    .fg(theme::text_primary())
                    .bg(theme::surface()))
                .block(Block::default()
                    .borders(Borders::ALL)
                    .border_style(Style::default().fg(theme::border()))
                    .border_type(BorderType::Rounded))
                .alignment(Alignment::Center);
            frame.render_widget(title, chunks[0]);
            
            // Content area with color samples
            let content_chunks = Layout::horizontal([
                Constraint::Percentage(50),
                Constraint::Percentage(50),
            ])
            .split(chunks[1]);
            
            // Left side - UI elements
            let ui_samples = vec![
                Line::from(vec![
                    Span::styled("Primary: ", Style::default().fg(theme::text_secondary())),
                    Span::styled("Blue accent color", Style::default().fg(theme::primary())),
                ]),
                Line::from(""),
                Line::from(vec![
                    Span::styled("Success: ", Style::default().fg(theme::text_secondary())),
                    Span::styled("✓ Operation complete", Style::default().fg(theme::success())),
                ]),
                Line::from(""),
                Line::from(vec![
                    Span::styled("Warning: ", Style::default().fg(theme::text_secondary())),
                    Span::styled("⚠ Check configuration", Style::default().fg(theme::warning())),
                ]),
                Line::from(""),
                Line::from(vec![
                    Span::styled("Error: ", Style::default().fg(theme::text_secondary())),
                    Span::styled("✗ Failed to load", Style::default().fg(theme::error())),
                ]),
                Line::from(""),
                Line::from(vec![
                    Span::styled("Info: ", Style::default().fg(theme::text_secondary())),
                    Span::styled("ℹ 5 extensions loaded", Style::default().fg(theme::info())),
                ]),
            ];
            
            let ui_block = Paragraph::new(ui_samples)
                .block(Block::default()
                    .title(" UI Elements ")
                    .borders(Borders::ALL)
                    .border_style(Style::default().fg(theme::border_focused()))
                    .border_type(BorderType::Rounded))
                .style(Style::default().fg(theme::text_primary()));
            frame.render_widget(ui_block, content_chunks[0]);
            
            // Right side - List example
            let list_items = vec![
                ListItem::new("Extension 1").style(Style::default().fg(theme::text_primary())),
                ListItem::new("Extension 2").style(Style::default().fg(theme::text_primary())),
                ListItem::new("Extension 3 (selected)").style(Style::default()
                    .fg(theme::highlight())
                    .add_modifier(Modifier::BOLD)),
                ListItem::new("Extension 4").style(Style::default().fg(theme::text_primary())),
                ListItem::new("Extension 5 (muted)").style(Style::default().fg(theme::text_muted())),
            ];
            
            let list = List::new(list_items)
                .block(Block::default()
                    .title(" Extension List ")
                    .borders(Borders::ALL)
                    .border_style(Style::default().fg(theme::border()))
                    .border_type(BorderType::Rounded))
                .highlight_style(Style::default().bg(theme::selection()));
            frame.render_widget(list, content_chunks[1]);
            
            // Status bar
            let status = Paragraph::new(" Press Enter to see next theme, q to quit ")
                .style(Style::default()
                    .fg(theme::text_muted())
                    .bg(theme::overlay()))
                .block(Block::default()
                    .borders(Borders::ALL)
                    .border_style(Style::default().fg(theme::border()))
                    .border_type(BorderType::Rounded))
                .alignment(Alignment::Center);
            frame.render_widget(status, chunks[2]);
        })?;
        
        // Wait for user input
        println!("\n{} theme rendered. Press Enter to continue...", name);
        let mut input = String::new();
        std::io::stdin().read_line(&mut input)?;
        
        if input.trim() == "q" {
            break;
        }
    }
    
    // Restore terminal
    terminal.clear()?;
    
    Ok(())
}