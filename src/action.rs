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
    Success(String),
    Help,

    // Extension management actions
    ViewExtensionDetails(String), // Extension ID
    ImportExtension,
    ResetImportDialog,       // Reset import dialog state
    CreateNewExtension,
    EditExtension(String),   // Extension ID
    DeleteExtension(String), // Extension ID
    RefreshExtensions,       // Reload extensions from storage

    // Navigation actions
    NavigateToExtensions,
    NavigateToProfiles,
    NavigateToSettings,
    NavigateBack,

    // Profile management actions
    ViewProfileDetails(String),// Profile ID
    CreateProfile,
    EditProfile(String),       // Profile ID
    DeleteProfile(String),     // Profile ID
    ConfirmDelete,             // Confirm deletion
    CancelDelete,              // Cancel deletion
    LaunchWithProfile(String), // Profile ID
    RefreshProfiles,           // Reload profiles from storage

    // Settings actions
    ChangeTheme(String),       // Theme name
    UpdateKeybinding(String, String), // Action name, key combination
    ResetKeybindings,          // Reset to defaults
    SaveSettings,              // Save settings to file
}
