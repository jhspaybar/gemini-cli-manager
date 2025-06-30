use color_eyre::Result;
use ratatui::{prelude::*, widgets::*};

use super::Component;
use crate::view::ViewType;

pub struct TabBar {
    current_view: ViewType,
    tabs: Vec<(String, ViewType)>,
}

impl TabBar {
    pub fn new() -> Self {
        Self {
            current_view: ViewType::ExtensionList,
            tabs: vec![
                ("Extensions".to_string(), ViewType::ExtensionList),
                ("Profiles".to_string(), ViewType::ProfileList),
            ],
        }
    }

    pub fn set_current_view(&mut self, view: ViewType) {
        self.current_view = view;
    }
}

impl Component for TabBar {
    fn update(&mut self, action: crate::action::Action) -> Result<Option<crate::action::Action>> {
        use crate::action::Action;
        
        // Update current view based on navigation actions
        match action {
            Action::NavigateToExtensions => {
                self.current_view = ViewType::ExtensionList;
            }
            Action::NavigateToProfiles => {
                self.current_view = ViewType::ProfileList;
            }
            Action::ViewExtensionDetails(_) => {
                self.current_view = ViewType::ExtensionDetail;
            }
            Action::EditProfile(_) => {
                // This could be either detail or edit view
                // ViewManager will set the correct view via set_current_view
            }
            Action::CreateProfile => {
                self.current_view = ViewType::ProfileCreate;
            }
            _ => {}
        }
        
        Ok(None)
    }

    fn draw(&mut self, frame: &mut Frame, area: Rect) -> Result<()> {
        // Don't draw if area is too small
        if area.height < 3 {
            return Ok(());
        }

        // Create tab titles with indicators
        let titles: Vec<Line> = self
            .tabs
            .iter()
            .map(|(title, view_type)| {
                let is_active = match (self.current_view, view_type) {
                    // Extensions tab is active for extension-related views
                    (ViewType::ExtensionList | ViewType::ExtensionDetail, ViewType::ExtensionList) => true,
                    // Profiles tab is active for all profile-related views
                    (ViewType::ProfileList | ViewType::ProfileDetail | ViewType::ProfileCreate | ViewType::ProfileEdit, ViewType::ProfileList) => true,
                    _ => false,
                };

                if is_active {
                    Line::from(vec![
                        Span::raw(" "),
                        Span::styled(
                            title,
                            Style::default()
                                .fg(Color::Yellow)
                                .add_modifier(Modifier::BOLD),
                        ),
                        Span::raw(" "),
                    ])
                } else {
                    Line::from(vec![
                        Span::raw(" "),
                        Span::styled(title, Style::default().fg(Color::DarkGray)),
                        Span::raw(" "),
                    ])
                }
            })
            .collect();

        // Create tabs widget
        let tabs = Tabs::new(titles)
            .block(
                Block::default()
                    .borders(Borders::ALL)
                    .border_type(BorderType::Rounded)
                    .border_style(Style::default().fg(Color::Cyan)),
            )
            .highlight_style(
                Style::default()
                    .bg(Color::DarkGray)
                    .add_modifier(Modifier::BOLD),
            )
            .divider(Span::raw(" │ "));

        // Select the active tab
        let selected = match self.current_view {
            ViewType::ExtensionList | ViewType::ExtensionDetail => 0,
            ViewType::ProfileList | ViewType::ProfileDetail | ViewType::ProfileCreate | ViewType::ProfileEdit => 1,
            _ => 0,
        };

        // Render with selection
        frame.render_widget(tabs.select(selected), area);

        // Add breadcrumb for detail/form views
        if matches!(self.current_view, ViewType::ExtensionDetail | ViewType::ProfileDetail | ViewType::ProfileCreate | ViewType::ProfileEdit) {
            let breadcrumb = match self.current_view {
                ViewType::ExtensionDetail => " > Extension Details",
                ViewType::ProfileDetail => " > Profile Details",
                ViewType::ProfileCreate => " > Create Profile",
                ViewType::ProfileEdit => " > Edit Profile",
                _ => "",
            };

            if !breadcrumb.is_empty() && area.width > 40 {
                let breadcrumb_area = Rect {
                    x: area.x + area.width.saturating_sub(breadcrumb.len() as u16 + 2),
                    y: area.y + 1,
                    width: breadcrumb.len() as u16,
                    height: 1,
                };

                frame.render_widget(
                    Paragraph::new(breadcrumb).style(Style::default().fg(Color::DarkGray)),
                    breadcrumb_area,
                );
            }
        }

        Ok(())
    }
}