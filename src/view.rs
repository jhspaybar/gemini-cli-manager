use std::collections::HashMap;
use std::sync::{Arc, RwLock};
use std::time::{Duration, Instant};

use color_eyre::Result;
use ratatui::prelude::*;
use tokio::sync::mpsc::UnboundedSender;

use crate::{
    action::Action,
    components::{
        Component,
        confirm_dialog::ConfirmDialog,
        extension_detail::ExtensionDetail,
        extension_form::ExtensionForm,
        extension_list::ExtensionList,
        import_dialog::ImportDialog,
        profile_detail::ProfileDetail,
        profile_form::ProfileForm,
        profile_list::ProfileList,
        settings_view::{Settings, UserSettings},
        tab_bar::TabBar,
    },
    config::Config,
    storage::Storage,
    theme,
};

#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
pub enum ViewType {
    ExtensionList,
    ExtensionDetail,
    ExtensionCreate,
    ExtensionEdit,
    ExtensionImport,
    ProfileList,
    ProfileDetail,
    ProfileCreate,
    ProfileEdit,
    ConfirmDelete,
    Settings,
}

pub struct ViewManager {
    current_view: ViewType,
    previous_view: Option<ViewType>,
    views: HashMap<ViewType, Box<dyn Component>>,
    action_tx: Option<UnboundedSender<Action>>,
    tab_bar: TabBar,
    storage: Storage,
    editing_profile_id: Option<String>,
    deleting_profile_id: Option<String>,
    editing_extension_id: Option<String>,
    deleting_extension_id: Option<String>,
    came_from_detail_view: bool, // Track if we came from detail view when editing
    error_message: Option<(String, Instant)>,
    success_message: Option<(String, Instant)>,
    message_display_duration: Duration,
    error_display_duration: Duration,
}

impl Default for ViewManager {
    fn default() -> Self {
        // Create a default storage instance
        let storage = Storage::default();
        Self::with_storage(storage)
    }
}

impl ViewManager {
    #[allow(dead_code)]
    pub fn new() -> Self {
        Self::default()
    }

    /// Test helper method - returns current view
    #[doc(hidden)]
    #[allow(dead_code)]
    pub fn current_view(&self) -> ViewType {
        self.current_view
    }

    /// Test helper method - checks if error is displayed
    #[doc(hidden)]
    #[allow(dead_code)]
    pub fn has_error(&self) -> bool {
        self.error_message.is_some()
    }

    pub fn with_storage(storage: Storage) -> Self {
        let mut views: HashMap<ViewType, Box<dyn Component>> = HashMap::new();

        // Initialize views with storage
        views.insert(
            ViewType::ExtensionList,
            Box::new(ExtensionList::with_storage(storage.clone())),
        );
        views.insert(
            ViewType::ExtensionDetail,
            Box::new(ExtensionDetail::with_storage(storage.clone())),
        );
        views.insert(
            ViewType::ExtensionCreate,
            Box::new(ExtensionForm::new(storage.clone())),
        );
        views.insert(
            ViewType::ExtensionImport,
            Box::new(ImportDialog::new(storage.clone())),
        );
        views.insert(
            ViewType::ProfileList,
            Box::new(ProfileList::with_storage(storage.clone())),
        );
        views.insert(
            ViewType::ProfileDetail,
            Box::new(ProfileDetail::with_storage(storage.clone())),
        );
        views.insert(
            ViewType::ProfileCreate,
            Box::new(ProfileForm::new(storage.clone())),
        );
        views.insert(ViewType::Settings, Box::new(Settings::new()));

        Self {
            current_view: ViewType::ExtensionList,
            previous_view: None,
            views,
            action_tx: None,
            tab_bar: TabBar::new(),
            storage,
            editing_profile_id: None,
            deleting_profile_id: None,
            editing_extension_id: None,
            deleting_extension_id: None,
            came_from_detail_view: false,
            error_message: None,
            success_message: None,
            message_display_duration: Duration::from_secs(3),
            error_display_duration: Duration::from_secs(10),
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

    pub fn register_settings_handler(&mut self, settings: Arc<RwLock<UserSettings>>) -> Result<()> {
        // Register settings for all views
        for (_, view) in self.views.iter_mut() {
            view.register_settings_handler(settings.clone())?;
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
            Action::CreateNewExtension => {
                // Clear any editing state when creating new
                self.editing_extension_id = None;
                self.came_from_detail_view = false;

                // Create a fresh ExtensionForm for creating new extensions
                let mut create_form = ExtensionForm::new(self.storage.clone());

                // Register action handler for the new form
                if let Some(tx) = &self.action_tx {
                    let _ = create_form.register_action_handler(tx.clone());
                }

                self.views
                    .insert(ViewType::ExtensionCreate, Box::new(create_form));
                self.navigate_to(ViewType::ExtensionCreate);
            }
            Action::ImportExtension => {
                // Navigate to the import dialog
                self.navigate_to(ViewType::ExtensionImport);
                // Send a reset action to the import dialog
                if let Some(tx) = &self.action_tx {
                    let _ = tx.send(Action::ResetImportDialog);
                }
            }
            Action::EditExtension(id) => {
                // Track where we came from
                self.came_from_detail_view = self.current_view == ViewType::ExtensionDetail;

                // Store the extension ID we're editing
                self.editing_extension_id = Some(id.clone());

                // Load the extension into the edit form
                if let Ok(extension) = self.storage.load_extension(id) {
                    // Create a new edit form with the extension data
                    let mut edit_form =
                        ExtensionForm::with_extension(self.storage.clone(), &extension);

                    // Register action handler for the new form
                    if let Some(tx) = &self.action_tx {
                        let _ = edit_form.register_action_handler(tx.clone());
                    }

                    self.views
                        .insert(ViewType::ExtensionEdit, Box::new(edit_form));
                    self.navigate_to(ViewType::ExtensionEdit);
                }
            }
            Action::DeleteExtension(id) => {
                // First check if any profiles reference this extension
                let profiles = self.storage.list_profiles().unwrap_or_default();
                let referenced_by: Vec<String> = profiles
                    .iter()
                    .filter(|p| p.extension_ids.contains(id))
                    .map(|p| p.name.clone())
                    .collect();

                if !referenced_by.is_empty() {
                    // Show error message
                    let message = format!(
                        "Cannot delete extension: it is referenced by {} profile(s): {}",
                        referenced_by.len(),
                        referenced_by.join(", ")
                    );
                    self.error_message = Some((message, Instant::now()));
                } else {
                    // Store the extension ID to delete
                    self.deleting_extension_id = Some(id.clone());

                    // Load the extension to get its name for the confirmation message
                    let message = if let Ok(extension) = self.storage.load_extension(id) {
                        format!(
                            "Are you sure you want to delete the extension '{}'?\nThis action cannot be undone.",
                            extension.name
                        )
                    } else {
                        "Are you sure you want to delete this extension?\nThis action cannot be undone.".to_string()
                    };

                    // Create confirmation dialog
                    let dialog = ConfirmDialog::new("Delete Extension", &message)
                        .with_actions(Action::ConfirmDelete, Action::CancelDelete);

                    self.views.insert(ViewType::ConfirmDelete, Box::new(dialog));
                    self.navigate_to(ViewType::ConfirmDelete);
                }
            }
            Action::ViewProfileDetails(_id) => {
                self.navigate_to(ViewType::ProfileDetail);
            }
            Action::EditProfile(id) => {
                // Track where we came from
                self.came_from_detail_view = self.current_view == ViewType::ProfileDetail;

                // Store the profile ID we're editing
                self.editing_profile_id = Some(id.clone());

                // Load the profile into the edit form
                if let Ok(profile) = self.storage.load_profile(id) {
                    // Create a new edit form with the profile data
                    let mut edit_form = ProfileForm::with_profile(self.storage.clone(), &profile);

                    // Register action handler for the new form
                    if let Some(tx) = &self.action_tx {
                        let _ = edit_form.register_action_handler(tx.clone());
                    }

                    self.views
                        .insert(ViewType::ProfileEdit, Box::new(edit_form));
                    self.navigate_to(ViewType::ProfileEdit);
                }
            }
            Action::CreateProfile => {
                // Clear any editing state when creating new
                self.editing_profile_id = None;
                self.came_from_detail_view = false;

                // Create a fresh ProfileForm for creating new profiles
                let mut create_form = ProfileForm::new(self.storage.clone());

                // Register action handler for the new form
                if let Some(tx) = &self.action_tx {
                    let _ = create_form.register_action_handler(tx.clone());
                }

                self.views
                    .insert(ViewType::ProfileCreate, Box::new(create_form));
                self.navigate_to(ViewType::ProfileCreate);
            }
            Action::NavigateBack => {
                // Smart navigation based on current view
                match self.current_view {
                    ViewType::ProfileDetail => {
                        // From profile detail, always go back to profile list
                        self.navigate_to(ViewType::ProfileList);
                    }
                    ViewType::ProfileEdit => {
                        // From edit, go back to where we came from
                        if self.came_from_detail_view && self.editing_profile_id.is_some() {
                            // We came from detail view, go back there
                            self.navigate_to(ViewType::ProfileDetail);
                        } else {
                            // We came from list or it's a new profile, go to list
                            self.navigate_to(ViewType::ProfileList);
                        }
                        self.editing_profile_id = None;
                        self.came_from_detail_view = false;
                    }
                    ViewType::ProfileCreate => {
                        // From create, always go back to list
                        self.navigate_to(ViewType::ProfileList);
                    }
                    ViewType::ExtensionDetail => {
                        // From extension detail, go back to extension list
                        self.navigate_to(ViewType::ExtensionList);
                    }
                    ViewType::ExtensionEdit => {
                        // From edit, go back to where we came from
                        if self.came_from_detail_view && self.editing_extension_id.is_some() {
                            // We came from detail view, go back there
                            self.navigate_to(ViewType::ExtensionDetail);
                        } else {
                            // We came from list, go to list
                            self.navigate_to(ViewType::ExtensionList);
                        }
                        self.editing_extension_id = None;
                        self.came_from_detail_view = false;
                    }
                    ViewType::ExtensionCreate => {
                        // From create, always go back to list
                        self.navigate_to(ViewType::ExtensionList);
                    }
                    ViewType::ExtensionImport => {
                        // From import, go back to extension list
                        self.navigate_to(ViewType::ExtensionList);
                    }
                    _ => {
                        // For other views, use the previous view
                        if let Some(prev) = self.previous_view {
                            self.navigate_to(prev);
                        }
                    }
                }
            }
            Action::NavigateToExtensions => {
                self.navigate_to(ViewType::ExtensionList);
            }
            Action::NavigateToProfiles => {
                // Clear any editing state when going to profile list
                self.editing_profile_id = None;
                self.navigate_to(ViewType::ProfileList);
            }
            Action::NavigateToSettings => {
                self.navigate_to(ViewType::Settings);
            }
            Action::DeleteProfile(id) => {
                // Store the profile ID to delete
                self.deleting_profile_id = Some(id.clone());

                // Load the profile to get its name for the confirmation message
                let message = if let Ok(profile) = self.storage.load_profile(id) {
                    format!(
                        "Are you sure you want to delete the profile '{}'?\nThis action cannot be undone.",
                        profile.name
                    )
                } else {
                    "Are you sure you want to delete this profile?\nThis action cannot be undone."
                        .to_string()
                };

                // Create confirmation dialog
                let dialog = ConfirmDialog::new("Delete Profile", &message)
                    .with_actions(Action::ConfirmDelete, Action::CancelDelete);

                self.views.insert(ViewType::ConfirmDelete, Box::new(dialog));
                self.navigate_to(ViewType::ConfirmDelete);
            }
            Action::ConfirmDelete => {
                // Check if we're deleting a profile or extension
                if let Some(id) = &self.deleting_profile_id {
                    // Delete profile
                    if let Err(e) = self.storage.delete_profile(id) {
                        // Send error action
                        if let Some(tx) = &self.action_tx {
                            let _ =
                                tx.send(Action::Error(format!("Failed to delete profile: {e}")));
                        }
                    } else {
                        // Send success notification and refresh action
                        if let Some(tx) = &self.action_tx {
                            let _ = tx
                                .send(Action::Success("Profile deleted successfully".to_string()));
                            let _ = tx.send(Action::RefreshProfiles);
                            let _ = tx.send(Action::Render);
                        }
                    }
                    // Clear deletion state
                    self.deleting_profile_id = None;

                    // If we were in profile detail view, go to list instead
                    if self.previous_view == Some(ViewType::ProfileDetail) {
                        self.navigate_to(ViewType::ProfileList);
                    } else if let Some(prev) = self.previous_view {
                        self.navigate_to(prev);
                    }
                } else if let Some(id) = &self.deleting_extension_id {
                    // Delete extension
                    if let Err(e) = self.storage.delete_extension(id) {
                        // Send error action
                        if let Some(tx) = &self.action_tx {
                            let _ =
                                tx.send(Action::Error(format!("Failed to delete extension: {e}")));
                        }
                    } else {
                        // Send success notification and refresh action
                        if let Some(tx) = &self.action_tx {
                            let _ = tx.send(Action::Success(
                                "Extension deleted successfully".to_string(),
                            ));
                            let _ = tx.send(Action::RefreshExtensions);
                            let _ = tx.send(Action::Render);
                        }
                    }
                    // Clear deletion state
                    self.deleting_extension_id = None;

                    // If we were in extension detail view, go to list instead
                    if self.previous_view == Some(ViewType::ExtensionDetail) {
                        self.navigate_to(ViewType::ExtensionList);
                    } else if let Some(prev) = self.previous_view {
                        self.navigate_to(prev);
                    }
                }
            }
            Action::CancelDelete => {
                // Clear deletion state and go back
                self.deleting_profile_id = None;
                self.deleting_extension_id = None;
                if let Some(prev) = self.previous_view {
                    self.navigate_to(prev);
                }
            }
            Action::Error(msg) => {
                self.error_message = Some((msg.clone(), Instant::now()));
            }
            Action::Success(msg) => {
                self.success_message = Some((msg.clone(), Instant::now()));
            }
            Action::Tick => {
                // Clear old error messages (using longer duration)
                if let Some((_, timestamp)) = &self.error_message
                    && timestamp.elapsed() > self.error_display_duration {
                        self.error_message = None;
                    }

                // Clear old success messages
                if let Some((_, timestamp)) = &self.success_message
                    && timestamp.elapsed() > self.message_display_duration {
                        self.success_message = None;
                    }
            }
            _ => {}
        }

        // Clone action for passing to components
        let action_clone = action.clone();

        // Update tab bar
        self.tab_bar.update(action_clone.clone())?;

        // Forward action to all views (they'll handle what's relevant to them)
        let mut result = None;
        for (_, view) in self.views.iter_mut() {
            if let Some(action) = view.update(action_clone.clone())? {
                result = Some(action);
            }
        }

        Ok(result)
    }

    pub fn draw(&mut self, frame: &mut Frame, area: Rect) -> Result<()> {
        use ratatui::layout::{Constraint, Direction, Layout};
        use ratatui::style::Modifier;
        use ratatui::widgets::{Block, BorderType, Borders, Clear, Paragraph, Wrap};

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

        // Draw error message if present
        if let Some((message, _)) = &self.error_message {
            let popup_area = self.centered_rect(60, 20, area);

            // Clear the area first
            frame.render_widget(Clear, popup_area);

            // Create error block with background
            let error_block = Block::default()
                .title(" Error ")
                .title_style(
                    Style::default()
                        .fg(theme::error())
                        .add_modifier(Modifier::BOLD),
                )
                .title_alignment(Alignment::Center)
                .borders(Borders::ALL)
                .border_type(BorderType::Rounded)
                .border_style(Style::default().fg(theme::error()))
                .style(Style::default().bg(theme::overlay()));

            let error_content = format!("✗ {message}\n\n(Press Esc to dismiss)");
            let error_text = Paragraph::new(error_content)
                .block(error_block)
                .style(
                    Style::default()
                        .fg(theme::error())
                        .bg(theme::overlay())
                        .add_modifier(Modifier::BOLD),
                )
                .alignment(Alignment::Center)
                .wrap(Wrap { trim: true });

            frame.render_widget(error_text, popup_area);
        }

        // Draw success message if present
        if let Some((message, _)) = &self.success_message {
            // Position in top-right corner
            let notification_width = 40.min(area.width - 4);
            let notification_height = 3;
            let notification_area = Rect {
                x: area.width.saturating_sub(notification_width + 2),
                y: 1,
                width: notification_width,
                height: notification_height,
            };

            // Clear the area first
            frame.render_widget(Clear, notification_area);

            // Create success block with background
            let success_block = Block::default()
                .title(" Success ")
                .title_style(
                    Style::default()
                        .fg(theme::success())
                        .add_modifier(Modifier::BOLD),
                )
                .borders(Borders::ALL)
                .border_type(BorderType::Rounded)
                .border_style(Style::default().fg(theme::success()))
                .style(Style::default().bg(theme::surface()));

            let success_text = Paragraph::new(format!("✓ {message}"))
                .style(
                    Style::default()
                        .fg(theme::success())
                        .bg(theme::surface())
                        .add_modifier(Modifier::BOLD),
                )
                .block(success_block)
                .wrap(Wrap { trim: true });

            frame.render_widget(success_text, notification_area);
        }

        Ok(())
    }

    /// Helper function to create a centered rect
    fn centered_rect(&self, percent_x: u16, percent_y: u16, area: Rect) -> Rect {
        use ratatui::layout::{Constraint, Flex, Layout};

        let vertical = Layout::vertical([Constraint::Percentage(percent_y)]).flex(Flex::Center);
        let horizontal = Layout::horizontal([Constraint::Percentage(percent_x)]).flex(Flex::Center);
        let [area] = vertical.areas(area);
        let [area] = horizontal.areas(area);
        area
    }

    pub fn handle_events(&mut self, event: Option<crate::tui::Event>) -> Result<Option<Action>> {
        use crossterm::event::KeyCode;

        // Handle error message dismissal
        if self.error_message.is_some()
            && let Some(crate::tui::Event::Key(key)) = &event
            && key.code == KeyCode::Esc {
                    self.error_message = None;
                    return Ok(Some(Action::Render));
                }

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

    fn register_settings_handler(&mut self, settings: Arc<RwLock<UserSettings>>) -> Result<()> {
        self.register_settings_handler(settings)
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
