use std::collections::HashMap;

use color_eyre::Result;
use ratatui::prelude::*;
use tokio::sync::mpsc::UnboundedSender;

use crate::{
    action::Action,
    components::{
        extension_detail::ExtensionDetail, extension_list::ExtensionList,
        profile_detail::ProfileDetail, profile_list::ProfileList, tab_bar::TabBar, Component,
    },
    config::Config,
};

#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
pub enum ViewType {
    ExtensionList,
    ExtensionDetail,
    ProfileList,
    ProfileDetail,
    Settings,
}

pub struct ViewManager {
    current_view: ViewType,
    previous_view: Option<ViewType>,
    views: HashMap<ViewType, Box<dyn Component>>,
    action_tx: Option<UnboundedSender<Action>>,
    tab_bar: TabBar,
}

impl ViewManager {
    pub fn new() -> Self {
        let mut views: HashMap<ViewType, Box<dyn Component>> = HashMap::new();
        
        // Initialize views
        views.insert(ViewType::ExtensionList, Box::new(ExtensionList::new()));
        views.insert(ViewType::ExtensionDetail, Box::new(ExtensionDetail::new()));
        views.insert(ViewType::ProfileList, Box::new(ProfileList::new()));
        views.insert(ViewType::ProfileDetail, Box::new(ProfileDetail::new()));
        
        Self {
            current_view: ViewType::ExtensionList,
            previous_view: None,
            views,
            action_tx: None,
            tab_bar: TabBar::new(),
        }
    }

    pub fn register_action_handler(&mut self, tx: UnboundedSender<Action>) -> Result<()> {
        self.action_tx = Some(tx.clone());
        
        // Register action handler for all views
        for (_, view) in self.views.iter_mut() {
            view.register_action_handler(tx.clone())?;
        }
        
        Ok(())
    }

    pub fn register_config_handler(&mut self, config: Config) -> Result<()> {
        // Register config for all views
        for (_, view) in self.views.iter_mut() {
            view.register_config_handler(config.clone())?;
        }
        
        Ok(())
    }

    pub fn init(&mut self, size: Size) -> Result<()> {
        // Initialize all views
        for (_, view) in self.views.iter_mut() {
            view.init(size)?;
        }
        
        Ok(())
    }

    pub fn update(&mut self, action: Action) -> Result<Option<Action>> {
        // Handle navigation actions
        match &action {
            Action::ViewExtensionDetails(_) => {
                self.navigate_to(ViewType::ExtensionDetail);
            }
            Action::EditProfile(_) => {
                self.navigate_to(ViewType::ProfileDetail);
            }
            Action::NavigateBack => {
                if let Some(prev) = self.previous_view {
                    self.navigate_to(prev);
                }
            }
            Action::NavigateToExtensions => {
                self.navigate_to(ViewType::ExtensionList);
            }
            Action::NavigateToProfiles => {
                self.navigate_to(ViewType::ProfileList);
            }
            Action::NavigateToSettings => {
                self.navigate_to(ViewType::Settings);
            }
            _ => {}
        }
        
        // Update tab bar
        self.tab_bar.update(action.clone())?;
        
        // Forward action to all views (they'll handle what's relevant to them)
        let mut result = None;
        for (_, view) in self.views.iter_mut() {
            if let Some(action) = view.update(action.clone())? {
                result = Some(action);
            }
        }
        
        Ok(result)
    }

    pub fn draw(&mut self, frame: &mut Frame, area: Rect) -> Result<()> {
        use ratatui::layout::{Constraint, Direction, Layout};
        
        // Split area into tab bar and content
        let chunks = Layout::default()
            .direction(Direction::Vertical)
            .constraints([
                Constraint::Length(3), // Tab bar
                Constraint::Min(0),    // Content
            ])
            .split(area);
        
        // Draw tab bar
        self.tab_bar.draw(frame, chunks[0])?;
        
        // Draw current view in remaining space
        if let Some(view) = self.views.get_mut(&self.current_view) {
            view.draw(frame, chunks[1])?;
        }
        
        Ok(())
    }

    pub fn handle_events(
        &mut self,
        event: Option<crate::tui::Event>,
    ) -> Result<Option<Action>> {
        // Forward events only to the current view
        if let Some(view) = self.views.get_mut(&self.current_view) {
            view.handle_events(event)
        } else {
            Ok(None)
        }
    }

    fn navigate_to(&mut self, view_type: ViewType) {
        if self.current_view != view_type {
            self.previous_view = Some(self.current_view);
            self.current_view = view_type;
            self.tab_bar.set_current_view(view_type);
        }
    }
}

impl Component for ViewManager {
    fn register_action_handler(&mut self, tx: UnboundedSender<Action>) -> Result<()> {
        self.register_action_handler(tx)
    }

    fn register_config_handler(&mut self, config: Config) -> Result<()> {
        self.register_config_handler(config)
    }

    fn init(&mut self, area: Size) -> Result<()> {
        self.init(area)
    }

    fn update(&mut self, action: Action) -> Result<Option<Action>> {
        self.update(action)
    }

    fn draw(&mut self, frame: &mut Frame, area: Rect) -> Result<()> {
        self.draw(frame, area)
    }

    fn handle_events(&mut self, event: Option<crate::tui::Event>) -> Result<Option<Action>> {
        self.handle_events(event)
    }
}