# Example: Researching Architectural Decision

Complete walkthrough of investigating "Why does Liberty use Omnichannel for feature flags instead of ezCater's feature flag packages?"

## Initial Question

**User asks**: "Why does Liberty use Omnichannel for feature flags instead of ezCater's feature flag packages like Eppo?"

**Research goal**: Understand the technical decision, alternatives considered, and rationale

## Step 1: Glean Search for Context

Start broad to understand what exists:

```
mcp__glean__company_search
query: "Liberty Omnichannel feature flags"
datasources: ["confluence", "slack", "gdrive"]
```

**Results found**:
- Confluence page: "Feature Flag Consolidation RFC"
- Slack discussions in #mobile-backend (March 2024)
- Meeting notes: "Platform Engineering Q1 Planning"

**Initial findings**:
- Migration happened Q1 2024
- Decision involved platform team and mobile team
- Multiple providers were evaluated

## Step 2: Read RFC in Confluence

Get detailed design document:

```bash
atl confluence search "feature flag consolidation RFC"
```

**Output**:
```
Found page: "Feature Flag Consolidation RFC" (ID: 123456789)
Author: Jane Smith
Updated: 2024-02-15
```

Read the RFC:

```bash
atl confluence page view 123456789 --output markdown
```

**Key information extracted**:
- **Problem**: Mobile (Liberty) and web (ez-rails) had different feature flag systems
  - Liberty used LaunchDarkly
  - ez-rails used mix of Eppo and custom flags
  - Inconsistent flag values between platforms caused UX issues

- **Requirements**:
  - Unified management across all applications
  - Real-time updates
  - Support for complex targeting rules
  - Integration with existing Backstage tooling

- **Alternatives evaluated**:
  1. **Eppo**: Great analytics, but limited mobile SDK support at time
  2. **LaunchDarkly**: Good mobile support, but expensive at ezCater scale
  3. **Omnichannel**: ezCater-built, unified across platforms, integrates with existing tools

- **Decision**: Omnichannel
  - Rationale: Already proven in ez-rails, mobile SDK available, cost-effective, team expertise

- **Trade-offs**:
  - Pro: Unified system, cost savings, team control
  - Con: Less analytics than Eppo, requires maintaining internal code

## Step 3: Find Jira Epic

Search for implementation tickets:

```bash
atl jira search "project = FX AND text ~ 'feature flag consolidation'"
```

**Results**:
```
FX-3456 - Feature Flag Consolidation (Epic)
FX-3457 - Implement Omnichannel in Liberty (Story)
FX-3458 - Migrate existing LaunchDarkly flags (Story)
FX-3459 - Deprecate LaunchDarkly dependency (Story)
```

View epic:

```bash
atl jira issue view FX-3456
```

**Key information**:
- Created: 2024-01-15
- Epic owner: Mobile team lead
- Business driver: Reduce operational complexity and cost
- Timeline: Q1 2024 (completed May 2024)
- Acceptance criteria included automated testing of flag consistency

## Step 4: Find Implementation PRs

Search GitHub for related PRs:

```bash
gh pr list --repo ezcater/liberty --search "is:merged omnichannel feature flag" --limit 20
```

**Results**:
```
#4521 - Add Omnichannel spike (merged 2024-01-20)
#4782 - Implement full Omnichannel integration (merged 2024-03-10)
#4893 - Migrate LaunchDarkly flags to Omnichannel (merged 2024-04-15)
#4921 - Remove LaunchDarkly dependency (merged 2024-05-05)
```

View initial spike PR:

```bash
gh pr view 4521 --repo ezcater/liberty --comments
```

**Key information**:
- PR description links to RFC and epic (FX-3456)
- Review comments discuss real-time update testing
- Mobile team validated SDK integration
- Performance testing showed acceptable latency

## Step 5: Examine Code Changes

Check git history for details:

```bash
# Assume Liberty checked out at ~/code/ezcater/liberty
cd ~/code/ezcater/liberty

git log --grep="FX-3456" --oneline
```

**Output**:
```
a1b2c3d Add Omnichannel client initialization
b2c3d4e Implement flag evaluation caching
c3d4e5f Add specs for Omnichannel integration
d4e5f6g Migrate existing flags to Omnichannel
e5f6g7h Remove LaunchDarkly configuration
```

View key commit:

```bash
git show a1b2c3d
```

**Details**:
- Added Omnichannel client with configuration
- Implemented fallback for network failures
- Added caching layer for performance
- Commit message references performance requirements from RFC

## Step 6: Check Slack Discussions

Review team discussions:

**From Glean results**, Slack conversations in #mobile-backend (March 2024) included:

- Concerns about real-time updates vs polling
- Discussion of caching strategy
- Questions about backward compatibility
- Decision to use 5-minute polling initially, optimize later

## Step 7: Synthesize Findings

### Timeline

1. **Q4 2023**: Mobile team identifies flag inconsistency issues
2. **January 2024**: Platform team evaluates options, creates RFC
3. **February 2024**: RFC approved, epic created (FX-3456)
4. **January-March 2024**: Spike and full implementation (PRs #4521, #4782)
5. **April-May 2024**: Migration and LaunchDarkly deprecation (PRs #4893, #4921)

### Stakeholders

- **Jane Smith** (Platform team): Wrote RFC, made technical recommendation
- **Mobile team lead**: Epic owner, business case
- **Implementation**: Mobile engineers (from git blame and PR authors)

### Decision Factors

**Why Omnichannel was chosen**:
1. **Unification**: Same system for mobile and web eliminated inconsistencies
2. **Cost**: Significant savings vs LaunchDarkly at ezCater's scale
3. **Integration**: Already integrated with Backstage and other ezCater tools
4. **Control**: Internal system, team can modify as needed
5. **Proven**: Already successful in ez-rails

**Why not Eppo**:
- Mobile SDK support was limited at evaluation time
- Would still require maintaining multiple systems temporarily
- Cost comparable to maintaining Omnichannel

**Why not LaunchDarkly**:
- High cost at ezCater scale
- External dependency for critical path
- Limited customization options

### Trade-offs Accepted

**Gained**:
- Flag consistency across platforms
- Lower operational cost
- Team control and customization ability
- Simplified tooling

**Sacrificed**:
- Advanced analytics (Eppo strength)
- Some enterprise features of LaunchDarkly
- Ongoing maintenance responsibility

## Final Report

### Executive Summary

Liberty uses Omnichannel for feature flags to unify flag management across all ezCater applications (mobile and web). The decision was made in Q1 2024 as part of a platform consolidation initiative that eliminated inconsistencies between Liberty's previous LaunchDarkly system and ez-rails' mixed flag providers.

### Key Findings

- **Unified system requirement**: Previous setup (LaunchDarkly in Liberty, Eppo/custom in ez-rails) caused flag value inconsistencies that led to UX issues between mobile and web platforms

- **Cost savings**: LaunchDarkly's per-seat pricing was expensive at ezCater's scale; Omnichannel as an internal system significantly reduces costs

- **Proven solution**: Omnichannel was already successfully used in ez-rails with mobile SDK support, reducing implementation risk

- **Three alternatives evaluated**: Eppo (limited mobile SDK), LaunchDarkly (high cost), Omnichannel (chosen for unification and cost)

- **Completed migration**: Full implementation from January-May 2024, LaunchDarkly fully deprecated by May 2024

### Technical Details

**Implementation** (from PR #4782):
- Omnichannel client initialization in `lib/feature_flags/omnichannel_client.rb`
- 5-minute polling interval for flag updates (vs real-time goal, performance trade-off)
- Local caching layer for offline resilience
- Backward-compatible API maintained during migration

**Integration points**:
- Backstage integration for flag management UI
- Automated testing validates consistency between Liberty and ez-rails
- Monitoring dashboards track flag evaluation performance

### References

**Confluence**:
- "Feature Flag Consolidation RFC" (ID: 123456789) - Complete evaluation and decision rationale

**Jira**:
- FX-3456 - Epic: Feature Flag Consolidation
- FX-3457, FX-3458, FX-3459 - Implementation stories

**GitHub PRs (ezcater/liberty)**:
- #4521 - Initial Omnichannel spike (Jan 2024)
- #4782 - Full Omnichannel integration (Mar 2024)
- #4893 - Flag migration (Apr 2024)
- #4921 - LaunchDarkly removal (May 2024)

**Git commits**:
- a1b2c3d - Add Omnichannel client initialization
- e5f6g7h - Remove LaunchDarkly configuration

**Glean sources**:
- Slack #mobile-backend discussions (March 2024) about implementation details

### Anticipated Follow-ups

**Q: Could we switch back to Eppo or another provider?**
A: Possible but significant effort. Would require reversing the unification benefits and reintroducing platform-specific flag systems. The RFC notes Eppo's mobile SDK has improved since evaluation, so it could be reconsidered if Omnichannel proves insufficient.

**Q: Are there any downsides to Omnichannel?**
A: Yes, noted in RFC:
- Maintenance burden on internal team
- Less sophisticated analytics than Eppo
- 5-minute polling vs real-time updates (performance trade-off)

**Q: Who maintains Omnichannel now?**
A: Platform team (per Confluence docs), with contributions from teams using it. Jane Smith is the primary maintainer according to git blame.

## Research Techniques Demonstrated

1. **Multi-source search**: Started with Glean to find RFC, epic, and discussions
2. **Document deep-dive**: Read complete RFC for detailed evaluation
3. **Ticket tracking**: Followed Jira epic to understand timeline and scope
4. **Code archaeology**: Used git log and PR history for implementation details
5. **Team context**: Checked Slack for informal discussions and concerns
6. **Synthesis**: Combined all sources into timeline with clear decision factors
7. **Anticipation**: Addressed obvious follow-up questions proactively

This example demonstrates how to investigate an architectural decision thoroughly using all available research tools.
