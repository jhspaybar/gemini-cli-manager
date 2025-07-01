# GITHUB.md - GitHub API Integration Guide

This guide provides best practices for integrating with GitHub's API, handling authentication, and following GitHub's conventions.

## Authentication

### Personal Access Tokens (Recommended)

```bash
# Set your GitHub token as an environment variable
export GITHUB_TOKEN="ghp_xxxxxxxxxxxxxxxxxxxx"

# Use in API requests
curl -H "Authorization: token $GITHUB_TOKEN" \
     https://api.github.com/user
```

### GitHub CLI Authentication

```bash
# Authenticate with GitHub CLI
gh auth login

# Use gh API commands
gh api user
gh api repos/:owner/:repo/issues
```

## API Best Practices

### Rate Limiting

```python
import requests
import time

class GitHubClient:
    def __init__(self, token):
        self.token = token
        self.session = requests.Session()
        self.session.headers.update({
            'Authorization': f'token {token}',
            'Accept': 'application/vnd.github.v3+json'
        })
    
    def request_with_retry(self, url, method='GET', **kwargs):
        """Make request with rate limit handling"""
        response = self.session.request(method, url, **kwargs)
        
        # Check rate limit
        remaining = int(response.headers.get('X-RateLimit-Remaining', 0))
        if remaining == 0:
            reset_time = int(response.headers.get('X-RateLimit-Reset', 0))
            sleep_time = reset_time - time.time() + 1
            if sleep_time > 0:
                print(f"Rate limit hit. Sleeping for {sleep_time} seconds")
                time.sleep(sleep_time)
                return self.request_with_retry(url, method, **kwargs)
        
        response.raise_for_status()
        return response
```

### Pagination

```python
def get_all_items(self, url):
    """Get all items handling pagination"""
    items = []
    
    while url:
        response = self.request_with_retry(url)
        items.extend(response.json())
        
        # Check for next page
        link_header = response.headers.get('Link', '')
        next_url = None
        
        for link in link_header.split(','):
            if 'rel="next"' in link:
                next_url = link.split('<')[1].split('>')[0]
                break
        
        url = next_url
    
    return items
```

## Common Operations

### Repository Management

```python
# Create a repository
def create_repo(self, name, description="", private=False):
    """Create a new repository"""
    data = {
        'name': name,
        'description': description,
        'private': private,
        'auto_init': True,  # Initialize with README
        'gitignore_template': 'Python',  # Add .gitignore
        'license_template': 'mit'  # Add license
    }
    
    response = self.request_with_retry(
        'https://api.github.com/user/repos',
        method='POST',
        json=data
    )
    return response.json()

# Clone repository
def clone_repo(repo_url, local_path):
    """Clone repository with progress"""
    import subprocess
    
    cmd = ['git', 'clone', '--progress', repo_url, local_path]
    process = subprocess.Popen(
        cmd,
        stdout=subprocess.PIPE,
        stderr=subprocess.STDOUT,
        universal_newlines=True
    )
    
    for line in process.stdout:
        print(line.strip())
    
    return process.wait() == 0
```

### Issue Management

```python
# Create an issue
def create_issue(self, owner, repo, title, body, labels=None, assignees=None):
    """Create a new issue"""
    data = {
        'title': title,
        'body': body,
        'labels': labels or [],
        'assignees': assignees or []
    }
    
    url = f'https://api.github.com/repos/{owner}/{repo}/issues'
    response = self.request_with_retry(url, method='POST', json=data)
    return response.json()

# Search issues
def search_issues(self, query):
    """Search issues with advanced query"""
    # Example queries:
    # - "is:open is:issue assignee:username"
    # - "is:pr is:open review-requested:username"
    # - "is:issue label:bug created:>2024-01-01"
    
    url = 'https://api.github.com/search/issues'
    params = {'q': query, 'per_page': 100}
    
    response = self.request_with_retry(url, params=params)
    return response.json()['items']
```

### Pull Request Workflow

```python
# Create a pull request
def create_pull_request(self, owner, repo, title, head, base, body=""):
    """Create a pull request"""
    data = {
        'title': title,
        'head': head,  # Branch name or "username:branch"
        'base': base,  # Target branch (usually 'main')
        'body': body,
        'draft': False
    }
    
    url = f'https://api.github.com/repos/{owner}/{repo}/pulls'
    response = self.request_with_retry(url, method='POST', json=data)
    return response.json()

# Review a pull request
def review_pull_request(self, owner, repo, pr_number, event, body=""):
    """Submit a pull request review"""
    # event: "APPROVE", "REQUEST_CHANGES", "COMMENT"
    data = {
        'body': body,
        'event': event
    }
    
    url = f'https://api.github.com/repos/{owner}/{repo}/pulls/{pr_number}/reviews'
    response = self.request_with_retry(url, method='POST', json=data)
    return response.json()
```

### GitHub Actions

```yaml
# .github/workflows/ci.yml
name: CI

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Python
      uses: actions/setup-python@v4
      with:
        python-version: '3.11'
    
    - name: Install dependencies
      run: |
        python -m pip install --upgrade pip
        pip install -r requirements.txt
    
    - name: Run tests
      run: |
        pytest tests/ -v --cov=src
    
    - name: Upload coverage
      uses: codecov/codecov-action@v3
```

### Webhooks

```python
# Webhook handler
from flask import Flask, request
import hmac
import hashlib

app = Flask(__name__)
WEBHOOK_SECRET = 'your-webhook-secret'

def verify_webhook_signature(payload, signature):
    """Verify GitHub webhook signature"""
    expected = hmac.new(
        WEBHOOK_SECRET.encode(),
        payload,
        hashlib.sha256
    ).hexdigest()
    
    return hmac.compare_digest(
        f"sha256={expected}",
        signature
    )

@app.route('/webhook', methods=['POST'])
def handle_webhook():
    signature = request.headers.get('X-Hub-Signature-256')
    if not verify_webhook_signature(request.data, signature):
        return 'Invalid signature', 401
    
    event = request.headers.get('X-GitHub-Event')
    payload = request.json
    
    if event == 'push':
        handle_push(payload)
    elif event == 'pull_request':
        handle_pull_request(payload)
    elif event == 'issues':
        handle_issue(payload)
    
    return 'OK', 200
```

## GraphQL API

```python
# GraphQL queries
def graphql_query(self, query, variables=None):
    """Execute GraphQL query"""
    url = 'https://api.github.com/graphql'
    data = {'query': query}
    if variables:
        data['variables'] = variables
    
    response = self.request_with_retry(url, method='POST', json=data)
    return response.json()

# Example: Get repository info
query = '''
query($owner: String!, $name: String!) {
  repository(owner: $owner, name: $name) {
    name
    description
    stargazerCount
    forkCount
    issues(states: OPEN) {
      totalCount
    }
    pullRequests(states: OPEN) {
      totalCount
    }
    releases(last: 1) {
      nodes {
        tagName
        publishedAt
      }
    }
  }
}
'''

variables = {'owner': 'octocat', 'name': 'hello-world'}
result = graphql_query(query, variables)
```

## Git Operations

### Commit Best Practices

```bash
# Good commit messages
git commit -m "feat: add user authentication module"
git commit -m "fix: resolve memory leak in data processor"
git commit -m "docs: update API documentation"
git commit -m "refactor: simplify error handling logic"

# Conventional commits
# Format: <type>(<scope>): <subject>
# Types: feat, fix, docs, style, refactor, test, chore
```

### Branch Management

```bash
# Create feature branch
git checkout -b feature/add-oauth-support

# Create from specific commit
git checkout -b fix/memory-leak abc123

# Push new branch
git push -u origin feature/add-oauth-support

# Delete local and remote branch
git branch -d feature/old-feature
git push origin --delete feature/old-feature
```

### Working with Forks

```bash
# Add upstream remote
git remote add upstream https://github.com/original/repo.git

# Sync fork with upstream
git fetch upstream
git checkout main
git merge upstream/main
git push origin main

# Create PR from fork
gh pr create --base upstream:main --head yourusername:feature-branch
```

## Security Best Practices

### Secret Scanning

```python
# Never commit secrets
# Use environment variables
import os

GITHUB_TOKEN = os.environ.get('GITHUB_TOKEN')
if not GITHUB_TOKEN:
    raise ValueError("GITHUB_TOKEN environment variable not set")

# Use .gitignore
# .env
# *.key
# *.pem
# config/secrets.yml
```

### Dependabot Configuration

```yaml
# .github/dependabot.yml
version: 2
updates:
  - package-ecosystem: "pip"
    directory: "/"
    schedule:
      interval: "weekly"
    open-pull-requests-limit: 10
    
  - package-ecosystem: "github-actions"
    directory: "/"
    schedule:
      interval: "weekly"
```

### Code Scanning

```yaml
# .github/workflows/codeql.yml
name: "CodeQL"

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]
  schedule:
    - cron: '30 5 * * 1'

jobs:
  analyze:
    name: Analyze
    runs-on: ubuntu-latest
    
    strategy:
      matrix:
        language: [ 'python' ]
    
    steps:
    - name: Checkout repository
      uses: actions/checkout@v3
    
    - name: Initialize CodeQL
      uses: github/codeql-action/init@v2
      with:
        languages: ${{ matrix.language }}
    
    - name: Autobuild
      uses: github/codeql-action/autobuild@v2
    
    - name: Perform CodeQL Analysis
      uses: github/codeql-action/analyze@v2
```

## GitHub CLI Examples

```bash
# Repository operations
gh repo create my-project --public
gh repo clone owner/repo
gh repo view owner/repo --web

# Issue operations
gh issue create --title "Bug report" --body "Description"
gh issue list --label "bug"
gh issue close 123

# PR operations
gh pr create --fill
gh pr list --state open
gh pr review 456 --approve
gh pr merge 456 --squash

# Workflow operations
gh workflow list
gh workflow run ci.yml
gh run list --workflow=ci.yml

# Release operations
gh release create v1.0.0 --notes "First release"
gh release download v1.0.0
```

## Error Handling

```python
class GitHubError(Exception):
    """Base exception for GitHub operations"""
    pass

class RateLimitError(GitHubError):
    """Rate limit exceeded"""
    pass

class NotFoundError(GitHubError):
    """Resource not found"""
    pass

def handle_github_error(response):
    """Handle GitHub API errors"""
    if response.status_code == 404:
        raise NotFoundError("Resource not found")
    elif response.status_code == 403:
        if 'rate limit' in response.text.lower():
            raise RateLimitError("Rate limit exceeded")
        raise GitHubError("Forbidden")
    elif response.status_code >= 400:
        raise GitHubError(f"API error: {response.text}")
```

## Best Practices Summary

1. **Always use authentication** for better rate limits
2. **Handle rate limiting** gracefully with retries
3. **Use pagination** for large result sets
4. **Verify webhooks** with signatures
5. **Never commit secrets** - use environment variables
6. **Follow conventional commits** for clear history
7. **Use GitHub Actions** for CI/CD
8. **Enable security features** like Dependabot and CodeQL
9. **Cache API responses** when appropriate
10. **Use GraphQL** for complex queries to reduce API calls

Remember: The GitHub API is powerful but has limits. Design your integrations to be respectful of rate limits and efficient in their API usage.