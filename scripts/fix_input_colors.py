#!/usr/bin/env python3
"""
Fix text input colors in form views
"""

import re

def fix_extension_form():
    with open('src/components/extension_form.rs', 'r') as f:
        content = f.read()
    
    # Fix context file name input - line 482 needs fg color
    content = re.sub(
        r'(\s+else \{\s*\n\s*Style::default\(\)\s*\n\s*\});',
        r' else {\n                Style::default().fg(theme::text_primary())\n            });',
        content,
        count=1
    )
    
    # Fix context content text - add style after line 499
    content = re.sub(
        r'(let context_content_text = Paragraph::new\(self\.context_content_input\.value\(\)\)\s*\n\s*\.wrap\(Wrap \{ trim: false \}\)\s*\n\s*\.scroll\(\(self\.context_scroll_offset, 0\)\));',
        r'let context_content_text = Paragraph::new(self.context_content_input.value())\n            .style(Style::default().fg(theme::text_primary()))\n            .wrap(Wrap { trim: false })\n            .scroll((self.context_scroll_offset, 0));',
        content
    )
    
    with open('src/components/extension_form.rs', 'w') as f:
        f.write(content)
    
    print("Fixed extension_form.rs")

def fix_profile_form():
    with open('src/components/profile_form.rs', 'r') as f:
        content = f.read()
    
    # Fix all Paragraph::new instances for input fields
    replacements = [
        # Name input
        (r'(let name_text = Paragraph::new\(self\.name_input\.value\(\)\));',
         r'let name_text = Paragraph::new(self.name_input.value())\n            .style(Style::default().fg(theme::text_primary()));'),
        
        # Description input
        (r'(let desc_text = Paragraph::new\(self\.description_input\.value\(\)\));',
         r'let desc_text = Paragraph::new(self.description_input.value())\n            .style(Style::default().fg(theme::text_primary()));'),
        
        # Working directory input
        (r'(let dir_text = Paragraph::new\(self\.working_directory_input\.value\(\)\));',
         r'let dir_text = Paragraph::new(self.working_directory_input.value())\n            .style(Style::default().fg(theme::text_primary()));'),
        
        # Tags input
        (r'(let tags_text = Paragraph::new\(self\.tags_input\.value\(\)\));',
         r'let tags_text = Paragraph::new(self.tags_input.value())\n            .style(Style::default().fg(theme::text_primary()));'),
    ]
    
    for pattern, replacement in replacements:
        content = re.sub(pattern, replacement, content)
    
    with open('src/components/profile_form.rs', 'w') as f:
        f.write(content)
    
    print("Fixed profile_form.rs")

if __name__ == '__main__':
    fix_extension_form()
    fix_profile_form()