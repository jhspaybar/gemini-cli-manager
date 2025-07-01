# GitHub Operations Guide

## Authentication

Always use environment variables for credentials:
```bash
export GITHUB_TOKEN="your-token-here"
```

Never hardcode tokens in code or commit them to repositories.

## Core Requirements

1. **Check rate limits** - Always monitor API rate limits before making requests
2. **Use conditional requests** - Include ETags to avoid unnecessary API calls
3. **Handle errors gracefully** - API calls can fail for many reasons
4. **Prefer GraphQL over REST** - More efficient for complex queries
5. **Use GitHub CLI when available** - `gh` command is often simpler than raw API

## Common Operations

### Using GitHub CLI
```bash
# List issues
gh issue list --repo owner/repo

# Create a pull request
gh pr create --title "Fix bug" --body "Description"

# View PR status
gh pr status

# Clone a repository
gh repo clone owner/repo
```

### API Best Practices
- Always include User-Agent header
- Use pagination for list endpoints (default: 30 items)
- Cache responses when appropriate
- Implement exponential backoff for retries

## Required Checks

Before any destructive operations:
- Verify repository ownership
- Check user permissions
- Confirm branch protection rules
- Validate webhook signatures

## Error Handling

Common errors to handle:
- 401: Invalid authentication
- 403: Rate limit exceeded or insufficient permissions
- 404: Resource not found
- 422: Validation failed

## What NOT to Do
- Never store tokens in code
- Never ignore rate limits
- Never assume API calls succeed
- Never bypass branch protection
- Never make bulk changes without confirmation