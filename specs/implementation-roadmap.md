# Implementation Roadmap

This document tracks the remaining features to be implemented based on our specification analysis.

## Overview

The Gemini CLI Manager has a solid foundation with basic profile and extension management, but several critical features remain unimplemented. This roadmap prioritizes features based on user impact and technical dependencies.

## Phase 1: Critical Features (High Priority)

### 1. Extension Detail View
**Status**: Not Implemented  
**Specification**: ui-ux-design.md  
**Current State**: Shows "Extension details view not yet implemented" message

**Requirements**:
- Full detail panel accessible via Enter/Space on extension list
- Display extension metadata (name, version, description, author)
- Show MCP server configurations
- Display current health status
- Show extension logs
- Configuration editing interface
- Start/Stop controls for MCP servers

### 2. MCP Server Process Management
**Status**: Not Implemented  
**Specification**: extension-architecture.md  
**Current State**: MCP configuration is parsed but servers cannot be started

**Requirements**:
- Start/stop MCP servers defined in extensions
- Process lifecycle management
- Health monitoring (ping/healthcheck)
- Automatic restart on failure
- Resource limit enforcement
- Log collection and streaming

### 3. Credential Security
**Status**: Not Implemented  
**Specification**: security-specification.md  
**Current State**: Credentials stored as plain text environment variables

**Requirements**:
- OS keychain integration (macOS Keychain, Windows Credential Manager, Linux Secret Service)
- Encrypted storage for sensitive data
- Secure credential injection into MCP processes
- Access control and audit logging
- Credential rotation support

## Phase 2: Enhanced Usability (Medium Priority)

### 4. Command Palette (Ctrl+K)
**Status**: Not Implemented  
**Specification**: ui-ux-design.md  

**Requirements**:
- Global command search
- Fuzzy matching
- Recent commands history
- Keyboard-only navigation
- Context-aware suggestions

### 5. Profile Templates
**Status**: Not Implemented  
**Specification**: profile-management.md  

**Built-in Templates Needed**:
- `web-frontend`: Node.js, React tools, browser automation
- `data-science`: Python, Jupyter, data processing tools
- `backend-api`: Database tools, API testing, monitoring
- `devops`: Kubernetes, Docker, cloud CLI tools
- `mobile`: React Native, Flutter, device management

### 6. Settings Editor
**Status**: Read-only implementation exists  
**Specification**: ui-ux-design.md  

**Requirements**:
- Edit Gemini CLI path
- Configure default profile
- Theme selection
- Auto-update toggle
- Extension directory configuration
- Keyboard shortcut customization

### 7. Extension Updates
**Status**: Method exists but returns "not implemented"  
**Specification**: extension-architecture.md  

**Requirements**:
- Version checking against remote sources
- Safe update process with rollback
- Update notifications in UI
- Batch updates
- Changelog display

### 8. Profile Inheritance
**Status**: Field exists but not implemented  
**Specification**: profile-management.md  

**Requirements**:
- Inherit settings from parent profiles
- Override specific values
- Multiple inheritance support
- Circular dependency detection
- Inheritance visualization

## Phase 3: Advanced Features (Lower Priority)

### 9. Extension Dependency Resolution
**Status**: Not Implemented  
**Specification**: extension-architecture.md  

**Requirements**:
- Parse dependency declarations
- Build dependency graph
- Check version compatibility
- Auto-install dependencies
- Conflict resolution

### 10. Profile Auto-Detection
**Status**: Field exists but not implemented  
**Specification**: profile-management.md  

**Requirements**:
- Pattern matching on working directory
- Automatic profile switching
- User confirmation option
- Override mechanism
- Pattern priority system

### 11. Import/Export Functionality
**Status**: Not Implemented  
**Specification**: profile-management.md  

**Requirements**:
- Export profiles to shareable format
- Import with validation
- Team sharing via Git
- Backup/restore functionality
- Migration tools

### 12. Extension Development Tools
**Status**: Not Implemented  
**Specification**: extension-architecture.md  

**Requirements**:
- Extension scaffolding
- Validation tools
- Testing framework
- Documentation generator
- Publishing tools

## Technical Debt Items

### Security Hardening
- Input validation improvements
- Command injection prevention
- Path traversal prevention (partially implemented)
- Audit logging system
- Permission model for extensions

### Performance Optimization
- Lazy loading of extensions
- Profile caching
- Startup time optimization
- Memory usage optimization

### Testing & Quality
- Integration tests
- End-to-end tests
- Performance benchmarks
- Security audit
- Accessibility testing

## Implementation Guidelines

### Priority Criteria
1. **User Impact**: How many users will benefit?
2. **Security**: Does it address a security concern?
3. **Dependencies**: What other features depend on this?
4. **Complexity**: How difficult is implementation?
5. **Risk**: What could go wrong?

### Development Process
1. Create detailed design document for each feature
2. Implement with comprehensive tests
3. Security review for sensitive features
4. Performance testing for resource-intensive features
5. User acceptance testing

### Success Metrics
- Feature completion rate
- Test coverage maintained above 80%
- No critical security issues
- Performance benchmarks met
- User satisfaction scores

## Notes

This roadmap is a living document and should be updated as:
- New requirements are discovered
- Priorities shift based on user feedback
- Technical constraints are identified
- Features are completed

Last Updated: {{current_date}}