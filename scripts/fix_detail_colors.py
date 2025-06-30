#!/usr/bin/env python3
"""
Fix text colors in detail views
"""

import re

# Read extension_detail.rs
with open('src/components/extension_detail.rs', 'r') as f:
    content = f.read()

# List of replacements
replacements = [
    # Fix value spans
    ('Span::raw(desc),', 'Span::styled(desc, Style::default().fg(theme::text_primary())),'),
    ('Span::raw(&extension.id),', 'Span::styled(&extension.id, Style::default().fg(theme::text_primary())),'),
    ('Span::raw(extension.metadata.imported_at', 'Span::styled(extension.metadata.imported_at'),
    ('.to_string()),', '.to_string(), Style::default().fg(theme::text_primary())),'),
    ('Span::raw(path),', 'Span::styled(path, Style::default().fg(theme::text_primary())),'),
    
    # Fix labels  
    ('Span::raw("    Type: "),', 'Span::styled("    Type: ", Style::default().fg(theme::text_secondary())),'),
    ('Span::raw("    Args: "),', 'Span::styled("    Args: ", Style::default().fg(theme::text_secondary())),'),
    ('Span::raw("    Env: "),', 'Span::styled("    Env: ", Style::default().fg(theme::text_secondary())),'),
    ('Span::raw("    Trust: "),', 'Span::styled("    Trust: ", Style::default().fg(theme::text_secondary())),'),
    
    # Fix values in server config
    ('Span::raw(url),', 'Span::styled(url, Style::default().fg(theme::text_primary())),'),
    ('Span::raw(cmd),', 'Span::styled(cmd, Style::default().fg(theme::text_primary())),'),
    ('Span::raw(args.join(" ")),', 'Span::styled(args.join(" "), Style::default().fg(theme::text_primary())),'),
    ('Span::raw(value),', 'Span::styled(value, Style::default().fg(theme::text_primary())),'),
    
    # Fix context file content
    ('Span::raw(line),', 'Span::styled(line, Style::default().fg(theme::text_primary())),'),
]

# Apply replacements
for old, new in replacements:
    content = content.replace(old, new)

# Write back
with open('src/components/extension_detail.rs', 'w') as f:
    f.write(content)

print("Fixed extension_detail.rs")

# Now fix profile_detail.rs
with open('src/components/profile_detail.rs', 'r') as f:
    content = f.read()

replacements = [
    # Fix value spans
    ('Span::raw(desc),', 'Span::styled(desc, Style::default().fg(theme::text_primary())),'),
    ('Span::raw(&profile.id),', 'Span::styled(&profile.id, Style::default().fg(theme::text_primary())),'),
    ('Span::raw(dir),', 'Span::styled(dir, Style::default().fg(theme::text_primary())),'),
    ('Span::raw(profile.metadata.created_at', 'Span::styled(profile.metadata.created_at'),
    ('.to_string()),', '.to_string(), Style::default().fg(theme::text_primary())),'),
    
    # Fix list items
    ('Line::from("  No extensions included")', 'Line::from(Span::styled("  No extensions included", Style::default().fg(theme::text_muted())))'),
    
    # Fix env var values
    ('Span::raw(display_value),', 'Span::styled(display_value, Style::default().fg(theme::text_primary())),'),
    
    # Fix summary lines
    ('Line::from(format!("  • {} extensions"', 'Line::from(Span::styled(format!("  • {} extensions"'),
    ('Line::from(format!("  • {} environment variables"', 'Line::from(Span::styled(format!("  • {} environment variables"'),
    ('Line::from(format!("  • {} MCP servers total"', 'Line::from(Span::styled(format!("  • {} MCP servers total"'),
    (', self.extensions.len()))', ', self.extensions.len()), Style::default().fg(theme::text_primary())))'),
    (', profile.environment_variables.len()))', ', profile.environment_variables.len()), Style::default().fg(theme::text_primary())))'),
    (', total_mcp_servers))', ', total_mcp_servers), Style::default().fg(theme::text_primary())))'),
]

# Apply replacements
for old, new in replacements:
    content = content.replace(old, new)

# Write back
with open('src/components/profile_detail.rs', 'w') as f:
    f.write(content)

print("Fixed profile_detail.rs")