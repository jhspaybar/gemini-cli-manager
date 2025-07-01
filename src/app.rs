use std::sync::{Arc, RwLock};

use color_eyre::Result;
use crossterm::event::KeyEvent;
use ratatui::prelude::Rect;
use tokio::sync::mpsc;
use tracing::{debug, info};

use crate::{
    action::Action,
    components::{Component, settings_view::{UserSettings, SettingsManager}},
    config::Config,
    storage::Storage,
    tui::{Event, Tui},
    view::ViewManager,
};

pub struct App {
    config: Config,
    components: Vec<Box<dyn Component>>,
    should_quit: bool,
    should_suspend: bool,
    last_tick_key_events: Vec<KeyEvent>,
    action_tx: mpsc::UnboundedSender<Action>,
    action_rx: mpsc::UnboundedReceiver<Action>,
    storage: Storage,
    settings: Arc<RwLock<UserSettings>>,
    in_form_view: bool,
}


impl App {
    pub fn new() -> Result<Self> {
        let (action_tx, action_rx) = mpsc::unbounded_channel();
        
        // Initialize storage
        let storage = Storage::new()?;
        storage.init()?;
        
        // Load settings from disk into shared memory
        let settings_manager = SettingsManager::new()?;
        let settings = Arc::new(RwLock::new(settings_manager.get_settings().clone()));
        
        // Apply saved theme from the loaded settings
        if let Ok(settings_lock) = settings.read() {
            if let Err(e) = crate::theme::set_theme_by_name(&settings_lock.theme) {
                debug!("Warning: Could not apply saved theme '{}': {}", settings_lock.theme, e);
            }
        }
        
        // Create view manager with storage
        let view_manager = ViewManager::with_storage(storage.clone());
        
        Ok(Self {
            components: vec![
                Box::new(view_manager),
            ],
            should_quit: false,
            should_suspend: false,
            config: Config::new()?,
            last_tick_key_events: Vec::new(),
            action_tx,
            action_rx,
            storage,
            settings,
            in_form_view: false,
        })
    }

    pub async fn run(&mut self) -> Result<()> {
        let mut tui = Tui::new()?
            // .mouse(true) // uncomment this line to enable mouse support
            .tick_rate(4.0)
            .frame_rate(60.0);
        tui.enter()?;

        for component in self.components.iter_mut() {
            component.register_action_handler(self.action_tx.clone())?;
        }
        for component in self.components.iter_mut() {
            component.register_config_handler(self.config.clone())?;
        }
        for component in self.components.iter_mut() {
            component.register_settings_handler(self.settings.clone())?;
        }
        for component in self.components.iter_mut() {
            component.init(tui.size()?)?;
        }

        let action_tx = self.action_tx.clone();
        loop {
            self.handle_events(&mut tui).await?;
            self.handle_actions(&mut tui)?;
            if self.should_suspend {
                tui.suspend()?;
                action_tx.send(Action::Resume)?;
                action_tx.send(Action::ClearScreen)?;
                // tui.mouse(true);
                tui.enter()?;
            } else if self.should_quit {
                tui.stop()?;
                break;
            }
        }
        tui.exit()?;
        Ok(())
    }

    async fn handle_events(&mut self, tui: &mut Tui) -> Result<()> {
        let Some(event) = tui.next_event().await else {
            return Ok(());
        };
        let action_tx = self.action_tx.clone();
        
        // First, let components handle the event
        for component in self.components.iter_mut() {
            if let Some(action) = component.handle_events(Some(event.clone()))? {
                action_tx.send(action)?;
            }
        }
        
        // Process system events (Tick, Render, Resize) always
        // But only process Key events if we're not in a form view
        match event {
            Event::Tick => action_tx.send(Action::Tick)?,
            Event::Render => action_tx.send(Action::Render)?,
            Event::Resize(x, y) => action_tx.send(Action::Resize(x, y))?,
            Event::Key(key) => {
                // Only process global keybindings if we're not in a form view
                if !self.in_form_view {
                    self.handle_key_event(key)?;
                }
            }
            Event::Quit => action_tx.send(Action::Quit)?,
            _ => {}
        }
        Ok(())
    }

    fn handle_key_event(&mut self, key: KeyEvent) -> Result<()> {
        let action_tx = self.action_tx.clone();
        let Some(keymap) = self.config.keybindings.get(&crate::config::Mode::Normal) else {
            return Ok(());
        };
        match keymap.get(&vec![key]) {
            Some(action) => {
                info!("Got action: {action:?}");
                action_tx.send(action.clone())?;
            }
            _ => {
                // If the key was not handled as a single key action,
                // then consider it for multi-key combinations.
                self.last_tick_key_events.push(key);

                // Check for multi-key combinations
                if let Some(action) = keymap.get(&self.last_tick_key_events) {
                    info!("Got action: {action:?}");
                    action_tx.send(action.clone())?;
                }
            }
        }
        Ok(())
    }

    fn handle_actions(&mut self, tui: &mut Tui) -> Result<()> {
        while let Ok(action) = self.action_rx.try_recv() {
            if action != Action::Tick && action != Action::Render {
                debug!("{action:?}");
            }
            match action.clone() {
                Action::Tick => {
                    self.last_tick_key_events.drain(..);
                }
                Action::Quit => self.should_quit = true,
                Action::Suspend => self.should_suspend = true,
                Action::Resume => self.should_suspend = false,
                Action::ClearScreen => tui.terminal.clear()?,
                Action::Resize(w, h) => self.handle_resize(tui, w, h)?,
                Action::Render => self.render(tui)?,
                Action::LaunchWithProfile(profile_id) => {
                    self.handle_launch_profile(profile_id, tui)?;
                }
                // Track when we're in form views
                Action::CreateNewExtension | Action::EditExtension(_) | 
                Action::CreateProfile | Action::EditProfile(_) => {
                    self.in_form_view = true;
                }
                Action::NavigateBack | Action::NavigateToExtensions | 
                Action::NavigateToProfiles | Action::NavigateToSettings => {
                    self.in_form_view = false;
                }
                _ => {}
            }
            for component in self.components.iter_mut() {
                if let Some(action) = component.update(action.clone())? {
                    self.action_tx.send(action)?
                };
            }
        }
        Ok(())
    }

    fn handle_resize(&mut self, tui: &mut Tui, w: u16, h: u16) -> Result<()> {
        tui.resize(Rect::new(0, 0, w, h))?;
        self.render(tui)?;
        Ok(())
    }

    fn render(&mut self, tui: &mut Tui) -> Result<()> {
        tui.draw(|frame| {
            // Fill the entire terminal with the theme background
            let area = frame.area();
            frame.render_widget(
                ratatui::widgets::Block::default()
                    .style(ratatui::style::Style::default().bg(crate::theme::background())),
                area,
            );
            
            // Then render components on top
            for component in self.components.iter_mut() {
                if let Err(err) = component.draw(frame, frame.area()) {
                    let _ = self
                        .action_tx
                        .send(Action::Error(format!("Failed to draw: {err:?}")));
                }
            }
        })?;
        Ok(())
    }

    fn handle_launch_profile(&mut self, profile_id: String, tui: &mut Tui) -> Result<()> {
        use crate::launcher::Launcher;
        
        // Get the profile from storage
        match self.storage.load_profile(&profile_id) {
            Ok(profile) => {
                // Exit TUI mode before launching
                tui.exit()?;
                
                // Display launch message
                println!("Preparing to launch profile: {}", profile.display_name());
                println!();
                
                // Launch the profile with storage
                let launcher = Launcher::with_storage(self.storage.clone());
                match launcher.launch_with_profile(&profile) {
                    Ok(_) => {
                        println!();
                        println!("✅ Gemini CLI session ended successfully.");
                        
                        // Small delay to let the user see the message
                        std::thread::sleep(std::time::Duration::from_millis(500));
                        
                        // Re-enter TUI mode
                        tui.enter()?;
                        self.action_tx.send(Action::ClearScreen)?;
                        self.action_tx.send(Action::Render)?;
                    }
                    Err(e) => {
                        eprintln!("❌ Error launching profile: {}", e);
                        
                        // Longer delay for errors so user can read the message
                        std::thread::sleep(std::time::Duration::from_secs(2));
                        
                        // Re-enter TUI mode
                        tui.enter()?;
                        self.action_tx.send(Action::ClearScreen)?;
                        self.action_tx.send(Action::Error(format!("Failed to launch profile: {}", e)))?;
                    }
                }
            }
            Err(e) => {
                // Re-enter TUI mode on error
                tui.enter()?;
                self.action_tx.send(Action::ClearScreen)?;
                self.action_tx.send(Action::Error(format!("Failed to load profile: {}", e)))?;
            }
        }
        
        Ok(())
    }
}

