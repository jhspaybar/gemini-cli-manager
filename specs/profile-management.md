# Profile Management Specification

## Overview

Profiles are the heart of the Gemini CLI Manager, allowing users to quickly switch between different configurations for various projects, clients, or environments. Each profile represents a complete configuration state including enabled extensions, environment variables, and MCP server settings.

## Profile System Architecture

### Profile Definition

```go
type Profile struct {
    ID          string                 `yaml:"id"`
    Name        string                 `yaml:"name"`
    Description string                 `yaml:"description"`
    Icon        string                 `yaml:"icon,omitempty"`
    Color       string                 `yaml:"color,omitempty"`
    
    // Configuration
    Extensions  []ExtensionRef         `yaml:"extensions"`
    Environment map[string]string      `yaml:"environment"`
    MCPServers  map[string]ServerConfig `yaml:"mcp_servers"`
    
    // Metadata
    CreatedAt   time.Time             `yaml:"created_at"`
    UpdatedAt   time.Time             `yaml:"updated_at"`
    LastUsed    *time.Time            `yaml:"last_used,omitempty"`
    UsageCount  int                   `yaml:"usage_count"`
    
    // Advanced
    Inherits    []string              `yaml:"inherits,omitempty"`
    Tags        []string              `yaml:"tags,omitempty"`
    AutoDetect  *AutoDetectRules      `yaml:"auto_detect,omitempty"`
}

type ExtensionRef struct {
    ID      string            `yaml:"id"`
    Enabled bool              `yaml:"enabled"`
    Config  map[string]interface{} `yaml:"config,omitempty"`
}

type ServerConfig struct {
    Enabled  bool              `yaml:"enabled"`
    Settings map[string]interface{} `yaml:"settings,omitempty"`
}

type AutoDetectRules struct {
    Patterns []string          `yaml:"patterns"`
    Priority int               `yaml:"priority"`
}
```

### Profile Storage

#### File Structure
```
~/.gemini/profiles/
├── default.yaml          # Default profile (always exists)
├── web-development.yaml
├── data-science.yaml
├── client-project-a.yaml
└── .backups/
    ├── web-development.yaml.20240115-100000
    └── web-development.yaml.20240114-100000
```

#### Example Profile File
```yaml
id: web-development
name: Web Development
description: Full-stack web development with React and Node.js
icon: "🌐"
color: "#61DAFB"

extensions:
  - id: typescript-tools
    enabled: true
    config:
      strictMode: true
      target: "ES2022"
  
  - id: react-helper
    enabled: true
    config:
      version: "18"
  
  - id: database-tools
    enabled: true
    config:
      default_connection: "postgresql://localhost/devdb"

environment:
  NODE_ENV: development
  REACT_APP_API_URL: http://localhost:3001
  DATABASE_URL: postgresql://localhost/devdb
  LOG_LEVEL: debug

mcp_servers:
  code-analyzer:
    enabled: true
    settings:
      model: gemini-2.0-flash
      temperature: 0.3
      
  typescript-lsp:
    enabled: true
    settings:
      diagnostics: true
      formatting: true

created_at: 2024-01-15T10:00:00Z
updated_at: 2024-01-15T14:30:00Z
last_used: 2024-01-15T14:30:00Z
usage_count: 42

inherits:
  - base-development

tags:
  - frontend
  - backend
  - typescript
  - react

auto_detect:
  patterns:
    - "**/package.json"
    - "**/tsconfig.json"
    - "**/.gemini-profile"
  priority: 10
```

## Profile Manager Implementation

### Core Operations

```go
type ProfileManager struct {
    storage     ProfileStorage
    validator   ProfileValidator
    resolver    ProfileResolver
    activeID    string
    cache       *ProfileCache
}

// CRUD Operations
func (pm *ProfileManager) Create(profile *Profile) error
func (pm *ProfileManager) Get(id string) (*Profile, error)
func (pm *ProfileManager) Update(id string, profile *Profile) error
func (pm *ProfileManager) Delete(id string) error
func (pm *ProfileManager) List() ([]*Profile, error)

// Profile Operations
func (pm *ProfileManager) Activate(id string) error
func (pm *ProfileManager) Deactivate() error
func (pm *ProfileManager) GetActive() (*Profile, error)
func (pm *ProfileManager) Clone(sourceID, newID string) error
func (pm *ProfileManager) Export(id string, path string) error
func (pm *ProfileManager) Import(path string) (*Profile, error)
```

### Profile Inheritance

```go
type ProfileResolver struct {
    profiles map[string]*Profile
}

func (pr *ProfileResolver) Resolve(id string) (*ResolvedProfile, error) {
    profile, err := pr.getProfile(id)
    if err != nil {
        return nil, err
    }
    
    // Start with empty resolved profile
    resolved := &ResolvedProfile{
        Base:        profile,
        Extensions:  make(map[string]ExtensionRef),
        Environment: make(map[string]string),
        MCPServers:  make(map[string]ServerConfig),
    }
    
    // Apply inheritance chain
    chain, err := pr.getInheritanceChain(profile)
    if err != nil {
        return nil, err
    }
    
    for _, parent := range chain {
        pr.mergeProfile(resolved, parent)
    }
    
    // Apply current profile (highest priority)
    pr.mergeProfile(resolved, profile)
    
    return resolved, nil
}

func (pr *ProfileResolver) mergeProfile(target *ResolvedProfile, source *Profile) {
    // Merge extensions (source overrides target)
    for _, ext := range source.Extensions {
        target.Extensions[ext.ID] = ext
    }
    
    // Merge environment variables
    for k, v := range source.Environment {
        target.Environment[k] = v
    }
    
    // Merge MCP servers
    for k, v := range source.MCPServers {
        target.MCPServers[k] = v
    }
}
```

### Profile Validation

```go
type ProfileValidator struct {
    extensionRegistry *ExtensionRegistry
}

func (pv *ProfileValidator) Validate(profile *Profile) ValidationResult {
    result := ValidationResult{Valid: true}
    
    // Validate basic fields
    if profile.Name == "" {
        result.AddError("Profile name is required")
    }
    
    // Validate extensions exist
    for _, ext := range profile.Extensions {
        if !pv.extensionRegistry.Exists(ext.ID) {
            result.AddError(fmt.Sprintf("Extension '%s' not found", ext.ID))
        }
    }
    
    // Validate inheritance
    if len(profile.Inherits) > 0 {
        if pv.hasCircularInheritance(profile) {
            result.AddError("Circular inheritance detected")
        }
    }
    
    // Validate auto-detect patterns
    if profile.AutoDetect != nil {
        for _, pattern := range profile.AutoDetect.Patterns {
            if _, err := filepath.Match(pattern, "test"); err != nil {
                result.AddError(fmt.Sprintf("Invalid pattern: %s", pattern))
            }
        }
    }
    
    return result
}
```

## Profile Switching

### Quick Switch Interface

```go
type ProfileSwitcher struct {
    manager   *ProfileManager
    ui        *ProfileSwitcherUI
    launcher  *GeminiLauncher
}

func (ps *ProfileSwitcher) ShowQuickSwitch() error {
    profiles, err := ps.manager.List()
    if err != nil {
        return err
    }
    
    // Sort by last used and usage count
    sort.Slice(profiles, func(i, j int) bool {
        if profiles[i].LastUsed != nil && profiles[j].LastUsed != nil {
            return profiles[i].LastUsed.After(*profiles[j].LastUsed)
        }
        return profiles[i].UsageCount > profiles[j].UsageCount
    })
    
    selected, err := ps.ui.SelectProfile(profiles)
    if err != nil {
        return err
    }
    
    return ps.Switch(selected.ID)
}

func (ps *ProfileSwitcher) Switch(profileID string) error {
    // 1. Deactivate current profile
    if err := ps.manager.Deactivate(); err != nil {
        return err
    }
    
    // 2. Activate new profile
    if err := ps.manager.Activate(profileID); err != nil {
        return err
    }
    
    // 3. Restart Gemini CLI if running
    if ps.launcher.IsRunning() {
        return ps.launcher.Restart()
    }
    
    return nil
}
```

### Auto-Detection

```go
type ProfileAutoDetector struct {
    manager  *ProfileManager
    patterns map[string]*CompiledPattern
}

func (pad *ProfileAutoDetector) DetectProfile(workingDir string) (*Profile, error) {
    profiles, err := pad.manager.List()
    if err != nil {
        return nil, err
    }
    
    var candidates []*ProfileMatch
    
    for _, profile := range profiles {
        if profile.AutoDetect == nil {
            continue
        }
        
        score := pad.calculateMatchScore(workingDir, profile.AutoDetect)
        if score > 0 {
            candidates = append(candidates, &ProfileMatch{
                Profile:  profile,
                Score:    score,
                Priority: profile.AutoDetect.Priority,
            })
        }
    }
    
    if len(candidates) == 0 {
        return nil, nil
    }
    
    // Sort by priority, then score
    sort.Slice(candidates, func(i, j int) bool {
        if candidates[i].Priority != candidates[j].Priority {
            return candidates[i].Priority > candidates[j].Priority
        }
        return candidates[i].Score > candidates[j].Score
    })
    
    return candidates[0].Profile, nil
}
```

## Profile Templates

### Built-in Templates

```go
var builtInTemplates = map[string]*ProfileTemplate{
    "web-frontend": {
        Name:        "Web Frontend",
        Description: "React/Vue/Angular development",
        Extensions:  []string{"typescript-tools", "eslint-helper", "prettier-format"},
        Environment: map[string]string{
            "NODE_ENV": "development",
        },
    },
    "python-data-science": {
        Name:        "Python Data Science",
        Description: "Jupyter, pandas, scikit-learn",
        Extensions:  []string{"python-tools", "jupyter-helper", "data-viz"},
        Environment: map[string]string{
            "PYTHONPATH": "./src",
        },
    },
    "go-backend": {
        Name:        "Go Backend",
        Description: "Go API development",
        Extensions:  []string{"go-tools", "database-tools", "api-testing"},
        Environment: map[string]string{
            "GO111MODULE": "on",
        },
    },
}

type ProfileTemplate struct {
    Name        string
    Description string
    Extensions  []string
    Environment map[string]string
    MCPServers  []string
}

func (pm *ProfileManager) CreateFromTemplate(templateID, profileID string) error {
    template, exists := builtInTemplates[templateID]
    if !exists {
        return fmt.Errorf("template '%s' not found", templateID)
    }
    
    profile := &Profile{
        ID:          profileID,
        Name:        template.Name,
        Description: template.Description,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }
    
    // Add extensions from template
    for _, extID := range template.Extensions {
        profile.Extensions = append(profile.Extensions, ExtensionRef{
            ID:      extID,
            Enabled: true,
        })
    }
    
    // Copy environment
    profile.Environment = make(map[string]string)
    for k, v := range template.Environment {
        profile.Environment[k] = v
    }
    
    return pm.Create(profile)
}
```

## Profile Sharing

### Export Format

```go
type ProfileExport struct {
    Version     string    `json:"version"`
    Profile     *Profile  `json:"profile"`
    Extensions  []string  `json:"required_extensions"`
    ExportedAt  time.Time `json:"exported_at"`
    ExportedBy  string    `json:"exported_by"`
    Checksum    string    `json:"checksum"`
}

func (pm *ProfileManager) Export(id string, path string) error {
    profile, err := pm.Get(id)
    if err != nil {
        return err
    }
    
    // Collect required extensions
    var requiredExts []string
    for _, ext := range profile.Extensions {
        requiredExts = append(requiredExts, ext.ID)
    }
    
    export := &ProfileExport{
        Version:    "1.0",
        Profile:    profile,
        Extensions: requiredExts,
        ExportedAt: time.Now(),
        ExportedBy: os.Getenv("USER"),
    }
    
    // Calculate checksum
    export.Checksum = pm.calculateChecksum(export)
    
    // Write to file
    data, err := json.MarshalIndent(export, "", "  ")
    if err != nil {
        return err
    }
    
    return os.WriteFile(path, data, 0644)
}
```

### Team Sharing

```go
type TeamProfileSync struct {
    repository string
    branch     string
    manager    *ProfileManager
}

func (tps *TeamProfileSync) Pull() error {
    // 1. Clone/pull repository
    // 2. Validate profiles
    // 3. Import non-conflicting profiles
    // 4. Report conflicts
}

func (tps *TeamProfileSync) Push(profileID string) error {
    // 1. Export profile
    // 2. Commit to repository
    // 3. Push to remote
}
```

## Profile Analytics

### Usage Tracking

```go
type ProfileAnalytics struct {
    storage AnalyticsStorage
}

type ProfileUsage struct {
    ProfileID    string
    SessionStart time.Time
    SessionEnd   time.Time
    Commands     []CommandLog
    Errors       []ErrorLog
}

func (pa *ProfileAnalytics) TrackUsage(profileID string) *UsageSession {
    return &UsageSession{
        profileID: profileID,
        startTime: time.Now(),
        analytics: pa,
    }
}

func (pa *ProfileAnalytics) GetInsights(profileID string) *ProfileInsights {
    usage := pa.storage.GetUsageData(profileID)
    
    return &ProfileInsights{
        TotalSessions:    len(usage),
        AverageSession:   pa.calculateAverageSession(usage),
        MostUsedFeatures: pa.extractTopFeatures(usage),
        ErrorRate:        pa.calculateErrorRate(usage),
        Recommendations:  pa.generateRecommendations(usage),
    }
}
```

## Profile Backup and Recovery

### Automatic Backups

```go
type ProfileBackup struct {
    manager  *ProfileManager
    storage  BackupStorage
    schedule *BackupSchedule
}

func (pb *ProfileBackup) BackupProfile(profileID string) error {
    profile, err := pb.manager.Get(profileID)
    if err != nil {
        return err
    }
    
    backup := &Backup{
        ProfileID: profileID,
        Profile:   profile,
        Timestamp: time.Now(),
        Version:   pb.calculateVersion(profileID),
    }
    
    return pb.storage.Store(backup)
}

func (pb *ProfileBackup) RestoreProfile(profileID string, version int) error {
    backup, err := pb.storage.Get(profileID, version)
    if err != nil {
        return err
    }
    
    // Create restoration point
    if err := pb.BackupProfile(profileID); err != nil {
        return err
    }
    
    return pb.manager.Update(profileID, backup.Profile)
}
```

### Conflict Resolution

```go
type ConflictResolver struct {
    strategies map[ConflictType]ResolutionStrategy
}

type ProfileConflict struct {
    Type     ConflictType
    Local    *Profile
    Remote   *Profile
    Details  map[string]interface{}
}

func (cr *ConflictResolver) Resolve(conflict ProfileConflict) (*Profile, error) {
    strategy, exists := cr.strategies[conflict.Type]
    if !exists {
        return nil, fmt.Errorf("no strategy for conflict type: %v", conflict.Type)
    }
    
    return strategy.Resolve(conflict)
}
```

## Performance Optimization

### Profile Caching

```go
type ProfileCache struct {
    profiles  *lru.Cache
    resolved  *lru.Cache
    analytics *lru.Cache
}

func NewProfileCache() *ProfileCache {
    return &ProfileCache{
        profiles:  lru.New(50),
        resolved:  lru.New(20),
        analytics: lru.New(100),
    }
}
```

### Lazy Loading

```go
type LazyProfile struct {
    id       string
    metadata *ProfileMetadata
    loader   ProfileLoader
    profile  *Profile
    mu       sync.RWMutex
}

func (lp *LazyProfile) Get() (*Profile, error) {
    lp.mu.RLock()
    if lp.profile != nil {
        lp.mu.RUnlock()
        return lp.profile, nil
    }
    lp.mu.RUnlock()
    
    return lp.load()
}
```

This comprehensive profile management system provides flexibility, security, and ease of use for managing different Gemini CLI configurations.