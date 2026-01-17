---
name: tdd-feature-development
description: This skill should be used when the user asks to "implement a feature", "build new functionality", "add a feature", invokes feature-dev, or discusses feature development workflows. Ensures test-driven development practices are followed by requiring tests to be written early and explicitly included as separate steps in implementation plans.
version: 0.1.0
---

# TDD Feature Development

This skill augments feature development workflows to ensure test-driven development (TDD) practices are consistently followed. It emphasizes writing tests early, planning them as explicit steps, and maintaining high test coverage.

## Purpose

When implementing features, tests should not be an afterthought. This skill ensures that:

1. Tests are planned as **explicit, separate TODO items** (not buried within implementation steps)
2. Tests are written **alongside code** using a TDD-ish approach (not all at the end)
3. Testing principles are followed throughout the development process
4. Code is designed to be testable from the start

## When to Use This Skill

Apply these guidelines when:
- Implementing any new feature or functionality
- Working through a feature development workflow (like feature-dev)
- Creating implementation plans or TODO lists for code changes
- Designing architecture for new features

## Core Testing Principles

### Write Tests Early

Tests are the foundation of code quality. Write tests **early** in the development process, not as a final step.

**TDD-ish Approach:** For core business logic, write tests alongside implementation:
1. Design the interface/API
2. Write the test for expected behavior
3. Implement the functionality
4. Refactor if needed
5. Move to next feature

This is not strict TDD (test-first always), but rather a pragmatic approach that ensures tests are written while the implementation context is fresh.

### Test as Separate TODO Items

When creating implementation plans or TODO lists:

**❌ WRONG:**
```
- Implement user authentication service (including tests)
- Add password hashing
- Create login endpoint
```

**✅ CORRECT:**
```
- Design user authentication service interface
- Write tests for user authentication service
- Implement user authentication service
- Write tests for password hashing
- Add password hashing functionality
- Write tests for login endpoint
- Create login endpoint
- Run all tests and verify 80% coverage
```

Each test-writing activity should be a **separate, explicit TODO item** so it's tracked and not forgotten.

### Test Behavior, Not Implementation (CRITICAL)

**Tests should verify WHAT the code does (behavior/features), not HOW it does it (implementation details).**

This is one of the most important testing principles. Tests should:
- Focus on the **public API** and observable behavior
- Verify **features and outcomes**, not internal mechanics
- Be resilient to refactoring (implementation changes shouldn't break tests)

**Good Examples:**
- Test that `authenticateUser()` returns a valid token for correct credentials
- Test that invalid credentials return an appropriate error
- Test that expired tokens are rejected

**Bad Examples:**
- Test that `authenticateUser()` calls specific internal helper methods
- Test that a private `_hashPassword()` function is invoked
- Test the order of internal operations

If the implementation changes but the behavior stays the same, tests should still pass.

### Don't Test Private Functions (CRITICAL)

**Never write tests for private/internal functions or methods.**

Test only the **public API** and observable behavior. Private functions are implementation details that:
- Should be tested indirectly through the public API
- May change frequently during refactoring
- Create brittle, implementation-coupled tests

If a private function feels like it needs its own tests, consider:
1. Is it complex enough to be extracted as its own public module?
2. Is the public API properly testing all the behaviors this private function enables?

**The solution is better design, not testing private functions.**

### Aim for 80% Test Coverage

Target 80% code coverage as a baseline for quality. This ensures most code paths are exercised without requiring 100% coverage on trivial code.

## Integration with Feature Development Workflows

### During Phase 1: Discovery

When understanding requirements, consider:
- What behaviors need to be tested?
- Are there edge cases that need test coverage?
- What's the testability of the proposed solution?

### During Phase 4: Architecture Design

When evaluating approaches, consider:
- Is this design testable?
- Can core logic be isolated for unit testing?
- Are dependencies injected to allow mocking?

Prefer architectures that enable easy testing. Code that's hard to test is often poorly designed.

### During Phase 5: Implementation

**CRITICAL:** When creating the implementation plan and TODO list:

1. **Break down each feature into implementation + test pairs**
2. **Add separate TODO items for writing tests**
3. **Place test items BEFORE or ALONGSIDE their implementation items**

For example, when implementing a new API endpoint:
```
- Design endpoint interface and response schema
- Write tests for successful request case
- Write tests for validation errors
- Write tests for authentication failures
- Implement endpoint handler
- Run tests and verify they pass
```

**NOT:**
```
- Implement endpoint
- Add tests later
```

4. **Include a final TODO for running tests and checking coverage**

### During Phase 6: Quality Review

Verify:
- Are tests written for all new functionality?
- Do tests cover edge cases and error conditions?
- Is test coverage at or above 80%?
- Are tests testing behavior rather than implementation?

## Project-Agnostic Guidelines

This skill does NOT specify:
- How to run tests (pytest, npm test, etc.) - that's project-specific
- Where test files should be located - that's project-specific
- What testing framework to use - that's project-specific

This skill DOES specify:
- Tests must be written early
- Tests must be planned as explicit TODO items
- Tests must be written alongside code, not at the end
- Test coverage should reach 80%
- Tests should focus on behavior, not implementation

## Key Reminders

1. **Tests are not optional** - They're part of the feature implementation
2. **Tests are not afterthoughts** - Plan them from the start
3. **Tests are separate steps** - Don't hide them in implementation TODOs
4. **Tests guide design** - Hard-to-test code is often bad design
5. **Run tests frequently** - Catch issues early

## Practical Workflow

When creating an implementation plan:

1. **Review the architecture/design**
2. **Identify all components that need implementation**
3. **For each component:**
   - Add TODO: "Write tests for [component]"
   - Add TODO: "Implement [component]"
4. **Add TODO: "Run full test suite and verify 80% coverage"**
5. **Begin implementation**, marking TODOs complete as you go

This ensures tests are visible, tracked, and completed alongside implementation rather than being forgotten or delayed until the end.

---

By following these guidelines, feature development naturally incorporates testing as a first-class activity, resulting in higher quality, more maintainable code with fewer bugs.
