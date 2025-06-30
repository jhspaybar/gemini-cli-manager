#!/usr/bin/env python3
"""
Script to update hardcoded colors to use theme functions
"""

import os
import re

# Mapping of hardcoded colors to theme functions
COLOR_MAPPINGS = {
    'Color::Yellow': 'theme::highlight()',
    'Color::White': 'theme::text_primary()',
    'Color::Green': 'theme::success()',
    'Color::DarkGray': 'theme::text_muted()',
    'Color::Gray': 'theme::text_secondary()',
    'Color::Cyan': 'theme::info()',
    'Color::Blue': 'theme::primary()',
    'Color::Red': 'theme::error()',
    'Color::Magenta': 'theme::accent()',
}

# Files to update
FILES = [
    'src/components/profile_list.rs',
    'src/components/profile_detail.rs',
    'src/components/extension_form.rs',
    'src/components/extension_detail.rs', 
    'src/components/profile_form.rs',
    'src/components/confirm_dialog.rs',
]

def update_file(filepath):
    print(f"Processing {filepath}...")
    
    with open(filepath, 'r') as f:
        content = f.read()
    
    # Check if theme import is already present
    if 'use crate::{' in content and 'theme' not in content:
        # Add theme to imports
        content = re.sub(
            r'(use crate::\{[^}]+)(}\;)',
            r'\1, theme\2',
            content
        )
    
    # Replace hardcoded colors
    for old_color, theme_func in COLOR_MAPPINGS.items():
        if old_color in content:
            content = content.replace(old_color, theme_func)
            print(f"  Replaced {old_color} with {theme_func}")
    
    with open(filepath, 'w') as f:
        f.write(content)

def main():
    for filepath in FILES:
        if os.path.exists(filepath):
            update_file(filepath)
        else:
            print(f"Warning: {filepath} not found")

if __name__ == '__main__':
    main()