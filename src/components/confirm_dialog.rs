use color_eyre::Result;
use ratatui::{prelude::*, widgets::*};
use tokio::sync::mpsc::UnboundedSender;

use super::Component;
use crate::{action::Action, config::Config, theme};

pub struct ConfirmDialog {
    command_tx: Option<UnboundedSender<Action>>,
    title: String,
    message: String,
    confirm_action: Option<Action>,
    cancel_action: Option<Action>,
    selected_button: bool, // false = Cancel, true = Confirm
}

impl ConfirmDialog {
    pub fn new(title: String, message: String) -> Self {
        Self {
            command_tx: None,
            title,
            message,
            confirm_action: None,
            cancel_action: None,
            selected_button: false, // Default to Cancel for safety
        }
    }
    
    pub fn with_actions(mut self, confirm: Action, cancel: Action) -> Self {
        self.confirm_action = Some(confirm);
        self.cancel_action = Some(cancel);
        self
    }
}

impl Component for ConfirmDialog {
    fn register_action_handler(&mut self, tx: UnboundedSender<Action>) -> Result<()> {
        self.command_tx = Some(tx);
        Ok(())
    }

    fn register_config_handler(&mut self, _config: Config) -> Result<()> {
        Ok(())
    }

    fn update(&mut self, _action: Action) -> Result<Option<Action>> {
        Ok(None)
    }

    fn draw(&mut self, frame: &mut Frame, area: Rect) -> Result<()> {
        // Calculate dialog size
        let dialog_width = 60.min(area.width - 4);
        let dialog_height = 10.min(area.height - 4);
        
        // Center the dialog
        let x = (area.width.saturating_sub(dialog_width)) / 2;
        let y = (area.height.saturating_sub(dialog_height)) / 2;
        
        let dialog_area = Rect {
            x: area.x + x,
            y: area.y + y,
            width: dialog_width,
            height: dialog_height,
        };
        
        // Clear the background
        let clear = Block::default()
            .style(Style::default().bg(theme::overlay()));
        frame.render_widget(clear, dialog_area);
        
        // Create the dialog block
        let block = Block::default()
            .title(format!(" {} ", self.title))
            .title_alignment(Alignment::Center)
            .borders(Borders::ALL)
            .border_type(BorderType::Double)
            .border_style(Style::default().fg(theme::error()));
        
        let inner = block.inner(dialog_area);
        frame.render_widget(block, dialog_area);
        
        // Layout: Message and buttons
        let chunks = Layout::default()
            .direction(Direction::Vertical)
            .margin(1)
            .constraints([
                Constraint::Min(3),    // Message
                Constraint::Length(3), // Buttons
            ])
            .split(inner);
        
        // Render message
        let message = Paragraph::new(self.message.as_str())
            .wrap(Wrap { trim: true })
            .alignment(Alignment::Center);
        frame.render_widget(message, chunks[0]);
        
        // Render buttons
        let button_area = Layout::default()
            .direction(Direction::Horizontal)
            .constraints([
                Constraint::Percentage(50),
                Constraint::Percentage(50),
            ])
            .split(chunks[1]);
        
        // Cancel button
        let cancel_style = if !self.selected_button {
            Style::default()
                .fg(theme::background())
                .bg(theme::text_primary())
                .add_modifier(Modifier::BOLD)
        } else {
            Style::default().fg(theme::text_muted())
        };
        
        let cancel_button = Paragraph::new(" Cancel (Esc) ")
            .style(cancel_style)
            .alignment(Alignment::Center)
            .block(
                Block::default()
                    .borders(Borders::ALL)
                    .border_type(BorderType::Rounded),
            );
        frame.render_widget(cancel_button, button_area[0]);
        
        // Confirm button
        let confirm_style = if self.selected_button {
            Style::default()
                .fg(theme::background())
                .bg(theme::error())
                .add_modifier(Modifier::BOLD)
        } else {
            Style::default().fg(theme::error())
        };
        
        let confirm_button = Paragraph::new(" Delete (Enter) ")
            .style(confirm_style)
            .alignment(Alignment::Center)
            .block(
                Block::default()
                    .borders(Borders::ALL)
                    .border_type(BorderType::Rounded),
            );
        frame.render_widget(confirm_button, button_area[1]);
        
        Ok(())
    }

    fn handle_events(&mut self, event: Option<crate::tui::Event>) -> Result<Option<Action>> {
        use crossterm::event::KeyCode;

        match event {
            Some(crate::tui::Event::Key(key)) => match key.code {
                KeyCode::Left | KeyCode::Right | KeyCode::Tab => {
                    self.selected_button = !self.selected_button;
                    Ok(Some(Action::Render))
                }
                KeyCode::Enter => {
                    if self.selected_button {
                        // Confirm
                        Ok(self.confirm_action.clone())
                    } else {
                        // Cancel
                        Ok(self.cancel_action.clone())
                    }
                }
                KeyCode::Esc => {
                    // Always cancel on Esc
                    Ok(self.cancel_action.clone())
                }
                KeyCode::Char('y') | KeyCode::Char('Y') => {
                    // Quick confirm
                    Ok(self.confirm_action.clone())
                }
                KeyCode::Char('n') | KeyCode::Char('N') => {
                    // Quick cancel
                    Ok(self.cancel_action.clone())
                }
                _ => Ok(None),
            },
            _ => Ok(None),
        }
    }
}