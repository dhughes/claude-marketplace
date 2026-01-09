# Information Synthesis and Reporting Guide

Guide for connecting research findings from multiple sources and presenting comprehensive, actionable results.

## Overview

Research at ezCater often involves gathering information from 4-5 different sources (Glean, Jira, Confluence, GitHub, Git). The value comes from synthesizing these disparate pieces into a coherent narrative that answers the original question and anticipates follow-ups.

## Synthesis Process

### 1. Organize Findings by Source

As you research, maintain a mental (or actual) structure of what you find:

**Glean findings**:
- Design docs about the feature
- Slack discussions about decisions
- Meeting notes mentioning considerations

**Jira findings**:
- Original epic or story
- Acceptance criteria
- Implementation sub-tasks
- Related bugs or improvements

**Confluence findings**:
- RFCs or ADRs explaining choices
- Architecture diagrams
- Process documentation

**GitHub findings**:
- Initial implementation PR
- Subsequent modification PRs
- Review comments revealing intent

**Git findings**:
- Commit messages explaining changes
- Code evolution showing iteration
- Blame revealing authors/owners

### 2. Create Timeline

Order findings chronologically to show evolution:

1. **Original need** (when/why feature was needed)
   - Source: Jira epic, Glean meeting notes
   - Key info: Business driver, timeline

2. **Design decisions** (how it would be built)
   - Source: Confluence RFC, Slack discussions
   - Key info: Alternatives considered, trade-offs

3. **Initial implementation** (first version)
   - Source: GitHub PR, git commits
   - Key info: What was built, who built it

4. **Iterations** (changes over time)
   - Source: Subsequent PRs, Jira improvements
   - Key info: What changed and why

5. **Current state** (today)
   - Source: Recent commits, current code
   - Key info: Latest version, recent changes

### 3. Identify Stakeholders

Note who was involved and their roles:

- **Business owner**: Who requested the feature (from Jira, Glean)
- **Technical lead**: Who made architecture decisions (from Confluence, Slack)
- **Implementers**: Who wrote the code (from git, GitHub)
- **Reviewers**: Who approved changes (from GitHub PRs)

**Why it matters**: Knowing stakeholders helps answer "who can I ask about this?" and provides confidence in findings.

### 4. Understand Business Context

Connect technical decisions to business needs:

- **Problem**: What business problem was being solved?
- **Constraints**: What limitations existed (time, resources, existing systems)?
- **Requirements**: What were the must-haves vs nice-to-haves?
- **Success criteria**: How was success measured?

**Sources**:
- Jira epic descriptions
- Confluence PRDs and project docs
- Glean meeting notes
- PR descriptions linking to business context

### 5. Recognize Trade-offs

Every decision involves trade-offs. Identify:

**What was gained**:
- Performance improvements
- Better maintainability
- New capabilities
- Simplified code

**What was sacrificed**:
- Other features delayed
- Increased complexity elsewhere
- Additional dependencies
- Technical debt accepted

**Why the trade-off made sense**:
- Business priority
- Time constraints
- Team expertise
- Existing architecture

**Sources**:
- RFC/ADR documents in Confluence
- PR review comments in GitHub
- Slack discussions in Glean
- Commit messages explaining choices

### 6. Cross-Reference and Verify

Confirm findings across multiple sources:

1. **Check consistency**: Do Jira, Confluence, and GitHub tell the same story?
2. **Resolve conflicts**: If sources disagree, check dates and authority
3. **Identify gaps**: What questions remain unanswered?
4. **Verify currency**: Is information current or outdated?

**Red flags**:
- Dates don't align across sources
- Design doc says one thing, code does another
- Discussion mentions concerns that weren't addressed
- Recent changes contradict original design

## Report Structure

### Executive Summary

**Purpose**: Answer the question directly in 2-3 sentences

**Format**:
```
[Answer to the original question]

[Key supporting fact 1]. [Key supporting fact 2].
```

**Example**:
```
Liberty uses Omnichannel for feature flags because it provides unified
management across all ezCater applications (mobile and web).

The decision was made in Q2 2024 as part of the platform consolidation
initiative. Omnichannel's real-time updates and centralized control
addressed limitations in the previous per-app flag systems.
```

### Key Findings

**Purpose**: Provide detailed evidence and context

**Format**:
```
- Finding 1: [Specific fact with evidence]
  - Source: [Where this came from]
  - Context: [Why this matters]

- Finding 2: [Another specific fact]
  - Source: [Evidence]
  - Context: [Significance]
```

**Example**:
```
- Initial implementation in January 2024 (PR #4521)
  - Source: GitHub PR ezcater/liberty#4521
  - Context: This was a spike/proof-of-concept to validate Omnichannel integration

- Design doc outlining the migration strategy
  - Source: Confluence page "Feature Flag Consolidation RFC"
  - Context: Documented trade-offs between Omnichannel, Eppo, and LaunchDarkly

- Previous system (per-app flags) caused inconsistencies
  - Source: Jira epic FX-3456, Slack #mobile-backend discussions
  - Context: Different flag values between web and mobile led to UX issues
```

### Historical Context (when relevant)

**Purpose**: Show evolution over time

**Format**:
```
Timeline:
- [Date]: [Event] - [Significance]
- [Date]: [Event] - [Significance]
```

**Example**:
```
Timeline:
- Q4 2023: Mobile team identifies flag inconsistencies between Liberty and ez-rails
- Jan 2024: Platform team evaluates Omnichannel as unified solution (Confluence RFC)
- Feb 2024: Spike implementation proves feasibility (PR #4521)
- Mar 2024: Full migration begins (Epic FX-3456)
- May 2024: Migration complete, per-app flags deprecated
```

### Technical Details (when needed)

**Purpose**: Explain how it works for technical audience

**Format**:
```
Architecture:
[Brief explanation of structure]

Key components:
- [Component 1]: [Purpose]
- [Component 2]: [Purpose]

Integration points:
- [How it connects to system A]
- [How it connects to system B]
```

**Include**:
- Code snippets (from git)
- Architecture diagrams (from Confluence)
- Configuration examples (from code)
- API contracts (from code or docs)

### References

**Purpose**: Enable user to verify and dig deeper

**Format**:
```
References:
- Jira: [TICKET-ID] - [Brief description]
- PRs: [REPO#NUM] - [What it implemented]
- Commits: [SHA] - [Significant change]
- Confluence: [Page title] (ID: [page-id])
- Glean: [Type of document found]
```

**Example**:
```
References:
- Jira: FX-3456 - Feature flag consolidation epic
- PRs:
  - ezcater/liberty#4521 - Initial Omnichannel spike
  - ezcater/liberty#4782 - Full integration
- Commits:
  - abc123d - Remove LaunchDarkly dependency
  - def456e - Add Omnichannel client initialization
- Confluence: "Feature Flag Consolidation RFC" (ID: 123456789)
- Glean: Slack discussions in #mobile-backend (March 2024)
```

## Quality Standards

### Specificity

✅ **Good**: "Initial implementation in PR #4521 merged on January 15, 2024"
❌ **Bad**: "It was implemented earlier this year"

✅ **Good**: "The team chose Omnichannel over Eppo because Omnichannel supported real-time updates"
❌ **Bad**: "They picked Omnichannel for some reason"

### Evidence

Always cite sources:

✅ **Good**: "According to the RFC in Confluence (ID: 123456789), three alternatives were evaluated"
❌ **Bad**: "Several alternatives were probably considered"

✅ **Good**: "Git blame shows this code was added by John Doe in commit abc123d"
❌ **Bad**: "Someone added this code at some point"

### Acknowledgment of Gaps

Be explicit about what you don't know:

✅ **Good**: "The RFC mentions performance concerns but doesn't provide specific benchmarks. Further investigation would require load testing results."
❌ **Bad**: [Silently omitting information you couldn't find]

✅ **Good**: "Design discussions likely happened in Slack, but messages from that period are not available. The Confluence RFC represents the final decision."
❌ **Bad**: [Presenting the RFC as the complete story when discussions are missing]

### Facts vs Speculation

Clearly distinguish:

✅ **Good**: "The code suggests this was a quick fix (small diff, no tests), though the PR description doesn't explicitly state this was temporary."
❌ **Bad**: "This was definitely a hack job."

✅ **Good**: "Based on the timing (2 weeks before deadline per Jira) and the PR comments ('will refactor later'), this appears to have been done under time pressure."
❌ **Bad**: "They rushed this because they were lazy."

## Anticipating Follow-ups

Think about what the user might ask next:

### Original Question: "Why was X chosen?"

**Anticipate**:
- "What were the alternatives?"
- "Are there any downsides?"
- "Could we switch to Y later?"
- "Who maintains X?"

**Address proactively**:
```
Alternatives considered: [List with brief pros/cons]
Current limitations: [Known issues from Jira/GitHub]
Migration feasibility: [Based on architecture complexity]
Ownership: [Team name from Confluence, key contributors from git]
```

### Original Question: "How does process X work?"

**Anticipate**:
- "Where is this implemented in code?"
- "What triggers this process?"
- "What happens if it fails?"
- "How do I test this?"

**Address proactively**:
```
Implementation: [File paths, key functions]
Triggers: [Events, cron jobs, user actions]
Error handling: [From code and Confluence docs]
Testing: [Test fixtures, dev environment setup]
```

### Original Question: "What's the history of feature Y?"

**Anticipate**:
- "Why was it built this way?"
- "Has it changed significantly?"
- "Are there known issues?"
- "What's planned next?"

**Address proactively**:
```
Original design rationale: [From RFC/ADR]
Major iterations: [PRs showing evolution]
Current issues: [Open bugs in Jira]
Future plans: [Upcoming epics, backlog]
```

## Common Synthesis Patterns

### Pattern 1: Architecture Decision

**Question type**: "Why was X chosen over Y?"

**Synthesis**:
1. Find alternatives considered (Confluence RFC, Slack)
2. List pros/cons of each (RFC, review comments)
3. Identify deciding factors (business needs, constraints)
4. Show implementation (initial PR)
5. Note any regrets/revisions (later PRs, Jira improvements)

### Pattern 2: Code Evolution

**Question type**: "How did feature X change over time?"

**Synthesis**:
1. Original implementation (git, GitHub PR)
2. Modifications and why (subsequent PRs, commits)
3. Bug fixes revealing issues (Jira, PRs)
4. Refactorings improving design (PRs with "refactor")
5. Current state and recent changes (git log recent)

### Pattern 3: Process Understanding

**Question type**: "How does process X work?"

**Synthesis**:
1. Process documentation (Confluence)
2. Code implementation (git, grep)
3. Example usage (test fixtures, staging environment)
4. Common issues (Jira bugs, Glean discussions)
5. Recent changes (git log, recent PRs)

### Pattern 4: Multi-System Integration

**Question type**: "How do systems A and B interact?"

**Synthesis**:
1. Architecture overview (Confluence diagrams)
2. API contracts (code, OpenAPI specs)
3. Integration PRs (GitHub)
4. Message flows (code, documentation)
5. Known integration issues (Jira, incidents)

## Presentation Tips

### Use Formatting

**Lists for facts**:
```
- Implementation: January 2024
- Author: John Doe
- PR: #4521
```

**Timeline for evolution**:
```
1. Q4 2023: Original need identified
2. Q1 2024: Design phase
3. Q2 2024: Implementation
```

**Tables for comparisons**:
```
| Option | Pros | Cons | Decision |
|--------|------|------|----------|
| X | Fast, simple | Limited | Rejected |
| Y | Full-featured | Complex | Chosen |
```

### Link to Evidence

Make it easy to verify:

```
The decision to use Omnichannel is documented in:
- Confluence: "Feature Flag Consolidation RFC" (page ID: 123456789)
- Jira: FX-3456
- PR: ezcater/liberty#4521
```

### Highlight Important Points

```
**Critical finding**: The original design assumed real-time updates, but
current implementation uses 5-minute polling due to infrastructure constraints.
```

## Summary

Effective synthesis and reporting requires:
1. **Organized gathering**: Track sources as you research
2. **Timeline creation**: Show evolution chronologically
3. **Stakeholder identification**: Note who was involved
4. **Context understanding**: Connect technical to business
5. **Trade-off recognition**: Explain what was gained and sacrificed
6. **Verification**: Cross-reference across sources
7. **Structured reporting**: Executive summary → details → references
8. **Quality standards**: Specific, cited, honest about gaps
9. **Anticipation**: Address obvious follow-up questions

The goal is not just to answer the question, but to provide the complete story that helps the user understand the "why" behind the "what" and make informed decisions going forward.
