# TDD Feature Development Plugin

Augments feature development workflows in Claude Code with test-driven development (TDD) best practices. Ensures tests are written early, explicitly planned as separate TODO items, and integrated throughout the development lifecycle.

## Overview

This plugin provides a skill that automatically composes with feature development workflows (like `feature-dev:feature-dev`) to inject TDD principles. Instead of treating tests as an afterthought or final step, this skill ensures testing is a first-class activity throughout feature implementation.

## What It Does

When you're implementing features, this skill automatically ensures:

1. **Tests are explicit TODO items** - Each test-writing task is tracked separately, not buried in implementation steps
2. **Tests are written early** - Tests are written alongside code using a TDD-ish approach, not all at the end
3. **Tests focus on behavior** - Emphasizes testing features and public APIs, not implementation details or private functions
4. **Coverage targets are met** - Aims for 80% test coverage as a quality baseline

## When It Activates

This skill automatically loads when you:
- Ask to "implement a feature" or "build new functionality"
- Invoke the `/feature-dev` command
- Discuss feature development workflows
- Create implementation plans or TODO lists

The skill composes with other feature development tools, augmenting them with testing best practices.

## Core Principles

### Write Tests Early

Tests are written **alongside** implementation, not as a final step. For core business logic:
1. Design the interface/API
2. Write the test for expected behavior
3. Implement the functionality
4. Refactor if needed
5. Move to next feature

### Tests as Separate TODO Items

**Wrong:**
```
- Implement user authentication (including tests)
- Add password hashing
```

**Right:**
```
- Design user authentication interface
- Write tests for user authentication
- Implement user authentication
- Write tests for password hashing
- Add password hashing
- Run all tests and verify 80% coverage
```

### Test Behavior, Not Implementation

**Critical:** Tests should verify **WHAT** the code does (features/behavior), not **HOW** it does it (implementation).

- ✅ Test that `authenticateUser()` returns a valid token for correct credentials
- ❌ Test that `authenticateUser()` calls specific internal helper methods

### Don't Test Private Functions

**Never write tests for private/internal functions.** Test only the public API. If a private function feels like it needs tests, that's a design smell—consider extracting it as a public module or ensuring the public API tests all its behaviors.

### Aim for 80% Coverage

Target 80% code coverage as a quality baseline without requiring 100% coverage on trivial code.

## Integration with Feature Development

### Phase 1: Discovery
- Consider what behaviors need testing
- Identify edge cases requiring test coverage
- Evaluate solution testability

### Phase 4: Architecture Design
- Prefer testable architectures
- Ensure dependencies can be mocked
- Isolate core logic for unit testing

### Phase 5: Implementation
**This is where the skill has the most impact:**
- Breaks features into implementation + test pairs
- Creates separate TODO items for each test
- Places test items before/alongside implementation items
- Includes final TODO for running tests and checking coverage

### Phase 6: Quality Review
- Verifies tests exist for all new functionality
- Checks edge case and error condition coverage
- Validates 80% test coverage target
- Confirms tests focus on behavior, not implementation

## Project-Agnostic

This skill is **project-agnostic** and works with any testing framework:
- Doesn't specify how to run tests (pytest, npm test, etc.)
- Doesn't dictate test file locations
- Doesn't prescribe testing frameworks

It focuses on the **principles** of when and what to test, not the mechanics.

## Example TODO List

When implementing a new API endpoint, the skill ensures your TODO list looks like:

```
✓ Understand requirements and existing patterns
✓ Design endpoint interface and response schema
⊙ Write tests for successful request case
⊙ Write tests for validation errors
⊙ Write tests for authentication failures
○ Implement endpoint handler
○ Run tests and verify they pass
○ Review code quality
```

Notice how test-writing is **explicit, visible, and tracked** throughout implementation.

## Installation

Install from the marketplace:

```bash
/plugins install tdd-feature-dev
```

Or manually:

```bash
git clone <marketplace-url>
cp -r plugins/tdd-feature-dev ~/.claude/plugins/
```

## Usage

Simply use feature development workflows as normal:

```bash
/feature-dev Add user authentication to the API
```

The TDD skill will automatically compose with feature-dev, ensuring testing best practices are followed throughout.

## Benefits

- **Higher code quality** - Tests catch bugs early
- **Better design** - Testable code is usually better designed
- **Explicit tracking** - Tests aren't forgotten or delayed
- **Maintainability** - Behavior-focused tests survive refactoring
- **Confidence** - 80% coverage provides quality baseline

## License

MIT

## Author

Doug Hughes (doug.hughes@ezcater.com)

## Version

0.1.0
