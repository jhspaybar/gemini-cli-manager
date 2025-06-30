#!/usr/bin/env python3
"""
Fix all text rendering to use explicit theme colors
"""

import re
import os

def fix_span_raw(content):
    """Fix Span::raw() to include text color"""
    # Match Span::raw(some_text) where some_text is not just whitespace
    pattern = r'Span::raw\(([^)]+)\)'
    
    def replace_span(match):
        text = match.group(1)
        # Skip if it's just whitespace or empty quotes
        if text.strip() in ['""', '""', '" "', '"  "', '"    "', '"\t"']:
            return match.group(0)
        # Skip if it already has a style
        if 'Span::styled' in match.group(0):
            return match.group(0)
        return f'Span::styled({text}, Style::default().fg(theme::text_primary()))'
    
    return re.sub(pattern, replace_span, content)

def fix_line_from_string(content):
    """Fix Line::from("string") to include text color"""
    # Match Line::from("some text") but not Line::from("")
    pattern = r'Line::from\("([^"]+)"\)'
    
    def replace_line(match):
        text = match.group(1)
        if text == "":  # Keep empty lines as-is
            return match.group(0)
        return f'Line::from(Span::styled("{text}", Style::default().fg(theme::text_primary())))'
    
    return re.sub(pattern, replace_line, content)

def fix_format_strings_in_lines(content):
    """Fix format! strings in Line::from"""
    # Match Line::from(format!(...))
    pattern = r'Line::from\((format!\([^)]+\))\)'
    
    def replace_format(match):
        format_expr = match.group(1)
        return f'Line::from(Span::styled({format_expr}, Style::default().fg(theme::text_primary())))'
    
    return re.sub(pattern, replace_format, content)

def process_file(filepath):
    """Process a single file"""
    print(f"Processing {filepath}...")
    
    with open(filepath, 'r') as f:
        content = f.read()
    
    original = content
    
    # Apply fixes
    content = fix_span_raw(content)
    content = fix_line_from_string(content)
    content = fix_format_strings_in_lines(content)
    
    if content != original:
        with open(filepath, 'w') as f:
            f.write(content)
        print(f"  Fixed text color issues")
    else:
        print(f"  No changes needed")

def main():
    # Focus on the problem files identified
    files = [
        'src/components/extension_detail.rs',
        'src/components/extension_form.rs',
        'src/components/profile_detail.rs',
        'src/components/profile_form.rs',
    ]
    
    for filepath in files:
        if os.path.exists(filepath):
            process_file(filepath)

if __name__ == '__main__':
    main()