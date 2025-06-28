# VHS Test Analysis Guide

This guide helps interpret VHS test results and identify issues from the generated GIFs and screenshots.

## How to Analyze Test Results

### 1. Locate Test Artifacts

After running a test, you'll find:
- **GIF file**: `output/<test-name>.gif` - The full interaction recording
- **PNG files**: `output/<test-name>-*.png` - Screenshots at key moments
- **Test spec**: `tests/<category>/<test>.test.md` - Expected behavior checklist

### 2. Review Process

1. **Open the test specification** (`.test.md` file)
2. **Play the GIF** alongside the specification
3. **Check each screenshot** against its corresponding checklist
4. **Note any deviations** from expected behavior

### 3. Common Issues to Look For

#### Visual Issues
- [ ] **Rendering artifacts**: Broken borders, missing characters
- [ ] **Theme inconsistencies**: Wrong colors, hardcoded values
- [ ] **Layout problems**: Misaligned elements, overflow
- [ ] **Missing UI elements**: Icons, borders, text

#### Behavioral Issues
- [ ] **Wrong navigation**: Keys don't work as expected
- [ ] **State problems**: UI doesn't update properly
- [ ] **Focus issues**: Wrong element highlighted
- [ ] **Timing problems**: Animations too fast/slow

#### Content Issues
- [ ] **Wrong text**: Incorrect labels or messages
- [ ] **Missing data**: Expected content not shown
- [ ] **Placeholder issues**: Generic instead of context-aware
- [ ] **Truncation**: Text cut off unexpectedly

### 4. Analysis Checklist Template

Use this template when reviewing tests:

```markdown
## Test Analysis: [Test Name]
Date: [Date]
Reviewer: [Name]

### Overall Result
- [ ] PASS
- [ ] FAIL
- [ ] PARTIAL (some issues)

### Findings

#### Screenshot 1: [Name]
- [ ] All checkpoints pass
Issues found:
- 

#### Screenshot 2: [Name]
- [ ] All checkpoints pass
Issues found:
- 

### Summary
- Total checkpoints: X
- Passed: Y
- Failed: Z

### Recommended Actions
1. 
2. 
```

## Automated Analysis Ideas

For future enhancement, consider:

1. **Image diff tools**: Compare screenshots against baseline
2. **OCR verification**: Extract text from screenshots
3. **Color analysis**: Verify theme compliance
4. **Performance metrics**: Measure response times

## Tips for Writing Good Test Specs

1. **Be specific**: "Border is blue" not "Border looks correct"
2. **Use coordinates**: "Top-left corner" not "somewhere on screen"
3. **Check states**: Active, inactive, hover, focused
4. **Verify transitions**: Not just final state
5. **Include negative cases**: What shouldn't happen

## Common Patterns

### Navigation Tests
- Start position
- Each movement
- Wrap-around behavior
- Selection feedback

### Form Tests
- Empty state
- Typing
- Validation
- Submission
- Error states

### Feature Tests
- Initial state
- User action
- Result state
- Edge cases

## Debugging Failed Tests

1. **Run test with slower timing**: Increase Sleep durations
2. **Add more screenshots**: Capture intermediate states
3. **Check terminal size**: Some issues are size-dependent
4. **Verify prerequisites**: Clean state, correct setup
5. **Run interactively**: Execute commands manually

## CI/CD Integration

For automated testing:

```yaml
# Example GitHub Action
- name: Run VHS Tests
  run: |
    cd vhs-tests
    make test-all
    
- name: Upload Results
  uses: actions/upload-artifact@v3
  with:
    name: vhs-test-results
    path: vhs-tests/output/
```

Then manually review artifacts or use image comparison tools.