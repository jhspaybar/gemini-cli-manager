# Security and Credential Management Specification

## Security Philosophy

The Gemini CLI Manager follows the principle of **Defense in Depth** with multiple layers of security:

1. **Least Privilege**: Extensions and MCP servers run with minimal required permissions
2. **Secure by Default**: Safe defaults with opt-in for advanced features
3. **Transparent Security**: Users understand what permissions they're granting
4. **Zero Trust**: Verify all inputs and validate all operations
5. **Secure Storage**: Credentials never stored in plain text

## Threat Model

### Identified Threats

1. **Malicious Extensions**
   - Code injection through settings.json
   - Path traversal in file operations
   - Command injection in MCP server execution
   - Data exfiltration through network access

2. **Credential Theft**
   - API keys exposed in configuration files
   - Environment variables leaked in logs
   - Credentials transmitted insecurely
   - Weak encryption of stored secrets

3. **Privilege Escalation**
   - Extensions accessing unauthorized resources
   - MCP servers escaping sandboxes
   - Profile manipulation to gain access

4. **Supply Chain Attacks**
   - Compromised extension repositories
   - Dependency vulnerabilities
   - Update mechanism hijacking

## Security Architecture

### Extension Sandboxing

```go
type SecuritySandbox struct {
    permissions  PermissionSet
    filesystem   *RestrictedFS
    network      *NetworkPolicy
    process      *ProcessPolicy
}

type PermissionSet struct {
    FileSystem   FilePermissions
    Network      NetworkPermissions
    Process      ProcessPermissions
    Environment  []string // Allowed env vars
}

type FilePermissions struct {
    ReadPaths    []string // Allowed read paths
    WritePaths   []string // Allowed write paths
    ExecutePaths []string // Allowed executable paths
    Restrictions FileRestrictions
}

type FileRestrictions struct {
    MaxFileSize      int64
    AllowedExtensions []string
    BlockedPaths     []string // Always blocked
}

func (sb *SecuritySandbox) ValidateOperation(op Operation) error {
    switch op.Type {
    case OpFileRead:
        return sb.validateFileRead(op.Path)
    case OpFileWrite:
        return sb.validateFileWrite(op.Path)
    case OpNetworkRequest:
        return sb.validateNetworkAccess(op.URL)
    case OpProcessSpawn:
        return sb.validateProcessSpawn(op.Command)
    }
    return ErrOperationNotAllowed
}
```

### Permission Model

```go
type ExtensionPermissions struct {
    // File System
    ReadFiles       []PathPattern `json:"read_files"`
    WriteFiles      []PathPattern `json:"write_files"`
    
    // Network
    NetworkAccess   bool          `json:"network_access"`
    AllowedHosts    []string      `json:"allowed_hosts"`
    
    // Process
    ExecuteCommands bool          `json:"execute_commands"`
    AllowedCommands []string      `json:"allowed_commands"`
    
    // Environment
    EnvironmentVars []string      `json:"environment_vars"`
    
    // Special
    SystemAPIs      []string      `json:"system_apis"`
}

type PathPattern struct {
    Pattern string
    Scope   PathScope
}

type PathScope string

const (
    ScopeUser      PathScope = "user"      // User home directory
    ScopeExtension PathScope = "extension" // Extension directory
    ScopeTemp      PathScope = "temp"      // Temporary files
    ScopeSystem    PathScope = "system"    // System paths (restricted)
)

// Example permission request in settings.json
{
    "permissions": {
        "read_files": [
            {"pattern": "**/*.ts", "scope": "user"},
            {"pattern": ".gemini/config.json", "scope": "user"}
        ],
        "network_access": true,
        "allowed_hosts": ["api.github.com", "*.googleapis.com"],
        "environment_vars": ["GEMINI_API_KEY", "GITHUB_TOKEN"]
    }
}
```

### Secure Process Execution

```go
type SecureProcessExecutor struct {
    sandbox     *SecuritySandbox
    monitor     *ProcessMonitor
    rateLimiter *RateLimiter
}

func (spe *SecureProcessExecutor) Execute(cmd ProcessCommand) (*Process, error) {
    // 1. Validate command against whitelist
    if err := spe.validateCommand(cmd); err != nil {
        return nil, err
    }
    
    // 2. Prepare sandboxed environment
    env := spe.prepareSandboxedEnv(cmd.Env)
    
    // 3. Set resource limits
    limits := &ResourceLimits{
        MaxMemory:     512 * 1024 * 1024, // 512MB
        MaxCPUPercent: 50,
        MaxFileSize:   10 * 1024 * 1024,  // 10MB
        MaxOpenFiles:  100,
    }
    
    // 4. Create process with restrictions
    proc := &Process{
        Command:   cmd.Command,
        Args:      cmd.Args,
        Env:       env,
        Dir:       spe.sandbox.filesystem.SandboxPath(cmd.Dir),
        Limits:    limits,
        Isolation: IsolationStrict,
    }
    
    // 5. Start with monitoring
    if err := proc.Start(); err != nil {
        return nil, err
    }
    
    spe.monitor.Track(proc)
    return proc, nil
}

type ProcessIsolation int

const (
    IsolationNone ProcessIsolation = iota
    IsolationBasic    // Restricted env, working directory
    IsolationStrict   // Namespace isolation, cgroups
    IsolationComplete // Full container isolation
)
```

## Credential Management

### Credential Storage Architecture

```go
type CredentialManager struct {
    store      CredentialStore
    encryptor  Encryptor
    validator  CredentialValidator
    auditor    AuditLogger
}

type CredentialStore interface {
    Store(key string, credential *Credential) error
    Retrieve(key string) (*Credential, error)
    Delete(key string) error
    List() ([]string, error)
}

type Credential struct {
    ID          string
    Type        CredentialType
    Value       []byte // Always encrypted
    Metadata    CredentialMetadata
    Permissions []string
}

type CredentialMetadata struct {
    CreatedAt    time.Time
    UpdatedAt    time.Time
    LastAccessed time.Time
    AccessCount  int
    ExpiresAt    *time.Time
    Description  string
}
```

### Platform-Specific Credential Stores

```go
// macOS Keychain implementation
type KeychainStore struct {
    service string
}

func (ks *KeychainStore) Store(key string, cred *Credential) error {
    item := keychain.NewItem()
    item.SetSecClass(keychain.SecClassGenericPassword)
    item.SetService(ks.service)
    item.SetAccount(key)
    item.SetData(cred.Value)
    item.SetAccessible(keychain.AccessibleWhenUnlocked)
    
    return keychain.AddItem(item)
}

// Linux Secret Service implementation
type SecretServiceStore struct {
    collection string
}

func (sss *SecretServiceStore) Store(key string, cred *Credential) error {
    service, err := secretservice.NewService()
    if err != nil {
        return err
    }
    
    session, err := service.OpenSession()
    if err != nil {
        return err
    }
    
    collection, err := service.GetCollection(sss.collection)
    if err != nil {
        return err
    }
    
    return collection.CreateItem(key, cred.Value, true)
}

// Windows Credential Manager implementation
type WindowsCredStore struct {
    target string
}

func (wcs *WindowsCredStore) Store(key string, cred *Credential) error {
    return wincred.Write(&wincred.Credential{
        TargetName: fmt.Sprintf("%s/%s", wcs.target, key),
        Type:       wincred.TypeGeneric,
        Persist:    wincred.PersistLocalMachine,
        CredentialBlob: cred.Value,
    })
}
```

### Encryption Layer

```go
type Encryptor struct {
    masterKey []byte
    algorithm EncryptionAlgorithm
}

type EncryptionAlgorithm interface {
    Encrypt(plaintext, key []byte) ([]byte, error)
    Decrypt(ciphertext, key []byte) ([]byte, error)
    GenerateKey() ([]byte, error)
}

// AES-256-GCM implementation
type AES256GCM struct{}

func (a *AES256GCM) Encrypt(plaintext, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    
    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, err
    }
    
    ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
    return ciphertext, nil
}

func (a *AES256GCM) Decrypt(ciphertext, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    
    nonceSize := gcm.NonceSize()
    if len(ciphertext) < nonceSize {
        return nil, errors.New("ciphertext too short")
    }
    
    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    return gcm.Open(nil, nonce, ciphertext, nil)
}
```

### Credential Access Control

```go
type CredentialAccessControl struct {
    policies map[string]*AccessPolicy
    auditor  *AuditLogger
}

type AccessPolicy struct {
    AllowedExtensions []string
    AllowedProfiles   []string
    RequireMFA        bool
    TimeRestrictions  *TimeWindow
    IPRestrictions    []net.IPNet
}

func (cac *CredentialAccessControl) CheckAccess(
    credID string,
    requester AccessRequester,
) error {
    policy, exists := cac.policies[credID]
    if !exists {
        policy = cac.defaultPolicy()
    }
    
    // Log access attempt
    cac.auditor.LogAccess(AccessEvent{
        CredentialID: credID,
        Requester:    requester,
        Timestamp:    time.Now(),
    })
    
    // Check extension allowlist
    if !contains(policy.AllowedExtensions, requester.ExtensionID) {
        return ErrAccessDenied
    }
    
    // Check time restrictions
    if policy.TimeRestrictions != nil {
        if !policy.TimeRestrictions.IsActive() {
            return ErrAccessOutsideTimeWindow
        }
    }
    
    // Check MFA requirement
    if policy.RequireMFA && !requester.MFAVerified {
        return ErrMFARequired
    }
    
    return nil
}
```

### Environment Variable Injection

```go
type SecureEnvironment struct {
    credentials *CredentialManager
    resolver    *VariableResolver
    sanitizer   *EnvSanitizer
}

func (se *SecureEnvironment) PrepareEnvironment(
    baseEnv map[string]string,
    requiredCreds []string,
) (map[string]string, error) {
    env := make(map[string]string)
    
    // Copy allowed base environment
    for k, v := range baseEnv {
        if se.sanitizer.IsAllowed(k) {
            env[k] = v
        }
    }
    
    // Inject credentials
    for _, credID := range requiredCreds {
        cred, err := se.credentials.Retrieve(credID)
        if err != nil {
            return nil, fmt.Errorf("retrieving credential %s: %w", credID, err)
        }
        
        // Decrypt and inject
        value, err := se.credentials.Decrypt(cred)
        if err != nil {
            return nil, fmt.Errorf("decrypting credential %s: %w", credID, err)
        }
        
        envKey := se.resolver.ResolveEnvKey(credID)
        env[envKey] = string(value)
    }
    
    return env, nil
}

type EnvSanitizer struct {
    blocklist []string
    allowlist []string
}

func (es *EnvSanitizer) IsAllowed(key string) bool {
    // Never allow sensitive system variables
    sensitiveVars := []string{
        "LD_PRELOAD",
        "DYLD_INSERT_LIBRARIES",
        "PATH", // Controlled separately
    }
    
    for _, blocked := range sensitiveVars {
        if key == blocked {
            return false
        }
    }
    
    return true
}
```

## Authentication and Authorization

### Multi-Factor Authentication

```go
type MFAProvider interface {
    GenerateChallenge(userID string) (*MFAChallenge, error)
    VerifyResponse(challenge *MFAChallenge, response string) error
}

type TOTPProvider struct {
    secrets map[string]string
}

func (tp *TOTPProvider) VerifyResponse(
    challenge *MFAChallenge,
    response string,
) error {
    secret := tp.secrets[challenge.UserID]
    
    totp := gotp.NewTOTP(secret)
    if !totp.Verify(response, time.Now().Unix()) {
        return ErrInvalidTOTP
    }
    
    return nil
}
```

### Session Management

```go
type SessionManager struct {
    store    SessionStore
    timeout  time.Duration
    verifier TokenVerifier
}

type Session struct {
    ID          string
    UserID      string
    ProfileID   string
    Permissions []string
    ExpiresAt   time.Time
    MFAVerified bool
}

func (sm *SessionManager) CreateSession(
    userID string,
    profileID string,
) (*Session, error) {
    session := &Session{
        ID:        generateSecureToken(),
        UserID:    userID,
        ProfileID: profileID,
        ExpiresAt: time.Now().Add(sm.timeout),
    }
    
    if err := sm.store.Save(session); err != nil {
        return nil, err
    }
    
    return session, nil
}
```

## Security Monitoring

### Audit Logging

```go
type AuditLogger struct {
    writer    io.Writer
    encryptor Encryptor
    signer    Signer
}

type AuditEvent struct {
    ID          string
    Type        EventType
    Actor       Actor
    Resource    Resource
    Action      Action
    Result      Result
    Timestamp   time.Time
    Metadata    map[string]interface{}
}

func (al *AuditLogger) LogEvent(event AuditEvent) error {
    // Serialize event
    data, err := json.Marshal(event)
    if err != nil {
        return err
    }
    
    // Sign for integrity
    signature := al.signer.Sign(data)
    
    // Encrypt sensitive data
    encrypted, err := al.encryptor.Encrypt(data, nil)
    if err != nil {
        return err
    }
    
    // Write to audit log
    entry := AuditLogEntry{
        Data:      encrypted,
        Signature: signature,
        Timestamp: event.Timestamp,
    }
    
    return al.write(entry)
}
```

### Anomaly Detection

```go
type AnomalyDetector struct {
    patterns    []AnomalyPattern
    threshold   float64
    alerter     Alerter
}

type AnomalyPattern struct {
    Name        string
    Description string
    Detector    func(events []AuditEvent) float64
}

func (ad *AnomalyDetector) Analyze(events []AuditEvent) {
    for _, pattern := range ad.patterns {
        score := pattern.Detector(events)
        if score > ad.threshold {
            ad.alerter.Alert(SecurityAlert{
                Pattern:     pattern.Name,
                Score:       score,
                Description: pattern.Description,
                Events:      events,
            })
        }
    }
}

// Example patterns
var anomalyPatterns = []AnomalyPattern{
    {
        Name: "credential_brute_force",
        Description: "Multiple failed credential access attempts",
        Detector: func(events []AuditEvent) float64 {
            // Count failed credential accesses in time window
        },
    },
    {
        Name: "privilege_escalation",
        Description: "Unusual permission requests",
        Detector: func(events []AuditEvent) float64 {
            // Detect unusual permission patterns
        },
    },
}
```

## Secure Update Mechanism

### Update Verification

```go
type UpdateManager struct {
    verifier    SignatureVerifier
    downloader  SecureDownloader
    installer   Installer
}

func (um *UpdateManager) ApplyUpdate(update UpdatePackage) error {
    // 1. Verify signature
    if err := um.verifier.Verify(update); err != nil {
        return fmt.Errorf("invalid update signature: %w", err)
    }
    
    // 2. Verify checksum
    if err := um.verifyChecksum(update); err != nil {
        return fmt.Errorf("checksum mismatch: %w", err)
    }
    
    // 3. Scan for malware (optional)
    if err := um.scanForMalware(update); err != nil {
        return fmt.Errorf("malware detected: %w", err)
    }
    
    // 4. Create backup
    if err := um.createBackup(); err != nil {
        return fmt.Errorf("backup failed: %w", err)
    }
    
    // 5. Apply update
    return um.installer.Install(update)
}
```

### Certificate Pinning

```go
type CertificatePinner struct {
    pins map[string][]string // hostname -> certificate hashes
}

func (cp *CertificatePinner) VerifyConnection(
    conn *tls.Conn,
) error {
    state := conn.ConnectionState()
    hostname := conn.RemoteAddr().String()
    
    expectedPins, exists := cp.pins[hostname]
    if !exists {
        return nil // No pinning for this host
    }
    
    for _, cert := range state.PeerCertificates {
        hash := sha256.Sum256(cert.Raw)
        pin := base64.StdEncoding.EncodeToString(hash[:])
        
        for _, expectedPin := range expectedPins {
            if pin == expectedPin {
                return nil // Valid pin found
            }
        }
    }
    
    return ErrCertificatePinMismatch
}
```

## Security Best Practices

### Input Validation

```go
type InputValidator struct {
    rules map[string]ValidationRule
}

type ValidationRule interface {
    Validate(input interface{}) error
}

// Common validation rules
var commonRules = map[string]ValidationRule{
    "extension_id": RegexRule{
        Pattern: "^[a-z0-9-]+$",
        MaxLen:  50,
    },
    "file_path": PathRule{
        AllowRelative: false,
        MaxDepth:      10,
        BlockPatterns: []string{"../", "..\\"},
    },
    "url": URLRule{
        AllowedSchemes: []string{"https"},
        AllowedHosts:   []string{"*.gemini.dev", "github.com"},
    },
}
```

### Rate Limiting

```go
type RateLimiter struct {
    limits map[string]*Limit
    store  RateLimitStore
}

type Limit struct {
    MaxRequests int
    Window      time.Duration
    BurstSize   int
}

func (rl *RateLimiter) CheckLimit(
    key string,
    operation string,
) error {
    limit, exists := rl.limits[operation]
    if !exists {
        return nil // No limit defined
    }
    
    count, err := rl.store.Increment(key, limit.Window)
    if err != nil {
        return err
    }
    
    if count > limit.MaxRequests {
        return ErrRateLimitExceeded
    }
    
    return nil
}
```

## Security Checklist

### For Extension Developers
- [ ] Declare all required permissions in settings.json
- [ ] Use minimal permissions necessary
- [ ] Never log sensitive information
- [ ] Validate all inputs
- [ ] Use secure communication (HTTPS)
- [ ] Handle errors gracefully without exposing internals

### For Users
- [ ] Review extension permissions before installation
- [ ] Use strong, unique credentials
- [ ] Enable MFA where available
- [ ] Regularly update extensions
- [ ] Monitor audit logs for suspicious activity
- [ ] Report security issues promptly

### For Maintainers
- [ ] Regular security audits
- [ ] Dependency vulnerability scanning
- [ ] Penetration testing
- [ ] Security training for contributors
- [ ] Incident response plan
- [ ] Regular security updates

This comprehensive security specification ensures that the Gemini CLI Manager maintains the highest security standards while remaining user-friendly and performant.