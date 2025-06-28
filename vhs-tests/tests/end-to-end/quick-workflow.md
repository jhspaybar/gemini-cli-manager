# Quick Workflow Test Specification

## Test: quick-workflow.tape

### Purpose
A streamlined version of the full workflow test, focusing on speed and key interactions without extensive screenshots.

### Test Flow

1. **Setup** - Start app with test state directory
2. **Install** - Quick extension installation
3. **Profile** - Create and configure profile in one flow
4. **Activate** - Set profile as active
5. **Launch** - Start Gemini with profile
6. **Exit** - Clean exit from both Gemini and manager

### Key Differences from Full Test

- Faster typing speed (100ms vs 75ms)
- Fewer screenshots (only GIF output)
- Minimal pauses between actions
- Streamlined navigation
- Same functional coverage

### Validation

This test validates:
- Core workflow can be completed quickly
- All integrations work correctly
- No blocking errors in the flow
- Clean state management

### Use Cases

- Quick regression testing
- CI/CD pipeline integration
- Demonstrating the workflow efficiently
- Testing after code changes

### Duration

Approximately 30 seconds vs 60+ seconds for full test