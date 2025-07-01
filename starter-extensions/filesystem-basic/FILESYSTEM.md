# FILESYSTEM.md - Safe Filesystem Operations Guide

This guide provides safe patterns and best practices for filesystem operations. Always prioritize security and data integrity when working with files and directories.

## Core Principles

1. **Always validate paths** - Never trust user input
2. **Check permissions** - Ensure proper access rights
3. **Handle errors gracefully** - Filesystem operations can fail
4. **Use atomic operations** - Prevent data corruption
5. **Clean up resources** - Close files and remove temporary files

## Safe Path Handling

### Path Validation

```bash
# Always resolve paths to prevent traversal attacks
realpath "$user_path" || exit 1

# Check if path is within allowed directory
allowed_dir="/home/user/workspace"
resolved_path=$(realpath "$user_path")
if [[ ! "$resolved_path" =~ ^"$allowed_dir" ]]; then
    echo "Error: Path outside allowed directory"
    exit 1
fi
```

### Path Sanitization

```python
import os
from pathlib import Path

def safe_join(base_path, user_path):
    """Safely join paths preventing directory traversal"""
    base = Path(base_path).resolve()
    target = (base / user_path).resolve()
    
    # Ensure target is within base directory
    if not str(target).startswith(str(base)):
        raise ValueError("Path traversal attempt detected")
    
    return target
```

## File Operations

### Reading Files Safely

```python
def read_file_safely(filepath, max_size=10*1024*1024):  # 10MB limit
    """Read file with size limit to prevent memory exhaustion"""
    filepath = Path(filepath).resolve()
    
    # Check file exists and is a regular file
    if not filepath.is_file():
        raise ValueError(f"Not a regular file: {filepath}")
    
    # Check file size
    size = filepath.stat().st_size
    if size > max_size:
        raise ValueError(f"File too large: {size} bytes")
    
    # Read with encoding error handling
    try:
        return filepath.read_text(encoding='utf-8', errors='replace')
    except Exception as e:
        raise IOError(f"Failed to read file: {e}")
```

### Writing Files Atomically

```python
import tempfile
import shutil

def write_file_atomically(filepath, content):
    """Write file atomically to prevent corruption"""
    filepath = Path(filepath).resolve()
    
    # Create temporary file in same directory
    temp_fd, temp_path = tempfile.mkstemp(
        dir=filepath.parent,
        prefix=f".{filepath.name}.",
        suffix=".tmp"
    )
    
    try:
        # Write to temporary file
        with os.fdopen(temp_fd, 'w', encoding='utf-8') as f:
            f.write(content)
            f.flush()
            os.fsync(f.fileno())  # Force write to disk
        
        # Set proper permissions
        os.chmod(temp_path, 0o644)
        
        # Atomic rename
        shutil.move(temp_path, filepath)
    except Exception:
        # Clean up on error
        try:
            os.unlink(temp_path)
        except:
            pass
        raise
```

### Creating Directories Safely

```python
def create_directory_safely(dirpath, mode=0o755):
    """Create directory with proper permissions"""
    dirpath = Path(dirpath).resolve()
    
    try:
        # Create with parents, exist_ok for idempotency
        dirpath.mkdir(parents=True, exist_ok=True, mode=mode)
        
        # Verify permissions
        actual_mode = dirpath.stat().st_mode & 0o777
        if actual_mode != mode:
            dirpath.chmod(mode)
    except Exception as e:
        raise IOError(f"Failed to create directory: {e}")
```

## Working with Temporary Files

### Secure Temporary Files

```python
import tempfile
import contextlib

@contextlib.contextmanager
def secure_temp_file(suffix=None, prefix=None):
    """Create secure temporary file that's automatically cleaned up"""
    fd, path = tempfile.mkstemp(suffix=suffix, prefix=prefix)
    try:
        # Close file descriptor but keep the file
        os.close(fd)
        yield path
    finally:
        # Always clean up
        try:
            os.unlink(path)
        except OSError:
            pass

# Usage
with secure_temp_file(suffix='.json') as temp_path:
    Path(temp_path).write_text(json.dumps(data))
    process_file(temp_path)
    # File automatically deleted when done
```

### Temporary Directories

```python
import tempfile
import shutil

@contextlib.contextmanager
def secure_temp_dir():
    """Create secure temporary directory"""
    temp_dir = tempfile.mkdtemp(prefix='workspace_')
    try:
        # Set restrictive permissions
        os.chmod(temp_dir, 0o700)
        yield temp_dir
    finally:
        # Recursive removal
        shutil.rmtree(temp_dir, ignore_errors=True)
```

## File Permissions and Security

### Checking Permissions

```python
def check_file_permissions(filepath):
    """Check if file has secure permissions"""
    filepath = Path(filepath).resolve()
    stat = filepath.stat()
    mode = stat.st_mode
    
    # Check for world-writable
    if mode & 0o002:
        raise PermissionError("File is world-writable")
    
    # Check for group-writable (optional)
    if mode & 0o020:
        print("Warning: File is group-writable")
    
    # Check ownership
    if stat.st_uid != os.getuid():
        raise PermissionError("File not owned by current user")
```

### Setting Secure Permissions

```python
def set_secure_permissions(filepath, is_executable=False):
    """Set secure file permissions"""
    filepath = Path(filepath).resolve()
    
    if is_executable:
        # Owner: rwx, Group: r-x, Other: r-x
        mode = 0o755
    else:
        # Owner: rw-, Group: r--, Other: r--
        mode = 0o644
    
    # For sensitive files, remove group/other access
    if 'private' in str(filepath) or 'secret' in str(filepath):
        mode = 0o600  # Owner only
    
    filepath.chmod(mode)
```

## Directory Traversal

### Safe Directory Walking

```python
def walk_directory_safely(root_path, max_depth=10):
    """Safely walk directory tree with depth limit"""
    root_path = Path(root_path).resolve()
    
    for dirpath, dirnames, filenames in os.walk(root_path):
        # Calculate depth
        depth = len(Path(dirpath).relative_to(root_path).parts)
        
        if depth >= max_depth:
            dirnames[:] = []  # Don't recurse deeper
            continue
        
        # Skip hidden directories
        dirnames[:] = [d for d in dirnames if not d.startswith('.')]
        
        for filename in filenames:
            # Skip hidden files
            if filename.startswith('.'):
                continue
            
            filepath = Path(dirpath) / filename
            yield filepath
```

### Finding Files Safely

```python
import fnmatch

def find_files_safely(root_path, pattern, max_results=1000):
    """Find files matching pattern with limits"""
    root_path = Path(root_path).resolve()
    results = []
    
    for filepath in walk_directory_safely(root_path):
        if fnmatch.fnmatch(filepath.name, pattern):
            results.append(filepath)
            
            if len(results) >= max_results:
                print(f"Warning: Result limit ({max_results}) reached")
                break
    
    return results
```

## File Locking

### Exclusive File Access

```python
import fcntl
import errno

def exclusive_file_access(filepath, operation):
    """Perform operation with exclusive file access"""
    filepath = Path(filepath).resolve()
    
    with open(filepath, 'r+b') as f:
        try:
            # Try to acquire exclusive lock (non-blocking)
            fcntl.flock(f.fileno(), fcntl.LOCK_EX | fcntl.LOCK_NB)
            
            try:
                # Perform operation
                return operation(f)
            finally:
                # Always release lock
                fcntl.flock(f.fileno(), fcntl.LOCK_UN)
                
        except IOError as e:
            if e.errno in (errno.EACCES, errno.EAGAIN):
                raise IOError("File is locked by another process")
            raise
```

## Monitoring File Changes

### Safe File Watching

```python
import time
from pathlib import Path

def watch_file_safely(filepath, callback, interval=1.0):
    """Watch file for changes safely"""
    filepath = Path(filepath).resolve()
    
    if not filepath.exists():
        raise FileNotFoundError(f"File not found: {filepath}")
    
    last_mtime = filepath.stat().st_mtime
    last_size = filepath.stat().st_size
    
    while True:
        try:
            current_stat = filepath.stat()
            current_mtime = current_stat.st_mtime
            current_size = current_stat.st_size
            
            # Check for changes
            if current_mtime != last_mtime or current_size != last_size:
                callback(filepath)
                last_mtime = current_mtime
                last_size = current_size
            
        except FileNotFoundError:
            # File was deleted
            callback(None)
            break
        except Exception as e:
            print(f"Error watching file: {e}")
        
        time.sleep(interval)
```

## Cleanup and Resource Management

### Ensuring Cleanup

```python
class FileManager:
    """Manage file resources with guaranteed cleanup"""
    
    def __init__(self):
        self.temp_files = []
        self.open_files = []
    
    def create_temp_file(self, **kwargs):
        """Create temp file tracked for cleanup"""
        fd, path = tempfile.mkstemp(**kwargs)
        self.temp_files.append(path)
        os.close(fd)
        return path
    
    def open_file(self, filepath, mode='r'):
        """Open file tracked for cleanup"""
        f = open(filepath, mode)
        self.open_files.append(f)
        return f
    
    def cleanup(self):
        """Clean up all resources"""
        # Close all open files
        for f in self.open_files:
            try:
                f.close()
            except:
                pass
        
        # Remove all temp files
        for path in self.temp_files:
            try:
                os.unlink(path)
            except:
                pass
    
    def __enter__(self):
        return self
    
    def __exit__(self, exc_type, exc_val, exc_tb):
        self.cleanup()
```

## Best Practices Summary

1. **Always validate and sanitize paths**
   - Use `realpath` or `Path.resolve()`
   - Check paths stay within allowed directories

2. **Handle errors gracefully**
   - Expect filesystem operations to fail
   - Provide meaningful error messages

3. **Use atomic operations**
   - Write to temp file then rename
   - Use file locking when needed

4. **Set appropriate permissions**
   - Avoid world-writable files
   - Use restrictive permissions for sensitive data

5. **Clean up resources**
   - Close files when done
   - Remove temporary files
   - Use context managers

6. **Implement limits**
   - File size limits
   - Directory depth limits  
   - Result count limits

7. **Be cautious with user input**
   - Never use user input directly in paths
   - Validate file types and content

8. **Monitor and log operations**
   - Log file access for security
   - Monitor for unusual patterns

Remember: Filesystem operations are privileged operations. Always assume they can fail and always validate inputs to prevent security vulnerabilities.