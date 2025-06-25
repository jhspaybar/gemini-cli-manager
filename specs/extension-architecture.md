# Extension Architecture Specification

## Overview

The extension system is the core of the Gemini CLI Manager, providing a flexible and secure way to manage MCP servers, custom tools, and configurations.

## Extension Structure

### Directory Layout
```
~/.gemini/extensions/
└── extension-name/
    ├── settings.json       # Required: Extension configuration
    ├── GEMINI.md          # Required: Documentation and guidance
    ├── mcp-servers/       # Optional: MCP server implementations
    │   ├── server1/
    │   └── server2/
    ├── tools/             # Optional: Custom tools
    ├── templates/         # Optional: Code templates
    └── resources/         # Optional: Additional resources
```

### settings.json Schema
```json
{
  "$schema": "https://gemini-cli.dev/schemas/extension-v1.json",
  "id": "extension-unique-id",
  "name": "Extension Display Name",
  "version": "1.0.0",
  "description": "Brief description of the extension",
  "author": {
    "name": "Author Name",
    "email": "author@example.com",
    "url": "https://github.com/author"
  },
  "repository": {
    "type": "git",
    "url": "https://github.com/author/extension.git"
  },
  "engines": {
    "gemini-cli": ">=1.0.0",
    "node": ">=18.0.0"
  },
  "categories": ["productivity", "development", "ai"],
  "keywords": ["mcp", "typescript", "tools"],
  "dependencies": {
    "other-extension-id": "^1.0.0"
  },
  "mcp_servers": {
    "server-name": {
      "displayName": "Server Display Name",
      "description": "What this server does",
      "command": "node",
      "args": ["./mcp-servers/server1/index.js"],
      "env": {
        "DEFAULT_MODEL": "gemini-2.0-flash"
      },
      "requiredEnv": ["API_KEY"],
      "capabilities": ["code-analysis", "text-generation"]
    }
  },
  "tools": {
    "tool-name": {
      "displayName": "Tool Display Name",
      "description": "What this tool does",
      "command": "python",
      "args": ["./tools/tool.py"],
      "input": "stdin",
      "output": "stdout"
    }
  },
  "configuration": {
    "properties": {
      "apiEndpoint": {
        "type": "string",
        "description": "API endpoint URL",
        "default": "https://api.example.com"
      },
      "timeout": {
        "type": "number",
        "description": "Request timeout in seconds",
        "default": 30,
        "minimum": 1,
        "maximum": 300
      }
    }
  },
  "activation": {
    "events": ["onStartup", "onCommand:extension.activate"]
  },
  "contributes": {
    "commands": [
      {
        "command": "extension.doSomething",
        "title": "Do Something Cool"
      }
    ],
    "keybindings": [
      {
        "command": "extension.doSomething",
        "key": "ctrl+shift+d"
      }
    ]
  }
}
```

### GEMINI.md Structure
```markdown
# Extension Name

## Overview
Brief description of what this extension provides.

## Features
- Feature 1: Description
- Feature 2: Description

## MCP Servers
### server-name
Description of the server and its capabilities.

**Required Environment Variables:**
- `API_KEY`: Your API key from...

**Example Usage:**
\```
gemini "Analyze this code for security issues"
\```

## Configuration
### Setting Name
- Type: string
- Default: "value"
- Description: What this setting controls

## Troubleshooting
Common issues and solutions.

## Contributing
How to contribute to this extension.
```

## Extension Lifecycle

### Discovery Phase
```go
type ExtensionDiscovery struct {
    scanner    *ExtensionScanner
    validator  *ExtensionValidator
    registry   *ExtensionRegistry
}

func (ed *ExtensionDiscovery) ScanDirectory(path string) ([]*Extension, error) {
    // 1. Find all subdirectories
    // 2. Check for settings.json
    // 3. Validate basic structure
    // 4. Register found extensions
}
```

### Validation Phase
```go
type ValidationResult struct {
    Valid    bool
    Errors   []ValidationError
    Warnings []ValidationWarning
}

type ExtensionValidator struct {
    schemaValidator *jsonschema.Validator
    fsChecker       *FileSystemChecker
}

func (v *ExtensionValidator) Validate(ext *Extension) ValidationResult {
    // 1. Validate settings.json against schema
    // 2. Check required files exist
    // 3. Validate MCP server configurations
    // 4. Check dependencies
    // 5. Verify no conflicts
}
```

### Loading Phase
```go
type ExtensionLoader struct {
    registry    *ExtensionRegistry
    resolver    *DependencyResolver
    configMgr   *ConfigurationManager
}

func (l *ExtensionLoader) Load(id string) (*LoadedExtension, error) {
    // 1. Resolve dependencies
    // 2. Load configuration
    // 3. Initialize MCP servers
    // 4. Register tools
    // 5. Apply activation events
}
```

### Runtime Management
```go
type ExtensionRuntime struct {
    loaded    map[string]*LoadedExtension
    processes map[string]*MCPProcess
    monitor   *HealthMonitor
}

func (r *ExtensionRuntime) Enable(id string) error {
    // 1. Load extension if not loaded
    // 2. Start MCP servers
    // 3. Register tools
    // 4. Update active profile
}

func (r *ExtensionRuntime) Disable(id string) error {
    // 1. Stop MCP servers gracefully
    // 2. Unregister tools
    // 3. Clean up resources
    // 4. Update active profile
}
```

## Dependency Resolution

### Dependency Graph
```go
type DependencyGraph struct {
    nodes map[string]*Extension
    edges map[string][]string
}

func (g *DependencyGraph) Resolve(rootID string) ([]string, error) {
    // Topological sort to determine load order
    // Detect circular dependencies
    // Return ordered list of extension IDs
}
```

### Version Compatibility
```go
type VersionConstraint struct {
    operator string // ">=", "^", "~", "="
    version  *semver.Version
}

func CheckCompatibility(required, actual string) (bool, error) {
    // Parse version constraint
    // Compare with actual version
    // Return compatibility status
}
```

## MCP Server Management

### Process Lifecycle
```go
type MCPProcess struct {
    ID          string
    Command     string
    Args        []string
    Env         map[string]string
    Process     *os.Process
    State       ProcessState
    HealthCheck HealthChecker
}

type ProcessManager struct {
    processes map[string]*MCPProcess
    monitor   *ProcessMonitor
}

func (pm *ProcessManager) Start(server MCPServerConfig) error {
    // 1. Prepare environment
    // 2. Start process
    // 3. Wait for ready signal
    // 4. Register with monitor
}

func (pm *ProcessManager) Stop(id string) error {
    // 1. Send graceful shutdown signal
    // 2. Wait for process to exit
    // 3. Force kill if timeout
    // 4. Clean up resources
}
```

### Health Monitoring
```go
type HealthChecker interface {
    Check(ctx context.Context) HealthStatus
}

type HealthStatus struct {
    Healthy       bool
    ResponseTime  time.Duration
    ErrorMessage  string
    LastChecked   time.Time
}

func (pm *ProcessMonitor) MonitorHealth(ctx context.Context) {
    // Periodic health checks
    // Restart unhealthy processes
    // Alert on repeated failures
}
```

## Security Model

### Extension Sandboxing
```go
type SecurityPolicy struct {
    AllowedPaths      []string
    AllowedCommands   []string
    NetworkAccess     bool
    EnvironmentVars   []string
}

func (sp *SecurityPolicy) Validate(ext *Extension) error {
    // Check extension doesn't exceed permissions
    // Validate command allowlist
    // Ensure no path traversal
}
```

### Credential Management
```go
type CredentialStore interface {
    Store(key string, value []byte) error
    Retrieve(key string) ([]byte, error)
    Delete(key string) error
}

type KeychainStore struct {
    serviceName string
}

func (ks *KeychainStore) Store(key string, value []byte) error {
    // Use OS keychain APIs
    // Encrypt if necessary
    // Set appropriate permissions
}
```

## Extension Registry

### Local Registry
```go
type LocalRegistry struct {
    extensions map[string]*Extension
    index      *SearchIndex
    cache      *RegistryCache
}

func (r *LocalRegistry) Search(query string) []*Extension {
    // Full-text search across metadata
    // Filter by categories
    // Sort by relevance
}
```

### Remote Registry (Future)
```go
type RemoteRegistry interface {
    Search(query string) ([]*Extension, error)
    Download(id string, version string) (io.Reader, error)
    GetMetadata(id string) (*ExtensionMetadata, error)
}
```

## Error Handling

### Error Types
```go
type ExtensionError struct {
    Type    ErrorType
    Message string
    Context map[string]interface{}
}

const (
    ErrorTypeValidation ErrorType = iota
    ErrorTypeDependency
    ErrorTypeRuntime
    ErrorTypeConfiguration
    ErrorTypeSecurity
)
```

### Recovery Strategies
```go
type RecoveryStrategy interface {
    Recover(err ExtensionError) error
}

type AutoRecovery struct {
    strategies map[ErrorType]RecoveryStrategy
}

func (ar *AutoRecovery) Handle(err ExtensionError) error {
    if strategy, exists := ar.strategies[err.Type]; exists {
        return strategy.Recover(err)
    }
    return err
}
```

## Performance Optimization

### Lazy Loading
```go
type LazyExtension struct {
    metadata *ExtensionMetadata
    loader   func() (*Extension, error)
    loaded   *Extension
    mu       sync.RWMutex
}

func (le *LazyExtension) Get() (*Extension, error) {
    le.mu.RLock()
    if le.loaded != nil {
        le.mu.RUnlock()
        return le.loaded, nil
    }
    le.mu.RUnlock()
    
    le.mu.Lock()
    defer le.mu.Unlock()
    
    if le.loaded == nil {
        ext, err := le.loader()
        if err != nil {
            return nil, err
        }
        le.loaded = ext
    }
    return le.loaded, nil
}
```

### Caching Strategy
```go
type ExtensionCache struct {
    metadata  *lru.Cache
    manifests *lru.Cache
    icons     *lru.Cache
}

func NewExtensionCache() *ExtensionCache {
    return &ExtensionCache{
        metadata:  lru.New(100),
        manifests: lru.New(50),
        icons:     lru.New(200),
    }
}
```

## Testing Framework

### Extension Testing
```go
type ExtensionTestSuite struct {
    extension *Extension
    harness   *TestHarness
}

func (ts *ExtensionTestSuite) TestValidation(t *testing.T) {
    // Test settings.json validity
    // Test file structure
    // Test dependencies
}

func (ts *ExtensionTestSuite) TestMCPServers(t *testing.T) {
    // Test server startup
    // Test communication
    // Test graceful shutdown
}
```

## Migration Support

### Version Migration
```go
type MigrationStep struct {
    FromVersion string
    ToVersion   string
    Migrate     func(old, new *Extension) error
}

type Migrator struct {
    steps []MigrationStep
}

func (m *Migrator) Migrate(ext *Extension, targetVersion string) error {
    // Find migration path
    // Execute migrations in order
    // Validate result
}
```

This architecture provides a robust, secure, and extensible foundation for managing Gemini CLI extensions.