use serde::{Deserialize, Serialize};
use strum::Display;

#[derive(Debug, Clone, PartialEq, Eq, Display, Serialize, Deserialize)]
pub enum Action {
    Tick,
    Render,
    Resize(u16, u16),
    Suspend,
    Resume,
    Quit,
    ClearScreen,
    Error(String),
    Help,
    
    // Extension management actions
    ViewExtensionDetails(String), // Extension ID
    ImportExtension,
    CreateNewExtension,
    EditExtension(String), // Extension ID
    DeleteExtension(String), // Extension ID
    
    // Navigation actions
    NavigateToExtensions,
    NavigateToProfiles,
    NavigateToSettings,
    NavigateBack,
    
    // Profile management actions
    CreateProfile,
    EditProfile(String), // Profile ID
    DeleteProfile(String), // Profile ID
    LaunchWithProfile(String), // Profile ID
}
