# Gemini CLI Manager Architecture Guide

This document outlines the architecture for implementing the Gemini CLI Manager using Rust and Ratatui.

## Table of Contents
1. [Application Architecture](#application-architecture)
2. [Core Components](#core-components)
3. [Data Flow](#data-flow)
4. [Module Structure](#module-structure)
5. [Implementation Roadmap](#implementation-roadmap)

## Application Architecture

### Overview

The Gemini CLI Manager follows a component-based architecture with message passing for communication:

```
┌─────────────────────────────────────────────────────────────┐
│                        Terminal UI                          │
├─────────────────────────────────────────────────────────────┤
│                    Component Layer                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐       │
│  │ Extension   │  │  Profile    │  │   Tool      │       │
│  │ Manager     │  │  Manager    │  │  Manager    │       │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘       │
│         │                 │                 │               │
├─────────┴─────────────────┴─────────────────┴───────────────┤
│                    Action Bus (MPSC)                        │
├─────────────────────────────────────────────────────────────┤
│                    Service Layer                            │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐       │
│  │ Extension   │  │  Profile    │  │   Config    │       │
│  │  Service    │  │  Service    │  │  Service    │       │
│  └─────────────┘  └─────────────┘  └─────────────┘       │
├─────────────────────────────────────────────────────────────┤
│                    Storage Layer                            │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐       │
│  │    JSON     │  │    TOML     │  │ File System │       │
│  │   Storage   │  │   Storage   │  │   Storage   │       │
│  └─────────────┘  └─────────────┘  └─────────────┘       │
└─────────────────────────────────────────────────────────────┘
```

### Key Design Principles

1. **Separation of Concerns**: UI components don't directly access storage
2. **Message Passing**: Components communicate through actions
3. **Async Operations**: I/O operations are non-blocking
4. **Testability**: Each layer can be tested independently
5. **Extensibility**: Easy to add new components and features

## Core Components

### 1. Extension Manager Component

Handles the UI for browsing, installing, and managing extensions.

```rust
// src/components/extension_manager.rs
pub struct ExtensionManager {
    extensions: Vec<Extension>,
    selected: usize,
    filter: String,
    view: ExtensionView,
}

pub enum ExtensionView {
    List,
    Details(String),
    Install,
}

impl ExtensionManager {
    pub fn handle_action(&mut self, action: Action) -> Option<Action> {
        match action {
            Action::NavigateDown => self.next_extension(),
            Action::NavigateUp => self.previous_extension(),
            Action::Select => Some(Action::ViewExtensionDetails(self.get_selected_id())),
            Action::InstallExtension(id) => Some(Action::StartInstall(id)),
            _ => None,
        }
    }
}
```

### 2. Profile Manager Component

Manages user profiles and environment configurations.

```rust
// src/components/profile_manager.rs
pub struct ProfileManager {
    profiles: Vec<Profile>,
    active_profile: Option<String>,
    edit_mode: bool,
    form: ProfileForm,
}

pub struct ProfileForm {
    name: String,
    env_vars: HashMap<String, String>,
    extensions: Vec<String>,
    validation_errors: Vec<String>,
}
```

### 3. Tool Manager Component

Interfaces with external tools and manages their execution.

```rust
// src/components/tool_manager.rs
pub struct ToolManager {
    available_tools: Vec<Tool>,
    running_tools: HashMap<String, ToolProcess>,
}

pub struct Tool {
    id: String,
    name: String,
    command: String,
    args: Vec<String>,
    env: HashMap<String, String>,
}
```

## Data Flow

### Action Flow

```
User Input → Event → Component → Action → Service → Storage
                           ↓                    ↓
                        UI Update ← Action ← Result
```

### Example: Installing an Extension

```rust
// 1. User presses Enter on an extension
KeyCode::Enter => {
    let action = Action::InstallExtension(extension_id);
    self.action_tx.send(action)?;
}

// 2. App receives action and delegates to service
Action::InstallExtension(id) => {
    let service = self.extension_service.clone();
    tokio::spawn(async move {
        match service.install_extension(&id).await {
            Ok(()) => Action::InstallComplete(id),
            Err(e) => Action::Error(e.to_string()),
        }
    });
}

// 3. Service performs installation
pub async fn install_extension(&self, id: &str) -> Result<()> {
    let manifest = self.fetch_manifest(id).await?;
    self.validate_manifest(&manifest)?;
    self.download_files(&manifest).await?;
    self.register_extension(manifest)?;
    Ok(())
}

// 4. UI updates based on result
Action::InstallComplete(id) => {
    self.show_notification("Extension installed successfully");
    self.refresh_extension_list();
}
```

## Module Structure

### Proposed Directory Layout

```
src/
├── main.rs                 # Entry point
├── app.rs                  # Main application orchestrator
├── action.rs               # Action types and routing
├── error.rs                # Error types
│
├── components/             # UI Components
│   ├── mod.rs
│   ├── extension_manager/
│   │   ├── mod.rs
│   │   ├── list.rs         # Extension list view
│   │   ├── details.rs      # Extension details view
│   │   └── install.rs      # Installation progress
│   ├── profile_manager/
│   │   ├── mod.rs
│   │   ├── list.rs
│   │   ├── form.rs
│   │   └── quick_switch.rs
│   └── common/
│       ├── mod.rs
│       ├── search_bar.rs
│       ├── status_bar.rs
│       └── help.rs
│
├── services/               # Business Logic
│   ├── mod.rs
│   ├── extension_service.rs
│   ├── profile_service.rs
│   ├── config_service.rs
│   └── launcher_service.rs
│
├── models/                 # Data Structures
│   ├── mod.rs
│   ├── extension.rs
│   ├── profile.rs
│   ├── tool.rs
│   └── config.rs
│
├── storage/                # Persistence Layer
│   ├── mod.rs
│   ├── json_store.rs
│   ├── toml_store.rs
│   └── file_store.rs
│
├── widgets/                # Custom Ratatui Widgets
│   ├── mod.rs
│   ├── form_field.rs
│   ├── multi_select.rs
│   └── file_picker.rs
│
└── utils/                  # Utilities
    ├── mod.rs
    ├── validation.rs
    ├── paths.rs
    └── shell.rs
```

### Key Modules

#### Models

```rust
// src/models/extension.rs
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Extension {
    pub id: String,
    pub name: String,
    pub version: String,
    pub description: String,
    pub author: String,
    pub homepage: Option<String>,
    pub tags: Vec<String>,
    pub dependencies: Vec<Dependency>,
    pub install_path: PathBuf,
}

// src/models/profile.rs
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Profile {
    pub name: String,
    pub description: Option<String>,
    pub env_vars: HashMap<String, String>,
    pub extensions: Vec<String>,
    pub tools: Vec<String>,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}
```

#### Services

```rust
// src/services/extension_service.rs
pub struct ExtensionService {
    storage: Arc<dyn ExtensionStorage>,
    registry: Arc<ExtensionRegistry>,
}

impl ExtensionService {
    pub async fn list_available(&self) -> Result<Vec<Extension>> {
        self.registry.fetch_catalog().await
    }

    pub async fn install(&self, id: &str) -> Result<()> {
        let extension = self.registry.get(id).await?;
        self.validate_dependencies(&extension)?;
        self.download_and_extract(&extension).await?;
        self.storage.save(&extension)?;
        Ok(())
    }

    pub async fn uninstall(&self, id: &str) -> Result<()> {
        let extension = self.storage.get(id)?;
        self.cleanup_files(&extension).await?;
        self.storage.remove(id)?;
        Ok(())
    }
}
```

## Implementation Roadmap

### Phase 1: Core Infrastructure (Current)
- [x] Basic TUI framework setup
- [x] Component trait and registration
- [x] Action-based message passing
- [x] Configuration system
- [ ] Error handling framework

### Phase 2: Extension Management
- [ ] Extension model and storage
- [ ] Extension list component
- [ ] Extension details view
- [ ] Search and filter functionality
- [ ] Mock extension registry

### Phase 3: Profile Management
- [ ] Profile model and storage
- [ ] Profile list and switching
- [ ] Profile creation/edit form
- [ ] Environment variable management
- [ ] Profile activation

### Phase 4: Tool Integration
- [ ] Tool execution framework
- [ ] Tool configuration
- [ ] Output capture and display
- [ ] Background process management

### Phase 5: Advanced Features
- [ ] Extension dependencies
- [ ] Auto-updates
- [ ] Import/export profiles
- [ ] Plugin system for extensions
- [ ] Network registry integration

### Phase 6: Polish and Testing
- [ ] Comprehensive error handling
- [ ] Performance optimization
- [ ] Integration tests
- [ ] Documentation
- [ ] Release packaging

## Best Practices

### Component Development

1. **Keep components focused**: One component = one responsibility
2. **Use stateful widgets**: For lists, tables, etc.
3. **Handle all possible actions**: Even if just returning None
4. **Validate input early**: In the component before sending actions

### Service Layer

1. **Make services async**: All I/O should be non-blocking
2. **Use dependency injection**: Services should depend on traits
3. **Handle errors gracefully**: Return Result types
4. **Cache when appropriate**: Reduce file system access

### Testing Strategy

```rust
// Component tests
#[test]
fn test_extension_navigation() {
    let mut manager = ExtensionManager::new(vec![
        Extension { id: "1", name: "Ext 1" },
        Extension { id: "2", name: "Ext 2" },
    ]);
    
    assert_eq!(manager.selected, 0);
    manager.next_extension();
    assert_eq!(manager.selected, 1);
}

// Service tests
#[tokio::test]
async fn test_extension_install() {
    let storage = MockExtensionStorage::new();
    let registry = MockExtensionRegistry::new();
    let service = ExtensionService::new(storage, registry);
    
    let result = service.install("test-ext").await;
    assert!(result.is_ok());
}

// Integration tests
#[tokio::test]
async fn test_full_extension_workflow() {
    let app = create_test_app().await;
    
    // Simulate user actions
    app.handle_key(KeyCode::Down);
    app.handle_key(KeyCode::Enter);
    
    // Verify state changes
    assert_eq!(app.current_view(), View::ExtensionDetails);
}
```

This architecture provides a solid foundation for building the Gemini CLI Manager with clear separation of concerns, testability, and room for future expansion.