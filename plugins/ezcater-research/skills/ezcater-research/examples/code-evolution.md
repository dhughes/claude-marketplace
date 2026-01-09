# Example: Tracing Code Evolution

Complete walkthrough of investigating "What's the history behind the authentication system refactor in ez-rails?"

## Initial Question

**User asks**: "What's the history behind the authentication system refactor in ez-rails?"

**Research goal**: Understand when, why, and how the authentication system changed over time

## Step 1: Find Current Authentication Code

Start by understanding current implementation:

```bash
# Assume ez-rails checked out at ~/code/ezcater/ez/rails
cd ~/code/ezcater/ez/rails

# Find authentication-related files
find app -name "*auth*" -type f | head -20
```

**Results**:
```
app/controllers/concerns/authenticatable.rb
app/services/authentication_service.rb
app/models/user_session.rb
spec/services/authentication_service_spec.rb
```

Key file: `app/services/authentication_service.rb`

## Step 2: Trace File History

See how the main file evolved:

```bash
git log --follow --oneline app/services/authentication_service.rb
```

**Output**:
```
f1e2d3c (2024-10-15) Refactor: Extract token validation logic
e2d3c4b (2024-09-20) Add support for API key authentication
d3c4b5a (2024-06-10) Major refactor: Simplify authentication flow
c4b5a6f (2024-02-15) Add JWT token support
b5a6f7e (2023-11-01) Extract authentication logic from controller
a6f7e8d (2023-08-15) Initial authentication service
```

Notice major refactor at `d3c4b5a` (June 2024)

## Step 3: Examine Major Refactor Commit

View the significant refactor:

```bash
git show d3c4b5a --stat
```

**Output**:
```
commit d3c4b5a
Author: Alex Johnson <alex@ezcater.com>
Date: 2024-06-10

Refactor authentication system for better testability and maintainability

This refactors the authentication system to:
- Separate concerns between authentication and authorization
- Introduce strategy pattern for different auth methods
- Improve error handling and logging
- Reduce complexity in controllers

Closes FX-2341

 app/services/authentication_service.rb          | 150 +++++++++---------
 app/services/auth/jwt_strategy.rb               |  45 ++++++
 app/services/auth/api_key_strategy.rb           |  38 +++++
 app/controllers/concerns/authenticatable.rb     |  62 +++----
 spec/services/authentication_service_spec.rb    | 120 ++++++++------
 8 files changed, 289 insertions(+), 184 deletions(-)
```

Key info: Refactor introduces strategy pattern, references ticket FX-2341

## Step 4: Find Related Jira Ticket

Get business context:

```bash
atl jira issue view FX-2341
```

**Key information**:
```
Title: Refactor Authentication System for Better Maintainability
Type: Technical Debt
Created: 2024-05-01
Completed: 2024-06-15

Description:
The current authentication system has become difficult to maintain and test.
Multiple authentication methods (sessions, JWT, API keys) are tangled together,
making it hard to add new auth methods or modify existing ones.

Goal: Refactor to use strategy pattern for different auth methods.

Acceptance Criteria:
- Each auth method in separate strategy class
- 100% test coverage maintained
- No breaking changes to existing API
- Performance impact < 5ms per request
```

**Comments include**:
- Discussion of strategy pattern vs factory pattern
- Performance benchmarks from testing
- Migration plan for existing code

## Step 5: Find Related PRs

Search GitHub for refactor PRs:

```bash
gh pr list --repo ezcater/ez-rails --search "is:merged authentication refactor" --limit 10
```

**Results**:
```
#5234 - Refactor authentication system (merged 2024-06-12)
#5123 - Add comprehensive auth tests (merged 2024-06-05)
#5089 - Extract JWT handling (merged 2024-05-28)
```

View main refactor PR:

```bash
gh pr view 5234 --repo ezcater/ez-rails --comments
```

**Key information from PR**:
- **Description**: Links to FX-2341, explains refactor motivation
- **Review comments**:
  - Security team approval required
  - Performance testing results shared
  - Discussion of backward compatibility
- **Files changed**: 15 files, mostly in `app/services/auth/`
- **Tests**: Added 200+ lines of new test coverage

## Step 6: Understand Pre-Refactor State

Look at file before refactor:

```bash
git show d3c4b5a~1:app/services/authentication_service.rb | head -50
```

**Observations**:
- Single large class (~400 lines)
- Multiple authentication methods in one file
- Complex conditional logic for different auth types
- Difficult to test individual auth methods

## Step 7: Search for Original Implementation

Find when authentication service was created:

```bash
git log --diff-filter=A --follow -- app/services/authentication_service.rb
```

**Result**:
```
a6f7e8d (2023-08-15) Initial authentication service
Author: Maria Garcia
```

View original implementation:

```bash
git show a6f7e8d
```

**Initial version**:
- Extracted from controllers (mentioned in commit message)
- Only handled session-based auth
- ~100 lines
- Basic structure

## Step 8: Trace Evolution Between Original and Refactor

List all commits between creation and refactor:

```bash
git log --oneline a6f7e8d..d3c4b5a -- app/services/authentication_service.rb
```

**Evolution timeline**:
1. **Aug 2023** (a6f7e8d): Initial extraction from controllers
2. **Nov 2023** (b5a6f7e): Enhanced with better error handling
3. **Feb 2024** (c4b5a6f): Added JWT token support for mobile API
4. **Jun 2024** (d3c4b5a): Major refactor with strategy pattern

## Step 9: Find Related Discussions

Search Glean for context:

```
mcp__glean__company_search
query: "authentication refactor ez-rails strategy pattern"
datasources: ["slack", "confluence"]
```

**Found**:
- Slack discussion in #backend-engineering (May 2024)
- Confluence page: "Auth System Technical Debt Review"

Slack discussion revealed:
- Team struggled with adding API key auth due to existing complexity
- Refactor proposed at architecture review meeting
- Decision to use strategy pattern based on Ruby best practices

## Step 10: Check Subsequent Changes

See what changed after refactor:

```bash
git log --since="2024-06-10" --oneline -- app/services/auth*
```

**Results**:
```
f1e2d3c (Oct 2024) Extract token validation logic
e2d3c4b (Sep 2024) Add support for API key authentication
```

These show the refactor enabled easier extension (API key auth added in Sep)

## Step 11: Identify Related Bug Fixes

Find bugs that motivated refactor:

```bash
atl jira search "project = FX AND type = Bug AND text ~ 'authentication' AND created >= '2024-01-01' AND created <= '2024-05-01'"
```

**Found bugs**:
- FX-2156: JWT tokens expiring incorrectly
- FX-2234: Session auth failing for API requests
- FX-2298: Difficult to add OAuth support

These bugs highlighted the complexity issues

## Synthesis: Complete History

### Timeline

1. **August 2023**: Initial extraction
   - Pulled authentication logic from controllers
   - Basic session-based auth only
   - Author: Maria Garcia
   - Motivation: Controller complexity reduction

2. **November 2023**: Enhancement
   - Better error handling
   - Improved logging
   - Small incremental improvement

3. **February 2024**: JWT support added
   - Mobile app required API authentication
   - JWT logic added to existing service
   - Started showing complexity issues (per PR comments)

4. **January-May 2024**: Problems accumulate
   - Multiple bugs filed (FX-2156, FX-2234, FX-2298)
   - Difficulty adding new features
   - Testing became challenging

5. **May 2024**: Refactor proposed
   - Architecture review identifies technical debt
   - Ticket FX-2341 created
   - Strategy pattern chosen for solution

6. **June 2024**: Major refactor
   - Strategy pattern implementation
   - Separate classes for each auth method
   - Improved testability
   - Backward compatible

7. **Post-refactor**: Easier evolution
   - API key auth added easily (September 2024)
   - Token validation extracted (October 2024)
   - Demonstrates improved maintainability

### Key Changes

**Original structure** (Aug 2023 - Jun 2024):
```ruby
# Single monolithic class
class AuthenticationService
  def authenticate(request)
    if session_auth?
      # session logic
    elsif jwt_auth?
      # JWT logic
    elsif api_key_auth?
      # API key logic
    end
  end
end
```

**Refactored structure** (Jun 2024 - present):
```ruby
# Strategy pattern
class AuthenticationService
  def authenticate(request)
    strategy = select_strategy(request)
    strategy.authenticate(request)
  end
end

# Separate strategy classes
class Auth::JwtStrategy; end
class Auth::SessionStrategy; end
class Auth::ApiKeyStrategy; end
```

### Motivation for Refactor

**Immediate triggers**:
- Bug FX-2234: Difficulty distinguishing between auth types
- Need to add OAuth support (future requirement)
- Testing complexity growing

**Underlying causes**:
- Gradual feature additions without redesign
- Multiple authentication methods in one class
- Tight coupling between auth types

**Business impact**:
- Mobile app auth issues
- Delayed OAuth integration for SSO
- Slow bug resolution due to test complexity

### Stakeholders

- **Maria Garcia**: Original implementation (2023)
- **Alex Johnson**: Refactor implementation (2024)
- **Backend team**: Code reviews and testing
- **Security team**: Approved changes
- **Mobile team**: Benefited from JWT improvements

## Final Report

### Executive Summary

The authentication system in ez-rails underwent a major refactor in June 2024 to address growing complexity and maintainability issues. Originally extracted from controllers as a simple session auth service in August 2023, it evolved to support multiple authentication methods (sessions, JWT, API keys), but the monolithic structure made testing and extending difficult. The refactor introduced a strategy pattern, separating each auth method into its own class, which has since enabled easier feature additions.

### Key Findings

- **Original implementation** (Aug 2023): Simple extraction from controllers for session-based auth (commit a6f7e8d by Maria Garcia)

- **Gradual complexity** (Aug 2023 - Feb 2024): JWT support added for mobile API, increasing code from ~100 to ~400 lines in single class

- **Problem accumulation** (Jan-May 2024): Multiple bugs filed (FX-2156, FX-2234, FX-2298) highlighting difficulty maintaining and testing the monolithic structure

- **Strategy pattern refactor** (Jun 2024): PR #5234 by Alex Johnson separated concerns, created Auth::JwtStrategy, Auth::SessionStrategy, Auth::ApiKeyStrategy classes

- **Successful outcome**: Post-refactor additions (API key auth in Sep 2024) demonstrated improved maintainability, as noted in PR #5267

### Evolution Diagram

```
Aug 2023          Feb 2024          Jun 2024          Present
Simple ---------> Complex --------> Refactored -----> Extended
Session only      + JWT             Strategy          + OAuth (planned)
~100 lines        ~400 lines        pattern           Easy additions
                  Growing bugs      Modular design
```

### Technical Debt Addressed

**Before refactor**:
- Single 400-line class
- Conditional logic for auth types
- Difficult to test individual methods
- Hard to add new auth types

**After refactor**:
- Separate strategy classes (~50 lines each)
- Clear separation of concerns
- Independent testing per strategy
- New auth methods easily added

### References

**Git commits**:
- a6f7e8d - Initial authentication service (Aug 2023)
- c4b5a6f - Add JWT token support (Feb 2024)
- d3c4b5a - Major refactor with strategy pattern (Jun 2024)

**Jira**:
- FX-2341 - Refactor authentication system (Epic)
- FX-2156, FX-2234, FX-2298 - Bugs that motivated refactor

**GitHub PRs**:
- ezcater/ez-rails#5234 - Authentication refactor (Jun 2024)
- ezcater/ez-rails#5267 - API key auth (Sep 2024, shows benefit)

**Confluence**:
- "Auth System Technical Debt Review" - Architecture analysis

**Glean**:
- Slack #backend-engineering discussions (May 2024)

### Anticipated Follow-ups

**Q: Are there similar refactors planned for other systems?**
A: The Confluence doc "Auth System Technical Debt Review" mentions authorization system as next candidate, but no ticket created yet.

**Q: Was there any performance impact from the refactor?**
A: Per FX-2341 acceptance criteria and PR comments, performance impact was < 2ms per request, well under the 5ms threshold.

**Q: What's preventing OAuth implementation now?**
A: Per recent Slack discussions, OAuth is on the Q1 2025 roadmap. The refactor removed the technical blocker; now it's a prioritization decision.

## Research Techniques Demonstrated

1. **Git history tracing**: Used `git log --follow` to track file evolution through renames
2. **Commit analysis**: Examined key commits with `git show` to understand changes
3. **Ticket correlation**: Connected git commits to Jira tickets via commit messages
4. **PR investigation**: Used `gh pr view` to get implementation context and reviews
5. **Timeline construction**: Built chronological narrative from commits, PRs, and tickets
6. **Pattern recognition**: Identified accumulation of technical debt over time
7. **Impact assessment**: Showed post-refactor benefits with subsequent easy additions

This example demonstrates how to trace code evolution from initial implementation through refactoring using git, GitHub, Jira, and Glean.
