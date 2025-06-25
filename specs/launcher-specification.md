# Launcher Script Specification

## Overview

The Gemini CLI Launcher is a sophisticated shell script and Go binary combination that manages the lifecycle of Gemini CLI sessions with proper extension loading, environment configuration, and profile management.

## Architecture

### Components

1. **Shell Wrapper** (`gemini-launcher`) - Entry point script
2. **Go Binary** (`gemini-launcher-core`) - Core logic and heavy lifting
3. **Session Manager** - Tracks running Gemini instances
4. **Environment Builder** - Constructs proper environment
5. **Extension Loader** - Prepares extensions for use

## Shell Wrapper Script

### Basic Structure

```bash
#!/usr/bin/env bash
# gemini-launcher - Gemini CLI Launcher with Profile Support

set -euo pipefail

# Constants
readonly LAUNCHER_VERSION="1.0.0"
readonly LAUNCHER_HOME="${GEMINI_HOME:-$HOME/.gemini}"
readonly LAUNCHER_BIN="$LAUNCHER_HOME/bin/gemini-launcher-core"
readonly LAUNCHER_CONFIG="$LAUNCHER_HOME/config/launcher.yaml"
readonly LAUNCHER_LOG="$LAUNCHER_HOME/logs/launcher.log"

# Ensure launcher is installed
if [[ ! -x "$LAUNCHER_BIN" ]]; then
    echo "Error: Gemini launcher core not found at $LAUNCHER_BIN" >&2
    echo "Please run: gemini-cli-manager install-launcher" >&2
    exit 1
fi

# Parse command line arguments
PROFILE=""
COMMAND=""
DEBUG=false
DETACHED=false
AUTO_DETECT=true

while [[ $# -gt 0 ]]; do
    case $1 in
        -p|--profile)
            PROFILE="$2"
            AUTO_DETECT=false
            shift 2
            ;;
        -c|--command)
            COMMAND="$2"
            shift 2
            ;;
        -d|--debug)
            DEBUG=true
            shift
            ;;
        --detached)
            DETACHED=true
            shift
            ;;
        --no-auto-detect)
            AUTO_DETECT=false
            shift
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        -v|--version)
            echo "gemini-launcher version $LAUNCHER_VERSION"
            exit 0
            ;;
        *)
            # Pass through to Gemini CLI
            break
            ;;
    esac
done

# Execute launcher core
exec "$LAUNCHER_BIN" \
    ${PROFILE:+--profile "$PROFILE"} \
    ${COMMAND:+--command "$COMMAND"} \
    ${DEBUG:+--debug} \
    ${DETACHED:+--detached} \
    ${AUTO_DETECT:+--auto-detect} \
    -- "$@"
```

### Helper Functions

```bash
show_help() {
    cat << EOF
Gemini CLI Launcher v$LAUNCHER_VERSION

Usage: gemini-launcher [OPTIONS] [-- GEMINI_ARGS]

Options:
    -p, --profile NAME      Use specific profile
    -c, --command CMD       Run command and exit
    -d, --debug             Enable debug logging
    --detached              Run in background
    --no-auto-detect        Disable profile auto-detection
    -h, --help              Show this help
    -v, --version           Show version

Environment Variables:
    GEMINI_HOME             Base directory (default: ~/.gemini)
    GEMINI_PROFILE          Default profile name
    GEMINI_AUTO_DETECT      Enable auto-detection (true/false)

Examples:
    # Launch with default/auto-detected profile
    gemini-launcher

    # Launch with specific profile
    gemini-launcher -p web-development

    # Run single command
    gemini-launcher -c "analyze this code" file.ts

    # Pass arguments to Gemini CLI
    gemini-launcher -- --model gemini-2.0-flash

EOF
}

# Installation helper
install_launcher() {
    local install_dir="${1:-/usr/local/bin}"
    
    if [[ ! -w "$install_dir" ]]; then
        echo "Error: Cannot write to $install_dir" >&2
        echo "Try: sudo $0 install" >&2
        exit 1
    fi
    
    cp "$0" "$install_dir/gemini-launcher"
    chmod +x "$install_dir/gemini-launcher"
    
    # Create convenient aliases
    cat << 'EOF' > "$install_dir/gcm"
#!/bin/bash
exec gemini-launcher "$@"
EOF
    chmod +x "$install_dir/gcm"
    
    echo "Installed gemini-launcher to $install_dir"
    echo "You can now use 'gemini-launcher' or 'gcm' commands"
}

# Handle special commands
case "${1:-}" in
    install)
        install_launcher "${2:-}"
        exit 0
        ;;
esac
```

## Go Binary Core (`gemini-launcher-core`)

### Main Structure

```go
package main

import (
    "github.com/spf13/cobra"
    "github.com/gemini-cli/launcher/internal/config"
    "github.com/gemini-cli/launcher/internal/launcher"
    "github.com/gemini-cli/launcher/internal/profile"
    "github.com/gemini-cli/launcher/internal/session"
)

type LauncherOptions struct {
    Profile      string
    Command      string
    Debug        bool
    Detached     bool
    AutoDetect   bool
    GeminiArgs   []string
}

func main() {
    var opts LauncherOptions
    
    rootCmd := &cobra.Command{
        Use:   "gemini-launcher-core",
        Short: "Core launcher for Gemini CLI",
        RunE: func(cmd *cobra.Command, args []string) error {
            return runLauncher(opts, args)
        },
    }
    
    flags := rootCmd.Flags()
    flags.StringVarP(&opts.Profile, "profile", "p", "", "Profile to use")
    flags.StringVarP(&opts.Command, "command", "c", "", "Command to execute")
    flags.BoolVarP(&opts.Debug, "debug", "d", false, "Enable debug mode")
    flags.BoolVar(&opts.Detached, "detached", false, "Run detached")
    flags.BoolVar(&opts.AutoDetect, "auto-detect", true, "Auto-detect profile")
    
    // Collect remaining args after --
    rootCmd.SetArgs(os.Args[1:])
    rootCmd.Execute()
}
```

### Launcher Implementation

```go
type Launcher struct {
    config         *config.Config
    profileManager *profile.Manager
    sessionManager *session.Manager
    envBuilder     *environment.Builder
    extensionMgr   *extension.Manager
}

func (l *Launcher) Launch(opts LauncherOptions) error {
    // 1. Determine profile
    activeProfile, err := l.determineProfile(opts)
    if err != nil {
        return fmt.Errorf("determining profile: %w", err)
    }
    
    // 2. Validate profile
    if err := l.validateProfile(activeProfile); err != nil {
        return fmt.Errorf("invalid profile: %w", err)
    }
    
    // 3. Prepare environment
    env, err := l.prepareEnvironment(activeProfile)
    if err != nil {
        return fmt.Errorf("preparing environment: %w", err)
    }
    
    // 4. Start MCP servers
    servers, err := l.startMCPServers(activeProfile)
    if err != nil {
        return fmt.Errorf("starting MCP servers: %w", err)
    }
    defer l.stopMCPServers(servers)
    
    // 5. Create session
    sess, err := l.createSession(activeProfile, opts)
    if err != nil {
        return fmt.Errorf("creating session: %w", err)
    }
    
    // 6. Launch Gemini CLI
    if opts.Detached {
        return l.launchDetached(sess, env, opts)
    }
    return l.launchInteractive(sess, env, opts)
}

func (l *Launcher) determineProfile(opts LauncherOptions) (*profile.Profile, error) {
    // Priority order:
    // 1. Command line flag
    // 2. Environment variable
    // 3. Auto-detection
    // 4. Last used
    // 5. Default
    
    if opts.Profile != "" {
        return l.profileManager.Get(opts.Profile)
    }
    
    if envProfile := os.Getenv("GEMINI_PROFILE"); envProfile != "" {
        return l.profileManager.Get(envProfile)
    }
    
    if opts.AutoDetect {
        if detected, err := l.autoDetectProfile(); err == nil && detected != nil {
            log.Debug("Auto-detected profile: %s", detected.Name)
            return detected, nil
        }
    }
    
    if lastUsed, err := l.profileManager.GetLastUsed(); err == nil {
        return lastUsed, nil
    }
    
    return l.profileManager.GetDefault()
}
```

### Environment Preparation

```go
type EnvironmentBuilder struct {
    baseEnv      map[string]string
    profileMgr   *profile.Manager
    credentialMgr *credential.Manager
}

func (eb *EnvironmentBuilder) Build(p *profile.Profile) ([]string, error) {
    env := make(map[string]string)
    
    // 1. Start with system environment (filtered)
    for _, e := range os.Environ() {
        if eb.isAllowedSystemVar(e) {
            parts := strings.SplitN(e, "=", 2)
            if len(parts) == 2 {
                env[parts[0]] = parts[1]
            }
        }
    }
    
    // 2. Apply profile environment
    for k, v := range p.Environment {
        env[k] = eb.expandVariables(v, env)
    }
    
    // 3. Inject credentials
    for _, credID := range p.RequiredCredentials {
        value, err := eb.credentialMgr.GetDecrypted(credID)
        if err != nil {
            return nil, fmt.Errorf("retrieving credential %s: %w", credID, err)
        }
        env[credID] = value
    }
    
    // 4. Add Gemini-specific variables
    env["GEMINI_PROFILE"] = p.Name
    env["GEMINI_PROFILE_ID"] = p.ID
    env["GEMINI_EXTENSIONS_PATH"] = filepath.Join(eb.baseEnv["GEMINI_HOME"], "extensions")
    env["GEMINI_MCP_SERVERS"] = eb.buildMCPServerList(p)
    
    // 5. Convert to slice
    var result []string
    for k, v := range env {
        result = append(result, fmt.Sprintf("%s=%s", k, v))
    }
    
    return result, nil
}

func (eb *EnvironmentBuilder) expandVariables(value string, env map[string]string) string {
    return os.Expand(value, func(key string) string {
        if val, ok := env[key]; ok {
            return val
        }
        return os.Getenv(key)
    })
}
```

### MCP Server Management

```go
type MCPServerManager struct {
    servers  map[string]*MCPServer
    monitor  *HealthMonitor
    launcher *ProcessLauncher
}

type MCPServer struct {
    ID          string
    Name        string
    Config      ServerConfig
    Process     *os.Process
    Port        int
    HealthCheck HealthChecker
    Status      ServerStatus
}

func (msm *MCPServerManager) StartServers(profile *profile.Profile) ([]*MCPServer, error) {
    var started []*MCPServer
    var errors []error
    
    // Start servers in parallel with proper ordering
    dag := msm.buildDependencyGraph(profile.MCPServers)
    groups := dag.TopologicalSort()
    
    for _, group := range groups {
        var wg sync.WaitGroup
        var mu sync.Mutex
        
        for _, serverID := range group {
            wg.Add(1)
            go func(id string) {
                defer wg.Done()
                
                server, err := msm.startServer(id, profile.MCPServers[id])
                mu.Lock()
                defer mu.Unlock()
                
                if err != nil {
                    errors = append(errors, fmt.Errorf("starting %s: %w", id, err))
                } else {
                    started = append(started, server)
                }
            }(serverID)
        }
        
        wg.Wait()
        
        if len(errors) > 0 {
            // Clean up started servers
            msm.stopServers(started)
            return nil, fmt.Errorf("failed to start servers: %v", errors)
        }
    }
    
    // Wait for all servers to be healthy
    if err := msm.waitForHealthy(started, 30*time.Second); err != nil {
        msm.stopServers(started)
        return nil, err
    }
    
    return started, nil
}

func (msm *MCPServerManager) startServer(id string, config ServerConfig) (*MCPServer, error) {
    // Find available port
    port, err := msm.findAvailablePort(config.PreferredPort)
    if err != nil {
        return nil, err
    }
    
    // Prepare command
    cmd := exec.Command(config.Command, config.Args...)
    cmd.Env = append(os.Environ(), 
        fmt.Sprintf("MCP_SERVER_PORT=%d", port),
        fmt.Sprintf("MCP_SERVER_ID=%s", id),
    )
    
    // Set up logging
    logFile, err := msm.createLogFile(id)
    if err != nil {
        return nil, err
    }
    cmd.Stdout = logFile
    cmd.Stderr = logFile
    
    // Start process
    if err := cmd.Start(); err != nil {
        return nil, err
    }
    
    server := &MCPServer{
        ID:      id,
        Name:    config.DisplayName,
        Config:  config,
        Process: cmd.Process,
        Port:    port,
        Status:  ServerStarting,
    }
    
    msm.servers[id] = server
    msm.monitor.Track(server)
    
    return server, nil
}
```

### Session Management

```go
type SessionManager struct {
    activeSessions map[string]*Session
    store          SessionStore
    monitor        *SessionMonitor
}

type Session struct {
    ID           string
    ProfileID    string
    StartTime    time.Time
    Environment  []string
    MCPServers   []MCPServerInfo
    GeminiPID    int
    Status       SessionStatus
    Detached     bool
}

func (sm *SessionManager) CreateSession(
    profile *profile.Profile,
    env []string,
    servers []*MCPServer,
) (*Session, error) {
    session := &Session{
        ID:        generateSessionID(),
        ProfileID: profile.ID,
        StartTime: time.Now(),
        Environment: env,
        Status:    SessionStarting,
    }
    
    // Record MCP server info
    for _, server := range servers {
        session.MCPServers = append(session.MCPServers, MCPServerInfo{
            ID:   server.ID,
            Port: server.Port,
            PID:  server.Process.Pid,
        })
    }
    
    // Persist session
    if err := sm.store.Save(session); err != nil {
        return nil, err
    }
    
    sm.activeSessions[session.ID] = session
    return session, nil
}

func (sm *SessionManager) MonitorSession(session *Session) {
    sm.monitor.Watch(session, func(event SessionEvent) {
        switch event.Type {
        case SessionEventStopped:
            sm.handleSessionStop(session)
        case SessionEventCrashed:
            sm.handleSessionCrash(session)
        case SessionEventHung:
            sm.handleSessionHung(session)
        }
    })
}
```

### Launch Modes

```go
func (l *Launcher) launchInteractive(
    sess *Session,
    env []string,
    opts LauncherOptions,
) error {
    // Build Gemini CLI command
    geminiPath := l.config.GeminiCLIPath
    if geminiPath == "" {
        geminiPath = "gemini"
    }
    
    args := []string{}
    
    // Add session tracking
    args = append(args, "--session-id", sess.ID)
    
    // Add command if provided
    if opts.Command != "" {
        args = append(args, "exec", opts.Command)
    }
    
    // Add user-provided args
    args = append(args, opts.GeminiArgs...)
    
    // Create command
    cmd := exec.Command(geminiPath, args...)
    cmd.Env = env
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    
    // Start Gemini CLI
    if err := cmd.Start(); err != nil {
        return fmt.Errorf("starting Gemini CLI: %w", err)
    }
    
    sess.GeminiPID = cmd.Process.Pid
    sess.Status = SessionRunning
    l.sessionManager.UpdateSession(sess)
    
    // Monitor and wait
    go l.sessionManager.MonitorSession(sess)
    
    // Wait for completion
    err := cmd.Wait()
    
    // Update session status
    if err != nil {
        sess.Status = SessionFailed
    } else {
        sess.Status = SessionCompleted
    }
    l.sessionManager.UpdateSession(sess)
    
    return err
}

func (l *Launcher) launchDetached(
    sess *Session,
    env []string,
    opts LauncherOptions,
) error {
    // Create detached process
    cmd := l.buildGeminiCommand(sess, opts)
    cmd.Env = env
    
    // Redirect output to log files
    logDir := filepath.Join(l.config.LogDir, sess.ID)
    os.MkdirAll(logDir, 0755)
    
    stdout, err := os.Create(filepath.Join(logDir, "stdout.log"))
    if err != nil {
        return err
    }
    stderr, err := os.Create(filepath.Join(logDir, "stderr.log"))
    if err != nil {
        return err
    }
    
    cmd.Stdout = stdout
    cmd.Stderr = stderr
    
    // Detach from current process group
    cmd.SysProcAttr = &syscall.SysProcAttr{
        Setpgid: true,
    }
    
    // Start process
    if err := cmd.Start(); err != nil {
        return err
    }
    
    sess.GeminiPID = cmd.Process.Pid
    sess.Status = SessionRunning
    sess.Detached = true
    l.sessionManager.UpdateSession(sess)
    
    // Monitor in background
    go l.sessionManager.MonitorSession(sess)
    
    fmt.Printf("Started Gemini CLI in background\n")
    fmt.Printf("Session ID: %s\n", sess.ID)
    fmt.Printf("Logs: %s\n", logDir)
    fmt.Printf("\nTo attach: gemini-launcher attach %s\n", sess.ID)
    fmt.Printf("To stop: gemini-launcher stop %s\n", sess.ID)
    
    return nil
}
```

### Auto-Detection

```go
type ProfileDetector struct {
    patterns map[string]*DetectionPattern
    scorer   *MatchScorer
}

type DetectionPattern struct {
    Files      []string
    Indicators map[string]string
    Weight     float64
}

func (pd *ProfileDetector) Detect(workDir string) (*profile.Profile, error) {
    var matches []ProfileMatch
    
    // Check each profile's detection patterns
    profiles, err := pd.profileManager.ListAll()
    if err != nil {
        return nil, err
    }
    
    for _, p := range profiles {
        if p.AutoDetect == nil {
            continue
        }
        
        score := pd.calculateScore(workDir, p.AutoDetect)
        if score > 0 {
            matches = append(matches, ProfileMatch{
                Profile: p,
                Score:   score,
            })
        }
    }
    
    if len(matches) == 0 {
        return nil, nil
    }
    
    // Sort by score
    sort.Slice(matches, func(i, j int) bool {
        return matches[i].Score > matches[j].Score
    })
    
    // Log detection results
    log.Debug("Profile detection results:")
    for _, m := range matches {
        log.Debug("  %s: %.2f", m.Profile.Name, m.Score)
    }
    
    return matches[0].Profile, nil
}

func (pd *ProfileDetector) calculateScore(
    workDir string,
    rules *AutoDetectRules,
) float64 {
    score := 0.0
    
    // Check file patterns
    for _, pattern := range rules.Patterns {
        matches, _ := filepath.Glob(filepath.Join(workDir, pattern))
        if len(matches) > 0 {
            score += 1.0
        }
    }
    
    // Check for .gemini-profile file
    profileFile := filepath.Join(workDir, ".gemini-profile")
    if data, err := os.ReadFile(profileFile); err == nil {
        if string(data) == rules.ProfileName {
            score += 10.0 // Strong indicator
        }
    }
    
    // Apply priority weight
    score *= float64(rules.Priority)
    
    return score
}
```

### Cleanup and Recovery

```go
type CleanupManager struct {
    sessionMgr *SessionManager
    serverMgr  *MCPServerManager
}

func (cm *CleanupManager) CleanupSession(sessionID string) error {
    session, err := cm.sessionMgr.GetSession(sessionID)
    if err != nil {
        return err
    }
    
    // Stop Gemini CLI process
    if session.GeminiPID > 0 {
        if proc, err := os.FindProcess(session.GeminiPID); err == nil {
            proc.Signal(syscall.SIGTERM)
            
            // Wait for graceful shutdown
            done := make(chan bool)
            go func() {
                proc.Wait()
                done <- true
            }()
            
            select {
            case <-done:
                // Graceful shutdown succeeded
            case <-time.After(10 * time.Second):
                // Force kill
                proc.Kill()
            }
        }
    }
    
    // Stop MCP servers
    for _, serverInfo := range session.MCPServers {
        cm.serverMgr.StopServer(serverInfo.ID)
    }
    
    // Update session status
    session.Status = SessionStopped
    cm.sessionMgr.UpdateSession(session)
    
    return nil
}

func (cm *CleanupManager) RecoverOrphanedSessions() error {
    sessions, err := cm.sessionMgr.GetActiveSessions()
    if err != nil {
        return err
    }
    
    for _, session := range sessions {
        // Check if process still exists
        if !cm.isProcessRunning(session.GeminiPID) {
            log.Info("Found orphaned session: %s", session.ID)
            cm.CleanupSession(session.ID)
        }
    }
    
    return nil
}
```

## Additional Commands

### Session Management Commands

```go
// List active sessions
type ListCommand struct{}

func (c *ListCommand) Run() error {
    sessions, err := sessionManager.GetActiveSessions()
    if err != nil {
        return err
    }
    
    table := tablewriter.NewWriter(os.Stdout)
    table.SetHeader([]string{"Session ID", "Profile", "Status", "Started", "PID"})
    
    for _, sess := range sessions {
        table.Append([]string{
            sess.ID[:8],
            sess.ProfileName,
            string(sess.Status),
            sess.StartTime.Format("15:04:05"),
            fmt.Sprintf("%d", sess.GeminiPID),
        })
    }
    
    table.Render()
    return nil
}

// Attach to detached session
type AttachCommand struct {
    SessionID string
}

func (c *AttachCommand) Run() error {
    session, err := sessionManager.GetSession(c.SessionID)
    if err != nil {
        return err
    }
    
    if !session.Detached {
        return fmt.Errorf("session is not detached")
    }
    
    // Use screen/tmux to attach
    return attachToSession(session)
}

// Stop a session
type StopCommand struct {
    SessionID string
    Force     bool
}

func (c *StopCommand) Run() error {
    return cleanupManager.CleanupSession(c.SessionID)
}
```

## Integration Points

### Shell Integration

```bash
# Bash completion
_gemini_launcher_completions() {
    local cur="${COMP_WORDS[COMP_CWORD]}"
    local prev="${COMP_WORDS[COMP_CWORD-1]}"
    
    case "$prev" in
        -p|--profile)
            COMPREPLY=($(compgen -W "$(gemini-launcher list-profiles)" -- "$cur"))
            ;;
        *)
            COMPREPLY=($(compgen -W "--profile --command --debug --help" -- "$cur"))
            ;;
    esac
}

complete -F _gemini_launcher_completions gemini-launcher gcm
```

### IDE Integration

```json
// VS Code tasks.json
{
    "version": "2.0.0",
    "tasks": [
        {
            "label": "Launch Gemini CLI",
            "type": "shell",
            "command": "gemini-launcher",
            "args": [
                "--profile", "${input:profile}"
            ],
            "problemMatcher": []
        }
    ],
    "inputs": [
        {
            "id": "profile",
            "type": "pickString",
            "description": "Select Gemini profile",
            "options": ["web-development", "data-science", "go-backend"]
        }
    ]
}
```

## Error Handling

### Common Error Scenarios

1. **Missing Dependencies**
   - Check for Gemini CLI installation
   - Verify extension dependencies
   - Validate MCP server binaries

2. **Port Conflicts**
   - Automatically find available ports
   - Report conflicts clearly
   - Provide manual override options

3. **Permission Issues**
   - Check file permissions
   - Validate credential access
   - Handle elevated privileges gracefully

4. **Network Issues**
   - Timeout handling for MCP servers
   - Retry logic with backoff
   - Offline mode support

## Performance Optimizations

1. **Parallel Startup**
   - Start MCP servers concurrently
   - Pre-warm frequently used servers
   - Cache environment preparations

2. **Resource Management**
   - Monitor memory usage
   - CPU throttling for background processes
   - Automatic cleanup of idle servers

3. **Fast Profile Switching**
   - Keep recent profiles in memory
   - Lazy load extension metadata
   - Incremental environment updates

This comprehensive launcher specification provides a robust, efficient, and user-friendly way to manage Gemini CLI sessions with full profile and extension support.