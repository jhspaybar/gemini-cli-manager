use gemini_cli_manager::theme::{self, ThemeFlavour};
use ratatui::{prelude::*, widgets::*};

fn main() -> color_eyre::Result<()> {
    color_eyre::install()?;
    
    // Set Mocha theme
    theme::set_flavour(ThemeFlavour::Mocha);
    
    let mut stdout = std::io::stdout();
    let backend = ratatui::backend::CrosstermBackend::new(&mut stdout);
    let mut terminal = Terminal::new(backend)?;
    
    terminal.clear()?;
    
    terminal.draw(|frame| {
        let area = frame.area();
        
        // Fill background with theme background
        frame.render_widget(
            Block::default().style(Style::default().bg(theme::background())),
            area,
        );
        
        let chunks = Layout::vertical([
            Constraint::Length(3),
            Constraint::Length(10),
            Constraint::Length(10),
            Constraint::Min(5),
        ])
        .split(area);
        
        // Title
        let title = Paragraph::new(" Contrast Test - Catppuccin Mocha ")
            .style(Style::default().fg(theme::text_primary()).bg(theme::surface()))
            .block(Block::default().borders(Borders::ALL).border_style(Style::default().fg(theme::border())))
            .alignment(Alignment::Center);
        frame.render_widget(title, chunks[0]);
        
        // Text contrast examples
        let text_examples = vec![
            Line::from(vec![
                Span::raw("Primary text: "),
                Span::styled("This is primary text", Style::default().fg(theme::text_primary())),
            ]),
            Line::from(vec![
                Span::raw("Secondary text: "),
                Span::styled("This is secondary text", Style::default().fg(theme::text_secondary())),
            ]),
            Line::from(vec![
                Span::raw("Muted text: "),
                Span::styled("This is muted text (improved contrast)", Style::default().fg(theme::text_muted())),
            ]),
            Line::from(vec![
                Span::raw("Disabled text: "),
                Span::styled("This is disabled text", Style::default().fg(theme::text_disabled())),
            ]),
        ];
        
        let text_block = Paragraph::new(text_examples)
            .block(Block::default()
                .title(" Text Contrast Levels ")
                .borders(Borders::ALL)
                .border_style(Style::default().fg(theme::border())))
            .style(Style::default().fg(theme::text_primary()));
        frame.render_widget(text_block, chunks[1]);
        
        // Selection examples
        let selection_items = vec![
            ListItem::new("Normal item").style(Style::default().fg(theme::text_primary())),
            ListItem::new("Selected item (with background)").style(Style::default().fg(theme::text_primary()).add_modifier(Modifier::BOLD)),
            ListItem::new("Another normal item").style(Style::default().fg(theme::text_primary())),
            ListItem::new("Muted item").style(Style::default().fg(theme::text_muted())),
        ];
        
        let mut list_state = ListState::default();
        list_state.select(Some(1));
        
        let selection_list = List::new(selection_items)
            .block(Block::default()
                .title(" Selection Contrast (Improved) ")
                .borders(Borders::ALL)
                .border_style(Style::default().fg(theme::border())))
            .highlight_style(Style::default().bg(theme::selection()))
            .highlight_symbol("│ ");
        
        frame.render_stateful_widget(selection_list, chunks[2], &mut list_state);
        
        // Color palette
        let palette = vec![
            Line::from(vec![
                Span::styled("█████", Style::default().fg(theme::primary())),
                Span::raw(" Primary   "),
                Span::styled("█████", Style::default().fg(theme::secondary())),
                Span::raw(" Secondary   "),
                Span::styled("█████", Style::default().fg(theme::accent())),
                Span::raw(" Accent"),
            ]),
            Line::from(vec![
                Span::styled("█████", Style::default().fg(theme::success())),
                Span::raw(" Success   "),
                Span::styled("█████", Style::default().fg(theme::warning())),
                Span::raw(" Warning    "),
                Span::styled("█████", Style::default().fg(theme::error())),
                Span::raw(" Error"),
            ]),
            Line::from(vec![
                Span::styled("█████", Style::default().bg(theme::selection())),
                Span::raw(" Selection "),
                Span::styled("█████", Style::default().bg(theme::surface())),
                Span::raw(" Surface    "),
                Span::styled("█████", Style::default().bg(theme::overlay())),
                Span::raw(" Overlay"),
            ]),
        ];
        
        let palette_block = Paragraph::new(palette)
            .block(Block::default()
                .title(" Color Palette ")
                .borders(Borders::ALL)
                .border_style(Style::default().fg(theme::border())))
            .style(Style::default().fg(theme::text_primary()));
        frame.render_widget(palette_block, chunks[3]);
    })?;
    
    println!("\n\nContrast improvements:");
    println!("- Muted text now uses overlay1 for better visibility");
    println!("- Selection background uses surface1 for better contrast");
    println!("- Selected items use primary text color with bold instead of yellow");
    println!("- Consistent selection styling across all lists");
    println!("\nPress Enter to exit...");
    
    let mut input = String::new();
    std::io::stdin().read_line(&mut input)?;
    
    terminal.clear()?;
    
    Ok(())
}