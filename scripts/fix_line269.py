#!/usr/bin/env python3

# Read the file
with open('src/components/profile_detail.rs', 'r') as f:
    lines = f.readlines()

# Fix line 269 (index 268)
if 'Line::from(format!' in lines[268]:
    lines[268] = '        content.push(Line::from(Span::styled(format!("  • {} MCP servers total", total_mcp_servers), Style::default().fg(theme::text_primary()))));\n'

# Write back
with open('src/components/profile_detail.rs', 'w') as f:
    f.writelines(lines)

print("Fixed line 269 in profile_detail.rs")